package helper

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	json "github.com/pquerna/ffjson/ffjson"

	"github.com/Gigamons/common/logger"
)

func CheeseStatus(rankedStatus int8) int8 {
	switch rankedStatus {
	case 0:
		return 1
	case 2:
		return 0
	case 3:
		return 3
	case 4:
		return -100
	case 5:
		return -2
	case 7:
		return 2
	case 8:
		return 4
	default:
		return 1
	}
}

type CheeseGull struct {
	RankedStatus  int8
	Query         string
	Page          int32
	PlayMode      int8
	_useSet       bool
	_beatmapID    int
	_beatmapSetID int
	Beatmap       []*Beatmap
}

func (c *CheeseGull) _search() {
	if strings.Contains(strings.ToLower(c.Query), "newest") || strings.Contains(strings.ToLower(c.Query), "top rated") || strings.Contains(strings.ToLower(c.Query), "most played") {
		c.Query = ""
	}
	rankedStatus := CheeseStatus(c.RankedStatus)
	rstatus := ""
	if rankedStatus < -100 {
		rstatus = ""
	} else {
		rstatus = strconv.Itoa(int(rankedStatus))
	}

	query := fmt.Sprintf("?mode=%s&amount=%v&offset=%v&status=%s&query=%s", func() string {
		fmt.Println(c.PlayMode)
		if c.PlayMode < 0 || c.PlayMode > 3 {
			return ""
		}
		return strconv.Itoa(int(c.PlayMode))
	}(), 100, c.Page*100, rstatus, url.QueryEscape(c.Query))

	fmt.Println(query)

	api, err := http.Get(os.Getenv("CHEESEGULL") + "/search" + query)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if api == nil {
		logger.Errorln("URL not Valid!")
		return
	}
	defer api.Body.Close()
	body, err := ioutil.ReadAll(api.Body)
	if err != nil {
		logger.Errorln(err)
		return
	}
	CheeseGullAnswer := []*Beatmap{}
	if err = json.Unmarshal(body, &CheeseGullAnswer); err != nil {
		logger.Errorln(err)
		return
	}
	c.Beatmap = CheeseGullAnswer
}

func (c *CheeseGull) GetSet(SetID int) *Beatmap {
	api, err := http.Get(os.Getenv("CHEESEGULL") + "/s/" + strconv.Itoa(SetID))
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	if api == nil {
		logger.Errorln("URL not Valid!")
		return nil
	}
	defer api.Body.Close()
	body, err := ioutil.ReadAll(api.Body)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	CheeseGullAnswer := NewBeatmap()
	if err = json.Unmarshal(body, CheeseGullAnswer); err != nil {
		logger.Errorln(err)
		return nil
	}

	return CheeseGullAnswer
}

func (c *CheeseGull) GetBeatmap(BeatmapID int) *ChildrenBeatmaps {
	api, err := http.Get(os.Getenv("CHEESEGULL") + "/b/" + strconv.Itoa(BeatmapID))
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	if api == nil {
		logger.Errorln("URL not Valid!")
		return nil
	}
	defer api.Body.Close()
	body, err := ioutil.ReadAll(api.Body)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	CheeseGullAnswer := NewChildrenBeatmap()
	if err = json.Unmarshal(body, &CheeseGullAnswer); err != nil {
		logger.Errorln(err)
		return nil
	}

	return CheeseGullAnswer
}

func (c *CheeseGull) GetBeatmapByHash(FileMD5 string) *ChildrenBeatmaps {
	api, err := http.Get(os.Getenv("CHEESEGULL") + "/hash/" + FileMD5)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	if api == nil {
		logger.Errorln("URL not Valid!")
		return nil
	}
	defer api.Body.Close()
	body, err := ioutil.ReadAll(api.Body)
	if err != nil {
		logger.Errorln(err)
		return nil
	}
	CheeseGullAnswer := NewChildrenBeatmap()
	if err = json.Unmarshal(body, &CheeseGullAnswer); err != nil {
		logger.Errorln(err)
		return nil
	}

	return CheeseGullAnswer
}

