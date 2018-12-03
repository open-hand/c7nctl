package common

import (
	"github.com/ugorji/go/codec"
	"fmt"
	"sync"
	"net/http"
	"bytes"
	"github.com/vinkdong/gox/log"
)

type Metrics struct {
	CPU         int64
	Memory      int64
	Province    string
	City        string
	Version     string
	Status      string
	ErrorMsg    []string
	CurrentApp  string
	Mux         sync.Mutex
}

const (
	metricsUrl = "http://get.choerodon.com.cn/api/v1/metrics"
)

func (m *Metrics) Send() {
	data := m.pack()
	client := http.Client{}
	req, err := http.NewRequest("POST", metricsUrl, bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	if resp.StatusCode != 200 {
		log.Errorf("send metrics failed with code: %d", resp.StatusCode)
	}
}

func (m *Metrics) pack() []byte{
	var (
		//v interface{} // value to decode/encode into
		b []byte
		mh codec.MsgpackHandle
	)

	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(m)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
	}
	return b
}
