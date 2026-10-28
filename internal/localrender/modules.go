package localrender

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/layoutcatalog"
)

type moduleRenderer func(r *renderer, b *layoutcatalog.Block) string

var moduleRenderers map[string]moduleRenderer

func init() {
	moduleRenderers = map[string]moduleRenderer{
		// opening
		"hero":          (*renderer).renderHero,
		"toc":           (*renderer).renderTOC,
		"cards":         (*renderer).renderCards,
		"part":          (*renderer).renderPart,
		"label-title":   (*renderer).renderLabelTitle,
		"section-title": (*renderer).renderSectionTitle,
		"epilogue":      (*renderer).renderEpilogue,
		"closing":       (*renderer).renderEpilogue,
		"cover-reveal":  (*renderer).renderCoverReveal,
		// judgment
		"verdict":      (*renderer).renderVerdict,
		"audience-fit": (*renderer).renderAudienceFit,
		"myth-fact":    (*renderer).renderMythFact,
		"manifesto":    (*renderer).renderManifesto,
		"bridge":       (*renderer).renderBridge,
		// infographic
		"metrics":     (*renderer).renderMetrics,
		"steps":       (*renderer).renderSteps,
		"timeline":    (*renderer).renderTimeline,
		"compare":     (*renderer).renderCompare,
		"infographic": (*renderer).renderInfographic,
		// evidence
		"quote":            (*renderer).renderQuote,
		"gallery":          (*renderer).renderGallery,
		"gallery-grid":     (*renderer).renderGalleryGrid,
		"gallery-story":    (*renderer).renderGalleryStory,
		"image-phone-shot": (*renderer).renderGalleryGrid,
		"figure-caption":   (*renderer).renderFigureCaption,
		"image-compare":    (*renderer).renderImageCompare,
		"image-annotate":   (*renderer).renderImageAnnotate,
		"image-text":       (*renderer).renderImageText,
		"image-steps":      (*renderer).renderImageSteps,
		"longimage":        (*renderer).renderGallery,
		// brand
		"author-card": (*renderer).renderAuthorCard,
		"subscribe":   (*renderer).renderSubscribe,
		"series":      (*renderer).renderSeries,
		"people":      (*renderer).renderPeople,
		// conversion
		"cta":       (*renderer).renderCTA,
		"summary":   (*renderer).renderSummary,
		"faq":       (*renderer).renderQA,
		"question":  (*renderer).renderQA,
		"checklist": (*renderer).renderChecklist,
		"notice":    (*renderer).renderNotice,
		"logos":     (*renderer).renderLogos,
		"cases":     (*renderer).renderCases,
		"pricing":   (*renderer).renderPricing,
		"specs":     (*renderer).renderSpecs,
		"toolbox":   (*renderer).renderToolbox,
		// sprint4
		"callout":          (*renderer).renderCallout,
		"definition":       (*renderer).renderDefinition,
		"quote-card":       (*renderer).renderQuoteCard,
		"tweet":            (*renderer).renderTweet,
		"stat-row":         (*renderer).renderStatRow,
		"resource-list":    (*renderer).renderResourceList,
		"comparison-table": (*renderer).renderComparisonTable,
		"changelog":        (*renderer).renderChangelog,
		// free layout
		"matrix":        (*renderer).renderMatrix,
		"dialogue-pair": (*renderer).renderDialogue,
		"dialogue":      (*renderer).renderDialogue,
		"flow":          (*renderer).renderFlow,
		"split":         (*renderer).renderSplit,
		// interactive modules render as complete static content
		"expand":            (*renderer).renderExpand,
		"svg-reveal":        (*renderer).renderReveal,
		"svg-swipe-gallery": (*renderer).renderGalleryGrid,
	}
}

// HasModule reports whether the local renderer draws the named module.
func HasModule(name string) bool {
	_, ok := moduleRenderers[name]
	return ok
}

func (r *renderer) renderBlock(b *layoutcatalog.Block) string {
	if render, ok := moduleRenderers[b.Name]; ok {
		if out := render(r, b); strings.TrimSpace(out) != "" {
			return out
		}
	}
	return r.renderFallback(b)
}

// renderFallback keeps unknown modules readable instead of leaking ::: syntax.
func (r *renderer) renderFallback(b *layoutcatalog.Block) string {
	fields, rest := parseFields(b.Body, declaredFields(b))
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, f := range fields {
		inner.WriteString(r.body("<strong>" + esc(f.key) + "</strong>：" + r.inline(f.value)))
	}
	inner.WriteString(r.markdownToHTML(strings.Join(rest, "\n")))
	return r.card(inner.String(), false)
}

func (r *renderer) fields(b *layoutcatalog.Block) fieldSet {
	fields, _ := parseFields(b.Body, declaredFields(b))
	return fields
}

// ---------- opening ----------

func (r *renderer) renderHero(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	variant := f.get("variant")
	var inner strings.Builder

	switch variant {
	case "masthead", "journal", "seal", "orbit":
		inner.WriteString(r.symbol(f.get("symbol"), "34px"))
		if brand := f.get("brand", "eyebrow"); brand != "" {
			inner.WriteString(section(css("font-size", "12px", "letter-spacing", "4px", "color", t.Primary, "text-align", "center", "margin", "0 0 10px"), esc(brand)))
		}
		inner.WriteString(section(css("font-size", "24px", "font-weight", "bold", "line-height", "1.4", "color", t.Text, "text-align", "center", "padding", "12px 0", "border-top", "2px solid "+t.Primary, "border-bottom", "2px solid "+t.Primary), r.inline(f.get("title"))))
		if sub := f.get("subtitle", "kicker"); sub != "" {
			inner.WriteString(section(css("font-size", "14px", "color", t.Muted, "text-align", "center", "margin", "12px 0 0"), r.inline(sub)))
		}
		if meta := f.get("meta"); meta != "" {
			inner.WriteString(section(css("font-size", "12px", "color", t.Muted, "text-align", "center", "margin", "8px 0 0"), esc(meta)))
		}
		return section(css("margin", "8px 0 32px", "padding", "24px 16px", "background", t.Soft, "border-radius", "12px"), inner.String())
	}

	if img := f.get("image"); img != "" {
		inner.WriteString(section(css("margin", "0 0 16px"), r.image(img, f.get("title"), "10px")))
	}
	top := []string{}
	if e := f.get("eyebrow"); e != "" {
		top = append(top, esc(e))
	}
	if m := f.get("meta"); m != "" {
		top = append(top, esc(m))
	}
	if len(top) > 0 {
		inner.WriteString(section(css("font-size", "12px", "letter-spacing", "2px", "color", t.Primary, "font-weight", "bold", "margin", "0 0 10px"), strings.Join(top, " · ")))
	}
	if variant == "story" {
		if k := f.get("kicker"); k != "" {
			inner.WriteString(section(css("font-size", "15px", "font-style", "italic", "color", t.Muted, "margin", "0 0 10px"), r.inline(k)))
		}
	} else if k := f.get("kicker"); k != "" {
		inner.WriteString(section(css("font-size", "14px", "color", t.Primary, "margin", "0 0 6px"), r.inline(k)))
	}
	inner.WriteString(section(css("font-size", "26px", "font-weight", "bold", "line-height", "1.35", "color", t.Text, "margin", "0 0 12px"), r.inline(f.get("title"))))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "15px", "line-height", "1.7", "color", t.Muted, "margin", "0 0 8px"), r.inline(sub)))
	}
	if points := f.list("points"); len(points) > 0 {
		var list strings.Builder
		for i, point := range points {
			list.WriteString(section(css("margin", "8px 0 0", "font-size", "14px", "color", t.Text), r.badge(fmt.Sprintf("%02d", i+1))+span(css("margin-left", "10px"), r.inline(point))))
		}
		inner.WriteString(section(css("margin", "12px 0 0"), list.String()))
	}
	inner.WriteString(r.pills(f.list("tags"), false))
	if brand := f.get("brand"); brand != "" {
		inner.WriteString(section(css("margin", "14px 0 0", "font-size", "12px", "color", t.Muted, "letter-spacing", "1px"), "— "+esc(brand)))
	}
	return section(css("margin", "8px 0 32px", "padding", "22px 20px", "border-left", "5px solid "+t.Primary, "background", t.Soft, "border-radius", "0 12px 12px 0"), inner.String())
}

