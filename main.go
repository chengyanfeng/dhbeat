package main

import (
	. "dhbeat/def"
	. "dhbeat/util"
	"github.com/fsnotify/fsnotify"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"os"
	"strings"
	"time"
)

func main() {
	LocalDb, _ = scribble.New("log", nil)
	var err error
	if err != nil {
		panic(err)
	}
	initProducer()
	scanFiles()
	go AutoSaveOffset()
	StartWatcher()
}

// 初始化 nats 连接
func initProducer() {
	var e error
	Nc, e = nats.Connect(NATS_HOST)
	if e != nil {
		Error(e)
	}
}

// 列出dir下面的所有log文件，加载每个文件的offset
func scanFiles() {
	files := DirTree(DIR, ".log", 100)
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
					Debug("modified file:", file)
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
	if e != nil {
		Error(e)
		return offset
	}
	step := int64(BLOCK_SIZE)
	half := ""
	for ptr := offset; ptr <= size; ptr += step {
		b := make([]byte, step)
		d, _ := f.ReadAt(b, offset)

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
				err := Nc.Publish(Q_NAME, []byte(line))
				if err != nil {
					Error(err)
					break
				}
				Cmap.Set(file, offset)
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
			Debug("SaveOffset", file, offset)
		}
	}
}
