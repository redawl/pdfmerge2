package types

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/model"
)

type FileList struct {
    widget.List
    DataList binding.UntypedList
    CheckedCount binding.Int
}

func newList(length func() int, createItem func() fyne.CanvasObject, updateItem func(int, fyne.CanvasObject)) *FileList {
	list := &FileList{}
	list.Length = length
	list.CreateItem = createItem
	list.UpdateItem = updateItem
	list.ExtendBaseWidget(list)

	return list
}
func newListWithData(data binding.DataList, createItem func() fyne.CanvasObject, updateItem func(binding.DataItem, fyne.CanvasObject)) *FileList {
	l := newList(
		data.Length,
		createItem,
        func(i int, o fyne.CanvasObject) {
			item, err := data.GetItem(i)
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
				return
			}
			updateItem(item, o)
		})

	data.AddListener(binding.NewDataListener(l.Refresh))
	return l
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
    dataList := binding.NewUntypedList()
    checkedCount := binding.NewInt()
    fileList := newListWithData(
        dataList,
        func() fyne.CanvasObject {
            check := NewDraggableCheck("template", func(b bool) {})
            
            return check
        },
        func(di binding.DataItem, co fyne.CanvasObject) {
            uriBinding := di.(binding.Untyped)

            value, err := uriBinding.Get()

            if err != nil {
                slog.Error("Error Getting Uri", "error", err)
                return
            }

            uriChecked := value.(*model.UriChecked)

            checkBox := co.(*DraggableCheck)

            checkBox.SetText(uriChecked.Uri.Name())
            checkBox.Checked = uriChecked.Checked
            checkBox.OnChanged = func(b bool) {
                if uriChecked.Checked != b {
                    if b {
                        IncrementIntBinding(checkedCount) 
                    } else {
                        DecrementIntBinding(checkedCount)
                    }
                }
                uriChecked.Checked = b
            }
            // Call refresh to ensure checkbox is updated
            // with visual state
            checkBox.Refresh() 
        },
    )

    fileList.DataList = dataList
    fileList.CheckedCount = checkedCount

    fileList.ExtendBaseWidget(fileList)
    fileList.Show()

    return fileList
}

func (fileList *FileList) AppendItem(uri fyne.URI) error {
    if !strings.HasSuffix(uri.Path(), ".pdf") {
        return errors.New(fmt.Sprintf("%s is not a pdf", uri.Path()))
    }

    fileList.DataList.Append(&model.UriChecked{
        Uri: uri,
        Checked: true,
    })

    IncrementIntBinding(fileList.CheckedCount)

    return nil
}

func (fileList *FileList) GetItem(index int) (*model.UriChecked, error) {
    value, err := fileList.DataList.GetValue(index)

    if err != nil {
        return nil, err
    }

    return value.(*model.UriChecked), err
}

