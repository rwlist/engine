package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rwlist/engine/pkg/conf"
	"github.com/rwlist/engine/pkg/jsonrpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	cfg, err := conf.ParseEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to parse config from env")
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		err2 := http.ListenAndServe(cfg.PrometheusBind, mux)
		if err2 != nil && err2 != http.ErrServerClosed {
			log.WithError(err2).Fatal("prometheus server error")
		}
	}()

	var handler jsonrpc.Handler
	// TODO: init handler

	transport := jsonrpc.NewHTTP(handler)

	mux := http.NewServeMux()
	mux.Handle("/api", transport)

	err = http.ListenAndServe(cfg.ServerBind, mux)
	log.WithError(err).Info("http server finished")
}
