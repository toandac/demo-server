package repoimpl

import (
	"demo-server/database"
	"demo-server/models"
	"demo-server/repository"
	"errors"
	"log"
	"sort"

	"github.com/zippoxer/bow"
)

type RecordRepoImpl struct {
	badger *database.BadgerDB
}

func NewRecordRepo(badger *database.BadgerDB) repository.RecordRepo {
	return &RecordRepoImpl{
		badger: badger,
	}
}

func (r *RecordRepoImpl) Get(id string, record *models.Record) error {
	if err := r.badger.BadgerDB.Bucket("records").Get(id, &record); err != nil {
		if !errors.Is(err, bow.ErrNotFound) {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (r *RecordRepoImpl) Put(record models.Record) error {
	if err := r.badger.BadgerDB.Bucket("records").Put(record); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *RecordRepoImpl) Iter(record models.Record) ([]models.Record, error) {
	var records []models.Record
	iter := r.badger.BadgerDB.Bucket("records").Iter()
	defer iter.Close()
	for iter.Next(&record) {
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].UpdatedAt.After(records[j].UpdatedAt)
	})

	return records, nil
}
