package service

import (
	"fmt"
	"html/template"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
	"github.com/fsnotify/fsnotify"
	"github.com/whisperverse/whisperverse/config"
	"github.com/whisperverse/whisperverse/model"
)

// Layout service manages the global site layout that is stored in a particular path of the
// filesystem.
type Layout struct {
	folder      config.Folder
	funcMap     template.FuncMap
	analytics   model.Layout
	appearance  model.Layout
	connections model.Layout
	domain      model.Layout
	global      model.Layout
	group       model.Layout
	profile     model.Layout
	topLevel    model.Layout
	user        model.Layout
}

// NewLayout returns a fully initialized Layout service.
func NewLayout(folder config.Folder, funcMap template.FuncMap) Layout {

	return Layout{
		folder:  folder,
		funcMap: funcMap,
	}
}

/*******************************************
 * LAYOUT ACCESSORS
 *******************************************/

func (service *Layout) Analytics() *model.Layout {
	return &service.analytics
}

func (service *Layout) Appearance() *model.Layout {
	return &service.appearance
}

func (service *Layout) Connections() *model.Layout {
	return &service.connections
}

func (service *Layout) Domain() *model.Layout {
	return &service.domain
}

func (service *Layout) Global() *model.Layout {
	return &service.global
}

func (service *Layout) Group() *model.Layout {
	return &service.group
}

func (service *Layout) Profile() *model.Layout {
	return &service.profile
}

func (service *Layout) TopLevel() *model.Layout {
	return &service.topLevel
}
func (service *Layout) User() *model.Layout {
	return &service.user
}

/*******************************************
 * FILE WATCHER
 *******************************************/

// fileNames returns a list of directories that are owned by the Layout service.
func (service *Layout) fileNames() []string {
	return []string{"analytics", "appearance", "connections", "domain", "global", "groups", "profiles", "toplevel", "users"}
}

// watch must be run as a goroutine, and constantly monitors the
// "Updates" channel for news that a template has been updated.
func (service *Layout) Watch() {

	const location = "service.Layout.Watch"

	// Only synchronize on folders that are configured to do so.
	if !service.folder.Sync {
		return
	}

	// Only synchronize on FILESYSTEM folders (for now)
	if service.folder.Adapter != "FILE" {
		return
	}

	// Create a new directory watcher
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		panic(err)
	}

	files := service.fileNames()

	// Use a separate counter because not all files will be included in the result
	for _, filename := range files {

		// Add all other directories into the Template service as Templates
		if err := service.loadFromFilesystem(filename); err != nil {
			derp.Report(derp.Wrap(err, location, "Error loading Layout from filesystem", filename))
			panic("Error loading Layout from Filesystem")
		}

		// Add fsnotify watchers for all other directories
		if err := watcher.Add(service.folder.Location + "/" + filename); err != nil {
			derp.Report(derp.Wrap(err, location, "Error adding file watcher to file", filename))
		}
	}

	// All Files Loaded.  Now Listen for Changes

	// Repeat indefinitely, listen and process file updates
	for {

		select {

		case event, ok := <-watcher.Events:

			if !ok {
				continue
			}

			filename := list.Slash(event.Name).RemoveLast().Last()

			if err := service.loadFromFilesystem(filename); err != nil {
				derp.Report(derp.Wrap(err, location, "Error loading changes to layout", event, filename))
				continue
			}

			fmt.Println("Updated layout: " + filename)

		case err, ok := <-watcher.Errors:

			if ok {
				derp.Report(derp.Wrap(err, location, "Error watching filesystem"))
			}
		}
	}
}

// loadFromFilesystem retrieves the template from the disk and parses it into
func (service *Layout) loadFromFilesystem(filename string) error {

	fs := GetFS(service.folder, filename)

	layout := model.NewLayout(filename, service.funcMap)

	// System folders (except for "static" and "global") have a schema.json file
	if filename != "global" {
		if err := loadModelFromFilesystem(fs, &layout); err != nil {
			return derp.Wrap(err, "service.layout.loadFromFilesystem", "Error loading Schema", fs, filename)
		}
	}

	if err := loadHTMLTemplateFromFilesystem(fs, layout.HTMLTemplate, service.funcMap); err != nil {
		return derp.Wrap(err, "service.layout.loadFromFilesystem", "Error loading Template", fs, filename)
	}

	switch filename {

	case "analytics":
		service.analytics = layout

	case "appearance":
		service.appearance = layout

	case "connections":
		service.connections = layout

	case "domain":
		service.domain = layout

	case "global":
		service.global = layout

	case "groups":
		service.group = layout

	case "profiles":
		service.profile = layout

	case "toplevel":
		service.topLevel = layout

	case "users":
		service.user = layout

	default:
		// This should never happen
		panic("Unrecognized layout: " + filename)
	}

	return nil
}
