package define

type FilterFieldSetting struct {
	FieldName string `json:"field_name"`
	Value     any    `json:"value"`
	KeyName   string `json:"key_name"`
}
type SortFieldSetting struct {
	FieldName string  `json:"field_name"`
	Value     float64 `json:"value"`
	KeyName   string  `json:"key_name"`
}

type DataAgreement struct {
	DocId       string                         `json:"doc_id"`
	CreateTime  int64                          `json:"create_time"`
	Refs        map[string]string              `json:"refs"`
	ManyRefs    map[string]string              `json:"many_refs"`
	SortFields  map[string]*SortFieldSetting   `json:"sort_fields"`
	FilterField map[string]*FilterFieldSetting `json:"filter_field"`
	Fields      map[string]any                 `json:"fields"`
}
