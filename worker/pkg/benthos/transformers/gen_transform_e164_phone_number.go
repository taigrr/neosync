
// Code generated by Neosync neosync_transformer_generator.go. DO NOT EDIT.
// source: transform_e164_phone_number.go

package transformers

import (
	"fmt"
	
)

type TransformE164PhoneNumber struct{}

type TransformE164PhoneNumberOpts struct {
	preserveLength bool
	maxLength *int64
}

func NewTransformE164PhoneNumber() *TransformE164PhoneNumber {
	return &TransformE164PhoneNumber{}
}

func (t *TransformE164PhoneNumber) GetJsTemplateData() (*TemplateData, error) {
	return &TemplateData{
		Name: "transformE164PhoneNumber",
		Description: "Transforms an existing E164 formatted phone number.",
		Example: "",
	}, nil
}

func (t *TransformE164PhoneNumber) ParseOptions(opts map[string]any) (any, error) {
	transformerOpts := &TransformE164PhoneNumberOpts{}

	if _, ok := opts["preserveLength"].(bool); !ok {
		return nil, fmt.Errorf("missing required argument. function: %s argument: %s", "transformE164PhoneNumber", "preserveLength")
	}
	preserveLength := opts["preserveLength"].(bool)
	transformerOpts.preserveLength = preserveLength

	var maxLength *int64
	if arg, ok := opts["maxLength"].(int64); ok {
		maxLength = &arg
	}
	transformerOpts.maxLength = maxLength

	return transformerOpts, nil
}
