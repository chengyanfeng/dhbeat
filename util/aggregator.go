package util

import (
	"sync"
)

type Aggregator struct {
	Cache P
	Lock  sync.Mutex
}

// 往聚合器里面增加数据
func (this *Aggregator) Add(ps ...P) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	for _, p := range ps {
		key := Md5(p["time_local"],
			p["spid"],
			p["pid"],
			p["dhbeat_hostname"])

		tmp := this.Cache[key]
		if IsEmpty(tmp) {
			p["key"] = key
			this.Cache[key] = p
		} else {
			v := tmp.(P)
			v["bytes_sent"] = ToInt(v["bytes_sent"]) + ToInt(p["bytes_sent"])
			v["request_time"] = ToFloat(v["request_time"]) + ToFloat(p["request_time"])
		}
	}
}

// 导出聚合后的数据，同时清空
func (this *Aggregator) Dump() []P {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	dump := []P{}
	for _, tmp := range this.Cache {
		v := tmp.(P)
		dump = append(dump, v.Copy())
	}
	this.Cache = P{}
	return dump
}

// 聚合器容量
func (this *Aggregator) Size() int {
	return len(this.Cache)
}
