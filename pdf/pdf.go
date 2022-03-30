package pdf

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

const dir = "dir"
const maximumFileSize = 47185920

func SplitPDF(path string) error {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = pdfcpu.Split(file, "", path, 1, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	inFiles := []string{}
	pages, err := pdfcpu.PageCount(file, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fileName := "_"
	fileExtension := ".pdf"
	totalSize := int64(0)

	for i := 1; i <= pages; i++ {
		fileName = strings.Trim(path, fileExtension) + "" + "_" + strconv.Itoa(i) + fileExtension

		fi, err := os.Stat(fileName)
		if err != nil {
			fmt.Println("Err", err)
			return err
		}
		size := fi.Size()

		totalSize = size + totalSize

		if totalSize > maximumFileSize {
			fmt.Println("page: ", i-1)
			fmt.Println(totalSize)
			break
		}

		inFiles = append(inFiles, fileName)

	}

	fmt.Println(inFiles)

	return nil
}
