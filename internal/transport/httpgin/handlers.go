package httpgin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/domain"
	"github.com/janaka/web-analyzer/pkg/humanizer"
	"github.com/janaka/web-analyzer/pkg/logger"
	"github.com/janaka/web-analyzer/pkg/validator"
)

type Analyzer interface {
	Analyze(ctx context.Context, url string) (*domain.Analysis, error)
}
type Repository interface {
	Save(ctx context.Context, a *domain.Analysis) error
	ListRecent(ctx context.Context, limit int) ([]domain.Analysis, error)
	GetByID(ctx context.Context, id uint) (*domain.Analysis, error)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, _ := h.repo.ListRecent(ctx, 10) // latest 10 entries
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Rows": rows,
	})
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

	res, analyzeErr := h.ana.Analyze(ctx, url)
	_ = h.repo.Save(ctx, res)

	// Compute human-friendly error message (if any)
	friendlyError := ""
	if analyzeErr != nil {
		friendlyError = humanizer.HTTPError(analyzeErr.Error(), 0)
	} else if res != nil && res.ErrorMessage != "" {
		friendlyError = humanizer.HTTPError(res.ErrorMessage, res.HTTPStatus)
	}
	fmt.Println("FriendlyError:", friendlyError)
	c.HTML(http.StatusOK, "result.html", gin.H{"Result": res, "FriendlyError": friendlyError})
}

// ViewAnalysis - show single analysis result (HTML reused by popup)
func (h *Handlers) ViewAnalysis(c *gin.Context) {
	idStr := c.Param("id")
	fmt.Println("🧩 idStr:", idStr)
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	row, err := h.repo.GetByID(ctx, uint(idUint))
	fmt.Println("🧩 row:", row)
	if err != nil || row == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Analysis not found.",
		})
		return
	}

	friendly := ""
	if row.ErrorMessage != "" {
		friendly = humanizer.HTTPError(row.ErrorMessage, row.HTTPStatus)
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Result":        row,
		"FriendlyError": friendly,
	})
}
