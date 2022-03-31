package main

import (
	"splliter/internal/pdf"
)

func main() {
	split := &pdf.Split{FileName: "assimil.pdf"}
	split.NewSplit()
}
