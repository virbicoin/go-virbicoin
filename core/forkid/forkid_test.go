// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package forkid

import (
	"bytes"
	"math"
	"testing"

	"github.com/virbicoin/go-virbicoin/common"
	"github.com/virbicoin/go-virbicoin/params"
	"github.com/virbicoin/go-virbicoin/rlp"
)

// TestCreation tests that different genesis and fork rule combinations result in
// the correct fork ID.
func TestCreation(t *testing.T) {
	type testcase struct {
		head uint64
		want ID
	}
	tests := []struct {
		config  *params.ChainConfig
		genesis common.Hash
		cases   []testcase
	}{
		// VirBiCoin mainnet, including reward reductions.
		{
			params.MainnetChainConfig,
			params.MainnetGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xf5da9347), Next: 4200000}},
				{4199999, ID{Hash: checksumToBytes(0xf5da9347), Next: 4200000}},
				{4200000, ID{Hash: checksumToBytes(0xb1f5e1df), Next: 6300000}},
				{6300000, ID{Hash: checksumToBytes(0xfeca00ef), Next: 8400000}},
				{8400000, ID{Hash: checksumToBytes(0x4d19c844), Next: 10500000}},
				{10500000, ID{Hash: checksumToBytes(0x2c145385), Next: 12600000}},
				{12600000, ID{Hash: checksumToBytes(0x38365c02), Next: 14700000}},
				{14700000, ID{Hash: checksumToBytes(0xcc3f4fb0), Next: 16800000}},
				{16800000, ID{Hash: checksumToBytes(0x0b9b3d84), Next: 100000000}},
				{100000000, ID{Hash: checksumToBytes(0xb822c4ee), Next: 500000000}},
				{500000000, ID{Hash: checksumToBytes(0x8277a1f9), Next: 0}},
			},
		},
		// Ropsten test cases
		{
			params.RopstenChainConfig,
			params.RopstenGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0x30c7ddbc), Next: 10}},            // Unsynced, last Frontier, Homestead and first Tangerine block
				{9, ID{Hash: checksumToBytes(0x30c7ddbc), Next: 10}},            // Last Tangerine block
				{10, ID{Hash: checksumToBytes(0x63760190), Next: 1700000}},      // First Spurious block
				{1699999, ID{Hash: checksumToBytes(0x63760190), Next: 1700000}}, // Last Spurious block
				{1700000, ID{Hash: checksumToBytes(0x3ea159c7), Next: 4230000}}, // First Byzantium block
				{4229999, ID{Hash: checksumToBytes(0x3ea159c7), Next: 4230000}}, // Last Byzantium block
				{4230000, ID{Hash: checksumToBytes(0x97b544f3), Next: 4939394}}, // First Constantinople block
				{4939393, ID{Hash: checksumToBytes(0x97b544f3), Next: 4939394}}, // Last Constantinople block
				{4939394, ID{Hash: checksumToBytes(0xd6e2149b), Next: 6485846}}, // First Petersburg block
				{6485845, ID{Hash: checksumToBytes(0xd6e2149b), Next: 6485846}}, // Last Petersburg block
				{6485846, ID{Hash: checksumToBytes(0x4bc66396), Next: 7117117}}, // First Istanbul block
				{7117116, ID{Hash: checksumToBytes(0x4bc66396), Next: 7117117}}, // Last Istanbul block
				{7117117, ID{Hash: checksumToBytes(0x6727ef90), Next: 0}},       // First Muir Glacier block
				{7500000, ID{Hash: checksumToBytes(0x6727ef90), Next: 0}},       // Future
			},
		},
		// Rinkeby test cases
		{
			params.RinkebyChainConfig,
			params.RinkebyGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0x3b8e0691), Next: 1}},             // Unsynced, last Frontier block
				{1, ID{Hash: checksumToBytes(0x60949295), Next: 2}},             // First and last Homestead block
				{2, ID{Hash: checksumToBytes(0x8bde40dd), Next: 3}},             // First and last Tangerine block
				{3, ID{Hash: checksumToBytes(0xcb3a64bb), Next: 1035301}},       // First Spurious block
				{1035300, ID{Hash: checksumToBytes(0xcb3a64bb), Next: 1035301}}, // Last Spurious block
				{1035301, ID{Hash: checksumToBytes(0x8d748b57), Next: 3660663}}, // First Byzantium block
				{3660662, ID{Hash: checksumToBytes(0x8d748b57), Next: 3660663}}, // Last Byzantium block
				{3660663, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // First Constantinople block
				{4321233, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // Last Constantinople block
				{4321234, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // First Petersburg block
				{5435344, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // Last Petersburg block
				{5435345, ID{Hash: checksumToBytes(0xcbdb8838), Next: 0}},       // First Istanbul block
				{6000000, ID{Hash: checksumToBytes(0xcbdb8838), Next: 0}},       // Future Istanbul block
			},
		},
		// Goerli test cases
		{
			params.GoerliChainConfig,
			params.GoerliGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xa3f5ab08), Next: 1561651}},       // Unsynced, last Frontier, Homestead, Tangerine, Spurious, Byzantium, Constantinople and first Petersburg block
				{1561650, ID{Hash: checksumToBytes(0xa3f5ab08), Next: 1561651}}, // Last Petersburg block
				{1561651, ID{Hash: checksumToBytes(0xc25efa5c), Next: 0}},       // First Istanbul block
				{2000000, ID{Hash: checksumToBytes(0xc25efa5c), Next: 0}},       // Future Istanbul block
			},
		},
	}
	for i, tt := range tests {
		for j, ttt := range tt.cases {
			if have := NewID(tt.config, tt.genesis, ttt.head); have != ttt.want {
				t.Errorf("test %d, case %d: fork ID mismatch: have %x, want %x", i, j, have, ttt.want)
			}
		}
	}
}

