package handler

import (
	"fmt"
	"net/http"

	"github.com/Gigamons/common/helpers"
	"github.com/Gigamons/common/logger"

	"github.com/Gigamons/common/tools/usertools"
)

func GETreplaycompressed(w http.ResponseWriter, r *http.Request) {
	validlogin := usertools.GetUser(usertools.GetUserID(r.URL.Query().Get("u"))).CheckPassword(r.URL.Query().Get("h"))
	if !validlogin {
		fmt.Fprint(w, "error: pass")
		return
	}
	ScoreID := r.URL.Query().Get("c")

	Query := "SELECT (SELECT Replay FROM replays WHERE scores.ReplayMD5 = ReplayMD5) FROM scores WHERE ScoreID = ?"
	var Replay []byte
	Row := helpers.DB.QueryRow(Query, ScoreID)
	err := Row.Scan(&Replay)
	if err != nil {
		logger.Errorln(err)
		return
	}
	w.Write(Replay)
}
