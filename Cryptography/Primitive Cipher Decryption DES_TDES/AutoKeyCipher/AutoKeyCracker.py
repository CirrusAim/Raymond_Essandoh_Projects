import os.path
from time import time
import argparse

alphabet="ABCDEFGHIJKLMNOPQRSTUVWXYZ"
M = len(alphabet)

bigramScores = open('BigramScores.txt', 'r')

bigrams = {}

def bigramScore(text):
    score = 0
    for i in range(len(text)-1):
        score += bigrams[text[i:i+2]]
    return score

def alphaToNumber(L):
    """ Convert letters to numbers according to APLHABET """
    return [alphabet.index(i) for i in L]


def numberToAlpha(L):
    """ Convert numbers to letters according to APLHABET """
    return [alphabet[i] for i in L]



def crackCipher(ctext,limit=10):
    """ Since we have no way of determining the length of the key we will have to test many different
    key lengths. By default we check every key length from 2 to 10.

    To find the best key we first assume it is the letter A repeated many times
    which represents no key at all. Then we change the first letter of the key
    and test if it looks better. After going through all 26 posibilities for the
    first letter we do the same for the second. Then the third and so on. 
    """
    # Our starting score
    outKey = ["A"]
    outScore = bigramScore(ctext)
    
    # Try every key length
    for klen in range(2,limit):
        
        # Start with the simplest key and its score
        bestKey = ["A"]*klen
        # AA
        bestScore = bigramScore(ctext)
    
        # For each position
        for i in range(klen):
            # Try every possible letter in that position
            for l in alphabet:
                tempKey = bestKey[:]
                tempKey[i] = l
                dtext = decrypt(ctext,tempKey)
                score = bigramScore(dtext)
                # If this is the best so far save it
                if score > bestScore:
                    bestKey = tempKey
                    bestScore = score
        
        # If this key length produced a better decode than any previous key
        # length we make it the one we will save.
        if bestScore > outScore:
            outScore = bestScore
            outKey = bestKey
    
    outKey = "".join(outKey)
    print("Best Key Found: {}".format(outKey))
    return decrypt(ctext,outKey)

def decrypt(text,key):
    """ Decode the text using the given key """    
    validptext(text)
        
    """ Convert the text to numbers """
    T = alphaToNumber(text)

    """ Conver the key to numbers """
    K = alphaToNumber(key) 

    out = []
    for keynum,textnum in zip(K,T):
        # Decode a letter then add it to the keystream
        out.append( (textnum-keynum) % M )
        K.append( out[-1] )

    return "".join(numberToAlpha(out))

# Ciphertext must each character from the ALPHABET.
def validptext(T):
    if type(T) != str:
        raise Exception("Plaintext must be a string")
    
    for i in T:
        if i not in alphabet:
            raise Exception("{} is not a valid plaintext character".format(i))


if __name__ == "__main__":

    keyLen = 10
    ciphertext = ""
    plaintext = ""
    bestKey = ""

    my_parser = argparse.ArgumentParser(prog='AutoKeyCracker', description="crackCipher the ciphertext using bigrams")
    my_parser.add_argument('-c', '--cipher', metavar='ciphertext', action='store', nargs="*", help="cipher text")
    my_parser.add_argument('-f', '--file', metavar='file', action='store', nargs="*", help="input cipher file")
    my_parser.add_argument('-k', '--key', metavar='key', action='store', help="Key to decode the cipher")
    my_parser.add_argument('-l', '--keylen', metavar='length', action='store', help="Max length of the key. Default = 10")
    parser =  my_parser.parse_args()

    if not parser.file and not parser.cipher:
        print("Either cipher file or cipher text has to be given to decode.\n\n")
        my_parser.print_help()
        exit()

    if parser.file and parser.cipher:
        print("Either of cipher file or cipher text can be given not both.\n\n")
        my_parser.print_help()
        exit()
    
    if parser.cipher:
        ciphertext = parser.cipher[0]
    elif parser.file:
        inpFile = os.path.abspath(parser.file[0])
        with open(inpFile,"r") as f:
            ciphertext = f.read()
        ciphertext = "".join(ciphertext.split())

    
    if not ciphertext :
        print("Empty cipher text is not valid")
        exit()
    
    print("\nThe ciphertext is:",ciphertext,"\n\n")
    
    if parser.key:
        plaintext = decrypt(ciphertext,parser.key.upper())
        print(plaintext.lower())
        exit()
    
    if parser.keylen:
        keyLen = int(parser.keylen)
    
    t1 = time()
    for line in bigramScores:
        L = line.split(" ")
        bigrams[L[0]] = int(L[1])

    
    plaintext = crackCipher(ciphertext,keyLen).lower()
    print(plaintext)
    t2 = time()
    print("Elapsed time for key len {} is {} secs".format(keyLen,t2-t1))
    