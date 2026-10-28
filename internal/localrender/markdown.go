package localrender

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var markdownEngine = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
)

const monoFont = `Menlo,Monaco,Consolas,'Courier New',monospace`

// markdownToHTML renders a Markdown fragment and inlines theme styles on every
// element so the result survives the WeChat editor.
func (r *renderer) markdownToHTML(source string) string {
	if strings.TrimSpace(source) == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := markdownEngine.Convert([]byte(source), &buf); err != nil {
		return "<p>" + html.EscapeString(source) + "</p>"
	}
	return r.styleFragment(buf.String())
}

// inline renders a short Markdown value (one field) and drops the wrapping
// paragraph so it can sit inside another element.
func (r *renderer) inline(source string) string {
	source = strings.TrimSpace(source)
	if source == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := markdownEngine.Convert([]byte(source), &buf); err != nil {
		return html.EscapeString(source)
	}
	nodes := parseFragment(buf.String())
	var kept []*html.Node
	for _, n := range nodes {
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) == "" {
			continue
		}
		kept = append(kept, n)
	}
	if len(kept) == 1 && kept[0].DataAtom == atom.P {
		var out strings.Builder
		for c := kept[0].FirstChild; c != nil; c = c.NextSibling {
			r.styleNode(c)
			_ = html.Render(&out, c)
		}
		return out.String()
	}
	return r.styleFragment(buf.String())
}

func parseFragment(fragment string) []*html.Node {
	context := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	nodes, err := html.ParseFragment(strings.NewReader(fragment), context)
	if err != nil {
		return nil
	}
	return nodes
}

func (r *renderer) styleFragment(fragment string) string {
	var out strings.Builder
	for _, n := range parseFragment(fragment) {
		r.styleNode(n)
		_ = html.Render(&out, n)
	}
	return out.String()
}

func (r *renderer) styleNode(n *html.Node) {
	if n.Type == html.ElementNode {
		if style := r.elementStyle(n); style != "" {
			setStyle(n, style)
		}
		if n.DataAtom == atom.Input {
			convertTaskCheckbox(n)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.styleNode(c)
	}
}

// convertTaskCheckbox swaps GFM task-list checkboxes, which WeChat strips, for
// plain glyphs.
func convertTaskCheckbox(n *html.Node) {
	checked := false
	for _, a := range n.Attr {
		if a.Key == "checked" {
			checked = true
		}
	}
	glyph := "☐ "
	if checked {
		glyph = "☑ "
	}
	n.Type = html.TextNode
	n.Data = glyph
	n.Attr = nil
	n.DataAtom = 0
}

func setStyle(n *html.Node, style string) {
	for i, a := range n.Attr {
		if a.Key == "style" {
			n.Attr[i].Val = style + ";" + a.Val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: "style", Val: style})
}

func hasAncestor(n *html.Node, tags ...atom.Atom) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		for _, tag := range tags {
			if p.DataAtom == tag {
				return true
			}
		}
	}
	return false
}

