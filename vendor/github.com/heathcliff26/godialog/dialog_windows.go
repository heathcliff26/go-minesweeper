//go:build windows

package godialog

import (
	"fmt"
	"log/slog"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// Maximum path length for Windows
	MAX_PATH uint32 = 260
)

// GetOpenFileName and GetSaveFileName specific error codes
const (
	FNERR_BUFFERTOOSMALL  = 0x3003
	FNERR_INVALIDFILENAME = 0x3002
	FNERR_SUBCLASSFAILURE = 0x3001
)

// GetOpenFileName and GetSaveFileName flags
const (
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_CREATEPROMPT         = 0x00002000
	OFN_DONTADDTORECENT      = 0x02000000
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLEINCLUDENOTIFY  = 0x00400000
	OFN_ENABLESIZING         = 0x00800000
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_EXPLORER             = 0x00080000
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_FORCESHOWHIDDEN      = 0x10000000
	OFN_HIDEREADONLY         = 0x00000004
	OFN_LONGNAMES            = 0x00200000
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_NOLONGNAMES          = 0x00040000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NOVALIDATE           = 0x00000100
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_READONLY             = 0x00000001
	OFN_SHAREAWARE           = 0x00004000
	OFN_SHOWHELP             = 0x00000010
)

type OPENFILENAME struct {
	LStructSize       uint32
	HwndOwner         uintptr
	HInstance         uintptr
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        unsafe.Pointer
	DwReserved        uint32
	FlagsEx           uint32
}

var (
	// Library
	libcomdlg32 *windows.LazyDLL

	// Functions
	commDlgExtendedError *windows.LazyProc
	getOpenFileName      *windows.LazyProc
	getSaveFileName      *windows.LazyProc
)

func init() {
	libcomdlg32 = windows.NewLazySystemDLL("comdlg32.dll")

	commDlgExtendedError = libcomdlg32.NewProc("CommDlgExtendedError")
	getOpenFileName = libcomdlg32.NewProc("GetOpenFileNameW")
	getSaveFileName = libcomdlg32.NewProc("GetSaveFileNameW")
}

// Show a file open dialog in a new window and return path.
func (fd *fileDialog) Open(title string, cb DialogCallback) {
	err := getOpenFileName.Find()
	if err != nil {
		if fd.fallback != nil {
			slog.Info("Failed to open windows native file dialog, using fallback", "error", err)
			fd.fallback.Open(title, fd.InitialDirectory(), fd.filters, cb)
		} else {
			go cb("", fmt.Errorf("cannot open file dialog: %w", err))
		}
		return
	}
	go fd.windowsFileOpen(title, cb)
}

// Show a file save dialog in a new window and return path.
func (fd *fileDialog) Save(title string, cb DialogCallback) {
	err := getSaveFileName.Find()
	if err != nil {
		if fd.fallback != nil {
			slog.Info("Failed to open windows native file dialog, using fallback", "error", err)
			fd.fallback.Save(title, fd.InitialDirectory(), fd.filters, cb)
		} else {
			go cb("", fmt.Errorf("cannot open file dialog: %w", err))
		}
		return
	}
	go fd.windowsFileSave(title, cb)
}

func convertFilterToUTF16(filters FileFilters) (*uint16, error) {
	res := make([]uint16, 0, len(filters)*10)
	for _, filter := range filters {
		utf16Name, err := windows.UTF16FromString(filter.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to convert filter name '%s' to UTF16: %w", filter.Description, err)
		}
		res = append(res, utf16Name...)

		var extensionsStr string
		for i, ext := range filter.Extensions {
			if i > 0 {
				extensionsStr += ";"
			}
			extensionsStr += "*" + ext
		}
		utf16Ext, err := windows.UTF16FromString(extensionsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert filter extensions '%s' to UTF16: %w", extensionsStr, err)
		}
		res = append(res, utf16Ext...)
	}
	res = append(res, 0) // Null-terminate the filter string

	return &res[0], nil
}

func (fd *fileDialog) prepareOpenFileName(title string) (*OPENFILENAME, error) {
	filterPtr, err := convertFilterToUTF16(fd.filters)
	if err != nil {
		return nil, fmt.Errorf("failed to convert filter to UTF16: %w", err)
	}
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return nil, fmt.Errorf("failed to convert title to UTF16: %w", err)
	}
	startLocationPtr, err := windows.UTF16PtrFromString(fd.InitialDirectory())
	if err != nil {
		return nil, fmt.Errorf("failed to convert start location to UTF16: %w", err)
	}

	var filePath [MAX_PATH]uint16 // MAX_PATH is 260

	ofn := &OPENFILENAME{
		LStructSize:     uint32(unsafe.Sizeof(OPENFILENAME{})),
		LpstrFilter:     filterPtr,
		LpstrFile:       (*uint16)(unsafe.Pointer(&filePath[0])),
		NMaxFile:        MAX_PATH,
		LpstrTitle:      titlePtr,
		LpstrInitialDir: startLocationPtr,
	}

	return ofn, nil
}

func (fd *fileDialog) windowsFileOpen(title string, cb DialogCallback) {
	ofn, err := fd.prepareOpenFileName(title)
	if err != nil {
		cb("", err)
		return
	}

	ofn.Flags = OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_HIDEREADONLY

	ret, _, _ := getOpenFileName.Call(uintptr(unsafe.Pointer(ofn)))
	if ret == 0 {
		cb("", getLastError())
		return
	}

	result := syscall.UTF16ToString((*[MAX_PATH]uint16)(unsafe.Pointer(ofn.LpstrFile))[:])
	cb(result, nil)
}

func (fd *fileDialog) windowsFileSave(title string, cb DialogCallback) {
	ofn, err := fd.prepareOpenFileName(title)
	if err != nil {
		cb("", err)
		return
	}

	ofn.Flags = OFN_OVERWRITEPROMPT | OFN_PATHMUSTEXIST

	ret, _, _ := getSaveFileName.Call(uintptr(unsafe.Pointer(ofn)))
	if ret == 0 {
		cb("", getLastError())
		return
	}

	result := syscall.UTF16ToString((*[MAX_PATH]uint16)(unsafe.Pointer(ofn.LpstrFile))[:])
	cb(result, nil)
}

func getLastError() error {
	ret, _, _ := commDlgExtendedError.Call()
	if ret == 0 {
		return nil
	}
	switch ret {
	case FNERR_BUFFERTOOSMALL:
		return fmt.Errorf("buffer too small")
	case FNERR_INVALIDFILENAME:
		return fmt.Errorf("invalid filename")
	case FNERR_SUBCLASSFAILURE:
		return fmt.Errorf("subclass failure")
	default:
		return fmt.Errorf("unknown error: %d", ret)
	}
}
