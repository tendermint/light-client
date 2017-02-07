package lightclient

import "sort"

// Storage has many implementation, based on security and sharing requirements
// like disk-backed, mem-backed, vault, db, etc.
type Storage interface {
	Put(name string, key []byte, info KeyInfo) error
	Get(name string) ([]byte, KeyInfo, error)
	List() ([]KeyInfo, error)
	Delete(name string) error
}

// KeyInfos is a wrapper to allows alphabetical sorting of the keys
type SortKeys []KeyInfo

func (k SortKeys) Len() int           { return len(k) }
func (k SortKeys) Less(i, j int) bool { return k[i].Name < k[j].Name }
func (k SortKeys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

func (k SortKeys) Sort() {
	if k != nil {
		sort.Sort(k)
	}
}
