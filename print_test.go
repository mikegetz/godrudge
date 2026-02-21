package godrudge

import "testing"

func TestPrint(t *testing.T) {
	c := NewClient()
	err := c.ParseRSS()
	c.PrintDrudge(true)
	if err != nil {
		t.Fatalf("Failed to parse RSS feed: %v", err)
	}
}
