package interpreter

import (
	"testing"
)

func TestJSONBasicFunctionality(t *testing.T) {
	// Test JSON namespace builtin exists
	jsonNamespaceFunc := builtins["JSON"]
	if jsonNamespaceFunc == nil {
		t.Fatal("JSON builtin function not found")
	}

	// Get the JSON namespace
	jsonNamespace := jsonNamespaceFunc.Fn()
	if jsonNamespace == nil {
		t.Fatal("JSON namespace returned nil")
	}
	
	_, ok := jsonNamespace.(*JSONNamespace)
	if !ok {
		t.Fatalf("expected JSONNamespace, got %T", jsonNamespace)
	}

	// Test parsing a simple JSON string via static method
	result := parseJSON(`{"name": "John", "age": 30}`)
	if result == nil {
		t.Fatal("parseJSON returned nil")
	}
	
	jsonObj, ok := result.(*JSON)
	if !ok {
		t.Fatalf("expected JSON object, got %T", result)
	}
	
	// Verify the parsed data is a hash
	hash, ok := jsonObj.Data.(*Hash)
	if !ok {
		t.Fatalf("expected Hash data, got %T", jsonObj.Data)
	}
	
	// Test accessing parsed data
	nameKey := CreateHashKey(&String{Value: "name"})
	nameValue, exists := hash.Pairs[nameKey]
	if !exists {
		t.Fatal("name key not found in parsed JSON")
	}
	
	nameStr, ok := nameValue.(*String)
	if !ok {
		t.Fatalf("expected String value for name, got %T", nameValue)
	}
	
	if nameStr.Value != "John" {
		t.Errorf("expected name 'John', got '%s'", nameStr.Value)
	}
}

func TestJSONStringifyBasic(t *testing.T) {
	// Test stringifying a string
	result, err := stringifyValue(&String{Value: "hello"})
	if err != nil {
		t.Fatalf("stringifyValue returned error: %s", err.Error())
	}
	
	expected := `"hello"`
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
	
	// Test stringifying an integer
	result, err = stringifyValue(&Integer{Value: 42})
	if err != nil {
		t.Fatalf("stringifyValue returned error: %s", err.Error())
	}
	
	if result != "42" {
		t.Errorf("expected '42', got '%s'", result)
	}
}

func TestJSONValueTypes(t *testing.T) {
	// Test JSON value type constant
	if JSON_VALUE != "JSON" {
		t.Errorf("expected JSON_VALUE to be 'JSON', got '%s'", JSON_VALUE)
	}
	
	// Test JSON method value type constant
	if JSON_METHOD_VALUE != "JSON_METHOD" {
		t.Errorf("expected JSON_METHOD_VALUE to be 'JSON_METHOD', got '%s'", JSON_METHOD_VALUE)
	}
	
	// Test JSON namespace value type constant
	if JSON_NAMESPACE_VALUE != "JSON_NAMESPACE" {
		t.Errorf("expected JSON_NAMESPACE_VALUE to be 'JSON_NAMESPACE', got '%s'", JSON_NAMESPACE_VALUE)
	}
	
	// Test JSON object creation and type
	data := &String{Value: "test"}
	jsonObj := &JSON{Data: data}
	
	if jsonObj.Type() != JSON_VALUE {
		t.Errorf("expected JSON type, got %s", jsonObj.Type())
	}
	
	// Test JSON method creation and type
	jsonMethod := &JSONMethod{JSON: jsonObj, Method: "get"}
	
	if jsonMethod.Type() != JSON_METHOD_VALUE {
		t.Errorf("expected JSON_METHOD type, got %s", jsonMethod.Type())
	}
	
	// Test JSON namespace creation and type
	jsonNamespace := &JSONNamespace{}
	
	if jsonNamespace.Type() != JSON_NAMESPACE_VALUE {
		t.Errorf("expected JSON_NAMESPACE type, got %s", jsonNamespace.Type())
	}
}

func TestJSONErrorHandling(t *testing.T) {
	// Test parsing invalid JSON
	result := parseJSON("{invalid json}")
	
	// Should return an error
	_, ok := result.(*Error)
	if !ok {
		t.Fatalf("expected Error result for invalid JSON, got %T", result)
	}
	
	// Test JSON namespace static method error handling
	jsonNamespace := &JSONNamespace{}
	
	// Test JSON.parse with wrong argument type
	parseResult := applyJSONNamespaceMethod(jsonNamespace, "parse", &Integer{Value: 42})
	_, ok = parseResult.(*Error)
	if !ok {
		t.Fatalf("expected Error result for wrong argument type in JSON.parse, got %T", parseResult)
	}
	
	// Test JSON.stringify with unsupported type (should not error, but test the path)
	stringifyResult := applyJSONNamespaceMethod(jsonNamespace, "stringify", &String{Value: "test"})
	str, ok := stringifyResult.(*String)
	if !ok {
		t.Fatalf("expected String result from JSON.stringify, got %T", stringifyResult)
	}
	
	if str.Value != `"test"` {
		t.Errorf("expected stringified \"test\", got %s", str.Value)
	}
}