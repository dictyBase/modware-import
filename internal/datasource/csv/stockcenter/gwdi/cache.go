package gwdi

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/snowflake"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ListCacher interface {
	StartBatch()
	AppendToBatch([]byte, []byte)
	CommitBatch() error
	Push([]byte, []byte) error
	PushAll([]byte, ...[]byte) error
	IterateByPrefix([]byte) iterator.Iterator
	CommonPrefixes() ([][]byte, error)
}

type IdMapper interface {
	Put(key []byte, value []byte) error
	Get(key []byte) (value []byte, err error)
	Remove(key []byte) error
	Exist(key []byte) (bool, error)
	Iterate() iterator.Iterator
}

type leveldbListCache struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	node  *snowflake.Node
}

func NewListCache() (ListCacher, error) {
	b := &leveldbListCache{}
	n, err := snowflake.NewNode(1)
	if err != nil {
		return b,
			fmt.Errorf("error in getting a node for random string generator %w", err)
	}
	b.node = n
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return b, fmt.Errorf("error in opening key value store %w", err)
	}
	b.db = db
	return b, nil
}

func (b *leveldbListCache) CommonPrefixes() ([][]byte, error) {
	pm := make(map[string]int)
	re := regexp.MustCompile(`^(.+)\/{LIST}\/.*$`)
	itr := b.db.NewIterator(nil, nil)
	for itr.Next() {
		key := itr.Key()
		m := re.FindStringSubmatch(string(key))
		if m == nil {
			continue
		}
		pm[m[1]] = 1
	}
	var p [][]byte
	itr.Release()
	if err := itr.Error(); err != nil {
		return p, fmt.Errorf("error with cache iteration %w", err)
	}
	for k, _ := range pm {
		p = append(p, []byte(k))
	}
	return p, nil
}

func (b *leveldbListCache) IterateByPrefix(prefix []byte) iterator.Iterator {
	p := append(prefix, []byte("/{LIST}/")...)
	return b.db.NewIterator(util.BytesPrefix(p), nil)
}

func (b *leveldbListCache) PushAll(k []byte, allv ...[]byte) error {
	txn, err := b.db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("error in creating transaction %w", err)
	}
	for _, v := range allv {
		key := fmt.Sprintf("%s/{LIST}/%s", string(k), b.node.Generate().String())
		if err := txn.Put([]byte(key), v, nil); err != nil {
			txn.Discard()
			return fmt.Errorf(
				"error in storing key %s with value %s",
				string(k), string(v),
			)
		}
	}
	return txn.Commit()
}

func (b *leveldbListCache) StartBatch() {
	b.batch = new(leveldb.Batch)
}

func (b *leveldbListCache) AppendToBatch(k, v []byte) {
	b.batch.Put(
		[]byte(fmt.Sprintf("%s/{LIST}/%s", string(k), b.node.Generate().String())),
		v,
	)
}

func (b *leveldbListCache) CommitBatch() error {
	return b.db.Write(b.batch, nil)
}

func (b *leveldbListCache) Push(k, v []byte) error {
	key := fmt.Sprintf("%s/{LIST}/%s", string(k), b.node.Generate().String())
	txn, err := b.db.OpenTransaction()
	if err != nil {
		return fmt.Errorf("error in creating transaction %w", err)
	}
	if err := txn.Put([]byte(key), v, nil); err != nil {
		txn.Discard()
		return fmt.Errorf(
			"error in storing key %s with value %s",
			string(k), string(v),
		)
	}
	return txn.Commit()
}

type leveldbMap struct {
	db *leveldb.DB
}

func NewStrainMap() (IdMapper, error) {
	m := &leveldbMap{}
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return m, fmt.Errorf("error in opening key value store %w", err)
	}
	m.db = db
	return m, nil
}

func (m *leveldbMap) Iterate() iterator.Iterator {
	return m.db.NewIterator(util.BytesPrefix([]byte("GWDI")), nil)
}

func (m *leveldbMap) Put(k, v []byte) error {
	err := m.db.Put(k, v, nil)
	if err != nil {
		return fmt.Errorf(
			"error in writing key value in storage key:%s value:%s",
			string(k), string(v),
		)
	}
	return nil
}

func (m *leveldbMap) Get(k []byte) ([]byte, error) {
	v, err := m.db.Get(k, nil)
	if err != nil {
		return v, fmt.Errorf(
			"error in getting key value from key:%s",
			string(v),
		)
	}
	return v, nil
}

func (m *leveldbMap) Remove(k []byte) error {
	err := m.db.Delete(k, nil)
	if err != nil {
		return fmt.Errorf(
			"error in removing record with key %s", string(k),
		)
	}
	return nil
}

func (m *leveldbMap) Exist(k []byte) (bool, error) {
	ok, err := m.db.Has(k, nil)
	if err != nil {
		return ok, fmt.Errorf(
			"error in checking for key %s",
			string(k),
		)
	}
	return ok, nil
}