func (r *renderer) renderTOC(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(firstNonEmpty(b.Caption, "目录")))
	for i, cells := range rows {
		index := firstNonEmpty(cell(cells, 0), fmt.Sprintf("%02d", i+1))
		text := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(cell(cells, 1)))
		if d := cell(cells, 2); d != "" {
			text += section(css("font-size", "13px", "color", t.Muted, "margin", "2px 0 0"), r.inline(d))
		}
		border := "1px dashed " + t.Border
		if i == len(rows)-1 {
			border = "0"
		}
		inner.WriteString(section(css("display", "flex", "align-items", "flex-start", "padding", "10px 0", "border-bottom", border),
			section(css("flex", "0 0 auto", "width", "44px", "font-size", "18px", "font-weight", "bold", "color", t.Primary, "font-family", monoFont), esc(index))+
				section(css("flex", "1", "min-width", "0"), text)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderCards(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		accent := r.toneAccent(cell(cells, 3))
		content := ""
		if l := cell(cells, 0); l != "" {
			content += section(css("font-size", "11px", "letter-spacing", "2px", "color", t.Primary, "font-weight", "bold"), esc(l))
		}
		content += section(css("font-size", "16px", "font-weight", "bold", "color", t.Text, "margin", "4px 0 0"), r.inline(cell(cells, 1)))
		content += r.body(cell(cells, 2))
		bg, border := t.cardBg(), "1px solid "+t.Border
		if accent {
			bg, border = t.Soft, "1px solid "+t.Primary
		}
		inner.WriteString(section(css("margin", "0 0 10px", "padding", "14px 16px", "border-radius", "10px", "background", bg, "border", border), content))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderPart(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	label := firstNonEmpty(f.get("label"), "PART")
	inner.WriteString(section(css("font-size", "12px", "letter-spacing", "4px", "color", t.Primary, "font-weight", "bold"), esc(label)+" "+esc(f.get("index"))))
	if idx := f.get("index"); idx != "" {
		inner.WriteString(section(css("font-size", "44px", "font-weight", "bold", "line-height", "1.1", "color", t.Primary, "margin", "6px 0", "font-family", monoFont), esc(idx)))
	}
	inner.WriteString(section(css("font-size", "20px", "font-weight", "bold", "color", t.Text, "line-height", "1.45"), r.inline(f.get("title"))))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "margin", "8px 0 0"), r.inline(sub)))
	}
	return section(css("margin", "40px 0 24px", "padding", "0 0 16px", "border-bottom", "2px solid "+t.Primary), inner.String())
}

func (r *renderer) renderLabelTitle(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	label := span(css("display", "inline-block", "padding", "2px 10px", "background", t.Primary, "color", t.onPrimary(), "font-size", "12px", "border-radius", "4px", "margin-right", "10px", "vertical-align", "2px"), esc(f.get("label")))
	return section(css("margin", "32px 0 16px", "font-size", "18px", "font-weight", "bold", "line-height", "1.5", "color", t.Text), label+r.inline(f.get("title")))
}

func (r *renderer) renderSectionTitle(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	title := r.inline(f.get("title"))
	glyph := symbolGlyph(f.get("symbol"))
	sub := ""
	if s := f.get("subtitle"); s != "" {
		sub = section(css("font-size", "13px", "font-weight", "normal", "color", t.Muted, "margin", "6px 0 0"), r.inline(s))
	}
	base := []string{"margin", "36px 0 18px", "font-size", "19px", "font-weight", "bold", "line-height", "1.45", "color", t.Text}

	switch f.get("variant") {
	case "divider":
		line := span(css("display", "inline-block", "width", "36px", "height", "1px", "background", t.Primary, "vertical-align", "middle", "margin", "0 10px"), "")
		g := ""
		if glyph != "" {
			g = span(css("color", t.Primary, "margin-right", "6px"), glyph)
		}
		return section(css(append(base, "text-align", "center")...), line+g+title+line+sub)
	case "numbered":
		idx := firstNonEmpty(f.get("index"), "01")
		return section(css(base...), span(css("font-size", "28px", "color", t.Primary, "font-family", monoFont, "margin-right", "10px", "vertical-align", "-3px"), esc(idx))+title+sub)
	case "frame":
		return section(css(append(base, "text-align", "center", "padding", "14px 12px", "border", "1px solid "+t.Primary, "border-radius", "8px")...), title+sub)
	case "focus":
		g := ""
		if glyph != "" {
			g = section(css("color", t.Primary, "font-size", "22px", "margin", "0 0 6px"), glyph)
		}
		return section(css(append(base, "text-align", "center")...), g+title+sub)
	case "vertical":
		g := ""
		if glyph != "" {
			g = span(css("color", t.Primary, "margin-right", "8px"), glyph)
		}
		return section(css(append(base, "border-left", "4px solid "+t.Primary, "padding", "2px 0 2px 12px")...), g+title+sub)
	default: // marker
		g := glyph
		if g == "" {
			g = "▍"
		}
		return section(css(base...), span(css("color", t.Primary, "margin-right", "8px"), g)+title+sub)
	}
}

func (r *renderer) renderEpilogue(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.symbol(firstNonEmpty(f.get("symbol"), "asterism"), "22px"))
	inner.WriteString(section(css("font-size", "18px", "font-weight", "bold", "line-height", "1.6", "color", t.Text, "text-align", "center"), r.inline(f.get("title"))))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "text-align", "center", "margin", "10px 0 0"), r.inline(sub)))
	}
	return section(css("margin", "40px 0", "padding", "24px 12px", "border-top", "1px solid "+t.Border, "border-bottom", "1px solid "+t.Border), inner.String())
}

func (r *renderer) renderCoverReveal(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	if img := f.get("image"); img != "" {
		inner.WriteString(section(css("margin", "0 0 16px"), r.image(img, f.get("alt"), "10px")))
	} else {
		inner.WriteString(r.symbol(f.get("symbol"), "40px"))
	}
	if l := f.get("label"); l != "" {
		inner.WriteString(section(css("font-size", "12px", "letter-spacing", "3px", "color", t.Primary, "text-align", "center", "margin", "0 0 8px"), esc(l)))
	}
	inner.WriteString(section(css("font-size", "24px", "font-weight", "bold", "color", t.Text, "text-align", "center", "line-height", "1.4"), r.inline(f.get("title"))))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "14px", "color", t.Muted, "text-align", "center", "margin", "10px 0 0"), r.inline(sub)))
	}
	return section(css("margin", "8px 0 32px", "padding", "36px 18px", "background", t.Soft, "border-radius", "14px"), inner.String())
}

// ---------- judgment ----------

func (r *renderer) renderVerdict(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(section(css("font-size", "12px", "letter-spacing", "2px", "color", t.onPrimary(), "font-weight", "bold", "opacity", "0.85", "margin", "0 0 8px"), esc(firstNonEmpty(f.get("eyebrow"), "结论"))+metaSuffix(f.get("meta"))))
	inner.WriteString(section(css("font-size", "19px", "font-weight", "bold", "line-height", "1.5", "color", t.onPrimary()), esc(f.get("title"))))
	if body := f.get("body"); body != "" {
		inner.WriteString(section(css("font-size", "14px", "line-height", "1.75", "color", t.onPrimary(), "opacity", "0.92", "margin", "10px 0 0"), esc(body)))
	}
	if note := f.get("note"); note != "" {
		inner.WriteString(section(css("font-size", "12px", "color", t.onPrimary(), "opacity", "0.75", "margin", "12px 0 0"), esc(note)))
	}
	return section(css("margin", "28px 0", "padding", "20px 20px", "background", t.Primary, "border-radius", "12px", "box-shadow", "0 6px 16px rgba(0,0,0,0.10)"), inner.String())
}

func metaSuffix(meta string) string {
	if meta == "" {
		return ""
	}
	return " · " + esc(meta)
}

func (r *renderer) checkList(items []string, mark, color string) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString(section(css("margin", "8px 0 0", "font-size", "14px", "line-height", "1.6", "color", r.theme.Text),
			span(css("color", color, "font-weight", "bold", "margin-right", "8px"), mark)+r.inline(item)))
	}
	return b.String()
}

func (r *renderer) renderAudienceFit(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.title(firstNonEmpty(f.get("title"), "这篇适合谁"), "17px"))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "margin", "0 0 12px"), r.inline(sub)))
	}
	fit := section(css("padding", "12px 14px", "background", t.tint(toneSuccess), "border-radius", "8px", "height", "100%", "box-sizing", "border-box"),
		section(css("font-size", "13px", "font-weight", "bold", "color", "#1f9d55"), "适合")+r.checkList(f.list("fit"), "✓", "#1f9d55"))
	blocks := []string{fit}
	if avoid := f.list("avoid"); len(avoid) > 0 {
		blocks = append(blocks, section(css("padding", "12px 14px", "background", t.tint(toneDanger), "border-radius", "8px", "height", "100%", "box-sizing", "border-box"),
			section(css("font-size", "13px", "font-weight", "bold", "color", "#d64545"), "不适合")+r.checkList(avoid, "✕", "#d64545")))
	}
	inner.WriteString(row(blocks...))
	inner.WriteString(r.note(f.get("note")))
	return r.card(inner.String(), false)
}

