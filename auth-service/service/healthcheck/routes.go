package healthcheck

import (
	"net/http"

	// "github.com/hoyci/ms-chat/auth-service/types"
	// "github.com/hoyci/ms-chat/auth-service/utils"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/hoyci/ms-chat/ms-core/utils"
)

type HealthCheckHandler struct {
	cfg types.Config
}

func NewHealthCheckHandler(cfg types.Config) *HealthCheckHandler {
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
