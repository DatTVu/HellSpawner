package hswidget

import (
	"fmt"
	image2 "image"
	"image/color"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

// PaletteGridState represents palette grid's state
type PaletteGridState struct {
	// nolint:unused,structcheck // will be used
	loading bool
	// nolint:unused,structcheck // will be used
	failure bool
	texture *giu.Texture
}

// Dispose cleans palette grids state
func (p *PaletteGridState) Dispose() {
	p.texture = nil
}

// PaletteGridWidget represents palette grids widget
type PaletteGridWidget struct {
	id     string
	colors *[256]d2interface.Color
}

// PaletteGrid creates a new palette grid's widget
func PaletteGrid(id string, colors *[256]d2interface.Color) *PaletteGridWidget {
	result := &PaletteGridWidget{
		id:     id,
		colors: colors,
	}

	return result
}

// Build build a new widget
func (p *PaletteGridWidget) Build() {
	var widget *giu.ImageWidget

	stateID := fmt.Sprintf("PaletteGridWidget_%s", p.id)

	state := giu.Context.GetState(stateID)
	if state == nil {
		widget = giu.Image(nil).Size(gridWidth*cellSize, gridHeight*cellSize)

		// Prevent multiple invocation to LoadImage.
		giu.Context.SetState(stateID, &PaletteGridState{})

		rgb := image2.NewRGBA(image2.Rect(0, 0, gridWidth*cellSize, gridHeight*cellSize))

		for y := 0; y < gridHeight*cellSize; y++ {
			if y%cellSize == 0 {
				continue
			}

			for x := 0; x < gridWidth*cellSize; x++ {
				if x%cellSize == 0 {
					continue
				}

				idx := (x / cellSize) + ((y / cellSize) * gridWidth)

				col := p.colors[idx]

				// nolint:gomnd // const
				rgb.Set(x, y, color.RGBA{R: col.R(), G: col.G(), B: col.B(), A: 255})
			}
		}

		go func() {
			texture, err := giu.NewTextureFromRgba(rgb)
			if err == nil {
				giu.Context.SetState(stateID, &PaletteGridState{texture: texture})
			}
		}()
	} else {
		imgState := state.(*PaletteGridState)
		widget = giu.Image(imgState.texture).Size(gridWidth*cellSize, gridHeight*cellSize)
	}

	widget.Build()
}
