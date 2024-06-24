package main

import (
    "fmt"
    "os"
    "flag"
    "strings"
    "net/url"
    "net/http"
    "io"
)

func scanFile(fileName string) ([]string, error) {
    urls, err := os.ReadFile(fileName)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(urls), "\n")

    return lines, nil
}

func parseUrl(urls []string) ([]*url.URL) {
    result := []*url.URL{}

    for _, str := range urls {
        url, err := url.Parse(str)
        if err == nil && url.Scheme != "" && url.Host != "" {
            result = append(result, url)
        }
    }

    return result
}

func writeGetResult(urls []*url.URL, outputDirectoryName string) (error){
    for _, url := range urls {
        response, err := http.Get(url.String())
        if err != nil {
            continue
        }

        body, err := io.ReadAll(response.Body)
        if err != nil {
            continue
        }

        content := string(body)
        
        fileName := fmt.Sprintf("%s/%s.txt", outputDirectoryName, url.Host)
        
        err = os.WriteFile(fileName, []byte(content), 0777)
        if err != nil {
            fmt.Println(err)
            continue
        }
    }

    return nil
}

func main() {
    srcPtr := flag.String("src", "", "Название исходного файла")
    dstPtr := flag.String("dst", "", "Название конечной директории")

    flag.Parse()
    
    inputFileName := *srcPtr
    outputDirectoryName := *dstPtr

    _ = os.MkdirAll(outputDirectoryName, 0777)
    
    urls, err := scanFile(inputFileName)
    if err!= nil {
        fmt.Println(err)
        return
    }

    result := parseUrl(urls)

    err = writeGetResult(result, outputDirectoryName)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Ok")
}