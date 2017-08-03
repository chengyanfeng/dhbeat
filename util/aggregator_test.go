package util

import "testing"

var aggr Aggregator

func init() {
	aggr = Aggregator{}
}

func TestAggregator_Add(t *testing.T) {
	p := P{"time": "2017/08/03 14:00:01", "spid": "xxx", "bw": 1.0}
	aggr.Add(p)
	size := aggr.Size()
	if size != 1 {
		t.Fatal()
	}
	p = P{"time": "2017/08/03 14:00:01", "spid": "xxx", "bw": 1.1}
	aggr.Add(p)
	size = aggr.Size()
	if size != 1 {
		t.Fatal()
	}
}

func TestAggregator_Dump(t *testing.T) {
	p := []P{
		{"time": "2017/08/03 14:00:01", "spid": "xxx", "bw": 1.0},
		{"time": "2017/08/03 14:00:01", "spid": "xxx", "bw": 1.1},
	}
	aggr.Add(p...)
	r := aggr.Dump()
	if len(r) != 1 {
		t.Fatal(r)
	}
	r0 := r[0]
	if r0["bw"] != 2.1 {
		t.Fatal()
	}
}
