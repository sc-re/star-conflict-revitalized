package variantdict

import (
	"fmt"
	"reflect"
	"starconflict/lib/bitreader"
	"starconflict/lib/bitwriter"
	"unicode"
)

type variantTags uint8

const (
	tagNil variantTags = iota
	tagI32
	tagU64
	tagunkwn
	tagF32
	tagStr
	tagDict
	tagunkwn2
	tagBool
)

func Marshal(in any) ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	if err := BwMarshal(bw, in); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

func BwMarshal(bw *bitwriter.Writer, in any) error {
	t := reflect.ValueOf(in)
	kind := t.Kind()
	switch kind {
	case reflect.Struct:
		if err := writeDict(bw, t); err != nil {
			return err
		}
	case reflect.Map:
		if err := writeMap(bw, t); err != nil {
			return err
		}
	case reflect.Slice:
		if err := writeSlice(bw, t); err != nil {
			return err
		}

	default:
		return fmt.Errorf("expected on of struct or map, got %s", t.Kind())
	}
	return nil
}

func Unmarshal(in []byte, out any) error {
	br := bitreader.NewReader(in)
	return BrUnmarshal(br, out)
}

func BrUnmarshal(br *bitreader.Reader, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("unmarshal: Expected non nil pointer for out, got: %v", v.Kind())
	}
	return readDict(br, v.Elem())
}

func writeMap(bw *bitwriter.Writer, in reflect.Value) error {
	length := in.Len()
	bw.WriteBeUint32(uint32(length))
	if length <= 0 {
		return nil
	}
	bw.WriteBool(false)
	iter := in.MapRange()
	for iter.Next() {
		bw.WriteCString(iter.Key().String())
		if err := writeValue(bw, iter.Value()); err != nil {
			return err
		}
	}
	return nil
}

func readStruct(br *bitreader.Reader, out reflect.Value, entryCount uint32) error {
	kind := out.Kind()
	if kind != reflect.Struct {
		return fmt.Errorf("v.Kind != struct, it is %v", kind)
	}
	for range entryCount {
		var field reflect.Value

		keyName, err := br.ReadCString()
		if err != nil {
			return err
		}
		runes := []rune(keyName)
		runes[0] = unicode.ToUpper(runes[0])
		keyName = string(runes)
		structField, fieldFound := out.Type().FieldByName(keyName)
		if !fieldFound {
			return fmt.Errorf("Failed to unmarshal: key %v not found in struct: %v", keyName, structField)
		}
		field = out.FieldByIndex(structField.Index)

		if err := readValue(br, field); err != nil {
			return err
		}
	}
	return nil
}

func readSlice(br *bitreader.Reader, out reflect.Value, entryCount uint32) error {
	kind := out.Kind()
	if kind != reflect.Slice {
		return fmt.Errorf("v.Kind != slice, it is %v", kind)
	}
	out.Grow(int(entryCount))
	elemType := out.Type().Elem()
	for range entryCount {
		v := reflect.New(elemType).Elem()
		if err := readValue(br, v); err != nil {
			return err
		}
		out.Set(reflect.Append(out, v))
		out = reflect.Append(out, v)
	}
	return nil
}

func readDict(br *bitreader.Reader, out reflect.Value) error {
	entryCount, err := br.ReadBeUint32()
	if err != nil {
		return err
	}
	indexedKeys, err := br.ReadBool()
	if err != nil {
		return err
	}
	if indexedKeys {
		return readSlice(br, out, entryCount)
	} else {
		return readStruct(br, out, entryCount)
	}
}

func writeDict(bw *bitwriter.Writer, in reflect.Value) error {
	length := in.NumField()
	bw.WriteBeUint32(uint32(length))
	if length <= 0 {
		return nil
	}
	bw.WriteBool(false)
	for field := range in.Type().Fields() {
		tag := field.Tag.Get("variantdict")
		key := field.Name
		writeStringKey(bw, key, tag)
		if err := writeValue(bw, in.FieldByIndex(field.Index)); err != nil {
			return err
		}
	}
	return nil
}

