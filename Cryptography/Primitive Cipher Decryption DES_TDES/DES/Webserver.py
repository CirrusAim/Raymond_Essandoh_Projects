from types import MethodType
from TripleSDES import TripleSDES
from flask import Flask
from flask import request
import os

app = Flask(__name__)

# Fixed Keys
key1 = 0b1000101110
key2 = 0b0110101110
tdes = TripleSDES()
dataLen = 8

def deciper(cipherString):
    dataList = [cipherString[i:i+dataLen] for i in range(0, len(cipherString), dataLen)]
    finalStr = ""
    for data in dataList:
        cp = int(data,2)
        pt = tdes.tripleDecrypt(key1,key2,cp)
        finalStr += chr(pt)
    return finalStr

@app.route('/index.js',methods=['POST','GET'])
@app.route('/',methods=['POST','GET'])
def decryptMessage():
    # if request.method == 'GET':
        # return "<p>Only POST requests allowed with following query params <b>cipher= binary-ciphertext<b> <p> "
    ciphertext = request.args.get('cipher')
    if not ciphertext:
        return "Cipher text requires as query param. E.g. /<url>?cipher=<ciphertext>"
    plaintext = deciper(ciphertext)
    return plaintext


if __name__ == '__main__':
    host = os.getenv('Triple SDES_HOST',default='localhost')
    port = os.getenv('Triple SDES_PORT',default='5000')
    print(host)
    print(port)
    app.run(host,port)
