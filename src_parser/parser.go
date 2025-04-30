package src_parser

import (
	"fmt"
	"log"
	"os"
)

func startSort(finalFilepathList *[]string, filepathList []string, blackList []string) {
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
				*finalFilepathList = append(*finalFilepathList, filepath+file.Name()+" "+data+"\n")
			} else if fileType == 0 && !contains(blackList, filepath) {
				filepathList = append(filepathList, filepath+file.Name()+"/")
			}
		}
		blackList = append(blackList, filepath)
		index := findItsIndex(filepathList, filepath)
		filepathList = append(filepathList[:index], filepathList[index+1:]...)
		startSort(finalFilepathList, filepathList, blackList)
	}
}

func ParseFolder(filepath string) []string {
	var filepathList []string
	var finalFilepathList []string

	if filepath != "" && filepath[len(filepath)-1] == '/' {
		filepathList = append(filepathList, filepath)
		startSort(&finalFilepathList, filepathList, []string{})
	}
	fmt.Println(finalFilepathList)
	return finalFilepathList
}
