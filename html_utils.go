package godrudge

import "golang.org/x/net/html"

func extractNodeAttr(n *html.Node, attrKey string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}
