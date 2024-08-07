
// Code generated by Neosync neosync_transformer_generator.go. DO NOT EDIT.
// source: transform_int64.go

package transformers

import (
	"fmt"
	
)

type TransformInt64 struct{}

type TransformInt64Opts struct {
	randomizationRangeMin int64
	randomizationRangeMax int64
}

func NewTransformInt64() *TransformInt64 {
	return &TransformInt64{}
}

func (t *TransformInt64) GetJsTemplateData() (*TemplateData, error) {
	return &TemplateData{
		Name: "transformInt64",
		Description: "Transforms an existing integer value.",
		Example: "",
	}, nil
}

func (t *TransformInt64) ParseOptions(opts map[string]any) (any, error) {
	transformerOpts := &TransformInt64Opts{}

	if _, ok := opts["randomizationRangeMin"].(int64); !ok {
		return nil, fmt.Errorf("missing required argument. function: %s argument: %s", "transformInt64", "randomizationRangeMin")
	}
	randomizationRangeMin := opts["randomizationRangeMin"].(int64)
	transformerOpts.randomizationRangeMin = randomizationRangeMin

	if _, ok := opts["randomizationRangeMax"].(int64); !ok {
		return nil, fmt.Errorf("missing required argument. function: %s argument: %s", "transformInt64", "randomizationRangeMax")
	}
	randomizationRangeMax := opts["randomizationRangeMax"].(int64)
	transformerOpts.randomizationRangeMax = randomizationRangeMax

	return transformerOpts, nil
}
