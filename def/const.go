package def

import (
	//"github.com/Shopify/sarama"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
	"github.com/orcaman/concurrent-map"
	"os"
	"time"
	"sync"
)

var Cmap cmap.ConcurrentMap = cmap.New()
var Lock cmap.ConcurrentMap = cmap.New()
var LocalDb *scribble.Driver
var NATS_HOST string
var Nc *nats.Conn
var Sc stan.Conn
var BLOCK_SIZE int64
var Q_NAME string
var DIR string
var CLUSTER_ID string
var CLIENT_ID string
var TYPE string
var UPTIME = time.Now().UnixNano() / int64(time.Millisecond)
var HOSTNAME, _ = os.Hostname()
var LD int64
var SR *sync.RWMutex
const (
	GENERAL_ERR int = 400
)

//type Conn interface {
//	// Publish
//	Publish(subject string, data []byte) error
//	PublishAsync(subject string, data []byte, ah AckHandler) (string, error)
//
//	// Subscribe
//	Subscribe(subject string, cb MsgHandler, opts ...SubscriptionOption) (Subscription, error)
//
//	// QueueSubscribe
//	QueueSubscribe(subject, qgroup string, cb MsgHandler, opts ...SubscriptionOption) (Subscription, error)
//
//	// Close
//	Close() error
//
//	// NatsConn returns the underlying NATS conn. Use this with care. For
//	// example, closing the wrapped NATS conn will put the NATS Streaming Conn
//	// in an invalid state.
//	NatsConn() *nats.Conn
//}
