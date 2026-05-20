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

func typeThreeResponse() []byte {
	return nil
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
		slog.Error("Unhandled CCmdStore Type", "Store Type", req.store_type)
	default:
		slog.Error("Invalid CCmdStore Request", "Store Type", req.store_type)
	}
}
