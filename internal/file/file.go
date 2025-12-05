package file

import (
	"bufio"
	"log"
	"os"
)

func Load(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open file '%s': %s", filename, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Fatalf("Unable to close file '%s': %s", filename, err)
		}
	}(file)

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Unable to read file '%s': %s", filename, err)
	}

	return lines
}
