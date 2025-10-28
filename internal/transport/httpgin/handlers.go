package httpgin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/domain"
	"github.com/janaka/web-analyzer/pkg/logger"
	"github.com/janaka/web-analyzer/pkg/validator"
)

type Analyzer interface {
	Analyze(ctx context.Context, url string) (*domain.Analysis, error)
}
type Repository interface {
	Save(ctx context.Context, a *domain.Analysis) error
	ListRecent(ctx context.Context, limit int) ([]domain.Analysis, error)
}

type Handlers struct {
	cfg  *config.Config
	log  logger.Logger
	ana  Analyzer
	repo Repository
}

func NewHandlers(cfg *config.Config, log logger.Logger, a Analyzer, r Repository) *Handlers {
	return &Handlers{cfg: cfg, log: log, ana: a, repo: r}
}

// -------- UI (templates) --------

func (h *Handlers) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (h *Handlers) AnalyzeForm(c *gin.Context) {
	raw := c.PostForm("url")
	url, err := validator.NormalizeURL(raw)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Invalid URL"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.RequestTimeoutSec)*time.Second)
	defer cancel()

	res, _ := h.ana.Analyze(ctx, url)
	_ = h.repo.Save(ctx, res) // best-effort save
	c.HTML(http.StatusOK, "result.html", gin.H{"Result": res})
}

// -------- JSON API --------

func (h *Handlers) AnalyzeJSON(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url required"})
		return
	}
	url, err := validator.NormalizeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.cfg.RequestTimeoutSec)*time.Second)
	defer cancel()

	res, err := h.ana.Analyze(ctx, url)
	if err != nil && res != nil && res.HTTPStatus == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": res.ErrorMessage})
		return
	}
	_ = h.repo.Save(ctx, res)
	c.JSON(http.StatusOK, res)
}

func (h *Handlers) ListAnalyses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	rows, err := h.repo.ListRecent(ctx, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, rows)
}
