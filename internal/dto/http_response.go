package dto

import (
	"sync"

	"github.com/yonisaka/go-boilerplate/pkg/msg"
)

var rsp *HTTPResponse
var oneRsp sync.Once

// HTTPResponse presentation contract object
type HTTPResponse struct {
	Code              int         `json:"code"`
	Message           interface{} `json:"message,omitempty"`
	Errors            interface{} `json:"errors,omitempty"`
	Data              interface{} `json:"data,omitempty"`
	Lang              string      `json:"-"`
	Meta              interface{} `json:"meta,omitempty"`
	isStream          bool
	streamContentType string
	dataStream        []byte
}

// MetaData represent meta data response for multi data
type MetaData struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"per_page"`
}

// MetaDataWithCount represent meta data response for multi data
type MetaDataWithCount struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"per_page"`
	Total uint64 `json:"total"`
}

// GetCode method to transform response name var to http status
func (r *HTTPResponse) GetCode() int {
	return msg.GetCode(r.Code)
}

// GetMessage method to transform response name var to message detail
func (r *HTTPResponse) GetMessage() string {
	return msg.Get(r.Code, r.Lang)
}

// GenerateMessage setter message
func (r *HTTPResponse) GenerateMessage() {
	if r.Message == nil {
		r.Message = msg.Get(r.Code, r.Lang)
	}
}

// WithCode setter response var name
func (r *HTTPResponse) WithCode(c int) *HTTPResponse {
	r.Code = c
	return r
}

// WithData setter data response
func (r *HTTPResponse) WithData(v interface{}) *HTTPResponse {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *HTTPResponse) WithError(v interface{}) *HTTPResponse {
	r.Errors = v
	return r
}

// WithMeta setter meta data response
func (r *HTTPResponse) WithMeta(v interface{}) *HTTPResponse {
	r.Meta = v
	return r
}

// WithMessage setter custom message response
func (r *HTTPResponse) WithMessage(v interface{}) *HTTPResponse {
	if v != nil {
		r.Message = v
	}

	return r
}

func (r *HTTPResponse) Error() *HTTPResponse {
	r.Code = 500
	return r
}

func (r *HTTPResponse) NotFound() *HTTPResponse {
	r.Code = 404
	return r
}

func (r *HTTPResponse) Success() *HTTPResponse {
	r.Code = 200
	return r
}

// WithIsStream setter custom hash response
func (r *HTTPResponse) WithIsStream(v bool) *HTTPResponse {
	r.isStream = v
	return r
}

func (r *HTTPResponse) WithStreamData(v []byte) *HTTPResponse {
	r.dataStream = v
	return r
}

func (r *HTTPResponse) WithContentType(v string) *HTTPResponse {
	r.streamContentType = v
	return r
}

func (r *HTTPResponse) ContentType() string {
	return r.streamContentType
}

func (r *HTTPResponse) IsStream() bool {
	return r.isStream
}

func (r *HTTPResponse) DataStream() []byte {
	return r.dataStream
}

// NewResponse initialize response
func NewResponse() *HTTPResponse {
	oneRsp.Do(func() {
		rsp = &HTTPResponse{}
	})

	// clone response
	x := *rsp

	return &x
}
