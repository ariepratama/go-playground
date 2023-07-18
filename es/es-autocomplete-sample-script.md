```
PUT completion-index-0001
{
  "mappings": {
    "properties": {
      "text_suggest": {"type": "completion"},
      "text": {"type": "text"}
    }
  }
}

POST completion-index-0001/_doc
{
  "text": "I want to break free",
  "text_suggest": "I want to break free"
}

POST completion-index-0001/_doc
{
  "text": "work work work work work",
  "text_suggest": "work work work work work"
}

POST completion-index-0001/_doc
{
  "text": "bohemian Raphsody",
  "text_suggest": "bohemian Raphsody"
}

POST completion-index-0001/_search
{
  "suggest": {
    "x": {
      "text": ["i"],
      "completion": {
        "field": "text_suggest"
      } 
    }
    
  }
}


POST completion-index-0001/_search
{
  "suggest": {
    "x": {
      "text": ["wrk"],
      "completion": {
        "field": "text_suggest",
        "fuzzy": {
          "fuzziness": 2
        }
      } 
    }
    
  }
}
```