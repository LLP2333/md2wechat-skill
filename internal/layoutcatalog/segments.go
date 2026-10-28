package layoutcatalog

import "strings"

// Block is one :::module directive found in a Markdown document.
type Block struct {
	Name    string
	Caption string
	Params  map[string]string
	Body    []string
	Line    int
	Spec    *LayoutSpec
}

// Segment is either plain Markdown or one layout block, in document order.
type Segment struct {
	Markdown string
	Block    *Block
}

// Segments splits Markdown into plain Markdown runs and :::module blocks using
// the same scanning rules as Validate. Lines that cannot be parsed as a block
// opener, and unterminated blocks, stay in the plain Markdown run.
func (c *Catalog) Segments(markdown string) []Segment {
	var segments []Segment
	var plain []string
	flush := func() {
		if len(plain) == 0 {
			return
		}
		segments = append(segments, Segment{Markdown: strings.Join(plain, "\n")})
		plain = nil
	}

	lines := strings.Split(markdown, "\n")
	var outsideFence markdownFence
	i := 0
	for i < len(lines) {
		line := strings.TrimRight(lines[i], "\r")
		if outsideFence.consume(line) {
			plain = append(plain, lines[i])
			i++
			continue
		}
		opener, err := parseBlockOpener(line)
		if err != nil {
			plain = append(plain, lines[i])
			i++
			continue
		}
		j := i + 1
		var body []string
		var bodyFence markdownFence
		depth := 0
		fenceAware := opener.Name == "expand"
		for j < len(lines) {
			content := strings.TrimRight(lines[j], "\r")
			if fenceAware && bodyFence.consume(content) {
				body = append(body, lines[j])
				j++
				continue
			}
			trimmed := strings.TrimSpace(content)
			if trimmed == ":::" {
				if !fenceAware || depth == 0 {
					break
				}
				depth--
			} else if fenceAware && nestedLayoutOpener(trimmed) {
				depth++
			}
			body = append(body, lines[j])
			j++
		}
		if j >= len(lines) {
			plain = append(plain, lines[i:]...)
			break
		}

		block := &Block{Name: opener.Name, Caption: opener.Caption, Params: opener.Params, Body: body, Line: i + 1}
		if spec, ok := c.Get(opener.Name); ok {
			block.Spec = spec
			if normalized, err := validateOpener(opener, spec.Opener); err == nil {
				block.Caption = normalized.Caption
				block.Params = normalized.Params
			}
		}
		if block.Params == nil {
			block.Params = map[string]string{}
		}
		flush()
		segments = append(segments, Segment{Block: block})
		i = j + 1
	}
	flush()
	return segments
}

// DeclaredFieldNames returns every field name the spec accepts in its body.
func (spec *LayoutSpec) DeclaredFieldNames() map[string]bool {
	names := map[string]bool{}
	if spec == nil {
		return names
	}
	for _, field := range allFieldSpecs(spec.Fields) {
		names[field.Name] = true
	}
	return names
}
