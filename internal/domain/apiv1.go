package domain

type APIV1Request struct {
	Count      bool     `query:"count"`
	Timepoint  *int      `query:"timepoint"`
	UIDs       []string `query:"uids"`
	Types      []string `query:"types"`
	Sort       string   `query:"sort"`
	Limit      int      `query:"limit"`
	Offset     int      `query:"offset"`
	PostNeuron string   `query:"post_neuron"`
	PreNeuron  string   `query:"pre_neuron"`
}
