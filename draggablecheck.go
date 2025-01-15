package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DraggableCheck struct {
    widget.Check
    initialX float32;
    initialY float32;
}

func NewDraggableCheck (text string, onChecked func (b bool)) (*DraggableCheck) {
    check := &DraggableCheck{}
    check.ExtendBaseWidget(check)

    check.Checked = false
    check.OnChanged = onChecked
    check.Text = text
    check.initialX = -1
    check.initialY = -1

    return check
}

func (check *DraggableCheck) Dragged(event *fyne.DragEvent) {
    if check.initialX == -1 && check.initialY == -1{
        check.initialX = check.Position().X
        check.initialY = check.Position().Y
    }
    check.Move(fyne.NewPos(check.Position().X + event.Dragged.DX, check.Position().Y + event.Dragged.DY))
}

func (check *DraggableCheck) DragEnd() {
    check.Move(fyne.NewPos(check.initialX, check.initialY))

    check.initialX = -1
    check.initialY = -1
}

