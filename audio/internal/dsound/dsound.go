//go:build windows

package dsound

import (
	"strings"
	"sync"
	"syscall"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/oakerr"
	"github.com/oov/directsound-go/dsound"
)

// A Config contains the API interfaces initalized by this package
type Config struct {
	Interface *dsound.IDirectSound
	Devices   map[*dsound.GUID]string
}

var cfg Config

var initLock sync.Mutex

// Init initializes directsound or returns an already intialized direct sound instance.
func Init() (Config, error) {
	initLock.Lock()
	defer initLock.Unlock()
	if cfg.Interface != nil {
		return cfg, nil
	}
	devices := make(map[*dsound.GUID]string)
	dsound.DirectSoundEnumerate(func(guid *dsound.GUID, description string, module string) bool {
		devices[guid] = description
		return true
	})
	if len(devices) == 0 {
		dlog.Error(dlog.NoAudioDevice)
		return cfg, oakerr.UnsupportedPlatform{
			Operation: "Init",
		}
	}
	cfg.Devices = devices
	// TODO: providing a GUID which is not nil appears to not succeed, even if those
	// GUIDs were returned by Enumerate above
	ds, err := dsound.DirectSoundCreate(nil)
	if err != nil {
		return cfg, err
	}

	user32 := syscall.NewLazyDLL("user32")
	getDesktopWindow := user32.NewProc("GetDesktopWindow")

	// Call() can return "The operation was completed successfully" as an error
	desktopWindow, _, err := getDesktopWindow.Call()
	if !strings.Contains(err.Error(), "success") {
		return cfg, err
	}
	err = ds.SetCooperativeLevel(syscall.Handle(desktopWindow), dsound.DSSCL_PRIORITY)
	if err != nil {
		return cfg, err
	}
	cfg.Interface = ds
	return cfg, nil
}
