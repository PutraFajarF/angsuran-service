package app

import (
	"angsuran-service/config"
	"angsuran-service/internal/controller/http"
	"angsuran-service/internal/usecase"
	"angsuran-service/pkg/httpserver"
	"angsuran-service/pkg/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func Run(cfg *config.Config) {
	fmt.Println("Start service-angsuran..")

	var err error
	l := logger.New(cfg)

	// Usecase
	angsuranUsecase := usecase.NewAngsuranUsecase(l, cfg)

	// HTTP Server
	handler := mux.NewRouter()
	http.NewRouter(handler, l, cfg, *angsuranUsecase)
	httpServer := httpserver.New(handler, cfg, httpserver.Port(cfg.HTTPServer.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
