package main

import "fmt"

func BFS(graph map[int][]int, start int) {
    visited := make(map[int]bool)
    queue := []int{start}

    visited[start] = true

    for len(queue) > 0 {
        // Remove the first element
        node := queue[0]
        queue = queue[1:]

        fmt.Print(node, " ")

        // Visit all unvisited neighbours
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}

func main() {
    graph := map[int][]int{
        0: {1, 2},
        1: {0, 3, 4},
        2: {0, 5},
        3: {1},
        4: {1, 5},
        5: {2, 4},
    }

    fmt.Print("BFS traversal: ")
    BFS(graph, 0)
    fmt.Println()
}
