package loaders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextLoader(t *testing.T) {
	var l Loader = NewTextLoader("text_test.go")
	docs, err := l.Load()
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.NotEmpty(t, docs[0].Metadata)
	assert.Contains(t, docs[0].Metadata["source"], "text_test.go")
	assert.Contains(t, docs[0].Content, "func TestTextLoader(t *testing.T)")
}
