package main

import (
	"log"
	"os"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

func generatePDF(
	title string,
	modelLink string,
	imagePath string,
	timeLvl int,
	difficultyLvl int,
	outPath string,
) (err error) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 })
	pdf.AddPage()
	err = pdf.AddTTFFont("Ubuntu", "assets/fonts/Ubuntu-Bold.ttf")
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = pdf.AddTTFFont("NotoSansJP", "assets/fonts/NotoSansJP-Regular.ttf")
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = pdf.SetFont("Ubuntu", "", 38)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = qrcode.WriteFile(modelLink, qrcode.Medium, 529, "assets/pictures/qr.png")
	if err != nil {
		log.Print(err.Error())
		return err
	}

	pageWidth := 595.0

	err = pdf.Image("assets/pictures/epigami.png", 510, 10, &gopdf.Rect{
		W: 73,
		H: 100,
	})
	if err != nil {
		log.Print(err.Error())
		return err
	}

	err = pdf.Image("assets/pictures/epitanime.png", 10, 735, &gopdf.Rect{
		W: 100,
		H: 97,
	})
	if err != nil {
		log.Print(err.Error())
		return err
	}

	epi := "Epitanime - Epigami"
	epiWidth, _ := pdf.MeasureTextWidth(epi)

    	pdf.SetY(80)
    	pdf.SetX((pageWidth - epiWidth) / 2)
    	pdf.Text(epi)

	err = pdf.SetFont("NotoSansJP", "", 28)
	if err != nil {
		log.Print(err.Error())
		return err
	}

    	titleWidth, _ := pdf.MeasureTextWidth(title)

    	pdf.SetY(135)
	pdf.SetX((pageWidth - titleWidth) / 2)
	pdf.Text(title)

	pdf.Image("assets/pictures/qr.png", 148.75, 155, nil)
	err = os.Remove("assets/pictures/qr.png")
	if err != nil {
		log.Print(err.Error())
		return err
	}

	origX := 297.5
	origY := 555.0

	xDiff := 277.5
	yDiff := 100.0

	var actualImagePath = imagePath
	var cleanup func()

	if isURL(imagePath) {
		actualImagePath, err = downloadImage(imagePath)
		if err != nil {
			return err
		}
		cleanup = func() { os.Remove(actualImagePath) }
	}

	if cleanup != nil {
		defer cleanup()
	}

	imagePath = actualImagePath

	w, h, e := getImageSize(imagePath)
	if e != nil {
		log.Print(err.Error())
		return e
	}

	if xDiff / yDiff <= w / h {
		h = h / w * 2 * xDiff
		w = 2 * xDiff
	} else {
		w = w / h * 2 * yDiff
		h = 2 * yDiff
	}

	origX -= w / 2
	origY -= h / 2

	pdf.Image(imagePath, origX, origY, &gopdf.Rect{
		W: w,
		H: h,
	})

	time := "Temps : "
	i := 0
	for ; i < timeLvl; i++ {
		time += "★"
	}
	for ; i < 5; i++ {
		time += "☆"
	}

    	timeWidth, _ := pdf.MeasureTextWidth(time)

    	pdf.SetY(720)
	pdf.SetX((pageWidth - timeWidth) / 2 + 8)
	pdf.Text(time)

	difficulty := "Difficulté : "
	i = 0
	for ; i < difficultyLvl; i++ {
		difficulty += "★"
	}
	for ; i < 5; i++ {
		difficulty += "☆"
	}

    	difficultyWidth, _ := pdf.MeasureTextWidth(difficulty)

    	pdf.SetY(770)
	pdf.SetX((pageWidth - difficultyWidth) / 2 - 7)
	pdf.Text(difficulty)

	pdf.WritePdf(outPath)

	return nil

}

