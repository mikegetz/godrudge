package godrudge

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractNodeAttr(t *testing.T) {
	tests := []struct {
		name     string
		htmlStr  string
		attrKey  string
		expected string
	}{
		{
			name:     "existing attribute",
			htmlStr:  `<div id="test" class="example"></div>`,
			attrKey:  "id",
			expected: "test",
		},
		{
			name:     "non-existing attribute",
			htmlStr:  `<div id="test" class="example"></div>`,
			attrKey:  "style",
			expected: "",
		},
		{
			name:     "empty attribute key",
			htmlStr:  `<div id="test" class="example"></div>`,
			attrKey:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tt.htmlStr))
			if err != nil {
				t.Fatalf("failed to parse html: %v", err)
			}

			var node *html.Node
			var f func(*html.Node)
			f = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" {
					node = n
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
			f(doc)

			if node == nil {
				t.Fatalf("failed to find div node")
			}

			got := extractNodeAttr(node, tt.attrKey)
			if got != tt.expected {
				t.Errorf("extractNodeAttr() = %v, want %v", got, tt.expected)
			}
		})
	}
}
