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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/pdf"
)

func NewHVBox(objects ...fyne.CanvasObject) fyne.CanvasObject {
    vBox := container.NewVBox()

    for i := 0; i < len(objects); i++ {
        vBox.Add(objects[i])
    }

    return container.NewHBox(vBox)
}

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

        fileListContainer.Add(canvas.NewText("PDFs to merge", nil))

        for i := 0; i < len(fileList); i++ {
            file := fileList[i]
            if strings.HasSuffix(file.Name(), ".pdf") {
                filesToMerge.PushFront(file.Path())
                newCheckbox := widget.NewCheck(file.Name(), func (checked bool) {
                    slog.Debug("checkbox was clicked")

                    if checked {
                        filesToMerge.PushFront(file.Path())
                    } else {
                        for elem := filesToMerge.Front(); elem != nil; elem = elem.Next() {
                            if elem.Value == file.Path() {
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
            slog.Error("Error merging pdfs", "error", err)
            errorDialog := dialog.NewError(err, myWindow)
            errorDialog.Show()
            return
        }

        if writer == nil {
            slog.Debug("User didn't select file to write to")
            return
        }

        if err := pdf.MergePdfs(*filesToMerge, writer.URI().Path()); err != nil {
            slog.Error("Error merging pdfs", "error", err)
            errorDialog := dialog.NewError(err, myWindow)
            errorDialog.Show()
        } else {
            slog.Info("PDF saved successfully")
            saveConfirmation := dialog.NewInformation("Success!", fmt.Sprintf("Saved merged pdf to %s successfully", writer.URI().Path()), myWindow)
            saveConfirmation.Show()
        }
    }, myWindow)

    mergePdfsButton := widget.NewButton("Merge pdfs", func() {
        saveFileDialog.Show()
    })

    chooseFolderButton := widget.NewButton("Find pdfs", func() {
        slog.Info("User clicked 'Find pdfs'")
        openFolderDialog.Show()
    })

    masterLayout := container.New(layout.NewVBoxLayout(),
        &canvas.Text{
            Text: "PDF merge utility",
            TextSize: 40,
        },
        container.NewHBox(
            NewHVBox(chooseFolderButton),
        ),
        fileListContainer,
        NewHVBox(mergePdfsButton),
    )

    myWindow.SetContent(masterLayout)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}
