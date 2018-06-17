package handler

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gigamons/Kokoro/helper"
	"github.com/Gigamons/common/consts"
	"github.com/Gigamons/common/logger"

	"github.com/Gigamons/common/helpers"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

type Achievement struct {
	ID          int32
	Name        string
	DisplayName string
	Description string
	Icon        []byte
}

func (a *Achievement) ToString() string {
	return strings.Replace(a.Name+"+"+a.DisplayName+"+"+a.Description, "\n", "\t", -1)
}

func (a *Achievement) FileIcon(x2 bool) []byte {
	if a.Icon == nil {
		return nil
	}
	f := bytes.NewBuffer(a.Icon)
	x := bytes.NewBuffer(nil)
	img, err := png.Decode(f)
	if err == nil {
		if x2 {
			img = resize.Resize(385, 420, img, resize.Bicubic)
		} else {
			img = resize.Resize(193, 210, img, resize.Bicubic)
		}
		png.Encode(x, img)
		return x.Bytes()
	}
	f = bytes.NewBuffer(a.Icon)
	img, err = jpeg.Decode(f)
	if err == nil {
		if x2 {
			img = resize.Resize(385, 420, img, resize.Bicubic)
		} else {
			img = resize.Resize(193, 210, img, resize.Bicubic)
		}
		jpeg.Encode(x, img, &jpeg.Options{Quality: 10})
		return x.Bytes()
	}
	f = bytes.NewBuffer(a.Icon)
	img, err = gif.Decode(f)
	if err == nil {
		if x2 {
			img = resize.Resize(385, 420, img, resize.Bicubic)
		} else {
			img = resize.Resize(193, 210, img, resize.Bicubic)
		}
		gif.Encode(x, img, &gif.Options{})
		return x.Bytes()
	}
	logger.Errorln("Unknown file type!")
	return nil
}

func GetAchievement(name string) *Achievement {
	a := &Achievement{}
	sql := "SELECT Name, DisplayName, Description, Icon, BitID FROM achievements WHERE Name = ?"
	helpers.DB.QueryRow(sql, name).Scan(&a.Name, &a.DisplayName, &a.Description, &a.Icon, &a.ID)
	return a
}

func (a *Achievement) AddToUser(u *consts.User) {
	sql := "UPDATE users SET achievements = achievements + ? WHERE id = ? AND (achievements & ? = 0)"
	helpers.DB.Exec(sql, a.ID, u.ID, a.ID)
}

func (a *Achievement) DisplayToUser(u *consts.User) {
	sql := "UPDATE users SET achievements_displayed = achievements_displayed + ? WHERE id = ? AND (achievements_displayed & ? = 0)"
	helpers.DB.Exec(sql, a.ID, u.ID, a.ID)
}

func AchievementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	achievement := vars["achievement"]
	x2 := strings.HasSuffix(achievement, "@2x")
	achievement = strings.Replace(achievement, "@2x", "", -1)
	if c, err := GetCache(achievement + strconv.FormatBool(x2)); err != nil && c != nil {
		w.Header().Set("Content-Type", "image/gif")
		w.Write(c)
		return
	}
	a := GetAchievement(achievement)
	w.Header().Set("Content-Type", "image/gif")
	o := a.FileIcon(x2)
	w.Write(o)
	SetCache(achievement+strconv.FormatBool(x2), o, 604800)
}

func ClaimAchievement(u *consts.User, sdata *scoredata, bm *helper.DBBeatmap) string {
	o := ""

	if u.AchievementsDisplayed&4 == 0 {
		arch := GetAchievement("welcome-to-gigamons")
		o = March(o, arch.ToString())
		arch.AddToUser(u)
		arch.DisplayToUser(u)
	}

	if u.AchievementsDisplayed&8 == 0 && strings.HasPrefix(strings.ToLower(bm.Title), "history maker") {
		arch := GetAchievement("history-maker")
		o = March(o, arch.ToString())
		arch.AddToUser(u)
		arch.DisplayToUser(u)
	}

	if u.AchievementsDisplayed&16 == 0 && strings.HasPrefix(strings.ToLower(bm.Title), "super nuko world") {
		arch := GetAchievement("nya-nya-nya")
		o = March(o, arch.ToString())
		arch.AddToUser(u)
		arch.DisplayToUser(u)
	}

	return o
}

func March(a ...string) string {
	o := ""
	for i := 0; i < len(a); i++ {
		if a[i] == "" {
			continue
		}
		o += a[i] + "/"
	}
	o = strings.Trim(o, "/")
	return o
}
