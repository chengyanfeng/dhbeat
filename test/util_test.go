package test

import (
	. "dhbeat/util"
	"testing"
	"time"
)

func TestBucketMinute(t *testing.T) {
	str := "2017-08-18 15:00:00"
	ti, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 60; i++ {
		nt := ti.Add(time.Duration(i) * time.Minute)
		Debug(BucketMinute(nt, 5))
	}
}
