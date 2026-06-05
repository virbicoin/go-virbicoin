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
	"enode://e8853097e51e155d6bac72b9ebff92ccdf5aa8ee9a845b92f6af969120b9ed017d6b79f231cf05f68e68237a824946f502b84914e0737e81c7a252a0174764f8@140.238.40.136:28329",  // bootnode-oci-ap-tokyo-1-ad-1-fd-1-01 (current)
	"enode://e8853097e51e155d6bac72b9ebff92ccdf5aa8ee9a845b92f6af969120b9ed017d6b79f231cf05f68e68237a824946f502b84914e0737e81c7a252a0174764f8@158.101.131.198:28329", // bootnode-oci-ap-tokyo-1-ad-1-fd-1-01 (new-ip, not yet active)
	"enode://141e23b0dc593a3e3474485d470df0057f2d710cd2d5dec733a7f2f8e6df81b034c4f5ff65e192af604c0d0f8bf2c53d2671190e3c9dcf78ada91e8810a43e54@158.101.135.30:28329",  // bootnode-oci-ap-tokyo-1-ad-1-fd-2-01
	"enode://77973b63ce0003ff40c1cfb10165a1045d156865ef682c65baf4501d09abe0f37705f9626ff2202a85be0f858935c82494d729ac35afa593e454a9c7a021c309@8.231.239.110:28329",   // bootnode-gcp-us-west1-a-01
	"enode://273bd300871544f39fcb3f0e9376c880cee1070db839da8ef6e1313a07cb6c11e6cdde0a78f1c8b1affb5b3f570f62c6cbddead3068aaf2e81b13f0509a0f737@13.202.197.147:28329",  // bootnode-aws-ap-south-1a-01
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
