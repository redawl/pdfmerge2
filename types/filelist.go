package types

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type FileList struct {
    widget.List
    DataList binding.URIList
}

func newList(length func() int, createItem func() fyne.CanvasObject, updateItem func(int, fyne.CanvasObject)) *FileList {
	list := &FileList{}
	list.Length = length
	list.CreateItem = createItem
	list.UpdateItem = updateItem
	list.ExtendBaseWidget(list)

	return list
}

func IncrementIntBinding(intBinding binding.Int) {
    value, err := intBinding.Get()
    if err != nil {
        slog.Error("Error!", "error", err)
        panic(err)
    }

    intBinding.Set(value + 1)
}

func DecrementIntBinding(intBinding binding.Int) {
    value, err := intBinding.Get()
    if err != nil {
        slog.Error("Error!", "error", err)
        panic(err)
    }

    intBinding.Set(value - 1)
}

func NewFileList () (*FileList) {
    dataList := binding.NewURIList()

    fileList := newList(
        dataList.Length,
        func() fyne.CanvasObject {
            fileItem := NewFileItem(nil, func() {})
            
            return fileItem
        },
        func(i int, o fyne.CanvasObject) {
			uri, err := dataList.GetValue(i)
            if err != nil {
                slog.Error("Error Getting Uri", "error", err)
                return
            }

            fileItem := o.(*FileItem)
            fileItem.SetUri(uri)
            fileItem.RemoveButton.OnTapped = func () {
                dataList.Remove(uri)
            }
		},
    )

	dataList.AddListener(binding.NewDataListener(fileList.Refresh))

    fileList.DataList = dataList

    fileList.ExtendBaseWidget(fileList)
    fileList.Show()

    return fileList
}

func (fileList *FileList) AppendItem(uri fyne.URI) error {
    if !strings.HasSuffix(uri.Path(), ".pdf") {
        return errors.New(fmt.Sprintf("%s is not a pdf", uri.Path()))
    }

    fileList.DataList.Append(uri)

    return nil
}

func (fileList *FileList) GetItem(index int) (fyne.URI, error) {
    value, err := fileList.DataList.GetValue(index)

    if err != nil {
        return nil, err
    }

    return value, err
}

