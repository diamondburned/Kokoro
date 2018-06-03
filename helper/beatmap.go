package helper

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
