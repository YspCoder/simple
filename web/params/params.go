package params

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/YspCoder/simple/common/dates"
	"github.com/YspCoder/simple/common/strs"
	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	decoder = schema.NewDecoder() // form, url, schema.
)

func init() {
	decoder.AddAliasTag("form", "json")
	decoder.ZeroEmpty(true)
	decoder.RegisterConverter(primitive.ObjectID{}, func(value string) reflect.Value {
		objID, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			return reflect.ValueOf(primitive.ObjectID{})
		}
		return reflect.ValueOf(objID)
	})
}

// param error
func paramError(name string) error {
	return errors.New(fmt.Sprintf("unable to find param value '%s'", name))
}

// ReadForm read object from FormData
func ReadForm(ctx iris.Context, obj interface{}) error {
	values := ctx.FormValues()
	if len(values) == 0 {
		return nil
	}
	return decoder.Decode(obj, values)
}

func Get(ctx iris.Context, name string) (string, bool) {
	str := ctx.FormValue(name)
	return str, str != ""
}

func GetInt64(c iris.Context, name string) (int64, bool) {
	str, ok := Get(c, name)
	if !ok {
		return 0, false
	}
	value, err := cast.ToInt64E(str)
	if err != nil {
		return 0, false
	}
	return value, true
}

func GetInt(c iris.Context, name string) (int, bool) {
	str, ok := Get(c, name)
	if !ok {
		return 0, false
	}
	value, err := cast.ToIntE(str)
	if err != nil {
		return 0, false
	}
	return value, true
}

func GetBool(c iris.Context, name string) (bool, bool) {
	str, ok := Get(c, name)
	if !ok {
		return false, false
	}
	value, err := cast.ToBoolE(str)
	if err != nil {
		return false, false
	}
	return value, true
}

func GetFloat32(c iris.Context, name string) (float32, bool) {
	str, ok := Get(c, name)
	if !ok {
		return 0, false
	}
	value, err := cast.ToFloat32E(str)
	if err != nil {
		return 0, false
	}
	return value, true
}

func GetFloat64(c iris.Context, name string) (float64, bool) {
	str, ok := Get(c, name)
	if !ok {
		return 0, false
	}
	value, err := cast.ToFloat64E(str)
	if err != nil {
		return 0, false
	}
	return value, true
}

func GetTime(ctx iris.Context, name string) *time.Time {
	value, _ := Get(ctx, name)
	if strs.IsBlank(value) {
		return nil
	}
	layouts := []string{dates.FmtDateTime, dates.FmtDate, dates.FmtDateTimeNoSeconds}
	for _, layout := range layouts {
		if ret, err := dates.Parse(value, layout); err == nil {
			return &ret
		}
	}
	return nil
}

func GetInt64Arr(c iris.Context, name string) []int64 {
	str, ok := Get(c, name)
	if ok {
		ss := strings.Split(str, ",")
		var ret []int64
		for _, s := range ss {
			i, err := cast.ToInt64E(s)
			if err == nil {
				ret = append(ret, i)
			}
		}
		return ret
	}
	return nil
}

func FormValue(ctx iris.Context, name string) string {
	return ctx.FormValue(name)
}

func FormValueRequired(ctx iris.Context, name string) (string, error) {
	str := FormValue(ctx, name)
	if len(str) == 0 {
		return "", errors.New("参数：" + name + "不能为空")
	}
	return str, nil
}

func FormValueDefault(ctx iris.Context, name, def string) string {
	return ctx.FormValueDefault(name, def)
}

func FormValueInt(ctx iris.Context, name string) (int, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.Atoi(str)
}

func FormValueIntDefault(ctx iris.Context, name string, def int) int {
	if v, err := FormValueInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64(ctx iris.Context, name string) (int64, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return 0, paramError(name)
	}
	return strconv.ParseInt(str, 10, 64)
}

func FormValueInt64Default(ctx iris.Context, name string, def int64) int64 {
	if v, err := FormValueInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func FormValueInt64Array(ctx iris.Context, name string) []int64 {
	str := ctx.FormValue(name)
	if str == "" {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []int64
	for _, v := range ss {
		item, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, item)
	}
	return ret
}

func FormValueStringArray(ctx iris.Context, name string) []string {
	str := ctx.FormValue(name)
	if len(str) == 0 {
		return nil
	}
	ss := strings.Split(str, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func FormValueBool(ctx iris.Context, name string) (bool, error) {
	str := ctx.FormValue(name)
	if str == "" {
		return false, paramError(name)
	}
	return strconv.ParseBool(str)
}

func FormValueBoolDefault(ctx iris.Context, name string, def bool) bool {
	str := ctx.FormValue(name)
	if str == "" {
		return def
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return def
	}
	return value
}

// 从请求中获取日期
func FormDate(ctx iris.Context, name string) *time.Time {
	value := FormValue(ctx, name)
	if strs.IsBlank(value) {
		return nil
	}
	layouts := []string{dates.FmtDateTime, dates.FmtDate, dates.FmtDateTimeNoSeconds}
	for _, layout := range layouts {
		if ret, err := dates.Parse(value, layout); err == nil {
			return &ret
		}
	}
	return nil
}

func GetPaging(ctx iris.Context) (page int64, limit int64) {
	page = FormValueInt64Default(ctx, "page", 1)
	limit = FormValueInt64Default(ctx, "limit", 20)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return
}
