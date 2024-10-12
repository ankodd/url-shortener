package save

import (
	"encoding/json"
	"github.com/ankodd/url-shortener/internal/storage"
	"github.com/ankodd/url-shortener/pkg/alias"
	"github.com/ankodd/url-shortener/pkg/response"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"net/url"
)

var (
	// InvalidURL error for invalid url
	InvalidURL = "invalid url"

	// URLDoesntContainScheme error for url doesn't contain scheme
	URLDoesntContainScheme = "url doesn't contain scheme"

	// AliasAlreadyExists error for alias already exists
	AliasAlreadyExists = "alias already exists"
)

// URLSaver is an interface for saving url to storage
type URLSaver interface {
	SaveURL(string, string) (int64, error)
}

// URL struct for parsing json
type URL struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

// Save handler for saving url to storage
//
// If alias is empty, it will be generated
func Save(saver URLSaver, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.WithGroup("save.Save")
		var savingURL URL

		if err := json.NewDecoder(r.Body).Decode(&savingURL); err != nil {
			response.Write(w, http.StatusBadRequest, err)
			log.Error("Failed to parse request body", slog.String("error", err.Error()))
			return
		}

		URL, err := url.Parse(savingURL.URL)
		if err != nil {
			response.Write(w, http.StatusBadRequest, errors.New(InvalidURL))
			log.Error(InvalidURL, slog.String("error", err.Error()))
			return
		}

		if URL.Scheme != "http" && URL.Scheme != "https" {
			response.Write(w, http.StatusBadRequest, errors.New(URLDoesntContainScheme))
			log.Error(URLDoesntContainScheme, slog.String("scheme", URL.Scheme))
			return
		}

		if savingURL.Alias == "" {
			savingURL.Alias = alias.Generate()
		}

		id, err := saver.SaveURL(savingURL.URL, savingURL.Alias)
		if err != nil {
			if err.Error() == storage.ErrAliasAlreadyExists {
				response.Write(w, http.StatusConflict, errors.New(AliasAlreadyExists))
				log.Error(AliasAlreadyExists, slog.String("error", err.Error()))
				return
			}

			response.Write(w, http.StatusInternalServerError, err)
			log.Error("Failed to save url", slog.String("error", err.Error()))
			return
		}

		log.Info(
			"Url saved",
			slog.Int64("id", id),
			slog.String("url", savingURL.URL),
			slog.String("alias", savingURL.Alias),
		)
		response.Write(w, http.StatusOK, nil, map[string]any{"id": id, "alias": savingURL.Alias})
	}
}
