package main

import (
	"net/http"

	"chi_api_rest_products/api_rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// REST routes for restaurant resource
	r.Route("/buyers", func(r chi.Router) {
		r.Post("/", api_rest.ChargeData)  // POST /Charge data
		r.Get("/", api_rest.GetBuyBuyers) // GET /Buyers

		r.Route("/{buyerId}", func(r chi.Router) {
			r.Get("/", api_rest.GetBuyerInformation) // GET /buyer/id_buyer
		})

	})

	http.ListenAndServe(":3000", r)
}
