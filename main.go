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
	"fyne.io/fyne/v2/storage"
	"github.com/redawl/pdfmerge/pdf"
	"github.com/redawl/pdfmerge/types"
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

    fileList := types.NewFileList()

    for _, arg := range flag.Args() {
        if strings.HasSuffix(arg, ".pdf") {
            newUri := storage.NewFileURI(arg)
            fileList.AppendItem(newUri)
        }
    }


    addFilesButton := pdf.AddFilesDialog(window, fileList)

    mergePdfsButton := pdf.SaveFileDialog(window, fileList)

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
            container.NewHBox(addFilesButton),
        ),
        container.NewHBox(mergePdfsButton),
        nil,
        nil,
        fileList,
    )

    window.SetContent(masterLayout)

    slog.Info("Started PDF merge")
    window.ShowAndRun()
}
