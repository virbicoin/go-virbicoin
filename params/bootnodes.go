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

import "github.com/emerauda/go-virbicoin/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// VirBiCoin Go Bootnodes
	"enode://a96b1e0d1e93987a81a08dcef602890d25e449ba2823b03bfb521cf4e72020e6af86ae87c6aae89b62d57a7ec380bced83108fb8810205db508bddf32f98f6e7@35.72.202.199:28329", // bootNode-AWS-ap-northeast-1a-001
	"enode://156c3707f0a68d1a38e398b32ba11529201d37be49f5990ca3294cadac2870f6c08b3db604dbdf306f93d1128589631356095c59121acab37e8be7c56ba75931@52.69.140.99:28329", // bootNode-AWS-ap-northeast-1d-001
	"enode://afb0a6b6a0a9b9881a384bea5c98f1997d0a0db7ae4ea89434e9b87ff65eb5e3f9a096d9ab757718f2915441f6d6ba8ab0ca9a6c452e1d7be9804bbcdf8ba662@13.208.46.232:28329", // bootNode-AWS-ap-northeast-3a-001
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
