// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import (
	"fmt"
	"math"
	"sort"

	"github.com/oakmound/oak/alg/floatgeom"
)

// Rtree represents an R-tree, a balanced search tree for storing and querying
// Space objects.  MinChildren/MaxChildren specify the minimum/maximum branching factors.
type Rtree struct {
	MinChildren int
	MaxChildren int
	root        *node
	size        int
	height      int
}

// NewTree creates a new R-tree instance.
func newTree(MinChildren, MaxChildren int) *Rtree {
	rt := Rtree{MinChildren: MinChildren, MaxChildren: MaxChildren}
	rt.height = 1
	rt.root = &node{}
	rt.root.entries = make([]entry, 0, MaxChildren)
	rt.root.leaf = true
	rt.root.level = 1
	return &rt
}

// node represents a tree node of an Rtree.
type node struct {
	parent  *node
	leaf    bool
	entries []entry
	level   int // node depth in the Rtree
}

// entry represents a Space index record stored in a tree node.
type entry struct {
	bb    floatgeom.Rect3 // bounding-box of all children of this entry
	child *node
	obj   *Space
}

// Insertion

// Insert inserts a Space object into the tree.  If insertion
// causes a leaf node to overflow, the tree is rebalanced automatically.
//
// Implemented per Section 3.2 of "R-trees: A Dynamic Index Structure for
// Space Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (tree *Rtree) Insert(obj *Space) {
	e := entry{obj.Location, nil, obj}
	tree.insert(e, 1)
	tree.size++
}

// insert adds the specified entry to the tree at the specified level.
func (tree *Rtree) insert(e entry, level int) {
	leaf := tree.chooseNode(tree.root, e, level)
	leaf.entries = append(leaf.entries, e)

	// update parent pointer if necessary
	if e.child != nil {
		e.child.parent = leaf
	}

	// split leaf if overflows
	var split *node
	if len(leaf.entries) > tree.MaxChildren {
		leaf, split = leaf.split(tree.MinChildren)
	}
	root, splitRoot := tree.adjustTree(leaf, split)
	if splitRoot != nil {
		oldRoot := root
		tree.height++
		tree.root = &node{
			parent: nil,
			level:  tree.height,
			entries: []entry{
				{bb: oldRoot.computeBoundingBox(), child: oldRoot},
				{bb: splitRoot.computeBoundingBox(), child: splitRoot},
			},
		}
		oldRoot.parent = tree.root
		splitRoot.parent = tree.root
	}
}

// chooseNode finds the node at the specified level to which e should be added.
func (tree *Rtree) chooseNode(n *node, e entry, level int) *node {
	if n.leaf || n.level == level {
		return n
	}

	// find the entry whose bb needs least enlargement to include obj
	diff := math.MaxFloat64
	var chosen entry
	var bb floatgeom.Rect3
	for _, en := range n.entries {
		bb = boundingBox(en.bb, e.bb)
		d := bb.Space() - en.bb.Space()
		if d < diff || (d == diff && en.bb.Space() < chosen.bb.Space()) {
			diff = d
			chosen = en
		}
	}

	return tree.chooseNode(chosen.child, e, level)
}

// adjustTree splits overflowing nodes and propagates the changes upwards.
func (tree *Rtree) adjustTree(n, nn *node) (*node, *node) {
	// Let the caller handle root adjustments.
	if n == tree.root {
		return n, nn
	}

	// Re-size the bounding box of n to account for lower-level changes.
	en := n.getEntry()
	en.bb = n.computeBoundingBox()

	// If nn is nil, then we're just propagating changes upwards.
	if nn == nil {
		return tree.adjustTree(n.parent, nil)
	}

	// Otherwise, these are two nodes resulting from a split.
	// n was reused as the "left" node, but we need to add nn to n.parent.
	enn := entry{nn.computeBoundingBox(), nn, nil}
	n.parent.entries = append(n.parent.entries, enn)

	// If the new entry overflows the parent, split the parent and propagate.
	if len(n.parent.entries) > tree.MaxChildren {
		return tree.adjustTree(n.parent.split(tree.MinChildren))
	}

	// Otherwise keep propagating changes upwards.
	return tree.adjustTree(n.parent, nil)
}

