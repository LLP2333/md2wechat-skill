package converter

import (
	"github.com/geekjourneyx/md2wechat-skill/internal/action"
	"github.com/geekjourneyx/md2wechat-skill/internal/localrender"
	"go.uber.org/zap"
)

// convertViaLocal 在本机渲染主题和高级排版模块，不调用任何远程服务
func (c *converter) convertViaLocal(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Mode:      ModeLocal,
		Theme:     req.Theme,
		Status:    action.StatusFailed,
		Action:    action.ActionConvert,
		Retryable: false,
		Success:   false,
	}

	theme, err := c.theme.ResolveThemeForMode(ModeLocal, req.Theme)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	html, err := localrender.Render(req.Markdown, localrender.Options{
		Theme:          theme.APITheme,
		FontSize:       req.FontSize,
		BackgroundType: req.BackgroundType,
	})
	if err != nil {
		result.Error = "local render failed: " + err.Error()
		c.log.Error("local conversion failed", zap.String("theme", req.Theme), zap.Error(err))
		return result
	}

	images := c.ExtractImages(req.Markdown)

	result.HTML = html
	result.Images = images
	result.Status = action.StatusCompleted
	result.Success = true

	c.log.Info("local conversion succeeded",
		zap.String("theme", req.Theme),
		zap.Int("image_count", len(images)))

	return result
}
