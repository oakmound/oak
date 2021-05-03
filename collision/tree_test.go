package collision

import (
	"math/rand"
	"testing"

	"github.com/oakmound/oak/v2/alg/floatgeom"
	"github.com/oakmound/oak/v2/alg/range/floatrange"
)

func TestNewTreeInvalidChildren(t *testing.T) {
	tree, err := NewCustomTree(20, 10)
	if err == nil {
		t.Fatalf("new tree with min children > max should have failed")
	}
	if tree != nil {
		t.Fatalf("new tree with min children > max should not have returned a tree")
	}
}

func TestTreeScene(t *testing.T) {
	tree, err := NewCustomTree(10, 20)
	if err != nil {
		t.Fatalf("unexpected error creating tree: %v", err)
	}

	// Empty tree checks
	if len(tree.Hits(NewSpace(0, 0, 10, 10, 0))) != 0 {
		t.Fatalf("empty tree had collision on new space")
	}
	if tree.HitLabel(NewUnassignedSpace(0, 0, 10, 10)) != nil {
		t.Fatalf("empty tree had collision on new unassigned space")
	}
	if len(tree.Hit(NewLabeledSpace(0, 0, 10, 10, 0), WithLabels(1))) != 0 {
		t.Fatalf("empty tree had collision on new labeled space")
	}

	// Positive hit checks
	s1 := NewFullSpace(0, 0, 10, 10, 1, 3)
	s2 := NewFullSpace(10, 10, 20, 20, 2, 4)
	tree.Add(s1, s2)
	if tree.Size() != 2 {
		t.Fatalf("tree with two additions did not have size 2")
	}
	if len(tree.Hits(NewSpace(5, 5, 1, 1, 0))) != 1 {
		t.Fatalf("Hits did not collide with s1")
	}
	if tree.HitLabel(NewSpace(15, 15, 1, 1, 0), 2) == nil {
		t.Fatalf("HitLabel did not collide with s2")
	}
	// Self-hit should not happen
	if len(tree.Hits(s1)) != 0 {
		t.Fatalf("tree allowed space to collide with itself")
	}

	// Filters
	if len(tree.Hit(NewSpace(0, 0, 100, 100, 0), WithoutLabels(2))) != 1 {
		t.Fatalf("Filtered Hits (1) did not collide with s1")
	}
	if len(tree.Hit(NewSpace(0, 0, 100, 100, 0), WithLabels(2))) != 1 {
		t.Fatalf("Filtered Hits (2) did not collide with s2")
	}
	if len(tree.Hit(NewSpace(0, 0, 100, 100, 0), WithoutCIDs(3))) != 1 {
		t.Fatalf("Filtered Hits (3) did not collide with s2")
	}
	if len(tree.Hit(NewSpace(0, 0, 100, 100, 0), FirstLabel(1))) != 1 {
		t.Fatalf("Filtered Hits (4) did not collide with s1")
	}
	if len(tree.Hit(NewSpace(0, 0, 100, 100, 0), FirstLabel(5))) != 0 {
		t.Fatalf("Filtered Hits (5) collided despite nonmatching filter")
	}

	// Update functions
	if tree.ShiftSpace(1, 1, s1) != nil {
		t.Fatalf("shift space failed")
	}
	if len(tree.Hits(NewSpace(0, 0, 1, 1, 0))) != 0 {
		t.Fatalf("hit did not hit s1 post shift")
	}
	if tree.UpdateSpace(20, 20, 5, 5, s2) != nil {
		t.Fatalf("update space failed")
	}
	if tree.HitLabel(NewSpace(21, 21, 20, 20, 0), 2) == nil {
		t.Fatalf("hit label did not hit s2 post update (1)")
	}
	if tree.UpdateSpaceRect(NewRect(40, 40, 5, 5), s2) != nil {
		t.Fatalf("update space rect failed")
	}
	if tree.HitLabel(NewSpace(40, 40, 20, 20, 0), 2) == nil {
		t.Fatalf("hit label did not hit s2 post update (2)")
	}

	// Removal, Clear
	if tree.Remove(s2) != 1 {
		t.Fatalf("remove space failed")
	}
	if tree.HitLabel(NewSpace(21, 21, 20, 20, 0), 2) != nil {
		t.Fatalf("hit label hit s2 after s2 was removed")
	}

	tree.Clear()

	if len(tree.Hits(NewSpace(0, 0, 100, 100, 0))) != 0 {
		t.Fatalf("hits post tree clear hit something")
	}
}

func TestUpdateSpaceNilSpace(t *testing.T) {
	if DefaultTree.ShiftSpace(0, 0, nil) == nil {
		t.Fatalf("shift space with nil space should have failed")
	}
	if DefaultTree.UpdateSpace(0, 0, 0, 0, nil) == nil {
		t.Fatalf("update space with nil space should have failed")
	}
	if DefaultTree.UpdateSpaceRect(floatgeom.NewRect3(0, 0, 0, 0, 0, 0), nil) == nil {
		t.Fatalf("update space rect with nil space should have failed")
	}
}

func TestUpdateSpaceNotExists(t *testing.T) {
	tree := NewTree()
	s := NewSpace(0, 0, 100, 100, 0)
	if tree.UpdateSpace(4, 4, 100, 100, s) == nil {
		t.Fatalf("update space should have not existed")
	}
}

func TestTreeStress(t *testing.T) {
	spaces := 100000
	tree, _ := NewCustomTree(3, 6)
	for i := 0; i < spaces; i++ {
		tree.Add(randomSpace())
	}
	if spaces != tree.Size() {
		t.Fatalf("tree did not have all spaces")
	}
	for i := 0; i < spaces/100; i++ {
		x := xRange.Poll()
		y := yRange.Poll()
		z := 0.0
		tree.NearestNeighbors(rand.Intn(10), floatgeom.Point3{x, y, z})
	}
	for i := 0; i < spaces; i++ {
		tree.Delete(randomSpace())
	}
}

func randomSpace() *Space {
	return NewUnassignedSpace(xRange.Poll(), yRange.Poll(), wRange.Poll(), hRange.Poll())
}

var (
	xRange = floatrange.NewLinear(0, 10000)
	yRange = floatrange.NewLinear(0, 10000)
	wRange = floatrange.NewLinear(1, 50)
	hRange = floatrange.NewLinear(1, 50)
)
