PUT geo-index-0001
{
  "mappings": {
    "properties": {
      "location": {"type": "geo_point"},
      "location_name": {"type": "text"}
    }
  }
}

POST geo-index-0001/_doc
{
  "location_name": "Time Square",
  "location": {
    "lat": 40.758049,
    "lon": -73.9878585
  }
}

POST geo-index-0001/_doc
{
  "location_name": "Val Cafe, Time Square",
  "location": {
    "lat": 40.75795,
    "lon": -73.9866843
  }
}

POST geo-index-0001/_doc
{
  "location_name": "Starbucks, Time Square",
  "location": {
    "lat": 40.7583137,
    "lon": -73.9868667
  }
}


POST geo-index-0001/_doc
{
  "location_name": "Gregorys Coffee, New York",
  "location": {
    "lat": 40.7342684,
    "lon": -74.0147306
  }
}

GET geo-index-0001/_search
{
  "query": {
    "geo_polygon": {
      "location": {
        "points": [
          [-73.9880304, 40.7597602], 
          [-73.9898329, 40.7572491], 
          [-73.9841359, 40.7548109], 
          [-73.9822691, 40.7573791], 
          [-73.9880304, 40.7597602]
        ]
      }
    }
  }  
}

GET geo-index-0001/_search
{
  "query": {
    "geo_shape": {
      "location": {
        "shape": "POLYGON ((-73.9880304 40.7597602, -73.9898329 40.7572491, -73.9841359 40.7548109, -73.9822691 40.7573791, -73.9880304 40.7597602))",
        "relation": "within"
      }
    }
  }  
}

GET geo-index-0001/_search
{
  "query": {
    "geo_distance": {
      "distance": "1km",
      "location": {
        "lat": 40.758049,
        "lon": -73.9878585 
      }
    }
  }
}

