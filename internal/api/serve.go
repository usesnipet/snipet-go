package api

import (
	"encoding/json"
	"errors"
	"net/http"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/infra/database"
)

func (a *Api) Serve(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var appErr *apperr.Error
			if errors.As(err, &appErr) {
				w.WriteHeader(appErr.StatusCode)
				json.NewEncoder(w).Encode(appErr)
				return
			} else {
				if err, ok := database.HandleDBError(err); ok {
					w.WriteHeader(err.StatusCode)
					json.NewEncoder(w).Encode(err)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(apperr.InternalServerError("internal server error"))
				return
			}
		}
	}
}
