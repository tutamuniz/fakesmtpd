package handler

type Handler interface {
	DoHelo(args string) bool
	DoMailFrom(args string) bool
	DoRcptTo(args string) bool
	DoData(d []byte) bool
}
