package util

import (
	"testing"
	"fmt"
)

var aggr Aggregator

func init() {
	aggr = Aggregator{}
}

func TestAggregator_Add(t *testing.T) {
	p := P{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21130","time_local":"2017-06-06T16:05:00+08:00"}
	aggr.Add(p)
	size := aggr.Size()
	if size != 1 {
		t.Fatal()
	}
	fmt.Println(size,aggr.Dump())
	p = P{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21130","time_local":"2017-06-06T16:05:00+08:00"}
	aggr.Add(p)
	size = aggr.Size()
	if size != 1 {
		t.Fatal()
	}
	fmt.Println(size,aggr.Dump())

	p = P{"bw":0.3,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21131","time_local":"2017-06-06T16:05:00+08:00"}
	aggr.Add(p)
	size = aggr.Size()
	if size != 2 {
	t.Fatal()
	}
	fmt.Println(size,aggr.Dump())
	p = P{"bw":0.3,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21131","time_local":"2017-06-06T16:05:00+08:00"}
	aggr.Add(p)
	size = aggr.Size()
	if size != 2 {
		t.Fatal()
	}
	fmt.Println(size,aggr.Dump())
}

func TestAggregator_Dump(t *testing.T) {
	p := []P{
		{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21130","time_local":"2017-06-06T16:05:00+08:00"},
		{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21131","time_local":"2017-06-06T16:05:00+08:00"},
		{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21130","time_local":"2017-06-06T16:05:00+08:00"},
		{"bw":0.1,"dhbeat_hostname":"acer-PC","pid":"8031006300","spid":"21130","time_local":"2017-06-06T16:05:00+08:00"},

	}
	aggr.Add(p...)
	r := aggr.Dump()
	if len(r) != 2 {
		t.Fatal(r)
	}
	r0 := r[0]
	if r0["bw"] != 0.2 {
//		t.Fatal()
	}
	// 遍历map
	//for k, v := range r {
	//	fmt.Println(k, v)
	//	if v != 0.2{
	//		t.Fatal()
	//	}
	//}
	arr := aggr.Dump()
	fmt.Println(arr)
}
