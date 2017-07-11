package main

import (
	. "dhbeat/def"
	. "dhbeat/util"
	"github.com/fsnotify/fsnotify"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"time"
)

var dir = "/data/log"

func main() {
	LocalDb, _ = scribble.New("log", nil)
	var err error
	if err != nil {
		panic(err)
	}
	initProducer()
	scanFiles()
	StartWatcher()
	// todo: auto save cmap
}

func initProducer() {
	var e error
	Nc, e = nats.Connect(NATS_HOST)
	if e != nil {
		Error(e)
	}
}

func scanFiles() {
	// todo 列出dir下面的所有log文件，加载每个文件的offset，如果没有则是新文件，如果offset对应的文件缺失说明被删除了
}

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
				if event.Op&fsnotify.Write == fsnotify.Write {
					file := event.Name
					Debug("modified file:", file)
					lines := ReadChanges(file)
					for _, line := range lines {
						Nc.Publish("log", []byte(line))
						// todo: save offset into cmap
					}
				} else if event.Op&fsnotify.Write == fsnotify.Remove {
					// todo: remove cmap[file]
				}
			case err := <-watcher.Errors:
				Error(err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		Error(err)
	}
	<-done
}

func ReadChanges(file string) []string {
	r := []string{}
	// todo: load offset from cmap,
	// todo: if cmp[file]==0, read from 0 to 10 lines
	//ofst := LoadOffset(file)

	return r
}

func LoadOffset(file string) int64 {
	i := int64(0)
	LocalDb.Read(file, "offset", &i)
	// todo: save offset into cmap
	//Debug("LoadOffset", i)
	return i
}

func SaveOffset(file string, offset int64) {
	LocalDb.Write(file, "offset", offset)
	Debug("SaveOffset", offset)
}

func AutoSaveOffset(file string) {
	for {
		time.Sleep(time.Duration(1 * time.Second))
		old := LoadOffset(file)
		tmp, _ := Cmap.Get(file)
		offset := ToInt64(tmp)
		if offset != old {
			SaveOffset(file, offset)
		}
	}
}
