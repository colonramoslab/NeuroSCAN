package models

type Synapse struct {
	ID          *int      `json:"id"`
	UID         *string   `json:"uid"`
	SynapseType *string   `json:"type"`
	Section     *string   `json:"section"`
	Position    *string   `json:"position"`
	NeuronPre   *Neuron   `json:"neuron_pre"`
	PostNeurons []*Neuron `json:"post_neurons"`
	Timepoint   *int      `json:"timepoint"`
	Filename    *string   `json:"filename"`
	FileHash    *string   `json:"file_hash"`
}

type SynapseData struct {
	neuronPre   string
	synapseType string
	section     string
	position    string
	neuronSite  int
	postNeurons []string
}

type SynapsePosition struct {
	section  string
	position string
	site     int
}
