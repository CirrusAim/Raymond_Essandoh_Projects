from fileformat_pb2 import *
from osmformat_pb2 import *
#from pyspark import *
import zlib
import time

def bintoint32(bytes):
  ret = 0
  bint = [b for b in bytearray(bytes)]
  for i in range(4):
    ret = ret + 16**(3-i)*int(bint[i])
  return ret

def parseBlob(blob):
  blobdata = PrimitiveBlock()
  blobdata.ParseFromString(zlib.decompress(blob.zlib_data))
  return parsePrimitiveBlock(blobdata)

# code for parsePrimitiveBlock ported from qedus' "osmpbf" go code on github:
## https://github.com/qedus/osmpbf/blob/master/decode_data.go
def parsePrimitiveBlock(primitiveBlock):
  nodes = parseDenseNodes(primitiveBlock)
  ways = parseWays(primitiveBlock)

  return nodes + ways

def parseWays(pb):
  outways = []
  ways = pb.primitivegroup[0].ways
  st = pb.stringtable.s
  #dateGranularity = pb.date_granularity

  for i in range(len(ways)):
    way = ways[i]
    tags = extractTags(st, way.keys, way.vals)

    #refs with delta encoding
    nodeID = 0
    nodeIDs = []
    for o in way.refs:
      nodeID = o + nodeID
      nodeIDs.append(nodeID)

    #info = extractInfo(st, way.info, dateGranularity)

    newway = {}
    newway["type"] = "way"
    newway["id"] = way.id
    newway["refs"] = nodeIDs
    newway["keys"] = []
    newway["vals"] = []
    for o in range(len(tags[0])):
      newway["keys"].append(tags[0][o])
      newway["vals"].append(tags[1][o])
    #newway["info"] = way.info
    outways.append(newway)

  return outways

def parseDenseNodes(pb):
  dense = pb.primitivegroup[0].dense
  st = pb.stringtable.s
  granularity = pb.granularity
  #dateGranularity = pb.date_granularity
  latOffset = pb.lat_offset
  lonOffset = pb.lon_offset

  outnodes = []
  deltaID = 0
  deltalat = 0
  deltalon = 0
  #state = denseInfoState()
  #tu = tagUnpacker(st, dense.keys_vals, 0)
  for i in range(len(dense.id)):
    deltaID += dense.id[i]
    deltalat += dense.lat[i]
    deltalon += dense.lon[i]
    latOff = (latOffset + granularity*deltalat)
    lonOff = (lonOffset + granularity*deltalon)
    #tags = tu.next()
    #info = extractDenseInfo(state, st, dense.denseinfo, i, dateGranularity)

    node = {}
    node["type"] = "node"
    node["id"] = deltaID
    node["lat"] = latOff
    node["lon"] = lonOff
    #node["keys"] = []
    #node["vals"] = []
    #for o in range(len(tags[0])):
    #  node["keys"].append(tags[0][o])
    #  node["vals"].append(tags[1][o])
    #node["info"]= info
    outnodes.append(node)
  return outnodes

def extractInfo(strTable, inf, dateGranularity):
  info = Info()
  info.visible = True

  if inf != None:
    info.version = inf.version
    #TODO: timestamp
    info.changeset = inf.changeset
    info.uid = inf.uid
    

#TODO: timestamp
def extractDenseInfo(state, strTable, inf, index, dateGranularity):
  info = Info()
  info.visible = True

  #versions
  if len(inf.version) > 0:
    info.version = inf.version[index]

  #timestamps
  if False: #len(inf.timestamp) > 0:
    state.timestamp = inf.timestamp[index] + state.timestamp
    millisec = time.Duration(state.timestamp * dateGranularity) * time.Millisecond
    info.timestamp = time.unix(0, millisec.Nanoseconds()).UTC()

  #changesets
  if len(inf.changeset) > 0:
    state.changeset = inf.changeset[index] + state.changeset
    info.changeset = state.changeset

  #uids
  if len(inf.uid) > 0:
    state.uid = inf.uid[index]

  #usersids
  if False: # len(inf.user_sid) > 0:
    state.userSid = str(int(inf.user_sid[index]) + int(state.userSid))
    info.user_sid = strTable[int(state.userSid)]

  #visible
  if len(inf.visible) > 0:
    info.visible = inf.visible[index]
    
  return info

