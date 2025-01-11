package main

import (
	"log/slog"
	"strings"
    "container/list"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
    slog.SetLogLoggerLevel(slog.LevelDebug)

    a := app.New()
    myWindow := a.NewWindow("PDF merge utility")

    filesToMerge := list.New()

    form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() {
            for elem := filesToMerge.Front(); elem.Next() != nil; elem = elem.Next() {
                slog.Debug("Would merge pdf", "name", elem.Value)
            }
	    },
        SubmitText: "Merge pdfs",
    }

    fileDialog := dialog.NewFolderOpen(func (reader fyne.ListableURI, err error) {
        fileList, err := reader.List()

        form.Append("", widget.NewLabel("PDFs to merge"))

        for i := 0; i < len(fileList); i++ {
            file := fileList[i].String()

            if strings.HasSuffix(file, ".pdf") {
                lastSlashIndex := strings.LastIndexAny(file, "/")

                form.Append("", widget.NewCheck(file[lastSlashIndex+1:], func (checked bool) {
                    slog.Debug("checkbox was clicked")

                    if checked {
                        filesToMerge.PushFront(file)
                    } else {
                        for elem := filesToMerge.Front(); elem.Next() != nil; elem = elem.Next() {
                            if elem.Value == file {
                                filesToMerge.Remove(elem)
                                break
                            }
                        }
                    }
                }))
            }

            slog.Info("Filename", "name", file)
        }
    }, myWindow)

    fileDialog.Hide()

    customButton := widget.NewButton("Choose folder", func() {
        slog.Debug("User clicked 'Choose folder")
        fileDialog.Show()
    })

    form.Append("", customButton)

    myWindow.SetContent(form)
    myWindow.ShowAndRun()
}
