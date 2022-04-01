#!/usr/bin/env python

from defs import *
import binascii
import zlib
import sys

#assuming this reads each line as a [key, val] pair
for line in sys.stdin:
  #each line is "(u'OSMData', '{zlib-encoded blob}') ", so indices 14:-3 fit to extract jsut the zlib data
  linedata = bytes(line[14:-3]) 
  blobstr = binascii.unhexlify(linedata)
  blob = Blob()
  blob.ParseFromString(blobstr)
  entities = parseBlob(blob)
  for entity in entities:
    #map way
    if entity["type"] == "way":
      #map wayNodes to their nodeID
      wayNodeVals = {"type": "node", "wayID": entity["id"], "ordinal": -1}
      for nodeID in entity["refs"]:
        wayNodeVals["ordinal"] = wayNodeVals["ordinal"] + 1
        print '%s\t%s' % (nodeID, wayNodeVals)
    
      #map way to its wayID
      wayID = entity["id"]
      wayVals = {"type": "way", "keys": entity["keys"], "vals": entity["vals"]}
      print '%s\t%s' % (wayID, wayVals)

    #map node
    elif entity["type"] == "node":
      nodeVals = {"type": "node", "lat": entity["lat"], "lon": entity["lon"]}
      print '%s\t%s' % (entity["id"], nodeVals)
  
    else:
      print '%s\t%s' % ("error: ", entity)
