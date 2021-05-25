package base

type FilePath struct {
	Title  string `json:"title"`
	Parent string `json:"parent"`
	Key    string `json:"key"`
	Leaf   bool   `json:"leaf"`
}
