package cmdstore

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"log/slog"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	catalog []byte
	once    sync.Once
)

const (
	catalogHash uint32 = 0xc9c95436
)

type ccmdStoreReq struct {
	store_type byte
	chunkCount uint16
	idk        uint32
	idk2       uint16
}

type ItemType byte

const (
	ItemTypeVessel ItemType = iota
	ItemTypeModule
	ItemTypeRocket
	ItemTypeDrug
	ItemTypeBundle
	ItemTypeResource
	ItemTypeBlueprint
	ItemTypeAvatar
	ItemTypeMotto
	ItemTypeJunk
)

type StoreItem struct {
	StoreItemId         uint32   `db:"StoreItemId"`
	CreditPrice         uint32   `db:"CreditPrice"`
	PremiumPrice        uint32   `db:"PremiumPrice"`
	TokenPrice          uint32   `db:"TokenPrice"`
	EventPrice          uint32   `db:"EventPrice"`
	BaseCreditsPrice    uint32   `db:"BaseCreditsPrice"`
	BasePremiumPrice    uint32   `db:"BasePremiumPrice"`
	BaseTokenPrice      uint32   `db:"BaseTokenPrice"`
	BaseEventPrice      uint32   `db:"BaseEventPrice"`
	TradePrice          uint32   `db:"TradePrice"`
	TradePremiumPrice   uint32   `db:"TradePremiumPrice"`
	Race                byte     `db:"Race"`
	RequiredRank        uint32   `db:"RequiredRank"`
	Stacks              bool     `db:"Stacks"`
	CantBeBought        bool     `db:"CantBeBought"`
	ItemType            ItemType `db:"ItemType"`
	ItemName            string   `db:"ItemName"`
	ItemFlags           byte     `db:"ItemFlags"`
	RequiredAccountAura string   `db:"RequiredAccountAura"`
	DeleteFromInventory bool     `db:"DeleteFromInventory"`
}

func initCatalogStore(db *sqlx.DB) {
	once.Do(func() {
		storeItems := []StoreItem{}
		db.Select(&storeItems, "SELECT * FROM store")
		tmp, _ := deflateStoreItems(storeItems)
		catalog = tmp.Bytes()
		slog.Debug("Fetched catalog from databse", "len(storeItems)", len(storeItems), "len(catalog)", len(catalog))
	})
}

func (req *ccmdStoreReq) parse(body []byte) {
	req.store_type = body[0]
	switch req.store_type {
	case 1:
		fallthrough
	case 2:
		req.chunkCount = binary.BigEndian.Uint16(body[1:])
		req.idk = binary.BigEndian.Uint32(body[3:])
		req.idk2 = binary.BigEndian.Uint16(body[7:])
	}
}

type goldPack struct {
	hidden       bool
	value        uint32
	idk          uint32
	licenseHours uint32
	weight       float32
	guid         string
	storeUrl     string
	arcCost      float32
	displayPrice float32
	idkCost      float32
	lang         string
	// ...
}

type creditPack struct {
	price      uint32
	value      uint32
	bonusValue uint32
	name       string
}

func reditPacks() []creditPack {
	return []creditPack{
		{20, 140000, 0, "TinyCreditPack"},
		{100, 702000, 38000, "SmallCreditPack"},
		{1000, 7000000, 700000, "MediumCreditPack"},
		{2500, 17500000, 2400000, "BigCreditPack"},
		{5000, 35000000, 6200000, "HudgeCreditPack"},
		{10000, 70000000, 17500000, "GreatCreditPack"},
	}
}

func writeIdkTable(bw *bitwriter.Writer) {
	table := [][3]uint32{
		{75, 25000, 0},
		{100, 25000, 0},
		{125, 50000, 0},
		{150, 75000, 0},
		{200, 100000, 0},
		{250, 100000, 0},
		{300, 100000, 0},
		{350, 100000, 0},
		{400, 200000, 0},
		{450, 200000, 0},
		{500, 300000, 0},
		{550, 300000, 0},
		{600, 400000, 0},
		{650, 400000, 0},
		{700, 500000, 0},
		{750, 500000, 0},
		{800, 600000, 0},
		{850, 600000, 0},
		{900, 700000, 0},
		{950, 700000, 0},
		{1000, 800000, 0},
		{1050, 800000, 0},
		{1100, 900000, 0},
		{1300, 1500000, 0},
		{1500, 2000000, 0},
		{1750, 3000000, 0},
		{2000, 4000000, 0},
		{2250, 5000000, 0},
		{2500, 6000000, 0},
		{2750, 0, 30000},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295}}
	for _, v := range table {
		for _, i := range v {
			bw.WriteBeUint32(i)
		}
	}
}

