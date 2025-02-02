package types

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FileItem struct {
    widget.BaseWidget
    Label widget.Label
    RemoveButton widget.Button
    MoveUpButton widget.Button
    MoveDownButton widget.Button
    Icon widget.FileIcon
    Uri fyne.URI
    initialX float32
    initialY float32
}

func NewFileItem (uri fyne.URI, onRemove func (), onMoveUp func(), onMoveDown func()) (*FileItem) {
    if uri == nil {
        uri = storage.NewFileURI("")
    }

    fileItem := &FileItem{
        Label: *widget.NewLabel(uri.Name()),
        RemoveButton: *widget.NewButton("Remove", onRemove),
        MoveUpButton: *widget.NewButtonWithIcon("", theme.MoveUpIcon(), onMoveUp),
        MoveDownButton: *widget.NewButtonWithIcon("", theme.MoveDownIcon(), onMoveDown),
        Icon: *widget.NewFileIcon(uri),
        Uri: uri,
        initialX: -1,
        initialY: -1,
    }

    fileItem.ExtendBaseWidget(fileItem)

    return fileItem
}

func (fileItem *FileItem) SetUri (uri fyne.URI) {
    fileItem.Label.SetText(uri.Name())
    fileItem.Label.Refresh()
    fileItem.Icon.SetURI(uri)
    fileItem.Icon.Refresh()
    fileItem.Uri = uri

    fileItem.Refresh()
}

func (fileItem *FileItem) CreateRenderer () fyne.WidgetRenderer {
    c := container.NewBorder(
        nil,
        nil,
        &fileItem.Icon,
        container.NewHBox(
            &fileItem.MoveUpButton,
            &fileItem.MoveDownButton,
            &fileItem.RemoveButton,
        ),
        &fileItem.Label,
    )

    return widget.NewSimpleRenderer(c)
}

