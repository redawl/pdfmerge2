package types

import (
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

func NewFileList () (*FileList) {
    dataList := binding.NewURIList()
    fileList := &FileList{}
    fileList = newList(
        dataList.Length,
        func() fyne.CanvasObject {
            fileItem := NewFileItem(nil, func() {}, func() {}, func() {})
            
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
                if err := dataList.Remove(uri); err != nil {
                    slog.Error("Error removing uri", "uri", uri, "error", err)
                    return
                }
            }

            fileItem.MoveUpButton.OnTapped = func() {
                if i != 0 {
                    uriItem, err := dataList.GetItem(i)
                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    uriItem2, err := dataList.GetItem(i-1)
                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    
                    uri, err := uriItem.(binding.URI).Get()

                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    uri2, err := uriItem2.(binding.URI).Get()

                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }

                    if err := dataList.SetValue(i, uri2); err != nil {
                        slog.Error("Error setting value", "error", err)
                        return
                    }

                    if err := dataList.SetValue(i-1, uri); err != nil {
                        slog.Error("Error setting value", "error", err)
                        return
                    }
                }
                fileList.Refresh()
            }
            fileItem.MoveDownButton.OnTapped = func() {
                if i != dataList.Length() - 1 {
                    uriItem, err := dataList.GetItem(i)
                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    uriItem2, err := dataList.GetItem(i+1)
                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    
                    uri, err := uriItem.(binding.URI).Get()

                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }
                    uri2, err := uriItem2.(binding.URI).Get()

                    if err != nil {
                        slog.Error("Error getting item", "error", err)
                        return
                    }

                    if err := dataList.SetValue(i, uri2); err != nil {
                        slog.Error("Error setting value", "error", err)
                        return
                    }

                    if err := dataList.SetValue(i+1, uri); err != nil {
                        slog.Error("Error setting value", "error", err)
                        return
                    }
                }

                fileList.Refresh()
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
        return fmt.Errorf("%s is not a pdf", uri.Path())
    }

    if err := fileList.DataList.Append(uri); err != nil {
        return err
    }

    return nil
}

func (fileList *FileList) GetItem(index int) (fyne.URI, error) {
    value, err := fileList.DataList.GetValue(index)

    if err != nil {
        return nil, err
    }

    return value, err
}

func (fileList *FileList) GetFileNames () ([]string, error) {
    pdfList := make([]string, fileList.DataList.Length())
    for i := range fileList.DataList.Length() {
        uri, err := fileList.GetItem(i)

        if err != nil {
            return nil, err
        }

        pdfList[i] = uri.Path()
    }

    return pdfList, nil
}

