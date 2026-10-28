package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
)

func TestPrepareWeChatSideEffectNamedAccountNeedsNoAPIKey(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
	})

	cfg = &config.Config{
		WechatAccounts: map[string]config.WechatAccount{
			"main": {AppID: "appid", Secret: "secret"},
		},
	}
	wechatAccountName = "main"

	if err := prepareWeChatSideEffect(); err != nil {
		t.Fatalf("prepareWeChatSideEffect() error = %v", err)
	}
	if cfg.WechatAppID != "appid" || cfg.WechatSecret != "secret" {
		t.Fatalf("selected credentials = %q/%q", cfg.WechatAppID, cfg.WechatSecret)
	}
}

func TestPrepareWeChatSideEffectHonorsFlagAfterAmbiguousLoad(t *testing.T) {
	oldCfg, oldLog, oldAccount := cfg, log, wechatAccountName
	t.Cleanup(func() {
		cfg = oldCfg
		log = oldLog
		wechatAccountName = oldAccount
	})

	home := t.TempDir()
	configDir := filepath.Join(home, ".config", "md2wechat")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(strings.TrimSpace(`
wechat:
  accounts:
    main:
      appid: wx-main
      secret: secret-main
    client-a:
      appid: wx-client
      secret: secret-client
`)+"\n"), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("WECHAT_APPID", "")
	t.Setenv("WECHAT_SECRET", "")
	t.Setenv("WECHAT_ACCOUNT", "")
	t.Setenv("MD2WECHAT_API_KEY", "")

	cfg = nil
	log = nil
	wechatAccountName = "main"

	if err := initConfig(); err != nil {
		t.Fatalf("initConfig() error = %v", err)
	}
	if err := prepareWeChatSideEffect(); err != nil {
		t.Fatalf("prepareWeChatSideEffect() error = %v", err)
	}
	if cfg.WechatAppID != "wx-main" || cfg.WechatSecret != "secret-main" {
		t.Fatalf("selected credentials = %q/%q", cfg.WechatAppID, cfg.WechatSecret)
	}
}

func TestRunInspectFailsClosedForMissingExplicitWechatAccount(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	oldUpload, oldDraft := inspectUpload, inspectDraft
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
		inspectUpload = oldUpload
		inspectDraft = oldDraft
	})

	cfg = &config.Config{
		WechatAppID:  "direct-appid",
		WechatSecret: "direct-secret",
		WechatAccounts: map[string]config.WechatAccount{
			"main": {AppID: "wx-main", Secret: "secret-main"},
		},
		MD2WechatAPIKey: "api-key",
	}
	wechatAccountName = "missing"
	inspectDraft = true
	inspectUpload = false

	markdownPath := filepath.Join(t.TempDir(), "article.md")
	if err := os.WriteFile(markdownPath, []byte("# Title\n"), 0600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	_, err := runInspect(markdownPath)
	cliErr, ok := err.(*cliError)
	if !ok || cliErr.Code != codeWechatAccountNotFound {
		t.Fatalf("runInspect() error = %#v", err)
	}
}

func TestPrepareWeChatSideEffectSideEffectPathsNeedNoAPIKey(t *testing.T) {
	cases := map[string]*config.Config{
		"direct legacy": {WechatAppID: "appid", WechatSecret: "secret"},
		"proxy url":     {WechatAppID: "appid", WechatSecret: "secret", WechatProxyURL: "https://proxy.example.com"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			oldCfg, oldAccount := cfg, wechatAccountName
			t.Cleanup(func() {
				cfg = oldCfg
				wechatAccountName = oldAccount
			})
			cfg = c
			wechatAccountName = ""
			if err := prepareWeChatSideEffect(); err != nil {
				t.Fatalf("prepareWeChatSideEffect() error = %v", err)
			}
		})
	}
}

