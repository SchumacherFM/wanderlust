package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var compressFileExt map[string]bool

func init() {
	compressFileExt = map[string]bool{
		"css": true,
		"js":  true,
		"eot": true,
		"svg": true,
	}
}

func operationEmbedGo(pkg *build.Package) {

	boxMap := findBoxes(pkg)

	// notify user when no calls to rice.FindBox are made (is this an error and therefore os.Exit(1) ?
	if len(boxMap) == 0 {
		fmt.Println("no calls to rice.FindBox() found")
		return
	}

	verbosef("\n")

	for boxname := range boxMap {
		// find path and filename for this box
		boxPath := filepath.Join(pkg.Dir, boxname)
		boxFilename := strings.Replace(boxname, "/", "-", -1)
		boxFilename = strings.Replace(boxFilename, "..", "back", -1)
		boxFilename = boxFilename + `.rice-box.go`

		// verbose info
		verbosef("embedding box '%s'\n", boxname)
		verbosef("\tto file %s\n", boxFilename)

		// create box datastructure (used by template)
		box := &boxDataType{
			Package: pkg.Name,
			BoxName: boxname,
			UnixNow: time.Now().Unix(),
			Files:   make([]*fileDataType, 0),
			Dirs:    make(map[string]*dirDataType),
		}

		// fill box datastructure with file data
		filepath.Walk(boxPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("error walking box: %s\n", err)
				os.Exit(1)
			}

			filename := strings.TrimPrefix(path, boxPath)
			filename = strings.Replace(filename, "\\", "/", -1)
			filename = strings.TrimPrefix(filename, "/")
			if info.IsDir() {
				dirData := &dirDataType{
					Identifier: "dir_" + nextIdentifier(),
					FileName:   filename,
					ModTime:    info.ModTime().Unix(),
					ChildFiles: make([]*fileDataType, 0),
					ChildDirs:  make([]*dirDataType, 0),
				}
				verbosef("\tincludes dir: '%s'\n", dirData.FileName)
				box.Dirs[dirData.FileName] = dirData

				// add tree entry (skip for root, it'll create a recursion)
				if dirData.FileName != "" {
					pathParts := strings.Split(dirData.FileName, "/")
					parentDir := box.Dirs[strings.Join(pathParts[:len(pathParts)-1], "/")]
					parentDir.ChildDirs = append(parentDir.ChildDirs, dirData)
				}
			} else {
				fileData := &fileDataType{
					Identifier: "file_" + nextIdentifier(),
					FileName:   filename,
					ModTime:    info.ModTime().Unix(),
				}
				verbosef("\tincludes file: '%s'\n", fileData.FileName)
				var isCompressed bool
				fileData.Content, isCompressed, err = getGzipFileContent(path)
				if err != nil {
					fmt.Printf("error reading file content while walking box: %s\n", err)
					os.Exit(1)
				}
				fileData.IsGzip = isCompressed

				// @todo add here the compressor for GZIP
				box.Files = append(box.Files, fileData)

				// add tree entry
				pathParts := strings.Split(fileData.FileName, "/")
				parentDir := box.Dirs[strings.Join(pathParts[:len(pathParts)-1], "/")]
				parentDir.ChildFiles = append(parentDir.ChildFiles, fileData)
			}
			return nil
		})

		embedSourceUnformated := bytes.NewBuffer(make([]byte, 0))

		// execute template to buffer
		err := tmplEmbeddedBox.Execute(embedSourceUnformated, box)
		if err != nil {
			log.Printf("error writing embedded box to file (template execute): %s\n", err)
			os.Exit(1)
		}

		// format the source code
		embedSource, err := format.Source(embedSourceUnformated.Bytes())
		if err != nil {
			log.Printf("error formatting embedSource: %s\n", err)
			os.Exit(1)
		}

		// create go file for box
		boxFile, err := os.Create(filepath.Join(pkg.Dir, boxFilename))
		if err != nil {
			log.Printf("error creating embedded box file: %s\n", err)
			os.Exit(1)
		}
		defer boxFile.Close()

		// write source to file
		_, err = io.Copy(boxFile, bytes.NewBuffer(embedSource))
		if err != nil {
			log.Printf("error writing embedSource to file: %s\n", err)
			os.Exit(1)
		}
	}
}

// returns gzip compressed data
func getGzipFileContent(filePath string) ([]byte, bool, error) {

	ext := path.Ext(filePath)
	_, isAllowed := compressFileExt[ext[1:]]

	content, err := ioutil.ReadFile(filePath)
	if nil != err {
		fmt.Printf("error reading file content while walking box: %s\n", err)
		os.Exit(1)
	}

	if false == isAllowed {
		return content, isAllowed, nil
	}

	var contentBuffer bytes.Buffer
	gzWriter, errGZ := gzip.NewWriterLevel(&contentBuffer, gzip.BestCompression)
	if nil != errGZ {
		fmt.Printf("error creating gz writer while walking box: %s\n", err)
		os.Exit(1)
	}
	defer gzWriter.Close()
	_, errGzW := gzWriter.Write(content)
	if nil != errGzW {
		fmt.Printf("error writing gz data while walking box: %s\n", err)
		os.Exit(1)
	}
	gzWriter.Flush()
	return contentBuffer.Bytes(), isAllowed, nil
}
