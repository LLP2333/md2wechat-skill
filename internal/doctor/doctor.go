package doctor

import (
	"fmt"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/converter"
	"github.com/geekjourneyx/md2wechat-skill/internal/layoutcatalog"
)

const (
	StatusPass = "pass"
	StatusWarn = "warn"
	StatusFail = "fail"

	OverallReady    = "ready"
	OverallDegraded = "degraded"
	OverallBlocked  = "blocked"
)

type Check struct {
	ID          string   `json:"id"`
	Status      string   `json:"status"`
	Blocking    bool     `json:"blocking"`
	Message     string   `json:"message"`
	NextActions []string `json:"next_actions,omitempty"`
}

type Readiness struct {
	FormatLocal    bool `json:"format_local"`
	FormatAPI      bool `json:"format_api"`
	AdvancedLayout bool `json:"advanced_layout"`
	Draft          bool `json:"draft"`
}

type Report struct {
	Overall   string    `json:"overall"`
	Live      bool      `json:"live"`
	Checks    []Check   `json:"checks"`
	Readiness Readiness `json:"readiness"`
}

func RunLocal(cfg *config.Config) Report {
	if cfg == nil {
		return LoadError(fmt.Errorf("configuration is nil"))
	}

	report := Report{Live: false}
	report.add(Check{
		ID:      "config.load",
		Status:  StatusPass,
		Message: "configuration loaded",
	})

	modeOK := report.checkMode(cfg)
	apiOK := report.checkAPIConfig(cfg, modeOK)
	themeOK := report.checkThemeCatalog(cfg, modeOK)
	layoutOK := report.checkLayoutCatalog()
	draftOK := report.checkWeChatConfig(cfg)

	report.Readiness.FormatLocal = modeOK && cfg.DefaultConvertMode == "local" && themeOK
	report.Readiness.FormatAPI = modeOK && cfg.DefaultConvertMode == "api" && apiOK && themeOK
	report.Readiness.AdvancedLayout = (report.Readiness.FormatLocal || report.Readiness.FormatAPI) && layoutOK
	report.Readiness.Draft = draftOK
	report.Overall = overallFromChecks(report.Checks)
	return report
}

func LoadError(err error) Report {
	report := Report{
		Live: false,
		Checks: []Check{{
			ID:       "config.load",
			Status:   StatusFail,
			Blocking: true,
			Message:  err.Error(),
			NextActions: []string{
				"Run `md2wechat config init` or fix the current configuration file.",
				"Run `md2wechat config validate --json` for the exact config error.",
			},
		}},
	}
	report.Overall = overallFromChecks(report.Checks)
	return report
}

func (r *Report) checkMode(cfg *config.Config) bool {
	mode := strings.TrimSpace(cfg.DefaultConvertMode)
	if mode == "local" || mode == "api" || mode == "ai" {
		r.add(Check{
			ID:      "config.mode",
			Status:  StatusPass,
			Message: fmt.Sprintf("default convert mode is %q", mode),
		})
		return true
	}
	r.add(Check{
		ID:       "config.mode",
		Status:   StatusFail,
		Blocking: true,
		Message:  fmt.Sprintf("default convert mode %q is invalid", mode),
		NextActions: []string{
			"Set `api.convert_mode: local` or `CONVERT_MODE=local` for free local rendering.",
			"Use `ai` only when you intentionally want AI-mode prompt output.",
		},
	})
	return false
}

func (r *Report) checkAPIConfig(cfg *config.Config, modeOK bool) bool {
	if modeOK && strings.TrimSpace(cfg.DefaultConvertMode) != "api" {
		r.add(Check{
			ID:      "api.config",
			Status:  StatusPass,
			Message: fmt.Sprintf("%s mode needs no API key", strings.TrimSpace(cfg.DefaultConvertMode)),
		})
		return false
	}
	missing := []string{}
	if strings.TrimSpace(cfg.MD2WechatAPIKey) == "" {
		missing = append(missing, "MD2WECHAT_API_KEY")
	}
	if strings.TrimSpace(cfg.MD2WechatBaseURL) == "" {
		missing = append(missing, "MD2WECHAT_BASE_URL")
	}
	if len(missing) == 0 {
		r.add(Check{
			ID:      "api.config",
			Status:  StatusPass,
			Message: "API conversion credentials are present",
		})
		return true
	}

	blocking := modeOK && cfg.DefaultConvertMode == "api"
	status := StatusWarn
	if blocking {
		status = StatusFail
	}
	r.add(Check{
		ID:       "api.config",
		Status:   status,
		Blocking: blocking,
		Message:  fmt.Sprintf("missing API conversion configuration: %s", strings.Join(missing, ", ")),
		NextActions: []string{
			"Set `MD2WECHAT_API_KEY` or `api.md2wechat_key`.",
			"Run `md2wechat config show --format json` to confirm the effective local config.",
		},
	})
	return false
}

