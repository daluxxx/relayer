package relayer

import (
	"sync"

	tmclient "github.com/cosmos/cosmos-sdk/x/ibc/07-tendermint/types"
)

// NewSyncHeaders returns a new instance of map[string]*tmclient.Header that can be easily
// kept "reasonably up to date"
func NewSyncHeaders(chains ...*Chain) (*SyncHeaders, error) {
	mp, err := UpdatesWithHeaders(chains...)
	if err != nil {
		return nil, err
	}
	return &SyncHeaders{hds: mp}, nil
}

// SyncHeaders is an instance of map[string]*tmclient.Header
// that can be kept "reasonably up to date" using it's Update method
type SyncHeaders struct {
	sync.Mutex

	hds map[string]*tmclient.Header
}

// Update the header for a given chain
func (uh *SyncHeaders) Update(c *Chain) error {
	hd, err := c.UpdateLiteWithHeader()
	if err != nil {
		return err
	}
	uh.Lock()
	defer uh.Unlock()
	uh.hds[c.ChainID] = hd
	return nil
}

// GetHeader returns the latest header for a given chainID
func (uh *SyncHeaders) GetHeader(chainID string) *tmclient.Header {
	uh.Lock()
	defer uh.Unlock()
	return uh.hds[chainID]
}

// GetHeight returns the latest height for a given chainID
func (uh *SyncHeaders) GetHeight(chainID string) uint64 {
	uh.Lock()
	defer uh.Unlock()
	return uh.hds[chainID].GetHeight()
}
