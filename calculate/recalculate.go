package calculate

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Gigamons/common/helpers"
	"github.com/Gigamons/oppai5"
)

func RecalculateUser(UserID int, scores bool) {
	if scores {
		RecalculateAllScores()
	}
	if UserID == 0 {
		IDs := getallusers()
		for i := 0; i < len(IDs); i++ {
			ID := IDs[i]
			if ID == 100 {
				continue
			}
			CalculateUser(ID, false, 0)
			CalculateUser(ID, false, 1)
			CalculateUser(ID, false, 2)
			CalculateUser(ID, false, 3)
		}
	} else {
		CalculateUser(UserID, false, 0)
		CalculateUser(UserID, false, 1)
		CalculateUser(UserID, false, 2)
		CalculateUser(UserID, false, 3)
	}
	RecalculateUserRX(UserID)
}

func getallusers() []int {
	var r []int
	query := "SELECT id FROM leaderboard"
	x, err := helpers.DB.Query(query)
	if err != nil {
		return r
	}
	for x.Next() {
		var i int
		err := x.Scan(&i)
		if err != nil {
			fmt.Println(err)
		}
		r = append(r, i)
	}
	return r
}

func RecalculateUserRX(UserID int) {
	if UserID == 0 {
		IDs := getallusers()
		for i := 0; i < len(IDs); i++ {
			ID := IDs[i]
			if ID == 100 {
				continue
			}
			CalculateUser(ID, true, 0)
			CalculateUser(ID, true, 1)
			CalculateUser(ID, true, 2)
			CalculateUser(ID, true, 3)
		}
	} else {
		CalculateUser(UserID, true, 0)
		CalculateUser(UserID, true, 1)
		CalculateUser(UserID, true, 2)
		CalculateUser(UserID, true, 3)
	}
}

func RecalculateAllScores() {
	query := `
	SELECT ScoreID, FileMD5,
	(
		SELECT beatmaps.BeatmapID
		FROM beatmaps
		WHERE beatmaps.FileMD5 = scores.FileMD5
	) as BeatmapID, MaxCombo, Count300, Count100, Count50, CountMiss, Mods
	FROM scores WHERE EXISTS
	(
		SELECT beatmaps.RankedStatus
		FROM beatmaps
		WHERE beatmaps.FileMD5 = scores.FileMD5 AND
		(
			RankedStatus = 2 OR RankedStatus = 3 OR RankedStatus = 4
		)
	)`
	rows, err := helpers.DB.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var beatmapid int
		var pp float64
		var filemd5 string
		var count300 uint16
		var count100 uint16
		var count50 uint16
		var countmiss uint16
		var mods uint32
		var maxcombo uint16
		err := rows.Scan(&id, &filemd5, &beatmapid, &maxcombo, &count300, &count100, &count50, &countmiss, &mods)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _, err = os.Stat(fmt.Sprintf("data/map/%s.osu", filemd5)); os.IsNotExist(err) {
			helpers.DownloadBeatmapbyName(strconv.Itoa(beatmapid))
		}
		f, err := os.OpenFile(fmt.Sprintf("data/map/%s.osu", filemd5), os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		Map := oppai.Parse(f)
		PPInfo := oppai.PPInfo(Map, &oppai.Parameters{
			Misses: countmiss,
			Combo:  maxcombo,
			Mods:   mods,
			N100:   count100,
			N300:   count300,
			N50:    count50,
		})
		pp = PPInfo.PP.Total

		_, err = helpers.DB.Exec("UPDATE scores SET PeppyPoints=? WHERE ScoreID=?", pp, id)
		if err != nil {
			fmt.Println(err)
		}
	}
}
