package file

import "bufio"
import "log"
import "os"
import "strings"

// Read a text file, and return content line by line.
func ReadLines(path string) []string {
	res := make([]string, 0, 100)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return res
	}

	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res
}

// Read a text file to string.
func ReadFile(path string) string {
	rows := ReadLines(path)
	return strings.Join(rows, "")
}

// Overwrite a file.
func WriteFile(path string, content string) {
	var file *os.File
	var err  error

	if  file, err = os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

func AppendToFile(path string, content string) {
	var file *os.File
	var err  error

	if  file, err = os.OpenFile(path, os.O_APPEND | os.O_CREATE |os.O_WRONLY, 0600); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

