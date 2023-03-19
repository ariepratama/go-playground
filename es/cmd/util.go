package cmd

import (
	"fmt"
	"net/http"
)

// cleanIndex call DELETE http method to remove the index given a name
func cleanIndex(esIndexName string) {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:9200/%s", esIndexName), nil)
	client := &http.Client{}
	client.Do(req)
}
