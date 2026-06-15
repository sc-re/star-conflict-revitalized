package session

import (
	"net"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Session struct {
	Uid          uint64
	sessionToken string
	seq          uint16
	Conn         net.Conn
	Db           *mongo.Client
}

func (s *Session) GetNextSeq() uint16 {
	s.seq += 1
	return s.seq
}
