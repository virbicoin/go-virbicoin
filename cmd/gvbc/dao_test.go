// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/virbicoin/go-virbicoin/common"
	"github.com/virbicoin/go-virbicoin/core/rawdb"
	"github.com/virbicoin/go-virbicoin/params"
)

// runGeth runs the given command and returns a struct with WaitExit method.
type CmdWrap struct {
	cmd *exec.Cmd
	t   *testing.T
}

func runGeth(t *testing.T, args ...string) *CmdWrap {
	cmd := exec.Command("gvbc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start geth: %v", err)
	}
	return &CmdWrap{cmd: cmd, t: t}
}

func (c *CmdWrap) WaitExit() {
	if err := c.cmd.Wait(); err != nil {
		c.t.Fatalf("gvbc exited with error: %v", err)
	}
}

// Genesis block for nodes which don't care about the DAO fork (i.e. not configured)
var daoOldGenesis = `{
	"alloc"      : {},
	"coinbase"   : "0x0000000000000000000000000000000000000000",
	"difficulty" : "0x20000",
	"extraData"  : "",
	"gasLimit"   : "0x2fefd8",
	"nonce"      : "0x0000000000000042",
	"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"timestamp"  : "0x00",
	"config"     : {
		"homesteadBlock" : 0
	}
}`

// Genesis block for nodes which actively oppose the DAO fork
var daoNoForkGenesis = `{
	"alloc"      : {},
	"coinbase"   : "0x0000000000000000000000000000000000000000",
	"difficulty" : "0x20000",
	"extraData"  : "",
	"gasLimit"   : "0x2fefd8",
	"nonce"      : "0x0000000000000042",
	"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"timestamp"  : "0x00",
	"config"     : {
		"homesteadBlock" : 0,
		"daoForkBlock"   : 0,
		"daoForkSupport" : false
	}
}`

// Genesis block for nodes which actively support the DAO fork
var daoProForkGenesis = `{
	"alloc"      : {},
	"coinbase"   : "0x0000000000000000000000000000000000000000",
	"difficulty" : "0x20000",
	"extraData"  : "",
	"gasLimit"   : "0x2fefd8",
	"nonce"      : "0x0000000000000042",
	"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"timestamp"  : "0x00",
	"config"     : {
		"homesteadBlock" : 0,
		"daoForkBlock"   : 0,
		"daoForkSupport" : true
	}
}`

var daoGenesisHash = common.HexToHash("5e1fc79cb4ffa4739177b5408045cd5d51c6cf766133f23f7cd72ee1f8d790e0")
var daoGenesisForkBlock = big.NewInt(0)

// TestDAOForkBlockNewChain tests that the DAO hard-fork number and the nodes support/opposition is correctly
// set in the database after various initialization procedures and invocations.
func TestDAOForkBlockNewChain(t *testing.T) {
	for i, arg := range []struct {
		genesis     string
		expectBlock *big.Int
		expectVote  bool
	}{
		// Test DAO Default Mainnet
		{"", params.MainnetChainConfig.DAOForkBlock, true},
		// test DAO Init Old Privnet
		{daoOldGenesis, nil, false},
		// test DAO Default No Fork Privnet
		{daoNoForkGenesis, daoGenesisForkBlock, false},
		// test DAO Default Pro Fork Privnet
		{daoProForkGenesis, daoGenesisForkBlock, true},
	} {
		testDAOForkBlockNewChain(t, i, arg.genesis, arg.expectBlock, arg.expectVote)
	}
}

func testDAOForkBlockNewChain(t *testing.T, test int, genesis string, expectBlock *big.Int, expectVote bool) {
	// Create a temporary data directory to use and inspect later
	datadir := tmpdir(t)
	defer os.RemoveAll(datadir)

	// Start a Gvbc instance with the requested flags set and immediately terminate
	if genesis != "" {
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", test, err)
		}
		runGeth(t, "--datadir", datadir, "--networkid", "1337", "init", json).WaitExit()
	} else {
		// Force chain initialization
		mainnetGenesisPath := filepath.Join("testdata", "test_genesis.json")
		runGeth(t, "--datadir", datadir, "--networkid", "1337", "init", mainnetGenesisPath).WaitExit()
		// After that, start the console
		args := []string{"--port", "0", "--networkid", "1337", "--maxpeers", "0", "--nodiscover", "--nat", "none", "--ipcdisable", "--datadir", datadir}
		runGeth(t, append(args, []string{"--exec", "2+2", "console"}...)...).WaitExit()
	}
	// Retrieve the DAO config flag from the database
	path := filepath.Join(datadir, "gvbc", "chaindata")
	db, err := rawdb.NewLevelDBDatabase(path, 0, 0, "")
	if err != nil {
		t.Fatalf("test %d: failed to open test database: %v", test, err)
	}
	defer db.Close()

	header := rawdb.ReadHeader(db, rawdb.ReadCanonicalHash(db, 0), 0)
	genesisHash := header.Hash()
	config := rawdb.ReadChainConfig(db, genesisHash)
	if config == nil {
		t.Errorf("test %d: failed to retrieve chain config: %v", test, err)
		return // we want to return here, the other checks can't make it past this point (nil panic).
	}
	// Validate the DAO hard-fork block number against the expected value
	if config.DAOForkBlock == nil {
		if expectBlock != nil {
			t.Errorf("test %d: dao hard-fork block mismatch: have nil, want %v", test, expectBlock)
		}
	} else if expectBlock == nil {
		t.Errorf("test %d: dao hard-fork block mismatch: have %v, want nil", test, config.DAOForkBlock)
	} else if config.DAOForkBlock.Cmp(expectBlock) != 0 {
		t.Errorf("test %d: dao hard-fork block mismatch: have %v, want %v", test, config.DAOForkBlock, expectBlock)
	}
	if config.DAOForkSupport != expectVote {
		t.Errorf("test %d: dao hard-fork support mismatch: have %v, want %v", test, config.DAOForkSupport, expectVote)
	}
}
