package localrender

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/layoutcatalog"
)

// Options configures one render.
type Options struct {
	Theme          string
	FontSize       string // small / medium / large
	BackgroundType string // default / grid / none
}

type renderer struct {
	theme    Theme
	fontSize string
	catalog  *layoutcatalog.Catalog
}

// ThemeError reports an unknown theme name.
type ThemeError struct{ Name string }

func (e *ThemeError) Error() string {
	return fmt.Sprintf("local renderer has no theme %q", e.Name)
}

// Render converts Markdown (including :::layout modules) into WeChat HTML.
//
// Markdown image references keep their original source in the output src so
// the upload flow can find and replace them, exactly like the remote API.
func Render(markdown string, opts Options) (string, error) {
	name := strings.TrimSpace(opts.Theme)
	if name == "" {
		name = "default"
	}
	theme, ok := ResolveTheme(name)
	if !ok {
		return "", &ThemeError{Name: name}
	}
	catalog, err := layoutcatalog.DefaultCatalog()
	if err != nil {
		return "", fmt.Errorf("load layout catalog: %w", err)
	}
	r := &renderer{theme: theme, fontSize: fontSizeFor(opts.FontSize), catalog: catalog}

	tokenized, images := tokenizeImages(markdown)

	var body strings.Builder
	for _, segment := range catalog.Segments(tokenized) {
		if segment.Block != nil {
			body.WriteString(r.renderBlock(segment.Block))
			continue
		}
		body.WriteString(r.markdownToHTML(segment.Markdown))
	}

	out := r.container(body.String(), opts.BackgroundType)
	return restoreImages(out, images), nil
}

func fontSizeFor(size string) string {
	switch strings.TrimSpace(size) {
	case "small":
		return "14px"
	case "large":
		return "17px"
	default:
		return "15px"
	}
}

func (r *renderer) container(inner, background string) string {
	t := r.theme
	bg := t.Bg
	extra := []string{}
	if strings.TrimSpace(background) == "grid" {
		if bg == "" {
			bg = "#ffffff"
		}
		extra = append(extra,
			"background-image", "linear-gradient(90deg,rgba(0,0,0,0.04) 1px,transparent 1px),linear-gradient(rgba(0,0,0,0.04) 1px,transparent 1px)",
			"background-size", "20px 20px")
	}
	padding := "8px 4px"
	if bg != "" {
		padding = "24px 16px"
	}
	style := css(append([]string{
		"font-family", t.fontFamily(),
		"font-size", r.fontSize,
		"color", t.Text,
		"line-height", "1.75",
		"letter-spacing", "0.5px",
		"word-wrap", "break-word",
		"padding", padding,
		"background-color", bg,
	}, extra...)...)
	return `<section style="` + attr(style) + `">` + inner + `</section>`
}

// imagePattern mirrors converter.ParseMarkdownImages so the n-th token here is
// the n-th ImageRef the publish flow will upload.
var imagePattern = regexp.MustCompile(`!\[[^\]]*\]\((__generate:[^)]+__|<[^>]+>|[^)\s]+)(?:\s+(?:"[^"]*"|'[^']*'|\([^)]*\)))?\)`)

const imageTokenPrefix = "md2wx-local-image-"

// tokenizeImages swaps every image destination for an ASCII token, so Markdown
// URL escaping cannot change the src the upload flow matches on.
func tokenizeImages(markdown string) (string, []string) {
	var originals []string
	out := imagePattern.ReplaceAllStringFunc(markdown, func(match string) string {
		sub := imagePattern.FindStringSubmatchIndex(match)
		if len(sub) < 4 {
			return match
		}
		ref := strings.TrimSpace(match[sub[2]:sub[3]])
		if strings.HasPrefix(ref, "<") && strings.HasSuffix(ref, ">") {
			ref = strings.TrimSpace(ref[1 : len(ref)-1])
		}
		if strings.HasPrefix(ref, "__generate:") && strings.HasSuffix(ref, "__") {
			ref = strings.TrimSuffix(strings.TrimPrefix(ref, "__generate:"), "__")
		}
		token := fmt.Sprintf("%s%d", imageTokenPrefix, len(originals))
		originals = append(originals, ref)
		return match[:sub[2]] + token + match[sub[3]:]
	})
	return out, originals
}

func restoreImages(out string, originals []string) string {
	// Replace higher indexes first so token-1 never clobbers token-10.
	for i := len(originals) - 1; i >= 0; i-- {
		out = strings.ReplaceAll(out, fmt.Sprintf("%s%d", imageTokenPrefix, i), html.EscapeString(originals[i]))
	}
	return out
}

func attr(value string) string {
	return html.EscapeString(value)
}

func esc(value string) string {
	return html.EscapeString(value)
}
