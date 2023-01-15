package samples

import (
	"fmt"
)

type (
	Repository interface {
		Get(key string) string
		Set(key, value string)
	}
	DbRepository struct {
		db map[string]string
	}
)

func NewDbRepository() Repository {
	repo := &DbRepository{db: make(map[string]string)}
	// initialize the db
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("%d", i)
		repo.Set(key, key)
	}
	return repo
}

func (r DbRepository) Get(key string) string {
	fmt.Printf("Getting %v from dbrepository...\n", key)
	return r.db[key]
}

func (r DbRepository) Set(key, value string) {
	fmt.Printf("Setting %v from dbrepository...\n", key)
	r.db[key] = value
}
