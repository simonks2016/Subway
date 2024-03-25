package Sorter

type insideSortStruct struct {
	DataId string  `json:"data_id"`
	Score  float64 `json:"score"`
}

type AscendingAlgorithm []insideSortStruct

func (s AscendingAlgorithm) Len() int {
	return len(s)
}
func (s AscendingAlgorithm) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}
func (s AscendingAlgorithm) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
