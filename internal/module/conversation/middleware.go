package conversation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/usesnipet/snipet/internal/api"
	"github.com/usesnipet/snipet/internal/module/client"
)

func ClientExistsMiddleware(service *client.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := chi.URLParam(r, "client_id")
			_, err := service.FindByID(r.Context(), clientID)
			if err != nil {
				api.WriteJSON(w, http.StatusNotFound, err.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
