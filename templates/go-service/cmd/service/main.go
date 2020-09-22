package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"gitlab.com/smive/dealer-portal/internal/api"
)

func main() {
	err := run()
	if err != nil {
		logrus.WithError(err).Error("Critical error occured")
		os.Exit(1)
	}
}

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func run() error {
	var cfg config
	err := parseConfig(&cfg, os.Args[1:])
	if err != nil {
		return err
	}

	err = configLogger(cfg)
	if err != nil {
		return fmt.Errorf("configuring logger failed: %w", err)
	}

	expvar.NewString("build").Set(build)

	logrus.WithField("version", build).Info("Starting Application initialization")
	defer logrus.Println("Application completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	logrus.WithField("config", out).Info("Config initialized")
	// TODO better formatting via clearing non visible fields from a copy of the configuration for this func - MR to github.com/ardanlabs/conf
	expvar.NewString("config").Set(out)

	// =========================================================================
	// Start Debug Service

	go func() {
		logrus.WithFields(map[string]interface{}{
			"addr":   cfg.DebugHost,
			"pprof":  fmt.Sprintf("%s/debug/pprof", cfg.DebugHost),
			"expvar": fmt.Sprintf("%s/debug/vars", cfg.DebugHost),
		}).Info("Debug Listening")
		err := http.ListenAndServe(cfg.DebugHost, http.DefaultServeMux)
		if err != nil {
			logrus.WithError(err).Error("Error occured in HTTP listener for debug endpoints")
			return
		}
		logrus.Info("Debug Listener closed")
	}()

	// =========================================================================
	// Start API Service

	logrus.WithField("addr", cfg.Web.APIHost).Info("Starting API")

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      api.NewRouter(cfg.Web.StaticFilesFolder),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logrus.WithField("addr", api.Addr).Info("API started")
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logrus.WithField("sig", sig).Info("Starting shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		// TODO is shutdown cancelling requests or waiting for them? Do I need to take care of anything?
		err := api.Shutdown(ctx)
		if err != nil {
			logrus.WithField("timeout", cfg.Web.ShutdownTimeout).WithError(err).Error("Graceful shutdown did not complete in timeout")
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}

type config struct {
	// TODO document ENV variables in Readme
	LogLevel  string `conf:"default:INFO"`
	LogFormat string `conf:"default:JSON"`

	Web struct {
		APIHost           string        `conf:"default:0.0.0.0:80"`
		ReadTimeout       time.Duration `conf:"default:5s"`
		WriteTimeout      time.Duration `conf:"default:5s"`
		ShutdownTimeout   time.Duration `conf:"default:5s"`
		StaticFilesFolder string        `conf:"default:./public"`
	}

	DebugHost string `conf:"default:0.0.0.0:4000"`
}

func configLogger(cfg config) error {
	switch cfg.LogFormat {
	case "TEXT":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "JSON":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		return fmt.Errorf("unsupported LOG_FORMAT %q", cfg.LogFormat)
	}

	switch cfg.LogLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		return fmt.Errorf("unsupported LOG_LEVEL %q", cfg.LogLevel)
	}

	return nil
}

func parseConfig(cfg *config, args []string) error {
	if err := conf.Parse(os.Args[1:], "DEALER_PORTAL", cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DEALER_PORTAL", cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}
	return nil
}
