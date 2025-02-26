package pdf

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/redawl/pdfmerge/types"
	nativedialog "github.com/tawesoft/golib/v2/dialog"
)

func SaveFileDialog(window fyne.Window, filesList *types.FileList) (*widget.Button) {
    saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error){
        if err != nil {
            slog.Error("Error merging pdfs", "error", err)
            errorDialog := dialog.NewError(err, window)
            errorDialog.Show()
            return
        }

        if writer == nil {
            slog.Debug("User didn't select file to write to")
            return
        }

        if err := MergePdfs(filesList, writer.URI().Path()); err != nil {
            slog.Error("Error merging pdfs", "error", err)
            errorDialog := dialog.NewError(err, window)
            errorDialog.Show()
        } else {
            slog.Debug("PDF saved successfully")
            saveConfirmation := dialog.NewInformation("Success!", fmt.Sprintf("Saved merged pdf to %s successfully", writer.URI().Path()), window)
            saveConfirmation.Show()
        }
    }, window)

    mergePdfsButton := widget.NewButton("Merge", func() {
        if filesList.Length() == 0 {
            slog.Debug("User clicked 'Merge' without selectnig any pdfs")
            errorDialog := dialog.NewError(errors.New("Select at least 1 pdf before clicking 'Merge'"), window)
            errorDialog.Show()
        } else {
            supported, err := nativedialog.Supported()

            if err != nil {
                slog.Error("Error checking native dialog support, falling back to fyne native", "error", err)
                saveFileDialog.Show()
                return;
            }

            if !supported.FilePicker {
                slog.Debug("Native filepicker not supported, falling back to fyne native")
                saveFileDialog.Show()
                return
            }

            saveFile, success, err := nativedialog.Save("merged.pdf")
            
            if err != nil {
                slog.Error("Error using native filepicker", "error", err)
            } else if !success {
                slog.Error("User clicked cancel or didn't pick any files")
            } else {
                if err := MergePdfs(filesList, saveFile); err != nil {
                    slog.Error("Error merging pdfs", "error", err)
                    errorDialog := dialog.NewError(err, window)
                    errorDialog.Show()
                } else {
                    slog.Debug("PDF saved successfully")
                    saveConfirmation := dialog.NewInformation("Success!", fmt.Sprintf("Saved merged pdf to %s successfully", saveFile), window)
                    saveConfirmation.Show()
                }
            }
        }
    })

    return mergePdfsButton
}

func AddFilesDialog(window fyne.Window, filesList *types.FileList) (*widget.Button) {
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
                filesList.AppendItem(file)
            }
        }

    }, window)

    addFilesButton := widget.NewButton("Add files", func() {
        slog.Debug("User clicked 'Add files'")

        supported, err := nativedialog.Supported()

        if err != nil {
            slog.Error("Error checking for native dialog support, falling back to fyne native", "error", err)
            openFolderDialog.Show()
            return
        }

        if !supported.FilePicker {
            slog.Debug("Native dialog not supported, falling back to fyne native")
            openFolderDialog.Show()
            return
        }

        fileList, success, err := nativedialog.FilePicker{
            FileTypes: [][2]string{
                {"PDF Document", "*.pdf"},
            },
        }.OpenMultiple()


        if err != nil {
            slog.Error("Error using native filepicker", "error", err)
        } else if !success {
            slog.Error("Native filepicker not supported, falling back to fyne native")
        } else {
            if len(fileList) == 0 {
                slog.Debug("User clicked cancel or didn't select a folder")
                return
            }

            for i := 0; i < len(fileList); i++ {
                file := fileList[i]
                if strings.HasSuffix(file, ".pdf") {
                    slog.Debug("Found pdf", "name", file)
                    filesList.AppendItem(storage.NewFileURI(file))
                }
            }
        }
    })

    return addFilesButton
}

