
from TripleSDES import TripleSDES
from time import time

allKeys = []
keyLen = 10

def generateAllPossibleKeys(keyLen, arr, i):
    """ Generate all binary strings of length equal to keyLen. 
        For example, for keyLen = 2, it will generate array = ['00', '01','10','11']
    """
    if i == keyLen:   
        key = ""
        for i in range(0,keyLen):
            key += str(arr[i])
        allKeys.append(key)
        return
    
    arr[i] = 0
    generateAllPossibleKeys(keyLen, arr, i + 1)
    arr[i] = 1
    generateAllPossibleKeys(keyLen, arr, i + 1)



def decryptTripleSDES(allKeys):
    """ It will read the encrypted file ctx2.txt and decrypt it using brute force Triple SDES. It will write the probable candidates to decrypted_ctx2.txt file"""
    t1 = time()
    tdes = TripleSDES()
    cipherString = ""
    inputFile = "ctx2.txt"
    print("Started decrypting cipher using Triple SDES from file : {}".format(inputFile))
    with open(inputFile,"r") as f:
        cipherString = f.read()
    
    keyDeciperMap = {}
    dataList = [cipherString[i:i+8] for i in range(0, len(cipherString), 8)]
    
    onePercent = int(len(allKeys)/100)
    percentCompleted = 0
    count = 0
    for key1 in allKeys:
        innerDecipherMap = {}
        binKey1 = int(key1,2)
        for key in allKeys:
            finalStr = ""
            binKey = int(key,2)
            for data in dataList:
                cp = int(data,2)
                pt = tdes.tripleDecrypt(binKey1,binKey,cp)
                finalStr += chr(pt)

            if finalStr.isascii() and finalStr.isprintable():
                innerDecipherMap[key] = finalStr
        if innerDecipherMap:
            keyDeciperMap[key1] = innerDecipherMap
        count += 1
        if count == onePercent:
            percentCompleted += 1
            print(str(percentCompleted) + "% completed",end='\r')
            count = 0
        
    outFile = "decrypted_" + inputFile
    with open(outFile,"w") as f :
        for key,val in keyDeciperMap.items():
            for k,v in val.items():
                f.write(key + "\t:\t"+ k + "\t:\t'"+ v + "'")

    print("Finished decrypting cipher using Triple SDES to file : {}".format(outFile))
    t2 = time()
    timeTaken = t2-t1/60
    print("Time taken to decrypt the cipher using TripleSDES : {:0.3f} minutes".format(timeTaken))




def decryptSDES(allKeys):
    """ It will read the encrypted file ctx1.txt and decrypt it using brute force SDES. It will write the probable candidates to decrypted_ctx2.txt file"""
    
    tdes = TripleSDES()
    cipherString = ""
    inputFile = "ctx1.txt"
    print("Started decrypting cipher using SDES from file : {}".format(inputFile))
    with open(inputFile,"r") as f:
        cipherString = f.read()
    
    keyDeciperMap = {}
    dataList = [cipherString[i:i+8] for i in range(0, len(cipherString), 8)]
    for key in allKeys:
        finalStr = ""
        binKey = int(key,2)
        for data in dataList:
            cp = int(data,2)
            pt = tdes.decrypt(binKey,cp)
            finalStr += chr(pt)
        
        if finalStr.isascii() and finalStr.isprintable():
            keyDeciperMap[key] = finalStr

    outFile = "decrypted_" + inputFile
    with open(outFile,"w") as f :
        for key,val in keyDeciperMap.items():
            f.write(key + "\t:\t'"+ val + "'")
    print("Finished decrypting cipher using SDES to file : {}".format(outFile))


if __name__ == "__main__":
    arr = [None] * keyLen
    generateAllPossibleKeys(keyLen, arr, 0)
    
    ask = True
    while ask:
        text = input("Enter s for SDES or t for Triple SDES to decrypt: ")
        if text =='s':
            decryptSDES(allKeys)
            break
        elif text == 't':
            decryptTripleSDES(allKeys)
            break
        else:
            print("Not a valid input")

