package txtparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"vrpSolution/internal/models"
)

// ParseTextFile reads a txt file and parses it into a slice of DrivingRoutes structs.
func ParseTextFile(filePath string) ([]models.DrivingRoute, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	var parsedLoads []models.DrivingRoute
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip the first line of headers

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		loadNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("error converting loadNumber: %s", line)
		}

		pickUpCoord, err := models.StringToCartesianCoordinate(parts[1])
		if err != nil {
			return nil, fmt.Errorf("error converting pickup: %s", line)
		}

		dropoffCoord, err := models.StringToCartesianCoordinate(parts[2])
		if err != nil {
			return nil, fmt.Errorf("error converting dropoff: %s", line)
		}

		dr := models.DrivingRoute{
			LoadNumber: loadNumber,
			PickUp:     pickUpCoord,
			DropOff:    dropoffCoord,
		}
		parsedLoads = append(parsedLoads, dr)
	}

	return parsedLoads, nil
}
