package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gigamons/common/logger"

	"github.com/Gigamons/Kokoro/helper"
)

func SearchDirect(w http.ResponseWriter, r *http.Request) {
	RankedStatus := r.URL.Query().Get("r")
	Query := r.URL.Query().Get("q")
	Page := r.URL.Query().Get("p")
	Mode := r.URL.Query().Get("m")
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")

	if len(RankedStatus) < 1 {
		return
	}
	if len(Page) < 1 {
		return
	}
	if len(Mode) < 1 {
		return
	}

	if b, err := GetCache(RankedStatus + Page + Mode + Query); len(b) > 0 {
		if err != nil {
			logger.Errorln(err)
		} else {
			w.Write(b)
			return
		}
	}

	Query = strings.Trim(Query, " ")

	var err error
	var rs int
	var p int
	var m int
	if rs, err = strconv.Atoi(RankedStatus); err != nil {
		return
	}
	if p, err = strconv.Atoi(Page); err != nil {
		return
	}
	if m, err = strconv.Atoi(Mode); err != nil {
		return
	}

	_Cheese := helper.CheeseGull{Query: Query, PlayMode: int8(m), RankedStatus: int8(rs), Page: int32(p)}

	out := _Cheese.ToDirect()
	fmt.Fprint(w, out)
	SetCache(RankedStatus+Page+Mode+Query, []byte(out), 60)
}

func GETDirectSet(w http.ResponseWriter, r *http.Request) {
	SetID := r.URL.Query().Get("s")
	BeatmapID := r.URL.Query().Get("b")

	if b, err := GetCache(SetID + BeatmapID); len(b) > 0 {
		if err != nil {
			logger.Errorln(err)
		} else {
			w.Write(b)
			return
		}
	}

	if SetID == "" {
		SetID = "0"
	}

	if BeatmapID == "" {
		BeatmapID = "0"
	}

	sid, err := strconv.Atoi(SetID)
	if err != nil {
		return
	}
	bid, err := strconv.Atoi(BeatmapID)
	if err != nil {
		return
	}

	_Cheese := helper.CheeseGull{}

	out := _Cheese.ToNP(sid, bid)

	SetCache(BeatmapID+SetID, []byte(out), 60)

	if out == "0" {
		return
	}
	fmt.Fprint(w, out)
}
