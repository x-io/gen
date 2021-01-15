// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/x-io/gen/core"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func (xmlBinding) Write(response *core.Response, obj interface{}) error {
	response.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if err := xml.NewEncoder(response).Encode(obj); err != nil {
		response.Header().Del("Content-Type")
		return err
	}
	return nil
}
