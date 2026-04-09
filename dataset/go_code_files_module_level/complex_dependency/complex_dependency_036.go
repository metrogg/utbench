func NewFromConfig(logger *logrus.Logger, conf Config) (*Server, error) {
	ret := &Server{}

	ret.Hostname = conf.Hostname
	ret.Tags = conf.Tags

	mappedTags := samplers.ParseTagSliceToMap(ret.Tags)

	ret.synchronizeInterval = conf.SynchronizeWithInterval

	ret.TagsAsMap = mappedTags
	ret.HistogramPercentiles = conf.Percentiles
	ret.HistogramAggregates.Value = 0
	for _, agg := range conf.Aggregates {
		ret.HistogramAggregates.Value += samplers.AggregatesLookup[agg]
	}
	ret.HistogramAggregates.Count = len(conf.Aggregates)

	var err error
	ret.interval, err = conf.ParseInterval()
	if err != nil {
		return ret, err
	}

	ret.stuckIntervals = conf.FlushWatchdogMissedFlushes

	transport := &http.Transport{
		IdleConnTimeout: ret.interval * 2, // If we're idle more than one interval something is up
	}

	ret.HTTPClient = &http.Client{
		// make sure that POSTs to datadog do not overflow the flush interval
		Timeout:   ret.interval * 9 / 10,
		Transport: transport,
	}

	stats, err := statsd.NewBuffered(conf.StatsAddress, 4096)
	if err != nil {
		return ret, err
	}
	stats.Namespace = "veneur."

	ret.Statsd = stats

	ret.SpanChan = make(chan *ssf.SSFSpan, conf.SpanChannelCapacity)
	ret.TraceClient, err = trace.NewChannelClient(ret.SpanChan,
		trace.ReportStatistics(stats, 1*time.Second, []string{"ssf_format:internal"}),
	)
	if err != nil {
		return ret, err
	}

	// nil is a valid sentry client that noops all methods, if there is no DSN
	// we can just leave it as nil
	if conf.SentryDsn != "" {
		ret.Sentry, err = raven.New(conf.SentryDsn)
		if err != nil {
			return ret, err
		}
	}

	if conf.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	mpf := 0
	if conf.MutexProfileFraction > 0 {
		mpf = runtime.SetMutexProfileFraction(conf.MutexProfileFraction)
	}

	log.WithFields(logrus.Fields{
		"MutexProfileFraction":         conf.MutexProfileFraction,
		"previousMutexProfileFraction": mpf,
	}).Info("Set mutex profile fraction")

	if conf.BlockProfileRate > 0 {
		runtime.SetBlockProfileRate(conf.BlockProfileRate)
	}
	log.WithField("BlockProfileRate", conf.BlockProfileRate).Info("Set block profile rate (nanoseconds)")

	if conf.EnableProfiling {
		ret.enableProfiling = true
	}

	// This is a check to ensure that we don't repeatedly add a hook
	// to the "global" log instance on repeated calls to `NewFromConfig`
	// such as those made in testing. By skipping this we avoid a race
	// condition in logrus discussed here:
	// https://github.com/sirupsen/logrus/issues/295
	if _, ok := logger.Hooks[logrus.FatalLevel]; !ok {
		logger.AddHook(sentryHook{
			c:        ret.Sentry,
			hostname: ret.Hostname,
			lv: []logrus.Level{
				logrus.ErrorLevel,
				logrus.FatalLevel,
				logrus.PanicLevel,
			},
		})
	}

	// After log hooks are configured, if any further errors are
	// found during setup we should use logger.Fatal or
	// logger.Panic, because that will give us breakage monitoring
	// through Sentry.
	numWorkers := 1
	if conf.NumWorkers > 1 {
		numWorkers = conf.NumWorkers
	}
	logger.WithField("number", numWorkers).Info("Preparing workers")
	// Allocate the slice, we'll fill it with workers later.
	ret.Workers = make([]*Worker, numWorkers)
	ret.numReaders = conf.NumReaders

	// Use the pre-allocated Workers slice to know how many to start.
	for i := range ret.Workers {
		ret.Workers[i] = NewWorker(i+1, ret.TraceClient, log, ret.Statsd)
		// do not close over loop index
		go func(w *Worker) {
			defer func() {
				ConsumePanic(ret.Sentry, ret.TraceClient, ret.Hostname, recover())
			}()
			w.Work()
		}(ret.Workers[i])
	}

	ret.EventWorker = NewEventWorker(ret.TraceClient, ret.Statsd)

	// Set up a span sink that extracts metrics from SSF spans and
	// reports them via the metric workers:
	processors := make([]ssfmetrics.Processor, len(ret.Workers))
	for i, w := range ret.Workers {
		processors[i] = w
	}
	metricSink, err := ssfmetrics.NewMetricExtractionSink(processors, conf.IndicatorSpanTimerName, conf.ObjectiveSpanTimerName, ret.TraceClient, log)
	if err != nil {
		return ret, err
	}
	ret.spanSinks = append(ret.spanSinks, metricSink)

	for _, addrStr := range conf.StatsdListenAddresses {
		addr, err := protocol.ResolveAddr(addrStr)
		if err != nil {
			return ret, err
		}
		ret.StatsdListenAddrs = append(ret.StatsdListenAddrs, addr)
	}
	for _, addrStr := range conf.SsfListenAddresses {
		addr, err := protocol.ResolveAddr(addrStr)
		if err != nil {
			return ret, err
		}
		ret.SSFListenAddrs = append(ret.SSFListenAddrs, addr)
	}

	ret.metricMaxLength = conf.MetricMaxLength
	ret.traceMaxLengthBytes = conf.TraceMaxLengthBytes
	ret.RcvbufBytes = conf.ReadBufferSizeBytes
	ret.HTTPAddr = conf.HTTPAddress
	ret.numListeningHTTP = new(int32)
	ret.ForwardAddr = conf.ForwardAddress

	if conf.TLSKey != "" {
		if conf.TLSCertificate == "" {
			err = errors.New("tls_key is set; must set tls_certificate")
			logger.WithError(err).Error("Improper TLS configuration")
			return ret, err
		}

		// load the TLS key and certificate
		var cert tls.Certificate
		cert, err = tls.X509KeyPair([]byte(conf.TLSCertificate), []byte(conf.TLSKey))
		if err != nil {
			logger.WithError(err).Error("Improper TLS configuration")
			return ret, err
		}

		clientAuthMode := tls.NoClientCert
		var clientCAs *x509.CertPool
		if conf.TLSAuthorityCertificate != "" {
			// load the authority; require clients to present certificated signed by this authority
			clientAuthMode = tls.RequireAndVerifyClientCert
			clientCAs = x509.NewCertPool()
			ok := clientCAs.AppendCertsFromPEM([]byte(conf.TLSAuthorityCertificate))
			if !ok {
				err = errors.New("tls_authority_certificate: Could not load any certificates")
				logger.WithError(err).Error("Improper TLS configuration")
				return ret, err
			}
		}

		ret.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   clientAuthMode,
			ClientCAs:    clientCAs,
		}
	}

	if conf.SignalfxAPIKey != "" {
		tracedHTTP := *ret.HTTPClient
		tracedHTTP.Transport = vhttp.NewTraceRoundTripper(tracedHTTP.Transport, ret.TraceClient, "signalfx")

		fallback := signalfx.NewClient(conf.SignalfxEndpointBase, conf.SignalfxAPIKey, &tracedHTTP)
		byTagClients := map[string]signalfx.DPClient{}
		for _, perTag := range conf.SignalfxPerTagAPIKeys {
			byTagClients[perTag.Name] = signalfx.NewClient(conf.SignalfxEndpointBase, perTag.APIKey, &tracedHTTP)
		}
		sfxSink, err := signalfx.NewSignalFxSink(conf.SignalfxHostnameTag, conf.Hostname, ret.TagsAsMap, log, fallback, conf.SignalfxVaryKeyBy, byTagClients, conf.SignalfxMetricNamePrefixDrops, conf.SignalfxMetricTagPrefixDrops, metricSink, conf.SignalfxFlushMaxPerBody)
		if err != nil {
			return ret, err
		}
		ret.metricSinks = append(ret.metricSinks, sfxSink)
	}
	if conf.DatadogAPIKey != "" && conf.DatadogAPIHostname != "" {
		ddSink, err := datadog.NewDatadogMetricSink(
			ret.interval.Seconds(), conf.DatadogFlushMaxPerBody, conf.Hostname, ret.Tags,
			conf.DatadogAPIHostname, conf.DatadogAPIKey, ret.HTTPClient, log,
		)
		if err != nil {
			return ret, err
		}
		ret.metricSinks = append(ret.metricSinks, ddSink)
	}

	// Configure tracing sinks
	if len(conf.SsfListenAddresses) > 0 {

		trace.Enable()

		// configure Datadog as a Span sink
		if conf.DatadogAPIKey != "" && conf.DatadogTraceAPIAddress != "" {
			ddSink, err := datadog.NewDatadogSpanSink(
				conf.DatadogTraceAPIAddress, conf.DatadogSpanBufferSize,
				ret.HTTPClient, log,
			)
			if err != nil {
				return ret, err
			}

			ret.spanSinks = append(ret.spanSinks, ddSink)
			logger.Info("Configured Datadog span sink")
		}

		if conf.XrayAddress != "" {
			if conf.XraySamplePercentage == 0 {
				log.Warn("XRay sample percentage is 0, no segments will be sent.")
			} else {

				annotationTags := make([]string, 0, len(conf.XrayAnnotationTags))
				for _, tag := range conf.XrayAnnotationTags {
					annotationTags = append(annotationTags, strings.Split(tag, ":")[0])
				}

				xraySink, err := xray.NewXRaySpanSink(conf.XrayAddress, conf.XraySamplePercentage, ret.TagsAsMap, annotationTags, log)
				if err != nil {
					return ret, err
				}
				ret.spanSinks = append(ret.spanSinks, xraySink)

				logger.WithFields(logrus.Fields{
					"sample_percentage":   conf.XraySamplePercentage,
					"num_annotation_tags": annotationTags,
				}).Info("Configured X-Ray span sink")
			}
		}

		// configure Lightstep as a Span Sink
		if conf.LightstepAccessToken != "" {

			var lsSink sinks.SpanSink
			lsSink, err = lightstep.NewLightStepSpanSink(
				conf.LightstepCollectorHost, conf.LightstepReconnectPeriod,
				conf.LightstepMaximumSpans, conf.LightstepNumClients,
				conf.LightstepAccessToken, log,
			)
			if err != nil {
				return ret, err
			}
			ret.spanSinks = append(ret.spanSinks, lsSink)

			logger.Info("Configured Lightstep span sink")
		}

		if (conf.SplunkHecToken != "" && conf.SplunkHecAddress == "") ||
			(conf.SplunkHecToken == "" && conf.SplunkHecAddress != "") {
			return ret, fmt.Errorf("both splunk_hec_address and splunk_hec_token need to be set!")
		}
		if conf.SplunkHecToken != "" && conf.SplunkHecAddress != "" {
			var sendTimeout, ingestTimeout, connLifetime, connJitter time.Duration
			if conf.SplunkHecSendTimeout != "" {
				sendTimeout, err = time.ParseDuration(conf.SplunkHecSendTimeout)
				if err != nil {
					return ret, err
				}
			}
			if conf.SplunkHecIngestTimeout != "" {
				ingestTimeout, err = time.ParseDuration(conf.SplunkHecIngestTimeout)
				if err != nil {
					return ret, err
				}
			}
			if conf.SplunkHecMaxConnectionLifetime != "" {
				connLifetime, err = time.ParseDuration(conf.SplunkHecMaxConnectionLifetime)
				if err != nil {
					return ret, err
				}
			}
			if conf.SplunkHecConnectionLifetimeJitter != "" {
				connJitter, err = time.ParseDuration(conf.SplunkHecConnectionLifetimeJitter)
				if err != nil {
					return ret, err
				}
			}

			sss, err := splunk.NewSplunkSpanSink(conf.SplunkHecAddress, conf.SplunkHecToken, conf.Hostname, conf.SplunkHecTLSValidateHostname, log, ingestTimeout, sendTimeout, conf.SplunkHecBatchSize, conf.SplunkHecSubmissionWorkers, conf.SplunkSpanSampleRate, connLifetime, connJitter)
			if err != nil {
				return ret, err
			}

			ret.spanSinks = append(ret.spanSinks, sss)
		}

		if conf.FalconerAddress != "" {
			falsink, err := falconer.NewSpanSink(context.Background(), conf.FalconerAddress, log, grpc.WithInsecure())
			if err != nil {
				return ret, err
			}

			ret.spanSinks = append(ret.spanSinks, falsink)
			logger.Info("Configured Falconer trace sink")
		}

		// Set up as many span workers as we need:
		ret.SpanWorkerGoroutines = 1
		if conf.NumSpanWorkers > 0 {
			ret.SpanWorkerGoroutines = conf.NumSpanWorkers
		}
	}

	if conf.KafkaBroker != "" {
		if conf.KafkaMetricTopic != "" || conf.KafkaCheckTopic != "" || conf.KafkaEventTopic != "" {
			kSink, err := kafka.NewKafkaMetricSink(
				log, ret.TraceClient, conf.KafkaBroker, conf.KafkaCheckTopic, conf.KafkaEventTopic,
				conf.KafkaMetricTopic, conf.KafkaMetricRequireAcks,
				conf.KafkaPartitioner, conf.KafkaRetryMax,
				conf.KafkaMetricBufferBytes, conf.KafkaMetricBufferMessages,
				conf.KafkaMetricBufferFrequency,
			)
			if err != nil {
				return ret, err
			}

			ret.metricSinks = append(ret.metricSinks, kSink)

			logger.Info("Configured Kafka metric sink")
		} else {
			logger.Warn("Kafka metric sink skipped due to missing metric, check and event topic")
		}

		if conf.KafkaSpanTopic != "" {
			sink, err := kafka.NewKafkaSpanSink(log, ret.TraceClient, conf.KafkaBroker, conf.KafkaSpanTopic,
				conf.KafkaPartitioner, conf.KafkaMetricRequireAcks, conf.KafkaRetryMax,
				conf.KafkaSpanBufferBytes, conf.KafkaSpanBufferMesages,
				conf.KafkaSpanBufferFrequency, conf.KafkaSpanSerializationFormat,
				conf.KafkaSpanSampleTag, conf.KafkaSpanSampleRatePercent,
			)
			if err != nil {
				return ret, err
			}

			ret.spanSinks = append(ret.spanSinks, sink)
			logger.Info("Configured Kafka span sink")
		} else {
			logger.Warn("Kafka span sink skipped due to missing span topic")
		}
	}

	{
		mtx := sync.Mutex{}
		if conf.DebugFlushedMetrics {
			ret.metricSinks = append(ret.metricSinks, debug.NewDebugMetricSink(&mtx, log))
		}
		if conf.DebugIngestedSpans {
			blackhole := debug.NewDebugSpanSink(&mtx, log)
			ret.spanSinks = append(ret.spanSinks, blackhole)
			logger.WithField("name", blackhole.Name()).Info("Starting logger debug sink")
		}
	}

	// After all sinks are initialized, set the list of tags to exclude
	setSinkExcludedTags(conf.TagsExclude, ret.metricSinks)

	var svc s3iface.S3API
	awsID := conf.AwsAccessKeyID
	awsSecret := conf.AwsSecretAccessKey
	if conf.AwsS3Bucket != "" {
		if len(awsID) > 0 && len(awsSecret) > 0 {
			sess, err := session.NewSession(&aws.Config{
				Region:      aws.String(conf.AwsRegion),
				Credentials: credentials.NewStaticCredentials(awsID, awsSecret, ""),
			})

			if err != nil {
				logger.Infof("error getting AWS session: %s", err)
				svc = nil
			} else {
				logger.Info("Successfully created AWS session")
				svc = s3.New(sess)
				plugin := &s3p.S3Plugin{
					Logger:   log,
					Svc:      svc,
					S3Bucket: conf.AwsS3Bucket,
					Hostname: ret.Hostname,
				}
				ret.registerPlugin(plugin)
			}
		} else {
			logger.Info("AWS S3 credentials not found. S3 plugin is disabled.")
		}
	} else {
		logger.Info("AWS S3 bucket not set. Skipping S3 Plugin initialization.")
	}

	if svc == nil {
		logger.Info("S3 archives are disabled")
	} else {
		logger.Info("S3 archives are enabled")
	}

	if conf.FlushFile != "" {
		localFilePlugin := &localfilep.Plugin{
			FilePath: conf.FlushFile,
			Logger:   log,
		}
		ret.registerPlugin(localFilePlugin)
		logger.Info(fmt.Sprintf("Local file logging to %s", conf.FlushFile))
	}

	// closed in Shutdown; Same approach and http.Shutdown
	ret.shutdown = make(chan struct{})

	// Don't emit keys into logs now that we're done with them.
	conf.SentryDsn = REDACTED
	conf.TLSKey = REDACTED
	conf.DatadogAPIKey = REDACTED
	conf.SignalfxAPIKey = REDACTED
	conf.LightstepAccessToken = REDACTED
	conf.AwsAccessKeyID = REDACTED
	conf.AwsSecretAccessKey = REDACTED

	ret.forwardUseGRPC = conf.ForwardUseGrpc

	// Setup the grpc server if it was configured
	ret.grpcListenAddress = conf.GrpcAddress
	if ret.grpcListenAddress != "" {
		// convert all the workers to the proper interface
		ingesters := make([]importsrv.MetricIngester, len(ret.Workers))
		for i, worker := range ret.Workers {
			ingesters[i] = worker
		}

		ret.grpcServer = importsrv.New(ingesters,
			importsrv.WithTraceClient(ret.TraceClient))
	}

	logger.WithField("config", conf).Debug("Initialized server")

	return ret, err
}
