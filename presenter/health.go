package presenter

import (
	"github.com/go-json-experiment/json"
	"github.com/y-yu/kindle-clock-go/domain/build"
	"log/slog"
	"net/http"
)

type HealthResponse struct {
	Version string `json:"version"`
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	out, err := json.Marshal(HealthResponse{
		Version: build.GetCommitHash(),
	})
	if err != nil {
		slog.Error("[HealthHandler#Handle] failed to marshal health response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		slog.Error("[HealthHandler#Handle] failed to write health response", "err", err)
	}
}
