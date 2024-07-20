package transformers

import (
	"errors"
	"fmt"
	"time"

	transformers_dataset "github.com/nucleuscloud/neosync/worker/pkg/benthos/transformers/data-sets"
	transformer_utils "github.com/nucleuscloud/neosync/worker/pkg/benthos/transformers/utils"
	"github.com/nucleuscloud/neosync/worker/pkg/rng"
	"github.com/warpstreamlabs/bento/public/bloblang"
)

// +neosyncTransformerBuilder:generate:generateFirstName

func init() {
	spec := bloblang.NewPluginSpec().
		Description("Generates a random first name.").
		Param(bloblang.NewInt64Param("max_length").Default(10000).Description("Specifies the maximum length for the generated data. This field ensures that the output does not exceed a certain number of characters.")).
		Param(bloblang.NewInt64Param("seed").Default(time.Now().UnixNano()).Description("An optional seed value used to generate deterministic outputs."))

	err := bloblang.RegisterFunctionV2("generate_first_name", spec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
		maxLength, err := args.GetInt64("max_length")
		if err != nil {
			return nil, err
		}
		seed, err := args.GetInt64("seed")
		if err != nil {
			return nil, err
		}

		randomizer := rng.New(seed)

		return func() (any, error) {
			output, err := generateRandomFirstName(randomizer, nil, maxLength)
			if err != nil {
				return nil, fmt.Errorf("unable to run generate_first_name: %w", err)
			}
			return output, nil
		}, nil
	})
	if err != nil {
		panic(err)
	}
}

func (t *GenerateFirstName) Generate(opts any) (any, error) {
	parsedOpts, ok := opts.(*GenerateFirstNameOpts)
	if !ok {
		return nil, errors.New("invalid parse opts")
	}

	return generateRandomFirstName(parsedOpts.randomizer, nil, parsedOpts.maxLength)
}

func generateRandomFirstName(randomizer rng.Rand, minLength *int64, maxLength int64) (string, error) {
	return transformer_utils.GenerateStringFromCorpus(
		randomizer,
		transformers_dataset.FirstNames,
		transformers_dataset.FirstNameMap,
		transformers_dataset.FirstNameIndices,
		minLength,
		maxLength,
		nil,
	)
}
