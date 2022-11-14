package cleanup

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type ParsedHTML struct {
	RootNode *html.Node
}

func NewParsedHTML(input io.Reader) (*ParsedHTML, error) {
	rootNode, err := html.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("call to html.Parse failed: %w", err)
	}
	return &ParsedHTML{rootNode}, nil
}

func (a *ParsedHTML) RenderHTML() (string, error) {
	if a.RootNode == nil {
		return "", fmt.Errorf("cannot render; root node is nil")
	}
	var buf strings.Builder
	err := html.Render(&buf, a.RootNode)
	if err != nil {
		return "", fmt.Errorf("call to html render failed: %w", err)
	}
	return buf.String(), nil
}

func (a *ParsedHTML) Text() string {
	var buf strings.Builder
	text(a.RootNode, &buf)
	return buf.String()
}

func text(n *html.Node, buf *strings.Builder) {
	if n == nil {
		return
	}
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	if n.FirstChild != nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			text(c, buf)
		}
	}
}

func (a *ParsedHTML) UpdateParsedFromString(inputHTML string) error {
	rootNode, err := html.Parse(strings.NewReader(inputHTML))
	if err != nil {
		return fmt.Errorf("call to html.Parse failed: %w", err)
	}
	a.RootNode = rootNode
	return nil
}

func (a *ParsedHTML) GoQueryDoc() *goquery.Document {
	return goquery.NewDocumentFromNode(a.RootNode)
}