func writeIdkOtherTable(bw *bitwriter.Writer) {
	table := [][3]uint32{
		{20, 1000000, 0},
		{30, 2000000, 0},
		{40, 4000000, 0},
		{50, 8000000, 0},
		{60, 10000000, 0},
		{70, 15000000, 0},
		{80, 20000000, 0},
		{90, 30000000, 0},
		{100, 40000000, 0},
		{110, 60000000, 0},
		{120, 80000000, 0},
		{130, 100000000, 0},
		{140, 150000000, 0},
		{150, 200000000, 0},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
		{0, 4294967295, 4294967295},
	}
	for _, v := range table {
		for _, i := range v {
			bw.WriteBeUint32(i)
		}
	}
}

func typeThreeResponse() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 27110))
	bw.BwWriteByte(3)
	bw.WriteBeUint32(10000) // idk
	bw.WriteCString("'Bundle_Chest_Gold_Reward'")
	bw.WriteBeUint32(400)
	bw.WriteBeUint32(1000)
	bw.WriteFloat32(0.10000000149011612)
	bw.WriteFloat32(0.5)
	bw.WriteFloat32(0.0005499999970197678)
	packs := reditPacks()
	bw.WriteBeUint32(uint32(len(packs)))
	for _, v := range packs {
		bw.WriteBeUint32(v.price)
		bw.WriteBeUint32(v.value)
		bw.WriteBeUint32(v.bonusValue)
		bw.WriteCString(v.name)
	}
	bw.WriteBeUint32(0)                  // gold packs count
	bw.WriteFloat32(0.15000000596046448) // idk
	bw.WriteFloat32(0.04000000283122063)
	bw.WriteFloat32(0.02500000223517418)
	bw.WriteFloat32(0.020000001415610313)
	bw.WriteFloat32(0.06599999964237213)
	bw.WriteFloat32(1.0)
	bw.WriteBeUint32(3000)
	bw.WriteBeUint32(500)
	bw.WriteFloat32(2000)
	writeIdkTable(bw)
	writeIdkOtherTable(bw)
	bw.BwWriteByte(1) // idk
	bw.WriteFloat32(0.003333332948386669)
	bw.WriteFloat32(1.0)
	bw.WriteBeUint32(1500)
	bw.WriteBeUint32(1000)
	bw.WriteBeUint32(1000)
	bw.WriteFloat32(0.0)
	arr := []uint32{1000000, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500, 1500}
	for _, v := range arr {
		bw.WriteBeUint32(v)
	}
	arr = []uint32{1000000, 3000000, 3000000, 3000000, 3000000, 3000000, 1000000, 1000000, 1000000, 1000000, 1000000, 1000000, 1000000, 1000000}
	for _, v := range arr {
		bw.WriteBeUint32(v)
	}
	sarr := []uint32{2, 1, 2, 2, 200000, 2, 5}
	for _, v := range sarr {
		bw.WriteBeUint32(v)
	}
	bw.WriteFloat32(95.0)
	bw.WriteFloat32(45.0)
	for range 25 {
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
	}
	bw.WriteFloat32(0.10000000149011612)
	bw.WriteFloat32(0.800000011920929)
	// three timestamps
	bw.WriteBeUint64(1772168400000)
	bw.WriteBeUint64(1780376400000)
	bw.WriteBeUint64(1777525200000)
	bw.WriteBeInt32(20000)
	bw.WriteBool(false)
	bw.WriteBeInt32(3)
	bw.WriteBeInt32(2)
	// XXX

	return bw.ReturnSlice()
}

func deflateStoreItems(si []StoreItem) (*bytes.Buffer, error) {
	tmp := &bytes.Buffer{}
	bw := bitwriter.NewWriter(make([]byte, 0, 4096))
	bw.WriteBeUint32(uint32(len(si)))
	for _, v := range si {
		serializeStoreItem(v, bw)
	}
	binary.Write(tmp, binary.BigEndian, uint32(len(bw.ReturnSlice())))
	deflate, _ := flate.NewWriter(tmp, flate.BestCompression)
	_, err := deflate.Write(bw.ReturnSlice()) // others call this an inflated empty list
	if err != nil {
		slog.Error("Failed to write to deflate", "err", err)
		return nil, err
	}
	if err := deflate.Close(); err != nil {
		slog.Error("Failed to flush deflate", "err", err)
		return nil, err
	}
	return tmp, nil
}

