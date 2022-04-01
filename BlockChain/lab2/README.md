# Lab 2: UTXO

| Lab 2:           | UTXO                         |
| ---------------- | ---------------------------- |
| Subject:         | DAT650 Blockchain Technology |
| Deadline:        | 21. SEP                      |
| Expected effort: | 2 weeks                      |
| Grading:         | Pass/fail                    |

## Table of Contents

- [Lab 2: UTXO](#lab-2-utxo)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Part 1](#part-1)
    - [UTXO](#utxo)
    - [Transaction Outputs](#transaction-outputs)
    - [Transaction Inputs](#transaction-inputs)
    - [The egg](#the-egg)
    - [Unspent Transaction Output Set](#unspent-transaction-output-set)
  - [Part 2](#part-2)
    - [Address](#address)
      - [Public-key Cryptography](#public-key-cryptography)
      - [Digital Signatures](#digital-signatures)
      - [Base58](#base58)
      - [Updating transactions](#updating-transactions)
      - [Implementing Signatures](#implementing-signatures)
  - [Demo Application](#demo-application)
  - [Lab Approval](#lab-approval)

## Introduction

We built a very simple blockchain prototype: it's just an array of blocks, with each block having a connection to the previous one.
The actual blockchain is much more complex though.
The purpose of the blockchain is to store transactions in a secure and reliable way in a distributed environment, so no one could modify them after they are created.
Until now our transaction implementation just store all the data in a byte array and don't contain any information that could help to identify the sender or the receiver of each transactions. Now it is time to make our blockchain a bit more useful!

## Part 1

If you've ever developed an application that manage payments, you would likely have some tables in your database like: `accounts` and `transactions`.
Accounts would store information about a user, including their personal information and balance, and a transaction would store information about money transferring from one account to another.
In Bitcoin, payments are realized in completely different way, there are no accounts or balances fields, or even tables.

Since blockchain is a public and open database, sensitive information about wallet owners should not be stored in the blockchain.
Thus, in Bitcoin, coins are not collected in accounts; Transactions do not transfer money from one address to another;
There's no field or attribute that holds account balance; There are only transactions. And the transactions are stored using the UTXO model.

### UTXO

In the UTXO model, a transaction is a combination of inputs and outputs, in our implementation we will change the previous definition of a transaction in the `transaction.go` file to the struct below:

```go
type Transaction struct {
  ID   []byte
  Vin  []TXInput
  Vout []TXOutput
}
```

[Inputs](#transaction-inputs) of a new transaction reference outputs of a previous transaction (there's an exception though, which we'll discuss [later](#the-egg)).
[Outputs](#transaction-outputs) are where coins are actually stored.
An output has an implied index number based on its location in the `TXOutput` array (i.e. the array index).
The [Figure 3](#btctransactions) illustrate the interconnection between transactions inputs and outputs.

![Bitcoin transaction propagation][btctransactions]

_(image extracted from https://bitcoin.org/en/blockchain-guide#introduction)_

Is important to notice that:

1. There are outputs that are not linked to inputs (i.e. unspent outputs).
2. In one transaction, inputs can reference outputs from multiple transactions (e.g. tx6 input uses tx4 and tx5 outputs).
3. An input must reference an output (with the exception of a [coinbase transaction](#the-egg)).

In Bitcoin transactions just `lock` values with a script mechanism, which can be `unlocked` only by the one who satisfies the locking logic.

### Transaction Outputs

It is outputs that store _coins_ (notice the Value field on the struct `TXOutput` below).
And storing means locking them with a "puzzle" in the `ScriptPubKey` field.
Internally, Bitcoin uses a scripting language called [Script](https://en.bitcoin.it/wiki/Script) to define outputs locking and unlocking logic.
The language is quite primitive and is not [Turing-Complete](https://en.wikipedia.org/wiki/Turing_completeness) (this is made intentionally, to avoid possible hacks and misuses), but we will not implement it in this course. Instead we will implement a very simple locking based on matching strings and we will improve it a bit in [part2](#part-2) to use cryptographic keys.
Having such scripting language means that Bitcoin can be used as a [contract](https://en.bitcoin.it/wiki/Contract) platform as well, but with much more limited functionalities than in the [Ethereum platform](https://github.com/ethereumbook/ethereumbook/blob/develop/07smart-contracts-solidity.asciidoc#what-is-a-smart-contract).

```go
type TXOutput struct {
	Value        int
	ScriptPubKey string
}
```

The `Value` field stores the number of [Satoshis](https://en.bitcoin.it/wiki/Satoshi_(unit)) which pays the conditional pubkey script. Anyone who can satisfy the conditions of that `ScriptPubKey` can spend up to the amount of satoshis paid to it.
A satoshi is a hundred millionth of a bitcoin (0.00000001 BTC), thus this is the smallest unit of currency in Bitcoin (like a cent).
In our implementation, `ScriptPubKey` will store an arbitrary string (e.g., user defined wallet address/identifier, like a username).

One important thing about outputs is that they are _indivisible_, which means that you cannot use only part of its value.
When an output is referenced in a new transaction, **it must be spent as a whole**.
And if its value is greater than required, a **change** is generated and sent back to the sender.
This is similar to a real-world situation when you pay, let's say, $10 for something that costs $8 and get a change of $2.

### Transaction Inputs

As mentioned earlier, an input references a previous output:

- `Txid` stores the ID of the transaction containing the output used.
- `OutIdx` stores an index of the particular output to be spent.
- `ScriptSig` is the data confirming the transaction authorization. In our case it is based on matching the `ScriptPubKey` field from the referenced output.

If the data is correct, the output can be unlocked, and its value can be used to generate new outputs; if it's not correct, the output cannot be referenced in the input.
This is the mechanism that guarantees that users cannot spend coins belonging to other people.

```go
type TXInput struct {
  Txid      []byte
  OutIdx    int
  ScriptSig string
}
```

Since we don't have addresses implemented yet, `ScriptSig` will store just an arbitrary user defined string (e.g., user name).

In summary, _Outputs_ are where _coins_ are stored.
Each output comes with an unlocking **script**, which determines the logic of **unlocking** the output.
Every new transaction **must** have at least one input and output.
An input **references** an output from a previous transaction and provides data (the `ScriptSig` field) that is used in the output's unlocking script to unlock the input and use its value to create new outputs.

But what came first: inputs or outputs?

### The egg

In Bitcoin, it's the egg that came before the chicken.
The _inputs-referencing-outputs_ logic is the classical "chicken or the egg" situation: inputs produce outputs and outputs make inputs possible.
And in Bitcoin, **outputs come before inputs**.

When a miner starts mining a block it adds a _coinbase_ transaction to it.
A _coinbase_ transaction is a special type of transactions, which doesn't require previously existing outputs, and is the first transaction in each block.
It creates outputs (i.e., "coins") out of nowhere.
The egg without a chicken.
This is the **reward** miners get for mining new blocks.
 
As you know, there's the genesis block in the beginning of a blockchain.
It's this block that generates the very first output in the blockchain.
And no previous outputs are required since there are no previous transactions and no such outputs.

A _coinbase_ transaction has only one input. In our implementation its `Txid` should be empty and OutIdx equals to `-1`.
Also, a _coinbase_ transaction doesn’t store a script in ScriptSig. Instead, arbitrary data is stored there.
Take a look [here](https://www.blockchain.com/btc/tx/4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b?show_adv=true) to see the message stored on the first coinbase transaction of Bitcoin blockchain.

We define our default coinbase message for the genesis block in the file `config.go`.
In the config file we also define the constant `BlockReward`, that represents the amount of reward a miner will earn when adding a new block to the chain.
In Bitcoin, this number is not stored anywhere and is calculated based on the total number of blocks and any transaction fees paid by transactions included in this block.
Mining the Bitcoin genesis block produced 50 BTC in 2009, and every 210000 blocks the reward is halved (approximately every four years).
In our implementation, we’ll store the reward as a constant (at least for now).

### Unspent Transaction Output Set

An Unspent Transaction Output (UTXO) is and output that can be spent as an input in a new transaction.
Unspent means that these outputs weren’t referenced in any inputs.
A transaction consumes previously recorded unspent transaction outputs and creates new transaction outputs that can be consumed by a future transaction.
This way, chunks of bitcoin value move forward from owner to owner in a chain of transactions consuming and creating UTXO, as showing the [Figure 4](#utxo).

![UTXO model][utxo]

The ownership of some output is given by the ability of the user to provide an input in which the `ScriptSig` that can unlock the referred output's `ScriptPubKey`, and consequently an authorization to spend it.
The lock mechanism will be implemented by the methods `CanBeUnlockedWith` and `CanUnlockOutputWith` defined in the `transaction_output.go` and `transaction_input.go` respectively.
For now, these functions should just compare the script fields with `unlockingData`, which is a string.
These functions will be improved in the next part of this lab, after we implement addresses based on private keys.

But how we can find all the transactions containing unspent outputs?

The simple solution is just check every block in the blockchain.
Starting by searching and collecting for outputs that can be unlocked by the address we’re searching unspent transaction outputs for.
But before taking those outputs, we need to check if it was previously referenced in an input (e.g., was already spent).

If an output were referenced in some input (i.e., their values were moved to other outputs), thus we cannot count them and we should ignore it.
Thus, after checking the desired outputs inside the transaction, we should gather all inputs that could unlock outputs locked with the provided address, since this will be used to check if some output was already spent.
But remember that this doesn’t apply to coinbase transactions, since they don’t unlock outputs.
Then this process continues for every transaction in the block, for all blocks, resulting in a list of transactions containing only unspent outputs (i.e. UTXO).

But as you probably imagined, since transactions are stored in blocks, iterate over each block in the blockchain database and checks every transaction in it can be very inefficient and time-consuming.
Currently the whole blockchain database takes around 230 Gb of disk space, with more than 592 blocks.

The solution to the problem is to have an index that stores only unspent outputs, which is exactly what Bitcoin full nodes do, by tracking all available and spendable outputs.
The Bitcoin blockchain stores data in different LevelDB [databases](https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage#Overview). The blocks are stored in the `blocks` database, and instead of store all the transactions, it stores what is called the `UTXO set` in the `chainstate` database to easily retrieve unspent transactions outputs.

The `UTXO set` is the collection of all unspent transaction outputs (UTXO) and it is used as a cache that is built from all blockchain transactions, and later used to calculate balance and validate new transactions.
The UTXO set grows as new UTXO is created and shrinks when UTXO is consumed.
Every transaction represents a change (state transition) in the UTXO set.

In this lab, as we are not implementing persistence, we will represent the UTXO Set as a map, where the keys are the transaction ID and the content is another map of the indexes of the outputs into the transaction to their proper Outputs (i.e. TXOutput).

```go
type UTXOSet map[string]map[int]TXOutput
```

But to check the balance of some user, we don’t need all unspent transactions output (i.e., UTXO set), but only those that can be unlocked by the known key.
But we can use the UTXO set to search for those UTXO, which is faster than search on the whole blockchain.
When the user's wallet application "receive" bitcoin, this actually means that the wallet has detected an UTXO that can be spent with one of the keys controlled by that wallet (i.e., unlocked by the user).
The account balance is the sum of values of all UTXO locked by the account address (i.e., that the user's wallet can spend).

Your task is implement all functions marked with `TODO(student)` on the lab code templates.
A small description of what each function should do is given below:

- `Block.HashTransactions`: Modify this method to use the result of the `Transaction.Serialize` method instead of just the `Data` field used in the previous lab, since now your transaction is more than just a byte slice.
- `Transaction.NewCoinbaseTX`: creates a new _coinbase_ transaction.
- `Transaction.NewUTXOTransaction`: creates a new transaction by getting a list of spendable outputs to be used and the current balance of the sender. It also checks if the sender sufficient funds to perform the transaction.
- `Transaction.IsCoinbase`: test if a transaction is _coinbase_.
- `Transaction.Hash`: hash the serialized copy of a transaction ignoring the ID field, since the return of this function is the new ID of a transaction. **Before** serialize the transaction, make sure to ignore any existent data in the ID field (i.e., set ID in the copy to []byte{}).
- `Transaction.Serialize`: encodes the transaction struct using the [gob](https://golang.org/pkg/encoding/gob/) library. This method should not ignore the ID of the transaction, the whole transaction fields should be serialized.
- `Blockchain.FindUTXOSet`: finds all unspent outputs by iterating over all blocks. Returns an UTXO set.
- `Blockchain.FindTransaction`: finds a transaction in the blockchain by its ID. It iterates over all blocks until finds it.
- `Blockchain.VerifyTransaction`: verifies if the referred inputs of a given transaction exist in the blockchain, discard invalid transactions that make reference to unknown inputs.
- `Blockchain.MineBlock`: creates a new block on the blockchain. This function should verify the existence of transactions inputs and mine a new block if there is a list of valid transactions. Remember that if a transaction is coinbase it is already valid. And that a block should contain only **one** coinbase transaction and it **must be the first** transaction in the list.

- `UTXOSet.FindSpendableOutputs`: this function is used when a new transaction is created, to finds the enough number of outputs holding the given amount to be spent. It uses the UTXO set instead of search in the blockchain.
- `UTXOSet.FindUTXO`: finds unspent outputs for a given address. Can be used to get the balance. It uses the UTXO set instead of search in the blockchain.
- `UTXOSet.Update`: updates the current UTXO set with a new set of transactions. Should ignore transactions that try to double-spend (i.e., use the same output twice).

[btctransactions]: btc-transactions.svg "Figure 3"
[utxo]: utxo.png "Figure 4"

## Part 2

As in Bitcoin there are no user accounts and your personal data (e.g., name, passport number) is not required and not stored anywhere.
But there still must be something that identifies you as the owner of transaction outputs (i.e. the owner of coins locked on these outputs).
And this is what Bitcoin addresses are needed for.

So far we've used arbitrary user defined strings as addresses, and the time has come to implement real addresses, as they're implemented in Bitcoin

###  Address

A [Bitcoin address](https://en.bitcoin.it/wiki/Address) is a string of digits and characters that can be shared with anyone who wants to send you money.
They are public and the most famous address in Bitcoin is `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`, what is the very first address, which allegedly belongs to Satoshi Nakamoto, the Bitcoin creator.

But addresses (despite being unique) are not something that identifies you as the owner of a _wallet_.
In fact, such addresses are a human readable representation of public keys, and consist of a string of numbers and letters, beginning with the digit "1".
In Bitcoin, your identity is a pair (or pairs) of private and public keys stored on your computer (or stored in some other place you have access to).
Bitcoin relies on a combination of cryptography algorithms to create these keys, and guarantee that no one else in the world can access your coins without getting physical access to your keys.

#### Public-key Cryptography

[Public-key cryptography](https://en.wikipedia.org/wiki/Public-key_cryptography) algorithms use pairs of keys: public keys and private keys.
Public keys are not sensitive and can be disclosed to anyone.
In contrast, private keys shouldn't be disclosed: no one but the owner should have access to them because it's private keys that serve as the identifier of the owner.
The private keys of an user are his identities in the world of cryptocurrencies.

In essence, a Bitcoin wallet is just a pair of such keys.
When you install a wallet application or use a Bitcoin client to generate a new address, a pair of keys is generated for you.
The one who controls the private key controls all the coins sent to this key in Bitcoin.

Private and public keys are just random sequences of bytes, thus they cannot be printed on the screen and read by a human.
That's why Bitcoin uses the [Base58](https://en.bitcoin.it/wiki/Base58Check_encoding) algorithm to convert public keys into a human readable string.
We will see how to implement it soon, but first we will see how does Bitcoin check the ownership of coins.

#### Digital Signatures

Bitcoin uses [digital signatures](https://en.wikipedia.org/wiki/Digital_signature) for verifying the authenticity of the transactions, making using of asymmetric cryptography to guarantee:

* that data wasn't modified while being transferred from a sender to a recipient;
* that data was created by a certain sender;
* that the sender cannot deny sending the data.

By applying a signing algorithm to data (i.e., signing the data), one gets a signature, which can later be verified.
Digital signing happens with the usage of a private key, and verification requires a public key.

In order to sign data we need the following things:

1. data to sign;
2. private key.

The operation of signing produces a signature, which is stored in transaction inputs.
In order to verify a signature, the following is required:

1. data that was signed;
2. the signature;
3. public key.

In simple terms, the verification process can be described as: check that this signature was obtained from this data with a private key used to generate the public key.

Note that digital signatures are not encryption, you cannot reconstruct the data from a signature.
This is similar to hashing: you run data through a hashing algorithm and get a unique representation of the data.
The difference between signatures and hashes is key pairs: they make signature verification possible.
But key pairs can also be used to encrypt data: a private key is used to encrypt, and a public key is used to decrypt the data.
Bitcoin doesn't use encryption algorithms though, only your wallet that can potentially encrypt your keys.

Every transaction input in Bitcoin is signed by the one who created the transaction.
Every transaction in Bitcoin must be verified before being put in a block.
Verification means (besides other procedures):

1. Checking that inputs have permission to use outputs from previous transactions.
2. Checking that the transaction signature is correct.

More information about digital signatures and how Bitcoin use it can be found [here](https://github.com/bitcoinbook/bitcoinbook/blob/develop/ch04.asciidoc#introduction).

In this lab we will not implement a wallet, but instead we will define some helper functions in `wallet.go` to manage keys and addresses.
The `newKeyPair` function is responsible to generate the key pair using the same algorithm used by Bitcoin to digitally sign transactions, the _ECDSA_ (Elliptic Curve Digital Signature Algorithm).
If you want to know more about elliptic curves and their use, you can check out this [tutorial](https://andrea.corbellini.name/2015/05/17/elliptic-curve-cryptography-a-gentle-introduction/).

The `newKeyPair` function will use the go standard package [elliptic](https://golang.org/pkg/crypto/elliptic/#P256) to create a P-256 elliptic curve, and for the purpose of this lab, you just need to know that this and other types of curves can be used to generate really big and random numbers on elliptic curves cryptosystems, and that the private and public keys are based on that numbers. 

The private key will be generated using the generated curve by calling the function `ecdsa.GenerateKey` from the go standard package [ecdsa](https://golang.org/pkg/crypto/ecdsa/#GenerateKey), and the public key will be generated from the private key.
In [elliptic curve based algorithms](https://github.com/bitcoinbook/bitcoinbook/blob/develop/ch04.asciidoc#elliptic-curve-cryptography-explained), public keys are points on a curve, thus, a public key is a combination of X, Y coordinates.
In Bitcoin, these coordinates are concatenated and form a public key.
So we will do the same in our implementation to generate the public key, concatenating the X and Y coordinates of the `ecdsa.PublicKey` field in the private key.

But we said previously that addresses are based on public keys, but how the addresses are generated?

#### Base58

We know so far that the address `18ZrMdiBrvtKhV9TqH4YXmBGDCwqxCBLwT` is a human-readable representation of a public key.
And if we decode it, here's what the public key looks like (as a sequence of bytes written in the hexadecimal system):

```
0052ff643039009a75c98ea449a3155b6629757c03fa7f72a2
```

Bitcoin uses the [Base58](https://en.bitcoin.it/wiki/Base58Check_encoding) algorithm to convert public keys into human readable format.
The algorithm is very similar to famous [Base64](https://en.wikipedia.org/wiki/Base64), but it uses shorter alphabet: some letters were removed from the alphabet to avoid some attacks that use letters similarity.
Thus, there are no these symbols: 0 (zero), O (capital o), I (capital i), l (lowercase L), because they look similar. Also, there are no + and / symbols.

The process of getting an address from a public key is done by the function `GetAddress` in `wallet.go` that follows the algorithm described [here](https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address), and is illustrated in the [Figure 5](#btcaddress).

![Bitcoin Address Generation][btcaddress]

Thus, the above mentioned decoded public key consists of three parts:
```
Version  Public key hash                           Checksum
00       52ff643039009a75c98ea449a3155b6629757c03  fa7f72a2
```

Since hashing functions are one way (i.e., they cannot be reversed), it's not possible to extract the public key from the hash.
But we can check if a public key was used to get the hash by running it thought the save hash functions and comparing the hashes.

As shown in [Figure 5](#btcaddress) the steps to convert a public key into a Base58 address are:

1. Take the public key (i.e., resulted of the concatenation of X and Y coordinates) and hash it twice, applying the [RIPEMD160](https://en.wikipedia.org/wiki/RIPEMD) hashing algorithm to the result of _SHA256_ of the public key.
   This will be done by calling the function `HashPubKey` shown in blue in the figure.
2. Prepend the version of the address generation algorithm to the hash (i.e., version + PubKeyHash combination).
3. Calculate the checksum by hashing the result of `step 2` with a double _SHA256_: `SHA256(SHA256(versionedPayload))`.
   From the resulting 32-byte hash, we take only the _first four bytes_.
   These four bytes serve as the error-checking code, or checksum, shown in purple in the figure.
   This operation will be done by the function `checksum` in the `wallet.go` file.
4. Append the checksum to the versioned payload (i.e., version + PubKeyHash + checksum combination).
5. Encode the `version + PubKeyHash + checksum` combination resulted of `step 4` with Base58.

As a result, you'll get a real Bitcoin address, you can even check its balance on [blockchain.info](https://www.blockchain.com/btc/address/18ZrMdiBrvtKhV9TqH4YXmBGDCwqxCBLwT).
But I can assure you that the balance will be 0 no matter how many times you generate a new address and check its balance.
This is why choosing proper public key cryptography algorithm is so crucial: considering private keys are random numbers, the chance of generating the same number must be as low as possible.
Ideally, it must be as low as "never".

Also, pay attention that you don't need to connect to a Bitcoin node to get an address.
The address generation algorithm utilizes a combination of open algorithms that are implemented in many programming languages and libraries.
More information about the address generation can be found [here](https://github.com/bitcoinbook/bitcoinbook/blob/develop/ch04.asciidoc#bitcoin-addresses).

Now it's your turn, go to the `wallet.go` and implement all the functions to generate keys and addresses marked with `TODO(student)` on the code template. These functions are:

- `newKeyPair`: generates a new key pair using the _P-256 curve_ and _ECDSA_ algorithm.
- `GetAddress`: computes the address based on the given public key (i.e., algorithm previously described)
- `HashPubKey`: computes the hash (RIPEMD160 + SHA256) of the public key (i.e., step 1).
- `checksum`: computes the checksum of a given versioned payload, applying a double SHA256 hash algorithm to it (i.e., step 3).
- `pubKeyToByte`: convert the generated `PublicKey` to a byte array concatenating the X and Y coordinates.
- `GetPubKeyHashFromAddress`: returns the hash of the public key ignoring the version and the checksum.
- `ValidateAddress`: checks if an address is valid by decoding the given address, extracting the current version, public key hash and the checksum, and re-computing the checksum. The address is valid only if the checksums match.

The function `ValidateAddress` can be used to prevent a mistyped bitcoin address from being accepted by the wallet software as a valid destination when transfer some coin, an error that would otherwise result in loss of funds. You can also use some [online bitcoin addresses validators](http://lenschulwitz.com/base58) to compare your implementation.

Now that you have implemented real addresses, let's modify the transaction inputs and outputs to use it.

#### Updating transactions

We will no longer use `ScriptPubKey` and `ScriptSig` fields, because we will not going to implement a scripting language for now.
Instead, we will implement the same outputs locking/unlocking and inputs signing logics as in Bitcoin, but we’ll do this in methods instead.
Thus the `ScriptSig` field of the `TXInput` struct will be split into `Signature` and `PubKey` fields, and `ScriptPubKey` field of the `TXOutput` struct will be renamed to `PubKeyHash`, as shown below.

```go
type TXInput struct {
	Txid      []byte
	OutIdx    int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}
```

We will replace the method `CanUnlockOutputWith` in the `transaction_input.go` by the `UsesKey` method, that will checks that an input uses a specific key to unlock an output.
Notice that inputs store raw public keys (i.e., not hashed), but the function takes a hashed one as a parameter and compares it with the hashed version of the `PubKey` on the input.
Before, we were just comparing strings.

We will do a similar replacement for the outputs, removing the function `CanBeUnlockedWith` and replacing by `IsLockedWithKey`, that will checks if the provided public key hash was used to lock the output.
This is a complementary function to `UsesKey`, and they're both used to build connections between transactions.

We will also add a new function to the `transaction_output.go` named `Lock`, which simply locks an output.
When we send coins to someone, we know only their address, thus the function takes an address as the only argument.
The address is then decoded and the public key hash is extracted from it and saved in the `PubKeyHash` field.
This mechanism "lock" the transaction output to a specific address, giving the ownership of the output to the specific address.

#### Implementing Signatures

Transactions must be signed because this is the only way in Bitcoin to guarantee that one cannot spend coins belonging to someone else.
If a signature is invalid, the transaction is considered invalid too and, thus, cannot be added to the blockchain.

We have all the pieces to implement transactions signing, except one thing: data to sign.
What parts of a transaction are actually signed? Or a transaction is signed as a whole? 
Choosing data to sign is quite important.

The thing is that data to be signed must contain information that identifies the data in a unique way.
For example, it makes no sense signing just output values because such signature won't consider the sender and the recipient.

Considering that transactions unlock previous outputs, redistribute their values, and lock new outputs, the following data must be signed:

* Public key hashes stored in unlocked outputs. This identifies "sender" of a transaction.
* Public key hashes stored in new, locked, outputs. This identifies "recipient" of a transaction.
* Values of new outputs.

In Bitcoin, locking/unlocking logic is stored in _scripts_, which are stored in _ScriptSig_ and _ScriptPubKey_ fields of inputs and outputs, respectively. Since Bitcoins allows different types of such scripts, **it signs the whole content of _ScriptPubKey_**.

As you can see, we don't need to sign the public keys stored in inputs.
Because of this, in Bitcoin, it's not a transaction that's signed, _but its trimmed copy with inputs storing ScriptPubKey from referenced outputs_.
A detailed process of getting a trimmed transaction copy is described [here](https://en.bitcoin.it/w/images/en/7/70/Bitcoin_OpCheckSig_InDetail.png).

We will implement the `Sign` function in the `transaction.go` file.
The function will take a private key and a map of previous transactions.
As mentioned above, in order to sign a transaction, we need to access the _outputs referenced in the inputs of the transaction_, thus we need the transactions that store these outputs.
But remember that _coinbase_ transactions are not signed because there are no real inputs in them.

So the first thing is to check if the transaction is coinbase, and if not, if it has valid inputs, which means, inputs that reference existent transactions, otherwise an error should be generated.
Then we will make a _trimmed copy_ of the transaction to be signed, as we will not sign transaction containing previous signatures.
The copy will include all the inputs and outputs of the original transaction, but the `TXInput.Signature` and `TXInput.PubKey` will be set to `nil` in the copy.

Next, we need to iterate over each input in the *copy*, and for each input, ensure that the `Signature` is `nil` and set the `PubKey` of the input with the `PubKeyHash` of the referenced output that is in the `prevTXs`. This is the data that will be signed, i.e., that will be serialized and used as input to the [ecdsa.Sign](https://golang.org/pkg/crypto/ecdsa/#PrivateKey.Sign) function from the go `crypto/ecdsa` library. The resulted _ECDSA_ signature is a pair of numbers (`r` and `s`), which will be concatenated as bytes and stored in the input's `Signature` field for each input in the `txCopy.Vin`.

Thus each input in the copy will be signed separately, since the all the `Signature` and `PubKey` field of the others inputs in the transaction are `nil` on each iteration, and only the one that is being signed have the `txCopy.Vin[in.txID].PubKey` set with the correspondent `PubKeyHash` of the output that is being spent. Note that to not affect further iterations, is a good idea to reset the `PubKey` field to `nil` in the end of each iteration for each input.

In the same file `transaction.go`, we will also create a function named `Verify` to check the signatures of transaction inputs from a list of transactions.
As before, there is no need to verify the signatures of _coinbase_ transaction, since they don't have one, so in this case, just return `true`.
If some of the inputs in the given list of transactions are invalid inputs, the transaction should be ignored, and an error generated.

Thus, we will first get a trimmed copy of the same transaction.
Next, we'll need to create the same curve that is used to generate key pairs (i.e., `elliptic.P256`).
And check the signature in each input, by iterating over the list of inputs in the trimmed copy (i.e., `txCopy.Vin`).
This piece is identical to the one in the `Sign` method, because during verification we need the same data what was signed, which means that we need to set the `Signature` field of the input in the copied transaction to `nil` and also `PubKey` to the `PubKeyHash` of the referenced output.

After that, we need to unpack the values of `TXInput.Signature` and `TXInput.PubKey` stored in the original transaction (*not the copied one*), since a signature is a pair of numbers and a public key is a pair of coordinates.
And we concatenated them earlier for storing, and now we need to unpack them to use it as parameters to the [ecdsa.Verify](https://golang.org/pkg/crypto/ecdsa/#Verify) function.

Thus, to extract `r` and `s` from the input `Signature`, as they are "big numbers" we will need to import the go [math/big](https://golang.org/pkg/math/big/#pkg-examples) package, and create two new Big Integers and assign to each one half of the Signature.
Remember that we concatenated them in the `Sign` function, so now we need to split it in the same order.
The same logic will be applied for the coordinates `X` and `Y` of the input `PubKey`.
Thus, after performing all extractions of the signature and public key from the input, we pass them to the `ecdsa.Verify` function and check if all inputs are verified, if so then return true, if at least one input fails verification, return false.

[btcaddress]: btc-address.png "Figure 5"

## Demo Application

Extend your client application of lab 1 to create two addresses and make a demo transaction between them. Show the correspondent content of the blocks where the transactions were added and the whole blockchain state in the end of the process. Need to be possible to see the content of the transactions. Create a command to get and show the current balance of both addresses based on the UTXO of each one.

Thus, create a `getBalance` in your application. This function should receive an address and the current UTXO set and return the balance of the address.

You should be able to demonstrate the follow test scenario:
1) create 3 addresses (`a`,`b`,`c`)
2) define `a` as the block miner and create the genesis with this address to fund it (it can be a hardcoded account if you prefer)
3) transfer rewarded (i.e. 10) coins of `a` to `b` address, now both addresses should have 10 coins since `a` mined the block for this transaction, and `c` should have 0
4) attempt to transfer coins using `c`, should fail and report that `c` has no funds
5) transfer 5 coins from `b` to `c`
6) the final result should be: 3 blocks (including the genesis), with balances:`a = 20`, `b = 5`, `c = 5`

Note: you may need to create some auxiliary functions in your client application to create addresses or perform any necessary setup.

## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours and present your solution. This lets you present the thought process behind your solution, and allows us to provide feedback on your solution then and there.
When you are ready to show your solution, reach out to a member of the teaching staff. It is expected that you can explain your code and show how it works. You may show your solution on a lab workstation or your own computer.

You should for this lab present a working demo of the application described in the previous section making a transaction between two addresses. 
You should demonstrate that your implementation fulfills the previously listed specification of each assignments part.
The task will be verified by a member of the teaching staff during lab hours.
