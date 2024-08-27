package proto

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type Response struct {
	statusCode int
	header     http.Header
	resByte    []byte
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) ResByte() []byte {
	return r.resByte
}

func (r *Response) JsonRawMessage() json.RawMessage {
	res := json.RawMessage{}
	json.Unmarshal(r.resByte, &res)
	return res
}

func (p *Response) MarshalJSON() ([]byte, error) {
	type TmpResponse struct {
		StatusCode int             `json:"status_code"`
		Header     http.Header     `json:"header"`
		Res        json.RawMessage `json:"res"`
	}

	var tmp TmpResponse
	tmp.Res = p.JsonRawMessage()
	tmp.StatusCode = p.statusCode
	tmp.Header = p.header

	return json.Marshal(tmp)
}

func (u *Response) UnmarshalJSON(b []byte) error {
	type TmpResponse struct {
		StatusCode int             `json:"status_code"`
		Header     http.Header     `json:"header"`
		Res        json.RawMessage `json:"res"`
	}

	var tmp TmpResponse
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	u.statusCode = tmp.StatusCode
	u.header = tmp.Header

	data, err := json.Marshal(tmp.Res)
	if err != nil {
		return err
	}
	u.resByte = data

	return nil
}

func NewResponse(statusCode int, header http.Header, readCloser io.ReadCloser) (res *Response, err error) {
	res = new(Response)
	defer readCloser.Close()
	var resByte []byte
	resByte, err = io.ReadAll(readCloser)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res.resByte = resByte
	res.header = header
	res.statusCode = statusCode
	return
}

func (r *Response) StatusErr() (err error) {
	if gjson.GetBytes(r.resByte, "status").Int() > 201 {
		err = errors.New(string(r.resByte))
		return
	}
	return
}
