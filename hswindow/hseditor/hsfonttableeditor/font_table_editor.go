// Package hsfonttableeditor represents fontTableEditor's window
package hsfonttableeditor

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
)

type fontTable map[rune]*fontGlyph

type fontGlyph struct {
	rune
	frameIndex int
	width      int
}

// FontTableEditor represents font table editor
type FontTableEditor struct {
	*hseditor.Editor
	fontTable
	rows g.Rows
}

// Create creates a new font table editor
func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	glyphs := make(fontTable)

	table := *data

	const (
		numHeaderBytes = 12
		bytesPerGlyph  = 14
	)

	for i := numHeaderBytes; i < len(table); i += bytesPerGlyph {
		chr := rune(binary.LittleEndian.Uint16(table[i : i+2]))

		glyphs[chr] = &fontGlyph{
			rune:       chr,
			frameIndex: int(binary.LittleEndian.Uint16(table[i+8 : i+10])),
			width:      int(table[i+3]),
		}
	}

	editor := &FontTableEditor{
		Editor:    hseditor.New(pathEntry, x, y, project),
		fontTable: glyphs,
	}

	return editor, nil
}

func (e *FontTableEditor) init() {
	e.rows = make(g.Rows, 0)

	e.rows = append(e.rows, g.Row(
		g.Label("Index"),
		g.Label("Character"),
		g.Label("Width (px)"),
	))

	// so that we can sort the glyphs
	glyphs := make([]*fontGlyph, len(e.fontTable))

	idx := 0

	for _, glyph := range e.fontTable {
		glyphs[idx] = glyph
		idx++
	}

	sort.Slice(glyphs, func(i, j int) bool {
		return glyphs[i].frameIndex < glyphs[j].frameIndex
	})

	for idx := range glyphs {
		e.rows = append(e.rows, e.makeGlyphLayout(glyphs[idx]))
	}
}

// Build builds a font table editor's window
func (e *FontTableEditor) Build() {
	if e.rows == nil {
		e.init()
		return
	}

	tableLayout := g.Layout{g.Child("").
		Border(false).
		Layout(
			g.Layout{
				g.FastTable("").Border(true).Rows(e.rows),
			},
		)}

	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Size(mainWindowW, mainWindowH).
		Layout(tableLayout)
}

func (e *FontTableEditor) makeGlyphLayout(glyph *fontGlyph) *g.RowWidget {
	if glyph == nil {
		return &g.RowWidget{}
	}

	row := g.Row(
		g.Label(fmt.Sprintf("%d", glyph.frameIndex)),
		g.Label(string(glyph.rune)),
		g.Label(fmt.Sprintf("%d", glyph.width)),
	)

	return row
}

// UpdateMainMenuLayout updates mainMenu layout's to it contain FontTableEditor's options
func (e *FontTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Table Editor").Layout(g.Layout{
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
func (e *FontTableEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *FontTableEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *FontTableEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
