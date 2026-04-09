func (c *Config) Client() (interface{}, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if !c.SkipRegionValidation {
		if err := awsbase.ValidateRegion(c.Region); err != nil {
			return nil, err
		}
	}

	log.Println("[INFO] Building AWS auth structure")
	awsbaseConfig := &awsbase.Config{
		AccessKey:               c.AccessKey,
		AssumeRoleARN:           c.AssumeRoleARN,
		AssumeRoleExternalID:    c.AssumeRoleExternalID,
		AssumeRolePolicy:        c.AssumeRolePolicy,
		AssumeRoleSessionName:   c.AssumeRoleSessionName,
		CredsFilename:           c.CredsFilename,
		DebugLogging:            logging.IsDebugOrHigher(),
		IamEndpoint:             c.Endpoints["iam"],
		Insecure:                c.Insecure,
		MaxRetries:              c.MaxRetries,
		Profile:                 c.Profile,
		Region:                  c.Region,
		SecretKey:               c.SecretKey,
		SkipCredsValidation:     c.SkipCredsValidation,
		SkipMetadataApiCheck:    c.SkipMetadataApiCheck,
		SkipRequestingAccountId: c.SkipRequestingAccountId,
		StsEndpoint:             c.Endpoints["sts"],
		Token:                   c.Token,
		UserAgentProducts: []*awsbase.UserAgentProduct{
			{Name: "APN", Version: "1.0"},
			{Name: "HashiCorp", Version: "1.0"},
			{Name: "Terraform", Version: terraform.VersionString()},
		},
	}

	sess, accountID, partition, err := awsbase.GetSessionWithAccountIDAndPartition(awsbaseConfig)
	if err != nil {
		return nil, err
	}

	if accountID == "" {
		log.Printf("[WARN] AWS account ID not found for provider. See https://www.terraform.io/docs/providers/aws/index.html#skip_requesting_account_id for implications.")
	}

	if err := awsbase.ValidateAccountID(accountID, c.AllowedAccountIds, c.ForbiddenAccountIds); err != nil {
		return nil, err
	}

	client := &AWSClient{
		accountid:                           accountID,
		acmconn:                             acm.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["acm"])})),
		acmpcaconn:                          acmpca.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["acmpca"])})),
		apigateway:                          apigateway.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["apigateway"])})),
		apigatewayv2conn:                    apigatewayv2.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["apigateway"])})),
		appautoscalingconn:                  applicationautoscaling.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["applicationautoscaling"])})),
		appmeshconn:                         appmesh.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["appmesh"])})),
		appsyncconn:                         appsync.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["appsync"])})),
		athenaconn:                          athena.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["athena"])})),
		autoscalingconn:                     autoscaling.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["autoscaling"])})),
		backupconn:                          backup.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["backup"])})),
		batchconn:                           batch.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["batch"])})),
		budgetconn:                          budgets.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["budgets"])})),
		cfconn:                              cloudformation.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudformation"])})),
		cloud9conn:                          cloud9.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloud9"])})),
		cloudfrontconn:                      cloudfront.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudfront"])})),
		cloudhsmv2conn:                      cloudhsmv2.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudhsm"])})),
		cloudsearchconn:                     cloudsearch.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudsearch"])})),
		cloudtrailconn:                      cloudtrail.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudtrail"])})),
		cloudwatchconn:                      cloudwatch.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudwatch"])})),
		cloudwatcheventsconn:                cloudwatchevents.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudwatchevents"])})),
		cloudwatchlogsconn:                  cloudwatchlogs.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cloudwatchlogs"])})),
		codebuildconn:                       codebuild.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["codebuild"])})),
		codecommitconn:                      codecommit.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["codecommit"])})),
		codedeployconn:                      codedeploy.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["codedeploy"])})),
		codepipelineconn:                    codepipeline.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["codepipeline"])})),
		cognitoconn:                         cognitoidentity.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cognitoidentity"])})),
		cognitoidpconn:                      cognitoidentityprovider.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cognitoidp"])})),
		configconn:                          configservice.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["configservice"])})),
		costandusagereportconn:              costandusagereportservice.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["cur"])})),
		datapipelineconn:                    datapipeline.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["datapipeline"])})),
		datasyncconn:                        datasync.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["datasync"])})),
		daxconn:                             dax.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["dax"])})),
		devicefarmconn:                      devicefarm.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["devicefarm"])})),
		dlmconn:                             dlm.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["dlm"])})),
		dmsconn:                             databasemigrationservice.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["dms"])})),
		docdbconn:                           docdb.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["docdb"])})),
		dsconn:                              directoryservice.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ds"])})),
		dxconn:                              directconnect.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["directconnect"])})),
		dynamodbconn:                        dynamodb.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["dynamodb"])})),
		ec2conn:                             ec2.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ec2"])})),
		ecrconn:                             ecr.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ecr"])})),
		ecsconn:                             ecs.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ecs"])})),
		efsconn:                             efs.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["efs"])})),
		eksconn:                             eks.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["eks"])})),
		elasticacheconn:                     elasticache.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["elasticache"])})),
		elasticbeanstalkconn:                elasticbeanstalk.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["elasticbeanstalk"])})),
		elastictranscoderconn:               elastictranscoder.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["elastictranscoder"])})),
		elbconn:                             elb.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["elb"])})),
		elbv2conn:                           elbv2.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["elb"])})),
		emrconn:                             emr.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["emr"])})),
		esconn:                              elasticsearch.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["es"])})),
		firehoseconn:                        firehose.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["firehose"])})),
		fmsconn:                             fms.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["fms"])})),
		fsxconn:                             fsx.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["fsx"])})),
		gameliftconn:                        gamelift.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["gamelift"])})),
		glacierconn:                         glacier.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["glacier"])})),
		globalacceleratorconn:               globalaccelerator.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["globalaccelerator"])})),
		glueconn:                            glue.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["glue"])})),
		guarddutyconn:                       guardduty.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["guardduty"])})),
		iamconn:                             iam.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["iam"])})),
		inspectorconn:                       inspector.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["inspector"])})),
		iotconn:                             iot.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["iot"])})),
		kafkaconn:                           kafka.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kafka"])})),
		kinesisanalyticsconn:                kinesisanalytics.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kinesisanalytics"])})),
		kinesisanalyticsv2conn:              kinesisanalyticsv2.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kinesisanalytics"])})),
		kinesisconn:                         kinesis.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kinesis"])})),
		kinesisvideoconn:                    kinesisvideo.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kinesisvideo"])})),
		kmsconn:                             kms.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kms"])})),
		lambdaconn:                          lambda.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["lambda"])})),
		lexmodelconn:                        lexmodelbuildingservice.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["lexmodels"])})),
		licensemanagerconn:                  licensemanager.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["licensemanager"])})),
		lightsailconn:                       lightsail.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["lightsail"])})),
		macieconn:                           macie.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["macie"])})),
		mediaconnectconn:                    mediaconnect.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mediaconnect"])})),
		mediaconvertconn:                    mediaconvert.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mediaconvert"])})),
		medialiveconn:                       medialive.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["medialive"])})),
		mediapackageconn:                    mediapackage.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mediapackage"])})),
		mediastoreconn:                      mediastore.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mediastore"])})),
		mediastoredataconn:                  mediastoredata.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mediastoredata"])})),
		mqconn:                              mq.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["mq"])})),
		neptuneconn:                         neptune.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["neptune"])})),
		opsworksconn:                        opsworks.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["opsworks"])})),
		organizationsconn:                   organizations.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["organizations"])})),
		partition:                           partition,
		pinpointconn:                        pinpoint.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["pinpoint"])})),
		pricingconn:                         pricing.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["pricing"])})),
		quicksightconn:                      quicksight.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["quicksight"])})),
		r53conn:                             route53.New(sess.Copy(&aws.Config{Region: aws.String("us-east-1"), Endpoint: aws.String(c.Endpoints["route53"])})),
		ramconn:                             ram.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ram"])})),
		rdsconn:                             rds.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["rds"])})),
		redshiftconn:                        redshift.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["redshift"])})),
		region:                              c.Region,
		resourcegroupsconn:                  resourcegroups.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["resourcegroups"])})),
		route53resolverconn:                 route53resolver.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["route53resolver"])})),
		s3conn:                              s3.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["s3"]), S3ForcePathStyle: aws.Bool(c.S3ForcePathStyle)})),
		s3controlconn:                       s3control.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["s3control"])})),
		sagemakerconn:                       sagemaker.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["sagemaker"])})),
		scconn:                              servicecatalog.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["servicecatalog"])})),
		sdconn:                              servicediscovery.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["servicediscovery"])})),
		secretsmanagerconn:                  secretsmanager.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["secretsmanager"])})),
		securityhubconn:                     securityhub.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["securityhub"])})),
		serverlessapplicationrepositoryconn: serverlessapplicationrepository.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["serverlessrepo"])})),
		sesConn:                             ses.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ses"])})),
		sfnconn:                             sfn.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["stepfunctions"])})),
		shieldconn:                          shield.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["shield"])})),
		simpledbconn:                        simpledb.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["sdb"])})),
		snsconn:                             sns.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["sns"])})),
		sqsconn:                             sqs.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["sqs"])})),
		ssmconn:                             ssm.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["ssm"])})),
		storagegatewayconn:                  storagegateway.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["storagegateway"])})),
		stsconn:                             sts.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["sts"])})),
		swfconn:                             swf.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["swf"])})),
		transferconn:                        transfer.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["transfer"])})),
		wafconn:                             waf.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["waf"])})),
		wafregionalconn:                     wafregional.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["wafregional"])})),
		worklinkconn:                        worklink.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["worklink"])})),
		workspacesconn:                      workspaces.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["workspaces"])})),
	}

	// Handle deprecated endpoint configurations
	if c.Endpoints["kinesis_analytics"] != "" {
		client.kinesisanalyticsconn = kinesisanalytics.New(sess.Copy(&aws.Config{Endpoint: aws.String(c.Endpoints["kinesis_analytics"])}))
	}
	if c.Endpoints["r53"] != "" {
		client.r53conn = route53.New(sess.Copy(&aws.Config{Region: aws.String("us-east-1"), Endpoint: aws.String(c.Endpoints["r53"])}))
	}

	// Workaround for https://github.com/aws/aws-sdk-go/issues/1376
	client.kinesisconn.Handlers.Retry.PushBack(func(r *request.Request) {
		if !strings.HasPrefix(r.Operation.Name, "Describe") && !strings.HasPrefix(r.Operation.Name, "List") {
			return
		}
		err, ok := r.Error.(awserr.Error)
		if !ok || err == nil {
			return
		}
		if err.Code() == kinesis.ErrCodeLimitExceededException {
			r.Retryable = aws.Bool(true)
		}
	})

	// Workaround for https://github.com/aws/aws-sdk-go/issues/1472
	client.appautoscalingconn.Handlers.Retry.PushBack(func(r *request.Request) {
		if !strings.HasPrefix(r.Operation.Name, "Describe") && !strings.HasPrefix(r.Operation.Name, "List") {
			return
		}
		err, ok := r.Error.(awserr.Error)
		if !ok || err == nil {
			return
		}
		if err.Code() == applicationautoscaling.ErrCodeFailedResourceAccessException {
			r.Retryable = aws.Bool(true)
		}
	})

	client.appsyncconn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "CreateGraphqlApi" {
			if isAWSErr(r.Error, appsync.ErrCodeConcurrentModificationException, "a GraphQL API creation is already in progress") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	// See https://github.com/aws/aws-sdk-go/pull/1276
	client.dynamodbconn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name != "PutItem" && r.Operation.Name != "UpdateItem" && r.Operation.Name != "DeleteItem" {
			return
		}
		if isAWSErr(r.Error, dynamodb.ErrCodeLimitExceededException, "Subscriber limit exceeded:") {
			r.Retryable = aws.Bool(true)
		}
	})

	client.ec2conn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "CreateVpnConnection" {
			if isAWSErr(r.Error, "VpnConnectionLimitExceeded", "maximum number of mutating objects has been reached") {
				r.Retryable = aws.Bool(true)
			}
		}

		if r.Operation.Name == "CreateVpnGateway" {
			if isAWSErr(r.Error, "VpnGatewayLimitExceeded", "maximum number of mutating objects has been reached") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.kinesisconn.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.Operation.Name == "CreateStream" {
			if isAWSErr(r.Error, kinesis.ErrCodeLimitExceededException, "simultaneously be in CREATING or DELETING") {
				r.Retryable = aws.Bool(true)
			}
		}
		if r.Operation.Name == "CreateStream" || r.Operation.Name == "DeleteStream" {
			if isAWSErr(r.Error, kinesis.ErrCodeLimitExceededException, "Rate exceeded for stream") {
				r.Retryable = aws.Bool(true)
			}
		}
	})

	client.storagegatewayconn.Handlers.Retry.PushBack(func(r *request.Request) {
		// InvalidGatewayRequestException: The specified gateway proxy network connection is busy.
		if isAWSErr(r.Error, storagegateway.ErrCodeInvalidGatewayRequestException, "The specified gateway proxy network connection is busy") {
			r.Retryable = aws.Bool(true)
		}
	})

	if !c.SkipGetEC2Platforms {
		supportedPlatforms, err := GetSupportedEC2Platforms(client.ec2conn)
		if err != nil {
			// We intentionally fail *silently* because there's a chance
			// user just doesn't have ec2:DescribeAccountAttributes permissions
			log.Printf("[WARN] Unable to get supported EC2 platforms: %s", err)
		} else {
			client.supportedplatforms = supportedPlatforms
		}
	}

	return client, nil
}
