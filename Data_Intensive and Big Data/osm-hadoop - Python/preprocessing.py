from pyspark import *
from fileformat_pb2 import *
from defs import bintoint32
import time
import binascii

# script config
country = "georgia"
inputpath = ""
outputpath = ""
if inputpath == "" or outputpath == "":
  print "please specify paths for input and output locations!"
  return

start_time = time.time()
f = open(filepath, "rb")

#spark setup
conf = SparkConf()
conf.setMaster('local')
conf.setAppName('osmPreprocessing')
sc = SparkContext(conf = conf)

blobs = []
blocksRead = 0
bytesRead = 0
outFileNum = 0
while True:
  print("block #{}".format(blocksRead))
  blocksRead += 1
  hlen = f.read(4)
  bytesRead += 4
  try:
    hlenint = bintoint32(hlen)
  except:
    break

  bheader = BlobHeader()
  bheader.ParseFromString(f.read(hlenint))
  bytesRead += hlenint

  rawblob = f.read(bheader.datasize)
  hexblob = binascii.hexlify(rawblob)
  bytesRead += bheader.datasize
  if bheader.type == "OSMData":
    blobrecord = (bheader.type, hexblob)
    blobs.append(blobrecord)

  # output must be written in multiple rounds due to memory limitations
  if bytesRead > 30000000:
    print("byte limit reached! outputting file #{}".format(outFileNum))
    rdd = sc.parallelize(blobs)
    rdd.saveAsTextFile(outpath + "/{}".format(outFileNum))
    rdd = None
    blobs = []
    bytesRead = 0
    outFileNum += 1

print("bytes read: {}".format(bytesRead))
print("EOF reached! outputting file #{}".format(outFileNum))
rdd = sc.parallelize(blobs)
rdd.saveAsTextFile(outpath + "/{}".format(outFileNum))
end_time = time.time()

print("time taken: {}".format(end_time - start_time))
