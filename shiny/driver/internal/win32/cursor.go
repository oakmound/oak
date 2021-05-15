package win32

import "sync"

var emptyCursor HCURSOR
var emptyCursorOnce sync.Once

// Create a custom cursor at run time.
func GetEmptyCursor() HCURSOR {
	emptyCursorOnce.Do(func() {
		andMASK := []byte{
			0xFF, 0xFF, 0xFF, 0xFF,
		}

		xorMASK := []byte{
			0x00, 0x00, 0x00, 0x00,
		}
		emptyCursor = CreateCursor(hThisInstance, // app. instance
			0, // horizontal position of hot spot
			0, // vertical position of hot spot
			// 0 width/height is unsupported in testing
			1, // cursor width
			1, // cursor height
			andMASK,
			xorMASK)
	})
	return emptyCursor
}

// TODO: Add image.Image to cursor conversion and setting functionality
// this can currently be done in oak by having a image follow the cursor around,
// but that will inherently not be as smooth as setting the OS cursor. (but more portable)
