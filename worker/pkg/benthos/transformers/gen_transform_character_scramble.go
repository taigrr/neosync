
// Code generated by Neosync neosync_transformer_generator.go. DO NOT EDIT.
// source: transform_character_scramble.go

package transformers

import (
)

type TransformCharacterScramble struct{}

type TransformCharacterScrambleOpts struct {
	userProvidedRegex *string
}

func NewTransformCharacterScramble() *TransformCharacterScramble {
	return &TransformCharacterScramble{}
}

func (t *TransformCharacterScramble) GetJsTemplateData() (*TemplateData, error) {
	return &TemplateData{
		Name: "transformCharacterScramble",
		Description: "Transforms an existing string value by scrambling the characters while maintaining the format.",
		Example: "",
	}, nil
}

func (t *TransformCharacterScramble) ParseOptions(opts map[string]any) (any, error) {
	transformerOpts := &TransformCharacterScrambleOpts{}

	var userProvidedRegex *string
	if arg, ok := opts["userProvidedRegex"].(string); ok {
		userProvidedRegex = &arg
	}
	transformerOpts.userProvidedRegex = userProvidedRegex

	return transformerOpts, nil
}