func (r *renderer) renderMythFact(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(firstNonEmpty(b.Caption, "误区与事实")))
	for _, cells := range rows {
		item := section(css("font-size", "14px", "color", t.Muted), span(css("color", "#d64545", "font-weight", "bold", "margin-right", "6px"), "✕ 误区")+span(css("text-decoration", "line-through"), r.inline(cell(cells, 0))))
		item += section(css("font-size", "15px", "color", t.Text, "font-weight", "bold", "margin", "6px 0 0"), span(css("color", "#1f9d55", "margin-right", "6px"), "✓ 事实")+r.inline(cell(cells, 1)))
		item += r.note(cell(cells, 2))
		inner.WriteString(section(css("padding", "12px 0", "border-bottom", "1px dashed "+t.Border), item))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderManifesto(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("label")))
	inner.WriteString(r.title(f.get("title"), "19px"))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted), r.inline(sub)))
	}
	inner.WriteString(r.body(f.get("body")))
	if believe := f.list("believe"); len(believe) > 0 {
		inner.WriteString(section(css("margin", "14px 0 0", "font-size", "13px", "font-weight", "bold", "color", t.Primary), "我相信"))
		inner.WriteString(r.checkList(believe, "●", t.Primary))
	}
	if against := f.list("against"); len(against) > 0 {
		inner.WriteString(section(css("margin", "14px 0 0", "font-size", "13px", "font-weight", "bold", "color", t.Muted), "我反对"))
		inner.WriteString(r.checkList(against, "○", t.Muted))
	}
	inner.WriteString(r.note(f.get("note")))
	return r.card(inner.String(), true)
}

func (r *renderer) renderBridge(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("label")))
	inner.WriteString(r.title(f.get("title"), "17px"))
	inner.WriteString(r.body(f.get("body")))
	if next := f.get("next"); next != "" {
		inner.WriteString(section(css("margin", "12px 0 0", "font-size", "14px", "font-weight", "bold", "color", t.Primary), "↓ "+r.inline(next)))
	}
	return section(css("margin", "28px 0", "padding", "16px 18px", "border-top", "1px solid "+t.Border, "border-bottom", "1px solid "+t.Border), inner.String())
}

// ---------- infographic ----------

func (r *renderer) renderMetrics(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var items []string
	for _, cells := range rows {
		accent := r.toneAccent(cell(cells, 3))
		bg, valueColor := t.Soft, t.Primary
		if accent {
			bg = t.Primary
			valueColor = t.onPrimary()
		}
		labelColor, noteColor := t.Muted, t.Muted
		if accent {
			labelColor, noteColor = t.onPrimary(), t.onPrimary()
		}
		item := section(css("font-size", "12px", "color", labelColor), esc(cell(cells, 0)))
		item += section(css("font-size", "26px", "font-weight", "bold", "line-height", "1.3", "color", valueColor, "margin", "4px 0"), esc(cell(cells, 1)))
		if n := cell(cells, 2); n != "" {
			item += section(css("font-size", "12px", "line-height", "1.5", "color", noteColor), r.inline(n))
		}
		items = append(items, section(css("padding", "14px", "border-radius", "10px", "background", bg, "height", "100%", "box-sizing", "border-box"), item))
	}
	columns := 2
	if len(items) == 3 {
		columns = 3
	}
	if len(items) == 1 {
		columns = 1
	}
	return section(css("margin", "24px 0"), r.caption(b.Caption)+grid(items, columns))
}

func (r *renderer) renderSteps(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for i, cells := range rows {
		index := firstNonEmpty(cell(cells, 0), fmt.Sprintf("%02d", i+1))
		content := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(cell(cells, 1)))
		if body := cell(cells, 2); body != "" {
			content += section(css("font-size", "13px", "color", t.Muted, "margin", "4px 0 0", "line-height", "1.7"), r.inline(body))
		}
		inner.WriteString(section(css("display", "flex", "align-items", "flex-start", "margin", "0 0 16px"),
			section(css("flex", "0 0 auto", "width", "36px", "text-align", "center"), r.badge(index))+
				section(css("flex", "1", "min-width", "0", "padding", "3px 0 0 10px"), content)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderTimeline(b *layoutcatalog.Block) string {
	rows, fields := parseRows(b)
	t := r.theme
	format := fields.get("format")
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		content := section(css("font-size", "12px", "color", t.Primary, "font-weight", "bold", "letter-spacing", "1px"), esc(cell(cells, 0)))
		bodyIdx, imgIdx := 2, 3
		if format == "meta-image-title-body" {
			imgIdx, bodyIdx = 1, 3
			if img := cell(cells, imgIdx); img != "" {
				content += section(css("margin", "6px 0"), r.image(img, cell(cells, 2), "6px"))
			}
			content += section(css("font-size", "15px", "font-weight", "bold", "color", t.Text, "margin", "2px 0 0"), r.inline(cell(cells, 2)))
		} else {
			content += section(css("font-size", "15px", "font-weight", "bold", "color", t.Text, "margin", "2px 0 0"), r.inline(cell(cells, 1)))
		}
		if format != "meta-dot-title" {
			if body := cell(cells, bodyIdx); body != "" {
				content += section(css("font-size", "13px", "color", t.Muted, "margin", "4px 0 0", "line-height", "1.7"), r.inline(body))
			}
			if format != "meta-image-title-body" {
				if img := cell(cells, imgIdx); img != "" {
					content += section(css("margin", "8px 0 0"), r.image(img, cell(cells, 1), "6px"))
				}
			}
		}
		dot := span(css("position", "absolute", "left", "-7px", "top", "4px", "width", "10px", "height", "10px", "border-radius", "5px", "background", t.Primary, "border", "2px solid "+t.cardBg()), "")
		inner.WriteString(section(css("position", "relative", "margin", "0 0 0 6px", "padding", "0 0 18px 18px", "border-left", "2px solid "+t.Soft), dot+content))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderCompare(b *layoutcatalog.Block) string {
	rows, fields := parseRows(b)
	t := r.theme
	headers := splitCells(b.Caption, "|")
	left := firstNonEmpty(fields.get("left-header"), cell(headers, 0), "Before")
	right := firstNonEmpty(fields.get("right-header"), cell(headers, 1), "After")
	head := row(
		section(css("padding", "8px 12px", "background", t.Border, "color", t.Text, "font-weight", "bold", "font-size", "14px", "border-radius", "8px 8px 0 0", "text-align", "center"), r.inline(left)),
		section(css("padding", "8px 12px", "background", t.Primary, "color", t.onPrimary(), "font-weight", "bold", "font-size", "14px", "border-radius", "8px 8px 0 0", "text-align", "center"), r.inline(right)),
	)
	var inner strings.Builder
	inner.WriteString(head)
	for _, cells := range rows {
		cellFor := func(title, body string, accent bool) string {
			bg := t.cardBg()
			if accent {
				bg = t.Soft
			}
			c := section(css("font-size", "14px", "font-weight", "bold", "color", t.Text), r.inline(title))
			if body != "" {
				c += section(css("font-size", "13px", "color", t.Muted, "margin", "4px 0 0", "line-height", "1.6"), r.inline(body))
			}
			return section(css("padding", "10px 12px", "background", bg, "border-bottom", "1px solid "+t.Border, "height", "100%", "box-sizing", "border-box"), c)
		}
		inner.WriteString(row(cellFor(cell(cells, 0), cell(cells, 1), false), cellFor(cell(cells, 2), cell(cells, 3), true)))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderInfographic(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	if accent := f.get("accent"); accent != "" {
		if color := accentColor(accent); color != "" {
			clone := *r
			clone.theme.Primary = color
			r = &clone
			t = r.theme
		}
	}
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("eyebrow")))
	inner.WriteString(r.title(f.get("title"), "20px"))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "14px", "color", t.Muted, "line-height", "1.7"), r.inline(sub)))
	}
	switch f.get("type") {
	case "pullquote":
		inner.WriteString(section(css("margin", "14px 0 0", "font-size", "18px", "line-height", "1.6", "color", t.Primary, "font-weight", "bold", "border-left", "3px solid "+t.Primary, "padding-left", "12px"), "“"+r.inline(f.get("quote"))+"”"))
	case "number":
		inner.WriteString(section(css("margin", "14px 0 0", "text-align", "center"),
			section(css("font-size", "44px", "font-weight", "bold", "color", t.Primary, "line-height", "1.2"), esc(f.get("value")))+
				section(css("font-size", "13px", "color", t.Muted), r.inline(f.get("label")))))
	case "contrast":
		inner.WriteString(section(css("margin", "14px 0 0"), row(
			section(css("padding", "12px", "background", t.Border, "border-radius", "8px", "text-align", "center", "font-size", "14px", "color", t.Muted, "height", "100%", "box-sizing", "border-box"), r.inline(f.get("left"))),
			section(css("padding", "12px", "background", t.Primary, "border-radius", "8px", "text-align", "center", "font-size", "14px", "color", t.onPrimary(), "font-weight", "bold", "height", "100%", "box-sizing", "border-box"), r.inline(f.get("right"))),
		)))
	case "formula":
		inner.WriteString(section(css("margin", "14px 0 0", "padding", "14px", "background", t.cardBg(), "border-radius", "8px", "text-align", "center", "font-size", "18px", "font-weight", "bold", "color", t.Primary, "letter-spacing", "1px"), r.inline(f.get("formula"))))
	case "tradeoff":
		for _, pair := range f.list("tradeoffs") {
			parts := strings.SplitN(pair, ":", 2)
			inner.WriteString(section(css("margin", "10px 0 0", "font-size", "14px", "color", t.Text), span(css("font-weight", "bold", "color", t.Primary), esc(cell(parts, 0)))+" ⇄ "+esc(cell(parts, 1))))
		}
	case "path", "evidence-chain", "micro-case":
		items := f.list("items")
		var chain strings.Builder
		for i, item := range items {
			chain.WriteString(section(css("margin", "10px 0 0", "font-size", "14px", "color", t.Text), r.badge(strconv.Itoa(i+1))+span(css("margin-left", "10px"), r.inline(item))))
			if i < len(items)-1 {
				chain.WriteString(section(css("margin", "2px 0 0 10px", "color", t.Primary, "font-size", "12px"), "↓"))
			}
		}
		inner.WriteString(chain.String())
	case "anatomy", "diagnosis":
		mark := "●"
		if f.get("type") == "diagnosis" {
			mark = "!"
		}
		inner.WriteString(r.checkList(f.list("items"), mark, t.Primary))
	default: // thesis
		if q := f.get("quote"); q != "" {
			inner.WriteString(r.body(q))
		}
		inner.WriteString(r.checkList(f.list("items"), "●", t.Primary))
	}
	inner.WriteString(r.note(f.get("note")))
	return section(css("margin", "28px 0", "padding", "20px 18px", "background", t.Soft, "border-radius", "12px", "border-top", "4px solid "+t.Primary), inner.String())
}

