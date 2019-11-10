package main

import (
	"path/filepath"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestSpecificationSchema(t *testing.T) {
	schemaLoader := gojsonschema.NewReferenceLoader("http://json-schema.org/draft-07/schema")
	documentLoader := gojsonschema.NewReferenceLoader("file://cfn-resource-specification-schema.json")

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid() {
		t.Error("The document is not valid. see errors :")
		for _, desc := range result.Errors() {
			t.Error(desc)
		}
	}
}

func TestSpecification(t *testing.T) {
	schema, err := filepath.Abs("cfn-resource-specification-schema.json")
	if err != nil {
		t.Fatal(err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schema)
	documentLoader := gojsonschema.NewReferenceLoader("file://cfn-resource-specification.json")

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid() {
		t.Error("The document is not valid. see errors :")
		for _, desc := range result.Errors() {
			t.Error(desc)
		}
	}
}
