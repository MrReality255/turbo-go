package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func ByteArray(items ...interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := WriteBytes(b, items...); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func CloseAfterWith[T any](x io.Closer, fct func() (T, error)) (result T, err error) {
	err = CloseAfter(x, func() error {
		d, err2 := fct()
		if err2 == nil {
			result = d
		}
		return err2
	})
	return
}

func CloseAfter(x io.Closer, fct func() error) error {
	var isClosed bool
	defer func() {
		if !isClosed {
			_ = x.Close()
		}
	}()

	if err := fct(); err != nil {
		return err
	}
	isClosed = true
	return x.Close()
}

func Dir(p string) ([]string, error) {
	return DirEx(p, true, true)
}

func DirEx(p string, wantDir bool, wantFiles bool) ([]string, error) {
	list, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	result := ArrayMapEx(
		list,
		func(item os.DirEntry) (string, bool) {
			return filepath.Join(p, item.Name()), (wantDir && item.IsDir()) || (wantFiles && !item.IsDir())
		},
	)
	return result, nil

}

func DirExists(dir string) bool {
	stat, err := os.Stat(dir)
	return err == nil && stat.IsDir()
}

func FileExists(file string) bool {
	stat, err := os.Stat(file)
	return err == nil && !stat.IsDir()
}

func MustByteArray(items ...interface{}) []byte {
	a, err := ByteArray(items...)
	if err != nil {
		panic(err)
	}
	return a
}

func FromBytes(b []byte, targets ...interface{}) error {
	return FromReader(bytes.NewReader(b), targets...)
}

func FromReader(r io.Reader, targets ...interface{}) error {
	for _, addr := range targets {
		if err := binary.Read(r, binary.BigEndian, addr); err != nil {
			return err
		}
	}
	return nil
}

func GetFileNameWithoutExt(filename string) string {
	var (
		n   = filepath.Base(filename)
		ext = filepath.Ext(filename)
	)

	return strings.TrimSuffix(n, ext)
}

func GetFileNameWithSuffix(path string, suffix string) string {
	var (
		ext  = filepath.Ext(path)
		dir  = filepath.Dir(path)
		base = filepath.Base(path)
		name = base[:len(base)-len(ext)]
	)

	return filepath.Join(dir, fmt.Sprintf("%v%v%v", name, suffix, ext))
}

func LoadDirJSON[T any](dir string, prepFct func() *T) ([]*T, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	r, err := ArrayMapExErr(
		files,
		func(item os.DirEntry) (*T, bool, error) {
			newValue := prepFct()
			err = LoadFromJSON(filepath.Join(dir, item.Name()), newValue)
			return newValue, true, err
		},
	)
	return r, err
}

func SaveToCSV[T any](filename string, data []T, header []string, rowCb func(row T) []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	return CloseAfter(file, func() error {
		writer := csv.NewWriter(file)
		defer writer.Flush()

		if len(header) > 0 {
			if err := writer.Write(header); err != nil {
				return err
			}
		}

		for _, item := range data {
			row := rowCb(item)
			if err := writer.Write(row); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	})
}

func LoadFromCSV[T any](filename string, rowCallback func([]string) *T, skipRows int) ([]*T, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	result := make([]*T, 0, 64000)

	if err := CloseAfter(file, func() error {
		var (
			rdr    = csv.NewReader(file)
			rowCnt = 0
		)
		for {
			row, err := rdr.Read()
			if IsErr(err, io.EOF) {
				return nil
			}

			if rowCnt < skipRows {
				rowCnt++
				continue
			}

			if r := rowCallback(row); r != nil {
				result = append(result, r)
			} else {
				return ErrInvalidFileFormat
			}
		}
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func LoadFromBin(filename string, ptr any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return CloseAfter(f, func() error {
		g := gob.NewDecoder(f)
		return g.Decode(ptr)
	})
}

func LoadFromJSON(file string, target interface{}) error {
	c, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return FromJSON(c, target)
}

func MkDir(targetDir string) error {
	return os.MkdirAll(targetDir, os.ModePerm)
}

func MkDirFor(targetFile string) error {
	return MkDir(filepath.Dir(targetFile))
}

func NewFromFileJSON[T any](file string) (*T, error) {
	var item T
	err := LoadFromJSON(file, &item)
	return &item, err
}

func SaveToBin(filename string, obj any) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	return CloseAfter(f, func() error {
		g := gob.NewEncoder(f)
		return g.Encode(obj)
	})
}

func WriteBytes(w io.Writer, items ...interface{}) error {
	for _, item := range items {
		switch item := item.(type) {
		case []byte:
			if _, err := w.Write(item); err != nil {
				return err
			}
		default:
			if err := binary.Write(w, binary.BigEndian, item); err != nil {
				return err
			}
		}
	}
	return nil
}
