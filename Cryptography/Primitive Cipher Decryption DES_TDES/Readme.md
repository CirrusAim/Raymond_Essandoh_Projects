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

## File Information: Poly-aplhabetic Ciphers
The directory contains the following files:
- `AutoKeyCracker.py` : This file the logic to find the key and decipher the ciphertext.
- `BigramScores.txt` : This contains the bigram scores of english language.

To crack the cipher run
```
python AutoKeyCracker.py
```

To get the  list of commmand line arguments supported, run
```
python AutoKeyCracker.py --help
```

## File Information: SDES and Triple SDES

The directory contains following files:
- `TripleSDES.py` : This has all the logic related to SDES and Triple SDES encryption and decryption. It also has a main function to run the basic sanity.
- `Webserver.py` : This contains the logic to create a simple websever using Flask and accept the POST request to decrypt the cipher and return a decrypted response.
    - `cipher.py`   : This file when run prompts the user to input any length of english text and then encrypt it to its corresponding binary values.
- `Task_SDES.py`: This generate the ciphertexts and plaintexts of task1 and task2 tables.
- `Task_SDES.py`: This generate the ciphertexts and plaintexts of task1 and task2 tables.
- `Decipher.py` : This contains the functions to crack the SDES and Triple SDES. It has two functions 
    -  `decryptSDES` :  This reads from `ctx1.txt` file, perform brute force SDES on the input and write the probable decrypted plaintext to `decrypted_ctx1.txt`.
    -  `decryptTripleSDES`:This reads from `ctx2.txt` file, perform brute force triple SDES on the input and write the probable decrypted plaintext to `decrypted_ctx2.txt`.

To Run the server and be able to serve the encryption & decryption over http, run 
```
python Webserver.py
```

To crack the ciphers, run
```
python Decipher.py
```

To generate ciphertext and plaintext of task1 and task2, run
```
python Task_SDES.py
```