func accentColor(name string) string {
	named := map[string]string{
		"green": "#2f9e6e", "emerald": "#10b981", "blue": "#2f6fdb", "sky": "#1e9bd7",
		"orange": "#e2722e", "amber": "#d9a441", "red": "#d64545", "rose": "#e11d63",
	}
	if c, ok := named[name]; ok {
		return c
	}
	if strings.HasPrefix(name, "#") && (len(name) == 4 || len(name) == 7) {
		return name
	}
	return ""
}

// ---------- evidence ----------

func (r *renderer) renderQuote(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	text := f.get("quote", "text")
	source := f.get("source", "author")
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("eyebrow")))
	if img := f.get("image", "img", "cover"); img != "" {
		inner.WriteString(section(css("margin", "0 0 12px"), r.image(img, f.get("image_alt", "alt"), "8px")))
	}
	if f.get("variant") == "brand" {
		inner.WriteString(section(css("font-size", "20px", "font-weight", "bold", "line-height", "1.6", "color", t.onPrimary(), "text-align", "center"), "“"+r.inline(text)+"”"))
		if brand := f.get("brand", "source"); brand != "" {
			inner.WriteString(section(css("font-size", "12px", "letter-spacing", "3px", "color", t.onPrimary(), "opacity", "0.8", "text-align", "center", "margin", "12px 0 0"), esc(brand)))
		}
		return section(css("margin", "28px 0", "padding", "28px 20px", "background", t.Primary, "border-radius", "12px"), inner.String())
	}
	inner.WriteString(section(css("font-size", "36px", "line-height", "1", "color", t.Primary, "font-family", "Georgia,serif", "height", "24px"), "“"))
	inner.WriteString(section(css("font-size", "17px", "line-height", "1.75", "color", t.Text, "font-weight", "bold"), r.inline(text)))
	attribution := []string{}
	if source != "" {
		attribution = append(attribution, esc(source))
	}
	if role := f.get("role", "context", "source_type"); role != "" {
		attribution = append(attribution, esc(role))
	}
	if brand := f.get("brand"); brand != "" && f.get("variant") != "brand" {
		attribution = append(attribution, esc(brand))
	}
	if len(attribution) > 0 {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "margin", "12px 0 0", "text-align", "right"), "—— "+strings.Join(attribution, " · ")))
	}
	inner.WriteString(r.note(f.get("note")))
	if f.get("variant") == "proof" {
		return r.card(inner.String(), true)
	}
	return section(css("margin", "28px 0", "padding", "18px 20px", "background", t.Soft, "border-radius", "12px"), inner.String())
}

func (r *renderer) captionText(text string, index int, style string) string {
	if strings.TrimSpace(text) == "" || style == "none" {
		return ""
	}
	t := r.theme
	prefix := ""
	switch style {
	case "numbered":
		prefix = fmt.Sprintf("图 %d  ", index+1)
	case "label":
		prefix = "▍"
	}
	return section(css("margin", "6px 0 0", "font-size", "12px", "line-height", "1.5", "color", t.Muted, "text-align", "center"), esc(prefix)+r.inline(text))
}

