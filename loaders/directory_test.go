package loaders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryLoader(t *testing.T) {
	var l Loader = NewDirectoryLoader(".")
	docs, err := l.Load()
	assert.NoError(t, err)
	assert.NotEmpty(t, docs)

	var content string
	sources := make([]string, 0, len(docs))
	for _, doc := range docs {
		sources = append(sources, doc.Metadata["source"])
		content += doc.Content
	}
	assert.Contains(t, sources, "directory_test.go")
	assert.Contains(t, sources, "text_test.go")
	assert.Contains(t, sources, "base.go")

	assert.Contains(t, content, "func TestDirectoryLoader(t *testing.T)")
	assert.Contains(t, content, "func TestTextLoader(t *testing.T)")
}
