package loaders

import (
	"io"
	"langchain/docstore"
	"os"

	log "github.com/sirupsen/logrus"
)

type TextLoader struct {
	filename string
}

func NewTextLoader(filename string) Loader {
	return &TextLoader{filename: filename}
}

func (l *TextLoader) Load() ([]docstore.Document, error) {
	log.Tracef("Loading file %s...", l.filename)
	f, err := os.Open(l.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return []docstore.Document{
		{
			Content: string(content),
			Metadata: map[string]string{
				"source": l.filename,
			},
		},
	}, nil
}
