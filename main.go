package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// scanFile считывает строки из файла и записывает в слайс
func scanFile(fileName string) ([]string, error) {
	lines, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	linesArr := strings.Split(string(lines), "\n")

	return linesArr, nil
}

// parseUrl отфильтровывает валидные URL от невалидных и записывает в слайс
func parseUrl(urls []string) []*url.URL {

	urlArr := []*url.URL{}

	for _, str := range urls {
		url, err := url.Parse(str)
		if err == nil && url.Scheme != "" && url.Host != "" {
			urlArr = append(urlArr, url)
		}
	}

	return urlArr
}

// getRequest выполняет get запрос по URL и возвращает ответ в виде string
func getRequest(url *url.URL) (string, error) {
	response, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	content := string(body)

	return content, nil
}

// upload получает данные из URL запросов и записывает в файлы
func upload(urls []*url.URL, outputDirectoryName string) []string {
	arrFileName := []string{}

	for _, url := range urls {
		content, err := getRequest(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fileName := fmt.Sprintf("%s/%s.txt", outputDirectoryName, url.Host)
		err = os.WriteFile(fileName, []byte(content), 0777)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		arrFileName = append(arrFileName, fileName)
	}

	return arrFileName
}

// printFileName выводит на экран названия файлов из слайса
func printFileName(fNames []string) {
	fmt.Println("\nПути к созданным файлам:")
	for _, fname := range fNames {
		fmt.Println(fname)
	}
}

func main() {
	srcPtr := flag.String("src", "", "Название исходного файла")
	dstPtr := flag.String("dst", "", "Название конечной директории")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: go run main.go --src=<путь_к_файлу> --dst=<путь_к_конечной_папке>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *srcPtr == "" || *dstPtr == "" {
		fmt.Println("Error: missing required flags.")
		flag.Usage()
		return
	}

	start := time.Now()
	fmt.Printf("Запуск программы: %v\n\n", start.Format("02.01.06 15:04:05"))

	inputFileName := *srcPtr
	outputDirectoryName := *dstPtr

	err := os.MkdirAll(outputDirectoryName, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	urls, err := scanFile(inputFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	urlsArray := parseUrl(urls)

	arrFileName := upload(urlsArray, outputDirectoryName)

	printFileName(arrFileName)

	fmt.Printf("\nКонец работы программы: %v\n", time.Now().Format("02.01.06 15:04:05"))
	fmt.Printf("Общее время работы: %v\n", time.Since(start))
}
