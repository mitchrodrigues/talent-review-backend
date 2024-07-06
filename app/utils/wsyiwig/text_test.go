package wsyiwig

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractText(t *testing.T) {
	// Define test cases
	tests := []struct {
		name         string
		jsonContent  string
		expectedText string
	}{
		{
			name: "Simple text extraction",
			jsonContent: `{
				"type": "doc",
				"content": [
					{
						"type": "paragraph",
						"content": [
							{
								"type": "text",
								"text": "Hello"
							},
							{
								"type": "text",
								"text": " "
							},
							{
								"type": "text",
								"text": "World"
							}
						]
					},
					{
						"type": "paragraph",
						"content": [
							{
								"type": "text",
								"text": "!"
							}
						]
					}
				]
			}`,
			expectedText: "Hello World!",
		},
		{
			name: "Empty content",
			jsonContent: `{
				"type": "doc",
				"content": []
			}`,
			expectedText: "",
		},
		{
			name: "Nested content",
			jsonContent: `{
				"type": "doc",
				"content": [
					{
						"type": "paragraph",
						"content": [
							{
								"type": "text",
								"text": "Nested "
							},
							{
								"type": "text",
								"text": "Content"
							}
						]
					},
					{
						"type": "paragraph",
						"content": [
							{
								"type": "text",
								"text": "!"
							}
						]
					}
				]
			}`,
			expectedText: "Nested Content!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var doc Node
			err := json.Unmarshal([]byte(tt.jsonContent), &doc)
			assert.NoError(t, err)
			actualText := ExtractText(doc.Content)
			assert.Equal(t, tt.expectedText, actualText)
		})
	}
}
