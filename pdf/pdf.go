package pdf

import (
    "fmt"
    "container/list"
    "log/slog"
    "errors"
    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func MergePdfs(inPdfs list.List, outPdf string) error {
    config := model.NewDefaultConfiguration()

    if len(outPdf) == 0 {
        slog.Info("User clicked 'Merge pdfs' before selecting save file location")
        return errors.New("Choose a save file location before clicking 'Merge pdfs'")
    } else if inPdfs.Front() == nil {
        slog.Info("User clicked 'Merge pdfs' before selecting pdfs")
        return errors.New("Select at least one pdf before clicking 'Merge pdfs'")
    }

    slice := []string{}
    for elem := inPdfs.Front(); elem != nil; elem = elem.Next() {
        slice = append(slice, elem.Value.(string))
    }

    if err := api.MergeCreateFile(slice, outPdf, false, config); err != nil {
        slog.Error("Error merging pdfs", "error", err)
        return errors.New(fmt.Sprintf("Error merging pdfs: Error: %s", err.Error()))
    }

    return nil
}

