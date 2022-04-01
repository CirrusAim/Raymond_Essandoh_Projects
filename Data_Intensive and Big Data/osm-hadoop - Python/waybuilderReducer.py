#!/usr/bin/env python

import shapely.geometry as gm
import shapely.wkt as wkt
from ast import literal_eval
import sys

def getgeom(ways, nodes):
  if len(ways) != 1:
    return "Error: expected single way, found {}".format(len(ways))

  #sort []nodes by node.ordinal
  nodes.sort(key=lambda x:x["ordinal"])

  #extract just coords from sorted []nodes
  coords = []
  for i in nodes:
    point = gm.Point(long(i["lat"]), long(i["lon"]))
    coords.append(point)
  
  if len(coords) > 1:
    geom = gm.asLineString(coords)
    wktGeom = wkt.dumps(geom)
    return wktGeom

  return "Error: too short geometry"

#main():
currentID = None
ways = []
nodes = []

#line: wayID\t(nodeVals OR wayVals)
#nodeVals: type, nodeID, wayID, ordinal, lat, lon
#wayVals: type, wayID, []keys, []vals
for line in sys.stdin:
  key, val = line.split("\t")
  try:
    vals = literal_eval(val)
  except:
    continue

  #when last node for a wayID has been passed: process and return output
  if key != currentID:
    geom = getgeom(ways, nodes)
    outval = geom
    if geom[:5] != "Error":
      outval = ways[0]
      outval["geometry"] = geom
    print '%s\t%s' % (currentID, outval)
    currentID = key
    ways = []
    nodes = []
      
  if vals["type"] == "node":
    nodes.append(vals)
  if vals["type"] == "way":
    ways.append(vals)