func (r *renderer) renderGallery(b *layoutcatalog.Block) string {
	items := parseImageLines(b.Body)
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for i, item := range items {
		inner.WriteString(section(css("margin", "0 0 12px"), r.image(item.src, item.alt, "8px")+r.captionText(firstNonEmpty(cell(item.parts, 0), item.alt), i, b.Params["caption_style"])))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func imageRadius(shape string) string {
	switch shape {
	case "square", "original":
		return "0"
	case "phone":
		return "18px"
	case "poster":
		return "4px"
	default:
		return "8px"
	}
}

func (r *renderer) renderGalleryGrid(b *layoutcatalog.Block) string {
	items := parseImageLines(b.Body)
	t := r.theme
	columns, _ := strconv.Atoi(b.Params["columns"])
	if columns <= 0 {
		columns = 2
		if len(items) >= 3 && len(items)%3 == 0 {
			columns = 3
		}
	}
	if columns > 3 {
		columns = 3
	}
	radius := imageRadius(b.Params["image_shape"])
	if b.Name == "image-phone-shot" && b.Params["image_shape"] == "" {
		radius = "18px"
	}
	cardStyle := b.Params["variant"] == "card"
	var cells []string
	for i, item := range items {
		img := r.image(item.src, item.alt, radius)
		if b.Name == "image-phone-shot" {
			img = section(css("padding", "6px", "background", "#1f1f1f", "border-radius", "22px"), img)
		}
		content := img + r.captionText(cell(item.parts, 0), i, b.Params["caption_style"])
		if cardStyle {
			content = section(css("padding", "8px", "background", t.cardBg(), "border", "1px solid "+t.Border, "border-radius", "10px", "height", "100%", "box-sizing", "border-box"), content)
		}
		cells = append(cells, content)
	}
	heading := r.caption(firstNonEmpty(b.Params["title"], b.Caption))
	return section(css("margin", "24px 0"), heading+grid(cells, columns))
}

func (r *renderer) renderGalleryStory(b *layoutcatalog.Block) string {
	items := parseImageLines(b.Body)
	t := r.theme
	var inner strings.Builder
	for _, item := range items {
		content := r.image(item.src, item.alt, "8px")
		if title := cell(item.parts, 0); title != "" {
			content += section(css("margin", "10px 0 0", "font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(title))
		}
		if desc := cell(item.parts, 1); desc != "" {
			content += section(css("margin", "4px 0 0", "font-size", "13px", "color", t.Muted, "line-height", "1.7"), r.inline(desc))
		}
		style := css("margin", "0 0 18px")
		if b.Params["variant"] == "card" {
			style = css("margin", "0 0 14px", "padding", "10px", "background", t.cardBg(), "border", "1px solid "+t.Border, "border-radius", "10px")
		}
		inner.WriteString(section(style, content))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderFigureCaption(b *layoutcatalog.Block) string {
	fields, rest := parseFields(b.Body, declaredFields(b))
	t := r.theme
	items := parseImageLines(rest)
	var inner strings.Builder
	for _, item := range items {
		inner.WriteString(r.image(item.src, item.alt, "8px"))
	}
	caption := fields.get("caption")
	style := b.Params["caption_style"]
	if caption != "" && style != "none" {
		prefix := ""
		switch style {
		case "numbered":
			prefix = "图 1  "
		case "label":
			prefix = "▍"
		}
		inner.WriteString(section(css("margin", "8px 0 0", "font-size", "13px", "color", t.Muted, "text-align", "center"), esc(prefix)+r.inline(caption)))
	}
	if source := fields.get("source"); source != "" {
		inner.WriteString(section(css("margin", "2px 0 0", "font-size", "12px", "color", t.Muted, "text-align", "center", "opacity", "0.8"), "来源："+r.inline(source)))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderImageCompare(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	side := func(title, img, alt string, accent bool) string {
		color := t.Muted
		if accent {
			color = t.Primary
		}
		return section(css("font-size", "13px", "font-weight", "bold", "color", color, "text-align", "center", "margin", "0 0 6px"), r.inline(title)) + r.image(img, alt, "8px")
	}
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("eyebrow")))
	inner.WriteString(r.title(f.get("title"), "17px"))
	inner.WriteString(section(css("margin", "12px 0 0"), row(
		side(firstNonEmpty(f.get("left_title"), "Before"), f.get("left_image"), f.get("left_alt"), false),
		side(firstNonEmpty(f.get("right_title"), "After"), f.get("right_image"), f.get("right_alt"), true),
	)))
	inner.WriteString(r.note(f.get("note")))
	return r.card(inner.String(), false)
}

func (r *renderer) renderImageAnnotate(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("eyebrow")))
	inner.WriteString(r.title(f.get("title"), "17px"))
	inner.WriteString(section(css("margin", "12px 0"), r.image(f.get("image"), f.get("alt"), "8px")))
	descs := f.all("point-description")
	for i, point := range f.all("point") {
		parts := splitCells(point, "|")
		index, label, desc := cell(parts, 0), cell(parts, 1), cell(parts, 2)
		if len(parts) < 3 {
			index, label, desc = fmt.Sprintf("%02d", i+1), cell(parts, 0), cell(parts, 1)
		}
		if desc == "" && i < len(descs) {
			desc = descs[i]
		}
		content := section(css("font-size", "14px", "font-weight", "bold", "color", t.Text), r.inline(label))
		if desc != "" {
			content += section(css("font-size", "13px", "color", t.Muted, "margin", "2px 0 0", "line-height", "1.6"), r.inline(desc))
		}
		inner.WriteString(section(css("display", "flex", "align-items", "flex-start", "margin", "10px 0 0"),
			section(css("flex", "0 0 auto", "margin-right", "10px"), r.badge(index))+section(css("flex", "1", "min-width", "0"), content)))
	}
	inner.WriteString(r.note(f.get("note")))
	return r.card(inner.String(), false)
}

func (r *renderer) renderImageText(b *layoutcatalog.Block) string {
	f := r.fields(b)
	text := r.eyebrow(f.get("eyebrow")) + r.title(f.get("title"), "16px") + r.body(f.get("body")) + r.note(f.get("note"))
	img := r.image(f.get("image"), f.get("alt"), "8px")
	var inner string
	if f.get("layout") == "left" {
		inner = row(img, text)
	} else {
		inner = row(text, img)
	}
	return r.card(inner, false)
}

func (r *renderer) renderImageSteps(b *layoutcatalog.Block) string {
	t := r.theme
	type step struct{ title, img, alt, desc, note string }
	var steps []step
	declared := declaredFields(b)
	for _, raw := range b.Body {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if key, value, ok := parseFieldLine(line, declared); ok {
			if key == "step" || len(steps) == 0 {
				steps = append(steps, step{})
			}
			s := &steps[len(steps)-1]
			switch key {
			case "step":
				s.title = value
			case "desc":
				s.desc = value
			case "note":
				s.note = value
			}
			continue
		}
		if items := parseImageLines([]string{line}); len(items) > 0 {
			if len(steps) == 0 {
				steps = append(steps, step{})
			}
			steps[len(steps)-1].img, steps[len(steps)-1].alt = items[0].src, items[0].alt
		}
	}
	columns, _ := strconv.Atoi(b.Params["columns"])
	if columns <= 0 || columns > 3 {
		columns = 1
	}
	var cells []string
	for i, s := range steps {
		head := section(css("font-size", "14px", "font-weight", "bold", "color", t.Text, "margin", "0 0 8px"), r.badge(strconv.Itoa(i+1))+span(css("margin-left", "8px"), r.inline(s.title)))
		content := head + r.image(s.img, s.alt, "8px")
		if s.desc != "" {
			content += section(css("font-size", "13px", "color", t.Text, "margin", "8px 0 0", "line-height", "1.6"), r.inline(s.desc))
		}
		content += r.note(s.note)
		cells = append(cells, section(css("padding", "10px", "background", t.cardBg(), "border", "1px solid "+t.Border, "border-radius", "10px", "height", "100%", "box-sizing", "border-box"), content))
	}
	return section(css("margin", "24px 0"), grid(cells, columns))
}

// ---------- brand ----------

func (r *renderer) renderAuthorCard(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	avatar := ""
	if src := f.get("avatar"); src != "" {
		avatar = `<img src="` + attr(src) + `" alt="` + attr(f.get("avatar_alt")) + `" style="` + attr(css("width", "56px", "height", "56px", "border-radius", "28px", "display", "block")) + `" />`
	} else {
		initial := []rune(firstNonEmpty(f.get("name"), "A"))
		glyph := string(initial[0])
		if g := symbolGlyph(f.get("symbol")); g != "" {
			glyph = g
		}
		avatar = section(css("width", "56px", "height", "56px", "line-height", "56px", "border-radius", "28px", "background", t.Primary, "color", t.onPrimary(), "text-align", "center", "font-size", "22px", "font-weight", "bold"), esc(glyph))
	}
	info := section(css("font-size", "17px", "font-weight", "bold", "color", t.Text), esc(f.get("name")))
	if role := f.get("role"); role != "" {
		info += section(css("font-size", "12px", "color", t.Primary, "margin", "2px 0 0"), esc(role))
	}
	var inner strings.Builder
	inner.WriteString(section(css("display", "flex", "align-items", "center"),
		section(css("flex", "0 0 auto", "margin-right", "14px"), avatar)+section(css("flex", "1", "min-width", "0"), info)))
	inner.WriteString(r.body(f.get("bio")))
	inner.WriteString(r.pills(f.list("tags"), false))
	inner.WriteString(r.note(f.get("note")))
	if link := f.get("link"); link != "" {
		inner.WriteString(section(css("margin", "8px 0 0", "font-size", "12px", "color", t.Primary), "🔗 "+esc(link)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderSubscribe(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.eyebrow(f.get("label")))
	inner.WriteString(r.title(f.get("title"), "17px"))
	inner.WriteString(r.body(f.get("subtitle", "description")))
	if qr := f.get("qrcode"); qr != "" {
		inner.WriteString(section(css("margin", "14px auto 0", "width", "140px"), r.image(qr, "二维码", "6px")))
	}
	inner.WriteString(r.buttons(f.get("primary", "cta"), f.get("secondary"), ""))
	inner.WriteString(r.note(f.get("note")))
	return section(css("margin", "28px 0", "padding", "20px 18px", "background", t.Soft, "border-radius", "12px", "text-align", "center"), inner.String())
}

func (r *renderer) buttons(primary, secondary, tertiary string) string {
	t := r.theme
	var b strings.Builder
	if primary != "" {
		b.WriteString(span(css("display", "inline-block", "margin", "6px 4px 0", "padding", "8px 20px", "border-radius", "999px", "background", t.Primary, "color", t.onPrimary(), "font-size", "14px", "font-weight", "bold"), esc(primary)))
	}
	for _, label := range []string{secondary, tertiary} {
		if label != "" {
			b.WriteString(span(css("display", "inline-block", "margin", "6px 4px 0", "padding", "7px 18px", "border-radius", "999px", "border", "1px solid "+t.Primary, "color", t.Primary, "font-size", "14px"), esc(label)))
		}
	}
	if b.Len() == 0 {
		return ""
	}
	return section(css("margin", "14px 0 0", "text-align", "center"), b.String())
}

func (r *renderer) renderSeries(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	issue := f.get("issue", "index")
	if total := f.get("total"); total != "" && issue != "" {
		issue += " / " + total
	}
	var inner strings.Builder
	head := esc(firstNonEmpty(f.get("name"), f.get("label"), "系列"))
	if issue != "" {
		head += span(css("float", "right", "font-family", monoFont), "No."+esc(issue))
	}
	inner.WriteString(section(css("font-size", "12px", "letter-spacing", "2px", "color", t.Primary, "font-weight", "bold", "padding", "0 0 8px", "border-bottom", "1px solid "+t.Border, "margin", "0 0 10px"), head))
	inner.WriteString(r.title(f.get("title"), "17px"))
	inner.WriteString(r.body(f.get("desc", "description")))
	inner.WriteString(r.pills(f.list("tags"), false))
	if next := f.get("next"); next != "" {
		inner.WriteString(section(css("margin", "12px 0 0", "font-size", "13px", "color", t.Primary), "→ "+r.inline(next)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderPeople(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		name := cell(cells, 0)
		initial := string([]rune(firstNonEmpty(name, "?"))[0])
		avatar := section(css("width", "40px", "height", "40px", "line-height", "40px", "border-radius", "20px", "background", t.Primary, "color", t.onPrimary(), "text-align", "center", "font-weight", "bold"), esc(initial))
		info := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), esc(name)+span(css("font-size", "12px", "font-weight", "normal", "color", t.Primary, "margin-left", "8px"), esc(cell(cells, 1))))
		if body := cell(cells, 2); body != "" {
			info += section(css("font-size", "13px", "color", t.Muted, "margin", "2px 0 0", "line-height", "1.6"), r.inline(body))
		}
		bg := t.cardBg()
		if r.toneAccent(cell(cells, 3)) {
			bg = t.Soft
		}
		inner.WriteString(section(css("display", "flex", "align-items", "center", "padding", "12px", "margin", "0 0 8px", "background", bg, "border", "1px solid "+t.Border, "border-radius", "10px"),
			section(css("flex", "0 0 auto", "margin-right", "12px"), avatar)+section(css("flex", "1", "min-width", "0"), info)))
	}
	return section(css("margin", "24px 0"), inner.String())
}

// ---------- conversion ----------

func (r *renderer) renderCTA(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(section(css("font-size", "17px", "font-weight", "bold", "line-height", "1.6", "color", t.Text, "text-align", "center"), r.inline(f.get("title"))))
	if sub := f.get("subtitle"); sub != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "text-align", "center", "margin", "6px 0 0"), r.inline(sub)))
	}
	if points := f.list("points"); len(points) > 0 {
		var list strings.Builder
		for i, p := range points {
			list.WriteString(section(css("margin", "6px 0 0", "font-size", "14px", "color", t.Text), r.badge(strconv.Itoa(i+1))+span(css("margin-left", "8px"), r.inline(p))))
		}
		inner.WriteString(section(css("margin", "12px auto 0", "display", "table"), list.String()))
	}
	inner.WriteString(r.buttons(f.get("primary"), f.get("secondary"), f.get("tertiary")))
	if note := f.get("note"); note != "" {
		inner.WriteString(section(css("font-size", "12px", "color", t.Muted, "text-align", "center", "margin", "12px 0 0"), r.inline(note)))
	}
	return section(css("margin", "32px 0", "padding", "22px 18px", "border", "2px solid "+t.Primary, "border-radius", "14px", "background", t.cardBg()), inner.String())
}

func (r *renderer) renderSummary(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.eyebrow(firstNonEmpty(f.get("eyebrow"), "总结")))
	inner.WriteString(r.title(f.get("title"), "17px"))
	if h := f.get("highlight"); h != "" {
		inner.WriteString(section(css("font-size", "18px", "font-weight", "bold", "line-height", "1.6", "color", t.Primary, "margin", "4px 0 0"), r.inline(h)))
	}
	inner.WriteString(r.body(f.get("body")))
	switch f.get("variant") {
	case "decision":
		if fit := f.list("fit"); len(fit) > 0 {
			inner.WriteString(section(css("margin", "12px 0 0", "font-size", "13px", "font-weight", "bold", "color", "#1f9d55"), "适合"))
			inner.WriteString(r.checkList(fit, "✓", "#1f9d55"))
		}
		if notFit := f.list("not_fit"); len(notFit) > 0 {
			inner.WriteString(section(css("margin", "12px 0 0", "font-size", "13px", "font-weight", "bold", "color", "#d64545"), "不适合"))
			inner.WriteString(r.checkList(notFit, "✕", "#d64545"))
		}
		if rec := f.get("recommendation"); rec != "" {
			inner.WriteString(section(css("margin", "14px 0 0", "padding", "10px 12px", "background", t.cardBg(), "border-radius", "8px", "font-size", "14px", "font-weight", "bold", "color", t.Primary), "建议："+r.inline(rec)))
		}
	default:
		items := f.list("items")
		for i, item := range items {
			inner.WriteString(section(css("margin", "10px 0 0", "font-size", "14px", "color", t.Text), r.badge(strconv.Itoa(i+1))+span(css("margin-left", "10px"), r.inline(item))))
		}
	}
	inner.WriteString(r.note(f.get("note")))
	return section(css("margin", "28px 0", "padding", "18px 18px", "background", t.Soft, "border-radius", "12px", "border-left", "4px solid "+t.Primary), inner.String())
}

func (r *renderer) renderQA(b *layoutcatalog.Block) string {
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, item := range parseDialogue(b.Body) {
		speaker := strings.ToUpper(item.speaker)
		if speaker == "Q" {
			inner.WriteString(section(css("margin", "14px 0 0", "font-size", "15px", "font-weight", "bold", "color", t.Text, "line-height", "1.6"), span(css("color", t.Primary, "margin-right", "8px", "font-family", monoFont), "Q")+r.inline(item.text)))
			continue
		}
		inner.WriteString(section(css("margin", "6px 0 0", "padding", "0 0 12px", "font-size", "14px", "color", t.Muted, "line-height", "1.75", "border-bottom", "1px dashed "+t.Border), span(css("color", t.Accent, "margin-right", "8px", "font-weight", "bold", "font-family", monoFont), "A")+r.inline(item.text)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderChecklist(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		mark, color := "☐", t.Muted
		switch cell(cells, 0) {
		case "done", "ok", "yes", "x":
			mark, color = "☑", "#1f9d55"
		case "warn", "warning", "risk":
			mark, color = "⚠", "#d9822b"
		case "no", "fail", "blocked":
			mark, color = "☒", "#d64545"
		}
		content := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(cell(cells, 1)))
		if d := cell(cells, 2); d != "" {
			content += section(css("font-size", "13px", "color", t.Muted, "margin", "2px 0 0"), r.inline(d))
		}
		inner.WriteString(section(css("display", "flex", "align-items", "flex-start", "padding", "8px 0"),
			section(css("flex", "0 0 auto", "width", "28px", "font-size", "18px", "line-height", "1.3", "color", color), mark)+section(css("flex", "1", "min-width", "0"), content)))
	}
	return r.card(inner.String(), false)
}

func noticeColors(tone string, t Theme) (string, string) {
	switch tone {
	case "fit", "success", "ok":
		return "#1f9d55", t.tint(toneSuccess)
	case "risk", "warning", "warn":
		return "#d9822b", t.tint(toneWarning)
	case "avoid", "danger", "error":
		return "#d64545", t.tint(toneDanger)
	default:
		return t.Primary, t.Soft
	}
}

func (r *renderer) renderNotice(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		color, bg := noticeColors(cell(cells, 0), t)
		content := section(css("font-size", "14px", "color", t.Text, "line-height", "1.7"), span(css("font-weight", "bold", "color", color, "margin-right", "8px"), esc(cell(cells, 1)))+r.inline(cell(cells, 2)))
		if n := cell(cells, 3); n != "" {
			content += section(css("font-size", "12px", "color", t.Muted, "margin", "4px 0 0"), r.inline(n))
		}
		inner.WriteString(section(css("margin", "0 0 8px", "padding", "10px 12px", "background", bg, "border-left", "3px solid "+color, "border-radius", "0 8px 8px 0"), content))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderLogos(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var cells []string
	for _, c := range rows {
		content := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text, "letter-spacing", "1px"), esc(cell(c, 0)))
		if n := cell(c, 1); n != "" {
			content += section(css("font-size", "12px", "color", t.Muted, "margin", "2px 0 0"), esc(n))
		}
		cells = append(cells, section(css("padding", "12px 8px", "text-align", "center", "border", "1px solid "+t.Border, "border-radius", "8px", "background", t.cardBg(), "height", "100%", "box-sizing", "border-box"), content))
	}
	columns := 3
	if len(cells) <= 2 || len(cells) == 4 {
		columns = 2
	}
	return section(css("margin", "24px 0"), r.caption(b.Caption)+grid(cells, columns))
}

func (r *renderer) renderCases(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, cells := range rows {
		content := ""
		if img := cell(cells, 3); img != "" && !r.toneAccent(img) && img != "default" {
			content += section(css("margin", "0 0 10px"), r.image(img, cell(cells, 0), "8px"))
		}
		content += section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(cell(cells, 0)))
		if res := cell(cells, 1); res != "" {
			content += section(css("font-size", "20px", "font-weight", "bold", "color", t.Primary, "margin", "4px 0"), r.inline(res))
		}
		content += r.body(cell(cells, 2))
		accent := r.toneAccent(cell(cells, len(cells)-1))
		bg, border := t.cardBg(), "1px solid "+t.Border
		if accent {
			bg, border = t.Soft, "1px solid "+t.Primary
		}
		inner.WriteString(section(css("margin", "0 0 10px", "padding", "14px 16px", "background", bg, "border", border, "border-radius", "10px"), content))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderPricing(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var cells []string
	for _, c := range rows {
		accent := r.toneAccent(cell(c, 3))
		fg, sub, bg, border := t.Text, t.Muted, t.cardBg(), "1px solid "+t.Border
		priceColor := t.Primary
		if accent {
			fg, sub, bg, border, priceColor = t.onPrimary(), t.onPrimary(), t.Primary, "1px solid "+t.Primary, t.onPrimary()
		}
		content := section(css("font-size", "14px", "font-weight", "bold", "color", fg), esc(cell(c, 0)))
		content += section(css("font-size", "26px", "font-weight", "bold", "color", priceColor, "margin", "6px 0"), esc(cell(c, 1)))
		for _, feature := range splitList(cell(c, 2), "/") {
			content += section(css("font-size", "13px", "color", sub, "margin", "4px 0 0"), "✓ "+r.inline(feature))
		}
		cells = append(cells, section(css("padding", "16px 14px", "background", bg, "border", border, "border-radius", "12px", "height", "100%", "box-sizing", "border-box", "text-align", "center"), content))
	}
	columns := min(len(cells), 3)
	return section(css("margin", "24px 0"), r.caption(b.Caption)+grid(cells, columns))
}

func (r *renderer) renderSpecs(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var tbody strings.Builder
	for i, c := range rows {
		bg := t.cardBg()
		if i%2 == 1 {
			bg = t.Soft
		}
		value := r.inline(cell(c, 1))
		if n := cell(c, 2); n != "" {
			value += section(css("font-size", "12px", "color", t.Muted, "margin", "2px 0 0"), r.inline(n))
		}
		tbody.WriteString(`<tr style="` + attr(css("background", bg)) + `">` +
			`<td style="` + attr(css("padding", "10px 12px", "border", "1px solid "+t.Border, "font-weight", "bold", "color", t.Primary, "width", "32%", "font-size", "14px", "vertical-align", "top")) + `">` + esc(cell(c, 0)) + `</td>` +
			`<td style="` + attr(css("padding", "10px 12px", "border", "1px solid "+t.Border, "color", t.Text, "font-size", "14px", "line-height", "1.6")) + `">` + value + `</td></tr>`)
	}
	table := `<table style="` + attr(css("width", "100%", "border-collapse", "collapse")) + `"><tbody>` + tbody.String() + `</tbody></table>`
	return section(css("margin", "24px 0"), r.caption(b.Caption)+table)
}

func (r *renderer) renderToolbox(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, c := range rows {
		content := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), span(css("display", "inline-block", "font-size", "11px", "padding", "1px 8px", "border-radius", "4px", "background", t.Soft, "color", t.Primary, "margin-right", "8px", "vertical-align", "1px"), esc(cell(c, 0)))+r.inline(cell(c, 1)))
		if d := cell(c, 2); d != "" {
			content += section(css("font-size", "13px", "color", t.Muted, "margin", "4px 0 0"), r.inline(d))
		}
		inner.WriteString(section(css("padding", "10px 0", "border-bottom", "1px dashed "+t.Border), content))
	}
	return r.card(inner.String(), false)
}

// ---------- sprint4 ----------

func (r *renderer) renderCallout(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	kind := firstNonEmpty(f.get("type"), b.Params["variant"], "info")
	icons := map[string]string{"info": "ℹ️", "warning": "⚠️", "success": "✅", "danger": "⛔"}
	colors := map[string][2]string{
		"info":    {t.Primary, t.Soft},
		"warning": {"#d9822b", t.tint(toneWarning)},
		"success": {"#1f9d55", t.tint(toneSuccess)},
		"danger":  {"#d64545", t.tint(toneDanger)},
	}
	pair, ok := colors[kind]
	if !ok {
		pair = colors["info"]
	}
	_, rest := parseFields(b.Body, declaredFields(b))
	body := r.inline(f.get("body"))
	if len(rest) > 0 {
		body += r.markdownToHTML(strings.Join(rest, "\n"))
	}
	return section(css("margin", "24px 0", "padding", "12px 14px", "background", pair[1], "border-left", "4px solid "+pair[0], "border-radius", "0 8px 8px 0", "font-size", "14px", "line-height", "1.75", "color", t.Text),
		span(css("margin-right", "8px"), icons[kind])+body)
}

func (r *renderer) renderDefinition(b *layoutcatalog.Block) string {
	obj, _ := parseJSON(b.Body).(map[string]any)
	t := r.theme
	label := firstNonEmpty(jsonString(obj, "termLabel"), "术语")
	inner := section(css("font-size", "11px", "letter-spacing", "2px", "color", t.Primary, "font-weight", "bold"), esc(label)) +
		section(css("font-size", "20px", "font-weight", "bold", "color", t.Text, "margin", "4px 0 6px"), esc(jsonString(obj, "term"))) +
		section(css("font-size", "14px", "color", t.Text, "line-height", "1.75"), esc(jsonString(obj, "def")))
	return section(css("margin", "24px 0", "padding", "16px 18px", "border", "1px solid "+t.Border, "border-left", "4px solid "+t.Primary, "border-radius", "0 10px 10px 0", "background", t.cardBg()), inner)
}

func (r *renderer) renderQuoteCard(b *layoutcatalog.Block) string {
	obj, _ := parseJSON(b.Body).(map[string]any)
	t := r.theme
	inner := section(css("font-size", "18px", "font-weight", "bold", "line-height", "1.7", "color", t.Text, "text-align", "center"), "「"+esc(jsonString(obj, "text"))+"」")
	attribution := []string{}
	for _, key := range []string{"author", "source"} {
		if v := jsonString(obj, key); v != "" {
			attribution = append(attribution, esc(v))
		}
	}
	if len(attribution) > 0 {
		inner += section(css("font-size", "12px", "color", t.Muted, "text-align", "center", "margin", "12px 0 0"), "—— "+strings.Join(attribution, " · "))
	}
	return section(css("margin", "28px 0", "padding", "26px 20px", "background", t.Soft, "border-radius", "12px"), inner)
}

func (r *renderer) renderTweet(b *layoutcatalog.Block) string {
	obj, _ := parseJSON(b.Body).(map[string]any)
	t := r.theme
	name := jsonString(obj, "name")
	avatar := ""
	if src := jsonString(obj, "avatar"); src != "" {
		avatar = `<img src="` + attr(src) + `" style="` + attr(css("width", "40px", "height", "40px", "border-radius", "20px", "display", "block")) + `" />`
	} else {
		avatar = section(css("width", "40px", "height", "40px", "line-height", "40px", "border-radius", "20px", "background", t.Primary, "color", t.onPrimary(), "text-align", "center", "font-weight", "bold"), esc(string([]rune(firstNonEmpty(name, "?"))[0])))
	}
	head := section(css("display", "flex", "align-items", "center"),
		section(css("flex", "0 0 auto", "margin-right", "10px"), avatar)+
			section(css("flex", "1", "min-width", "0"), section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), esc(name))+section(css("font-size", "12px", "color", t.Muted), esc(jsonString(obj, "handle")))))
	inner := head + section(css("font-size", "15px", "line-height", "1.75", "color", t.Text, "margin", "10px 0 0"), esc(jsonString(obj, "text")))
	if ts := jsonString(obj, "timestamp"); ts != "" {
		inner += section(css("font-size", "12px", "color", t.Muted, "margin", "8px 0 0"), esc(ts))
	}
	return r.card(inner, false)
}

func (r *renderer) renderStatRow(b *layoutcatalog.Block) string {
	t := r.theme
	var cells []string
	for _, obj := range jsonObjects(parseJSON(b.Body)) {
		value := esc(jsonString(obj, "value"))
		if unit := jsonString(obj, "unit"); unit != "" {
			value += span(css("font-size", "13px", "margin-left", "2px"), esc(unit))
		}
		content := section(css("font-size", "26px", "font-weight", "bold", "color", t.Primary, "line-height", "1.3"), value) +
			section(css("font-size", "12px", "color", t.Muted, "margin", "2px 0 0"), esc(jsonString(obj, "label")))
		if n := jsonString(obj, "note"); n != "" {
			content += section(css("font-size", "11px", "color", t.Muted, "margin", "2px 0 0"), esc(n))
		}
		cells = append(cells, section(css("padding", "12px 6px", "text-align", "center", "background", t.Soft, "border-radius", "10px", "height", "100%", "box-sizing", "border-box"), content))
	}
	if len(cells) == 0 {
		return ""
	}
	return section(css("margin", "24px 0"), grid(cells, min(len(cells), 4)))
}

func (r *renderer) renderResourceList(b *layoutcatalog.Block) string {
	t := r.theme
	var inner strings.Builder
	for _, obj := range jsonObjects(parseJSON(b.Body)) {
		icon := firstNonEmpty(jsonString(obj, "icon"), string([]rune(firstNonEmpty(jsonString(obj, "name"), "•"))[0]))
		badge := section(css("min-width", "36px", "height", "36px", "line-height", "36px", "padding", "0 6px", "border-radius", "8px", "background", t.Primary, "color", t.onPrimary(), "text-align", "center", "font-size", "13px", "font-weight", "bold", "box-sizing", "border-box"), esc(icon))
		content := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), esc(jsonString(obj, "name")))
		if d := jsonString(obj, "desc"); d != "" {
			content += section(css("font-size", "13px", "color", t.Muted, "margin", "2px 0 0"), esc(d))
		}
		if u := jsonString(obj, "url"); u != "" {
			content += section(css("font-size", "12px", "color", t.Primary, "margin", "2px 0 0", "word-break", "break-all"), esc(u))
		}
		inner.WriteString(section(css("display", "flex", "align-items", "center", "padding", "10px 0", "border-bottom", "1px dashed "+t.Border),
			section(css("flex", "0 0 auto", "margin-right", "12px"), badge)+section(css("flex", "1", "min-width", "0"), content)))
	}
	return r.card(inner.String(), false)
}

