package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"gopkg.in/russross/blackfriday.v2"
	bf "gopkg.in/russross/blackfriday.v2"
)

var (
	// Search replace pattern for camel-cased words.
	validPage = regexp.MustCompile("([A-Z][a-z]+[A-Z][a-zA-Z]+)")
)

func main() {
	// Example string demonstrating functionality inherited from existing blackfriday.HTMLRenderer in addition to some better wiki-linking.
	input := []byte("*This* is CamelCase. This is http://CamelCaseInsideAUrl.com")
	r := NewRenderer()
	output := bf.Run(input, blackfriday.WithRenderer(r))
	fmt.Println(string(output))
}

type Renderer struct {
	htmlRenderer *bf.HTMLRenderer
}

func NewRenderer() *Renderer {
	// Compose custom Renderer of the existing blackfriday.HTMLRenderer.
	hr := bf.NewHTMLRenderer(bf.HTMLRendererParameters{
		Flags: bf.CommonHTMLFlags,
	})
	return &Renderer{
		htmlRenderer: hr,
	}
}

func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	// Only touch camel-cased words when they are not inside of a Link node.
	if node.Parent != nil && node.Parent.Type != bf.Link {
		if entering {
			md := r.wikLink(node.Literal)
			if md != nil {
				w.Write(md)
				return bf.SkipChildren
			}
		}
	}

	// Pipe the node back through stock blackfriday.HTMLRenderer.
	return r.htmlRenderer.RenderNode(w, node, entering)
}

func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
	r.htmlRenderer.RenderHeader(w, ast)
}

func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
	r.htmlRenderer.RenderHeader(w, ast)
}

func (r *Renderer) wikLink(b []byte) []byte {
	// We don't care about anything less than two characters b/c that wouldn't qualify as camel-case.
	if len(b) <= 2 {
		return nil
	}

	markdown := validPage.ReplaceAll(
		b,
		[]byte(fmt.Sprintf("[$1](%s$1)", "/view/")),
	)

	if bytes.Compare(b, markdown) == 0 {
		return nil
	}

	return markdown
}
