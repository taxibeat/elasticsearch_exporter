package collector

// IndicesSettingsResponse is a representation of Elasticsearch Indices Settings
type IndicesSettingsResponse map[string]Index

// Index is a representation of Elasticsearch Index Settings
type Index struct {
	Settings Settings `json:"settings"`
}

type Settings struct {
	IndexInfo IndexInfo `json:"index"`
}

type IndexInfo struct {
	Blocks Blocks `json:"blocks"`
}

type Blocks struct {
	ReadOnly string `json:"read_only_allow_delete"`
}
