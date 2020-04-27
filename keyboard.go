package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var Vkeys = map[string]uint16{
	"0":          0x30,
	"1":          0x31,
	"2":          0x32,
	"3":          0x33,
	"4":          0x34,
	"5":          0x35,
	"6":          0x36,
	"7":          0x37,
	"8":          0x38,
	"9":          0x39,
	"A":          0x41,
	"B":          0x42,
	"C":          0x43,
	"D":          0x44,
	"E":          0x45,
	"F":          0x46,
	"G":          0x47,
	"H":          0x48,
	"I":          0x49,
	"J":          0x4A,
	"K":          0x4B,
	"L":          0x4C,
	"M":          0x4D,
	"N":          0x4E,
	"O":          0x4F,
	"P":          0x50,
	"Q":          0x51,
	"R":          0x52,
	"S":          0x53,
	"T":          0x54,
	"U":          0x55,
	"V":          0x56,
	"W":          0x57,
	"X":          0x58,
	"Y":          0x59,
	"Z":          0x5A,
	"BACK":       0x08,
	"TAB":        0x09,
	"CLEAR":      0x0C,
	"RETURN":     0x0D,
	"SHIFT":      0x10,
	"CONTROL":    0x11,
	"MENU":       0x12,
	"PAUSE":      0x13,
	"CAPITAL":    0x14,
	"KANA":       0x15,
	"HANGEUL":    0x15, /* old name - should be here for compatibility */
	"HANGUL":     0x15,
	"JUNJA":      0x17,
	"FINAL":      0x18,
	"HANJA":      0x19,
	"KANJI":      0x19,
	"ESCAPE":     0x1B,
	"CONVERT":    0x1C,
	"NONCONVERT": 0x1D,
	"ACCEPT":     0x1E,
	"MODECHANGE": 0x1F,
	"SPACE":      0x20,
	"PRIOR":      0x21,
	"NEXT":       0x22,
	"END":        0x23,
	"HOME":       0x24,
	"LEFT":       0x25,
	"UP":         0x26,
	"RIGHT":      0x27,
	"DOWN":       0x28,
	"SELECT":     0x29,
	"PRINT":      0x2A,
	"EXECUTE":    0x2B,
	"SNAPSHOT":   0x2C,
	"INSERT":     0x2D,
	"DELETE":     0x2E,
	"HELP":       0x2F,
	"LWIN":       0x5B,
	"RWIN":       0x5C,
	"APPS":       0x5D,
	"SLEEP":      0x5F,
	"NUMPAD0":    0x60,
	"NUMPAD1":    0x61,
	"NUMPAD2":    0x62,
	"NUMPAD3":    0x63,
	"NUMPAD4":    0x64,
	"NUMPAD5":    0x65,
	"NUMPAD6":    0x66,
	"NUMPAD7":    0x67,
	"NUMPAD8":    0x68,
	"NUMPAD9":    0x69,
	"MULTIPLY":   0x6A,
	"ADD":        0x6B,
	"SEPARATOR":  0x6C,
	"SUBTRACT":   0x6D,
	"DECIMAL":    0x6E,
	"DIVIDE":     0x6F,
	"F1":         0x70,
	"F2":         0x71,
	"F3":         0x72,
	"F4":         0x73,
	"F5":         0x74,
	"F6":         0x75,
	"F7":         0x76,
	"F8":         0x77,
	"F9":         0x78,
	"F10":        0x79,
	"F11":        0x7A,
	"F12":        0x7B,
	"F13":        0x7C,
	"F14":        0x7D,
	"F15":        0x7E,
	"F16":        0x7F,
	"F17":        0x80,
	"F18":        0x81,
	"F19":        0x82,
	"F20":        0x83,
	"F21":        0x84,
	"F22":        0x85,
	"F23":        0x86,
	"F24":        0x87,
}
var ScanKeys = map[string]uint16{
	"ESC":          0x0401, //
	"1":            0x0402, //
	"2":            0x0403, //
	"3":            0x0404, //
	"4":            0x0405, //
	"5":            0x0406, //
	"6":            0x0407, //
	"7":            0x0408, //
	"8":            0x0409, //
	"9":            0x040A, //
	"0":            0x040B, //
	"MINUS":        0x040C, // (* - on main keyboard *)
	"EQUALS":       0x040D, //
	"BACK":         0x040E, // (* backspace *)
	"TAB":          0x040F, //
	"Q":            0x0410, //
	"W":            0x0411, //
	"E":            0x0412, //
	"R":            0x0413, //
	"T":            0x0414, //
	"Y":            0x0415, //
	"U":            0x0416, //
	"I":            0x0417, //
	"O":            0x0418, //
	"P":            0x0419, //
	"LBRACKET":     0x041A, //
	"RBRACKET":     0x041B, //
	"RETURN":       0x041C, // (* Enter on main keyboard *)
	"LCONTROL":     0x041D, //
	"A":            0x041E, //
	"S":            0x041F, //
	"D":            0x0420, //
	"F":            0x0421, //
	"G":            0x0422, //
	"H":            0x0423, //
	"J":            0x0424, //
	"K":            0x0425, //
	"L":            0x0426, //
	"SEMICOLON":    0x0427, //
	"APOSTROPHE":   0x0428, //
	"GRAVE":        0x0429, // (* accent grave *)
	"LSHIFT":       0x042A, //
	"BACKSLASH":    0x042B, //
	"Z":            0x042C, //
	"X":            0x042D, //
	"C":            0x042E, //
	"V":            0x042F, //
	"B":            0x0430, //
	"N":            0x0431, //
	"M":            0x0432, //
	"COMMA":        0x0433, //
	"PERIOD":       0x0434, // (* . on main keyboard *)
	"SLASH":        0x0435, // (* / on main keyboard *)
	"RSHIFT":       0x0436, //
	"MULTIPLY":     0x0437, // (* * on numeric keypad *)
	"LMENU":        0x0438, // (* left Alt *)
	"SPACE":        0x0439, //
	"CAPITAL":      0x043A, //
	"F1":           0x043B, //
	"F2":           0x043C, //
	"F3":           0x043D, //
	"F4":           0x043E, //
	"F5":           0x043F, //
	"F6":           0x0440, //
	"F7":           0x0441, //
	"F8":           0x0442, //
	"F9":           0x0443, //
	"F10":          0x0444, //
	"NUMLOCK":      0x0445, //
	"SCROLL":       0x0446, // (* Scroll Lock *)
	"NUMPAD7":      0x0447, //
	"NUMPAD8":      0x0448, //
	"NUMPAD9":      0x0449, //
	"SUBTRACT":     0x044A, // (* - on numeric keypad *)
	"NUMPAD4":      0x044B, //
	"NUMPAD5":      0x044C, //
	"NUMPAD6":      0x044D, //
	"ADD":          0x044E, // (* + on numeric keypad *)
	"NUMPAD1":      0x044F, //
	"NUMPAD2":      0x0450, //
	"NUMPAD3":      0x0451, //
	"NUMPAD0":      0x0452, //
	"DECIMAL":      0x0453, // (* . on numeric keypad *)
	"OEM_102":      0x0456, // (* < > | on UK/Germany keyboards *)
	"F11":          0x0457, //
	"F12":          0x0458, //
	"F13":          0x0464, // (* (NEC PC98) *)
	"F14":          0x0465, // (* (NEC PC98) *)
	"F15":          0x0466, // (* (NEC PC98) *)
	"KANA":         0x0470, // (* (Japanese keyboard) *)
	"ABNT_C1":      0x0473, // (* / ? on Portugese (Brazilian) keyboards *)
	"CONVERT":      0x0479, // (* (Japanese keyboard) *)
	"NOCONVERT":    0x047B, // (* (Japanese keyboard) *)
	"YEN":          0x047D, // (* (Japanese keyboard) *)
	"ABNT_C2":      0x047E, // (* Numpad . on Portugese (Brazilian) keyboards *)
	"NUMPADEQUALS": 0x048D, // (* = on numeric keypad (NEC PC98) *)
	"PREVTRACK":    0x0490, // (* Previous Track (DIK_CIRCUMFLEX on Japanese keyboard) *)
	"AT":           0x0491, // (* (NEC PC98) *)
	"COLON":        0x0492, // (* (NEC PC98) *)
	"UNDERLINE":    0x0493, // (* (NEC PC98) *)
	"KANJI":        0x0494, // (* (Japanese keyboard) *)
	"STOP":         0x0495, // (* (NEC PC98) *)
	"AX":           0x0496, // (* (Japan AX) *)
	"UNLABELED":    0x0497, // (* (J3100) *)
	"NEXTTRACK":    0x0499, // (* Next Track *)
	"NUMPADENTER":  0x049C, // (* Enter on numeric keypad *)
	"RCONTROL":     0x049D, //
	"MUTE":         0x04A0, // (* Mute *)
	"CALCULATOR":   0x04A1, // (* Calculator *)
	"PLAYPAUSE":    0x04A2, // (* Play / Pause *)
	"MEDIASTOP":    0x04A4, // (* Media Stop *)
	"VOLUMEDOWN":   0x04AE, // (* Volume - *)
	"VOLUMEUP":     0x04B0, // (* Volume + *)
	"WEBHOME":      0x04B2, // (* Web home *)
	"NUMPADCOMMA":  0x04B3, // (* , on numeric keypad (NEC PC98) *)
	"DIVIDE":       0x04B5, // (* / on numeric keypad *)
	"SYSRQ":        0x04B7, //
	"RMENU":        0x04B8, // (* right Alt *)
	"PAUSE":        0x04C5, // (* Pause *)
	"HOME":         0x04C7, // (* Home on arrow keypad *)
	"UP":           0x04C8, // (* UpArrow on arrow keypad *)
	"PRIOR":        0x04C9, // (* PgUp on arrow keypad *)
	"LEFT":         0x04CB, // (* LeftArrow on arrow keypad *)
	"RIGHT":        0x04CD, // (* RightArrow on arrow keypad *)
	"END":          0x04CF, // (* End on arrow keypad *)
	"DOWN":         0x04D0, // (* DownArrow on arrow keypad *)
	"NEXT":         0x04D1, // (* PgDn on arrow keypad *)
	"INSERT":       0x04D2, // (* Insert on arrow keypad *)
	"DELETE":       0x04D3, // (* Delete on arrow keypad *)
	"LWIN":         0x04DB, // (* Left Windows key *)
	"RWIN":         0x04DC, // (* Right Windows key *)
	"APPS":         0x04DD, // (* AppMenu key *)
	"POWER":        0x04DE, // (* System Power *)
	"SLEEP":        0x04DF, // (* System Sleep *)
	"WAKE":         0x04E3, // (* System Wake *)
	"WEBSEARCH":    0x04E5, // (* Web Search *)
	"WEBFAVORITES": 0x04E6, // (* Web Favorites *)
	"WEBREFRESH":   0x04E7, // (* Web Refresh *)
	"WEBSTOP":      0x04E8, // (* Web Stop *)
	"WEBFORWARD":   0x04E9, // (* Web Forward *)
	"WEBBACK":      0x04EA, // (* Web Back *)
	"MYCOMPUTER":   0x04EB, // (* My Computer *)
	"MAIL":         0x04EC, // (* Mail *)
	"MEDIASELECT":  0x04ED, // (* Media Select *)
}

