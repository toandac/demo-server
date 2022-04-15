package database

import (
	"net/url"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/zippoxer/bow"
)

type BadgerDB struct {
	BadgerDB *bow.DB
}

func (bg *BadgerDB) NewBadgerDB(dbDSN *url.URL) {
	var err error
	bg.BadgerDB, err = bow.Open(dbDSN.Path, bow.SetBadgerOptions(
		badger.DefaultOptions(dbDSN.Path).
			WithTableLoadingMode(options.FileIO).
			WithValueLogLoadingMode(options.FileIO).
			WithNumVersionsToKeep(1).
			WithNumLevelZeroTables(1).
			WithNumLevelZeroTablesStall(2),
	))

	if err != nil {
		panic(err)
	}
}

func (bg *BadgerDB) Close() {
	bg.BadgerDB.Close()
}
