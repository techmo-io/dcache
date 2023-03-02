package main

import (
	"testing"
)

func TestCluster_HashLookup(t *testing.T) {
	config := ClusterConfig{
		HashSpaceSize:     20,
		HashPointsPerNode: 2,
	}

	c := NewCluster(config)

	c.AddNode(&Node{
		UUID:       "a",
		HashPoints: []int{5, 13},
	})
	c.AddNode(&Node{
		UUID:       "b",
		HashPoints: []int{8, 18},
	})

	expected := []string{"a", "a", "a", "a", "a", "a", "b", "b", "b", "a", "a", "a", "a", "a", "b", "b", "b", "b", "b", "a"}
	if !slicesAreEqual(c.HashLookup, expected) {
		t.Errorf("calculated hash lookup != expected. %v != %v", c.HashLookup, expected)
	}
}

func TestCluster_HashLookupIsValid(t *testing.T) {
	config := ClusterConfig{
		HashSpaceSize:     20,
		HashPointsPerNode: 2,
	}

	c := NewCluster(config)

	c.AddNode(&Node{
		UUID:       "a",
		HashPoints: nil,
	})
	c.AddNode(&Node{
		UUID:       "b",
		HashPoints: nil,
	})

	if len(c.HashLookup) != config.HashSpaceSize {
		t.Errorf("lookup has invalid size.  got != want. %v != %v", len(c.HashLookup), c.Config.HashSpaceSize)
	}

	var hasA, hasB bool
	for _, UUID := range c.HashLookup {
		if UUID != "a" && UUID != "b" {
			t.Errorf("lookup has invalid Nodes UUID = %v. expected a or b", UUID)
		}
		if UUID == "a" {
			hasA = true
		}
		if UUID == "b" {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Errorf("lookup does not include all nodes, hasA=%v, hasB=%v", hasA, hasB)
	}

}

func slicesAreEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
