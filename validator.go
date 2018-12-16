package govalidator

import (
	"fmt"
	_ "net/http"
	"strings"
)

type (
	// MapData represents basic data structure for govalidator Rules and Messages
	MapData map[string][]string

	// Options describes configuration option for validator
	Options struct {
		Data            interface{} // Data represents structure for JSON body
		Request         map[string]interface{}
		RequiredDefault bool    // RequiredDefault represents if all the fields are by default required or not
		Rules           MapData // Rules represents rules for form-data/x-url-encoded/query params data
		Messages        MapData // Messages represents custom/localize message for rules
		TagIdentifier   string  // TagIdentifier represents struct tag identifier, e.g: json or validate etc
		FormSize        int64   //Form represents the multipart forom data max memory size in bytes
	}

	// Validator represents a validator with options
	Validator struct {
		Opts Options // Opts contains all the options for validator
	}
)

// New return a new validator object using provided options
func New(opts Options) *Validator {
	return &Validator{Opts: opts}
}

// getMessage return if a custom message exist against the field name and rule
// if not available it return an empty string
func (v *Validator) getCustomMessage(field, rule string) string {
	if msgList, ok := v.Opts.Messages[field]; ok {
		for _, m := range msgList {
			//if rules has params, remove params. e.g: between:3,5 would be between
			if strings.Contains(rule, ":") {
				rule = strings.Split(rule, ":")[0]
			}
			if strings.HasPrefix(m, rule+":") {
				return strings.TrimPrefix(m, rule+":")
			}
		}
	}
	return ""
}

// SetDefaultRequired change the required behavior of fields
// Default value if false
// If SetDefaultRequired set to true then it will mark all the field in the rules list as required
func (v *Validator) SetDefaultRequired(required bool) {
	v.Opts.RequiredDefault = required
}

// Validate validate request data like form-data, x-www-form-urlencoded and query params
// see example in README.md file
// ref: https://github.com/thedevsaddam/govalidator#example
func (v *Validator) Validate() MapData {
	// if request object and rules not passed rise a panic
	if len(v.Opts.Rules) == 0 || v.Opts.Request == nil {
		panic(errValidateArgsMismatch)
	}
	errsBag := MapData{}

	// get non required rules
	nr := v.getNonRequiredFields()

	for field, rules := range v.Opts.Rules {
		if _, ok := nr[field]; ok {
			continue
		}
		for _, rule := range rules {
			if !isRuleExist(rule) {
				panic(fmt.Errorf("govalidator: %s is not a valid rule", rule))
			}
			msg := v.getCustomMessage(field, rule)

			// validate if custom rules exist
			reqVal := v.Opts.Request[field]
			validateCustomRules(field, rule, msg, reqVal, errsBag)

		}
	}

	return errsBag
}

// getNonRequiredFields remove non required rules fields from rules if requiredDefault field is false
// and if the input data is empty for this field
func (v *Validator) getNonRequiredFields() map[string]struct{} {
	inputs := v.Opts.Request
	nr := make(map[string]struct{})
	if !v.Opts.RequiredDefault {
		for k, r := range v.Opts.Rules {
			if _, ok := inputs[k]; !ok {
				if !isContainRequiredField(r) {
					nr[k] = struct{}{}
				}
			}
		}
	}
	return nr
}
