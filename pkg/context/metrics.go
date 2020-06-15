package context

import (
	"bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"net"
	"net/http"
	"sync"
)

type Metrics struct {
	CPU        int64
	Memory     int64
	Province   string
	City       string
	Version    string
	Status     string
	ErrorMsg   []string
	CurrentApp string
	Mux        sync.Mutex
	Ip         string
	Mail       string
}

// TODO is move to pkg consts ?
const (
	metricsUrl = "http://get.choerodon.com.cn/api/v1/metrics"
	ipAddr     = "ns1.dnspod.net:6666"
)

func (m *Metrics) Send() {
	log.Debug("sending metrics...")
	data := m.pack()
	client := http.Client{}
	req, err := http.NewRequest("POST", metricsUrl, bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	}
	m.Ip = "127.0.0.1"
	resp, err := client.Do(req)
	if err != nil {
		log.Debug(err)
	}
	if resp.StatusCode != 200 {
		log.Debugf("send metrics failed with code: %d", resp.StatusCode)
	}
}

func (m *Metrics) pack() []byte {
	var (
		//v interface{} // value to decode/encode into
		b  []byte
		mh codec.MsgpackHandle
	)

	enc := codec.NewEncoderBytes(&b, &mh)

	err := enc.Encode(m)
	if err != nil {
		log.Error(err)
	}
	return b
}

func GetPublicIP() string {
	conn, err := net.Dial("tcp", ipAddr)
	if err != nil {
		log.Error(err)
		return "127.0.0.1"
	}
	defer conn.Close()

	ip, _ := bufio.NewReader(conn).ReadString('\n')
	return ip
}