func (r *Report) checkThemeCatalog(cfg *config.Config, modeOK bool) bool {
	tm := converter.NewThemeManager()
	if err := tm.LoadThemes(); err != nil {
		r.add(Check{
			ID:       "theme.catalog",
			Status:   StatusFail,
			Blocking: true,
			Message:  fmt.Sprintf("theme catalog failed to load: %v", err),
			NextActions: []string{
				"Run `md2wechat themes list --json` to inspect theme catalog errors.",
			},
		})
		return false
	}
	if !modeOK {
		r.add(Check{
			ID:      "theme.catalog",
			Status:  StatusWarn,
			Message: "theme compatibility was skipped because convert mode is invalid",
		})
		return false
	}
	theme, err := tm.ResolveThemeForMode(converter.ConvertMode(cfg.DefaultConvertMode), cfg.DefaultTheme)
	if err != nil {
		r.add(Check{
			ID:       "theme.catalog",
			Status:   StatusFail,
			Blocking: true,
			Message:  err.Error(),
			NextActions: []string{
				"Run `md2wechat themes list --json` and choose a selectable theme for the current mode.",
				"Set `DEFAULT_THEME` or `api.default_theme` to a compatible theme.",
			},
		})
		return false
	}
	r.add(Check{
		ID:      "theme.catalog",
		Status:  StatusPass,
		Message: fmt.Sprintf("default theme %q is selectable for %s mode", theme.Name, cfg.DefaultConvertMode),
	})
	return true
}

func (r *Report) checkLayoutCatalog() bool {
	cat, err := layoutcatalog.DefaultCatalog()
	if err != nil {
		r.add(Check{
			ID:       "layout.catalog",
			Status:   StatusFail,
			Blocking: true,
			Message:  fmt.Sprintf("layout catalog failed to load: %v", err),
			NextActions: []string{
				"Run `md2wechat layout list --json` to inspect layout catalog errors.",
			},
		})
		return false
	}
	modules := cat.ListFiltered(layoutcatalog.ListFilter{})
	if len(modules) == 0 {
		r.add(Check{
			ID:       "layout.catalog",
			Status:   StatusFail,
			Blocking: true,
			Message:  "layout catalog loaded with zero modules",
			NextActions: []string{
				"Reinstall md2wechat or restore the built-in layout assets.",
			},
		})
		return false
	}
	r.add(Check{
		ID:      "layout.catalog",
		Status:  StatusPass,
		Message: fmt.Sprintf("layout catalog loaded with %d modules", len(modules)),
	})
	return true
}

func (r *Report) checkWeChatConfig(cfg *config.Config) bool {
	missing := []string{}
	if strings.TrimSpace(cfg.WechatAppID) == "" {
		missing = append(missing, "WECHAT_APPID")
	}
	if strings.TrimSpace(cfg.WechatSecret) == "" {
		missing = append(missing, "WECHAT_SECRET")
	}
	if len(missing) == 0 {
		r.add(Check{
			ID:      "wechat.config",
			Status:  StatusPass,
			Message: "WeChat draft credentials are present",
		})
		return true
	}
	r.add(Check{
		ID:      "wechat.config",
		Status:  StatusWarn,
		Message: fmt.Sprintf("draft creation is not ready; missing %s", strings.Join(missing, ", ")),
		NextActions: []string{
			"Set `WECHAT_APPID` and `WECHAT_SECRET` before creating WeChat drafts.",
			"Skip draft commands if you only need local convert or preview.",
		},
	})
	return false
}

func (r *Report) add(check Check) {
	r.Checks = append(r.Checks, check)
}

func overallFromChecks(checks []Check) string {
	degraded := false
	for _, check := range checks {
		if check.Status == StatusFail && check.Blocking {
			return OverallBlocked
		}
		if check.Status == StatusWarn || check.Status == StatusFail {
			degraded = true
		}
	}
	if degraded {
		return OverallDegraded
	}
	return OverallReady
}
