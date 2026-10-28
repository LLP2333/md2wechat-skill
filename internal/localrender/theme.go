// Package localrender renders Markdown and :::layout modules into WeChat
// Official Account HTML entirely on the local machine. Every style is inlined
// because the WeChat editor strips <style> blocks and external CSS.
package localrender

import (
	"sort"
	"strings"
)

// Series controls how headings and accents are decorated.
const (
	SeriesClassic = "classic"
	SeriesMinimal = "minimal"
	SeriesFocus   = "focus"
	SeriesElegant = "elegant"
	SeriesBold    = "bold"
)

// Theme is the full set of visual tokens used by the renderer.
type Theme struct {
	Name    string
	Series  string
	Primary string // headings, accents, links
	Accent  string // secondary accent used for highlights
	Soft    string // tinted card background
	Text    string
	Muted   string
	Border  string
	Bg      string // article background; empty means transparent
	CodeBg  string
	CodeFg  string
	Serif   bool
}

type palette struct {
	primary, accent, soft string
}

var colorPalettes = map[string]palette{
	"gold":   {"#b7862f", "#d9a441", "#fbf5e8"},
	"green":  {"#2f9e6e", "#44b883", "#eef8f3"},
	"blue":   {"#2f6fdb", "#4f8bef", "#eef4fd"},
	"orange": {"#e2722e", "#f08f4f", "#fdf2ea"},
	"red":    {"#d64545", "#e86666", "#fdeeee"},
	"navy":   {"#23395d", "#3c5a8a", "#eef1f6"},
	"gray":   {"#5f6b7a", "#7d8896", "#f3f4f6"},
	"sky":    {"#1e9bd7", "#47b4e6", "#ecf7fc"},
}

func baseTheme(name, series, primary, accent, soft string) Theme {
	return Theme{
		Name:    name,
		Series:  series,
		Primary: primary,
		Accent:  accent,
		Soft:    soft,
		Text:    "#3f3f3f",
		Muted:   "#8a8f99",
		Border:  "#e6e8eb",
		CodeBg:  "#f6f8fa",
		CodeFg:  "#24292f",
	}
}

func namedThemes() map[string]Theme {
	themes := map[string]Theme{}
	for _, series := range []string{SeriesMinimal, SeriesFocus, SeriesElegant, SeriesBold} {
		for color, p := range colorPalettes {
			name := series + "-" + color
			themes[name] = baseTheme(name, series, p.primary, p.accent, p.soft)
		}
	}

	add := func(t Theme) { themes[t.Name] = t }

	add(baseTheme("default", SeriesClassic, "#e2725b", "#f0a07a", "#fdf3ef"))
	add(baseTheme("bytedance", SeriesMinimal, "#3370ff", "#5b8cff", "#eef3ff"))
	apple := baseTheme("apple", SeriesBold, "#5e5ce6", "#bf5af2", "#f3f2fe")
	apple.Text = "#1d1d1f"
	add(apple)
	add(baseTheme("sports", SeriesBold, "#ff5a1f", "#ffb020", "#fff1ea"))
	chinese := baseTheme("chinese", SeriesFocus, "#a23b2a", "#c8a15a", "#f8f1e6")
	chinese.Serif = true
	chinese.Bg = "#fbf7ef"
	add(chinese)
	cyber := baseTheme("cyber", SeriesElegant, "#00d1ff", "#ff2bd6", "#15193a")
	cyber.Bg = "#0d1026"
	cyber.Text = "#e4e7f5"
	cyber.Muted = "#8b93b8"
	cyber.Border = "#2a3060"
	cyber.CodeBg = "#161a3d"
	cyber.CodeFg = "#c8f6ff"
	add(cyber)

	add(baseTheme("sspai-red", SeriesMinimal, "#d71a1b", "#ef4b4c", "#fdeeee"))
	add(baseTheme("wechat-native", SeriesClassic, "#07c160", "#36d17f", "#edf9f2"))
	nyt := baseTheme("nyt-classic", SeriesFocus, "#1a1a1a", "#7a6a4f", "#f1ead8")
	nyt.Serif = true
	nyt.Bg = "#f7f3e8"
	nyt.Text = "#222222"
	nyt.Border = "#d8cfb8"
	add(nyt)
	github := baseTheme("github-readme", SeriesMinimal, "#0969da", "#1f883d", "#f6f8fa")
	github.Text = "#1f2328"
	github.Border = "#d0d7de"
	add(github)
	mint := baseTheme("mint-fresh", SeriesElegant, "#20a58a", "#5cc9b0", "#e9f8f4")
	mint.Bg = "#f4fbf9"
	add(mint)
	sunset := baseTheme("sunset-amber", SeriesElegant, "#d9822b", "#f2b25c", "#fdf1e2")
	sunset.Bg = "#fffaf3"
	sunset.Text = "#4a3a2c"
	add(sunset)
	ink := baseTheme("ink-minimal", SeriesFocus, "#111111", "#555555", "#f4f4f4")
	ink.Serif = true
	ink.Text = "#222222"
	add(ink)
	lavender := baseTheme("lavender-dream", SeriesElegant, "#8a6fd1", "#c09ef0", "#f4f0fd")
	lavender.Bg = "#faf8ff"
	add(lavender)
	coffee := baseTheme("coffee-house", SeriesClassic, "#7b4b2a", "#b07a4f", "#f5ece3")
	coffee.Bg = "#fbf7f2"
	coffee.Text = "#3e2c20"
	add(coffee)
	add(baseTheme("bauhaus-primary", SeriesBold, "#d62828", "#1d4ed8", "#fdf6dd"))

	return themes
}

var builtinThemes = namedThemes()

// ResolveTheme returns the named theme, or false when it is unknown.
func ResolveTheme(name string) (Theme, bool) {
	t, ok := builtinThemes[strings.TrimSpace(name)]
	return t, ok
}

// ThemeNames lists every theme the local renderer can draw.
func ThemeNames() []string {
	names := make([]string, 0, len(builtinThemes))
	for name := range builtinThemes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (t Theme) fontFamily() string {
	if t.Serif {
		return `'Songti SC','STSong','Noto Serif SC',serif`
	}
	return `-apple-system,BlinkMacSystemFont,'Helvetica Neue','PingFang SC','Hiragino Sans GB','Microsoft YaHei',Arial,sans-serif`
}

func (t Theme) dark() bool {
	return t.Name == "cyber"
}

// cardBg is the background for module cards; it stays readable on dark themes.
func (t Theme) cardBg() string {
	if t.dark() {
		return "#15193a"
	}
	return "#ffffff"
}

func (t Theme) onPrimary() string {
	if t.dark() {
		return "#0d1026"
	}
	return "#ffffff"
}

const (
	toneSuccess = "success"
	toneWarning = "warning"
	toneDanger  = "danger"
)

// tint is the background behind status text; dark themes need dark tints so
// their light body text stays readable.
func (t Theme) tint(tone string) string {
	if t.dark() {
		switch tone {
		case toneSuccess:
			return "#12302a"
		case toneWarning:
			return "#33280f"
		default:
			return "#3a1724"
		}
	}
	switch tone {
	case toneSuccess:
		return "#eef8f1"
	case toneWarning:
		return "#fdf4e7"
	default:
		return "#fdf0f0"
	}
}
