package main

// Fixed block timestamp
const TestBlockTime int64 = 1563897484

var testTransactions = map[string]*Transaction{
	"tx0": {
		ID:   Hex2Bytes("30f2e93d7c139e7766fb80b3cb0150e0e764946bb7e4d7d54d69b53f0b1a6af1"),
		Data: []byte("tx0 data"),
	},
	"tx1": {
		ID:   Hex2Bytes("d0e97f315f21504514f5d0f4da36292dcaa7f6c01b6444d9a141df9da830f524"),
		Data: []byte("tx1 data"),
	},
	"tx2": {
		ID:   Hex2Bytes("ae9d98e7ceca0294bfe5094d30c89878d6caa290c03de1146c537e8e1a916e7d"),
		Data: []byte("tx2 data"),
	},
	"tx3": {
		ID:   Hex2Bytes("b80e1479a33da9d4aa8e838a72823fe298c261d885d8f58f8d3821e4ab1630d7"),
		Data: []byte("tx3 data"),
	},
	"tx4": {
		ID:   Hex2Bytes("3264002e4c57bb822b50891118a78c75b5c7c2c678b39c09c192e784de2cdf5e"),
		Data: []byte("tx4 data"),
	},
	"tx5": {
		ID:   Hex2Bytes("b70800696a8d3c255e957accf287730e1cc8e97b0af3066fdf845d2dc67eb212"),
		Data: []byte("tx5 data"),
	},
}

var testBlockchainData = map[string]*Block{
	"block0": {
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: nil,
		Hash:          Hex2Bytes("00b8075f4a34f54c1cf0c7f6ec9605a52161ee21e974abb4fa8a39ab7553049a"),
		Nonce:         9,
	},
	"block1": {
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx1"]},
		PrevBlockHash: Hex2Bytes("00b8075f4a34f54c1cf0c7f6ec9605a52161ee21e974abb4fa8a39ab7553049a"),
		Hash:          Hex2Bytes("00940171f20a13b9fd2cdf2c5866023c9ba876cf219951c853905bbff18af962"),
		Nonce:         711,
	},
	"block2": {
		Timestamp: TestBlockTime,
		Transactions: []*Transaction{
			testTransactions["tx3"],
			testTransactions["tx2"],
		},
		PrevBlockHash: Hex2Bytes("00940171f20a13b9fd2cdf2c5866023c9ba876cf219951c853905bbff18af962"),
		Hash:          Hex2Bytes("00532582bda0ba7fe5c313fe1175bc1f8df8df17e1291c09a098b8adf564b84c"),
		Nonce:         14,
	},
	"block3": {
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx4"]},
		PrevBlockHash: Hex2Bytes("00532582bda0ba7fe5c313fe1175bc1f8df8df17e1291c09a098b8adf564b84c"),
		Hash:          Hex2Bytes("00700132df3e9a5045bfa9d67c0a6d0a8feee6df3409424ce624cfb9a0679b43"),
		Nonce:         517,
	},
	"block4": {
		Timestamp:     TestBlockTime,
		Transactions:  []*Transaction{testTransactions["tx5"]},
		PrevBlockHash: Hex2Bytes("00700132df3e9a5045bfa9d67c0a6d0a8feee6df3409424ce624cfb9a0679b43"),
		Hash:          Hex2Bytes("001af7a5a3e3d57d81c22e6d0cf221416049372e8b4de12b41c4da1ee808000d"),
		Nonce:         24,
	},
}
