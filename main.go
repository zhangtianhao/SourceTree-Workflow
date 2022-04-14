package main

import (
	"fmt"
	"howett.net/plist"
	"io/ioutil"
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

func splitMatch(title string) string {
	fragments := strings.FieldsFunc(title, splitString)
	if len(fragments) > 1 {
		return title + " " + strings.Join(fragments, " ")
	}
	return title
}

func getList() *AlfredList {
	bs, _ := ioutil.ReadFile("/Users/user/Library/Application Support/SourceTree/browser.plist")
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
