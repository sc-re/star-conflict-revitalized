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
	t := reflect.ValueOf(in)
	kind := t.Kind()
	switch kind {
	case reflect.Struct:
		if err := writeDict(bw, t); err != nil {
			return nil, err
		}
	case reflect.Map:
		if err := writeMap(bw, t); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("expected on of struct or map, got %s", t.Kind())
	}
	return bw.ReturnSlice(), nil
}

func writeMap(bw *bitwriter.Writer, in reflect.Value) error {
	len := in.Len()
	if err := bw.WriteBeUint32(uint32(len)); err != nil {
		return err
	}
	if len <= 0 {
		return nil
	}
	bw.WriteBool(false)
	iter := in.MapRange()
	for iter.Next() {
		if err := writeStringKey(bw, iter.Key().String(), ""); err != nil {
			return err
		}
		if err := writeValue(bw, iter.Value()); err != nil {
			return err
		}
	}
	return nil
}

func writeDict(bw *bitwriter.Writer, in reflect.Value) error {
	len := reflect.TypeOf(in).NumField()
	if err := bw.WriteBeUint32(uint32(len)); err != nil {
		return err
	}
	if len <= 0 {
		return nil
	}
	bw.WriteBool(false)
	for field := range in.Type().Fields() {
		tag := field.Tag.Get("variantdict")
		key := field.Name
		if err := writeStringKey(bw, key, tag); err != nil {
			return err
		}
		if err := writeValue(bw, in.FieldByIndex(field.Index)); err != nil {
			return nil
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

func writeStringKey(bw *bitwriter.Writer, key string, tag string) error {
	if tag != "" {
		bw.WriteCString(tag)
	} else {
		runes := []rune(key)
		runes[0] = unicode.ToLower(runes[0])
		bw.WriteCString(string(runes))
	}
	return nil
}

func writeIntKey(bw *bitwriter.Writer, key uint32) error {
	return nil
}

func writeValue(bw *bitwriter.Writer, v reflect.Value) error {
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		if err := writeBool(bw, v.Bool()); err != nil {
			return err
		}
	case reflect.Int32:
		if err := writeInt32(bw, int32(v.Int())); err != nil {
			return err
		}
	case reflect.Uint64:
		if err := writeUint64(bw, v.Uint()); err != nil {
			return err
		}
	case reflect.String:
		if err := writeString(bw, v.String()); err != nil {
			return err
		}
	case reflect.Slice:
		if err := writeSlice(bw, v); err != nil {
			return err
		}
	case reflect.Struct:
		if err := writeDict(bw, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled datatype: %v", v.Kind())
	}
	return nil
}

func writeInt32(bw *bitwriter.Writer, v int32) error {
	if err := bw.WriteByte(byte(tagI32)); err != nil {
		return err
	}
	if err := bw.WriteBeInt32(v); err != nil {
		return err
	}
	return nil
}

func writeUint64(bw *bitwriter.Writer, v uint64) error {
	if err := bw.WriteByte(byte(tagU64)); err != nil {
		return err
	}
	if err := bw.WriteBeUint64(v); err != nil {
		return err
	}
	return nil
}

func writeString(bw *bitwriter.Writer, v string) error {
	if err := bw.WriteByte(byte(tagStr)); err != nil {
		return err
	}
	if err := bw.WriteCString(v); err != nil {
		return err
	}
	return nil
}

func writeStruct(bw *bitwriter.Writer, v any) error {
	if err := bw.WriteByte(byte(tagDict)); err != nil {
		return err
	}
	if err := bw.WriteBool(false); err != nil {
		return err
	}
	return nil
}

func writeBool(bw *bitwriter.Writer, v bool) error {
	if err := bw.WriteByte(byte(tagBool)); err != nil {
		return err
	}
	if err := bw.WriteBool(v); err != nil {
		return err
	}
	return nil
}
