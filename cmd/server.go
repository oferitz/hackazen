package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (app *application) serve() error {
	srv := app.server
	cfg := app.config

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server", map[string]string{
			"signal": s.String(),
		})

		err := srv.Shutdown()
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks")

		app.wg.Wait()
		shutdownError <- nil
	}()

	err := srv.Listen(fmt.Sprintf("localhost:%d", cfg.Int("app.port")))
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server")

	return nil
}
