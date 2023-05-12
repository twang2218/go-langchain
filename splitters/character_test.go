package splitters

import (
	"fmt"
	"langchain/docstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharacterSpliter(t *testing.T) {
	sample_doc := docstore.Document{
		Content: "The Maowusu Desert in Northern China was transformed from a barren wasteland into a thriving ecosystem through a comprehensive desertification control project.",
		Metadata: map[string]string{
			"source": "test",
		},
	}

	// No overlap
	var splitter Splitter = &CharacterSplitter{Separator: " ", ChunkSize: 50, ChunkOverlap: 0}

	docs, err := splitter.Split(sample_doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, docs)
	assert.Equal(t, 4, len(docs), fmt.Sprintf("Expected 4 chunks, got %v", docs))
	assert.Equal(t, "The Maowusu Desert in Northern China was", docs[0].Content)
	assert.Equal(t, "transformed from a barren wasteland into a", docs[1].Content)
	assert.Equal(t, "thriving ecosystem through a comprehensive", docs[2].Content)
	assert.Equal(t, "desertification control project.", docs[3].Content)

	// With overlap by 10 characters
	splitter = &CharacterSplitter{Separator: " ", ChunkSize: 50, ChunkOverlap: 10}

	docs, err = splitter.Split(sample_doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, docs)
	assert.Equal(t, 4, len(docs))
	assert.Equal(t, "The Maowusu Desert in Northern China was", docs[0].Content)
	assert.Equal(t, "China was transformed from a barren wasteland into", docs[1].Content)
	assert.Equal(t, "into a thriving ecosystem through a comprehensive", docs[2].Content)
	assert.Equal(t, "desertification control project.", docs[3].Content)
}