// static void dummy(void) { }
type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uint64
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   uint64
}

var dll = syscall.NewLazyDLL("user32.dll")
var sendInputProc = dll.NewProc("SendInput")
var mapVirtualKeyProc = dll.NewProc("MapVirtualKeyA")

// SendKeyPress send key press down, wait and send key release
func SendKeyPress(Key string) (ok bool, err error) {

	if ok, err := SendInput(false, Key); !ok {
		return false, fmt.Errorf("%s, Sendkeypress: Error sending key press", err)
	}

	time.Sleep(10 * time.Millisecond)

	if ok, err := SendInput(true, Key); !ok {
		return false, fmt.Errorf("%s, Sendkeypress: Error sending key release", err)
	}

	return true, nil

}

// SendInput send the keypress of the 'key' or release
func SendInput(release bool, key string) (ok bool, err error) {
	var i input
	i.inputType = 1 //INPUT_KEYBOARD
	/*wScan, ok := ScanKeys[key] //0x041E // virtual key code for a
	if !ok {
		return false, fmt.Errorf("Sendinput: Wrong ScanKeys specified: %s", key)
	}
	i.ki.wScan = wScan*/

	vVk, ok := Vkeys[key]

	i.ki.wVk = vVk
	if !ok {
		return false, fmt.Errorf("Sendinput: Wrong vKey specified: %s", key)
	}

	wScan, _, _ := mapVirtualKeyProc.Call(uintptr(vVk), uintptr(0))

	i.ki.wScan = uint16(wScan)

	if release == true {
		i.ki.dwFlags = 0x0002
	} else {
		i.ki.dwFlags = 0
	}
	ret, _, err := sendInputProc.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
	log.Printf("ret: %v error: %v", ret, err)
	return true, nil

}
