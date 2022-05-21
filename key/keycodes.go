package key

import "golang.org/x/mobile/event/key"

// Code is the identity of a key relative to a notional "standard" keyboard.
// It is a straight copy of mobile package's key codes cleaned up for ease of binding in oak.
// See AllKeys for string mappers.
type Code = key.Code

const (
	Unknown Code = 0
	A       Code = 4
	B       Code = 5
	C       Code = 6
	D       Code = 7
	E       Code = 8
	F       Code = 9
	G       Code = 10
	H       Code = 11
	I       Code = 12
	J       Code = 13
	K       Code = 14
	L       Code = 15
	M       Code = 16
	N       Code = 17
	O       Code = 18
	P       Code = 19
	Q       Code = 20
	R       Code = 21
	S       Code = 22
	T       Code = 23
	U       Code = 24
	V       Code = 25
	W       Code = 26
	X       Code = 27
	Y       Code = 28
	Z       Code = 29

	Num1 Code = 30
	Num2 Code = 31
	Num3 Code = 32
	Num4 Code = 33
	Num5 Code = 34
	Num6 Code = 35
	Num7 Code = 36
	Num8 Code = 37
	Num9 Code = 38
	Num0 Code = 39

	ReturnEnter        Code = 40
	Escape             Code = 41
	DeleteBackspace    Code = 42
	Tab                Code = 43
	Spacebar           Code = 44
	HyphenMinus        Code = 45
	EqualSign          Code = 46
	LeftSquareBracket  Code = 47
	RightSquareBracket Code = 48
	Backslash          Code = 49
	Semicolon          Code = 51
	Apostrophe         Code = 52
	GraveAccent        Code = 53
	Comma              Code = 54
	FullStop           Code = 55
	Slash              Code = 56
	CapsLock           Code = 57

	F1  Code = 58
	F2  Code = 59
	F3  Code = 60
	F4  Code = 61
	F5  Code = 62
	F6  Code = 63
	F7  Code = 64
	F8  Code = 65
	F9  Code = 66
	F10 Code = 67
	F11 Code = 68
	F12 Code = 69

	Pause         Code = 72
	Insert        Code = 73
	Home          Code = 74
	PageUp        Code = 75
	DeleteForward Code = 76
	End           Code = 77
	PageDown      Code = 78

	RightArrow Code = 79
	LeftArrow  Code = 80
	DownArrow  Code = 81
	UpArrow    Code = 82

	KeypadNumLock     Code = 83
	KeypadSlash       Code = 84
	KeypadAsterisk    Code = 85
	KeypadHyphenMinus Code = 86
	KeypadPlusSign    Code = 87
	KeypadEnter       Code = 88
	Keypad1           Code = 89
	Keypad2           Code = 90
	Keypad3           Code = 91
	Keypad4           Code = 92
	Keypad5           Code = 93
	Keypad6           Code = 94
	Keypad7           Code = 95
	Keypad8           Code = 96
	Keypad9           Code = 97
	Keypad0           Code = 98
	KeypadFullStop    Code = 99
	KeypadEqualSign   Code = 103

	F13 Code = 104
	F14 Code = 105
	F15 Code = 106
	F16 Code = 107
	F17 Code = 108
	F18 Code = 109
	F19 Code = 110
	F20 Code = 111
	F21 Code = 112
	F22 Code = 113
	F23 Code = 114
	F24 Code = 115

	Help Code = 117

	Mute       Code = 127
	VolumeUp   Code = 128
	VolumeDown Code = 129

	LeftControl  Code = 224
	LeftShift    Code = 225
	LeftAlt      Code = 226
	LeftGUI      Code = 227
	RightControl Code = 228
	RightShift   Code = 229
	RightAlt     Code = 230
	RightGUI     Code = 231
	Compose      Code = 0x10000
)
