// Package riff reads and umarshalls RIFF files
package riff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// A Reader is a bytes reader with some helper functions to read IDs, Lens, and Data
// from RIFF files.
type Reader struct {
	*bytes.Reader
}

// NewReader returns an initial Reader
func NewReader(data []byte) *Reader {
	return &Reader{
		Reader: bytes.NewReader(data),
	}
}

// Print prints a reader without any knowledge of the structure of the reader,
// so all values will be []bytes.
// It assumes the reader has not advanced at all. Todo: Change that
func (r *Reader) Print() {
	deepPrint(r, " ", -1)
}

func deepPrint(r *Reader, prefix string, readLimit int) {
	var err error
	var typ string
	var l uint32
	var data []byte
	var isList bool
	var read int
	for err == nil && readLimit == -1 || read < readLimit {
		typ, l, isList, err = r.NextIDLen()
		// There will be a bogus byte at the end of some prints.
		if l%2 != 0 {
			l++
		}
		read += 8
		if err == nil {
			fmt.Print(prefix, typ, " Length:", l)
			if isList {
				typ2, err2 := r.NextID()
				read += 4
				if err2 == nil {
					fmt.Println(prefix+"  ", typ2)
					deepPrint(r, prefix+"    ", int(l)-4)
				} else {
					fmt.Println(prefix, err2)
				}
			} else if l < 40 {
				data = make([]byte, l)
				r.Read(data)
				fmt.Println(" Content:", data)
			} else {
				r.Seek(int64(l), io.SeekCurrent)
				fmt.Println(" Long Content")
			}
			read += int(l)
		}
	}
	if err != nil && err != io.EOF {
		fmt.Println(prefix, err)
	}
}

// Unmarshal is a mirror of json.Unmarshal, for RIFF files
func Unmarshal(data []byte, v interface{}) error {
	return NewReader(data).unmarshal(v)
}

func (r *Reader) unmarshal(v interface{}) error {
	// Mirrors json.unmarshal
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Invalid Unmarshal Struct")
	}
	// The first ID in the riff should be RIFF
	id, err := r.NextID()
	if err != nil {
		return err
	}
	if id != "RIFF" {
		return errors.New("RIFF format must begin with RIFF")
	}
	ln, err := r.NextLen()
	if err != nil {
		return err
	}
	// The next ID identifies this file type. We don't want it.
	_, err = r.NextID()
	if err != nil {
		return err
	}
	_, err = r.chunks(reflect.Indirect(rv), int(ln))
	return err
}

// NextID returns the next four byte sof the reader as a string
func (r *Reader) NextID() (string, error) {
	id := make([]byte, 4)
	l, err := r.Reader.Read(id)
	if l != 4 || err != nil {
		return "", errors.New("RIFF missing expected ID")
	}
	return string(id), nil
}

// NextIDLen returns NextID and NextLen
func (r *Reader) NextIDLen() (string, uint32, bool, error) {
	id, err := r.NextID()
	if err != nil {
		return "", 0, false, err
	}
	ln, err := r.NextLen()
	if err != nil {
		return "", 0, false, err
	}
	return id, ln, id == "LIST" || id == "RIFF", nil
}

// NextLen returns the next four bytes of the reader as a length.
func (r *Reader) NextLen() (uint32, error) {
	var ln uint32
	err := binary.Read(r.Reader, binary.LittleEndian, &ln)
	if err != nil {
		return ln, errors.New("RIFF missing expected length")
	}
	return ln, nil
}

func (r *Reader) chunks(rv reflect.Value, inLength int) (reflect.Value, error) {
	// Find chunkId in rv
	// If it can't be found, ignore it as a value the user does not want
	switch rv.Kind() {
	case reflect.Struct:
		return rv, r.structChunks(rv, inLength)
	case reflect.Slice:
		return r.sliceChunks(rv, inLength)
	default:
		return reflect.Value{}, errors.New("Unsupported unmarshal type")
	}
}

func (r *Reader) sliceChunks(rv reflect.Value, inLength int) (reflect.Value, error) {

	slTy := rv.Type()
	ty := slTy.Elem()
	newSlice := reflect.MakeSlice(slTy, 0, 10000)
	for inLength > 0 {
		_, ln, isList, err := r.NextIDLen()
		if err != nil {
			return reflect.Value{}, err
		}
		if !isList {
			return reflect.Value{}, errors.New("Slice structs need to be LISTs")
		}
		ln -= 4
		inLength -= 4
		if inLength <= 0 {
			break
		}
		_, err = r.NextID()
		if err != nil {
			return reflect.Value{}, err
		}

		inLength -= 8
		if inLength <= 0 {
			break
		}
		newStruct := reflect.New(ty)
		err = r.structChunks(reflect.Indirect(newStruct), int(ln))
		if err != nil {
			return reflect.Value{}, err
		}
		newSlice = reflect.Append(newSlice, reflect.Indirect(newStruct))
		if ln%2 != 0 {
			r.Reader.ReadByte()
			inLength--
		}
		inLength -= int(ln)
	}
	return newSlice, nil
}

