package main

import (
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
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/model"
	"github.com/redawl/pdfmerge/pdf"
)

func setupLogging(debugEnabled bool) {
    fileWriter, err := os.OpenFile(fmt.Sprintf("%s/pdfmerge.log", os.TempDir()), os.O_RDWR | os.O_CREATE, 0666)

    if err != nil {
        panic("Couldn't open log file")
    }

    logWriter := io.MultiWriter(fileWriter, os.Stdout)
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
    window := a.NewWindow("PDF merge utility")
    window.Resize(fyne.NewSize(800, 600))

    filesToMerge := binding.NewUntypedList()

    for _, arg := range flag.Args() {
        if strings.HasSuffix(arg, ".pdf") {
            newUri := storage.NewFileURI(arg)
            filesToMerge.Append(&model.UriChecked{
                Uri: newUri,
                Checked: true,
            })
        }
    }

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
            // Call refresh to ensure checkbox is updated 
            // with visual state
            checkBox.Refresh()
        },
    )

    addFilesButton, fileCountLabel := pdf.AddFilesDialog(window, filesToMerge)

    mergePdfsButton := pdf.SaveFileDialog(window, filesToMerge)

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
            container.NewHBox(addFilesButton, fileCountLabel),
        ),
        container.NewHBox(mergePdfsButton),
        nil,
        nil,
        fileListContainer,
    )

    window.SetContent(masterLayout)

    slog.Info("Started PDF merge")
    window.ShowAndRun()
}
