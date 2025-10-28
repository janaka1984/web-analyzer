package httpgin

import (
	"github.com/gin-gonic/gin"
	"github.com/janaka/web-analyzer/internal/adapter/httpfetch"
	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/repository"
	"github.com/janaka/web-analyzer/internal/service/analyzer"
	"github.com/janaka/web-analyzer/pkg/logger"
)

func BuildRouter(cfg *config.Config, log logger.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.LoadHTMLGlob("internal/view/templates/*")

	// Wire dependencies
	fetch := httpfetch.New(cfg.HTTPDialTimeoutSec, cfg.HTTPTLSTimeoutSec, cfg.HTTPReqTimeoutSec)
	ana := analyzer.NewDefaultAnalyzer(fetch, cfg.MaxLinkWorkers, cfg.MaxLinksPerPage)
	repo := repository.NewAnalysisRepository(cfg, log)

	h := NewHandlers(cfg, log, ana, repo)

	// UI
	r.GET("/", h.Index)               // form page
	r.POST("/analyze", h.AnalyzeForm) // form submit -> results page

	// API JSON
	api := r.Group("/api/v1")
	{
		api.POST("/analyze", h.AnalyzeJSON)
		api.GET("/analyses", h.ListAnalyses)
	}

	// health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	return r
}
