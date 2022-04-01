#!/usr/bin/env python

from defs import *
from ast import literal_eval as lit
import shapely.wkt as wkt
import sys

def printpoint(pt, color): 
  print '%s,%s\t%s' % (int(pt[0]), int(pt[1]), color)

# for each (way: tags, geometry), output a pixel(x,y z)
# x,y = lat/lon
# z = pixel color value, based on tags
for line in sys.stdin:
  key, val = line.split("\t")
  try:
    way = lit(val)
  except:
    print '%s\t%s' % (key, "Error: {}".format(val))
    continue

  geom = wkt.loads(way["geometry"])

  # if the way is a highway, extract the type of highway
  # then assign the way a pixel color value based on way type, or 0 (black) if it's not a highway
  roadtype = ""
  for i in range(len(way["keys"])):
    if way["keys"][i] == "highway":
      roadtype = way["vals"][i]
  color = wayMap(roadtype)

  # for each point in the geometry of the way, "draw" a line from the previous point to that point
  # drawing consists of outputting a line with (lat,lon): color 
  if len(geom.coords) > 1 and color != 0:
    frompt = geom.coords[0]
    for i in range(1, len(geom.coords)):
      topt = geom.coords[i]
      drawline(frompt[0], frompt[1], topt[0], topt[1], 10**6, printpoint, color)
      frompt = topt
