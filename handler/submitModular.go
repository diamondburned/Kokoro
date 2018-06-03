package handler

import (
	"fmt"
	"net/http"

	"github.com/Gigamons/common/logger"
)

func POSTSubmitModular(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		w.WriteHeader(500)
		logger.Error(err.Error())
		fmt.Println(err)
		return
	}

	//OsuVersion := r.MultipartForm["osuver"]
	//Security := r.MultipartForm["s"]
	//iv := r.MultipartForm["iv"]
	//Score := r.MultipartForm["score"]
	//Password := r.MultipartForm["pass"]

	fmt.Println(r.FormValue("score"))
}
