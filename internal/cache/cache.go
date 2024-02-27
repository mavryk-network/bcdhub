package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/karlseguin/ccache"
	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/microcosm-cc/bluemonday"
)

// Cache -
type Cache struct {
	*ccache.Cache
	rpc noderpc.INode

	accounts  account.Repository
	contracts contract.Repository
	protocols protocol.Repository
	sanitizer *bluemonday.Policy
}

// NewCache -
func NewCache(rpc noderpc.INode, accounts account.Repository, contracts contract.Repository, protocols protocol.Repository) *Cache {
	sanitizer := bluemonday.UGCPolicy()
	sanitizer.AllowAttrs("em")
	return &Cache{
		ccache.New(ccache.Configure().MaxSize(1000)),
		rpc,
		accounts,
		contracts,
		protocols,
		sanitizer,
	}
}

// ContractTags -
func (cache *Cache) ContractTags(ctx context.Context, address string) (types.Tags, error) {
	if !bcd.IsContract(address) {
		return 0, nil
	}

	key := fmt.Sprintf("contract:%s", address)
	item, err := cache.Fetch(key, time.Minute*10, func() (interface{}, error) {
		c, err := cache.contracts.Get(ctx, address)
		if err != nil {
			return 0, err
		}
		return c.Tags, nil
	})
	if err != nil {
		cache.Delete(key)
		return 0, err
	}
	return item.Value().(types.Tags), nil
}

// TezosBalance -
func (cache *Cache) TezosBalance(ctx context.Context, address string, level int64) (int64, error) {
	key := fmt.Sprintf("tezos_balance:%s:%d", address, level)
	item, err := cache.Fetch(key, 30*time.Second, func() (interface{}, error) {
		return cache.rpc.GetContractBalance(ctx, address, level)
	})
	if err != nil {
		cache.Delete(key)
		return 0, err
	}
	return item.Value().(int64), nil
}

// StorageTypeBytes -
func (cache *Cache) StorageTypeBytes(ctx context.Context, address, symLink string) ([]byte, error) {
	if !bcd.IsContract(address) {
		return nil, nil
	}

	key := fmt.Sprintf("storage:%s", address)
	item, err := cache.Fetch(key, 5*time.Minute, func() (interface{}, error) {
		return cache.contracts.ScriptPart(ctx, address, symLink, consts.STORAGE)
	})
	if err != nil {
		cache.Delete(key)
		return nil, err
	}
	return item.Value().([]byte), nil
}

// ProtocolByID -
func (cache *Cache) ProtocolByID(ctx context.Context, id int64) (protocol.Protocol, error) {
	key := fmt.Sprintf("protocol_id:%d", id)
	item, err := cache.Fetch(key, time.Hour, func() (interface{}, error) {
		return cache.protocols.GetByID(ctx, id)
	})
	if err != nil {
		cache.Delete(key)
		return protocol.Protocol{}, err
	}
	return item.Value().(protocol.Protocol), nil
}

func (cache *Cache) Script(ctx context.Context, address, symLink string) (contract.Script, error) {
	key := fmt.Sprintf("script:%s:%s", address, symLink)
	item, err := cache.Fetch(key, time.Hour, func() (interface{}, error) {
		return cache.contracts.Script(ctx, address, symLink)
	})
	if err != nil {
		cache.Delete(key)
		return contract.Script{}, err
	}
	return item.Value().(contract.Script), nil
}

func (cache *Cache) ScriptBytes(ctx context.Context, address, symLink string) ([]byte, error) {
	key := fmt.Sprintf("script_bytes:%s:%s", address, symLink)
	item, err := cache.Fetch(key, time.Hour, func() (interface{}, error) {
		script, err := cache.contracts.Script(ctx, address, symLink)
		if err != nil {
			return script, err
		}
		return script.Full()
	})
	if err != nil {
		cache.Delete(key)
		return nil, err
	}
	return item.Value().([]byte), nil
}
