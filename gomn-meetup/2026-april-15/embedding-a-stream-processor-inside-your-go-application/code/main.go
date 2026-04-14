package main

import (
	"context"
	"fmt"

	"github.com/warpstreamlabs/bento/public/bloblang"
	"github.com/warpstreamlabs/bento/public/service"

	_ "github.com/warpstreamlabs/bento/public/components/all"
)

func main() {
	ctx := context.Background()
	env := service.NewEnvironment()
	blobEnv := bloblang.NewEnvironment()

	env.UseBloblangEnvironment(blobEnv)

	builder := env.NewStreamBuilder()

	builder.SetYAML(`
logger:
  level: DEBUG
  format: logfmt # You can also set this to 'json' if you prefer structured logs
  add_timestamp: true

input:
  generate:
    count: 1000
    interval: 1s
    mapping: |
      root = if random_int() % 2 == 0 {
        {
          "type": "foo",
          "foo": "xyz"
        }
      } else {
        {
          "type": "bar",
          "bar": "abc"
        }
      }
output:
  kafka:
    topic: "output_topic"
	`)

	stream, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("error building: %w", err))
	}

	err = stream.Run(ctx)
	if err != nil {
		panic(fmt.Errorf("error running: %w", err))
	}
}
