### Fit

Fit is a toolkit for exploring and manipulating datasets. Go has [many](https://github.com/gonum) [great](https://github.com/montanaflynn/stats) 
statistical libraries but is largely unreprented in the data science world.
Fit aims to be a general purpose tool for handling data [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load),
analytics, and visualization.

Note that this project is *largely unfinished* and unsuitable for any practical purpose.


#### Running


    # Start the Fit HTTP server
    fit server
    # Load sample data
    fit load --name LakeHuron sample_data/LakeHuron.csv
    # List datasets
    fit ls

    NAME      ROWS  COLS  COLUMNS          
    LakeHuron 98    3     [ time LakeHuron]

    # Query
    fit query -n "LakeHuron,time,LakeHuron"

    [time LakeHuron]

    Dims(98, 2)
      ⎡  1875  580.38⎤
      ⎢  1876  581.86⎥
      ⎢  1877  580.97⎥
      ⎢  1878   580.8⎥
      ⎢  1879  579.79⎥
     .
     .
     .
      ⎢  1968  578.52⎥
      ⎢  1969  579.74⎥
      ⎢  1970  579.31⎥
      ⎢  1971  579.89⎥
      ⎣  1972  579.96⎦


    fit --json query -n "LakeHuron,time,LakeHuron" |jq .

    {
    "Name": "QueryResult",
    "Columns": [
      "time",
      "LakeHuron"
    ],
    "Stats": {
      "Rows": 98,
      "Columns": 2
    },
    "Mtx": [
      1875,
      580.38 ....

    # Open your web browser and perform the same query:
    # http://localhost:8000/explore?q=LakeHuron,time&q=LakeHuron,LakeHuron


