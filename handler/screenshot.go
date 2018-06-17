package handler

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Gigamons/common/helpers"

	"github.com/Gigamons/common/tools/usertools"
	"github.com/gorilla/mux"
)

func GETScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ScreenShot := vars["screenshot"]
	if helpers.NotExists("data/screenshots/" + ScreenShot) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404, Screenshot not found!")
	} else {
		http.ServeFile(w, r, "data/screenshots/"+ScreenShot)
	}
}

func POSTScreenshot(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	UserName := r.FormValue("u")
	Password := r.FormValue("p")
	if usertools.GetUser(usertools.GetUserID(UserName)).CheckPassword(Password) {
		Screenshot, _, err := r.FormFile("ss")
		if err != nil {
			return
		}
		defer Screenshot.Close()
		f, err := ioutil.ReadAll(Screenshot)
		if err != nil {
			return
		}
		ha, err := helpers.MD5(f)
		if err != nil {
			return
		}
		ioutil.WriteFile("data/screenshots/"+hex.EncodeToString(ha)[:8], f, 0644)
		w.Write([]byte("https://osu.gigamons.de/ss/" + hex.EncodeToString(ha)[:8]))
	} else {
		w.Write([]byte("No"))
	}
}
