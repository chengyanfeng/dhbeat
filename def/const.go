package def

import (
	//"github.com/Shopify/sarama"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"github.com/orcaman/concurrent-map"
	"gopkg.in/robfig/cron.v2"
	"time"
)

var Cron *cron.Cron
var Cmap cmap.ConcurrentMap = cmap.New()
var LocalDb *scribble.Driver
var NATS_HOST = "nats://nats.datahunter.cn:4222"
var Nc *nats.Conn
var UPTIME = time.Now().UnixNano() / int64(time.Millisecond)

const (
	Cname string = "job"
)

const (
	GENERAL_ERR int = 400
)
