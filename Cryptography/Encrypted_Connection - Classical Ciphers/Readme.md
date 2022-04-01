# Encryption

This file will contain the following information.

- Installation of required packages to run the code
- File Information

## Installation
The package include a **requirements.txt** file. Run the file using following command
```sh
pip install -r requirements.txt
```

This should install all the dependencies and requirements.

## File Information: Part 1
The directory contains the following files:
- `DHKeyExchange.py` : This file has the logic for selecting the cyclic group and creating keys for DH Key exchange.
- `TestAliceBob.py` : This contains test programs for Alice and Bob. Alice will only encrypt and Bob will decrypt the encrypted data
- `TripleSDES.py` : Simple Triple DES implementation

To run the Alice and Bob file:
```
python TestAliceBob.py <Alice/Bob>
```

## File Information: Part2

The directory contains following files:
- `DHKeyExchange.py` : This file has the logic for selecting the cyclic group and creating keys for DH Key exchange.
- `Server.py` : Contains logic of connecting Alice and Bob via internet over http.
- `TripleSDES.py` : Simple Triple DES implementation

To Run the server and be able to serve the encryption & decryption over http, run 
```
python Server.py <Alice/Bob> <host:port> <targetHost:port> [number1-number2]

e.g. python Server.py Alice localhost:5600 localhost:5601 5000-8000
python Server.py Bob  localhost:5601 localhost:5600 5000-8000
```
where, <br />
`Alice/Bob` : Run the program as Alice or Bob <br />
`host:port` : hostname port string separated by colon <br />
`targetHost:port` : destination hostname port string <br />
`number1-number2` : Optional, range of number to select P,Q and cyclic group <br />


## Endpoints for Part2
`/getpub` : will return public key for user
`/sendMessage` : send message. POST supported

