package hsaboutdialog

import (
	"image/color"
	"io/ioutil"
	"log"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

type AboutDialog struct {
	*hsdialog.Dialog
	titleFont   imgui.Font
	regularFont imgui.Font
	fixedFont   imgui.Font
	credits     string
	license     string
}

func Create(regularFont, titleFont, fixedFont imgui.Font) (*AboutDialog, error) {
	result := &AboutDialog{
		Dialog:      hsdialog.New("About HellSpawner"),
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}
	var err error
	var data []byte

	if data, err = ioutil.ReadFile("LICENSE"); err != nil {
		log.Fatal(err)
	}
	result.license = string(data)

	if data, err = ioutil.ReadFile("CONTRIBUTORS"); err != nil {
		log.Fatal(err)
	}
	result.credits = string(data)

	return result, nil
}

func (a *AboutDialog) Build() {
	a.IsOpen(&a.Visible).Layout(g.Layout{
		g.Line(
			g.ImageWithFile("hsassets/images/d2logo.png").Size(256, 256),
			g.Child("AboutHellSpawnerLayout").Size(500, -1).Layout(g.Layout{
				g.Label("HellSpawner").Color(&color.RGBA{R: 255, G: 255, B: 255, A: 255}).Font(&a.titleFont),
				g.Label("The OpenDiablo 2 Toolset").Color(&color.RGBA{R: 255, G: 255, B: 255, A: 255}).Font(&a.regularFont),
				g.Label("Local Build").Color(&color.RGBA{R: 255, G: 255, B: 255, A: 255}).Font(&a.fixedFont),
				g.Separator(),
				g.TabBar("AboutHellSpawnerTabBar").Flags(g.TabBarFlagsNoCloseWithMiddleMouseButton).Layout(g.Layout{
					g.TabItem("Credits##AboutHellSpawner").Layout(g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("##AboutHellSpawnerCredits", &a.credits).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { g.PopFont() }),
					}),
					g.TabItem("Licenses##AboutHellSpawner").Layout(g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("##AboutHellSpawnerLicense", &a.license).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { g.PopFont() }),
					}),
				}),
			}),
		),
	})
}
