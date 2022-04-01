#!/usr/bin/env python

from pyspark import *
from defs import *
from ast import literal_eval
from operator import itemgetter
import sys

currentID = None
outVals = {}

for line in sys.stdin:
  #key, val = ID, vals
  key, val = line.split("\t")
  vals = literal_eval(val)

  #if way
  if vals["type"] == "way":
    print '%s\t%s' % (key, vals)

  #if node/waynode
  if vals["type"] == "node":
    # when all nodes for one ID have been handled, output naand start processing next node
    if key != currentID:
      outVals["nodeID"] = currentID
      try:
        print '%s\t%s' % (outVals["wayID"], outVals)
      except:
        pass
      outVals = {}
      currentID = key
    for i in ["type", "ordinal", "wayID", "lat", "lon"]:
      try:
        outVals[i] = vals[i]
      except:
        continue
  
if key == currentID:
  try:
    print '%s\t%s' % (outVals["wayID"], outVals)
  except:
    pass
