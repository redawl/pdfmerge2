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

    pdfList := make([]string, inPdfs.Length())

    for i := 0; i < inPdfs.Length(); i++ {
        uri, err := inPdfs.GetItem(i)

        if err != nil {
            return err
        }

        pdfList[i] = uri.Path()
    }

    if len(pdfList) == 0 {
        slog.Info("No pdfs where checked in the list")
        return errors.New("Error occurred when saving file. Check log for more info")
    }

    slog.Debug("Merging pdfs", "inPdfs", pdfList, "outPdf", outPdf)
    if err := api.MergeCreateFile(pdfList, outPdf, false, config); err != nil {
        slog.Error("Error merging pdfs", "error", err)
        return fmt.Errorf("Error merging pdfs: Error: %s", err.Error())
    }

    return nil
}

