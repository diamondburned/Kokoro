package handler

import (
	"fmt"
	"net/http"
)

func POSTosuerror(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.RawPath)
	r.ParseMultipartForm(0)
}
