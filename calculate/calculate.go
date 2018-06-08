package calculate

import (
	"database/sql"
	"fmt"
	"math"
	"sort"

	"github.com/Gigamons/common/consts"
	"github.com/Gigamons/common/helpers"
	"github.com/Gigamons/common/logger"
)

func CalculateUser(UserID int, relaxing bool, playMode int8) {
	var TotalPP float64
	Query := "SELECT PeppyPoints FROM scores WHERE UserID = ? AND PlayMode = ? "
	if relaxing {
		Query += "AND (Mods & 128 > 0) "
	} else {
		Query += "AND (Mods & 128 < 1) "
	}
	Query += "GROUP BY FileMD5 ORDER BY MAX(Score) DESC"
	rows, err := helpers.DB.Query(Query, UserID, playMode)
	if err != nil {
		logger.Error(err.Error())
	}

	pp := rowstoarray(rows)

	sort.Sort(sort.Reverse(sort.Float64Slice(pp)))

	for i := 0; i < len(pp); i++ {
		PeppyPoints := pp[i]
		TotalPP += PeppyPoints * math.Pow(0.95, float64(i))
	}
	m := func() string {
		if relaxing {
			return "_rx"
		} else {
			return ""
		}
	}()
	helpers.DB.Exec("UPDATE leaderboard"+m+" SET pp_"+consts.ToPlaymodeString(playMode)+"= ? WHERE id = ?", TotalPP, UserID)
}

func rowstoarray(r *sql.Rows) []float64 {
	var o []float64
	defer r.Close()
	for r.Next() {
		var i float64
		err := r.Scan(&i)
		if err != nil {
			fmt.Println(err)
		}
		o = append(o, i)
	}
	return o
}
