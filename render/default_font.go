package render

import (
	_ "embed"
)

// Oak ships with a free font to enable text display without needing to set up
// a font on the user's machine to import. This is that font. It is embedded into
// the Go code to ensure it is not stripped from the code by vendoring, for example.
// The file is called luxisr.ttf.
// Todo: consider shipping with a smaller-sized font (would be a breaking change)

//go:embed luxisr.ttf
var luxisrTTF []byte
