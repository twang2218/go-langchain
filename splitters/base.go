package splitters

import "langchain/docstore"

type Splitter interface {
	Split(doc docstore.Document) ([]docstore.Document, error)
}

type FuncNewSplitter func() Splitter
