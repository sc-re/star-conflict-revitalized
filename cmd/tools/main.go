package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"starconflict/lib/bitreader"
	"starconflict/lib/cmdstore"

	"github.com/alecthomas/kong"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func readStoreToByte(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open %v: %v", path, err)
	}
	defer f.Close()
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatalf("Failed to get filesize: %v", err)
	}
	log.Printf("File size: %v", size)
	f.Seek(0, io.SeekStart)
	buf := make([]byte, size)
	if _, err := f.Read(buf); err != nil {
		log.Fatalf("Failed to read file into buffer: %v", err)
	}
	return buf
}

func parseStore(store []byte) []cmdstore.StoreItem {
	ret := []cmdstore.StoreItem{}
	br := bitreader.NewReader(store[12:])
	itemCount, _ := br.ReadBeUint32()
	for range itemCount {
		var err error
		item := cmdstore.StoreItem{}
		item.StoreItemId, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.CreditPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.PremiumPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.TokenPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.EventPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.BaseCreditsPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.BasePremiumPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.BaseTokenPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.BaseEventPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.TradePrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.TradePremiumPrice, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.Race, err = br.ReadByte()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.RequiredRank, err = br.ReadBeUint32()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.Stacks, err = br.ReadBool()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.CantBeBought, err = br.ReadBool()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		if b, err := br.ReadByte(); err != nil {
			log.Fatalf("Failed to read: %v", err)
		} else {
			item.ItemType = cmdstore.ItemType(b)
		}
		item.ItemName, err = br.ReadCString()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.ItemFlags, err = br.ReadByte()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.RequiredAccountAura, err = br.ReadCString()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		item.DeleteFromInventory, err = br.ReadBool()
		if err != nil {
			log.Fatalf("Failed to read: %v", err)
		}
		ret = append(ret, item)
	}

	return ret
}

func writeStoreToDatabase(store []cmdstore.StoreItem, databasePath string) {
	db, err := sqlx.Connect("sqlite3", databasePath+"?cache=shared&mode=memory")
	if err != nil {
		log.Fatalf("Failed to connect to databse %v", err)
	}
	defer db.Close()
	s := 0
	end := 0
	for s < len(store) {
		end += 1000
		if end > len(store) {
			end = len(store)
		}
		log.Printf("Inserting store[%v:%v]", s, end)
		_, err = db.NamedExec(`INSERT INTO store (StoreItemId, CreditPrice, PremiumPrice, TokenPrice, EventPrice, BaseCreditsPrice, BasePremiumPrice, BaseTokenPrice, BaseEventPrice, TradePrice, TradePremiumPrice, Race, RequiredRank, Stacks, CantBeBought, ItemType, ItemName, ItemFlags, RequiredAccountAura, DeleteFromInventory) VALUES (:StoreItemId, :CreditPrice, :PremiumPrice, :TokenPrice, :EventPrice, :BaseCreditsPrice, :BasePremiumPrice, :BaseTokenPrice, :BaseEventPrice, :TradePrice, :TradePremiumPrice, :Race, :RequiredRank, :Stacks, :CantBeBought, :ItemType, :ItemName, :ItemFlags, :RequiredAccountAura, :DeleteFromInventory)`, store[s:end])
		if err != nil {
			log.Fatalf("Failed to insert store: %v", err)
		}
		s = end
	}
}

func writeStoreCacheToDatabase(storeCachePath string, databasePath string) {
	if storeCachePath == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get homedir: %v", err)
		}
		storeCachePath = filepath.Join(homedir, ".local/share/starconflict/store_cache.dat")
	}
	storeBytes := readStoreToByte(storeCachePath)
	store := parseStore(storeBytes)
	bruh := map[uint32]bool{}
	for _, v := range store {
		if bruh[v.StoreItemId] {
			log.Fatalf("Duplicate store item: %v", v.StoreItemId)
		}
		bruh[v.StoreItemId] = true
	}

	writeStoreToDatabase(store, databasePath)
}

func (db *DatabaseCmd) Run() error {
	writeStoreCacheToDatabase(db.StoreCachePath, db.DatabasePath)
	return nil
}

type DatabaseCmd struct {
	StoreCachePath string `cmd:"" optional:"" help:"Path to store_cache.dat"`
	DatabasePath   string `cmd:"" optional:"" default:"sc.db" env:"SC_DB_PATH" help:"Path to sqlite Database"`
}

var CLI struct {
	Db DatabaseCmd `cmd:"" help:"Database subcommand"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if err := ctx.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
