package assets

type Asset struct {
	Name string  `json:"name"`
	Path string  `json:"path"`
	Font string  `json:"font"`
	Size float64 `json:"size"`
}

type Assets struct {
	Images []Asset `json:"images"`
	Fonts  []Asset `json:"fonts"`
	Raws   []Asset `json:"raws"`
	Faces  []Asset `json:"faces"`
}
