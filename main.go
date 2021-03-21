// Package main package contains the main application logic
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"sync"
)

var (
	filePath       *string
	wg             = sync.WaitGroup{}
	mutex          = sync.Mutex{}
	filesSlice     []FileStruct
	duplicateFiles []string
)

type FileStruct struct {
	FileName string
	FileSize int64
}

func init() {
	filePath = flag.String("filePath", "/Users/gregory", "the file path")
	flag.Parse()
}

func main() {
	wg.Add(1)
	go readDirectory(*filePath)
	wg.Wait()

	findEqualFiles()
	wg.Wait()

	if len(duplicateFiles) == 0 {
		fmt.Println("No repeating files")
	} else {
		fmt.Println("Find repeating files:")
		for _, fileRepeat := range duplicateFiles {
			fmt.Println(fileRepeat)
		}
	}
}

// readDirectory reads current and nested directories and puts all regular files in filesSlice
func readDirectory(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if f.Mode()&(1<<2) != 0 {
				wg.Add(1)
				go readDirectory(getDirectoryPath(path, f.Name()))
			}
		} else {
			mutex.Lock()
			filesSlice = append(filesSlice, FileStruct{FileName: f.Name(), FileSize: f.Size()})
			mutex.Unlock()
		}
	}
	wg.Done()
}

// getDirectoryPath returns full file path
func getDirectoryPath(path string, directory string) string {
	sb := strings.Builder{}
	sb.WriteString(path)
	sb.WriteString("/")
	sb.WriteString(directory)
	return sb.String()
}

// findEqualFiles finds equals files in filesSlice and puts the name of repeating file into duplicateFiles
func findEqualFiles() {
	for i := 0; i < len(filesSlice)-1; i++ {
		wg.Add(1)
		go func() {
			for j := len(filesSlice) - 1; j > i; j-- {
				if reflect.DeepEqual(filesSlice[i], filesSlice[j]) {
					putFile(filesSlice[i].FileName)
				}
			}
			wg.Done()
		}()
	}
}

// putFile puts file name in duplicateFiles if fileName is unique in this slice
func putFile(fileName string) {
	if notContains(fileName) {
		mutex.Lock()
		duplicateFiles = append(duplicateFiles, fileName)
		mutex.Unlock()
	}
}

// notContains returns false if duplicate slice contains input fileName and returns true if not
func notContains(fileName string) bool {
	for _, fn := range duplicateFiles {
		if fn == fileName {
			return false
		}
	}
	return true
}