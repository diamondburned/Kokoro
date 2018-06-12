package handler

import (
	"net/http"
	"net/url"

	"github.com/Gigamons/common/helpers"
	"github.com/Gigamons/common/logger"
)

func GETUpdates(w http.ResponseWriter, r *http.Request) {
	Action := r.URL.Query().Get("action")
	ReleaseStream := r.URL.Query().Get("stream")

	if c, err := GetCache(Action + "Updater" + ReleaseStream); err == nil && len(c) > 0 {
		w.Write(c)
		return
	}

	uri := url.URL{Host: "osu.ppy.sh", Path: "/web/check-updates.php"}

	b, err := helpers.Download("https:" + uri.String() + "?" + r.URL.RawQuery)
	if err != nil {
		logger.Errorln(err)
		return
	}
	SetCache(Action+"Updater"+ReleaseStream, b, 86400)
	w.Write(b)
}