// getEntry returns a pointer to the entry for the node n from n's parent.
func (n *node) getEntry() *entry {
	var e *entry
	for i := range n.parent.entries {
		if n.parent.entries[i].child == n {
			e = &n.parent.entries[i]
			break
		}
	}
	return e
}

// computeBoundingBox finds the MBR of the children of n.
func (n *node) computeBoundingBox() (bb floatgeom.Rect3) {
	childBoxes := make([]floatgeom.Rect3, len(n.entries))
	for i, e := range n.entries {
		childBoxes[i] = e.bb
	}
	bb = boundingBoxN(childBoxes...)
	return
}

// split splits a node into two groups while attempting to minimize the
// bounding-box area of the resulting groups.
func (n *node) split(minGroupSize int) (left, right *node) {
	// find the initial split
	l, r := n.pickSeeds()
	leftSeed, rightSeed := n.entries[l], n.entries[r]

	// get the entries to be divided between left and right
	remaining := append(n.entries[:l], n.entries[l+1:r]...)
	remaining = append(remaining, n.entries[r+1:]...)

	// setup the new split nodes, but re-use n as the left node
	left = n
	left.entries = []entry{leftSeed}
	right = &node{
		parent:  n.parent,
		leaf:    n.leaf,
		level:   n.level,
		entries: []entry{rightSeed},
	}

	// TODO
	if rightSeed.child != nil {
		rightSeed.child.parent = right
	}
	if leftSeed.child != nil {
		leftSeed.child.parent = left
	}

	// distribute all of n's old entries into left and right.
	for len(remaining) > 0 {
		next := pickNext(left, right, remaining)
		e := remaining[next]

		if len(remaining)+len(left.entries) <= minGroupSize {
			assign(e, left)
		} else if len(remaining)+len(right.entries) <= minGroupSize {
			assign(e, right)
		} else {
			assignGroup(e, left, right)
		}

		remaining = append(remaining[:next], remaining[next+1:]...)
	}

	return
}

func assign(e entry, group *node) {
	if e.child != nil {
		e.child.parent = group
	}
	group.entries = append(group.entries, e)
}

// assignGroup chooses one of two groups to which a node should be added.
func assignGroup(e entry, left, right *node) {
	leftBB := left.computeBoundingBox()
	rightBB := right.computeBoundingBox()
	leftEnlarged := boundingBox(leftBB, e.bb)
	rightEnlarged := boundingBox(rightBB, e.bb)

	// first, choose the group that needs the least enlargement
	leftDiff := leftEnlarged.Space() - leftBB.Space()
	rightDiff := rightEnlarged.Space() - rightBB.Space()
	if diff := leftDiff - rightDiff; diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group that has smaller area
	if diff := leftBB.Space() - rightBB.Space(); diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group with fewer entries
	if diff := len(left.entries) - len(right.entries); diff <= 0 {
		assign(e, left)
		return
	}
	assign(e, right)
}

// pickSeeds chooses two child entries of n to start a split.
func (n *node) pickSeeds() (int, int) {
	left, right := 0, 1
	maxWastedSpace := -1.0
	for i, e1 := range n.entries {
		for j, e2 := range n.entries[i+1:] {
			d := boundingBox(e1.bb, e2.bb).Space() - e1.bb.Space() - e2.bb.Space()
			if d > maxWastedSpace {
				maxWastedSpace = d
				left, right = i, j+i+1
			}
		}
	}
	return left, right
}