// structChunks reads chunks and matches them to fields on rv (which is a struct)
// structChunks sets the fields of rv to be the output it gets
func (r *Reader) structChunks(rv reflect.Value, inLength int) error {
	chunkID, ln, isList, err := r.NextIDLen()
	if err != nil {
		return err
	}
	if isList {
		ln -= 4
		inLength -= 4
		chunkID, err = r.NextID()
		if err != nil {
			return err
		}
	}
	inLength -= 8
	ty := reflect.TypeOf(rv.Interface())
	fields := make([]reflect.Value, rv.NumField())
	fieldTags := make([]reflect.StructTag, rv.NumField())
	for i := range fields {
		fields[i] = rv.Field(i)
		fieldTags[i] = ty.Field(i).Tag
	}
	i := 0
	for inLength > 0 {
		tag := fieldTags[i].Get("riff")
		//spew.Dump(fields[i])
		if tag == chunkID {
			// get contents from recursive call
			var content reflect.Value
			if isList {
				content, err = r.chunks(fields[i], int(ln))
			} else {
				content, err = r.fieldValue(fields[i], ln)
			}
			if err != nil {
				return err
			}
			inLength -= int(ln)

			fields[i].Set(content)
			// if length is odd read one more
			if ln%2 != 0 {
				r.Reader.ReadByte()
				inLength--
			}
			if inLength <= 0 {
				return nil
			}
			// next id
			chunkID, ln, isList, err = r.NextIDLen()
			if err != nil {
				return err
			}
			if isList {
				ln -= 4
				inLength -= 4
				chunkID, err = r.NextID()
				if err != nil {
					return err
				}
			}
			inLength -= 8
			i = -1
		}
		if inLength <= 0 {
			return nil
		}
		i++
		if i >= len(fields) {
			// Skip this id
			// if length is odd read one more
			if ln%2 != 0 {
				ln++
			}
			_, err = r.Reader.Seek(int64(ln), io.SeekCurrent)
			if err != nil {
				return err
			}
			inLength -= int(ln)
			if inLength <= 0 {
				return nil
			}
			// next id
			chunkID, ln, isList, err = r.NextIDLen()
			if err != nil {
				return err
			}
			if isList {
				ln -= 4
				inLength -= 4
				chunkID, err = r.NextID()
				if err != nil {
					return err
				}
			}
			inLength -= 8
			i = 0
		}
	}
	return nil
}

// Todo: the switch here should change to some separate functions, there's some
// repetition here that is not necessary.
func (r *Reader) fieldValue(rv reflect.Value, ln uint32) (reflect.Value, error) {
	switch rv.Kind() {
	case reflect.Struct:
		st := rv.Addr().Interface()
		err := binary.Read(r.Reader, binary.LittleEndian, st)
		if err != nil {
			// Something on this struct has an undefined size
			// Read each field in part by part.
		}
		return reflect.Indirect(reflect.ValueOf(st)), err
	case reflect.String:
		data := make([]byte, ln)
		n, err := r.Reader.Read(data)
		if n != int(ln) {
			return reflect.Value{}, errors.New("Insufficient data found in RIFF data block")
		}
		return reflect.ValueOf(string(data)), err
	case reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.Uint8:
			data := make([]byte, ln)
			n, err := r.Reader.Read(data)
			if n != int(ln) {
				return reflect.Value{}, errors.New("Insufficient data found in RIFF data block")
			}
			return reflect.ValueOf(data), err
		default:
			return reflect.Value{}, errors.New("Unsupported type in input struct")
		}
	case reflect.Uint32:
		if ln != 4 {
			return reflect.Value{}, errors.New("Invalid length for uint32: " + strconv.Itoa(int(ln)))
		}
		data := make([]byte, ln)
		n, err := r.Reader.Read(data)
		if n != int(ln) {
			return reflect.Value{}, errors.New("Insufficient data found in RIFF data block")
		}
		if err != nil {
			return reflect.Value{}, err
		}
		val, n := binary.Uvarint(data)
		if n <= 0 {
			return reflect.Value{}, errors.New("Unable to decode int64 from data")
		}
		val32 := uint32(val)
		return reflect.ValueOf(val32), nil
	case reflect.Int64:
		if ln != 8 {
			return reflect.Value{}, errors.New("Invalid length for int64: " + strconv.Itoa(int(ln)))
		}
		data := make([]byte, ln)
		n, err := r.Reader.Read(data)
		if n != int(ln) {
			return reflect.Value{}, errors.New("Insufficient data found in RIFF data block")
		}
		if err != nil {
			return reflect.Value{}, err
		}
		val, n := binary.Varint(data)
		if n <= 0 {
			return reflect.Value{}, errors.New("Unable to decode int64 from data")
		}
		return reflect.ValueOf(val), nil
	}
	return reflect.Value{}, nil
}
