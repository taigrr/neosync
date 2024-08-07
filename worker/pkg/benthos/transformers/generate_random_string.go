package transformers

import (
	"errors"
	"fmt"

	transformer_utils "github.com/nucleuscloud/neosync/worker/pkg/benthos/transformers/utils"
	"github.com/warpstreamlabs/bento/public/bloblang"
)

// +neosyncTransformerBuilder:generate:generateRandomString

func init() {
	spec := bloblang.NewPluginSpec().
		Description("Creates a randomly ordered alphanumeric string with a default length of 10 unless the String Length parameter are defined.").
		Param(bloblang.NewInt64Param("min").Description("Specifies the minimum length for the generated string.")).
		Param(bloblang.NewInt64Param("max").Description("Specifies the maximum length for the generated string."))

	err := bloblang.RegisterFunctionV2("generate_string", spec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
		min, err := args.GetInt64("min")
		if err != nil {
			return nil, err
		}

		max, err := args.GetInt64("max")
		if err != nil {
			return nil, err
		}

		return func() (any, error) {
			out, err := transformer_utils.GenerateRandomStringWithInclusiveBounds(min, max)
			if err != nil {
				return nil, fmt.Errorf("unable to run generate_string: %w", err)
			}
			return out, nil
		}, nil
	})

	if err != nil {
		panic(err)
	}
}

func (t *GenerateRandomString) Generate(opts any) (any, error) {
	parsedOpts, ok := opts.(*GenerateRandomStringOpts)
	if !ok {
		return nil, errors.New("invalid parse opts")
	}

	return transformer_utils.GenerateRandomStringWithInclusiveBounds(parsedOpts.min, parsedOpts.max)
}
