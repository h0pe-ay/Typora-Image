package main

import (
	"bytes"
	"fmt"
	"bufio"
	"regexp"
	"mime/multipart"
	"net/http"
	"os"
	"io"
	"encoding/json"
	"strings"
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
	req.Header.Set("Authorization", "*******")

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


func main(){

	//扫描文件
	targetFile := "babku分析.md"
	file , err := os.Open(targetFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

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
			url := upload(path)
			fmt.Println(path)
			if url != ""{
				line = strings.Replace(line, match[1], url, -1)
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