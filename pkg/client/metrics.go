package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"
	"io/ioutil"
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
	metricsUrl = "http://localhost:8080/api/v1/metrics"
	ipAddr     = "ns1.dnspod.net:6666"
)

func (m *Metrics) Send() {
	log.Debug("sending metrics...")
	contentType := "application/json;charset=utf-8"
	b, err := json.Marshal(m)
	if err != nil {
		log.Println("json format error:", err)
		return
	}

	body := bytes.NewBuffer(b)
	resp, err := http.Post(consts.MetricsUrl, contentType, body)
	if err != nil {
		log.Println("Post failed:", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Debugf("send metrics failed with code: %d", resp.StatusCode)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", err)
		return
	}
	log.Println("content:", string(content))

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
