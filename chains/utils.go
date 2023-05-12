package chains

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

func ChainApply(c Chain, ctx context.Context, inputs []map[string]string) ([]string, error) {
	var outputs []string

	for _, input := range inputs {
		output, err := c.Run(ctx, input)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func ChainApplyAsync(c Chain, ctx context.Context, inputs []map[string]string) ([]string, error) {
	// logrus.Tracef("ChainApplyAsync: %d inputs", len(inputs))

	//	parallelize
	type result struct {
		index  int // keep the order
		output string
	}
	ch := make(chan result, len(inputs))
	wg := sync.WaitGroup{}
	wg.Add(len(inputs))

	//	producers
	for i, input := range inputs {
		go func(index int, input map[string]string) {
			defer wg.Done()

			output, err := c.Run(ctx, input)
			if err != nil {
				logrus.Warnf("error running chain: %v", err)
			}
			ch <- result{index, output}
		}(i, input)
	}

	wg.Wait()
	close(ch)

	//	consumer
	outputs := make([]string, len(inputs))
	for output := range ch {
		if output.index >= len(inputs) {
			logrus.Warnf("index out of range: %d", output.index)
			continue
		}
		outputs[output.index] = output.output
	}

	return outputs, nil
}
