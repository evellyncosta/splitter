package pdf

import (
	"fmt"
	"os"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

const dir = "dir"
const maximumFileSize = 52428800

func SplitPDF(path string) error {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = pdfcpu.Split(file, "dir", path, 1, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	inFiles := []string{"assimil_1.pdf", "assimil_2.pdf", "assimil_3.pdf"}
	pdfcpu.MergeCreateFile(inFiles, "merged.pdf", nil)

	return nil
}

func GetTotalSplits(path string) error {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return err
	}

	pages, _ := pdfcpu.PageCount(file, nil)
	fmt.Println(pages)

	//arr := []string{"-100"}

	return nil
}
