	"github.com/m3db/m3/src/dbnode/encoding/proto"
	xtchannel "github.com/m3db/m3/src/dbnode/x/tchannel"
	xconfig "github.com/m3db/m3/src/x/config"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	"github.com/m3db/m3/src/x/instrument"
	xlog "github.com/m3db/m3/src/x/log"
	"github.com/m3db/m3/src/x/pool"
	xsync "github.com/m3db/m3/src/x/sync"
	"github.com/jhump/protoreflect/desc"
	var schema *desc.MessageDescriptor
	if cfg.Proto != nil {
		logger.Info("Probuf data mode enabled")
		schema, err = proto.ParseProtoSchema(cfg.Proto.SchemaFilePath)
		if err != nil {
			logger.Fatalf("error parsing protobuffer schema: %v", err)
		}
	}

	// Set the series cache policy.
	seriesCachePolicy := cfg.Cache.SeriesConfiguration().Policy
	opts = opts.SetSeriesCachePolicy(seriesCachePolicy)

	// Apply pooling options.
	opts = withEncodingAndPoolingOptions(cfg, logger, schema, opts, cfg.PoolingPolicy)
		},
		func(opts client.AdminOptions) client.AdminOptions {
			if cfg.Proto != nil {
				return opts.SetEncodingProto(
					schema,
					encoding.NewOptions(),
				).(client.AdminOptions)
			}
			return opts
		},
	)
	schema *desc.MessageDescriptor,
		if schema != nil {
			enc := proto.NewEncoder(time.Time{}, encodingOpts)
			enc.SetSchema(schema)
			return enc
		}

		if schema != nil {
			return proto.NewIterator(r, schema, encodingOpts)
		}
		SetReaderIteratorPool(iteratorPool).
	queryResultsPool := index.NewQueryResultsPool(
		poolOptions(policy.IndexResultsPool, scope.SubScope("index-query-results-pool")))
	aggregateQueryResultsPool := index.NewAggregateResultsPool(
		poolOptions(policy.IndexResultsPool, scope.SubScope("index-aggregate-results-pool")))
		SetQueryResultsPool(queryResultsPool).
		SetAggregateResultsPool(aggregateQueryResultsPool)
	queryResultsPool.Init(func() index.QueryResults {
		// NB(r): Need to initialize after setting the index opts so
		// it sees the same reference of the options as is set for the DB.
		return index.NewQueryResults(nil, index.QueryResultsOptions{}, indexOpts)
	})
	aggregateQueryResultsPool.Init(func() index.AggregateResults {
		return index.NewAggregateResults(nil, index.AggregateResultsOptions{}, indexOpts)