package main

import (
	"encoding/json"
	"log"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"io"
)

const (
	storageFilename = "storage.json"
	serverIP = "185.126.115.49"
)

func openFileStorage() *os.File{
	file, err := os.Open(storageFilename)
	if err != nil {
		log.Printf("Failed to open %s file\n", storageFilename)
		os.Exit(1)
	}
	return file
}

func generateUniqueURL(sourceURL string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, sourceURL)
	if err != nil {
		log.Println("Failed to generate")
		os.Exit(0)
	}
	return hex.EncodeToString(hasher.Sum(nil))

}

func addURLtoJSON(sourceURL string) {
	file := openFileStorage()
	defer file.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Failed to decode JSON")
		os.Exit(1)
	}
	uniqueURL := generateUniqueURL(sourceURL)
	shortURL := fmt.Sprintf("http://%s/%s", serverIP, uniqueURL)
	data[shortURL] = sourceURL

	updatedData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Failed to marshal data")
		os.Exit(1)
	}

	err = os.WriteFile(storageFilename, updatedData, 0644)
	if err != nil {
		log.Println("Failed to write updated data")
		os.Exit(1)
	}
}

func readURLFromJSON(shortURL string) string {
	file := openFileStorage()
	defer file.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Failed to decode JSON")
		os.Exit(1)
	}

	if url, ok := data[shortURL].(string); ok {
		return url
	}

	return ""
}

func main() {
	
}

