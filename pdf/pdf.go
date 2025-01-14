package pdf

import (
	"errors"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2/data/binding"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
    pdfMergeModel "github.com/redawl/pdfmerge/model"
)

func MergePdfs(inPdfs binding.UntypedList, outPdf string) error {
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
        elem, err := inPdfs.GetValue(i)

        if err != nil {
            return err
        }

        uriChecked := elem.(*pdfMergeModel.UriChecked)
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

