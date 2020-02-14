package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Message is a simple generic array of key values
type KV struct {
	K string `json:"k"`
	V string `json:"v"`
}

func NewKV(k, v string) *KV {
	return &KV{k, v}
}

type Message struct {
	Values []KV `json:"values"`
}

func NewMessage(k, v string) (msg *Message) {
	msg = &Message{}
	msg.Values = append(msg.Values, *NewKV(k, v))
	return msg
}

func (msg *Message) Add(k, v string) {
	msg.Values = append(msg.Values, *NewKV(k, v))
}

func (msg *Message) Marshal() []byte {
	b, err := json.Marshal(msg.Values)
	if err != nil {
		log.Errorln(err)
		return nil
	}
	return b
}

func MessageFromBuffer(buf []byte) *Message {

	msg := &Message{}
	if err := json.Unmarshal(buf, &msg.Values); err != nil {
		log.Errorf("failed to unmarshal JSON %v", err)
		return nil
	}
	return msg
}
