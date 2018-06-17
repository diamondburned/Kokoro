package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/Gigamons/Kokoro/api"
	"github.com/Gigamons/Kokoro/handler"
	"github.com/Gigamons/common/logger"
	"github.com/gorilla/mux"
)

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("------------ERROR------------")
				fmt.Println(err)
				fmt.Println("---------ERROR TRACE---------")
				fmt.Println(string(debug.Stack()))
				fmt.Println("----------END ERROR----------")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func unknownWeb() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("Method\t%s Path\t%s\n", r.Method, r.URL.Path)
	})
}

func webFolder(w http.ResponseWriter, r *http.Request) {
	logger.Debugln(r.URL.RawQuery)
	w.Write([]byte("not yet"))
}

// StartServer starts our HTTP Server.
func StartServer(host string, port int16) {
	r := mux.NewRouter()
	r.Use(middleWare)

	r.HandleFunc("/{avatar}", handler.GETAvatar)
	r.HandleFunc("/ss/{screenshot}", handler.GETScreenshot)
	r.HandleFunc("/web/osu-screenshot.php", handler.POSTScreenshot)

	r.HandleFunc("/web/osu-search.php", handler.SearchDirect)
	r.HandleFunc("/web/osu-search-set.php", handler.GETDirectSet)

	r.HandleFunc("/web/osu-submit-modular.php", handler.POSTSubmitModular)
	r.HandleFunc("/web/osu-osz2-getscores.php", handler.GETScoreboard)
	r.HandleFunc("/web/osu-getreplay.php", handler.GETreplaycompressed)

	r.HandleFunc("/web/check-updates.php", handler.GETUpdates)
	r.HandleFunc("/web/osu-error.php", handler.POSTosuerror)
	r.HandleFunc("/web/{web}", webFolder)

	// Achievements
	r.HandleFunc("/images/medals-client/{achievement}.png", handler.AchievementHandler)
	r.HandleFunc("/images/medals-client/{achievement}.gif", handler.AchievementHandler)

	r.HandleFunc("/api/kokoro/addbeatmap", api.AddMap)
	r.NotFoundHandler = unknownWeb()

	logger.Infof("Kokoro is listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%v", host, port), r))
}
