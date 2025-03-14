package pdf

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func MergePdfs(inPdfs []string, outPdf string) error {
    config := model.NewDefaultConfiguration()

    if len(outPdf) == 0 {
        slog.Error("No outPdf name was specified")
        return errors.New("Error occurred when saving file. Check log for more info")
    } else if len(inPdfs) == 0 {
        slog.Error("No inPdfs where specified")
        return errors.New("Error occurred when saving file. Check log for more info")
    }

    slog.Debug("Merging pdfs", "inPdfs", inPdfs, "outPdf", outPdf)
    if err := api.MergeCreateFile(inPdfs, outPdf, false, config); err != nil {
        return fmt.Errorf("Error merging pdfs: Error: %s", err.Error())
    }

    return nil
}

