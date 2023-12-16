package rest

import (
	"fmt"
	"net/http"

	"moneyman-writer-go/internal/core"
	x "moneyman-writer-go/internal/utils/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux *chi.Mux
}

func MakeRouter(svc *core.Service) *Router {
	r := chi.NewRouter()
	controller := NewRestController(svc)

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/healthz", controller.HandleHealth)
	r.Post("/gcs/transactions", controller.HandleGcsTransactionsUploadedEvent)

	return &Router{
		mux: r,
	}
}

func (r *Router) Serve(port int) {
	listenAddr := ":" + fmt.Sprint(port)
	x.Logger().Infow("starting server", "port", port)
	err := http.ListenAndServe(listenAddr, r.mux)
	if err != nil {
		x.Logger().Fatalw("failed to start server", "error", err)
	}
}
