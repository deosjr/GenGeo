package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"os"
	
	m "github.com/deosjr/GRayT/src/model"
	"github.com/deosjr/GenGeo/gen"
)

// reads in a list of bicubic bezier patches
// NOTE: patch/vertex count is 1-based
// NOTE: line values are comma-separated
func LoadPatches(filename string) ([]gen.ParametricSurface, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return loadPatches(scanner)
}

func loadPatches(scanner *bufio.Scanner) ([]gen.ParametricSurface, error) {
	numPatches, err := readNum(scanner)
	if err != nil {
		return nil, err
	}
	rawPatches := [][]int64{}
	var i int64
	for i=0;i<numPatches; i++ {
		raw, err := readRawPatch(scanner)
		if err != nil {
			return nil, err
		}
		rawPatches = append(rawPatches, raw)
	}
	numVertices, err := readNum(scanner)
	if err != nil {
		return nil, err
	}
	vertexMap := map[int64]m.Vector{}
	for i=0; i<numVertices; i++ {
		scanner.Scan()
		line := scanner.Text()
		line = strings.TrimSpace(line)
		v, err := readVertex(strings.Split(line, ","))
		if err != nil {
			return nil, err
		}
		vertexMap[i+1] = v
	}

	surfaces := make([]gen.ParametricSurface, len(rawPatches))
	for i, r := range rawPatches {
		controlPoints := make([]m.Vector, len(r))
		for j, n := range r {
			v, ok := vertexMap[n]
			if !ok {
				return nil, fmt.Errorf("Vertex not in map: %d", n)
			}
			controlPoints[j] = v
		}
		surfaces[i] = gen.NewBicubicBezierPatch(controlPoints)
	}
	return surfaces, nil
}

func readNum(scanner *bufio.Scanner) (int64, error) {
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	return strconv.ParseInt(line, 10, 64)
}

func readRawPatch(scanner *bufio.Scanner) ([]int64, error) {
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	nums := strings.Split(line, ",")
	if len(nums) != 16 {
		return nil, fmt.Errorf("Invalid patch: %v", nums)
	}
	list := make([]int64, len(nums))
	for i, n := range nums {
		num, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return nil, err
		}
		list[i] = num
	}
	return list, nil
}

func readVertex(coordinates []string) (m.Vector, error) {
	if len(coordinates) != 3 {
		return m.Vector{}, fmt.Errorf("Invalid coordinates: %v", coordinates)
	}
	p1, err := strconv.ParseFloat(coordinates[0], 32)
	if err != nil {
		return m.Vector{}, err
	}
	p2, err := strconv.ParseFloat(coordinates[1], 32)
	if err != nil {
		return m.Vector{}, err
	}
	p3, err := strconv.ParseFloat(coordinates[2], 32)
	if err != nil {
		return m.Vector{}, err
	}
	return m.Vector{float32(p1), float32(p2), float32(p3)}, nil
}