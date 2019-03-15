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

	f.WriteLine("// GENERATED FILE: DO NOT EDIT!")
	f.WriteLine(``)
	f.WriteLine(`syntax = "proto3";`)
	f.WriteLine(`import "google/api/annotations.proto";`)
	f.WriteLine(`package ` + renderer.Package + `;`)
	f.WriteLine(``)

	defineRPCservice(f, renderer)

	return f.Bytes(), nil
}

func renderRPCsignature(f *LineWriter, methodName string) {
	requestName := methodName + "Request"
	responseName := methodName + "Response"
	f.WriteLine(`  rpc ` + methodName + ` (` + requestName + `) ` + `returns` + ` (` + responseName + `) {`)
}

func renderOptions(f *LineWriter) {
	f.WriteLine(`    option (google.api.http) = {`)
	f.WriteLine(`      get: "/path"`)
	f.WriteLine(`    };`)
}

func defineRPCservice(f *LineWriter, renderer *Renderer) {

	f.WriteLine(`service ` + strings.Title(renderer.Package) + `{`)
	for _, method := range renderer.Model.Methods {
		renderRPCsignature(f, method.Name)
		renderOptions(f)
		f.WriteLine(`  }`) // Closing bracket of method
		f.WriteLine(``)
	}
	f.WriteLine(`}`)
}
