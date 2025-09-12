package server

import (
	"fmt"
	"net/http"

	"ailap-backend/internal/config"
	"ailap-backend/internal/router"
	"ailap-backend/internal/utils"
)

func StartHTTPServer() error {
	cfg := config.Get()
	r := router.New()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	utils.GetLogger().Sugar().Infow("http server listening", "port", cfg.HTTPPort)
	return srv.ListenAndServe()
}
