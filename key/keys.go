package key

// This lists the keys sent through oak's input events.
// This list is not used internally by oak, but was generated from
// the expected output from x/mobile/key.
//
// These strings are sent as payloads to Key.Down and Key.Up events,
// and through "KeyDown"+$a, "KeyUp"+$a for any $a in the const.
const (
	UnknownStr = "Unknown"

	AStr = "A"
	BStr = "B"
	CStr = "C"
	DStr = "D"
	EStr = "E"
	FStr = "F"
	GStr = "G"
	HStr = "H"
	IStr = "I"
	JStr = "J"
	KStr = "K"
	LStr = "L"
	MStr = "M"
	NStr = "N"
	OStr = "O"
	PStr = "P"
	QStr = "Q"
	RStr = "R"
	SStr = "S"
	TStr = "T"
	UStr = "U"
	VStr = "V"
	WStr = "W"
	XStr = "X"
	YStr = "Y"
	ZStr = "Z"

	OneStr   = "1"
	TwoStr   = "2"
	ThreeStr = "3"
	FourStr  = "4"
	FiveStr  = "5"
	SixStr   = "6"
	SevenStr = "7"
	EightStr = "8"
	NineStr  = "9"
	ZeroStr  = "0"

	ReturnEnterStr        = "ReturnEnter"
	Enter                 = ReturnEnter
	EscapeStr             = "Escape"
	DeleteBackspaceStr    = "DeleteBackspace"
	TabStr                = "Tab"
	SpacebarStr           = "Spacebar"
	HyphenMinusStr        = "HyphenMinus"        //-
	EqualSignStr          = "EqualSign"          //=
	LeftSquareBracketStr  = "LeftSquareBracket"  //[
	RightSquareBracketStr = "RightSquareBracket" //]
	BackslashStr          = "Backslash"          //\
	SemicolonStr          = "Semicolon"          //;
	ApostropheStr         = "Apostrophe"         //'
	GraveAccentStr        = "GraveAccent"        //`
	CommaStr              = "Comma"              //,
	FullStopStr           = "FullStop"           //.
	Period                = "FullStop"
	SlashStr              = "Slash" ///
	CapsLockStr           = "CapsLock"

	F1Str  = "F1"
	F2Str  = "F2"
	F3Str  = "F3"
	F4Str  = "F4"
	F5Str  = "F5"
	F6Str  = "F6"
	F7Str  = "F7"
	F8Str  = "F8"
	F9Str  = "F9"
	F10Str = "F10"
	F11Str = "F11"
	F12Str = "F12"

	PauseStr         = "Pause"
	InsertStr        = "Insert"
	HomeStr          = "Home"
	PageUpStr        = "PageUp"
	DeleteForwardStr = "DeleteForward"
	EndStr           = "End"
	PageDownStr      = "PageDown"

	RightArrowStr = "RightArrow"
	LeftArrowStr  = "LeftArrow"
	DownArrowStr  = "DownArrow"
	UpArrowStr    = "UpArrow"

	KeypadNumLockStr     = "KeypadNumLock"
	KeypadSlashStr       = "KeypadSlash"       ///
	KeypadAsteriskStr    = "KeypadAsterisk"    //*
	KeypadHyphenMinusStr = "KeypadHyphenMinus" //-
	KeypadPlusSignStr    = "KeypadPlusSign"    //+
	KeypadEnterStr       = "KeypadEnter"
	Keypad1Str           = "Keypad1"
	Keypad2Str           = "Keypad2"
	Keypad3Str           = "Keypad3"
	Keypad4Str           = "Keypad4"
	Keypad5Str           = "Keypad5"
	Keypad6Str           = "Keypad6"
	Keypad7Str           = "Keypad7"
	Keypad8Str           = "Keypad8"
	Keypad9Str           = "Keypad9"
	Keypad0Str           = "Keypad0"
	KeypadFullStopStr    = "KeypadFullStop" //.
	KeypadPeriod         = "KeypadFullStop"
	KeypadEqualSignStr   = "KeypadEqualSign" //=

	F13Str = "F13"
	F14Str = "F14"
	F15Str = "F15"
	F16Str = "F16"
	F17Str = "F17"
	F18Str = "F18"
	F19Str = "F19"
	F20Str = "F20"
	F21Str = "F21"
	F22Str = "F22"
	F23Str = "F23"
	F24Str = "F24"

	HelpStr = "Help"

	MuteStr       = "Mute"
	VolumeUpStr   = "VolumeUp"
	VolumeDownStr = "VolumeDown"

	LeftControlStr  = "LeftControl"
	LeftShiftStr    = "LeftShift"
	LeftAltStr      = "LeftAlt"
	LeftGUIStr      = "LeftGUI"
	RightControlStr = "RightControl"
	RightShiftStr   = "RightShift"
	RightAltStr     = "RightAlt"
	RightGUIStr     = "RightGUI"
)

