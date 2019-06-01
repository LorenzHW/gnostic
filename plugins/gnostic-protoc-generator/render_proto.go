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
	renderRequestParameters(f, renderer)
	renderResponses(f, renderer)
	renderComponents(f, renderer)

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
	f.WriteLine(``)
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

func renderRequestParameters(f *LineWriter, renderer *Renderer) {
	for _, t := range renderer.Model.Types {
		if strings.Contains(t.Name, "Parameters") {
			lines := make([]string, 0)
			protobufFieldNumberCounter := 1
			for _, field := range t.Fields {
				if field.Kind == surface.FieldKind_REFERENCE {
					// We got a reference to a parameter
					// According to: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62
					// it is not allowed to have non-primitive types. hence we flatten it:
					fieldType := getTypeForName(field.Type, renderer.Model.Types)
					for _, f := range fieldType.Fields {
						l := getFieldLine(f.Kind, f.NativeType, f.Name, protobufFieldNumberCounter)
						protobufFieldNumberCounter += 1
						lines = append(lines, l)
					}
				} else {
					l := getFieldLine(field.Kind, field.NativeType, field.Name, protobufFieldNumberCounter)
					protobufFieldNumberCounter += 1
					lines = append(lines, l)
				}
			}
			renderMessages(f, t.Name, lines)
		}
	}
	f.WriteLine(``)
}

func renderResponses(f *LineWriter, renderer *Renderer) {
	for _, t := range renderer.Model.Types {
		if strings.Contains(t.Name, "Responses") {
			lines := make([]string, 0)
			fields := removeDuplicates(t.Fields)
			for i, field := range fields {
				l := getFieldLine(field.Kind, field.NativeType, field.FieldName, i+1)
				lines = append(lines, l)
			}
			renderMessages(f, t.Name, lines)
		}
	}
	f.WriteLine(``)
}

func renderComponents(f *LineWriter, renderer *Renderer) {
	for _, t := range renderer.Model.Types {
		if !strings.Contains(t.Name, "Parameters") && !strings.Contains(t.Name, "Responses") {
			lines := make([]string, 0)
			for i, field := range t.Fields {
				l := getFieldLine(field.Kind, field.NativeType, field.Name, i+1)
				lines = append(lines, l)
			}
			renderMessages(f, t.Name, lines)
			f.WriteLine(``)
		}
	}
}

func renderMessages(f *LineWriter, messageName string, fields []string) {
	f.WriteLine(`message ` + messageName + ` {`)
	for _, field := range fields {
		f.WriteLine(field)
	}
	f.WriteLine(`}`)
}

func getFieldLine(kind surface.FieldKind, nativeType string, fieldName string, counter int) (l string) {
	if kind == surface.FieldKind_ARRAY {
		l = `  ` + `repeated ` + nativeType + ` ` + fieldName + ` = ` + strconv.Itoa(counter) + `;`
	} else {
		l = `  ` + nativeType + ` ` + fieldName + ` = ` + strconv.Itoa(counter) + `;`
	}
	return l
}

func getTypeForName(name string, types []*surface.Type) (t *surface.Type) {
	for _, t := range types {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func removeDuplicates(fields []*surface.Field) (result []surface.Field) {
	// It is possible that the OpenAPI description has responses with multiple
	// formats (e.g.: application/json and application/xml), if that is the
	// case we only want to create on field inside our protobuf file,
	// therefore we remove the duplicates
	alreadyAdded := make(map[string]bool)
	for _, f := range fields {
		if !alreadyAdded[f.Name] {
			result = append(result, *f)
			alreadyAdded[f.Name] = true
		}
	}
	return result
}

// TODO: Referencing a $ref parameter does not work: In surface model now --> render correctly (native fields)
// TODO: Print description / summary of field names
// TODO: Add requestbodies
// TODO: documentation of functions

// WATCH OUT FOR:
// The path template may refer to one or more fields in the gRPC request message, as long
// as each field is a non-repeated field with a primitive (non-message) type.
// see: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L62
// AND:
// Note that fields which are mapped to URL query parameters must have a
// primitive type or a repeated primitive type or a non-repeated message type.
// see: https://github.com/googleapis/googleapis/blob/a8ee1416f4c588f2ab92da72e7c1f588c784d3e6/google/api/http.proto#L119

//TODO: handle enum
//TODO: Watch out for: https://github.com/googleapis/googleapis/blob/152dabdfea620675c2db2f2a74878572324e8fd2/google/api/http.proto#L308
//TODO: Not sure if possible
