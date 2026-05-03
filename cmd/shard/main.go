package main

import (
	"flag"
	"log"
	"net"
	"runtime/debug"
	"starconflict/lib/asyncreq"
	"starconflict/lib/auth"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type session struct {
	uid          uint64
	sessionToken string
	seq          uint16
	conn         net.Conn
}

func (session *session) handleAuthentication() error {
	key, group := auth.CreateClientChallenge()
	session.seq += 1
	if err := auth.SendChallenge(session.conn, key, session.seq); err != nil {
		return err
	}
	for {
		hdr, body, err := protocol.ParseNextMessage(session.conn)
		if err != nil {
			return err
		}
		if hdr.CommandType == types.CCMD_AUTH_REQUEST {
			_ = body
			log.Printf("Got a %s", hdr.CommandType)
			valid, dr, err := auth.Authenticate(body, key, group)
			if valid {
				session.seq += 1
				auth.SendAuthAck(session.conn, session.seq, hdr.Sequence)
			} else {
				session.conn.Write(protocol.MakeDisconnectMessage(dr))
			}
			if err != nil {
				log.Printf("Error %v", err)
				return err
			}
			return nil
		}
		log.Printf("Unhandled message of type: %s", hdr.CommandType)
	}
}

func (session *session) handleMainLoop() error {
	for {
		hdr, body, err := protocol.ParseNextMessage(session.conn)
		_ = body
		if err != nil {
			return err
		}
		if hdr.Special {
			// TODO: Do something
			continue
		}
		if hdr.CommandType == types.CSCMD_ASYNC_REQ {
			session.seq += 1
			go asyncreq.HandleAsyncReq(hdr, body, session.seq, session.conn, session.uid)
			continue
		}
		if hdr.CommandType == types.SCMD_KEEP_ALIVE {
			// TODO: Keep Alive
			continue
		}
		log.Printf("Unhandled message of type: %s", hdr.CommandType)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("handle failed: %v", err)
			debug.PrintStack()
		}
	}()
	session := session{}
	session.seq = 0
	session.conn = conn
	session.handleAuthentication()
	session.handleMainLoop()
}

func main() {
	listenAddress := flag.String("listen", "127.0.0.1:3802", "Address to listen on")
	flag.Parse()

	listen, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("shard listening on :%s", *listenAddress)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn)
	}
}
