package main

import (
	"fmt"
	"github.com/zhangtianhao/SourceTree-Workflow/homedir"
	"howett.net/plist"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"unicode"
)

// SourceTreePlist {'$version': 100000, '$archiver': 'NSKeyedArchiver', '$top': {'root': Uid(1)}}
type SourceTreePlist struct {
	Version  uint                   `plist:"$version"`
	Archiver string                 `plist:"$archiver"`
	Top      map[string]interface{} `plist:"$top"`
	Objects  []interface{}          `plist:"$objects"`
}

func splitString(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

// k8s-docker-desktop-for-mac
// k8s docker desktop for mac
// k8s-docker k8s-docker-desktop k8s-docker-desktop-for k8s-docker-desktop-for-mac
// docker-desktop docker-desktop-for docker-desktop-for-mac
// desktop-for desktop-for-mac
// for-mac
func splitMatch(title string) string {
	fragments := strings.FieldsFunc(title, splitString)
	length := len(fragments)
	for i := 0; i < length; i++ {
		temp := fragments[i]
		for j := i + 1; j < length; j++ {
			temp = temp + "-" + fragments[j]
			fragments = append(fragments, temp)
		}
	}
	return strings.Join(fragments, " ")
}

func getList() *AlfredList {
	dir, _ := homedir.Dir()
	bs, _ := ioutil.ReadFile(path.Join(dir, "/Library/Application Support/SourceTree/browser.plist"))
	var stp SourceTreePlist
	_, _ = plist.Unmarshal(bs, &stp)
	alfredSlice := make([]AlfredItem, 0, len(stp.Objects))
	var tempName string
	for _, obj := range stp.Objects {
		typeName := reflect.TypeOf(obj).Name()
		if typeName == "string" {
			item := reflect.ValueOf(obj).String()
			if item[:1] == "/" {
				match := splitMatch(tempName)
				alfredSlice = append(alfredSlice, AlfredItem{
					Title:    tempName,
					Subtitle: item,
					Arg:      item,
					Match:    match,
				})
			} else {
				tempName = item
			}
		}
	}
	return &AlfredList{Items: alfredSlice}
}

func main() {
	alfredList := getList()
	fmt.Print(alfredList.ToJson())
}
