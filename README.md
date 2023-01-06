# Blockchain app in go

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
- [ ] Merkle Tree
- [ ] Network Module

## Merkle tree

Merkel tree allow users to be part of the decentralized blockchain network without storing the full node (copy of the blockchain).
The user will be able to carry a ligth blockchain node which doesn't need to download the entirely blockchain and doesn't verify blocks or
transactions. It finds transactions in blocks and it's linked to a full node which will retrieve the necessary data for it. The idea of the tree
is to allow users to run a ton of ligth wallets nodes on top of just one full node. (For more information of why this is used refer to the
[Bitcoin's SVP White Paper](https://wiki.bitcoinsv.io/index.php/Simplified_Payment_Verification))