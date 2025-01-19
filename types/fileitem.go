package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type FileItem struct {
    widget.BaseWidget
    Label widget.Label
    RemoveButton widget.Button
    Icon widget.FileIcon
    Uri fyne.URI
    initialX float32
    initialY float32
    OnDragged func(event *fyne.DragEvent)
}

func NewFileItem (uri fyne.URI, onRemove func ()) (*FileItem) {
    if uri == nil {
        uri = storage.NewFileURI("/pdf.pdf")
    }
    fileItem := &FileItem{
        Label: *widget.NewLabel(uri.Name()),
        RemoveButton: *widget.NewButton("Remove", onRemove),
        Icon: *widget.NewFileIcon(uri),
        Uri: uri,
        initialX: -1,
        initialY: -1,
    }

    fileItem.ExtendBaseWidget(fileItem)

    return fileItem
}

func (fileItem *FileItem) SetUri (uri fyne.URI) {
    fileItem.Label = *widget.NewLabel(uri.Name())
    fileItem.Icon.SetURI(uri)
    fileItem.Uri = uri

    fileItem.Refresh()
}

func (fileItem *FileItem) CreateRenderer () fyne.WidgetRenderer {
    c := container.NewBorder(nil, nil, &fileItem.Icon, &fileItem.RemoveButton, &fileItem.Label)

    return widget.NewSimpleRenderer(c)
}

func (fileItem *FileItem) Dragged(event *fyne.DragEvent) {
    if fileItem.OnDragged != nil {
        fileItem.OnDragged(event)
    } else {
        if fileItem.initialX == -1 && fileItem.initialY == -1{
            fileItem.initialX = fileItem.Position().X
            fileItem.initialY = fileItem.Position().Y
        }
        fileItem.Move(fyne.NewPos(fileItem.Position().X + event.Dragged.DX, fileItem.Position().Y + event.Dragged.DY))
    }
}

func (fileItem *FileItem) DragEnd() {
    fileItem.Move(fyne.NewPos(fileItem.initialX, fileItem.initialY))

    fileItem.initialX = -1
    fileItem.initialY = -1
}

