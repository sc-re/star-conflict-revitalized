package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"runtime/debug"
	"sync"

	"starconflict/lib/asyncreq"
	"starconflict/lib/auth"
	"starconflict/lib/protocol"
	"starconflict/lib/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE user (
	mail TEXT NOT NULL UNIQUE,
	nickname TEXT NOT NULL UNIQUE,
	uid INT8 NOT NULL PRIMARY KEY,
	zone INT2,
	password TEXT
);
`

type session struct {
	uid          uint64
	sessionToken string
	seq          uint16
	conn         net.Conn
	db           *sqlx.DB
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
			log.Printf("Got a %s", hdr.CommandType)
			slog.Info("Recieved", "commandType", hdr.CommandType)
			var valid bool
			var disconnectReason types.MasterServerDisconnectReason
			valid, session.uid, disconnectReason, err = auth.Authenticate(body, session.db, key, group)
			if valid {
				session.seq += 1
				auth.SendAuthAck(session.conn, session.db, session.seq, hdr.Sequence, session.uid)
				connectionMap.Store(session.uid, session)
			} else {
				session.conn.Write(protocol.MakeDisconnectMessage(disconnectReason))
			}
			if err != nil {
				slog.Warn("failed authentication", "error", err)
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

func handle(conn net.Conn, db *sqlx.DB, connectionMap *sync.Map) {
	defer conn.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("handle failed: %v", err)
			debug.PrintStack()
		}
	}()
	session := session{}
	session.db = db
	session.seq = 0
	session.conn = conn
	if err := session.handleAuthentication(connectionMap); err != nil {
		return
	}
	session.handleMainLoop()
}

func main() {
	listenAddress := flag.String("listen", "127.0.0.1:3802", "Address to listen on")
	logLevel := flag.String("loglevel", "info", "Set loglevel [debug/info/warn/error]")
	flag.Parse()

	connectionMap := &sync.Map{}

	slogLevel := slog.LevelDebug
	if err := slogLevel.UnmarshalText([]byte(*logLevel)); err != nil {
		log.Fatalf("Invalid logLevel %v: %v", logLevel, err)
	}
	slogOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, slogOptions))
	slog.SetDefault(logger)

	db, err := sqlx.Connect("sqlite3", "sc.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatalf("Failed to connect to databse %v", err)
	}
	defer db.Close()
	//db.MustExec(schema)

	listen, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer listen.Close()
	slog.Info("TCP Socket opened", "listenAddress", *listenAddress)

	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("accept failed", "error", err)
			continue
		}
		go handle(conn, db, connectionMap)
	}
}
