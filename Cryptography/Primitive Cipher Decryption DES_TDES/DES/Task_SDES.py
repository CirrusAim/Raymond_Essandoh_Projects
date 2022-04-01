from TripleSDES import TripleSDES

def test_sdes_enc():
    print("Running Simple DES task 1 Encryption")
    tdes = TripleSDES()
    encSDES =  ((0b0000000000,0b00000000),
            (0b0000011111,0b11111111),
            (0b0010011111,0b11111100),
            (0b0010011111,0b10100101))
    for item in encSDES:
        cipher = tdes.encrypt(item[0],item[1])
        print("cipher : {0:08b}".format(cipher))
    print("Finished  Simple DES task 1 Encryption")
    
def test_sdes_dec():
    print("Running Simple DES task 1 Decryption")
    tdes = TripleSDES()
    decSDES =  ((0b1111111111,0b00001111),
            (0b0000011111,0b01000011),
            (0b1000101110,0b00011100),
            (0b1000101110,0b11000010))
    for item in decSDES:
        plaintext = tdes.decrypt(item[0],item[1])
        print("plaintext : {0:08b}".format(plaintext))
    print("Finished Simple DES task 1 Decryption")
    
def test_triple_sdes_enc():
    print("Running Triple SDES task 1 Encryption")
    tdes = TripleSDES()
    encSDES =  ((0b1000101110,0b0110101110,0b11010111),
            (0b1000101110,0b0110101110,0b10101010),
            (0b1111111111,0b1111111111,0b00000000),
            (0b0000000000,0b0000000000,0b01010010))
    for item in encSDES:
        cipher = tdes.tripleEncrypt(item[0],item[1],item[2])
        print("cipher : {0:08b}".format(cipher))
    print("Finished Triple SDES task 1 Encryption")
    
def test_triple_sdes_dec():
    print("Running Triple SDES task 1 Decryption")
    tdes = TripleSDES()
    decSDES =  ((0b1000101110,0b0110101110,0b11100110),
            (0b1011101111,0b0110101110,0b01010000),
            (0b1111111111,0b1111111111,0b00000100),
            (0b0000000000,0b0000000000,0b11110000))
    for item in decSDES:
        plaintext = tdes.tripleDecrypt(item[0],item[1],item[2])
        print("plaintext : {0:08b}".format(plaintext))
    print("Finished Triple SDES task 1 Decryption")

def task1():
    print("Task1 table Started")
    test_sdes_enc()
    test_sdes_dec()
    print("Task1 table Finished")

def task2():
    print("Task2 table Started")
    test_triple_sdes_enc()
    test_triple_sdes_dec()
    print("Task2 table Finished")

tdes = TripleSDES

if __name__ == '__main__':
    task1()
    task2()


