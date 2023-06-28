package dto

import (
	"sync"

	"github.com/yonisaka/go-boilerplate/pkg/msg"
)

var rsp *HttpResponse
var oneRsp sync.Once

// HttpResponse presentation contract object
type HttpResponse struct {
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
func (r *HttpResponse) GetCode() int {
	return msg.GetCode(r.Code)
}

// GetMessage method to transform response name var to message detail
func (r *HttpResponse) GetMessage() string {
	return msg.Get(r.Code, r.Lang)
}

// GenerateMessage setter message
func (r *HttpResponse) GenerateMessage() {
	if r.Message == nil {
		r.Message = msg.Get(r.Code, r.Lang)
	}
}

// WithCode setter response var name
func (r *HttpResponse) WithCode(c int) *HttpResponse {
	r.Code = c
	return r
}

// WithData setter data response
func (r *HttpResponse) WithData(v interface{}) *HttpResponse {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *HttpResponse) WithError(v interface{}) *HttpResponse {
	r.Errors = v
	return r
}

// WithMeta setter meta data response
func (r *HttpResponse) WithMeta(v interface{}) *HttpResponse {
	r.Meta = v
	return r
}

// WithMessage setter custom message response
func (r *HttpResponse) WithMessage(v interface{}) *HttpResponse {
	if v != nil {
		r.Message = v
	}
	return r
}

func (r *HttpResponse) Error() *HttpResponse {
	r.Code = 500
	return r
}

func (r *HttpResponse) NotFound() *HttpResponse {
	r.Code = 404
	return r
}

func (r *HttpResponse) Success() *HttpResponse {
	r.Code = 200
	return r
}

// WithIsStream setter custom hash response
func (r *HttpResponse) WithIsStream(v bool) *HttpResponse {
	r.isStream = v
	return r
}

func (r *HttpResponse) WithStreamData(v []byte) *HttpResponse {
	r.dataStream = v
	return r
}

func (r *HttpResponse) WithContentType(v string) *HttpResponse {
	r.streamContentType = v
	return r
}

func (r *HttpResponse) ContentType() string {
	return r.streamContentType
}

func (r *HttpResponse) IsStream() bool {
	return r.isStream
}

func (r *HttpResponse) DataStream() []byte {
	return r.dataStream
}

// NewResponse initialize response
func NewResponse() *HttpResponse {
	oneRsp.Do(func() {
		rsp = &HttpResponse{}
	})

	// clone response
	x := *rsp

	return &x
}
