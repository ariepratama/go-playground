package main

import (
	"context"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"os"
)

const EnvHzServer = "HZ_SERVER"
const EnvHzClusterName = "HZ_CLUSTER_NAME"

func main() {
	hzServer := os.Getenv(EnvHzServer)
	hzClusterName := os.Getenv(EnvHzClusterName)

	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses(fmt.Sprintf("%s:5701", hzServer))
	config.Cluster.Name = hzClusterName

	ctx := context.TODO()
	client, err := hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	sampleMap, err := client.GetMap(ctx, "sample_map")
	if err != nil {
		panic(err)
	}
	key := "sample_key"
	err = sampleMap.Set(ctx, key, "sample_value")
	if err != nil {
		panic(err)
	}

	val, err := sampleMap.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Got value: ", val)
}
