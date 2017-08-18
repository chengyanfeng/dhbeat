package util

import (
	. "dhbeat/def"
	"net/url"
	"os"
	"strings"
	"sync"
)

var Aggr = Aggregator{Cache: P{}, Lock: sync.Mutex{}}

type LogParser struct {
}

func (this *LogParser) Parse(msg string) P {
	segs := ToSegs(msg)
	p := P{}
	if len(segs) < 14 {
		Error("Invalid msg", len(segs), msg)
		return p
	}
	p["request_time"] = ToFloat(segs[1]) // 服务时间
	p["bytes_sent"] = ToInt(segs[6]) * 8 // 发送长度
	p["url"] = segs[8]                   // 完整请求链接
	this.ParseUrl(p)
	min, _ := ToTime(segs[0])
	p["time_local"] = BucketMinute(min, 5)
	p["dhbeat_hostname"] = HOSTNAME
	return p
}

func ToSegs(msg string) []string {
	tmp := strings.Split(msg, " ")
	segs := []string{}
	token := ""
	wait := false
	for _, v := range tmp {
		if StartsWith(v, `"`) || StartsWith(v, `[`) {
			if EndsWith(v, `"`) {
				token = Replace(v, []string{`"`}, "")
				segs = append(segs, token)
			} else {
				wait = true
				token = v[1:]
			}
		} else {
			if wait {
				if EndsWith(v, `"`) || EndsWith(v, `]`) {
					token = JoinStr(token, " ", v)
					wait = false
					token = token[0 : len(token)-1]
					segs = append(segs, token)
					token = ""
				} else {
					token = JoinStr(token, " ", v)
				}
			} else {
				segs = append(segs, v)
			}
		}
	}
	return segs
}

// 解析url
func (this *LogParser) ParseUrl(p P) {
	u, err := url.Parse(ToString(p["url"]))
	if err != nil {
		Error(err)
		return
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		Error(err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			Error("ParseUrl", err, m)
		}
	}()
	p["spid"] = m["spid"][0] //客户标识（产品）可暂时作为产品id
	p["pid"] = m["pid"][0]   //产品id
	delete(p, "url")
}

// 处理文件，从offset开始读取，分行，然后给Aggregator
func ProcFile(file string) int64 {
	Debug("ProcFile", file)
	tmp, _ := Cmap.Get(file)
	offset := ToInt64(tmp)
	_, b := Lock.Get(file)
	if b {
		return offset
	}
	Lock.Set(file, 1)
	defer Lock.Remove(file)
	lines := []string{}
	size := FileSize(file)
	f, e := os.Open(file)
	defer f.Close()
	if e != nil {
		Error(e)
		return offset
	}
	step := int64(BLOCK_SIZE)
	half := ""
	//由于目前日志文件名不改变，则整点将offset归零
	if offset > size {
		offset = 0
	}
	for ptr := offset; ptr <= size; ptr += step {
		b := make([]byte, step)
		d, err := f.ReadAt(b, offset)
		if err != nil {
			Error(err)
		}
		body := string(b[:d])
		lines = strings.Split(body, "\n")
		if !IsEmpty(half) {
			lines[0] = half + lines[0]
			half = ""
		}
		if !EndsWith(body, "\n") {
			half = lines[len(lines)-1]
			if len(lines) > 1 {
				lines = lines[0 : len(lines)-1]
			}
		} else {
			half = ""
		}
		offset += int64(d)
		parser := LogParser{}
		for _, line := range lines {
			if Trim(line) != "" {
				p := parser.Parse(line)
				Aggr.Add(p)
				Cmap.Set(file, offset)
			}
		}
	}
	return offset
}
