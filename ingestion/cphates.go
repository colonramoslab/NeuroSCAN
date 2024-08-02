package main

import (
	"database/sql"
	"errors"
	"github.com/charmbracelet/log"
	"strconv"
	"strings"
)

type CphateNode struct {
	id             int
	uid            string
	cphateId       int
	cluster        int
	clusterCount   int
	iteration      int
	iterationCount int
	serial         int
	neurons        []string
}

type Cphate struct {
	id                 int
	timepoint          int
	developmentalStage sql.NullInt64
	filename           string
	fileHash           string
	nodes              []CphateNode
}

// GetCphate gets the CPHATE by timepoint and returns it
func (n *Neuroscan) GetCphate(timepoint int) (Cphate, error) {
	var cphate Cphate

	err := n.connPool.QueryRow(n.context, "SELECT id, timepoint, filename, file_hash FROM cphates WHERE timepoint = $1", timepoint).Scan(&cphate.id, &cphate.timepoint, &cphate.filename, &cphate.fileHash)
	if err != nil {
		return cphate, err
	}

	var cphateNodes []CphateNode
	rows, err := n.connPool.Query(n.context, "SELECT id, uid, cphate_id, cluster, cluster_count, iteration, iteration_count, serial FROM cphate_nodes WHERE cphate_id = $1", cphate.id)

	if err != nil {
		return cphate, err
	}

	defer rows.Close()

	for rows.Next() {
		var cphateNode CphateNode
		err = rows.Scan(&cphateNode.id, &cphateNode.uid, &cphateNode.cphateId, &cphateNode.cluster, &cphateNode.clusterCount, &cphateNode.iteration, &cphateNode.iterationCount, &cphateNode.serial)
		if err != nil {
			return cphate, err
		}

		var nodeNeurons []string

		neurons, err := n.connPool.Query(n.context, "SELECT neuron_id FROM cphate_node_neurons WHERE cphate_node_id = $1", cphateNode.id)

		if err != nil {
			return cphate, err
		}

		defer neurons.Close()

		for neurons.Next() {
			var neuron string
			err = neurons.Scan(&neuron)
			if err != nil {
				return cphate, err
			}
			nodeNeurons = append(nodeNeurons, neuron)
		}

		cphateNodes = append(cphateNodes, cphateNode)
	}

	return cphate, nil
}

// GetCphateNode gets a CPHATE node by UID and cphate ID
func (n *Neuroscan) GetCphateNode(uid string, cphateID int) (CphateNode, error) {
	var cphateNode CphateNode

	err := n.connPool.QueryRow(n.context, "SELECT id, uid, cphate_id, cluster, cluster_count, iteration, iteration_count, serial FROM cphate_nodes WHERE cphate_id = $1 AND uid = $2", cphateID, uid).Scan(&cphateNode.id, &cphateNode.uid, &cphateNode.cphateId, &cphateNode.cluster, &cphateNode.clusterCount, &cphateNode.iteration, &cphateNode.iterationCount, &cphateNode.serial)
	if err != nil {
		return cphateNode, err
	}

	var nodeNeurons []string

	neurons, err := n.connPool.Query(n.context, "SELECT neuron_id FROM cphate_node_neurons WHERE cphate_node_id = $1", cphateNode.id)

	if err != nil {
		return cphateNode, err
	}

	defer neurons.Close()

	for neurons.Next() {
		var neuron string
		err = neurons.Scan(&neuron)
		if err != nil {
			return cphateNode, err
		}
		nodeNeurons = append(nodeNeurons, neuron)
	}

	return cphateNode, nil
}

// writeToDb writes the CPHATE to the database
func (cphate Cphate) writeToDB(n *Neuroscan) {
	exists, err := n.CphateExists(cphate.timepoint)

	if err != nil {
		log.Error("Error checking if CPHATE exists", "err", err)
		return
	}

	if n.skipExisting && exists {
		log.Debug("CPHATE exists, skipping", "timepoint", cphate.timepoint)
		return
	}

	if exists {
		err = n.DeleteCphate(cphate.timepoint)
		if err != nil {
			log.Error("Error deleting existing CPHATE", "err", err)
		}
	}

	err = n.CreateCphate(cphate.timepoint, cphate.developmentalStage, cphate.filename, cphate.fileHash, cphate.nodes)

	if err != nil {
		log.Error("Error creating CPHATE", "err", err)
		return
	}

	n.IncrementCphate()

	log.Info("CPHATE created", "timepoint", cphate.timepoint)
}

