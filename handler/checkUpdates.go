package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GETUpdates(w http.ResponseWriter, r *http.Request) {
	Action := r.URL.Query().Get("action")
	ReleaseStream := r.URL.Query().Get("stream")

	if c, err := GetCache(Action + "Updater" + ReleaseStream); err == nil && len(c) > 0 {
		w.Write(c)
		return
	}

	uri := url.URL{Host: "osu.ppy.sh", Path: "/web/check-updates.php"}

	h, err := http.Get("https:" + uri.String() + "?" + r.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer h.Body.Close()
	b, err := ioutil.ReadAll(h.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	SetCache(Action+"Updater"+ReleaseStream, b, 86400)
	w.Write(b)
}
