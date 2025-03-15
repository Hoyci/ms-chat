package healthcheck

import (
	"net/http"

	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"github.com/hoyci/ms-chat/message-service/config"
	"github.com/hoyci/ms-chat/message-service/types"
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
	coreUtils.WriteJSON(w, http.StatusOK, &types.HealthcheckResponse{
		Status: "available",
		SystemInfo: map[string]string{
			"environment": h.cfg.Environment,
		},
	})
}
