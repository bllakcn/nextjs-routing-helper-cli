package main

import (
	"github.com/bllakcn/nextjs-routing-helper-cli/cmd"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	figure.NewFigure("Next.js Routing Helper", "rectangles", true).Print()
	cmd.Execute()
}
