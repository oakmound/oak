package win32

import "unsafe"

// From MSDN: Windows Data Types
// http://msdn.microsoft.com/en-us/library/s3f49ktz.aspx
// http://msdn.microsoft.com/en-us/library/windows/desktop/aa383751.aspx
type (
	ATOM            uint16
	BOOL            int32
	COLORREF        uint32
	DWM_FRAME_COUNT uint64
	DWORD           uint32
	HACCEL          HANDLE
	HANDLE          uintptr
	HBITMAP         HANDLE
	HBRUSH          HANDLE
	HCURSOR         HANDLE
	HDC             HANDLE
	HDROP           HANDLE
	HDWP            HANDLE
	HENHMETAFILE    HANDLE
	HFONT           HANDLE
	HGDIOBJ         HANDLE
	HGLOBAL         HANDLE
	HGLRC           HANDLE
	HHOOK           HANDLE
	HICON           HANDLE
	HIMAGELIST      HANDLE
	HINSTANCE       HANDLE
	HKEY            HANDLE
	HKL             HANDLE
	HMENU           HANDLE
	HMODULE         HANDLE
	HMONITOR        HANDLE
	HPEN            HANDLE
	HRESULT         int32
	HRGN            HANDLE
	HRSRC           HANDLE
	HTHUMBNAIL      HANDLE
	HWND            HANDLE
	LPARAM          uintptr
	LPCVOID         unsafe.Pointer
	LRESULT         uintptr
	PVOID           unsafe.Pointer
	QPC_TIME        uint64
	ULONG_PTR       uintptr
	WPARAM          uintptr
	TRACEHANDLE     uintptr
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162897.aspx
type RECT struct {
	Left, Top, Right, Bottom int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms633577.aspx
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   *uint16
	ClassName  *uint16
	IconSm     HICON
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// Window style constants
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms632600(v=vs.85).aspx
const (
	// The window is an overlapped window.
	// An overlapped window has a title bar and a border. Same as the WS_TILED style.
	WS_OVERLAPPED = 0x00000000
	// The windows is a pop-up window. This style cannot be used with the WS_CHILD style.
	WS_POPUP = 0x80000000
	// The window is a child window. A window with this style cannot have a menu bar.
	// This style cannot be used with the WS_POPUP style.
	WS_CHILD = 0x40000000
	// The window is initially minimized. Same as the WS_ICONIC style.
	WS_MINIMIZE = 0x20000000
	// The window is initially visible.
	// This style can be turned on and off by using the ShowWindow or SetWindowPos function.
	WS_VISIBLE = 0x10000000
	// The window is initially disabled. A disabled window cannot receive input from the user.
	// To change this after a window has been created, use the EnableWindow function
	WS_DISABLED = 0x08000000
	// Clips child windows relative to each other; that is, when a particular child window receives
	// a WM_PAINT message, the WS_CLIPSIBLINGS style clips all other overlapping child windows out
	// of the region of the child window to be updated. If WS_CLIPSIBLINGS is not specified and
	// child windows overlap, it is possible, when drawing within the client area of a child window,
	// to draw within the client area of a neighboring child window.
	WS_CLIPSIBLINGS = 0x04000000
	// Excludes the area occupied by child windows when drawing occurs within the parent window.
	// This style is used when creating the parent window.
	WS_CLIPCHILDREN = 0x02000000
	// The window is initially maximized.
	WS_MAXIMIZE = 0x01000000
	// The window has a title bar (includes the WS_BORDER style).
	WS_CAPTION = 0x00C00000
	// The window has a thin-line border.
	WS_BORDER = 0x00800000
	// The window has a border of a style typically used with dialog boxes.
	// A window with this style cannot have a title bar.
	WS_DLGFRAME = 0x00400000
	// The window has a vertical scroll bar.
	WS_VSCROLL = 0x00200000
	// The window has a horizontal scroll bar.
	WS_HSCROLL = 0x00100000
	// The window has a window menu on its title bar. The WS_CAPTION style must also be specified.
	WS_SYSMENU = 0x00080000
	// The window has a sizing border. Same as the WS_SIZEBOX style.
	WS_THICKFRAME = 0x00040000
	// The window is the first control of a group of controls.
	// The group consists of this first control and all controls defined after it,
	// up to the next control with the WS_GROUP style. The first control in each group usually has
	// the WS_TABSTOP style so that the user can move from group to group. The user can subsequently
	// change the keyboard focus from one control in the group to the next control in the group by
	// using the direction keys.
	//
	// You can turn this style on and off to change dialog box navigation.
	// To change this style after a window has been created, use the SetWindowLong function.
	WS_GROUP = 0x00020000
	// The window is a control that can receive the keyboard focus when the user presses the TAB key.
	// Pressing the TAB key changes the keyboard focus to the next control with the WS_TABSTOP style.
	// You can turn this style on and off to change dialog box navigation.
	// To change this style after a window has been created, use the SetWindowLong function.
	// For user-created windows and modeless dialogs to work with tab stops, alter the message loop
	// to call the IsDialogMessage function.
	WS_TABSTOP = 0x00010000
	// The window has a minimize button. Cannot be combined with the WS_EX_CONTEXTHELP style.
	// The WS_SYSMENU style must also be specified.
	WS_MINIMIZEBOX = 0x00020000
	// The window has a maximize button. Cannot be combined with the WS_EX_CONTEXTHELP style.
	// The WS_SYSMENU style must also be specified.
	WS_MAXIMIZEBOX = 0x00010000
	// The window is an overlapped window. An overlapped window has a title bar and a border.
	// Same as the WS_OVERLAPPED style.
	WS_TILED = 0x00000000
	// The window is initially minimized. Same as the WS_MINIMIZE style.
	WS_ICONIC = 0x20000000
	// The window has a sizing border. Same as the WS_THICKFRAME style.
	WS_SIZEBOX = 0x00040000
	// The window is an overlapped window. Same as the WS_OVERLAPPEDWINDOW style.
	WS_TILEDWINDOW = WS_OVERLAPPEDWINDOW
	// The window is an overlapped window. Same as the WS_TILEDWINDOW style.
	WS_OVERLAPPEDWINDOW = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
	// The window is a pop-up window. The WS_CAPTION and WS_POPUPWINDOW styles must be
	// combined to make the window menu visible.
	WS_POPUPWINDOW = 0x80000000 | 0x00800000 | 0x00080000
	// Same as the WS_CHILD style.
	WS_CHILDWINDOW = 0x40000000
)

// Extended window style constants

const (
	// The window accepts drag-drop files.
	WS_EX_ACCEPTFILES = 0x00000010
	// Forces a top-level window onto the taskbar when the window is visible.
	WS_EX_APPWINDOW = 0x00040000
	// The window has a border with a sunken edge.
	WS_EX_CLIENTEDGE = 0x00000200
	// Paints all descendants of a window in bottom-to-top painting order using double-buffering.
	// For more information, see Remarks. This cannot be used if the window has a class style of either CS_OWNDC or CS_CLASSDC.
	// Windows 2000:  This style is not supported.
	WS_EX_COMPOSITED = 0x02000000
	// The title bar of the window includes a question mark. When the user clicks the question mark,
	// the cursor changes to a question mark with a pointer. If the user then clicks a child window,
	// the child receives a WM_HELP message. The child window should pass the message to the parent window procedure,
	// which should call the WinHelp function using the HELP_WM_HELP command. The Help application displays
	// a pop-up window that typically contains help for the child window.
	// WS_EX_CONTEXTHELP cannot be used with the WS_MAXIMIZEBOX or WS_MINIMIZEBOX styles.
	WS_EX_CONTEXTHELP = 0x00000400
	// The window itself contains child windows that should take part in dialog box navigation.
	// If this style is specified, the dialog manager recurses into children of this window when
	// performing navigation operations such as handling the TAB key, an arrow key, or a keyboard mnemonic.
	WS_EX_CONTROLPARENT = 0x00010000
	// The window has a double border; the window can, optionally, be created with a title bar
	// by specifying the WS_CAPTION style in the dwStyle parameter.
	WS_EX_DLGMODALFRAME = 0x00000001
	// The window is a layered window. This style cannot be used if the window has a
	// class style of either CS_OWNDC or CS_CLASSDC.
	// Windows 8:  The WS_EX_LAYERED style is supported for top-level windows and
	// child windows. Previous Windows versions support WS_EX_LAYERED only for top-level windows.
	WS_EX_LAYERED = 0x00080000
	// If the shell language is Hebrew, Arabic, or another language that supports reading
	// order alignment, the horizontal origin of the window is on the right edge.
	// Increasing horizontal values advance to the left.
	WS_EX_LAYOUTRTL = 0x00400000
	// The window has generic left-aligned properties. This is the default.
	WS_EX_LEFT = 0x00000000
	// If the shell language is Hebrew, Arabic, or another language that supports reading order
	// alignment, the vertical scroll bar (if present) is to the left of the client area.
	// For other languages, the style is ignored.
	WS_EX_LEFTSCROLLBAR = 0x00004000
	// The window text is displayed using left-to-right reading-order properties. This is the default.
	WS_EX_LTRREADING = 0x00000000
	// The window is a MDI child window.
	WS_EX_MDICHILD = 0x00000040
	// A top-level window created with this style does not become the foreground window when the
	// user clicks it. The system does not bring this window to the foreground when the user
	// minimizes or closes the foreground window.
	// To activate the window, use the SetActiveWindow or SetForegroundWindow function.
	// The window does not appear on the taskbar by default. To force the window to appear
	// on the taskbar, use the WS_EX_APPWINDOW style.
	WS_EX_NOACTIVATE = 0x08000000
	// The window does not pass its window layout to its child windows.
	WS_EX_NOINHERITLAYOUT = 0x00100000
	// The child window created with this style does not send the WM_PARENTNOTIFY message
	// to its parent window when it is created or destroyed.
	WS_EX_NOPARENTNOTIFY = 0x00000004
	// The window does not render to a redirection surface. This is for windows that do not
	// have visible content or that use mechanisms other than surfaces to provide their visual.
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000
	// The window is an overlapped window.
	WS_EX_OVERLAPPEDWINDOW = (WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE)
	// The window is palette window, which is a modeless dialog box that presents an array of commands.
	WS_EX_PALETTEWINDOW = (WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST)
	// The window has generic "right-aligned" properties. This depends on the window class.
	// This style has an effect only if the shell language is Hebrew, Arabic, or another language
	// that supports reading-order alignment; otherwise, the style is ignored.
	// Using the WS_EX_RIGHT style for static or edit controls has the same effect as using the SS_RIGHT
	// or ES_RIGHT style, respectively. Using this style with button controls has the same effect as
	// using BS_RIGHT and BS_RIGHTBUTTON styles.
	WS_EX_RIGHT = 0x00001000
	//The vertical scroll bar (if present) is to the right of the client area. This is the default.
	WS_EX_RIGHTSCROLLBAR = 0x00000000
	// If the shell language is Hebrew, Arabic, or another language that supports reading-order alignment,
	// the window text is displayed using right-to-left reading-order properties.
	// For other languages, the style is ignored.
	WS_EX_RTLREADING = 0x00002000
	// The window has a three-dimensional border style intended to be used for items that do not accept user input.
	WS_EX_STATICEDGE = 0x00020000
	// The window is intended to be used as a floating toolbar.
	// A tool window has a title bar that is shorter than a normal title bar, and the window title is drawn
	// using a smaller font. A tool window does not appear in the taskbar or in the dialog that
	// appears when the user presses ALT+TAB. If a tool window has a system menu, its icon is not displayed
	// on the title bar. However, you can display the system menu by right-clicking or by typing ALT+SPACE.
	WS_EX_TOOLWINDOW = 0x00000080
	// The window should be placed above all non-topmost windows and should stay above them,
	// even when the window is deactivated. To add or remove this style, use the SetWindowPos function.
	WS_EX_TOPMOST = 0x00000008
	// The window should not be painted until siblings beneath the window (that were created by the same thread)
	// have been painted. The window appears transparent because the bits of underlying sibling windows have
	// already been painted.
	// To achieve transparency without these restrictions, use the SetWindowRgn function.
	WS_EX_TRANSPARENT = 0x00000020
	// The window has a border with a raised edge.
	WS_EX_WINDOWEDGE = 0x00000100
)

// Window message constants
const (
	WM_APP                    = 32768
	WM_ACTIVATE               = 6
	WM_ACTIVATEAPP            = 28
	WM_AFXFIRST               = 864
	WM_AFXLAST                = 895
	WM_ASKCBFORMATNAME        = 780
	WM_CANCELJOURNAL          = 75
	WM_CANCELMODE             = 31
	WM_CAPTURECHANGED         = 533
	WM_CHANGECBCHAIN          = 781
	WM_CHAR                   = 258
	WM_CHARTOITEM             = 47
	WM_CHILDACTIVATE          = 34
	WM_CLEAR                  = 771
	WM_CLOSE                  = 16
	WM_COMMAND                = 273
	WM_COMMNOTIFY             = 68 /* OBSOLETE */
	WM_COMPACTING             = 65
	WM_COMPAREITEM            = 57
	WM_CONTEXTMENU            = 123
	WM_COPY                   = 769
	WM_COPYDATA               = 74
	WM_CREATE                 = 1
	WM_CTLCOLORBTN            = 309
	WM_CTLCOLORDLG            = 310
	WM_CTLCOLOREDIT           = 307
	WM_CTLCOLORLISTBOX        = 308
	WM_CTLCOLORMSGBOX         = 306
	WM_CTLCOLORSCROLLBAR      = 311
	WM_CTLCOLORSTATIC         = 312
	WM_CUT                    = 768
	WM_DEADCHAR               = 259
	WM_DELETEITEM             = 45
	WM_DESTROY                = 2
	WM_DESTROYCLIPBOARD       = 775
	WM_DEVICECHANGE           = 537
	WM_DEVMODECHANGE          = 27
	WM_DISPLAYCHANGE          = 126
	WM_DRAWCLIPBOARD          = 776
	WM_DRAWITEM               = 43
	WM_DROPFILES              = 563
	WM_ENABLE                 = 10
	WM_ENDSESSION             = 22
	WM_ENTERIDLE              = 289
	WM_ENTERMENULOOP          = 529
	WM_ENTERSIZEMOVE          = 561
	WM_ERASEBKGND             = 20
	WM_EXITMENULOOP           = 530
	WM_EXITSIZEMOVE           = 562
	WM_FONTCHANGE             = 29
	WM_GETDLGCODE             = 135
	WM_GETFONT                = 49
	WM_GETHOTKEY              = 51
	WM_GETICON                = 127
	WM_GETMINMAXINFO          = 36
	WM_GETTEXT                = 13
	WM_GETTEXTLENGTH          = 14
	WM_HANDHELDFIRST          = 856
	WM_HANDHELDLAST           = 863
	WM_HELP                   = 83
	WM_HOTKEY                 = 786
	WM_HSCROLL                = 276
	WM_HSCROLLCLIPBOARD       = 782
	WM_ICONERASEBKGND         = 39
	WM_INITDIALOG             = 272
	WM_INITMENU               = 278
	WM_INITMENUPOPUP          = 279
	WM_INPUT                  = 0x00FF
	WM_INPUTLANGCHANGE        = 81
	WM_INPUTLANGCHANGEREQUEST = 80
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_KILLFOCUS              = 8
	WM_MDIACTIVATE            = 546
	WM_MDICASCADE             = 551
	WM_MDICREATE              = 544
	WM_MDIDESTROY             = 545
	WM_MDIGETACTIVE           = 553
	WM_MDIICONARRANGE         = 552
	WM_MDIMAXIMIZE            = 549
	WM_MDINEXT                = 548
	WM_MDIREFRESHMENU         = 564
	WM_MDIRESTORE             = 547
	WM_MDISETMENU             = 560
	WM_MDITILE                = 550
	WM_MEASUREITEM            = 44
	WM_GETOBJECT              = 0x003D
	WM_CHANGEUISTATE          = 0x0127
	WM_UPDATEUISTATE          = 0x0128
	WM_QUERYUISTATE           = 0x0129
	WM_UNINITMENUPOPUP        = 0x0125
	WM_MENURBUTTONUP          = 290
	WM_MENUCOMMAND            = 0x0126
	WM_MENUGETOBJECT          = 0x0124
	WM_MENUDRAG               = 0x0123
	WM_APPCOMMAND             = 0x0319
	WM_MENUCHAR               = 288
	WM_MENUSELECT             = 287
	WM_MOVE                   = 3
	WM_MOVING                 = 534
	WM_NCACTIVATE             = 134
	WM_NCCALCSIZE             = 131
	WM_NCCREATE               = 129
	WM_NCDESTROY              = 130
	WM_NCHITTEST              = 132
	WM_NCLBUTTONDBLCLK        = 163
	WM_NCLBUTTONDOWN          = 161
	WM_NCLBUTTONUP            = 162
	WM_NCMBUTTONDBLCLK        = 169
	WM_NCMBUTTONDOWN          = 167
	WM_NCMBUTTONUP            = 168
	WM_NCXBUTTONDOWN          = 171
	WM_NCXBUTTONUP            = 172
	WM_NCXBUTTONDBLCLK        = 173
	WM_NCMOUSEHOVER           = 0x02A0
	WM_NCMOUSELEAVE           = 0x02A2
	WM_NCMOUSEMOVE            = 160
	WM_NCPAINT                = 133
	WM_NCRBUTTONDBLCLK        = 166
	WM_NCRBUTTONDOWN          = 164
	WM_NCRBUTTONUP            = 165
	WM_NEXTDLGCTL             = 40
	WM_NEXTMENU               = 531
	WM_NOTIFY                 = 78
	WM_NOTIFYFORMAT           = 85
	WM_NULL                   = 0
	WM_PAINT                  = 15
	WM_PAINTCLIPBOARD         = 777
	WM_PAINTICON              = 38
	WM_PALETTECHANGED         = 785
	WM_PALETTEISCHANGING      = 784
	WM_PARENTNOTIFY           = 528
	WM_PASTE                  = 770
	WM_PENWINFIRST            = 896
	WM_PENWINLAST             = 911
	WM_POWER                  = 72
	WM_POWERBROADCAST         = 536
	WM_PRINT                  = 791
	WM_PRINTCLIENT            = 792
	WM_QUERYDRAGICON          = 55
	WM_QUERYENDSESSION        = 17
	WM_QUERYNEWPALETTE        = 783
	WM_QUERYOPEN              = 19
	WM_QUEUESYNC              = 35
	WM_QUIT                   = 18
	WM_RENDERALLFORMATS       = 774
	WM_RENDERFORMAT           = 773
	WM_SETCURSOR              = 32
	WM_SETFOCUS               = 7
	WM_SETFONT                = 48
	WM_SETHOTKEY              = 50
	WM_SETICON                = 128
	WM_SETREDRAW              = 11
	WM_SETTEXT                = 12
	WM_SETTINGCHANGE          = 26
	WM_SHOWWINDOW             = 24
	WM_SIZE                   = 5
	WM_SIZECLIPBOARD          = 779
	WM_SIZING                 = 532
	WM_SPOOLERSTATUS          = 42
	WM_STYLECHANGED           = 125
	WM_STYLECHANGING          = 124
	WM_SYSCHAR                = 262
	WM_SYSCOLORCHANGE         = 21
	WM_SYSCOMMAND             = 274
	WM_SYSDEADCHAR            = 263
	WM_SYSKEYDOWN             = 260
	WM_SYSKEYUP               = 261
	WM_TCARD                  = 82
	WM_THEMECHANGED           = 794
	WM_TIMECHANGE             = 30
	WM_TIMER                  = 275
	WM_UNDO                   = 772
	WM_USER                   = 1024
	WM_USERCHANGED            = 84
	WM_VKEYTOITEM             = 46
	WM_VSCROLL                = 277
	WM_VSCROLLCLIPBOARD       = 778
	WM_WINDOWPOSCHANGED       = 71
	WM_WINDOWPOSCHANGING      = 70
	WM_WININICHANGE           = 26
	WM_KEYFIRST               = 256
	WM_KEYLAST                = 264
	WM_SYNCPAINT              = 136
	WM_MOUSEACTIVATE          = 33
	WM_MOUSEMOVE              = 512
	WM_LBUTTONDOWN            = 513
	WM_LBUTTONUP              = 514
	WM_LBUTTONDBLCLK          = 515
	WM_RBUTTONDOWN            = 516
	WM_RBUTTONUP              = 517
	WM_RBUTTONDBLCLK          = 518
	WM_MBUTTONDOWN            = 519
	WM_MBUTTONUP              = 520
	WM_MBUTTONDBLCLK          = 521
	WM_MOUSEWHEEL             = 522
	WM_XBUTTONDOWN            = 523
	WM_XBUTTONUP              = 524
	WM_XBUTTONDBLCLK          = 525
	WM_MOUSEHWHEEL            = 526
	WM_MOUSEFIRST             = 512
	WM_MOUSELAST              = 526
	WM_MOUSEHOVER             = 0x2A1
	WM_MOUSELEAVE             = 0x2A3
	WM_CLIPBOARDUPDATE        = 0x031D
)

// Predefined brushes constants
const (
	COLOR_3DDKSHADOW              = 21
	COLOR_3DFACE                  = 15
	COLOR_3DHILIGHT               = 20
	COLOR_3DHIGHLIGHT             = 20
	COLOR_3DLIGHT                 = 22
	COLOR_BTNHILIGHT              = 20
	COLOR_3DSHADOW                = 16
	COLOR_ACTIVEBORDER            = 10
	COLOR_ACTIVECAPTION           = 2
	COLOR_APPWORKSPACE            = 12
	COLOR_BACKGROUND              = 1
	COLOR_DESKTOP                 = 1
	COLOR_BTNFACE                 = 15
	COLOR_BTNHIGHLIGHT            = 20
	COLOR_BTNSHADOW               = 16
	COLOR_BTNTEXT                 = 18
	COLOR_CAPTIONTEXT             = 9
	COLOR_GRAYTEXT                = 17
	COLOR_HIGHLIGHT               = 13
	COLOR_HIGHLIGHTTEXT           = 14
	COLOR_INACTIVEBORDER          = 11
	COLOR_INACTIVECAPTION         = 3
	COLOR_INACTIVECAPTIONTEXT     = 19
	COLOR_INFOBK                  = 24
	COLOR_INFOTEXT                = 23
	COLOR_MENU                    = 4
	COLOR_MENUTEXT                = 7
	COLOR_SCROLLBAR               = 0
	COLOR_WINDOW                  = 5
	COLOR_WINDOWFRAME             = 6
	COLOR_WINDOWTEXT              = 8
	COLOR_HOTLIGHT                = 26
	COLOR_GRADIENTACTIVECAPTION   = 27
	COLOR_GRADIENTINACTIVECAPTION = 28
)

// Predefined window handles
const (
	HWND_BROADCAST = HWND(0xFFFF)
	HWND_BOTTOM    = HWND(1)
	HWND_NOTOPMOST = ^HWND(1) // -2
	HWND_TOP       = HWND(0)
	HWND_TOPMOST   = ^HWND(0) // -1
	HWND_DESKTOP   = HWND(0)
	HWND_MESSAGE   = ^HWND(2) // -3
)

const (
	IMAGE_BITMAP = 0
	IMAGE_ICON   = 1
	IMAGE_CURSOR = 2

	LR_DEFAULTCOLOR     = 0x00000000
	LR_MONOCHROME       = 0x00000001
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADTRANSPARENT  = 0x00000020
	LR_DEFAULTSIZE      = 0x00000040
	LR_VGACOLOR         = 0x00000080
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_CREATEDIBSECTION = 0x00002000
	LR_SHARED           = 0x00008000
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa373931.aspx
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type Shell_NotifyAction uint32

const (
	NIM_ADD        Shell_NotifyAction = 0x0
	NIM_MODIFY     Shell_NotifyAction = 0x1
	NIM_DELETE     Shell_NotifyAction = 0x2
	NIM_SETFOCUS   Shell_NotifyAction = 0x3
	NIM_SETVERSION Shell_NotifyAction = 0x4
)

const (
	NIF_MESSAGE  = 0x00000001
	NIF_ICON     = 0x00000002
	NIF_TIP      = 0x00000004
	NIF_STATE    = 0x00000008
	NIF_INFO     = 0x00000010
	NIF_GUID     = 0x00000020
	NIF_REALTIME = 0x00000040
	NIF_SHOWTIP  = 0x00000080
)

const (
	NIIF_NONE               = 0x00000000
	NIIF_INFO               = 0x00000001
	NIIF_WARNING            = 0x00000002
	NIIF_ERROR              = 0x00000003
	NIIF_USER               = 0x00000004
	NIIF_NOSOUND            = 0x00000010
	NIIF_LARGE_ICON         = 0x00000020
	NIIF_RESPECT_QUIET_TIME = 0x00000080
	NIIF_ICON_MASK          = 0x0000000F
)

type NOTIFYICONDATA struct {
	CbSize           uint32
	HWnd             HWND
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GUIDItem         GUID
	HBalloonIcon     uintptr
}

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_DRAWFRAME      = SWP_FRAMECHANGED
	SWP_NOREPOSITION   = SWP_NOOWNERZORDER
	SWP_DEFERERASE     = 0x2000
	SWP_ASYNCWINDOWPOS = 0x4000
)

// WM_SYSCOMMAND wParams
const (
	SC_MAXIMIZE = 0xF030
	SC_RESTORE  = 0xF120
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145065.aspx
type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

const (
	MONITOR_DEFAULTTONULL    = 0x00000000
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002

	MONITORINFOF_PRIMARY = 0x00000001
)
