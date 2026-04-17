func (api objectAPIHandlers) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "CopyObject")

	defer logger.AuditLog(w, r, "CopyObject", mustGetClaimsFromToken(r))

	objectAPI := api.ObjectAPI()
	if objectAPI == nil {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrServerNotInitialized), r.URL, guessIsBrowserReq(r))
		return
	}
	if crypto.S3KMS.IsRequested(r.Header) {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrNotImplemented), r.URL, guessIsBrowserReq(r)) // SSE-KMS is not supported
		return
	}
	if !api.EncryptionEnabled() && (hasServerSideEncryptionHeader(r.Header) || crypto.SSECopy.IsRequested(r.Header)) {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrNotImplemented), r.URL, guessIsBrowserReq(r))
		return
	}
	vars := mux.Vars(r)
	dstBucket := vars["bucket"]
	dstObject := vars["object"]

	if s3Error := checkRequestAuthType(ctx, r, policy.PutObjectAction, dstBucket, dstObject); s3Error != ErrNone {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(s3Error), r.URL, guessIsBrowserReq(r))
		return
	}

	// TODO: Reject requests where body/payload is present, for now we don't even read it.

	// Read escaped copy source path to check for parameters.
	cpSrcPath := r.Header.Get("X-Amz-Copy-Source")

	// Check https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectVersioning.html
	// Regardless of whether you have enabled versioning, each object in your bucket
	// has a version ID. If you have not enabled versioning, Amazon S3 sets the value
	// of the version ID to null. If you have enabled versioning, Amazon S3 assigns a
	// unique version ID value for the object.
	if u, err := url.Parse(cpSrcPath); err == nil {
		// Check if versionId query param was added, if yes then check if
		// its non "null" value, we should error out since we do not support
		// any versions other than "null".
		if vid := u.Query().Get("versionId"); vid != "" && vid != "null" {
			writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrNoSuchVersion), r.URL, guessIsBrowserReq(r))
			return
		}
		// Note that url.Parse does the unescaping
		cpSrcPath = u.Path
	}
	if vid := r.Header.Get("X-Amz-Copy-Source-Version-Id"); vid != "" {
		// Check if versionId header was added, if yes then check if
		// its non "null" value, we should error out since we do not support
		// any versions other than "null".
		if vid != "null" {
			writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrNoSuchVersion), r.URL, guessIsBrowserReq(r))
			return
		}
	}

	srcBucket, srcObject := path2BucketAndObject(cpSrcPath)
	// If source object is empty or bucket is empty, reply back invalid copy source.
	if srcObject == "" || srcBucket == "" {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrInvalidCopySource), r.URL, guessIsBrowserReq(r))
		return
	}

	if s3Error := checkRequestAuthType(ctx, r, policy.GetObjectAction, srcBucket, srcObject); s3Error != ErrNone {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(s3Error), r.URL, guessIsBrowserReq(r))
		return
	}

	// Check if metadata directive is valid.
	if !isMetadataDirectiveValid(r.Header) {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrInvalidMetadataDirective), r.URL, guessIsBrowserReq(r))
		return
	}
	// This request header needs to be set prior to setting ObjectOptions
	if globalAutoEncryption && !crypto.SSEC.IsRequested(r.Header) {
		r.Header.Add(crypto.SSEHeader, crypto.SSEAlgorithmAES256)
	}

	var srcOpts, dstOpts ObjectOptions
	srcOpts, err := copySrcOpts(ctx, r, srcBucket, srcObject)
	if err != nil {
		logger.LogIf(ctx, err)
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
		return
	}
	// convert copy src encryption options for GET calls
	var getOpts = ObjectOptions{}
	getSSE := encrypt.SSE(srcOpts.ServerSideEncryption)
	if getSSE != srcOpts.ServerSideEncryption {
		getOpts.ServerSideEncryption = getSSE
	}
	dstOpts, err = copyDstOpts(ctx, r, dstBucket, dstObject, nil)
	if err != nil {
		logger.LogIf(ctx, err)
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
		return
	}

	cpSrcDstSame := isStringEqual(pathJoin(srcBucket, srcObject), pathJoin(dstBucket, dstObject))

	// Deny if WORM is enabled. If operation is key rotation of SSE-S3 encrypted object
	// allow the operation
	if globalWORMEnabled && !(cpSrcDstSame && crypto.S3.IsRequested(r.Header)) {
		if _, err = objectAPI.GetObjectInfo(ctx, dstBucket, dstObject, dstOpts); err == nil {
			writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrMethodNotAllowed), r.URL, guessIsBrowserReq(r))
			return
		}
	}

	getObjectNInfo := objectAPI.GetObjectNInfo
	if api.CacheAPI() != nil {
		getObjectNInfo = api.CacheAPI().GetObjectNInfo
	}

	var lock = noLock
	if !cpSrcDstSame {
		lock = readLock
	}
	checkCopyPrecondFn := func(o ObjectInfo, encETag string) bool {
		return checkCopyObjectPreconditions(ctx, w, r, o, encETag)
	}
	getOpts.CheckCopyPrecondFn = checkCopyPrecondFn
	srcOpts.CheckCopyPrecondFn = checkCopyPrecondFn
	var rs *HTTPRangeSpec
	gr, err := getObjectNInfo(ctx, srcBucket, srcObject, rs, r.Header, lock, getOpts)
	if err != nil {
		if isErrPreconditionFailed(err) {
			return
		}
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
		return
	}
	defer gr.Close()
	srcInfo := gr.ObjInfo

	/// maximum Upload size for object in a single CopyObject operation.
	if isMaxObjectSize(srcInfo.Size) {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrEntityTooLarge), r.URL, guessIsBrowserReq(r))
		return
	}

	// Deny if WORM is enabled, and it is not a SSE-S3 -> SSE-S3 key rotation or if metadata replacement is requested.
	if globalWORMEnabled && cpSrcDstSame && (!crypto.S3.IsEncrypted(srcInfo.UserDefined) || isMetadataReplace(r.Header)) {
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrMethodNotAllowed), r.URL, guessIsBrowserReq(r))
		return
	}
	// We have to copy metadata only if source and destination are same.
	// this changes for encryption which can be observed below.
	if cpSrcDstSame {
		srcInfo.metadataOnly = true
	}

	var reader io.Reader
	var length = srcInfo.Size

	// Set the actual size to the decrypted size if encrypted.
	actualSize := srcInfo.Size
	if crypto.IsEncrypted(srcInfo.UserDefined) {
		actualSize, err = srcInfo.DecryptedSize()
		if err != nil {
			writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
			return
		}
		length = actualSize
	}

	// Check if the destination bucket is on a remote site, this code only gets executed
	// when federation is enabled, ie when globalDNSConfig is non 'nil'.
	//
	// This function is similar to isRemoteCallRequired but specifically for COPY object API
	// if destination and source are same we do not need to check for destnation bucket
	// to exist locally.
	var isRemoteCopyRequired = func(ctx context.Context, srcBucket, dstBucket string, objAPI ObjectLayer) bool {
		if globalDNSConfig == nil {
			return false
		}
		if srcBucket == dstBucket {
			return false
		}
		_, err := objAPI.GetBucketInfo(ctx, dstBucket)
		return err == toObjectErr(errVolumeNotFound, dstBucket)
	}

	var compressMetadata map[string]string
	// No need to compress for remote etcd calls
	// Pass the decompressed stream to such calls.
	isCompressed := objectAPI.IsCompressionSupported() && isCompressible(r.Header, srcObject) && !isRemoteCopyRequired(ctx, srcBucket, dstBucket, objectAPI)
	if isCompressed {
		compressMetadata = make(map[string]string, 2)
		// Preserving the compression metadata.
		compressMetadata[ReservedMetadataPrefix+"compression"] = compressionAlgorithmV1
		compressMetadata[ReservedMetadataPrefix+"actual-size"] = strconv.FormatInt(actualSize, 10)
		// Remove all source encrypted related metadata to
		// avoid copying them in target object.
		crypto.RemoveInternalEntries(srcInfo.UserDefined)

		reader = newSnappyCompressReader(gr)
		length = -1
	} else {
		// Remove the metadata for remote calls.
		delete(srcInfo.UserDefined, ReservedMetadataPrefix+"compression")
		delete(srcInfo.UserDefined, ReservedMetadataPrefix+"actual-size")
		reader = gr
	}

	srcInfo.Reader, err = hash.NewReader(reader, length, "", "", actualSize)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
		return
	}

	rawReader := srcInfo.Reader
	pReader := NewPutObjReader(srcInfo.Reader, nil, nil)

	var encMetadata = make(map[string]string)
	if objectAPI.IsEncryptionSupported() && !isCompressed {
		// Encryption parameters not applicable for this object.
		if !crypto.IsEncrypted(srcInfo.UserDefined) && crypto.SSECopy.IsRequested(r.Header) {
			writeErrorResponse(ctx, w, toAPIError(ctx, errInvalidEncryptionParameters), r.URL, guessIsBrowserReq(r))
			return
		}
		// Encryption parameters not present for this object.
		if crypto.SSEC.IsEncrypted(srcInfo.UserDefined) && !crypto.SSECopy.IsRequested(r.Header) {
			writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrInvalidSSECustomerAlgorithm), r.URL, guessIsBrowserReq(r))
			return
		}

		var oldKey, newKey, objEncKey []byte
		sseCopyS3 := crypto.S3.IsEncrypted(srcInfo.UserDefined)
		sseCopyC := crypto.SSEC.IsEncrypted(srcInfo.UserDefined) && crypto.SSECopy.IsRequested(r.Header)
		sseC := crypto.SSEC.IsRequested(r.Header)
		sseS3 := crypto.S3.IsRequested(r.Header)

		isSourceEncrypted := sseCopyC || sseCopyS3
		isTargetEncrypted := sseC || sseS3

		if sseC {
			newKey, err = ParseSSECustomerRequest(r)
			if err != nil {
				writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
				return
			}
		}

		// If src == dst and either
		// - the object is encrypted using SSE-C and two different SSE-C keys are present
		// - the object is encrypted using SSE-S3 and the SSE-S3 header is present
		// than execute a key rotation.
		var keyRotation bool
		if cpSrcDstSame && ((sseCopyC && sseC) || (sseS3 && sseCopyS3)) {
			if sseCopyC && sseC {
				oldKey, err = ParseSSECopyCustomerRequest(r.Header, srcInfo.UserDefined)
				if err != nil {
					writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
					return
				}
			}

			for k, v := range srcInfo.UserDefined {
				if hasPrefix(k, ReservedMetadataPrefix) {
					encMetadata[k] = v
				}
			}

			// In case of SSE-S3 oldKey and newKey aren't used - the KMS manages the keys.
			if err = rotateKey(oldKey, newKey, srcBucket, srcObject, encMetadata); err != nil {
				writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
				return
			}

			// Since we are rotating the keys, make sure to update the metadata.
			srcInfo.metadataOnly = true
			keyRotation = true
		} else {
			if isSourceEncrypted || isTargetEncrypted {
				// We are not only copying just metadata instead
				// we are creating a new object at this point, even
				// if source and destination are same objects.
				if !keyRotation {
					srcInfo.metadataOnly = false
				}
			}

			// Calculate the size of the target object
			var targetSize int64

			switch {
			case !isSourceEncrypted && !isTargetEncrypted:
				targetSize = srcInfo.Size
			case isSourceEncrypted && isTargetEncrypted:
				objInfo := ObjectInfo{Size: actualSize}
				targetSize = objInfo.EncryptedSize()
			case !isSourceEncrypted && isTargetEncrypted:
				targetSize = srcInfo.EncryptedSize()
			case isSourceEncrypted && !isTargetEncrypted:
				targetSize, _ = srcInfo.DecryptedSize()
			}

			if isTargetEncrypted {
				reader, objEncKey, err = newEncryptReader(srcInfo.Reader, newKey, dstBucket, dstObject, encMetadata, sseS3)
				if err != nil {
					writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
					return
				}
			}

			if isSourceEncrypted {
				// Remove all source encrypted related metadata to
				// avoid copying them in target object.
				crypto.RemoveInternalEntries(srcInfo.UserDefined)
			}

			srcInfo.Reader, err = hash.NewReader(reader, targetSize, "", "", targetSize) // do not try to verify encrypted content
			if err != nil {
				writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
				return
			}

			pReader = NewPutObjReader(rawReader, srcInfo.Reader, objEncKey)
		}
	}

	srcInfo.PutObjReader = pReader

	srcInfo.UserDefined, err = getCpObjMetadataFromHeader(ctx, r, srcInfo.UserDefined)
	if err != nil {
		writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
		return
	}

	// Store the preserved compression metadata.
	for k, v := range compressMetadata {
		srcInfo.UserDefined[k] = v
	}

	// We need to preserve the encryption headers set in EncryptRequest,
	// so we do not want to override them, copy them instead.
	for k, v := range encMetadata {
		srcInfo.UserDefined[k] = v
	}

	// Ensure that metadata does not contain sensitive information
	crypto.RemoveSensitiveEntries(srcInfo.UserDefined)
	// Check if x-amz-metadata-directive was not set to REPLACE and source,
	// desination are same objects. Apply this restriction also when
	// metadataOnly is true indicating that we are not overwriting the object.
	// if encryption is enabled we do not need explicit "REPLACE" metadata to
	// be enabled as well - this is to allow for key-rotation.
	if !isMetadataReplace(r.Header) && srcInfo.metadataOnly && !crypto.IsEncrypted(srcInfo.UserDefined) {
		// If x-amz-metadata-directive is not set to REPLACE then we need
		// to error out if source and destination are same.
		writeErrorResponse(ctx, w, errorCodes.ToAPIErr(ErrInvalidCopyDest), r.URL, guessIsBrowserReq(r))
		return
	}

	var objInfo ObjectInfo

	if isRemoteCopyRequired(ctx, srcBucket, dstBucket, objectAPI) {
		var dstRecords []dns.SrvRecord
		dstRecords, err = globalDNSConfig.Get(dstBucket)
		if err != nil {
			writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
			return
		}

		// Send PutObject request to appropriate instance (in federated deployment)
		client, rerr := getRemoteInstanceClient(r, getHostFromSrv(dstRecords))
		if rerr != nil {
			writeErrorResponse(ctx, w, toAPIError(ctx, rerr), r.URL, guessIsBrowserReq(r))
			return
		}
		remoteObjInfo, rerr := client.PutObject(dstBucket, dstObject, srcInfo.Reader,
			srcInfo.Size, "", "", srcInfo.UserDefined, dstOpts.ServerSideEncryption)
		if rerr != nil {
			writeErrorResponse(ctx, w, toAPIError(ctx, rerr), r.URL, guessIsBrowserReq(r))
			return
		}
		objInfo.ETag = remoteObjInfo.ETag
		objInfo.ModTime = remoteObjInfo.LastModified
	} else {
		// Copy source object to destination, if source and destination
		// object is same then only metadata is updated.
		objInfo, err = objectAPI.CopyObject(ctx, srcBucket, srcObject, dstBucket, dstObject, srcInfo, srcOpts, dstOpts)
		if err != nil {
			writeErrorResponse(ctx, w, toAPIError(ctx, err), r.URL, guessIsBrowserReq(r))
			return
		}
	}

	response := generateCopyObjectResponse(getDecryptedETag(r.Header, objInfo, false), objInfo.ModTime)
	encodedSuccessResponse := encodeResponse(response)

	// Write success response.
	writeSuccessResponseXML(w, encodedSuccessResponse)

	if objInfo.IsCompressed() {
		objInfo.Size = actualSize
	}

	// Notify object created event.
	sendEvent(eventArgs{
		EventName:    event.ObjectCreatedCopy,
		BucketName:   dstBucket,
		Object:       objInfo,
		ReqParams:    extractReqParams(r),
		RespElements: extractRespElements(w),
		UserAgent:    r.UserAgent(),
		Host:         handlers.GetSourceIP(r),
	})
}
