
import math
from random import randint
import sys
import signal

class DHKeyExchange:
    iterCount = 0
    def __init__(self,lower=2000,upper=4000) -> None:
        self.P = None
        self.Q = None  
        self.M = None  
        self.QPrimes = []
        for i in range(lower,upper+1):
            if self.isPrime(i):
                self.QPrimes.append(i)
        
        self.QPrimes = [ele for ele in reversed(self.QPrimes)]

        """ Find P such that P = 2Q + 1 and isPrime(P) == True         
            Start from biggest prime 
        """
        for q in self.QPrimes:
            p = 2*q + 1
            if self.isPrime(p):
                self.P = p
                self.Q = q
                break
        
        """ Find generator of the multiplicative group of integers modulo P. """
        self.G = self.primRoots(self.P)
        self.M = self.P * self.Q
        print("P : ",self.P, "\tQ : ",self.Q)
        print("Cyclic group selected G : n ∈ G | n=[1,2,...{}]".format(self.P-1))

    def isPrime(self,num):
        """ Check if a number is prime """
        sq = math.ceil(math.sqrt(num))
        for i in range(2,sq):
            if(num%i)==0:
                return False
        return True

    
    def primRoots(self,modulo):
        """This will select Largest primitive root (generator) """
        coprime_set = {num for num in range(1, modulo) if math.gcd(num, modulo) == 1}
        for g in range(modulo-1,1,-1):
            actual_set = set(pow(g, powers, modulo) for powers in range (1, modulo))
            if coprime_set == actual_set:
                return g 
    
    def genPubPriKeys(self,excludeKeys = []):
        """ Generate public and private keys from P by selecting random numbers from range 1-P """
        while(True):
            priKey = randint(1,self.P-1)
            if priKey not in excludeKeys:
                break
        
        pubKey = (self.G ** priKey) % self.P
        return pubKey,priKey

    def genSharedKey(self,pubKeyDest,priKey):
        return (pubKeyDest ** priKey) % self.P
    
    def genSecretKey(self,seed,count=10):
        """ Generate K1 and k2 for TripleDES """
        secretKey = ''
        X = 0
        x = pow(seed,2,self.M)
        for n in range(0, count):
            x = pow(x,2,self.M)
            X = bin(x)
            secretKey += str(X[len(X) - 1])
        
        return secretKey
