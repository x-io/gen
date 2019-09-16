// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
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

func (textBinding) Write(response core.Response, obj interface{}) error {
	response.WriteHeader(200)
	response.WriteString(obj.(string))
	response.Header().Set("Content-Type", "application/text; charset=utf-8")
	return nil
}
