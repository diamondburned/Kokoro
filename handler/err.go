package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/Gigamons/common/logger"
)

func POSTosuerror(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	logger.Debugln(string(BodyReader(r)))
	logger.Debugln(r.Form)
}

func BodyReader(r *http.Request) []byte {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	defer r.Body.Close()
	return b
}
