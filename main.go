package main

import "splliter/pdf"

func main() {
	split := &pdf.Split{FileName: "assimil.pdf"}
	split.NewSplit()
}
