package models

type Contact struct {
	ID        *int    `json:"id"`
	UID       *string `json:"uid"`
	Weight    *int    `json:"weight"`
	NeuronA   *Neuron `json:"neuron_a"`
	NeuronB   *Neuron `json:"neuron_b"`
	Filename  *string `json:"filename"`
	FileHash  *string `json:"file_hash"`
	Timepoint *int    `json:"timepoint"`
}
