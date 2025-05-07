//go:build linux

package filedialog

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	ObjectName = "org.freedesktop.portal.Desktop"
	ObjectPath = "/org/freedesktop/portal/desktop"

	FileChooserBase     = "org.freedesktop.portal.FileChooser"
	FileChooserOpenFile = ".OpenFile"
	FileChooserSaveFile = ".SaveFile"
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

func convertFilters(filters FileFilter) []freedesktopFilter {
	var result []freedesktopFilter
	for name, rules := range filters {
		var filterRules []freedesktopFilterRule
		for _, rule := range rules {
			filterRules = append(filterRules, freedesktopFilterRule{Type: 0, Pattern: "*" + rule})
		}
		result = append(result, freedesktopFilter{Name: name, Rules: filterRules})
	}
	return result
}

func dbusFileChooser(method string, title string, startLocation string, filters FileFilter) error {
	freedesktopFilters := convertFilters(filters)

	currentFolder := make([]byte, len(startLocation)+1)
	copy(currentFolder, startLocation)

	options := map[string]dbus.Variant{
		"modal":          dbus.MakeVariant(true),
		"current_folder": dbus.MakeVariant(currentFolder),
		"filters":        dbus.MakeVariant(freedesktopFilters),
	}

	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}

	obj := conn.Object(ObjectName, ObjectPath)
	err = obj.Call(FileChooserBase+method, 0, "", title, options).Err
	if err != nil {
		return fmt.Errorf("failed to call %s on dbus: %w", method, err)
	}
	return nil
}

func dbusWaitForResponse() (string, error) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		return "", fmt.Errorf("failed to connect to session bus: %w", err)
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(ObjectPath),
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
	return path, nil
}
