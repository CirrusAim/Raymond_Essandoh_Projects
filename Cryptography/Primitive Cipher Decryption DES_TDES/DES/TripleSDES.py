import unittest

class TripleSDES:
    KeyLength = 10
    SubKeyLength = 8
    DataLength = 8
    FLength = 4
    
    # Tables for initial and final permutations (b1, b2, b3, ... b8)
    IPtable = (2, 6, 3, 1, 4, 8, 5, 7)
    FPtable = (4, 1, 3, 5, 7, 2, 8, 6)
    
    # Tables for subkey generation (k1, k2, k3, ... k10)
    P10table = (3, 5, 2, 7, 4, 10, 1, 9, 8, 6)
    P8table = (6, 3, 7, 4, 8, 5, 10, 9)
    
    # Tables for the fk function
    EPtable = (4, 1, 2, 3, 2, 3, 4, 1)
    S0table = (1, 0, 3, 2, 3, 2, 1, 0, 0, 2, 1, 3, 3, 1, 3, 2)
    S1table = (0, 1, 2, 3, 2, 0, 1, 3, 3, 0, 1, 0, 2, 1, 0, 3)
    P4table = (2, 4, 3, 1)
    
    def __init__(self):
        pass

    def perm(self,inputByte, permTable):
        """Permute input byte according to permutation table"""
        outputByte = 0
        for index, elem in enumerate(permTable):
            if index >= elem:
                outputByte |= (inputByte & (128 >> (elem - 1))) >> (index - (elem - 1))
            else:
                outputByte |= (inputByte & (128 >> (elem - 1))) << ((elem - 1) - index)
        return outputByte
    
    def ip(self,inputByte):
        """Perform the initial permutation on data"""
        return self.perm(inputByte, self.IPtable)
    
    def fp(self,inputByte):
        """Perform the final permutation on data"""
        return self.perm(inputByte, self.FPtable)
    
    def swapNibbles(self,inputByte):
        """Swap the two nibbles of data"""
        return (inputByte << 4 | inputByte >> 4) & 0xff
    
    def keyGen(self,key):
        """Generate the two required subkeys"""
        def leftShift(keyBitList):
            """Perform a circular left shift on the first and second five bits"""
            shiftedKey = [None] * self.KeyLength
            shiftedKey[0:9] = keyBitList[1:10]
            shiftedKey[4] = keyBitList[0]
            shiftedKey[9] = keyBitList[5]
            return shiftedKey
    
        # Converts input key (integer) into a list of binary digits
        keyList = [(key & 1 << i) >> i for i in reversed(range(self.KeyLength))]

        permKeyList = [None] * self.KeyLength
        for index, elem in enumerate(self.P10table):
            permKeyList[index] = keyList[elem - 1]
        shiftedOnceKey = leftShift(permKeyList)
        shiftedTwiceKey = leftShift(leftShift(shiftedOnceKey))
        subKey1 = subKey2 = 0
        for index, elem in enumerate(self.P8table):
            subKey1 += (128 >> index) * shiftedOnceKey[elem - 1]
            subKey2 += (128 >> index) * shiftedTwiceKey[elem - 1]
        return (subKey1, subKey2)
    
    def feistelFun(self,sKey, rightNibble):
        aux = sKey ^ self.perm(self.swapNibbles(rightNibble), self.EPtable)
        index1 = ((aux & 0x80) >> 4) + ((aux & 0x40) >> 5) + ((aux & 0x20) >> 5) + ((aux & 0x10) >> 2)
        index2 = ((aux & 0x08) >> 0) + ((aux & 0x04) >> 1) + ((aux & 0x02) >> 1) + ((aux & 0x01) << 2)
        sboxOutputs = self.swapNibbles((self.S0table[index1] << 2) + self.S1table[index2])
        return self.perm(sboxOutputs, self.P4table)

    def fk(self,subKey, inputData):
        """Apply Feistel function on data with given subkey"""
        leftNibble, rightNibble = inputData & 0xf0, inputData & 0x0f
        return (leftNibble ^ self.feistelFun(subKey, rightNibble)) | rightNibble
    
    def encrypt(self,key, plaintext):
        """Encrypt plaintext with given key"""
        data = self.fk(self.keyGen(key)[0], self.ip(plaintext))
        cipher =  self.fp(self.fk(self.keyGen(key)[1], self.swapNibbles(data)))
        return cipher
    
    def decrypt(self,key, ciphertext):
        """Decrypt ciphertext with given key"""
        data = self.fk(self.keyGen(key)[1], self.ip(ciphertext))
        plaintext =  self.fp(self.fk(self.keyGen(key)[0], self.swapNibbles(data)))  
        return plaintext
    
    def tripleEncrypt(self,key1, key2, plaintext):
        """ Enc = Enc(key1, Dec(k2, Enc(k1, plaintext))) """
        iteration1 = self.encrypt(key1,plaintext)
        iteration2 = self.decrypt(key2, iteration1)
        cipher = self.encrypt(key1,iteration2)
        return cipher
    def tripleDecrypt(self,key1,key2,cipher):
        """ Dec = Dec(key1, Enc(k2, Dec(k1, plaintext))) """
        iteration1 = self.decrypt(key1,cipher)
        iteration2 = self.encrypt(key2, iteration1)
        plaintext = self.decrypt(key1,iteration2)
        return plaintext



class TestStringMethods(unittest.TestCase):
    def test_verify_implementation(self):
        tdes = TripleSDES()
        self.assertEqual(tdes.encrypt(0b0000000000,0b10101010), 0b00010001)
        self.assertEqual(tdes.encrypt(0b1110001110,0b10101010), 0b11001010)
        self.assertEqual(tdes.encrypt(0b1110001110,0b01010101), 0b01110000)
        self.assertEqual(tdes.encrypt(0b1111111111,0b10101010), 0b00000100)


if __name__ == '__main__':
    unittest.main()