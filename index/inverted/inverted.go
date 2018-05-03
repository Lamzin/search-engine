package inverted

import (
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/boltdb/bolt"

	"github.com/lamzin/search-engine/index/model/doc"
)

const (
	shards    = 16
	queueSize = 1000
)

type InvertedIndex struct {
	IndexPath string

	queues         []chan saveTaks
	workerFinished chan bool
	open           bool
}

type saveTaks struct {
	ID    int
	Token string
}

func NewInvertedIndex(indexPath string) (*InvertedIndex, error) {
	var queues []chan saveTaks
	for i := 0; i < shards; i++ {
		queues = append(queues, make(chan saveTaks, queueSize))
	}

	index := InvertedIndex{
		IndexPath:      indexPath,
		queues:         queues,
		workerFinished: make(chan bool, shards),
		open:           true,
	}
	for i := 0; i < shards; i++ {
		go index.worker(i)
	}
	return &index, nil
}

func (i *InvertedIndex) Close() {
	i.open = false
	for j := 0; j < shards; j++ {
		<-i.workerFinished
	}
}

// func (i *InvertedIndex) GetDocIDs(token string) ([]int, error) {
// 	_, filePath := common.FilePath(token)

// 	lines, err := common.ReadFile(filepath.Join(i.IndexPath, filePath))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(lines) != 1 {
// 		return nil, fmt.Errorf("does not contain only one line")
// 	}

// 	numbers := strings.Split(strings.TrimSpace(lines[0]), " ")
// 	var arr []int
// 	for _, number := range numbers {
// 		n, err := strconv.ParseInt(number, 10, 32)
// 		if err != nil {
// 			return nil, err
// 		}
// 		arr = append(arr, (int)(n))
// 	}
// 	return arr, nil
// }

// func (i *InvertedIndex) AddToken(info *doc.DocInfo, token string) error {
// 	_, filePath := common.FilePath(token)
// 	return common.AppendFile(filepath.Join(i.IndexPath, filePath), fmt.Sprintf("%d ", info.ID))
// }

func (i *InvertedIndex) AddToken(info *doc.DocInfo, token string) error {
	i.queues[hash(token)] <- saveTaks{ID: info.ID, Token: token}
	return nil
}

// func (i *InvertedIndex) worker() {
// 	for i.open {
// 		l := len(i.queue)
// 		if l == 0 {
// 			time.Sleep(5 * time.Second)
// 		}
// 		log.Printf("Will process %d\n", l)

// 		queue := map[string][]int{}
// 		for index := 0; index < l; index++ {
// 			task := <-i.queue
// 			queue[task.Token] = append(queue[task.Token], task.ID)
// 		}

// 		err := i.db.Batch(func(tx *bolt.Tx) error {
// 			b := tx.Bucket([]byte("tokens_to_doc"))

// 			for token, docIDs := range queue {
// 				var ids string
// 				if key := b.Get([]byte(token)); key != nil {
// 					ids = string(key)
// 				}
// 				for _, id := range docIDs {
// 					ids += fmt.Sprintf("%d ", id)
// 				}
// 				if err := b.Put([]byte(token), []byte(ids)); err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		})
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 	}
// 	i.workerFinished <- true
// }

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return (int)(h.Sum32() % shards)
}

func (i *InvertedIndex) worker(workerIndex int) {
	db, err := bolt.Open(i.IndexPath+fmt.Sprintf("%d.bolt.db", workerIndex), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i.open {
		l := len(i.queues[workerIndex])
		if l == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("Worker %d: will process %d\n", workerIndex, l)

		queue := map[string][]int{}
		for index := 0; index < l; index++ {
			task := <-i.queues[workerIndex]
			queue[task.Token] = append(queue[task.Token], task.ID)
		}

		err := db.Batch(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("tokens_to_doc"))
			if err != nil {
				return err
			}

			for token, docIDs := range queue {
				var ids string
				if key := b.Get([]byte(token)); key != nil {
					ids = string(key)
				}
				for _, id := range docIDs {
					ids += fmt.Sprintf("%d ", id)
				}
				if err := b.Put([]byte(token), []byte(ids)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Println(err.Error())
		}
	}
	i.workerFinished <- true
}
