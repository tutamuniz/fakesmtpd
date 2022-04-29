package handler

import (
	"bytes"
	"fmt"
)

type Simple struct{}

func (s *Simple) DoHelo(args string) bool {
	fmt.Println(args)
	return true
}

func (s *Simple) DoMailFrom(args string) bool {
	fmt.Println(args)
	return true
}

func (s *Simple) DoRcptTo(args string) bool {
	fmt.Println(args)
	return true
}

func (s *Simple) DoData(d []byte) bool {
	buff := bytes.NewBuffer(d)
	fmt.Println(buff.String())
	return true
}

func (s Simple) String() string {
	return "Simple Handler"
}
