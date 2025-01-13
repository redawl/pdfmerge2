package main

import (
	"container/list"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/pdf"
)

func setupLogging(debugEnabled bool) {
    fileWriter, err := os.OpenFile(fmt.Sprintf("%s/%s", os.TempDir(),"pdfmerge.log"), os.O_RDWR, 0666)

    if err != nil {
        slog.Error("Couldn't open log file", "error", err)
        panic("Couldn't open log file")
    }

    defer fileWriter.Close()

    logWriter := io.MultiWriter(os.Stdout, fileWriter)
    logLevel := slog.LevelInfo
    if debugEnabled {
        logLevel = slog.LevelDebug
    }
    logger := slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
        Level: logLevel,
    }))
    slog.SetDefault(logger)
}

func main() {
    debugEnabled := flag.Bool("d", false, "Enable debug logging to TEMP_DIR/pdfmerge.log and stdout")
    flag.Parse()
    setupLogging(*debugEnabled)


    a := app.New()
    myWindow := a.NewWindow("PDF merge utility")

    filesToMerge := list.New()

    saveFileLocation := widget.NewEntry()

    fileListContainer := container.NewVBox()

    openFolderDialog := dialog.NewFolderOpen(func (reader fyne.ListableURI, err error) {
        if err != nil {
            slog.Error("Error occurred during selection of folder", "error", err)
            return
        } else if reader == nil {
            slog.Debug("User clicked cancel or didn't select a folder")
            return
        }

        fileList, err := reader.List()

        fileListContainer.RemoveAll()

        fileListContainer.Add(widget.NewLabel("PDFs to merge"))

        for i := 0; i < len(fileList); i++ {
            file := fileList[i].String()

            if strings.HasSuffix(file, ".pdf") {
                filePath := file[7:]
                lastSlashIndex := strings.LastIndexAny(file, "/")

                newCheckbox := widget.NewCheck(file[lastSlashIndex+1:], func (checked bool) {
                    slog.Debug("checkbox was clicked")

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
                })

                newCheckbox.Checked = true

                fileListContainer.Add(newCheckbox)
            }

            slog.Debug("Found pdf", "name", file)
        }

        fileListContainer.Show()
    }, myWindow)

    saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error){
        if err != nil {
            slog.Error("Error occurred during selection of save file", "error", err)
            return
        } else if writer == nil {
            slog.Debug("User clicked cancel or didn't select a file")
            return
        }
        saveLocation := writer.URI().String()

        filepath := saveLocation[7:]
        if strings.HasSuffix(filepath, ".pdf") {
            saveFileLocation.SetText(filepath)
        } else {
            saveFileLocation.SetText(filepath + ".pdf")
        }
        saveFileLocation.Show()
    }, myWindow)

    mergePdfsButton := widget.NewButton("Merge pdfs", func() {
        if err := pdf.MergePdfs(*filesToMerge, saveFileLocation.Text); err != nil {
            slog.Error("Error merging pdfs", "error", err)
            errorDialog := dialog.NewError(err, myWindow)
            errorDialog.Show()
        } else {
            slog.Info("PDF saved successfully")
            saveConfirmation := dialog.NewInformation("Success!", fmt.Sprintf("Saved merged pdf to %s successfully", saveFileLocation.Text), myWindow)
            saveConfirmation.Show()
        }
    })

    chooseFolderButton := widget.NewButton("Choose folder", func() {
        slog.Info("User clicked 'Choose folder'")
        openFolderDialog.Show()
    })

    chooseSaveFileButton := widget.NewButton("Create save file", func() {
        slog.Info("User clicked 'Create save file'")
        saveFileDialog.Show()
    })

    masterLayout := container.New(layout.NewVBoxLayout(),
        container.NewGridWithColumns(2,
            container.NewVBox(chooseFolderButton), fileListContainer,
            container.NewVBox(chooseSaveFileButton), container.NewVBox(saveFileLocation),
            container.NewVBox(mergePdfsButton),
        ),
    )

    myWindow.SetContent(masterLayout)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}
