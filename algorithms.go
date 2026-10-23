package main

import (
	"fmt"
	"sort"
)

type Step struct {
	Values  []int  `json:"values,omitempty"`
	Indices []int  `json:"indices,omitempty"`
	Action  string `json:"action"`
	Message string `json:"message"`
	Low     int    `json:"low"`
	High    int    `json:"high"`
	Mid     int    `json:"mid"`
}

type Trace struct {
	Steps       []Step `json:"steps"`
	Comparisons int    `json:"comparisons"`
	Swaps       int    `json:"swaps"`
	FoundIndex  int    `json:"foundIndex"`
}

func SortValues(algorithm string, input []int) (Trace, error) {
	values := append([]int(nil), input...)
	trace := Trace{FoundIndex: -1}
	record := func(action, message string, indices ...int) {
		trace.Steps = append(trace.Steps, Step{
			Values:  append([]int(nil), values...),
			Indices: append([]int(nil), indices...),
			Action:  action,
			Message: message,
		})
	}
	compare := func(i, j int) {
		trace.Comparisons++
		record("compare", fmt.Sprintf("Compare %d and %d.", values[i], values[j]), i, j)
	}
	swap := func(i, j int) {
		values[i], values[j] = values[j], values[i]
		trace.Swaps++
		record("swap", fmt.Sprintf("Swap positions %d and %d.", i, j), i, j)
	}
	record("start", "Start with the input values.")

	switch algorithm {
	case "bubble":
		for end := len(values) - 1; end > 0; end-- {
			swapped := false
			for i := 0; i < end; i++ {
				compare(i, i+1)
				if values[i] > values[i+1] {
					swap(i, i+1)
					swapped = true
				}
			}
			if !swapped {
				break
			}
		}
	case "selection":
		for start := 0; start < len(values)-1; start++ {
			minimum := start
			for i := start + 1; i < len(values); i++ {
				compare(minimum, i)
				if values[i] < values[minimum] {
					minimum = i
				}
			}
			if minimum != start {
				swap(start, minimum)
			}
		}
	case "insertion":
		for i := 1; i < len(values); i++ {
			for j := i; j > 0; j-- {
				compare(j-1, j)
				if values[j-1] <= values[j] {
					break
				}
				swap(j-1, j)
			}
		}
	case "merge":
		var mergeSort func(int, int)
		mergeSort = func(start, end int) {
			if end-start < 2 {
				return
			}
			middle := start + (end-start)/2
			mergeSort(start, middle)
			mergeSort(middle, end)
			left := append([]int(nil), values[start:middle]...)
			right := append([]int(nil), values[middle:end]...)
			i, j, k := 0, 0, start
			for i < len(left) && j < len(right) {
				trace.Comparisons++
				record("compare", fmt.Sprintf("Compare %d and %d before merging.", left[i], right[j]), start+i, middle+j)
				if left[i] <= right[j] {
					values[k] = left[i]
					i++
				} else {
					values[k] = right[j]
					j++
				}
				record("write", fmt.Sprintf("Write %d at position %d.", values[k], k), k)
				k++
			}
			for i < len(left) {
				values[k] = left[i]
				record("write", fmt.Sprintf("Write %d at position %d.", values[k], k), k)
				i++
				k++
			}
			for j < len(right) {
				values[k] = right[j]
				record("write", fmt.Sprintf("Write %d at position %d.", values[k], k), k)
				j++
				k++
			}
		}
		mergeSort(0, len(values))
	case "quick":
		var quickSort func(int, int)
		quickSort = func(low, high int) {
			if low >= high {
				return
			}
			pivot := values[high]
			record("pivot", fmt.Sprintf("Choose %d as the pivot.", pivot), high)
			boundary := low
			for i := low; i < high; i++ {
				trace.Comparisons++
				record("compare", fmt.Sprintf("Compare %d with pivot %d.", values[i], pivot), i, high)
				if values[i] < pivot {
					if boundary != i {
						swap(boundary, i)
					}
					boundary++
				}
			}
			if boundary != high {
				swap(boundary, high)
			}
			quickSort(low, boundary-1)
			quickSort(boundary+1, high)
		}
		quickSort(0, len(values)-1)
	default:
		return Trace{}, fmt.Errorf("unknown sorting algorithm %q", algorithm)
	}
	record("done", "The values are sorted.")
	return trace, nil
}

func SearchValues(algorithm string, values []int, target int) (Trace, error) {
	trace := Trace{FoundIndex: -1}
	record := func(action, message string, index, low, high, mid int) {
		trace.Steps = append(trace.Steps, Step{
			Values:  append([]int(nil), values...),
			Indices: []int{index},
			Action:  action,
			Message: message,
			Low:     low,
			High:    high,
			Mid:     mid,
		})
	}

	switch algorithm {
	case "linear":
		record("start", fmt.Sprintf("Scan from left to right for %d.", target), -1, 0, len(values)-1, -1)
		for i, value := range values {
			trace.Comparisons++
			record("compare", fmt.Sprintf("Compare %d with target %d.", value, target), i, 0, len(values)-1, i)
			if value == target {
				trace.FoundIndex = i
				record("found", fmt.Sprintf("Found %d at index %d.", target, i), i, 0, len(values)-1, i)
				return trace, nil
			}
		}
		record("not-found", fmt.Sprintf("%d is not in the list.", target), -1, 0, len(values)-1, -1)
	case "binary":
		if !sort.IntsAreSorted(values) {
			return Trace{}, fmt.Errorf("binary search requires values in ascending order")
		}
		low, high := 0, len(values)-1
		record("start", fmt.Sprintf("Search the sorted list for %d.", target), -1, low, high, -1)
		for low <= high {
			mid := low + (high-low)/2
			trace.Comparisons++
			record("compare", fmt.Sprintf("Compare middle value %d with target %d.", values[mid], target), mid, low, high, mid)
			if values[mid] == target {
				trace.FoundIndex = mid
				record("found", fmt.Sprintf("Found %d at index %d.", target, mid), mid, low, high, mid)
				return trace, nil
			}
			if values[mid] < target {
				low = mid + 1
			} else {
				high = mid - 1
			}
			record("range", fmt.Sprintf("Keep searching indices %d through %d.", low, high), -1, low, high, -1)
		}
		record("not-found", fmt.Sprintf("%d is not in the list.", target), -1, low, high, -1)
	default:
		return Trace{}, fmt.Errorf("unknown searching algorithm %q", algorithm)
	}
	return trace, nil
}
