package doc

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lamzin/search-engine/index/common"
)

func NewDocManager(dataRoot string) (*DocManager, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}
	return &DocManager{
		docFolerRoot:  dataRoot,
		fileNameRegex: reg,
	}, nil
}

type DocManager struct {
	docFolerRoot string

	fileNameRegex *regexp.Regexp

	docInfoList []DocInfo
}

func (data *DocManager) getDocPath(doc Doc) (string, string) {
	fileName := data.fileNameRegex.ReplaceAllString(strings.ToLower(doc.Name), "")
	if len(fileName) > 64 {
		fileName = fileName[:64]
	}

	fileFolder := "_short"
	if len(fileName) > 2 {
		fileFolder = strings.ToLower(fileName[:2])
	}
	return fileFolder, filepath.Join(fileFolder, fileName+".txt")
}

func (data *DocManager) DumpDocument(doc Doc) error {
	fileFolder, filePath := data.getDocPath(doc)

	errFolder := os.Mkdir(filepath.Join(data.docFolerRoot, fileFolder), os.ModePerm)
	if errFolder != nil && !strings.Contains(errFolder.Error(), "file exists") {
		return errFolder
	}

	file, err := os.Create(filepath.Join(data.docFolerRoot, filePath))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(doc.MustCompress()); err != nil {
		return err
	}

	data.docInfoList = append(data.docInfoList, DocInfo{
		ID:   len(data.docInfoList),
		Name: doc.Name,
		Path: filePath,
	})
	return nil
}

func (data *DocManager) Close() error {
	file, err := os.Create(data.docInfoFilePath())
	if err != nil {
		return err
	}

	for _, doc := range data.docInfoList {
		if _, err := file.WriteString(doc.String() + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (data *DocManager) GetAllList() ([]DocInfo, error) {
	if len(data.docInfoList) > 0 {
		return data.docInfoList, nil
	}

	lines, err := common.ReadFile(data.docInfoFilePath())
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		doc, err := DocInfoFromString(line)
		if err != nil {
			return nil, err
		}
		data.docInfoList = append(data.docInfoList, *doc)
	}
	return data.docInfoList, nil
}

func (data *DocManager) docInfoFilePath() string {
	return filepath.Join(data.docFolerRoot, "docs.txt")
}

func (data *DocManager) GetByInfo(info DocInfo) (*Doc, error) {
	lines, err := common.ReadFile(filepath.Join(data.docFolerRoot, info.Path))
	if err != nil {
		return nil, err
	}
	return &Doc{
		DocInfo: info,
		Lines:   lines,
	}, nil
}
