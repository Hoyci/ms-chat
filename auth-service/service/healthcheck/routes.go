package healthcheck

import (
	"net/http"

	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/hoyci/ms-chat/core/utils"
)

type HealthCheckHandler struct {
	cfg config.Config
}

func NewHealthCheckHandler(cfg config.Config) *HealthCheckHandler {
	return &HealthCheckHandler{
		cfg: cfg,
	}
}

func (h *HealthCheckHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, &types.HealthcheckResponse{
		Status: "available",
		SystemInfo: map[string]string{
			"environment": h.cfg.Environment,
		},
	})
}
