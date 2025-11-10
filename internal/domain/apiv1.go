package domain

type APIV1Request struct {
	Count      bool     `query:"count"`
	Timepoint  *int     `query:"timepoint" param:"timepoint"`
	ULID       string   `query:"ulid" param:"ulid"`
	UID        string   `query:"uid" param:"uid"`
	UIDs       []string `query:"uid"`
	Types      []string `query:"type"`
	Sort       string   `query:"sort"`
	Limit      int      `query:"limit"`
	Offset     int      `query:"start"`
	PostNeuron string   `query:"post_neuron"`
	PreNeuron  string   `query:"pre_neuron"`
}
