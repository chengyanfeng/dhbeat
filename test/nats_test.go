package test

import (
	. "dhbeat/def"
	"github.com/nats-io/go-nats"
	"testing"
)

var log_sample = `2017/07/18 11:00:00,22125,m10004160,altmac.local,3.066,44168128`

func Test_Nats(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatal(err)
	}
	err = nc.Publish(Q_NAME, []byte(log_sample))
	if err != nil {
		t.Fatal(err)
	}
}
