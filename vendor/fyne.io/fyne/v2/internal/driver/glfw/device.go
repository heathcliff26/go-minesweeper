package glfw

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

type glDevice struct {
}

// Declare conformity with Device
var _ fyne.Device = (*glDevice)(nil)

func (*glDevice) Locale() fyne.Locale {
	return lang.SystemLocale()
}

func (*glDevice) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationHorizontalLeft // TODO should we consider the monitor orientation or topmost window?
}

func (*glDevice) HasKeyboard() bool {
	return true // TODO actually check - we could be in tablet mode
}

func (*glDevice) IsBrowser() bool {
	return runtime.GOARCH == "js" || runtime.GOOS == "js"
}
