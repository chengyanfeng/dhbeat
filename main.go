package main

import (
	. "dhbeat/controller"
	. "dhbeat/def"
	. "dhbeat/task"
	. "dhbeat/util"
	"encoding/json"
	. "github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nanobox-io/golang-scribble"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	beego.BConfig.Listen.HTTPPort = ToInt(Trim(os.Getenv("port")), 6117)
	beego.BConfig.RecoverPanic = true
	beego.BConfig.EnableErrorsShow = true
	beego.AutoRouter(&ApiController{})
	//beego.SetLogger("file", `{"filename":"logs/run.log"}`)
	//beego.BeeLogger.SetLogFuncCallDepth(4)
	LocalDb, _ = scribble.New("log", nil)
	var err error
	// todo 配置通过文件读取
	Stream, err = gorm.Open("postgres", "host=pipeline1 user=soooner dbname=dh sslmode=disable password=")
	defer Stream.Close()
	if err != nil {
		panic(err)
	}
	//Citus, err = gorm.Open("postgres", "host=citus user=postgres dbname=postgres sslmode=disable password=")
	//defer Citus.Close()
	//if err != nil {
	//	panic(err)
	//}
	initConsumer()
	defer func() {
		if err := KafkaConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	go consume("log_topic", ProcLog)
	//go consume("ws_topic", ProcWs)
	beego.Run()
}

func consume(topic string, f func(msg *ConsumerMessage)) {
	go AutoSaveOffset(topic)
	offset := LoadOffset(topic) + 1
	if offset < 2 {
		offset = OffsetNewest
	}
	// todo 要考虑分区消费，便于并行处理
	logConsumer, err := KafkaConsumer.ConsumePartition(topic, 0, offset)
	if err != nil {
		Error(err)
		logConsumer, err = KafkaConsumer.ConsumePartition(topic, 0, OffsetNewest)
	}

	defer func() {
		if err := logConsumer.Close(); err != nil {
			Error(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	Debug("Consuming", topic)
ConsumerLoop:
	for {
		select {
		case msg := <-logConsumer.Messages():
			Dhq <- func() {
				f(msg)
			}
		case <-signals:
			break ConsumerLoop
		}
	}
}

func initConsumer() {
	var err error
	KafkaConsumer, err = NewConsumer([]string{"kafka1:19092", "kafka2:19092", "kafka3:19092"}, nil)
	if err != nil {
		panic(err)
	}
}

func ProcLog(msg *ConsumerMessage) {
	parser := LogParser{}
	t := string(msg.Value)
	//由于t为json格式，下面处理json转换成map，取出message对应的值，然后进行解析
	msg1 := []byte(t)
	dat := make(map[string]interface{})
	if err := json.Unmarshal(msg1, &dat); err != nil {
		panic(err)
	}

	p := parser.Parse(dat["message"].(string))

	//Debug("Consumed ", msg.Offset, string(msg.Value))
	err := InsertDb(p)
	if err != nil {
		Error(err)
	} else {
		Cmap.Set(msg.Topic, msg.Offset)
	}
}

func InsertDb(p *P) (e error) {
	v := *p

	//todo
	//Debug(v)
	e = Stream.Exec(`insert into s_log (time_local,request_time,remote_addr,status,err_code,request_length,bytes_sent,request_method,http_referer,http_user_agent,cache_status,http_range,sent_http_content_range,uri,userip,spid,pid,spport,lsttm,vkey,userid,portalid,spip,sdtfrom,tradeid,enkey,st,bw) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		v["time_local"], v["request_time"], v["remote_addr"], v["status"], v["err_code"], v["request_length"], v["bytes_sent"], v["request_method"], v["http_referer"], v["http_user_agent"], v["cache_status"], v["http_range"], v["sent_http_content_range"], v["uri"], v["userip"], v["spid"], v["pid"], v["spport"], v["lsttm"], v["vkey"], v["userid"], v["portalid"], v["spip"], v["sdtfrom"], v["tradeid"], v["enkey"], v["st"], v["bw"]).Error

	if e != nil {
		return
	}
	//e = Citus.Exec(`insert into t_userlog (msg) values (?)`,
	//	v["msg"]).Error
	//if e != nil {
	//	return
	//}
	return
}

func LoadOffset(topic string) int64 {
	i := int64(0)
	LocalDb.Read(topic, "offset", &i)
	//Debug("LoadOffset", i)
	return i
}

func SaveOffset(topic string, offset int64) {
	LocalDb.Write(topic, "offset", offset)
	Debug("SaveOffset", offset)
}

func AutoSaveOffset(topic string) {
	for {
		time.Sleep(time.Duration(1 * time.Second))
		old := LoadOffset(topic)
		tmp, _ := Cmap.Get(topic)
		offset := ToInt64(tmp)
		if offset != old {
			SaveOffset(topic, offset)
		}
	}
}
