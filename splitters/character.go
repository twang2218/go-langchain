package splitters

import (
	"langchain/docstore"
	"strings"
)

type CharacterSplitter struct {
	Separator    string
	ChunkSize    int
	ChunkOverlap int
}

func NewCharacterSplitter() Splitter {
	return &CharacterSplitter{
		Separator:    "\n",
		ChunkSize:    1000,
		ChunkOverlap: 0,
	}
}

func (s *CharacterSplitter) Split(doc docstore.Document) ([]docstore.Document, error) {
	splits := strings.Split(doc.Content, s.Separator)
	chunks := make([]string, 0, len(splits))

	var begin int
	for i := 0; i < len(splits); i++ {
		current_length := strings.Join(splits[begin:i+1], s.Separator)
		if len(current_length) > s.ChunkSize {
			chunks = append(chunks, strings.Join(splits[begin:i], s.Separator))
			//	calculate Overlap backward steps
			for j := i; j >= 0; j-- {
				if len(strings.Join(splits[j:i], s.Separator)) > s.ChunkOverlap {
					begin = j + 1
					break
				}
			}
		}
	}
	if begin < len(splits) {
		chunks = append(chunks, strings.Join(splits[begin:], s.Separator))
	}

	var docs []docstore.Document
	for _, chunk := range chunks {
		if len(chunk) > 0 {
			docs = append(docs, docstore.Document{
				Content:  chunk,
				Metadata: doc.Metadata,
			})
		}
	}

	return docs, nil
}