// CphateExists checks if a CPHATE exists by timepoint
func (n *Neuroscan) CphateExists(timepoint int) (bool, error) {
	var count int

	err := n.connPool.QueryRow(n.context, "SELECT COUNT(*) FROM cphates WHERE timepoint = $1", timepoint).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateCphate creates a CPHATE in the database
func (n *Neuroscan) CreateCphate(timepoint int, devStage sql.NullInt64, filename string, fileHash string, nodes []CphateNode) error {
	exists, err := n.CphateExists(timepoint)

	if err != nil {
		return err
	}

	// if the cphate exists, delete it and it's corresponding nodes and the node neurons
	if exists {
		log.Error("CPHATE already exists, cannot create", "timepoint", timepoint)
		return nil
	}

	_, err = n.connPool.Exec(n.context, "INSERT INTO cphates (timepoint, filename, file_hash) VALUES ($1, $2, $3)", timepoint, filename, fileHash)

	if err != nil {
		return err
	}

	newCphate, err := n.GetCphate(timepoint)

	if err != nil {
		return err
	}

	for _, node := range nodes {
		_, err := n.connPool.Exec(n.context, "INSERT INTO cphate_nodes (cphate_id, uid, cluster, cluster_count, iteration, iteration_count, serial) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", newCphate.id, node.uid, node.cluster, node.clusterCount, node.iteration, node.iterationCount, node.serial)

		if err != nil {
			return err
		}

		newNode, err := n.GetCphateNode(node.uid, newCphate.id)

		if err != nil {
			return err
		}

		for _, neuronNode := range node.neurons {
			// we need to get the neuron ID by the neuron UID
			neuron, err := n.GetNeuron(neuronNode, timepoint, devStage)

			if err != nil {
				log.Error("Error getting neuron", "err", err)
				continue
			}

			_, err = n.connPool.Exec(n.context, "INSERT INTO cphate_node_neurons (cphate_node_id, neuron_id) VALUES ($1, $2)", newNode.id, neuron.id)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteCphate deletes a CPHATE by timepoint
func (n *Neuroscan) DeleteCphate(timepoint int) error {
	_, err := n.connPool.Exec(n.context, "DELETE FROM cphate_node_neurons WHERE cphate_node_id IN (SELECT id FROM cphate_nodes WHERE cphate_id = (SELECT id FROM cphates WHERE timepoint = $1))", timepoint)
	if err != nil {
		return err
	}

	_, err = n.connPool.Exec(n.context, "DELETE FROM cphate_nodes WHERE cphate_id = (SELECT id FROM cphates WHERE timepoint = $1)", timepoint)
	if err != nil {
		return err
	}

	_, err = n.connPool.Exec(n.context, "DELETE FROM cphates WHERE timepoint = $1", timepoint)
	if err != nil {
		return err
	}

	return nil
}

// parseCphateNode processes a CPHATE node name string and return a CPHATE node object
// some examples of nodes are:
// ADAL_RIFL_SIBVL_ADAR_SIBVR-i14/26-c8/45-s2178
// ADAL-i11/26-c11/172-s1780
// VB1_ALA_DVA_DVC_SABD_RMED_RID_RIH_IL1L_IL2L_ADAL_AIAL_RIAL_AUAL_AVAL_AWAL_AIBL_RIBL_URBL_AVBL_AWBL_RICL_PVCL_AWCL_IL1DL_IL2DL_SAADL_SIADL_URADL_SIBDL_SMBDL_RMDDL_SMDDL_AFDL_RMDL_CEPDL_OLQDL_AVDL_URYDL_ADEL_RMEL_ASEL_AVEL_ADFL_RIFL_RMFL_AVFL_BAGL_RIGL_RMGL_ASGL_RMHL_ASHL_AVHL_ASIL_ASJL_AVJL_ASKL_AVKL_ADLL_OLLL_AIML_RIML_ALML_CANL_AINL_ALNL_RIPL_FLPL_PVPL_SDQL_PVQL_BDUL_IL1VL_IL2VL_AVL_SAAVL_SIAVL_URAVL_SIBVL_SMBVL_RMDVL_SMDVL_RIVL_CEPVL_OLQVL_URYVL_URXL_AIYL_AIZL_AVM_IL1R_IL2R_ADAR_AIAR_RIAR_AUAR_AVAR_AWAR_AIBR_RIBR_URBR_AVBR_AWBR_RICR_PVCR_AWCR_IL1DR_IL2DR_SAADR_SIADR_URADR_SIBDR_SMBDR_RMDDR_SMDDR_AFDR_RMDR_CEPDR_OLQDR_AVDR_URYDR_ADER_RMER_ASER_AVER_ADFR_RIFR_AVFR_BAGR_RIGR_RMGR_ASGR_RMHR_ASHR_AVHR_RIR_ASIR_ASJR_AVJR_ASKR_AVKR_ADLR_OLLR_AIMR_RIMR_ALMR_CANR_AINR_ALNR_RIPR_FLPR_PVPR_AQR_SDQR_PVQR_BDUR_IL1VR_IL2VR_SAAVR_SIAVR_URAVR_SIBVR_SMBVR_RMDVR_SMDVR_RIVR_PVR_CEPVR_OLQVR_URYVR_URXR_AIYR_AIZR_RIS_PVT_RMEV-i26/26-c1/1-s2351
func parseCphateNode(n *Neuroscan, node string) (CphateNode, error) {
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

	for _, neuron := range neurons {
		cphateNode.neurons = append(cphateNode.neurons, neuron)
	}

	return cphateNode, nil
}

// parseCphate parses the CPHATE file and returns a CPHATE object
func parseCphate(n *Neuroscan, filePath string) (Cphate, error) {
	fileMetas, err := FilePathParse(filePath)

	if err != nil {
		log.Error("Error parsing file path", "err", err)
		return Cphate{}, err
	}

	devStage, err := n.GetDevStageByUID(fileMetas[0].developmentalStage)

	if err != nil {
		log.Error("Failed to get developmental stage by UID", "error", err)
		return Cphate{}, err
	}

	var cphateNodes []CphateNode

	// loop over the fileMetas and create a CPHATE node for each
	for _, fileMeta := range fileMetas {
		cphateNode, err := parseCphateNode(n, fileMeta.uid)

		if err != nil {
			log.Error("Error parsing CPHATE node", "err", err)
			continue
		}

		cphateNodes = append(cphateNodes, cphateNode)
	}

	cphate := Cphate{
		timepoint:          fileMetas[0].timepoint,
		developmentalStage: devStage.id,
		filename:           fileMetas[0].filename,
		fileHash:           fileMetas[0].filehash,
		nodes:              cphateNodes,
	}

	return cphate, nil
}

// ProcessCphate processes the CPHATE file
func ProcessCphate(n *Neuroscan, filePath string) {
	cphate, err := parseCphate(n, filePath)

	if err != nil {
		log.Error("Failed to parse CPHATE", "error", err)
		return
	}

	cphate.writeToDB(n)
}
