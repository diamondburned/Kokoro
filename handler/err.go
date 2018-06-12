package handler

import (
	"net/http"

	"github.com/Gigamons/common/logger"
)

func POSTosuerror(w http.ResponseWriter, r *http.Request) {
	logger.Debugln(r.URL.RawPath)
	r.ParseMultipartForm(0)
}
