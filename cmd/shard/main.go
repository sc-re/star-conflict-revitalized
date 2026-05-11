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
	"sync"
)

type session struct {
	uid          uint64
	sessionToken string
	seq          uint16
	conn         net.Conn
}

func (session *session) handleAuthentication(connectionMap *sync.Map) error {
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
				connectionMap.Store(session.uid, session)
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
	session.uid = 10
	//go userprofilenotification.SendUserProfileNotificationOnlineState(session.conn, session.uid, userprofilenotification.USER_STATE_ONLINE)
	go asyncreq.Send_ac_vessel_strip_improper_battle(session.conn)
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
			log.Printf("Recieved SCMD_KEEP_ALIVE")
			// TODO: Keep Alive
			continue
		}
		log.Printf("Unhandled message of type: %s", hdr.CommandType)
	}
}

func handle(conn net.Conn, connectionMap *sync.Map) {
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
	session.handleAuthentication(connectionMap)
	session.handleMainLoop()
}

func main() {
	listenAddress := flag.String("listen", "127.0.0.1:3802", "Address to listen on")
	flag.Parse()

	connectionMap := &sync.Map{}

	listen, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listen.Close()
	log.Printf("shard listening on :%s", *listenAddress)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(conn, connectionMap)
	}
}
