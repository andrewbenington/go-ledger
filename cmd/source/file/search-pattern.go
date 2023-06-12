package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

type SearchPattern struct {
	Directories      []string        `yaml:"directories"`
	FileNameKeywords []string        `yaml:"filename_keywords"`
	FileNamePatterns []regexp.Regexp `yaml:"filename_patterns"`
}

var (
	posixHomeGlobPattern = regexp.MustCompile("^~")
	runtimeOS            = runtime.GOOS
)

func (f *SearchPattern) FindMatchingFiles() ([]string, error) {
	filePaths := []string{}
	for _, directory := range f.Directories {
		dirFiles, err := f.findFilesInDirectory(directory)
		if err != nil {
			return []string{}, err
		}
		filePaths = append(filePaths, dirFiles...)
	}
	if len(filePaths) == 0 {
		fmt.Println("No matching files found")
	}
	return filePaths, nil
}

func (f *SearchPattern) findFilesInDirectory(directory string) ([]string, error) {
	matchingFiles := []string{}
	globbedDir := globDirectory(directory)
	dirFiles, err := os.ReadDir(globbedDir)
	if err != nil {
		return []string{}, fmt.Errorf("error reading directory %s: %w", globbedDir, err)
	}
	keywordPattern, err := f.keywordRegExp()
	if err != nil {
		return []string{}, fmt.Errorf("keyword error: %w", err)
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			continue
		}
		if len(f.FileNameKeywords) != 0 {
			if keywordPattern.MatchString(file.Name()) {
				fmt.Printf("Found file %s matching keyword in directory %s\n", file.Name(), globbedDir)
				matchingFiles = append(matchingFiles, filepath.Join(globbedDir, file.Name()))
				continue
			}
		}
		for _, pattern := range f.FileNamePatterns {
			if pattern.MatchString(file.Name()) {
				fmt.Printf("Found file %s matching pattern %s in directory %s\n", file.Name(), pattern.String(), globbedDir)
				matchingFiles = append(matchingFiles, filepath.Join(globbedDir, file.Name()))
				break
			}
		}
	}
	return matchingFiles, nil
}

func (f *SearchPattern) keywordRegExp() (*regexp.Regexp, error) {
	if len(f.FileNameKeywords) == 0 {
		return nil, nil
	}
	reString := ""
	for i, keyword := range f.FileNameKeywords {
		if i > 0 {
			reString = fmt.Sprintf("%s|", reString)
		}
		reString = fmt.Sprintf("%s%s", reString, keyword)
	}
	return regexp.Compile(reString)
}

func globDirectory(directory string) string {
	if runtimeOS == "linux" || runtimeOS == "darwin" {
		homePath := os.Getenv("HOME")
		return posixHomeGlobPattern.ReplaceAllString(directory, homePath)
	}
	return directory
}
