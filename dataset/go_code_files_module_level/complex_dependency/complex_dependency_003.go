func InitWithProcAddrFunc(getProcAddr func(name string) unsafe.Pointer) error {
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
	gpApplyFramebufferAttachmentCMAAINTEL = (C.GPAPPLYFRAMEBUFFERATTACHMENTCMAAINTEL)(getProcAddr("glApplyFramebufferAttachmentCMAAINTEL"))
	gpAttachShader = (C.GPATTACHSHADER)(getProcAddr("glAttachShader"))
	if gpAttachShader == nil {
		return errors.New("glAttachShader")
	}
	gpBeginConditionalRender = (C.GPBEGINCONDITIONALRENDER)(getProcAddr("glBeginConditionalRender"))
	if gpBeginConditionalRender == nil {
		return errors.New("glBeginConditionalRender")
	}
	gpBeginConditionalRenderNV = (C.GPBEGINCONDITIONALRENDERNV)(getProcAddr("glBeginConditionalRenderNV"))
	gpBeginPerfMonitorAMD = (C.GPBEGINPERFMONITORAMD)(getProcAddr("glBeginPerfMonitorAMD"))
	gpBeginPerfQueryINTEL = (C.GPBEGINPERFQUERYINTEL)(getProcAddr("glBeginPerfQueryINTEL"))
	gpBeginQuery = (C.GPBEGINQUERY)(getProcAddr("glBeginQuery"))
	if gpBeginQuery == nil {
		return errors.New("glBeginQuery")
	}
	gpBeginQueryIndexed = (C.GPBEGINQUERYINDEXED)(getProcAddr("glBeginQueryIndexed"))
	if gpBeginQueryIndexed == nil {
		return errors.New("glBeginQueryIndexed")
	}
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
	gpBindBuffersBase = (C.GPBINDBUFFERSBASE)(getProcAddr("glBindBuffersBase"))
	gpBindBuffersRange = (C.GPBINDBUFFERSRANGE)(getProcAddr("glBindBuffersRange"))
	gpBindFragDataLocation = (C.GPBINDFRAGDATALOCATION)(getProcAddr("glBindFragDataLocation"))
	if gpBindFragDataLocation == nil {
		return errors.New("glBindFragDataLocation")
	}
	gpBindFragDataLocationIndexed = (C.GPBINDFRAGDATALOCATIONINDEXED)(getProcAddr("glBindFragDataLocationIndexed"))
	if gpBindFragDataLocationIndexed == nil {
		return errors.New("glBindFragDataLocationIndexed")
	}
	gpBindFramebuffer = (C.GPBINDFRAMEBUFFER)(getProcAddr("glBindFramebuffer"))
	if gpBindFramebuffer == nil {
		return errors.New("glBindFramebuffer")
	}
	gpBindImageTexture = (C.GPBINDIMAGETEXTURE)(getProcAddr("glBindImageTexture"))
	gpBindImageTextures = (C.GPBINDIMAGETEXTURES)(getProcAddr("glBindImageTextures"))
	gpBindMultiTextureEXT = (C.GPBINDMULTITEXTUREEXT)(getProcAddr("glBindMultiTextureEXT"))
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
	gpBindSamplers = (C.GPBINDSAMPLERS)(getProcAddr("glBindSamplers"))
	gpBindTexture = (C.GPBINDTEXTURE)(getProcAddr("glBindTexture"))
	if gpBindTexture == nil {
		return errors.New("glBindTexture")
	}
	gpBindTextureUnit = (C.GPBINDTEXTUREUNIT)(getProcAddr("glBindTextureUnit"))
	gpBindTextures = (C.GPBINDTEXTURES)(getProcAddr("glBindTextures"))
	gpBindTransformFeedback = (C.GPBINDTRANSFORMFEEDBACK)(getProcAddr("glBindTransformFeedback"))
	if gpBindTransformFeedback == nil {
		return errors.New("glBindTransformFeedback")
	}
	gpBindVertexArray = (C.GPBINDVERTEXARRAY)(getProcAddr("glBindVertexArray"))
	if gpBindVertexArray == nil {
		return errors.New("glBindVertexArray")
	}
	gpBindVertexBuffer = (C.GPBINDVERTEXBUFFER)(getProcAddr("glBindVertexBuffer"))
	gpBindVertexBuffers = (C.GPBINDVERTEXBUFFERS)(getProcAddr("glBindVertexBuffers"))
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
	gpBlendEquationSeparate = (C.GPBLENDEQUATIONSEPARATE)(getProcAddr("glBlendEquationSeparate"))
	if gpBlendEquationSeparate == nil {
		return errors.New("glBlendEquationSeparate")
	}
	gpBlendEquationSeparatei = (C.GPBLENDEQUATIONSEPARATEI)(getProcAddr("glBlendEquationSeparatei"))
	if gpBlendEquationSeparatei == nil {
		return errors.New("glBlendEquationSeparatei")
	}
	gpBlendEquationSeparateiARB = (C.GPBLENDEQUATIONSEPARATEIARB)(getProcAddr("glBlendEquationSeparateiARB"))
	gpBlendEquationi = (C.GPBLENDEQUATIONI)(getProcAddr("glBlendEquationi"))
	if gpBlendEquationi == nil {
		return errors.New("glBlendEquationi")
	}
	gpBlendEquationiARB = (C.GPBLENDEQUATIONIARB)(getProcAddr("glBlendEquationiARB"))
	gpBlendFunc = (C.GPBLENDFUNC)(getProcAddr("glBlendFunc"))
	if gpBlendFunc == nil {
		return errors.New("glBlendFunc")
	}
	gpBlendFuncSeparate = (C.GPBLENDFUNCSEPARATE)(getProcAddr("glBlendFuncSeparate"))
	if gpBlendFuncSeparate == nil {
		return errors.New("glBlendFuncSeparate")
	}
	gpBlendFuncSeparatei = (C.GPBLENDFUNCSEPARATEI)(getProcAddr("glBlendFuncSeparatei"))
	if gpBlendFuncSeparatei == nil {
		return errors.New("glBlendFuncSeparatei")
	}
	gpBlendFuncSeparateiARB = (C.GPBLENDFUNCSEPARATEIARB)(getProcAddr("glBlendFuncSeparateiARB"))
	gpBlendFunci = (C.GPBLENDFUNCI)(getProcAddr("glBlendFunci"))
	if gpBlendFunci == nil {
		return errors.New("glBlendFunci")
	}
	gpBlendFunciARB = (C.GPBLENDFUNCIARB)(getProcAddr("glBlendFunciARB"))
	gpBlendParameteriNV = (C.GPBLENDPARAMETERINV)(getProcAddr("glBlendParameteriNV"))
	gpBlitFramebuffer = (C.GPBLITFRAMEBUFFER)(getProcAddr("glBlitFramebuffer"))
	if gpBlitFramebuffer == nil {
		return errors.New("glBlitFramebuffer")
	}
	gpBlitNamedFramebuffer = (C.GPBLITNAMEDFRAMEBUFFER)(getProcAddr("glBlitNamedFramebuffer"))
	gpBufferAddressRangeNV = (C.GPBUFFERADDRESSRANGENV)(getProcAddr("glBufferAddressRangeNV"))
	gpBufferData = (C.GPBUFFERDATA)(getProcAddr("glBufferData"))
	if gpBufferData == nil {
		return errors.New("glBufferData")
	}
	gpBufferPageCommitmentARB = (C.GPBUFFERPAGECOMMITMENTARB)(getProcAddr("glBufferPageCommitmentARB"))
	gpBufferStorage = (C.GPBUFFERSTORAGE)(getProcAddr("glBufferStorage"))
	gpBufferSubData = (C.GPBUFFERSUBDATA)(getProcAddr("glBufferSubData"))
	if gpBufferSubData == nil {
		return errors.New("glBufferSubData")
	}
	gpCallCommandListNV = (C.GPCALLCOMMANDLISTNV)(getProcAddr("glCallCommandListNV"))
	gpCheckFramebufferStatus = (C.GPCHECKFRAMEBUFFERSTATUS)(getProcAddr("glCheckFramebufferStatus"))
	if gpCheckFramebufferStatus == nil {
		return errors.New("glCheckFramebufferStatus")
	}
	gpCheckNamedFramebufferStatus = (C.GPCHECKNAMEDFRAMEBUFFERSTATUS)(getProcAddr("glCheckNamedFramebufferStatus"))
	gpCheckNamedFramebufferStatusEXT = (C.GPCHECKNAMEDFRAMEBUFFERSTATUSEXT)(getProcAddr("glCheckNamedFramebufferStatusEXT"))
	gpClampColor = (C.GPCLAMPCOLOR)(getProcAddr("glClampColor"))
	if gpClampColor == nil {
		return errors.New("glClampColor")
	}
	gpClear = (C.GPCLEAR)(getProcAddr("glClear"))
	if gpClear == nil {
		return errors.New("glClear")
	}
	gpClearBufferData = (C.GPCLEARBUFFERDATA)(getProcAddr("glClearBufferData"))
	gpClearBufferSubData = (C.GPCLEARBUFFERSUBDATA)(getProcAddr("glClearBufferSubData"))
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
	gpClearDepth = (C.GPCLEARDEPTH)(getProcAddr("glClearDepth"))
	if gpClearDepth == nil {
		return errors.New("glClearDepth")
	}
	gpClearDepthf = (C.GPCLEARDEPTHF)(getProcAddr("glClearDepthf"))
	if gpClearDepthf == nil {
		return errors.New("glClearDepthf")
	}
	gpClearNamedBufferData = (C.GPCLEARNAMEDBUFFERDATA)(getProcAddr("glClearNamedBufferData"))
	gpClearNamedBufferDataEXT = (C.GPCLEARNAMEDBUFFERDATAEXT)(getProcAddr("glClearNamedBufferDataEXT"))
	gpClearNamedBufferSubData = (C.GPCLEARNAMEDBUFFERSUBDATA)(getProcAddr("glClearNamedBufferSubData"))
	gpClearNamedBufferSubDataEXT = (C.GPCLEARNAMEDBUFFERSUBDATAEXT)(getProcAddr("glClearNamedBufferSubDataEXT"))
	gpClearNamedFramebufferfi = (C.GPCLEARNAMEDFRAMEBUFFERFI)(getProcAddr("glClearNamedFramebufferfi"))
	gpClearNamedFramebufferfv = (C.GPCLEARNAMEDFRAMEBUFFERFV)(getProcAddr("glClearNamedFramebufferfv"))
	gpClearNamedFramebufferiv = (C.GPCLEARNAMEDFRAMEBUFFERIV)(getProcAddr("glClearNamedFramebufferiv"))
	gpClearNamedFramebufferuiv = (C.GPCLEARNAMEDFRAMEBUFFERUIV)(getProcAddr("glClearNamedFramebufferuiv"))
	gpClearStencil = (C.GPCLEARSTENCIL)(getProcAddr("glClearStencil"))
	if gpClearStencil == nil {
		return errors.New("glClearStencil")
	}
	gpClearTexImage = (C.GPCLEARTEXIMAGE)(getProcAddr("glClearTexImage"))
	gpClearTexSubImage = (C.GPCLEARTEXSUBIMAGE)(getProcAddr("glClearTexSubImage"))
	gpClientAttribDefaultEXT = (C.GPCLIENTATTRIBDEFAULTEXT)(getProcAddr("glClientAttribDefaultEXT"))
	gpClientWaitSync = (C.GPCLIENTWAITSYNC)(getProcAddr("glClientWaitSync"))
	if gpClientWaitSync == nil {
		return errors.New("glClientWaitSync")
	}
	gpClipControl = (C.GPCLIPCONTROL)(getProcAddr("glClipControl"))
	gpColorFormatNV = (C.GPCOLORFORMATNV)(getProcAddr("glColorFormatNV"))
	gpColorMask = (C.GPCOLORMASK)(getProcAddr("glColorMask"))
	if gpColorMask == nil {
		return errors.New("glColorMask")
	}
	gpColorMaski = (C.GPCOLORMASKI)(getProcAddr("glColorMaski"))
	if gpColorMaski == nil {
		return errors.New("glColorMaski")
	}
	gpCommandListSegmentsNV = (C.GPCOMMANDLISTSEGMENTSNV)(getProcAddr("glCommandListSegmentsNV"))
	gpCompileCommandListNV = (C.GPCOMPILECOMMANDLISTNV)(getProcAddr("glCompileCommandListNV"))
	gpCompileShader = (C.GPCOMPILESHADER)(getProcAddr("glCompileShader"))
	if gpCompileShader == nil {
		return errors.New("glCompileShader")
	}
	gpCompileShaderIncludeARB = (C.GPCOMPILESHADERINCLUDEARB)(getProcAddr("glCompileShaderIncludeARB"))
	gpCompressedMultiTexImage1DEXT = (C.GPCOMPRESSEDMULTITEXIMAGE1DEXT)(getProcAddr("glCompressedMultiTexImage1DEXT"))
	gpCompressedMultiTexImage2DEXT = (C.GPCOMPRESSEDMULTITEXIMAGE2DEXT)(getProcAddr("glCompressedMultiTexImage2DEXT"))
	gpCompressedMultiTexImage3DEXT = (C.GPCOMPRESSEDMULTITEXIMAGE3DEXT)(getProcAddr("glCompressedMultiTexImage3DEXT"))
	gpCompressedMultiTexSubImage1DEXT = (C.GPCOMPRESSEDMULTITEXSUBIMAGE1DEXT)(getProcAddr("glCompressedMultiTexSubImage1DEXT"))
	gpCompressedMultiTexSubImage2DEXT = (C.GPCOMPRESSEDMULTITEXSUBIMAGE2DEXT)(getProcAddr("glCompressedMultiTexSubImage2DEXT"))
	gpCompressedMultiTexSubImage3DEXT = (C.GPCOMPRESSEDMULTITEXSUBIMAGE3DEXT)(getProcAddr("glCompressedMultiTexSubImage3DEXT"))
	gpCompressedTexImage1D = (C.GPCOMPRESSEDTEXIMAGE1D)(getProcAddr("glCompressedTexImage1D"))
	if gpCompressedTexImage1D == nil {
		return errors.New("glCompressedTexImage1D")
	}
	gpCompressedTexImage2D = (C.GPCOMPRESSEDTEXIMAGE2D)(getProcAddr("glCompressedTexImage2D"))
	if gpCompressedTexImage2D == nil {
		return errors.New("glCompressedTexImage2D")
	}
	gpCompressedTexImage3D = (C.GPCOMPRESSEDTEXIMAGE3D)(getProcAddr("glCompressedTexImage3D"))
	if gpCompressedTexImage3D == nil {
		return errors.New("glCompressedTexImage3D")
	}
	gpCompressedTexSubImage1D = (C.GPCOMPRESSEDTEXSUBIMAGE1D)(getProcAddr("glCompressedTexSubImage1D"))
	if gpCompressedTexSubImage1D == nil {
		return errors.New("glCompressedTexSubImage1D")
	}
	gpCompressedTexSubImage2D = (C.GPCOMPRESSEDTEXSUBIMAGE2D)(getProcAddr("glCompressedTexSubImage2D"))
	if gpCompressedTexSubImage2D == nil {
		return errors.New("glCompressedTexSubImage2D")
	}
	gpCompressedTexSubImage3D = (C.GPCOMPRESSEDTEXSUBIMAGE3D)(getProcAddr("glCompressedTexSubImage3D"))
	if gpCompressedTexSubImage3D == nil {
		return errors.New("glCompressedTexSubImage3D")
	}
	gpCompressedTextureImage1DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE1DEXT)(getProcAddr("glCompressedTextureImage1DEXT"))
	gpCompressedTextureImage2DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE2DEXT)(getProcAddr("glCompressedTextureImage2DEXT"))
	gpCompressedTextureImage3DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE3DEXT)(getProcAddr("glCompressedTextureImage3DEXT"))
	gpCompressedTextureSubImage1D = (C.GPCOMPRESSEDTEXTURESUBIMAGE1D)(getProcAddr("glCompressedTextureSubImage1D"))
	gpCompressedTextureSubImage1DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE1DEXT)(getProcAddr("glCompressedTextureSubImage1DEXT"))
	gpCompressedTextureSubImage2D = (C.GPCOMPRESSEDTEXTURESUBIMAGE2D)(getProcAddr("glCompressedTextureSubImage2D"))
	gpCompressedTextureSubImage2DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE2DEXT)(getProcAddr("glCompressedTextureSubImage2DEXT"))
	gpCompressedTextureSubImage3D = (C.GPCOMPRESSEDTEXTURESUBIMAGE3D)(getProcAddr("glCompressedTextureSubImage3D"))
	gpCompressedTextureSubImage3DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE3DEXT)(getProcAddr("glCompressedTextureSubImage3DEXT"))
	gpConservativeRasterParameterfNV = (C.GPCONSERVATIVERASTERPARAMETERFNV)(getProcAddr("glConservativeRasterParameterfNV"))
	gpConservativeRasterParameteriNV = (C.GPCONSERVATIVERASTERPARAMETERINV)(getProcAddr("glConservativeRasterParameteriNV"))
	gpCopyBufferSubData = (C.GPCOPYBUFFERSUBDATA)(getProcAddr("glCopyBufferSubData"))
	if gpCopyBufferSubData == nil {
		return errors.New("glCopyBufferSubData")
	}
	gpCopyImageSubData = (C.GPCOPYIMAGESUBDATA)(getProcAddr("glCopyImageSubData"))
	gpCopyMultiTexImage1DEXT = (C.GPCOPYMULTITEXIMAGE1DEXT)(getProcAddr("glCopyMultiTexImage1DEXT"))
	gpCopyMultiTexImage2DEXT = (C.GPCOPYMULTITEXIMAGE2DEXT)(getProcAddr("glCopyMultiTexImage2DEXT"))
	gpCopyMultiTexSubImage1DEXT = (C.GPCOPYMULTITEXSUBIMAGE1DEXT)(getProcAddr("glCopyMultiTexSubImage1DEXT"))
	gpCopyMultiTexSubImage2DEXT = (C.GPCOPYMULTITEXSUBIMAGE2DEXT)(getProcAddr("glCopyMultiTexSubImage2DEXT"))
	gpCopyMultiTexSubImage3DEXT = (C.GPCOPYMULTITEXSUBIMAGE3DEXT)(getProcAddr("glCopyMultiTexSubImage3DEXT"))
	gpCopyNamedBufferSubData = (C.GPCOPYNAMEDBUFFERSUBDATA)(getProcAddr("glCopyNamedBufferSubData"))
	gpCopyPathNV = (C.GPCOPYPATHNV)(getProcAddr("glCopyPathNV"))
	gpCopyTexImage1D = (C.GPCOPYTEXIMAGE1D)(getProcAddr("glCopyTexImage1D"))
	if gpCopyTexImage1D == nil {
		return errors.New("glCopyTexImage1D")
	}
	gpCopyTexImage2D = (C.GPCOPYTEXIMAGE2D)(getProcAddr("glCopyTexImage2D"))
	if gpCopyTexImage2D == nil {
		return errors.New("glCopyTexImage2D")
	}
	gpCopyTexSubImage1D = (C.GPCOPYTEXSUBIMAGE1D)(getProcAddr("glCopyTexSubImage1D"))
	if gpCopyTexSubImage1D == nil {
		return errors.New("glCopyTexSubImage1D")
	}
	gpCopyTexSubImage2D = (C.GPCOPYTEXSUBIMAGE2D)(getProcAddr("glCopyTexSubImage2D"))
	if gpCopyTexSubImage2D == nil {
		return errors.New("glCopyTexSubImage2D")
	}
	gpCopyTexSubImage3D = (C.GPCOPYTEXSUBIMAGE3D)(getProcAddr("glCopyTexSubImage3D"))
	if gpCopyTexSubImage3D == nil {
		return errors.New("glCopyTexSubImage3D")
	}
	gpCopyTextureImage1DEXT = (C.GPCOPYTEXTUREIMAGE1DEXT)(getProcAddr("glCopyTextureImage1DEXT"))
	gpCopyTextureImage2DEXT = (C.GPCOPYTEXTUREIMAGE2DEXT)(getProcAddr("glCopyTextureImage2DEXT"))
	gpCopyTextureSubImage1D = (C.GPCOPYTEXTURESUBIMAGE1D)(getProcAddr("glCopyTextureSubImage1D"))
	gpCopyTextureSubImage1DEXT = (C.GPCOPYTEXTURESUBIMAGE1DEXT)(getProcAddr("glCopyTextureSubImage1DEXT"))
	gpCopyTextureSubImage2D = (C.GPCOPYTEXTURESUBIMAGE2D)(getProcAddr("glCopyTextureSubImage2D"))
	gpCopyTextureSubImage2DEXT = (C.GPCOPYTEXTURESUBIMAGE2DEXT)(getProcAddr("glCopyTextureSubImage2DEXT"))
	gpCopyTextureSubImage3D = (C.GPCOPYTEXTURESUBIMAGE3D)(getProcAddr("glCopyTextureSubImage3D"))
	gpCopyTextureSubImage3DEXT = (C.GPCOPYTEXTURESUBIMAGE3DEXT)(getProcAddr("glCopyTextureSubImage3DEXT"))
	gpCoverFillPathInstancedNV = (C.GPCOVERFILLPATHINSTANCEDNV)(getProcAddr("glCoverFillPathInstancedNV"))
	gpCoverFillPathNV = (C.GPCOVERFILLPATHNV)(getProcAddr("glCoverFillPathNV"))
	gpCoverStrokePathInstancedNV = (C.GPCOVERSTROKEPATHINSTANCEDNV)(getProcAddr("glCoverStrokePathInstancedNV"))
	gpCoverStrokePathNV = (C.GPCOVERSTROKEPATHNV)(getProcAddr("glCoverStrokePathNV"))
	gpCoverageModulationNV = (C.GPCOVERAGEMODULATIONNV)(getProcAddr("glCoverageModulationNV"))
	gpCoverageModulationTableNV = (C.GPCOVERAGEMODULATIONTABLENV)(getProcAddr("glCoverageModulationTableNV"))
	gpCreateBuffers = (C.GPCREATEBUFFERS)(getProcAddr("glCreateBuffers"))
	gpCreateCommandListsNV = (C.GPCREATECOMMANDLISTSNV)(getProcAddr("glCreateCommandListsNV"))
	gpCreateFramebuffers = (C.GPCREATEFRAMEBUFFERS)(getProcAddr("glCreateFramebuffers"))
	gpCreatePerfQueryINTEL = (C.GPCREATEPERFQUERYINTEL)(getProcAddr("glCreatePerfQueryINTEL"))
	gpCreateProgram = (C.GPCREATEPROGRAM)(getProcAddr("glCreateProgram"))
	if gpCreateProgram == nil {
		return errors.New("glCreateProgram")
	}
	gpCreateProgramPipelines = (C.GPCREATEPROGRAMPIPELINES)(getProcAddr("glCreateProgramPipelines"))
	gpCreateQueries = (C.GPCREATEQUERIES)(getProcAddr("glCreateQueries"))
	gpCreateRenderbuffers = (C.GPCREATERENDERBUFFERS)(getProcAddr("glCreateRenderbuffers"))
	gpCreateSamplers = (C.GPCREATESAMPLERS)(getProcAddr("glCreateSamplers"))
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
	gpCreateStatesNV = (C.GPCREATESTATESNV)(getProcAddr("glCreateStatesNV"))
	gpCreateSyncFromCLeventARB = (C.GPCREATESYNCFROMCLEVENTARB)(getProcAddr("glCreateSyncFromCLeventARB"))
	gpCreateTextures = (C.GPCREATETEXTURES)(getProcAddr("glCreateTextures"))
	gpCreateTransformFeedbacks = (C.GPCREATETRANSFORMFEEDBACKS)(getProcAddr("glCreateTransformFeedbacks"))
	gpCreateVertexArrays = (C.GPCREATEVERTEXARRAYS)(getProcAddr("glCreateVertexArrays"))
	gpCullFace = (C.GPCULLFACE)(getProcAddr("glCullFace"))
	if gpCullFace == nil {
		return errors.New("glCullFace")
	}
	gpDebugMessageCallback = (C.GPDEBUGMESSAGECALLBACK)(getProcAddr("glDebugMessageCallback"))
	gpDebugMessageCallbackARB = (C.GPDEBUGMESSAGECALLBACKARB)(getProcAddr("glDebugMessageCallbackARB"))
	gpDebugMessageCallbackKHR = (C.GPDEBUGMESSAGECALLBACKKHR)(getProcAddr("glDebugMessageCallbackKHR"))
	gpDebugMessageControl = (C.GPDEBUGMESSAGECONTROL)(getProcAddr("glDebugMessageControl"))
	gpDebugMessageControlARB = (C.GPDEBUGMESSAGECONTROLARB)(getProcAddr("glDebugMessageControlARB"))
	gpDebugMessageControlKHR = (C.GPDEBUGMESSAGECONTROLKHR)(getProcAddr("glDebugMessageControlKHR"))
	gpDebugMessageInsert = (C.GPDEBUGMESSAGEINSERT)(getProcAddr("glDebugMessageInsert"))
	gpDebugMessageInsertARB = (C.GPDEBUGMESSAGEINSERTARB)(getProcAddr("glDebugMessageInsertARB"))
	gpDebugMessageInsertKHR = (C.GPDEBUGMESSAGEINSERTKHR)(getProcAddr("glDebugMessageInsertKHR"))
	gpDeleteBuffers = (C.GPDELETEBUFFERS)(getProcAddr("glDeleteBuffers"))
	if gpDeleteBuffers == nil {
		return errors.New("glDeleteBuffers")
	}
	gpDeleteCommandListsNV = (C.GPDELETECOMMANDLISTSNV)(getProcAddr("glDeleteCommandListsNV"))
	gpDeleteFramebuffers = (C.GPDELETEFRAMEBUFFERS)(getProcAddr("glDeleteFramebuffers"))
	if gpDeleteFramebuffers == nil {
		return errors.New("glDeleteFramebuffers")
	}
	gpDeleteNamedStringARB = (C.GPDELETENAMEDSTRINGARB)(getProcAddr("glDeleteNamedStringARB"))
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
	gpDeleteRenderbuffers = (C.GPDELETERENDERBUFFERS)(getProcAddr("glDeleteRenderbuffers"))
	if gpDeleteRenderbuffers == nil {
		return errors.New("glDeleteRenderbuffers")
	}
	gpDeleteSamplers = (C.GPDELETESAMPLERS)(getProcAddr("glDeleteSamplers"))
	if gpDeleteSamplers == nil {
		return errors.New("glDeleteSamplers")
	}
	gpDeleteShader = (C.GPDELETESHADER)(getProcAddr("glDeleteShader"))
	if gpDeleteShader == nil {
		return errors.New("glDeleteShader")
	}
	gpDeleteStatesNV = (C.GPDELETESTATESNV)(getProcAddr("glDeleteStatesNV"))
	gpDeleteSync = (C.GPDELETESYNC)(getProcAddr("glDeleteSync"))
	if gpDeleteSync == nil {
		return errors.New("glDeleteSync")
	}
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
	gpDepthFunc = (C.GPDEPTHFUNC)(getProcAddr("glDepthFunc"))
	if gpDepthFunc == nil {
		return errors.New("glDepthFunc")
	}
	gpDepthMask = (C.GPDEPTHMASK)(getProcAddr("glDepthMask"))
	if gpDepthMask == nil {
		return errors.New("glDepthMask")
	}
	gpDepthRange = (C.GPDEPTHRANGE)(getProcAddr("glDepthRange"))
	if gpDepthRange == nil {
		return errors.New("glDepthRange")
	}
	gpDepthRangeArrayv = (C.GPDEPTHRANGEARRAYV)(getProcAddr("glDepthRangeArrayv"))
	if gpDepthRangeArrayv == nil {
		return errors.New("glDepthRangeArrayv")
	}
	gpDepthRangeIndexed = (C.GPDEPTHRANGEINDEXED)(getProcAddr("glDepthRangeIndexed"))
	if gpDepthRangeIndexed == nil {
		return errors.New("glDepthRangeIndexed")
	}
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
	gpDisableClientStateIndexedEXT = (C.GPDISABLECLIENTSTATEINDEXEDEXT)(getProcAddr("glDisableClientStateIndexedEXT"))
	gpDisableClientStateiEXT = (C.GPDISABLECLIENTSTATEIEXT)(getProcAddr("glDisableClientStateiEXT"))
	gpDisableIndexedEXT = (C.GPDISABLEINDEXEDEXT)(getProcAddr("glDisableIndexedEXT"))
	gpDisableVertexArrayAttrib = (C.GPDISABLEVERTEXARRAYATTRIB)(getProcAddr("glDisableVertexArrayAttrib"))
	gpDisableVertexArrayAttribEXT = (C.GPDISABLEVERTEXARRAYATTRIBEXT)(getProcAddr("glDisableVertexArrayAttribEXT"))
	gpDisableVertexArrayEXT = (C.GPDISABLEVERTEXARRAYEXT)(getProcAddr("glDisableVertexArrayEXT"))
	gpDisableVertexAttribArray = (C.GPDISABLEVERTEXATTRIBARRAY)(getProcAddr("glDisableVertexAttribArray"))
	if gpDisableVertexAttribArray == nil {
		return errors.New("glDisableVertexAttribArray")
	}
	gpDisablei = (C.GPDISABLEI)(getProcAddr("glDisablei"))
	if gpDisablei == nil {
		return errors.New("glDisablei")
	}
	gpDispatchCompute = (C.GPDISPATCHCOMPUTE)(getProcAddr("glDispatchCompute"))
	gpDispatchComputeGroupSizeARB = (C.GPDISPATCHCOMPUTEGROUPSIZEARB)(getProcAddr("glDispatchComputeGroupSizeARB"))
	gpDispatchComputeIndirect = (C.GPDISPATCHCOMPUTEINDIRECT)(getProcAddr("glDispatchComputeIndirect"))
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
	gpDrawArraysInstancedARB = (C.GPDRAWARRAYSINSTANCEDARB)(getProcAddr("glDrawArraysInstancedARB"))
	gpDrawArraysInstancedBaseInstance = (C.GPDRAWARRAYSINSTANCEDBASEINSTANCE)(getProcAddr("glDrawArraysInstancedBaseInstance"))
	gpDrawArraysInstancedEXT = (C.GPDRAWARRAYSINSTANCEDEXT)(getProcAddr("glDrawArraysInstancedEXT"))
	gpDrawBuffer = (C.GPDRAWBUFFER)(getProcAddr("glDrawBuffer"))
	if gpDrawBuffer == nil {
		return errors.New("glDrawBuffer")
	}
	gpDrawBuffers = (C.GPDRAWBUFFERS)(getProcAddr("glDrawBuffers"))
	if gpDrawBuffers == nil {
		return errors.New("glDrawBuffers")
	}
	gpDrawCommandsAddressNV = (C.GPDRAWCOMMANDSADDRESSNV)(getProcAddr("glDrawCommandsAddressNV"))
	gpDrawCommandsNV = (C.GPDRAWCOMMANDSNV)(getProcAddr("glDrawCommandsNV"))
	gpDrawCommandsStatesAddressNV = (C.GPDRAWCOMMANDSSTATESADDRESSNV)(getProcAddr("glDrawCommandsStatesAddressNV"))
	gpDrawCommandsStatesNV = (C.GPDRAWCOMMANDSSTATESNV)(getProcAddr("glDrawCommandsStatesNV"))
	gpDrawElements = (C.GPDRAWELEMENTS)(getProcAddr("glDrawElements"))
	if gpDrawElements == nil {
		return errors.New("glDrawElements")
	}
	gpDrawElementsBaseVertex = (C.GPDRAWELEMENTSBASEVERTEX)(getProcAddr("glDrawElementsBaseVertex"))
	if gpDrawElementsBaseVertex == nil {
		return errors.New("glDrawElementsBaseVertex")
	}
	gpDrawElementsIndirect = (C.GPDRAWELEMENTSINDIRECT)(getProcAddr("glDrawElementsIndirect"))
	if gpDrawElementsIndirect == nil {
		return errors.New("glDrawElementsIndirect")
	}
	gpDrawElementsInstanced = (C.GPDRAWELEMENTSINSTANCED)(getProcAddr("glDrawElementsInstanced"))
	if gpDrawElementsInstanced == nil {
		return errors.New("glDrawElementsInstanced")
	}
	gpDrawElementsInstancedARB = (C.GPDRAWELEMENTSINSTANCEDARB)(getProcAddr("glDrawElementsInstancedARB"))
	gpDrawElementsInstancedBaseInstance = (C.GPDRAWELEMENTSINSTANCEDBASEINSTANCE)(getProcAddr("glDrawElementsInstancedBaseInstance"))
	gpDrawElementsInstancedBaseVertex = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEX)(getProcAddr("glDrawElementsInstancedBaseVertex"))
	if gpDrawElementsInstancedBaseVertex == nil {
		return errors.New("glDrawElementsInstancedBaseVertex")
	}
	gpDrawElementsInstancedBaseVertexBaseInstance = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEXBASEINSTANCE)(getProcAddr("glDrawElementsInstancedBaseVertexBaseInstance"))
	gpDrawElementsInstancedEXT = (C.GPDRAWELEMENTSINSTANCEDEXT)(getProcAddr("glDrawElementsInstancedEXT"))
	gpDrawRangeElements = (C.GPDRAWRANGEELEMENTS)(getProcAddr("glDrawRangeElements"))
	if gpDrawRangeElements == nil {
		return errors.New("glDrawRangeElements")
	}
	gpDrawRangeElementsBaseVertex = (C.GPDRAWRANGEELEMENTSBASEVERTEX)(getProcAddr("glDrawRangeElementsBaseVertex"))
	if gpDrawRangeElementsBaseVertex == nil {
		return errors.New("glDrawRangeElementsBaseVertex")
	}
	gpDrawTransformFeedback = (C.GPDRAWTRANSFORMFEEDBACK)(getProcAddr("glDrawTransformFeedback"))
	if gpDrawTransformFeedback == nil {
		return errors.New("glDrawTransformFeedback")
	}
	gpDrawTransformFeedbackInstanced = (C.GPDRAWTRANSFORMFEEDBACKINSTANCED)(getProcAddr("glDrawTransformFeedbackInstanced"))
	gpDrawTransformFeedbackStream = (C.GPDRAWTRANSFORMFEEDBACKSTREAM)(getProcAddr("glDrawTransformFeedbackStream"))
	if gpDrawTransformFeedbackStream == nil {
		return errors.New("glDrawTransformFeedbackStream")
	}
	gpDrawTransformFeedbackStreamInstanced = (C.GPDRAWTRANSFORMFEEDBACKSTREAMINSTANCED)(getProcAddr("glDrawTransformFeedbackStreamInstanced"))
	gpDrawVkImageNV = (C.GPDRAWVKIMAGENV)(getProcAddr("glDrawVkImageNV"))
	gpEGLImageTargetTexStorageEXT = (C.GPEGLIMAGETARGETTEXSTORAGEEXT)(getProcAddr("glEGLImageTargetTexStorageEXT"))
	gpEGLImageTargetTextureStorageEXT = (C.GPEGLIMAGETARGETTEXTURESTORAGEEXT)(getProcAddr("glEGLImageTargetTextureStorageEXT"))
	gpEdgeFlagFormatNV = (C.GPEDGEFLAGFORMATNV)(getProcAddr("glEdgeFlagFormatNV"))
	gpEnable = (C.GPENABLE)(getProcAddr("glEnable"))
	if gpEnable == nil {
		return errors.New("glEnable")
	}
	gpEnableClientStateIndexedEXT = (C.GPENABLECLIENTSTATEINDEXEDEXT)(getProcAddr("glEnableClientStateIndexedEXT"))
	gpEnableClientStateiEXT = (C.GPENABLECLIENTSTATEIEXT)(getProcAddr("glEnableClientStateiEXT"))
	gpEnableIndexedEXT = (C.GPENABLEINDEXEDEXT)(getProcAddr("glEnableIndexedEXT"))
	gpEnableVertexArrayAttrib = (C.GPENABLEVERTEXARRAYATTRIB)(getProcAddr("glEnableVertexArrayAttrib"))
	gpEnableVertexArrayAttribEXT = (C.GPENABLEVERTEXARRAYATTRIBEXT)(getProcAddr("glEnableVertexArrayAttribEXT"))
	gpEnableVertexArrayEXT = (C.GPENABLEVERTEXARRAYEXT)(getProcAddr("glEnableVertexArrayEXT"))
	gpEnableVertexAttribArray = (C.GPENABLEVERTEXATTRIBARRAY)(getProcAddr("glEnableVertexAttribArray"))
	if gpEnableVertexAttribArray == nil {
		return errors.New("glEnableVertexAttribArray")
	}
	gpEnablei = (C.GPENABLEI)(getProcAddr("glEnablei"))
	if gpEnablei == nil {
		return errors.New("glEnablei")
	}
	gpEndConditionalRender = (C.GPENDCONDITIONALRENDER)(getProcAddr("glEndConditionalRender"))
	if gpEndConditionalRender == nil {
		return errors.New("glEndConditionalRender")
	}
	gpEndConditionalRenderNV = (C.GPENDCONDITIONALRENDERNV)(getProcAddr("glEndConditionalRenderNV"))
	gpEndPerfMonitorAMD = (C.GPENDPERFMONITORAMD)(getProcAddr("glEndPerfMonitorAMD"))
	gpEndPerfQueryINTEL = (C.GPENDPERFQUERYINTEL)(getProcAddr("glEndPerfQueryINTEL"))
	gpEndQuery = (C.GPENDQUERY)(getProcAddr("glEndQuery"))
	if gpEndQuery == nil {
		return errors.New("glEndQuery")
	}
	gpEndQueryIndexed = (C.GPENDQUERYINDEXED)(getProcAddr("glEndQueryIndexed"))
	if gpEndQueryIndexed == nil {
		return errors.New("glEndQueryIndexed")
	}
	gpEndTransformFeedback = (C.GPENDTRANSFORMFEEDBACK)(getProcAddr("glEndTransformFeedback"))
	if gpEndTransformFeedback == nil {
		return errors.New("glEndTransformFeedback")
	}
	gpEvaluateDepthValuesARB = (C.GPEVALUATEDEPTHVALUESARB)(getProcAddr("glEvaluateDepthValuesARB"))
	gpFenceSync = (C.GPFENCESYNC)(getProcAddr("glFenceSync"))
	if gpFenceSync == nil {
		return errors.New("glFenceSync")
	}
	gpFinish = (C.GPFINISH)(getProcAddr("glFinish"))
	if gpFinish == nil {
		return errors.New("glFinish")
	}
	gpFlush = (C.GPFLUSH)(getProcAddr("glFlush"))
	if gpFlush == nil {
		return errors.New("glFlush")
	}
	gpFlushMappedBufferRange = (C.GPFLUSHMAPPEDBUFFERRANGE)(getProcAddr("glFlushMappedBufferRange"))
	if gpFlushMappedBufferRange == nil {
		return errors.New("glFlushMappedBufferRange")
	}
	gpFlushMappedNamedBufferRange = (C.GPFLUSHMAPPEDNAMEDBUFFERRANGE)(getProcAddr("glFlushMappedNamedBufferRange"))
	gpFlushMappedNamedBufferRangeEXT = (C.GPFLUSHMAPPEDNAMEDBUFFERRANGEEXT)(getProcAddr("glFlushMappedNamedBufferRangeEXT"))
	gpFogCoordFormatNV = (C.GPFOGCOORDFORMATNV)(getProcAddr("glFogCoordFormatNV"))
	gpFragmentCoverageColorNV = (C.GPFRAGMENTCOVERAGECOLORNV)(getProcAddr("glFragmentCoverageColorNV"))
	gpFramebufferDrawBufferEXT = (C.GPFRAMEBUFFERDRAWBUFFEREXT)(getProcAddr("glFramebufferDrawBufferEXT"))
	gpFramebufferDrawBuffersEXT = (C.GPFRAMEBUFFERDRAWBUFFERSEXT)(getProcAddr("glFramebufferDrawBuffersEXT"))
	gpFramebufferFetchBarrierEXT = (C.GPFRAMEBUFFERFETCHBARRIEREXT)(getProcAddr("glFramebufferFetchBarrierEXT"))
	gpFramebufferParameteri = (C.GPFRAMEBUFFERPARAMETERI)(getProcAddr("glFramebufferParameteri"))
	gpFramebufferReadBufferEXT = (C.GPFRAMEBUFFERREADBUFFEREXT)(getProcAddr("glFramebufferReadBufferEXT"))
	gpFramebufferRenderbuffer = (C.GPFRAMEBUFFERRENDERBUFFER)(getProcAddr("glFramebufferRenderbuffer"))
	if gpFramebufferRenderbuffer == nil {
		return errors.New("glFramebufferRenderbuffer")
	}
	gpFramebufferSampleLocationsfvARB = (C.GPFRAMEBUFFERSAMPLELOCATIONSFVARB)(getProcAddr("glFramebufferSampleLocationsfvARB"))
	gpFramebufferSampleLocationsfvNV = (C.GPFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glFramebufferSampleLocationsfvNV"))
	gpFramebufferTexture = (C.GPFRAMEBUFFERTEXTURE)(getProcAddr("glFramebufferTexture"))
	if gpFramebufferTexture == nil {
		return errors.New("glFramebufferTexture")
	}
	gpFramebufferTexture1D = (C.GPFRAMEBUFFERTEXTURE1D)(getProcAddr("glFramebufferTexture1D"))
	if gpFramebufferTexture1D == nil {
		return errors.New("glFramebufferTexture1D")
	}
	gpFramebufferTexture2D = (C.GPFRAMEBUFFERTEXTURE2D)(getProcAddr("glFramebufferTexture2D"))
	if gpFramebufferTexture2D == nil {
		return errors.New("glFramebufferTexture2D")
	}
	gpFramebufferTexture3D = (C.GPFRAMEBUFFERTEXTURE3D)(getProcAddr("glFramebufferTexture3D"))
	if gpFramebufferTexture3D == nil {
		return errors.New("glFramebufferTexture3D")
	}
	gpFramebufferTextureARB = (C.GPFRAMEBUFFERTEXTUREARB)(getProcAddr("glFramebufferTextureARB"))
	gpFramebufferTextureFaceARB = (C.GPFRAMEBUFFERTEXTUREFACEARB)(getProcAddr("glFramebufferTextureFaceARB"))
	gpFramebufferTextureLayer = (C.GPFRAMEBUFFERTEXTURELAYER)(getProcAddr("glFramebufferTextureLayer"))
	if gpFramebufferTextureLayer == nil {
		return errors.New("glFramebufferTextureLayer")
	}
	gpFramebufferTextureLayerARB = (C.GPFRAMEBUFFERTEXTURELAYERARB)(getProcAddr("glFramebufferTextureLayerARB"))
	gpFramebufferTextureMultiviewOVR = (C.GPFRAMEBUFFERTEXTUREMULTIVIEWOVR)(getProcAddr("glFramebufferTextureMultiviewOVR"))
	gpFrontFace = (C.GPFRONTFACE)(getProcAddr("glFrontFace"))
	if gpFrontFace == nil {
		return errors.New("glFrontFace")
	}
	gpGenBuffers = (C.GPGENBUFFERS)(getProcAddr("glGenBuffers"))
	if gpGenBuffers == nil {
		return errors.New("glGenBuffers")
	}
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
	gpGenRenderbuffers = (C.GPGENRENDERBUFFERS)(getProcAddr("glGenRenderbuffers"))
	if gpGenRenderbuffers == nil {
		return errors.New("glGenRenderbuffers")
	}
	gpGenSamplers = (C.GPGENSAMPLERS)(getProcAddr("glGenSamplers"))
	if gpGenSamplers == nil {
		return errors.New("glGenSamplers")
	}
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
	gpGenerateMipmap = (C.GPGENERATEMIPMAP)(getProcAddr("glGenerateMipmap"))
	if gpGenerateMipmap == nil {
		return errors.New("glGenerateMipmap")
	}
	gpGenerateMultiTexMipmapEXT = (C.GPGENERATEMULTITEXMIPMAPEXT)(getProcAddr("glGenerateMultiTexMipmapEXT"))
	gpGenerateTextureMipmap = (C.GPGENERATETEXTUREMIPMAP)(getProcAddr("glGenerateTextureMipmap"))
	gpGenerateTextureMipmapEXT = (C.GPGENERATETEXTUREMIPMAPEXT)(getProcAddr("glGenerateTextureMipmapEXT"))
	gpGetActiveAtomicCounterBufferiv = (C.GPGETACTIVEATOMICCOUNTERBUFFERIV)(getProcAddr("glGetActiveAtomicCounterBufferiv"))
	gpGetActiveAttrib = (C.GPGETACTIVEATTRIB)(getProcAddr("glGetActiveAttrib"))
	if gpGetActiveAttrib == nil {
		return errors.New("glGetActiveAttrib")
	}
	gpGetActiveSubroutineName = (C.GPGETACTIVESUBROUTINENAME)(getProcAddr("glGetActiveSubroutineName"))
	if gpGetActiveSubroutineName == nil {
		return errors.New("glGetActiveSubroutineName")
	}
	gpGetActiveSubroutineUniformName = (C.GPGETACTIVESUBROUTINEUNIFORMNAME)(getProcAddr("glGetActiveSubroutineUniformName"))
	if gpGetActiveSubroutineUniformName == nil {
		return errors.New("glGetActiveSubroutineUniformName")
	}
	gpGetActiveSubroutineUniformiv = (C.GPGETACTIVESUBROUTINEUNIFORMIV)(getProcAddr("glGetActiveSubroutineUniformiv"))
	if gpGetActiveSubroutineUniformiv == nil {
		return errors.New("glGetActiveSubroutineUniformiv")
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
	gpGetActiveUniformName = (C.GPGETACTIVEUNIFORMNAME)(getProcAddr("glGetActiveUniformName"))
	if gpGetActiveUniformName == nil {
		return errors.New("glGetActiveUniformName")
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
	gpGetBooleanIndexedvEXT = (C.GPGETBOOLEANINDEXEDVEXT)(getProcAddr("glGetBooleanIndexedvEXT"))
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
	gpGetBufferParameterui64vNV = (C.GPGETBUFFERPARAMETERUI64VNV)(getProcAddr("glGetBufferParameterui64vNV"))
	gpGetBufferPointerv = (C.GPGETBUFFERPOINTERV)(getProcAddr("glGetBufferPointerv"))
	if gpGetBufferPointerv == nil {
		return errors.New("glGetBufferPointerv")
	}
	gpGetBufferSubData = (C.GPGETBUFFERSUBDATA)(getProcAddr("glGetBufferSubData"))
	if gpGetBufferSubData == nil {
		return errors.New("glGetBufferSubData")
	}
	gpGetCommandHeaderNV = (C.GPGETCOMMANDHEADERNV)(getProcAddr("glGetCommandHeaderNV"))
	gpGetCompressedMultiTexImageEXT = (C.GPGETCOMPRESSEDMULTITEXIMAGEEXT)(getProcAddr("glGetCompressedMultiTexImageEXT"))
	gpGetCompressedTexImage = (C.GPGETCOMPRESSEDTEXIMAGE)(getProcAddr("glGetCompressedTexImage"))
	if gpGetCompressedTexImage == nil {
		return errors.New("glGetCompressedTexImage")
	}
	gpGetCompressedTextureImage = (C.GPGETCOMPRESSEDTEXTUREIMAGE)(getProcAddr("glGetCompressedTextureImage"))
	gpGetCompressedTextureImageEXT = (C.GPGETCOMPRESSEDTEXTUREIMAGEEXT)(getProcAddr("glGetCompressedTextureImageEXT"))
	gpGetCompressedTextureSubImage = (C.GPGETCOMPRESSEDTEXTURESUBIMAGE)(getProcAddr("glGetCompressedTextureSubImage"))
	gpGetCoverageModulationTableNV = (C.GPGETCOVERAGEMODULATIONTABLENV)(getProcAddr("glGetCoverageModulationTableNV"))
	gpGetDebugMessageLog = (C.GPGETDEBUGMESSAGELOG)(getProcAddr("glGetDebugMessageLog"))
	gpGetDebugMessageLogARB = (C.GPGETDEBUGMESSAGELOGARB)(getProcAddr("glGetDebugMessageLogARB"))
	gpGetDebugMessageLogKHR = (C.GPGETDEBUGMESSAGELOGKHR)(getProcAddr("glGetDebugMessageLogKHR"))
	gpGetDoubleIndexedvEXT = (C.GPGETDOUBLEINDEXEDVEXT)(getProcAddr("glGetDoubleIndexedvEXT"))
	gpGetDoublei_v = (C.GPGETDOUBLEI_V)(getProcAddr("glGetDoublei_v"))
	if gpGetDoublei_v == nil {
		return errors.New("glGetDoublei_v")
	}
	gpGetDoublei_vEXT = (C.GPGETDOUBLEI_VEXT)(getProcAddr("glGetDoublei_vEXT"))
	gpGetDoublev = (C.GPGETDOUBLEV)(getProcAddr("glGetDoublev"))
	if gpGetDoublev == nil {
		return errors.New("glGetDoublev")
	}
	gpGetError = (C.GPGETERROR)(getProcAddr("glGetError"))
	if gpGetError == nil {
		return errors.New("glGetError")
	}
	gpGetFirstPerfQueryIdINTEL = (C.GPGETFIRSTPERFQUERYIDINTEL)(getProcAddr("glGetFirstPerfQueryIdINTEL"))
	gpGetFloatIndexedvEXT = (C.GPGETFLOATINDEXEDVEXT)(getProcAddr("glGetFloatIndexedvEXT"))
	gpGetFloati_v = (C.GPGETFLOATI_V)(getProcAddr("glGetFloati_v"))
	if gpGetFloati_v == nil {
		return errors.New("glGetFloati_v")
	}
	gpGetFloati_vEXT = (C.GPGETFLOATI_VEXT)(getProcAddr("glGetFloati_vEXT"))
	gpGetFloatv = (C.GPGETFLOATV)(getProcAddr("glGetFloatv"))
	if gpGetFloatv == nil {
		return errors.New("glGetFloatv")
	}
	gpGetFragDataIndex = (C.GPGETFRAGDATAINDEX)(getProcAddr("glGetFragDataIndex"))
	if gpGetFragDataIndex == nil {
		return errors.New("glGetFragDataIndex")
	}
	gpGetFragDataLocation = (C.GPGETFRAGDATALOCATION)(getProcAddr("glGetFragDataLocation"))
	if gpGetFragDataLocation == nil {
		return errors.New("glGetFragDataLocation")
	}
	gpGetFramebufferAttachmentParameteriv = (C.GPGETFRAMEBUFFERATTACHMENTPARAMETERIV)(getProcAddr("glGetFramebufferAttachmentParameteriv"))
	if gpGetFramebufferAttachmentParameteriv == nil {
		return errors.New("glGetFramebufferAttachmentParameteriv")
	}
	gpGetFramebufferParameteriv = (C.GPGETFRAMEBUFFERPARAMETERIV)(getProcAddr("glGetFramebufferParameteriv"))
	gpGetFramebufferParameterivEXT = (C.GPGETFRAMEBUFFERPARAMETERIVEXT)(getProcAddr("glGetFramebufferParameterivEXT"))
	gpGetGraphicsResetStatus = (C.GPGETGRAPHICSRESETSTATUS)(getProcAddr("glGetGraphicsResetStatus"))
	gpGetGraphicsResetStatusARB = (C.GPGETGRAPHICSRESETSTATUSARB)(getProcAddr("glGetGraphicsResetStatusARB"))
	gpGetGraphicsResetStatusKHR = (C.GPGETGRAPHICSRESETSTATUSKHR)(getProcAddr("glGetGraphicsResetStatusKHR"))
	gpGetImageHandleARB = (C.GPGETIMAGEHANDLEARB)(getProcAddr("glGetImageHandleARB"))
	gpGetImageHandleNV = (C.GPGETIMAGEHANDLENV)(getProcAddr("glGetImageHandleNV"))
	gpGetInteger64i_v = (C.GPGETINTEGER64I_V)(getProcAddr("glGetInteger64i_v"))
	if gpGetInteger64i_v == nil {
		return errors.New("glGetInteger64i_v")
	}
	gpGetInteger64v = (C.GPGETINTEGER64V)(getProcAddr("glGetInteger64v"))
	if gpGetInteger64v == nil {
		return errors.New("glGetInteger64v")
	}
	gpGetIntegerIndexedvEXT = (C.GPGETINTEGERINDEXEDVEXT)(getProcAddr("glGetIntegerIndexedvEXT"))
	gpGetIntegeri_v = (C.GPGETINTEGERI_V)(getProcAddr("glGetIntegeri_v"))
	if gpGetIntegeri_v == nil {
		return errors.New("glGetIntegeri_v")
	}
	gpGetIntegerui64i_vNV = (C.GPGETINTEGERUI64I_VNV)(getProcAddr("glGetIntegerui64i_vNV"))
	gpGetIntegerui64vNV = (C.GPGETINTEGERUI64VNV)(getProcAddr("glGetIntegerui64vNV"))
	gpGetIntegerv = (C.GPGETINTEGERV)(getProcAddr("glGetIntegerv"))
	if gpGetIntegerv == nil {
		return errors.New("glGetIntegerv")
	}
	gpGetInternalformatSampleivNV = (C.GPGETINTERNALFORMATSAMPLEIVNV)(getProcAddr("glGetInternalformatSampleivNV"))
	gpGetInternalformati64v = (C.GPGETINTERNALFORMATI64V)(getProcAddr("glGetInternalformati64v"))
	gpGetInternalformativ = (C.GPGETINTERNALFORMATIV)(getProcAddr("glGetInternalformativ"))
	gpGetMultiTexEnvfvEXT = (C.GPGETMULTITEXENVFVEXT)(getProcAddr("glGetMultiTexEnvfvEXT"))
	gpGetMultiTexEnvivEXT = (C.GPGETMULTITEXENVIVEXT)(getProcAddr("glGetMultiTexEnvivEXT"))
	gpGetMultiTexGendvEXT = (C.GPGETMULTITEXGENDVEXT)(getProcAddr("glGetMultiTexGendvEXT"))
	gpGetMultiTexGenfvEXT = (C.GPGETMULTITEXGENFVEXT)(getProcAddr("glGetMultiTexGenfvEXT"))
	gpGetMultiTexGenivEXT = (C.GPGETMULTITEXGENIVEXT)(getProcAddr("glGetMultiTexGenivEXT"))
	gpGetMultiTexImageEXT = (C.GPGETMULTITEXIMAGEEXT)(getProcAddr("glGetMultiTexImageEXT"))
	gpGetMultiTexLevelParameterfvEXT = (C.GPGETMULTITEXLEVELPARAMETERFVEXT)(getProcAddr("glGetMultiTexLevelParameterfvEXT"))
	gpGetMultiTexLevelParameterivEXT = (C.GPGETMULTITEXLEVELPARAMETERIVEXT)(getProcAddr("glGetMultiTexLevelParameterivEXT"))
	gpGetMultiTexParameterIivEXT = (C.GPGETMULTITEXPARAMETERIIVEXT)(getProcAddr("glGetMultiTexParameterIivEXT"))
	gpGetMultiTexParameterIuivEXT = (C.GPGETMULTITEXPARAMETERIUIVEXT)(getProcAddr("glGetMultiTexParameterIuivEXT"))
	gpGetMultiTexParameterfvEXT = (C.GPGETMULTITEXPARAMETERFVEXT)(getProcAddr("glGetMultiTexParameterfvEXT"))
	gpGetMultiTexParameterivEXT = (C.GPGETMULTITEXPARAMETERIVEXT)(getProcAddr("glGetMultiTexParameterivEXT"))
	gpGetMultisamplefv = (C.GPGETMULTISAMPLEFV)(getProcAddr("glGetMultisamplefv"))
	if gpGetMultisamplefv == nil {
		return errors.New("glGetMultisamplefv")
	}
	gpGetNamedBufferParameteri64v = (C.GPGETNAMEDBUFFERPARAMETERI64V)(getProcAddr("glGetNamedBufferParameteri64v"))
	gpGetNamedBufferParameteriv = (C.GPGETNAMEDBUFFERPARAMETERIV)(getProcAddr("glGetNamedBufferParameteriv"))
	gpGetNamedBufferParameterivEXT = (C.GPGETNAMEDBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedBufferParameterivEXT"))
	gpGetNamedBufferParameterui64vNV = (C.GPGETNAMEDBUFFERPARAMETERUI64VNV)(getProcAddr("glGetNamedBufferParameterui64vNV"))
	gpGetNamedBufferPointerv = (C.GPGETNAMEDBUFFERPOINTERV)(getProcAddr("glGetNamedBufferPointerv"))
	gpGetNamedBufferPointervEXT = (C.GPGETNAMEDBUFFERPOINTERVEXT)(getProcAddr("glGetNamedBufferPointervEXT"))
	gpGetNamedBufferSubData = (C.GPGETNAMEDBUFFERSUBDATA)(getProcAddr("glGetNamedBufferSubData"))
	gpGetNamedBufferSubDataEXT = (C.GPGETNAMEDBUFFERSUBDATAEXT)(getProcAddr("glGetNamedBufferSubDataEXT"))
	gpGetNamedFramebufferAttachmentParameteriv = (C.GPGETNAMEDFRAMEBUFFERATTACHMENTPARAMETERIV)(getProcAddr("glGetNamedFramebufferAttachmentParameteriv"))
	gpGetNamedFramebufferAttachmentParameterivEXT = (C.GPGETNAMEDFRAMEBUFFERATTACHMENTPARAMETERIVEXT)(getProcAddr("glGetNamedFramebufferAttachmentParameterivEXT"))
	gpGetNamedFramebufferParameteriv = (C.GPGETNAMEDFRAMEBUFFERPARAMETERIV)(getProcAddr("glGetNamedFramebufferParameteriv"))
	gpGetNamedFramebufferParameterivEXT = (C.GPGETNAMEDFRAMEBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedFramebufferParameterivEXT"))
	gpGetNamedProgramLocalParameterIivEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERIIVEXT)(getProcAddr("glGetNamedProgramLocalParameterIivEXT"))
	gpGetNamedProgramLocalParameterIuivEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERIUIVEXT)(getProcAddr("glGetNamedProgramLocalParameterIuivEXT"))
	gpGetNamedProgramLocalParameterdvEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERDVEXT)(getProcAddr("glGetNamedProgramLocalParameterdvEXT"))
	gpGetNamedProgramLocalParameterfvEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERFVEXT)(getProcAddr("glGetNamedProgramLocalParameterfvEXT"))
	gpGetNamedProgramStringEXT = (C.GPGETNAMEDPROGRAMSTRINGEXT)(getProcAddr("glGetNamedProgramStringEXT"))
	gpGetNamedProgramivEXT = (C.GPGETNAMEDPROGRAMIVEXT)(getProcAddr("glGetNamedProgramivEXT"))
	gpGetNamedRenderbufferParameteriv = (C.GPGETNAMEDRENDERBUFFERPARAMETERIV)(getProcAddr("glGetNamedRenderbufferParameteriv"))
	gpGetNamedRenderbufferParameterivEXT = (C.GPGETNAMEDRENDERBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedRenderbufferParameterivEXT"))
	gpGetNamedStringARB = (C.GPGETNAMEDSTRINGARB)(getProcAddr("glGetNamedStringARB"))
	gpGetNamedStringivARB = (C.GPGETNAMEDSTRINGIVARB)(getProcAddr("glGetNamedStringivARB"))
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
	gpGetPointerIndexedvEXT = (C.GPGETPOINTERINDEXEDVEXT)(getProcAddr("glGetPointerIndexedvEXT"))
	gpGetPointeri_vEXT = (C.GPGETPOINTERI_VEXT)(getProcAddr("glGetPointeri_vEXT"))
	gpGetPointerv = (C.GPGETPOINTERV)(getProcAddr("glGetPointerv"))
	gpGetPointervKHR = (C.GPGETPOINTERVKHR)(getProcAddr("glGetPointervKHR"))
	gpGetProgramBinary = (C.GPGETPROGRAMBINARY)(getProcAddr("glGetProgramBinary"))
	if gpGetProgramBinary == nil {
		return errors.New("glGetProgramBinary")
	}
	gpGetProgramInfoLog = (C.GPGETPROGRAMINFOLOG)(getProcAddr("glGetProgramInfoLog"))
	if gpGetProgramInfoLog == nil {
		return errors.New("glGetProgramInfoLog")
	}
	gpGetProgramInterfaceiv = (C.GPGETPROGRAMINTERFACEIV)(getProcAddr("glGetProgramInterfaceiv"))
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
	gpGetProgramResourceLocation = (C.GPGETPROGRAMRESOURCELOCATION)(getProcAddr("glGetProgramResourceLocation"))
	gpGetProgramResourceLocationIndex = (C.GPGETPROGRAMRESOURCELOCATIONINDEX)(getProcAddr("glGetProgramResourceLocationIndex"))
	gpGetProgramResourceName = (C.GPGETPROGRAMRESOURCENAME)(getProcAddr("glGetProgramResourceName"))
	gpGetProgramResourcefvNV = (C.GPGETPROGRAMRESOURCEFVNV)(getProcAddr("glGetProgramResourcefvNV"))
	gpGetProgramResourceiv = (C.GPGETPROGRAMRESOURCEIV)(getProcAddr("glGetProgramResourceiv"))
	gpGetProgramStageiv = (C.GPGETPROGRAMSTAGEIV)(getProcAddr("glGetProgramStageiv"))
	if gpGetProgramStageiv == nil {
		return errors.New("glGetProgramStageiv")
	}
	gpGetProgramiv = (C.GPGETPROGRAMIV)(getProcAddr("glGetProgramiv"))
	if gpGetProgramiv == nil {
		return errors.New("glGetProgramiv")
	}
	gpGetQueryBufferObjecti64v = (C.GPGETQUERYBUFFEROBJECTI64V)(getProcAddr("glGetQueryBufferObjecti64v"))
	gpGetQueryBufferObjectiv = (C.GPGETQUERYBUFFEROBJECTIV)(getProcAddr("glGetQueryBufferObjectiv"))
	gpGetQueryBufferObjectui64v = (C.GPGETQUERYBUFFEROBJECTUI64V)(getProcAddr("glGetQueryBufferObjectui64v"))
	gpGetQueryBufferObjectuiv = (C.GPGETQUERYBUFFEROBJECTUIV)(getProcAddr("glGetQueryBufferObjectuiv"))
	gpGetQueryIndexediv = (C.GPGETQUERYINDEXEDIV)(getProcAddr("glGetQueryIndexediv"))
	if gpGetQueryIndexediv == nil {
		return errors.New("glGetQueryIndexediv")
	}
	gpGetQueryObjecti64v = (C.GPGETQUERYOBJECTI64V)(getProcAddr("glGetQueryObjecti64v"))
	if gpGetQueryObjecti64v == nil {
		return errors.New("glGetQueryObjecti64v")
	}
	gpGetQueryObjectiv = (C.GPGETQUERYOBJECTIV)(getProcAddr("glGetQueryObjectiv"))
	if gpGetQueryObjectiv == nil {
		return errors.New("glGetQueryObjectiv")
	}
	gpGetQueryObjectui64v = (C.GPGETQUERYOBJECTUI64V)(getProcAddr("glGetQueryObjectui64v"))
	if gpGetQueryObjectui64v == nil {
		return errors.New("glGetQueryObjectui64v")
	}
	gpGetQueryObjectuiv = (C.GPGETQUERYOBJECTUIV)(getProcAddr("glGetQueryObjectuiv"))
	if gpGetQueryObjectuiv == nil {
		return errors.New("glGetQueryObjectuiv")
	}
	gpGetQueryiv = (C.GPGETQUERYIV)(getProcAddr("glGetQueryiv"))
	if gpGetQueryiv == nil {
		return errors.New("glGetQueryiv")
	}
	gpGetRenderbufferParameteriv = (C.GPGETRENDERBUFFERPARAMETERIV)(getProcAddr("glGetRenderbufferParameteriv"))
	if gpGetRenderbufferParameteriv == nil {
		return errors.New("glGetRenderbufferParameteriv")
	}
	gpGetSamplerParameterIiv = (C.GPGETSAMPLERPARAMETERIIV)(getProcAddr("glGetSamplerParameterIiv"))
	if gpGetSamplerParameterIiv == nil {
		return errors.New("glGetSamplerParameterIiv")
	}
	gpGetSamplerParameterIuiv = (C.GPGETSAMPLERPARAMETERIUIV)(getProcAddr("glGetSamplerParameterIuiv"))
	if gpGetSamplerParameterIuiv == nil {
		return errors.New("glGetSamplerParameterIuiv")
	}
	gpGetSamplerParameterfv = (C.GPGETSAMPLERPARAMETERFV)(getProcAddr("glGetSamplerParameterfv"))
	if gpGetSamplerParameterfv == nil {
		return errors.New("glGetSamplerParameterfv")
	}
	gpGetSamplerParameteriv = (C.GPGETSAMPLERPARAMETERIV)(getProcAddr("glGetSamplerParameteriv"))
	if gpGetSamplerParameteriv == nil {
		return errors.New("glGetSamplerParameteriv")
	}
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
	gpGetStageIndexNV = (C.GPGETSTAGEINDEXNV)(getProcAddr("glGetStageIndexNV"))
	gpGetString = (C.GPGETSTRING)(getProcAddr("glGetString"))
	if gpGetString == nil {
		return errors.New("glGetString")
	}
	gpGetStringi = (C.GPGETSTRINGI)(getProcAddr("glGetStringi"))
	if gpGetStringi == nil {
		return errors.New("glGetStringi")
	}
	gpGetSubroutineIndex = (C.GPGETSUBROUTINEINDEX)(getProcAddr("glGetSubroutineIndex"))
	if gpGetSubroutineIndex == nil {
		return errors.New("glGetSubroutineIndex")
	}
	gpGetSubroutineUniformLocation = (C.GPGETSUBROUTINEUNIFORMLOCATION)(getProcAddr("glGetSubroutineUniformLocation"))
	if gpGetSubroutineUniformLocation == nil {
		return errors.New("glGetSubroutineUniformLocation")
	}
	gpGetSynciv = (C.GPGETSYNCIV)(getProcAddr("glGetSynciv"))
	if gpGetSynciv == nil {
		return errors.New("glGetSynciv")
	}
	gpGetTexImage = (C.GPGETTEXIMAGE)(getProcAddr("glGetTexImage"))
	if gpGetTexImage == nil {
		return errors.New("glGetTexImage")
	}
	gpGetTexLevelParameterfv = (C.GPGETTEXLEVELPARAMETERFV)(getProcAddr("glGetTexLevelParameterfv"))
	if gpGetTexLevelParameterfv == nil {
		return errors.New("glGetTexLevelParameterfv")
	}
	gpGetTexLevelParameteriv = (C.GPGETTEXLEVELPARAMETERIV)(getProcAddr("glGetTexLevelParameteriv"))
	if gpGetTexLevelParameteriv == nil {
		return errors.New("glGetTexLevelParameteriv")
	}
	gpGetTexParameterIiv = (C.GPGETTEXPARAMETERIIV)(getProcAddr("glGetTexParameterIiv"))
	if gpGetTexParameterIiv == nil {
		return errors.New("glGetTexParameterIiv")
	}
	gpGetTexParameterIuiv = (C.GPGETTEXPARAMETERIUIV)(getProcAddr("glGetTexParameterIuiv"))
	if gpGetTexParameterIuiv == nil {
		return errors.New("glGetTexParameterIuiv")
	}
	gpGetTexParameterfv = (C.GPGETTEXPARAMETERFV)(getProcAddr("glGetTexParameterfv"))
	if gpGetTexParameterfv == nil {
		return errors.New("glGetTexParameterfv")
	}
	gpGetTexParameteriv = (C.GPGETTEXPARAMETERIV)(getProcAddr("glGetTexParameteriv"))
	if gpGetTexParameteriv == nil {
		return errors.New("glGetTexParameteriv")
	}
	gpGetTextureHandleARB = (C.GPGETTEXTUREHANDLEARB)(getProcAddr("glGetTextureHandleARB"))
	gpGetTextureHandleNV = (C.GPGETTEXTUREHANDLENV)(getProcAddr("glGetTextureHandleNV"))
	gpGetTextureImage = (C.GPGETTEXTUREIMAGE)(getProcAddr("glGetTextureImage"))
	gpGetTextureImageEXT = (C.GPGETTEXTUREIMAGEEXT)(getProcAddr("glGetTextureImageEXT"))
	gpGetTextureLevelParameterfv = (C.GPGETTEXTURELEVELPARAMETERFV)(getProcAddr("glGetTextureLevelParameterfv"))
	gpGetTextureLevelParameterfvEXT = (C.GPGETTEXTURELEVELPARAMETERFVEXT)(getProcAddr("glGetTextureLevelParameterfvEXT"))
	gpGetTextureLevelParameteriv = (C.GPGETTEXTURELEVELPARAMETERIV)(getProcAddr("glGetTextureLevelParameteriv"))
	gpGetTextureLevelParameterivEXT = (C.GPGETTEXTURELEVELPARAMETERIVEXT)(getProcAddr("glGetTextureLevelParameterivEXT"))
	gpGetTextureParameterIiv = (C.GPGETTEXTUREPARAMETERIIV)(getProcAddr("glGetTextureParameterIiv"))
	gpGetTextureParameterIivEXT = (C.GPGETTEXTUREPARAMETERIIVEXT)(getProcAddr("glGetTextureParameterIivEXT"))
	gpGetTextureParameterIuiv = (C.GPGETTEXTUREPARAMETERIUIV)(getProcAddr("glGetTextureParameterIuiv"))
	gpGetTextureParameterIuivEXT = (C.GPGETTEXTUREPARAMETERIUIVEXT)(getProcAddr("glGetTextureParameterIuivEXT"))
	gpGetTextureParameterfv = (C.GPGETTEXTUREPARAMETERFV)(getProcAddr("glGetTextureParameterfv"))
	gpGetTextureParameterfvEXT = (C.GPGETTEXTUREPARAMETERFVEXT)(getProcAddr("glGetTextureParameterfvEXT"))
	gpGetTextureParameteriv = (C.GPGETTEXTUREPARAMETERIV)(getProcAddr("glGetTextureParameteriv"))
	gpGetTextureParameterivEXT = (C.GPGETTEXTUREPARAMETERIVEXT)(getProcAddr("glGetTextureParameterivEXT"))
	gpGetTextureSamplerHandleARB = (C.GPGETTEXTURESAMPLERHANDLEARB)(getProcAddr("glGetTextureSamplerHandleARB"))
	gpGetTextureSamplerHandleNV = (C.GPGETTEXTURESAMPLERHANDLENV)(getProcAddr("glGetTextureSamplerHandleNV"))
	gpGetTextureSubImage = (C.GPGETTEXTURESUBIMAGE)(getProcAddr("glGetTextureSubImage"))
	gpGetTransformFeedbackVarying = (C.GPGETTRANSFORMFEEDBACKVARYING)(getProcAddr("glGetTransformFeedbackVarying"))
	if gpGetTransformFeedbackVarying == nil {
		return errors.New("glGetTransformFeedbackVarying")
	}
	gpGetTransformFeedbacki64_v = (C.GPGETTRANSFORMFEEDBACKI64_V)(getProcAddr("glGetTransformFeedbacki64_v"))
	gpGetTransformFeedbacki_v = (C.GPGETTRANSFORMFEEDBACKI_V)(getProcAddr("glGetTransformFeedbacki_v"))
	gpGetTransformFeedbackiv = (C.GPGETTRANSFORMFEEDBACKIV)(getProcAddr("glGetTransformFeedbackiv"))
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
	gpGetUniformSubroutineuiv = (C.GPGETUNIFORMSUBROUTINEUIV)(getProcAddr("glGetUniformSubroutineuiv"))
	if gpGetUniformSubroutineuiv == nil {
		return errors.New("glGetUniformSubroutineuiv")
	}
	gpGetUniformdv = (C.GPGETUNIFORMDV)(getProcAddr("glGetUniformdv"))
	if gpGetUniformdv == nil {
		return errors.New("glGetUniformdv")
	}
	gpGetUniformfv = (C.GPGETUNIFORMFV)(getProcAddr("glGetUniformfv"))
	if gpGetUniformfv == nil {
		return errors.New("glGetUniformfv")
	}
	gpGetUniformi64vARB = (C.GPGETUNIFORMI64VARB)(getProcAddr("glGetUniformi64vARB"))
	gpGetUniformi64vNV = (C.GPGETUNIFORMI64VNV)(getProcAddr("glGetUniformi64vNV"))
	gpGetUniformiv = (C.GPGETUNIFORMIV)(getProcAddr("glGetUniformiv"))
	if gpGetUniformiv == nil {
		return errors.New("glGetUniformiv")
	}
	gpGetUniformui64vARB = (C.GPGETUNIFORMUI64VARB)(getProcAddr("glGetUniformui64vARB"))
	gpGetUniformui64vNV = (C.GPGETUNIFORMUI64VNV)(getProcAddr("glGetUniformui64vNV"))
	gpGetUniformuiv = (C.GPGETUNIFORMUIV)(getProcAddr("glGetUniformuiv"))
	if gpGetUniformuiv == nil {
		return errors.New("glGetUniformuiv")
	}
	gpGetVertexArrayIndexed64iv = (C.GPGETVERTEXARRAYINDEXED64IV)(getProcAddr("glGetVertexArrayIndexed64iv"))
	gpGetVertexArrayIndexediv = (C.GPGETVERTEXARRAYINDEXEDIV)(getProcAddr("glGetVertexArrayIndexediv"))
	gpGetVertexArrayIntegeri_vEXT = (C.GPGETVERTEXARRAYINTEGERI_VEXT)(getProcAddr("glGetVertexArrayIntegeri_vEXT"))
	gpGetVertexArrayIntegervEXT = (C.GPGETVERTEXARRAYINTEGERVEXT)(getProcAddr("glGetVertexArrayIntegervEXT"))
	gpGetVertexArrayPointeri_vEXT = (C.GPGETVERTEXARRAYPOINTERI_VEXT)(getProcAddr("glGetVertexArrayPointeri_vEXT"))
	gpGetVertexArrayPointervEXT = (C.GPGETVERTEXARRAYPOINTERVEXT)(getProcAddr("glGetVertexArrayPointervEXT"))
	gpGetVertexArrayiv = (C.GPGETVERTEXARRAYIV)(getProcAddr("glGetVertexArrayiv"))
	gpGetVertexAttribIiv = (C.GPGETVERTEXATTRIBIIV)(getProcAddr("glGetVertexAttribIiv"))
	if gpGetVertexAttribIiv == nil {
		return errors.New("glGetVertexAttribIiv")
	}
	gpGetVertexAttribIuiv = (C.GPGETVERTEXATTRIBIUIV)(getProcAddr("glGetVertexAttribIuiv"))
	if gpGetVertexAttribIuiv == nil {
		return errors.New("glGetVertexAttribIuiv")
	}
	gpGetVertexAttribLdv = (C.GPGETVERTEXATTRIBLDV)(getProcAddr("glGetVertexAttribLdv"))
	if gpGetVertexAttribLdv == nil {
		return errors.New("glGetVertexAttribLdv")
	}
	gpGetVertexAttribLi64vNV = (C.GPGETVERTEXATTRIBLI64VNV)(getProcAddr("glGetVertexAttribLi64vNV"))
	gpGetVertexAttribLui64vARB = (C.GPGETVERTEXATTRIBLUI64VARB)(getProcAddr("glGetVertexAttribLui64vARB"))
	gpGetVertexAttribLui64vNV = (C.GPGETVERTEXATTRIBLUI64VNV)(getProcAddr("glGetVertexAttribLui64vNV"))
	gpGetVertexAttribPointerv = (C.GPGETVERTEXATTRIBPOINTERV)(getProcAddr("glGetVertexAttribPointerv"))
	if gpGetVertexAttribPointerv == nil {
		return errors.New("glGetVertexAttribPointerv")
	}
	gpGetVertexAttribdv = (C.GPGETVERTEXATTRIBDV)(getProcAddr("glGetVertexAttribdv"))
	if gpGetVertexAttribdv == nil {
		return errors.New("glGetVertexAttribdv")
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
	gpGetnCompressedTexImageARB = (C.GPGETNCOMPRESSEDTEXIMAGEARB)(getProcAddr("glGetnCompressedTexImageARB"))
	gpGetnTexImageARB = (C.GPGETNTEXIMAGEARB)(getProcAddr("glGetnTexImageARB"))
	gpGetnUniformdvARB = (C.GPGETNUNIFORMDVARB)(getProcAddr("glGetnUniformdvARB"))
	gpGetnUniformfv = (C.GPGETNUNIFORMFV)(getProcAddr("glGetnUniformfv"))
	gpGetnUniformfvARB = (C.GPGETNUNIFORMFVARB)(getProcAddr("glGetnUniformfvARB"))
	gpGetnUniformfvKHR = (C.GPGETNUNIFORMFVKHR)(getProcAddr("glGetnUniformfvKHR"))
	gpGetnUniformi64vARB = (C.GPGETNUNIFORMI64VARB)(getProcAddr("glGetnUniformi64vARB"))
	gpGetnUniformiv = (C.GPGETNUNIFORMIV)(getProcAddr("glGetnUniformiv"))
	gpGetnUniformivARB = (C.GPGETNUNIFORMIVARB)(getProcAddr("glGetnUniformivARB"))
	gpGetnUniformivKHR = (C.GPGETNUNIFORMIVKHR)(getProcAddr("glGetnUniformivKHR"))
	gpGetnUniformui64vARB = (C.GPGETNUNIFORMUI64VARB)(getProcAddr("glGetnUniformui64vARB"))
	gpGetnUniformuiv = (C.GPGETNUNIFORMUIV)(getProcAddr("glGetnUniformuiv"))
	gpGetnUniformuivARB = (C.GPGETNUNIFORMUIVARB)(getProcAddr("glGetnUniformuivARB"))
	gpGetnUniformuivKHR = (C.GPGETNUNIFORMUIVKHR)(getProcAddr("glGetnUniformuivKHR"))
	gpHint = (C.GPHINT)(getProcAddr("glHint"))
	if gpHint == nil {
		return errors.New("glHint")
	}
	gpIndexFormatNV = (C.GPINDEXFORMATNV)(getProcAddr("glIndexFormatNV"))
	gpInsertEventMarkerEXT = (C.GPINSERTEVENTMARKEREXT)(getProcAddr("glInsertEventMarkerEXT"))
	gpInterpolatePathsNV = (C.GPINTERPOLATEPATHSNV)(getProcAddr("glInterpolatePathsNV"))
	gpInvalidateBufferData = (C.GPINVALIDATEBUFFERDATA)(getProcAddr("glInvalidateBufferData"))
	gpInvalidateBufferSubData = (C.GPINVALIDATEBUFFERSUBDATA)(getProcAddr("glInvalidateBufferSubData"))
	gpInvalidateFramebuffer = (C.GPINVALIDATEFRAMEBUFFER)(getProcAddr("glInvalidateFramebuffer"))
	gpInvalidateNamedFramebufferData = (C.GPINVALIDATENAMEDFRAMEBUFFERDATA)(getProcAddr("glInvalidateNamedFramebufferData"))
	gpInvalidateNamedFramebufferSubData = (C.GPINVALIDATENAMEDFRAMEBUFFERSUBDATA)(getProcAddr("glInvalidateNamedFramebufferSubData"))
	gpInvalidateSubFramebuffer = (C.GPINVALIDATESUBFRAMEBUFFER)(getProcAddr("glInvalidateSubFramebuffer"))
	gpInvalidateTexImage = (C.GPINVALIDATETEXIMAGE)(getProcAddr("glInvalidateTexImage"))
	gpInvalidateTexSubImage = (C.GPINVALIDATETEXSUBIMAGE)(getProcAddr("glInvalidateTexSubImage"))
	gpIsBuffer = (C.GPISBUFFER)(getProcAddr("glIsBuffer"))
	if gpIsBuffer == nil {
		return errors.New("glIsBuffer")
	}
	gpIsBufferResidentNV = (C.GPISBUFFERRESIDENTNV)(getProcAddr("glIsBufferResidentNV"))
	gpIsCommandListNV = (C.GPISCOMMANDLISTNV)(getProcAddr("glIsCommandListNV"))
	gpIsEnabled = (C.GPISENABLED)(getProcAddr("glIsEnabled"))
	if gpIsEnabled == nil {
		return errors.New("glIsEnabled")
	}
	gpIsEnabledIndexedEXT = (C.GPISENABLEDINDEXEDEXT)(getProcAddr("glIsEnabledIndexedEXT"))
	gpIsEnabledi = (C.GPISENABLEDI)(getProcAddr("glIsEnabledi"))
	if gpIsEnabledi == nil {
		return errors.New("glIsEnabledi")
	}
	gpIsFramebuffer = (C.GPISFRAMEBUFFER)(getProcAddr("glIsFramebuffer"))
	if gpIsFramebuffer == nil {
		return errors.New("glIsFramebuffer")
	}
	gpIsImageHandleResidentARB = (C.GPISIMAGEHANDLERESIDENTARB)(getProcAddr("glIsImageHandleResidentARB"))
	gpIsImageHandleResidentNV = (C.GPISIMAGEHANDLERESIDENTNV)(getProcAddr("glIsImageHandleResidentNV"))
	gpIsNamedBufferResidentNV = (C.GPISNAMEDBUFFERRESIDENTNV)(getProcAddr("glIsNamedBufferResidentNV"))
	gpIsNamedStringARB = (C.GPISNAMEDSTRINGARB)(getProcAddr("glIsNamedStringARB"))
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
	gpIsRenderbuffer = (C.GPISRENDERBUFFER)(getProcAddr("glIsRenderbuffer"))
	if gpIsRenderbuffer == nil {
		return errors.New("glIsRenderbuffer")
	}
	gpIsSampler = (C.GPISSAMPLER)(getProcAddr("glIsSampler"))
	if gpIsSampler == nil {
		return errors.New("glIsSampler")
	}
	gpIsShader = (C.GPISSHADER)(getProcAddr("glIsShader"))
	if gpIsShader == nil {
		return errors.New("glIsShader")
	}
	gpIsStateNV = (C.GPISSTATENV)(getProcAddr("glIsStateNV"))
	gpIsSync = (C.GPISSYNC)(getProcAddr("glIsSync"))
	if gpIsSync == nil {
		return errors.New("glIsSync")
	}
	gpIsTexture = (C.GPISTEXTURE)(getProcAddr("glIsTexture"))
	if gpIsTexture == nil {
		return errors.New("glIsTexture")
	}
	gpIsTextureHandleResidentARB = (C.GPISTEXTUREHANDLERESIDENTARB)(getProcAddr("glIsTextureHandleResidentARB"))
	gpIsTextureHandleResidentNV = (C.GPISTEXTUREHANDLERESIDENTNV)(getProcAddr("glIsTextureHandleResidentNV"))
	gpIsTransformFeedback = (C.GPISTRANSFORMFEEDBACK)(getProcAddr("glIsTransformFeedback"))
	if gpIsTransformFeedback == nil {
		return errors.New("glIsTransformFeedback")
	}
	gpIsVertexArray = (C.GPISVERTEXARRAY)(getProcAddr("glIsVertexArray"))
	if gpIsVertexArray == nil {
		return errors.New("glIsVertexArray")
	}
	gpLabelObjectEXT = (C.GPLABELOBJECTEXT)(getProcAddr("glLabelObjectEXT"))
	gpLineWidth = (C.GPLINEWIDTH)(getProcAddr("glLineWidth"))
	if gpLineWidth == nil {
		return errors.New("glLineWidth")
	}
	gpLinkProgram = (C.GPLINKPROGRAM)(getProcAddr("glLinkProgram"))
	if gpLinkProgram == nil {
		return errors.New("glLinkProgram")
	}
	gpListDrawCommandsStatesClientNV = (C.GPLISTDRAWCOMMANDSSTATESCLIENTNV)(getProcAddr("glListDrawCommandsStatesClientNV"))
	gpLogicOp = (C.GPLOGICOP)(getProcAddr("glLogicOp"))
	if gpLogicOp == nil {
		return errors.New("glLogicOp")
	}
	gpMakeBufferNonResidentNV = (C.GPMAKEBUFFERNONRESIDENTNV)(getProcAddr("glMakeBufferNonResidentNV"))
	gpMakeBufferResidentNV = (C.GPMAKEBUFFERRESIDENTNV)(getProcAddr("glMakeBufferResidentNV"))
	gpMakeImageHandleNonResidentARB = (C.GPMAKEIMAGEHANDLENONRESIDENTARB)(getProcAddr("glMakeImageHandleNonResidentARB"))
	gpMakeImageHandleNonResidentNV = (C.GPMAKEIMAGEHANDLENONRESIDENTNV)(getProcAddr("glMakeImageHandleNonResidentNV"))
	gpMakeImageHandleResidentARB = (C.GPMAKEIMAGEHANDLERESIDENTARB)(getProcAddr("glMakeImageHandleResidentARB"))
	gpMakeImageHandleResidentNV = (C.GPMAKEIMAGEHANDLERESIDENTNV)(getProcAddr("glMakeImageHandleResidentNV"))
	gpMakeNamedBufferNonResidentNV = (C.GPMAKENAMEDBUFFERNONRESIDENTNV)(getProcAddr("glMakeNamedBufferNonResidentNV"))
	gpMakeNamedBufferResidentNV = (C.GPMAKENAMEDBUFFERRESIDENTNV)(getProcAddr("glMakeNamedBufferResidentNV"))
	gpMakeTextureHandleNonResidentARB = (C.GPMAKETEXTUREHANDLENONRESIDENTARB)(getProcAddr("glMakeTextureHandleNonResidentARB"))
	gpMakeTextureHandleNonResidentNV = (C.GPMAKETEXTUREHANDLENONRESIDENTNV)(getProcAddr("glMakeTextureHandleNonResidentNV"))
	gpMakeTextureHandleResidentARB = (C.GPMAKETEXTUREHANDLERESIDENTARB)(getProcAddr("glMakeTextureHandleResidentARB"))
	gpMakeTextureHandleResidentNV = (C.GPMAKETEXTUREHANDLERESIDENTNV)(getProcAddr("glMakeTextureHandleResidentNV"))
	gpMapBuffer = (C.GPMAPBUFFER)(getProcAddr("glMapBuffer"))
	if gpMapBuffer == nil {
		return errors.New("glMapBuffer")
	}
	gpMapBufferRange = (C.GPMAPBUFFERRANGE)(getProcAddr("glMapBufferRange"))
	if gpMapBufferRange == nil {
		return errors.New("glMapBufferRange")
	}
	gpMapNamedBuffer = (C.GPMAPNAMEDBUFFER)(getProcAddr("glMapNamedBuffer"))
	gpMapNamedBufferEXT = (C.GPMAPNAMEDBUFFEREXT)(getProcAddr("glMapNamedBufferEXT"))
	gpMapNamedBufferRange = (C.GPMAPNAMEDBUFFERRANGE)(getProcAddr("glMapNamedBufferRange"))
	gpMapNamedBufferRangeEXT = (C.GPMAPNAMEDBUFFERRANGEEXT)(getProcAddr("glMapNamedBufferRangeEXT"))
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
	gpMaxShaderCompilerThreadsARB = (C.GPMAXSHADERCOMPILERTHREADSARB)(getProcAddr("glMaxShaderCompilerThreadsARB"))
	gpMaxShaderCompilerThreadsKHR = (C.GPMAXSHADERCOMPILERTHREADSKHR)(getProcAddr("glMaxShaderCompilerThreadsKHR"))
	gpMemoryBarrier = (C.GPMEMORYBARRIER)(getProcAddr("glMemoryBarrier"))
	gpMemoryBarrierByRegion = (C.GPMEMORYBARRIERBYREGION)(getProcAddr("glMemoryBarrierByRegion"))
	gpMinSampleShading = (C.GPMINSAMPLESHADING)(getProcAddr("glMinSampleShading"))
	if gpMinSampleShading == nil {
		return errors.New("glMinSampleShading")
	}
	gpMinSampleShadingARB = (C.GPMINSAMPLESHADINGARB)(getProcAddr("glMinSampleShadingARB"))
	gpMultiDrawArrays = (C.GPMULTIDRAWARRAYS)(getProcAddr("glMultiDrawArrays"))
	if gpMultiDrawArrays == nil {
		return errors.New("glMultiDrawArrays")
	}
	gpMultiDrawArraysIndirect = (C.GPMULTIDRAWARRAYSINDIRECT)(getProcAddr("glMultiDrawArraysIndirect"))
	gpMultiDrawArraysIndirectBindlessCountNV = (C.GPMULTIDRAWARRAYSINDIRECTBINDLESSCOUNTNV)(getProcAddr("glMultiDrawArraysIndirectBindlessCountNV"))
	gpMultiDrawArraysIndirectBindlessNV = (C.GPMULTIDRAWARRAYSINDIRECTBINDLESSNV)(getProcAddr("glMultiDrawArraysIndirectBindlessNV"))
	gpMultiDrawArraysIndirectCountARB = (C.GPMULTIDRAWARRAYSINDIRECTCOUNTARB)(getProcAddr("glMultiDrawArraysIndirectCountARB"))
	gpMultiDrawElements = (C.GPMULTIDRAWELEMENTS)(getProcAddr("glMultiDrawElements"))
	if gpMultiDrawElements == nil {
		return errors.New("glMultiDrawElements")
	}
	gpMultiDrawElementsBaseVertex = (C.GPMULTIDRAWELEMENTSBASEVERTEX)(getProcAddr("glMultiDrawElementsBaseVertex"))
	if gpMultiDrawElementsBaseVertex == nil {
		return errors.New("glMultiDrawElementsBaseVertex")
	}
	gpMultiDrawElementsIndirect = (C.GPMULTIDRAWELEMENTSINDIRECT)(getProcAddr("glMultiDrawElementsIndirect"))
	gpMultiDrawElementsIndirectBindlessCountNV = (C.GPMULTIDRAWELEMENTSINDIRECTBINDLESSCOUNTNV)(getProcAddr("glMultiDrawElementsIndirectBindlessCountNV"))
	gpMultiDrawElementsIndirectBindlessNV = (C.GPMULTIDRAWELEMENTSINDIRECTBINDLESSNV)(getProcAddr("glMultiDrawElementsIndirectBindlessNV"))
	gpMultiDrawElementsIndirectCountARB = (C.GPMULTIDRAWELEMENTSINDIRECTCOUNTARB)(getProcAddr("glMultiDrawElementsIndirectCountARB"))
	gpMultiTexBufferEXT = (C.GPMULTITEXBUFFEREXT)(getProcAddr("glMultiTexBufferEXT"))
	gpMultiTexCoordPointerEXT = (C.GPMULTITEXCOORDPOINTEREXT)(getProcAddr("glMultiTexCoordPointerEXT"))
	gpMultiTexEnvfEXT = (C.GPMULTITEXENVFEXT)(getProcAddr("glMultiTexEnvfEXT"))
	gpMultiTexEnvfvEXT = (C.GPMULTITEXENVFVEXT)(getProcAddr("glMultiTexEnvfvEXT"))
	gpMultiTexEnviEXT = (C.GPMULTITEXENVIEXT)(getProcAddr("glMultiTexEnviEXT"))
	gpMultiTexEnvivEXT = (C.GPMULTITEXENVIVEXT)(getProcAddr("glMultiTexEnvivEXT"))
	gpMultiTexGendEXT = (C.GPMULTITEXGENDEXT)(getProcAddr("glMultiTexGendEXT"))
	gpMultiTexGendvEXT = (C.GPMULTITEXGENDVEXT)(getProcAddr("glMultiTexGendvEXT"))
	gpMultiTexGenfEXT = (C.GPMULTITEXGENFEXT)(getProcAddr("glMultiTexGenfEXT"))
	gpMultiTexGenfvEXT = (C.GPMULTITEXGENFVEXT)(getProcAddr("glMultiTexGenfvEXT"))
	gpMultiTexGeniEXT = (C.GPMULTITEXGENIEXT)(getProcAddr("glMultiTexGeniEXT"))
	gpMultiTexGenivEXT = (C.GPMULTITEXGENIVEXT)(getProcAddr("glMultiTexGenivEXT"))
	gpMultiTexImage1DEXT = (C.GPMULTITEXIMAGE1DEXT)(getProcAddr("glMultiTexImage1DEXT"))
	gpMultiTexImage2DEXT = (C.GPMULTITEXIMAGE2DEXT)(getProcAddr("glMultiTexImage2DEXT"))
	gpMultiTexImage3DEXT = (C.GPMULTITEXIMAGE3DEXT)(getProcAddr("glMultiTexImage3DEXT"))
	gpMultiTexParameterIivEXT = (C.GPMULTITEXPARAMETERIIVEXT)(getProcAddr("glMultiTexParameterIivEXT"))
	gpMultiTexParameterIuivEXT = (C.GPMULTITEXPARAMETERIUIVEXT)(getProcAddr("glMultiTexParameterIuivEXT"))
	gpMultiTexParameterfEXT = (C.GPMULTITEXPARAMETERFEXT)(getProcAddr("glMultiTexParameterfEXT"))
	gpMultiTexParameterfvEXT = (C.GPMULTITEXPARAMETERFVEXT)(getProcAddr("glMultiTexParameterfvEXT"))
	gpMultiTexParameteriEXT = (C.GPMULTITEXPARAMETERIEXT)(getProcAddr("glMultiTexParameteriEXT"))
	gpMultiTexParameterivEXT = (C.GPMULTITEXPARAMETERIVEXT)(getProcAddr("glMultiTexParameterivEXT"))
	gpMultiTexRenderbufferEXT = (C.GPMULTITEXRENDERBUFFEREXT)(getProcAddr("glMultiTexRenderbufferEXT"))
	gpMultiTexSubImage1DEXT = (C.GPMULTITEXSUBIMAGE1DEXT)(getProcAddr("glMultiTexSubImage1DEXT"))
	gpMultiTexSubImage2DEXT = (C.GPMULTITEXSUBIMAGE2DEXT)(getProcAddr("glMultiTexSubImage2DEXT"))
	gpMultiTexSubImage3DEXT = (C.GPMULTITEXSUBIMAGE3DEXT)(getProcAddr("glMultiTexSubImage3DEXT"))
	gpNamedBufferData = (C.GPNAMEDBUFFERDATA)(getProcAddr("glNamedBufferData"))
	gpNamedBufferDataEXT = (C.GPNAMEDBUFFERDATAEXT)(getProcAddr("glNamedBufferDataEXT"))
	gpNamedBufferPageCommitmentARB = (C.GPNAMEDBUFFERPAGECOMMITMENTARB)(getProcAddr("glNamedBufferPageCommitmentARB"))
	gpNamedBufferPageCommitmentEXT = (C.GPNAMEDBUFFERPAGECOMMITMENTEXT)(getProcAddr("glNamedBufferPageCommitmentEXT"))
	gpNamedBufferStorage = (C.GPNAMEDBUFFERSTORAGE)(getProcAddr("glNamedBufferStorage"))
	gpNamedBufferStorageEXT = (C.GPNAMEDBUFFERSTORAGEEXT)(getProcAddr("glNamedBufferStorageEXT"))
	gpNamedBufferSubData = (C.GPNAMEDBUFFERSUBDATA)(getProcAddr("glNamedBufferSubData"))
	gpNamedBufferSubDataEXT = (C.GPNAMEDBUFFERSUBDATAEXT)(getProcAddr("glNamedBufferSubDataEXT"))
	gpNamedCopyBufferSubDataEXT = (C.GPNAMEDCOPYBUFFERSUBDATAEXT)(getProcAddr("glNamedCopyBufferSubDataEXT"))
	gpNamedFramebufferDrawBuffer = (C.GPNAMEDFRAMEBUFFERDRAWBUFFER)(getProcAddr("glNamedFramebufferDrawBuffer"))
	gpNamedFramebufferDrawBuffers = (C.GPNAMEDFRAMEBUFFERDRAWBUFFERS)(getProcAddr("glNamedFramebufferDrawBuffers"))
	gpNamedFramebufferParameteri = (C.GPNAMEDFRAMEBUFFERPARAMETERI)(getProcAddr("glNamedFramebufferParameteri"))
	gpNamedFramebufferParameteriEXT = (C.GPNAMEDFRAMEBUFFERPARAMETERIEXT)(getProcAddr("glNamedFramebufferParameteriEXT"))
	gpNamedFramebufferReadBuffer = (C.GPNAMEDFRAMEBUFFERREADBUFFER)(getProcAddr("glNamedFramebufferReadBuffer"))
	gpNamedFramebufferRenderbuffer = (C.GPNAMEDFRAMEBUFFERRENDERBUFFER)(getProcAddr("glNamedFramebufferRenderbuffer"))
	gpNamedFramebufferRenderbufferEXT = (C.GPNAMEDFRAMEBUFFERRENDERBUFFEREXT)(getProcAddr("glNamedFramebufferRenderbufferEXT"))
	gpNamedFramebufferSampleLocationsfvARB = (C.GPNAMEDFRAMEBUFFERSAMPLELOCATIONSFVARB)(getProcAddr("glNamedFramebufferSampleLocationsfvARB"))
	gpNamedFramebufferSampleLocationsfvNV = (C.GPNAMEDFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glNamedFramebufferSampleLocationsfvNV"))
	gpNamedFramebufferTexture = (C.GPNAMEDFRAMEBUFFERTEXTURE)(getProcAddr("glNamedFramebufferTexture"))
	gpNamedFramebufferTexture1DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE1DEXT)(getProcAddr("glNamedFramebufferTexture1DEXT"))
	gpNamedFramebufferTexture2DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE2DEXT)(getProcAddr("glNamedFramebufferTexture2DEXT"))
	gpNamedFramebufferTexture3DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE3DEXT)(getProcAddr("glNamedFramebufferTexture3DEXT"))
	gpNamedFramebufferTextureEXT = (C.GPNAMEDFRAMEBUFFERTEXTUREEXT)(getProcAddr("glNamedFramebufferTextureEXT"))
	gpNamedFramebufferTextureFaceEXT = (C.GPNAMEDFRAMEBUFFERTEXTUREFACEEXT)(getProcAddr("glNamedFramebufferTextureFaceEXT"))
	gpNamedFramebufferTextureLayer = (C.GPNAMEDFRAMEBUFFERTEXTURELAYER)(getProcAddr("glNamedFramebufferTextureLayer"))
	gpNamedFramebufferTextureLayerEXT = (C.GPNAMEDFRAMEBUFFERTEXTURELAYEREXT)(getProcAddr("glNamedFramebufferTextureLayerEXT"))
	gpNamedProgramLocalParameter4dEXT = (C.GPNAMEDPROGRAMLOCALPARAMETER4DEXT)(getProcAddr("glNamedProgramLocalParameter4dEXT"))
	gpNamedProgramLocalParameter4dvEXT = (C.GPNAMEDPROGRAMLOCALPARAMETER4DVEXT)(getProcAddr("glNamedProgramLocalParameter4dvEXT"))
	gpNamedProgramLocalParameter4fEXT = (C.GPNAMEDPROGRAMLOCALPARAMETER4FEXT)(getProcAddr("glNamedProgramLocalParameter4fEXT"))
	gpNamedProgramLocalParameter4fvEXT = (C.GPNAMEDPROGRAMLOCALPARAMETER4FVEXT)(getProcAddr("glNamedProgramLocalParameter4fvEXT"))
	gpNamedProgramLocalParameterI4iEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERI4IEXT)(getProcAddr("glNamedProgramLocalParameterI4iEXT"))
	gpNamedProgramLocalParameterI4ivEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERI4IVEXT)(getProcAddr("glNamedProgramLocalParameterI4ivEXT"))
	gpNamedProgramLocalParameterI4uiEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERI4UIEXT)(getProcAddr("glNamedProgramLocalParameterI4uiEXT"))
	gpNamedProgramLocalParameterI4uivEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERI4UIVEXT)(getProcAddr("glNamedProgramLocalParameterI4uivEXT"))
	gpNamedProgramLocalParameters4fvEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERS4FVEXT)(getProcAddr("glNamedProgramLocalParameters4fvEXT"))
	gpNamedProgramLocalParametersI4ivEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERSI4IVEXT)(getProcAddr("glNamedProgramLocalParametersI4ivEXT"))
	gpNamedProgramLocalParametersI4uivEXT = (C.GPNAMEDPROGRAMLOCALPARAMETERSI4UIVEXT)(getProcAddr("glNamedProgramLocalParametersI4uivEXT"))
	gpNamedProgramStringEXT = (C.GPNAMEDPROGRAMSTRINGEXT)(getProcAddr("glNamedProgramStringEXT"))
	gpNamedRenderbufferStorage = (C.GPNAMEDRENDERBUFFERSTORAGE)(getProcAddr("glNamedRenderbufferStorage"))
	gpNamedRenderbufferStorageEXT = (C.GPNAMEDRENDERBUFFERSTORAGEEXT)(getProcAddr("glNamedRenderbufferStorageEXT"))
	gpNamedRenderbufferStorageMultisample = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLE)(getProcAddr("glNamedRenderbufferStorageMultisample"))
	gpNamedRenderbufferStorageMultisampleCoverageEXT = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLECOVERAGEEXT)(getProcAddr("glNamedRenderbufferStorageMultisampleCoverageEXT"))
	gpNamedRenderbufferStorageMultisampleEXT = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLEEXT)(getProcAddr("glNamedRenderbufferStorageMultisampleEXT"))
	gpNamedStringARB = (C.GPNAMEDSTRINGARB)(getProcAddr("glNamedStringARB"))
	gpNormalFormatNV = (C.GPNORMALFORMATNV)(getProcAddr("glNormalFormatNV"))
	gpObjectLabel = (C.GPOBJECTLABEL)(getProcAddr("glObjectLabel"))
	gpObjectLabelKHR = (C.GPOBJECTLABELKHR)(getProcAddr("glObjectLabelKHR"))
	gpObjectPtrLabel = (C.GPOBJECTPTRLABEL)(getProcAddr("glObjectPtrLabel"))
	gpObjectPtrLabelKHR = (C.GPOBJECTPTRLABELKHR)(getProcAddr("glObjectPtrLabelKHR"))
	gpPatchParameterfv = (C.GPPATCHPARAMETERFV)(getProcAddr("glPatchParameterfv"))
	if gpPatchParameterfv == nil {
		return errors.New("glPatchParameterfv")
	}
	gpPatchParameteri = (C.GPPATCHPARAMETERI)(getProcAddr("glPatchParameteri"))
	if gpPatchParameteri == nil {
		return errors.New("glPatchParameteri")
	}
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
	gpPixelStoref = (C.GPPIXELSTOREF)(getProcAddr("glPixelStoref"))
	if gpPixelStoref == nil {
		return errors.New("glPixelStoref")
	}
	gpPixelStorei = (C.GPPIXELSTOREI)(getProcAddr("glPixelStorei"))
	if gpPixelStorei == nil {
		return errors.New("glPixelStorei")
	}
	gpPointAlongPathNV = (C.GPPOINTALONGPATHNV)(getProcAddr("glPointAlongPathNV"))
	gpPointParameterf = (C.GPPOINTPARAMETERF)(getProcAddr("glPointParameterf"))
	if gpPointParameterf == nil {
		return errors.New("glPointParameterf")
	}
	gpPointParameterfv = (C.GPPOINTPARAMETERFV)(getProcAddr("glPointParameterfv"))
	if gpPointParameterfv == nil {
		return errors.New("glPointParameterfv")
	}
	gpPointParameteri = (C.GPPOINTPARAMETERI)(getProcAddr("glPointParameteri"))
	if gpPointParameteri == nil {
		return errors.New("glPointParameteri")
	}
	gpPointParameteriv = (C.GPPOINTPARAMETERIV)(getProcAddr("glPointParameteriv"))
	if gpPointParameteriv == nil {
		return errors.New("glPointParameteriv")
	}
	gpPointSize = (C.GPPOINTSIZE)(getProcAddr("glPointSize"))
	if gpPointSize == nil {
		return errors.New("glPointSize")
	}
	gpPolygonMode = (C.GPPOLYGONMODE)(getProcAddr("glPolygonMode"))
	if gpPolygonMode == nil {
		return errors.New("glPolygonMode")
	}
	gpPolygonOffset = (C.GPPOLYGONOFFSET)(getProcAddr("glPolygonOffset"))
	if gpPolygonOffset == nil {
		return errors.New("glPolygonOffset")
	}
	gpPolygonOffsetClamp = (C.GPPOLYGONOFFSETCLAMP)(getProcAddr("glPolygonOffsetClamp"))
	gpPolygonOffsetClampEXT = (C.GPPOLYGONOFFSETCLAMPEXT)(getProcAddr("glPolygonOffsetClampEXT"))
	gpPopDebugGroup = (C.GPPOPDEBUGGROUP)(getProcAddr("glPopDebugGroup"))
	gpPopDebugGroupKHR = (C.GPPOPDEBUGGROUPKHR)(getProcAddr("glPopDebugGroupKHR"))
	gpPopGroupMarkerEXT = (C.GPPOPGROUPMARKEREXT)(getProcAddr("glPopGroupMarkerEXT"))
	gpPrimitiveBoundingBoxARB = (C.GPPRIMITIVEBOUNDINGBOXARB)(getProcAddr("glPrimitiveBoundingBoxARB"))
	gpPrimitiveRestartIndex = (C.GPPRIMITIVERESTARTINDEX)(getProcAddr("glPrimitiveRestartIndex"))
	if gpPrimitiveRestartIndex == nil {
		return errors.New("glPrimitiveRestartIndex")
	}
	gpProgramBinary = (C.GPPROGRAMBINARY)(getProcAddr("glProgramBinary"))
	if gpProgramBinary == nil {
		return errors.New("glProgramBinary")
	}
	gpProgramParameteri = (C.GPPROGRAMPARAMETERI)(getProcAddr("glProgramParameteri"))
	if gpProgramParameteri == nil {
		return errors.New("glProgramParameteri")
	}
	gpProgramParameteriARB = (C.GPPROGRAMPARAMETERIARB)(getProcAddr("glProgramParameteriARB"))
	gpProgramParameteriEXT = (C.GPPROGRAMPARAMETERIEXT)(getProcAddr("glProgramParameteriEXT"))
	gpProgramPathFragmentInputGenNV = (C.GPPROGRAMPATHFRAGMENTINPUTGENNV)(getProcAddr("glProgramPathFragmentInputGenNV"))
	gpProgramUniform1d = (C.GPPROGRAMUNIFORM1D)(getProcAddr("glProgramUniform1d"))
	if gpProgramUniform1d == nil {
		return errors.New("glProgramUniform1d")
	}
	gpProgramUniform1dEXT = (C.GPPROGRAMUNIFORM1DEXT)(getProcAddr("glProgramUniform1dEXT"))
	gpProgramUniform1dv = (C.GPPROGRAMUNIFORM1DV)(getProcAddr("glProgramUniform1dv"))
	if gpProgramUniform1dv == nil {
		return errors.New("glProgramUniform1dv")
	}
	gpProgramUniform1dvEXT = (C.GPPROGRAMUNIFORM1DVEXT)(getProcAddr("glProgramUniform1dvEXT"))
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
	gpProgramUniform1i64ARB = (C.GPPROGRAMUNIFORM1I64ARB)(getProcAddr("glProgramUniform1i64ARB"))
	gpProgramUniform1i64NV = (C.GPPROGRAMUNIFORM1I64NV)(getProcAddr("glProgramUniform1i64NV"))
	gpProgramUniform1i64vARB = (C.GPPROGRAMUNIFORM1I64VARB)(getProcAddr("glProgramUniform1i64vARB"))
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
	gpProgramUniform1ui64ARB = (C.GPPROGRAMUNIFORM1UI64ARB)(getProcAddr("glProgramUniform1ui64ARB"))
	gpProgramUniform1ui64NV = (C.GPPROGRAMUNIFORM1UI64NV)(getProcAddr("glProgramUniform1ui64NV"))
	gpProgramUniform1ui64vARB = (C.GPPROGRAMUNIFORM1UI64VARB)(getProcAddr("glProgramUniform1ui64vARB"))
	gpProgramUniform1ui64vNV = (C.GPPROGRAMUNIFORM1UI64VNV)(getProcAddr("glProgramUniform1ui64vNV"))
	gpProgramUniform1uiEXT = (C.GPPROGRAMUNIFORM1UIEXT)(getProcAddr("glProgramUniform1uiEXT"))
	gpProgramUniform1uiv = (C.GPPROGRAMUNIFORM1UIV)(getProcAddr("glProgramUniform1uiv"))
	if gpProgramUniform1uiv == nil {
		return errors.New("glProgramUniform1uiv")
	}
	gpProgramUniform1uivEXT = (C.GPPROGRAMUNIFORM1UIVEXT)(getProcAddr("glProgramUniform1uivEXT"))
	gpProgramUniform2d = (C.GPPROGRAMUNIFORM2D)(getProcAddr("glProgramUniform2d"))
	if gpProgramUniform2d == nil {
		return errors.New("glProgramUniform2d")
	}
	gpProgramUniform2dEXT = (C.GPPROGRAMUNIFORM2DEXT)(getProcAddr("glProgramUniform2dEXT"))
	gpProgramUniform2dv = (C.GPPROGRAMUNIFORM2DV)(getProcAddr("glProgramUniform2dv"))
	if gpProgramUniform2dv == nil {
		return errors.New("glProgramUniform2dv")
	}
	gpProgramUniform2dvEXT = (C.GPPROGRAMUNIFORM2DVEXT)(getProcAddr("glProgramUniform2dvEXT"))
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
	gpProgramUniform2i64ARB = (C.GPPROGRAMUNIFORM2I64ARB)(getProcAddr("glProgramUniform2i64ARB"))
	gpProgramUniform2i64NV = (C.GPPROGRAMUNIFORM2I64NV)(getProcAddr("glProgramUniform2i64NV"))
	gpProgramUniform2i64vARB = (C.GPPROGRAMUNIFORM2I64VARB)(getProcAddr("glProgramUniform2i64vARB"))
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
	gpProgramUniform2ui64ARB = (C.GPPROGRAMUNIFORM2UI64ARB)(getProcAddr("glProgramUniform2ui64ARB"))
	gpProgramUniform2ui64NV = (C.GPPROGRAMUNIFORM2UI64NV)(getProcAddr("glProgramUniform2ui64NV"))
	gpProgramUniform2ui64vARB = (C.GPPROGRAMUNIFORM2UI64VARB)(getProcAddr("glProgramUniform2ui64vARB"))
	gpProgramUniform2ui64vNV = (C.GPPROGRAMUNIFORM2UI64VNV)(getProcAddr("glProgramUniform2ui64vNV"))
	gpProgramUniform2uiEXT = (C.GPPROGRAMUNIFORM2UIEXT)(getProcAddr("glProgramUniform2uiEXT"))
	gpProgramUniform2uiv = (C.GPPROGRAMUNIFORM2UIV)(getProcAddr("glProgramUniform2uiv"))
	if gpProgramUniform2uiv == nil {
		return errors.New("glProgramUniform2uiv")
	}
	gpProgramUniform2uivEXT = (C.GPPROGRAMUNIFORM2UIVEXT)(getProcAddr("glProgramUniform2uivEXT"))
	gpProgramUniform3d = (C.GPPROGRAMUNIFORM3D)(getProcAddr("glProgramUniform3d"))
	if gpProgramUniform3d == nil {
		return errors.New("glProgramUniform3d")
	}
	gpProgramUniform3dEXT = (C.GPPROGRAMUNIFORM3DEXT)(getProcAddr("glProgramUniform3dEXT"))
	gpProgramUniform3dv = (C.GPPROGRAMUNIFORM3DV)(getProcAddr("glProgramUniform3dv"))
	if gpProgramUniform3dv == nil {
		return errors.New("glProgramUniform3dv")
	}
	gpProgramUniform3dvEXT = (C.GPPROGRAMUNIFORM3DVEXT)(getProcAddr("glProgramUniform3dvEXT"))
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
	gpProgramUniform3i64ARB = (C.GPPROGRAMUNIFORM3I64ARB)(getProcAddr("glProgramUniform3i64ARB"))
	gpProgramUniform3i64NV = (C.GPPROGRAMUNIFORM3I64NV)(getProcAddr("glProgramUniform3i64NV"))
	gpProgramUniform3i64vARB = (C.GPPROGRAMUNIFORM3I64VARB)(getProcAddr("glProgramUniform3i64vARB"))
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
	gpProgramUniform3ui64ARB = (C.GPPROGRAMUNIFORM3UI64ARB)(getProcAddr("glProgramUniform3ui64ARB"))
	gpProgramUniform3ui64NV = (C.GPPROGRAMUNIFORM3UI64NV)(getProcAddr("glProgramUniform3ui64NV"))
	gpProgramUniform3ui64vARB = (C.GPPROGRAMUNIFORM3UI64VARB)(getProcAddr("glProgramUniform3ui64vARB"))
	gpProgramUniform3ui64vNV = (C.GPPROGRAMUNIFORM3UI64VNV)(getProcAddr("glProgramUniform3ui64vNV"))
	gpProgramUniform3uiEXT = (C.GPPROGRAMUNIFORM3UIEXT)(getProcAddr("glProgramUniform3uiEXT"))
	gpProgramUniform3uiv = (C.GPPROGRAMUNIFORM3UIV)(getProcAddr("glProgramUniform3uiv"))
	if gpProgramUniform3uiv == nil {
		return errors.New("glProgramUniform3uiv")
	}
	gpProgramUniform3uivEXT = (C.GPPROGRAMUNIFORM3UIVEXT)(getProcAddr("glProgramUniform3uivEXT"))
	gpProgramUniform4d = (C.GPPROGRAMUNIFORM4D)(getProcAddr("glProgramUniform4d"))
	if gpProgramUniform4d == nil {
		return errors.New("glProgramUniform4d")
	}
	gpProgramUniform4dEXT = (C.GPPROGRAMUNIFORM4DEXT)(getProcAddr("glProgramUniform4dEXT"))
	gpProgramUniform4dv = (C.GPPROGRAMUNIFORM4DV)(getProcAddr("glProgramUniform4dv"))
	if gpProgramUniform4dv == nil {
		return errors.New("glProgramUniform4dv")
	}
	gpProgramUniform4dvEXT = (C.GPPROGRAMUNIFORM4DVEXT)(getProcAddr("glProgramUniform4dvEXT"))
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
	gpProgramUniform4i64ARB = (C.GPPROGRAMUNIFORM4I64ARB)(getProcAddr("glProgramUniform4i64ARB"))
	gpProgramUniform4i64NV = (C.GPPROGRAMUNIFORM4I64NV)(getProcAddr("glProgramUniform4i64NV"))
	gpProgramUniform4i64vARB = (C.GPPROGRAMUNIFORM4I64VARB)(getProcAddr("glProgramUniform4i64vARB"))
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
	gpProgramUniform4ui64ARB = (C.GPPROGRAMUNIFORM4UI64ARB)(getProcAddr("glProgramUniform4ui64ARB"))
	gpProgramUniform4ui64NV = (C.GPPROGRAMUNIFORM4UI64NV)(getProcAddr("glProgramUniform4ui64NV"))
	gpProgramUniform4ui64vARB = (C.GPPROGRAMUNIFORM4UI64VARB)(getProcAddr("glProgramUniform4ui64vARB"))
	gpProgramUniform4ui64vNV = (C.GPPROGRAMUNIFORM4UI64VNV)(getProcAddr("glProgramUniform4ui64vNV"))
	gpProgramUniform4uiEXT = (C.GPPROGRAMUNIFORM4UIEXT)(getProcAddr("glProgramUniform4uiEXT"))
	gpProgramUniform4uiv = (C.GPPROGRAMUNIFORM4UIV)(getProcAddr("glProgramUniform4uiv"))
	if gpProgramUniform4uiv == nil {
		return errors.New("glProgramUniform4uiv")
	}
	gpProgramUniform4uivEXT = (C.GPPROGRAMUNIFORM4UIVEXT)(getProcAddr("glProgramUniform4uivEXT"))
	gpProgramUniformHandleui64ARB = (C.GPPROGRAMUNIFORMHANDLEUI64ARB)(getProcAddr("glProgramUniformHandleui64ARB"))
	gpProgramUniformHandleui64NV = (C.GPPROGRAMUNIFORMHANDLEUI64NV)(getProcAddr("glProgramUniformHandleui64NV"))
	gpProgramUniformHandleui64vARB = (C.GPPROGRAMUNIFORMHANDLEUI64VARB)(getProcAddr("glProgramUniformHandleui64vARB"))
	gpProgramUniformHandleui64vNV = (C.GPPROGRAMUNIFORMHANDLEUI64VNV)(getProcAddr("glProgramUniformHandleui64vNV"))
	gpProgramUniformMatrix2dv = (C.GPPROGRAMUNIFORMMATRIX2DV)(getProcAddr("glProgramUniformMatrix2dv"))
	if gpProgramUniformMatrix2dv == nil {
		return errors.New("glProgramUniformMatrix2dv")
	}
	gpProgramUniformMatrix2dvEXT = (C.GPPROGRAMUNIFORMMATRIX2DVEXT)(getProcAddr("glProgramUniformMatrix2dvEXT"))
	gpProgramUniformMatrix2fv = (C.GPPROGRAMUNIFORMMATRIX2FV)(getProcAddr("glProgramUniformMatrix2fv"))
	if gpProgramUniformMatrix2fv == nil {
		return errors.New("glProgramUniformMatrix2fv")
	}
	gpProgramUniformMatrix2fvEXT = (C.GPPROGRAMUNIFORMMATRIX2FVEXT)(getProcAddr("glProgramUniformMatrix2fvEXT"))
	gpProgramUniformMatrix2x3dv = (C.GPPROGRAMUNIFORMMATRIX2X3DV)(getProcAddr("glProgramUniformMatrix2x3dv"))
	if gpProgramUniformMatrix2x3dv == nil {
		return errors.New("glProgramUniformMatrix2x3dv")
	}
	gpProgramUniformMatrix2x3dvEXT = (C.GPPROGRAMUNIFORMMATRIX2X3DVEXT)(getProcAddr("glProgramUniformMatrix2x3dvEXT"))
	gpProgramUniformMatrix2x3fv = (C.GPPROGRAMUNIFORMMATRIX2X3FV)(getProcAddr("glProgramUniformMatrix2x3fv"))
	if gpProgramUniformMatrix2x3fv == nil {
		return errors.New("glProgramUniformMatrix2x3fv")
	}
	gpProgramUniformMatrix2x3fvEXT = (C.GPPROGRAMUNIFORMMATRIX2X3FVEXT)(getProcAddr("glProgramUniformMatrix2x3fvEXT"))
	gpProgramUniformMatrix2x4dv = (C.GPPROGRAMUNIFORMMATRIX2X4DV)(getProcAddr("glProgramUniformMatrix2x4dv"))
	if gpProgramUniformMatrix2x4dv == nil {
		return errors.New("glProgramUniformMatrix2x4dv")
	}
	gpProgramUniformMatrix2x4dvEXT = (C.GPPROGRAMUNIFORMMATRIX2X4DVEXT)(getProcAddr("glProgramUniformMatrix2x4dvEXT"))
	gpProgramUniformMatrix2x4fv = (C.GPPROGRAMUNIFORMMATRIX2X4FV)(getProcAddr("glProgramUniformMatrix2x4fv"))
	if gpProgramUniformMatrix2x4fv == nil {
		return errors.New("glProgramUniformMatrix2x4fv")
	}
	gpProgramUniformMatrix2x4fvEXT = (C.GPPROGRAMUNIFORMMATRIX2X4FVEXT)(getProcAddr("glProgramUniformMatrix2x4fvEXT"))
	gpProgramUniformMatrix3dv = (C.GPPROGRAMUNIFORMMATRIX3DV)(getProcAddr("glProgramUniformMatrix3dv"))
	if gpProgramUniformMatrix3dv == nil {
		return errors.New("glProgramUniformMatrix3dv")
	}
	gpProgramUniformMatrix3dvEXT = (C.GPPROGRAMUNIFORMMATRIX3DVEXT)(getProcAddr("glProgramUniformMatrix3dvEXT"))
	gpProgramUniformMatrix3fv = (C.GPPROGRAMUNIFORMMATRIX3FV)(getProcAddr("glProgramUniformMatrix3fv"))
	if gpProgramUniformMatrix3fv == nil {
		return errors.New("glProgramUniformMatrix3fv")
	}
	gpProgramUniformMatrix3fvEXT = (C.GPPROGRAMUNIFORMMATRIX3FVEXT)(getProcAddr("glProgramUniformMatrix3fvEXT"))
	gpProgramUniformMatrix3x2dv = (C.GPPROGRAMUNIFORMMATRIX3X2DV)(getProcAddr("glProgramUniformMatrix3x2dv"))
	if gpProgramUniformMatrix3x2dv == nil {
		return errors.New("glProgramUniformMatrix3x2dv")
	}
	gpProgramUniformMatrix3x2dvEXT = (C.GPPROGRAMUNIFORMMATRIX3X2DVEXT)(getProcAddr("glProgramUniformMatrix3x2dvEXT"))
	gpProgramUniformMatrix3x2fv = (C.GPPROGRAMUNIFORMMATRIX3X2FV)(getProcAddr("glProgramUniformMatrix3x2fv"))
	if gpProgramUniformMatrix3x2fv == nil {
		return errors.New("glProgramUniformMatrix3x2fv")
	}
	gpProgramUniformMatrix3x2fvEXT = (C.GPPROGRAMUNIFORMMATRIX3X2FVEXT)(getProcAddr("glProgramUniformMatrix3x2fvEXT"))
	gpProgramUniformMatrix3x4dv = (C.GPPROGRAMUNIFORMMATRIX3X4DV)(getProcAddr("glProgramUniformMatrix3x4dv"))
	if gpProgramUniformMatrix3x4dv == nil {
		return errors.New("glProgramUniformMatrix3x4dv")
	}
	gpProgramUniformMatrix3x4dvEXT = (C.GPPROGRAMUNIFORMMATRIX3X4DVEXT)(getProcAddr("glProgramUniformMatrix3x4dvEXT"))
	gpProgramUniformMatrix3x4fv = (C.GPPROGRAMUNIFORMMATRIX3X4FV)(getProcAddr("glProgramUniformMatrix3x4fv"))
	if gpProgramUniformMatrix3x4fv == nil {
		return errors.New("glProgramUniformMatrix3x4fv")
	}
	gpProgramUniformMatrix3x4fvEXT = (C.GPPROGRAMUNIFORMMATRIX3X4FVEXT)(getProcAddr("glProgramUniformMatrix3x4fvEXT"))
	gpProgramUniformMatrix4dv = (C.GPPROGRAMUNIFORMMATRIX4DV)(getProcAddr("glProgramUniformMatrix4dv"))
	if gpProgramUniformMatrix4dv == nil {
		return errors.New("glProgramUniformMatrix4dv")
	}
	gpProgramUniformMatrix4dvEXT = (C.GPPROGRAMUNIFORMMATRIX4DVEXT)(getProcAddr("glProgramUniformMatrix4dvEXT"))
	gpProgramUniformMatrix4fv = (C.GPPROGRAMUNIFORMMATRIX4FV)(getProcAddr("glProgramUniformMatrix4fv"))
	if gpProgramUniformMatrix4fv == nil {
		return errors.New("glProgramUniformMatrix4fv")
	}
	gpProgramUniformMatrix4fvEXT = (C.GPPROGRAMUNIFORMMATRIX4FVEXT)(getProcAddr("glProgramUniformMatrix4fvEXT"))
	gpProgramUniformMatrix4x2dv = (C.GPPROGRAMUNIFORMMATRIX4X2DV)(getProcAddr("glProgramUniformMatrix4x2dv"))
	if gpProgramUniformMatrix4x2dv == nil {
		return errors.New("glProgramUniformMatrix4x2dv")
	}
	gpProgramUniformMatrix4x2dvEXT = (C.GPPROGRAMUNIFORMMATRIX4X2DVEXT)(getProcAddr("glProgramUniformMatrix4x2dvEXT"))
	gpProgramUniformMatrix4x2fv = (C.GPPROGRAMUNIFORMMATRIX4X2FV)(getProcAddr("glProgramUniformMatrix4x2fv"))
	if gpProgramUniformMatrix4x2fv == nil {
		return errors.New("glProgramUniformMatrix4x2fv")
	}
	gpProgramUniformMatrix4x2fvEXT = (C.GPPROGRAMUNIFORMMATRIX4X2FVEXT)(getProcAddr("glProgramUniformMatrix4x2fvEXT"))
	gpProgramUniformMatrix4x3dv = (C.GPPROGRAMUNIFORMMATRIX4X3DV)(getProcAddr("glProgramUniformMatrix4x3dv"))
	if gpProgramUniformMatrix4x3dv == nil {
		return errors.New("glProgramUniformMatrix4x3dv")
	}
	gpProgramUniformMatrix4x3dvEXT = (C.GPPROGRAMUNIFORMMATRIX4X3DVEXT)(getProcAddr("glProgramUniformMatrix4x3dvEXT"))
	gpProgramUniformMatrix4x3fv = (C.GPPROGRAMUNIFORMMATRIX4X3FV)(getProcAddr("glProgramUniformMatrix4x3fv"))
	if gpProgramUniformMatrix4x3fv == nil {
		return errors.New("glProgramUniformMatrix4x3fv")
	}
	gpProgramUniformMatrix4x3fvEXT = (C.GPPROGRAMUNIFORMMATRIX4X3FVEXT)(getProcAddr("glProgramUniformMatrix4x3fvEXT"))
	gpProgramUniformui64NV = (C.GPPROGRAMUNIFORMUI64NV)(getProcAddr("glProgramUniformui64NV"))
	gpProgramUniformui64vNV = (C.GPPROGRAMUNIFORMUI64VNV)(getProcAddr("glProgramUniformui64vNV"))
	gpProvokingVertex = (C.GPPROVOKINGVERTEX)(getProcAddr("glProvokingVertex"))
	if gpProvokingVertex == nil {
		return errors.New("glProvokingVertex")
	}
	gpPushClientAttribDefaultEXT = (C.GPPUSHCLIENTATTRIBDEFAULTEXT)(getProcAddr("glPushClientAttribDefaultEXT"))
	gpPushDebugGroup = (C.GPPUSHDEBUGGROUP)(getProcAddr("glPushDebugGroup"))
	gpPushDebugGroupKHR = (C.GPPUSHDEBUGGROUPKHR)(getProcAddr("glPushDebugGroupKHR"))
	gpPushGroupMarkerEXT = (C.GPPUSHGROUPMARKEREXT)(getProcAddr("glPushGroupMarkerEXT"))
	gpQueryCounter = (C.GPQUERYCOUNTER)(getProcAddr("glQueryCounter"))
	if gpQueryCounter == nil {
		return errors.New("glQueryCounter")
	}
	gpRasterSamplesEXT = (C.GPRASTERSAMPLESEXT)(getProcAddr("glRasterSamplesEXT"))
	gpReadBuffer = (C.GPREADBUFFER)(getProcAddr("glReadBuffer"))
	if gpReadBuffer == nil {
		return errors.New("glReadBuffer")
	}
	gpReadPixels = (C.GPREADPIXELS)(getProcAddr("glReadPixels"))
	if gpReadPixels == nil {
		return errors.New("glReadPixels")
	}
	gpReadnPixels = (C.GPREADNPIXELS)(getProcAddr("glReadnPixels"))
	gpReadnPixelsARB = (C.GPREADNPIXELSARB)(getProcAddr("glReadnPixelsARB"))
	gpReadnPixelsKHR = (C.GPREADNPIXELSKHR)(getProcAddr("glReadnPixelsKHR"))
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
	gpRenderbufferStorageMultisampleCoverageNV = (C.GPRENDERBUFFERSTORAGEMULTISAMPLECOVERAGENV)(getProcAddr("glRenderbufferStorageMultisampleCoverageNV"))
	gpResolveDepthValuesNV = (C.GPRESOLVEDEPTHVALUESNV)(getProcAddr("glResolveDepthValuesNV"))
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
	gpSamplerParameterIiv = (C.GPSAMPLERPARAMETERIIV)(getProcAddr("glSamplerParameterIiv"))
	if gpSamplerParameterIiv == nil {
		return errors.New("glSamplerParameterIiv")
	}
	gpSamplerParameterIuiv = (C.GPSAMPLERPARAMETERIUIV)(getProcAddr("glSamplerParameterIuiv"))
	if gpSamplerParameterIuiv == nil {
		return errors.New("glSamplerParameterIuiv")
	}
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
	gpScissorArrayv = (C.GPSCISSORARRAYV)(getProcAddr("glScissorArrayv"))
	if gpScissorArrayv == nil {
		return errors.New("glScissorArrayv")
	}
	gpScissorIndexed = (C.GPSCISSORINDEXED)(getProcAddr("glScissorIndexed"))
	if gpScissorIndexed == nil {
		return errors.New("glScissorIndexed")
	}
	gpScissorIndexedv = (C.GPSCISSORINDEXEDV)(getProcAddr("glScissorIndexedv"))
	if gpScissorIndexedv == nil {
		return errors.New("glScissorIndexedv")
	}
	gpSecondaryColorFormatNV = (C.GPSECONDARYCOLORFORMATNV)(getProcAddr("glSecondaryColorFormatNV"))
	gpSelectPerfMonitorCountersAMD = (C.GPSELECTPERFMONITORCOUNTERSAMD)(getProcAddr("glSelectPerfMonitorCountersAMD"))
	gpShaderBinary = (C.GPSHADERBINARY)(getProcAddr("glShaderBinary"))
	if gpShaderBinary == nil {
		return errors.New("glShaderBinary")
	}
	gpShaderSource = (C.GPSHADERSOURCE)(getProcAddr("glShaderSource"))
	if gpShaderSource == nil {
		return errors.New("glShaderSource")
	}
	gpShaderStorageBlockBinding = (C.GPSHADERSTORAGEBLOCKBINDING)(getProcAddr("glShaderStorageBlockBinding"))
	gpSignalVkFenceNV = (C.GPSIGNALVKFENCENV)(getProcAddr("glSignalVkFenceNV"))
	gpSignalVkSemaphoreNV = (C.GPSIGNALVKSEMAPHORENV)(getProcAddr("glSignalVkSemaphoreNV"))
	gpSpecializeShaderARB = (C.GPSPECIALIZESHADERARB)(getProcAddr("glSpecializeShaderARB"))
	gpStateCaptureNV = (C.GPSTATECAPTURENV)(getProcAddr("glStateCaptureNV"))
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
	gpTexBuffer = (C.GPTEXBUFFER)(getProcAddr("glTexBuffer"))
	if gpTexBuffer == nil {
		return errors.New("glTexBuffer")
	}
	gpTexBufferARB = (C.GPTEXBUFFERARB)(getProcAddr("glTexBufferARB"))
	gpTexBufferRange = (C.GPTEXBUFFERRANGE)(getProcAddr("glTexBufferRange"))
	gpTexCoordFormatNV = (C.GPTEXCOORDFORMATNV)(getProcAddr("glTexCoordFormatNV"))
	gpTexImage1D = (C.GPTEXIMAGE1D)(getProcAddr("glTexImage1D"))
	if gpTexImage1D == nil {
		return errors.New("glTexImage1D")
	}
	gpTexImage2D = (C.GPTEXIMAGE2D)(getProcAddr("glTexImage2D"))
	if gpTexImage2D == nil {
		return errors.New("glTexImage2D")
	}
	gpTexImage2DMultisample = (C.GPTEXIMAGE2DMULTISAMPLE)(getProcAddr("glTexImage2DMultisample"))
	if gpTexImage2DMultisample == nil {
		return errors.New("glTexImage2DMultisample")
	}
	gpTexImage3D = (C.GPTEXIMAGE3D)(getProcAddr("glTexImage3D"))
	if gpTexImage3D == nil {
		return errors.New("glTexImage3D")
	}
	gpTexImage3DMultisample = (C.GPTEXIMAGE3DMULTISAMPLE)(getProcAddr("glTexImage3DMultisample"))
	if gpTexImage3DMultisample == nil {
		return errors.New("glTexImage3DMultisample")
	}
	gpTexPageCommitmentARB = (C.GPTEXPAGECOMMITMENTARB)(getProcAddr("glTexPageCommitmentARB"))
	gpTexParameterIiv = (C.GPTEXPARAMETERIIV)(getProcAddr("glTexParameterIiv"))
	if gpTexParameterIiv == nil {
		return errors.New("glTexParameterIiv")
	}
	gpTexParameterIuiv = (C.GPTEXPARAMETERIUIV)(getProcAddr("glTexParameterIuiv"))
	if gpTexParameterIuiv == nil {
		return errors.New("glTexParameterIuiv")
	}
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
	gpTexStorage1D = (C.GPTEXSTORAGE1D)(getProcAddr("glTexStorage1D"))
	gpTexStorage2D = (C.GPTEXSTORAGE2D)(getProcAddr("glTexStorage2D"))
	gpTexStorage2DMultisample = (C.GPTEXSTORAGE2DMULTISAMPLE)(getProcAddr("glTexStorage2DMultisample"))
	gpTexStorage3D = (C.GPTEXSTORAGE3D)(getProcAddr("glTexStorage3D"))
	gpTexStorage3DMultisample = (C.GPTEXSTORAGE3DMULTISAMPLE)(getProcAddr("glTexStorage3DMultisample"))
	gpTexSubImage1D = (C.GPTEXSUBIMAGE1D)(getProcAddr("glTexSubImage1D"))
	if gpTexSubImage1D == nil {
		return errors.New("glTexSubImage1D")
	}
	gpTexSubImage2D = (C.GPTEXSUBIMAGE2D)(getProcAddr("glTexSubImage2D"))
	if gpTexSubImage2D == nil {
		return errors.New("glTexSubImage2D")
	}
	gpTexSubImage3D = (C.GPTEXSUBIMAGE3D)(getProcAddr("glTexSubImage3D"))
	if gpTexSubImage3D == nil {
		return errors.New("glTexSubImage3D")
	}
	gpTextureBarrier = (C.GPTEXTUREBARRIER)(getProcAddr("glTextureBarrier"))
	gpTextureBarrierNV = (C.GPTEXTUREBARRIERNV)(getProcAddr("glTextureBarrierNV"))
	gpTextureBuffer = (C.GPTEXTUREBUFFER)(getProcAddr("glTextureBuffer"))
	gpTextureBufferEXT = (C.GPTEXTUREBUFFEREXT)(getProcAddr("glTextureBufferEXT"))
	gpTextureBufferRange = (C.GPTEXTUREBUFFERRANGE)(getProcAddr("glTextureBufferRange"))
	gpTextureBufferRangeEXT = (C.GPTEXTUREBUFFERRANGEEXT)(getProcAddr("glTextureBufferRangeEXT"))
	gpTextureImage1DEXT = (C.GPTEXTUREIMAGE1DEXT)(getProcAddr("glTextureImage1DEXT"))
	gpTextureImage2DEXT = (C.GPTEXTUREIMAGE2DEXT)(getProcAddr("glTextureImage2DEXT"))
	gpTextureImage3DEXT = (C.GPTEXTUREIMAGE3DEXT)(getProcAddr("glTextureImage3DEXT"))
	gpTexturePageCommitmentEXT = (C.GPTEXTUREPAGECOMMITMENTEXT)(getProcAddr("glTexturePageCommitmentEXT"))
	gpTextureParameterIiv = (C.GPTEXTUREPARAMETERIIV)(getProcAddr("glTextureParameterIiv"))
	gpTextureParameterIivEXT = (C.GPTEXTUREPARAMETERIIVEXT)(getProcAddr("glTextureParameterIivEXT"))
	gpTextureParameterIuiv = (C.GPTEXTUREPARAMETERIUIV)(getProcAddr("glTextureParameterIuiv"))
	gpTextureParameterIuivEXT = (C.GPTEXTUREPARAMETERIUIVEXT)(getProcAddr("glTextureParameterIuivEXT"))
	gpTextureParameterf = (C.GPTEXTUREPARAMETERF)(getProcAddr("glTextureParameterf"))
	gpTextureParameterfEXT = (C.GPTEXTUREPARAMETERFEXT)(getProcAddr("glTextureParameterfEXT"))
	gpTextureParameterfv = (C.GPTEXTUREPARAMETERFV)(getProcAddr("glTextureParameterfv"))
	gpTextureParameterfvEXT = (C.GPTEXTUREPARAMETERFVEXT)(getProcAddr("glTextureParameterfvEXT"))
	gpTextureParameteri = (C.GPTEXTUREPARAMETERI)(getProcAddr("glTextureParameteri"))
	gpTextureParameteriEXT = (C.GPTEXTUREPARAMETERIEXT)(getProcAddr("glTextureParameteriEXT"))
	gpTextureParameteriv = (C.GPTEXTUREPARAMETERIV)(getProcAddr("glTextureParameteriv"))
	gpTextureParameterivEXT = (C.GPTEXTUREPARAMETERIVEXT)(getProcAddr("glTextureParameterivEXT"))
	gpTextureRenderbufferEXT = (C.GPTEXTURERENDERBUFFEREXT)(getProcAddr("glTextureRenderbufferEXT"))
	gpTextureStorage1D = (C.GPTEXTURESTORAGE1D)(getProcAddr("glTextureStorage1D"))
	gpTextureStorage1DEXT = (C.GPTEXTURESTORAGE1DEXT)(getProcAddr("glTextureStorage1DEXT"))
	gpTextureStorage2D = (C.GPTEXTURESTORAGE2D)(getProcAddr("glTextureStorage2D"))
	gpTextureStorage2DEXT = (C.GPTEXTURESTORAGE2DEXT)(getProcAddr("glTextureStorage2DEXT"))
	gpTextureStorage2DMultisample = (C.GPTEXTURESTORAGE2DMULTISAMPLE)(getProcAddr("glTextureStorage2DMultisample"))
	gpTextureStorage2DMultisampleEXT = (C.GPTEXTURESTORAGE2DMULTISAMPLEEXT)(getProcAddr("glTextureStorage2DMultisampleEXT"))
	gpTextureStorage3D = (C.GPTEXTURESTORAGE3D)(getProcAddr("glTextureStorage3D"))
	gpTextureStorage3DEXT = (C.GPTEXTURESTORAGE3DEXT)(getProcAddr("glTextureStorage3DEXT"))
	gpTextureStorage3DMultisample = (C.GPTEXTURESTORAGE3DMULTISAMPLE)(getProcAddr("glTextureStorage3DMultisample"))
	gpTextureStorage3DMultisampleEXT = (C.GPTEXTURESTORAGE3DMULTISAMPLEEXT)(getProcAddr("glTextureStorage3DMultisampleEXT"))
	gpTextureSubImage1D = (C.GPTEXTURESUBIMAGE1D)(getProcAddr("glTextureSubImage1D"))
	gpTextureSubImage1DEXT = (C.GPTEXTURESUBIMAGE1DEXT)(getProcAddr("glTextureSubImage1DEXT"))
	gpTextureSubImage2D = (C.GPTEXTURESUBIMAGE2D)(getProcAddr("glTextureSubImage2D"))
	gpTextureSubImage2DEXT = (C.GPTEXTURESUBIMAGE2DEXT)(getProcAddr("glTextureSubImage2DEXT"))
	gpTextureSubImage3D = (C.GPTEXTURESUBIMAGE3D)(getProcAddr("glTextureSubImage3D"))
	gpTextureSubImage3DEXT = (C.GPTEXTURESUBIMAGE3DEXT)(getProcAddr("glTextureSubImage3DEXT"))
	gpTextureView = (C.GPTEXTUREVIEW)(getProcAddr("glTextureView"))
	gpTransformFeedbackBufferBase = (C.GPTRANSFORMFEEDBACKBUFFERBASE)(getProcAddr("glTransformFeedbackBufferBase"))
	gpTransformFeedbackBufferRange = (C.GPTRANSFORMFEEDBACKBUFFERRANGE)(getProcAddr("glTransformFeedbackBufferRange"))
	gpTransformFeedbackVaryings = (C.GPTRANSFORMFEEDBACKVARYINGS)(getProcAddr("glTransformFeedbackVaryings"))
	if gpTransformFeedbackVaryings == nil {
		return errors.New("glTransformFeedbackVaryings")
	}
	gpTransformPathNV = (C.GPTRANSFORMPATHNV)(getProcAddr("glTransformPathNV"))
	gpUniform1d = (C.GPUNIFORM1D)(getProcAddr("glUniform1d"))
	if gpUniform1d == nil {
		return errors.New("glUniform1d")
	}
	gpUniform1dv = (C.GPUNIFORM1DV)(getProcAddr("glUniform1dv"))
	if gpUniform1dv == nil {
		return errors.New("glUniform1dv")
	}
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
	gpUniform1i64ARB = (C.GPUNIFORM1I64ARB)(getProcAddr("glUniform1i64ARB"))
	gpUniform1i64NV = (C.GPUNIFORM1I64NV)(getProcAddr("glUniform1i64NV"))
	gpUniform1i64vARB = (C.GPUNIFORM1I64VARB)(getProcAddr("glUniform1i64vARB"))
	gpUniform1i64vNV = (C.GPUNIFORM1I64VNV)(getProcAddr("glUniform1i64vNV"))
	gpUniform1iv = (C.GPUNIFORM1IV)(getProcAddr("glUniform1iv"))
	if gpUniform1iv == nil {
		return errors.New("glUniform1iv")
	}
	gpUniform1ui = (C.GPUNIFORM1UI)(getProcAddr("glUniform1ui"))
	if gpUniform1ui == nil {
		return errors.New("glUniform1ui")
	}
	gpUniform1ui64ARB = (C.GPUNIFORM1UI64ARB)(getProcAddr("glUniform1ui64ARB"))
	gpUniform1ui64NV = (C.GPUNIFORM1UI64NV)(getProcAddr("glUniform1ui64NV"))
	gpUniform1ui64vARB = (C.GPUNIFORM1UI64VARB)(getProcAddr("glUniform1ui64vARB"))
	gpUniform1ui64vNV = (C.GPUNIFORM1UI64VNV)(getProcAddr("glUniform1ui64vNV"))
	gpUniform1uiv = (C.GPUNIFORM1UIV)(getProcAddr("glUniform1uiv"))
	if gpUniform1uiv == nil {
		return errors.New("glUniform1uiv")
	}
	gpUniform2d = (C.GPUNIFORM2D)(getProcAddr("glUniform2d"))
	if gpUniform2d == nil {
		return errors.New("glUniform2d")
	}
	gpUniform2dv = (C.GPUNIFORM2DV)(getProcAddr("glUniform2dv"))
	if gpUniform2dv == nil {
		return errors.New("glUniform2dv")
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
	gpUniform2i64ARB = (C.GPUNIFORM2I64ARB)(getProcAddr("glUniform2i64ARB"))
	gpUniform2i64NV = (C.GPUNIFORM2I64NV)(getProcAddr("glUniform2i64NV"))
	gpUniform2i64vARB = (C.GPUNIFORM2I64VARB)(getProcAddr("glUniform2i64vARB"))
	gpUniform2i64vNV = (C.GPUNIFORM2I64VNV)(getProcAddr("glUniform2i64vNV"))
	gpUniform2iv = (C.GPUNIFORM2IV)(getProcAddr("glUniform2iv"))
	if gpUniform2iv == nil {
		return errors.New("glUniform2iv")
	}
	gpUniform2ui = (C.GPUNIFORM2UI)(getProcAddr("glUniform2ui"))
	if gpUniform2ui == nil {
		return errors.New("glUniform2ui")
	}
	gpUniform2ui64ARB = (C.GPUNIFORM2UI64ARB)(getProcAddr("glUniform2ui64ARB"))
	gpUniform2ui64NV = (C.GPUNIFORM2UI64NV)(getProcAddr("glUniform2ui64NV"))
	gpUniform2ui64vARB = (C.GPUNIFORM2UI64VARB)(getProcAddr("glUniform2ui64vARB"))
	gpUniform2ui64vNV = (C.GPUNIFORM2UI64VNV)(getProcAddr("glUniform2ui64vNV"))
	gpUniform2uiv = (C.GPUNIFORM2UIV)(getProcAddr("glUniform2uiv"))
	if gpUniform2uiv == nil {
		return errors.New("glUniform2uiv")
	}
	gpUniform3d = (C.GPUNIFORM3D)(getProcAddr("glUniform3d"))
	if gpUniform3d == nil {
		return errors.New("glUniform3d")
	}
	gpUniform3dv = (C.GPUNIFORM3DV)(getProcAddr("glUniform3dv"))
	if gpUniform3dv == nil {
		return errors.New("glUniform3dv")
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
	gpUniform3i64ARB = (C.GPUNIFORM3I64ARB)(getProcAddr("glUniform3i64ARB"))
	gpUniform3i64NV = (C.GPUNIFORM3I64NV)(getProcAddr("glUniform3i64NV"))
	gpUniform3i64vARB = (C.GPUNIFORM3I64VARB)(getProcAddr("glUniform3i64vARB"))
	gpUniform3i64vNV = (C.GPUNIFORM3I64VNV)(getProcAddr("glUniform3i64vNV"))
	gpUniform3iv = (C.GPUNIFORM3IV)(getProcAddr("glUniform3iv"))
	if gpUniform3iv == nil {
		return errors.New("glUniform3iv")
	}
	gpUniform3ui = (C.GPUNIFORM3UI)(getProcAddr("glUniform3ui"))
	if gpUniform3ui == nil {
		return errors.New("glUniform3ui")
	}
	gpUniform3ui64ARB = (C.GPUNIFORM3UI64ARB)(getProcAddr("glUniform3ui64ARB"))
	gpUniform3ui64NV = (C.GPUNIFORM3UI64NV)(getProcAddr("glUniform3ui64NV"))
	gpUniform3ui64vARB = (C.GPUNIFORM3UI64VARB)(getProcAddr("glUniform3ui64vARB"))
	gpUniform3ui64vNV = (C.GPUNIFORM3UI64VNV)(getProcAddr("glUniform3ui64vNV"))
	gpUniform3uiv = (C.GPUNIFORM3UIV)(getProcAddr("glUniform3uiv"))
	if gpUniform3uiv == nil {
		return errors.New("glUniform3uiv")
	}
	gpUniform4d = (C.GPUNIFORM4D)(getProcAddr("glUniform4d"))
	if gpUniform4d == nil {
		return errors.New("glUniform4d")
	}
	gpUniform4dv = (C.GPUNIFORM4DV)(getProcAddr("glUniform4dv"))
	if gpUniform4dv == nil {
		return errors.New("glUniform4dv")
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
	gpUniform4i64ARB = (C.GPUNIFORM4I64ARB)(getProcAddr("glUniform4i64ARB"))
	gpUniform4i64NV = (C.GPUNIFORM4I64NV)(getProcAddr("glUniform4i64NV"))
	gpUniform4i64vARB = (C.GPUNIFORM4I64VARB)(getProcAddr("glUniform4i64vARB"))
	gpUniform4i64vNV = (C.GPUNIFORM4I64VNV)(getProcAddr("glUniform4i64vNV"))
	gpUniform4iv = (C.GPUNIFORM4IV)(getProcAddr("glUniform4iv"))
	if gpUniform4iv == nil {
		return errors.New("glUniform4iv")
	}
	gpUniform4ui = (C.GPUNIFORM4UI)(getProcAddr("glUniform4ui"))
	if gpUniform4ui == nil {
		return errors.New("glUniform4ui")
	}
	gpUniform4ui64ARB = (C.GPUNIFORM4UI64ARB)(getProcAddr("glUniform4ui64ARB"))
	gpUniform4ui64NV = (C.GPUNIFORM4UI64NV)(getProcAddr("glUniform4ui64NV"))
	gpUniform4ui64vARB = (C.GPUNIFORM4UI64VARB)(getProcAddr("glUniform4ui64vARB"))
	gpUniform4ui64vNV = (C.GPUNIFORM4UI64VNV)(getProcAddr("glUniform4ui64vNV"))
	gpUniform4uiv = (C.GPUNIFORM4UIV)(getProcAddr("glUniform4uiv"))
	if gpUniform4uiv == nil {
		return errors.New("glUniform4uiv")
	}
	gpUniformBlockBinding = (C.GPUNIFORMBLOCKBINDING)(getProcAddr("glUniformBlockBinding"))
	if gpUniformBlockBinding == nil {
		return errors.New("glUniformBlockBinding")
	}
	gpUniformHandleui64ARB = (C.GPUNIFORMHANDLEUI64ARB)(getProcAddr("glUniformHandleui64ARB"))
	gpUniformHandleui64NV = (C.GPUNIFORMHANDLEUI64NV)(getProcAddr("glUniformHandleui64NV"))
	gpUniformHandleui64vARB = (C.GPUNIFORMHANDLEUI64VARB)(getProcAddr("glUniformHandleui64vARB"))
	gpUniformHandleui64vNV = (C.GPUNIFORMHANDLEUI64VNV)(getProcAddr("glUniformHandleui64vNV"))
	gpUniformMatrix2dv = (C.GPUNIFORMMATRIX2DV)(getProcAddr("glUniformMatrix2dv"))
	if gpUniformMatrix2dv == nil {
		return errors.New("glUniformMatrix2dv")
	}
	gpUniformMatrix2fv = (C.GPUNIFORMMATRIX2FV)(getProcAddr("glUniformMatrix2fv"))
	if gpUniformMatrix2fv == nil {
		return errors.New("glUniformMatrix2fv")
	}
	gpUniformMatrix2x3dv = (C.GPUNIFORMMATRIX2X3DV)(getProcAddr("glUniformMatrix2x3dv"))
	if gpUniformMatrix2x3dv == nil {
		return errors.New("glUniformMatrix2x3dv")
	}
	gpUniformMatrix2x3fv = (C.GPUNIFORMMATRIX2X3FV)(getProcAddr("glUniformMatrix2x3fv"))
	if gpUniformMatrix2x3fv == nil {
		return errors.New("glUniformMatrix2x3fv")
	}
	gpUniformMatrix2x4dv = (C.GPUNIFORMMATRIX2X4DV)(getProcAddr("glUniformMatrix2x4dv"))
	if gpUniformMatrix2x4dv == nil {
		return errors.New("glUniformMatrix2x4dv")
	}
	gpUniformMatrix2x4fv = (C.GPUNIFORMMATRIX2X4FV)(getProcAddr("glUniformMatrix2x4fv"))
	if gpUniformMatrix2x4fv == nil {
		return errors.New("glUniformMatrix2x4fv")
	}
	gpUniformMatrix3dv = (C.GPUNIFORMMATRIX3DV)(getProcAddr("glUniformMatrix3dv"))
	if gpUniformMatrix3dv == nil {
		return errors.New("glUniformMatrix3dv")
	}
	gpUniformMatrix3fv = (C.GPUNIFORMMATRIX3FV)(getProcAddr("glUniformMatrix3fv"))
	if gpUniformMatrix3fv == nil {
		return errors.New("glUniformMatrix3fv")
	}
	gpUniformMatrix3x2dv = (C.GPUNIFORMMATRIX3X2DV)(getProcAddr("glUniformMatrix3x2dv"))
	if gpUniformMatrix3x2dv == nil {
		return errors.New("glUniformMatrix3x2dv")
	}
	gpUniformMatrix3x2fv = (C.GPUNIFORMMATRIX3X2FV)(getProcAddr("glUniformMatrix3x2fv"))
	if gpUniformMatrix3x2fv == nil {
		return errors.New("glUniformMatrix3x2fv")
	}
	gpUniformMatrix3x4dv = (C.GPUNIFORMMATRIX3X4DV)(getProcAddr("glUniformMatrix3x4dv"))
	if gpUniformMatrix3x4dv == nil {
		return errors.New("glUniformMatrix3x4dv")
	}
	gpUniformMatrix3x4fv = (C.GPUNIFORMMATRIX3X4FV)(getProcAddr("glUniformMatrix3x4fv"))
	if gpUniformMatrix3x4fv == nil {
		return errors.New("glUniformMatrix3x4fv")
	}
	gpUniformMatrix4dv = (C.GPUNIFORMMATRIX4DV)(getProcAddr("glUniformMatrix4dv"))
	if gpUniformMatrix4dv == nil {
		return errors.New("glUniformMatrix4dv")
	}
	gpUniformMatrix4fv = (C.GPUNIFORMMATRIX4FV)(getProcAddr("glUniformMatrix4fv"))
	if gpUniformMatrix4fv == nil {
		return errors.New("glUniformMatrix4fv")
	}
	gpUniformMatrix4x2dv = (C.GPUNIFORMMATRIX4X2DV)(getProcAddr("glUniformMatrix4x2dv"))
	if gpUniformMatrix4x2dv == nil {
		return errors.New("glUniformMatrix4x2dv")
	}
	gpUniformMatrix4x2fv = (C.GPUNIFORMMATRIX4X2FV)(getProcAddr("glUniformMatrix4x2fv"))
	if gpUniformMatrix4x2fv == nil {
		return errors.New("glUniformMatrix4x2fv")
	}
	gpUniformMatrix4x3dv = (C.GPUNIFORMMATRIX4X3DV)(getProcAddr("glUniformMatrix4x3dv"))
	if gpUniformMatrix4x3dv == nil {
		return errors.New("glUniformMatrix4x3dv")
	}
	gpUniformMatrix4x3fv = (C.GPUNIFORMMATRIX4X3FV)(getProcAddr("glUniformMatrix4x3fv"))
	if gpUniformMatrix4x3fv == nil {
		return errors.New("glUniformMatrix4x3fv")
	}
	gpUniformSubroutinesuiv = (C.GPUNIFORMSUBROUTINESUIV)(getProcAddr("glUniformSubroutinesuiv"))
	if gpUniformSubroutinesuiv == nil {
		return errors.New("glUniformSubroutinesuiv")
	}
	gpUniformui64NV = (C.GPUNIFORMUI64NV)(getProcAddr("glUniformui64NV"))
	gpUniformui64vNV = (C.GPUNIFORMUI64VNV)(getProcAddr("glUniformui64vNV"))
	gpUnmapBuffer = (C.GPUNMAPBUFFER)(getProcAddr("glUnmapBuffer"))
	if gpUnmapBuffer == nil {
		return errors.New("glUnmapBuffer")
	}
	gpUnmapNamedBuffer = (C.GPUNMAPNAMEDBUFFER)(getProcAddr("glUnmapNamedBuffer"))
	gpUnmapNamedBufferEXT = (C.GPUNMAPNAMEDBUFFEREXT)(getProcAddr("glUnmapNamedBufferEXT"))
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
	gpVertexArrayAttribBinding = (C.GPVERTEXARRAYATTRIBBINDING)(getProcAddr("glVertexArrayAttribBinding"))
	gpVertexArrayAttribFormat = (C.GPVERTEXARRAYATTRIBFORMAT)(getProcAddr("glVertexArrayAttribFormat"))
	gpVertexArrayAttribIFormat = (C.GPVERTEXARRAYATTRIBIFORMAT)(getProcAddr("glVertexArrayAttribIFormat"))
	gpVertexArrayAttribLFormat = (C.GPVERTEXARRAYATTRIBLFORMAT)(getProcAddr("glVertexArrayAttribLFormat"))
	gpVertexArrayBindVertexBufferEXT = (C.GPVERTEXARRAYBINDVERTEXBUFFEREXT)(getProcAddr("glVertexArrayBindVertexBufferEXT"))
	gpVertexArrayBindingDivisor = (C.GPVERTEXARRAYBINDINGDIVISOR)(getProcAddr("glVertexArrayBindingDivisor"))
	gpVertexArrayColorOffsetEXT = (C.GPVERTEXARRAYCOLOROFFSETEXT)(getProcAddr("glVertexArrayColorOffsetEXT"))
	gpVertexArrayEdgeFlagOffsetEXT = (C.GPVERTEXARRAYEDGEFLAGOFFSETEXT)(getProcAddr("glVertexArrayEdgeFlagOffsetEXT"))
	gpVertexArrayElementBuffer = (C.GPVERTEXARRAYELEMENTBUFFER)(getProcAddr("glVertexArrayElementBuffer"))
	gpVertexArrayFogCoordOffsetEXT = (C.GPVERTEXARRAYFOGCOORDOFFSETEXT)(getProcAddr("glVertexArrayFogCoordOffsetEXT"))
	gpVertexArrayIndexOffsetEXT = (C.GPVERTEXARRAYINDEXOFFSETEXT)(getProcAddr("glVertexArrayIndexOffsetEXT"))
	gpVertexArrayMultiTexCoordOffsetEXT = (C.GPVERTEXARRAYMULTITEXCOORDOFFSETEXT)(getProcAddr("glVertexArrayMultiTexCoordOffsetEXT"))
	gpVertexArrayNormalOffsetEXT = (C.GPVERTEXARRAYNORMALOFFSETEXT)(getProcAddr("glVertexArrayNormalOffsetEXT"))
	gpVertexArraySecondaryColorOffsetEXT = (C.GPVERTEXARRAYSECONDARYCOLOROFFSETEXT)(getProcAddr("glVertexArraySecondaryColorOffsetEXT"))
	gpVertexArrayTexCoordOffsetEXT = (C.GPVERTEXARRAYTEXCOORDOFFSETEXT)(getProcAddr("glVertexArrayTexCoordOffsetEXT"))
	gpVertexArrayVertexAttribBindingEXT = (C.GPVERTEXARRAYVERTEXATTRIBBINDINGEXT)(getProcAddr("glVertexArrayVertexAttribBindingEXT"))
	gpVertexArrayVertexAttribDivisorEXT = (C.GPVERTEXARRAYVERTEXATTRIBDIVISOREXT)(getProcAddr("glVertexArrayVertexAttribDivisorEXT"))
	gpVertexArrayVertexAttribFormatEXT = (C.GPVERTEXARRAYVERTEXATTRIBFORMATEXT)(getProcAddr("glVertexArrayVertexAttribFormatEXT"))
	gpVertexArrayVertexAttribIFormatEXT = (C.GPVERTEXARRAYVERTEXATTRIBIFORMATEXT)(getProcAddr("glVertexArrayVertexAttribIFormatEXT"))
	gpVertexArrayVertexAttribIOffsetEXT = (C.GPVERTEXARRAYVERTEXATTRIBIOFFSETEXT)(getProcAddr("glVertexArrayVertexAttribIOffsetEXT"))
	gpVertexArrayVertexAttribLFormatEXT = (C.GPVERTEXARRAYVERTEXATTRIBLFORMATEXT)(getProcAddr("glVertexArrayVertexAttribLFormatEXT"))
	gpVertexArrayVertexAttribLOffsetEXT = (C.GPVERTEXARRAYVERTEXATTRIBLOFFSETEXT)(getProcAddr("glVertexArrayVertexAttribLOffsetEXT"))
	gpVertexArrayVertexAttribOffsetEXT = (C.GPVERTEXARRAYVERTEXATTRIBOFFSETEXT)(getProcAddr("glVertexArrayVertexAttribOffsetEXT"))
	gpVertexArrayVertexBindingDivisorEXT = (C.GPVERTEXARRAYVERTEXBINDINGDIVISOREXT)(getProcAddr("glVertexArrayVertexBindingDivisorEXT"))
	gpVertexArrayVertexBuffer = (C.GPVERTEXARRAYVERTEXBUFFER)(getProcAddr("glVertexArrayVertexBuffer"))
	gpVertexArrayVertexBuffers = (C.GPVERTEXARRAYVERTEXBUFFERS)(getProcAddr("glVertexArrayVertexBuffers"))
	gpVertexArrayVertexOffsetEXT = (C.GPVERTEXARRAYVERTEXOFFSETEXT)(getProcAddr("glVertexArrayVertexOffsetEXT"))
	gpVertexAttrib1d = (C.GPVERTEXATTRIB1D)(getProcAddr("glVertexAttrib1d"))
	if gpVertexAttrib1d == nil {
		return errors.New("glVertexAttrib1d")
	}
	gpVertexAttrib1dv = (C.GPVERTEXATTRIB1DV)(getProcAddr("glVertexAttrib1dv"))
	if gpVertexAttrib1dv == nil {
		return errors.New("glVertexAttrib1dv")
	}
	gpVertexAttrib1f = (C.GPVERTEXATTRIB1F)(getProcAddr("glVertexAttrib1f"))
	if gpVertexAttrib1f == nil {
		return errors.New("glVertexAttrib1f")
	}
	gpVertexAttrib1fv = (C.GPVERTEXATTRIB1FV)(getProcAddr("glVertexAttrib1fv"))
	if gpVertexAttrib1fv == nil {
		return errors.New("glVertexAttrib1fv")
	}
	gpVertexAttrib1s = (C.GPVERTEXATTRIB1S)(getProcAddr("glVertexAttrib1s"))
	if gpVertexAttrib1s == nil {
		return errors.New("glVertexAttrib1s")
	}
	gpVertexAttrib1sv = (C.GPVERTEXATTRIB1SV)(getProcAddr("glVertexAttrib1sv"))
	if gpVertexAttrib1sv == nil {
		return errors.New("glVertexAttrib1sv")
	}
	gpVertexAttrib2d = (C.GPVERTEXATTRIB2D)(getProcAddr("glVertexAttrib2d"))
	if gpVertexAttrib2d == nil {
		return errors.New("glVertexAttrib2d")
	}
	gpVertexAttrib2dv = (C.GPVERTEXATTRIB2DV)(getProcAddr("glVertexAttrib2dv"))
	if gpVertexAttrib2dv == nil {
		return errors.New("glVertexAttrib2dv")
	}
	gpVertexAttrib2f = (C.GPVERTEXATTRIB2F)(getProcAddr("glVertexAttrib2f"))
	if gpVertexAttrib2f == nil {
		return errors.New("glVertexAttrib2f")
	}
	gpVertexAttrib2fv = (C.GPVERTEXATTRIB2FV)(getProcAddr("glVertexAttrib2fv"))
	if gpVertexAttrib2fv == nil {
		return errors.New("glVertexAttrib2fv")
	}
	gpVertexAttrib2s = (C.GPVERTEXATTRIB2S)(getProcAddr("glVertexAttrib2s"))
	if gpVertexAttrib2s == nil {
		return errors.New("glVertexAttrib2s")
	}
	gpVertexAttrib2sv = (C.GPVERTEXATTRIB2SV)(getProcAddr("glVertexAttrib2sv"))
	if gpVertexAttrib2sv == nil {
		return errors.New("glVertexAttrib2sv")
	}
	gpVertexAttrib3d = (C.GPVERTEXATTRIB3D)(getProcAddr("glVertexAttrib3d"))
	if gpVertexAttrib3d == nil {
		return errors.New("glVertexAttrib3d")
	}
	gpVertexAttrib3dv = (C.GPVERTEXATTRIB3DV)(getProcAddr("glVertexAttrib3dv"))
	if gpVertexAttrib3dv == nil {
		return errors.New("glVertexAttrib3dv")
	}
	gpVertexAttrib3f = (C.GPVERTEXATTRIB3F)(getProcAddr("glVertexAttrib3f"))
	if gpVertexAttrib3f == nil {
		return errors.New("glVertexAttrib3f")
	}
	gpVertexAttrib3fv = (C.GPVERTEXATTRIB3FV)(getProcAddr("glVertexAttrib3fv"))
	if gpVertexAttrib3fv == nil {
		return errors.New("glVertexAttrib3fv")
	}
	gpVertexAttrib3s = (C.GPVERTEXATTRIB3S)(getProcAddr("glVertexAttrib3s"))
	if gpVertexAttrib3s == nil {
		return errors.New("glVertexAttrib3s")
	}
	gpVertexAttrib3sv = (C.GPVERTEXATTRIB3SV)(getProcAddr("glVertexAttrib3sv"))
	if gpVertexAttrib3sv == nil {
		return errors.New("glVertexAttrib3sv")
	}
	gpVertexAttrib4Nbv = (C.GPVERTEXATTRIB4NBV)(getProcAddr("glVertexAttrib4Nbv"))
	if gpVertexAttrib4Nbv == nil {
		return errors.New("glVertexAttrib4Nbv")
	}
	gpVertexAttrib4Niv = (C.GPVERTEXATTRIB4NIV)(getProcAddr("glVertexAttrib4Niv"))
	if gpVertexAttrib4Niv == nil {
		return errors.New("glVertexAttrib4Niv")
	}
	gpVertexAttrib4Nsv = (C.GPVERTEXATTRIB4NSV)(getProcAddr("glVertexAttrib4Nsv"))
	if gpVertexAttrib4Nsv == nil {
		return errors.New("glVertexAttrib4Nsv")
	}
	gpVertexAttrib4Nub = (C.GPVERTEXATTRIB4NUB)(getProcAddr("glVertexAttrib4Nub"))
	if gpVertexAttrib4Nub == nil {
		return errors.New("glVertexAttrib4Nub")
	}
	gpVertexAttrib4Nubv = (C.GPVERTEXATTRIB4NUBV)(getProcAddr("glVertexAttrib4Nubv"))
	if gpVertexAttrib4Nubv == nil {
		return errors.New("glVertexAttrib4Nubv")
	}
	gpVertexAttrib4Nuiv = (C.GPVERTEXATTRIB4NUIV)(getProcAddr("glVertexAttrib4Nuiv"))
	if gpVertexAttrib4Nuiv == nil {
		return errors.New("glVertexAttrib4Nuiv")
	}
	gpVertexAttrib4Nusv = (C.GPVERTEXATTRIB4NUSV)(getProcAddr("glVertexAttrib4Nusv"))
	if gpVertexAttrib4Nusv == nil {
		return errors.New("glVertexAttrib4Nusv")
	}
	gpVertexAttrib4bv = (C.GPVERTEXATTRIB4BV)(getProcAddr("glVertexAttrib4bv"))
	if gpVertexAttrib4bv == nil {
		return errors.New("glVertexAttrib4bv")
	}
	gpVertexAttrib4d = (C.GPVERTEXATTRIB4D)(getProcAddr("glVertexAttrib4d"))
	if gpVertexAttrib4d == nil {
		return errors.New("glVertexAttrib4d")
	}
	gpVertexAttrib4dv = (C.GPVERTEXATTRIB4DV)(getProcAddr("glVertexAttrib4dv"))
	if gpVertexAttrib4dv == nil {
		return errors.New("glVertexAttrib4dv")
	}
	gpVertexAttrib4f = (C.GPVERTEXATTRIB4F)(getProcAddr("glVertexAttrib4f"))
	if gpVertexAttrib4f == nil {
		return errors.New("glVertexAttrib4f")
	}
	gpVertexAttrib4fv = (C.GPVERTEXATTRIB4FV)(getProcAddr("glVertexAttrib4fv"))
	if gpVertexAttrib4fv == nil {
		return errors.New("glVertexAttrib4fv")
	}
	gpVertexAttrib4iv = (C.GPVERTEXATTRIB4IV)(getProcAddr("glVertexAttrib4iv"))
	if gpVertexAttrib4iv == nil {
		return errors.New("glVertexAttrib4iv")
	}
	gpVertexAttrib4s = (C.GPVERTEXATTRIB4S)(getProcAddr("glVertexAttrib4s"))
	if gpVertexAttrib4s == nil {
		return errors.New("glVertexAttrib4s")
	}
	gpVertexAttrib4sv = (C.GPVERTEXATTRIB4SV)(getProcAddr("glVertexAttrib4sv"))
	if gpVertexAttrib4sv == nil {
		return errors.New("glVertexAttrib4sv")
	}
	gpVertexAttrib4ubv = (C.GPVERTEXATTRIB4UBV)(getProcAddr("glVertexAttrib4ubv"))
	if gpVertexAttrib4ubv == nil {
		return errors.New("glVertexAttrib4ubv")
	}
	gpVertexAttrib4uiv = (C.GPVERTEXATTRIB4UIV)(getProcAddr("glVertexAttrib4uiv"))
	if gpVertexAttrib4uiv == nil {
		return errors.New("glVertexAttrib4uiv")
	}
	gpVertexAttrib4usv = (C.GPVERTEXATTRIB4USV)(getProcAddr("glVertexAttrib4usv"))
	if gpVertexAttrib4usv == nil {
		return errors.New("glVertexAttrib4usv")
	}
	gpVertexAttribBinding = (C.GPVERTEXATTRIBBINDING)(getProcAddr("glVertexAttribBinding"))
	gpVertexAttribDivisor = (C.GPVERTEXATTRIBDIVISOR)(getProcAddr("glVertexAttribDivisor"))
	if gpVertexAttribDivisor == nil {
		return errors.New("glVertexAttribDivisor")
	}
	gpVertexAttribDivisorARB = (C.GPVERTEXATTRIBDIVISORARB)(getProcAddr("glVertexAttribDivisorARB"))
	gpVertexAttribFormat = (C.GPVERTEXATTRIBFORMAT)(getProcAddr("glVertexAttribFormat"))
	gpVertexAttribFormatNV = (C.GPVERTEXATTRIBFORMATNV)(getProcAddr("glVertexAttribFormatNV"))
	gpVertexAttribI1i = (C.GPVERTEXATTRIBI1I)(getProcAddr("glVertexAttribI1i"))
	if gpVertexAttribI1i == nil {
		return errors.New("glVertexAttribI1i")
	}
	gpVertexAttribI1iv = (C.GPVERTEXATTRIBI1IV)(getProcAddr("glVertexAttribI1iv"))
	if gpVertexAttribI1iv == nil {
		return errors.New("glVertexAttribI1iv")
	}
	gpVertexAttribI1ui = (C.GPVERTEXATTRIBI1UI)(getProcAddr("glVertexAttribI1ui"))
	if gpVertexAttribI1ui == nil {
		return errors.New("glVertexAttribI1ui")
	}
	gpVertexAttribI1uiv = (C.GPVERTEXATTRIBI1UIV)(getProcAddr("glVertexAttribI1uiv"))
	if gpVertexAttribI1uiv == nil {
		return errors.New("glVertexAttribI1uiv")
	}
	gpVertexAttribI2i = (C.GPVERTEXATTRIBI2I)(getProcAddr("glVertexAttribI2i"))
	if gpVertexAttribI2i == nil {
		return errors.New("glVertexAttribI2i")
	}
	gpVertexAttribI2iv = (C.GPVERTEXATTRIBI2IV)(getProcAddr("glVertexAttribI2iv"))
	if gpVertexAttribI2iv == nil {
		return errors.New("glVertexAttribI2iv")
	}
	gpVertexAttribI2ui = (C.GPVERTEXATTRIBI2UI)(getProcAddr("glVertexAttribI2ui"))
	if gpVertexAttribI2ui == nil {
		return errors.New("glVertexAttribI2ui")
	}
	gpVertexAttribI2uiv = (C.GPVERTEXATTRIBI2UIV)(getProcAddr("glVertexAttribI2uiv"))
	if gpVertexAttribI2uiv == nil {
		return errors.New("glVertexAttribI2uiv")
	}
	gpVertexAttribI3i = (C.GPVERTEXATTRIBI3I)(getProcAddr("glVertexAttribI3i"))
	if gpVertexAttribI3i == nil {
		return errors.New("glVertexAttribI3i")
	}
	gpVertexAttribI3iv = (C.GPVERTEXATTRIBI3IV)(getProcAddr("glVertexAttribI3iv"))
	if gpVertexAttribI3iv == nil {
		return errors.New("glVertexAttribI3iv")
	}
	gpVertexAttribI3ui = (C.GPVERTEXATTRIBI3UI)(getProcAddr("glVertexAttribI3ui"))
	if gpVertexAttribI3ui == nil {
		return errors.New("glVertexAttribI3ui")
	}
	gpVertexAttribI3uiv = (C.GPVERTEXATTRIBI3UIV)(getProcAddr("glVertexAttribI3uiv"))
	if gpVertexAttribI3uiv == nil {
		return errors.New("glVertexAttribI3uiv")
	}
	gpVertexAttribI4bv = (C.GPVERTEXATTRIBI4BV)(getProcAddr("glVertexAttribI4bv"))
	if gpVertexAttribI4bv == nil {
		return errors.New("glVertexAttribI4bv")
	}
	gpVertexAttribI4i = (C.GPVERTEXATTRIBI4I)(getProcAddr("glVertexAttribI4i"))
	if gpVertexAttribI4i == nil {
		return errors.New("glVertexAttribI4i")
	}
	gpVertexAttribI4iv = (C.GPVERTEXATTRIBI4IV)(getProcAddr("glVertexAttribI4iv"))
	if gpVertexAttribI4iv == nil {
		return errors.New("glVertexAttribI4iv")
	}
	gpVertexAttribI4sv = (C.GPVERTEXATTRIBI4SV)(getProcAddr("glVertexAttribI4sv"))
	if gpVertexAttribI4sv == nil {
		return errors.New("glVertexAttribI4sv")
	}
	gpVertexAttribI4ubv = (C.GPVERTEXATTRIBI4UBV)(getProcAddr("glVertexAttribI4ubv"))
	if gpVertexAttribI4ubv == nil {
		return errors.New("glVertexAttribI4ubv")
	}
	gpVertexAttribI4ui = (C.GPVERTEXATTRIBI4UI)(getProcAddr("glVertexAttribI4ui"))
	if gpVertexAttribI4ui == nil {
		return errors.New("glVertexAttribI4ui")
	}
	gpVertexAttribI4uiv = (C.GPVERTEXATTRIBI4UIV)(getProcAddr("glVertexAttribI4uiv"))
	if gpVertexAttribI4uiv == nil {
		return errors.New("glVertexAttribI4uiv")
	}
	gpVertexAttribI4usv = (C.GPVERTEXATTRIBI4USV)(getProcAddr("glVertexAttribI4usv"))
	if gpVertexAttribI4usv == nil {
		return errors.New("glVertexAttribI4usv")
	}
	gpVertexAttribIFormat = (C.GPVERTEXATTRIBIFORMAT)(getProcAddr("glVertexAttribIFormat"))
	gpVertexAttribIFormatNV = (C.GPVERTEXATTRIBIFORMATNV)(getProcAddr("glVertexAttribIFormatNV"))
	gpVertexAttribIPointer = (C.GPVERTEXATTRIBIPOINTER)(getProcAddr("glVertexAttribIPointer"))
	if gpVertexAttribIPointer == nil {
		return errors.New("glVertexAttribIPointer")
	}
	gpVertexAttribL1d = (C.GPVERTEXATTRIBL1D)(getProcAddr("glVertexAttribL1d"))
	if gpVertexAttribL1d == nil {
		return errors.New("glVertexAttribL1d")
	}
	gpVertexAttribL1dv = (C.GPVERTEXATTRIBL1DV)(getProcAddr("glVertexAttribL1dv"))
	if gpVertexAttribL1dv == nil {
		return errors.New("glVertexAttribL1dv")
	}
	gpVertexAttribL1i64NV = (C.GPVERTEXATTRIBL1I64NV)(getProcAddr("glVertexAttribL1i64NV"))
	gpVertexAttribL1i64vNV = (C.GPVERTEXATTRIBL1I64VNV)(getProcAddr("glVertexAttribL1i64vNV"))
	gpVertexAttribL1ui64ARB = (C.GPVERTEXATTRIBL1UI64ARB)(getProcAddr("glVertexAttribL1ui64ARB"))
	gpVertexAttribL1ui64NV = (C.GPVERTEXATTRIBL1UI64NV)(getProcAddr("glVertexAttribL1ui64NV"))
	gpVertexAttribL1ui64vARB = (C.GPVERTEXATTRIBL1UI64VARB)(getProcAddr("glVertexAttribL1ui64vARB"))
	gpVertexAttribL1ui64vNV = (C.GPVERTEXATTRIBL1UI64VNV)(getProcAddr("glVertexAttribL1ui64vNV"))
	gpVertexAttribL2d = (C.GPVERTEXATTRIBL2D)(getProcAddr("glVertexAttribL2d"))
	if gpVertexAttribL2d == nil {
		return errors.New("glVertexAttribL2d")
	}
	gpVertexAttribL2dv = (C.GPVERTEXATTRIBL2DV)(getProcAddr("glVertexAttribL2dv"))
	if gpVertexAttribL2dv == nil {
		return errors.New("glVertexAttribL2dv")
	}
	gpVertexAttribL2i64NV = (C.GPVERTEXATTRIBL2I64NV)(getProcAddr("glVertexAttribL2i64NV"))
	gpVertexAttribL2i64vNV = (C.GPVERTEXATTRIBL2I64VNV)(getProcAddr("glVertexAttribL2i64vNV"))
	gpVertexAttribL2ui64NV = (C.GPVERTEXATTRIBL2UI64NV)(getProcAddr("glVertexAttribL2ui64NV"))
	gpVertexAttribL2ui64vNV = (C.GPVERTEXATTRIBL2UI64VNV)(getProcAddr("glVertexAttribL2ui64vNV"))
	gpVertexAttribL3d = (C.GPVERTEXATTRIBL3D)(getProcAddr("glVertexAttribL3d"))
	if gpVertexAttribL3d == nil {
		return errors.New("glVertexAttribL3d")
	}
	gpVertexAttribL3dv = (C.GPVERTEXATTRIBL3DV)(getProcAddr("glVertexAttribL3dv"))
	if gpVertexAttribL3dv == nil {
		return errors.New("glVertexAttribL3dv")
	}
	gpVertexAttribL3i64NV = (C.GPVERTEXATTRIBL3I64NV)(getProcAddr("glVertexAttribL3i64NV"))
	gpVertexAttribL3i64vNV = (C.GPVERTEXATTRIBL3I64VNV)(getProcAddr("glVertexAttribL3i64vNV"))
	gpVertexAttribL3ui64NV = (C.GPVERTEXATTRIBL3UI64NV)(getProcAddr("glVertexAttribL3ui64NV"))
	gpVertexAttribL3ui64vNV = (C.GPVERTEXATTRIBL3UI64VNV)(getProcAddr("glVertexAttribL3ui64vNV"))
	gpVertexAttribL4d = (C.GPVERTEXATTRIBL4D)(getProcAddr("glVertexAttribL4d"))
	if gpVertexAttribL4d == nil {
		return errors.New("glVertexAttribL4d")
	}
	gpVertexAttribL4dv = (C.GPVERTEXATTRIBL4DV)(getProcAddr("glVertexAttribL4dv"))
	if gpVertexAttribL4dv == nil {
		return errors.New("glVertexAttribL4dv")
	}
	gpVertexAttribL4i64NV = (C.GPVERTEXATTRIBL4I64NV)(getProcAddr("glVertexAttribL4i64NV"))
	gpVertexAttribL4i64vNV = (C.GPVERTEXATTRIBL4I64VNV)(getProcAddr("glVertexAttribL4i64vNV"))
	gpVertexAttribL4ui64NV = (C.GPVERTEXATTRIBL4UI64NV)(getProcAddr("glVertexAttribL4ui64NV"))
	gpVertexAttribL4ui64vNV = (C.GPVERTEXATTRIBL4UI64VNV)(getProcAddr("glVertexAttribL4ui64vNV"))
	gpVertexAttribLFormat = (C.GPVERTEXATTRIBLFORMAT)(getProcAddr("glVertexAttribLFormat"))
	gpVertexAttribLFormatNV = (C.GPVERTEXATTRIBLFORMATNV)(getProcAddr("glVertexAttribLFormatNV"))
	gpVertexAttribLPointer = (C.GPVERTEXATTRIBLPOINTER)(getProcAddr("glVertexAttribLPointer"))
	if gpVertexAttribLPointer == nil {
		return errors.New("glVertexAttribLPointer")
	}
	gpVertexAttribP1ui = (C.GPVERTEXATTRIBP1UI)(getProcAddr("glVertexAttribP1ui"))
	if gpVertexAttribP1ui == nil {
		return errors.New("glVertexAttribP1ui")
	}
	gpVertexAttribP1uiv = (C.GPVERTEXATTRIBP1UIV)(getProcAddr("glVertexAttribP1uiv"))
	if gpVertexAttribP1uiv == nil {
		return errors.New("glVertexAttribP1uiv")
	}
	gpVertexAttribP2ui = (C.GPVERTEXATTRIBP2UI)(getProcAddr("glVertexAttribP2ui"))
	if gpVertexAttribP2ui == nil {
		return errors.New("glVertexAttribP2ui")
	}
	gpVertexAttribP2uiv = (C.GPVERTEXATTRIBP2UIV)(getProcAddr("glVertexAttribP2uiv"))
	if gpVertexAttribP2uiv == nil {
		return errors.New("glVertexAttribP2uiv")
	}
	gpVertexAttribP3ui = (C.GPVERTEXATTRIBP3UI)(getProcAddr("glVertexAttribP3ui"))
	if gpVertexAttribP3ui == nil {
		return errors.New("glVertexAttribP3ui")
	}
	gpVertexAttribP3uiv = (C.GPVERTEXATTRIBP3UIV)(getProcAddr("glVertexAttribP3uiv"))
	if gpVertexAttribP3uiv == nil {
		return errors.New("glVertexAttribP3uiv")
	}
	gpVertexAttribP4ui = (C.GPVERTEXATTRIBP4UI)(getProcAddr("glVertexAttribP4ui"))
	if gpVertexAttribP4ui == nil {
		return errors.New("glVertexAttribP4ui")
	}
	gpVertexAttribP4uiv = (C.GPVERTEXATTRIBP4UIV)(getProcAddr("glVertexAttribP4uiv"))
	if gpVertexAttribP4uiv == nil {
		return errors.New("glVertexAttribP4uiv")
	}
	gpVertexAttribPointer = (C.GPVERTEXATTRIBPOINTER)(getProcAddr("glVertexAttribPointer"))
	if gpVertexAttribPointer == nil {
		return errors.New("glVertexAttribPointer")
	}
	gpVertexBindingDivisor = (C.GPVERTEXBINDINGDIVISOR)(getProcAddr("glVertexBindingDivisor"))
	gpVertexFormatNV = (C.GPVERTEXFORMATNV)(getProcAddr("glVertexFormatNV"))
	gpViewport = (C.GPVIEWPORT)(getProcAddr("glViewport"))
	if gpViewport == nil {
		return errors.New("glViewport")
	}
	gpViewportArrayv = (C.GPVIEWPORTARRAYV)(getProcAddr("glViewportArrayv"))
	if gpViewportArrayv == nil {
		return errors.New("glViewportArrayv")
	}
	gpViewportIndexedf = (C.GPVIEWPORTINDEXEDF)(getProcAddr("glViewportIndexedf"))
	if gpViewportIndexedf == nil {
		return errors.New("glViewportIndexedf")
	}
	gpViewportIndexedfv = (C.GPVIEWPORTINDEXEDFV)(getProcAddr("glViewportIndexedfv"))
	if gpViewportIndexedfv == nil {
		return errors.New("glViewportIndexedfv")
	}
	gpViewportPositionWScaleNV = (C.GPVIEWPORTPOSITIONWSCALENV)(getProcAddr("glViewportPositionWScaleNV"))
	gpViewportSwizzleNV = (C.GPVIEWPORTSWIZZLENV)(getProcAddr("glViewportSwizzleNV"))
	gpWaitSync = (C.GPWAITSYNC)(getProcAddr("glWaitSync"))
	if gpWaitSync == nil {
		return errors.New("glWaitSync")
	}
	gpWaitVkSemaphoreNV = (C.GPWAITVKSEMAPHORENV)(getProcAddr("glWaitVkSemaphoreNV"))
	gpWeightPathsNV = (C.GPWEIGHTPATHSNV)(getProcAddr("glWeightPathsNV"))
	gpWindowRectanglesEXT = (C.GPWINDOWRECTANGLESEXT)(getProcAddr("glWindowRectanglesEXT"))
	return nil
}
