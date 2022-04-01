from TripleSDES import TripleSDES

key1 = 0b1000101110
key2 = 0b0110101110

tdes = TripleSDES()

def encrypt(cipherString):
    res = []
    for data in cipherString:
        cp = int(data,2)
        pt = tdes.tripleEncrypt(key1,key2,cp)
        res.append(chr(pt))
    return res

def string2bits(s=''):
    return [bin(ord(x))[2:].zfill(8) for x in s]


inpStr = input("Enter string : ")

bistr = string2bits(inpStr)
res = encrypt(bistr)

binres = string2bits(res)
print("".join(i for i in binres))