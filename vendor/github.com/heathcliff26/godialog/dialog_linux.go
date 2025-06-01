//go:build linux

package godialog

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	DBusObjectName = "org.freedesktop.portal.Desktop"
	DBusObjectPath = "/org/freedesktop/portal/desktop"

	DBusFileChooserBase     = "org.freedesktop.portal.FileChooser"
	DBusFileChooserOpenFile = ".OpenFile"
	DBusFileChooserSaveFile = ".SaveFile"
)

type freedesktopFilterRule struct {
	Type    uint32
	Pattern string
}

// Filter specifies a filter containing various rules for allowed files.
type freedesktopFilter struct {
	Name  string
	Rules []freedesktopFilterRule
}

// Show a file open dialog in a new window and return path.
func (fd *fileDialog) Open(title string, cb DialogCallback) {
	go fd.open(title, cb)
}

// Show a file save dialog in a new window and return path.
func (fd *fileDialog) Save(title string, cb DialogCallback) {
	go fd.save(title, cb)
}

// The actual implementation for Open. Should be run in a goroutine.
func (fd *fileDialog) open(title string, cb DialogCallback) {
	err := fd.dbusFileChooser(DBusFileChooserOpenFile, title)
	if err != nil {
		if fd.fallback != nil {
			slog.Info("Failed to open linux native file dialog, using fallback", "error", err)
			fd.fallback.Open(title, fd.InitialDirectory(), fd.filters, cb)
		} else {
			cb("", fmt.Errorf("cannot open file dialog: %w", err))
		}
		return
	}

	cb(dbusWaitForResponse())
}

// The actual implementation for Save. Should be run in a goroutine.
func (fd *fileDialog) save(title string, cb DialogCallback) {
	err := fd.dbusFileChooser(DBusFileChooserSaveFile, title)
	if err != nil {
		if fd.fallback != nil {
			slog.Info("Failed to open linux native file dialog, using fallback", "error", err)
			fd.fallback.Save(title, fd.InitialDirectory(), fd.filters, cb)
		} else {
			cb("", fmt.Errorf("cannot open file dialog: %w", err))
		}
		return
	}

	cb(dbusWaitForResponse())
}

// Call freedesktop via dbus to show a file chooser dialog.
func (fd *fileDialog) dbusFileChooser(method string, title string) error {
	freedesktopFilters := convertFiltersToFreedesktopFilter(fd.filters)

	currentFolder := make([]byte, len(fd.InitialDirectory())+1)
	copy(currentFolder, fd.InitialDirectory())

	options := map[string]dbus.Variant{
		"modal":          dbus.MakeVariant(true),
		"current_folder": dbus.MakeVariant(currentFolder),
		"filters":        dbus.MakeVariant(freedesktopFilters),
	}

	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}

	obj := conn.Object(DBusObjectName, DBusObjectPath)
	err = obj.Call(DBusFileChooserBase+method, 0, "", title, options).Err
	if err != nil {
		return fmt.Errorf("failed to call %s on dbus: %w", method, err)
	}
	return nil
}

func convertFiltersToFreedesktopFilter(filters FileFilters) []freedesktopFilter {
	var result []freedesktopFilter
	for _, filter := range filters {
		var filterRules []freedesktopFilterRule
		for _, rule := range filter.Extensions {
			filterRules = append(filterRules, freedesktopFilterRule{Type: 0, Pattern: "*" + rule})
		}
		result = append(result, freedesktopFilter{Name: filter.Description, Rules: filterRules})
	}
	return result
}

// Wait for the response from the file chooser dialog.
func dbusWaitForResponse() (string, error) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		return "", fmt.Errorf("failed to connect to session bus: %w", err)
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(DBusObjectPath),
		dbus.WithMatchInterface("org.freedesktop.portal.Request"),
		dbus.WithMatchMember("Response"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to subscribe to response signal: %w", err)
	}
	c := make(chan *dbus.Signal)
	conn.Signal(c)

	res := <-c
	if len(res.Body) < 2 {
		return "", fmt.Errorf("invalid response from dbus: %v", res.Body)
	}
	if res.Body[0].(uint32) != 0 {
		// User cancelled the dialog
		return "", nil
	}
	uris := res.Body[1].(map[string]dbus.Variant)["uris"].Value().([]string)
	if len(uris) == 0 {
		return "", nil
	}

	path, _ := strings.CutPrefix(uris[0], "file://")
	path, err = url.PathUnescape(path)
	if err != nil {
		return "", fmt.Errorf("failed to unescape path: %w", err)
	}
	return path, nil
}