func writeSlice(bw *bitwriter.Writer, in reflect.Value) error {
	bw.WriteBeUint32(uint32(in.Len()))
	bw.WriteBool(true)
	for i := 0; i < in.Len(); i++ {
		if err := writeValue(bw, in.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func writeStringKey(bw *bitwriter.Writer, key string, tag string) {
	if tag != "" {
		bw.WriteCString(tag)
	} else {
		runes := []rune(key)
		runes[0] = unicode.ToLower(runes[0])
		bw.WriteCString(string(runes))
	}
}

func readString(br *bitreader.Reader, v reflect.Value) error {
	str, err := br.ReadCString()
	if err != nil {
		return err
	}
	v.SetString(str)
	return nil
}

func readFloat32(br *bitreader.Reader, v reflect.Value) error {
	f, err := br.ReadFloat32()
	if err != nil {
		return err
	}
	v.SetFloat(float64(f))
	return nil
}

func readBool(br *bitreader.Reader, v reflect.Value) error {
	b, err := br.ReadBool()
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

func readInt32(br *bitreader.Reader, v reflect.Value) error {
	i, err := br.ReadBeInt32()
	if err != nil {
		return err
	}
	v.SetInt(int64(i))
	return nil
}

func readUint64(br *bitreader.Reader, v reflect.Value) error {
	u, err := br.ReadBeUint64()
	if err != nil {
		return err
	}
	v.SetUint(u)
	return nil
}

func readValue(br *bitreader.Reader, v reflect.Value) error {
	tmp, err := br.ReadByte()
	if err != nil {
		return err
	}
	kind := v.Kind()
	tag := variantTags(tmp)
	switch tag {
	case tagDict:
		switch kind {
		case reflect.Struct:
			fallthrough
		case reflect.Slice:
			return readDict(br, v)
		default:
			return fmt.Errorf("Got tag Dict, struct expetcts %v", kind)
		}

	case tagF32:
		if kind != reflect.Float32 {
			return fmt.Errorf("Got tag Flaot32, struct expects %v", kind)
		}
		return readFloat32(br, v)
	case tagStr:
		if kind != reflect.String {
			return fmt.Errorf("Got tag string, struct expects %v", kind)
		}
		return readString(br, v)
	case tagBool:
		if kind != reflect.Bool {
			return fmt.Errorf("Got tag bool, struct expects %v", kind)
		}
		return readBool(br, v)
	case tagI32:
		if kind != reflect.Int32 {
			return fmt.Errorf("Got tag int32, struct expects %v", kind)
		}
		return readInt32(br, v)
	case tagU64:
		if kind != reflect.Uint64 {
			return fmt.Errorf("Got tag uint64, struct expects %v", kind)
		}
		return readUint64(br, v)
	default:
		return fmt.Errorf("Unhandled tag type: %v", tag)
	}
}

func writeValue(bw *bitwriter.Writer, v reflect.Value) error {
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		writeBool(bw, v.Bool())
	case reflect.Int32:
		writeInt32(bw, int32(v.Int()))
	case reflect.Uint64:
		writeUint64(bw, v.Uint())
	case reflect.Float32:
		writeFloat32(bw, float32(v.Float()))
	case reflect.String:
		writeString(bw, v.String())
	case reflect.Slice:
		if err := writeNestedSlice(bw, v); err != nil {
			return err
		}
	case reflect.Struct:
		if err := writeNestedDict(bw, v); err != nil {
			return err
		}
	case reflect.Map:
		if err := writeNestedMap(bw, v); err != nil {
			return err
		}
	case reflect.Pointer:
		if err := writeValue(bw, v.Elem()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled datatype: %v", v.Kind())
	}
	return nil
}

func writeInt32(bw *bitwriter.Writer, v int32) {
	bw.BwWriteByte(byte(tagI32))
	bw.WriteBeInt32(v)
}

func writeUint64(bw *bitwriter.Writer, v uint64) {
	bw.BwWriteByte(byte(tagU64))
	bw.WriteBeUint64(v)
}

func writeFloat32(bw *bitwriter.Writer, v float32) {
	bw.BwWriteByte(byte(tagF32))
	bw.WriteFloat32(v)
}

func writeString(bw *bitwriter.Writer, v string) {
	bw.BwWriteByte(byte(tagStr))
	bw.WriteCString(v)
}

func writeNestedMap(bw *bitwriter.Writer, v reflect.Value) error {
	bw.BwWriteByte(byte(tagDict))
	if err := writeMap(bw, v); err != nil {
		return err
	}
	return nil
}

func writeNestedDict(bw *bitwriter.Writer, v reflect.Value) error {
	bw.BwWriteByte(byte(tagDict))
	if err := writeDict(bw, v); err != nil {
		return err
	}
	return nil
}

func writeNestedSlice(bw *bitwriter.Writer, v reflect.Value) error {
	bw.BwWriteByte(byte(tagDict))
	if err := writeSlice(bw, v); err != nil {
		return err
	}
	return nil
}

func writeBool(bw *bitwriter.Writer, v bool) {
	bw.BwWriteByte(byte(tagBool))
	bw.WriteBool(v)
}
