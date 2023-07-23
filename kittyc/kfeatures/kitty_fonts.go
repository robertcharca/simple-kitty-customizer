package kfeatures

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/robertcharca/skittyc/kittyc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type urlFont struct {
	url string
}

func (u urlFont) verifyUrlFontDownload() (int, bool) {
	resp, err := http.Get(u.url)
	if err != nil {
		log.Fatalln(err)
	}
	
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return resp.StatusCode, true
	}

	return resp.StatusCode, false
}

func (z urlFont) verifyZipFontDowload() bool {
	if strings.Contains(z.url, ".zip") {
		return true
	}

	return false
}

// verifyFontDownload: compares three download alternatives and checks if the status is between 200 and 299
func verifyFontDownload(font string) (bool, string, bool) {
	var (
		corrFont string
		firstUrl string
		secondUrl string
	)

	corrFont = strings.ReplaceAll(font, " ", "-")

	firstUrl = "https://www.1001fonts.com/download/" + corrFont + ".zip"
	secondUrl = "https://www.fontsquirrel.com/fonts/download/" + corrFont
	
	urlFistAlt := urlFont{firstUrl}
	urlSecondAlt := urlFont{secondUrl}
	urlThirdAlt := urlFont{font}

	if _, resp := urlFistAlt.verifyUrlFontDownload(); resp {
		return true, firstUrl, urlFistAlt.verifyZipFontDowload()
	} else if _, resp := urlSecondAlt.verifyUrlFontDownload(); resp {
		return true, secondUrl, urlSecondAlt.verifyZipFontDowload()
	} else if _, resp := urlThirdAlt.verifyUrlFontDownload(); resp {
		return true, font, urlThirdAlt.verifyZipFontDowload()
	}

	return false, "", false
}

func downloadFontZip(font string) (string, bool, string) {
	var (
		fileName string
		fontsPath string
	)

	homePath, _ := os.UserHomeDir()
	fontsPath = homePath + "/.local/share/fonts/"

	// Verifying if the font can be downloaded 
	fontStatus, download, zip := verifyFontDownload(font)

	if !fontStatus {
		fmt.Println("This font cannot be downloaded.")
		return "", false, ""
	}

	// Getting the download url
	fileURL, err := url.Parse(download)
	if err != nil {
		log.Fatalln(err)
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")

	// Verifying if it's a .zip file for creating it
	if zip {
		fileName = segments[len(segments)-1]
	} else {
		fileName = segments[len(segments)-1] + ".zip"
	}
	
	file, err := os.Create(fontsPath + fileName)
	if err != nil {
		log.Fatalln(err)
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	resp, err := client.Get(download)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)

	return fileName, true, fontsPath
}

func DownloadNewFont(font string) string {
	var fontName string

	//Getting the file name, verifying if it was downloaded and getting the font path
	file, downloaded, path := downloadFontZip(font)

	if downloaded {
		kittyc.UnzipFile(file, path)
		fmt.Println("Unzipped file. Check it out!")
		fontName = strings.ReplaceAll(file, ".zip", "")
		newFN := strings.ReplaceAll(fontName, "-", " ")

		return newFN 
	} else {
		fmt.Println("Problem. Check it out!")
	}

	return ""
}

func SetFontComparing(font string) {
	var lowerFonts []string

	entryFont := strings.ToLower(font)
	editedFonts := kittyc.ListAllFonts()	

	for _, v := range(editedFonts) {
		lower := strings.ToLower(v)
		lowerFonts = append(lowerFonts, lower)
	}	

	fonts, _ := kittyc.SearchingSimilarValues(lowerFonts, entryFont)

	if !fonts {
		fmt.Println("Does this font exist?")
	} else {
		SetNewFont(cases.Title(language.English, cases.NoLower).String(entryFont))
		fmt.Println("Implemented font. Check it out")
	}
}

func SetNewFont(font string) {
	fontAttribute := "font_family"

	fontValue := fontAttribute + " " + font
	
	if !kittyc.ModifyingAtLine(fontAttribute, fontValue) {
		kittyc.WritingAtLine("# Fonts", fontValue)
	} 
}
