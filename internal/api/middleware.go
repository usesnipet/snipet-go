package api

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler
