package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSortValues(t *testing.T) {
	input := []int{7, 3, 9, 3, -1, 5}
	want := []int{-1, 3, 3, 5, 7, 9}
	for _, algorithm := range []string{"bubble", "selection", "insertion", "merge", "quick"} {
		t.Run(algorithm, func(t *testing.T) {
			trace, err := SortValues(algorithm, input)
			if err != nil {
				t.Fatal(err)
			}
			if got := trace.Steps[len(trace.Steps)-1].Values; !reflect.DeepEqual(got, want) {
				t.Fatalf("final values = %v, want %v", got, want)
			}
			if !reflect.DeepEqual(input, []int{7, 3, 9, 3, -1, 5}) {
				t.Fatalf("SortValues changed its input: %v", input)
			}
		})
	}
}

func TestSearchValues(t *testing.T) {
	tests := []struct {
		name      string
		algorithm string
		values    []int
		target    int
		wantIndex int
	}{
		{"linear found", "linear", []int{8, 2, 5}, 2, 1},
		{"linear missing", "linear", []int{8, 2, 5}, 4, -1},
		{"binary found", "binary", []int{1, 2, 5, 8, 13}, 8, 3},
		{"binary missing", "binary", []int{1, 2, 5, 8, 13}, 4, -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trace, err := SearchValues(test.algorithm, test.values, test.target)
			if err != nil {
				t.Fatal(err)
			}
			if trace.FoundIndex != test.wantIndex {
				t.Fatalf("found index = %d, want %d", trace.FoundIndex, test.wantIndex)
			}
		})
	}
}

func TestBinarySearchRejectsUnsortedValues(t *testing.T) {
	if _, err := SearchValues("binary", []int{2, 1}, 1); err == nil {
		t.Fatal("expected unsorted input error")
	}
}

func TestSortHandler(t *testing.T) {
	body, err := json.Marshal(algorithmRequest{Algorithm: "merge", Values: []int{3, 1, 2}})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/sort", bytes.NewReader(body))
	response := httptest.NewRecorder()
	sortHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusOK, response.Body)
	}
}

func TestSortHandlerRejectsUnknownAlgorithm(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/sort", bytes.NewBufferString(`{"algorithm":"bogus","values":[3,1]}`))
	response := httptest.NewRecorder()
	sortHandler(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
