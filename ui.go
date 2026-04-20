package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
	"strings"
)

func main() {
	a := app.NewWithID("fr.epigami.model-manager")
	w := a.NewWindow("PDF Generator")
	w.Resize(fyne.NewSize(420, 420))
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title")
	linkEntry := widget.NewEntry()
	linkEntry.SetPlaceHolder("https://...")
	imageLabel := widget.NewLabel("No image selected")
	imageURLEntry := widget.NewEntry()
	imageURLEntry.SetPlaceHolder("Or use URL")
	var imagePath string
	timeSlider := widget.NewSlider(1, 5)
	timeSlider.Step = 1
	timeSlider.Value = 1
	diffSlider := widget.NewSlider(1, 5)
	diffSlider.Step = 1
	diffSlider.Value = 1
	outLabel := widget.NewLabel("No output file")
	var outPath string

	imgFilter := storage.NewExtensionFileFilter([]string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp",
	})

	selectImage := widget.NewButton("Browse locally...", func() {
		fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil {
				return
			}
			imagePath = r.URI().Path()
			imageLabel.SetText(imagePath)
		}, w)
		fd.SetFilter(imgFilter)
		fd.Show()
	})

	selectOut := widget.NewButton("Save PDF As", func() {
		dialog.ShowFileSave(func(wc fyne.URIWriteCloser, err error) {
			if wc == nil {
				return
			}
			outPath = wc.URI().Path()
			if strings.ToLower(filepath.Ext(outPath)) != ".pdf" {
				outPath += ".pdf"
			}
			outLabel.SetText(outPath)
		}, w)
	})

	generate := widget.NewButton("Generate PDF", func() {
		timeLvl := int(timeSlider.Value)
		diffLvl := int(diffSlider.Value)

		if imageURLEntry.Text != "" {
			imagePath = imageURLEntry.Text
		}

		err := generatePDF(
			titleEntry.Text,
			linkEntry.Text,
			imagePath,
			timeLvl,
			diffLvl,
			outPath,
		)

		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Success", "PDF generated successfully!", w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		titleEntry,

		widget.NewLabelWithStyle("Model link", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		linkEntry,

		widget.NewLabelWithStyle("Image", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		selectImage,
		imageLabel,
		imageURLEntry,

		widget.NewLabelWithStyle("Time level", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		timeSlider,

		widget.NewLabelWithStyle("Difficulty level", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		diffSlider,

		widget.NewLabelWithStyle("Output path", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		selectOut,
		outLabel,

		generate,
	))

	w.ShowAndRun()
}