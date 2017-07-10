package main

import (
	. "dhbeat/def"
	. "dhbeat/util"
	"github.com/fsnotify/fsnotify"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"time"
)

func main() {
	LocalDb, _ = scribble.New("log", nil)
	var err error
	if err != nil {
		panic(err)
	}
	initProducer()
	StartWatcher()
}

func initProducer() {
	var e error
	Nc, e = nats.Connect(NATS_HOST)
	if e != nil {
		Error(e)
	}
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
					}
				}
			case err := <-watcher.Errors:
				Error(err)
			}
		}
	}()

	err = watcher.Add("/data/log")
	if err != nil {
		Error(err)
	}
	<-done
}

func ReadChanges(file string) []string {
	r := []string{}
	//ofst := LoadOffset(file)

	return r
}

func LoadOffset(topic string) int64 {
	i := int64(0)
	LocalDb.Read(topic, "offset", &i)
	//Debug("LoadOffset", i)
	return i
}

func SaveOffset(topic string, offset int64) {
	LocalDb.Write(topic, "offset", offset)
	Debug("SaveOffset", offset)
}

func AutoSaveOffset(topic string) {
	for {
		time.Sleep(time.Duration(1 * time.Second))
		old := LoadOffset(topic)
		tmp, _ := Cmap.Get(topic)
		offset := ToInt64(tmp)
		if offset != old {
			SaveOffset(topic, offset)
		}
	}
}
