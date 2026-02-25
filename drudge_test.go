package godrudge

import (
	"fmt"
	"os"
	"testing"
)

var c *Client

func TestMain(m *testing.M) {
	c = NewClient()
	err := c.ParseRSS()
	if err != nil {
		fmt.Printf("Failed to parse RSS feed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

// this test is just to print
func TestPrint(t *testing.T) {
	fmt.Println(c.Page)
}

func TestColumnHeadlines(t *testing.T) {
	if len(c.Page.HeadlineColumns) == 0 {
		t.Fatal("No headline columns found")
	}

	for i := 0; i < len(c.Page.HeadlineColumns); i++ {
		if len(c.Page.HeadlineColumns[i]) == 0 {
			t.Fatalf("No headlines found in column %d", i)
		}
	}
}

func TestMainHeadlines(t *testing.T) {
	if len(c.Page.MainHeadlines) == 0 {
		t.Fatal("No main headlines found")
	}
}

func TestTopHeadlines(t *testing.T) {
	if len(c.Page.TopHeadlines) == 0 {
		t.Fatal("No top headlines found")
	}
}
