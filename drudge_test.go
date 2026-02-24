package godrudge

import (
	"fmt"
	"testing"
)

// TODO: write real tests
func TestPrint(t *testing.T) {
	c := NewClient()
	err := c.ParseRSS()
	fmt.Println(c.Page)
	if err != nil {
		t.Fatalf("Failed to parse RSS feed: %v", err)
	}
}