// pickNext chooses an entry to be added to an entry group.
func pickNext(left, right *node, entries []entry) (next int) {
	maxDiff := -1.0
	leftBB := left.computeBoundingBox()
	rightBB := right.computeBoundingBox()
	for i, e := range entries {
		d1 := boundingBox(leftBB, e.bb).Space() - leftBB.Space()
		d2 := boundingBox(rightBB, e.bb).Space() - rightBB.Space()
		d := math.Abs(d1 - d2)
		if d > maxDiff {
			maxDiff = d
			next = i
		}
	}
	return
}

// Deletion

// Delete removes an object from the tree.  If the object is not found, ok
// is false; otherwise ok is true.
//
// Implemented per Section 3.3 of "R-trees: A Dynamic Index Structure for
// Space Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (tree *Rtree) Delete(obj *Space) bool {
	n := tree.findLeaf(tree.root, obj)
	if n == nil {
		return false
	}

	ind := -1
	for i, e := range n.entries {
		if e.obj == obj {
			ind = i
		}
	}
	if ind < 0 {
		return false
	}

	n.entries = append(n.entries[:ind], n.entries[ind+1:]...)

	tree.condenseTree(n)
	tree.size--

	if !tree.root.leaf && len(tree.root.entries) == 1 {
		tree.root = tree.root.entries[0].child
	}

	tree.height = tree.root.level

	return true
}

// findLeaf finds the leaf node containing obj.
func (tree *Rtree) findLeaf(n *node, obj *Space) *node {
	if n.leaf {
		return n
	}
	// if not leaf, search all candidate subtrees
	for _, e := range n.entries {
		if e.bb.ContainsRect(obj.Location) {
			leaf := tree.findLeaf(e.child, obj)
			if leaf == nil {
				continue
			}
			// check if the leaf actually contains the object
			for _, leafEntry := range leaf.entries {
				if leafEntry.obj == obj {
					return leaf
				}
			}
		}
	}
	return nil
}

// condenseTree deletes underflowing nodes and propagates the changes upwards.
func (tree *Rtree) condenseTree(n *node) {
	deleted := []*node{}

	for n != tree.root {
		if len(n.entries) < tree.MinChildren {
			// remove n from parent entries
			entries := []entry{}
			for _, e := range n.parent.entries {
				if e.child != n {
					entries = append(entries, e)
				}
			}
			if len(n.parent.entries) == len(entries) {
				// todo: don't panic
				panic(fmt.Errorf("Failed to remove entry from parent"))
			}
			n.parent.entries = entries

			// only add n to deleted if it still has children
			if len(n.entries) > 0 {
				deleted = append(deleted, n)
			}
		} else {
			// just a child entry deletion, no underflow
			n.getEntry().bb = n.computeBoundingBox()
		}
		n = n.parent
	}

	for _, n := range deleted {
		// reinsert entry so that it will remain at the same level as before
		e := entry{n.computeBoundingBox(), n, nil}
		tree.insert(e, n.level+1)
	}
}

// Searching

// SearchIntersect returns all objects that intersect the specified rectangle.
//
// Implemented per Section 3.1 of "R-trees: A Dynamic Index Structure for
// Space Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (tree *Rtree) SearchIntersect(bb floatgeom.Rect3) []*Space {
	return tree.searchIntersect(tree.root, bb, []*Space{})
}

func (tree *Rtree) searchIntersect(n *node, bb floatgeom.Rect3, results []*Space) []*Space {
	for _, e := range n.entries {
		if e.bb.Intersects(bb) {
			if n.leaf {
				results = append(results, e.obj)
			} else {
				results = tree.searchIntersect(e.child, bb, results)
			}
		}
	}
	return results
}

// NearestNeighbor returns the closest object to the specified point.
// Implemented per "Nearest Neighbor Queries" by Roussopoulos et al
func (tree *Rtree) NearestNeighbor(p floatgeom.Point3) *Space {
	obj, _ := tree.nearestNeighbor(p, tree.root, math.MaxFloat64, nil)
	return obj
}

// utilities for sorting slices of entries

