package session

import (
	"github.com/jmoiron/sqlx"
	"net"
)

type Session struct {
	Uid          uint64
	sessionToken string
	seq          uint16
	Conn         net.Conn
	Db           *sqlx.DB
}

func (s *Session) GetNextSeq() uint16 {
	s.seq += 1
	return s.seq
}
