package api

import (
	"net/http"

	"github.com/Gigamons/common/logger"

	"github.com/pquerna/ffjson/ffjson"
)

type Err struct {
	StatusCode int
	Message    string
}

func finish(w http.ResponseWriter, StatusCode int, Message string) {
	b, err := ffjson.Marshal(Err{StatusCode: StatusCode, Message: Message})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500"))
		logger.Errorln(err)
		return
	}
	w.WriteHeader(StatusCode)
	w.Write(b)
}
