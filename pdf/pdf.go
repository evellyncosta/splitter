package pdf

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

const dir = "tmp/"
const maximumFileSize = 47185920
const fileExtension = ".pdf"

func SplitPDF(fileName string) error {
	os.Chdir(dir)
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = pdfcpu.Split(file, "", fileName, 1, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	totalPages, err := pdfcpu.PageCount(file, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		fmt.Println("Err", err)
		return err
	}
	currentPage := 1

	for i := 1; i < totalPages; i++ {
		partsNames, lastPage, err := extractPartsNames(currentPage, totalPages, fileName, i)
		if err != nil {
			fmt.Println(err)
			return err
		}
		mergePages(partsNames, fileName, i)
		if err != nil {
			fmt.Println(err)
			return err
		}

		currentPage = lastPage

		isLastPage, err := isLastPage(partsNames, totalPages)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if isLastPage {
			break
		}

	}

	return nil
}

func extractPartsNames(currentPage int, totalPages int, fileName string, partNumber int) ([]string, int, error) {
	inFiles := []string{}
	completeFileName := "_"
	totalSize := int64(0)
	for i := currentPage; i <= totalPages; i++ {
		completeFileName = strings.Trim(fileName, fileExtension) + "" + "_" + strconv.Itoa(i) + fileExtension

		fi, err := os.Stat(completeFileName)
		if err != nil {
			fmt.Println("Err", err)
			return nil, 0, err
		}
		size := fi.Size()

		totalSize = size + totalSize
		currentPage = i
		if totalSize > maximumFileSize {
			break
		}

		inFiles = append(inFiles, completeFileName)
	}
	partNumber++

	return inFiles, currentPage, nil
}

func mergePages(parts []string, fileName string, partNumber int) error {
	fileName = strings.Trim(fileName, fileExtension)
	fileName = fileName + "part" + strconv.Itoa(partNumber) + fileExtension
	err := pdfcpu.MergeCreateFile(parts, fileName, nil)
	if err != nil {
		return err
	}

	return nil
}

func isLastPage(parts []string, totalPage int) (bool, error) {
	lastPart := parts[len(parts)-1]
	re := regexp.MustCompile("[0-9]+")
	page := re.FindAllString(lastPart, -1)
	pageConverted, err := strconv.Atoi(page[0])
	if err != nil {
		return false, err
	}
	if pageConverted == totalPage {
		return true, nil
	}

	return false, nil
}
