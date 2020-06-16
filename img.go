package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/abhinababasu/facethumbnail"
)

func main() {
	srcFolder := flag.String("src", "", "Source directory with images")
	dstFolder := flag.String("dst", "", "Destination directory where the album will generated")
	aspectRatio := flag.String("ratio", "9:16", "Target aspect ratio as w:h format. e.g. 9:16")
	face := flag.Bool("face", false, "Detect faces in images for generating thumbnails")
	verbose := flag.Bool("v", true, "Verbose logging mode")

	flag.Parse()

	if *srcFolder == "" || *dstFolder == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	if *face {
		log.Print("Will attempt to detect faces")
	}

	w, h, err := parseAspectRatio(*aspectRatio)
	if err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	GenerateImagesIntoDir(*srcFolder, *dstFolder, uint(w), uint(h), *face)
}

func parseAspectRatio(s string) (int, int, error) {
	re := regexp.MustCompile(":")
	ss := re.Split(s, -1)
	if len(ss) != 2 {
		return 0, 0, fmt.Errorf("Invalid aspect ratio")
	}

	w, err := strconv.Atoi(ss[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid aspect ratio")
	}

	h, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid aspect ratio")
	}

	return w, h, nil
}

func GenerateImagesIntoDir(srcFolder, dstFolder string, w, h uint, detectFace bool) (int, error) {
	log.Printf("Enumerating folder %v\n", srcFolder)

	if _, err := os.Stat(dstFolder); os.IsNotExist(err) {
		os.Mkdir(dstFolder, 0666)
	}

	files, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		return 0, fmt.Errorf("Failed to enumerate folder %v", srcFolder)
	}

	var fd facethumbnail.FaceDetector

	if detectFace {
		pwd, _ := os.Getwd()
		cascadeFile := path.Join(pwd, "facefinder")
		if _, err := os.Stat(cascadeFile); err != nil {
			return 0, fmt.Errorf("Cascade file not found for face detection")
		}

		fd = facethumbnail.GetFaceDetector(cascadeFile)
		fd.Init(-1, -1)
	}

	i := 0
	var wg sync.WaitGroup

	for _, file := range files {
		if file.IsDir() {
			log.Printf("Sub-dir not supported, skipping %v\n", file.Name())
		} else {
			wg.Add(1)
			go func(index int, filename string) {
				defer wg.Done()
				srcPath := filepath.Join(srcFolder, filename)
				dstPath := filepath.Join(dstFolder, filename)

				facethumbnail.ResizeToAspectRatio(fd, srcPath, dstPath, w, h)
				log.Printf(">>>>> Done %v", dstPath)
			}(i, file.Name())
			i++
		}
	}

	wg.Wait()

	return i, nil
}

// CreateDirIfNotExist creates all dirs in the path if does not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