// TestValidation tests that a local peer correctly validates and accepts a remote
// fork ID.
func TestValidation(t *testing.T) {
	tests := []struct {
		head uint64
		id   ID
		err  error
	}{
		// Same fork, with and without a future fork announcement.
		{8400001, ID{Hash: checksumToBytes(0x4d19c844), Next: 0}, nil},
		{8400001, ID{Hash: checksumToBytes(0x4d19c844), Next: math.MaxUint64}, nil},
		{8399999, ID{Hash: checksumToBytes(0xfeca00ef), Next: 8400000}, nil},
		{8399999, ID{Hash: checksumToBytes(0xfeca00ef), Next: 0}, nil},
		{8399999, ID{Hash: checksumToBytes(0xfeca00ef), Next: math.MaxUint64}, nil},

		// Peers syncing across a known fork in either direction.
		{8400001, ID{Hash: checksumToBytes(0xfeca00ef), Next: 8400000}, nil},
		{8400001, ID{Hash: checksumToBytes(0xb1f5e1df), Next: 6300000}, nil},
		{8399999, ID{Hash: checksumToBytes(0x4d19c844), Next: 0}, nil},
		{6299999, ID{Hash: checksumToBytes(0xfeca00ef), Next: 0}, nil},

		// A peer on the previous fork must know the local activation block.
		{8400000, ID{Hash: checksumToBytes(0xfeca00ef), Next: 8400000}, nil},
		{8400000, ID{Hash: checksumToBytes(0xfeca00ef), Next: 0}, ErrRemoteStale},
		{8400001, ID{Hash: checksumToBytes(0xfeca00ef), Next: 0}, ErrRemoteStale},

		// Unknown checksums and remote-only forks already reached locally.
		{8400001, ID{Hash: checksumToBytes(0x5cddc0e1), Next: 0}, ErrLocalIncompatibleOrStale},
		{8399999, ID{Hash: checksumToBytes(0x5cddc0e1), Next: 0}, ErrLocalIncompatibleOrStale},
		{8400001, ID{Hash: checksumToBytes(0xafec6b27), Next: 0}, ErrLocalIncompatibleOrStale},
		{588888888, ID{Hash: checksumToBytes(0x8277a1f9), Next: 588888888}, ErrLocalIncompatibleOrStale},
		{8399999, ID{Hash: checksumToBytes(0xfeca00ef), Next: 8399999}, ErrLocalIncompatibleOrStale},
	}
	for i, tt := range tests {
		filter := newFilter(params.MainnetChainConfig, params.MainnetGenesisHash, func() uint64 { return tt.head })
		if err := filter(tt.id); err != tt.err {
			t.Errorf("test %d: validation error mismatch: have %v, want %v", i, err, tt.err)
		}
	}
}

// Tests that IDs are properly RLP encoded (specifically important because we
// use uint32 to store the hash, but we need to encode it as [4]byte).
func TestEncoding(t *testing.T) {
	tests := []struct {
		id   ID
		want []byte
	}{
		{ID{Hash: checksumToBytes(0), Next: 0}, common.Hex2Bytes("c6840000000080")},
		{ID{Hash: checksumToBytes(0xdeadbeef), Next: 0xBADDCAFE}, common.Hex2Bytes("ca84deadbeef84baddcafe")},
		{ID{Hash: checksumToBytes(math.MaxUint32), Next: math.MaxUint64}, common.Hex2Bytes("ce84ffffffff88ffffffffffffffff")},
	}
	for i, tt := range tests {
		have, err := rlp.EncodeToBytes(tt.id)
		if err != nil {
			t.Errorf("test %d: failed to encode forkid: %v", i, err)
			continue
		}
		if !bytes.Equal(have, tt.want) {
			t.Errorf("test %d: RLP mismatch: have %x, want %x", i, have, tt.want)
		}
	}
}
