package test

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"testing"

	"github.com/consensys/go-corset/pkg/util"
)

func Test_HashSet_01(t *testing.T) {
	items := []uint64{1, 2, 3, 4, 3, 2, 1}
	check_HashSet(t, items)
}

func Test_HashSet_02(t *testing.T) {
	items := generateRandomInputs(10, 32)
	check_HashSet(t, items)
}

func Test_HashSet_03(t *testing.T) {
	items := generateRandomInputs(100, 32)
	check_HashSet(t, items)
}

func Test_HashSet_04(t *testing.T) {
	items := generateRandomInputs(1000, 32)
	check_HashSet(t, items)
}

func Test_HashSet_05(t *testing.T) {
	items := generateRandomInputs(100000, 32)
	check_HashSet(t, items)
}

func Test_HashSet_08(t *testing.T) {
	items := generateRandomInputs(100000, 64)
	check_HashSet(t, items)
}

func Test_HashSet_09(t *testing.T) {
	items := generateRandomInputs(100000, 128)
	check_HashSet(t, items)
}

// ===================================================================
// Test Helpers
// ===================================================================

func check_HashSet(t *testing.T, items []uint64) {
	set := util.NewHashSet[testKey](0)
	dups := uint(0)
	// Insert items
	for _, item := range items {
		if set.Insert(testKey{item}) {
			// Duplicate item inserted
			dups++
		}
	}
	// Sort items
	sort.Slice(items, func(i, j int) bool {
		return items[i] < items[j]
	})
	//
	count := uint(0)
	// Count unique items
	for i := 0; i < len(items); i++ {
		if i == 0 || items[i-1] != items[i] {
			count++
		}
	}
	// Sanity check number of unique items
	if set.Size() != count {
		t.Errorf("expected %d unique items, got %d: %s", count, set.Size(), set.String())
	}
	// Sanity check duplicates calculation
	if count+dups != uint(len(items)) {
		t.Errorf("incorrect number of duplicates %d: %s", dups, set.String())
	}
	// Sanity check containership
	for _, ith := range items {
		if !set.Contains(testKey{ith}) {
			t.Errorf("missing item %d: %s", ith, set.String())
		}
	}
}

func generateRandomInputs(n, m uint) []uint64 {
	items := make([]uint64, n)

	for i := uint(0); i < n; i++ {
		items[i] = uint64(rand.UintN(m))
	}

	return items
}

// A simple wrapper around a uint64.  This is deliberately broken to ensure a
// relatively limited spread of hash values.  This helps to ensure that we get
// some collisions.
type testKey struct {
	value uint64
}

// Equals compares two Uint64Keys to check whether they represent the same
// underlying byte array (or not).
func (p testKey) Equals(other testKey) bool {
	return p.value == other.value
}

// Hash generates a 64-bit hashcode from the underlying value.
func (p testKey) Hash() uint64 {
	// This is a deliberate act to limit the qualitfy of this hash function.
	return p.value % 16
}

func (p testKey) String() string {
	return fmt.Sprintf("%d", p.value)
}
