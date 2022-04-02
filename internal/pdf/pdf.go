package pdf

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

const maximumFileSize = 47185920
const fileExtension = ".pdf"

type Splitter interface {
	NewSplit(fileName string) error
}

type Split struct {
	Parts    []*Part
	FileName string
}

type Part struct {
	Name       string
	ParentName string
	PartNumber int
	TotalPages int
	LastPage   int
	Pages      []string
}

func (s *Split) NewSplit() error {
	tempDir := os.TempDir()
	os.Chdir(tempDir)
	file, err := os.Open(s.FileName)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = pdfcpu.Split(file, "", s.FileName, 1, nil)
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
		part, err := getPartInfo(currentPage, totalPages, s.FileName, i)
		if err != nil {
			fmt.Println(err)
			return err
		}

		newPart(part)
		if err != nil {
			fmt.Println(err)
			return err
		}

		currentPage = part.LastPage

		isLastPage, err := isLastPage(part.Pages, totalPages)
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

func getPartInfo(currentPage int, totalPages int, fileName string, partNumber int) (*Part, error) {
	inFiles := []string{}
	completeFileName := "_"
	totalSize := int64(0)
	for i := currentPage; i <= totalPages; i++ {
		completeFileName = strings.Trim(fileName, fileExtension) + "" + "_" + strconv.Itoa(i) + fileExtension

		fi, err := os.Stat(completeFileName)

		if err != nil {
			fmt.Println("Err", err)
			return nil, err
		}
		size := fi.Size()

		totalSize = size + totalSize
		currentPage = i
		if totalSize > maximumFileSize {
			break
		}

		inFiles = append(inFiles, completeFileName)
	}
	part := &Part{
		PartNumber: partNumber,
		LastPage:   currentPage,
		Pages:      inFiles,
		ParentName: fileName,
	}
	return part, nil
}

func newPart(part *Part) error {
	fileName := strings.Trim(part.ParentName, fileExtension)
	part.Name = fileName + "part" + strconv.Itoa(part.PartNumber) + fileExtension
	err := pdfcpu.MergeCreateFile(part.Pages, part.Name, nil)
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
