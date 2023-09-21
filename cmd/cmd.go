package cmd

import (
	"StoryGoAPI/config"
	myecho "StoryGoAPI/controller/echo"
	"StoryGoAPI/repository"
	"StoryGoAPI/service"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {
	repo, err := repository.NewRepo(cfg)
	if err != nil {
		return err
	}

	srv := service.NewService(cfg, repo)

	rest := myecho.NewRest(cfg, srv)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-ch

		if err := rest.Stop(); err != nil {
			panic(err)
		}
	}()

	if cfg.Server.Debug != "" {
		go func() {
			server := http.Server{
				Addr: cfg.Server.Debug,
				//Handler: http.DefaultServeMux, // DefaultServeMux is served by default if no handler is provided
			}
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)

			}
		}()
	}

	if err = rest.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil

}
