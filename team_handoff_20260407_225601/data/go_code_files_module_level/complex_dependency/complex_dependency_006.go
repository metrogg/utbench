func UnmarshalMessage(msg *Message) (interface{}, error) {
	var v easyjson.Unmarshaler
	switch msg.Method {
	case CommandAccessibilityDisable:
		return emptyVal, nil

	case CommandAccessibilityEnable:
		return emptyVal, nil

	case CommandAccessibilityGetPartialAXTree:
		v = new(accessibility.GetPartialAXTreeReturns)

	case CommandAccessibilityGetFullAXTree:
		v = new(accessibility.GetFullAXTreeReturns)

	case CommandAnimationDisable:
		return emptyVal, nil

	case CommandAnimationEnable:
		return emptyVal, nil

	case CommandAnimationGetCurrentTime:
		v = new(animation.GetCurrentTimeReturns)

	case CommandAnimationGetPlaybackRate:
		v = new(animation.GetPlaybackRateReturns)

	case CommandAnimationReleaseAnimations:
		return emptyVal, nil

	case CommandAnimationResolveAnimation:
		v = new(animation.ResolveAnimationReturns)

	case CommandAnimationSeekAnimations:
		return emptyVal, nil

	case CommandAnimationSetPaused:
		return emptyVal, nil

	case CommandAnimationSetPlaybackRate:
		return emptyVal, nil

	case CommandAnimationSetTiming:
		return emptyVal, nil

	case EventAnimationAnimationCanceled:
		v = new(animation.EventAnimationCanceled)

	case EventAnimationAnimationCreated:
		v = new(animation.EventAnimationCreated)

	case EventAnimationAnimationStarted:
		v = new(animation.EventAnimationStarted)

	case CommandApplicationCacheEnable:
		return emptyVal, nil

	case CommandApplicationCacheGetApplicationCacheForFrame:
		v = new(applicationcache.GetApplicationCacheForFrameReturns)

	case CommandApplicationCacheGetFramesWithManifests:
		v = new(applicationcache.GetFramesWithManifestsReturns)

	case CommandApplicationCacheGetManifestForFrame:
		v = new(applicationcache.GetManifestForFrameReturns)

	case EventApplicationCacheApplicationCacheStatusUpdated:
		v = new(applicationcache.EventApplicationCacheStatusUpdated)

	case EventApplicationCacheNetworkStateUpdated:
		v = new(applicationcache.EventNetworkStateUpdated)

	case CommandAuditsGetEncodedResponse:
		v = new(audits.GetEncodedResponseReturns)

	case CommandBackgroundServiceStartObserving:
		return emptyVal, nil

	case CommandBackgroundServiceStopObserving:
		return emptyVal, nil

	case CommandBackgroundServiceSetRecording:
		return emptyVal, nil

	case CommandBackgroundServiceClearEvents:
		return emptyVal, nil

	case EventBackgroundServiceRecordingStateChanged:
		v = new(backgroundservice.EventRecordingStateChanged)

	case EventBackgroundServiceBackgroundServiceEventReceived:
		v = new(backgroundservice.EventBackgroundServiceEventReceived)

	case CommandBrowserGrantPermissions:
		return emptyVal, nil

	case CommandBrowserResetPermissions:
		return emptyVal, nil

	case CommandBrowserClose:
		return emptyVal, nil

	case CommandBrowserCrash:
		return emptyVal, nil

	case CommandBrowserCrashGpuProcess:
		return emptyVal, nil

	case CommandBrowserGetVersion:
		v = new(browser.GetVersionReturns)

	case CommandBrowserGetBrowserCommandLine:
		v = new(browser.GetBrowserCommandLineReturns)

	case CommandBrowserGetHistograms:
		v = new(browser.GetHistogramsReturns)

	case CommandBrowserGetHistogram:
		v = new(browser.GetHistogramReturns)

	case CommandBrowserGetWindowBounds:
		v = new(browser.GetWindowBoundsReturns)

	case CommandBrowserGetWindowForTarget:
		v = new(browser.GetWindowForTargetReturns)

	case CommandBrowserSetWindowBounds:
		return emptyVal, nil

	case CommandBrowserSetDockTile:
		return emptyVal, nil

	case CommandCSSAddRule:
		v = new(css.AddRuleReturns)

	case CommandCSSCollectClassNames:
		v = new(css.CollectClassNamesReturns)

	case CommandCSSCreateStyleSheet:
		v = new(css.CreateStyleSheetReturns)

	case CommandCSSDisable:
		return emptyVal, nil

	case CommandCSSEnable:
		return emptyVal, nil

	case CommandCSSForcePseudoState:
		return emptyVal, nil

	case CommandCSSGetBackgroundColors:
		v = new(css.GetBackgroundColorsReturns)

	case CommandCSSGetComputedStyleForNode:
		v = new(css.GetComputedStyleForNodeReturns)

	case CommandCSSGetInlineStylesForNode:
		v = new(css.GetInlineStylesForNodeReturns)

	case CommandCSSGetMatchedStylesForNode:
		v = new(css.GetMatchedStylesForNodeReturns)

	case CommandCSSGetMediaQueries:
		v = new(css.GetMediaQueriesReturns)

	case CommandCSSGetPlatformFontsForNode:
		v = new(css.GetPlatformFontsForNodeReturns)

	case CommandCSSGetStyleSheetText:
		v = new(css.GetStyleSheetTextReturns)

	case CommandCSSSetEffectivePropertyValueForNode:
		return emptyVal, nil

	case CommandCSSSetKeyframeKey:
		v = new(css.SetKeyframeKeyReturns)

	case CommandCSSSetMediaText:
		v = new(css.SetMediaTextReturns)

	case CommandCSSSetRuleSelector:
		v = new(css.SetRuleSelectorReturns)

	case CommandCSSSetStyleSheetText:
		v = new(css.SetStyleSheetTextReturns)

	case CommandCSSSetStyleTexts:
		v = new(css.SetStyleTextsReturns)

	case CommandCSSStartRuleUsageTracking:
		return emptyVal, nil

	case CommandCSSStopRuleUsageTracking:
		v = new(css.StopRuleUsageTrackingReturns)

	case CommandCSSTakeCoverageDelta:
		v = new(css.TakeCoverageDeltaReturns)

	case EventCSSFontsUpdated:
		v = new(css.EventFontsUpdated)

	case EventCSSMediaQueryResultChanged:
		v = new(css.EventMediaQueryResultChanged)

	case EventCSSStyleSheetAdded:
		v = new(css.EventStyleSheetAdded)

	case EventCSSStyleSheetChanged:
		v = new(css.EventStyleSheetChanged)

	case EventCSSStyleSheetRemoved:
		v = new(css.EventStyleSheetRemoved)

	case CommandCacheStorageDeleteCache:
		return emptyVal, nil

	case CommandCacheStorageDeleteEntry:
		return emptyVal, nil

	case CommandCacheStorageRequestCacheNames:
		v = new(cachestorage.RequestCacheNamesReturns)

	case CommandCacheStorageRequestCachedResponse:
		v = new(cachestorage.RequestCachedResponseReturns)

	case CommandCacheStorageRequestEntries:
		v = new(cachestorage.RequestEntriesReturns)

	case CommandCastEnable:
		return emptyVal, nil

	case CommandCastDisable:
		return emptyVal, nil

	case CommandCastSetSinkToUse:
		return emptyVal, nil

	case CommandCastStartTabMirroring:
		return emptyVal, nil

	case CommandCastStopCasting:
		return emptyVal, nil

	case EventCastSinksUpdated:
		v = new(cast.EventSinksUpdated)

	case EventCastIssueUpdated:
		v = new(cast.EventIssueUpdated)

	case CommandDOMCollectClassNamesFromSubtree:
		v = new(dom.CollectClassNamesFromSubtreeReturns)

	case CommandDOMCopyTo:
		v = new(dom.CopyToReturns)

	case CommandDOMDescribeNode:
		v = new(dom.DescribeNodeReturns)

	case CommandDOMDisable:
		return emptyVal, nil

	case CommandDOMDiscardSearchResults:
		return emptyVal, nil

	case CommandDOMEnable:
		return emptyVal, nil

	case CommandDOMFocus:
		return emptyVal, nil

	case CommandDOMGetAttributes:
		v = new(dom.GetAttributesReturns)

	case CommandDOMGetBoxModel:
		v = new(dom.GetBoxModelReturns)

	case CommandDOMGetContentQuads:
		v = new(dom.GetContentQuadsReturns)

	case CommandDOMGetDocument:
		v = new(dom.GetDocumentReturns)

	case CommandDOMGetFlattenedDocument:
		v = new(dom.GetFlattenedDocumentReturns)

	case CommandDOMGetNodeForLocation:
		v = new(dom.GetNodeForLocationReturns)

	case CommandDOMGetOuterHTML:
		v = new(dom.GetOuterHTMLReturns)

	case CommandDOMGetRelayoutBoundary:
		v = new(dom.GetRelayoutBoundaryReturns)

	case CommandDOMGetSearchResults:
		v = new(dom.GetSearchResultsReturns)

	case CommandDOMMarkUndoableState:
		return emptyVal, nil

	case CommandDOMMoveTo:
		v = new(dom.MoveToReturns)

	case CommandDOMPerformSearch:
		v = new(dom.PerformSearchReturns)

	case CommandDOMPushNodeByPathToFrontend:
		v = new(dom.PushNodeByPathToFrontendReturns)

	case CommandDOMPushNodesByBackendIdsToFrontend:
		v = new(dom.PushNodesByBackendIdsToFrontendReturns)

	case CommandDOMQuerySelector:
		v = new(dom.QuerySelectorReturns)

	case CommandDOMQuerySelectorAll:
		v = new(dom.QuerySelectorAllReturns)

	case CommandDOMRedo:
		return emptyVal, nil

	case CommandDOMRemoveAttribute:
		return emptyVal, nil

	case CommandDOMRemoveNode:
		return emptyVal, nil

	case CommandDOMRequestChildNodes:
		return emptyVal, nil

	case CommandDOMRequestNode:
		v = new(dom.RequestNodeReturns)

	case CommandDOMResolveNode:
		v = new(dom.ResolveNodeReturns)

	case CommandDOMSetAttributeValue:
		return emptyVal, nil

	case CommandDOMSetAttributesAsText:
		return emptyVal, nil

	case CommandDOMSetFileInputFiles:
		return emptyVal, nil

	case CommandDOMGetFileInfo:
		v = new(dom.GetFileInfoReturns)

	case CommandDOMSetInspectedNode:
		return emptyVal, nil

	case CommandDOMSetNodeName:
		v = new(dom.SetNodeNameReturns)

	case CommandDOMSetNodeValue:
		return emptyVal, nil

	case CommandDOMSetOuterHTML:
		return emptyVal, nil

	case CommandDOMUndo:
		return emptyVal, nil

	case CommandDOMGetFrameOwner:
		v = new(dom.GetFrameOwnerReturns)

	case EventDOMAttributeModified:
		v = new(dom.EventAttributeModified)

	case EventDOMAttributeRemoved:
		v = new(dom.EventAttributeRemoved)

	case EventDOMCharacterDataModified:
		v = new(dom.EventCharacterDataModified)

	case EventDOMChildNodeCountUpdated:
		v = new(dom.EventChildNodeCountUpdated)

	case EventDOMChildNodeInserted:
		v = new(dom.EventChildNodeInserted)

	case EventDOMChildNodeRemoved:
		v = new(dom.EventChildNodeRemoved)

	case EventDOMDistributedNodesUpdated:
		v = new(dom.EventDistributedNodesUpdated)

	case EventDOMDocumentUpdated:
		v = new(dom.EventDocumentUpdated)

	case EventDOMInlineStyleInvalidated:
		v = new(dom.EventInlineStyleInvalidated)

	case EventDOMPseudoElementAdded:
		v = new(dom.EventPseudoElementAdded)

	case EventDOMPseudoElementRemoved:
		v = new(dom.EventPseudoElementRemoved)

	case EventDOMSetChildNodes:
		v = new(dom.EventSetChildNodes)

	case EventDOMShadowRootPopped:
		v = new(dom.EventShadowRootPopped)

	case EventDOMShadowRootPushed:
		v = new(dom.EventShadowRootPushed)

	case CommandDOMDebuggerGetEventListeners:
		v = new(domdebugger.GetEventListenersReturns)

	case CommandDOMDebuggerRemoveDOMBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerRemoveEventListenerBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerRemoveInstrumentationBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerRemoveXHRBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerSetDOMBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerSetEventListenerBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerSetInstrumentationBreakpoint:
		return emptyVal, nil

	case CommandDOMDebuggerSetXHRBreakpoint:
		return emptyVal, nil

	case CommandDOMSnapshotDisable:
		return emptyVal, nil

	case CommandDOMSnapshotEnable:
		return emptyVal, nil

	case CommandDOMSnapshotCaptureSnapshot:
		v = new(domsnapshot.CaptureSnapshotReturns)

	case CommandDOMStorageClear:
		return emptyVal, nil

	case CommandDOMStorageDisable:
		return emptyVal, nil

	case CommandDOMStorageEnable:
		return emptyVal, nil

	case CommandDOMStorageGetDOMStorageItems:
		v = new(domstorage.GetDOMStorageItemsReturns)

	case CommandDOMStorageRemoveDOMStorageItem:
		return emptyVal, nil

	case CommandDOMStorageSetDOMStorageItem:
		return emptyVal, nil

	case EventDOMStorageDomStorageItemAdded:
		v = new(domstorage.EventDomStorageItemAdded)

	case EventDOMStorageDomStorageItemRemoved:
		v = new(domstorage.EventDomStorageItemRemoved)

	case EventDOMStorageDomStorageItemUpdated:
		v = new(domstorage.EventDomStorageItemUpdated)

	case EventDOMStorageDomStorageItemsCleared:
		v = new(domstorage.EventDomStorageItemsCleared)

	case CommandDatabaseDisable:
		return emptyVal, nil

	case CommandDatabaseEnable:
		return emptyVal, nil

	case CommandDatabaseExecuteSQL:
		v = new(database.ExecuteSQLReturns)

	case CommandDatabaseGetDatabaseTableNames:
		v = new(database.GetDatabaseTableNamesReturns)

	case EventDatabaseAddDatabase:
		v = new(database.EventAddDatabase)

	case CommandDebuggerContinueToLocation:
		return emptyVal, nil

	case CommandDebuggerDisable:
		return emptyVal, nil

	case CommandDebuggerEnable:
		v = new(debugger.EnableReturns)

	case CommandDebuggerEvaluateOnCallFrame:
		v = new(debugger.EvaluateOnCallFrameReturns)

	case CommandDebuggerGetPossibleBreakpoints:
		v = new(debugger.GetPossibleBreakpointsReturns)

	case CommandDebuggerGetScriptSource:
		v = new(debugger.GetScriptSourceReturns)

	case CommandDebuggerGetStackTrace:
		v = new(debugger.GetStackTraceReturns)

	case CommandDebuggerPause:
		return emptyVal, nil

	case CommandDebuggerPauseOnAsyncCall:
		return emptyVal, nil

	case CommandDebuggerRemoveBreakpoint:
		return emptyVal, nil

	case CommandDebuggerRestartFrame:
		v = new(debugger.RestartFrameReturns)

	case CommandDebuggerResume:
		return emptyVal, nil

	case CommandDebuggerSearchInContent:
		v = new(debugger.SearchInContentReturns)

	case CommandDebuggerSetAsyncCallStackDepth:
		return emptyVal, nil

	case CommandDebuggerSetBlackboxPatterns:
		return emptyVal, nil

	case CommandDebuggerSetBlackboxedRanges:
		return emptyVal, nil

	case CommandDebuggerSetBreakpoint:
		v = new(debugger.SetBreakpointReturns)

	case CommandDebuggerSetBreakpointByURL:
		v = new(debugger.SetBreakpointByURLReturns)

	case CommandDebuggerSetBreakpointOnFunctionCall:
		v = new(debugger.SetBreakpointOnFunctionCallReturns)

	case CommandDebuggerSetBreakpointsActive:
		return emptyVal, nil

	case CommandDebuggerSetPauseOnExceptions:
		return emptyVal, nil

	case CommandDebuggerSetReturnValue:
		return emptyVal, nil

	case CommandDebuggerSetScriptSource:
		v = new(debugger.SetScriptSourceReturns)

	case CommandDebuggerSetSkipAllPauses:
		return emptyVal, nil

	case CommandDebuggerSetVariableValue:
		return emptyVal, nil

	case CommandDebuggerStepInto:
		return emptyVal, nil

	case CommandDebuggerStepOut:
		return emptyVal, nil

	case CommandDebuggerStepOver:
		return emptyVal, nil

	case EventDebuggerBreakpointResolved:
		v = new(debugger.EventBreakpointResolved)

	case EventDebuggerPaused:
		v = new(debugger.EventPaused)

	case EventDebuggerResumed:
		v = new(debugger.EventResumed)

	case EventDebuggerScriptFailedToParse:
		v = new(debugger.EventScriptFailedToParse)

	case EventDebuggerScriptParsed:
		v = new(debugger.EventScriptParsed)

	case CommandDeviceOrientationClearDeviceOrientationOverride:
		return emptyVal, nil

	case CommandDeviceOrientationSetDeviceOrientationOverride:
		return emptyVal, nil

	case CommandEmulationCanEmulate:
		v = new(emulation.CanEmulateReturns)

	case CommandEmulationClearDeviceMetricsOverride:
		return emptyVal, nil

	case CommandEmulationClearGeolocationOverride:
		return emptyVal, nil

	case CommandEmulationResetPageScaleFactor:
		return emptyVal, nil

	case CommandEmulationSetFocusEmulationEnabled:
		return emptyVal, nil

	case CommandEmulationSetCPUThrottlingRate:
		return emptyVal, nil

	case CommandEmulationSetDefaultBackgroundColorOverride:
		return emptyVal, nil

	case CommandEmulationSetDeviceMetricsOverride:
		return emptyVal, nil

	case CommandEmulationSetScrollbarsHidden:
		return emptyVal, nil

	case CommandEmulationSetDocumentCookieDisabled:
		return emptyVal, nil

	case CommandEmulationSetEmitTouchEventsForMouse:
		return emptyVal, nil

	case CommandEmulationSetEmulatedMedia:
		return emptyVal, nil

	case CommandEmulationSetGeolocationOverride:
		return emptyVal, nil

	case CommandEmulationSetPageScaleFactor:
		return emptyVal, nil

	case CommandEmulationSetScriptExecutionDisabled:
		return emptyVal, nil

	case CommandEmulationSetTouchEmulationEnabled:
		return emptyVal, nil

	case CommandEmulationSetVirtualTimePolicy:
		v = new(emulation.SetVirtualTimePolicyReturns)

	case CommandEmulationSetUserAgentOverride:
		return emptyVal, nil

	case EventEmulationVirtualTimeBudgetExpired:
		v = new(emulation.EventVirtualTimeBudgetExpired)

	case CommandFetchDisable:
		return emptyVal, nil

	case CommandFetchEnable:
		return emptyVal, nil

	case CommandFetchFailRequest:
		return emptyVal, nil

	case CommandFetchFulfillRequest:
		return emptyVal, nil

	case CommandFetchContinueRequest:
		return emptyVal, nil

	case CommandFetchContinueWithAuth:
		return emptyVal, nil

	case CommandFetchGetResponseBody:
		v = new(fetch.GetResponseBodyReturns)

	case CommandFetchTakeResponseBodyAsStream:
		v = new(fetch.TakeResponseBodyAsStreamReturns)

	case EventFetchRequestPaused:
		v = new(fetch.EventRequestPaused)

	case EventFetchAuthRequired:
		v = new(fetch.EventAuthRequired)

	case CommandHeadlessExperimentalBeginFrame:
		v = new(headlessexperimental.BeginFrameReturns)

	case CommandHeadlessExperimentalDisable:
		return emptyVal, nil

	case CommandHeadlessExperimentalEnable:
		return emptyVal, nil

	case EventHeadlessExperimentalNeedsBeginFramesChanged:
		v = new(headlessexperimental.EventNeedsBeginFramesChanged)

	case CommandHeapProfilerAddInspectedHeapObject:
		return emptyVal, nil

	case CommandHeapProfilerCollectGarbage:
		return emptyVal, nil

	case CommandHeapProfilerDisable:
		return emptyVal, nil

	case CommandHeapProfilerEnable:
		return emptyVal, nil

	case CommandHeapProfilerGetHeapObjectID:
		v = new(heapprofiler.GetHeapObjectIDReturns)

	case CommandHeapProfilerGetObjectByHeapObjectID:
		v = new(heapprofiler.GetObjectByHeapObjectIDReturns)

	case CommandHeapProfilerGetSamplingProfile:
		v = new(heapprofiler.GetSamplingProfileReturns)

	case CommandHeapProfilerStartSampling:
		return emptyVal, nil

	case CommandHeapProfilerStartTrackingHeapObjects:
		return emptyVal, nil

	case CommandHeapProfilerStopSampling:
		v = new(heapprofiler.StopSamplingReturns)

	case CommandHeapProfilerStopTrackingHeapObjects:
		return emptyVal, nil

	case CommandHeapProfilerTakeHeapSnapshot:
		return emptyVal, nil

	case EventHeapProfilerAddHeapSnapshotChunk:
		v = new(heapprofiler.EventAddHeapSnapshotChunk)

	case EventHeapProfilerHeapStatsUpdate:
		v = new(heapprofiler.EventHeapStatsUpdate)

	case EventHeapProfilerLastSeenObjectID:
		v = new(heapprofiler.EventLastSeenObjectID)

	case EventHeapProfilerReportHeapSnapshotProgress:
		v = new(heapprofiler.EventReportHeapSnapshotProgress)

	case EventHeapProfilerResetProfiles:
		v = new(heapprofiler.EventResetProfiles)

	case CommandIOClose:
		return emptyVal, nil

	case CommandIORead:
		v = new(io.ReadReturns)

	case CommandIOResolveBlob:
		v = new(io.ResolveBlobReturns)

	case CommandIndexedDBClearObjectStore:
		return emptyVal, nil

	case CommandIndexedDBDeleteDatabase:
		return emptyVal, nil

	case CommandIndexedDBDeleteObjectStoreEntries:
		return emptyVal, nil

	case CommandIndexedDBDisable:
		return emptyVal, nil

	case CommandIndexedDBEnable:
		return emptyVal, nil

	case CommandIndexedDBRequestData:
		v = new(indexeddb.RequestDataReturns)

	case CommandIndexedDBGetMetadata:
		v = new(indexeddb.GetMetadataReturns)

	case CommandIndexedDBRequestDatabase:
		v = new(indexeddb.RequestDatabaseReturns)

	case CommandIndexedDBRequestDatabaseNames:
		v = new(indexeddb.RequestDatabaseNamesReturns)

	case CommandInputDispatchKeyEvent:
		return emptyVal, nil

	case CommandInputInsertText:
		return emptyVal, nil

	case CommandInputDispatchMouseEvent:
		return emptyVal, nil

	case CommandInputDispatchTouchEvent:
		return emptyVal, nil

	case CommandInputEmulateTouchFromMouseEvent:
		return emptyVal, nil

	case CommandInputSetIgnoreInputEvents:
		return emptyVal, nil

	case CommandInputSynthesizePinchGesture:
		return emptyVal, nil

	case CommandInputSynthesizeScrollGesture:
		return emptyVal, nil

	case CommandInputSynthesizeTapGesture:
		return emptyVal, nil

	case CommandInspectorDisable:
		return emptyVal, nil

	case CommandInspectorEnable:
		return emptyVal, nil

	case EventInspectorDetached:
		v = new(inspector.EventDetached)

	case EventInspectorTargetCrashed:
		v = new(inspector.EventTargetCrashed)

	case EventInspectorTargetReloadedAfterCrash:
		v = new(inspector.EventTargetReloadedAfterCrash)

	case CommandLayerTreeCompositingReasons:
		v = new(layertree.CompositingReasonsReturns)

	case CommandLayerTreeDisable:
		return emptyVal, nil

	case CommandLayerTreeEnable:
		return emptyVal, nil

	case CommandLayerTreeLoadSnapshot:
		v = new(layertree.LoadSnapshotReturns)

	case CommandLayerTreeMakeSnapshot:
		v = new(layertree.MakeSnapshotReturns)

	case CommandLayerTreeProfileSnapshot:
		v = new(layertree.ProfileSnapshotReturns)

	case CommandLayerTreeReleaseSnapshot:
		return emptyVal, nil

	case CommandLayerTreeReplaySnapshot:
		v = new(layertree.ReplaySnapshotReturns)

	case CommandLayerTreeSnapshotCommandLog:
		v = new(layertree.SnapshotCommandLogReturns)

	case EventLayerTreeLayerPainted:
		v = new(layertree.EventLayerPainted)

	case EventLayerTreeLayerTreeDidChange:
		v = new(layertree.EventLayerTreeDidChange)

	case CommandLogClear:
		return emptyVal, nil

	case CommandLogDisable:
		return emptyVal, nil

	case CommandLogEnable:
		return emptyVal, nil

	case CommandLogStartViolationsReport:
		return emptyVal, nil

	case CommandLogStopViolationsReport:
		return emptyVal, nil

	case EventLogEntryAdded:
		v = new(log.EventEntryAdded)

	case CommandMemoryGetDOMCounters:
		v = new(memory.GetDOMCountersReturns)

	case CommandMemoryPrepareForLeakDetection:
		return emptyVal, nil

	case CommandMemoryForciblyPurgeJavaScriptMemory:
		return emptyVal, nil

	case CommandMemorySetPressureNotificationsSuppressed:
		return emptyVal, nil

	case CommandMemorySimulatePressureNotification:
		return emptyVal, nil

	case CommandMemoryStartSampling:
		return emptyVal, nil

	case CommandMemoryStopSampling:
		return emptyVal, nil

	case CommandMemoryGetAllTimeSamplingProfile:
		v = new(memory.GetAllTimeSamplingProfileReturns)

	case CommandMemoryGetBrowserSamplingProfile:
		v = new(memory.GetBrowserSamplingProfileReturns)

	case CommandMemoryGetSamplingProfile:
		v = new(memory.GetSamplingProfileReturns)

	case CommandNetworkClearBrowserCache:
		return emptyVal, nil

	case CommandNetworkClearBrowserCookies:
		return emptyVal, nil

	case CommandNetworkContinueInterceptedRequest:
		return emptyVal, nil

	case CommandNetworkDeleteCookies:
		return emptyVal, nil

	case CommandNetworkDisable:
		return emptyVal, nil

	case CommandNetworkEmulateNetworkConditions:
		return emptyVal, nil

	case CommandNetworkEnable:
		return emptyVal, nil

	case CommandNetworkGetAllCookies:
		v = new(network.GetAllCookiesReturns)

	case CommandNetworkGetCertificate:
		v = new(network.GetCertificateReturns)

	case CommandNetworkGetCookies:
		v = new(network.GetCookiesReturns)

	case CommandNetworkGetResponseBody:
		v = new(network.GetResponseBodyReturns)

	case CommandNetworkGetRequestPostData:
		v = new(network.GetRequestPostDataReturns)

	case CommandNetworkGetResponseBodyForInterception:
		v = new(network.GetResponseBodyForInterceptionReturns)

	case CommandNetworkTakeResponseBodyForInterceptionAsStream:
		v = new(network.TakeResponseBodyForInterceptionAsStreamReturns)

	case CommandNetworkReplayXHR:
		return emptyVal, nil

	case CommandNetworkSearchInResponseBody:
		v = new(network.SearchInResponseBodyReturns)

	case CommandNetworkSetBlockedURLS:
		return emptyVal, nil

	case CommandNetworkSetBypassServiceWorker:
		return emptyVal, nil

	case CommandNetworkSetCacheDisabled:
		return emptyVal, nil

	case CommandNetworkSetCookie:
		v = new(network.SetCookieReturns)

	case CommandNetworkSetCookies:
		return emptyVal, nil

	case CommandNetworkSetDataSizeLimitsForTest:
		return emptyVal, nil

	case CommandNetworkSetExtraHTTPHeaders:
		return emptyVal, nil

	case CommandNetworkSetRequestInterception:
		return emptyVal, nil

	case EventNetworkDataReceived:
		v = new(network.EventDataReceived)

	case EventNetworkEventSourceMessageReceived:
		v = new(network.EventEventSourceMessageReceived)

	case EventNetworkLoadingFailed:
		v = new(network.EventLoadingFailed)

	case EventNetworkLoadingFinished:
		v = new(network.EventLoadingFinished)

	case EventNetworkRequestIntercepted:
		v = new(network.EventRequestIntercepted)

	case EventNetworkRequestServedFromCache:
		v = new(network.EventRequestServedFromCache)

	case EventNetworkRequestWillBeSent:
		v = new(network.EventRequestWillBeSent)

	case EventNetworkResourceChangedPriority:
		v = new(network.EventResourceChangedPriority)

	case EventNetworkSignedExchangeReceived:
		v = new(network.EventSignedExchangeReceived)

	case EventNetworkResponseReceived:
		v = new(network.EventResponseReceived)

	case EventNetworkWebSocketClosed:
		v = new(network.EventWebSocketClosed)

	case EventNetworkWebSocketCreated:
		v = new(network.EventWebSocketCreated)

	case EventNetworkWebSocketFrameError:
		v = new(network.EventWebSocketFrameError)

	case EventNetworkWebSocketFrameReceived:
		v = new(network.EventWebSocketFrameReceived)

	case EventNetworkWebSocketFrameSent:
		v = new(network.EventWebSocketFrameSent)

	case EventNetworkWebSocketHandshakeResponseReceived:
		v = new(network.EventWebSocketHandshakeResponseReceived)

	case EventNetworkWebSocketWillSendHandshakeRequest:
		v = new(network.EventWebSocketWillSendHandshakeRequest)

	case CommandOverlayDisable:
		return emptyVal, nil

	case CommandOverlayEnable:
		return emptyVal, nil

	case CommandOverlayGetHighlightObjectForTest:
		v = new(overlay.GetHighlightObjectForTestReturns)

	case CommandOverlayHideHighlight:
		return emptyVal, nil

	case CommandOverlayHighlightFrame:
		return emptyVal, nil

	case CommandOverlayHighlightNode:
		return emptyVal, nil

	case CommandOverlayHighlightQuad:
		return emptyVal, nil

	case CommandOverlayHighlightRect:
		return emptyVal, nil

	case CommandOverlaySetInspectMode:
		return emptyVal, nil

	case CommandOverlaySetShowAdHighlights:
		return emptyVal, nil

	case CommandOverlaySetPausedInDebuggerMessage:
		return emptyVal, nil

	case CommandOverlaySetShowDebugBorders:
		return emptyVal, nil

	case CommandOverlaySetShowFPSCounter:
		return emptyVal, nil

	case CommandOverlaySetShowPaintRects:
		return emptyVal, nil

	case CommandOverlaySetShowScrollBottleneckRects:
		return emptyVal, nil

	case CommandOverlaySetShowHitTestBorders:
		return emptyVal, nil

	case CommandOverlaySetShowViewportSizeOnResize:
		return emptyVal, nil

	case EventOverlayInspectNodeRequested:
		v = new(overlay.EventInspectNodeRequested)

	case EventOverlayNodeHighlightRequested:
		v = new(overlay.EventNodeHighlightRequested)

	case EventOverlayScreenshotRequested:
		v = new(overlay.EventScreenshotRequested)

	case EventOverlayInspectModeCanceled:
		v = new(overlay.EventInspectModeCanceled)

	case CommandPageAddScriptToEvaluateOnNewDocument:
		v = new(page.AddScriptToEvaluateOnNewDocumentReturns)

	case CommandPageBringToFront:
		return emptyVal, nil

	case CommandPageCaptureScreenshot:
		v = new(page.CaptureScreenshotReturns)

	case CommandPageCaptureSnapshot:
		v = new(page.CaptureSnapshotReturns)

	case CommandPageCreateIsolatedWorld:
		v = new(page.CreateIsolatedWorldReturns)

	case CommandPageDisable:
		return emptyVal, nil

	case CommandPageEnable:
		return emptyVal, nil

	case CommandPageGetAppManifest:
		v = new(page.GetAppManifestReturns)

	case CommandPageGetInstallabilityErrors:
		v = new(page.GetInstallabilityErrorsReturns)

	case CommandPageGetFrameTree:
		v = new(page.GetFrameTreeReturns)

	case CommandPageGetLayoutMetrics:
		v = new(page.GetLayoutMetricsReturns)

	case CommandPageGetNavigationHistory:
		v = new(page.GetNavigationHistoryReturns)

	case CommandPageResetNavigationHistory:
		return emptyVal, nil

	case CommandPageGetResourceContent:
		v = new(page.GetResourceContentReturns)

	case CommandPageGetResourceTree:
		v = new(page.GetResourceTreeReturns)

	case CommandPageHandleJavaScriptDialog:
		return emptyVal, nil

	case CommandPageNavigate:
		v = new(page.NavigateReturns)

	case CommandPageNavigateToHistoryEntry:
		return emptyVal, nil

	case CommandPagePrintToPDF:
		v = new(page.PrintToPDFReturns)

	case CommandPageReload:
		return emptyVal, nil

	case CommandPageRemoveScriptToEvaluateOnNewDocument:
		return emptyVal, nil

	case CommandPageScreencastFrameAck:
		return emptyVal, nil

	case CommandPageSearchInResource:
		v = new(page.SearchInResourceReturns)

	case CommandPageSetAdBlockingEnabled:
		return emptyVal, nil

	case CommandPageSetBypassCSP:
		return emptyVal, nil

	case CommandPageSetFontFamilies:
		return emptyVal, nil

	case CommandPageSetFontSizes:
		return emptyVal, nil

	case CommandPageSetDocumentContent:
		return emptyVal, nil

	case CommandPageSetDownloadBehavior:
		return emptyVal, nil

	case CommandPageSetLifecycleEventsEnabled:
		return emptyVal, nil

	case CommandPageStartScreencast:
		return emptyVal, nil

	case CommandPageStopLoading:
		return emptyVal, nil

	case CommandPageCrash:
		return emptyVal, nil

	case CommandPageClose:
		return emptyVal, nil

	case CommandPageSetWebLifecycleState:
		return emptyVal, nil

	case CommandPageStopScreencast:
		return emptyVal, nil

	case CommandPageSetProduceCompilationCache:
		return emptyVal, nil

	case CommandPageAddCompilationCache:
		return emptyVal, nil

	case CommandPageClearCompilationCache:
		return emptyVal, nil

	case CommandPageGenerateTestReport:
		return emptyVal, nil

	case CommandPageWaitForDebugger:
		return emptyVal, nil

	case EventPageDomContentEventFired:
		v = new(page.EventDomContentEventFired)

	case EventPageFrameAttached:
		v = new(page.EventFrameAttached)

	case EventPageFrameDetached:
		v = new(page.EventFrameDetached)

	case EventPageFrameNavigated:
		v = new(page.EventFrameNavigated)

	case EventPageFrameResized:
		v = new(page.EventFrameResized)

	case EventPageFrameRequestedNavigation:
		v = new(page.EventFrameRequestedNavigation)

	case EventPageFrameStartedLoading:
		v = new(page.EventFrameStartedLoading)

	case EventPageFrameStoppedLoading:
		v = new(page.EventFrameStoppedLoading)

	case EventPageDownloadWillBegin:
		v = new(page.EventDownloadWillBegin)

	case EventPageInterstitialHidden:
		v = new(page.EventInterstitialHidden)

	case EventPageInterstitialShown:
		v = new(page.EventInterstitialShown)

	case EventPageJavascriptDialogClosed:
		v = new(page.EventJavascriptDialogClosed)

	case EventPageJavascriptDialogOpening:
		v = new(page.EventJavascriptDialogOpening)

	case EventPageLifecycleEvent:
		v = new(page.EventLifecycleEvent)

	case EventPageLoadEventFired:
		v = new(page.EventLoadEventFired)

	case EventPageNavigatedWithinDocument:
		v = new(page.EventNavigatedWithinDocument)

	case EventPageScreencastFrame:
		v = new(page.EventScreencastFrame)

	case EventPageScreencastVisibilityChanged:
		v = new(page.EventScreencastVisibilityChanged)

	case EventPageWindowOpen:
		v = new(page.EventWindowOpen)

	case EventPageCompilationCacheProduced:
		v = new(page.EventCompilationCacheProduced)

	case CommandPerformanceDisable:
		return emptyVal, nil

	case CommandPerformanceEnable:
		return emptyVal, nil

	case CommandPerformanceSetTimeDomain:
		return emptyVal, nil

	case CommandPerformanceGetMetrics:
		v = new(performance.GetMetricsReturns)

	case EventPerformanceMetrics:
		v = new(performance.EventMetrics)

	case CommandProfilerDisable:
		return emptyVal, nil

	case CommandProfilerEnable:
		return emptyVal, nil

	case CommandProfilerGetBestEffortCoverage:
		v = new(profiler.GetBestEffortCoverageReturns)

	case CommandProfilerSetSamplingInterval:
		return emptyVal, nil

	case CommandProfilerStart:
		return emptyVal, nil

	case CommandProfilerStartPreciseCoverage:
		return emptyVal, nil

	case CommandProfilerStartTypeProfile:
		return emptyVal, nil

	case CommandProfilerStop:
		v = new(profiler.StopReturns)

	case CommandProfilerStopPreciseCoverage:
		return emptyVal, nil

	case CommandProfilerStopTypeProfile:
		return emptyVal, nil

	case CommandProfilerTakePreciseCoverage:
		v = new(profiler.TakePreciseCoverageReturns)

	case CommandProfilerTakeTypeProfile:
		v = new(profiler.TakeTypeProfileReturns)

	case EventProfilerConsoleProfileFinished:
		v = new(profiler.EventConsoleProfileFinished)

	case EventProfilerConsoleProfileStarted:
		v = new(profiler.EventConsoleProfileStarted)

	case CommandRuntimeAwaitPromise:
		v = new(runtime.AwaitPromiseReturns)

	case CommandRuntimeCallFunctionOn:
		v = new(runtime.CallFunctionOnReturns)

	case CommandRuntimeCompileScript:
		v = new(runtime.CompileScriptReturns)

	case CommandRuntimeDisable:
		return emptyVal, nil

	case CommandRuntimeDiscardConsoleEntries:
		return emptyVal, nil

	case CommandRuntimeEnable:
		return emptyVal, nil

	case CommandRuntimeEvaluate:
		v = new(runtime.EvaluateReturns)

	case CommandRuntimeGetIsolateID:
		v = new(runtime.GetIsolateIDReturns)

	case CommandRuntimeGetHeapUsage:
		v = new(runtime.GetHeapUsageReturns)

	case CommandRuntimeGetProperties:
		v = new(runtime.GetPropertiesReturns)

	case CommandRuntimeGlobalLexicalScopeNames:
		v = new(runtime.GlobalLexicalScopeNamesReturns)

	case CommandRuntimeQueryObjects:
		v = new(runtime.QueryObjectsReturns)

	case CommandRuntimeReleaseObject:
		return emptyVal, nil

	case CommandRuntimeReleaseObjectGroup:
		return emptyVal, nil

	case CommandRuntimeRunIfWaitingForDebugger:
		return emptyVal, nil

	case CommandRuntimeRunScript:
		v = new(runtime.RunScriptReturns)

	case CommandRuntimeSetCustomObjectFormatterEnabled:
		return emptyVal, nil

	case CommandRuntimeSetMaxCallStackSizeToCapture:
		return emptyVal, nil

	case CommandRuntimeTerminateExecution:
		return emptyVal, nil

	case CommandRuntimeAddBinding:
		return emptyVal, nil

	case CommandRuntimeRemoveBinding:
		return emptyVal, nil

	case EventRuntimeBindingCalled:
		v = new(runtime.EventBindingCalled)

	case EventRuntimeConsoleAPICalled:
		v = new(runtime.EventConsoleAPICalled)

	case EventRuntimeExceptionRevoked:
		v = new(runtime.EventExceptionRevoked)

	case EventRuntimeExceptionThrown:
		v = new(runtime.EventExceptionThrown)

	case EventRuntimeExecutionContextCreated:
		v = new(runtime.EventExecutionContextCreated)

	case EventRuntimeExecutionContextDestroyed:
		v = new(runtime.EventExecutionContextDestroyed)

	case EventRuntimeExecutionContextsCleared:
		v = new(runtime.EventExecutionContextsCleared)

	case EventRuntimeInspectRequested:
		v = new(runtime.EventInspectRequested)

	case CommandSecurityDisable:
		return emptyVal, nil

	case CommandSecurityEnable:
		return emptyVal, nil

	case CommandSecuritySetIgnoreCertificateErrors:
		return emptyVal, nil

	case EventSecuritySecurityStateChanged:
		v = new(security.EventSecurityStateChanged)

	case CommandServiceWorkerDeliverPushMessage:
		return emptyVal, nil

	case CommandServiceWorkerDisable:
		return emptyVal, nil

	case CommandServiceWorkerDispatchSyncEvent:
		return emptyVal, nil

	case CommandServiceWorkerEnable:
		return emptyVal, nil

	case CommandServiceWorkerInspectWorker:
		return emptyVal, nil

	case CommandServiceWorkerSetForceUpdateOnPageLoad:
		return emptyVal, nil

	case CommandServiceWorkerSkipWaiting:
		return emptyVal, nil

	case CommandServiceWorkerStartWorker:
		return emptyVal, nil

	case CommandServiceWorkerStopAllWorkers:
		return emptyVal, nil

	case CommandServiceWorkerStopWorker:
		return emptyVal, nil

	case CommandServiceWorkerUnregister:
		return emptyVal, nil

	case CommandServiceWorkerUpdateRegistration:
		return emptyVal, nil

	case EventServiceWorkerWorkerErrorReported:
		v = new(serviceworker.EventWorkerErrorReported)

	case EventServiceWorkerWorkerRegistrationUpdated:
		v = new(serviceworker.EventWorkerRegistrationUpdated)

	case EventServiceWorkerWorkerVersionUpdated:
		v = new(serviceworker.EventWorkerVersionUpdated)

	case CommandStorageClearDataForOrigin:
		return emptyVal, nil

	case CommandStorageGetUsageAndQuota:
		v = new(storage.GetUsageAndQuotaReturns)

	case CommandStorageTrackCacheStorageForOrigin:
		return emptyVal, nil

	case CommandStorageTrackIndexedDBForOrigin:
		return emptyVal, nil

	case CommandStorageUntrackCacheStorageForOrigin:
		return emptyVal, nil

	case CommandStorageUntrackIndexedDBForOrigin:
		return emptyVal, nil

	case EventStorageCacheStorageContentUpdated:
		v = new(storage.EventCacheStorageContentUpdated)

	case EventStorageCacheStorageListUpdated:
		v = new(storage.EventCacheStorageListUpdated)

	case EventStorageIndexedDBContentUpdated:
		v = new(storage.EventIndexedDBContentUpdated)

	case EventStorageIndexedDBListUpdated:
		v = new(storage.EventIndexedDBListUpdated)

	case CommandSystemInfoGetInfo:
		v = new(systeminfo.GetInfoReturns)

	case CommandSystemInfoGetProcessInfo:
		v = new(systeminfo.GetProcessInfoReturns)

	case CommandTargetActivateTarget:
		return emptyVal, nil

	case CommandTargetAttachToTarget:
		v = new(target.AttachToTargetReturns)

	case CommandTargetAttachToBrowserTarget:
		v = new(target.AttachToBrowserTargetReturns)

	case CommandTargetCloseTarget:
		v = new(target.CloseTargetReturns)

	case CommandTargetExposeDevToolsProtocol:
		return emptyVal, nil

	case CommandTargetCreateBrowserContext:
		v = new(target.CreateBrowserContextReturns)

	case CommandTargetGetBrowserContexts:
		v = new(target.GetBrowserContextsReturns)

	case CommandTargetCreateTarget:
		v = new(target.CreateTargetReturns)

	case CommandTargetDetachFromTarget:
		return emptyVal, nil

	case CommandTargetDisposeBrowserContext:
		return emptyVal, nil

	case CommandTargetGetTargetInfo:
		v = new(target.GetTargetInfoReturns)

	case CommandTargetGetTargets:
		v = new(target.GetTargetsReturns)

	case CommandTargetSendMessageToTarget:
		return emptyVal, nil

	case CommandTargetSetAutoAttach:
		return emptyVal, nil

	case CommandTargetSetDiscoverTargets:
		return emptyVal, nil

	case CommandTargetSetRemoteLocations:
		return emptyVal, nil

	case EventTargetAttachedToTarget:
		v = new(target.EventAttachedToTarget)

	case EventTargetDetachedFromTarget:
		v = new(target.EventDetachedFromTarget)

	case EventTargetReceivedMessageFromTarget:
		v = new(target.EventReceivedMessageFromTarget)

	case EventTargetTargetCreated:
		v = new(target.EventTargetCreated)

	case EventTargetTargetDestroyed:
		v = new(target.EventTargetDestroyed)

	case EventTargetTargetCrashed:
		v = new(target.EventTargetCrashed)

	case EventTargetTargetInfoChanged:
		v = new(target.EventTargetInfoChanged)

	case CommandTetheringBind:
		return emptyVal, nil

	case CommandTetheringUnbind:
		return emptyVal, nil

	case EventTetheringAccepted:
		v = new(tethering.EventAccepted)

	case CommandTracingEnd:
		return emptyVal, nil

	case CommandTracingGetCategories:
		v = new(tracing.GetCategoriesReturns)

	case CommandTracingRecordClockSyncMarker:
		return emptyVal, nil

	case CommandTracingRequestMemoryDump:
		v = new(tracing.RequestMemoryDumpReturns)

	case CommandTracingStart:
		return emptyVal, nil

	case EventTracingBufferUsage:
		v = new(tracing.EventBufferUsage)

	case EventTracingDataCollected:
		v = new(tracing.EventDataCollected)

	case EventTracingTracingComplete:
		v = new(tracing.EventTracingComplete)

	case CommandWebAudioEnable:
		return emptyVal, nil

	case CommandWebAudioDisable:
		return emptyVal, nil

	case CommandWebAudioGetRealtimeData:
		v = new(webaudio.GetRealtimeDataReturns)

	case EventWebAudioContextCreated:
		v = new(webaudio.EventContextCreated)

	case EventWebAudioContextDestroyed:
		v = new(webaudio.EventContextDestroyed)

	case EventWebAudioContextChanged:
		v = new(webaudio.EventContextChanged)

	default:
		return nil, cdp.ErrUnknownCommandOrEvent(msg.Method)
	}

	var buf easyjson.RawMessage
	switch {
	case msg.Params != nil:
		buf = msg.Params

	case msg.Result != nil:
		buf = msg.Result

	default:
		return nil, cdp.ErrMsgMissingParamsOrResult
	}

	err := easyjson.Unmarshal(buf, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
