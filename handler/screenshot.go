package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Gigamons/common/tools/usertools"
	"github.com/gorilla/mux"
)

func GETScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ScreenShot := vars["screenshot"]
	if _, err := os.Stat("data/screenshots/" + ScreenShot); os.IsNotExist(err) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404, Screenshot not found!")
	} else {
		b, err := ioutil.ReadFile("data/screenshots/" + ScreenShot)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "500, Internal Server Error!")
		}
		w.WriteHeader(200)
		w.Write(b)
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
		h := md5.New()
		h.Write(f)
		ha := h.Sum(nil)
		ioutil.WriteFile("data/screenshots/"+hex.EncodeToString(ha)[:8], f, 0644)
		w.Write([]byte("https://gigamons.de/ss/" + hex.EncodeToString(ha)[:8]))
	} else {
		w.Write([]byte("No"))
	}
}
