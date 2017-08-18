package main

import (
	. "dhbeat/def"
	. "dhbeat/util"
	"github.com/fsnotify/fsnotify"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"strconv"
	"strings"
	"time"
)

func main() {
	initConf()
	LocalDb, _ = scribble.New("log", nil)
	var err error
	if err != nil {
		panic(err)
	}
	initProducer()
	scanFiles()
	go AutoSaveOffset()
	go AutoDump()
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
	TYPE = config["type"]
}

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
		if strings.Index(file, "soooner_cache.log") > -1 {
			offset := LoadOffset(file)
			Cmap.Set(file, offset)
			ProcFile(file)
		}
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
		}
	}
}

// 自动定时dump aggr

func AutoDump() {

	for {
		data := aggr.Dump()
		for _, v := range data {
			log := ToString(v["time_local"]) + "|" + ToString(v["spid"]) + "|" + ToString(v["pid"]) + "|" + ToString(v["dhbeat_hostname"]) + "|" + ToString(v["bw"])
			err := Nc.Publish(Q_NAME, []byte(log))
			if err != nil {
				Error(err)
				break
			}
			Debug(log)
		}
	}
}
