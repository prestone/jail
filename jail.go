package jail

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/OneOfOne/xxhash"
	bolt "github.com/coreos/bbolt"
	"github.com/niubaoshu/gotiny"
)

const (
	shards = 1024
	sleep  = time.Minute * 10
)

func New(db ...string) (a *Jail) {
	a = new(Jail)
	a.hide.list = make([]map[int]int64, shards)
	a.hide.reason = make(map[interface{}]string)
	for x := 0; x < shards; x++ {
		a.hide.list[x] = make(map[int]int64)
	}
	if len(db) > 0 {
		r, err := bolt.Open(db[0], os.ModePerm, bolt.DefaultOptions)
		if err != nil {
			fmt.Println(err)
			return
		}
		a.db = r
		a.initdb()
	}
	return
}

type Jail struct {
	hide struct {
		sync.Mutex
		list   []map[int]int64
		reason map[interface{}]string
	}
	db    *bolt.DB
	count int
}

//int
func intBytes(i int) (r []byte) {
	r = make([]byte, 8)
	binary.LittleEndian.PutUint64(r, uint64(i))
	return
}

func bytesInt(b []byte) (i int) {
	if b == nil {
		return 0
	}
	return int(binary.LittleEndian.Uint64(b))
}

func marshal(id int, times int64, reason string) (bid, body []byte) {
	bid = intBytes(id)
	body = gotiny.Marshal(&id, &times, &reason)
	return
}

func unmarshal(v []byte) (id int, times int64, reason string) {
	gotiny.Unmarshal(v, &id, &times, &reason)
	return
}

func (a *Jail) Check(id interface{}) bool {
	sid := convert(id)
	a.hide.Lock()
	defer a.hide.Unlock()
	return a.hide.list[sid%shards][sid] > time.Now().Unix()
}

func (a *Jail) Reason(id interface{}) string {
	return a.hide.reason[convert(id)]
}

func (a *Jail) initdb() {
	var expired [][]byte
	a.db.Update(func(tx *bolt.Tx) (err error) {
		b, err := tx.CreateBucketIfNotExists([]byte("id"))
		now := time.Now().Unix()
		b.ForEach(func(k, v []byte) (err error) {
			id, until, reason := unmarshal(v)
			if until < now {
				expired = append(expired, k)
			}
			a.hide.list[id%shards][id] = until
			a.hide.reason[id] = reason
			a.count++
			return
		})
		for _, x := range expired {
			b.Delete(x)
		}
		return
	})
}

func (a *Jail) put(id int) {
	a.db.Update(func(tx *bolt.Tx) (err error) {
		b, err := tx.CreateBucketIfNotExists([]byte("id"))
		a.hide.Lock()
		b.Put(marshal(id, a.hide.list[id%shards][id], a.hide.reason[id]))
		a.hide.Unlock()
		return
	})
}

func (a *Jail) Put(id interface{}, d time.Duration, reason ...string) {
	sid := convert(id)
	a.hide.Lock()
	a.hide.list[sid%shards][sid] = time.Now().Add(d).Unix()
	if len(reason) > 0 {
		a.hide.reason[sid] = reason[0]
	}
	a.hide.Unlock()
	a.count++
	//has db
	if a.db != nil {
		go a.put(sid)
	}

}

func (a *Jail) Count() int {
	return a.count
}

func (a *Jail) delete(id int) {
	delete(a.hide.list[id%shards], id)
	delete(a.hide.reason, id)
	a.count--

	//has db
	if a.db != nil {
		a.db.Update(func(tx *bolt.Tx) (err error) {
			b, err := tx.CreateBucketIfNotExists([]byte("id"))
			b.Delete(intBytes(id))
			return
		})
	}
}

func (a *Jail) Delete(id interface{}) {
	sid := convert(id)
	a.hide.Lock()
	a.delete(sid)
	if a.db != nil {
		a.db.Update(func(tx *bolt.Tx) (err error) {
			b, err := tx.CreateBucketIfNotExists([]byte("id"))
			b.Delete(intBytes(sid))
			return
		})
	}
	a.hide.Unlock()
}

func (a *Jail) cleaner() {
	for {
		time.Sleep(sleep)
		a.hide.Lock()
		now := time.Now().Unix()
		for _, users := range a.hide.list {
			for id, times := range users {
				if times < now {
					a.delete(id)
					continue
				}
			}
		}
		a.hide.Unlock()
	}
}

func stringhash(v string) int {
	return int(xxhash.Checksum32([]byte(v)))
}

func convert(id interface{}) int {
	switch id.(type) {
	case int:
		return id.(int)
	case int8:
		return int(id.(int8))
	case int16:
		return int(id.(int16))
	case int32:
		return int(id.(int32))
	case int64:
		return int(id.(int64))
	case uint:
		return int(id.(uint))
	case uint8:
		return int(id.(uint8))
	case uint16:
		return int(id.(uint16))
	case uint32:
		return int(id.(uint32))
	case uint64:
		return int(id.(uint64))
	case string:
		return stringhash(id.(string))
	case []byte:
		return stringhash(string(id.([]byte)))
	default:
		return stringhash(fmt.Sprint(id))
	}
}
