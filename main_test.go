package main

import (
	. "dhbeat/util"
	"testing"
)

func init() {
	initProducer()
}

func TestProcFile(t *testing.T) {
	i := ProcFile("/data/log/test.csv")
	Debug(i)
}

func BenchmarkProcFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProcFile("/data/log/test.csv")
	}
}
