package chains

import "context"

type Chain interface {
	Run(ctx context.Context, inputs map[string]string) (string, error)
	RunText(ctx context.Context, input string) (string, error)
	Apply(ctx context.Context, inputs []map[string]string) ([]string, error)
}
