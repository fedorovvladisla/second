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
	// Регистрируем обработчики с точными путями
	http.HandleFunc("/api/rv/", reverseHandler)
	http.HandleFunc("/", dateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func reverseHandler(w http.ResponseWriter, r *http.Request) {
	// Удаляем trailing slash если есть
	path := strings.TrimSuffix(r.URL.Path, "/")

	// Извлекаем ввод после /api/rv/
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[3] == "" {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}
	input := parts[3]

	// Проверяем валидность ввода
	if matched, _ := regexp.MatchString("^[a-z]+$", input); !matched {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Возвращаем перевернутую строку
	fmt.Fprint(w, reverseString(input))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func dateHandler(w http.ResponseWriter, r *http.Request) {
	// Удаляем trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	// Проверяем формат даты в URL
	now := time.Now()
	expectedPath := "/" + now.Format("020106")

	if path != expectedPath {
		http.NotFound(w, r)
		return
	}

	// Формируем JSON ответ
	response := map[string]string{
		"date":  now.Format("02-01-2006"),
		"login": "fedorovvlad",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
