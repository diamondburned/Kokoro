package api

import (
	"io/ioutil"
	"net/http"

	"github.com/Gigamons/common/logger"
)

func BodyReader(r *http.Request) []byte {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	defer r.Body.Close()
	return b
}
