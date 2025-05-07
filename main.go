package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/api/rv/", reverseHandler)

	http.HandleFunc("/", dateHandler)

	port := "5000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Сервер запущен на порту %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func reverseHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/rv/"
	path := r.URL.Path

	if !strings.HasPrefix(path, prefix) {
		http.NotFound(w, r)
		return
	}
	input := path[len(prefix):]
	if input == "" {
		http.NotFound(w, r)
		return
	}

	matched, err := regexp.MatchString("^[a-z]+$", input)
	if err != nil || !matched {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	reversed := reverseString(input)
	fmt.Fprint(w, reversed)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func dateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	expectedPath := "/" + now.Format("020106")
	if r.URL.Path != expectedPath {
		http.NotFound(w, r)
		return
	}

	response := map[string]string{
		"date":  now.Format("02-01-2006"),
		"login": "fedorovvlad",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
