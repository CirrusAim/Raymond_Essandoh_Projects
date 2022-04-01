# py-osm-hadoop
Generate rasterized maps from OSM map files, using Hadoop and MapReduce for data processing

The implementation of this project is based on certain other github projects:
* General processing outline and functionality of major processing steps: http://github.com/willtemperley/osm-hadoop
* Extracting nodes/ways from compressed OSM data: http://github.com/qedus/osmpbf/blob/master/decode_tag.go

The code is missing several quality-of-life features, scripts not based on hadoop streaming must be edited to change input- and output locations.

## Steps to generate a map from an osm.pbf file:

# Step 0: Preprocessing
The .osm.pbf file format is not hadoop-friendly, the file must be split so that each block is on its own line. 
To perform this step, edit the "inputpath" and "outputpath" variables in the preprocessing.py file, then run it with 
```
python preprocessing.py
```

# Step 1: Join nodes and waynodes
OSM nodes contain (lat,lon) coordinates and OSM ways contain a list of nodeIDs that make up the road. Map-side, this step extracts nodes and ways from the compressed .osm.pbf file, writes out nodes and ways with their ID as key, and also writes (nodeID: parentWay, wayOrdinal) for each node that is part of a way. Reduce-side the nodes and waynodes are combined to fully specified way nodes, i.e. nodes with info on their location and on their parent way, and written out with the ID of the parent way as key.
To perform this step we use hadoop streaming:
```
hadoop jar /path/to/hadoop-stremaing.jar \
  -mapper nodejoinMapper.py \
  -reducer nodejoinReducer.py \
  -input /hdfs-path/to/preprocessed/osmfile \
  -output /hdfs-path/to/nodejoin/output \
  -file nodejoinMapper.py \
  -file nodejoinReducer.py
```

# Step 2: Build ways
In this step, ways and waynodes are combined so that each way now also has a properly ordered list of coordinates. The map-side sorts the ways and nodes based on wayID so that reduce-side can handle all nodes related to one way in one go. Reduce-side a shapely.geometry LineString is generated from the list of nodes related to each way, then appended to the way and written out again.
This step is also performed via hadoop streaming:
```
hadoop jar /path/to/hadoop-stremaing.jar \
  -mapper waybuilderMapper.py \
  -reducer waybuilderReducer.py \
  -input /hdfs-path/to/nodejoin/output \
  -output /hdfs-path/to/waybuild/output \
  -file waybuilderMapper.py \
  -file waybuilderReducer.py
```

# Step 3: Rasterize
Rasterization consists of "drawing" lines between each point on each way built in the previous step. "Drawing" ocnsists of outputting (x,y),z values for each point in the line, where x and y are latitude and longitude coordinates and together serve as the key, and z is the raster value which determines which color the corresponding pixel should be. Map-side draws every pixel for every line for every way, and reduce-side filters the pixel so only the pixel value of the most interest is written out.
hadoop streaming command:
```
hadoop jar /path/to/hadoop-stremaing.jar \
  -mapper rasterizerMapper.py \
  -reducer rasterizerReducer.py \
  -input /hdfs-path/to/waybuild/output \
  -output /hdfs-path/to/rasterizer/output \
  -file rasterizerMapper.py \
  -file rasterizerReducer.py
```

# Step 4: Extract the raster
This step uses Spark to read the rasterization output and filter it for coordinates within a user-defined bounding box, the raster values for each point are then added to an array at the right indices and written to a .png file with PyPNG.
To perform this step, first edit the "conf" dict at the top of the extractraster.py file to specify a bounding box by its minimum and maximum x- and y-coordinates, and specify input- and output paths like in step 0. Once the script is configured, execute it:
```
python extractraster.py
```
