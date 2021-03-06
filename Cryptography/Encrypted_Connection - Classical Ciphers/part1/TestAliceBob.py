from DHKeyExchange import DHKeyExchange
from TripleSDES import TripleSDES
import sys

dataLen = 8

def deciper(cipherString):
    dataList = [cipherString[i:i+dataLen] for i in range(0, len(cipherString), dataLen)]
    finalStr = ""
    for data in dataList:
        cp = int(data,2)
        pt = tdes.decrypt(bkey,cp)
        finalStr += chr(pt)
    return finalStr

def encrypt(cipherString):
    res = []
    for data in cipherString:
        cp = int(data,2)
        pt = tdes.encrypt(bkey,cp)
        res.append(chr(pt))
    return res

def string2bits(s=''):
    return [bin(ord(x))[2:].zfill(8) for x in s]



# start the program
if len(sys.argv) < 2:
    print("Usage : python ",__file__, " <Alice/Bob>")
    exit(1)
    
from TripleSDES import TripleSDES
name = sys.argv[1].upper()

print("User : ",name)
if name == "ALICE":
    src = "Alice"
    dest = "Bob"
else:
    src = "Bob"
    dest = "Alice"

dhk = DHKeyExchange()
pubKey,priKey = dhk.genPubPriKeys()
print("{}'s : Public key : ".format(src),pubKey, " Private key : ",priKey)

pubKeyDest = input("Enter {}'s public key : ".format(dest))

pubKeyDest = int(pubKeyDest)

sharedKey = dhk.genSharedKey(pubKeyDest,priKey)

secretKey = dhk.genSecretKey(sharedKey,10)

print(src," : Shared key : ", sharedKey)

bkey = int(secretKey,2)

print(src," : secret key generated by BBS generator : ", bkey)

tdes = TripleSDES()

while(True):
    if name =="ALICE":
        text = input('Message to Bob: ')
        bistr = string2bits(text)
        res = encrypt(bistr)
        binres = string2bits(res)
        print("".join(i for i in binres))
    else:
        text = input('Message to decrypt from Alice: ')
        # bistr = string2bits(text)
        res = deciper(text)
        print(res)