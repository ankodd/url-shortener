package redirect

import (
	"github.com/ankodd/url-shortener/internal/storage"
	"github.com/ankodd/url-shortener/pkg/response"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

var (
	// AliasNotFound is a map for alias not found
	AliasNotFound = "Alias not found"
)

// URLGetter is an interface for getting url from storage
type URLGetter interface {
	GetUrl(alias string) (string, error)
}

// Redirect is a handler for redirecting to url by alias
//
// It returns http.StatusNotFound if alias not found
func Redirect(getter URLGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.WithGroup("redirect.Redirect")
		alias := mux.Vars(r)["alias"]

		url, err := getter.GetUrl(alias)
		if err != nil {
			if err.Error() == storage.ErrAliasNotFound {
				response.Write(w, http.StatusNotFound, errors.New(AliasNotFound))
				log.Error("Alias not found", slog.String("alias", alias))
				return
			}

			response.Write(w, http.StatusInternalServerError, nil)
			log.Error("Failed to get url", slog.String("alias", alias), slog.String("error", err.Error()))
			return
		}

		log.Info("Redirected", slog.String("alias", alias), slog.String("url", url))
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
