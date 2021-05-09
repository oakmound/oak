package win32

import (
	"crypto/rand"
	"unsafe"
)

func boolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}

// MakeGUID allocates a GUID from a [16]byte array. If the array
// is uninitialized a random byte array will be used.
func MakeGUID(guid [16]byte) GUID {
	if guid == ([16]byte{}) {
		rand.Read(guid[:])
	}
	// Implementation from hallazzang/go-windows-programming. Would have
	// reimplemented but it's so simple that its hard to do so.
	return *(*GUID)(unsafe.Pointer(&guid[0]))
}
