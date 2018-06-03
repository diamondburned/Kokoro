package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gigamons/Kokoro/helper"
	"github.com/Gigamons/common/consts"
	"github.com/Gigamons/common/logger"
	"github.com/Gigamons/common/tools/usertools"
)

type Scoreboard struct {
	User              consts.User
	ScoreboardType    int
	ScoreboardVersion int
	Beatmap           *helper.ChildrenBeatmaps
	PlayMode          int
	Mods              int
	Scores            []*Score
}

type Score struct {
	ScoreID   int
	UserID    int
	FileMD5   string
	ScoreMD5  string
	ReplayMD5 string
	Score     int
	MaxCombo  int
	PlayMode  int
	Mods      int
	Count300  int
	Count100  int
	Count50   int
	CountGeki int
	CountKatu int
	CountMiss int
	Date      time.Time
	Accuracy  float64
	PP        float32
}

func GETScoreboard(w http.ResponseWriter, r *http.Request) {
	FileMD5 := r.URL.Query().Get("c")
	pm, err := strconv.Atoi(r.URL.Query().Get("m"))
	if err != nil {
		logger.Debug("Error while parsing mode")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}
	sbt, err := strconv.Atoi(r.URL.Query().Get("v"))
	if err != nil {
		logger.Debug("Error while parsing Scoreboard Type")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}
	sbv, err := strconv.Atoi(r.URL.Query().Get("vv"))
	if err != nil {
		logger.Debug("Error while parsing Scoreboard Version")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}
	mods, err := strconv.Atoi(r.URL.Query().Get("mods"))
	if err != nil {
		logger.Debug("Error while parsing Mods")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}
	UserID := usertools.GetUserID(r.URL.Query().Get("us"))
	if UserID < 0 {
		logger.Debug("User Doesn't exists")
		return
	}
	User := usertools.GetUser(UserID)
	if User == nil {
		logger.Debug("User Doesn't exists")
		return
	}
	if !User.CheckPassword(r.URL.Query().Get("ha")) {
		logger.Debug("User exists, but Password doesn't match.")
		return
	}

	_Cheese := helper.CheeseGull{}
	BM := _Cheese.GetBeatmapByHash(FileMD5)
	if BM == nil {
		logger.Debug("Beatmap not found.")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}
	Set := _Cheese.GetSet(int(BM.ParentSetID))
	if Set == nil {
		logger.Debug("BeatmapSet not found.")
		fmt.Fprintf(w, "%v|false", consts.LatestPending)
		return
	}

	helper.AddBeatmap(Set)

	ScoreBoard := Scoreboard{Beatmap: BM, ScoreboardVersion: sbv, Mods: mods, PlayMode: pm, ScoreboardType: sbt}
	ScoreBoard.DisplayScoreboard(w)
}

func (sb *Scoreboard) DisplayScoreboard(w http.ResponseWriter) {

	fmt.Fprint(w, sb.Beatmap.GetHeader(0))
}
