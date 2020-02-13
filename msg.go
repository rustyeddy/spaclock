package main

// Message is a simple generic array of key values
type Value struct {
	Key   string      `json:"key"`
	Value interface{} `json:"val"`
}

type Message []Value