// AllKeys is the set of all defined key codes to their Codes
var AllKeys = map[string]Code{
	UnknownStr: Unknown,

	AStr: A,
	BStr: B,
	CStr: C,
	DStr: D,
	EStr: E,
	FStr: F,
	GStr: G,
	HStr: H,
	IStr: I,
	JStr: J,
	KStr: K,
	LStr: L,
	MStr: M,
	NStr: N,
	OStr: O,
	PStr: P,
	QStr: Q,
	RStr: R,
	SStr: S,
	TStr: T,
	UStr: U,
	VStr: V,
	WStr: W,
	XStr: X,
	YStr: Y,
	ZStr: Z,

	OneStr:   Num1,
	TwoStr:   Num2,
	ThreeStr: Num3,
	FourStr:  Num4,
	FiveStr:  Num5,
	SixStr:   Num6,
	SevenStr: Num7,
	EightStr: Num8,
	NineStr:  Num9,
	ZeroStr:  Num0,

	ReturnEnterStr:        ReturnEnter,
	EscapeStr:             Escape,
	DeleteBackspaceStr:    DeleteBackspace,
	TabStr:                Tab,
	SpacebarStr:           Spacebar,
	HyphenMinusStr:        HyphenMinus,
	EqualSignStr:          EqualSign,
	LeftSquareBracketStr:  LeftSquareBracket,
	RightSquareBracketStr: RightSquareBracket,
	BackslashStr:          Backslash,
	SemicolonStr:          Semicolon,
	ApostropheStr:         Apostrophe,
	GraveAccentStr:        GraveAccent,
	CommaStr:              Comma,
	FullStopStr:           FullStop,
	SlashStr:              Slash,
	CapsLockStr:           CapsLock,

	F1Str:  F1,
	F2Str:  F2,
	F3Str:  F3,
	F4Str:  F4,
	F5Str:  F5,
	F6Str:  F6,
	F7Str:  F7,
	F8Str:  F8,
	F9Str:  F9,
	F10Str: F10,
	F11Str: F11,
	F12Str: F12,

	PauseStr:         Pause,
	InsertStr:        Insert,
	HomeStr:          Home,
	PageUpStr:        PageUp,
	DeleteForwardStr: DeleteForward,
	EndStr:           End,
	PageDownStr:      PageDown,

	RightArrowStr: RightArrow,
	LeftArrowStr:  LeftArrow,
	DownArrowStr:  DownArrow,
	UpArrowStr:    UpArrow,

	KeypadNumLockStr:     KeypadNumLock,
	KeypadSlashStr:       KeypadSlash,
	KeypadAsteriskStr:    KeypadAsterisk,
	KeypadHyphenMinusStr: KeypadHyphenMinus,
	KeypadPlusSignStr:    KeypadPlusSign,
	KeypadEnterStr:       KeypadEnter,
	Keypad1Str:           Keypad1,
	Keypad2Str:           Keypad2,
	Keypad3Str:           Keypad3,
	Keypad4Str:           Keypad4,
	Keypad5Str:           Keypad5,
	Keypad6Str:           Keypad6,
	Keypad7Str:           Keypad7,
	Keypad8Str:           Keypad8,
	Keypad9Str:           Keypad9,
	Keypad0Str:           Keypad0,
	KeypadFullStopStr:    KeypadFullStop,
	KeypadEqualSignStr:   KeypadEqualSign,

	F13Str: F13,
	F14Str: F14,
	F15Str: F15,
	F16Str: F16,
	F17Str: F17,
	F18Str: F18,
	F19Str: F19,
	F20Str: F20,
	F21Str: F21,
	F22Str: F22,
	F23Str: F23,
	F24Str: F24,

	HelpStr: Help,

	MuteStr:       Mute,
	VolumeUpStr:   VolumeUp,
	VolumeDownStr: VolumeDown,

	LeftControlStr:  LeftControl,
	LeftShiftStr:    LeftShift,
	LeftAltStr:      LeftAlt,
	LeftGUIStr:      LeftGUI,
	RightControlStr: RightControl,
	RightShiftStr:   RightShift,
	RightAltStr:     RightAlt,
	RightGUIStr:     RightGUI,
}
