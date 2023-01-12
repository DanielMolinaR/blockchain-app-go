# Blockchain app in go

This project aims to touch ground with the basic logic of blockchain.

The blockchain developed here is centralized but it has a Network of different type of nodes that can intereact with the main node.

## Roadmap

- [x] Base blockchain app setup
- [x] Proof of Work
- [x] BadgerDB
- [x] CLI
- [x] Transactions
- [x] Basic Wallet
- [x] Persistance
- [x] Digital Signatures
- [x] UTXO Persintance layer
- [x] Merkle Tree
- [x] Network Module

## Merkle tree

Merkel tree is used in this project because it allows users to be part of the decentralized blockchain network without storing the full node (copy of the blockchain).
The user will be able to carry a ligth blockchain node which doesn't need to download the entirely blockchain and doesn't verify blocks or
transactions. It finds transactions in blocks and it's linked to a full node which will retrieve the necessary data for it. The idea of the tree
is to allow users to run a ton of ligth wallets nodes on top of just one full node. (For more information of why this is used refer to the
[Bitcoin's SVP White Paper](https://wiki.bitcoinsv.io/index.php/Simplified_Payment_Verification))

## How to use it

To interact with this blockchain project you need <code>Go version 1.13+</code>.

To make it work you will need at least to terminals open at the same time. First, in bot of them you need to set up an ENV variable which represents the port that the node will be listening to.

<code>NODE_ID=3000</code>

Once setted the port you need to create a wallet for each terminal by typing:

<code>go run main.go createwallet</code>

After that, you will have to select one of the wallets as the main blockchain. Then you will have to run:

<code>go run main.go createblockchain --address {wallet_adress} </code>

So all the nodes star with the same genesis block we will need to create the necesary folder for it from the actual blockchain folder created by the precious command

<code>cp tmp && cp -R blocks_{Main node id} blocks_{node_id}</code>

Once done that you can start the nodes from each of the terminals so they can download and update the blockchain.

<code>go run main.go startnode</code>

You can also start de node as a miner with the ***-miner*** flag followed by the wallet address.

Other commands that can be run are:

<code>
go run main.go send -from {wallet_address_1} -to {wallet_address_2} -amount 10 
<br>
go run main.go send -from {wallet_address_1} -to {wallet_address_2} -amount 10 -mine
<br>
go run main.go getbalance --address {wallet_address}
<br>
go run main.go printchain
<br>


</code>