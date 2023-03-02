package main

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
)

const DefaultHashSpaceSize = 10000
const DefaultHashPointsPerNode = 4

type Host struct {
	Name           string `json:"name"`
	Port           string `json:"port"`
	ControllerPort string `json:"controllerPort"`
}

func (h Host) string() string {
	return fmt.Sprintf("%s:%s", h.Name, h.Port)
}

type Node struct {
	UUID       string `json:"uuid"`
	HashPoints []int  `json:"hashPoints"`
	Host       Host   `json:"host"`
}

func NewNode(host Host) *Node {
	return &Node{
		Host: host,
	}
}

type Cluster struct {
	Nodes      []*Node       `json:"nodes"`
	HashLookup []string      `json:"hashLookup"` // point to Nodes uuid for all hash key space
	Config     ClusterConfig `json:"config"`
}

type ClusterConfig struct {
	HashSpaceSize     int
	HashPointsPerNode int
}

var DefaultClusterConfig = ClusterConfig{
	HashSpaceSize:     100,
	HashPointsPerNode: 3,
}

func NewCluster(config ClusterConfig) *Cluster {
	c := new(Cluster)
	c.Config = config

	return c
}

func (c *Cluster) AddNode(n *Node) {
	if n.UUID == "" {
		n.UUID = uuid.NewString()
	}
	c.Nodes = append(c.Nodes, n)
	c.generatePointsFor(n)
	c.generateHashLookup()
}

func (c *Cluster) generatePointsFor(n *Node) {
	existingPoints := c.getAllHashPoints()

	for len(n.HashPoints) < c.Config.HashPointsPerNode {
		p := rand.Intn(c.Config.HashSpaceSize)
		if _, exists := existingPoints[p]; exists {
			continue
		}
		n.HashPoints = append(n.HashPoints, p)
		existingPoints[p] = ""
	}
}

func (c *Cluster) getAllHashPoints() map[int]string {
	m := make(map[int]string)

	for _, n := range c.Nodes {
		for _, p := range n.HashPoints {
			m[p] = n.UUID
		}
	}
	return m
}

func (c *Cluster) generateHashLookup() {
	points := c.getAllHashPoints()
	lookup := make([]string, c.Config.HashSpaceSize)
	for p, UUID := range points {
		lookup[p] = UUID
	}

	lastEncounteredUUID := ""
	// second round takes care of the highest range
	for round := 0; round < 2; round++ {
		i := c.Config.HashSpaceSize - 1
		for i >= 0 {
			if u := lookup[i]; u == "" {
				lookup[i] = lastEncounteredUUID
			} else {
				lastEncounteredUUID = lookup[i]
			}
			i--
		}
	}

	c.HashLookup = lookup
}
