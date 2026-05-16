package variantdict

import (
	"fmt"
	"reflect"
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
	default:
		return fmt.Errorf("expected on of struct or map, got %s", t.Kind())
	}
	return nil
}

func writeMap(bw *bitwriter.Writer, in reflect.Value) error {
	len := in.Len()
	bw.WriteBeUint32(uint32(len))
	if len <= 0 {
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

func writeDict(bw *bitwriter.Writer, in reflect.Value) error {
	len := reflect.TypeOf(in).NumField()
	bw.WriteBeUint32(uint32(len))
	if len <= 0 {
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
	for i := 0; i < in.Len(); i++ {
		if err := writeIntKey(bw, uint32(i)); err != nil {
			return err
		}
		if err := writeValue(bw, reflect.ValueOf(in.Index(i))); err != nil {
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

func writeIntKey(bw *bitwriter.Writer, key uint32) error {
	return nil
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
		if err := writeSlice(bw, v); err != nil {
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

func writeBool(bw *bitwriter.Writer, v bool) {
	bw.BwWriteByte(byte(tagBool))
	bw.WriteBool(v)
}