def extractTags(strTable, keyIDs, valIDs):
  keys = []
  vals = []
  for i in keyIDs:
    keys.append(strTable[i])
  for i in valIDs:
    vals.append(strTable[i])
  return [keys, vals]

class denseInfoState:
  def __init__(self):
    self.timestamp = 0
    self.changeset = 0
    self.uid = 0
    self.userSid = 0

class tagUnpacker:
  def __init__(self, st, keysVals, index):
    self.strTable = st
    self.keysVals = keysVals
    self.index = index
  
  def next(self):
    keys = []
    vals = []
    while(self.index < len(self.keysVals)):
      keyID = self.keysVals[self.index]
      self.index += 1
      if keyID == 0:
        break

      valID = self.keysVals[self.index]
      self.index += 1

      keys.append(self.strTable[keyID])
      vals.append(self.strTable[valID])
    return (keys, vals)

def wayMap(strType): 
  map = {
    "morotway": 15,
    "trunk": 14,
    "railway": 13,
    "primary": 12,
    "secondary": 11,
    "tertiary": 10,
    "motorway link": 9,
    "primary link": 8,
    "unclassified": 7,
    "road": 6,
    "residential": 5,
    "service": 3,
    "track": 2,
    "pedestrian": 1,
  }
  try: 
    return map[strType]
  except:
    return 0

def drawline(x1,y1,x2,y2, scale, func, color):
  dx = abs(x2 - x1)
  dy = abs(y2 - y1)

  if dy < dx:
    if x2 < x1:
      return drawline(x2, y2, x1, y1, scale, func, color)
    elif y1 < y2:
      return drawplusx(round(x1/scale), round(y1/scale), round(x2/scale), round(y2/scale), func, color)
    else:
      return drawminusx(round(x1/scale), round(y1/scale), round(x2/scale), round(y2/scale), func, color)

  else:
    if y2 < y1:
      return drawline(x2, y2, x1, y1, scale, func, color)
    elif x1 < x2:
      return drawplusy(round(x1/scale), round(y1/scale), round(x2/scale), round(y2/scale), func, color)
    else:
      return drawminusy(round(x1/scale), round(y1/scale), round(x2/scale), round(y2/scale), func, color)

def drawplusx(x1,y1,x2,y2, func, color):
  dx = x2 - x1
  dy = y2 - y1

  x = x1
  y = y1
  err = 0

  while (x <= x2):
    pt = (x,y)
    func(pt, color)
    x += 1
    err += dy
    if(2*err >= dx):
      y += 1
      err -= dx

def drawminusx(x1,y1,x2,y2, func, color):
  dx = x2 - x1
  dy = y2 - y1

  x = x1
  y = y1
  err = 0

  while (x <= x2):
    pt = (x,y)
    func(pt, color)
    x += 1
    err += dy
    if(2*err < dx):
      y -= 1
      err += dx

def drawplusy(x1,y1,x2,y2, func, color):
  dx = x2 - x1
  dy = y2 - y1

  x = x1
  y = y1
  err = 0

  while (y <= y2):
    pt = (x,y)
    func(pt, color)
    y += 1
    err += dx
    if(2*err >= dy):
      x += 1
      err -= dy

def drawminusy(x1,y1,x2,y2, func, color):
  dx = x2 - x1
  dy = y2 - y1

  x = x1
  y = y1
  err = 0

  while (y <= y2):
    pt = (x,y)
    func(pt, color)
    y += 1
    err += dx
    if(2*err < dy):
      x -= 1
      err += dy

# raster extraction

# line in raster/{country}/part-00000 = "ycoord,xcoord\tcolorint"
# coords are y,x rather than x,y due to confusion between latitude and longitude,
# code has been modified to compensate
def filtercoords(line, minx, miny, maxx, maxy):
  try:
    coords = line.split("\t")[0].split(",")
    x = int(coords[1])
    y = int(coords[0])
    if minx <= x and x <= maxx and miny <= y and y <= maxy:
      return True
    return False 
  except:
    return False

def addtoarr(line, arr, minx=0, maxy=0):
  try:
    coords, val = line.split("\t")
    x = int(coords.split(",")[1])
    y = int(coords.split(",")[0])
    arrlenx = len(arr[0])
    arrleny = len(arr)
    arr[maxy - y][x - minx] = int(val)*17
    return True
  except:
    return False
