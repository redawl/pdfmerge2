package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
    "errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/model"
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
    debugEnabled := flag.Bool("d", false, fmt.Sprintf("Enable debug logging to %s/pdfmerge.log and stdout", os.TempDir()))
    flag.Parse()
    setupLogging(*debugEnabled)

    a := app.New()
    myWindow := a.NewWindow("PDF merge utility")

    filesToMerge := binding.NewUntypedList()

    fileListContainer := widget.NewListWithData(filesToMerge,
        func() fyne.CanvasObject {
            return NewDraggableCheck("template", func(b bool) {})
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
                uriChecked.Checked = b
            }
        },
    )

    openFolderDialog := dialog.NewFolderOpen(func (reader fyne.ListableURI, err error) {
        if err != nil {
            slog.Error("Error occurred during selection of folder", "error", err)
            return
        } else if reader == nil {
            slog.Debug("User clicked cancel or didn't select a folder")
            return
        }

        fileList, err := reader.List()

        if err != nil {
            slog.Debug("Error occurred when retrieving file list", "error", err)
            return
        }

        for i := 0; i < len(fileList); i++ {
            file := fileList[i]
            if strings.HasSuffix(file.Name(), ".pdf") {
                slog.Debug("Found pdf", "name", file.Name())
                filesToMerge.Append(&model.UriChecked{
                    Uri: file,
                    Checked: true,
                })
            }
        }

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

        if err := pdf.MergePdfs(filesToMerge, writer.URI().Path()); err != nil {
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
        checkedCount := 0

        for i := 0; i < filesToMerge.Length(); i++ {
            value, err := filesToMerge.GetValue(i)

            if err != nil {
                slog.Error("Error validating files list", "error", err)
                continue
            }

            if value.(*model.UriChecked).Checked {
                checkedCount++
            }
        }

        if checkedCount == 0 {
            slog.Info("User clicked 'Merge pdfs' without selectnig any pdfs")
            errorDialog := dialog.NewError(errors.New("Please select at least 1 pdf before clicking 'Merge pdfs'"), myWindow)
            errorDialog.Show()
        } else {
            saveFileDialog.Show()
        }
    })

    chooseFolderButton := widget.NewButton("Find pdfs", func() {
        slog.Info("User clicked 'Find pdfs'")
        openFolderDialog.Show()
    })

    headerText := &canvas.Text{
        Text: "PDF merge utility",
        TextSize: 40,
    }

    headerIcon := &canvas.Image{
        Resource: a.Metadata().Icon,
        FillMode: canvas.ImageFillContain,
        ScaleMode: canvas.ImageScaleFastest,
    }

    headerIcon.SetMinSize(fyne.NewSize(headerText.MinSize().Height, headerText.MinSize().Height))

    masterLayout := container.NewBorder(
        container.NewVBox(
            container.NewHBox(
                headerIcon,
                headerText,
            ),
            container.NewHBox(chooseFolderButton),
        ),
        container.NewHBox(mergePdfsButton),
        nil,
        nil,
        fileListContainer,
    )
    myWindow.SetContent(masterLayout)
    myWindow.Resize(fyne.NewSize(800, 600))
    myWindow.ShowAndRun()
}
