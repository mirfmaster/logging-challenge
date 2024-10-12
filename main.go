package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.With().
			Str("request_id", uuid.New().String()).
			Str("path", r.URL.Path).
			Str("url", r.URL.String()).
			Str("method", r.Method).
			Logger()

		ctx := log.WithContext(r.Context())

		// calculate time elapsed
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		elapsed := time.Since(start)

		log.Info().Float64("elapsed_ms", float64(elapsed.Nanoseconds()/1000000)).Msg("request processed")
	})
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		oscall := <-ch
		log.Warn().Msgf("system call:%+v", oscall)
		cancel()
	}()

	r := mux.NewRouter()
	r.Use(middleware)

	r.HandleFunc("/", handler)

	// start: set up any of your logger configuration here if necessary
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	defer f.Close()

	mw := zerolog.MultiLevelWriter(os.Stdout, f)
	log.Logger = zerolog.New(mw).With().Timestamp().Logger()

	log.Info().
		Msg("Starting the server on port http://localhost:8080")
	// end: set up any of your logger configuration here

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen and serve http server")
		}
	}()
	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server gracefully")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Ctx(ctx).With().Str("func", "handler").Logger()
	name := r.URL.Query().Get("name")

	log.
		Debug().
		Str("name", name).
		Msg("Processing request to the endpoint")

	res, err := greeting(ctx, name)
	if err != nil {
		log.Error().Err(err).Msg("Failed from greetings")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Processing is finished")
	w.Write([]byte(res))
}

func greeting(ctx context.Context, name string) (string, error) {
	log := log.Ctx(ctx).With().Str("func", "greeting").Logger()

	if len(name) < 5 {
		log.Warn().Msg("Name is too short")
		return fmt.Sprintf("Hello %s! Your name is to short\n", name), nil
	}

	log.Debug().Msg("Name is long enough")
	return fmt.Sprintf("Hi %s", name), nil
}
