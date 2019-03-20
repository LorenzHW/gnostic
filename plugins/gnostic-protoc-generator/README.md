# Protoc Generator Plugin

This directory contains a `gnostic` plugin that can be used to generate a protocol buffer specification for an API with an OpenAPI description. This protocol buffer specification can further be used to generate a gRPC API.

Run inside this directory:

    go build

Then run to generate a petstore.proto file.

    ./gnostic-protoc-generator -input development-test-data/petstore.pb -output development-test-data/generated-by-plugin/protoc/petstore