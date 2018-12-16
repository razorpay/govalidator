package govalidator

import (
	"encoding/json"
	"testing"
)

func TestValidator_SetDefaultRequired(t *testing.T) {
	v := New(Options{})
	v.SetDefaultRequired(true)
	if !v.Opts.RequiredDefault {
		t.Error("SetDefaultRequired failed")
	}
}

func TestValidator_Validate(t *testing.T) {
	params := map[string]interface{}{
		"name":     "John Doe",
		"username": "jhondoe",
		"email":    "john@mail.com",
		"zip":      "8233",
	}

	rulesList := MapData{
		"name":  []string{"required"},
		"age":   []string{"between:5,16"},
		"email": []string{"email"},
		"zip":   []string{"digits:4"},
	}

	opts := Options{
		Request: params,
		Rules:   rulesList,
	}
	v := New(opts)
	validationError := v.Validate()
	if len(validationError) > 0 {
		t.Log(validationError)
		t.Error("Validate failed to validate correct inputs!")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Validate did not panic")
		}
	}()

	v1 := New(Options{Rules: MapData{}})
	v1.Validate()
}

func Benchmark_Validate(b *testing.B) {
	params := map[string]interface{}{
		"name":     "John Doe",
		"username": "jhondoe",
		"email":    "john@mail.com",
		"zip":      "8233",
	}
	rulesList := MapData{
		"name":  []string{"required"},
		"age":   []string{"numeric_between:18,60"},
		"email": []string{"email"},
		"zip":   []string{"digits:4"},
	}

	opts := Options{
		Request: params,
		Rules:   rulesList,
	}
	v := New(opts)
	for n := 0; n < b.N; n++ {
		v.Validate()
	}
}

//============ validate json test ====================

func TestValidator_ValidateJson(t *testing.T) {
	type User struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Age     int    `json:"age"`
		Zip     string `json:"zip"`
		Color   int    `json:"color"`
	}

	postUser := User{
		Name:    "",
		Email:   "inalid email",
		Address: "",
		Age:     1,
		Zip:     "122",
		Color:   5,
	}

	rules := MapData{
		"name":    []string{"required"},
		"email":   []string{"email"},
		"address": []string{"required", "between:3,5"},
		"age":     []string{"bool"},
		"zip":     []string{"len:4"},
		"color":   []string{"min:10"},
	}

	var user User

	body, _ := json.Marshal(postUser)

	var req map[string]interface{}
	json.Unmarshal(body, &req)

	opts := Options{
		Request: req,
		Data:    &user,
		Rules:   rules,
	}

	vd := New(opts)
	validationErr := vd.Validate()
	if len(validationErr) != 5 {
		t.Error("Validate failed")
	}
}

func TestValidator_Validate_NULLValue(t *testing.T) {
	type User struct {
		Name   string `json:"name"`
		Count  Int    `json:"count"`
		Option Int    `json:"option"`
		Active Bool   `json:"active"`
	}

	rules := MapData{
		"name":   []string{"required"},
		"count":  []string{"required"},
		"option": []string{"required"},
		"active": []string{"required"},
	}

	postUser := map[string]interface{}{
		"name":   "John Doe",
		"count":  0,
		"option": nil,
		"active": nil,
	}

	var user User
	body, _ := json.Marshal(postUser)
	var req map[string]interface{}
	json.Unmarshal(body, &req)

	opts := Options{
		Request: req,
		Data:    &user,
		Rules:   rules,
	}

	vd := New(opts)
	validationErr := vd.Validate()
	if len(validationErr) < 1 {
		t.Error("Validate failed")
	}
}

func TestValidator_Validate_panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Validate did not panic")
		}
	}()

	opts := Options{}

	vd := New(opts)
	validationErr := vd.Validate()
	if len(validationErr) != 5 {
		t.Error("Validate failed")
	}
}
