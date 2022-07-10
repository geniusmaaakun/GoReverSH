package pkgclient

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

//outdir
/*
	if runtime.GOOS == "windows" {
		fpth = `C:\Windows\Temp\`
	} else {
		fpth = `/tmp/`
	}
*/

//TODO outputDirを指定
//スクリーンショットを撮影
//conn net.Conn
func Getscreenshot(outdir string) ([]string, error) {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(outdir); os.IsNotExist(err) {
			if err2 := os.MkdirAll(outdir, 0755); err2 != nil {
				log.Fatalf("Could not create the path %s", fpth)
			}
		}
		fileName := fmt.Sprintf("Scr-%d-%dx%d.png", i, bounds.Dx(), bounds.Dy())
		fullpath := filepath.Join(outdir, fileName)
		filenames = append(filenames, fullpath)
		file, err := os.Create(fullpath)
		if err != nil {
			return nil, err
		}

		defer file.Close()
		png.Encode(file, img)
	}
	return filenames, nil
}
