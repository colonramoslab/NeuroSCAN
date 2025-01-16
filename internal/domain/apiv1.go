package domain

import (
	"fmt"
	"strings"
)

type APIV1Request struct {
	Count      bool     `query:"count"`
	Timepoint  int      `query:"timepoint"`
	UIDs       []string `query:"uids"`
	Types      []string `query:"types"`
	Sort       string   `query:"sort"`
	Limit      int      `query:"limit"`
	Offset     int      `query:"offset"`
	PostNeuron string   `query:"post_neuron"`
	PreNeuron  string   `query:"pre_neuron"`
}

func (r *APIV1Request) ToPostgresQuery() (string, []interface{}) {
	queryParts := []string{"where 1=1"}
	args := []interface{}{}

	args = append(args, r.Timepoint)
	queryParts = append(queryParts, fmt.Sprintf("timepoint = $%d", len(args)))

	if len(r.UIDs) > 0 {
		// we need to build a query where UID is like or, looping over the UIDs, wrapping them in % and adding them to the array[]
		uidArray := []string{}
		for _, uid := range r.UIDs {
			uidArray = append(uidArray, fmt.Sprintf("%%%s%%", uid))
		}
		args = append(args, uidArray)
		queryParts = append(queryParts, fmt.Sprintf("uid ILIKE ANY($%d)", len(args)))
	}

	if len(r.Types) > 0 {
		args = append(args, r.Types)
		queryParts = append(queryParts, fmt.Sprintf("synapse_type IN ($%d)", len(args)))
	}

	// if r.PostNeuron != "" {
	// 	queryParts = append(queryParts, fmt.Sprintf("post_neuron = %s", r.PostNeuron))
	// }

	// if r.PreNeuron != "" {
	// 	queryParts = append(queryParts, fmt.Sprintf("pre_neuron = %s", r.PreNeuron))
	// }

	query := strings.Join(queryParts, " AND ")

	// if count is true, return the query and args before adding the sort and limit
	if r.Count {
		return query, args
	}

	if r.Sort != "" {
		// split by ":", first part is the field, second part is the direction
		parts := strings.Split(r.Sort, ":")

		if len(parts) == 2 {

			// if the second part is not asc or desc, default to asc
			if parts[1] != "asc" && parts[1] != "desc" {
				parts[1] = "asc"
			}

			query += fmt.Sprintf(" order by %s %s", parts[0], parts[1])
		}
	}

	if r.Limit > 0 {
		args = append(args, r.Limit)
		query += fmt.Sprintf(" limit $%d", len(args))
	} else {
		query += " limit 30"
	}

	if r.Offset > 0 {
		args = append(args, r.Offset)
		query += fmt.Sprintf(" offset $%d", len(args))
	}

	return query, args
}
