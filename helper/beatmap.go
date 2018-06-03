package helper

import (
	"fmt"

	"github.com/Gigamons/common/consts"
	"github.com/Gigamons/common/logger"

	"github.com/Gigamons/common/helpers"
)

type Beatmap struct {
	SetID            int64              `json:"SetID"`
	ChildrenBeatmaps []ChildrenBeatmaps `json:"ChildrenBeatmaps"`
	RankedStatus     int8               `json:"RankedStatus"`
	ApprovedDate     string             `json:"ApprovedDate"`
	LastUpdate       string             `json:"LastUpdate"`
	LastChecked      string             `json:"LastChecked"`
	Artist           string             `json:"Artist"`
	Title            string             `json:"Title"`
	Creator          string             `json:"Creator"`
	Source           string             `json:"Source"`
	Tags             string             `json:"Tags"`
	HasVideo         bool               `json:"HasVideo"`
	Genre            int8               `json:"Genre"`
	Language         int8               `json:"Language"`
	Favourites       int32              `json:"Favourites"`
}

type ChildrenBeatmaps struct {
	BeatmapID        int64   `json:"BeatmapID"`
	ParentSetID      int64   `json:"ParentSetID"`
	DiffName         string  `json:"DiffName"`
	FileMD5          string  `json:"FileMD5"`
	Mode             int8    `json:"Mode"`
	BPM              float32 `json:"BPM"`
	CS               float32 `json:"CS"`
	AR               float32 `json:"AR"`
	OD               float32 `json:"OD"`
	HP               float32 `json:"HP"`
	TotalLength      int32   `json:"TotalLength"`
	HitLength        int32   `json:"HitLength"`
	PlayCount        int64   `json:"PlayCount"`
	PassCount        int64   `json:"PassCount"`
	MaxCombo         int32   `json:"MaxCombo"`
	DifficultyRating float64 `json:"DifficultyRating"`
}

func NewBeatmap() *Beatmap {
	return &Beatmap{}
}

func NewChildrenBeatmap() *ChildrenBeatmaps {
	return &ChildrenBeatmaps{}
}

func BeatmapExists(FileMD5 string, BeatmapID int) (bool, bool) {
	var exists, needupdate bool
	helpers.DB.QueryRow("SELECT * FROM beatmaps WHERE FileMD5 = ?", FileMD5).Scan(&exists)
	helpers.DB.QueryRow("SELECT * FROM beatmaps WHERE BeatmapID = ?", BeatmapID).Scan(&needupdate)
	return exists, needupdate
}