func TestPrepareWeChatSideEffectStillRequiresWeChatCredentials(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
	})
	cfg = &config.Config{MD2WechatAPIKey: "api-key"}
	wechatAccountName = ""

	err := prepareWeChatSideEffect()
	cliErr, ok := err.(*cliError)
	if !ok || cliErr.Code != codeConfigInvalid {
		t.Fatalf("prepareWeChatSideEffect() error = %#v", err)
	}
}

func TestValidateConvertConfigFailsClosedForMissingExplicitWechatAccountWithoutSideEffect(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	oldMode, oldAPIKey := convertMode, convertAPIKey
	oldUpload, oldDraft := convertUpload, convertDraft
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
		convertMode = oldMode
		convertAPIKey = oldAPIKey
		convertUpload = oldUpload
		convertDraft = oldDraft
	})

	cfg = &config.Config{
		WechatAppID:        "direct-appid",
		WechatSecret:       "direct-secret",
		MD2WechatAPIKey:    "api-key",
		DefaultConvertMode: "api",
		WechatAccounts: map[string]config.WechatAccount{
			"main": {AppID: "wx-main", Secret: "secret-main"},
		},
	}
	wechatAccountName = "missing"
	convertMode = "api"
	convertAPIKey = ""
	convertUpload = false
	convertDraft = false

	err := validateConvertConfig()
	cliErr, ok := err.(*cliError)
	if !ok || cliErr.Code != codeWechatAccountNotFound {
		t.Fatalf("validateConvertConfig() error = %#v", err)
	}
}

func TestRunInspectFailsClosedForMissingExplicitWechatAccountWithoutTarget(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	oldUpload, oldDraft := inspectUpload, inspectDraft
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
		inspectUpload = oldUpload
		inspectDraft = oldDraft
	})

	cfg = &config.Config{
		WechatAppID:     "direct-appid",
		WechatSecret:    "direct-secret",
		MD2WechatAPIKey: "api-key",
		WechatAccounts: map[string]config.WechatAccount{
			"main": {AppID: "wx-main", Secret: "secret-main"},
		},
	}
	wechatAccountName = "missing"
	inspectDraft = false
	inspectUpload = false

	markdownPath := filepath.Join(t.TempDir(), "article.md")
	if err := os.WriteFile(markdownPath, []byte("# Title\n"), 0600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	_, err := runInspect(markdownPath)
	cliErr, ok := err.(*cliError)
	if !ok || cliErr.Code != codeWechatAccountNotFound {
		t.Fatalf("runInspect() error = %#v", err)
	}
}

func TestRunCreateImagePostDryRunFailsClosedForMissingExplicitWechatAccount(t *testing.T) {
	oldCfg, oldAccount := cfg, wechatAccountName
	oldTitle, oldImages := imagePostTitle, imagePostImages
	oldFromMD, oldDryRun := imagePostFromMD, imagePostDryRun
	oldContent, oldOutput := imagePostContent, imagePostOutput
	t.Cleanup(func() {
		cfg = oldCfg
		wechatAccountName = oldAccount
		imagePostTitle = oldTitle
		imagePostImages = oldImages
		imagePostFromMD = oldFromMD
		imagePostDryRun = oldDryRun
		imagePostContent = oldContent
		imagePostOutput = oldOutput
	})

	cfg = &config.Config{
		WechatAppID:  "direct-appid",
		WechatSecret: "direct-secret",
		WechatAccounts: map[string]config.WechatAccount{
			"main": {AppID: "wx-main", Secret: "secret-main"},
		},
	}
	wechatAccountName = "missing"
	imagePostTitle = "Title"
	imagePostImages = "image.png"
	imagePostFromMD = ""
	imagePostDryRun = true
	imagePostContent = ""
	imagePostOutput = ""

	_, err := runCreateImagePost()
	cliErr, ok := err.(*cliError)
	if !ok || cliErr.Code != codeWechatAccountNotFound {
		t.Fatalf("runCreateImagePost() error = %#v", err)
	}
}
