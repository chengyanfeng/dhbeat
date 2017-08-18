package test

import (
	"github.com/nats-io/go-nats"
	"testing"
)

var log_sample = `[08/May/2017:14:17:40 +0800] 0.588 61.148.243.45 200 75973 GET /if5ax/258/50/92/letv-uts/14/ver_00_22-1099556292-avc-2997249-aac-128000-2588360-1015099561-f1f58b59596cfeb9ce276b9068e51e48-1493874820974.m3u8?srcpara=leditafl%3D79220d924f43edd5%26ledituid%3D721621445%26leditcid%3D29241531%26leditcip%3D61.148.243.45%26leditfl%3Dnull%26crypt%3D31aa7f2e3279%26b%3D3137%26nlh%3D4096%26nlt%3D60%26bf%3D84%26p2p%3D1%26video_type%3Dmp4%26termid%3D2%26tss%3Dios%26platid%3D3%26splatid%3D341%26its%3D0%26qos%3D4%26fcheck%3D0%26amltag%3D4701%26mltag%3D4701%26proxy%3D1039176407%2C1039176737%2C467476916%26lsbv%3D2ga%26uid%3D156270520%26keyitem%3DWne1u_MCbvVKI8_TgM51eJFXewfs4uPxt0sGV1DC1PBfQTrtiLwkJTFP2K0.%26ntm%3D1494242400%26nkey%3De03364bacc690279215792d025e8263c%26nkey2%3D6de88b138f70bfe44467a6bb56cdf83e%26geo%3DCN-1-5-2%26cvid%3D368598997241%26dname%3Dmobile%26hwtype%3Diphone%26iscpn%3Df9050%26key%3D64d92246028b8d6fbfb253ca5376147e%26m3v%3D3%26mmsid%3D64602912%26ostype%3Dmacos%26p1%3D0%26p2%3D00%26payff%3D0%26pcode%3D010210000%26pid%3D10031263%26playid%3D0%26sign%3Dmb%26tag%3Dmobile%26tm%3D1494224228%26uinfo%3DAAAAAAAAAAC7AT1nBBKQjKhwQOl69GzHSXtspK1HsCoJ5RfVr9OkYQ%3D%3D%26uuid%3D43537108-DEF7-42DA-9B71-4B2EEE864D4F1494224227578%26version%3D7.0%26vid%3D29241531%26vtype%3D52%26errc%3D0%26gn%3D1190%26vrtmcd%3D106%26buss%3D4701%26cips%3D61.148.243.45&tag1=1&videoname=5aSW56eR6aOO5LqRMzQ%3D&videoid=29241531&apptype=app&userid=18601013472&userip=61.148.243.45&spid=22125&pid=m29241531&preview=1&portalid=368&spip=61.240.143.159&spport=80&tradeid=732bdf19ae22d64c83ef343c92bd262b&lsttm=20170508201710&enkey=a815aedc3d88f4eb8511fcb2cbb390e0 - "AppleCoreMedia/1.0.0.13F69 (iPhone; U; CPU OS 9_3_2 like Mac OS X; zh_cn)" "MISS" "-" "-"`

func Test_Nats(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		t.Fatal(err)
	}
	err = nc.Publish("test", []byte(log_sample))
	if err != nil {
		t.Fatal(err)
	}
}

func Benchmark_Nats(b *testing.B) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		err = nc.Publish("test", []byte(log_sample))
		if err != nil {
			b.Fatal(err)
		}
	}
}
