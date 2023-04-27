package webserver

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oats/backend/web/website"
)

func (s *WebServer) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != homeUrl {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.l.Debug("webserver - handling home request")

	response, err := website.Templates.ReadFile(website.HomeTemplate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		s.l.Error("webserver - could not read template for home page", zap.Error(err))
		return
	}

	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		s.l.Error("webserver - could not write response",
			zap.String("url", homeUrl),
			zap.Error(err))
		return
	}
}