type entrySlice struct {
	entries []entry
	dists   []float64
}

func (s entrySlice) Len() int { return len(s.entries) }

func (s entrySlice) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
	s.dists[i], s.dists[j] = s.dists[j], s.dists[i]
}

func (s entrySlice) Less(i, j int) bool {
	return s.dists[i] < s.dists[j]
}

func sortEntries(p floatgeom.Point3, entries []entry) ([]entry, []float64) {
	sorted := make([]entry, len(entries))
	dists := make([]float64, len(entries))
	for i := 0; i < len(entries); i++ {
		sorted[i] = entries[i]
		dists[i] = minDist(p, entries[i].bb)
	}
	sort.Sort(entrySlice{sorted, dists})
	return sorted, dists
}

func pruneEntries(p floatgeom.Point3, entries []entry, minDists []float64) []entry {
	minMinMaxDist := math.MaxFloat64
	for i := range entries {
		minMaxDist := minMaxDist(p, entries[i].bb)
		if minMaxDist < minMinMaxDist {
			minMinMaxDist = minMaxDist
		}
	}
	// remove all entries with minDist > minMinMaxDist
	pruned := []entry{}
	for i := range entries {
		if minDists[i] <= minMinMaxDist {
			pruned = append(pruned, entries[i])
		}
	}
	return pruned
}

func (tree *Rtree) nearestNeighbor(p floatgeom.Point3, n *node, d float64, nearest *Space) (*Space, float64) {
	if n.leaf {
		for _, e := range n.entries {
			dist := math.Sqrt(minDist(p, e.bb))
			if dist < d {
				d = dist
				nearest = e.obj
			}
		}
	} else {
		branches, dists := sortEntries(p, n.entries)
		branches = pruneEntries(p, branches, dists)
		for _, e := range branches {
			subNearest, dist := tree.nearestNeighbor(p, e.child, d, nearest)
			if dist < d {
				d = dist
				nearest = subNearest
			}
		}
	}

	return nearest, d
}

// NearestNeighbors returns the k nearest neighbors in the rtree to the input point.
func (tree *Rtree) NearestNeighbors(k int, p floatgeom.Point3) []*Space {
	dists := make([]float64, k)
	objs := make([]*Space, k)
	for i := 0; i < k; i++ {
		dists[i] = math.MaxFloat64
		objs[i] = nil
	}
	objs, _ = tree.nearestNeighbors(k, p, tree.root, dists, objs)
	return objs
}

// insert obj into nearest and return the first k elements in increasing order.
func insertNearest(k int, dists []float64, nearest []*Space, dist float64, obj *Space) ([]float64, []*Space) {
	i := 0
	for i < k && dist >= dists[i] {
		i++
	}
	if i >= k {
		return dists, nearest
	}

	left, right := dists[:i], dists[i:k-1]
	updatedDists := make([]float64, k)
	copy(updatedDists, left)
	updatedDists[i] = dist
	copy(updatedDists[i+1:], right)

	leftObjs, rightObjs := nearest[:i], nearest[i:k-1]
	updatedNearest := make([]*Space, k)
	copy(updatedNearest, leftObjs)
	updatedNearest[i] = obj
	copy(updatedNearest[i+1:], rightObjs)

	return updatedDists, updatedNearest
}

func (tree *Rtree) nearestNeighbors(k int, p floatgeom.Point3, n *node, dists []float64, nearest []*Space) ([]*Space, []float64) {
	if n.leaf {
		for _, e := range n.entries {
			dist := math.Sqrt(minDist(p, e.bb))
			dists, nearest = insertNearest(k, dists, nearest, dist, e.obj)
		}
	} else {
		branches, branchDists := sortEntries(p, n.entries)
		branches = pruneEntries(p, branches, branchDists)
		for _, e := range branches {
			nearest, dists = tree.nearestNeighbors(k, p, e.child, dists, nearest)
		}
	}
	return nearest, dists
}
