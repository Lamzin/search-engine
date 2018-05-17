package indexwriter

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	startBuffer = 4
	maxFiles    = 26
	queueSize   = 1024 * 10
)

type LexemeInfo struct {
	Positions  []int
	LastLength int
}

type lexemeWriteTask struct {
	DocID    int
	Position int
}

type IndexDBWriter struct {
	indexPath string

	lexemeInfo map[string]*LexemeInfo

	fileSizes  []int
	files      []*os.File
	queues     []chan lexemeWriteTask
	workerStop chan bool
}

func NewIndexDBWriter(indexPath string) *IndexDBWriter {
	w := &IndexDBWriter{
		indexPath:  indexPath,
		lexemeInfo: make(map[string]*LexemeInfo, 0),
	}
	w.start()
	return w
}

func (w *IndexDBWriter) start() {
	w.fileSizes = make([]int, maxFiles)
	w.files = make([]*os.File, maxFiles)
	w.workerStop = make(chan bool, maxFiles)
	for i := 0; i < maxFiles; i++ {
		w.queues = append(w.queues, make(chan lexemeWriteTask, queueSize))
		go w.worker(i)
	}
}

func (w *IndexDBWriter) worker(fileInt int) {
	filePath := filepath.Join(w.indexPath, strconv.Itoa(fileInt))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	w.files[fileInt] = file

	for task := range w.queues[fileInt] {
		// fmt.Println(fileInt, len(w.queues[fileInt]))

		bytes, err := bigEndian.Compress([]int{task.DocID})
		if err != nil {
			panic(err)
		}
		_, err = file.WriteAt(bytes, int64(task.Position))
		if err != nil {
			panic(err)
		}
	}

	w.workerStop <- true
}

func (w *IndexDBWriter) findFileAndPosition(lexeme string) (file int, info *LexemeInfo) {
	info, ok := w.lexemeInfo[lexeme]
	if !ok {
		info = &LexemeInfo{
			Positions:  []int{w.fileSizes[0]},
			LastLength: 0,
		}
		w.lexemeInfo[lexeme] = info
	}

	file = len(info.Positions) - 1
	bufferLength := startBuffer << (uint)(file)
	if info.LastLength == bufferLength {
		file++
		info.LastLength = 0
		for file >= len(w.fileSizes) {
			w.fileSizes = append(w.fileSizes, 0)
		}
		info.Positions = append(info.Positions, w.fileSizes[file])
	} else if info.LastLength > bufferLength {
		panic("length more that allowed buffer size")
	}
	return file, info
}

func (w *IndexDBWriter) writeDocID(docID int, fileInt int, position int) error {
	if fileInt == len(w.files) || w.files[fileInt] == nil {
		filePath := filepath.Join(w.indexPath, strconv.Itoa(fileInt))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		w.files = append(w.files, file)
	}
	file := w.files[fileInt]

	bytes, err := bigEndian.Compress([]int{docID})
	if err != nil {
		return err
	}

	_, err = file.WriteAt(bytes, int64(position))
	return err
}

func (w *IndexDBWriter) AddLexeme(docID int, lexeme string) error {
	file, info := w.findFileAndPosition(lexeme)

	// err := w.writeDocID(docID, file, info.Positions[len(info.Positions)-1])
	// if err != nil {
	// 	return err
	// }

	w.queues[file] <- lexemeWriteTask{
		DocID:    docID,
		Position: info.Positions[len(info.Positions)-1],
	}

	info.LastLength += 4
	if w.fileSizes[file] == info.Positions[len(info.Positions)-1] {
		w.fileSizes[file] += 4
	}
	return nil
}

func (w *IndexDBWriter) Close() error {
	for i := 0; i < maxFiles; i++ {
		close(w.queues[i])
	}
	for _, file := range w.files {
		file.Close()
	}
	for i := 0; i < maxFiles; i++ {
		<-w.workerStop
	}
	return nil
}
