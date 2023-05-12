package llms

import "context"

type LLM interface {
	// Chat 与用户进行对话
	Chat(ctx context.Context, prompt string) (string, error)
}
