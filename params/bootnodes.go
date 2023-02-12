// Copyright 2015 The go-ethereum Authors
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

package params

import "github.com/ethereum/go-ethereum/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// VirBiCoin Go Bootnodes
	"enode://c83eb39cd5c7e535eeeff8479a351369ee185983df21137589aee99caa48c070ba2949ac930b91aca7acdd07faa46a0f961e1382fc15d26d5599d1dc999769e0@35.72.202.199:28329", // bootnode-aws-ap-northeast-1c-001
	"enode://1a005cc477f38ca5a8dfa46c96149b23f7868b5dd4d532940bbe4a211d041e488d4faea841401724037ed0a9b986129dc99cfd17071c304c45f976d321c89f62@13.208.86.138:28329", // bootnode-aws-ap-northeast-3c-001
}

// RopstenBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var RopstenBootnodes = []string{}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Görli test network.
var GoerliBootnodes = []string{}

// YoloV2Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// YOLOv2 ephemeral test network.
var YoloV2Bootnodes = []string{}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "mainnet"
	case RopstenGenesisHash:
		net = "ropsten"
	case RinkebyGenesisHash:
		net = "rinkeby"
	case GoerliGenesisHash:
		net = "goerli"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
