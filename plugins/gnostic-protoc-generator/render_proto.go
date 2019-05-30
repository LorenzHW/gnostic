// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strconv"
	"strings"

	surface "github.com/googleapis/gnostic/surface"
)

// ParameterList returns a string representation of a method's parameters
func ParameterList(parametersType *surface.Type) string {
	result := ""
	if parametersType != nil {
		for _, field := range parametersType.Fields {
			result += field.ParameterName + " " + field.NativeType + "," + "\n"
		}
	}
	return result
}

func (renderer *Renderer) RenderProto() ([]byte, error) {
	f := NewLineWriter()

	// TODO: print license
	f.WriteLine("// GENERATED FILE: DO NOT EDIT!")
	f.WriteLine(``)
	f.WriteLine(`syntax = "proto3";`)
	f.WriteLine(``)
	f.WriteLine(`import "google/api/annotations.proto";`)
	f.WriteLine(`import "google/protobuf/empty.proto";`)
	f.WriteLine(``)
	f.WriteLine(`package ` + renderer.Package + `;`)
	f.WriteLine(``)

	renderRPCservice(f, renderer)
	f.WriteLine(``)
	renderMessages(f, renderer)

	return f.Bytes(), nil
}

func renderRPCservice(f *LineWriter, renderer *Renderer) {
	f.WriteLine(`service ` + strings.Title(renderer.Package) + ` {`)
	for _, method := range renderer.Model.Methods {
		renderRPCsignature(f, method)
		renderOptions(f, method)
		f.WriteLine(`  }`) // Closing bracket of method
		f.WriteLine(``)
	}
	f.WriteLine(`}`) // Closing bracket of RPC service
}

func renderRPCsignature(f *LineWriter, method *surface.Method) {
	if method.ResponsesTypeName == "" {
		method.ResponsesTypeName = "google.protobuf.Empty"
	}

	if method.ParametersTypeName == "" {
		method.ParametersTypeName = "google.protobuf.Empty"
	}

	f.WriteLine(`  rpc ` + method.Name + ` (` + method.ParametersTypeName + `) ` + `returns` + ` (` + method.ResponsesTypeName + `) {`)
}

func renderOptions(f *LineWriter, method *surface.Method) {
	f.WriteLine(`    option (google.api.http) = {`)
	f.WriteLine(`      ` + strings.ToLower(method.Method) + `: "` + method.Path + `"`)
	f.WriteLine(`    };`)
}

func renderMessages(f *LineWriter, renderer *Renderer) {
	for _, t := range renderer.Model.Types {
		f.WriteLine(`message ` + t.Name + ` {`)
		for i, field := range t.Fields {
			messageFieldName := field.Name
			//TODO: handle enum

			// Field is a HTTP 200 response
			if messageFieldName == "200" {
				// TODO: Better name for field. If it is a non primitive type (e.g.: pet) name field like that?
				// TODO: If there is also application/xml --> this will be rendered twice inside .proto
				messageFieldName = "response"
			}

			if field.Kind == surface.FieldKind_ARRAY {
				f.WriteLine(`  ` + `repeated ` + field.NativeType + ` ` + messageFieldName + ` = ` + strconv.Itoa(i+1) + `;`)
			} else {
				f.WriteLine(`  ` + field.NativeType + ` ` + messageFieldName + ` = ` + strconv.Itoa(i+1) + `;`)
			}
		}
		f.WriteLine(`}`)
		f.WriteLine(``)
	}

}
