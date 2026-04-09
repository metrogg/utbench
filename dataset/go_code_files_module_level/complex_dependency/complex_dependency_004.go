func InitWithProcAddrFunc(getProcAddr func(name string) unsafe.Pointer) error {
	gpAcquireKeyedMutexWin32EXT = (C.GPACQUIREKEYEDMUTEXWIN32EXT)(getProcAddr("glAcquireKeyedMutexWin32EXT"))
	gpActiveProgramEXT = (C.GPACTIVEPROGRAMEXT)(getProcAddr("glActiveProgramEXT"))
	gpActiveShaderProgram = (C.GPACTIVESHADERPROGRAM)(getProcAddr("glActiveShaderProgram"))
	if gpActiveShaderProgram == nil {
		return errors.New("glActiveShaderProgram")
	}
	gpActiveShaderProgramEXT = (C.GPACTIVESHADERPROGRAMEXT)(getProcAddr("glActiveShaderProgramEXT"))
	gpActiveTexture = (C.GPACTIVETEXTURE)(getProcAddr("glActiveTexture"))
	if gpActiveTexture == nil {
		return errors.New("glActiveTexture")
	}
	gpAlphaFuncQCOM = (C.GPALPHAFUNCQCOM)(getProcAddr("glAlphaFuncQCOM"))
	gpApplyFramebufferAttachmentCMAAINTEL = (C.GPAPPLYFRAMEBUFFERATTACHMENTCMAAINTEL)(getProcAddr("glApplyFramebufferAttachmentCMAAINTEL"))
	gpAttachShader = (C.GPATTACHSHADER)(getProcAddr("glAttachShader"))
	if gpAttachShader == nil {
		return errors.New("glAttachShader")
	}
	gpBeginConditionalRenderNV = (C.GPBEGINCONDITIONALRENDERNV)(getProcAddr("glBeginConditionalRenderNV"))
	gpBeginPerfMonitorAMD = (C.GPBEGINPERFMONITORAMD)(getProcAddr("glBeginPerfMonitorAMD"))
	gpBeginPerfQueryINTEL = (C.GPBEGINPERFQUERYINTEL)(getProcAddr("glBeginPerfQueryINTEL"))
	gpBeginQuery = (C.GPBEGINQUERY)(getProcAddr("glBeginQuery"))
	if gpBeginQuery == nil {
		return errors.New("glBeginQuery")
	}
	gpBeginQueryEXT = (C.GPBEGINQUERYEXT)(getProcAddr("glBeginQueryEXT"))
	gpBeginTransformFeedback = (C.GPBEGINTRANSFORMFEEDBACK)(getProcAddr("glBeginTransformFeedback"))
	if gpBeginTransformFeedback == nil {
		return errors.New("glBeginTransformFeedback")
	}
	gpBindAttribLocation = (C.GPBINDATTRIBLOCATION)(getProcAddr("glBindAttribLocation"))
	if gpBindAttribLocation == nil {
		return errors.New("glBindAttribLocation")
	}
	gpBindBuffer = (C.GPBINDBUFFER)(getProcAddr("glBindBuffer"))
	if gpBindBuffer == nil {
		return errors.New("glBindBuffer")
	}
	gpBindBufferBase = (C.GPBINDBUFFERBASE)(getProcAddr("glBindBufferBase"))
	if gpBindBufferBase == nil {
		return errors.New("glBindBufferBase")
	}
	gpBindBufferRange = (C.GPBINDBUFFERRANGE)(getProcAddr("glBindBufferRange"))
	if gpBindBufferRange == nil {
		return errors.New("glBindBufferRange")
	}
	gpBindFragDataLocationEXT = (C.GPBINDFRAGDATALOCATIONEXT)(getProcAddr("glBindFragDataLocationEXT"))
	gpBindFragDataLocationIndexedEXT = (C.GPBINDFRAGDATALOCATIONINDEXEDEXT)(getProcAddr("glBindFragDataLocationIndexedEXT"))
	gpBindFramebuffer = (C.GPBINDFRAMEBUFFER)(getProcAddr("glBindFramebuffer"))
	if gpBindFramebuffer == nil {
		return errors.New("glBindFramebuffer")
	}
	gpBindImageTexture = (C.GPBINDIMAGETEXTURE)(getProcAddr("glBindImageTexture"))
	if gpBindImageTexture == nil {
		return errors.New("glBindImageTexture")
	}
	gpBindProgramPipeline = (C.GPBINDPROGRAMPIPELINE)(getProcAddr("glBindProgramPipeline"))
	if gpBindProgramPipeline == nil {
		return errors.New("glBindProgramPipeline")
	}
	gpBindProgramPipelineEXT = (C.GPBINDPROGRAMPIPELINEEXT)(getProcAddr("glBindProgramPipelineEXT"))
	gpBindRenderbuffer = (C.GPBINDRENDERBUFFER)(getProcAddr("glBindRenderbuffer"))
	if gpBindRenderbuffer == nil {
		return errors.New("glBindRenderbuffer")
	}
	gpBindSampler = (C.GPBINDSAMPLER)(getProcAddr("glBindSampler"))
	if gpBindSampler == nil {
		return errors.New("glBindSampler")
	}
	gpBindTexture = (C.GPBINDTEXTURE)(getProcAddr("glBindTexture"))
	if gpBindTexture == nil {
		return errors.New("glBindTexture")
	}
	gpBindTransformFeedback = (C.GPBINDTRANSFORMFEEDBACK)(getProcAddr("glBindTransformFeedback"))
	if gpBindTransformFeedback == nil {
		return errors.New("glBindTransformFeedback")
	}
	gpBindVertexArray = (C.GPBINDVERTEXARRAY)(getProcAddr("glBindVertexArray"))
	if gpBindVertexArray == nil {
		return errors.New("glBindVertexArray")
	}
	gpBindVertexArrayOES = (C.GPBINDVERTEXARRAYOES)(getProcAddr("glBindVertexArrayOES"))
	gpBindVertexBuffer = (C.GPBINDVERTEXBUFFER)(getProcAddr("glBindVertexBuffer"))
	if gpBindVertexBuffer == nil {
		return errors.New("glBindVertexBuffer")
	}
	gpBlendBarrierKHR = (C.GPBLENDBARRIERKHR)(getProcAddr("glBlendBarrierKHR"))
	gpBlendBarrierNV = (C.GPBLENDBARRIERNV)(getProcAddr("glBlendBarrierNV"))
	gpBlendColor = (C.GPBLENDCOLOR)(getProcAddr("glBlendColor"))
	if gpBlendColor == nil {
		return errors.New("glBlendColor")
	}
	gpBlendEquation = (C.GPBLENDEQUATION)(getProcAddr("glBlendEquation"))
	if gpBlendEquation == nil {
		return errors.New("glBlendEquation")
	}
	gpBlendEquationEXT = (C.GPBLENDEQUATIONEXT)(getProcAddr("glBlendEquationEXT"))
	gpBlendEquationSeparate = (C.GPBLENDEQUATIONSEPARATE)(getProcAddr("glBlendEquationSeparate"))
	if gpBlendEquationSeparate == nil {
		return errors.New("glBlendEquationSeparate")
	}
	gpBlendEquationSeparateiEXT = (C.GPBLENDEQUATIONSEPARATEIEXT)(getProcAddr("glBlendEquationSeparateiEXT"))
	gpBlendEquationSeparateiOES = (C.GPBLENDEQUATIONSEPARATEIOES)(getProcAddr("glBlendEquationSeparateiOES"))
	gpBlendEquationiEXT = (C.GPBLENDEQUATIONIEXT)(getProcAddr("glBlendEquationiEXT"))
	gpBlendEquationiOES = (C.GPBLENDEQUATIONIOES)(getProcAddr("glBlendEquationiOES"))
	gpBlendFunc = (C.GPBLENDFUNC)(getProcAddr("glBlendFunc"))
	if gpBlendFunc == nil {
		return errors.New("glBlendFunc")
	}
	gpBlendFuncSeparate = (C.GPBLENDFUNCSEPARATE)(getProcAddr("glBlendFuncSeparate"))
	if gpBlendFuncSeparate == nil {
		return errors.New("glBlendFuncSeparate")
	}
	gpBlendFuncSeparateiEXT = (C.GPBLENDFUNCSEPARATEIEXT)(getProcAddr("glBlendFuncSeparateiEXT"))
	gpBlendFuncSeparateiOES = (C.GPBLENDFUNCSEPARATEIOES)(getProcAddr("glBlendFuncSeparateiOES"))
	gpBlendFunciEXT = (C.GPBLENDFUNCIEXT)(getProcAddr("glBlendFunciEXT"))
	gpBlendFunciOES = (C.GPBLENDFUNCIOES)(getProcAddr("glBlendFunciOES"))
	gpBlendParameteriNV = (C.GPBLENDPARAMETERINV)(getProcAddr("glBlendParameteriNV"))
	gpBlitFramebuffer = (C.GPBLITFRAMEBUFFER)(getProcAddr("glBlitFramebuffer"))
	if gpBlitFramebuffer == nil {
		return errors.New("glBlitFramebuffer")
	}
	gpBlitFramebufferANGLE = (C.GPBLITFRAMEBUFFERANGLE)(getProcAddr("glBlitFramebufferANGLE"))
	gpBlitFramebufferNV = (C.GPBLITFRAMEBUFFERNV)(getProcAddr("glBlitFramebufferNV"))
	gpBufferData = (C.GPBUFFERDATA)(getProcAddr("glBufferData"))
	if gpBufferData == nil {
		return errors.New("glBufferData")
	}
	gpBufferStorageEXT = (C.GPBUFFERSTORAGEEXT)(getProcAddr("glBufferStorageEXT"))
	gpBufferStorageExternalEXT = (C.GPBUFFERSTORAGEEXTERNALEXT)(getProcAddr("glBufferStorageExternalEXT"))
	gpBufferStorageMemEXT = (C.GPBUFFERSTORAGEMEMEXT)(getProcAddr("glBufferStorageMemEXT"))
	gpBufferSubData = (C.GPBUFFERSUBDATA)(getProcAddr("glBufferSubData"))
	if gpBufferSubData == nil {
		return errors.New("glBufferSubData")
	}
	gpCheckFramebufferStatus = (C.GPCHECKFRAMEBUFFERSTATUS)(getProcAddr("glCheckFramebufferStatus"))
	if gpCheckFramebufferStatus == nil {
		return errors.New("glCheckFramebufferStatus")
	}
	gpClear = (C.GPCLEAR)(getProcAddr("glClear"))
	if gpClear == nil {
		return errors.New("glClear")
	}
	gpClearBufferfi = (C.GPCLEARBUFFERFI)(getProcAddr("glClearBufferfi"))
	if gpClearBufferfi == nil {
		return errors.New("glClearBufferfi")
	}
	gpClearBufferfv = (C.GPCLEARBUFFERFV)(getProcAddr("glClearBufferfv"))
	if gpClearBufferfv == nil {
		return errors.New("glClearBufferfv")
	}
	gpClearBufferiv = (C.GPCLEARBUFFERIV)(getProcAddr("glClearBufferiv"))
	if gpClearBufferiv == nil {
		return errors.New("glClearBufferiv")
	}
	gpClearBufferuiv = (C.GPCLEARBUFFERUIV)(getProcAddr("glClearBufferuiv"))
	if gpClearBufferuiv == nil {
		return errors.New("glClearBufferuiv")
	}
	gpClearColor = (C.GPCLEARCOLOR)(getProcAddr("glClearColor"))
	if gpClearColor == nil {
		return errors.New("glClearColor")
	}
	gpClearDepthf = (C.GPCLEARDEPTHF)(getProcAddr("glClearDepthf"))
	if gpClearDepthf == nil {
		return errors.New("glClearDepthf")
	}
	gpClearPixelLocalStorageuiEXT = (C.GPCLEARPIXELLOCALSTORAGEUIEXT)(getProcAddr("glClearPixelLocalStorageuiEXT"))
	gpClearStencil = (C.GPCLEARSTENCIL)(getProcAddr("glClearStencil"))
	if gpClearStencil == nil {
		return errors.New("glClearStencil")
	}
	gpClearTexImageEXT = (C.GPCLEARTEXIMAGEEXT)(getProcAddr("glClearTexImageEXT"))
	gpClearTexSubImageEXT = (C.GPCLEARTEXSUBIMAGEEXT)(getProcAddr("glClearTexSubImageEXT"))
	gpClientWaitSync = (C.GPCLIENTWAITSYNC)(getProcAddr("glClientWaitSync"))
	if gpClientWaitSync == nil {
		return errors.New("glClientWaitSync")
	}
	gpClientWaitSyncAPPLE = (C.GPCLIENTWAITSYNCAPPLE)(getProcAddr("glClientWaitSyncAPPLE"))
	gpClipControlEXT = (C.GPCLIPCONTROLEXT)(getProcAddr("glClipControlEXT"))
	gpColorMask = (C.GPCOLORMASK)(getProcAddr("glColorMask"))
	if gpColorMask == nil {
		return errors.New("glColorMask")
	}
	gpColorMaskiEXT = (C.GPCOLORMASKIEXT)(getProcAddr("glColorMaskiEXT"))
	gpColorMaskiOES = (C.GPCOLORMASKIOES)(getProcAddr("glColorMaskiOES"))
	gpCompileShader = (C.GPCOMPILESHADER)(getProcAddr("glCompileShader"))
	if gpCompileShader == nil {
		return errors.New("glCompileShader")
	}
	gpCompressedTexImage2D = (C.GPCOMPRESSEDTEXIMAGE2D)(getProcAddr("glCompressedTexImage2D"))
	if gpCompressedTexImage2D == nil {
		return errors.New("glCompressedTexImage2D")
	}
	gpCompressedTexImage3D = (C.GPCOMPRESSEDTEXIMAGE3D)(getProcAddr("glCompressedTexImage3D"))
	if gpCompressedTexImage3D == nil {
		return errors.New("glCompressedTexImage3D")
	}
	gpCompressedTexImage3DOES = (C.GPCOMPRESSEDTEXIMAGE3DOES)(getProcAddr("glCompressedTexImage3DOES"))
	gpCompressedTexSubImage2D = (C.GPCOMPRESSEDTEXSUBIMAGE2D)(getProcAddr("glCompressedTexSubImage2D"))
	if gpCompressedTexSubImage2D == nil {
		return errors.New("glCompressedTexSubImage2D")
	}
	gpCompressedTexSubImage3D = (C.GPCOMPRESSEDTEXSUBIMAGE3D)(getProcAddr("glCompressedTexSubImage3D"))
	if gpCompressedTexSubImage3D == nil {
		return errors.New("glCompressedTexSubImage3D")
	}
	gpCompressedTexSubImage3DOES = (C.GPCOMPRESSEDTEXSUBIMAGE3DOES)(getProcAddr("glCompressedTexSubImage3DOES"))
	gpConservativeRasterParameteriNV = (C.GPCONSERVATIVERASTERPARAMETERINV)(getProcAddr("glConservativeRasterParameteriNV"))
	gpCopyBufferSubData = (C.GPCOPYBUFFERSUBDATA)(getProcAddr("glCopyBufferSubData"))
	if gpCopyBufferSubData == nil {
		return errors.New("glCopyBufferSubData")
	}
	gpCopyBufferSubDataNV = (C.GPCOPYBUFFERSUBDATANV)(getProcAddr("glCopyBufferSubDataNV"))
	gpCopyImageSubDataEXT = (C.GPCOPYIMAGESUBDATAEXT)(getProcAddr("glCopyImageSubDataEXT"))
	gpCopyImageSubDataOES = (C.GPCOPYIMAGESUBDATAOES)(getProcAddr("glCopyImageSubDataOES"))
	gpCopyPathNV = (C.GPCOPYPATHNV)(getProcAddr("glCopyPathNV"))
	gpCopyTexImage2D = (C.GPCOPYTEXIMAGE2D)(getProcAddr("glCopyTexImage2D"))
	if gpCopyTexImage2D == nil {
		return errors.New("glCopyTexImage2D")
	}
	gpCopyTexSubImage2D = (C.GPCOPYTEXSUBIMAGE2D)(getProcAddr("glCopyTexSubImage2D"))
	if gpCopyTexSubImage2D == nil {
		return errors.New("glCopyTexSubImage2D")
	}
	gpCopyTexSubImage3D = (C.GPCOPYTEXSUBIMAGE3D)(getProcAddr("glCopyTexSubImage3D"))
	if gpCopyTexSubImage3D == nil {
		return errors.New("glCopyTexSubImage3D")
	}
	gpCopyTexSubImage3DOES = (C.GPCOPYTEXSUBIMAGE3DOES)(getProcAddr("glCopyTexSubImage3DOES"))
	gpCopyTextureLevelsAPPLE = (C.GPCOPYTEXTURELEVELSAPPLE)(getProcAddr("glCopyTextureLevelsAPPLE"))
	gpCoverFillPathInstancedNV = (C.GPCOVERFILLPATHINSTANCEDNV)(getProcAddr("glCoverFillPathInstancedNV"))
	gpCoverFillPathNV = (C.GPCOVERFILLPATHNV)(getProcAddr("glCoverFillPathNV"))
	gpCoverStrokePathInstancedNV = (C.GPCOVERSTROKEPATHINSTANCEDNV)(getProcAddr("glCoverStrokePathInstancedNV"))
	gpCoverStrokePathNV = (C.GPCOVERSTROKEPATHNV)(getProcAddr("glCoverStrokePathNV"))
	gpCoverageMaskNV = (C.GPCOVERAGEMASKNV)(getProcAddr("glCoverageMaskNV"))
	gpCoverageModulationNV = (C.GPCOVERAGEMODULATIONNV)(getProcAddr("glCoverageModulationNV"))
	gpCoverageModulationTableNV = (C.GPCOVERAGEMODULATIONTABLENV)(getProcAddr("glCoverageModulationTableNV"))
	gpCoverageOperationNV = (C.GPCOVERAGEOPERATIONNV)(getProcAddr("glCoverageOperationNV"))
	gpCreateMemoryObjectsEXT = (C.GPCREATEMEMORYOBJECTSEXT)(getProcAddr("glCreateMemoryObjectsEXT"))
	gpCreatePerfQueryINTEL = (C.GPCREATEPERFQUERYINTEL)(getProcAddr("glCreatePerfQueryINTEL"))
	gpCreateProgram = (C.GPCREATEPROGRAM)(getProcAddr("glCreateProgram"))
	if gpCreateProgram == nil {
		return errors.New("glCreateProgram")
	}
	gpCreateShader = (C.GPCREATESHADER)(getProcAddr("glCreateShader"))
	if gpCreateShader == nil {
		return errors.New("glCreateShader")
	}
	gpCreateShaderProgramEXT = (C.GPCREATESHADERPROGRAMEXT)(getProcAddr("glCreateShaderProgramEXT"))
	gpCreateShaderProgramv = (C.GPCREATESHADERPROGRAMV)(getProcAddr("glCreateShaderProgramv"))
	if gpCreateShaderProgramv == nil {
		return errors.New("glCreateShaderProgramv")
	}
	gpCreateShaderProgramvEXT = (C.GPCREATESHADERPROGRAMVEXT)(getProcAddr("glCreateShaderProgramvEXT"))
	gpCullFace = (C.GPCULLFACE)(getProcAddr("glCullFace"))
	if gpCullFace == nil {
		return errors.New("glCullFace")
	}
	gpDebugMessageCallback = (C.GPDEBUGMESSAGECALLBACK)(getProcAddr("glDebugMessageCallback"))
	gpDebugMessageCallbackKHR = (C.GPDEBUGMESSAGECALLBACKKHR)(getProcAddr("glDebugMessageCallbackKHR"))
	gpDebugMessageControl = (C.GPDEBUGMESSAGECONTROL)(getProcAddr("glDebugMessageControl"))
	gpDebugMessageControlKHR = (C.GPDEBUGMESSAGECONTROLKHR)(getProcAddr("glDebugMessageControlKHR"))
	gpDebugMessageInsert = (C.GPDEBUGMESSAGEINSERT)(getProcAddr("glDebugMessageInsert"))
	gpDebugMessageInsertKHR = (C.GPDEBUGMESSAGEINSERTKHR)(getProcAddr("glDebugMessageInsertKHR"))
	gpDeleteBuffers = (C.GPDELETEBUFFERS)(getProcAddr("glDeleteBuffers"))
	if gpDeleteBuffers == nil {
		return errors.New("glDeleteBuffers")
	}
	gpDeleteFencesNV = (C.GPDELETEFENCESNV)(getProcAddr("glDeleteFencesNV"))
	gpDeleteFramebuffers = (C.GPDELETEFRAMEBUFFERS)(getProcAddr("glDeleteFramebuffers"))
	if gpDeleteFramebuffers == nil {
		return errors.New("glDeleteFramebuffers")
	}
	gpDeleteMemoryObjectsEXT = (C.GPDELETEMEMORYOBJECTSEXT)(getProcAddr("glDeleteMemoryObjectsEXT"))
	gpDeletePathsNV = (C.GPDELETEPATHSNV)(getProcAddr("glDeletePathsNV"))
	gpDeletePerfMonitorsAMD = (C.GPDELETEPERFMONITORSAMD)(getProcAddr("glDeletePerfMonitorsAMD"))
	gpDeletePerfQueryINTEL = (C.GPDELETEPERFQUERYINTEL)(getProcAddr("glDeletePerfQueryINTEL"))
	gpDeleteProgram = (C.GPDELETEPROGRAM)(getProcAddr("glDeleteProgram"))
	if gpDeleteProgram == nil {
		return errors.New("glDeleteProgram")
	}
	gpDeleteProgramPipelines = (C.GPDELETEPROGRAMPIPELINES)(getProcAddr("glDeleteProgramPipelines"))
	if gpDeleteProgramPipelines == nil {
		return errors.New("glDeleteProgramPipelines")
	}
	gpDeleteProgramPipelinesEXT = (C.GPDELETEPROGRAMPIPELINESEXT)(getProcAddr("glDeleteProgramPipelinesEXT"))
	gpDeleteQueries = (C.GPDELETEQUERIES)(getProcAddr("glDeleteQueries"))
	if gpDeleteQueries == nil {
		return errors.New("glDeleteQueries")
	}
	gpDeleteQueriesEXT = (C.GPDELETEQUERIESEXT)(getProcAddr("glDeleteQueriesEXT"))
	gpDeleteRenderbuffers = (C.GPDELETERENDERBUFFERS)(getProcAddr("glDeleteRenderbuffers"))
	if gpDeleteRenderbuffers == nil {
		return errors.New("glDeleteRenderbuffers")
	}
	gpDeleteSamplers = (C.GPDELETESAMPLERS)(getProcAddr("glDeleteSamplers"))
	if gpDeleteSamplers == nil {
		return errors.New("glDeleteSamplers")
	}
	gpDeleteSemaphoresEXT = (C.GPDELETESEMAPHORESEXT)(getProcAddr("glDeleteSemaphoresEXT"))
	gpDeleteShader = (C.GPDELETESHADER)(getProcAddr("glDeleteShader"))
	if gpDeleteShader == nil {
		return errors.New("glDeleteShader")
	}
	gpDeleteSync = (C.GPDELETESYNC)(getProcAddr("glDeleteSync"))
	if gpDeleteSync == nil {
		return errors.New("glDeleteSync")
	}
	gpDeleteSyncAPPLE = (C.GPDELETESYNCAPPLE)(getProcAddr("glDeleteSyncAPPLE"))
	gpDeleteTextures = (C.GPDELETETEXTURES)(getProcAddr("glDeleteTextures"))
	if gpDeleteTextures == nil {
		return errors.New("glDeleteTextures")
	}
	gpDeleteTransformFeedbacks = (C.GPDELETETRANSFORMFEEDBACKS)(getProcAddr("glDeleteTransformFeedbacks"))
	if gpDeleteTransformFeedbacks == nil {
		return errors.New("glDeleteTransformFeedbacks")
	}
	gpDeleteVertexArrays = (C.GPDELETEVERTEXARRAYS)(getProcAddr("glDeleteVertexArrays"))
	if gpDeleteVertexArrays == nil {
		return errors.New("glDeleteVertexArrays")
	}
	gpDeleteVertexArraysOES = (C.GPDELETEVERTEXARRAYSOES)(getProcAddr("glDeleteVertexArraysOES"))
	gpDepthFunc = (C.GPDEPTHFUNC)(getProcAddr("glDepthFunc"))
	if gpDepthFunc == nil {
		return errors.New("glDepthFunc")
	}
	gpDepthMask = (C.GPDEPTHMASK)(getProcAddr("glDepthMask"))
	if gpDepthMask == nil {
		return errors.New("glDepthMask")
	}
	gpDepthRangeArrayfvNV = (C.GPDEPTHRANGEARRAYFVNV)(getProcAddr("glDepthRangeArrayfvNV"))
	gpDepthRangeArrayfvOES = (C.GPDEPTHRANGEARRAYFVOES)(getProcAddr("glDepthRangeArrayfvOES"))
	gpDepthRangeIndexedfNV = (C.GPDEPTHRANGEINDEXEDFNV)(getProcAddr("glDepthRangeIndexedfNV"))
	gpDepthRangeIndexedfOES = (C.GPDEPTHRANGEINDEXEDFOES)(getProcAddr("glDepthRangeIndexedfOES"))
	gpDepthRangef = (C.GPDEPTHRANGEF)(getProcAddr("glDepthRangef"))
	if gpDepthRangef == nil {
		return errors.New("glDepthRangef")
	}
	gpDetachShader = (C.GPDETACHSHADER)(getProcAddr("glDetachShader"))
	if gpDetachShader == nil {
		return errors.New("glDetachShader")
	}
	gpDisable = (C.GPDISABLE)(getProcAddr("glDisable"))
	if gpDisable == nil {
		return errors.New("glDisable")
	}
	gpDisableDriverControlQCOM = (C.GPDISABLEDRIVERCONTROLQCOM)(getProcAddr("glDisableDriverControlQCOM"))
	gpDisableVertexAttribArray = (C.GPDISABLEVERTEXATTRIBARRAY)(getProcAddr("glDisableVertexAttribArray"))
	if gpDisableVertexAttribArray == nil {
		return errors.New("glDisableVertexAttribArray")
	}
	gpDisableiEXT = (C.GPDISABLEIEXT)(getProcAddr("glDisableiEXT"))
	gpDisableiNV = (C.GPDISABLEINV)(getProcAddr("glDisableiNV"))
	gpDisableiOES = (C.GPDISABLEIOES)(getProcAddr("glDisableiOES"))
	gpDiscardFramebufferEXT = (C.GPDISCARDFRAMEBUFFEREXT)(getProcAddr("glDiscardFramebufferEXT"))
	gpDispatchCompute = (C.GPDISPATCHCOMPUTE)(getProcAddr("glDispatchCompute"))
	if gpDispatchCompute == nil {
		return errors.New("glDispatchCompute")
	}
	gpDispatchComputeIndirect = (C.GPDISPATCHCOMPUTEINDIRECT)(getProcAddr("glDispatchComputeIndirect"))
	if gpDispatchComputeIndirect == nil {
		return errors.New("glDispatchComputeIndirect")
	}
	gpDrawArrays = (C.GPDRAWARRAYS)(getProcAddr("glDrawArrays"))
	if gpDrawArrays == nil {
		return errors.New("glDrawArrays")
	}
	gpDrawArraysIndirect = (C.GPDRAWARRAYSINDIRECT)(getProcAddr("glDrawArraysIndirect"))
	if gpDrawArraysIndirect == nil {
		return errors.New("glDrawArraysIndirect")
	}
	gpDrawArraysInstanced = (C.GPDRAWARRAYSINSTANCED)(getProcAddr("glDrawArraysInstanced"))
	if gpDrawArraysInstanced == nil {
		return errors.New("glDrawArraysInstanced")
	}
	gpDrawArraysInstancedANGLE = (C.GPDRAWARRAYSINSTANCEDANGLE)(getProcAddr("glDrawArraysInstancedANGLE"))
	gpDrawArraysInstancedBaseInstanceEXT = (C.GPDRAWARRAYSINSTANCEDBASEINSTANCEEXT)(getProcAddr("glDrawArraysInstancedBaseInstanceEXT"))
	gpDrawArraysInstancedEXT = (C.GPDRAWARRAYSINSTANCEDEXT)(getProcAddr("glDrawArraysInstancedEXT"))
	gpDrawArraysInstancedNV = (C.GPDRAWARRAYSINSTANCEDNV)(getProcAddr("glDrawArraysInstancedNV"))
	gpDrawBuffers = (C.GPDRAWBUFFERS)(getProcAddr("glDrawBuffers"))
	if gpDrawBuffers == nil {
		return errors.New("glDrawBuffers")
	}
	gpDrawBuffersEXT = (C.GPDRAWBUFFERSEXT)(getProcAddr("glDrawBuffersEXT"))
	gpDrawBuffersIndexedEXT = (C.GPDRAWBUFFERSINDEXEDEXT)(getProcAddr("glDrawBuffersIndexedEXT"))
	gpDrawBuffersNV = (C.GPDRAWBUFFERSNV)(getProcAddr("glDrawBuffersNV"))
	gpDrawElements = (C.GPDRAWELEMENTS)(getProcAddr("glDrawElements"))
	if gpDrawElements == nil {
		return errors.New("glDrawElements")
	}
	gpDrawElementsBaseVertexEXT = (C.GPDRAWELEMENTSBASEVERTEXEXT)(getProcAddr("glDrawElementsBaseVertexEXT"))
	gpDrawElementsBaseVertexOES = (C.GPDRAWELEMENTSBASEVERTEXOES)(getProcAddr("glDrawElementsBaseVertexOES"))
	gpDrawElementsIndirect = (C.GPDRAWELEMENTSINDIRECT)(getProcAddr("glDrawElementsIndirect"))
	if gpDrawElementsIndirect == nil {
		return errors.New("glDrawElementsIndirect")
	}
	gpDrawElementsInstanced = (C.GPDRAWELEMENTSINSTANCED)(getProcAddr("glDrawElementsInstanced"))
	if gpDrawElementsInstanced == nil {
		return errors.New("glDrawElementsInstanced")
	}
	gpDrawElementsInstancedANGLE = (C.GPDRAWELEMENTSINSTANCEDANGLE)(getProcAddr("glDrawElementsInstancedANGLE"))
	gpDrawElementsInstancedBaseInstanceEXT = (C.GPDRAWELEMENTSINSTANCEDBASEINSTANCEEXT)(getProcAddr("glDrawElementsInstancedBaseInstanceEXT"))
	gpDrawElementsInstancedBaseVertexBaseInstanceEXT = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEXBASEINSTANCEEXT)(getProcAddr("glDrawElementsInstancedBaseVertexBaseInstanceEXT"))
	gpDrawElementsInstancedBaseVertexEXT = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEXEXT)(getProcAddr("glDrawElementsInstancedBaseVertexEXT"))
	gpDrawElementsInstancedBaseVertexOES = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEXOES)(getProcAddr("glDrawElementsInstancedBaseVertexOES"))
	gpDrawElementsInstancedEXT = (C.GPDRAWELEMENTSINSTANCEDEXT)(getProcAddr("glDrawElementsInstancedEXT"))
	gpDrawElementsInstancedNV = (C.GPDRAWELEMENTSINSTANCEDNV)(getProcAddr("glDrawElementsInstancedNV"))
	gpDrawRangeElements = (C.GPDRAWRANGEELEMENTS)(getProcAddr("glDrawRangeElements"))
	if gpDrawRangeElements == nil {
		return errors.New("glDrawRangeElements")
	}
	gpDrawRangeElementsBaseVertexEXT = (C.GPDRAWRANGEELEMENTSBASEVERTEXEXT)(getProcAddr("glDrawRangeElementsBaseVertexEXT"))
	gpDrawRangeElementsBaseVertexOES = (C.GPDRAWRANGEELEMENTSBASEVERTEXOES)(getProcAddr("glDrawRangeElementsBaseVertexOES"))
	gpDrawTransformFeedbackEXT = (C.GPDRAWTRANSFORMFEEDBACKEXT)(getProcAddr("glDrawTransformFeedbackEXT"))
	gpDrawTransformFeedbackInstancedEXT = (C.GPDRAWTRANSFORMFEEDBACKINSTANCEDEXT)(getProcAddr("glDrawTransformFeedbackInstancedEXT"))
	gpDrawVkImageNV = (C.GPDRAWVKIMAGENV)(getProcAddr("glDrawVkImageNV"))
	gpEGLImageTargetRenderbufferStorageOES = (C.GPEGLIMAGETARGETRENDERBUFFERSTORAGEOES)(getProcAddr("glEGLImageTargetRenderbufferStorageOES"))
	gpEGLImageTargetTexStorageEXT = (C.GPEGLIMAGETARGETTEXSTORAGEEXT)(getProcAddr("glEGLImageTargetTexStorageEXT"))
	gpEGLImageTargetTexture2DOES = (C.GPEGLIMAGETARGETTEXTURE2DOES)(getProcAddr("glEGLImageTargetTexture2DOES"))
	gpEGLImageTargetTextureStorageEXT = (C.GPEGLIMAGETARGETTEXTURESTORAGEEXT)(getProcAddr("glEGLImageTargetTextureStorageEXT"))
	gpEnable = (C.GPENABLE)(getProcAddr("glEnable"))
	if gpEnable == nil {
		return errors.New("glEnable")
	}
	gpEnableDriverControlQCOM = (C.GPENABLEDRIVERCONTROLQCOM)(getProcAddr("glEnableDriverControlQCOM"))
	gpEnableVertexAttribArray = (C.GPENABLEVERTEXATTRIBARRAY)(getProcAddr("glEnableVertexAttribArray"))
	if gpEnableVertexAttribArray == nil {
		return errors.New("glEnableVertexAttribArray")
	}
	gpEnableiEXT = (C.GPENABLEIEXT)(getProcAddr("glEnableiEXT"))
	gpEnableiNV = (C.GPENABLEINV)(getProcAddr("glEnableiNV"))
	gpEnableiOES = (C.GPENABLEIOES)(getProcAddr("glEnableiOES"))
	gpEndConditionalRenderNV = (C.GPENDCONDITIONALRENDERNV)(getProcAddr("glEndConditionalRenderNV"))
	gpEndPerfMonitorAMD = (C.GPENDPERFMONITORAMD)(getProcAddr("glEndPerfMonitorAMD"))
	gpEndPerfQueryINTEL = (C.GPENDPERFQUERYINTEL)(getProcAddr("glEndPerfQueryINTEL"))
	gpEndQuery = (C.GPENDQUERY)(getProcAddr("glEndQuery"))
	if gpEndQuery == nil {
		return errors.New("glEndQuery")
	}
	gpEndQueryEXT = (C.GPENDQUERYEXT)(getProcAddr("glEndQueryEXT"))
	gpEndTilingQCOM = (C.GPENDTILINGQCOM)(getProcAddr("glEndTilingQCOM"))
	gpEndTransformFeedback = (C.GPENDTRANSFORMFEEDBACK)(getProcAddr("glEndTransformFeedback"))
	if gpEndTransformFeedback == nil {
		return errors.New("glEndTransformFeedback")
	}
	gpExtGetBufferPointervQCOM = (C.GPEXTGETBUFFERPOINTERVQCOM)(getProcAddr("glExtGetBufferPointervQCOM"))
	gpExtGetBuffersQCOM = (C.GPEXTGETBUFFERSQCOM)(getProcAddr("glExtGetBuffersQCOM"))
	gpExtGetFramebuffersQCOM = (C.GPEXTGETFRAMEBUFFERSQCOM)(getProcAddr("glExtGetFramebuffersQCOM"))
	gpExtGetProgramBinarySourceQCOM = (C.GPEXTGETPROGRAMBINARYSOURCEQCOM)(getProcAddr("glExtGetProgramBinarySourceQCOM"))
	gpExtGetProgramsQCOM = (C.GPEXTGETPROGRAMSQCOM)(getProcAddr("glExtGetProgramsQCOM"))
	gpExtGetRenderbuffersQCOM = (C.GPEXTGETRENDERBUFFERSQCOM)(getProcAddr("glExtGetRenderbuffersQCOM"))
	gpExtGetShadersQCOM = (C.GPEXTGETSHADERSQCOM)(getProcAddr("glExtGetShadersQCOM"))
	gpExtGetTexLevelParameterivQCOM = (C.GPEXTGETTEXLEVELPARAMETERIVQCOM)(getProcAddr("glExtGetTexLevelParameterivQCOM"))
	gpExtGetTexSubImageQCOM = (C.GPEXTGETTEXSUBIMAGEQCOM)(getProcAddr("glExtGetTexSubImageQCOM"))
	gpExtGetTexturesQCOM = (C.GPEXTGETTEXTURESQCOM)(getProcAddr("glExtGetTexturesQCOM"))
	gpExtIsProgramBinaryQCOM = (C.GPEXTISPROGRAMBINARYQCOM)(getProcAddr("glExtIsProgramBinaryQCOM"))
	gpExtTexObjectStateOverrideiQCOM = (C.GPEXTTEXOBJECTSTATEOVERRIDEIQCOM)(getProcAddr("glExtTexObjectStateOverrideiQCOM"))
	gpFenceSync = (C.GPFENCESYNC)(getProcAddr("glFenceSync"))
	if gpFenceSync == nil {
		return errors.New("glFenceSync")
	}
	gpFenceSyncAPPLE = (C.GPFENCESYNCAPPLE)(getProcAddr("glFenceSyncAPPLE"))
	gpFinish = (C.GPFINISH)(getProcAddr("glFinish"))
	if gpFinish == nil {
		return errors.New("glFinish")
	}
	gpFinishFenceNV = (C.GPFINISHFENCENV)(getProcAddr("glFinishFenceNV"))
	gpFlush = (C.GPFLUSH)(getProcAddr("glFlush"))
	if gpFlush == nil {
		return errors.New("glFlush")
	}
	gpFlushMappedBufferRange = (C.GPFLUSHMAPPEDBUFFERRANGE)(getProcAddr("glFlushMappedBufferRange"))
	if gpFlushMappedBufferRange == nil {
		return errors.New("glFlushMappedBufferRange")
	}
	gpFlushMappedBufferRangeEXT = (C.GPFLUSHMAPPEDBUFFERRANGEEXT)(getProcAddr("glFlushMappedBufferRangeEXT"))
	gpFragmentCoverageColorNV = (C.GPFRAGMENTCOVERAGECOLORNV)(getProcAddr("glFragmentCoverageColorNV"))
	gpFramebufferFetchBarrierEXT = (C.GPFRAMEBUFFERFETCHBARRIEREXT)(getProcAddr("glFramebufferFetchBarrierEXT"))
	gpFramebufferFetchBarrierQCOM = (C.GPFRAMEBUFFERFETCHBARRIERQCOM)(getProcAddr("glFramebufferFetchBarrierQCOM"))
	gpFramebufferFoveationConfigQCOM = (C.GPFRAMEBUFFERFOVEATIONCONFIGQCOM)(getProcAddr("glFramebufferFoveationConfigQCOM"))
	gpFramebufferFoveationParametersQCOM = (C.GPFRAMEBUFFERFOVEATIONPARAMETERSQCOM)(getProcAddr("glFramebufferFoveationParametersQCOM"))
	gpFramebufferParameteri = (C.GPFRAMEBUFFERPARAMETERI)(getProcAddr("glFramebufferParameteri"))
	if gpFramebufferParameteri == nil {
		return errors.New("glFramebufferParameteri")
	}
	gpFramebufferPixelLocalStorageSizeEXT = (C.GPFRAMEBUFFERPIXELLOCALSTORAGESIZEEXT)(getProcAddr("glFramebufferPixelLocalStorageSizeEXT"))
	gpFramebufferRenderbuffer = (C.GPFRAMEBUFFERRENDERBUFFER)(getProcAddr("glFramebufferRenderbuffer"))
	if gpFramebufferRenderbuffer == nil {
		return errors.New("glFramebufferRenderbuffer")
	}
	gpFramebufferSampleLocationsfvNV = (C.GPFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glFramebufferSampleLocationsfvNV"))
	gpFramebufferTexture2D = (C.GPFRAMEBUFFERTEXTURE2D)(getProcAddr("glFramebufferTexture2D"))
	if gpFramebufferTexture2D == nil {
		return errors.New("glFramebufferTexture2D")
	}
	gpFramebufferTexture2DDownsampleIMG = (C.GPFRAMEBUFFERTEXTURE2DDOWNSAMPLEIMG)(getProcAddr("glFramebufferTexture2DDownsampleIMG"))
	gpFramebufferTexture2DMultisampleEXT = (C.GPFRAMEBUFFERTEXTURE2DMULTISAMPLEEXT)(getProcAddr("glFramebufferTexture2DMultisampleEXT"))
	gpFramebufferTexture2DMultisampleIMG = (C.GPFRAMEBUFFERTEXTURE2DMULTISAMPLEIMG)(getProcAddr("glFramebufferTexture2DMultisampleIMG"))
	gpFramebufferTexture3DOES = (C.GPFRAMEBUFFERTEXTURE3DOES)(getProcAddr("glFramebufferTexture3DOES"))
	gpFramebufferTextureEXT = (C.GPFRAMEBUFFERTEXTUREEXT)(getProcAddr("glFramebufferTextureEXT"))
	gpFramebufferTextureLayer = (C.GPFRAMEBUFFERTEXTURELAYER)(getProcAddr("glFramebufferTextureLayer"))
	if gpFramebufferTextureLayer == nil {
		return errors.New("glFramebufferTextureLayer")
	}
	gpFramebufferTextureLayerDownsampleIMG = (C.GPFRAMEBUFFERTEXTURELAYERDOWNSAMPLEIMG)(getProcAddr("glFramebufferTextureLayerDownsampleIMG"))
	gpFramebufferTextureMultisampleMultiviewOVR = (C.GPFRAMEBUFFERTEXTUREMULTISAMPLEMULTIVIEWOVR)(getProcAddr("glFramebufferTextureMultisampleMultiviewOVR"))
	gpFramebufferTextureMultiviewOVR = (C.GPFRAMEBUFFERTEXTUREMULTIVIEWOVR)(getProcAddr("glFramebufferTextureMultiviewOVR"))
	gpFramebufferTextureOES = (C.GPFRAMEBUFFERTEXTUREOES)(getProcAddr("glFramebufferTextureOES"))
	gpFrontFace = (C.GPFRONTFACE)(getProcAddr("glFrontFace"))
	if gpFrontFace == nil {
		return errors.New("glFrontFace")
	}
	gpGenBuffers = (C.GPGENBUFFERS)(getProcAddr("glGenBuffers"))
	if gpGenBuffers == nil {
		return errors.New("glGenBuffers")
	}
	gpGenFencesNV = (C.GPGENFENCESNV)(getProcAddr("glGenFencesNV"))
	gpGenFramebuffers = (C.GPGENFRAMEBUFFERS)(getProcAddr("glGenFramebuffers"))
	if gpGenFramebuffers == nil {
		return errors.New("glGenFramebuffers")
	}
	gpGenPathsNV = (C.GPGENPATHSNV)(getProcAddr("glGenPathsNV"))
	gpGenPerfMonitorsAMD = (C.GPGENPERFMONITORSAMD)(getProcAddr("glGenPerfMonitorsAMD"))
	gpGenProgramPipelines = (C.GPGENPROGRAMPIPELINES)(getProcAddr("glGenProgramPipelines"))
	if gpGenProgramPipelines == nil {
		return errors.New("glGenProgramPipelines")
	}
	gpGenProgramPipelinesEXT = (C.GPGENPROGRAMPIPELINESEXT)(getProcAddr("glGenProgramPipelinesEXT"))
	gpGenQueries = (C.GPGENQUERIES)(getProcAddr("glGenQueries"))
	if gpGenQueries == nil {
		return errors.New("glGenQueries")
	}
	gpGenQueriesEXT = (C.GPGENQUERIESEXT)(getProcAddr("glGenQueriesEXT"))
	gpGenRenderbuffers = (C.GPGENRENDERBUFFERS)(getProcAddr("glGenRenderbuffers"))
	if gpGenRenderbuffers == nil {
		return errors.New("glGenRenderbuffers")
	}
	gpGenSamplers = (C.GPGENSAMPLERS)(getProcAddr("glGenSamplers"))
	if gpGenSamplers == nil {
		return errors.New("glGenSamplers")
	}
	gpGenSemaphoresEXT = (C.GPGENSEMAPHORESEXT)(getProcAddr("glGenSemaphoresEXT"))
	gpGenTextures = (C.GPGENTEXTURES)(getProcAddr("glGenTextures"))
	if gpGenTextures == nil {
		return errors.New("glGenTextures")
	}
	gpGenTransformFeedbacks = (C.GPGENTRANSFORMFEEDBACKS)(getProcAddr("glGenTransformFeedbacks"))
	if gpGenTransformFeedbacks == nil {
		return errors.New("glGenTransformFeedbacks")
	}
	gpGenVertexArrays = (C.GPGENVERTEXARRAYS)(getProcAddr("glGenVertexArrays"))
	if gpGenVertexArrays == nil {
		return errors.New("glGenVertexArrays")
	}
	gpGenVertexArraysOES = (C.GPGENVERTEXARRAYSOES)(getProcAddr("glGenVertexArraysOES"))
	gpGenerateMipmap = (C.GPGENERATEMIPMAP)(getProcAddr("glGenerateMipmap"))
	if gpGenerateMipmap == nil {
		return errors.New("glGenerateMipmap")
	}
	gpGetActiveAttrib = (C.GPGETACTIVEATTRIB)(getProcAddr("glGetActiveAttrib"))
	if gpGetActiveAttrib == nil {
		return errors.New("glGetActiveAttrib")
	}
	gpGetActiveUniform = (C.GPGETACTIVEUNIFORM)(getProcAddr("glGetActiveUniform"))
	if gpGetActiveUniform == nil {
		return errors.New("glGetActiveUniform")
	}
	gpGetActiveUniformBlockName = (C.GPGETACTIVEUNIFORMBLOCKNAME)(getProcAddr("glGetActiveUniformBlockName"))
	if gpGetActiveUniformBlockName == nil {
		return errors.New("glGetActiveUniformBlockName")
	}
	gpGetActiveUniformBlockiv = (C.GPGETACTIVEUNIFORMBLOCKIV)(getProcAddr("glGetActiveUniformBlockiv"))
	if gpGetActiveUniformBlockiv == nil {
		return errors.New("glGetActiveUniformBlockiv")
	}
	gpGetActiveUniformsiv = (C.GPGETACTIVEUNIFORMSIV)(getProcAddr("glGetActiveUniformsiv"))
	if gpGetActiveUniformsiv == nil {
		return errors.New("glGetActiveUniformsiv")
	}
	gpGetAttachedShaders = (C.GPGETATTACHEDSHADERS)(getProcAddr("glGetAttachedShaders"))
	if gpGetAttachedShaders == nil {
		return errors.New("glGetAttachedShaders")
	}
	gpGetAttribLocation = (C.GPGETATTRIBLOCATION)(getProcAddr("glGetAttribLocation"))
	if gpGetAttribLocation == nil {
		return errors.New("glGetAttribLocation")
	}
	gpGetBooleani_v = (C.GPGETBOOLEANI_V)(getProcAddr("glGetBooleani_v"))
	if gpGetBooleani_v == nil {
		return errors.New("glGetBooleani_v")
	}
	gpGetBooleanv = (C.GPGETBOOLEANV)(getProcAddr("glGetBooleanv"))
	if gpGetBooleanv == nil {
		return errors.New("glGetBooleanv")
	}
	gpGetBufferParameteri64v = (C.GPGETBUFFERPARAMETERI64V)(getProcAddr("glGetBufferParameteri64v"))
	if gpGetBufferParameteri64v == nil {
		return errors.New("glGetBufferParameteri64v")
	}
	gpGetBufferParameteriv = (C.GPGETBUFFERPARAMETERIV)(getProcAddr("glGetBufferParameteriv"))
	if gpGetBufferParameteriv == nil {
		return errors.New("glGetBufferParameteriv")
	}
	gpGetBufferPointerv = (C.GPGETBUFFERPOINTERV)(getProcAddr("glGetBufferPointerv"))
	if gpGetBufferPointerv == nil {
		return errors.New("glGetBufferPointerv")
	}
	gpGetBufferPointervOES = (C.GPGETBUFFERPOINTERVOES)(getProcAddr("glGetBufferPointervOES"))
	gpGetCoverageModulationTableNV = (C.GPGETCOVERAGEMODULATIONTABLENV)(getProcAddr("glGetCoverageModulationTableNV"))
	gpGetDebugMessageLog = (C.GPGETDEBUGMESSAGELOG)(getProcAddr("glGetDebugMessageLog"))
	gpGetDebugMessageLogKHR = (C.GPGETDEBUGMESSAGELOGKHR)(getProcAddr("glGetDebugMessageLogKHR"))
	gpGetDriverControlStringQCOM = (C.GPGETDRIVERCONTROLSTRINGQCOM)(getProcAddr("glGetDriverControlStringQCOM"))
	gpGetDriverControlsQCOM = (C.GPGETDRIVERCONTROLSQCOM)(getProcAddr("glGetDriverControlsQCOM"))
	gpGetError = (C.GPGETERROR)(getProcAddr("glGetError"))
	if gpGetError == nil {
		return errors.New("glGetError")
	}
	gpGetFenceivNV = (C.GPGETFENCEIVNV)(getProcAddr("glGetFenceivNV"))
	gpGetFirstPerfQueryIdINTEL = (C.GPGETFIRSTPERFQUERYIDINTEL)(getProcAddr("glGetFirstPerfQueryIdINTEL"))
	gpGetFloati_vNV = (C.GPGETFLOATI_VNV)(getProcAddr("glGetFloati_vNV"))
	gpGetFloati_vOES = (C.GPGETFLOATI_VOES)(getProcAddr("glGetFloati_vOES"))
	gpGetFloatv = (C.GPGETFLOATV)(getProcAddr("glGetFloatv"))
	if gpGetFloatv == nil {
		return errors.New("glGetFloatv")
	}
	gpGetFragDataIndexEXT = (C.GPGETFRAGDATAINDEXEXT)(getProcAddr("glGetFragDataIndexEXT"))
	gpGetFragDataLocation = (C.GPGETFRAGDATALOCATION)(getProcAddr("glGetFragDataLocation"))
	if gpGetFragDataLocation == nil {
		return errors.New("glGetFragDataLocation")
	}
	gpGetFramebufferAttachmentParameteriv = (C.GPGETFRAMEBUFFERATTACHMENTPARAMETERIV)(getProcAddr("glGetFramebufferAttachmentParameteriv"))
	if gpGetFramebufferAttachmentParameteriv == nil {
		return errors.New("glGetFramebufferAttachmentParameteriv")
	}
	gpGetFramebufferParameteriv = (C.GPGETFRAMEBUFFERPARAMETERIV)(getProcAddr("glGetFramebufferParameteriv"))
	if gpGetFramebufferParameteriv == nil {
		return errors.New("glGetFramebufferParameteriv")
	}
	gpGetFramebufferPixelLocalStorageSizeEXT = (C.GPGETFRAMEBUFFERPIXELLOCALSTORAGESIZEEXT)(getProcAddr("glGetFramebufferPixelLocalStorageSizeEXT"))
	gpGetGraphicsResetStatus = (C.GPGETGRAPHICSRESETSTATUS)(getProcAddr("glGetGraphicsResetStatus"))
	gpGetGraphicsResetStatusEXT = (C.GPGETGRAPHICSRESETSTATUSEXT)(getProcAddr("glGetGraphicsResetStatusEXT"))
	gpGetGraphicsResetStatusKHR = (C.GPGETGRAPHICSRESETSTATUSKHR)(getProcAddr("glGetGraphicsResetStatusKHR"))
	gpGetImageHandleNV = (C.GPGETIMAGEHANDLENV)(getProcAddr("glGetImageHandleNV"))
	gpGetInteger64i_v = (C.GPGETINTEGER64I_V)(getProcAddr("glGetInteger64i_v"))
	if gpGetInteger64i_v == nil {
		return errors.New("glGetInteger64i_v")
	}
	gpGetInteger64v = (C.GPGETINTEGER64V)(getProcAddr("glGetInteger64v"))
	if gpGetInteger64v == nil {
		return errors.New("glGetInteger64v")
	}
	gpGetInteger64vAPPLE = (C.GPGETINTEGER64VAPPLE)(getProcAddr("glGetInteger64vAPPLE"))
	gpGetIntegeri_v = (C.GPGETINTEGERI_V)(getProcAddr("glGetIntegeri_v"))
	if gpGetIntegeri_v == nil {
		return errors.New("glGetIntegeri_v")
	}
	gpGetIntegeri_vEXT = (C.GPGETINTEGERI_VEXT)(getProcAddr("glGetIntegeri_vEXT"))
	gpGetIntegerv = (C.GPGETINTEGERV)(getProcAddr("glGetIntegerv"))
	if gpGetIntegerv == nil {
		return errors.New("glGetIntegerv")
	}
	gpGetInternalformatSampleivNV = (C.GPGETINTERNALFORMATSAMPLEIVNV)(getProcAddr("glGetInternalformatSampleivNV"))
	gpGetInternalformativ = (C.GPGETINTERNALFORMATIV)(getProcAddr("glGetInternalformativ"))
	if gpGetInternalformativ == nil {
		return errors.New("glGetInternalformativ")
	}
	gpGetMemoryObjectParameterivEXT = (C.GPGETMEMORYOBJECTPARAMETERIVEXT)(getProcAddr("glGetMemoryObjectParameterivEXT"))
	gpGetMultisamplefv = (C.GPGETMULTISAMPLEFV)(getProcAddr("glGetMultisamplefv"))
	if gpGetMultisamplefv == nil {
		return errors.New("glGetMultisamplefv")
	}
	gpGetNextPerfQueryIdINTEL = (C.GPGETNEXTPERFQUERYIDINTEL)(getProcAddr("glGetNextPerfQueryIdINTEL"))
	gpGetObjectLabel = (C.GPGETOBJECTLABEL)(getProcAddr("glGetObjectLabel"))
	gpGetObjectLabelEXT = (C.GPGETOBJECTLABELEXT)(getProcAddr("glGetObjectLabelEXT"))
	gpGetObjectLabelKHR = (C.GPGETOBJECTLABELKHR)(getProcAddr("glGetObjectLabelKHR"))
	gpGetObjectPtrLabel = (C.GPGETOBJECTPTRLABEL)(getProcAddr("glGetObjectPtrLabel"))
	gpGetObjectPtrLabelKHR = (C.GPGETOBJECTPTRLABELKHR)(getProcAddr("glGetObjectPtrLabelKHR"))
	gpGetPathCommandsNV = (C.GPGETPATHCOMMANDSNV)(getProcAddr("glGetPathCommandsNV"))
	gpGetPathCoordsNV = (C.GPGETPATHCOORDSNV)(getProcAddr("glGetPathCoordsNV"))
	gpGetPathDashArrayNV = (C.GPGETPATHDASHARRAYNV)(getProcAddr("glGetPathDashArrayNV"))
	gpGetPathLengthNV = (C.GPGETPATHLENGTHNV)(getProcAddr("glGetPathLengthNV"))
	gpGetPathMetricRangeNV = (C.GPGETPATHMETRICRANGENV)(getProcAddr("glGetPathMetricRangeNV"))
	gpGetPathMetricsNV = (C.GPGETPATHMETRICSNV)(getProcAddr("glGetPathMetricsNV"))
	gpGetPathParameterfvNV = (C.GPGETPATHPARAMETERFVNV)(getProcAddr("glGetPathParameterfvNV"))
	gpGetPathParameterivNV = (C.GPGETPATHPARAMETERIVNV)(getProcAddr("glGetPathParameterivNV"))
	gpGetPathSpacingNV = (C.GPGETPATHSPACINGNV)(getProcAddr("glGetPathSpacingNV"))
	gpGetPerfCounterInfoINTEL = (C.GPGETPERFCOUNTERINFOINTEL)(getProcAddr("glGetPerfCounterInfoINTEL"))
	gpGetPerfMonitorCounterDataAMD = (C.GPGETPERFMONITORCOUNTERDATAAMD)(getProcAddr("glGetPerfMonitorCounterDataAMD"))
	gpGetPerfMonitorCounterInfoAMD = (C.GPGETPERFMONITORCOUNTERINFOAMD)(getProcAddr("glGetPerfMonitorCounterInfoAMD"))
	gpGetPerfMonitorCounterStringAMD = (C.GPGETPERFMONITORCOUNTERSTRINGAMD)(getProcAddr("glGetPerfMonitorCounterStringAMD"))
	gpGetPerfMonitorCountersAMD = (C.GPGETPERFMONITORCOUNTERSAMD)(getProcAddr("glGetPerfMonitorCountersAMD"))
	gpGetPerfMonitorGroupStringAMD = (C.GPGETPERFMONITORGROUPSTRINGAMD)(getProcAddr("glGetPerfMonitorGroupStringAMD"))
	gpGetPerfMonitorGroupsAMD = (C.GPGETPERFMONITORGROUPSAMD)(getProcAddr("glGetPerfMonitorGroupsAMD"))
	gpGetPerfQueryDataINTEL = (C.GPGETPERFQUERYDATAINTEL)(getProcAddr("glGetPerfQueryDataINTEL"))
	gpGetPerfQueryIdByNameINTEL = (C.GPGETPERFQUERYIDBYNAMEINTEL)(getProcAddr("glGetPerfQueryIdByNameINTEL"))
	gpGetPerfQueryInfoINTEL = (C.GPGETPERFQUERYINFOINTEL)(getProcAddr("glGetPerfQueryInfoINTEL"))
	gpGetPointerv = (C.GPGETPOINTERV)(getProcAddr("glGetPointerv"))
	gpGetPointervKHR = (C.GPGETPOINTERVKHR)(getProcAddr("glGetPointervKHR"))
	gpGetProgramBinary = (C.GPGETPROGRAMBINARY)(getProcAddr("glGetProgramBinary"))
	if gpGetProgramBinary == nil {
		return errors.New("glGetProgramBinary")
	}
	gpGetProgramBinaryOES = (C.GPGETPROGRAMBINARYOES)(getProcAddr("glGetProgramBinaryOES"))
	gpGetProgramInfoLog = (C.GPGETPROGRAMINFOLOG)(getProcAddr("glGetProgramInfoLog"))
	if gpGetProgramInfoLog == nil {
		return errors.New("glGetProgramInfoLog")
	}
	gpGetProgramInterfaceiv = (C.GPGETPROGRAMINTERFACEIV)(getProcAddr("glGetProgramInterfaceiv"))
	if gpGetProgramInterfaceiv == nil {
		return errors.New("glGetProgramInterfaceiv")
	}
	gpGetProgramPipelineInfoLog = (C.GPGETPROGRAMPIPELINEINFOLOG)(getProcAddr("glGetProgramPipelineInfoLog"))
	if gpGetProgramPipelineInfoLog == nil {
		return errors.New("glGetProgramPipelineInfoLog")
	}
	gpGetProgramPipelineInfoLogEXT = (C.GPGETPROGRAMPIPELINEINFOLOGEXT)(getProcAddr("glGetProgramPipelineInfoLogEXT"))
	gpGetProgramPipelineiv = (C.GPGETPROGRAMPIPELINEIV)(getProcAddr("glGetProgramPipelineiv"))
	if gpGetProgramPipelineiv == nil {
		return errors.New("glGetProgramPipelineiv")
	}
	gpGetProgramPipelineivEXT = (C.GPGETPROGRAMPIPELINEIVEXT)(getProcAddr("glGetProgramPipelineivEXT"))
	gpGetProgramResourceIndex = (C.GPGETPROGRAMRESOURCEINDEX)(getProcAddr("glGetProgramResourceIndex"))
	if gpGetProgramResourceIndex == nil {
		return errors.New("glGetProgramResourceIndex")
	}
	gpGetProgramResourceLocation = (C.GPGETPROGRAMRESOURCELOCATION)(getProcAddr("glGetProgramResourceLocation"))
	if gpGetProgramResourceLocation == nil {
		return errors.New("glGetProgramResourceLocation")
	}
	gpGetProgramResourceLocationIndexEXT = (C.GPGETPROGRAMRESOURCELOCATIONINDEXEXT)(getProcAddr("glGetProgramResourceLocationIndexEXT"))
	gpGetProgramResourceName = (C.GPGETPROGRAMRESOURCENAME)(getProcAddr("glGetProgramResourceName"))
	if gpGetProgramResourceName == nil {
		return errors.New("glGetProgramResourceName")
	}
	gpGetProgramResourcefvNV = (C.GPGETPROGRAMRESOURCEFVNV)(getProcAddr("glGetProgramResourcefvNV"))
	gpGetProgramResourceiv = (C.GPGETPROGRAMRESOURCEIV)(getProcAddr("glGetProgramResourceiv"))
	if gpGetProgramResourceiv == nil {
		return errors.New("glGetProgramResourceiv")
	}
	gpGetProgramiv = (C.GPGETPROGRAMIV)(getProcAddr("glGetProgramiv"))
	if gpGetProgramiv == nil {
		return errors.New("glGetProgramiv")
	}
	gpGetQueryObjecti64vEXT = (C.GPGETQUERYOBJECTI64VEXT)(getProcAddr("glGetQueryObjecti64vEXT"))
	gpGetQueryObjectivEXT = (C.GPGETQUERYOBJECTIVEXT)(getProcAddr("glGetQueryObjectivEXT"))
	gpGetQueryObjectui64vEXT = (C.GPGETQUERYOBJECTUI64VEXT)(getProcAddr("glGetQueryObjectui64vEXT"))
	gpGetQueryObjectuiv = (C.GPGETQUERYOBJECTUIV)(getProcAddr("glGetQueryObjectuiv"))
	if gpGetQueryObjectuiv == nil {
		return errors.New("glGetQueryObjectuiv")
	}
	gpGetQueryObjectuivEXT = (C.GPGETQUERYOBJECTUIVEXT)(getProcAddr("glGetQueryObjectuivEXT"))
	gpGetQueryiv = (C.GPGETQUERYIV)(getProcAddr("glGetQueryiv"))
	if gpGetQueryiv == nil {
		return errors.New("glGetQueryiv")
	}
	gpGetQueryivEXT = (C.GPGETQUERYIVEXT)(getProcAddr("glGetQueryivEXT"))
	gpGetRenderbufferParameteriv = (C.GPGETRENDERBUFFERPARAMETERIV)(getProcAddr("glGetRenderbufferParameteriv"))
	if gpGetRenderbufferParameteriv == nil {
		return errors.New("glGetRenderbufferParameteriv")
	}
	gpGetSamplerParameterIivEXT = (C.GPGETSAMPLERPARAMETERIIVEXT)(getProcAddr("glGetSamplerParameterIivEXT"))
	gpGetSamplerParameterIivOES = (C.GPGETSAMPLERPARAMETERIIVOES)(getProcAddr("glGetSamplerParameterIivOES"))
	gpGetSamplerParameterIuivEXT = (C.GPGETSAMPLERPARAMETERIUIVEXT)(getProcAddr("glGetSamplerParameterIuivEXT"))
	gpGetSamplerParameterIuivOES = (C.GPGETSAMPLERPARAMETERIUIVOES)(getProcAddr("glGetSamplerParameterIuivOES"))
	gpGetSamplerParameterfv = (C.GPGETSAMPLERPARAMETERFV)(getProcAddr("glGetSamplerParameterfv"))
	if gpGetSamplerParameterfv == nil {
		return errors.New("glGetSamplerParameterfv")
	}
	gpGetSamplerParameteriv = (C.GPGETSAMPLERPARAMETERIV)(getProcAddr("glGetSamplerParameteriv"))
	if gpGetSamplerParameteriv == nil {
		return errors.New("glGetSamplerParameteriv")
	}
	gpGetSemaphoreParameterui64vEXT = (C.GPGETSEMAPHOREPARAMETERUI64VEXT)(getProcAddr("glGetSemaphoreParameterui64vEXT"))
	gpGetShaderInfoLog = (C.GPGETSHADERINFOLOG)(getProcAddr("glGetShaderInfoLog"))
	if gpGetShaderInfoLog == nil {
		return errors.New("glGetShaderInfoLog")
	}
	gpGetShaderPrecisionFormat = (C.GPGETSHADERPRECISIONFORMAT)(getProcAddr("glGetShaderPrecisionFormat"))
	if gpGetShaderPrecisionFormat == nil {
		return errors.New("glGetShaderPrecisionFormat")
	}
	gpGetShaderSource = (C.GPGETSHADERSOURCE)(getProcAddr("glGetShaderSource"))
	if gpGetShaderSource == nil {
		return errors.New("glGetShaderSource")
	}
	gpGetShaderiv = (C.GPGETSHADERIV)(getProcAddr("glGetShaderiv"))
	if gpGetShaderiv == nil {
		return errors.New("glGetShaderiv")
	}
	gpGetString = (C.GPGETSTRING)(getProcAddr("glGetString"))
	if gpGetString == nil {
		return errors.New("glGetString")
	}
	gpGetStringi = (C.GPGETSTRINGI)(getProcAddr("glGetStringi"))
	if gpGetStringi == nil {
		return errors.New("glGetStringi")
	}
	gpGetSynciv = (C.GPGETSYNCIV)(getProcAddr("glGetSynciv"))
	if gpGetSynciv == nil {
		return errors.New("glGetSynciv")
	}
	gpGetSyncivAPPLE = (C.GPGETSYNCIVAPPLE)(getProcAddr("glGetSyncivAPPLE"))
	gpGetTexLevelParameterfv = (C.GPGETTEXLEVELPARAMETERFV)(getProcAddr("glGetTexLevelParameterfv"))
	if gpGetTexLevelParameterfv == nil {
		return errors.New("glGetTexLevelParameterfv")
	}
	gpGetTexLevelParameteriv = (C.GPGETTEXLEVELPARAMETERIV)(getProcAddr("glGetTexLevelParameteriv"))
	if gpGetTexLevelParameteriv == nil {
		return errors.New("glGetTexLevelParameteriv")
	}
	gpGetTexParameterIivEXT = (C.GPGETTEXPARAMETERIIVEXT)(getProcAddr("glGetTexParameterIivEXT"))
	gpGetTexParameterIivOES = (C.GPGETTEXPARAMETERIIVOES)(getProcAddr("glGetTexParameterIivOES"))
	gpGetTexParameterIuivEXT = (C.GPGETTEXPARAMETERIUIVEXT)(getProcAddr("glGetTexParameterIuivEXT"))
	gpGetTexParameterIuivOES = (C.GPGETTEXPARAMETERIUIVOES)(getProcAddr("glGetTexParameterIuivOES"))
	gpGetTexParameterfv = (C.GPGETTEXPARAMETERFV)(getProcAddr("glGetTexParameterfv"))
	if gpGetTexParameterfv == nil {
		return errors.New("glGetTexParameterfv")
	}
	gpGetTexParameteriv = (C.GPGETTEXPARAMETERIV)(getProcAddr("glGetTexParameteriv"))
	if gpGetTexParameteriv == nil {
		return errors.New("glGetTexParameteriv")
	}
	gpGetTextureHandleIMG = (C.GPGETTEXTUREHANDLEIMG)(getProcAddr("glGetTextureHandleIMG"))
	gpGetTextureHandleNV = (C.GPGETTEXTUREHANDLENV)(getProcAddr("glGetTextureHandleNV"))
	gpGetTextureSamplerHandleIMG = (C.GPGETTEXTURESAMPLERHANDLEIMG)(getProcAddr("glGetTextureSamplerHandleIMG"))
	gpGetTextureSamplerHandleNV = (C.GPGETTEXTURESAMPLERHANDLENV)(getProcAddr("glGetTextureSamplerHandleNV"))
	gpGetTransformFeedbackVarying = (C.GPGETTRANSFORMFEEDBACKVARYING)(getProcAddr("glGetTransformFeedbackVarying"))
	if gpGetTransformFeedbackVarying == nil {
		return errors.New("glGetTransformFeedbackVarying")
	}
	gpGetTranslatedShaderSourceANGLE = (C.GPGETTRANSLATEDSHADERSOURCEANGLE)(getProcAddr("glGetTranslatedShaderSourceANGLE"))
	gpGetUniformBlockIndex = (C.GPGETUNIFORMBLOCKINDEX)(getProcAddr("glGetUniformBlockIndex"))
	if gpGetUniformBlockIndex == nil {
		return errors.New("glGetUniformBlockIndex")
	}
	gpGetUniformIndices = (C.GPGETUNIFORMINDICES)(getProcAddr("glGetUniformIndices"))
	if gpGetUniformIndices == nil {
		return errors.New("glGetUniformIndices")
	}
	gpGetUniformLocation = (C.GPGETUNIFORMLOCATION)(getProcAddr("glGetUniformLocation"))
	if gpGetUniformLocation == nil {
		return errors.New("glGetUniformLocation")
	}
	gpGetUniformfv = (C.GPGETUNIFORMFV)(getProcAddr("glGetUniformfv"))
	if gpGetUniformfv == nil {
		return errors.New("glGetUniformfv")
	}
	gpGetUniformi64vNV = (C.GPGETUNIFORMI64VNV)(getProcAddr("glGetUniformi64vNV"))
	gpGetUniformiv = (C.GPGETUNIFORMIV)(getProcAddr("glGetUniformiv"))
	if gpGetUniformiv == nil {
		return errors.New("glGetUniformiv")
	}
	gpGetUniformuiv = (C.GPGETUNIFORMUIV)(getProcAddr("glGetUniformuiv"))
	if gpGetUniformuiv == nil {
		return errors.New("glGetUniformuiv")
	}
	gpGetUnsignedBytei_vEXT = (C.GPGETUNSIGNEDBYTEI_VEXT)(getProcAddr("glGetUnsignedBytei_vEXT"))
	gpGetUnsignedBytevEXT = (C.GPGETUNSIGNEDBYTEVEXT)(getProcAddr("glGetUnsignedBytevEXT"))
	gpGetVertexAttribIiv = (C.GPGETVERTEXATTRIBIIV)(getProcAddr("glGetVertexAttribIiv"))
	if gpGetVertexAttribIiv == nil {
		return errors.New("glGetVertexAttribIiv")
	}
	gpGetVertexAttribIuiv = (C.GPGETVERTEXATTRIBIUIV)(getProcAddr("glGetVertexAttribIuiv"))
	if gpGetVertexAttribIuiv == nil {
		return errors.New("glGetVertexAttribIuiv")
	}
	gpGetVertexAttribPointerv = (C.GPGETVERTEXATTRIBPOINTERV)(getProcAddr("glGetVertexAttribPointerv"))
	if gpGetVertexAttribPointerv == nil {
		return errors.New("glGetVertexAttribPointerv")
	}
	gpGetVertexAttribfv = (C.GPGETVERTEXATTRIBFV)(getProcAddr("glGetVertexAttribfv"))
	if gpGetVertexAttribfv == nil {
		return errors.New("glGetVertexAttribfv")
	}
	gpGetVertexAttribiv = (C.GPGETVERTEXATTRIBIV)(getProcAddr("glGetVertexAttribiv"))
	if gpGetVertexAttribiv == nil {
		return errors.New("glGetVertexAttribiv")
	}
	gpGetVkProcAddrNV = (C.GPGETVKPROCADDRNV)(getProcAddr("glGetVkProcAddrNV"))
	gpGetnUniformfv = (C.GPGETNUNIFORMFV)(getProcAddr("glGetnUniformfv"))
	gpGetnUniformfvEXT = (C.GPGETNUNIFORMFVEXT)(getProcAddr("glGetnUniformfvEXT"))
	gpGetnUniformfvKHR = (C.GPGETNUNIFORMFVKHR)(getProcAddr("glGetnUniformfvKHR"))
	gpGetnUniformiv = (C.GPGETNUNIFORMIV)(getProcAddr("glGetnUniformiv"))
	gpGetnUniformivEXT = (C.GPGETNUNIFORMIVEXT)(getProcAddr("glGetnUniformivEXT"))
	gpGetnUniformivKHR = (C.GPGETNUNIFORMIVKHR)(getProcAddr("glGetnUniformivKHR"))
	gpGetnUniformuiv = (C.GPGETNUNIFORMUIV)(getProcAddr("glGetnUniformuiv"))
	gpGetnUniformuivKHR = (C.GPGETNUNIFORMUIVKHR)(getProcAddr("glGetnUniformuivKHR"))
	gpHint = (C.GPHINT)(getProcAddr("glHint"))
	if gpHint == nil {
		return errors.New("glHint")
	}
	gpImportMemoryFdEXT = (C.GPIMPORTMEMORYFDEXT)(getProcAddr("glImportMemoryFdEXT"))
	gpImportMemoryWin32HandleEXT = (C.GPIMPORTMEMORYWIN32HANDLEEXT)(getProcAddr("glImportMemoryWin32HandleEXT"))
	gpImportMemoryWin32NameEXT = (C.GPIMPORTMEMORYWIN32NAMEEXT)(getProcAddr("glImportMemoryWin32NameEXT"))
	gpImportSemaphoreFdEXT = (C.GPIMPORTSEMAPHOREFDEXT)(getProcAddr("glImportSemaphoreFdEXT"))
	gpImportSemaphoreWin32HandleEXT = (C.GPIMPORTSEMAPHOREWIN32HANDLEEXT)(getProcAddr("glImportSemaphoreWin32HandleEXT"))
	gpImportSemaphoreWin32NameEXT = (C.GPIMPORTSEMAPHOREWIN32NAMEEXT)(getProcAddr("glImportSemaphoreWin32NameEXT"))
	gpInsertEventMarkerEXT = (C.GPINSERTEVENTMARKEREXT)(getProcAddr("glInsertEventMarkerEXT"))
	gpInterpolatePathsNV = (C.GPINTERPOLATEPATHSNV)(getProcAddr("glInterpolatePathsNV"))
	gpInvalidateFramebuffer = (C.GPINVALIDATEFRAMEBUFFER)(getProcAddr("glInvalidateFramebuffer"))
	if gpInvalidateFramebuffer == nil {
		return errors.New("glInvalidateFramebuffer")
	}
	gpInvalidateSubFramebuffer = (C.GPINVALIDATESUBFRAMEBUFFER)(getProcAddr("glInvalidateSubFramebuffer"))
	if gpInvalidateSubFramebuffer == nil {
		return errors.New("glInvalidateSubFramebuffer")
	}
	gpIsBuffer = (C.GPISBUFFER)(getProcAddr("glIsBuffer"))
	if gpIsBuffer == nil {
		return errors.New("glIsBuffer")
	}
	gpIsEnabled = (C.GPISENABLED)(getProcAddr("glIsEnabled"))
	if gpIsEnabled == nil {
		return errors.New("glIsEnabled")
	}
	gpIsEnablediEXT = (C.GPISENABLEDIEXT)(getProcAddr("glIsEnablediEXT"))
	gpIsEnablediNV = (C.GPISENABLEDINV)(getProcAddr("glIsEnablediNV"))
	gpIsEnablediOES = (C.GPISENABLEDIOES)(getProcAddr("glIsEnablediOES"))
	gpIsFenceNV = (C.GPISFENCENV)(getProcAddr("glIsFenceNV"))
	gpIsFramebuffer = (C.GPISFRAMEBUFFER)(getProcAddr("glIsFramebuffer"))
	if gpIsFramebuffer == nil {
		return errors.New("glIsFramebuffer")
	}
	gpIsImageHandleResidentNV = (C.GPISIMAGEHANDLERESIDENTNV)(getProcAddr("glIsImageHandleResidentNV"))
	gpIsMemoryObjectEXT = (C.GPISMEMORYOBJECTEXT)(getProcAddr("glIsMemoryObjectEXT"))
	gpIsPathNV = (C.GPISPATHNV)(getProcAddr("glIsPathNV"))
	gpIsPointInFillPathNV = (C.GPISPOINTINFILLPATHNV)(getProcAddr("glIsPointInFillPathNV"))
	gpIsPointInStrokePathNV = (C.GPISPOINTINSTROKEPATHNV)(getProcAddr("glIsPointInStrokePathNV"))
	gpIsProgram = (C.GPISPROGRAM)(getProcAddr("glIsProgram"))
	if gpIsProgram == nil {
		return errors.New("glIsProgram")
	}
	gpIsProgramPipeline = (C.GPISPROGRAMPIPELINE)(getProcAddr("glIsProgramPipeline"))
	if gpIsProgramPipeline == nil {
		return errors.New("glIsProgramPipeline")
	}
	gpIsProgramPipelineEXT = (C.GPISPROGRAMPIPELINEEXT)(getProcAddr("glIsProgramPipelineEXT"))
	gpIsQuery = (C.GPISQUERY)(getProcAddr("glIsQuery"))
	if gpIsQuery == nil {
		return errors.New("glIsQuery")
	}
	gpIsQueryEXT = (C.GPISQUERYEXT)(getProcAddr("glIsQueryEXT"))
	gpIsRenderbuffer = (C.GPISRENDERBUFFER)(getProcAddr("glIsRenderbuffer"))
	if gpIsRenderbuffer == nil {
		return errors.New("glIsRenderbuffer")
	}
	gpIsSampler = (C.GPISSAMPLER)(getProcAddr("glIsSampler"))
	if gpIsSampler == nil {
		return errors.New("glIsSampler")
	}
	gpIsSemaphoreEXT = (C.GPISSEMAPHOREEXT)(getProcAddr("glIsSemaphoreEXT"))
	gpIsShader = (C.GPISSHADER)(getProcAddr("glIsShader"))
	if gpIsShader == nil {
		return errors.New("glIsShader")
	}
	gpIsSync = (C.GPISSYNC)(getProcAddr("glIsSync"))
	if gpIsSync == nil {
		return errors.New("glIsSync")
	}
	gpIsSyncAPPLE = (C.GPISSYNCAPPLE)(getProcAddr("glIsSyncAPPLE"))
	gpIsTexture = (C.GPISTEXTURE)(getProcAddr("glIsTexture"))
	if gpIsTexture == nil {
		return errors.New("glIsTexture")
	}
	gpIsTextureHandleResidentNV = (C.GPISTEXTUREHANDLERESIDENTNV)(getProcAddr("glIsTextureHandleResidentNV"))
	gpIsTransformFeedback = (C.GPISTRANSFORMFEEDBACK)(getProcAddr("glIsTransformFeedback"))
	if gpIsTransformFeedback == nil {
		return errors.New("glIsTransformFeedback")
	}
	gpIsVertexArray = (C.GPISVERTEXARRAY)(getProcAddr("glIsVertexArray"))
	if gpIsVertexArray == nil {
		return errors.New("glIsVertexArray")
	}
	gpIsVertexArrayOES = (C.GPISVERTEXARRAYOES)(getProcAddr("glIsVertexArrayOES"))
	gpLabelObjectEXT = (C.GPLABELOBJECTEXT)(getProcAddr("glLabelObjectEXT"))
	gpLineWidth = (C.GPLINEWIDTH)(getProcAddr("glLineWidth"))
	if gpLineWidth == nil {
		return errors.New("glLineWidth")
	}
	gpLinkProgram = (C.GPLINKPROGRAM)(getProcAddr("glLinkProgram"))
	if gpLinkProgram == nil {
		return errors.New("glLinkProgram")
	}
	gpMakeImageHandleNonResidentNV = (C.GPMAKEIMAGEHANDLENONRESIDENTNV)(getProcAddr("glMakeImageHandleNonResidentNV"))
	gpMakeImageHandleResidentNV = (C.GPMAKEIMAGEHANDLERESIDENTNV)(getProcAddr("glMakeImageHandleResidentNV"))
	gpMakeTextureHandleNonResidentNV = (C.GPMAKETEXTUREHANDLENONRESIDENTNV)(getProcAddr("glMakeTextureHandleNonResidentNV"))
	gpMakeTextureHandleResidentNV = (C.GPMAKETEXTUREHANDLERESIDENTNV)(getProcAddr("glMakeTextureHandleResidentNV"))
	gpMapBufferOES = (C.GPMAPBUFFEROES)(getProcAddr("glMapBufferOES"))
	gpMapBufferRange = (C.GPMAPBUFFERRANGE)(getProcAddr("glMapBufferRange"))
	if gpMapBufferRange == nil {
		return errors.New("glMapBufferRange")
	}
	gpMapBufferRangeEXT = (C.GPMAPBUFFERRANGEEXT)(getProcAddr("glMapBufferRangeEXT"))
	gpMatrixFrustumEXT = (C.GPMATRIXFRUSTUMEXT)(getProcAddr("glMatrixFrustumEXT"))
	gpMatrixLoad3x2fNV = (C.GPMATRIXLOAD3X2FNV)(getProcAddr("glMatrixLoad3x2fNV"))
	gpMatrixLoad3x3fNV = (C.GPMATRIXLOAD3X3FNV)(getProcAddr("glMatrixLoad3x3fNV"))
	gpMatrixLoadIdentityEXT = (C.GPMATRIXLOADIDENTITYEXT)(getProcAddr("glMatrixLoadIdentityEXT"))
	gpMatrixLoadTranspose3x3fNV = (C.GPMATRIXLOADTRANSPOSE3X3FNV)(getProcAddr("glMatrixLoadTranspose3x3fNV"))
	gpMatrixLoadTransposedEXT = (C.GPMATRIXLOADTRANSPOSEDEXT)(getProcAddr("glMatrixLoadTransposedEXT"))
	gpMatrixLoadTransposefEXT = (C.GPMATRIXLOADTRANSPOSEFEXT)(getProcAddr("glMatrixLoadTransposefEXT"))
	gpMatrixLoaddEXT = (C.GPMATRIXLOADDEXT)(getProcAddr("glMatrixLoaddEXT"))
	gpMatrixLoadfEXT = (C.GPMATRIXLOADFEXT)(getProcAddr("glMatrixLoadfEXT"))
	gpMatrixMult3x2fNV = (C.GPMATRIXMULT3X2FNV)(getProcAddr("glMatrixMult3x2fNV"))
	gpMatrixMult3x3fNV = (C.GPMATRIXMULT3X3FNV)(getProcAddr("glMatrixMult3x3fNV"))
	gpMatrixMultTranspose3x3fNV = (C.GPMATRIXMULTTRANSPOSE3X3FNV)(getProcAddr("glMatrixMultTranspose3x3fNV"))
	gpMatrixMultTransposedEXT = (C.GPMATRIXMULTTRANSPOSEDEXT)(getProcAddr("glMatrixMultTransposedEXT"))
	gpMatrixMultTransposefEXT = (C.GPMATRIXMULTTRANSPOSEFEXT)(getProcAddr("glMatrixMultTransposefEXT"))
	gpMatrixMultdEXT = (C.GPMATRIXMULTDEXT)(getProcAddr("glMatrixMultdEXT"))
	gpMatrixMultfEXT = (C.GPMATRIXMULTFEXT)(getProcAddr("glMatrixMultfEXT"))
	gpMatrixOrthoEXT = (C.GPMATRIXORTHOEXT)(getProcAddr("glMatrixOrthoEXT"))
	gpMatrixPopEXT = (C.GPMATRIXPOPEXT)(getProcAddr("glMatrixPopEXT"))
	gpMatrixPushEXT = (C.GPMATRIXPUSHEXT)(getProcAddr("glMatrixPushEXT"))
	gpMatrixRotatedEXT = (C.GPMATRIXROTATEDEXT)(getProcAddr("glMatrixRotatedEXT"))
	gpMatrixRotatefEXT = (C.GPMATRIXROTATEFEXT)(getProcAddr("glMatrixRotatefEXT"))
	gpMatrixScaledEXT = (C.GPMATRIXSCALEDEXT)(getProcAddr("glMatrixScaledEXT"))
	gpMatrixScalefEXT = (C.GPMATRIXSCALEFEXT)(getProcAddr("glMatrixScalefEXT"))
	gpMatrixTranslatedEXT = (C.GPMATRIXTRANSLATEDEXT)(getProcAddr("glMatrixTranslatedEXT"))
	gpMatrixTranslatefEXT = (C.GPMATRIXTRANSLATEFEXT)(getProcAddr("glMatrixTranslatefEXT"))
	gpMaxShaderCompilerThreadsKHR = (C.GPMAXSHADERCOMPILERTHREADSKHR)(getProcAddr("glMaxShaderCompilerThreadsKHR"))
	gpMemoryBarrier = (C.GPMEMORYBARRIER)(getProcAddr("glMemoryBarrier"))
	if gpMemoryBarrier == nil {
		return errors.New("glMemoryBarrier")
	}
	gpMemoryBarrierByRegion = (C.GPMEMORYBARRIERBYREGION)(getProcAddr("glMemoryBarrierByRegion"))
	if gpMemoryBarrierByRegion == nil {
		return errors.New("glMemoryBarrierByRegion")
	}
	gpMemoryObjectParameterivEXT = (C.GPMEMORYOBJECTPARAMETERIVEXT)(getProcAddr("glMemoryObjectParameterivEXT"))
	gpMinSampleShadingOES = (C.GPMINSAMPLESHADINGOES)(getProcAddr("glMinSampleShadingOES"))
	gpMultiDrawArraysEXT = (C.GPMULTIDRAWARRAYSEXT)(getProcAddr("glMultiDrawArraysEXT"))
	gpMultiDrawArraysIndirectEXT = (C.GPMULTIDRAWARRAYSINDIRECTEXT)(getProcAddr("glMultiDrawArraysIndirectEXT"))
	gpMultiDrawElementsBaseVertexEXT = (C.GPMULTIDRAWELEMENTSBASEVERTEXEXT)(getProcAddr("glMultiDrawElementsBaseVertexEXT"))
	gpMultiDrawElementsEXT = (C.GPMULTIDRAWELEMENTSEXT)(getProcAddr("glMultiDrawElementsEXT"))
	gpMultiDrawElementsIndirectEXT = (C.GPMULTIDRAWELEMENTSINDIRECTEXT)(getProcAddr("glMultiDrawElementsIndirectEXT"))
	gpNamedBufferStorageExternalEXT = (C.GPNAMEDBUFFERSTORAGEEXTERNALEXT)(getProcAddr("glNamedBufferStorageExternalEXT"))
	gpNamedBufferStorageMemEXT = (C.GPNAMEDBUFFERSTORAGEMEMEXT)(getProcAddr("glNamedBufferStorageMemEXT"))
	gpNamedFramebufferSampleLocationsfvNV = (C.GPNAMEDFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glNamedFramebufferSampleLocationsfvNV"))
	gpObjectLabel = (C.GPOBJECTLABEL)(getProcAddr("glObjectLabel"))
	gpObjectLabelKHR = (C.GPOBJECTLABELKHR)(getProcAddr("glObjectLabelKHR"))
	gpObjectPtrLabel = (C.GPOBJECTPTRLABEL)(getProcAddr("glObjectPtrLabel"))
	gpObjectPtrLabelKHR = (C.GPOBJECTPTRLABELKHR)(getProcAddr("glObjectPtrLabelKHR"))
	gpPatchParameteriEXT = (C.GPPATCHPARAMETERIEXT)(getProcAddr("glPatchParameteriEXT"))
	gpPatchParameteriOES = (C.GPPATCHPARAMETERIOES)(getProcAddr("glPatchParameteriOES"))
	gpPathCommandsNV = (C.GPPATHCOMMANDSNV)(getProcAddr("glPathCommandsNV"))
	gpPathCoordsNV = (C.GPPATHCOORDSNV)(getProcAddr("glPathCoordsNV"))
	gpPathCoverDepthFuncNV = (C.GPPATHCOVERDEPTHFUNCNV)(getProcAddr("glPathCoverDepthFuncNV"))
	gpPathDashArrayNV = (C.GPPATHDASHARRAYNV)(getProcAddr("glPathDashArrayNV"))
	gpPathGlyphIndexArrayNV = (C.GPPATHGLYPHINDEXARRAYNV)(getProcAddr("glPathGlyphIndexArrayNV"))
	gpPathGlyphIndexRangeNV = (C.GPPATHGLYPHINDEXRANGENV)(getProcAddr("glPathGlyphIndexRangeNV"))
	gpPathGlyphRangeNV = (C.GPPATHGLYPHRANGENV)(getProcAddr("glPathGlyphRangeNV"))
	gpPathGlyphsNV = (C.GPPATHGLYPHSNV)(getProcAddr("glPathGlyphsNV"))
	gpPathMemoryGlyphIndexArrayNV = (C.GPPATHMEMORYGLYPHINDEXARRAYNV)(getProcAddr("glPathMemoryGlyphIndexArrayNV"))
	gpPathParameterfNV = (C.GPPATHPARAMETERFNV)(getProcAddr("glPathParameterfNV"))
	gpPathParameterfvNV = (C.GPPATHPARAMETERFVNV)(getProcAddr("glPathParameterfvNV"))
	gpPathParameteriNV = (C.GPPATHPARAMETERINV)(getProcAddr("glPathParameteriNV"))
	gpPathParameterivNV = (C.GPPATHPARAMETERIVNV)(getProcAddr("glPathParameterivNV"))
	gpPathStencilDepthOffsetNV = (C.GPPATHSTENCILDEPTHOFFSETNV)(getProcAddr("glPathStencilDepthOffsetNV"))
	gpPathStencilFuncNV = (C.GPPATHSTENCILFUNCNV)(getProcAddr("glPathStencilFuncNV"))
	gpPathStringNV = (C.GPPATHSTRINGNV)(getProcAddr("glPathStringNV"))
	gpPathSubCommandsNV = (C.GPPATHSUBCOMMANDSNV)(getProcAddr("glPathSubCommandsNV"))
	gpPathSubCoordsNV = (C.GPPATHSUBCOORDSNV)(getProcAddr("glPathSubCoordsNV"))
	gpPauseTransformFeedback = (C.GPPAUSETRANSFORMFEEDBACK)(getProcAddr("glPauseTransformFeedback"))
	if gpPauseTransformFeedback == nil {
		return errors.New("glPauseTransformFeedback")
	}
	gpPixelStorei = (C.GPPIXELSTOREI)(getProcAddr("glPixelStorei"))
	if gpPixelStorei == nil {
		return errors.New("glPixelStorei")
	}
	gpPointAlongPathNV = (C.GPPOINTALONGPATHNV)(getProcAddr("glPointAlongPathNV"))
	gpPolygonModeNV = (C.GPPOLYGONMODENV)(getProcAddr("glPolygonModeNV"))
	gpPolygonOffset = (C.GPPOLYGONOFFSET)(getProcAddr("glPolygonOffset"))
	if gpPolygonOffset == nil {
		return errors.New("glPolygonOffset")
	}
	gpPolygonOffsetClampEXT = (C.GPPOLYGONOFFSETCLAMPEXT)(getProcAddr("glPolygonOffsetClampEXT"))
	gpPopDebugGroup = (C.GPPOPDEBUGGROUP)(getProcAddr("glPopDebugGroup"))
	gpPopDebugGroupKHR = (C.GPPOPDEBUGGROUPKHR)(getProcAddr("glPopDebugGroupKHR"))
	gpPopGroupMarkerEXT = (C.GPPOPGROUPMARKEREXT)(getProcAddr("glPopGroupMarkerEXT"))
	gpPrimitiveBoundingBoxEXT = (C.GPPRIMITIVEBOUNDINGBOXEXT)(getProcAddr("glPrimitiveBoundingBoxEXT"))
	gpPrimitiveBoundingBoxOES = (C.GPPRIMITIVEBOUNDINGBOXOES)(getProcAddr("glPrimitiveBoundingBoxOES"))
	gpProgramBinary = (C.GPPROGRAMBINARY)(getProcAddr("glProgramBinary"))
	if gpProgramBinary == nil {
		return errors.New("glProgramBinary")
	}
	gpProgramBinaryOES = (C.GPPROGRAMBINARYOES)(getProcAddr("glProgramBinaryOES"))
	gpProgramParameteri = (C.GPPROGRAMPARAMETERI)(getProcAddr("glProgramParameteri"))
	if gpProgramParameteri == nil {
		return errors.New("glProgramParameteri")
	}
	gpProgramParameteriEXT = (C.GPPROGRAMPARAMETERIEXT)(getProcAddr("glProgramParameteriEXT"))
	gpProgramPathFragmentInputGenNV = (C.GPPROGRAMPATHFRAGMENTINPUTGENNV)(getProcAddr("glProgramPathFragmentInputGenNV"))
	gpProgramUniform1f = (C.GPPROGRAMUNIFORM1F)(getProcAddr("glProgramUniform1f"))
	if gpProgramUniform1f == nil {
		return errors.New("glProgramUniform1f")
	}
	gpProgramUniform1fEXT = (C.GPPROGRAMUNIFORM1FEXT)(getProcAddr("glProgramUniform1fEXT"))
	gpProgramUniform1fv = (C.GPPROGRAMUNIFORM1FV)(getProcAddr("glProgramUniform1fv"))
	if gpProgramUniform1fv == nil {
		return errors.New("glProgramUniform1fv")
	}
	gpProgramUniform1fvEXT = (C.GPPROGRAMUNIFORM1FVEXT)(getProcAddr("glProgramUniform1fvEXT"))
	gpProgramUniform1i = (C.GPPROGRAMUNIFORM1I)(getProcAddr("glProgramUniform1i"))
	if gpProgramUniform1i == nil {
		return errors.New("glProgramUniform1i")
	}
	gpProgramUniform1i64NV = (C.GPPROGRAMUNIFORM1I64NV)(getProcAddr("glProgramUniform1i64NV"))
	gpProgramUniform1i64vNV = (C.GPPROGRAMUNIFORM1I64VNV)(getProcAddr("glProgramUniform1i64vNV"))
	gpProgramUniform1iEXT = (C.GPPROGRAMUNIFORM1IEXT)(getProcAddr("glProgramUniform1iEXT"))
	gpProgramUniform1iv = (C.GPPROGRAMUNIFORM1IV)(getProcAddr("glProgramUniform1iv"))
	if gpProgramUniform1iv == nil {
		return errors.New("glProgramUniform1iv")
	}
	gpProgramUniform1ivEXT = (C.GPPROGRAMUNIFORM1IVEXT)(getProcAddr("glProgramUniform1ivEXT"))
	gpProgramUniform1ui = (C.GPPROGRAMUNIFORM1UI)(getProcAddr("glProgramUniform1ui"))
	if gpProgramUniform1ui == nil {
		return errors.New("glProgramUniform1ui")
	}
	gpProgramUniform1ui64NV = (C.GPPROGRAMUNIFORM1UI64NV)(getProcAddr("glProgramUniform1ui64NV"))
	gpProgramUniform1ui64vNV = (C.GPPROGRAMUNIFORM1UI64VNV)(getProcAddr("glProgramUniform1ui64vNV"))
	gpProgramUniform1uiEXT = (C.GPPROGRAMUNIFORM1UIEXT)(getProcAddr("glProgramUniform1uiEXT"))
	gpProgramUniform1uiv = (C.GPPROGRAMUNIFORM1UIV)(getProcAddr("glProgramUniform1uiv"))
	if gpProgramUniform1uiv == nil {
		return errors.New("glProgramUniform1uiv")
	}
	gpProgramUniform1uivEXT = (C.GPPROGRAMUNIFORM1UIVEXT)(getProcAddr("glProgramUniform1uivEXT"))
	gpProgramUniform2f = (C.GPPROGRAMUNIFORM2F)(getProcAddr("glProgramUniform2f"))
	if gpProgramUniform2f == nil {
		return errors.New("glProgramUniform2f")
	}
	gpProgramUniform2fEXT = (C.GPPROGRAMUNIFORM2FEXT)(getProcAddr("glProgramUniform2fEXT"))
	gpProgramUniform2fv = (C.GPPROGRAMUNIFORM2FV)(getProcAddr("glProgramUniform2fv"))
	if gpProgramUniform2fv == nil {
		return errors.New("glProgramUniform2fv")
	}
	gpProgramUniform2fvEXT = (C.GPPROGRAMUNIFORM2FVEXT)(getProcAddr("glProgramUniform2fvEXT"))
	gpProgramUniform2i = (C.GPPROGRAMUNIFORM2I)(getProcAddr("glProgramUniform2i"))
	if gpProgramUniform2i == nil {
		return errors.New("glProgramUniform2i")
	}
	gpProgramUniform2i64NV = (C.GPPROGRAMUNIFORM2I64NV)(getProcAddr("glProgramUniform2i64NV"))
	gpProgramUniform2i64vNV = (C.GPPROGRAMUNIFORM2I64VNV)(getProcAddr("glProgramUniform2i64vNV"))
	gpProgramUniform2iEXT = (C.GPPROGRAMUNIFORM2IEXT)(getProcAddr("glProgramUniform2iEXT"))
	gpProgramUniform2iv = (C.GPPROGRAMUNIFORM2IV)(getProcAddr("glProgramUniform2iv"))
	if gpProgramUniform2iv == nil {
		return errors.New("glProgramUniform2iv")
	}
	gpProgramUniform2ivEXT = (C.GPPROGRAMUNIFORM2IVEXT)(getProcAddr("glProgramUniform2ivEXT"))
	gpProgramUniform2ui = (C.GPPROGRAMUNIFORM2UI)(getProcAddr("glProgramUniform2ui"))
	if gpProgramUniform2ui == nil {
		return errors.New("glProgramUniform2ui")
	}
	gpProgramUniform2ui64NV = (C.GPPROGRAMUNIFORM2UI64NV)(getProcAddr("glProgramUniform2ui64NV"))
	gpProgramUniform2ui64vNV = (C.GPPROGRAMUNIFORM2UI64VNV)(getProcAddr("glProgramUniform2ui64vNV"))
	gpProgramUniform2uiEXT = (C.GPPROGRAMUNIFORM2UIEXT)(getProcAddr("glProgramUniform2uiEXT"))
	gpProgramUniform2uiv = (C.GPPROGRAMUNIFORM2UIV)(getProcAddr("glProgramUniform2uiv"))
	if gpProgramUniform2uiv == nil {
		return errors.New("glProgramUniform2uiv")
	}
	gpProgramUniform2uivEXT = (C.GPPROGRAMUNIFORM2UIVEXT)(getProcAddr("glProgramUniform2uivEXT"))
	gpProgramUniform3f = (C.GPPROGRAMUNIFORM3F)(getProcAddr("glProgramUniform3f"))
	if gpProgramUniform3f == nil {
		return errors.New("glProgramUniform3f")
	}
	gpProgramUniform3fEXT = (C.GPPROGRAMUNIFORM3FEXT)(getProcAddr("glProgramUniform3fEXT"))
	gpProgramUniform3fv = (C.GPPROGRAMUNIFORM3FV)(getProcAddr("glProgramUniform3fv"))
	if gpProgramUniform3fv == nil {
		return errors.New("glProgramUniform3fv")
	}
	gpProgramUniform3fvEXT = (C.GPPROGRAMUNIFORM3FVEXT)(getProcAddr("glProgramUniform3fvEXT"))
	gpProgramUniform3i = (C.GPPROGRAMUNIFORM3I)(getProcAddr("glProgramUniform3i"))
	if gpProgramUniform3i == nil {
		return errors.New("glProgramUniform3i")
	}
	gpProgramUniform3i64NV = (C.GPPROGRAMUNIFORM3I64NV)(getProcAddr("glProgramUniform3i64NV"))
	gpProgramUniform3i64vNV = (C.GPPROGRAMUNIFORM3I64VNV)(getProcAddr("glProgramUniform3i64vNV"))
	gpProgramUniform3iEXT = (C.GPPROGRAMUNIFORM3IEXT)(getProcAddr("glProgramUniform3iEXT"))
	gpProgramUniform3iv = (C.GPPROGRAMUNIFORM3IV)(getProcAddr("glProgramUniform3iv"))
	if gpProgramUniform3iv == nil {
		return errors.New("glProgramUniform3iv")
	}
	gpProgramUniform3ivEXT = (C.GPPROGRAMUNIFORM3IVEXT)(getProcAddr("glProgramUniform3ivEXT"))
	gpProgramUniform3ui = (C.GPPROGRAMUNIFORM3UI)(getProcAddr("glProgramUniform3ui"))
	if gpProgramUniform3ui == nil {
		return errors.New("glProgramUniform3ui")
	}
	gpProgramUniform3ui64NV = (C.GPPROGRAMUNIFORM3UI64NV)(getProcAddr("glProgramUniform3ui64NV"))
	gpProgramUniform3ui64vNV = (C.GPPROGRAMUNIFORM3UI64VNV)(getProcAddr("glProgramUniform3ui64vNV"))
	gpProgramUniform3uiEXT = (C.GPPROGRAMUNIFORM3UIEXT)(getProcAddr("glProgramUniform3uiEXT"))
	gpProgramUniform3uiv = (C.GPPROGRAMUNIFORM3UIV)(getProcAddr("glProgramUniform3uiv"))
	if gpProgramUniform3uiv == nil {
		return errors.New("glProgramUniform3uiv")
	}
	gpProgramUniform3uivEXT = (C.GPPROGRAMUNIFORM3UIVEXT)(getProcAddr("glProgramUniform3uivEXT"))
	gpProgramUniform4f = (C.GPPROGRAMUNIFORM4F)(getProcAddr("glProgramUniform4f"))
	if gpProgramUniform4f == nil {
		return errors.New("glProgramUniform4f")
	}
	gpProgramUniform4fEXT = (C.GPPROGRAMUNIFORM4FEXT)(getProcAddr("glProgramUniform4fEXT"))
	gpProgramUniform4fv = (C.GPPROGRAMUNIFORM4FV)(getProcAddr("glProgramUniform4fv"))
	if gpProgramUniform4fv == nil {
		return errors.New("glProgramUniform4fv")
	}
	gpProgramUniform4fvEXT = (C.GPPROGRAMUNIFORM4FVEXT)(getProcAddr("glProgramUniform4fvEXT"))
	gpProgramUniform4i = (C.GPPROGRAMUNIFORM4I)(getProcAddr("glProgramUniform4i"))
	if gpProgramUniform4i == nil {
		return errors.New("glProgramUniform4i")
	}
	gpProgramUniform4i64NV = (C.GPPROGRAMUNIFORM4I64NV)(getProcAddr("glProgramUniform4i64NV"))
	gpProgramUniform4i64vNV = (C.GPPROGRAMUNIFORM4I64VNV)(getProcAddr("glProgramUniform4i64vNV"))
	gpProgramUniform4iEXT = (C.GPPROGRAMUNIFORM4IEXT)(getProcAddr("glProgramUniform4iEXT"))
	gpProgramUniform4iv = (C.GPPROGRAMUNIFORM4IV)(getProcAddr("glProgramUniform4iv"))
	if gpProgramUniform4iv == nil {
		return errors.New("glProgramUniform4iv")
	}
	gpProgramUniform4ivEXT = (C.GPPROGRAMUNIFORM4IVEXT)(getProcAddr("glProgramUniform4ivEXT"))
	gpProgramUniform4ui = (C.GPPROGRAMUNIFORM4UI)(getProcAddr("glProgramUniform4ui"))
	if gpProgramUniform4ui == nil {
		return errors.New("glProgramUniform4ui")
	}
	gpProgramUniform4ui64NV = (C.GPPROGRAMUNIFORM4UI64NV)(getProcAddr("glProgramUniform4ui64NV"))
	gpProgramUniform4ui64vNV = (C.GPPROGRAMUNIFORM4UI64VNV)(getProcAddr("glProgramUniform4ui64vNV"))
	gpProgramUniform4uiEXT = (C.GPPROGRAMUNIFORM4UIEXT)(getProcAddr("glProgramUniform4uiEXT"))
	gpProgramUniform4uiv = (C.GPPROGRAMUNIFORM4UIV)(getProcAddr("glProgramUniform4uiv"))
	if gpProgramUniform4uiv == nil {
		return errors.New("glProgramUniform4uiv")
	}
	gpProgramUniform4uivEXT = (C.GPPROGRAMUNIFORM4UIVEXT)(getProcAddr("glProgramUniform4uivEXT"))
	gpProgramUniformHandleui64IMG = (C.GPPROGRAMUNIFORMHANDLEUI64IMG)(getProcAddr("glProgramUniformHandleui64IMG"))
	gpProgramUniformHandleui64NV = (C.GPPROGRAMUNIFORMHANDLEUI64NV)(getProcAddr("glProgramUniformHandleui64NV"))
	gpProgramUniformHandleui64vIMG = (C.GPPROGRAMUNIFORMHANDLEUI64VIMG)(getProcAddr("glProgramUniformHandleui64vIMG"))
	gpProgramUniformHandleui64vNV = (C.GPPROGRAMUNIFORMHANDLEUI64VNV)(getProcAddr("glProgramUniformHandleui64vNV"))
	gpProgramUniformMatrix2fv = (C.GPPROGRAMUNIFORMMATRIX2FV)(getProcAddr("glProgramUniformMatrix2fv"))
	if gpProgramUniformMatrix2fv == nil {
		return errors.New("glProgramUniformMatrix2fv")
	}
	gpProgramUniformMatrix2fvEXT = (C.GPPROGRAMUNIFORMMATRIX2FVEXT)(getProcAddr("glProgramUniformMatrix2fvEXT"))
	gpProgramUniformMatrix2x3fv = (C.GPPROGRAMUNIFORMMATRIX2X3FV)(getProcAddr("glProgramUniformMatrix2x3fv"))
	if gpProgramUniformMatrix2x3fv == nil {
		return errors.New("glProgramUniformMatrix2x3fv")
	}
	gpProgramUniformMatrix2x3fvEXT = (C.GPPROGRAMUNIFORMMATRIX2X3FVEXT)(getProcAddr("glProgramUniformMatrix2x3fvEXT"))
	gpProgramUniformMatrix2x4fv = (C.GPPROGRAMUNIFORMMATRIX2X4FV)(getProcAddr("glProgramUniformMatrix2x4fv"))
	if gpProgramUniformMatrix2x4fv == nil {
		return errors.New("glProgramUniformMatrix2x4fv")
	}
	gpProgramUniformMatrix2x4fvEXT = (C.GPPROGRAMUNIFORMMATRIX2X4FVEXT)(getProcAddr("glProgramUniformMatrix2x4fvEXT"))
	gpProgramUniformMatrix3fv = (C.GPPROGRAMUNIFORMMATRIX3FV)(getProcAddr("glProgramUniformMatrix3fv"))
	if gpProgramUniformMatrix3fv == nil {
		return errors.New("glProgramUniformMatrix3fv")
	}
	gpProgramUniformMatrix3fvEXT = (C.GPPROGRAMUNIFORMMATRIX3FVEXT)(getProcAddr("glProgramUniformMatrix3fvEXT"))
	gpProgramUniformMatrix3x2fv = (C.GPPROGRAMUNIFORMMATRIX3X2FV)(getProcAddr("glProgramUniformMatrix3x2fv"))
	if gpProgramUniformMatrix3x2fv == nil {
		return errors.New("glProgramUniformMatrix3x2fv")
	}
	gpProgramUniformMatrix3x2fvEXT = (C.GPPROGRAMUNIFORMMATRIX3X2FVEXT)(getProcAddr("glProgramUniformMatrix3x2fvEXT"))
	gpProgramUniformMatrix3x4fv = (C.GPPROGRAMUNIFORMMATRIX3X4FV)(getProcAddr("glProgramUniformMatrix3x4fv"))
	if gpProgramUniformMatrix3x4fv == nil {
		return errors.New("glProgramUniformMatrix3x4fv")
	}
	gpProgramUniformMatrix3x4fvEXT = (C.GPPROGRAMUNIFORMMATRIX3X4FVEXT)(getProcAddr("glProgramUniformMatrix3x4fvEXT"))
	gpProgramUniformMatrix4fv = (C.GPPROGRAMUNIFORMMATRIX4FV)(getProcAddr("glProgramUniformMatrix4fv"))
	if gpProgramUniformMatrix4fv == nil {
		return errors.New("glProgramUniformMatrix4fv")
	}
	gpProgramUniformMatrix4fvEXT = (C.GPPROGRAMUNIFORMMATRIX4FVEXT)(getProcAddr("glProgramUniformMatrix4fvEXT"))
	gpProgramUniformMatrix4x2fv = (C.GPPROGRAMUNIFORMMATRIX4X2FV)(getProcAddr("glProgramUniformMatrix4x2fv"))
	if gpProgramUniformMatrix4x2fv == nil {
		return errors.New("glProgramUniformMatrix4x2fv")
	}
	gpProgramUniformMatrix4x2fvEXT = (C.GPPROGRAMUNIFORMMATRIX4X2FVEXT)(getProcAddr("glProgramUniformMatrix4x2fvEXT"))
	gpProgramUniformMatrix4x3fv = (C.GPPROGRAMUNIFORMMATRIX4X3FV)(getProcAddr("glProgramUniformMatrix4x3fv"))
	if gpProgramUniformMatrix4x3fv == nil {
		return errors.New("glProgramUniformMatrix4x3fv")
	}
	gpProgramUniformMatrix4x3fvEXT = (C.GPPROGRAMUNIFORMMATRIX4X3FVEXT)(getProcAddr("glProgramUniformMatrix4x3fvEXT"))
	gpPushDebugGroup = (C.GPPUSHDEBUGGROUP)(getProcAddr("glPushDebugGroup"))
	gpPushDebugGroupKHR = (C.GPPUSHDEBUGGROUPKHR)(getProcAddr("glPushDebugGroupKHR"))
	gpPushGroupMarkerEXT = (C.GPPUSHGROUPMARKEREXT)(getProcAddr("glPushGroupMarkerEXT"))
	gpQueryCounterEXT = (C.GPQUERYCOUNTEREXT)(getProcAddr("glQueryCounterEXT"))
	gpRasterSamplesEXT = (C.GPRASTERSAMPLESEXT)(getProcAddr("glRasterSamplesEXT"))
	gpReadBuffer = (C.GPREADBUFFER)(getProcAddr("glReadBuffer"))
	if gpReadBuffer == nil {
		return errors.New("glReadBuffer")
	}
	gpReadBufferIndexedEXT = (C.GPREADBUFFERINDEXEDEXT)(getProcAddr("glReadBufferIndexedEXT"))
	gpReadBufferNV = (C.GPREADBUFFERNV)(getProcAddr("glReadBufferNV"))
	gpReadPixels = (C.GPREADPIXELS)(getProcAddr("glReadPixels"))
	if gpReadPixels == nil {
		return errors.New("glReadPixels")
	}
	gpReadnPixels = (C.GPREADNPIXELS)(getProcAddr("glReadnPixels"))
	gpReadnPixelsEXT = (C.GPREADNPIXELSEXT)(getProcAddr("glReadnPixelsEXT"))
	gpReadnPixelsKHR = (C.GPREADNPIXELSKHR)(getProcAddr("glReadnPixelsKHR"))
	gpReleaseKeyedMutexWin32EXT = (C.GPRELEASEKEYEDMUTEXWIN32EXT)(getProcAddr("glReleaseKeyedMutexWin32EXT"))
	gpReleaseShaderCompiler = (C.GPRELEASESHADERCOMPILER)(getProcAddr("glReleaseShaderCompiler"))
	if gpReleaseShaderCompiler == nil {
		return errors.New("glReleaseShaderCompiler")
	}
	gpRenderbufferStorage = (C.GPRENDERBUFFERSTORAGE)(getProcAddr("glRenderbufferStorage"))
	if gpRenderbufferStorage == nil {
		return errors.New("glRenderbufferStorage")
	}
	gpRenderbufferStorageMultisample = (C.GPRENDERBUFFERSTORAGEMULTISAMPLE)(getProcAddr("glRenderbufferStorageMultisample"))
	if gpRenderbufferStorageMultisample == nil {
		return errors.New("glRenderbufferStorageMultisample")
	}
	gpRenderbufferStorageMultisampleANGLE = (C.GPRENDERBUFFERSTORAGEMULTISAMPLEANGLE)(getProcAddr("glRenderbufferStorageMultisampleANGLE"))
	gpRenderbufferStorageMultisampleAPPLE = (C.GPRENDERBUFFERSTORAGEMULTISAMPLEAPPLE)(getProcAddr("glRenderbufferStorageMultisampleAPPLE"))
	gpRenderbufferStorageMultisampleEXT = (C.GPRENDERBUFFERSTORAGEMULTISAMPLEEXT)(getProcAddr("glRenderbufferStorageMultisampleEXT"))
	gpRenderbufferStorageMultisampleIMG = (C.GPRENDERBUFFERSTORAGEMULTISAMPLEIMG)(getProcAddr("glRenderbufferStorageMultisampleIMG"))
	gpRenderbufferStorageMultisampleNV = (C.GPRENDERBUFFERSTORAGEMULTISAMPLENV)(getProcAddr("glRenderbufferStorageMultisampleNV"))
	gpResolveDepthValuesNV = (C.GPRESOLVEDEPTHVALUESNV)(getProcAddr("glResolveDepthValuesNV"))
	gpResolveMultisampleFramebufferAPPLE = (C.GPRESOLVEMULTISAMPLEFRAMEBUFFERAPPLE)(getProcAddr("glResolveMultisampleFramebufferAPPLE"))
	gpResumeTransformFeedback = (C.GPRESUMETRANSFORMFEEDBACK)(getProcAddr("glResumeTransformFeedback"))
	if gpResumeTransformFeedback == nil {
		return errors.New("glResumeTransformFeedback")
	}
	gpSampleCoverage = (C.GPSAMPLECOVERAGE)(getProcAddr("glSampleCoverage"))
	if gpSampleCoverage == nil {
		return errors.New("glSampleCoverage")
	}
	gpSampleMaski = (C.GPSAMPLEMASKI)(getProcAddr("glSampleMaski"))
	if gpSampleMaski == nil {
		return errors.New("glSampleMaski")
	}
	gpSamplerParameterIivEXT = (C.GPSAMPLERPARAMETERIIVEXT)(getProcAddr("glSamplerParameterIivEXT"))
	gpSamplerParameterIivOES = (C.GPSAMPLERPARAMETERIIVOES)(getProcAddr("glSamplerParameterIivOES"))
	gpSamplerParameterIuivEXT = (C.GPSAMPLERPARAMETERIUIVEXT)(getProcAddr("glSamplerParameterIuivEXT"))
	gpSamplerParameterIuivOES = (C.GPSAMPLERPARAMETERIUIVOES)(getProcAddr("glSamplerParameterIuivOES"))
	gpSamplerParameterf = (C.GPSAMPLERPARAMETERF)(getProcAddr("glSamplerParameterf"))
	if gpSamplerParameterf == nil {
		return errors.New("glSamplerParameterf")
	}
	gpSamplerParameterfv = (C.GPSAMPLERPARAMETERFV)(getProcAddr("glSamplerParameterfv"))
	if gpSamplerParameterfv == nil {
		return errors.New("glSamplerParameterfv")
	}
	gpSamplerParameteri = (C.GPSAMPLERPARAMETERI)(getProcAddr("glSamplerParameteri"))
	if gpSamplerParameteri == nil {
		return errors.New("glSamplerParameteri")
	}
	gpSamplerParameteriv = (C.GPSAMPLERPARAMETERIV)(getProcAddr("glSamplerParameteriv"))
	if gpSamplerParameteriv == nil {
		return errors.New("glSamplerParameteriv")
	}
	gpScissor = (C.GPSCISSOR)(getProcAddr("glScissor"))
	if gpScissor == nil {
		return errors.New("glScissor")
	}
	gpScissorArrayvNV = (C.GPSCISSORARRAYVNV)(getProcAddr("glScissorArrayvNV"))
	gpScissorArrayvOES = (C.GPSCISSORARRAYVOES)(getProcAddr("glScissorArrayvOES"))
	gpScissorIndexedNV = (C.GPSCISSORINDEXEDNV)(getProcAddr("glScissorIndexedNV"))
	gpScissorIndexedOES = (C.GPSCISSORINDEXEDOES)(getProcAddr("glScissorIndexedOES"))
	gpScissorIndexedvNV = (C.GPSCISSORINDEXEDVNV)(getProcAddr("glScissorIndexedvNV"))
	gpScissorIndexedvOES = (C.GPSCISSORINDEXEDVOES)(getProcAddr("glScissorIndexedvOES"))
	gpSelectPerfMonitorCountersAMD = (C.GPSELECTPERFMONITORCOUNTERSAMD)(getProcAddr("glSelectPerfMonitorCountersAMD"))
	gpSemaphoreParameterui64vEXT = (C.GPSEMAPHOREPARAMETERUI64VEXT)(getProcAddr("glSemaphoreParameterui64vEXT"))
	gpSetFenceNV = (C.GPSETFENCENV)(getProcAddr("glSetFenceNV"))
	gpShaderBinary = (C.GPSHADERBINARY)(getProcAddr("glShaderBinary"))
	if gpShaderBinary == nil {
		return errors.New("glShaderBinary")
	}
	gpShaderSource = (C.GPSHADERSOURCE)(getProcAddr("glShaderSource"))
	if gpShaderSource == nil {
		return errors.New("glShaderSource")
	}
	gpSignalSemaphoreEXT = (C.GPSIGNALSEMAPHOREEXT)(getProcAddr("glSignalSemaphoreEXT"))
	gpSignalVkFenceNV = (C.GPSIGNALVKFENCENV)(getProcAddr("glSignalVkFenceNV"))
	gpSignalVkSemaphoreNV = (C.GPSIGNALVKSEMAPHORENV)(getProcAddr("glSignalVkSemaphoreNV"))
	gpStartTilingQCOM = (C.GPSTARTTILINGQCOM)(getProcAddr("glStartTilingQCOM"))
	gpStencilFillPathInstancedNV = (C.GPSTENCILFILLPATHINSTANCEDNV)(getProcAddr("glStencilFillPathInstancedNV"))
	gpStencilFillPathNV = (C.GPSTENCILFILLPATHNV)(getProcAddr("glStencilFillPathNV"))
	gpStencilFunc = (C.GPSTENCILFUNC)(getProcAddr("glStencilFunc"))
	if gpStencilFunc == nil {
		return errors.New("glStencilFunc")
	}
	gpStencilFuncSeparate = (C.GPSTENCILFUNCSEPARATE)(getProcAddr("glStencilFuncSeparate"))
	if gpStencilFuncSeparate == nil {
		return errors.New("glStencilFuncSeparate")
	}
	gpStencilMask = (C.GPSTENCILMASK)(getProcAddr("glStencilMask"))
	if gpStencilMask == nil {
		return errors.New("glStencilMask")
	}
	gpStencilMaskSeparate = (C.GPSTENCILMASKSEPARATE)(getProcAddr("glStencilMaskSeparate"))
	if gpStencilMaskSeparate == nil {
		return errors.New("glStencilMaskSeparate")
	}
	gpStencilOp = (C.GPSTENCILOP)(getProcAddr("glStencilOp"))
	if gpStencilOp == nil {
		return errors.New("glStencilOp")
	}
	gpStencilOpSeparate = (C.GPSTENCILOPSEPARATE)(getProcAddr("glStencilOpSeparate"))
	if gpStencilOpSeparate == nil {
		return errors.New("glStencilOpSeparate")
	}
	gpStencilStrokePathInstancedNV = (C.GPSTENCILSTROKEPATHINSTANCEDNV)(getProcAddr("glStencilStrokePathInstancedNV"))
	gpStencilStrokePathNV = (C.GPSTENCILSTROKEPATHNV)(getProcAddr("glStencilStrokePathNV"))
	gpStencilThenCoverFillPathInstancedNV = (C.GPSTENCILTHENCOVERFILLPATHINSTANCEDNV)(getProcAddr("glStencilThenCoverFillPathInstancedNV"))
	gpStencilThenCoverFillPathNV = (C.GPSTENCILTHENCOVERFILLPATHNV)(getProcAddr("glStencilThenCoverFillPathNV"))
	gpStencilThenCoverStrokePathInstancedNV = (C.GPSTENCILTHENCOVERSTROKEPATHINSTANCEDNV)(getProcAddr("glStencilThenCoverStrokePathInstancedNV"))
	gpStencilThenCoverStrokePathNV = (C.GPSTENCILTHENCOVERSTROKEPATHNV)(getProcAddr("glStencilThenCoverStrokePathNV"))
	gpSubpixelPrecisionBiasNV = (C.GPSUBPIXELPRECISIONBIASNV)(getProcAddr("glSubpixelPrecisionBiasNV"))
	gpTestFenceNV = (C.GPTESTFENCENV)(getProcAddr("glTestFenceNV"))
	gpTexBufferEXT = (C.GPTEXBUFFEREXT)(getProcAddr("glTexBufferEXT"))
	gpTexBufferOES = (C.GPTEXBUFFEROES)(getProcAddr("glTexBufferOES"))
	gpTexBufferRangeEXT = (C.GPTEXBUFFERRANGEEXT)(getProcAddr("glTexBufferRangeEXT"))
	gpTexBufferRangeOES = (C.GPTEXBUFFERRANGEOES)(getProcAddr("glTexBufferRangeOES"))
	gpTexImage2D = (C.GPTEXIMAGE2D)(getProcAddr("glTexImage2D"))
	if gpTexImage2D == nil {
		return errors.New("glTexImage2D")
	}
	gpTexImage3D = (C.GPTEXIMAGE3D)(getProcAddr("glTexImage3D"))
	if gpTexImage3D == nil {
		return errors.New("glTexImage3D")
	}
	gpTexImage3DOES = (C.GPTEXIMAGE3DOES)(getProcAddr("glTexImage3DOES"))
	gpTexPageCommitmentEXT = (C.GPTEXPAGECOMMITMENTEXT)(getProcAddr("glTexPageCommitmentEXT"))
	gpTexParameterIivEXT = (C.GPTEXPARAMETERIIVEXT)(getProcAddr("glTexParameterIivEXT"))
	gpTexParameterIivOES = (C.GPTEXPARAMETERIIVOES)(getProcAddr("glTexParameterIivOES"))
	gpTexParameterIuivEXT = (C.GPTEXPARAMETERIUIVEXT)(getProcAddr("glTexParameterIuivEXT"))
	gpTexParameterIuivOES = (C.GPTEXPARAMETERIUIVOES)(getProcAddr("glTexParameterIuivOES"))
	gpTexParameterf = (C.GPTEXPARAMETERF)(getProcAddr("glTexParameterf"))
	if gpTexParameterf == nil {
		return errors.New("glTexParameterf")
	}
	gpTexParameterfv = (C.GPTEXPARAMETERFV)(getProcAddr("glTexParameterfv"))
	if gpTexParameterfv == nil {
		return errors.New("glTexParameterfv")
	}
	gpTexParameteri = (C.GPTEXPARAMETERI)(getProcAddr("glTexParameteri"))
	if gpTexParameteri == nil {
		return errors.New("glTexParameteri")
	}
	gpTexParameteriv = (C.GPTEXPARAMETERIV)(getProcAddr("glTexParameteriv"))
	if gpTexParameteriv == nil {
		return errors.New("glTexParameteriv")
	}
	gpTexStorage1DEXT = (C.GPTEXSTORAGE1DEXT)(getProcAddr("glTexStorage1DEXT"))
	gpTexStorage2D = (C.GPTEXSTORAGE2D)(getProcAddr("glTexStorage2D"))
	if gpTexStorage2D == nil {
		return errors.New("glTexStorage2D")
	}
	gpTexStorage2DEXT = (C.GPTEXSTORAGE2DEXT)(getProcAddr("glTexStorage2DEXT"))
	gpTexStorage2DMultisample = (C.GPTEXSTORAGE2DMULTISAMPLE)(getProcAddr("glTexStorage2DMultisample"))
	if gpTexStorage2DMultisample == nil {
		return errors.New("glTexStorage2DMultisample")
	}
	gpTexStorage3D = (C.GPTEXSTORAGE3D)(getProcAddr("glTexStorage3D"))
	if gpTexStorage3D == nil {
		return errors.New("glTexStorage3D")
	}
	gpTexStorage3DEXT = (C.GPTEXSTORAGE3DEXT)(getProcAddr("glTexStorage3DEXT"))
	gpTexStorage3DMultisampleOES = (C.GPTEXSTORAGE3DMULTISAMPLEOES)(getProcAddr("glTexStorage3DMultisampleOES"))
	gpTexStorageMem1DEXT = (C.GPTEXSTORAGEMEM1DEXT)(getProcAddr("glTexStorageMem1DEXT"))
	gpTexStorageMem2DEXT = (C.GPTEXSTORAGEMEM2DEXT)(getProcAddr("glTexStorageMem2DEXT"))
	gpTexStorageMem2DMultisampleEXT = (C.GPTEXSTORAGEMEM2DMULTISAMPLEEXT)(getProcAddr("glTexStorageMem2DMultisampleEXT"))
	gpTexStorageMem3DEXT = (C.GPTEXSTORAGEMEM3DEXT)(getProcAddr("glTexStorageMem3DEXT"))
	gpTexStorageMem3DMultisampleEXT = (C.GPTEXSTORAGEMEM3DMULTISAMPLEEXT)(getProcAddr("glTexStorageMem3DMultisampleEXT"))
	gpTexSubImage2D = (C.GPTEXSUBIMAGE2D)(getProcAddr("glTexSubImage2D"))
	if gpTexSubImage2D == nil {
		return errors.New("glTexSubImage2D")
	}
	gpTexSubImage3D = (C.GPTEXSUBIMAGE3D)(getProcAddr("glTexSubImage3D"))
	if gpTexSubImage3D == nil {
		return errors.New("glTexSubImage3D")
	}
	gpTexSubImage3DOES = (C.GPTEXSUBIMAGE3DOES)(getProcAddr("glTexSubImage3DOES"))
	gpTextureFoveationParametersQCOM = (C.GPTEXTUREFOVEATIONPARAMETERSQCOM)(getProcAddr("glTextureFoveationParametersQCOM"))
	gpTextureStorage1DEXT = (C.GPTEXTURESTORAGE1DEXT)(getProcAddr("glTextureStorage1DEXT"))
	gpTextureStorage2DEXT = (C.GPTEXTURESTORAGE2DEXT)(getProcAddr("glTextureStorage2DEXT"))
	gpTextureStorage3DEXT = (C.GPTEXTURESTORAGE3DEXT)(getProcAddr("glTextureStorage3DEXT"))
	gpTextureStorageMem1DEXT = (C.GPTEXTURESTORAGEMEM1DEXT)(getProcAddr("glTextureStorageMem1DEXT"))
	gpTextureStorageMem2DEXT = (C.GPTEXTURESTORAGEMEM2DEXT)(getProcAddr("glTextureStorageMem2DEXT"))
	gpTextureStorageMem2DMultisampleEXT = (C.GPTEXTURESTORAGEMEM2DMULTISAMPLEEXT)(getProcAddr("glTextureStorageMem2DMultisampleEXT"))
	gpTextureStorageMem3DEXT = (C.GPTEXTURESTORAGEMEM3DEXT)(getProcAddr("glTextureStorageMem3DEXT"))
	gpTextureStorageMem3DMultisampleEXT = (C.GPTEXTURESTORAGEMEM3DMULTISAMPLEEXT)(getProcAddr("glTextureStorageMem3DMultisampleEXT"))
	gpTextureViewEXT = (C.GPTEXTUREVIEWEXT)(getProcAddr("glTextureViewEXT"))
	gpTextureViewOES = (C.GPTEXTUREVIEWOES)(getProcAddr("glTextureViewOES"))
	gpTransformFeedbackVaryings = (C.GPTRANSFORMFEEDBACKVARYINGS)(getProcAddr("glTransformFeedbackVaryings"))
	if gpTransformFeedbackVaryings == nil {
		return errors.New("glTransformFeedbackVaryings")
	}
	gpTransformPathNV = (C.GPTRANSFORMPATHNV)(getProcAddr("glTransformPathNV"))
	gpUniform1f = (C.GPUNIFORM1F)(getProcAddr("glUniform1f"))
	if gpUniform1f == nil {
		return errors.New("glUniform1f")
	}
	gpUniform1fv = (C.GPUNIFORM1FV)(getProcAddr("glUniform1fv"))
	if gpUniform1fv == nil {
		return errors.New("glUniform1fv")
	}
	gpUniform1i = (C.GPUNIFORM1I)(getProcAddr("glUniform1i"))
	if gpUniform1i == nil {
		return errors.New("glUniform1i")
	}
	gpUniform1i64NV = (C.GPUNIFORM1I64NV)(getProcAddr("glUniform1i64NV"))
	gpUniform1i64vNV = (C.GPUNIFORM1I64VNV)(getProcAddr("glUniform1i64vNV"))
	gpUniform1iv = (C.GPUNIFORM1IV)(getProcAddr("glUniform1iv"))
	if gpUniform1iv == nil {
		return errors.New("glUniform1iv")
	}
	gpUniform1ui = (C.GPUNIFORM1UI)(getProcAddr("glUniform1ui"))
	if gpUniform1ui == nil {
		return errors.New("glUniform1ui")
	}
	gpUniform1ui64NV = (C.GPUNIFORM1UI64NV)(getProcAddr("glUniform1ui64NV"))
	gpUniform1ui64vNV = (C.GPUNIFORM1UI64VNV)(getProcAddr("glUniform1ui64vNV"))
	gpUniform1uiv = (C.GPUNIFORM1UIV)(getProcAddr("glUniform1uiv"))
	if gpUniform1uiv == nil {
		return errors.New("glUniform1uiv")
	}
	gpUniform2f = (C.GPUNIFORM2F)(getProcAddr("glUniform2f"))
	if gpUniform2f == nil {
		return errors.New("glUniform2f")
	}
	gpUniform2fv = (C.GPUNIFORM2FV)(getProcAddr("glUniform2fv"))
	if gpUniform2fv == nil {
		return errors.New("glUniform2fv")
	}
	gpUniform2i = (C.GPUNIFORM2I)(getProcAddr("glUniform2i"))
	if gpUniform2i == nil {
		return errors.New("glUniform2i")
	}
	gpUniform2i64NV = (C.GPUNIFORM2I64NV)(getProcAddr("glUniform2i64NV"))
	gpUniform2i64vNV = (C.GPUNIFORM2I64VNV)(getProcAddr("glUniform2i64vNV"))
	gpUniform2iv = (C.GPUNIFORM2IV)(getProcAddr("glUniform2iv"))
	if gpUniform2iv == nil {
		return errors.New("glUniform2iv")
	}
	gpUniform2ui = (C.GPUNIFORM2UI)(getProcAddr("glUniform2ui"))
	if gpUniform2ui == nil {
		return errors.New("glUniform2ui")
	}
	gpUniform2ui64NV = (C.GPUNIFORM2UI64NV)(getProcAddr("glUniform2ui64NV"))
	gpUniform2ui64vNV = (C.GPUNIFORM2UI64VNV)(getProcAddr("glUniform2ui64vNV"))
	gpUniform2uiv = (C.GPUNIFORM2UIV)(getProcAddr("glUniform2uiv"))
	if gpUniform2uiv == nil {
		return errors.New("glUniform2uiv")
	}
	gpUniform3f = (C.GPUNIFORM3F)(getProcAddr("glUniform3f"))
	if gpUniform3f == nil {
		return errors.New("glUniform3f")
	}
	gpUniform3fv = (C.GPUNIFORM3FV)(getProcAddr("glUniform3fv"))
	if gpUniform3fv == nil {
		return errors.New("glUniform3fv")
	}
	gpUniform3i = (C.GPUNIFORM3I)(getProcAddr("glUniform3i"))
	if gpUniform3i == nil {
		return errors.New("glUniform3i")
	}
	gpUniform3i64NV = (C.GPUNIFORM3I64NV)(getProcAddr("glUniform3i64NV"))
	gpUniform3i64vNV = (C.GPUNIFORM3I64VNV)(getProcAddr("glUniform3i64vNV"))
	gpUniform3iv = (C.GPUNIFORM3IV)(getProcAddr("glUniform3iv"))
	if gpUniform3iv == nil {
		return errors.New("glUniform3iv")
	}
	gpUniform3ui = (C.GPUNIFORM3UI)(getProcAddr("glUniform3ui"))
	if gpUniform3ui == nil {
		return errors.New("glUniform3ui")
	}
	gpUniform3ui64NV = (C.GPUNIFORM3UI64NV)(getProcAddr("glUniform3ui64NV"))
	gpUniform3ui64vNV = (C.GPUNIFORM3UI64VNV)(getProcAddr("glUniform3ui64vNV"))
	gpUniform3uiv = (C.GPUNIFORM3UIV)(getProcAddr("glUniform3uiv"))
	if gpUniform3uiv == nil {
		return errors.New("glUniform3uiv")
	}
	gpUniform4f = (C.GPUNIFORM4F)(getProcAddr("glUniform4f"))
	if gpUniform4f == nil {
		return errors.New("glUniform4f")
	}
	gpUniform4fv = (C.GPUNIFORM4FV)(getProcAddr("glUniform4fv"))
	if gpUniform4fv == nil {
		return errors.New("glUniform4fv")
	}
	gpUniform4i = (C.GPUNIFORM4I)(getProcAddr("glUniform4i"))
	if gpUniform4i == nil {
		return errors.New("glUniform4i")
	}
	gpUniform4i64NV = (C.GPUNIFORM4I64NV)(getProcAddr("glUniform4i64NV"))
	gpUniform4i64vNV = (C.GPUNIFORM4I64VNV)(getProcAddr("glUniform4i64vNV"))
	gpUniform4iv = (C.GPUNIFORM4IV)(getProcAddr("glUniform4iv"))
	if gpUniform4iv == nil {
		return errors.New("glUniform4iv")
	}
	gpUniform4ui = (C.GPUNIFORM4UI)(getProcAddr("glUniform4ui"))
	if gpUniform4ui == nil {
		return errors.New("glUniform4ui")
	}
	gpUniform4ui64NV = (C.GPUNIFORM4UI64NV)(getProcAddr("glUniform4ui64NV"))
	gpUniform4ui64vNV = (C.GPUNIFORM4UI64VNV)(getProcAddr("glUniform4ui64vNV"))
	gpUniform4uiv = (C.GPUNIFORM4UIV)(getProcAddr("glUniform4uiv"))
	if gpUniform4uiv == nil {
		return errors.New("glUniform4uiv")
	}
	gpUniformBlockBinding = (C.GPUNIFORMBLOCKBINDING)(getProcAddr("glUniformBlockBinding"))
	if gpUniformBlockBinding == nil {
		return errors.New("glUniformBlockBinding")
	}
	gpUniformHandleui64IMG = (C.GPUNIFORMHANDLEUI64IMG)(getProcAddr("glUniformHandleui64IMG"))
	gpUniformHandleui64NV = (C.GPUNIFORMHANDLEUI64NV)(getProcAddr("glUniformHandleui64NV"))
	gpUniformHandleui64vIMG = (C.GPUNIFORMHANDLEUI64VIMG)(getProcAddr("glUniformHandleui64vIMG"))
	gpUniformHandleui64vNV = (C.GPUNIFORMHANDLEUI64VNV)(getProcAddr("glUniformHandleui64vNV"))
	gpUniformMatrix2fv = (C.GPUNIFORMMATRIX2FV)(getProcAddr("glUniformMatrix2fv"))
	if gpUniformMatrix2fv == nil {
		return errors.New("glUniformMatrix2fv")
	}
	gpUniformMatrix2x3fv = (C.GPUNIFORMMATRIX2X3FV)(getProcAddr("glUniformMatrix2x3fv"))
	if gpUniformMatrix2x3fv == nil {
		return errors.New("glUniformMatrix2x3fv")
	}
	gpUniformMatrix2x3fvNV = (C.GPUNIFORMMATRIX2X3FVNV)(getProcAddr("glUniformMatrix2x3fvNV"))
	gpUniformMatrix2x4fv = (C.GPUNIFORMMATRIX2X4FV)(getProcAddr("glUniformMatrix2x4fv"))
	if gpUniformMatrix2x4fv == nil {
		return errors.New("glUniformMatrix2x4fv")
	}
	gpUniformMatrix2x4fvNV = (C.GPUNIFORMMATRIX2X4FVNV)(getProcAddr("glUniformMatrix2x4fvNV"))
	gpUniformMatrix3fv = (C.GPUNIFORMMATRIX3FV)(getProcAddr("glUniformMatrix3fv"))
	if gpUniformMatrix3fv == nil {
		return errors.New("glUniformMatrix3fv")
	}
	gpUniformMatrix3x2fv = (C.GPUNIFORMMATRIX3X2FV)(getProcAddr("glUniformMatrix3x2fv"))
	if gpUniformMatrix3x2fv == nil {
		return errors.New("glUniformMatrix3x2fv")
	}
	gpUniformMatrix3x2fvNV = (C.GPUNIFORMMATRIX3X2FVNV)(getProcAddr("glUniformMatrix3x2fvNV"))
	gpUniformMatrix3x4fv = (C.GPUNIFORMMATRIX3X4FV)(getProcAddr("glUniformMatrix3x4fv"))
	if gpUniformMatrix3x4fv == nil {
		return errors.New("glUniformMatrix3x4fv")
	}
	gpUniformMatrix3x4fvNV = (C.GPUNIFORMMATRIX3X4FVNV)(getProcAddr("glUniformMatrix3x4fvNV"))
	gpUniformMatrix4fv = (C.GPUNIFORMMATRIX4FV)(getProcAddr("glUniformMatrix4fv"))
	if gpUniformMatrix4fv == nil {
		return errors.New("glUniformMatrix4fv")
	}
	gpUniformMatrix4x2fv = (C.GPUNIFORMMATRIX4X2FV)(getProcAddr("glUniformMatrix4x2fv"))
	if gpUniformMatrix4x2fv == nil {
		return errors.New("glUniformMatrix4x2fv")
	}
	gpUniformMatrix4x2fvNV = (C.GPUNIFORMMATRIX4X2FVNV)(getProcAddr("glUniformMatrix4x2fvNV"))
	gpUniformMatrix4x3fv = (C.GPUNIFORMMATRIX4X3FV)(getProcAddr("glUniformMatrix4x3fv"))
	if gpUniformMatrix4x3fv == nil {
		return errors.New("glUniformMatrix4x3fv")
	}
	gpUniformMatrix4x3fvNV = (C.GPUNIFORMMATRIX4X3FVNV)(getProcAddr("glUniformMatrix4x3fvNV"))
	gpUnmapBuffer = (C.GPUNMAPBUFFER)(getProcAddr("glUnmapBuffer"))
	if gpUnmapBuffer == nil {
		return errors.New("glUnmapBuffer")
	}
	gpUnmapBufferOES = (C.GPUNMAPBUFFEROES)(getProcAddr("glUnmapBufferOES"))
	gpUseProgram = (C.GPUSEPROGRAM)(getProcAddr("glUseProgram"))
	if gpUseProgram == nil {
		return errors.New("glUseProgram")
	}
	gpUseProgramStages = (C.GPUSEPROGRAMSTAGES)(getProcAddr("glUseProgramStages"))
	if gpUseProgramStages == nil {
		return errors.New("glUseProgramStages")
	}
	gpUseProgramStagesEXT = (C.GPUSEPROGRAMSTAGESEXT)(getProcAddr("glUseProgramStagesEXT"))
	gpUseShaderProgramEXT = (C.GPUSESHADERPROGRAMEXT)(getProcAddr("glUseShaderProgramEXT"))
	gpValidateProgram = (C.GPVALIDATEPROGRAM)(getProcAddr("glValidateProgram"))
	if gpValidateProgram == nil {
		return errors.New("glValidateProgram")
	}
	gpValidateProgramPipeline = (C.GPVALIDATEPROGRAMPIPELINE)(getProcAddr("glValidateProgramPipeline"))
	if gpValidateProgramPipeline == nil {
		return errors.New("glValidateProgramPipeline")
	}
	gpValidateProgramPipelineEXT = (C.GPVALIDATEPROGRAMPIPELINEEXT)(getProcAddr("glValidateProgramPipelineEXT"))
	gpVertexAttrib1f = (C.GPVERTEXATTRIB1F)(getProcAddr("glVertexAttrib1f"))
	if gpVertexAttrib1f == nil {
		return errors.New("glVertexAttrib1f")
	}
	gpVertexAttrib1fv = (C.GPVERTEXATTRIB1FV)(getProcAddr("glVertexAttrib1fv"))
	if gpVertexAttrib1fv == nil {
		return errors.New("glVertexAttrib1fv")
	}
	gpVertexAttrib2f = (C.GPVERTEXATTRIB2F)(getProcAddr("glVertexAttrib2f"))
	if gpVertexAttrib2f == nil {
		return errors.New("glVertexAttrib2f")
	}
	gpVertexAttrib2fv = (C.GPVERTEXATTRIB2FV)(getProcAddr("glVertexAttrib2fv"))
	if gpVertexAttrib2fv == nil {
		return errors.New("glVertexAttrib2fv")
	}
	gpVertexAttrib3f = (C.GPVERTEXATTRIB3F)(getProcAddr("glVertexAttrib3f"))
	if gpVertexAttrib3f == nil {
		return errors.New("glVertexAttrib3f")
	}
	gpVertexAttrib3fv = (C.GPVERTEXATTRIB3FV)(getProcAddr("glVertexAttrib3fv"))
	if gpVertexAttrib3fv == nil {
		return errors.New("glVertexAttrib3fv")
	}
	gpVertexAttrib4f = (C.GPVERTEXATTRIB4F)(getProcAddr("glVertexAttrib4f"))
	if gpVertexAttrib4f == nil {
		return errors.New("glVertexAttrib4f")
	}
	gpVertexAttrib4fv = (C.GPVERTEXATTRIB4FV)(getProcAddr("glVertexAttrib4fv"))
	if gpVertexAttrib4fv == nil {
		return errors.New("glVertexAttrib4fv")
	}
	gpVertexAttribBinding = (C.GPVERTEXATTRIBBINDING)(getProcAddr("glVertexAttribBinding"))
	if gpVertexAttribBinding == nil {
		return errors.New("glVertexAttribBinding")
	}
	gpVertexAttribDivisor = (C.GPVERTEXATTRIBDIVISOR)(getProcAddr("glVertexAttribDivisor"))
	if gpVertexAttribDivisor == nil {
		return errors.New("glVertexAttribDivisor")
	}
	gpVertexAttribDivisorANGLE = (C.GPVERTEXATTRIBDIVISORANGLE)(getProcAddr("glVertexAttribDivisorANGLE"))
	gpVertexAttribDivisorEXT = (C.GPVERTEXATTRIBDIVISOREXT)(getProcAddr("glVertexAttribDivisorEXT"))
	gpVertexAttribDivisorNV = (C.GPVERTEXATTRIBDIVISORNV)(getProcAddr("glVertexAttribDivisorNV"))
	gpVertexAttribFormat = (C.GPVERTEXATTRIBFORMAT)(getProcAddr("glVertexAttribFormat"))
	if gpVertexAttribFormat == nil {
		return errors.New("glVertexAttribFormat")
	}
	gpVertexAttribI4i = (C.GPVERTEXATTRIBI4I)(getProcAddr("glVertexAttribI4i"))
	if gpVertexAttribI4i == nil {
		return errors.New("glVertexAttribI4i")
	}
	gpVertexAttribI4iv = (C.GPVERTEXATTRIBI4IV)(getProcAddr("glVertexAttribI4iv"))
	if gpVertexAttribI4iv == nil {
		return errors.New("glVertexAttribI4iv")
	}
	gpVertexAttribI4ui = (C.GPVERTEXATTRIBI4UI)(getProcAddr("glVertexAttribI4ui"))
	if gpVertexAttribI4ui == nil {
		return errors.New("glVertexAttribI4ui")
	}
	gpVertexAttribI4uiv = (C.GPVERTEXATTRIBI4UIV)(getProcAddr("glVertexAttribI4uiv"))
	if gpVertexAttribI4uiv == nil {
		return errors.New("glVertexAttribI4uiv")
	}
	gpVertexAttribIFormat = (C.GPVERTEXATTRIBIFORMAT)(getProcAddr("glVertexAttribIFormat"))
	if gpVertexAttribIFormat == nil {
		return errors.New("glVertexAttribIFormat")
	}
	gpVertexAttribIPointer = (C.GPVERTEXATTRIBIPOINTER)(getProcAddr("glVertexAttribIPointer"))
	if gpVertexAttribIPointer == nil {
		return errors.New("glVertexAttribIPointer")
	}
	gpVertexAttribPointer = (C.GPVERTEXATTRIBPOINTER)(getProcAddr("glVertexAttribPointer"))
	if gpVertexAttribPointer == nil {
		return errors.New("glVertexAttribPointer")
	}
	gpVertexBindingDivisor = (C.GPVERTEXBINDINGDIVISOR)(getProcAddr("glVertexBindingDivisor"))
	if gpVertexBindingDivisor == nil {
		return errors.New("glVertexBindingDivisor")
	}
	gpViewport = (C.GPVIEWPORT)(getProcAddr("glViewport"))
	if gpViewport == nil {
		return errors.New("glViewport")
	}
	gpViewportArrayvNV = (C.GPVIEWPORTARRAYVNV)(getProcAddr("glViewportArrayvNV"))
	gpViewportArrayvOES = (C.GPVIEWPORTARRAYVOES)(getProcAddr("glViewportArrayvOES"))
	gpViewportIndexedfNV = (C.GPVIEWPORTINDEXEDFNV)(getProcAddr("glViewportIndexedfNV"))
	gpViewportIndexedfOES = (C.GPVIEWPORTINDEXEDFOES)(getProcAddr("glViewportIndexedfOES"))
	gpViewportIndexedfvNV = (C.GPVIEWPORTINDEXEDFVNV)(getProcAddr("glViewportIndexedfvNV"))
	gpViewportIndexedfvOES = (C.GPVIEWPORTINDEXEDFVOES)(getProcAddr("glViewportIndexedfvOES"))
	gpViewportPositionWScaleNV = (C.GPVIEWPORTPOSITIONWSCALENV)(getProcAddr("glViewportPositionWScaleNV"))
	gpViewportSwizzleNV = (C.GPVIEWPORTSWIZZLENV)(getProcAddr("glViewportSwizzleNV"))
	gpWaitSemaphoreEXT = (C.GPWAITSEMAPHOREEXT)(getProcAddr("glWaitSemaphoreEXT"))
	gpWaitSync = (C.GPWAITSYNC)(getProcAddr("glWaitSync"))
	if gpWaitSync == nil {
		return errors.New("glWaitSync")
	}
	gpWaitSyncAPPLE = (C.GPWAITSYNCAPPLE)(getProcAddr("glWaitSyncAPPLE"))
	gpWaitVkSemaphoreNV = (C.GPWAITVKSEMAPHORENV)(getProcAddr("glWaitVkSemaphoreNV"))
	gpWeightPathsNV = (C.GPWEIGHTPATHSNV)(getProcAddr("glWeightPathsNV"))
	gpWindowRectanglesEXT = (C.GPWINDOWRECTANGLESEXT)(getProcAddr("glWindowRectanglesEXT"))
	return nil
}
