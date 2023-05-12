package loaders

import (
	"langchain/docstore"
)

type Loader interface {
	Load() ([]docstore.Document, error)
}

type FuncNewFileLoader func(filename string) Loader
