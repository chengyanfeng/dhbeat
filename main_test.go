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
