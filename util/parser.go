package util

import (
	"net/url"
	"strings"
	"math"
	"time"
	. "dhbeat/def"
)

type LogParser struct {

}

func (this *LogParser) Parse(msg string) ([]P) {
	//var LD int64 = 1
	// todo: 解析msg并入库，需要考虑数据计算，如带宽、开始时间，需要考虑多值返回，用于多stream入库
	ap := []P{}
	tmp := strings.Split(msg, " ")
	seg := []interface{}{}
	token := ""
	wait := false
	for _, v := range tmp {
		if StartsWith(v, `"`) || StartsWith(v, `[`) {
			if EndsWith(v, `"`) {
				token = Replace(v, []string{`"`}, "")
				seg = append(seg, token)
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
					seg = append(seg, token)
					token = ""
				} else {
					token = JoinStr(token, " ", v)
				}
			} else {
				seg = append(seg, v)
			}
		}
	}
	//Debug(len(seg), JsonEncode(seg))
	p := P{}
	if len(seg) != 16 {
		Error("Invalid msg", msg)
		return ap
	}
	if len(seg) == 16 && strings.Contains(seg[15].(string),"soooner_cache.log") {
	//	p["time_local_origin"], _ = ToTime(ToString(seg[0])) // 本地时间
		p["st"] = seg[0]
		//s, _ := ToTime(ToString(seg[0]))
	//	p["time_local"] = p["time_local_origin"].(time.Time).Format("2006-01-02") //只要年月日
		//fmt.Println(strings.Split(s.Format("2006-01-02"),"-")[0])
		p["request_time"] = ToFloat(seg[1])            // 服务时间
	//	p["remote_addr"] = seg[2]                  // 远端地址（客户端地址）
	//	p["status"] = seg[3]                    // HTTP回复状态码
	//	p["err_code"] = seg[4]                    // 错误码
	//	p["request_length"] = seg[5]                      // 请求长度
		p["bytes_sent"] = seg[6]                    // 发送长度
	//	p["request_method"] = seg[7]                  // 请求方式
		p["url"] = seg[8]                     // 完整请求链接
	//	p["http_referer"] = seg[9]                   // HTTP_REFERER
	//	p["http_user_agent"] = seg[10]                     // USERAGENT
	//	p["cache_status"] = seg[11]                    // 缓存状态（MISS HIT IOTHROUGH）
	//	p["dhbeat_hostname"] = seg[14]
		this.ParseUrl(p)
		this.ParseBandwidth(p)
	}
	if p != nil{
		p1 := p
		local_time, _ := ToTime(ToString(p1["st"]))

		dur := ToFloat(p1["request_time"])
		delete(p1, "request_time")
		//切分成5分钟时间段

		//将时间粒度分钟转换成秒单位
		t := int(LD * 60)
		var n int
		if int(math.Ceil(dur))%t == 0 {
			n = 0
		} else {
			n = 1
		}
		for i := 0; i < int(math.Ceil(dur))/t+n; i ++ {
			p0 := P{}
			fiveTime := time.Date(local_time.Year(), local_time.Month(), local_time.Day(), local_time.Hour(), (local_time.Minute()/int(LD)-i+1)*int(LD), 00, 00, time.Local)
			p0["time_local"] = fiveTime
			p0["spid"] = p1["spid"]
			p0["pid"] = p1["pid"]
			p0["bw"] = p1["bw"]
			p0["dhbeat_hostname"] = HOSTNAME
			//fmt.Println(p0["bw"])
			ap = append(ap,p0)
		}
	}

	return ap
}

// 解析url
func (this *LogParser) ParseUrl(p P) {
//	aes :=Aes{}
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
	//p["uri"] = u.Path                //文件地址
	//if m["userip"] != nil{
	//	p["userip"] = m["userip"][0]         //客户ip
//	}else {
	//	p["userip"] = nil
	//}

	p["spid"] = m["spid"][0]       //客户标识（产品）可暂时作为产品id
	p["pid"] = m["pid"][0]         //产品id
	//p["spport"] = m["spport"][0]     //sp端口
	//if strings.Count(string(m["userid"][0]),"")-1==32 {
	//	key :=[]byte("ac22273abb2f4960")
	//	userid,_ :=aes.CBCDecrypter(key, m["userid"][0])	//解密userid
	//	p["userid"] = strings.Split(userid, "\u0005")[0]       //用户id
//	}else {
	//	p["userid"] = m["userid"][0]			//用户id
//	}
	//p["portalid"] = m["portalid"][0]           //门户id，用以归类客户
//	p["spip"] = m["spip"][0] //服务商ip
	delete(p, "url")
}

func (this *LogParser) ParseBandwidth(p P) {
	// todo
	// todo
	//统计时间粒度
//	t := LD*60
	 var t int64 = 60
	var bw float64

	if  p["request_time"] == 0.00 {
		//	st = et
		bw = 0.00
	}else {
		dur := p["request_time"].(float64)

		//切分成5分钟平均带宽
		var d float64

		if dur > float64(LD*60) {
			var n int64
			if int64(math.Ceil(dur))%t == 0 {
				n = 0
			} else {
				n = 1
			}
			d = float64(t*int64(math.Ceil(dur))/t+n)
			bw = ToFloat(p["bytes_sent"]) * 8 / d
		}else{
			d = float64(t)
			bw = ToFloat(p["bytes_sent"]) * 8 / d
		}

	}
	p["bw"] = bw
	delete(p, "time_local_origin")
	delete(p, "bytes_sent")
}
