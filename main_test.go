package main

import (
	. "dhbeat/util"
	"testing"
)

func init() {
	initProducer()
}

func TestProcFile(t *testing.T) {
	i := ProcFile("E:/data/soooner_cache.log")
	Debug(i)
}

func BenchmarkProcFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProcFile("E:/data/soooner_cache.log")
	}
}
