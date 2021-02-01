// Package hsdc6editor represents a dc6 editor window
package hsdc6editor

import (
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// DC6Editor represents a dc6 editor
type DC6Editor struct {
	*hseditor.Editor
	dc6           *d2dc6.DC6
	textureLoader *hscommon.TextureLoader
}

// Create creates a new dc6 editor
func Create(textureLoader *hscommon.TextureLoader, pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &DC6Editor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		dc6:           dc6,
		textureLoader: textureLoader,
	}

	return result, nil
}

// Build builds a new dc6 editor
func (e *DC6Editor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.DC6Viewer(e.textureLoader, e.Path.GetUniqueID(), e.dc6),
	})
}

// UpdateMainMenuLayout updates main menu to it contain DC6's editor menu
func (e *DC6Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DC6 Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// GenerateSaveData generates save data
func (e *DC6Editor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves editor's data
func (e *DC6Editor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DC6Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
