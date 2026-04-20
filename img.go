package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func getImageSize(path string) (w, h float64, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return float64(cfg.Width), float64(cfg.Height), nil
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func downloadImage(imgURL string) (string, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	tmpFile, err := os.CreateTemp("", "epigami-img-*"+filepath.Ext(imgURL))

	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = io.Copy(tmpFile, resp.Body)

	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