func AddBeatmap(BM *Beatmap) {
	for i := 0; i < len(BM.ChildrenBeatmaps); i++ {
		child := BM.ChildrenBeatmaps[i]
		exists, needupdate := BeatmapExists(child.FileMD5, int(child.BeatmapID))
		if exists {
			continue
		}
		if needupdate {
			UpdateBeatmap(BM)
			continue
		}

		fmt.Println(BM.RankedStatus + 1)

		_, err := helpers.DB.Exec(`
			INSERT INTO beatmaps 
			(
				SetID,
				BeatmapID,
				FileMD5,
				RankedStatus,
				RankedDate,
				Artist,
				Title,
				Creator,
				LastUpdate,
				Difficulty,
				CS,
				OD,
				AR,
				HP,
				BPM,
				HitLength,
				DiffName,
				PlayMode,
				MaxCombo
			)
			VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			BM.SetID,
			child.BeatmapID,
			child.FileMD5,
			FixRankedStatus(BM.RankedStatus),
			BM.ApprovedDate,
			BM.Artist,
			BM.Title,
			BM.Creator,
			BM.LastUpdate,
			child.DifficultyRating,
			child.CS, child.OD, child.AR, child.HP,
			child.BPM,
			child.HitLength,
			child.DiffName,
			child.Mode,
			child.MaxCombo,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func UpdateBeatmap(BM *Beatmap) {
	_Cheese := CheeseGull{}
	for i := 0; i < len(BM.ChildrenBeatmaps); i++ {
		child := BM.ChildrenBeatmaps[i]
		if exists, _ := BeatmapExists(child.FileMD5, int(child.BeatmapID)); exists {
			continue
		}
		NewChild := _Cheese.GetBeatmap(int(child.BeatmapID))

		_, err := helpers.DB.Exec(`UPDATE beatmaps SET
			(
				SetID = ?,
				BeatmapID = ?,
				FileMD5 = ?,
				RankedStatus = ?,
				RankedDate = ?,
				Artist = ?,
				Title = ?,
				Creator = ?,
				LastUpdate = ?,
				Difficulty = ?,
				CS = ?,
				OD = ?,
				AR = ?,
				HP = ?,
				BPM = ?,
				HitLength = ?,
				DiffName = ?,
				PlayMode = ?,
				MaxCombo = ?
			) WHERE BeatmapID = ?
			`,
			BM.SetID,
			NewChild.BeatmapID,
			NewChild.FileMD5,
			FixRankedStatus(BM.RankedStatus),
			BM.ApprovedDate,
			BM.Artist,
			BM.Title,
			BM.Creator,
			BM.LastUpdate,
			NewChild.DifficultyRating,
			NewChild.CS, NewChild.OD, NewChild.AR, NewChild.HP,
			NewChild.BPM,
			NewChild.HitLength,
			NewChild.DiffName,
			NewChild.Mode,
			NewChild.MaxCombo,
			NewChild.BeatmapID,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (bm *ChildrenBeatmaps) GetParent() *Beatmap {
	_Cheese := CheeseGull{}
	return _Cheese.GetSet(int(bm.ParentSetID))
}

type DBBeatmap struct {
	SetID        int
	BeatmapID    int
	FileMD5      string
	RankedStatus int
	RankedDate   string
	Artist       string
	Title        string
	Creator      string
	LastUpdate   string
	Difficulty   int
	CS           int
	OD           int
	AR           int
	HP           int
	BPM          int
	HitLength    int
	DiffName     string
	PlayMode     int
	MaxCombo     int
}

func GetBeatmapofDB(BeatmapID int) *DBBeatmap {

	rows, err := helpers.DB.Query(`
		SELECT 
			SetID,
			BeatmapID,
			FileMD5,
			RankedStatus,
			RankedDate,
			Artist,
			Title,
			Creator,
			LastUpdate,
			Difficulty,
			CS,
			OD,
			AR,
			HP,
			BPM,
			HitLength,
			DiffName,
			PlayMode,
			MaxCombo,
		FROM beatmaps WHERE BeatmapID = ?`, BeatmapID)

	if err != nil {
		return nil
	}

	bmDB := &DBBeatmap{}

	for rows.Next() {
		err := rows.Scan(
			bmDB.SetID,
			bmDB.BeatmapID,
			bmDB.FileMD5,
			bmDB.RankedStatus,
			bmDB.RankedDate,
			bmDB.Artist,
			bmDB.Title,
			bmDB.Creator,
			bmDB.LastUpdate,
			bmDB.Difficulty,
			bmDB.CS,
			bmDB.OD,
			bmDB.AR,
			bmDB.HP,
			bmDB.BPM,
			bmDB.HitLength,
			bmDB.DiffName,
			bmDB.PlayMode,
			bmDB.MaxCombo,
		)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	defer rows.Close()

	return bmDB
}

func (bm *ChildrenBeatmaps) GetHeader(TotalScores int) string {
	db := GetBeatmapofDB(int(bm.BeatmapID))
	if db == nil {
		return fmt.Sprintf("%v|false", consts.LatestPending)
	}
	return fmt.Sprintf("%v|true|%v|%v|%v\n%v\n%s\n%v\n", db.RankedStatus, bm.BeatmapID, bm.ParentSetID, TotalScores, 0, db.Title, 10.00)
}

func FixRankedStatus(r int8) int8 {
	out := r
	if r == consts.NotSubmited {
		out = consts.LatestPending
	} else if r == consts.NeedUpdate {
		out = consts.Ranked
	} else if r == consts.Qualified {
		out = consts.Loved
	} else if r == consts.Ranked {
		out = consts.Approved
	} else if r == consts.Unknown {
		out = consts.LatestPending
	}
	return out
}
