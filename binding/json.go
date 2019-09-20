// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"

	"net/http"

	"github.com/x-io/gen/core"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func (jsonBinding) Write(response core.Response, obj interface{}) error {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(response).Encode(obj); err != nil {
		response.Header().Del("Content-Type")
		return err
	}
	return nil
}
