package web

import (
	"errors"

	"github.com/YspCoder/simple/common/structs"
)

type JsonResult struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func Json(code int, message string, data interface{}, success bool) *JsonResult {
	return &JsonResult{
		Code:    code,
		Msg:     message,
		Data:    data,
		Success: success,
	}
}

func JsonData(data interface{}) *JsonResult {
	return &JsonResult{
		Code:    0,
		Data:    data,
		Success: true,
	}
}

func JsonItemList(data []interface{}) *JsonResult {
	return &JsonResult{
		Code:    0,
		Data:    data,
		Success: true,
	}
}

func JsonPageData(results interface{}, total int64) *JsonResult {
	return JsonData(&PageResult{
		Results: results,
		Total:   total,
	})
}

func JsonCursorData(results interface{}, cursor string, hasMore bool) *JsonResult {
	return JsonData(&CursorResult{
		Results: results,
		Cursor:  cursor,
		HasMore: hasMore,
	})
}

func JsonSuccess() *JsonResult {
	return &JsonResult{
		Code:    0,
		Data:    nil,
		Success: true,
	}
}

func JsonError(err error) *JsonResult {
	var e *CodeError
	if errors.As(err, &e) {
		return &JsonResult{
			Code:    e.Code,
			Msg:     e.Msg,
			Data:    e.Data,
			Success: false,
		}
	}
	return &JsonResult{
		Code:    0,
		Msg:     err.Error(),
		Data:    nil,
		Success: false,
	}
}

func JsonErrorMsg(message string) *JsonResult {
	return &JsonResult{
		Code:    0,
		Msg:     message,
		Data:    nil,
		Success: false,
	}
}
func JsonErrorCode(code int, message string) *JsonResult {
	return &JsonResult{
		Code:    code,
		Msg:     message,
		Data:    nil,
		Success: false,
	}
}

func JsonErrorData(code int, message string, data interface{}) *JsonResult {
	return &JsonResult{
		Code:    code,
		Msg:     message,
		Data:    data,
		Success: false,
	}
}

type RspBuilder struct {
	Data map[string]interface{}
}

func NewEmptyRspBuilder() *RspBuilder {
	return &RspBuilder{Data: make(map[string]interface{})}
}

func NewRspBuilder(obj interface{}) *RspBuilder {
	return NewRspBuilderExcludes(obj)
}

func NewRspBuilderExcludes(obj interface{}, excludes ...string) *RspBuilder {
	return &RspBuilder{Data: structs.StructToMap(obj, excludes...)}
}

func (builder *RspBuilder) Put(key string, value interface{}) *RspBuilder {
	builder.Data[key] = value
	return builder
}

func (builder *RspBuilder) Build() map[string]interface{} {
	return builder.Data
}

func (builder *RspBuilder) JsonResult() *JsonResult {
	return JsonData(builder.Data)
}

func ConvertList[T any](results []T, conv func(item T) map[string]interface{}) (list []map[string]interface{}) {
	for _, item := range results {
		if ret := conv(item); ret != nil {
			list = append(list, ret)
		}
	}
	return
}
