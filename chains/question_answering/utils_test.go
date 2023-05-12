package question_answering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDocsFromInputs(t *testing.T) {
	inputs := map[string]string{
		KeyDocuments: "[\"doc1\", \"doc2\"]",
	}
	docs, err := GetDocsFromInputs(inputs, KeyDocuments)
	assert.NoError(t, err)
	assert.Len(t, docs, 2)
	assert.Equal(t, "doc1", docs[0])
	assert.Equal(t, "doc2", docs[1])
}

func TestPutDocsToInputs(t *testing.T) {
	docs := []string{"doc1", "doc2"}
	inputs := map[string]string{}
	err := PutDocsToInputs(docs, inputs, KeyDocuments)
	assert.NoError(t, err)
	assert.NotEmpty(t, inputs[KeyDocuments])
	assert.Equal(t, "[\"doc1\",\"doc2\"]", inputs[KeyDocuments])
}