func (r *renderer) renderComparisonTable(b *layoutcatalog.Block) string {
	obj, _ := parseJSON(b.Body).(map[string]any)
	t := r.theme
	side := func(key, mark, color string, accent bool) string {
		sideObj, _ := obj[key].(map[string]any)
		bg := t.cardBg()
		if accent {
			bg = t.Soft
		}
		content := section(css("font-size", "15px", "font-weight", "bold", "color", color, "margin", "0 0 6px"), esc(jsonString(sideObj, "title")))
		for _, item := range jsonStrings(sideObj, "items") {
			content += section(css("font-size", "14px", "color", t.Text, "margin", "6px 0 0"), span(css("color", color, "margin-right", "6px"), mark)+esc(item))
		}
		return section(css("padding", "14px", "background", bg, "border", "1px solid "+t.Border, "border-radius", "10px", "height", "100%", "box-sizing", "border-box"), content)
	}
	return section(css("margin", "24px 0"), row(side("left", "✓", "#1f9d55", true), side("right", "✕", "#d64545", false)))
}

func (r *renderer) renderChangelog(b *layoutcatalog.Block) string {
	obj, _ := parseJSON(b.Body).(map[string]any)
	t := r.theme
	var inner strings.Builder
	inner.WriteString(section(css("font-size", "18px", "font-weight", "bold", "color", t.Text), esc(jsonString(obj, "version"))+span(css("font-size", "12px", "font-weight", "normal", "color", t.Muted, "margin-left", "10px"), esc(jsonString(obj, "date")))))
	groups := []struct{ key, label, color string }{
		{"added", "新增", "#1f9d55"}, {"changed", "变更", t.Primary}, {"fixed", "修复", "#d9822b"}, {"removed", "移除", "#d64545"},
	}
	for _, g := range groups {
		items := jsonStrings(obj, g.key)
		if len(items) == 0 {
			continue
		}
		inner.WriteString(section(css("margin", "12px 0 4px"), span(css("display", "inline-block", "padding", "1px 8px", "border-radius", "4px", "background", g.color, "color", "#ffffff", "font-size", "12px"), g.label)))
		for _, item := range items {
			inner.WriteString(section(css("font-size", "14px", "color", t.Text, "margin", "4px 0 0 4px"), "• "+esc(item)))
		}
	}
	return r.card(inner.String(), false)
}

