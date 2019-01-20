package collector

// AllSettingsResponse is a representation of a Elasticsearch Cluster Settings
type AllSettingsResponse map[string]Indexes

// Cluster is a representation of a Elasticsearch Cluster Settings
type Indexes struct {
	Settings Settings `json:"settings"`
}

// Routing is a representation of a Elasticsearch Cluster shard routing configuration
type Settings struct {
	Index Index `json:"index"`
}

// Allocation is a representation of a Elasticsearch Cluster shard routing allocation settings
type Index struct {
	Blocks Blocks `json:"blocks"`
}

// Allocation is a representation of a Elasticsearch Cluster shard routing allocation settings
type Blocks struct {
	ReadOnly string `json:"read_only_allow_delete"`
}
