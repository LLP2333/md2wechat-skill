package localrender

import "strings"

var symbolGlyphs = map[string]string{
	"spark-solid":        "✦",
	"spark-outline":      "✧",
	"diamond-solid":      "◆",
	"diamond-outline":    "◇",
	"reference-mark":     "※",
	"asterism":           "⁂",
	"double-circle":      "◎",
	"circle":             "○",
	"square-solid":       "■",
	"square-outline":     "□",
	"star":               "★",
	"infinity":           "∞",
	"mountain":           "△",
	"concentric-circles": "◉",
	"nested-diamonds":    "◈",
	"four-petals":        "✤",
	"lens":               "⦿",
	"orbits":             "⊚",
	"archway":            "∩",
	"rounded-seal":       "❖",
	"four-point-star":    "✶",
	"honeycomb":          "⬡",
	"mirrored-waves":     "≋",
	"open-book":          "❦",
}

func symbolGlyph(name string) string {
	return symbolGlyphs[strings.TrimSpace(name)]
}

func section(style, inner string) string {
	return `<section style="` + attr(style) + `">` + inner + `</section>`
}

func span(style, inner string) string {
	return `<span style="` + attr(style) + `">` + inner + `</span>`
}

// card is the standard module container.
func (r *renderer) card(inner string, accent bool) string {
	t := r.theme
	border := "1px solid " + t.Border
	bg := t.cardBg()
	if accent {
		border = "1px solid " + t.Primary
		bg = t.Soft
	}
	return section(css("margin", "24px 0", "padding", "18px 18px", "border", border, "border-radius", "12px", "background", bg, "box-sizing", "border-box"), inner)
}

// softCard is a tinted container with no border.
func (r *renderer) softCard(inner string) string {
	return section(css("margin", "24px 0", "padding", "18px 18px", "border-radius", "12px", "background", r.theme.Soft, "box-sizing", "border-box"), inner)
}

func (r *renderer) caption(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	t := r.theme
	return section(css("margin", "0 0 12px", "font-size", "15px", "font-weight", "bold", "color", t.Primary, "letter-spacing", "1px"),
		span(css("display", "inline-block", "width", "4px", "height", "14px", "background", t.Primary, "border-radius", "2px", "margin-right", "8px", "vertical-align", "-1px"), "")+r.inline(text))
}

func (r *renderer) eyebrow(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	return section(css("margin", "0 0 8px", "font-size", "12px", "letter-spacing", "2px", "color", r.theme.Primary, "font-weight", "bold"), r.inline(text))
}

func (r *renderer) title(text, size string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	return section(css("margin", "0 0 8px", "font-size", size, "font-weight", "bold", "line-height", "1.45", "color", r.theme.Text), r.inline(text))
}

func (r *renderer) body(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	return section(css("margin", "6px 0 0", "font-size", "14px", "line-height", "1.75", "color", r.theme.Text), r.inline(text))
}

func (r *renderer) note(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	return section(css("margin", "10px 0 0", "font-size", "12px", "line-height", "1.6", "color", r.theme.Muted), r.inline(text))
}

func (r *renderer) pill(text string, filled bool) string {
	t := r.theme
	if filled {
		return span(css("display", "inline-block", "margin", "4px 6px 0 0", "padding", "2px 10px", "border-radius", "999px", "font-size", "12px", "background", t.Primary, "color", t.onPrimary()), r.inline(text))
	}
	return span(css("display", "inline-block", "margin", "4px 6px 0 0", "padding", "2px 10px", "border-radius", "999px", "font-size", "12px", "background", t.Soft, "color", t.Primary), r.inline(text))
}

func (r *renderer) pills(items []string, filled bool) string {
	if len(items) == 0 {
		return ""
	}
	var b strings.Builder
	for _, item := range items {
		b.WriteString(r.pill(item, filled))
	}
	return section(css("margin", "10px 0 0", "line-height", "1.6"), b.String())
}

func (r *renderer) badge(text string) string {
	t := r.theme
	return span(css("display", "inline-block", "min-width", "28px", "height", "28px", "line-height", "28px", "text-align", "center", "border-radius", "14px", "background", t.Primary, "color", t.onPrimary(), "font-size", "13px", "font-weight", "bold", "padding", "0 6px", "box-sizing", "border-box"), esc(text))
}

func (r *renderer) image(src, alt, radius string) string {
	if strings.TrimSpace(src) == "" {
		return ""
	}
	if radius == "" {
		radius = "8px"
	}
	return `<img src="` + attr(src) + `" alt="` + attr(alt) + `" style="` + attr(css("max-width", "100%", "width", "100%", "height", "auto", "display", "block", "margin", "0 auto", "border-radius", radius)) + `" />`
}

// row lays children side by side; WeChat keeps inline flex styles.
func row(children ...string) string {
	var b strings.Builder
	for i, child := range children {
		margin := "0"
		if i > 0 {
			margin = "0 0 0 10px"
		}
		b.WriteString(section(css("flex", "1", "min-width", "0", "margin", margin, "box-sizing", "border-box"), child))
	}
	return section(css("display", "flex", "align-items", "stretch", "margin", "0"), b.String())
}

// grid wraps items into rows of n columns.
func grid(items []string, columns int) string {
	if columns < 1 {
		columns = 1
	}
	var b strings.Builder
	for start := 0; start < len(items); start += columns {
		end := start + columns
		cells := make([]string, 0, columns)
		cells = append(cells, items[start:min(end, len(items))]...)
		for len(cells) < columns {
			cells = append(cells, "")
		}
		b.WriteString(section(css("margin", "0 0 10px"), row(cells...)))
	}
	return b.String()
}

func (r *renderer) toneAccent(tone string) bool {
	switch strings.TrimSpace(tone) {
	case "accent", "brand", "primary", "highlight":
		return true
	}
	return false
}

func (r *renderer) symbol(name, size string) string {
	glyph := symbolGlyph(name)
	if glyph == "" {
		return ""
	}
	return section(css("margin", "0 0 10px", "font-size", size, "line-height", "1", "color", r.theme.Primary, "text-align", "center"), glyph)
}
