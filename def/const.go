package def

import (
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"github.com/orcaman/concurrent-map"
	"os"
)

var Cmap cmap.ConcurrentMap = cmap.New()
var Lock cmap.ConcurrentMap = cmap.New()
var LocalDb *scribble.Driver
var NATS_HOST string
var Nc *nats.Conn
var BLOCK_SIZE int64
var DIR string
var TYPE string
var HOSTNAME, _ = os.Hostname()
var Q_NAME = "bw"

const (
	GENERAL_ERR int = 400
)
