package cmdstore

import (
	"bytes"
	"compress/flate"
	_ "compress/zlib"
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

const (
	catalogHash uint32 = 0xc9c95436
)

type ccmdStoreReq struct {
	store_type byte
	idk        uint64
}

type itemType byte

const (
	itemTypeVessel itemType = iota
	itemTypeModule
	itemTypeRocket
	itemTypeDrug
	itemTypeBundle
	itemTypeResource
	itemTypeBlueprint
	itemTypeAvatar
	itemTypeMotto
	itemTypeJunk
)

type storeItem struct {
	storeItemId         uint32
	creditPrice         uint32
	premiumPrice        uint32
	tokenPrice          uint32
	eventPrice          uint32
	baseCreditsPrice    uint32
	basePremiumPrice    uint32
	baseTokenPrice      uint32
	baseEventPrice      uint32
	tradePrice          uint32
	tradePremiumPrice   uint32
	race                byte
	requiredRank        uint32
	stacks              bool
	cantBeBought        bool
	itemType            itemType
	itemName            string
	itemFlags           byte
	requiredAccountAura string
	deleteFromInventory bool
}

func (req *ccmdStoreReq) parse(body []byte) {
	req.store_type = body[0]
	if req.store_type == 2 {
		req.idk = binary.BigEndian.Uint64(body[1:])
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

func serializeStoreItem(si storeItem, bw *bitwriter.Writer) {
	bw.WriteBeUint32(si.storeItemId)
	bw.WriteBeUint32(si.creditPrice)
	bw.WriteBeUint32(si.premiumPrice)
	bw.WriteBeUint32(si.tokenPrice)
	bw.WriteBeUint32(si.eventPrice)
	bw.WriteBeUint32(si.baseCreditsPrice)
	bw.WriteBeUint32(si.basePremiumPrice)
	bw.WriteBeUint32(si.baseTokenPrice)
	bw.WriteBeUint32(si.baseEventPrice)
	bw.WriteBeUint32(si.tradePrice)
	bw.WriteBeUint32(si.tradePremiumPrice)
	bw.BwWriteByte(si.race)
	bw.WriteBeUint32(si.requiredRank)
	bw.WriteBool(si.stacks)
	bw.WriteBool(si.cantBeBought)
	bw.BwWriteByte(byte(si.itemType))
	bw.WriteCString(si.itemName)
	bw.BwWriteByte(si.itemFlags)
	bw.WriteCString(si.requiredAccountAura)
	bw.WriteBool(si.deleteFromInventory)
}

func typeTwoResponse() []byte {
	var tmp bytes.Buffer

	entries := []storeItem{
		{
			storeItemId:         210170,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          2400,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      2400,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_09",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210169,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          2400,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      2400,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_08",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210168,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          2400,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      2400,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_07",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210167,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          600,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      600,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_06",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210166,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          600,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      600,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_05",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210165,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          900,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      900,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_04",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210164,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          640,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      640,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_03",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210163,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          200,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      200,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_02",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210162,
			creditPrice:         0,
			premiumPrice:        0,
			tokenPrice:          150,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    0,
			baseTokenPrice:      150,
			baseEventPrice:      0,
			tradePrice:          1,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_CorpEconomy_2_01",
			itemFlags:           8,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210544,
			creditPrice:         0,
			premiumPrice:        325,
			tokenPrice:          0,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    325,
			baseTokenPrice:      0,
			baseEventPrice:      0,
			tradePrice:          13000,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_Parts_T5_Karud_Unique",
			itemFlags:           0,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210545,
			creditPrice:         0,
			premiumPrice:        1460,
			tokenPrice:          0,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    1460,
			baseTokenPrice:      0,
			baseEventPrice:      0,
			tradePrice:          58400,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_Parts_T5_Karud_Unique_x5",
			itemFlags:           0,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210542,
			creditPrice:         0,
			premiumPrice:        228,
			tokenPrice:          0,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    325,
			baseTokenPrice:      0,
			baseEventPrice:      0,
			tradePrice:          9120,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_Parts_Race2_L_T5_Francisca",
			itemFlags:           12,
			requiredAccountAura: "",
			deleteFromInventory: false},
		{storeItemId: 210543,
			creditPrice:         0,
			premiumPrice:        1022,
			tokenPrice:          0,
			eventPrice:          0,
			baseCreditsPrice:    0,
			basePremiumPrice:    1460,
			baseTokenPrice:      0,
			baseEventPrice:      0,
			tradePrice:          40880,
			tradePremiumPrice:   0,
			race:                5,
			requiredRank:        0,
			stacks:              false,
			cantBeBought:        false,
			itemType:            4,
			itemName:            "Bundle_Parts_Race2_L_T5_Francisca_x5",
			itemFlags:           0,
			requiredAccountAura: "",
			deleteFromInventory: false},
	}

	bw := bitwriter.NewWriter(make([]byte, 0, 200))
	bw.WriteBeUint32(uint32(len(entries)))
	for _, v := range entries {
		serializeStoreItem(v, bw)
	}

	deflate, _ := flate.NewWriter(&tmp, flate.BestCompression)
	_, err := deflate.Write(bw.ReturnSlice()) // others call this an inflated empty list
	if err != nil {
		slog.Error("Failed to write to deflate", "err", err)
		return nil
	}
	if err := deflate.Close(); err != nil {
		slog.Error("Failed to flush deflate", "err", err)
		return nil
	}

	ret := make([]byte, 0, 13+tmp.Len())
	ret = append(ret, 2)
	ret = binary.BigEndian.AppendUint32(ret, uint32(tmp.Len()+4))
	ret = binary.BigEndian.AppendUint32(ret, 0)
	ret = binary.BigEndian.AppendUint32(ret, uint32(len(bw.ReturnSlice())))
	ret = append(ret, tmp.Bytes()...)

	return ret
}

func hashResponse() []byte {
	return []byte{0x00, 0xc9, 0xc9, 0x54, 0x36}
}

func HandleCCmdStore(hdr *protocol.Header, body []byte, seq uint16, conn net.Conn) {
	req := ccmdStoreReq{}
	req.parse(body)
	switch req.store_type {
	case 0:
		resp := hashResponse()
		conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	case 1:
		slog.Error("Unhandled CCmdStore Type", "Store Type", req.store_type)
	case 2:
		resp := typeTwoResponse()
		conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	case 3:
		resp := typeThreeResponse()
		conn.Write(protocol.MakeMessage(types.SCMD_STORE, seq, hdr.Sequence, resp))
	default:
		slog.Error("Invalid CCmdStore Request", "Store Type", req.store_type)
	}
}
