package converter

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/action"
	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"go.uber.org/zap"
)

func TestLocalModeIsDefaultAndNeedsNoAPIKey(t *testing.T) {
	conv := NewConverter(&config.Config{}, zap.NewNop())
	result := conv.Convert(&ConvertRequest{
		Markdown: "---\ntitle: 标题\n---\n\n## 小节\n\n:::verdict\ntitle: 本地也能渲染高级模块\n:::\n",
	})
	if !result.Success || result.Status != action.StatusCompleted {
		t.Fatalf("Convert() = %#v", result)
	}
	if result.Mode != ModeLocal {
		t.Fatalf("mode = %q, want %q", result.Mode, ModeLocal)
	}
	if strings.Contains(result.HTML, ":::verdict") || strings.Contains(result.HTML, "title: 标题") {
		t.Fatalf("HTML leaked raw syntax or front matter:\n%s", result.HTML)
	}
	if !strings.Contains(result.HTML, "本地也能渲染高级模块") {
		t.Fatalf("module content missing:\n%s", result.HTML)
	}
}

func TestLocalModeAcceptsEverySelectableAPITheme(t *testing.T) {
	tm := NewThemeManager()
	if err := tm.LoadThemes(); err != nil {
		t.Fatalf("LoadThemes() error = %v", err)
	}
	count := 0
	for _, theme := range tm.ListThemeDefinitions() {
		if theme.Type != "api" || !theme.Selectable() {
			continue
		}
		count++
		if _, err := tm.ResolveThemeForMode(ModeLocal, theme.Name); err != nil {
			t.Errorf("theme %q rejected in local mode: %v", theme.Name, err)
		}
	}
	if count < 48 {
		t.Fatalf("selectable API themes = %d, want at least 48", count)
	}
}

func TestLocalModeRejectsAIThemes(t *testing.T) {
	_, err := NewThemeManager().ResolveThemeForMode(ModeLocal, "autumn-warm")
	if !IsThemeModeMismatch(err) {
		t.Fatalf("ResolveThemeForMode(local, autumn-warm) error = %v, want mode mismatch", err)
	}
}

// TestLocalHTMLMatchesUploadPlaceholderFlow keeps local output compatible with
// the publish flow: every ParseMarkdownImages ref must map onto its own <img>.
func TestLocalHTMLMatchesUploadPlaceholderFlow(t *testing.T) {
	markdown := strings.Join([]string{
		"![a](./a.png)",
		"![b](<./b c.png>)",
		"![c](https://example.com/c.png)",
		"![d](__generate:an orange cat__)",
		":::gallery",
		"![e](https://example.com/e.png)",
		":::",
		":::image-text",
		"title: 字段图不属于 Markdown 图片",
		"image: https://example.com/field.png",
		":::",
	}, "\n\n")
	conv := NewConverter(&config.Config{}, zap.NewNop())
	result := conv.Convert(&ConvertRequest{Markdown: markdown, Mode: ModeLocal})
	if !result.Success {
		t.Fatalf("Convert() error = %s", result.Error)
	}
	if len(result.Images) != 5 {
		t.Fatalf("images = %d, want 5", len(result.Images))
	}

	withPlaceholders := InsertImagePlaceholders(result.HTML, result.Images)
	for _, img := range result.Images {
		if strings.Count(withPlaceholders, img.Placeholder) != 1 {
			t.Fatalf("placeholder %s count = %d\n%s", img.Placeholder, strings.Count(withPlaceholders, img.Placeholder), withPlaceholders)
		}
	}
	if !strings.Contains(withPlaceholders, `src="https://example.com/field.png"`) {
		t.Fatalf("field image was consumed by a Markdown placeholder:\n%s", withPlaceholders)
	}
}
