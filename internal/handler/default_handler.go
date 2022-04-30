package handler

type Default struct{}

func (s *Default) DoHelo(args string) bool {
	return true
}

func (s *Default) DoMailFrom(args string) bool {
	return true
}

func (s *Default) DoRcptTo(args string) bool {
	return true
}

func (s *Default) DoData(d []byte) bool {
	return true
}

func (s Default) String() string {
	return "Default Handler"
}
