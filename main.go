package main

import (
	. "dhbeat/def"
	. "dhbeat/util"
	"github.com/fsnotify/fsnotify"
	"github.com/nanobox-io/golang-scribble"
	"os"
	"strconv"
	"strings"
	"time"
	//"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats"
	"sync"
	"fmt"
)

func main() {
	SR = new(sync.RWMutex)
	initConf()
	Debug("initConf()...........")
	LocalDb, _ = scribble.New("log", nil)
	Debug("LocalDb..............")
	var err error
	if err != nil {
		panic(err)
	}
	initProducer()
	Debug("initProducer()..............")
	scanFiles()
	Debug("scanFiles()..............")
	go AutoSaveOffset()
	Debug("go AutoSaveOffset()..............")
	go AutoDump()
	Debug("go AutoDump()..............")
	StartWatcher()
}

var aggr = Aggregator{}

//初始化配置文件
func initConf() {
	myConfig := new(Config)
	config := myConfig.InitConfig("./", "DhBeat.ini", "nats")
	NATS_HOST = config["nats_host"]
	BLOCK_SIZE, _ = strconv.ParseInt(config["block_size"], 10, 64)
	Q_NAME = config["q_name"]
	DIR = config["dir"]
	CLUSTER_ID = config["cluster_id"]
	CLIENT_ID = config["client_id"]
	TYPE = config["type"]
	//LD ,_  = strconv.Atoi(config["ld"])
	LD ,_ = strconv.ParseInt(config["ld"], 10, 64)
}

//var ah stan.AckHandler

//初始化 nats-streaming 连接
func initProducer() {
	var err error
	Nc, err = nats.Connect(NATS_HOST)
	//Sc, err = stan.Connect(CLUSTER_ID, CLIENT_ID, stan.NatsURL(NATS_HOST))
	//Sc, err = stan.Connect(CLUSTER_ID, CLIENT_ID, stan.NatsURL("nats://127.0.0.1:4222"))
	if err != nil {
		panic(err)
	}

}

// 列出dir下面的所有log文件，加载每个文件的offset
func scanFiles() {
	files := DirTree(DIR, TYPE, 100)
	for _, file := range files {
		offset := LoadOffset(file)
		Cmap.Set(file, offset)
		ProcFile(file)
	}
}

// 开始监听文件变化（修改、删除）
func StartWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Error(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				file := event.Name
				if event.Op&fsnotify.Write == fsnotify.Write {
					//Debug("modified file:", file)
					ProcFile(file)
				} else if event.Op == fsnotify.Remove {
					Debug("delete file:", file)
					Cmap.Remove(file)
				}
			case err := <-watcher.Errors:
				Error(err)
			}
		}
	}()

	err = watcher.Add(DIR)
	if err != nil {
		Error(err)
	}
	<-done
}

// 处理文件，从offset开始读取，分行，然后扔给nats
func ProcFile(file string) int64 {
	var n int64
	tmp, _ := Cmap.Get(file)
	offset := ToInt64(tmp)
	_, b := Lock.Get(file)
	if b {
		return offset
	}
	Lock.Set(file, 1)
	defer Lock.Remove(file)
	lines := []string{}
	size := FileSize(file)
	f, e := os.Open(file)
	defer f.Close()
	if e != nil {
		Error(e)
		return offset
	}
	step := int64(BLOCK_SIZE)
	half := ""
	t0 := time.Now().UnixNano() / int64(time.Millisecond)
//	Debug("offset =", offset, "size =", size)
	//由于目前日志文件名不改变，则整点将offset归零
	if offset > size {
		offset = 0
	}

	for ptr := offset; ptr <= size; ptr += step {
		b := make([]byte, step)
		var d int
		d, _ = f.ReadAt(b, offset)
		body := string(b[:d])
		lines = strings.Split(body, "\n")
		if !IsEmpty(half) {
			lines[0] = half + lines[0]
			half = ""
		}
		if !EndsWith(body, "\n") {
			half = lines[len(lines)-1]
			if len(lines) > 1 {
				lines = lines[0 : len(lines)-1]
			}
		} else {
			half = ""
		}
		offset += int64(d)

		for _, line := range lines {

			if Trim(line) != "" {
				n++
				log := line + " " + HOSTNAME + " " + file
				fmt.Println(log)
				//err := Sc.Publish(Q_NAME, []byte(log))
				parser := LogParser{}
				p := parser.Parse(log)
				aggr.Add(p...)
			//	app := aggr.Dump()
			//	Debug(len(app))
				//err := Nc.Publish(Q_NAME, []byte(log))
				//
				//if err != nil {
				//	Error(err)
				//	break
				//}
				Cmap.Set(file, offset)

			}
		}
		//测试发布数据时间
		if offset == size {
			t1 := time.Now().UnixNano() / int64(time.Millisecond)
			if t1-t0 != 0 {
			//	Debug("共发送", n, "条数据，共花时间", t1-t0, "毫秒,发布速度为", n/(t1-t0)*1000, "条/秒")
			}
		}
	}
	return offset
}

// 加载本地保存的文件 offset
func LoadOffset(file string) int64 {
	i := int64(0)
	LocalDb.Read(file, "offset", &i)
	Debug("LoadOffset", file, i)
	return i
}

// 自动定时保存 offset
func AutoSaveOffset() {
	for {
		time.Sleep(time.Duration(1 * time.Second))
		for v := range Cmap.Iter() {
			file := v.Key
			offset := ToInt64(v.Val)
			LocalDb.Write(file, "offset", offset)
			//Debug("SaveOffset", file, offset)
		}
	}
}

// 自动定时dump aggr

func AutoDump() {

	for {
		//	time.Sleep(time.Duration(300 * time.Second))
		timestamp := time.Now().Unix()
		if timestamp%(LD*60) == 0 {
			SR.Lock()
			data := aggr.Dump()
			for _, v := range data {
				log := ToString(v["time_local"]) + "|" + ToString(v["spid"]) + "|" + ToString(v["pid"]) + "|" + ToString(v["dhbeat_hostname"]) + "|" + ToString(v["bw"])
				//	log := key + "|" + ToString(value)
				//	Nc.Publish(Q_NAME, []byte(JsonEncode(v)))
				err := Nc.Publish(Q_NAME, []byte(log))
				if err != nil {
					Error(err)
					break
				}
				//	fmt.Println(log)
				Debug(log)
			}
			//删除已发布的日志
			aggr.Ju = []P{}
			aggr.Lm = make(map[string]float64)

			SR.Unlock()
		}
	}
}