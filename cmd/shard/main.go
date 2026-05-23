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
	"starconflict/lib/cmdstore"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
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

func HandleAuthentication(session *session.Session, connectionMap *sync.Map) error {
	key, group := auth.CreateClientChallenge()
	if err := auth.SendChallenge(session.Conn, key, session.GetNextSeq()); err != nil {
		return err
	}
	for {
		hdr, body, err := protocol.ParseNextMessage(session.Conn)
		if err != nil {
			return err
		}
		if hdr.CommandType == types.CCMD_AUTH_REQUEST {
			slog.Info("Recieved", "commandType", hdr.CommandType)
			var valid bool
			var disconnectReason types.MasterServerDisconnectReason
			valid, session.Uid, disconnectReason, err = auth.Authenticate(body, session.Db, key, group)
			if valid {
				auth.SendAuthAck(session.Conn, session.Db, session.GetNextSeq(), hdr.Sequence, session.Uid)
				connectionMap.Store(session.Uid, session)
			} else {
				session.Conn.Write(protocol.MakeDisconnectMessage(disconnectReason))
			}
			if err != nil {
				slog.Warn("failed authentication", "error", err)
				return err
			}
			return nil
		}
		slog.Warn("Unhandled message", "type", hdr.CommandType)
	}
}

func HandleMainLoop(session *session.Session) error {
	//go userprofilenotification.SendUserProfileNotificationOnlineState(session.conn, session.uid, userprofilenotification.USER_STATE_ONLINE)
	go asyncreq.Send_ac_vessel_strip_improper_battle(session.Conn)
	go asyncreq.SendAcPlayerPush(session)
	for {
		hdr, body, err := protocol.ParseNextMessage(session.Conn)
		_ = body
		if err != nil {
			return err
		}
		if hdr.Special {
			// TODO: Do something
			continue
		}
		switch hdr.CommandType {
		case types.CSCMD_ASYNC_REQ:
			go asyncreq.HandleAsyncReq(hdr, body, session.GetNextSeq(), session)
			continue
		case types.SCMD_KEEP_ALIVE:
			slog.Warn("Recieved SCMD_KEEP_ALIVE")
			// TODO: Keep Alive
			continue
		case types.CCMD_STORE:
			go cmdstore.HandleCCmdStore(hdr, body, session.GetNextSeq(), session.Conn)
			continue
		}

		slog.Warn("Unhandled message", "type", hdr.CommandType)
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
	session := session.Session{}
	session.Db = db
	session.Conn = conn
	if err := HandleAuthentication(&session, connectionMap); err != nil {
		return
	}
	HandleMainLoop(&session)
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

	db, err := sqlx.Connect("sqlite3", "/home/john/Projects/star-conflict-revitalized/sc.db?cache=shared&mode=memory")
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
