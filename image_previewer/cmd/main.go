package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/config"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/http"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/logger"
	srv "github.com/Serega2780/otus-golang-final-project/image_previewer/internal/service"
)

var (
	configFile string
	wg         sync.WaitGroup
)

func init() {
	flag.StringVar(&configFile, "config", "/tmp/config.yaml",
		"Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := config.ReadConfig(configFile)
	log := logger.New(cfg.Logger)
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	service := srv.NewImageProcessingService(log, cfg.Cache)

	httpServer := http.NewServer(ctx, log, cfg.HTTP, service)

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(ctx)
	}()

	wg.Wait()
}
