package http

import (
	"angsuran-service/config"
	"angsuran-service/internal/usecase"
	"angsuran-service/pkg/logger"

	"github.com/gorilla/mux"
)

func NewRouter(r *mux.Router, l *logger.Logger, cfg *config.Config, au usecase.AngsuranUsecase) {
	{
		NewAngsuranRoutes(r, l, cfg, &au)
	}
}
