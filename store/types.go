package store

import (
	"github.com/cosmos/cosmos-sdk/abciapp"
)

// Import cosmos-sdk/types/store.go for convenience.
// nolint
type Store = abciapp.Store
type Committer = abciapp.Committer
type CommitStore = abciapp.CommitStore
type MultiStore = abciapp.MultiStore
type CacheMultiStore = abciapp.CacheMultiStore
type CommitMultiStore = abciapp.CommitMultiStore
type KVStore = abciapp.KVStore
type KVPair = abciapp.KVPair
type Iterator = abciapp.Iterator
type CacheKVStore = abciapp.CacheKVStore
type CommitKVStore = abciapp.CommitKVStore
type CacheWrapper = abciapp.CacheWrapper
type CacheWrap = abciapp.CacheWrap
type CommitID = abciapp.CommitID
type StoreKey = abciapp.StoreKey
type StoreType = abciapp.StoreType
type Queryable = abciapp.Queryable
