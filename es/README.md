## Getting Started
I used homebrew installation for elasticsearch since using docker seems to be counterproductive for this tutorial.
Reference: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/brew.html

After installing, you should start the elastic search via
```
brew services start elasticsearch-full
```

then after waiting about 1 minute, try to do
```
curl localhost:9200
```

it should return json response, like
```
{
  "name" : "Macbook.local",
  "cluster_name" : "elasticsearch_fawef",
  "cluster_uuid" : "k6OTjAJ_SjSiDniGAc4e2Q",
  "version" : {
    "number" : "7.17.4",
    "build_flavor" : "default",
    "build_type" : "tar",
    "build_hash" : "79878662c54c886ae89206c685d9f1051a9d6411",
    "build_date" : "2022-05-18T18:04:20.964345128Z",
    "build_snapshot" : false,
    "lucene_version" : "8.11.1",
    "minimum_wire_compatibility_version" : "6.8.0",
    "minimum_index_compatibility_version" : "6.0.0-beta1"
  },
  "tagline" : "You Know, for Search"
}
```