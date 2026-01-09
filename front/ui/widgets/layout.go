package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomVBox struct{}

func (b *CustomVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		
		childSize := o.MinSize()
		
		w = max(childSize.Width, w)
		h += childSize.Height + theme.Padding()
	}
	return fyne.NewSize(w, h)
}

func (b *CustomVBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		
		minSize := o.MinSize()
		o.Resize(fyne.NewSize(containerSize.Width, minSize.Height))
		o.Move(pos)
		
		pos = pos.Add(fyne.NewPos(0, minSize.Height+theme.Padding()))
	}
}