// ---------- free layout ----------

func (r *renderer) renderMatrix(b *layoutcatalog.Block) string {
	rows, _ := parseRows(b)
	t := r.theme
	headers := splitList(b.Params["headers"], ",")
	var table strings.Builder
	table.WriteString(`<table style="` + attr(css("width", "100%", "border-collapse", "collapse", "font-size", "13px")) + `">`)
	if len(headers) > 0 {
		table.WriteString("<thead><tr>")
		for _, h := range headers {
			table.WriteString(`<th style="` + attr(css("padding", "8px 6px", "background", t.Primary, "color", t.onPrimary(), "border", "1px solid "+t.Primary, "text-align", "center")) + `">` + esc(h) + `</th>`)
		}
		table.WriteString("</tr></thead>")
	}
	table.WriteString("<tbody>")
	for i, c := range rows {
		bg := t.cardBg()
		if i%2 == 1 {
			bg = t.Soft
		}
		table.WriteString(`<tr style="` + attr(css("background", bg)) + `">`)
		for j, v := range c {
			style := css("padding", "8px 6px", "border", "1px solid "+t.Border, "text-align", "center", "color", t.Text)
			if j == 0 {
				style = css("padding", "8px 6px", "border", "1px solid "+t.Border, "font-weight", "bold", "color", t.Text)
			}
			table.WriteString(`<td style="` + attr(style) + `">` + r.inline(v) + `</td>`)
		}
		table.WriteString("</tr>")
	}
	table.WriteString("</tbody></table>")
	return section(css("margin", "24px 0", "overflow-x", "auto"), table.String())
}

