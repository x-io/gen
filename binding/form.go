// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"

	"github.com/x-io/gen/core"
)

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(32 << 10) // 32 MB
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return validate(obj)
}

func (formBinding) Write(response core.Response, obj interface{}) error {
	switch data := obj.(type) {
	case string:
		response.WriteString(data)
		break
	case []byte:
		response.Write(data)
		break
	case *bytes.Buffer:
		response.Write(data.Bytes())
		break
	default:
		response.WriteString("类型转换异常")
	}
	response.Header().Set("Content-Type", "application/text; charset=utf-8")
	return nil
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return validate(obj)
}

func (formPostBinding) Write(response core.Response, obj interface{}) error {
	switch data := obj.(type) {
	case string:
		response.WriteString(data)
		break
	case []byte:
		response.Write(data)
		break
	case *bytes.Buffer:
		response.Write(data.Bytes())
		break
	default:
		response.WriteString("类型转换异常")
	}
	response.Header().Set("Content-Type", "application/text; charset=utf-8")
	return nil
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(32 << 10); err != nil {
		return err
	}
	if err := mapForm(obj, req.MultipartForm.Value); err != nil {
		return err
	}
	return validate(obj)
}

func (formMultipartBinding) Write(response core.Response, obj interface{}) error {
	return nil
}