func (c *CheeseGull) ToDirect() string {
	c._search()
	OutputString := ""

	if len(c.Beatmap) >= 100 {
		OutputString += "101"
	} else {
		OutputString += strconv.Itoa(len(c.Beatmap))
	}

	OutputString += "\n"

	logger.Debugln(OutputString)

	if len(c.Beatmap) > 0 {
		for i := 0; i < len(c.Beatmap); i++ {
			BMSet := c.Beatmap[i]
			MaxDiff := 0.0
			for x := 0; x < len(BMSet.ChildrenBeatmaps); x++ {
				BM := BMSet.ChildrenBeatmaps[x]
				if BM.DifficultyRating > MaxDiff {
					MaxDiff = BM.DifficultyRating
				}
			}
			MaxDiff += 3
			// BeatmapSetID | Artist | Title | Creator | RankedStatus | MaxDiff | LastUpdate | SetID | SetID | HasVideo | 0 | 1234 | VideoLength
			OutputString += fmt.Sprintf("%v.osz|%s|%s|%s|%v|%.2f|%s|%v|%v|%v|0|1234|%s|",
				BMSet.SetID,
				BMSet.Artist,
				BMSet.Title,
				BMSet.Creator,
				BMSet.RankedStatus,
				MaxDiff,
				BMSet.LastUpdate.Format("2006-01-02T15:04:05Z"),
				BMSet.SetID,
				BMSet.SetID,
				func() int {
					if BMSet.HasVideo {
						return 1
					}
					return 0
				}(),
				func() string {
					if BMSet.HasVideo {
						return "4321"
					}
					return "0"
				}(),
			)
			for x := 0; x < len(BMSet.ChildrenBeatmaps); x++ {
				BM := BMSet.ChildrenBeatmaps[x]
				OutputString += fmt.Sprintf("%s (%.2f★~%v♫~AR%v~OD%v~CS%v~HP%v~%vm%vs)@%v,",
					strings.Replace(BM.DiffName, "@", "", -1),
					BM.DifficultyRating,
					BM.BPM,
					BM.AR,
					BM.OD,
					BM.CS,
					BM.HP,
					math.Floor(float64(BM.TotalLength)/float64(60)),
					BM.TotalLength%60,
					BM.Mode,
				)
			}
			if last := len(OutputString) - 1; last >= 0 && OutputString[last] == ',' {
				OutputString = OutputString[:last] + "|"
			}
			OutputString += "\n"
		}
	}

	if len(c.Beatmap) <= 0 {
		OutputString = "-1\nNo beatmaps found!"
	}
	if len(c.Beatmap) <= 0 && c.Query == "" {
		OutputString = "-1\nWhoops, looks like osu!direct is down!"
	}
	return OutputString
}

func (c *CheeseGull) ToNP(SetID int, BeatmapID int) string {
	Set := c.GetSet(SetID)
	Beatmap := c.GetBeatmap(BeatmapID)
	OutputString := ""

	if Set == nil && Beatmap == nil {
		return "0"
	}

	if Beatmap != nil {
		Set = c.GetSet(int(Beatmap.ParentSetID))
	}

	if Set != nil {
		OutputString = fmt.Sprintf("%v.osz|%s|%s|%s|%v|10.00|%s|%v|%v|%s|0|1234|%s\r\n",
			Set.SetID,
			Set.Artist,
			Set.Title,
			Set.Creator,
			Set.RankedStatus,
			Set.LastUpdate,
			Set.SetID,
			Set.SetID,
			func() string {
				if Set.HasVideo {
					return "1"
				}
				return "0"
			}(),
			func() string {
				if Set.HasVideo {
					return "4321"
				}
				return ""
			}(),
		)
	}
	return OutputString
}