func (r *renderer) renderDialogue(b *layoutcatalog.Block) string {
	t := r.theme
	turns := parseDialogue(b.Body)
	leftName, rightName := b.Params["left"], b.Params["right"]
	var firstSpeaker string
	if len(turns) > 0 {
		firstSpeaker = turns[0].speaker
	}
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	for _, item := range turns {
		right := item.speaker != firstSpeaker
		if b.Name == "dialogue-pair" {
			right = strings.EqualFold(item.speaker, "E") || strings.EqualFold(item.speaker, "R") || strings.EqualFold(item.speaker, "B")
		}
		name := item.speaker
		if b.Name == "dialogue-pair" {
			if right && rightName != "" {
				name = rightName
			} else if !right && leftName != "" {
				name = leftName
			}
		}
		bubbleBg, bubbleFg, align := t.Soft, t.Text, "left"
		if right {
			bubbleBg, bubbleFg, align = t.Primary, t.onPrimary(), "right"
		}
		inner.WriteString(section(css("margin", "10px 0", "text-align", align),
			section(css("font-size", "12px", "color", t.Muted, "margin", "0 4px 4px"), esc(name))+
				section(css("display", "inline-block", "max-width", "82%", "padding", "10px 14px", "border-radius", "14px", "background", bubbleBg, "color", bubbleFg, "font-size", "14px", "line-height", "1.7", "text-align", "left"), r.inline(item.text))))
	}
	return section(css("margin", "24px 0"), inner.String())
}

func (r *renderer) renderFlow(b *layoutcatalog.Block) string {
	t := r.theme
	var steps []string
	for _, line := range b.Body {
		steps = append(steps, splitList(line, "→")...)
	}
	var inner strings.Builder
	inner.WriteString(r.caption(b.Caption))
	var chips strings.Builder
	for i, step := range steps {
		chips.WriteString(span(css("display", "inline-block", "margin", "4px 0", "padding", "6px 12px", "border-radius", "6px", "background", t.Soft, "color", t.Primary, "font-size", "14px", "font-weight", "bold", "border", "1px solid "+t.Primary), r.inline(step)))
		if i < len(steps)-1 {
			chips.WriteString(span(css("display", "inline-block", "margin", "0 6px", "color", t.Primary, "font-size", "14px"), "→"))
		}
	}
	inner.WriteString(section(css("line-height", "2.2", "text-align", "center"), chips.String()))
	return r.card(inner.String(), false)
}

func (r *renderer) renderSplit(b *layoutcatalog.Block) string {
	t := r.theme
	left, right, _ := splitBody(b.Body, "---")
	sideStyle := css("padding", "12px 14px", "background", t.Soft, "border-radius", "10px", "height", "100%", "box-sizing", "border-box")
	return section(css("margin", "24px 0"), row(
		section(sideStyle, r.markdownToHTML(strings.Join(left, "\n"))),
		section(sideStyle, r.markdownToHTML(strings.Join(right, "\n"))),
	))
}

// ---------- interactive (static) ----------

func (r *renderer) renderExpand(b *layoutcatalog.Block) string {
	head, rest, _ := splitBody(b.Body, "---")
	f, _ := parseFields(head, declaredFields(b))
	t := r.theme
	var inner strings.Builder
	label := firstNonEmpty(f.get("label"), "展开阅读")
	inner.WriteString(section(css("font-size", "12px", "color", t.Primary, "letter-spacing", "1px", "margin", "0 0 4px"), "▾ "+esc(label)))
	inner.WriteString(r.title(f.get("title"), "16px"))
	if s := f.get("summary"); s != "" {
		inner.WriteString(section(css("font-size", "13px", "color", t.Muted, "margin", "0 0 10px"), r.inline(s)))
	}
	inner.WriteString(section(css("padding", "10px 0 0", "border-top", "1px dashed "+t.Border), r.markdownToHTML(strings.Join(rest, "\n"))))
	return r.card(inner.String(), false)
}

func (r *renderer) renderReveal(b *layoutcatalog.Block) string {
	f := r.fields(b)
	t := r.theme
	inner := section(css("font-size", "15px", "font-weight", "bold", "color", t.Text), r.inline(f.get("question"))) +
		section(css("margin", "10px 0 0", "padding", "10px 12px", "background", t.cardBg(), "border-radius", "8px", "font-size", "15px", "color", t.Primary, "font-weight", "bold"), "答案："+r.inline(f.get("answer")))
	return r.softCard(inner)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
