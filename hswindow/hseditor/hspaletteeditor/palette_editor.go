package hspaletteeditor

import (
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// PaletteEditor represents a palette editor
type PaletteEditor struct {
	*hseditor.Editor
	palette d2interface.Palette
}

// Create creates a new palette editor
func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		Editor:  hseditor.New(pathEntry, x, y, project),
		palette: palette,
	}

	return result, nil
}

// Build builds a palette editor
func (e *PaletteEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.PaletteGrid(e.GetID()+"_grid", e.palette.GetColors()),
	})
}

// UpdateMainMenuLayout updates a main menu layout to it contain pallete editor's options
func (e *PaletteEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Editor").Layout(g.Layout{
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

// GenerateSaveData generates data to be saved
func (e *PaletteEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves editor
func (e *PaletteEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides palette editor
func (e *PaletteEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
