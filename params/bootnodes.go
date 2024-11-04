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

import "github.com/virbicoin/go-virbicoin/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{

    // VirBiCoin Official Go Bootnodes
	"enode://e8853097e51e155d6bac72b9ebff92ccdf5aa8ee9a845b92f6af969120b9ed017d6b79f231cf05f68e68237a824946f502b84914e0737e81c7a252a0174764f8@35.72.202.199:28329", // bootnode-aws-ap-northeast-1a-01
	"enode://94fe52c478d13d584afe04102bd00c4fb0f229194d2163d9564f6c609b8f62614cc339866ab2473aed95c43d622ece71cbf1e87cb9fba65351a2c3db45be2074@52.69.140.99:28329", // bootnode-aws-ap-northeast-1d-01
	"enode://141e23b0dc593a3e3474485d470df0057f2d710cd2d5dec733a7f2f8e6df81b034c4f5ff65e192af604c0d0f8bf2c53d2671190e3c9dcf78ada91e8810a43e54@13.208.46.232:28329", // bootnode-aws-ap-northeast-3a-01
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

// YoloV3Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// YOLOv3 ephemeral test network.
// TODO: Set Yolov3 bootnodes
var YoloV3Bootnodes = []string{}

var V5Bootnodes = []string{}

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
