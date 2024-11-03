package node

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

/*
Add a node with the given ID and return a pointer to it
*/
func (config *NodesYaml) AddNode(nodeID string) (*Node, error) {
	newNode := NewNode(nodeID)
	wwlog.Verbose("Adding new node: %s", nodeID)
	if _, ok := config.Nodes[nodeID]; ok {
		return nil, errors.New("nodename already exists: " + nodeID)
	} else {
		config.Nodes[nodeID] = &newNode
	}
	return &newNode, nil
}

/*
delete node with the given id
*/
func (config *NodesYaml) DelNode(nodeID string) error {
	if _, ok := config.Nodes[nodeID]; !ok {
		return errors.New("nodename does not exist: " + nodeID)
	}

	wwlog.Verbose("Deleting node: %s", nodeID)
	delete(config.Nodes, nodeID)

	return nil
}

/*
set node for the node with id the values of vals
*/
func (config *NodesYaml) SetNode(nodeID string, vals Node) error {
	node, ok := config.Nodes[nodeID]
	if !ok {
		return ErrNotFound
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(vals)
	if err != nil {
		return err
	}
	err = dec.Decode(node)
	return err
}

/*
set profile for the node with id the values of vals
*/
func (config *NodesYaml) SetProfile(profileId string, vals Profile) error {
	profile, ok := config.NodeProfiles[profileId]
	if !ok {
		return ErrNotFound
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(vals)
	if err != nil {
		return err
	}
	err = dec.Decode(profile)
	return err
}

/*
Add a node with the given ID and return a pointer to it
*/
func (config *NodesYaml) AddProfile(profileId string) (*Profile, error) {
	profile := EmptyProfile()
	wwlog.Verbose("adding new profile: %s", profileId)
	if _, ok := config.NodeProfiles[profileId]; ok {
		return nil, errors.New("profile already exists: " + profileId)
	} else {
		config.NodeProfiles[profileId] = &profile
	}
	return &profile, nil
}

/*
delete node with the given id
*/
func (config *NodesYaml) DelProfile(nodeID string) error {
	if _, ok := config.Nodes[nodeID]; !ok {
		return errors.New("profile does not exist: " + nodeID)
	}

	wwlog.Verbose("deleting profile: %s", nodeID)
	delete(config.Nodes, nodeID)

	return nil
}

/*
Write the the NodeYaml to disk.
*/
func (config *NodesYaml) Persist() error {
	out, dumpErr := config.Dump()
	if dumpErr != nil {
		wwlog.Error("%s", dumpErr)
		return dumpErr
	}
	file, err := os.OpenFile(ConfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		wwlog.Error("%s", err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(string(out))
	if err != nil {
		return err
	}
	wwlog.Debug("persisted: %s", ConfigFile)
	return nil
}

/*
Dump returns a YAML document representing the nodeDb
instance. Passes through any errors generated by yaml.Marshal.
*/
func (config *NodesYaml) Dump() ([]byte, error) {
	// flatten out profiles and nodes
	for _, val := range config.NodeProfiles {
		val.Flatten()
	}
	for _, val := range config.Nodes {
		val.Flatten()
	}
	var buf bytes.Buffer
	// Run through encoder
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(config)
	return buf.Bytes(), err //yaml.Marshal(config)
}
