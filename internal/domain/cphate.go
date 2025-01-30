package domain

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"neuroscan/internal/toolshed"
)

const CphateULIDPrefix = "cph"
const CphateMetaItemULIDPrefix = "cph_mi"

type CphateNode struct {
	id             int
	ulid string
	uid            string
	cphateId       int
	cluster        int
	clusterCount   int
	iteration      int
	iterationCount int
	serial         int
	neurons        []string
}

type CphateMeta []CphateMetaItem

type CphateMetaItem struct {
	I       int            `json:"i"`
	C       int            `json:"c"`
	Neurons []string       `json:"neurons"`
	ObjFile string         `json:"objFile"`
	Color   toolshed.Color `json:"color"`
	ULID 	string         `json:"ulid"`
}

type Cphate struct {
	ID        int        `json:"-"`
	ULID      string     `json:"id"`
	UID       string     `json:"uid"`
	Timepoint int        `json:"timepoint"`
	Structure CphateMeta `json:"structure"`
}

func (c *Cphate) Parse(dirPath string) error {

	timepoint, err := toolshed.GetTimepoint(dirPath)
	if err != nil {
		return errors.New("error getting timepoint: " + err.Error())
	}

	c.Timepoint = timepoint

	timepointString := strconv.Itoa(timepoint)

	ulid := toolshed.CreateULID(CphateULIDPrefix)
	c.ULID = ulid
	c.UID = "CPHATE " + timepointString

	var cphateMetaItems []CphateMetaItem

	err = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return errors.New("error walking directory: " + err.Error())
		}

		if d.IsDir() {
			return nil
		}

		if !toolshed.ValidExtension(path, []string{".gltf"}) {
			return nil
		}

		fileMetas, err := toolshed.FilePathParse(path)

		if err != nil {
			return errors.New("error parsing file path: " + err.Error())
		}

		// loop over the fileMetas and create a CPHATE node for each
		for _, fileMeta := range fileMetas {
			//cphateNode, err := parseCphateNode(n, fileMeta.uid)
			cphateMetaItem, err := buildCphateMetaItem(fileMeta.UID, fileMeta.Filename, fileMeta.Color)

			if err != nil {
				continue
			}

			cphateMetaItems = append(cphateMetaItems, cphateMetaItem)

			//cphateNodes = append(cphateNodes, cphateNode)
		}

		return nil
	})

	if err != nil {
		return errors.New("error walking directory: " + err.Error())
	}

	c.Structure = cphateMetaItems

	return nil
}

func buildCphateMetaItem(node string, filename string, color toolshed.Color) (CphateMetaItem, error) {
	// split the node string by the "-" character
	// the first part is the neuron names
	// the second part is the cluster, iteration, and serial
	// if we don't have a "-" character, then return an error
	parts := strings.SplitN(node, "-", 2)

	if len(parts) != 2 {
		return CphateMetaItem{}, errors.New("invalid cphate node")
	}

	// split the first part by the "_" character
	// this will give us the neuron names
	neurons := strings.Split(parts[0], "_")

	// split the second part by the "-" character
	// this will give us the cluster, iteration, and serial
	clusterParts := strings.Split(parts[1], "-")

	// loop over the cluster parts
	// check if they contain i, c, or s, as it determines if it is an iteration, cluster, or serial

	var cluster, iteration int

	for _, part := range clusterParts {
		if strings.Contains(part, "i") {
			part = strings.Replace(part, "i", "", 1)
			iterationParts := strings.Split(part, "/")
			iteration, _ = strconv.Atoi(iterationParts[0])
		} else if strings.Contains(part, "c") {
			part = strings.Replace(part, "c", "", 1)
			clusterParts := strings.Split(part, "/")
			cluster, _ = strconv.Atoi(clusterParts[0])
		}
	}

	ulid := toolshed.CreateULID(CphateMetaItemULIDPrefix)

	cphateMetaItem := CphateMetaItem{
		I:       iteration,
		C:       cluster,
		Neurons: neurons,
		ObjFile: filename,
		Color:   color,
		ULID:    ulid,
	}

	return cphateMetaItem, nil
}

func parseCphateNode(node string) (CphateNode, error) {
	// split the node string by the "-" character
	// the first part is the neuron names
	// the second part is the cluster, iteration, and serial
	// if we don't have a "-" character, then return an error
	parts := strings.SplitN(node, "-", 2)

	if len(parts) != 2 {
		return CphateNode{}, errors.New("invalid cphate node")
	}

	// split the first part by the "_" character
	// this will give us the neuron names
	neurons := strings.Split(parts[0], "_")

	// split the second part by the "-" character
	// this will give us the cluster, iteration, and serial
	clusterParts := strings.Split(parts[1], "-")

	// loop over the cluster parts
	// check if they contain i, c, or s, as it determines if it is an iteration, cluster, or serial

	var cluster, clusterCount, iteration, iterationCount, serial int

	for _, part := range clusterParts {
		if strings.Contains(part, "i") {
			part = strings.Replace(part, "i", "", 1)
			iterationParts := strings.Split(part, "/")
			iteration, _ = strconv.Atoi(iterationParts[0])

			if len(iterationParts) > 1 {
				iterationCount, _ = strconv.Atoi(iterationParts[1])
			}
		} else if strings.Contains(part, "c") {
			part = strings.Replace(part, "c", "", 1)
			clusterParts := strings.Split(part, "/")
			cluster, _ = strconv.Atoi(clusterParts[0])

			if len(clusterParts) > 1 {
				clusterCount, _ = strconv.Atoi(clusterParts[1])
			}
		} else if strings.Contains(part, "s") {
			part = strings.Replace(part, "s", "", 1)
			serial, _ = strconv.Atoi(part[1:])
		}
	}

	cphateNode := CphateNode{
		uid:            node,
		cluster:        cluster,
		clusterCount:   clusterCount,
		iteration:      iteration,
		iterationCount: iterationCount,
		serial:         serial,
	}

	cphateNode.neurons = append(cphateNode.neurons, neurons...)

	return cphateNode, nil
}
