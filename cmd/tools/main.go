package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"starconflict/lib/auth"
	"starconflict/lib/bitreader"
	"starconflict/lib/cmdstore"
	"starconflict/lib/dbtypes"

	"github.com/alecthomas/kong"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func writeStoreCacheToDatabase(storeCachePath string, storeCol *mongo.Collection) {
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

	if _, err := storeCol.InsertMany(context.TODO(), store); err != nil {
		log.Fatalf("Failed to insert store: %v", err)
	}
}

func (db *DatabaseCmd) Run() error {

	if db.MongoURI == "" {
		return errors.New("no MongoDB URI provided, set --mongo-uri or the MONGODB_URI environment variable")
	}

	ctx := context.TODO()
	client, err := mongo.Connect(options.Client().ApplyURI(db.MongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from database: %v", err)
		}
	}()

	store := client.Database("cosmosim").Collection("store")

	writeStoreCacheToDatabase(db.StoreCachePath, store)
	return nil
}

type DatabaseCmd struct {
	StoreCachePath string `cmd:"" optional:"" help:"Path to store_cache.dat"`
	MongoURI       string `optional:"" env:"MONGODB_URI" help:"MongoDB connection URI"`
}

func nextUid(ctx context.Context, client *mongo.Client) (uint64, error) {
	counters := client.Database("cosmosim").Collection("counters")
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var counter dbtypes.Counters
	err := counters.FindOneAndUpdate(ctx, bson.M{"_id": "uid"}, bson.M{"$inc": bson.M{"c": int64(1)}}, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}
	return uint64(counter.C), nil
}

func (c *CreateAccountCmd) Run() error {
	if c.MongoURI == "" {
		return errors.New("no MongoDB URI provided, set --mongo-uri or the MONGODB_URI environment variable")
	}

	ctx := context.TODO()
	client, err := mongo.Connect(options.Client().ApplyURI(c.MongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from database: %v", err)
		}
	}()

	accounts := client.Database("cosmosim").Collection("accounts")

	err = accounts.FindOne(ctx, bson.M{"username": c.Username}).Err()
	if err == nil {
		return fmt.Errorf("an account with username %q already exists", c.Username)
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check for an existing account: %w", err)
	}

	hashed, err := auth.HashPassword([]byte(c.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	uid, err := nextUid(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to allocate uid: %w", err)
	}

	nickname := c.Nickname
	if nickname == "" {
		nickname = c.Username
	}

	account := dbtypes.Accounts{
		Uid:      uid,
		Puid:     uid,
		Username: c.Username,
		Password: hashed,
		NickName: nickname,
	}
	if _, err := accounts.InsertOne(ctx, account); err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	log.Printf("Created account %q (nickname %q) with uid %v", c.Username, nickname, uid)
	return nil
}

type CreateAccountCmd struct {
	Username string `arg:"" help:"Username used to log in"`
	Password string `arg:"" help:"Password for the account"`
	Nickname string `optional:"" help:"In-game nickname (defaults to the username)"`
	MongoURI string `optional:"" env:"MONGODB_URI" help:"MongoDB connection URI"`
}

var CLI struct {
	Db            DatabaseCmd      `cmd:"" help:"Database subcommand"`
	CreateAccount CreateAccountCmd `cmd:"" help:"Create a new account"`
}

func main() {
	ctx := kong.Parse(&CLI)
	if err := ctx.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
