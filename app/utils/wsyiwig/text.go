package wsyiwig

import "encoding/json"

// Define structs to represent the Tiptap JSON structure
type Node struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Content []Node `json:"content,omitempty"`
	Attrs   Attrs  `json:"attrs,omitempty"`
	Marks   []Mark `json:"marks,omitempty"`
}

type Attrs struct {
	Level int `json:"level,omitempty"`
}

type Mark struct {
	Type string `json:"type"`
}

// Function to recursively extract text content from Tiptap JSON
func ExtractText(nodes []Node) string {
	var text string
	for _, node := range nodes {
		if node.Type == "text" {
			text += node.Text
		}
		if len(node.Content) > 0 {
			text += ExtractText(node.Content)
		}
	}
	return text
}

// Function to extract text content from a JSON string
func ExtractTextFromJSON(jsonStr string) (string, error) {
	var root Node
	err := json.Unmarshal([]byte(jsonStr), &root)
	if err != nil {
		return "", err
	}
	return ExtractText(root.Content), nil
}
