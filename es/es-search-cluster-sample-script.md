PUT ny-clusters-001
{
  "mappings": {
    "properties": {
      "cluster": {"type": "geo_shape"},
      "cluster_name": {"type": "text"}
    }
  }
}

POST ny-clusters-001/_doc
{
  "cluster_name": "NY City Hall",
  "cluster": "POLYGON ((-74.0080488 40.7120072, -74.0076089 40.7116819, -74.005608 40.7122105, -74.0051145 40.7123813, -74.0047819 40.7125521, -74.0042508 40.7131335, -74.0062947 40.7141257, -74.0080488 40.7120072))"
}

POST ny-clusters-001/_doc
{
  "cluster_name": "NY City Hall Cluster 1",
  "cluster": "POLYGON ((-74.0099852 40.712137, -74.0086549 40.7115352, -74.0063911 40.7141538, -74.0078931 40.7148532, -74.0099852 40.712137))"
}

POST ny-clusters-001/_doc
{
  "cluster_name": "NY City Hall Cluster 2",
  "cluster": "POLYGON ((-74.0048139 40.7124379, -74.0014129 40.7099168, -73.9983445 40.7132512, -74.0020995 40.715089, -74.0025072 40.7142921, -74.0033977 40.7136903, -74.0048139 40.7124379))"
}

GET ny-clusters-001/_search
{
  "query": {
    "geo_distance": {
      "distance": "71m",
      "cluster": {
        "lat": 40.71182,
        "lon": -74.00531
      }
    }
  }
}

GET ny-clusters-001/_search
{
  "query": {
    "geo_distance": {
      "distance": "1km",
      "cluster": {
        "lat": 40.71182,
        "lon": -74.00531
      }
    }
  }
}




