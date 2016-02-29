package main

import "fmt"
import "io/ioutil"
import "os"
import "path"
import "strings"

var Num int = 0
var FileNum int = 0

func traverse(folder string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if (file.IsDir()) {
			traverse(path.Join(folder, file.Name()))
		} else {
			if strings.EqualFold(file.Name(), ".DS_Store") {
				Num++
				os.Remove(path.Join(folder, file.Name()))
				fmt.Printf("%s\t%d\n", folder, Num)
			}
			FileNum++
		}
	}
}

func main() {
	traverse("/Users/Lixing")
	fmt.Println(FileNum)
}


