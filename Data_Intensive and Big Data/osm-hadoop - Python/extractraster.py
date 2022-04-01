from pyspark import *
from fileformat_pb2 import *
from defs import bintoint32, filtercoords, addtoarr
import png
import time
import binascii

start_time = time.time()

# y comes before x due to confusion between what latitude and longitude represent earlier in implementation
conf = {
  "country": "georgia3",
  "inpath": "",
  "outpath": "",
  "miny": 40000,
  "maxy": 50000,
  "minx": 39990,
  "maxx": 47000,
}

#spark setup
sconf = SparkConf()
sconf.setMaster('local')
sconf.setAppName('osmPreprocessing')
sc = SparkContext(conf = sconf)

filepath = conf["inpath"]
rasterrdd = sc.textFile(filepath)
rasterarr = []
dy = conf["maxy"] - conf["miny"]
dx = conf["maxx"] - conf["minx"]
for i in range(dy):
  rasterarr.append([0]*dx)

filteredraster = rasterrdd.filter(lambda x: filtercoords(x, conf["minx"], conf["miny"], conf["maxx"], conf["maxy"])).collect()
for line in filteredraster:
  addtoarr(line, rasterarr, conf["minx"], conf["miny"])

# output png image
outpath = conf["outpath"]
png.from_array(rasterarr, 'L').save(outpath)

end_time = time.time()
print("time taken: {}".format(end_time - start_time))
