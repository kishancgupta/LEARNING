package main

import "fmt"

func BFS(graph map[int][]int, start int) {
	visited := make(map[int]bool)
	queue := []int{start}

	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		fmt.Print(node, " ")

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}
