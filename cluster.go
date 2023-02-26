package main

import (
	"github.com/google/uuid"
	"math/rand"
)

const DefaultHashSpaceSize = 10000
const DefaultHashPointsPerNode = 4

type Host struct {
	Address string
	Port    int
}

type Node struct {
	UUID       string
	HashPoints []int
	Host       Host
}

type Cluster struct {
	nodes      []*Node
	hashLookup []string // point to node uuid for all hash key space
	config     ClusterConfig
}

type ClusterConfig struct {
	HashSpaceSize     int
	HashPointsPerNode int
}

func NewCluster(config ClusterConfig) *Cluster {
	c := new(Cluster)
	c.config = config

	return c
}

func (c *Cluster) AddNode(n *Node) {
	if n.UUID == "" {
		n.UUID = uuid.NewString()
	}
	c.nodes = append(c.nodes, n)
	c.generatePointsFor(n)
	c.generateHashLookup()
}

func (c *Cluster) generatePointsFor(n *Node) {
	existingPoints := c.getAllHashPoints()

	for len(n.HashPoints) < c.config.HashPointsPerNode {
		p := rand.Intn(c.config.HashSpaceSize)
		if _, exists := existingPoints[p]; exists {
			continue
		}
		n.HashPoints = append(n.HashPoints, p)
		existingPoints[p] = ""
	}
}

func (c *Cluster) getAllHashPoints() map[int]string {
	m := make(map[int]string)

	for _, n := range c.nodes {
		for _, p := range n.HashPoints {
			m[p] = n.UUID
		}
	}
	return m
}

func (c *Cluster) generateHashLookup() {
	points := c.getAllHashPoints()
	lookup := make([]string, c.config.HashSpaceSize)
	for p, UUID := range points {
		lookup[p] = UUID
	}

	lastEncounteredUUID := ""
	// second round takes care of the highest range
	for round := 0; round < 2; round++ {
		i := c.config.HashSpaceSize - 1
		for i >= 0 {
			if u := lookup[i]; u == "" {
				lookup[i] = lastEncounteredUUID
			} else {
				lastEncounteredUUID = lookup[i]
			}
			i--
		}
	}

	c.hashLookup = lookup
}
