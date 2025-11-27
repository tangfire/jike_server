package middleware

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/proto"
)

// 统一响应格式
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// 成功响应
func success(data interface{}) *Response {
	return &Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// 转换错误为统一格式
func convertError(err error) *Response {
	// 如果是 Kratos 错误
	if kratosErr := errors.FromError(err); kratosErr != nil {
		return &Response{
			Code:      int(kratosErr.Code),
			Message:   kratosErr.Message,
			Data:      nil,
			Timestamp: time.Now().Unix(),
		}
	}

	// 其他错误
	return &Response{
		Code:      500,
		Message:   err.Error(),
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}

// 自定义 HTTP 响应编码器
func ResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 如果已经是 Response 类型，直接编码
	if resp, ok := v.(*Response); ok {
		return json.NewEncoder(w).Encode(resp)
	}

	// 检查是否为 protobuf 消息
	if pm, ok := v.(proto.Message); ok {
		//使用 protojson 序列化 protobuf 消息，保留零值字段
		marshalOptions := protojson.MarshalOptions{
			EmitUnpopulated: true, // 保留零值字段
			UseProtoNames:   true, // 使用 proto 字段名
		}

		// 先序列化 protobuf 数据
		protoData, err := marshalOptions.Marshal(pm)
		if err != nil {
			return err
		}

		// 将序列化后的 JSON 字节转换为 interface{}
		var data interface{}
		if err := json.Unmarshal(protoData, &data); err != nil {
			return err
		}

		// 创建成功响应
		resp := success(data)
		//resp := success(pm)
		return json.NewEncoder(w).Encode(resp)
	}

	// 其他类型包装为成功响应
	resp := success(v)
	return json.NewEncoder(w).Encode(resp)
}

// 自定义错误编码器
func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp := convertError(err)

	// 设置 HTTP 状态码
	var statusCode int
	switch resp.Code {
	case 0:
		statusCode = http.StatusOK
	case 400:
		statusCode = http.StatusBadRequest
	case 401:
		statusCode = http.StatusUnauthorized
	case 403:
		statusCode = http.StatusForbidden
	case 404:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
