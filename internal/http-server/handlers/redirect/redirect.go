package redirect

import (
	"log/slog"
	"net/http"
	"url-shortener/lib/api/response"
	"url-shortener/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		url, err := urlGetter.GetUrl(alias)
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("Couldn't get url"))

			return
		}

		log.Info("got url", slog.String("url", url))

		// redirect to found url
		http.Redirect(w, r, url, http.StatusFound)
	}
}