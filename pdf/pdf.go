package pdf

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/redawl/pdfmerge/types"
)

func MergePdfs(inPdfs *types.FileList, outPdf string) error {
    config := model.NewDefaultConfiguration()

    if len(outPdf) == 0 {
        slog.Error("No outPdf name was specified")
        return errors.New("Error occurred when saving file. Check log for more info")
    } else if inPdfs.Length() == 0 {
        slog.Info("No inPdfs where specified")
        return errors.New("Error occurred when saving file. Check log for more info")
    }

    slice := []string{}
    for i := 0; i < inPdfs.Length(); i++ {
        uriChecked, err := inPdfs.GetItem(i)

        if err != nil {
            return err
        }

        if uriChecked.Checked {
            slice = append(slice, uriChecked.Uri.Path())
        }
    }

    if len(slice) == 0 {
        slog.Info("No pdfs where checked in the list")
        return errors.New("Error occurred when saving file. Check log for more info")
    }

    slog.Debug("Merging pdfs", "inPdfs", slice, "outPdf", outPdf)
    if err := api.MergeCreateFile(slice, outPdf, false, config); err != nil {
        slog.Error("Error merging pdfs", "error", err)
        return errors.New(fmt.Sprintf("Error merging pdfs: Error: %s", err.Error()))
    }

    return nil
}

