package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var mdFilePtr = flag.String("md", "xx.md", "input your md file with absolute path")

func main() {
	flag.Parse()

	mdFile := *mdFilePtr

	data, err := os.ReadFile(mdFile)
	if err != nil {
		panic(fmt.Errorf("read file error: %v", err))
	}
	fmt.Println("read file done")

	backupFile := fmt.Sprintf("%s.bak.%.2x", mdFile, md5.Sum(data))
	err = os.WriteFile(backupFile, data, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("write bak file error: %v", err))
	}
	fmt.Println("backup file done:", backupFile)
	fmt.Println("------")

	dir, fileName := path.Split(mdFile)
	fmt.Printf("dir: %s, file name: %s\n", dir, fileName)

	fileNameWithoutSuffix := strings.TrimRight(fileName, ".md")
	fmt.Println("file name without suffix:", fileNameWithoutSuffix)
	fmt.Println("------")

	assetsDir := filepath.Join(dir, fileNameWithoutSuffix)
	err = os.RemoveAll(assetsDir)
	if err != nil {
		panic(fmt.Errorf("remove image folder failed: %v", err))
	}
	fmt.Println("remove asset dir(if exist) done:", assetsDir)

	err = os.Mkdir(assetsDir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("create asset folder failed: %v", err))
	}
	fmt.Println("create asset dir done:", assetsDir)
	fmt.Println("------")

	re := regexp.MustCompile("(?m)\\]\\(http.*?\\.(png|jpg|bmp)\\)")
	urls := re.FindAllString(string(data), -1)
	fmt.Println("urls:", urls)

	for i, url := range urls {
		uu := url[2 : len(url)-1]
		fetchImage(assetsDir, uu, i)
	}
	fmt.Println("fetch images done")
	fmt.Println("------")

	counter := 0
	data = []byte(re.ReplaceAllStringFunc(string(data), func(m string) string {
		ret := "](" + fmt.Sprintf("%v", counter) + ".png)"
		counter++
		return ret
	}))

	err = ioutil.WriteFile(mdFile, data, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("write file failed: %v", err))
	}
	fmt.Println("update md file done:", mdFile)
}

func fetchImage(dir string, url string, index int) error {
	response, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	defer response.Body.Close()

	picName := filepath.Join(dir, fmt.Sprintf("%v.png", index))
	file, err := os.Create(fmt.Sprintf(picName))
	if err != nil {
		panic(fmt.Sprintf("create file failed: %v", err))
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(fmt.Sprintf("download file failed: %v", err))
	}
	return nil
}
