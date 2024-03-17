package main

import (
	"bytes"
	"fmt"
	"bufio"
	"regexp"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"io"
	"encoding/json"
	"strings"
	"path/filepath"
	"time"
)

func upload(imagePath string) string{

	url := "https://sm.ms/api/v2/upload"
	stringParam := "json"

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	imagePart, err := writer.CreateFormFile("smfile", imagePath)
	if (err != nil) {
		fmt.Println("Error creating form file:", err)
		return ""
	}


	_, err = io.Copy(imagePart, file)
	if (err != nil) {
		fmt.Println("Error writing file data to form field: ", err)
		return ""
	}

	err = writer.WriteField("format", stringParam)
	if err != nil {
		fmt.Println("Error writing form field:", err)
		return ""
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing form writer:", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "xxx")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	fmt.Println("Response:", string(responseBody))	

	responseJson := make(map[string]interface{})
	err = json.Unmarshal([]byte(responseBody), &responseJson)
	if err != nil{
		fmt.Println("Error Json", err)
		return ""
	}

	data, ok := responseJson["data"].(map[string]interface{})
	if !ok {
		fmt.Println("data does not exit")
		return ""
	}

	urlData, ok := data["url"].(string)
	if !ok { 
		fmt.Println("url does not exit")
		return ""
	}
	fmt.Println(urlData)
	return urlData
}

func download(imageUrl string, directoryPath string) string{

	proxyURL, err := url.Parse("http://xxx:xxx")
	if err != nil{
		fmt.Println("Error Set Proxy")
		return ""
	}

	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	//发出Http请求
	resp, err := http.Get(imageUrl)
	if err != nil{
		fmt.Println("Error Http", err)
		return ""
	}
	defer resp.Body.Close()

	fileName := filepath.Base(imageUrl)
	file, err := os.Create(directoryPath + "\\" + fileName)
	if err != nil{
		fmt.Println("Error Create File", err)
		return ""
	}
	defer file.Close()

	_ , err = io.Copy(file, resp.Body)
	if err != nil{
		fmt.Println("Error Copy Body ", err)
		return ""
	}
	return fileName
}


func main(){

	//扫描文件
	argCount := len(os.Args) 
	if argCount < 3{
		fmt.Println("Usage: main.exe path [upload|download]")
		return
	}
	
	targetFile := os.Args[1]
	file , err := os.Open(targetFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	op := "upload"
	if (os.Args[2] == "download"){
		op = "download"
	}

	//获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil{
		fmt.Println("Error Get Current Dir")
		return
	}
	fmt.Println("Current Dir:", currentDir)

	index := strings.LastIndex(targetFile, "\\")
	if index != -1{
		newDir := targetFile[:index]
		err = os.Chdir(newDir)
		if err != nil{
			fmt.Println("Error Change Dir", err)
			return
		}
	}

	currentDir , err = os.Getwd()
	if err != nil{
		fmt.Println("Error Get Current Dir", err)
		return
	}
	fmt.Println("Current Dir:", currentDir)


	scanner := bufio.NewScanner(file)
	pattern := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	var newLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line){
			fmt.Println("Match found:", line)
			re := regexp.MustCompile(`\(([^\(\)]*)\)`)
			match := re.FindStringSubmatch(line)
			path  := match[1]
			if op == "upload"{
				time.Sleep(3 * time.Second)
				url := upload(path)
				fmt.Println(path)
				if url != ""{
					line = strings.Replace(line, match[1], url, -1)
				}
			} else if op == "download"{
				directoryPath := "image"
				err = os.MkdirAll(directoryPath, 0755)
				if err != nil{
					fmt.Println("Error Mkdir", err)
					return 
				}
				imageLocalPath := download(path, directoryPath)
				fmt.Println(imageLocalPath)
				if imageLocalPath != ""{
					line = strings.Replace(line, match[1], directoryPath + "\\" + imageLocalPath, -1)
				}
			}
			
		}
		newLines = append(newLines, line)
	}
	fmt.Println(newLines)

	newFile, err := os.Create(targetFile + ".bak")
	if err != nil{
		fmt.Println("Create New File Error", err)
		return
	}
	defer newFile.Close()

	for _ , newLine := range newLines {
		_, err := fmt.Fprintln(newFile, newLine)
		if err != nil{
			fmt.Println("Error Write", err)
			return
		}
	}
}