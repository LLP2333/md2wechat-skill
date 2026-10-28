package localrender

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/layoutcatalog"
)

func allSpecs(t *testing.T) []*layoutcatalog.LayoutSpec {
	t.Helper()
	catalog, err := layoutcatalog.DefaultCatalog()
	if err != nil {
		t.Fatalf("DefaultCatalog() error = %v", err)
	}
	specs := catalog.ListFiltered(layoutcatalog.ListFilter{Lifecycle: layoutcatalog.LifecycleRecommended})
	specs = append(specs, catalog.ListFiltered(layoutcatalog.ListFilter{Lifecycle: layoutcatalog.LifecycleCompatibility})...)
	if len(specs) == 0 {
		t.Fatal("layout catalog is empty")
	}
	return specs
}

// TestEveryCatalogModuleHasLocalRenderer fails when a module YAML is added
// without a matching local renderer, so local mode never silently degrades.
func TestEveryCatalogModuleHasLocalRenderer(t *testing.T) {
	for _, spec := range allSpecs(t) {
		if !HasModule(spec.Name) {
			t.Errorf("module %q has no local renderer", spec.Name)
		}
	}
}

// TestCatalogExamplesRenderWithoutRawSyntax renders every canonical example
// and variant witness the catalog ships, across theme series.
func TestCatalogExamplesRenderWithoutRawSyntax(t *testing.T) {
	themes := []string{"default", "minimal-blue", "focus-red", "elegant-gold", "bold-navy", "cyber"}
	for _, spec := range allSpecs(t) {
		examples := []string{spec.Example}
		for _, variant := range spec.Variants {
			examples = append(examples, variant.Example)
		}
		for i, example := range examples {
			if strings.TrimSpace(example) == "" {
				continue
			}
			for _, theme := range themes {
				out, err := Render(example, Options{Theme: theme})
				if err != nil {
					t.Fatalf("%s example %d theme %s: Render() error = %v", spec.Name, i, theme, err)
				}
				if strings.Contains(out, ":::"+spec.Name) {
					t.Errorf("%s example %d theme %s left raw syntax:\n%s", spec.Name, i, theme, out)
				}
				if strings.Contains(out, "<style") || strings.Contains(out, "<script") || strings.Contains(out, "class=") {
					t.Errorf("%s example %d theme %s emitted non-inline styling", spec.Name, i, theme)
				}
			}
		}
	}
}

func TestLayoutE2EFileRendersEveryModule(t *testing.T) {
	src := readFile(t, "../../examples/layout-e2e-test.md")
	out, err := Render(src, Options{Theme: "elegant-blue"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, spec := range allSpecs(t) {
		if strings.Contains(out, "\n:::"+spec.Name) || strings.Contains(out, ">:::"+spec.Name) {
			t.Errorf("module %q left raw syntax in e2e output", spec.Name)
		}
	}
}

func TestEveryAPIThemeNameResolves(t *testing.T) {
	// Names from internal/assets/builtin/themes/api.yaml.
	names := []string{
		"default", "bytedance", "apple", "sports", "chinese", "cyber",
		"sspai-red", "wechat-native", "nyt-classic", "github-readme", "mint-fresh",
		"sunset-amber", "ink-minimal", "lavender-dream", "coffee-house", "bauhaus-primary",
	}
	for _, series := range []string{"minimal", "focus", "elegant", "bold"} {
		for _, color := range []string{"gold", "green", "blue", "orange", "red", "navy", "gray", "sky"} {
			names = append(names, series+"-"+color)
		}
	}
	if len(names) != 48 {
		t.Fatalf("expected 48 theme names, got %d", len(names))
	}
	for _, name := range names {
		if _, ok := ResolveTheme(name); !ok {
			t.Errorf("theme %q does not resolve locally", name)
		}
	}
	if len(ThemeNames()) != 48 {
		t.Errorf("ThemeNames() = %d, want 48", len(ThemeNames()))
	}
}

func TestUnknownThemeIsAnError(t *testing.T) {
	_, err := Render("# hi", Options{Theme: "does-not-exist"})
	var themeErr *ThemeError
	if !errors.As(err, &themeErr) {
		t.Fatalf("Render() error = %v, want ThemeError", err)
	}
}

// TestImageSourcesSurviveRendering keeps src byte-identical to the Markdown
// reference, which is what the upload flow matches on.
func TestImageSourcesSurviveRendering(t *testing.T) {
	src := strings.Join([]string{
		"![本地](./images/封面.png)",
		"![角括号](<./images/a b.png>)",
		"![远程](https://example.com/a.png?x=1&y=2)",
		"![生成](__generate:A cute orange cat__)",
		":::gallery",
		"![模块图](https://example.com/g.png)",
		":::",
	}, "\n\n")
	out, err := Render(src, Options{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{
		`src="./images/封面.png"`,
		`src="./images/a b.png"`,
		`src="https://example.com/a.png?x=1&amp;y=2"`,
		`src="A cute orange cat"`,
		`src="https://example.com/g.png"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %s\n%s", want, out)
		}
	}
	if strings.Contains(out, imageTokenPrefix) {
		t.Errorf("image token leaked into output:\n%s", out)
	}
}

func TestFieldValuesAreEscaped(t *testing.T) {
	src := ":::verdict\ntitle: <script>alert(1)</script>\nbody: a & b\n:::"
	out, err := Render(src, Options{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if strings.Contains(out, "<script>") {
		t.Fatalf("field value was not escaped:\n%s", out)
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Fatalf("escaped title missing:\n%s", out)
	}
}

func TestDirectivesInsideCodeFencesStayCode(t *testing.T) {
	src := "```markdown\n:::hero\ntitle: 示例\n:::\n```"
	out, err := Render(src, Options{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(out, ":::hero") || !strings.Contains(out, "<pre") {
		t.Fatalf("fenced directive should render as code:\n%s", out)
	}
}

func TestUnknownModuleFallsBackToCard(t *testing.T) {
	src := ":::future-module\ntitle: 未来模块\n正文内容\n:::"
	out, err := Render(src, Options{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if strings.Contains(out, ":::future-module") || !strings.Contains(out, "未来模块") || !strings.Contains(out, "正文内容") {
		t.Fatalf("unknown module fallback output:\n%s", out)
	}
}

func TestMarkdownElementsGetInlineStyles(t *testing.T) {
	src := "## 标题\n\n段落 **加粗** `code`\n\n> 引用\n\n- 列表\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n```go\nfmt.Println(1)\n```"
	out, err := Render(src, Options{Theme: "bold-blue", FontSize: "large"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, tag := range []string{"<h2 style=", "<p style=", "<strong style=", "<code style=", "<blockquote style=", "<li style=", "<table style=", "<td style=", "<pre style="} {
		if !strings.Contains(out, tag) {
			t.Errorf("missing inline style on %s\n%s", tag, out)
		}
	}
	if !strings.Contains(out, "font-size:17px") {
		t.Errorf("large font size not applied")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
