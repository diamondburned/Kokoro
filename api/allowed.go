package api

import (
	"net/http"

	"github.com/Gigamons/Kokoro/constants"
	"github.com/Gigamons/common/helpers"
)

func isAllowed(w http.ResponseWriter, r *http.Request) bool {
	key := r.URL.Query().Get("key")
	var conf constants.Config
	helpers.GetConfig("config", &conf)
	if key == "" || key != conf.API.APIKey {
		finish(w, 403, "Not Allowed!")
		return false
	}
	return true
}
