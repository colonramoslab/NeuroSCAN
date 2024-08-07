package models

type Cphate struct {
	ID                 *int          `json:"id"`
	Timepoint          *int          `json:"timepoint"`
	DevelopmentalStage *DevStage     `json:"developmental_stage"`
	Filename           *string       `json:"filename"`
	FileHash           *string       `json:"file_hash"`
	Nodes              []*CphateNode `json:"nodes"`
}

type CphateNode struct {
	ID             *int      `json:"id"`
	UID            *string   `json:"uid"`
	CphateID       *int      `json:"cphate_id"`
	Cluster        *int      `json:"cluster"`
	ClusterCount   *int      `json:"cluster_count"`
	Iteration      *int      `json:"iteration"`
	IterationCount *int      `json:"iteration_count"`
	Serial         *int      `json:"serial"`
	Neurons        []*Neuron `json:"neurons"`
}
