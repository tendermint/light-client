package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/abci/example/dummy"
	"github.com/tendermint/tendermint/types"
)

// Make sure status is correct (we connect properly)
func TestStatus(t *testing.T) {
	c := GetClient()
	status, err := c.Status()
	if assert.Nil(t, err) {
		assert.Equal(t, GetConfig().GetString("chain_id"), status.NodeInfo.Network)
	}
}

// Make some app checks
func TestAppCalls(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	c := GetClient()
	_, err := c.Block(1)
	assert.NotNil(err) // no block yet
	k, v, tx := TestTxKV()
	_, err = c.BroadcastTxCommit(tx)
	require.Nil(err)
	// wait before querying
	time.Sleep(time.Second)
	qres, err := c.ABCIQuery(k)
	if assert.Nil(err) {
		data := dummy.QueryResult{}
		// make sure dummy app return value...
		err = json.Unmarshal(qres.Result.Data, &data)
		if assert.Nil(err) {
			assert.Equal(string(v), data.Value)
			assert.True(data.Exists)
		}
	}
	// and we can even check the block is added
	_, err = c.Block(1)
	assert.Nil(err) // now it's good :)
}

// run most calls just to make sure no syntax errors
func TestNoErrors(t *testing.T) {
	assert := assert.New(t)
	c := GetClient()
	_, err := c.NetInfo()
	assert.Nil(err)
	_, err = c.BlockchainInfo(0, 4)
	assert.Nil(err)
	// TODO: check with a valid height
	_, err = c.Block(1000)
	assert.NotNil(err)
	// maybe this is an error???
	// _, err = c.DialSeeds([]string{"one", "two"})
	// assert.Nil(err)
	gen, err := c.Genesis()
	if assert.Nil(err) {
		assert.Equal(GetConfig().GetString("chain_id"), gen.Genesis.ChainID)
	}
}

func TestSubscriptions(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	c := GetClient()
	err := c.StartWebsocket()
	require.Nil(err)
	defer c.StopWebsocket()

	// subscribe to a transaction event
	_, _, tx := TestTxKV()
	// This DOES NOT cause an error!
	// eventType := types.EventStringNewBlock()
	// this causes a panic in tendermint core!!!
	eventType := types.EventStringTx(types.Tx(tx))
	c.Subscribe(eventType)
	read := 0

	// set up a listener
	r, e := c.GetEventChannels()
	go func() {
		// read one event in the background
		select {
		case <-r:
			// TODO: actually parse this or something
			read += 1
		case err := <-e:
			panic(err)
		}
	}()

	// make sure nothing has happened yet.
	assert.Equal(0, read)

	// send a tx and wait for it to propogate
	_, err = c.BroadcastTxCommit(tx)
	assert.Nil(err, string(tx))
	// wait before querying
	time.Sleep(time.Second)

	// now make sure the event arrived
	assert.Equal(1, read)
}
