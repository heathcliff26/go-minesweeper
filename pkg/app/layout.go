package app

import (
	"fyne.io/fyne/v2"
)

// Acts like stack, but enforces a fixed aspect ratio
type FixedRatioLayout struct {
	width, height float32
}

func NewFixedRatioLayout(width, height int) *FixedRatioLayout {
	return &FixedRatioLayout{
		width:  float32(width),
		height: float32(height),
	}
}

// Return the maximum size of the given objects
func (frl *FixedRatioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, o := range objects {
		minSize = minSize.Max(o.MinSize())
	}
	return minSize
}

// Position and resize the given objects to fit the given size, with offset for the x axis
func (frl *FixedRatioLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	target := frl.calculateSize(size)
	xOffset := (size.Width - target.Width) / 2

	for _, obj := range objects {
		obj.Resize(target)
		obj.Move(fyne.NewPos(xOffset, 0))
	}
}

func (frl *FixedRatioLayout) calculateSize(size fyne.Size) fyne.Size {
	targetAspect := float32(frl.width / frl.height)
	aspect := size.Width / size.Height

	if aspect > targetAspect {
		size.Width = size.Height * targetAspect
	} else {
		size.Height = size.Width / targetAspect
	}
	return size
}
