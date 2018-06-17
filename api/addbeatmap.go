package api

import (
	"net/http"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/Gigamons/common/helpers"
)

func AddMap(w http.ResponseWriter, r *http.Request) {
	if !isAllowed(w, r) {
		return
	}

	body := BodyReader(r)
	var s []string
	ffjson.Unmarshal(body, &s)
	sql := "SELECT BeatmapID FROM beatmaps WHERE BeatmapFile = ?"
	helpers.DB.Query(sql)
}
