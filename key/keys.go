package key

// This lists the keys sent through oak's input events.
// This list is not used internally by oak, but was generated from
// the expected output from x/mobile/key.
// todo: write a go generate script to perform said generation
//
// These strings are sent as payloads to Key.Down and Key.Up events,
// and through "KeyDown"+$a, "KeyUp"+$a for any $a in the const.
const (
	Unknown = "Unknown"

	A = "A"
	B = "B"
	C = "C"
	D = "D"
	E = "E"
	F = "F"
	G = "G"
	H = "H"
	I = "I"
	J = "J"
	K = "K"
	L = "L"
	M = "M"
	N = "N"
	O = "O"
	P = "P"
	Q = "Q"
	R = "R"
	S = "S"
	T = "T"
	U = "U"
	V = "V"
	W = "W"
	X = "X"
	Y = "Y"
	Z = "Z"

	One   = "1"
	Two   = "2"
	Three = "3"
	Four  = "4"
	Five  = "5"
	Six   = "6"
	Seven = "7"
	Eight = "8"
	Nine  = "9"
	Zero  = "0"

	ReturnEnter        = "ReturnEnter"
	Escape             = "Escape"
	DeleteBackspace    = "DeleteBackspace"
	Tab                = "Tab"
	Spacebar           = "Spacebar"
	HyphenMinus        = "HyphenMinus"        //-
	EqualSign          = "EqualSign"          //=
	LeftSquareBracket  = "LeftSquareBracket"  //[
	RightSquareBracket = "RightSquareBracket" //]
	Backslash          = "Backslash"          //\
	Semicolon          = "Semicolon"          //;
	Apostrophe         = "Apostrophe"         //'
	GraveAccent        = "GraveAccent"        //`
	Comma              = "Comma"              //,
	FullStop           = "FullStop"           //.
	Slash              = "Slash"              ///
	CapsLock           = "CapsLock"

	F1  = "F1"
	F2  = "F2"
	F3  = "F3"
	F4  = "F4"
	F5  = "F5"
	F6  = "F6"
	F7  = "F7"
	F8  = "F8"
	F9  = "F9"
	F10 = "F10"
	F11 = "F11"
	F12 = "F12"

	Pause         = "Pause"
	Insert        = "Insert"
	Home          = "Home"
	PageUp        = "PageUp"
	DeleteForward = "DeleteForward"
	End           = "End"
	PageDown      = "PageDown"

	RightArrow = "RightArrow"
	LeftArrow  = "LeftArrow"
	DownArrow  = "DownArrow"
	UpArrow    = "UpArrow"

	KeypadNumLock     = "KeypadNumLock"
	KeypadSlash       = "KeypadSlash"       ///
	KeypadAsterisk    = "KeypadAsterisk"    //*
	KeypadHyphenMinus = "KeypadHyphenMinus" //-
	KeypadPlusSign    = "KeypadPlusSign"    //+
	KeypadEnter       = "KeypadEnter"
	Keypad1           = "Keypad1"
	Keypad2           = "Keypad2"
	Keypad3           = "Keypad3"
	Keypad4           = "Keypad4"
	Keypad5           = "Keypad5"
	Keypad6           = "Keypad6"
	Keypad7           = "Keypad7"
	Keypad8           = "Keypad8"
	Keypad9           = "Keypad9"
	Keypad0           = "Keypad0"
	KeypadFullStop    = "KeypadFullStop"  //.
	KeypadEqualSign   = "KeypadEqualSign" //=

	F13 = "F13"
	F14 = "F14"
	F15 = "F15"
	F16 = "F16"
	F17 = "F17"
	F18 = "F18"
	F19 = "F19"
	F20 = "F20"
	F21 = "F21"
	F22 = "F22"
	F23 = "F23"
	F24 = "F24"

	Help = "Help"

	Mute       = "Mute"
	VolumeUp   = "VolumeUp"
	VolumeDown = "VolumeDown"

	LeftControl  = "LeftControl"
	LeftShift    = "LeftShift"
	LeftAlt      = "LeftAlt"
	LeftGUI      = "LeftGUI"
	RightControl = "RightControl"
	RightShift   = "RightShift"
	RightAlt     = "RightAlt"
	RightGUI     = "RightGUI"
)
