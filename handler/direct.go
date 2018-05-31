package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gigamons/Kokoro/helper"
)

func SearchDirect(w http.ResponseWriter, r *http.Request) {
	RankedStatus := r.URL.Query()["r"]
	Query := r.URL.Query()["q"]
	Page := r.URL.Query()["p"]
	Mode := r.URL.Query()["m"]
	fmt.Println(r.URL.Query())

	if len(RankedStatus) < 1 {
		return
	}
	if len(Query) < 1 {
		return
	}
	if len(Page) < 1 {
		return
	}
	if len(Mode) < 1 {
		return
	}

	var err error
	var rs int
	var p int
	var m int
	if rs, err = strconv.Atoi(RankedStatus[0]); err != nil {
		return
	}
	if p, err = strconv.Atoi(Page[0]); err != nil {
		return
	}
	if m, err = strconv.Atoi(RankedStatus[0]); err != nil {
		return
	}

	_Cheese := helper.CheeseGull{Query: Query[0], PlayMode: int8(m), RankedStatus: int8(rs), Page: int32(p)}

	w.Write([]byte(_Cheese.ToDirect()))
}