func serializeStoreItem(si StoreItem, bw *bitwriter.Writer) {
	bw.WriteBeUint32(si.StoreItemId)
	bw.WriteBeUint32(si.CreditPrice)
	bw.WriteBeUint32(si.PremiumPrice)
	bw.WriteBeUint32(si.TokenPrice)
	bw.WriteBeUint32(si.EventPrice)
	bw.WriteBeUint32(si.BaseCreditsPrice)
	bw.WriteBeUint32(si.BasePremiumPrice)
	bw.WriteBeUint32(si.BaseTokenPrice)
	bw.WriteBeUint32(si.BaseEventPrice)
	bw.WriteBeUint32(si.TradePrice)
	bw.WriteBeUint32(si.TradePremiumPrice)
	bw.BwWriteByte(si.Race)
	bw.WriteBeUint32(si.RequiredRank)
	bw.WriteBool(si.Stacks)
	bw.WriteBool(si.CantBeBought)
	bw.BwWriteByte(byte(si.ItemType))
	bw.WriteCString(si.ItemName)
	bw.BwWriteByte(si.ItemFlags)
	bw.WriteCString(si.RequiredAccountAura)
	bw.WriteBool(si.DeleteFromInventory)
}

func (req *ccmdStoreReq) typeOneResponse() []byte {
	ret := make([]byte, 0, 65545)
	chunkSize := 65536
	begin := min((chunkSize * int(req.chunkCount)), len(catalog))
	end := min((chunkSize * int(req.chunkCount+1)), len(catalog))
	ret = append(ret, 1)
	ret = binary.BigEndian.AppendUint32(ret, uint32(len(catalog)))
	ret = binary.BigEndian.AppendUint32(ret, uint32(begin))
	ret = append(ret, catalog[begin:end]...)
	return ret
}

func typeTwoResponse() []byte {

	entries := []StoreItem{
		{
			StoreItemId:         210170,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          2400,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      2400,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_09",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210169,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          2400,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      2400,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_08",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210168,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          2400,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      2400,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_07",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210167,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          600,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      600,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_06",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210166,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          600,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      600,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_05",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210165,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          900,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      900,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_04",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210164,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          640,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      640,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_03",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210163,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          200,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      200,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_02",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210162,
			CreditPrice:         0,
			PremiumPrice:        0,
			TokenPrice:          150,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    0,
			BaseTokenPrice:      150,
			BaseEventPrice:      0,
			TradePrice:          1,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_CorpEconomy_2_01",
			ItemFlags:           8,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210544,
			CreditPrice:         0,
			PremiumPrice:        325,
			TokenPrice:          0,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    325,
			BaseTokenPrice:      0,
			BaseEventPrice:      0,
			TradePrice:          13000,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_Parts_T5_Karud_Unique",
			ItemFlags:           0,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210545,
			CreditPrice:         0,
			PremiumPrice:        1460,
			TokenPrice:          0,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    1460,
			BaseTokenPrice:      0,
			BaseEventPrice:      0,
			TradePrice:          58400,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_Parts_T5_Karud_Unique_x5",
			ItemFlags:           0,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210542,
			CreditPrice:         0,
			PremiumPrice:        228,
			TokenPrice:          0,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    325,
			BaseTokenPrice:      0,
			BaseEventPrice:      0,
			TradePrice:          9120,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_Parts_Race2_L_T5_Francisca",
			ItemFlags:           12,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
		{StoreItemId: 210543,
			CreditPrice:         0,
			PremiumPrice:        1022,
			TokenPrice:          0,
			EventPrice:          0,
			BaseCreditsPrice:    0,
			BasePremiumPrice:    1460,
			BaseTokenPrice:      0,
			BaseEventPrice:      0,
			TradePrice:          40880,
			TradePremiumPrice:   0,
			Race:                5,
			RequiredRank:        0,
			Stacks:              false,
			CantBeBought:        false,
			ItemType:            4,
			ItemName:            "Bundle_Parts_Race2_L_T5_Francisca_x5",
			ItemFlags:           0,
			RequiredAccountAura: "",
			DeleteFromInventory: false},
	}

	tmp, _ := deflateStoreItems(entries)

	ret := make([]byte, 0, 9+tmp.Len())
	ret = append(ret, 2)
	ret = binary.BigEndian.AppendUint32(ret, uint32(tmp.Len()))
	ret = binary.BigEndian.AppendUint32(ret, 0)
	ret = append(ret, tmp.Bytes()...)

	return ret
}

func hashResponse() []byte {
	return []byte{0x00, 0xc9, 0xc9, 0x54, 0x36}
}

func HandleCCmdStore(hdr *protocol.Header, body []byte, seq uint16, session *session.Session) {
	req := ccmdStoreReq{}
	initCatalogStore(session.Db)
	req.parse(body)
	switch req.store_type {
	case 0:
		resp := hashResponse()
		session.Conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	case 1:
		resp := req.typeOneResponse()
		session.Conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	case 2:
		resp := typeTwoResponse()
		session.Conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	case 3:
		resp := typeThreeResponse()
		session.Conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	default:
		slog.Error("Invalid CCmdStore Request", "Store Type", req.store_type)
	}
}
