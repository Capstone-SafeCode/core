package main

import (
	"flag"
	"log"
	"os"
)

func startSort(filepathList []string, blackList []string) {
	if len(filepathList) <= 0 {
		return
	} else {
		filepath := filepathList[0]
		files, err := os.ReadDir(filepath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fileType, data := isItFile(filepath + file.Name())
			if fileType == 1 {
				writeInTXT(filepath + file.Name() + " " + data + "\n")
			} else if fileType == 0 && !contains(blackList, filepath) {
				filepathList = append(filepathList, filepath+file.Name()+"/")
			}
		}
		blackList = append(blackList, filepath)
		index := findItsIndex(filepathList, filepath)
		filepathList = append(filepathList[:index], filepathList[index+1:]...)
		startSort(filepathList, blackList)
	}
}

func main() {
	var filepath string
	var filepathList []string

	flag.StringVar(&filepath, "path", "default", "Dossier où commence l'analyse")
	flag.Parse()

	if filepath != "default" && filepath[len(filepath)-1] == '/' {
		filepathList = append(filepathList, filepath)
		os.Truncate("to_analyse.txt", 0)
		startSort(filepathList, []string{})
	} else {
		help()
	}
}
