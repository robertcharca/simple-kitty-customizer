package kittyc

import (
	"bufio"	
	"fmt"
	"io"
	"os"
	"strings"
)

func WritingAtLine(section string, addedLines []string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)			
	if err != nil {
		fmt.Println(err.Error())
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)		
	
	var lines []string	
	var newLines []string

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)	
		
		if strings.Contains(line, section) {	
			for _, val := range addedLines {
				newLines = append(newLines, val)
			} 

			lines = append(lines, newLines...)
		}	
	}

	scanningFile := scanner.Err()
	if scanningFile != nil {
		fmt.Println(scanningFile.Error())
	}

	/*
		Truncate: changes the size of the file. In this case, it change the last size to 0, 
		which means it's completely empty.
	*/
	truncate := file.Truncate(0) 
	if truncate != nil {
		fmt.Println(err.Error())
	}

	/*
		Seek: sets a new offset according to a "principal" offset. In this case, the new
		offset is 0 based on the "start" offset, which also is 0. This means the offset will
		be positioned at the complete beginning of the file.
	*/
	_, offsErr := file.Seek(0, io.SeekStart)
	if offsErr != nil {
		fmt.Println(err.Error())
	}
	
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	
	flush := writer.Flush()
	if flush != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Everything's great")
}

func ModifyingAtLine (oldLine, newLine string) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)			
	if err != nil {
		fmt.Println(err.Error())
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	foundLine := false	

	for scanner.Scan() {
		line := scanner.Text()

		possibleLines := strings.Fields(line)
		for _, val := range possibleLines {
			if val == oldLine {
				line = newLine
				foundLine = true
			}
		}	

		lines = append(lines, line)
	}

	scanningFile := scanner.Err()
	if scanningFile != nil {
		fmt.Println(scanningFile.Error())
		return false
	}

	if !foundLine {
		fmt.Println("Value not found in the file")
		return false
	}
	
	truncate := file.Truncate(0) 
	if truncate != nil {
		fmt.Println(err.Error())
		return false
	}
	
	_, offsErr := file.Seek(0, io.SeekStart)
	if offsErr != nil {
		fmt.Println(err.Error())
		return false
	}
	
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	
	flush := writer.Flush()
	if flush != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("Value updated")
	return true
}

func ModifyMultipleLines(oldLine , newLines []string) bool {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	foundLine := false

	counter := 0

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := line

		if len(oldLine) == 1 {
			if strings.Contains(line, oldLine[0]) {
				if counter < len(newLines) {
					modifiedLine = newLines[counter]
					counter++
				} else {
					continue
				}
				
				foundLine = true
			}
		}

		if len(oldLine) > 1 {
			if counter < len(newLines) {
				if strings.Contains(line, oldLine[counter]) {	
					modifiedLine = newLines[counter]
					foundLine = true
					counter++
				}
			}	
		}

		lines = append(lines, modifiedLine)
	}
	fmt.Println("counter ==> ", counter)

	if !foundLine {
		fmt.Println("Value not found in the file")
		return false
	}

	if err := file.Truncate(0); err != nil {
		fmt.Println(err.Error())
		return false
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		fmt.Println(err.Error())
		return false
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("Values updated")
	return true
}
