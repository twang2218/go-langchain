package loaders

import (
	"langchain/docstore"
	"path/filepath"
)

type DirectoryLoader struct {
	Path              string            //	文档所在的目录
	Glob              string            //	文件名的正则表达式
	FuncNewFileLoader FuncNewFileLoader //	文件加载器
}

func NewDirectoryLoader(path string) Loader {
	return &DirectoryLoader{
		Path:              path,
		Glob:              "*",
		FuncNewFileLoader: NewTextLoader,
	}
}

func (l *DirectoryLoader) Load() ([]docstore.Document, error) {
	if l.Glob == "" {
		l.Glob = "*"
	}
	if l.FuncNewFileLoader == nil {
		l.FuncNewFileLoader = NewTextLoader
	}

	f := filepath.Join(l.Path, l.Glob)
	files, err := filepath.Glob(f)
	if err != nil {
		return nil, err
	}

	docs := make([]docstore.Document, 0, len(files))

	for _, file := range files {
		loader := l.FuncNewFileLoader(file)
		dd, err := loader.Load()
		if err != nil {
			return nil, err
		}
		docs = append(docs, dd...)
	}
	return docs, nil
}
