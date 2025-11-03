package httpgin

import (
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/janaka/web-analyzer/internal/adapter/httpfetch"
	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/repository"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
	"github.com/janaka/web-analyzer/pkg/humanizer"
	"github.com/janaka/web-analyzer/pkg/logger"
)

func BuildRouter(cfg *config.Config, log logger.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// r.SetFuncMap(template.FuncMap{
	// 	"add": func(a, b int) int { return a + b },
	// })
	r.SetFuncMap(template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"formatDate": humanizer.FormatDate,
	})

	r.LoadHTMLGlob("internal/view/templates/*")

	// Wire dependencies
	fetch := httpfetch.New(cfg.HTTPDialTimeoutSec, cfg.HTTPTLSTimeoutSec, cfg.HTTPReqTimeoutSec)
	ana := analyzer.NewDefaultAnalyzer(fetch, cfg.MaxLinkWorkers, cfg.MaxLinksPerPage)
	repo := repository.NewAnalysisRepository(cfg, log)

	h := NewHandlers(cfg, log, ana, repo)

	// UI
	r.GET("/", h.Index)               // form page
	r.POST("/analyze", h.AnalyzeForm) // form submit -> results page

	// API v1
	api := r.Group("/api/v1")
	{
		api.GET("/analysis/:id", h.ViewAnalysis)
	}

	// health check
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	return r
}
