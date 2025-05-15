// mc-mirror.go
//
// Программа: mc-mirror (Maven Central Mirror)
// Назначение: Прокси и кэширующее зеркало Maven Central с поддержкой GET и HEAD-запросов.
// Использование: go run mc-mirror.go [-port=PORT]
// По умолчанию использует порт 8080 и сохраняет артефакты в ./storage/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	localCacheDir = "./storage" // Папка локального кэша
	upstreamURL   = "https://repo1.maven.org/maven2" // Апстрим — Maven Central
)

func main() {
	// Флаг командной строки: -port (по умолчанию 8080)
	port := flag.Int("port", 8080, "Порт для запуска сервера (по умолчанию 8080)")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	http.HandleFunc("/", proxyHandler)

	log.Printf("mc-mirror: сервер запущен на http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// proxyHandler обрабатывает HTTP-запросы: GET и HEAD
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path

	if strings.HasSuffix(requestPath, "/") {
		http.NotFound(w, r)
		return
	}

	localPath := filepath.Join(localCacheDir, requestPath)

	// === HEAD-запрос (только проверка наличия файла) ===
	if r.Method == http.MethodHead {
		if fileExists(localPath) {
			http.ServeFile(w, r, localPath)
			return
		}

		resp, err := http.Head(upstreamURL + requestPath)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.NotFound(w, r)
			return
		}

		copyHeaders(w.Header(), resp.Header)
		w.WriteHeader(http.StatusOK)
		return
	}

	// === GET-запрос (скачивание файла) ===
	if r.Method == http.MethodGet {
		if fileExists(localPath) {
			log.Printf("Отдаём из кэша: %s\n", localPath)
			http.ServeFile(w, r, localPath)
			return
		}

		resp, err := http.Get(upstreamURL + requestPath)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Ошибка загрузки из апстрима", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if err := saveToCache(localPath, resp.Body); err != nil {
			log.Printf("Ошибка сохранения в кэш: %v", err)
			http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
			return
		}

		log.Printf("Загружено и закешировано: %s\n", localPath)
		http.ServeFile(w, r, localPath)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Проверка наличия локального файла
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// Сохраняет артефакт в локальный кэш
func saveToCache(path string, data io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, data)
	return err
}

// Копирует заголовки из ответа апстрима в локальный ответ
func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