func (r *renderer) elementStyle(n *html.Node) string {
	t := r.theme
	switch n.DataAtom {
	case atom.H1:
		return r.headingStyle(1)
	case atom.H2:
		return r.headingStyle(2)
	case atom.H3:
		return r.headingStyle(3)
	case atom.H4, atom.H5, atom.H6:
		return css("font-size", "16px", "font-weight", "bold", "color", t.Text, "margin", "24px 0 12px", "line-height", "1.5")
	case atom.P:
		if hasAncestor(n, atom.Blockquote, atom.Li, atom.Td, atom.Th) {
			return css("margin", "0", "line-height", "1.75", "color", "inherit")
		}
		return r.paragraphStyle()
	case atom.Strong, atom.B:
		return css("font-weight", "bold", "color", t.Primary)
	case atom.Em, atom.I:
		return css("font-style", "italic")
	case atom.Del, atom.S:
		return css("text-decoration", "line-through", "color", t.Muted)
	case atom.A:
		return css("color", t.Primary, "text-decoration", "none", "border-bottom", "1px dashed "+t.Primary, "word-break", "break-all")
	case atom.Code:
		if hasAncestor(n, atom.Pre) {
			return css("font-family", monoFont, "color", t.CodeFg, "background", "none", "font-size", "13px", "white-space", "pre")
		}
		return css("font-family", monoFont, "font-size", "90%", "color", t.Primary, "background", t.Soft, "padding", "2px 6px", "border-radius", "4px", "word-break", "break-all")
	case atom.Pre:
		return css("background", t.CodeBg, "border-radius", "8px", "padding", "16px", "margin", "20px 0", "overflow-x", "auto", "line-height", "1.6", "font-size", "13px", "border", "1px solid "+t.Border)
	case atom.Blockquote:
		return css("margin", "20px 0", "padding", "12px 16px", "border-left", "4px solid "+t.Primary, "background", t.Soft, "color", t.Muted, "font-size", r.fontSize, "line-height", "1.75", "border-radius", "0 6px 6px 0")
	case atom.Ul, atom.Ol:
		if hasAncestor(n, atom.Li) {
			return css("margin", "4px 0", "padding-left", "20px", "color", t.Text)
		}
		return css("margin", "16px 0", "padding-left", "24px", "color", t.Text, "font-size", r.fontSize, "line-height", "1.75")
	case atom.Li:
		return css("margin", "6px 0", "line-height", "1.75")
	case atom.Table:
		return css("width", "100%", "border-collapse", "collapse", "margin", "20px 0", "font-size", "14px", "color", t.Text)
	case atom.Th:
		return css("border", "1px solid "+t.Border, "padding", "8px 10px", "background", t.Soft, "color", t.Primary, "font-weight", "bold", "text-align", "left")
	case atom.Td:
		return css("border", "1px solid "+t.Border, "padding", "8px 10px", "line-height", "1.6")
	case atom.Hr:
		return css("border", "0", "border-top", "1px solid "+t.Border, "margin", "32px 0", "height", "0")
	case atom.Img:
		return css("max-width", "100%", "height", "auto", "display", "block", "margin", "20px auto", "border-radius", "6px")
	}
	return ""
}

func (r *renderer) paragraphStyle() string {
	return css("margin", "0 0 16px", "line-height", "1.75", "font-size", r.fontSize, "color", r.theme.Text, "letter-spacing", "0.5px", "text-align", "justify")
}

func (r *renderer) headingStyle(level int) string {
	t := r.theme
	size := map[int]string{1: "22px", 2: "19px", 3: "17px"}[level]
	base := []string{"font-size", size, "font-weight", "bold", "line-height", "1.45", "margin", "32px 0 16px"}
	if level == 3 {
		base = []string{"font-size", size, "font-weight", "bold", "line-height", "1.5", "margin", "24px 0 12px"}
	}
	switch t.Series {
	case SeriesMinimal:
		color := t.Primary
		if level == 3 {
			color = t.Text
		}
		return css(append(base, "color", color)...)
	case SeriesFocus:
		if level == 3 {
			return css(append(base, "color", t.Primary, "text-align", "center")...)
		}
		return css(append(base, "color", t.Primary, "text-align", "center", "border-top", "1px solid "+t.Primary, "border-bottom", "1px solid "+t.Primary, "padding", "10px 0")...)
	case SeriesElegant:
		width := map[int]string{1: "6px", 2: "4px", 3: "3px"}[level]
		return css(append(base, "color", t.Primary, "border-left", width+" solid "+t.Primary, "padding", "6px 12px", "background", "linear-gradient(90deg,"+t.Soft+",rgba(255,255,255,0))", "border-radius", "0 6px 6px 0")...)
	case SeriesBold:
		if level == 3 {
			return css(append(base, "color", t.Primary)...)
		}
		return css(append(base, "display", "table", "color", t.onPrimary(), "background", t.Primary, "padding", "6px 16px", "border-radius", "8px", "box-shadow", "0 4px 12px rgba(0,0,0,0.12)")...)
	default: // classic
		if level == 3 {
			return css(append(base, "color", t.Primary)...)
		}
		return css(append(base, "display", "table", "color", t.Primary, "border-bottom", "2px solid "+t.Primary, "padding", "0 2px 6px")...)
	}
}

// css joins property/value pairs into an inline style string.
func css(pairs ...string) string {
	var b strings.Builder
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(';')
		}
		b.WriteString(pairs[i])
		b.WriteByte(':')
		b.WriteString(pairs[i+1])
	}
	return b.String()
}
