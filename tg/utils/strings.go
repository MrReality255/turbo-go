package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12/context"
	"golang.org/x/exp/constraints"
)

type IStrings interface {
	Add(items ...string)
	Addf(item string, args ...any)
	Content() []string
	Sorted() []string
	Join(sep string) string
	SortJoin(sep string) string
}

type StringList struct {
	mx      sync.Mutex
	items   []string
	itemMap map[string]bool
}

func FromJSON(src interface{}, ptr interface{}) error {
	switch src := src.(type) {
	case *context.Context:
		b, err := src.GetBody()
		if err != nil {
			return nil
		}
		return FromJSON(b, ptr)
	case io.Reader:
		return json.NewDecoder(src).Decode(ptr)
	case []byte:
		return json.Unmarshal(src, ptr)
	case string:
		return json.Unmarshal([]byte(src), ptr)
	default:
		panic(fmt.Errorf("unable to convert %T to json", src))
	}
}

func ToJSONB(obj interface{}) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func ToJSON(obj interface{}) string {
	return string(ToJSONB(obj))
}

func SaveToJSON(file string, obj interface{}, readable bool) error {
	var (
		b   []byte
		err error
	)

	if readable {
		b, err = json.MarshalIndent(obj, "", "  ")
	} else {
		b, err = json.Marshal(obj)
	}

	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0666)
}

func GetGUID() string {
	return uuid.New().String()
}

func GetToken() string {
	return strings.Replace(GetGUID(), "-", "", -1)
}

func Hash(input string, count int) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(input))
	return int(hasher.Sum32() % uint32(count))
}

func EvaluateAsStr(x interface{}) string {
	switch x := x.(type) {
	case nil:
		return ""
	case func() string:
		return x()
	case string:
		return x
	default:
		return fmt.Sprintf("%v", x)
	}
}

func ParseJSON[T any](src any) (*T, error) {
	var t T
	err := FromJSON(src, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func StrToIntDef(str string, defaultValue int) int {
	p, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(p)
}

func StrToInt64Def(str string, defaultValue int64) int64 {
	p, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultValue
	}
	return p
}

func StrToNrDef[A constraints.Integer](str string, def A) A {
	return A(StrToInt64Def(str, int64(def)))
}

func StrToWordDef(str string, def uint16) uint16 {
	return StrToNrDef(str, def)
}

func NewStringList(maxCount int, isDistinct bool, initItems ...string) IStrings {
	list := &StringList{
		items: make([]string, 0, maxCount),
		itemMap: IfThenFct(
			isDistinct,
			func() map[string]bool {
				return make(map[string]bool, maxCount)
			},
			ConstFct[map[string]bool](nil),
		),
	}

	list.Add(initItems...)
	return list
}

func (sl *StringList) Add(items ...string) {
	sl.mx.Lock()
	defer sl.mx.Unlock()

	for _, item := range items {
		canAdd := sl.itemMap == nil || !sl.itemMap[item]
		if canAdd {
			sl.items = append(sl.items, item)
			if sl.itemMap != nil {
				sl.itemMap[item] = true
			}
		}
	}
}

func (sl *StringList) Addf(item string, args ...any) {
	sl.Add(fmt.Sprintf(item, args...))
}

func (sl *StringList) Content() []string {
	sl.mx.Lock()
	defer sl.mx.Unlock()
	return sl.items
}

func (sl *StringList) Sorted() []string {
	return ArraySort(sl.Content(), func(item1 string, item2 string) bool {
		return item1 < item2
	})
}

func (sl *StringList) Join(sep string) string {
	return strings.Join(sl.Content(), sep)
}

func (sl *StringList) SortJoin(sep string) string {
	return strings.Join(sl.Sorted(), sep)
}
