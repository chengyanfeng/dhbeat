package test

import (
	. "dhbeat/util"
	"testing"
)

func TestAggregator_Add(t *testing.T) {
	p := P{"bytes_sent": 100, "request_time": 0.1, "dhbeat_hostname": "acer-PC", "pid": "8031006300", "spid": "21130", "time_local": "2017-06-06T16:05:00+08:00"}
	Aggr.Add(p)
	size := Aggr.Size()
	if size != 1 {
		t.Fatal()
	}
	Aggr.Add(p)
	size = Aggr.Size()
	if size != 1 {
		t.Fatal()
	}
	p2 := p.Copy()
	p2["spid"] = "123"
	Aggr.Add(p2)
	size = Aggr.Size()
	if size != 2 {
		t.Fatal()
	}
}

func TestAggregator_Dump(t *testing.T) {
	p := P{"bytes_sent": 100, "request_time": 0.1, "dhbeat_hostname": "acer-PC", "pid": "8031006300", "spid": "21130", "time_local": "2017-06-06T16:05:00+08:00"}
	ps := []P{p.Copy(), p.Copy(), p.Copy(), p.Copy()}
	Aggr.Add(ps...)
	r := Aggr.Dump()
	if len(r) != 1 {
		t.Fatal(r)
	}
	r0 := r[0]
	if r0["bytes_sent"] != 400 {
		t.Fatal()
	}
	if r0["request_time"] != 0.4 {
		t.Fatal()
	}
}
