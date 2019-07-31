package logic

import (
	"math"
)

// A circle with a center and a radius
type RadialFence struct{
	Center [2]float64	`json:"center"`
	Radius float64	`json:"radius"`
}

// Helper function for radial distance
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Determine the distance between 2 coordinates on the globe.
func radialDistance(c1, c2 [2]float64) float64 {
	lat1 := degreesToRadians(c1[0])
	lon1 := degreesToRadians(c1[1])
	lat2 := degreesToRadians(c2[0])
	lon2 := degreesToRadians(c2[1])

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := c * 6371 // Earths radius in KM

	return distance
}

// Determines if a coordinate lies within a RadialFence.
func InRadius(coordinate [2]float64, fence RadialFence, ) bool{
	if radialDistance(fence.Center, coordinate) <= fence.Radius {
		return true
	}
	return false
}

// Creates a line from 2 coordinates.
func makeLine(c1 , c2 [2]float64) [2][2]float64 {
	return [2][2]float64{c1, c2}
}

// Helper function to determine if a point intersects a line.
func intersectsLine(coordinate [2]float64, line [2][2]float64) float64 {
	gradient := gradient(line)

	// Determine the high and low points of the line
	yHighIndex := 0
	yHigh := line[0][1]
	if line[1][1] > yHigh {
		yHighIndex = 1
		yHigh = line[1][1]
	}

	// Now we have yHighIndex.
	var xHighIndex int
	if gradient > 0 {
		xHighIndex = yHighIndex
	} else if gradient < 0 {
		xHighIndex = 1 - yHighIndex
	} else {
		//GRADIENT IS ZERO => Consider as no intersections (even if ray covers line)
		return 0
	}

	if coordinate[1] > line[yHighIndex][1] || coordinate[1] < line[1-yHighIndex][1] {
		return 0
	}

	//BUG this returns 0.5 for points to the left as well.
	if coordinate[1] == line[0][1] || coordinate[1] == line[1][1] {
		return handleVerticeIntersection(coordinate, line, yHighIndex)
	}

	if coordinate[0] <= line[xHighIndex][0] && coordinate[0] <= line[1-xHighIndex][0] {
		return 1
	}

	deltaX := coordinate[0]-line[1-xHighIndex][0]
	yLine := line[1-xHighIndex][1] + gradient*deltaX
	return intersectsLineToRight(yLine, coordinate[1], gradient)
}

// Helper function to determine if intersection should count given a point where the ray from a line intersects the tip of a line.
// Intersection should only count when the line is to the right of the point AND the ray intersects the top of the line.
// E.g. (5, 5) intersects Line{(10, 5), (15, 3)} and Line{(10, 3), (15, 5)} but not Line{(10,10), (15, 5)}
func handleVerticeIntersection(coordinate [2]float64, line [2][2]float64, yHighIndex int) float64 {
	if coordinate[1] == line[yHighIndex][1] {
		if coordinate[0] > line[yHighIndex][0] {
			return 0
		} else {
			return 1
		}
	} else {
		return 0
	}
}

// Helper function to determine whether a point (xPoint, yPoint) intersects a line to its right.
// yLine refers to the y value of the line at xPoint.
func intersectsLineToRight (yLine float64, yPoint float64, gradient float64) float64 {
	if yLine == yPoint {
		return 1
	} else if yLine < yPoint {
		if gradient > 0 {
			return 1
		} else {
			return 0
		}
	} else {
		if gradient > 0 {
			return 0
		} else {
			return 1
		}
	}
}

// Helper function to generate a list of lines given a list of coordinates that make up a polygon
// Coordinates must be given in anticlockwise fashion.
func generateLines(coordinates [][2]float64) [][2][2]float64 {
	var lines [][2][2]float64

	for index, coordinate := range coordinates {
		newLine := [2][2]float64{coordinate, coordinates[(index+1)%len(coordinates)]}
		lines = append(lines, newLine)
	}
	return lines
}

// Helper function to determine the gradient of a line
func gradient(line [2][2]float64) float64 {
	return (line[1][1]-line[0][1])/(line[1][0]-line[0][0])
}

// Given a point and a list of coordinates that make up a polygon, determine the number of lines that the point (extended) would intersect
// Helper function for inOrOut
func countIntersects(point [2]float64, coordinates [][2]float64) float64 {
	result := 0.0
	lines := generateLines(coordinates)
	for _, line := range lines {
		result += intersectsLine(point, line)
	}
	return result
}

// Given a point and a list of coordinates that make up a polygon, determine if the point lies inside the polygon.
// Coordinates must be ordered in anticlockwise fashion.
func InPoly(point [2]float64, coordinates [][2]float64) bool {
	intersects := int(countIntersects(point, coordinates))
	if intersects%2 > 0 {
		return true
	} else {
		return false
	}

}
