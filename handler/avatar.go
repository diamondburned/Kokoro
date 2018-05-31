package handler

import (
	"net/http"
	"os"

	identicon "github.com/dgryski/go-identicon"
	"github.com/gorilla/mux"
)

func GETAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	avatar := vars["avatar"]
	if avatar == "" {
		avatar = "\x94\x82\xfc\x21fs9)01y.-cvm"
	}
	w.Header().Set("Content-Type", "image/png")
	if _, err := os.Stat("data/avatars/" + avatar); os.IsNotExist(err) {
		icon := identicon.New7x7([]byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
		w.Write(icon.Render([]byte(avatar)))
	}
}
