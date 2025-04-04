//go:build integration
// +build integration

// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package integration

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/m3db/m3/src/cluster/services"
	"github.com/m3db/m3/src/cluster/shard"
	"github.com/m3db/m3/src/dbnode/integration/fake"
	"github.com/m3db/m3/src/dbnode/integration/generate"
	"github.com/m3db/m3/src/dbnode/namespace"
	"github.com/m3db/m3/src/dbnode/retention"
	"github.com/m3db/m3/src/dbnode/topology"
	"github.com/m3db/m3/src/dbnode/topology/testutil"
	"github.com/m3db/m3/src/x/ident"
	xtest "github.com/m3db/m3/src/x/test"
	xtime "github.com/m3db/m3/src/x/time"
)

func TestClusterDecommissionOneNode(t *testing.T) {
	testClusterDecommissionOneNode(t)
}

func testClusterDecommissionOneNode(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// Test setups.
	log := xtest.NewLogger(t)
	namesp, err := namespace.NewMetadata(testNamespaces[0],
		namespace.NewOptions().SetRetentionOptions(
			retention.NewOptions().
				SetRetentionPeriod(6*time.Hour).
				SetBlockSize(2*time.Hour).
				SetBufferPast(10*time.Minute).
				SetBufferFuture(2*time.Minute)))
	require.NoError(t, err)
	opts := NewTestOptions(t).
		SetNamespaces([]namespace.Metadata{namesp}).
		// Prevent snapshotting from happening too frequently to allow for the
		// possibility of a snapshot occurring after the shard set is assigned,
		// but not after the node finishes bootstrapping.
		SetTickMinimumInterval(5 * time.Second)

	minShard := uint32(0)
	maxShard := uint32(opts.NumShards()) - 1
	midShard := (maxShard - minShard) / 2

	instances := struct {
		start []services.ServiceInstance
		add   []services.ServiceInstance
		added []services.ServiceInstance
	}{
		start: []services.ServiceInstance{
			node(t, 0, concatShards(
				newClusterShardsRange(0, 5, shard.Available),
				concatShards(newClusterShardsRange(8, 9, shard.Available), newClusterShardsRange(11, 11, shard.Available)))),
			node(t, 1, concatShards(
				newClusterShardsRange(0, 3, shard.Available),
				concatShards(newClusterShardsRange(6, 8, shard.Available),
					newClusterShardsRange(10, 11, shard.Available)))),
			node(t, 2, concatShards(
				newClusterShardsRange(0, 1, shard.Available),
				concatShards(newClusterShardsRange(4, 7, shard.Available),
					newClusterShardsRange(9, 11, shard.Available)))),
			node(t, 3, concatShards(
				newClusterShardsRange(2, 3, shard.Available),
				newClusterShardsRange(4, 10, shard.Available))),
		},

		add: []services.ServiceInstance{
			node(t, 0, concatShards(
				newClusterShardsRange(0, 5, shard.Available),
				concatShards(newClusterShardsRange(8, 9, shard.Available),
					concatShards(newClusterShardsRange(11, 11, shard.Available),
						concatShards(newClusterShardsRange(6, 7, shard.Initializing),
							newClusterShardsRange(10, 10, shard.Initializing)))))),
			node(t, 1, concatShards(
				newClusterShardsRange(0, 3, shard.Available), concatShards(
					newClusterShardsRange(6, 8, shard.Available), concatShards(
						newClusterShardsRange(10, 11, shard.Available), concatShards(
							newClusterShardsRange(4, 5, shard.Initializing),
							newClusterShardsRange(9, 9, shard.Initializing)))))),
			node(t, 2, concatShards(newClusterShardsRange(0, 1, shard.Available),
				concatShards(newClusterShardsRange(4, 7, shard.Available),
					concatShards(newClusterShardsRange(9, 11, shard.Available),
						concatShards(newClusterShardsRange(2, 3, shard.Initializing),
							newClusterShardsRange(8, 8, shard.Initializing)))))),
			node(t, 3, concatShards(
				newClusterShardsRange(2, 3, shard.Leaving),
				newClusterShardsRange(4, 10, shard.Leaving))),
		},

		added: []services.ServiceInstance{
			node(t, 0, newClusterShardsRange(minShard, maxShard, shard.Available)),
			node(t, 1, newClusterShardsRange(minShard, maxShard, shard.Available)),
			node(t, 2, newClusterShardsRange(minShard, maxShard, shard.Available)),
		},
	}

	svc := fake.NewM3ClusterService().
		SetInstances(instances.start).
		SetReplication(services.NewServiceReplication().SetReplicas(3)).
		SetSharding(services.NewServiceSharding().SetNumShards(opts.NumShards()))

	svcs := fake.NewM3ClusterServices()
	svcs.RegisterService("m3db", svc)

	topoOpts := topology.NewDynamicOptions().
		SetConfigServiceClient(fake.NewM3ClusterClient(svcs, nil))
	topoInit := topology.NewDynamicInitializer(topoOpts)
	setupOpts := []BootstrappableTestSetupOptions{
		{
			DisablePeersBootstrapper: true,
			TopologyInitializer:      topoInit,
		},
		{
			DisablePeersBootstrapper: false,
			TopologyInitializer:      topoInit,
		},
		{
			DisablePeersBootstrapper: false,
			TopologyInitializer:      topoInit,
		},
		{
			DisablePeersBootstrapper: false,
			TopologyInitializer:      topoInit,
		},
	}
	setups, closeFunction := NewDefaultBootstrappableTestSetups(t, opts, setupOpts)
	defer closeFunction()

	// Write test data for first node.
	topo, err := topoInit.Init()
	require.NoError(t, err)
	ids := []idShard{}

	// Boilerplate code to find two ID's that hash to the first half of the
	// shards, and one ID that hashes to the second half of the shards.
	shardSet := topo.Get().ShardSet()
	i := 0
	numFirstHalf := 0
	numSecondHalf := 0
	for {
		if numFirstHalf == 2 {
			break
		}
		idStr := strconv.Itoa(i)
		shard := shardSet.Lookup(ident.StringID(idStr))
		if shard < midShard && numFirstHalf < 2 {
			ids = append(ids, idShard{str: idStr, shard: shard})
			numFirstHalf++
		}
		if shard > midShard && numSecondHalf < 1 {
			ids = append(ids, idShard{idStr, shard})
			numSecondHalf++
		}
		i++
	}

	for _, id := range ids {
		// Verify IDs will map to halves of the shard space.
		require.Equal(t, id.shard, shardSet.Lookup(ident.StringID(id.str)))
	}

	var (
		now        = setups[0].NowFn()()
		blockStart = now
		blockSize  = namesp.Options().RetentionOptions().BlockSize()
		seriesMaps = generate.BlocksByStart([]generate.BlockConfig{
			{IDs: []string{ids[0].str, ids[1].str}, NumPoints: 180, Start: blockStart.Add(-blockSize)},
			{IDs: []string{ids[0].str, ids[1].str}, NumPoints: 90, Start: blockStart},
		})
	)
	err = writeTestDataToDisk(namesp, setups[0], seriesMaps, 0)
	require.NoError(t, err)

	// Prepare verification of data on nodes.
	expectedSeriesMaps := make([]map[xtime.UnixNano]generate.SeriesBlock, 2)
	expectedSeriesIDs := make([]map[string]struct{}, 2)
	for i := range expectedSeriesMaps {
		expectedSeriesMaps[i] = make(map[xtime.UnixNano]generate.SeriesBlock)
		expectedSeriesIDs[i] = make(map[string]struct{})
	}
	for start, series := range seriesMaps {
		list := make([]generate.SeriesBlock, 2)
		for j := range series {
			if shardSet.Lookup(series[j].ID) < midShard+1 {
				list[0] = append(list[0], series[j])
			} else {
				list[1] = append(list[1], series[j])
			}
		}
		for i := range expectedSeriesMaps {
			if len(list[i]) > 0 {
				expectedSeriesMaps[i][start] = list[i]
			}
		}
	}
	for i := range expectedSeriesMaps {
		for _, series := range expectedSeriesMaps[i] {
			for _, elem := range series {
				expectedSeriesIDs[i][elem.ID.String()] = struct{}{}
			}
		}
	}
	require.Equal(t, 2, len(expectedSeriesIDs[0]))

	// Start the first server with filesystem bootstrapper.
	require.NoError(t, setups[0].StartServer())

	// Start the 2nd server
	require.NoError(t, setups[1].StartServer())

	// Start the 3rd server
	require.NoError(t, setups[2].StartServer())

	// Start the 4th server
	require.NoError(t, setups[3].StartServer())

	log.Debug("servers are now up")

	// Stop the servers on test completion.
	defer func() {
		setups.parallel(func(s TestSetup) {
			require.NoError(t, s.StopServer())
		})
		log.Debug("servers are now down")
	}()

	// Bootstrap the new shards.
	log.Debug("re-sharding to decommission shards")
	svc.SetInstances(instances.add)
	svcs.NotifyServiceUpdate("m3db")
	go func() {
		for {
			time.Sleep(time.Second)
			for _, setup := range setups {
				now = now.Add(time.Second)
				setup.SetNowFn(now)
			}
		}
	}()

	// Generate some new data that will be written to the node while peer streaming is taking place
	// to make sure that the data that is streamed in and the data that is received while streaming
	// is going on are both handled correctly. In addition, this will ensure that we hold onto both
	// sets of data durably after topology changes and that the node can be properly bootstrapped
	// from just the filesystem and commitlog in a later portion of the test.
	seriesToWriteDuringPeerStreaming := []string{
		"series_after_bootstrap1",
		"series_after_bootstrap2",
	}
	// Ensure that the new series belong that we're going to write belong to the host that is peer
	// streaming data.
	for _, seriesName := range seriesToWriteDuringPeerStreaming {
		shard := shardSet.Lookup(ident.StringID(seriesName))
		require.True(t, shard > midShard,
			fmt.Sprintf("series: %s does not shard to second host", seriesName))
	}
	seriesReceivedDuringPeerStreaming := generate.BlocksByStart([]generate.BlockConfig{
		{IDs: seriesToWriteDuringPeerStreaming, NumPoints: 90, Start: blockStart},
	})
	// Merge the newly generated series into the expected series map.
	for blockStart, series := range seriesReceivedDuringPeerStreaming {
		expectedSeriesMaps[1][blockStart] = append(expectedSeriesMaps[1][blockStart], series...)
	}

	// Spin up a background goroutine to issue the writes to the node while its streaming data
	// from its peer.
	doneWritingWhilePeerStreaming := make(chan struct{})
	go func() {
		for _, testData := range seriesReceivedDuringPeerStreaming {
			err := setups[1].WriteBatch(namesp.ID(), testData)
			// We expect consistency errors because we're only running with
			// R.F = 2 and one node is leaving and one node is joining for
			// each of the shards that is changing hands.
			if err != nil {
				panic(err)
			}
		}
		doneWritingWhilePeerStreaming <- struct{}{}
	}()

	log.Debug("waiting for shards to be bootstrapped")
	waitUntilHasBootstrappedShardsExactly(setups[0].DB(), testutil.Uint32Range(minShard, maxShard))
	waitUntilHasBootstrappedShardsExactly(setups[1].DB(), testutil.Uint32Range(minShard, maxShard))
	waitUntilHasBootstrappedShardsExactly(setups[2].DB(), testutil.Uint32Range(minShard, maxShard))

	log.Debug("waiting for background writes to complete")
	<-doneWritingWhilePeerStreaming
	log.Debug("Background writes completed")

	log.Debug("waiting for shards to be marked initialized")

	allMarkedAvailable := func(
		fakePlacementService fake.M3ClusterPlacementService,
		instanceID string,
	) bool {
		markedAvailable := fakePlacementService.InstanceShardsMarkedAvailable()

		if len(markedAvailable) != 3 {
			return false
		}

		marked := shard.NewShards(nil)
		for _, id := range markedAvailable[instanceID] {
			marked.Add(shard.NewShard(id).SetState(shard.Available))
		}
		return true
	}

	fps := svcs.FakePlacementService()

	for !allMarkedAvailable(fps, "testhost1") {
		time.Sleep(100 * time.Millisecond)
	}
	log.Debug("all shards marked as initialized")

	// Shed the old shards from the first node
	log.Debug("resharding to shed shards from decomissioning node")
	svc.SetInstances(instances.added)
	svcs.NotifyServiceUpdate("m3db")
	waitUntilHasBootstrappedShardsExactly(setups[0].DB(), testutil.Uint32Range(minShard, maxShard))
	waitUntilHasBootstrappedShardsExactly(setups[1].DB(), testutil.Uint32Range(minShard, maxShard))
	waitUntilHasBootstrappedShardsExactly(setups[2].DB(), testutil.Uint32Range(minShard, maxShard))

	log.Debug("verifying data in servers matches expected data set")
	// Verify in-memory data match what we expect
	verifySeriesMaps(t, setups[1], namesp.ID(), expectedSeriesMaps[1])
}
