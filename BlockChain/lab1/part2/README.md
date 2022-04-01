# Lab 1: Creating a Blockchain

| Lab 1:           | Creating a Blockchain        |
| ---------------- | ---------------------------- |
| Subject:         | DAT650 Blockchain Technology |
| Deadline:        | 07. SEP                      |
| Expected effort: | 2 weeks                      |
| Grading:         | Pass/fail                    |

## Table of Contents

- [Lab 1: Creating a Blockchain](#lab-1-creating-a-blockchain)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Part 1](#part-1)
    - [Blocks](#blocks)
    - [Blockchain](#blockchain)
  - [Part 2](#part-2)
    - [Merkle Tree](#merkle-tree)
  - [Part 3](#part-3)
    - [Proof-Of-Work](#proof-of-work)
  - [Part 4](#part-4)
    - [Command line client](#command-line-client)
    - [Merkle tree benchmarks](#merkle-tree-benchmarks)
  - [Lab Approval](#lab-approval)

## Introduction

The main objective of this lab assignment is to build a simplified blockchain.
A blockchain is basically a distributed database of records. 
What makes it unique is that it’s uses cryptographic hash functions to create a 
tamper-proof mechanism of committed transactions through distributed consensus.
Most blockchains are permissionless, which means that they allow public membership of nodes,
often implemented on top of a peer-to-peer network, allowing a public distributed
database, i.e. everyone who uses the database has a full or partial copy of it.
A new record can be added only after the consensus between the other keepers of the database.
Also, it’s blockchain that made crypto-currencies and smart contracts possible.

This lab consist of four parts. Each part will be explained in more detail in their own sections.

1. **The chain of blocks:** Implement a chain of blocks as an ordered, back-linked list data structure.
   Use the provided skeleton code and unit tests.

2. **Efficient transactions and blocks verification:** Implement a efficient way
   to verify membership of certain transactions in a block using [Merkle Trees](https://en.wikipedia.org/wiki/Merkle_tree). 
   Use the provided skeleton code and unit tests.

3. **Proof-of-work**: Implement a proof-of-work mechanism to add new blocks in the blockchain.

4. **Command line client and benchmarks**: Develope a command-line tool to test your blockchain and benchmark your Merkle tree implementation.

For each part of the assignment you should copy your implementation of the previous part (and sorry for that, it may be a bit inconvenient to have to do it, but we are working on improve that). But **do not copy the tests**, they can differ from each part, copy only your implementation. If you prefer, you can create a new branch for each part of the assignment.

## Part 1

### Blocks
In blockchain it’s the blocks that store valuable information.
For example, Bitcoin blocks store [transactions](https://en.bitcoin.it/wiki/Transaction), the essence of any crypto-currency. 
Besides this, a block contains some technical information, like its version, current timestamp and the hash of the previous block.
In this assignment we will not implement the block as it’s described in current deployed blockchains or Bitcoin specifications, instead we’ll use a simplified version of it, which contains only significant information for learn purposes. Our block definition is defined in the `block.go` file and has the following structure:

```go
type Block struct {
	Timestamp     int64          
	Transactions  []*Transaction 
	PrevBlockHash []byte         
	Hash          []byte         
}

type Transaction struct {
	Data []byte
}
```

_Timestamp_ is the current timestamp (when the block is created), _Transactions_ is the actual valuable information containing in the block, _PrevBlockHash_ stores the hash of the previous block, and _Hash_ is the hash of the block.
In Bitcoin specification _Timestamp_, _PrevBlockHash_, and _Hash_ are block headers, which form a separate data structure, and _Transactions_ is a separate data structure (for now, our transaction is only a Two-dimensional slice of bytes contain the data to be stored). You can read more about how transactions are implemented [here](https://bitcoin.org/en/transactions-guide#introduction).

Each block is linked to the previous one using a hash function.
The way hashes are calculates is very important feature of blockchain, and it’s this feature that makes blockchain secure.
The thing is that calculating a hash is a computationally difficult operation, it takes some time even on fast computers.
This is an intentional architectural design of blockchain systems, which makes adding new blocks difficult, thus preventing their modification after they’re added.
We’ll discuss and implement this mechanism in the [part 3](#part-3) of this lab.

For now, you will just take block fields (i.e. headers), concatenate them, and calculate a SHA-256 hash on the concatenated combination. To do that, use the `SetHash` function. Feed the `PrevBlockHash`, `Transactions`, and `Timestamp` into the hash in this order.

To compute the SHA-256 checksum of the data you can use the [Sum256](https://golang.org/pkg/crypto/sha256/#Sum256) function from the go crypto package. The function receives as input parameters a byte slice, so you will need to convert each field of the header to bytes and then concatenate it. To convert the `timestamp`, you can use the function `IntToHex` in the file `utils.go`, this function uses the package [binary](https://golang.org/pkg/encoding/binary/) from the go standard library to write the binary representation of some data.

We also want all transactions in a block to be uniquely identified by a single hash.
To achieve this, you will get each transaction, concatenate them, and get a hash of the concatenated combination.
This hashing mechanism of providing unique representation of data will be given by the `HashTransactions` function, that will take all transactions of a block and return the hash of it to be included in the block _Hash_.

### Blockchain

Now let’s implement a blockchain.
In its essence blockchain is just a database with certain structure: it’s an ordered, back-linked list. 
Which means that blocks are stored in the insertion order and that each block is linked to the previous one.
This structure allows to quickly get the latest block in a chain and to get a block by its hash.

In Golang this structure can be implemented by using an array and a map: the array would keep ordered hashes (arrays are ordered in Go), and the map would keep hash to block pairs (maps are unordered).
But for now, in your blockchain prototype you just need to use an array as shown below.

```go
type Blockchain struct {
	blocks []*Block
}
```

As every block need to be linked to the previous one, to add a new block we need an existing block, but there’re no blocks in the blockchain on the beginning.
So, in any blockchain there must be at least one block, and such block is the first in the chain and is called genesis block.

Your task is implement all functions marked with `TODO(student)` in the file `blockchain.go`.
These functions are:
 - `NewBlockchain`: This function should creates a new blockchain initializing a Genesis block with the 
   hardcoded data `genesisCoinbaseData`.
   You can use the function `NewGenesisBlock` of the `block.go` to create the Genesis block.
 - `addBlock`: This function should get the previous block hash and add a new block linking it to the previous.
 - `GetGenesisBlock`: This function should return the Genesis block.
 - `CurrentBlock`: This function should return the last block.
 - `GetBlock`: This function should return a block based on its hash. 

## Part 2

### Merkle Tree

Until now we are using hashing as a mechanism of providing a unique representation of data, which give to us
an easy way to verify data integrity, i.e. if any of the transaction data in a block changes the root hash will change, and tampering is identified.
We did that in the function `HashTransactions` in the `block.go` file, by getting each transaction in a block, concatenate them in a specific order and applied SHA256 to the concatenated combination.
But besides uniquely identify all the transactions in a block by a single hash, for efficiency, we also want to be able to easily verify if some transaction is in the block without requiring to have all the block transactions.

[Merkle trees](https://xlinux.nist.gov/dads/HTML/MerkleTree.html) are used by [Bitcoin](https://bitcoin.org/bitcoin.pdf) to obtain transactions hash, which is then saved in block headers and is considered by the proof-of-work system.
The benefit of Merkle trees is that a node can verify membership of certain transaction without downloading the whole block, just using the transaction hash, the root hash of the merkle tree, and a set of intermediate hashes necessary to reconstruct the merkle path for that transaction, which is know as merkle proof.
The Merkle path is simply the set of hashes from the transaction at the leaf node to the Merkle root.
A Merkle proof is a way of proving that a certain transaction is part of a merkle tree without requiring any of the other transactions to be exposed, just the hashes.
Each hash in the proof is the sibling of the hash in the path at the same level in the tree.

This optimization mechanism is crucial for the successful adoption of Bitcoin or any [permissionless blockchain](https://eprint.iacr.org/2017/375.pdf).
For example, the full Bitcoin database (i.e., blockchain) currently takes [more than 230 Gb of disk space](https://www.blockchain.com/charts/blocks-size).
Because of the decentralized nature of Bitcoin, every node in the network must be independent and self-sufficient, i.e. every node in the network must store a full copy of the blockchain.
With many people starting using Bitcoin, this rule becomes more difficult to follow: it’s not likely that everyone will run a full node.
Also, since nodes are full-fledged participants of the network, they have responsibilities: they must verify transactions and blocks.
Also, there’s certain internet traffic required to interact with other nodes and download new blocks.

The above mechanism also enables SPV (Simple Payment Verification) in Bitcoin, allowing the creation of "light clients" that only store block headers (which includes the Merkle root) to perform transaction inclusion verifications.
Thus a light client doesn’t verify blocks and transactions, instead, it finds transactions in blocks (to verify payments) and maintain a connection with a full node to retrieve just necessary data.
This mechanism allows having multiple light nodes with running just one full node, but can also impose some centralization, since incentive less nodes to maintain the state consistency of the database.

A Merkle tree is built for each block, and it starts with leaves (the bottom of the tree), where a leaf is a transaction hash (Bitcoin uses double SHA256 hashing).
In a [Perfect Binary Merkle Tree](https://xlinux.nist.gov/dads/HTML/perfectBinaryTree.html), as shown in the [Figure 1](#pmtree), every interior node has two children and all leaves have the same depth, but not every block contains an even number of transactions.
In case there is an odd number of transactions, the hash of the last transaction is duplicated (in the [Tree](https://github.com/bitcoin/bitcoin/blob/d0f81a96d9c158a9226dc946bdd61d48c4d42959/src/consensus/merkle.cpp#L8), not in the block!) to form a [Full Binary Merkle Tree](https://xlinux.nist.gov/dads//HTML/fullBinaryTree.html), in which every node has either 0 or 2 children.
This is shown in [Figure 2](#fmtree), where the nodes `23AF` and `5101` were duplicated during the process of build the tree.

![Perfect Binary Merkle Tree][pmtree]

Moving from the bottom up, leaves are grouped in pairs, their hashes are concatenated, and a new hash is obtained from the concatenated hashes. 
The SHA256 hash is represented by the arrows in the figure.
The new hashes form new tree nodes.
This process is repeated until there’s just one node, which is called the root of the tree.
The root hash is then used as the unique representation of the transactions, is saved in block headers, and is used in the proof-of-work system.

Considering the example in [Figure 1](#pmtree).
The numbers inside the nodes represent the first 4 bytes of the hash of the transaction of that node.
Only leaf nodes store hash of real transactions, the internal nodes store the hash of its children.
The merkle path from the transaction `TX3` to the root hash `38C4` is shown by the _yellow nodes_ on [Figure 1](#pmtree).

The _blue nodes_ shows the set of the intermediate nodes (i.e, merkle proof) that can be used as proof to recreate the merkle path from the `TX3` to the root.
Thus, given `TX3`, the root hash `38C4` and the respective _blue nodes_: `D2B8`, `64B0` and `4A3B`, in this order and altogether with their respective orientations on the tree (i.e, left or right side), is possible to show that `TX3` exists in the tree by hashing it with the intermediate nodes until find the same root.
The same logic can be applied for the [Figure 2](#fmtree).

![Full Binary Merkle Tree][fmtree]

Thus, your task is to develop a Binary Merkle Tree by implementing all functions marked with `TODO(student)` in the `merkle_tree.go` file and change the function `HashTransactions` in the `block.go` to use it. 
These functions are:
 - `HashTransactions`: This function need to be changed in the `block.go` to take in consideration a merkle root instead of just the hash of all transactions.
 - `NewMerkleTree`: This function should creates a new Merkle tree from a sequence of data by using the `NewMerkleNode` function.
 - `NewMerkleNode`: This function should create a node on the merkle tree, the node can be a leaf node, which store the hash of the data or a inner node, which is a hash of its children.
 - `MerklePath`: This function should returns a list of nodes' hashes and indexes (nodes' positions: left or right) required to reconstruct the inclusion proof of a given hash.
 - `VerifyProof`: This function verifies the inclusion of a hash in the merkle tree by taking a hash and its merkle path and reconstructing the root hash of the merkle tree.

Remember to copy your implementation for the first part, but not the tests. The tests in `block_test.go`, for example, are different from the ones on the first part, since it now takes the merkle root hash in consideration.

For more information about the concept of Merkle Trees, and the [Bitcoin implementation](https://en.bitcoin.it/wiki/Protocol_documentation#Merkle_Trees) and its difference with the Ethereum implementation, please read [this](https://blog.ethereum.org/2015/11/15/merkling-in-ethereum/?source=post_page) article.

[pmtree]: perfect-merkle-tree.png "Figure 1"
[fmtree]: full-merkle-tree.png "Figure 2"

## Part 3

In our current blockchain implementation adding new blocks is easy and fast, but in real blockchain adding new blocks requires some work: one has to perform some heavy computations before getting permission to add block (this mechanism is called Proof-Of-Work).

### Proof-Of-Work

One of the key innovations of Bitcoin was to use a Proof-Of-Work algorithm to conduct a global "election" every 10 minutes (adjusted by the target difficulty), allowing the decentralized network to achieve consensus about the state of transactions.
The "elected" peer, i.e., the one that solves PoW for a block and publishes it to a majority of the peers in the network faster than the others, has granted the permission to write data on the blockchain, adding his published block to it.
And thus, consistently replicating the state in a decentralized network.

The miner that accomplish this task following the protocol receives a reward for his hard work to create a block and solve the PoW puzzle.
This is how new coins are introduced in the ecosystem.
The reward is the coinbase transaction that is added by the miner as the first transaction in the block, and there is only one coinbase transaction per block.
When a peer receives a new block, it will validate the block by checking, among other things, the block hash resulted from the PoW, the target difficult, the block timestamp, the block size, if and only if the first transaction in the block is a coinbase transaction, etc.

The miner will only "collect" its reward if his published block ends up on the longest chain because just like every other transaction, the coinbase transaction will only be accepted by other peers and included in the longest chain if a majority of miners decide to do so.
That's the key idea behind the Bitcoin incentive mechanism.
If most of the network is following the longest valid chain rule, all other peers are incentivized to continue to follow that rule.

In this lab, we will implement the Proof-Of-Work algorithm similar to the one used in Bitcoin.
Bitcoin uses [Hashcash](https://en.wikipedia.org/wiki/Hashcash), a Proof-of-Work algorithm that was initially developed to prevent email spam.
The goal of such work is to find a hash for a block, that should be computationally hard to produce.
And it's this hash that serves as a proof which should be easy for others to verify its validity.
Thus, finding a proof is the actual work.

In order to create a block, the peer that proposes that block is required to find a number, or _​nonce​_, such that when concatenating the header fields of the block with the _nonce_ and producing a hash of it, this hash output should be a number that falls into a target space that is quite small in relation to the much larger output space of that hash function.
We can define such a target space as any value falling below a certain target value, this is the _difficulty of mining_ a block.

Thus, this is a brute force algorithm: you change the _nonce_, calculate a new hash, check if it is smaller then the _target_, if not, increment the _nonce_, calculate the new hash, etc.
That's why the PoW is computationally expensive.

In the original Hashcash implementation, the target difficulty sounds like "find the hash in which the first 20 bits are zeros".
In Bitcoin, the "target bits" is added to the block header to inform the difficulty at which the block was mined, and this difficulty is adjusted from time to time, because, by design, a block must be generated every 10 minutes, despite computation power possibly increasing over the time and the number of miners in the network.

As we will not implement an adjustable mining difficulty now, we will just define the difficulty as a global constant `TARGETBITS` in the file `proof-of-work.go`.

We will set the default value of `TARGETBITS` to 8 to make the block creation fast on the tests, our goal is to have a _target_ that takes less than 256 bits in memory.
So our target will be calculated using the following formula: `2^(256-TARGETBITS)`.
The bigger the `TARGETBITS` the smaller will be the `target` number, and consequently, it's more difficult to find a proper hash.
Think of a _target_ as the upper boundary of a range: if a number (a hash) is lower than the boundary, it's valid, and vice versa.
Lowering the boundary will result in fewer valid numbers, and thus, more difficult will be the work required to find a valid one.

With the `TARGETBITS` equals to 8, you will be required to find a block hash in which the number representation is less than 2^248. Or in other words, a hash in which the first 8 bits are zeros.

For example, suppose that in the first iteration, after hash a block using the nonce value of 1, you obtained the hash `73d40a0510b6327d0fbcd4a2baf6e7a70f2de174ad2c84538a7b09320e9db3f2`, which converted to big integer representation is equals to `52390618318831801638175855856716822591931229920359547228571203746793472766962`.
As mentioned before, the default target difficulty is `2^248 == 452312848583266388373324160190187140051835877600158453279131187530910662656`.
Thus, the hash above isn't a valid PoW solution, since it is bigger than the target, i.e., `52390618318831801638175855856716822591931229920359547228571203746793472766962 > 452312848583266388373324160190187140051835877600158453279131187530910662656`. 

However, if you continue hashing the same header data just incrementing the nonce, let's say until the nonce value 59, you can eventually find a hash that is smaller than the target, like the hash `00d4eeaee903dce5468d4c6975376dfbc4c45ea1bc6c5bbbfd8e13b26aaf6e3b`, which can be represented by the big integer number `376218908933626769012171312496768664868580826658885427967934344392923377211` and is a valid solution.

Ok, so let's do the work!
Your task is to implement all functions marked with `TODO(student)` on the lab code templates.
A small description of what each function should do is given below:

- `ProofOfWork.NewProofOfWork`: Creates a new proof of work containing the given block and calculates the target difficulty based on the `TARGETBITS`. The _target_ should be a [big number](https://golang.org/pkg/math/big/).
- `ProofOfWork.setupHeader`: Prepares the header of the block by concatenating the `block.PrevBlockHash`, Merkle root hash of the `block.Transactions`, `block.Timestamp` and the `TARGETBITS` in this order.
- `ProofOfWork.addNonce`: Adds a nonce to the prepared header.
- `ProofOfWork.Run`: Performs the Proof-Of-Work algorithm. You should make a brute-force loop incrementing the nonce and hashing it with the prepared header using the [SHA256](https://golang.org/pkg/crypto/sha256/) until you find a valid hash or you reach the defined `maxNonce`. Remember that a valid PoW hash is the one that is smaller than the _target_, so you need to be able to compare them, converting the hash bytes to a big number.
- `ProofOfWork.Validate`: Validates a block's Proof-Of-Work. Currently, this function just validates if the block header hash value is less than the current target and the reconstructed header hash is equal to the mined block hash. It ignores validation errors like if the block timestamp is in the future.
  
- `Block.Mine`: Replace the function `SetHash()` for a new one called Mine that create a ProofOfWork and sets the block hash and nonce based on the result of the `Run()`.
  
- `Blockchain.ValidateBlock`: Validates a block after mining or before adding it to the blockchain (in case of receiving it from another peer). Currently, this function should perform the following validations:
  1. Check if the new block has transactions.
  2. Check if the Proof-Of-Work for that block is valid.

## Part 4

### Command line client

Until now we don't have any interface to interact with the blockchain.
Your task on this part is to create a command-line client application that will interact with your blockchain implementation.

You are free to choose any package of your preference to do this task, making a loop and reading the commands from the standard input, or using libraries like [promptui](https://github.com/manifoldco/promptui), it's up to you.

Independent of your choice, the command-line application should offer the functionalities below.
But feel free to add some more functionalities or use other names for the commands.

* create-blockchain: Creates a blockchain initializing it with the genesis block.
* add-transaction: Adds a transaction data (e.g. an input string) to a buffer but does not create a block.
* mine-block: Adds a new mined block to the blockchain committing the transactions in the buffer.
* print-chain: Prints all the blocks of the blockchain.
* print-block: Prints all transactions in a block (based on its HASH) and the block information (Hash, Nonce, Timestamp, etc).
* print-transaction: Prints the contents of a transaction based on the given ID.

Your application should be able to create transactions with the input data given by the user.
You can store the transaction in a temporary buffer in your application which can be used to create a block later (after mining it).
The application should also have a command to start mining and display the hash of the mined block in the hexadecimal format.
You will be requested to display the block and transactions information.

### Merkle tree benchmarks

Extend the merkle tree test file with two benchmarks. 
Check this post for a [tutorial](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) on writing benchmarks in Go.

* Write a benchmark function for the creation of Merkle trees. 
Benchmark the creation of a Merkle tree with 100, 1000 and 10000 data items.
Each data item should be unique, but of constant size.
* Write a benchmark function for the validation of merkle proofs.
Benchmark how long it takes to validate a Merkle proof from trees including 100, 1000 and 10000 data item.


## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours and present your solution. This lets you present the thought process behind your solution, and allows us to provide feedback on your solution then and there.
When you are ready to show your solution, reach out to a member of the teaching staff. It is expected that you can explain your code and show how it works. You may show your solution on a lab workstation or your own computer.

You should for this lab present a working demo of the application described in the previous section making a command-line client.
You should demonstrate that your implementation fulfills the previously listed specification of each assignments part.
The task will be verified by a member of the teaching staff during lab hours.
