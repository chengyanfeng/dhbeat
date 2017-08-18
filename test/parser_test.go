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
