package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ByteArray(items ...interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := WriteBytes(b, items...); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
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
