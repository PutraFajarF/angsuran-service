package http

import (
	"angsuran-service/config"
	"angsuran-service/internal/usecase"
	"angsuran-service/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type AngsuranRoutes struct {
	l   *logger.Logger
	cfg *config.Config
	au  usecase.IAngsuranUsecase
}

func NewAngsuranRoutes(r *mux.Router, l *logger.Logger, cfg *config.Config, au usecase.IAngsuranUsecase) {
	angsuranRoutes := &AngsuranRoutes{l, cfg, au}

	group := r.PathPrefix("/v1").Subrouter()
	group.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Ok Service Angsuran Running.."))
	}).Methods(http.MethodGet)
	group.HandleFunc("/angsuran", angsuranRoutes.CalculateAngsuranHandler).Methods(http.MethodPost)
}
