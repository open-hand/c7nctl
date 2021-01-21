package client

import (
	"github.com/ugorji/go/codec"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestPack(t *testing.T) {
	m := Metrics{
		CPU:      8,
		Memory:   1024 * 1024 * 1024 * 16,
		Province: "Shanghai",
		City:     "Shanghai",
		ErrorMsg: []string{"error 01", "error 02"},
	}

	data := m.pack()
	var mh codec.MsgpackHandle

	dec := codec.NewDecoderBytes(data, &mh)

	m2 := &Metrics{}
	err2 := dec.Decode(m2)
	if err2 != nil {
		t.Fatal("pack metrics failed")
	}
	dec = codec.NewDecoderBytes(data, &mh)
	if m2.ErrorMsg[0] != "error 01" {
		t.Fatal("unpack metrics failed")
	}
}

func TestSend(t *testing.T) {

	m := Metrics{
		CPU:      8,
		Memory:   1024 * 1024 * 1024 * 16,
		Province: "Shanghai",
		City:     "Shanghai",
		ErrorMsg: []string{"error 01", "error 02"},
	}
	m.Send()
}

func TestGetPublicIP(t *testing.T) {
	log.Info(GetPublicIP())
}
