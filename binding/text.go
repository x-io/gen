// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"net/http"

	"github.com/x-io/gen/core"
)

type textBinding struct{}

func (textBinding) Name() string {
	return "text"
}

func (textBinding) Bind(req *http.Request, obj interface{}) error {

	return validate(obj)
}

func (textBinding) Write(response *core.Response, obj interface{}) error {
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
