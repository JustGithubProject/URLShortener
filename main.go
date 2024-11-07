package main

import (
	"encoding/json"
	"log"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
)

const (
	storageFilename = "storage.json"
	serverIP = "localhost:8080"
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


const AddForm = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Добавить URL</title>
		<style>
			body {
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				font-family: Arial, sans-serif;
				background-color: #f4f4f9;
			}
			.container {
				text-align: center;
				padding: 20px;
				background-color: white;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				border-radius: 8px;
				width: 100%;
				max-width: 400px;
			}
			h1 {
				color: #333;
			}
			label {
				display: block;
				margin-bottom: 10px;
				font-size: 16px;
				color: #555;
			}
			input[type="text"] {
				width: 100%;
				padding: 10px;
				margin-bottom: 20px;
				border: 1px solid #ccc;
				border-radius: 4px;
				font-size: 16px;
			}
			button {
				width: 100%;
				padding: 10px;
				background-color: #4CAF50;
				color: white;
				border: none;
				border-radius: 4px;
				font-size: 16px;
				cursor: pointer;
			}
			button:hover {
				background-color: #45a049;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Добавить URL</h1>
			<form method="POST" action="/add">
				<label for="url">URL:</label>
				<input type="text" id="url" name="url" required>
				<button type="submit">Отправить</button>
			</form>
		</div>
	</body>
	</html>
`

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	fullURL := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
	log.Println("fullURL=", fullURL)

	redirectTo := readURLFromJSON(fullURL)

	http.Redirect(w, r, redirectTo, http.StatusMovedPermanently)
}

func addURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	addURLtoJSON(url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func displayFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, AddForm)
}


func main() {
	http.HandleFunc("/", redirectURLHandler)
	http.HandleFunc("/add", addURLHandler)
	http.HandleFunc("/form", displayFormHandler)
	http.ListenAndServe(":8080", nil)

}

