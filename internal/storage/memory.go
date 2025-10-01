package storage

import "sync"

type DBEntry struct {
	MessageID   int
	Movies      []string
	Page        int
	RequesterID int64 // ID of user who requested, used to restrict paging
}

var (
	mu  sync.RWMutex
	db  = make(map[string]DBEntry)
)

func Save(query string, entry DBEntry) {
	mu.Lock()
	defer mu.Unlock()
	db[query] = entry
}

func Get(query string) (DBEntry, bool) {
	mu.RLock()
	defer mu.RUnlock()
	e, ok := db[query]
	return e, ok
}

func Delete(query string) {
	mu.Lock()
	defer mu.Unlock()
	delete(db, query)
}
