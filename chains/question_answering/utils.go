package question_answering

import "encoding/json"

// GetDocsFromInputs parse the input docs from the inputs map
func GetDocsFromInputs(inputs map[string]string, key string) ([]string, error) {
	docs_str := inputs[key]
	var docs []string
	err := json.Unmarshal([]byte(docs_str), &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

// PutDocsToInputs put the docs to the inputs map
func PutDocsToInputs(docs []string, inputs map[string]string, key string) error {
	docs_str, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	inputs[key] = string(docs_str)
	return nil
}
