package test

import (
	. "dhbeat/util"
	"fmt"
	"testing"
)

func TestLogParser_Parse(t *testing.T) {
	parser := LogParser{}
	msg := ReadFile("sample.txt")

	//msg := ""
	p := parser.Parse(msg)
	Debug(JsonEncode(p))
	fmt.Println(p)
}

func BenchmarkLogParser_Parse(b *testing.B) {
	msg := ReadFile("sample.txt")
	for i := 0; i < b.N; i++ {
		parser := LogParser{}
		parser.Parse(msg)
	}
}

func TestLogParser_ProcFile_2(t *testing.T) {
	offset := ProcFile("/data/soooner/soooner_cache.log")
	Debug(offset)
	Debug(Aggr.Count, Aggr.Size())
	ps := Aggr.Dump()
	for _, v := range ps {
		//Debug(JsonEncode(v))
		log := JoinStr(v["time_local"],
			",",
			v["spid"],
			",",
			v["pid"],
			",",
			v["dhbeat_hostname"],
			",",
			v["request_time"],
			",",
			v["bytes_sent"])
		Debug(log)
	}
}
