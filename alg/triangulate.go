package alg

// TriangulateConvex takes a face, in the form of a slice
// of indices, and outputs those indicies split into triangles
// based on the assumption that the face is convex. This involves
// forming pairs of indices and drawing an edge back to the first
// index repeatedly.
//
// If given less than 3 indices, returns an empty slice.
//
// Example input: [0,1,2,3,4]
// Example output: [[0,1,2][0,2,3][0,3,4]]
//
// Visual Example:
//      ____0____
//     4         1
//      \       /
//       3-----2
//        - - -
//      ____0____
//     4   / \   1
//      \ /   \ /
//       3-----2
//
// This makes additional assumptions that the points represented
// by the indices are coplanar, and that there are no holes present
// in the face.
func TriangulateConvex(face []int) [][3]int {
	if len(face) < 3 {
		return [][3]int{}
	}
	tris := make([][3]int, len(face)-2)
	for i := 0; i < len(tris); i++ {
		tris[i][0] = face[0]
		tris[i][1] = face[i+1]
		tris[i][2] = face[i+2]
	}
	return tris
}
