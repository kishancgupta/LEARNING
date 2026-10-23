package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

const maxValues = 100

//go:embed web
var webFiles embed.FS

type algorithmRequest struct {
	Algorithm string `json:"algorithm"`
	Values    []int  `json:"values"`
	Target    int    `json:"target"`
}

func main() {
	web, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/algorithms", algorithmsHandler)
	mux.HandleFunc("/api/sort", sortHandler)
	mux.HandleFunc("/api/search", searchHandler)
	mux.Handle("/", http.FileServer(http.FS(web)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Sorting & Searching Lab listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func algorithmsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string][]map[string]string{
		"sorting": {
			{"id": "bubble", "name": "Bubble sort", "complexity": "O(n²)"},
			{"id": "selection", "name": "Selection sort", "complexity": "O(n²)"},
			{"id": "insertion", "name": "Insertion sort", "complexity": "O(n²)"},
			{"id": "merge", "name": "Merge sort", "complexity": "O(n log n)"},
			{"id": "quick", "name": "Quick sort", "complexity": "O(n log n) average"},
		},
		"searching": {
			{"id": "linear", "name": "Linear search", "complexity": "O(n)"},
			{"id": "binary", "name": "Binary search", "complexity": "O(log n)"},
		},
	})
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	request, err := decodeRequest(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(request.Values) == 0 || len(request.Values) > maxValues {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("provide between 1 and %d values", maxValues))
		return
	}
	trace, err := SortValues(request.Algorithm, request.Values)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, trace)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	request, err := decodeRequest(w, r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(request.Values) == 0 || len(request.Values) > maxValues {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("provide between 1 and %d values", maxValues))
		return
	}
	trace, err := SearchValues(request.Algorithm, request.Values, request.Target)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, trace)
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (algorithmRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request algorithmRequest
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return request, errors.New("request must contain a single JSON object")
	}
	return request, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
