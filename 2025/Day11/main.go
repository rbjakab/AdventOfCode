package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type Graph map[string][]string

func readGraph(filename string) Graph {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open input: %v", err)
	}
	defer file.Close()

	graph := make(Graph)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			log.Fatalf("invalid line: %q", line)
		}

		from := parts[0]
		to := strings.Fields(parts[1])
		graph[from] = to
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan input: %v", err)
	}

	return graph
}

func countPathsDFS(graph Graph, src, dst string, memo map[string]int) int {
	if src == dst {
		return 1
	}

	if v, ok := memo[src]; ok {
		return v
	}

	total := 0
	for _, next := range graph[src] {
		total += countPathsDFS(graph, next, dst, memo)
	}

	memo[src] = total
	return total
}

func hasPath(graph Graph, src, dst string, visited map[string]bool) bool {
	if src == dst {
		return true
	}

	if visited[src] {
		return false
	}
	visited[src] = true

	for _, next := range graph[src] {
		if hasPath(graph, next, dst, visited) {
			return true
		}
	}

	return false
}

func countPaths(graph Graph, src, dst string) int {
	return countPathsDFS(graph, src, dst, make(map[string]int))
}

func countPathsVia(graph Graph, src, dst string, viaA, viaB string) int {
	// Determine order
	if hasPath(graph, viaA, viaB, make(map[string]bool)) {
		return pathsThrough(graph, src, viaA, viaB, dst)
	}
	return pathsThrough(graph, src, viaB, viaA, dst)
}

func pathsThrough(graph Graph, src, mid1, mid2, dst string) int {
	p1 := countPathsDFS(graph, src, mid1, make(map[string]int))
	p2 := countPathsDFS(graph, mid1, mid2, make(map[string]int))
	p3 := countPathsDFS(graph, mid2, dst, make(map[string]int))
	return p1 * p2 * p3
}

func main() {
	start := time.Now()
	graph := readGraph("input.txt")

	log.Println("=== PART 1 ===")
	part1 := countPaths(graph, "you", "out")
	log.Printf("Paths from you to out: %d", part1)

	log.Println("=== PART 2 ===")
	part2 := countPathsVia(graph, "svr", "out", "dac", "fft")
	log.Printf("Paths from svr to out via dac & fft: %d", part2)

	log.Printf("Execution time: %s", time.Since(start))
}
