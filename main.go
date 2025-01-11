package main

import (
	"container/list"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
    "errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
    fileWriter, err := os.Create(fmt.Sprintf("%s/%s", os.TempDir(),"pdfmerge.log"))

    if err != nil {
        slog.Error("Couldn't open log file", "error", err)
        panic("Couldn't open log file")
    }

    logWriter := io.MultiWriter(os.Stdout, fileWriter)

    logger := slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    slog.SetDefault(logger)

    config := model.NewDefaultConfiguration()
    a := app.New()
    myWindow := a.NewWindow("PDF merge utility")

    filesToMerge := list.New()

    saveFileLocation := widget.NewLabel("")
    saveFileLocation.Hide()

    form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() {
            if len(saveFileLocation.Text) == 0 {
                slog.Info("User clicked 'Merge pdfs' before selecting save file location")
                errorDialog := dialog.NewError(errors.New("Choose a save file location before clicking 'Merge pdfs'"), myWindow)
                errorDialog.Show()
                return
            } else if filesToMerge.Front() == nil {
                slog.Info("User clicked 'Merge pdfs' before selecting pdfs")
                errorDialog := dialog.NewError(errors.New("Select at least one pdf before clicking 'Merge pdfs'"), myWindow)
                errorDialog.Show()
                return
            }
            
            slice := []string{}
            for elem := filesToMerge.Front(); elem != nil; elem = elem.Next() {
                slice = append(slice, elem.Value.(string))
            }

            if err := api.MergeCreateFile(slice, saveFileLocation.Text, false, config); err != nil {
                slog.Error("Error merging pdfs", "error", err)
                errorDialog := dialog.NewError(errors.New(fmt.Sprintf("Error merging pdfs: Error: %s", err.Error())), myWindow)
                errorDialog.Show()
            } else {
                slog.Info("PDF saved successfully")
                saveConfirmation := dialog.NewInformation("Success!", fmt.Sprintf("Saved merged pdf to %s successfully", saveFileLocation.Text), myWindow)
                saveConfirmation.Show()
            }
	    },
        SubmitText: "Merge pdfs",
    }

    openFolderDialog := dialog.NewFolderOpen(func (reader fyne.ListableURI, err error) {
        fileList, err := reader.List()

        form.Append("", widget.NewLabel("PDFs to merge"))

        for i := 0; i < len(fileList); i++ {
            file := fileList[i].String()

            if strings.HasSuffix(file, ".pdf") {
                filePath := file[7:]
                lastSlashIndex := strings.LastIndexAny(file, "/")

                form.Append("", widget.NewCheck(file[lastSlashIndex+1:], func (checked bool) {
                    slog.Info("checkbox was clicked")

                    if checked {
                        filesToMerge.PushFront(filePath)
                    } else {
                        for elem := filesToMerge.Front(); elem != nil; elem = elem.Next() {
                            if elem.Value == filePath {
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

    saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error){
        saveLocation := writer.URI().String()

        filepath := saveLocation[7:]
        saveFileLocation.SetText(filepath)
        saveFileLocation.Show()
    }, myWindow)

    openFolderDialog.Hide()
    saveFileDialog.Hide()

    chooseFolderButton := widget.NewButton("Choose folder", func() {
        slog.Info("User clicked 'Choose folder'")
        openFolderDialog.Show()
    })

    chooseSaveFileButton := widget.NewButton("Create save file", func() {
        slog.Info("User clicked 'Create save file'")
        saveFileDialog.Show()
    })

    form.Append("", chooseFolderButton)
    form.Append("", saveFileLocation)
    form.Append("", chooseSaveFileButton)

    myWindow.SetContent(form)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()

    // Cleanup
    if err := fileWriter.Close(); err != nil {
        slog.Error("Error closing logfile", "error", err)
        return
    }
}
