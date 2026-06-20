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
	"time"

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

func DateTimeToStr(x time.Time) string {
	return x.Format("2006-01-02 15:04:05")
}

func StrToIntDef(str string, defaultValue int) int {
	p, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(p)
}

func StrToTime(str string) time.Time {
	return StrToTimeLoc(str, time.Local)
}

func StrToTimeUTC(str string) time.Time {
	return StrToTimeLoc(str, time.UTC)
}

func StrToTimeLoc(str string, loc *time.Location) time.Time {
	formats := []string{
		"2006-01-02 15:04:05",
		"02.01.2006 15:04:05",
		"2006-01-02",
		"02.01.2006",
		"15:04:05",
	}

	for _, f := range formats {
		v, err := time.ParseInLocation(f, str, loc)
		if err == nil {
			return v
		}
	}
	return time.Time{}
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

func CreateTable[T any](
	rows []T,
	cols []string,
	rowFct func(row T) []any,
) string {
	if len(cols) == 0 {
		return ""
	}

	// prepare data -> convert everything to strings
	strRows := ArrayMap(
		rows,
		func(item T) []string {
			var (
				row    = rowFct(item)
				strRow = make([]string, len(cols))
			)
			for i, v := range row {
				strRow[i] = fmt.Sprintf("%v", v)
			}
			return strRow
		},
	)

	colWidths := ArrayMapIdx(
		cols,
		func(item string, idx int) int {
			return Max(
				append([]int{len(item)},
					ArrayMap(
						strRows,
						func(item []string) int {
							return len(item[idx])
						},
					)...,
				)...,
			)
		},
	)

	separator := strings.Join(
		ArrayMap(
			colWidths,
			func(w int) string {
				return strings.Repeat("-", w+2)
			},
		), "+")

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("+%v+\n", separator))
	sb.WriteString(
		fmt.Sprintf("|%v|",
			strings.Join(ArrayMapIdx(cols, func(item string, idx int) string {
				padding := colWidths[idx] - len(item)
				return fmt.Sprintf(" %v%v ", item, strings.Repeat(" ", padding))
			}), "+")),
	)
	sb.WriteString(fmt.Sprintf("\n+%v+\n", separator))
	for _, row := range strRows {
		sb.WriteString(
			fmt.Sprintf("|%v|\n",
				strings.Join(ArrayMapIdx(row, func(item string, idx int) string {
					padding := colWidths[idx] - len(item)
					return fmt.Sprintf(" %v%v ", item, strings.Repeat(" ", padding))
				}), "+")),
		)
	}

	return sb.String()
}
