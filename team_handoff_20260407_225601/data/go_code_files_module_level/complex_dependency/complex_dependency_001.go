func InitWithProcAddrFunc(getProcAddr func(name string) unsafe.Pointer) error {
	gpAccum = (C.GPACCUM)(getProcAddr("glAccum"))
	if gpAccum == nil {
		return errors.New("glAccum")
	}
	gpAccumxOES = (C.GPACCUMXOES)(getProcAddr("glAccumxOES"))
	gpAcquireKeyedMutexWin32EXT = (C.GPACQUIREKEYEDMUTEXWIN32EXT)(getProcAddr("glAcquireKeyedMutexWin32EXT"))
	gpActiveProgramEXT = (C.GPACTIVEPROGRAMEXT)(getProcAddr("glActiveProgramEXT"))
	gpActiveShaderProgram = (C.GPACTIVESHADERPROGRAM)(getProcAddr("glActiveShaderProgram"))
	if gpActiveShaderProgram == nil {
		return errors.New("glActiveShaderProgram")
	}
	gpActiveShaderProgramEXT = (C.GPACTIVESHADERPROGRAMEXT)(getProcAddr("glActiveShaderProgramEXT"))
	gpActiveStencilFaceEXT = (C.GPACTIVESTENCILFACEEXT)(getProcAddr("glActiveStencilFaceEXT"))
	gpActiveTexture = (C.GPACTIVETEXTURE)(getProcAddr("glActiveTexture"))
	if gpActiveTexture == nil {
		return errors.New("glActiveTexture")
	}
	gpActiveTextureARB = (C.GPACTIVETEXTUREARB)(getProcAddr("glActiveTextureARB"))
	gpActiveVaryingNV = (C.GPACTIVEVARYINGNV)(getProcAddr("glActiveVaryingNV"))
	gpAlphaFragmentOp1ATI = (C.GPALPHAFRAGMENTOP1ATI)(getProcAddr("glAlphaFragmentOp1ATI"))
	gpAlphaFragmentOp2ATI = (C.GPALPHAFRAGMENTOP2ATI)(getProcAddr("glAlphaFragmentOp2ATI"))
	gpAlphaFragmentOp3ATI = (C.GPALPHAFRAGMENTOP3ATI)(getProcAddr("glAlphaFragmentOp3ATI"))
	gpAlphaFunc = (C.GPALPHAFUNC)(getProcAddr("glAlphaFunc"))
	if gpAlphaFunc == nil {
		return errors.New("glAlphaFunc")
	}
	gpAlphaFuncxOES = (C.GPALPHAFUNCXOES)(getProcAddr("glAlphaFuncxOES"))
	gpAlphaToCoverageDitherControlNV = (C.GPALPHATOCOVERAGEDITHERCONTROLNV)(getProcAddr("glAlphaToCoverageDitherControlNV"))
	gpApplyFramebufferAttachmentCMAAINTEL = (C.GPAPPLYFRAMEBUFFERATTACHMENTCMAAINTEL)(getProcAddr("glApplyFramebufferAttachmentCMAAINTEL"))
	gpApplyTextureEXT = (C.GPAPPLYTEXTUREEXT)(getProcAddr("glApplyTextureEXT"))
	gpAreProgramsResidentNV = (C.GPAREPROGRAMSRESIDENTNV)(getProcAddr("glAreProgramsResidentNV"))
	gpAreTexturesResident = (C.GPARETEXTURESRESIDENT)(getProcAddr("glAreTexturesResident"))
	if gpAreTexturesResident == nil {
		return errors.New("glAreTexturesResident")
	}
	gpAreTexturesResidentEXT = (C.GPARETEXTURESRESIDENTEXT)(getProcAddr("glAreTexturesResidentEXT"))
	gpArrayElement = (C.GPARRAYELEMENT)(getProcAddr("glArrayElement"))
	if gpArrayElement == nil {
		return errors.New("glArrayElement")
	}
	gpArrayElementEXT = (C.GPARRAYELEMENTEXT)(getProcAddr("glArrayElementEXT"))
	gpArrayObjectATI = (C.GPARRAYOBJECTATI)(getProcAddr("glArrayObjectATI"))
	gpAsyncMarkerSGIX = (C.GPASYNCMARKERSGIX)(getProcAddr("glAsyncMarkerSGIX"))
	gpAttachObjectARB = (C.GPATTACHOBJECTARB)(getProcAddr("glAttachObjectARB"))
	gpAttachShader = (C.GPATTACHSHADER)(getProcAddr("glAttachShader"))
	if gpAttachShader == nil {
		return errors.New("glAttachShader")
	}
	gpBegin = (C.GPBEGIN)(getProcAddr("glBegin"))
	if gpBegin == nil {
		return errors.New("glBegin")
	}
	gpBeginConditionalRender = (C.GPBEGINCONDITIONALRENDER)(getProcAddr("glBeginConditionalRender"))
	if gpBeginConditionalRender == nil {
		return errors.New("glBeginConditionalRender")
	}
	gpBeginConditionalRenderNV = (C.GPBEGINCONDITIONALRENDERNV)(getProcAddr("glBeginConditionalRenderNV"))
	gpBeginConditionalRenderNVX = (C.GPBEGINCONDITIONALRENDERNVX)(getProcAddr("glBeginConditionalRenderNVX"))
	gpBeginFragmentShaderATI = (C.GPBEGINFRAGMENTSHADERATI)(getProcAddr("glBeginFragmentShaderATI"))
	gpBeginOcclusionQueryNV = (C.GPBEGINOCCLUSIONQUERYNV)(getProcAddr("glBeginOcclusionQueryNV"))
	gpBeginPerfMonitorAMD = (C.GPBEGINPERFMONITORAMD)(getProcAddr("glBeginPerfMonitorAMD"))
	gpBeginPerfQueryINTEL = (C.GPBEGINPERFQUERYINTEL)(getProcAddr("glBeginPerfQueryINTEL"))
	gpBeginQuery = (C.GPBEGINQUERY)(getProcAddr("glBeginQuery"))
	if gpBeginQuery == nil {
		return errors.New("glBeginQuery")
	}
	gpBeginQueryARB = (C.GPBEGINQUERYARB)(getProcAddr("glBeginQueryARB"))
	gpBeginQueryIndexed = (C.GPBEGINQUERYINDEXED)(getProcAddr("glBeginQueryIndexed"))
	if gpBeginQueryIndexed == nil {
		return errors.New("glBeginQueryIndexed")
	}
	gpBeginTransformFeedback = (C.GPBEGINTRANSFORMFEEDBACK)(getProcAddr("glBeginTransformFeedback"))
	if gpBeginTransformFeedback == nil {
		return errors.New("glBeginTransformFeedback")
	}
	gpBeginTransformFeedbackEXT = (C.GPBEGINTRANSFORMFEEDBACKEXT)(getProcAddr("glBeginTransformFeedbackEXT"))
	gpBeginTransformFeedbackNV = (C.GPBEGINTRANSFORMFEEDBACKNV)(getProcAddr("glBeginTransformFeedbackNV"))
	gpBeginVertexShaderEXT = (C.GPBEGINVERTEXSHADEREXT)(getProcAddr("glBeginVertexShaderEXT"))
	gpBeginVideoCaptureNV = (C.GPBEGINVIDEOCAPTURENV)(getProcAddr("glBeginVideoCaptureNV"))
	gpBindAttribLocation = (C.GPBINDATTRIBLOCATION)(getProcAddr("glBindAttribLocation"))
	if gpBindAttribLocation == nil {
		return errors.New("glBindAttribLocation")
	}
	gpBindAttribLocationARB = (C.GPBINDATTRIBLOCATIONARB)(getProcAddr("glBindAttribLocationARB"))
	gpBindBuffer = (C.GPBINDBUFFER)(getProcAddr("glBindBuffer"))
	if gpBindBuffer == nil {
		return errors.New("glBindBuffer")
	}
	gpBindBufferARB = (C.GPBINDBUFFERARB)(getProcAddr("glBindBufferARB"))
	gpBindBufferBase = (C.GPBINDBUFFERBASE)(getProcAddr("glBindBufferBase"))
	if gpBindBufferBase == nil {
		return errors.New("glBindBufferBase")
	}
	gpBindBufferBaseEXT = (C.GPBINDBUFFERBASEEXT)(getProcAddr("glBindBufferBaseEXT"))
	gpBindBufferBaseNV = (C.GPBINDBUFFERBASENV)(getProcAddr("glBindBufferBaseNV"))
	gpBindBufferOffsetEXT = (C.GPBINDBUFFEROFFSETEXT)(getProcAddr("glBindBufferOffsetEXT"))
	gpBindBufferOffsetNV = (C.GPBINDBUFFEROFFSETNV)(getProcAddr("glBindBufferOffsetNV"))
	gpBindBufferRange = (C.GPBINDBUFFERRANGE)(getProcAddr("glBindBufferRange"))
	if gpBindBufferRange == nil {
		return errors.New("glBindBufferRange")
	}
	gpBindBufferRangeEXT = (C.GPBINDBUFFERRANGEEXT)(getProcAddr("glBindBufferRangeEXT"))
	gpBindBufferRangeNV = (C.GPBINDBUFFERRANGENV)(getProcAddr("glBindBufferRangeNV"))
	gpBindBuffersBase = (C.GPBINDBUFFERSBASE)(getProcAddr("glBindBuffersBase"))
	if gpBindBuffersBase == nil {
		return errors.New("glBindBuffersBase")
	}
	gpBindBuffersRange = (C.GPBINDBUFFERSRANGE)(getProcAddr("glBindBuffersRange"))
	if gpBindBuffersRange == nil {
		return errors.New("glBindBuffersRange")
	}
	gpBindFragDataLocation = (C.GPBINDFRAGDATALOCATION)(getProcAddr("glBindFragDataLocation"))
	if gpBindFragDataLocation == nil {
		return errors.New("glBindFragDataLocation")
	}
	gpBindFragDataLocationEXT = (C.GPBINDFRAGDATALOCATIONEXT)(getProcAddr("glBindFragDataLocationEXT"))
	gpBindFragDataLocationIndexed = (C.GPBINDFRAGDATALOCATIONINDEXED)(getProcAddr("glBindFragDataLocationIndexed"))
	if gpBindFragDataLocationIndexed == nil {
		return errors.New("glBindFragDataLocationIndexed")
	}
	gpBindFragmentShaderATI = (C.GPBINDFRAGMENTSHADERATI)(getProcAddr("glBindFragmentShaderATI"))
	gpBindFramebuffer = (C.GPBINDFRAMEBUFFER)(getProcAddr("glBindFramebuffer"))
	if gpBindFramebuffer == nil {
		return errors.New("glBindFramebuffer")
	}
	gpBindFramebufferEXT = (C.GPBINDFRAMEBUFFEREXT)(getProcAddr("glBindFramebufferEXT"))
	gpBindImageTexture = (C.GPBINDIMAGETEXTURE)(getProcAddr("glBindImageTexture"))
	if gpBindImageTexture == nil {
		return errors.New("glBindImageTexture")
	}
	gpBindImageTextureEXT = (C.GPBINDIMAGETEXTUREEXT)(getProcAddr("glBindImageTextureEXT"))
	gpBindImageTextures = (C.GPBINDIMAGETEXTURES)(getProcAddr("glBindImageTextures"))
	if gpBindImageTextures == nil {
		return errors.New("glBindImageTextures")
	}
	gpBindLightParameterEXT = (C.GPBINDLIGHTPARAMETEREXT)(getProcAddr("glBindLightParameterEXT"))
	gpBindMaterialParameterEXT = (C.GPBINDMATERIALPARAMETEREXT)(getProcAddr("glBindMaterialParameterEXT"))
	gpBindMultiTextureEXT = (C.GPBINDMULTITEXTUREEXT)(getProcAddr("glBindMultiTextureEXT"))
	gpBindParameterEXT = (C.GPBINDPARAMETEREXT)(getProcAddr("glBindParameterEXT"))
	gpBindProgramARB = (C.GPBINDPROGRAMARB)(getProcAddr("glBindProgramARB"))
	gpBindProgramNV = (C.GPBINDPROGRAMNV)(getProcAddr("glBindProgramNV"))
	gpBindProgramPipeline = (C.GPBINDPROGRAMPIPELINE)(getProcAddr("glBindProgramPipeline"))
	if gpBindProgramPipeline == nil {
		return errors.New("glBindProgramPipeline")
	}
	gpBindProgramPipelineEXT = (C.GPBINDPROGRAMPIPELINEEXT)(getProcAddr("glBindProgramPipelineEXT"))
	gpBindRenderbuffer = (C.GPBINDRENDERBUFFER)(getProcAddr("glBindRenderbuffer"))
	if gpBindRenderbuffer == nil {
		return errors.New("glBindRenderbuffer")
	}
	gpBindRenderbufferEXT = (C.GPBINDRENDERBUFFEREXT)(getProcAddr("glBindRenderbufferEXT"))
	gpBindSampler = (C.GPBINDSAMPLER)(getProcAddr("glBindSampler"))
	if gpBindSampler == nil {
		return errors.New("glBindSampler")
	}
	gpBindSamplers = (C.GPBINDSAMPLERS)(getProcAddr("glBindSamplers"))
	if gpBindSamplers == nil {
		return errors.New("glBindSamplers")
	}
	gpBindTexGenParameterEXT = (C.GPBINDTEXGENPARAMETEREXT)(getProcAddr("glBindTexGenParameterEXT"))
	gpBindTexture = (C.GPBINDTEXTURE)(getProcAddr("glBindTexture"))
	if gpBindTexture == nil {
		return errors.New("glBindTexture")
	}
	gpBindTextureEXT = (C.GPBINDTEXTUREEXT)(getProcAddr("glBindTextureEXT"))
	gpBindTextureUnit = (C.GPBINDTEXTUREUNIT)(getProcAddr("glBindTextureUnit"))
	if gpBindTextureUnit == nil {
		return errors.New("glBindTextureUnit")
	}
	gpBindTextureUnitParameterEXT = (C.GPBINDTEXTUREUNITPARAMETEREXT)(getProcAddr("glBindTextureUnitParameterEXT"))
	gpBindTextures = (C.GPBINDTEXTURES)(getProcAddr("glBindTextures"))
	if gpBindTextures == nil {
		return errors.New("glBindTextures")
	}
	gpBindTransformFeedback = (C.GPBINDTRANSFORMFEEDBACK)(getProcAddr("glBindTransformFeedback"))
	if gpBindTransformFeedback == nil {
		return errors.New("glBindTransformFeedback")
	}
	gpBindTransformFeedbackNV = (C.GPBINDTRANSFORMFEEDBACKNV)(getProcAddr("glBindTransformFeedbackNV"))
	gpBindVertexArray = (C.GPBINDVERTEXARRAY)(getProcAddr("glBindVertexArray"))
	if gpBindVertexArray == nil {
		return errors.New("glBindVertexArray")
	}
	gpBindVertexArrayAPPLE = (C.GPBINDVERTEXARRAYAPPLE)(getProcAddr("glBindVertexArrayAPPLE"))
	gpBindVertexBuffer = (C.GPBINDVERTEXBUFFER)(getProcAddr("glBindVertexBuffer"))
	if gpBindVertexBuffer == nil {
		return errors.New("glBindVertexBuffer")
	}
	gpBindVertexBuffers = (C.GPBINDVERTEXBUFFERS)(getProcAddr("glBindVertexBuffers"))
	if gpBindVertexBuffers == nil {
		return errors.New("glBindVertexBuffers")
	}
	gpBindVertexShaderEXT = (C.GPBINDVERTEXSHADEREXT)(getProcAddr("glBindVertexShaderEXT"))
	gpBindVideoCaptureStreamBufferNV = (C.GPBINDVIDEOCAPTURESTREAMBUFFERNV)(getProcAddr("glBindVideoCaptureStreamBufferNV"))
	gpBindVideoCaptureStreamTextureNV = (C.GPBINDVIDEOCAPTURESTREAMTEXTURENV)(getProcAddr("glBindVideoCaptureStreamTextureNV"))
	gpBinormal3bEXT = (C.GPBINORMAL3BEXT)(getProcAddr("glBinormal3bEXT"))
	gpBinormal3bvEXT = (C.GPBINORMAL3BVEXT)(getProcAddr("glBinormal3bvEXT"))
	gpBinormal3dEXT = (C.GPBINORMAL3DEXT)(getProcAddr("glBinormal3dEXT"))
	gpBinormal3dvEXT = (C.GPBINORMAL3DVEXT)(getProcAddr("glBinormal3dvEXT"))
	gpBinormal3fEXT = (C.GPBINORMAL3FEXT)(getProcAddr("glBinormal3fEXT"))
	gpBinormal3fvEXT = (C.GPBINORMAL3FVEXT)(getProcAddr("glBinormal3fvEXT"))
	gpBinormal3iEXT = (C.GPBINORMAL3IEXT)(getProcAddr("glBinormal3iEXT"))
	gpBinormal3ivEXT = (C.GPBINORMAL3IVEXT)(getProcAddr("glBinormal3ivEXT"))
	gpBinormal3sEXT = (C.GPBINORMAL3SEXT)(getProcAddr("glBinormal3sEXT"))
	gpBinormal3svEXT = (C.GPBINORMAL3SVEXT)(getProcAddr("glBinormal3svEXT"))
	gpBinormalPointerEXT = (C.GPBINORMALPOINTEREXT)(getProcAddr("glBinormalPointerEXT"))
	gpBitmap = (C.GPBITMAP)(getProcAddr("glBitmap"))
	if gpBitmap == nil {
		return errors.New("glBitmap")
	}
	gpBitmapxOES = (C.GPBITMAPXOES)(getProcAddr("glBitmapxOES"))
	gpBlendBarrierKHR = (C.GPBLENDBARRIERKHR)(getProcAddr("glBlendBarrierKHR"))
	gpBlendBarrierNV = (C.GPBLENDBARRIERNV)(getProcAddr("glBlendBarrierNV"))
	gpBlendColor = (C.GPBLENDCOLOR)(getProcAddr("glBlendColor"))
	if gpBlendColor == nil {
		return errors.New("glBlendColor")
	}
	gpBlendColorEXT = (C.GPBLENDCOLOREXT)(getProcAddr("glBlendColorEXT"))
	gpBlendColorxOES = (C.GPBLENDCOLORXOES)(getProcAddr("glBlendColorxOES"))
	gpBlendEquation = (C.GPBLENDEQUATION)(getProcAddr("glBlendEquation"))
	if gpBlendEquation == nil {
		return errors.New("glBlendEquation")
	}
	gpBlendEquationEXT = (C.GPBLENDEQUATIONEXT)(getProcAddr("glBlendEquationEXT"))
	gpBlendEquationIndexedAMD = (C.GPBLENDEQUATIONINDEXEDAMD)(getProcAddr("glBlendEquationIndexedAMD"))
	gpBlendEquationSeparate = (C.GPBLENDEQUATIONSEPARATE)(getProcAddr("glBlendEquationSeparate"))
	if gpBlendEquationSeparate == nil {
		return errors.New("glBlendEquationSeparate")
	}
	gpBlendEquationSeparateEXT = (C.GPBLENDEQUATIONSEPARATEEXT)(getProcAddr("glBlendEquationSeparateEXT"))
	gpBlendEquationSeparateIndexedAMD = (C.GPBLENDEQUATIONSEPARATEINDEXEDAMD)(getProcAddr("glBlendEquationSeparateIndexedAMD"))
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
	gpBlendFuncIndexedAMD = (C.GPBLENDFUNCINDEXEDAMD)(getProcAddr("glBlendFuncIndexedAMD"))
	gpBlendFuncSeparate = (C.GPBLENDFUNCSEPARATE)(getProcAddr("glBlendFuncSeparate"))
	if gpBlendFuncSeparate == nil {
		return errors.New("glBlendFuncSeparate")
	}
	gpBlendFuncSeparateEXT = (C.GPBLENDFUNCSEPARATEEXT)(getProcAddr("glBlendFuncSeparateEXT"))
	gpBlendFuncSeparateINGR = (C.GPBLENDFUNCSEPARATEINGR)(getProcAddr("glBlendFuncSeparateINGR"))
	gpBlendFuncSeparateIndexedAMD = (C.GPBLENDFUNCSEPARATEINDEXEDAMD)(getProcAddr("glBlendFuncSeparateIndexedAMD"))
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
	gpBlitFramebufferEXT = (C.GPBLITFRAMEBUFFEREXT)(getProcAddr("glBlitFramebufferEXT"))
	gpBlitNamedFramebuffer = (C.GPBLITNAMEDFRAMEBUFFER)(getProcAddr("glBlitNamedFramebuffer"))
	if gpBlitNamedFramebuffer == nil {
		return errors.New("glBlitNamedFramebuffer")
	}
	gpBufferAddressRangeNV = (C.GPBUFFERADDRESSRANGENV)(getProcAddr("glBufferAddressRangeNV"))
	gpBufferData = (C.GPBUFFERDATA)(getProcAddr("glBufferData"))
	if gpBufferData == nil {
		return errors.New("glBufferData")
	}
	gpBufferDataARB = (C.GPBUFFERDATAARB)(getProcAddr("glBufferDataARB"))
	gpBufferPageCommitmentARB = (C.GPBUFFERPAGECOMMITMENTARB)(getProcAddr("glBufferPageCommitmentARB"))
	gpBufferParameteriAPPLE = (C.GPBUFFERPARAMETERIAPPLE)(getProcAddr("glBufferParameteriAPPLE"))
	gpBufferStorage = (C.GPBUFFERSTORAGE)(getProcAddr("glBufferStorage"))
	if gpBufferStorage == nil {
		return errors.New("glBufferStorage")
	}
	gpBufferStorageExternalEXT = (C.GPBUFFERSTORAGEEXTERNALEXT)(getProcAddr("glBufferStorageExternalEXT"))
	gpBufferStorageMemEXT = (C.GPBUFFERSTORAGEMEMEXT)(getProcAddr("glBufferStorageMemEXT"))
	gpBufferSubData = (C.GPBUFFERSUBDATA)(getProcAddr("glBufferSubData"))
	if gpBufferSubData == nil {
		return errors.New("glBufferSubData")
	}
	gpBufferSubDataARB = (C.GPBUFFERSUBDATAARB)(getProcAddr("glBufferSubDataARB"))
	gpCallCommandListNV = (C.GPCALLCOMMANDLISTNV)(getProcAddr("glCallCommandListNV"))
	gpCallList = (C.GPCALLLIST)(getProcAddr("glCallList"))
	if gpCallList == nil {
		return errors.New("glCallList")
	}
	gpCallLists = (C.GPCALLLISTS)(getProcAddr("glCallLists"))
	if gpCallLists == nil {
		return errors.New("glCallLists")
	}
	gpCheckFramebufferStatus = (C.GPCHECKFRAMEBUFFERSTATUS)(getProcAddr("glCheckFramebufferStatus"))
	if gpCheckFramebufferStatus == nil {
		return errors.New("glCheckFramebufferStatus")
	}
	gpCheckFramebufferStatusEXT = (C.GPCHECKFRAMEBUFFERSTATUSEXT)(getProcAddr("glCheckFramebufferStatusEXT"))
	gpCheckNamedFramebufferStatus = (C.GPCHECKNAMEDFRAMEBUFFERSTATUS)(getProcAddr("glCheckNamedFramebufferStatus"))
	if gpCheckNamedFramebufferStatus == nil {
		return errors.New("glCheckNamedFramebufferStatus")
	}
	gpCheckNamedFramebufferStatusEXT = (C.GPCHECKNAMEDFRAMEBUFFERSTATUSEXT)(getProcAddr("glCheckNamedFramebufferStatusEXT"))
	gpClampColor = (C.GPCLAMPCOLOR)(getProcAddr("glClampColor"))
	if gpClampColor == nil {
		return errors.New("glClampColor")
	}
	gpClampColorARB = (C.GPCLAMPCOLORARB)(getProcAddr("glClampColorARB"))
	gpClear = (C.GPCLEAR)(getProcAddr("glClear"))
	if gpClear == nil {
		return errors.New("glClear")
	}
	gpClearAccum = (C.GPCLEARACCUM)(getProcAddr("glClearAccum"))
	if gpClearAccum == nil {
		return errors.New("glClearAccum")
	}
	gpClearAccumxOES = (C.GPCLEARACCUMXOES)(getProcAddr("glClearAccumxOES"))
	gpClearBufferData = (C.GPCLEARBUFFERDATA)(getProcAddr("glClearBufferData"))
	if gpClearBufferData == nil {
		return errors.New("glClearBufferData")
	}
	gpClearBufferSubData = (C.GPCLEARBUFFERSUBDATA)(getProcAddr("glClearBufferSubData"))
	if gpClearBufferSubData == nil {
		return errors.New("glClearBufferSubData")
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
	gpClearColorIiEXT = (C.GPCLEARCOLORIIEXT)(getProcAddr("glClearColorIiEXT"))
	gpClearColorIuiEXT = (C.GPCLEARCOLORIUIEXT)(getProcAddr("glClearColorIuiEXT"))
	gpClearColorxOES = (C.GPCLEARCOLORXOES)(getProcAddr("glClearColorxOES"))
	gpClearDepth = (C.GPCLEARDEPTH)(getProcAddr("glClearDepth"))
	if gpClearDepth == nil {
		return errors.New("glClearDepth")
	}
	gpClearDepthdNV = (C.GPCLEARDEPTHDNV)(getProcAddr("glClearDepthdNV"))
	gpClearDepthf = (C.GPCLEARDEPTHF)(getProcAddr("glClearDepthf"))
	if gpClearDepthf == nil {
		return errors.New("glClearDepthf")
	}
	gpClearDepthfOES = (C.GPCLEARDEPTHFOES)(getProcAddr("glClearDepthfOES"))
	gpClearDepthxOES = (C.GPCLEARDEPTHXOES)(getProcAddr("glClearDepthxOES"))
	gpClearIndex = (C.GPCLEARINDEX)(getProcAddr("glClearIndex"))
	if gpClearIndex == nil {
		return errors.New("glClearIndex")
	}
	gpClearNamedBufferData = (C.GPCLEARNAMEDBUFFERDATA)(getProcAddr("glClearNamedBufferData"))
	if gpClearNamedBufferData == nil {
		return errors.New("glClearNamedBufferData")
	}
	gpClearNamedBufferDataEXT = (C.GPCLEARNAMEDBUFFERDATAEXT)(getProcAddr("glClearNamedBufferDataEXT"))
	gpClearNamedBufferSubData = (C.GPCLEARNAMEDBUFFERSUBDATA)(getProcAddr("glClearNamedBufferSubData"))
	if gpClearNamedBufferSubData == nil {
		return errors.New("glClearNamedBufferSubData")
	}
	gpClearNamedBufferSubDataEXT = (C.GPCLEARNAMEDBUFFERSUBDATAEXT)(getProcAddr("glClearNamedBufferSubDataEXT"))
	gpClearNamedFramebufferfi = (C.GPCLEARNAMEDFRAMEBUFFERFI)(getProcAddr("glClearNamedFramebufferfi"))
	if gpClearNamedFramebufferfi == nil {
		return errors.New("glClearNamedFramebufferfi")
	}
	gpClearNamedFramebufferfv = (C.GPCLEARNAMEDFRAMEBUFFERFV)(getProcAddr("glClearNamedFramebufferfv"))
	if gpClearNamedFramebufferfv == nil {
		return errors.New("glClearNamedFramebufferfv")
	}
	gpClearNamedFramebufferiv = (C.GPCLEARNAMEDFRAMEBUFFERIV)(getProcAddr("glClearNamedFramebufferiv"))
	if gpClearNamedFramebufferiv == nil {
		return errors.New("glClearNamedFramebufferiv")
	}
	gpClearNamedFramebufferuiv = (C.GPCLEARNAMEDFRAMEBUFFERUIV)(getProcAddr("glClearNamedFramebufferuiv"))
	if gpClearNamedFramebufferuiv == nil {
		return errors.New("glClearNamedFramebufferuiv")
	}
	gpClearStencil = (C.GPCLEARSTENCIL)(getProcAddr("glClearStencil"))
	if gpClearStencil == nil {
		return errors.New("glClearStencil")
	}
	gpClearTexImage = (C.GPCLEARTEXIMAGE)(getProcAddr("glClearTexImage"))
	if gpClearTexImage == nil {
		return errors.New("glClearTexImage")
	}
	gpClearTexSubImage = (C.GPCLEARTEXSUBIMAGE)(getProcAddr("glClearTexSubImage"))
	if gpClearTexSubImage == nil {
		return errors.New("glClearTexSubImage")
	}
	gpClientActiveTexture = (C.GPCLIENTACTIVETEXTURE)(getProcAddr("glClientActiveTexture"))
	if gpClientActiveTexture == nil {
		return errors.New("glClientActiveTexture")
	}
	gpClientActiveTextureARB = (C.GPCLIENTACTIVETEXTUREARB)(getProcAddr("glClientActiveTextureARB"))
	gpClientActiveVertexStreamATI = (C.GPCLIENTACTIVEVERTEXSTREAMATI)(getProcAddr("glClientActiveVertexStreamATI"))
	gpClientAttribDefaultEXT = (C.GPCLIENTATTRIBDEFAULTEXT)(getProcAddr("glClientAttribDefaultEXT"))
	gpClientWaitSync = (C.GPCLIENTWAITSYNC)(getProcAddr("glClientWaitSync"))
	if gpClientWaitSync == nil {
		return errors.New("glClientWaitSync")
	}
	gpClipControl = (C.GPCLIPCONTROL)(getProcAddr("glClipControl"))
	if gpClipControl == nil {
		return errors.New("glClipControl")
	}
	gpClipPlane = (C.GPCLIPPLANE)(getProcAddr("glClipPlane"))
	if gpClipPlane == nil {
		return errors.New("glClipPlane")
	}
	gpClipPlanefOES = (C.GPCLIPPLANEFOES)(getProcAddr("glClipPlanefOES"))
	gpClipPlanexOES = (C.GPCLIPPLANEXOES)(getProcAddr("glClipPlanexOES"))
	gpColor3b = (C.GPCOLOR3B)(getProcAddr("glColor3b"))
	if gpColor3b == nil {
		return errors.New("glColor3b")
	}
	gpColor3bv = (C.GPCOLOR3BV)(getProcAddr("glColor3bv"))
	if gpColor3bv == nil {
		return errors.New("glColor3bv")
	}
	gpColor3d = (C.GPCOLOR3D)(getProcAddr("glColor3d"))
	if gpColor3d == nil {
		return errors.New("glColor3d")
	}
	gpColor3dv = (C.GPCOLOR3DV)(getProcAddr("glColor3dv"))
	if gpColor3dv == nil {
		return errors.New("glColor3dv")
	}
	gpColor3f = (C.GPCOLOR3F)(getProcAddr("glColor3f"))
	if gpColor3f == nil {
		return errors.New("glColor3f")
	}
	gpColor3fVertex3fSUN = (C.GPCOLOR3FVERTEX3FSUN)(getProcAddr("glColor3fVertex3fSUN"))
	gpColor3fVertex3fvSUN = (C.GPCOLOR3FVERTEX3FVSUN)(getProcAddr("glColor3fVertex3fvSUN"))
	gpColor3fv = (C.GPCOLOR3FV)(getProcAddr("glColor3fv"))
	if gpColor3fv == nil {
		return errors.New("glColor3fv")
	}
	gpColor3hNV = (C.GPCOLOR3HNV)(getProcAddr("glColor3hNV"))
	gpColor3hvNV = (C.GPCOLOR3HVNV)(getProcAddr("glColor3hvNV"))
	gpColor3i = (C.GPCOLOR3I)(getProcAddr("glColor3i"))
	if gpColor3i == nil {
		return errors.New("glColor3i")
	}
	gpColor3iv = (C.GPCOLOR3IV)(getProcAddr("glColor3iv"))
	if gpColor3iv == nil {
		return errors.New("glColor3iv")
	}
	gpColor3s = (C.GPCOLOR3S)(getProcAddr("glColor3s"))
	if gpColor3s == nil {
		return errors.New("glColor3s")
	}
	gpColor3sv = (C.GPCOLOR3SV)(getProcAddr("glColor3sv"))
	if gpColor3sv == nil {
		return errors.New("glColor3sv")
	}
	gpColor3ub = (C.GPCOLOR3UB)(getProcAddr("glColor3ub"))
	if gpColor3ub == nil {
		return errors.New("glColor3ub")
	}
	gpColor3ubv = (C.GPCOLOR3UBV)(getProcAddr("glColor3ubv"))
	if gpColor3ubv == nil {
		return errors.New("glColor3ubv")
	}
	gpColor3ui = (C.GPCOLOR3UI)(getProcAddr("glColor3ui"))
	if gpColor3ui == nil {
		return errors.New("glColor3ui")
	}
	gpColor3uiv = (C.GPCOLOR3UIV)(getProcAddr("glColor3uiv"))
	if gpColor3uiv == nil {
		return errors.New("glColor3uiv")
	}
	gpColor3us = (C.GPCOLOR3US)(getProcAddr("glColor3us"))
	if gpColor3us == nil {
		return errors.New("glColor3us")
	}
	gpColor3usv = (C.GPCOLOR3USV)(getProcAddr("glColor3usv"))
	if gpColor3usv == nil {
		return errors.New("glColor3usv")
	}
	gpColor3xOES = (C.GPCOLOR3XOES)(getProcAddr("glColor3xOES"))
	gpColor3xvOES = (C.GPCOLOR3XVOES)(getProcAddr("glColor3xvOES"))
	gpColor4b = (C.GPCOLOR4B)(getProcAddr("glColor4b"))
	if gpColor4b == nil {
		return errors.New("glColor4b")
	}
	gpColor4bv = (C.GPCOLOR4BV)(getProcAddr("glColor4bv"))
	if gpColor4bv == nil {
		return errors.New("glColor4bv")
	}
	gpColor4d = (C.GPCOLOR4D)(getProcAddr("glColor4d"))
	if gpColor4d == nil {
		return errors.New("glColor4d")
	}
	gpColor4dv = (C.GPCOLOR4DV)(getProcAddr("glColor4dv"))
	if gpColor4dv == nil {
		return errors.New("glColor4dv")
	}
	gpColor4f = (C.GPCOLOR4F)(getProcAddr("glColor4f"))
	if gpColor4f == nil {
		return errors.New("glColor4f")
	}
	gpColor4fNormal3fVertex3fSUN = (C.GPCOLOR4FNORMAL3FVERTEX3FSUN)(getProcAddr("glColor4fNormal3fVertex3fSUN"))
	gpColor4fNormal3fVertex3fvSUN = (C.GPCOLOR4FNORMAL3FVERTEX3FVSUN)(getProcAddr("glColor4fNormal3fVertex3fvSUN"))
	gpColor4fv = (C.GPCOLOR4FV)(getProcAddr("glColor4fv"))
	if gpColor4fv == nil {
		return errors.New("glColor4fv")
	}
	gpColor4hNV = (C.GPCOLOR4HNV)(getProcAddr("glColor4hNV"))
	gpColor4hvNV = (C.GPCOLOR4HVNV)(getProcAddr("glColor4hvNV"))
	gpColor4i = (C.GPCOLOR4I)(getProcAddr("glColor4i"))
	if gpColor4i == nil {
		return errors.New("glColor4i")
	}
	gpColor4iv = (C.GPCOLOR4IV)(getProcAddr("glColor4iv"))
	if gpColor4iv == nil {
		return errors.New("glColor4iv")
	}
	gpColor4s = (C.GPCOLOR4S)(getProcAddr("glColor4s"))
	if gpColor4s == nil {
		return errors.New("glColor4s")
	}
	gpColor4sv = (C.GPCOLOR4SV)(getProcAddr("glColor4sv"))
	if gpColor4sv == nil {
		return errors.New("glColor4sv")
	}
	gpColor4ub = (C.GPCOLOR4UB)(getProcAddr("glColor4ub"))
	if gpColor4ub == nil {
		return errors.New("glColor4ub")
	}
	gpColor4ubVertex2fSUN = (C.GPCOLOR4UBVERTEX2FSUN)(getProcAddr("glColor4ubVertex2fSUN"))
	gpColor4ubVertex2fvSUN = (C.GPCOLOR4UBVERTEX2FVSUN)(getProcAddr("glColor4ubVertex2fvSUN"))
	gpColor4ubVertex3fSUN = (C.GPCOLOR4UBVERTEX3FSUN)(getProcAddr("glColor4ubVertex3fSUN"))
	gpColor4ubVertex3fvSUN = (C.GPCOLOR4UBVERTEX3FVSUN)(getProcAddr("glColor4ubVertex3fvSUN"))
	gpColor4ubv = (C.GPCOLOR4UBV)(getProcAddr("glColor4ubv"))
	if gpColor4ubv == nil {
		return errors.New("glColor4ubv")
	}
	gpColor4ui = (C.GPCOLOR4UI)(getProcAddr("glColor4ui"))
	if gpColor4ui == nil {
		return errors.New("glColor4ui")
	}
	gpColor4uiv = (C.GPCOLOR4UIV)(getProcAddr("glColor4uiv"))
	if gpColor4uiv == nil {
		return errors.New("glColor4uiv")
	}
	gpColor4us = (C.GPCOLOR4US)(getProcAddr("glColor4us"))
	if gpColor4us == nil {
		return errors.New("glColor4us")
	}
	gpColor4usv = (C.GPCOLOR4USV)(getProcAddr("glColor4usv"))
	if gpColor4usv == nil {
		return errors.New("glColor4usv")
	}
	gpColor4xOES = (C.GPCOLOR4XOES)(getProcAddr("glColor4xOES"))
	gpColor4xvOES = (C.GPCOLOR4XVOES)(getProcAddr("glColor4xvOES"))
	gpColorFormatNV = (C.GPCOLORFORMATNV)(getProcAddr("glColorFormatNV"))
	gpColorFragmentOp1ATI = (C.GPCOLORFRAGMENTOP1ATI)(getProcAddr("glColorFragmentOp1ATI"))
	gpColorFragmentOp2ATI = (C.GPCOLORFRAGMENTOP2ATI)(getProcAddr("glColorFragmentOp2ATI"))
	gpColorFragmentOp3ATI = (C.GPCOLORFRAGMENTOP3ATI)(getProcAddr("glColorFragmentOp3ATI"))
	gpColorMask = (C.GPCOLORMASK)(getProcAddr("glColorMask"))
	if gpColorMask == nil {
		return errors.New("glColorMask")
	}
	gpColorMaskIndexedEXT = (C.GPCOLORMASKINDEXEDEXT)(getProcAddr("glColorMaskIndexedEXT"))
	gpColorMaski = (C.GPCOLORMASKI)(getProcAddr("glColorMaski"))
	if gpColorMaski == nil {
		return errors.New("glColorMaski")
	}
	gpColorMaterial = (C.GPCOLORMATERIAL)(getProcAddr("glColorMaterial"))
	if gpColorMaterial == nil {
		return errors.New("glColorMaterial")
	}
	gpColorP3ui = (C.GPCOLORP3UI)(getProcAddr("glColorP3ui"))
	if gpColorP3ui == nil {
		return errors.New("glColorP3ui")
	}
	gpColorP3uiv = (C.GPCOLORP3UIV)(getProcAddr("glColorP3uiv"))
	if gpColorP3uiv == nil {
		return errors.New("glColorP3uiv")
	}
	gpColorP4ui = (C.GPCOLORP4UI)(getProcAddr("glColorP4ui"))
	if gpColorP4ui == nil {
		return errors.New("glColorP4ui")
	}
	gpColorP4uiv = (C.GPCOLORP4UIV)(getProcAddr("glColorP4uiv"))
	if gpColorP4uiv == nil {
		return errors.New("glColorP4uiv")
	}
	gpColorPointer = (C.GPCOLORPOINTER)(getProcAddr("glColorPointer"))
	if gpColorPointer == nil {
		return errors.New("glColorPointer")
	}
	gpColorPointerEXT = (C.GPCOLORPOINTEREXT)(getProcAddr("glColorPointerEXT"))
	gpColorPointerListIBM = (C.GPCOLORPOINTERLISTIBM)(getProcAddr("glColorPointerListIBM"))
	gpColorPointervINTEL = (C.GPCOLORPOINTERVINTEL)(getProcAddr("glColorPointervINTEL"))
	gpColorSubTable = (C.GPCOLORSUBTABLE)(getProcAddr("glColorSubTable"))
	gpColorSubTableEXT = (C.GPCOLORSUBTABLEEXT)(getProcAddr("glColorSubTableEXT"))
	gpColorTable = (C.GPCOLORTABLE)(getProcAddr("glColorTable"))
	gpColorTableEXT = (C.GPCOLORTABLEEXT)(getProcAddr("glColorTableEXT"))
	gpColorTableParameterfv = (C.GPCOLORTABLEPARAMETERFV)(getProcAddr("glColorTableParameterfv"))
	gpColorTableParameterfvSGI = (C.GPCOLORTABLEPARAMETERFVSGI)(getProcAddr("glColorTableParameterfvSGI"))
	gpColorTableParameteriv = (C.GPCOLORTABLEPARAMETERIV)(getProcAddr("glColorTableParameteriv"))
	gpColorTableParameterivSGI = (C.GPCOLORTABLEPARAMETERIVSGI)(getProcAddr("glColorTableParameterivSGI"))
	gpColorTableSGI = (C.GPCOLORTABLESGI)(getProcAddr("glColorTableSGI"))
	gpCombinerInputNV = (C.GPCOMBINERINPUTNV)(getProcAddr("glCombinerInputNV"))
	gpCombinerOutputNV = (C.GPCOMBINEROUTPUTNV)(getProcAddr("glCombinerOutputNV"))
	gpCombinerParameterfNV = (C.GPCOMBINERPARAMETERFNV)(getProcAddr("glCombinerParameterfNV"))
	gpCombinerParameterfvNV = (C.GPCOMBINERPARAMETERFVNV)(getProcAddr("glCombinerParameterfvNV"))
	gpCombinerParameteriNV = (C.GPCOMBINERPARAMETERINV)(getProcAddr("glCombinerParameteriNV"))
	gpCombinerParameterivNV = (C.GPCOMBINERPARAMETERIVNV)(getProcAddr("glCombinerParameterivNV"))
	gpCombinerStageParameterfvNV = (C.GPCOMBINERSTAGEPARAMETERFVNV)(getProcAddr("glCombinerStageParameterfvNV"))
	gpCommandListSegmentsNV = (C.GPCOMMANDLISTSEGMENTSNV)(getProcAddr("glCommandListSegmentsNV"))
	gpCompileCommandListNV = (C.GPCOMPILECOMMANDLISTNV)(getProcAddr("glCompileCommandListNV"))
	gpCompileShader = (C.GPCOMPILESHADER)(getProcAddr("glCompileShader"))
	if gpCompileShader == nil {
		return errors.New("glCompileShader")
	}
	gpCompileShaderARB = (C.GPCOMPILESHADERARB)(getProcAddr("glCompileShaderARB"))
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
	gpCompressedTexImage1DARB = (C.GPCOMPRESSEDTEXIMAGE1DARB)(getProcAddr("glCompressedTexImage1DARB"))
	gpCompressedTexImage2D = (C.GPCOMPRESSEDTEXIMAGE2D)(getProcAddr("glCompressedTexImage2D"))
	if gpCompressedTexImage2D == nil {
		return errors.New("glCompressedTexImage2D")
	}
	gpCompressedTexImage2DARB = (C.GPCOMPRESSEDTEXIMAGE2DARB)(getProcAddr("glCompressedTexImage2DARB"))
	gpCompressedTexImage3D = (C.GPCOMPRESSEDTEXIMAGE3D)(getProcAddr("glCompressedTexImage3D"))
	if gpCompressedTexImage3D == nil {
		return errors.New("glCompressedTexImage3D")
	}
	gpCompressedTexImage3DARB = (C.GPCOMPRESSEDTEXIMAGE3DARB)(getProcAddr("glCompressedTexImage3DARB"))
	gpCompressedTexSubImage1D = (C.GPCOMPRESSEDTEXSUBIMAGE1D)(getProcAddr("glCompressedTexSubImage1D"))
	if gpCompressedTexSubImage1D == nil {
		return errors.New("glCompressedTexSubImage1D")
	}
	gpCompressedTexSubImage1DARB = (C.GPCOMPRESSEDTEXSUBIMAGE1DARB)(getProcAddr("glCompressedTexSubImage1DARB"))
	gpCompressedTexSubImage2D = (C.GPCOMPRESSEDTEXSUBIMAGE2D)(getProcAddr("glCompressedTexSubImage2D"))
	if gpCompressedTexSubImage2D == nil {
		return errors.New("glCompressedTexSubImage2D")
	}
	gpCompressedTexSubImage2DARB = (C.GPCOMPRESSEDTEXSUBIMAGE2DARB)(getProcAddr("glCompressedTexSubImage2DARB"))
	gpCompressedTexSubImage3D = (C.GPCOMPRESSEDTEXSUBIMAGE3D)(getProcAddr("glCompressedTexSubImage3D"))
	if gpCompressedTexSubImage3D == nil {
		return errors.New("glCompressedTexSubImage3D")
	}
	gpCompressedTexSubImage3DARB = (C.GPCOMPRESSEDTEXSUBIMAGE3DARB)(getProcAddr("glCompressedTexSubImage3DARB"))
	gpCompressedTextureImage1DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE1DEXT)(getProcAddr("glCompressedTextureImage1DEXT"))
	gpCompressedTextureImage2DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE2DEXT)(getProcAddr("glCompressedTextureImage2DEXT"))
	gpCompressedTextureImage3DEXT = (C.GPCOMPRESSEDTEXTUREIMAGE3DEXT)(getProcAddr("glCompressedTextureImage3DEXT"))
	gpCompressedTextureSubImage1D = (C.GPCOMPRESSEDTEXTURESUBIMAGE1D)(getProcAddr("glCompressedTextureSubImage1D"))
	if gpCompressedTextureSubImage1D == nil {
		return errors.New("glCompressedTextureSubImage1D")
	}
	gpCompressedTextureSubImage1DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE1DEXT)(getProcAddr("glCompressedTextureSubImage1DEXT"))
	gpCompressedTextureSubImage2D = (C.GPCOMPRESSEDTEXTURESUBIMAGE2D)(getProcAddr("glCompressedTextureSubImage2D"))
	if gpCompressedTextureSubImage2D == nil {
		return errors.New("glCompressedTextureSubImage2D")
	}
	gpCompressedTextureSubImage2DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE2DEXT)(getProcAddr("glCompressedTextureSubImage2DEXT"))
	gpCompressedTextureSubImage3D = (C.GPCOMPRESSEDTEXTURESUBIMAGE3D)(getProcAddr("glCompressedTextureSubImage3D"))
	if gpCompressedTextureSubImage3D == nil {
		return errors.New("glCompressedTextureSubImage3D")
	}
	gpCompressedTextureSubImage3DEXT = (C.GPCOMPRESSEDTEXTURESUBIMAGE3DEXT)(getProcAddr("glCompressedTextureSubImage3DEXT"))
	gpConservativeRasterParameterfNV = (C.GPCONSERVATIVERASTERPARAMETERFNV)(getProcAddr("glConservativeRasterParameterfNV"))
	gpConservativeRasterParameteriNV = (C.GPCONSERVATIVERASTERPARAMETERINV)(getProcAddr("glConservativeRasterParameteriNV"))
	gpConvolutionFilter1D = (C.GPCONVOLUTIONFILTER1D)(getProcAddr("glConvolutionFilter1D"))
	gpConvolutionFilter1DEXT = (C.GPCONVOLUTIONFILTER1DEXT)(getProcAddr("glConvolutionFilter1DEXT"))
	gpConvolutionFilter2D = (C.GPCONVOLUTIONFILTER2D)(getProcAddr("glConvolutionFilter2D"))
	gpConvolutionFilter2DEXT = (C.GPCONVOLUTIONFILTER2DEXT)(getProcAddr("glConvolutionFilter2DEXT"))
	gpConvolutionParameterf = (C.GPCONVOLUTIONPARAMETERF)(getProcAddr("glConvolutionParameterf"))
	gpConvolutionParameterfEXT = (C.GPCONVOLUTIONPARAMETERFEXT)(getProcAddr("glConvolutionParameterfEXT"))
	gpConvolutionParameterfv = (C.GPCONVOLUTIONPARAMETERFV)(getProcAddr("glConvolutionParameterfv"))
	gpConvolutionParameterfvEXT = (C.GPCONVOLUTIONPARAMETERFVEXT)(getProcAddr("glConvolutionParameterfvEXT"))
	gpConvolutionParameteri = (C.GPCONVOLUTIONPARAMETERI)(getProcAddr("glConvolutionParameteri"))
	gpConvolutionParameteriEXT = (C.GPCONVOLUTIONPARAMETERIEXT)(getProcAddr("glConvolutionParameteriEXT"))
	gpConvolutionParameteriv = (C.GPCONVOLUTIONPARAMETERIV)(getProcAddr("glConvolutionParameteriv"))
	gpConvolutionParameterivEXT = (C.GPCONVOLUTIONPARAMETERIVEXT)(getProcAddr("glConvolutionParameterivEXT"))
	gpConvolutionParameterxOES = (C.GPCONVOLUTIONPARAMETERXOES)(getProcAddr("glConvolutionParameterxOES"))
	gpConvolutionParameterxvOES = (C.GPCONVOLUTIONPARAMETERXVOES)(getProcAddr("glConvolutionParameterxvOES"))
	gpCopyBufferSubData = (C.GPCOPYBUFFERSUBDATA)(getProcAddr("glCopyBufferSubData"))
	if gpCopyBufferSubData == nil {
		return errors.New("glCopyBufferSubData")
	}
	gpCopyColorSubTable = (C.GPCOPYCOLORSUBTABLE)(getProcAddr("glCopyColorSubTable"))
	gpCopyColorSubTableEXT = (C.GPCOPYCOLORSUBTABLEEXT)(getProcAddr("glCopyColorSubTableEXT"))
	gpCopyColorTable = (C.GPCOPYCOLORTABLE)(getProcAddr("glCopyColorTable"))
	gpCopyColorTableSGI = (C.GPCOPYCOLORTABLESGI)(getProcAddr("glCopyColorTableSGI"))
	gpCopyConvolutionFilter1D = (C.GPCOPYCONVOLUTIONFILTER1D)(getProcAddr("glCopyConvolutionFilter1D"))
	gpCopyConvolutionFilter1DEXT = (C.GPCOPYCONVOLUTIONFILTER1DEXT)(getProcAddr("glCopyConvolutionFilter1DEXT"))
	gpCopyConvolutionFilter2D = (C.GPCOPYCONVOLUTIONFILTER2D)(getProcAddr("glCopyConvolutionFilter2D"))
	gpCopyConvolutionFilter2DEXT = (C.GPCOPYCONVOLUTIONFILTER2DEXT)(getProcAddr("glCopyConvolutionFilter2DEXT"))
	gpCopyImageSubData = (C.GPCOPYIMAGESUBDATA)(getProcAddr("glCopyImageSubData"))
	if gpCopyImageSubData == nil {
		return errors.New("glCopyImageSubData")
	}
	gpCopyImageSubDataNV = (C.GPCOPYIMAGESUBDATANV)(getProcAddr("glCopyImageSubDataNV"))
	gpCopyMultiTexImage1DEXT = (C.GPCOPYMULTITEXIMAGE1DEXT)(getProcAddr("glCopyMultiTexImage1DEXT"))
	gpCopyMultiTexImage2DEXT = (C.GPCOPYMULTITEXIMAGE2DEXT)(getProcAddr("glCopyMultiTexImage2DEXT"))
	gpCopyMultiTexSubImage1DEXT = (C.GPCOPYMULTITEXSUBIMAGE1DEXT)(getProcAddr("glCopyMultiTexSubImage1DEXT"))
	gpCopyMultiTexSubImage2DEXT = (C.GPCOPYMULTITEXSUBIMAGE2DEXT)(getProcAddr("glCopyMultiTexSubImage2DEXT"))
	gpCopyMultiTexSubImage3DEXT = (C.GPCOPYMULTITEXSUBIMAGE3DEXT)(getProcAddr("glCopyMultiTexSubImage3DEXT"))
	gpCopyNamedBufferSubData = (C.GPCOPYNAMEDBUFFERSUBDATA)(getProcAddr("glCopyNamedBufferSubData"))
	if gpCopyNamedBufferSubData == nil {
		return errors.New("glCopyNamedBufferSubData")
	}
	gpCopyPathNV = (C.GPCOPYPATHNV)(getProcAddr("glCopyPathNV"))
	gpCopyPixels = (C.GPCOPYPIXELS)(getProcAddr("glCopyPixels"))
	if gpCopyPixels == nil {
		return errors.New("glCopyPixels")
	}
	gpCopyTexImage1D = (C.GPCOPYTEXIMAGE1D)(getProcAddr("glCopyTexImage1D"))
	if gpCopyTexImage1D == nil {
		return errors.New("glCopyTexImage1D")
	}
	gpCopyTexImage1DEXT = (C.GPCOPYTEXIMAGE1DEXT)(getProcAddr("glCopyTexImage1DEXT"))
	gpCopyTexImage2D = (C.GPCOPYTEXIMAGE2D)(getProcAddr("glCopyTexImage2D"))
	if gpCopyTexImage2D == nil {
		return errors.New("glCopyTexImage2D")
	}
	gpCopyTexImage2DEXT = (C.GPCOPYTEXIMAGE2DEXT)(getProcAddr("glCopyTexImage2DEXT"))
	gpCopyTexSubImage1D = (C.GPCOPYTEXSUBIMAGE1D)(getProcAddr("glCopyTexSubImage1D"))
	if gpCopyTexSubImage1D == nil {
		return errors.New("glCopyTexSubImage1D")
	}
	gpCopyTexSubImage1DEXT = (C.GPCOPYTEXSUBIMAGE1DEXT)(getProcAddr("glCopyTexSubImage1DEXT"))
	gpCopyTexSubImage2D = (C.GPCOPYTEXSUBIMAGE2D)(getProcAddr("glCopyTexSubImage2D"))
	if gpCopyTexSubImage2D == nil {
		return errors.New("glCopyTexSubImage2D")
	}
	gpCopyTexSubImage2DEXT = (C.GPCOPYTEXSUBIMAGE2DEXT)(getProcAddr("glCopyTexSubImage2DEXT"))
	gpCopyTexSubImage3D = (C.GPCOPYTEXSUBIMAGE3D)(getProcAddr("glCopyTexSubImage3D"))
	if gpCopyTexSubImage3D == nil {
		return errors.New("glCopyTexSubImage3D")
	}
	gpCopyTexSubImage3DEXT = (C.GPCOPYTEXSUBIMAGE3DEXT)(getProcAddr("glCopyTexSubImage3DEXT"))
	gpCopyTextureImage1DEXT = (C.GPCOPYTEXTUREIMAGE1DEXT)(getProcAddr("glCopyTextureImage1DEXT"))
	gpCopyTextureImage2DEXT = (C.GPCOPYTEXTUREIMAGE2DEXT)(getProcAddr("glCopyTextureImage2DEXT"))
	gpCopyTextureSubImage1D = (C.GPCOPYTEXTURESUBIMAGE1D)(getProcAddr("glCopyTextureSubImage1D"))
	if gpCopyTextureSubImage1D == nil {
		return errors.New("glCopyTextureSubImage1D")
	}
	gpCopyTextureSubImage1DEXT = (C.GPCOPYTEXTURESUBIMAGE1DEXT)(getProcAddr("glCopyTextureSubImage1DEXT"))
	gpCopyTextureSubImage2D = (C.GPCOPYTEXTURESUBIMAGE2D)(getProcAddr("glCopyTextureSubImage2D"))
	if gpCopyTextureSubImage2D == nil {
		return errors.New("glCopyTextureSubImage2D")
	}
	gpCopyTextureSubImage2DEXT = (C.GPCOPYTEXTURESUBIMAGE2DEXT)(getProcAddr("glCopyTextureSubImage2DEXT"))
	gpCopyTextureSubImage3D = (C.GPCOPYTEXTURESUBIMAGE3D)(getProcAddr("glCopyTextureSubImage3D"))
	if gpCopyTextureSubImage3D == nil {
		return errors.New("glCopyTextureSubImage3D")
	}
	gpCopyTextureSubImage3DEXT = (C.GPCOPYTEXTURESUBIMAGE3DEXT)(getProcAddr("glCopyTextureSubImage3DEXT"))
	gpCoverFillPathInstancedNV = (C.GPCOVERFILLPATHINSTANCEDNV)(getProcAddr("glCoverFillPathInstancedNV"))
	gpCoverFillPathNV = (C.GPCOVERFILLPATHNV)(getProcAddr("glCoverFillPathNV"))
	gpCoverStrokePathInstancedNV = (C.GPCOVERSTROKEPATHINSTANCEDNV)(getProcAddr("glCoverStrokePathInstancedNV"))
	gpCoverStrokePathNV = (C.GPCOVERSTROKEPATHNV)(getProcAddr("glCoverStrokePathNV"))
	gpCoverageModulationNV = (C.GPCOVERAGEMODULATIONNV)(getProcAddr("glCoverageModulationNV"))
	gpCoverageModulationTableNV = (C.GPCOVERAGEMODULATIONTABLENV)(getProcAddr("glCoverageModulationTableNV"))
	gpCreateBuffers = (C.GPCREATEBUFFERS)(getProcAddr("glCreateBuffers"))
	if gpCreateBuffers == nil {
		return errors.New("glCreateBuffers")
	}
	gpCreateCommandListsNV = (C.GPCREATECOMMANDLISTSNV)(getProcAddr("glCreateCommandListsNV"))
	gpCreateFramebuffers = (C.GPCREATEFRAMEBUFFERS)(getProcAddr("glCreateFramebuffers"))
	if gpCreateFramebuffers == nil {
		return errors.New("glCreateFramebuffers")
	}
	gpCreateMemoryObjectsEXT = (C.GPCREATEMEMORYOBJECTSEXT)(getProcAddr("glCreateMemoryObjectsEXT"))
	gpCreatePerfQueryINTEL = (C.GPCREATEPERFQUERYINTEL)(getProcAddr("glCreatePerfQueryINTEL"))
	gpCreateProgram = (C.GPCREATEPROGRAM)(getProcAddr("glCreateProgram"))
	if gpCreateProgram == nil {
		return errors.New("glCreateProgram")
	}
	gpCreateProgramObjectARB = (C.GPCREATEPROGRAMOBJECTARB)(getProcAddr("glCreateProgramObjectARB"))
	gpCreateProgramPipelines = (C.GPCREATEPROGRAMPIPELINES)(getProcAddr("glCreateProgramPipelines"))
	if gpCreateProgramPipelines == nil {
		return errors.New("glCreateProgramPipelines")
	}
	gpCreateQueries = (C.GPCREATEQUERIES)(getProcAddr("glCreateQueries"))
	if gpCreateQueries == nil {
		return errors.New("glCreateQueries")
	}
	gpCreateRenderbuffers = (C.GPCREATERENDERBUFFERS)(getProcAddr("glCreateRenderbuffers"))
	if gpCreateRenderbuffers == nil {
		return errors.New("glCreateRenderbuffers")
	}
	gpCreateSamplers = (C.GPCREATESAMPLERS)(getProcAddr("glCreateSamplers"))
	if gpCreateSamplers == nil {
		return errors.New("glCreateSamplers")
	}
	gpCreateShader = (C.GPCREATESHADER)(getProcAddr("glCreateShader"))
	if gpCreateShader == nil {
		return errors.New("glCreateShader")
	}
	gpCreateShaderObjectARB = (C.GPCREATESHADEROBJECTARB)(getProcAddr("glCreateShaderObjectARB"))
	gpCreateShaderProgramEXT = (C.GPCREATESHADERPROGRAMEXT)(getProcAddr("glCreateShaderProgramEXT"))
	gpCreateShaderProgramv = (C.GPCREATESHADERPROGRAMV)(getProcAddr("glCreateShaderProgramv"))
	if gpCreateShaderProgramv == nil {
		return errors.New("glCreateShaderProgramv")
	}
	gpCreateShaderProgramvEXT = (C.GPCREATESHADERPROGRAMVEXT)(getProcAddr("glCreateShaderProgramvEXT"))
	gpCreateStatesNV = (C.GPCREATESTATESNV)(getProcAddr("glCreateStatesNV"))
	gpCreateSyncFromCLeventARB = (C.GPCREATESYNCFROMCLEVENTARB)(getProcAddr("glCreateSyncFromCLeventARB"))
	gpCreateTextures = (C.GPCREATETEXTURES)(getProcAddr("glCreateTextures"))
	if gpCreateTextures == nil {
		return errors.New("glCreateTextures")
	}
	gpCreateTransformFeedbacks = (C.GPCREATETRANSFORMFEEDBACKS)(getProcAddr("glCreateTransformFeedbacks"))
	if gpCreateTransformFeedbacks == nil {
		return errors.New("glCreateTransformFeedbacks")
	}
	gpCreateVertexArrays = (C.GPCREATEVERTEXARRAYS)(getProcAddr("glCreateVertexArrays"))
	if gpCreateVertexArrays == nil {
		return errors.New("glCreateVertexArrays")
	}
	gpCullFace = (C.GPCULLFACE)(getProcAddr("glCullFace"))
	if gpCullFace == nil {
		return errors.New("glCullFace")
	}
	gpCullParameterdvEXT = (C.GPCULLPARAMETERDVEXT)(getProcAddr("glCullParameterdvEXT"))
	gpCullParameterfvEXT = (C.GPCULLPARAMETERFVEXT)(getProcAddr("glCullParameterfvEXT"))
	gpCurrentPaletteMatrixARB = (C.GPCURRENTPALETTEMATRIXARB)(getProcAddr("glCurrentPaletteMatrixARB"))
	gpDebugMessageCallback = (C.GPDEBUGMESSAGECALLBACK)(getProcAddr("glDebugMessageCallback"))
	if gpDebugMessageCallback == nil {
		return errors.New("glDebugMessageCallback")
	}
	gpDebugMessageCallbackAMD = (C.GPDEBUGMESSAGECALLBACKAMD)(getProcAddr("glDebugMessageCallbackAMD"))
	gpDebugMessageCallbackARB = (C.GPDEBUGMESSAGECALLBACKARB)(getProcAddr("glDebugMessageCallbackARB"))
	gpDebugMessageCallbackKHR = (C.GPDEBUGMESSAGECALLBACKKHR)(getProcAddr("glDebugMessageCallbackKHR"))
	gpDebugMessageControl = (C.GPDEBUGMESSAGECONTROL)(getProcAddr("glDebugMessageControl"))
	if gpDebugMessageControl == nil {
		return errors.New("glDebugMessageControl")
	}
	gpDebugMessageControlARB = (C.GPDEBUGMESSAGECONTROLARB)(getProcAddr("glDebugMessageControlARB"))
	gpDebugMessageControlKHR = (C.GPDEBUGMESSAGECONTROLKHR)(getProcAddr("glDebugMessageControlKHR"))
	gpDebugMessageEnableAMD = (C.GPDEBUGMESSAGEENABLEAMD)(getProcAddr("glDebugMessageEnableAMD"))
	gpDebugMessageInsert = (C.GPDEBUGMESSAGEINSERT)(getProcAddr("glDebugMessageInsert"))
	if gpDebugMessageInsert == nil {
		return errors.New("glDebugMessageInsert")
	}
	gpDebugMessageInsertAMD = (C.GPDEBUGMESSAGEINSERTAMD)(getProcAddr("glDebugMessageInsertAMD"))
	gpDebugMessageInsertARB = (C.GPDEBUGMESSAGEINSERTARB)(getProcAddr("glDebugMessageInsertARB"))
	gpDebugMessageInsertKHR = (C.GPDEBUGMESSAGEINSERTKHR)(getProcAddr("glDebugMessageInsertKHR"))
	gpDeformSGIX = (C.GPDEFORMSGIX)(getProcAddr("glDeformSGIX"))
	gpDeformationMap3dSGIX = (C.GPDEFORMATIONMAP3DSGIX)(getProcAddr("glDeformationMap3dSGIX"))
	gpDeformationMap3fSGIX = (C.GPDEFORMATIONMAP3FSGIX)(getProcAddr("glDeformationMap3fSGIX"))
	gpDeleteAsyncMarkersSGIX = (C.GPDELETEASYNCMARKERSSGIX)(getProcAddr("glDeleteAsyncMarkersSGIX"))
	gpDeleteBuffers = (C.GPDELETEBUFFERS)(getProcAddr("glDeleteBuffers"))
	if gpDeleteBuffers == nil {
		return errors.New("glDeleteBuffers")
	}
	gpDeleteBuffersARB = (C.GPDELETEBUFFERSARB)(getProcAddr("glDeleteBuffersARB"))
	gpDeleteCommandListsNV = (C.GPDELETECOMMANDLISTSNV)(getProcAddr("glDeleteCommandListsNV"))
	gpDeleteFencesAPPLE = (C.GPDELETEFENCESAPPLE)(getProcAddr("glDeleteFencesAPPLE"))
	gpDeleteFencesNV = (C.GPDELETEFENCESNV)(getProcAddr("glDeleteFencesNV"))
	gpDeleteFragmentShaderATI = (C.GPDELETEFRAGMENTSHADERATI)(getProcAddr("glDeleteFragmentShaderATI"))
	gpDeleteFramebuffers = (C.GPDELETEFRAMEBUFFERS)(getProcAddr("glDeleteFramebuffers"))
	if gpDeleteFramebuffers == nil {
		return errors.New("glDeleteFramebuffers")
	}
	gpDeleteFramebuffersEXT = (C.GPDELETEFRAMEBUFFERSEXT)(getProcAddr("glDeleteFramebuffersEXT"))
	gpDeleteLists = (C.GPDELETELISTS)(getProcAddr("glDeleteLists"))
	if gpDeleteLists == nil {
		return errors.New("glDeleteLists")
	}
	gpDeleteMemoryObjectsEXT = (C.GPDELETEMEMORYOBJECTSEXT)(getProcAddr("glDeleteMemoryObjectsEXT"))
	gpDeleteNamedStringARB = (C.GPDELETENAMEDSTRINGARB)(getProcAddr("glDeleteNamedStringARB"))
	gpDeleteNamesAMD = (C.GPDELETENAMESAMD)(getProcAddr("glDeleteNamesAMD"))
	gpDeleteObjectARB = (C.GPDELETEOBJECTARB)(getProcAddr("glDeleteObjectARB"))
	gpDeleteOcclusionQueriesNV = (C.GPDELETEOCCLUSIONQUERIESNV)(getProcAddr("glDeleteOcclusionQueriesNV"))
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
	gpDeleteProgramsARB = (C.GPDELETEPROGRAMSARB)(getProcAddr("glDeleteProgramsARB"))
	gpDeleteProgramsNV = (C.GPDELETEPROGRAMSNV)(getProcAddr("glDeleteProgramsNV"))
	gpDeleteQueries = (C.GPDELETEQUERIES)(getProcAddr("glDeleteQueries"))
	if gpDeleteQueries == nil {
		return errors.New("glDeleteQueries")
	}
	gpDeleteQueriesARB = (C.GPDELETEQUERIESARB)(getProcAddr("glDeleteQueriesARB"))
	gpDeleteQueryResourceTagNV = (C.GPDELETEQUERYRESOURCETAGNV)(getProcAddr("glDeleteQueryResourceTagNV"))
	gpDeleteRenderbuffers = (C.GPDELETERENDERBUFFERS)(getProcAddr("glDeleteRenderbuffers"))
	if gpDeleteRenderbuffers == nil {
		return errors.New("glDeleteRenderbuffers")
	}
	gpDeleteRenderbuffersEXT = (C.GPDELETERENDERBUFFERSEXT)(getProcAddr("glDeleteRenderbuffersEXT"))
	gpDeleteSamplers = (C.GPDELETESAMPLERS)(getProcAddr("glDeleteSamplers"))
	if gpDeleteSamplers == nil {
		return errors.New("glDeleteSamplers")
	}
	gpDeleteSemaphoresEXT = (C.GPDELETESEMAPHORESEXT)(getProcAddr("glDeleteSemaphoresEXT"))
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
	gpDeleteTexturesEXT = (C.GPDELETETEXTURESEXT)(getProcAddr("glDeleteTexturesEXT"))
	gpDeleteTransformFeedbacks = (C.GPDELETETRANSFORMFEEDBACKS)(getProcAddr("glDeleteTransformFeedbacks"))
	if gpDeleteTransformFeedbacks == nil {
		return errors.New("glDeleteTransformFeedbacks")
	}
	gpDeleteTransformFeedbacksNV = (C.GPDELETETRANSFORMFEEDBACKSNV)(getProcAddr("glDeleteTransformFeedbacksNV"))
	gpDeleteVertexArrays = (C.GPDELETEVERTEXARRAYS)(getProcAddr("glDeleteVertexArrays"))
	if gpDeleteVertexArrays == nil {
		return errors.New("glDeleteVertexArrays")
	}
	gpDeleteVertexArraysAPPLE = (C.GPDELETEVERTEXARRAYSAPPLE)(getProcAddr("glDeleteVertexArraysAPPLE"))
	gpDeleteVertexShaderEXT = (C.GPDELETEVERTEXSHADEREXT)(getProcAddr("glDeleteVertexShaderEXT"))
	gpDepthBoundsEXT = (C.GPDEPTHBOUNDSEXT)(getProcAddr("glDepthBoundsEXT"))
	gpDepthBoundsdNV = (C.GPDEPTHBOUNDSDNV)(getProcAddr("glDepthBoundsdNV"))
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
	gpDepthRangedNV = (C.GPDEPTHRANGEDNV)(getProcAddr("glDepthRangedNV"))
	gpDepthRangef = (C.GPDEPTHRANGEF)(getProcAddr("glDepthRangef"))
	if gpDepthRangef == nil {
		return errors.New("glDepthRangef")
	}
	gpDepthRangefOES = (C.GPDEPTHRANGEFOES)(getProcAddr("glDepthRangefOES"))
	gpDepthRangexOES = (C.GPDEPTHRANGEXOES)(getProcAddr("glDepthRangexOES"))
	gpDetachObjectARB = (C.GPDETACHOBJECTARB)(getProcAddr("glDetachObjectARB"))
	gpDetachShader = (C.GPDETACHSHADER)(getProcAddr("glDetachShader"))
	if gpDetachShader == nil {
		return errors.New("glDetachShader")
	}
	gpDetailTexFuncSGIS = (C.GPDETAILTEXFUNCSGIS)(getProcAddr("glDetailTexFuncSGIS"))
	gpDisable = (C.GPDISABLE)(getProcAddr("glDisable"))
	if gpDisable == nil {
		return errors.New("glDisable")
	}
	gpDisableClientState = (C.GPDISABLECLIENTSTATE)(getProcAddr("glDisableClientState"))
	if gpDisableClientState == nil {
		return errors.New("glDisableClientState")
	}
	gpDisableClientStateIndexedEXT = (C.GPDISABLECLIENTSTATEINDEXEDEXT)(getProcAddr("glDisableClientStateIndexedEXT"))
	gpDisableClientStateiEXT = (C.GPDISABLECLIENTSTATEIEXT)(getProcAddr("glDisableClientStateiEXT"))
	gpDisableIndexedEXT = (C.GPDISABLEINDEXEDEXT)(getProcAddr("glDisableIndexedEXT"))
	gpDisableVariantClientStateEXT = (C.GPDISABLEVARIANTCLIENTSTATEEXT)(getProcAddr("glDisableVariantClientStateEXT"))
	gpDisableVertexArrayAttrib = (C.GPDISABLEVERTEXARRAYATTRIB)(getProcAddr("glDisableVertexArrayAttrib"))
	if gpDisableVertexArrayAttrib == nil {
		return errors.New("glDisableVertexArrayAttrib")
	}
	gpDisableVertexArrayAttribEXT = (C.GPDISABLEVERTEXARRAYATTRIBEXT)(getProcAddr("glDisableVertexArrayAttribEXT"))
	gpDisableVertexArrayEXT = (C.GPDISABLEVERTEXARRAYEXT)(getProcAddr("glDisableVertexArrayEXT"))
	gpDisableVertexAttribAPPLE = (C.GPDISABLEVERTEXATTRIBAPPLE)(getProcAddr("glDisableVertexAttribAPPLE"))
	gpDisableVertexAttribArray = (C.GPDISABLEVERTEXATTRIBARRAY)(getProcAddr("glDisableVertexAttribArray"))
	if gpDisableVertexAttribArray == nil {
		return errors.New("glDisableVertexAttribArray")
	}
	gpDisableVertexAttribArrayARB = (C.GPDISABLEVERTEXATTRIBARRAYARB)(getProcAddr("glDisableVertexAttribArrayARB"))
	gpDisablei = (C.GPDISABLEI)(getProcAddr("glDisablei"))
	if gpDisablei == nil {
		return errors.New("glDisablei")
	}
	gpDispatchCompute = (C.GPDISPATCHCOMPUTE)(getProcAddr("glDispatchCompute"))
	if gpDispatchCompute == nil {
		return errors.New("glDispatchCompute")
	}
	gpDispatchComputeGroupSizeARB = (C.GPDISPATCHCOMPUTEGROUPSIZEARB)(getProcAddr("glDispatchComputeGroupSizeARB"))
	gpDispatchComputeIndirect = (C.GPDISPATCHCOMPUTEINDIRECT)(getProcAddr("glDispatchComputeIndirect"))
	if gpDispatchComputeIndirect == nil {
		return errors.New("glDispatchComputeIndirect")
	}
	gpDrawArrays = (C.GPDRAWARRAYS)(getProcAddr("glDrawArrays"))
	if gpDrawArrays == nil {
		return errors.New("glDrawArrays")
	}
	gpDrawArraysEXT = (C.GPDRAWARRAYSEXT)(getProcAddr("glDrawArraysEXT"))
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
	if gpDrawArraysInstancedBaseInstance == nil {
		return errors.New("glDrawArraysInstancedBaseInstance")
	}
	gpDrawArraysInstancedEXT = (C.GPDRAWARRAYSINSTANCEDEXT)(getProcAddr("glDrawArraysInstancedEXT"))
	gpDrawBuffer = (C.GPDRAWBUFFER)(getProcAddr("glDrawBuffer"))
	if gpDrawBuffer == nil {
		return errors.New("glDrawBuffer")
	}
	gpDrawBuffers = (C.GPDRAWBUFFERS)(getProcAddr("glDrawBuffers"))
	if gpDrawBuffers == nil {
		return errors.New("glDrawBuffers")
	}
	gpDrawBuffersARB = (C.GPDRAWBUFFERSARB)(getProcAddr("glDrawBuffersARB"))
	gpDrawBuffersATI = (C.GPDRAWBUFFERSATI)(getProcAddr("glDrawBuffersATI"))
	gpDrawCommandsAddressNV = (C.GPDRAWCOMMANDSADDRESSNV)(getProcAddr("glDrawCommandsAddressNV"))
	gpDrawCommandsNV = (C.GPDRAWCOMMANDSNV)(getProcAddr("glDrawCommandsNV"))
	gpDrawCommandsStatesAddressNV = (C.GPDRAWCOMMANDSSTATESADDRESSNV)(getProcAddr("glDrawCommandsStatesAddressNV"))
	gpDrawCommandsStatesNV = (C.GPDRAWCOMMANDSSTATESNV)(getProcAddr("glDrawCommandsStatesNV"))
	gpDrawElementArrayAPPLE = (C.GPDRAWELEMENTARRAYAPPLE)(getProcAddr("glDrawElementArrayAPPLE"))
	gpDrawElementArrayATI = (C.GPDRAWELEMENTARRAYATI)(getProcAddr("glDrawElementArrayATI"))
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
	if gpDrawElementsInstancedBaseInstance == nil {
		return errors.New("glDrawElementsInstancedBaseInstance")
	}
	gpDrawElementsInstancedBaseVertex = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEX)(getProcAddr("glDrawElementsInstancedBaseVertex"))
	if gpDrawElementsInstancedBaseVertex == nil {
		return errors.New("glDrawElementsInstancedBaseVertex")
	}
	gpDrawElementsInstancedBaseVertexBaseInstance = (C.GPDRAWELEMENTSINSTANCEDBASEVERTEXBASEINSTANCE)(getProcAddr("glDrawElementsInstancedBaseVertexBaseInstance"))
	if gpDrawElementsInstancedBaseVertexBaseInstance == nil {
		return errors.New("glDrawElementsInstancedBaseVertexBaseInstance")
	}
	gpDrawElementsInstancedEXT = (C.GPDRAWELEMENTSINSTANCEDEXT)(getProcAddr("glDrawElementsInstancedEXT"))
	gpDrawMeshArraysSUN = (C.GPDRAWMESHARRAYSSUN)(getProcAddr("glDrawMeshArraysSUN"))
	gpDrawPixels = (C.GPDRAWPIXELS)(getProcAddr("glDrawPixels"))
	if gpDrawPixels == nil {
		return errors.New("glDrawPixels")
	}
	gpDrawRangeElementArrayAPPLE = (C.GPDRAWRANGEELEMENTARRAYAPPLE)(getProcAddr("glDrawRangeElementArrayAPPLE"))
	gpDrawRangeElementArrayATI = (C.GPDRAWRANGEELEMENTARRAYATI)(getProcAddr("glDrawRangeElementArrayATI"))
	gpDrawRangeElements = (C.GPDRAWRANGEELEMENTS)(getProcAddr("glDrawRangeElements"))
	if gpDrawRangeElements == nil {
		return errors.New("glDrawRangeElements")
	}
	gpDrawRangeElementsBaseVertex = (C.GPDRAWRANGEELEMENTSBASEVERTEX)(getProcAddr("glDrawRangeElementsBaseVertex"))
	if gpDrawRangeElementsBaseVertex == nil {
		return errors.New("glDrawRangeElementsBaseVertex")
	}
	gpDrawRangeElementsEXT = (C.GPDRAWRANGEELEMENTSEXT)(getProcAddr("glDrawRangeElementsEXT"))
	gpDrawTextureNV = (C.GPDRAWTEXTURENV)(getProcAddr("glDrawTextureNV"))
	gpDrawTransformFeedback = (C.GPDRAWTRANSFORMFEEDBACK)(getProcAddr("glDrawTransformFeedback"))
	if gpDrawTransformFeedback == nil {
		return errors.New("glDrawTransformFeedback")
	}
	gpDrawTransformFeedbackInstanced = (C.GPDRAWTRANSFORMFEEDBACKINSTANCED)(getProcAddr("glDrawTransformFeedbackInstanced"))
	if gpDrawTransformFeedbackInstanced == nil {
		return errors.New("glDrawTransformFeedbackInstanced")
	}
	gpDrawTransformFeedbackNV = (C.GPDRAWTRANSFORMFEEDBACKNV)(getProcAddr("glDrawTransformFeedbackNV"))
	gpDrawTransformFeedbackStream = (C.GPDRAWTRANSFORMFEEDBACKSTREAM)(getProcAddr("glDrawTransformFeedbackStream"))
	if gpDrawTransformFeedbackStream == nil {
		return errors.New("glDrawTransformFeedbackStream")
	}
	gpDrawTransformFeedbackStreamInstanced = (C.GPDRAWTRANSFORMFEEDBACKSTREAMINSTANCED)(getProcAddr("glDrawTransformFeedbackStreamInstanced"))
	if gpDrawTransformFeedbackStreamInstanced == nil {
		return errors.New("glDrawTransformFeedbackStreamInstanced")
	}
	gpDrawVkImageNV = (C.GPDRAWVKIMAGENV)(getProcAddr("glDrawVkImageNV"))
	gpEGLImageTargetTexStorageEXT = (C.GPEGLIMAGETARGETTEXSTORAGEEXT)(getProcAddr("glEGLImageTargetTexStorageEXT"))
	gpEGLImageTargetTextureStorageEXT = (C.GPEGLIMAGETARGETTEXTURESTORAGEEXT)(getProcAddr("glEGLImageTargetTextureStorageEXT"))
	gpEdgeFlag = (C.GPEDGEFLAG)(getProcAddr("glEdgeFlag"))
	if gpEdgeFlag == nil {
		return errors.New("glEdgeFlag")
	}
	gpEdgeFlagFormatNV = (C.GPEDGEFLAGFORMATNV)(getProcAddr("glEdgeFlagFormatNV"))
	gpEdgeFlagPointer = (C.GPEDGEFLAGPOINTER)(getProcAddr("glEdgeFlagPointer"))
	if gpEdgeFlagPointer == nil {
		return errors.New("glEdgeFlagPointer")
	}
	gpEdgeFlagPointerEXT = (C.GPEDGEFLAGPOINTEREXT)(getProcAddr("glEdgeFlagPointerEXT"))
	gpEdgeFlagPointerListIBM = (C.GPEDGEFLAGPOINTERLISTIBM)(getProcAddr("glEdgeFlagPointerListIBM"))
	gpEdgeFlagv = (C.GPEDGEFLAGV)(getProcAddr("glEdgeFlagv"))
	if gpEdgeFlagv == nil {
		return errors.New("glEdgeFlagv")
	}
	gpElementPointerAPPLE = (C.GPELEMENTPOINTERAPPLE)(getProcAddr("glElementPointerAPPLE"))
	gpElementPointerATI = (C.GPELEMENTPOINTERATI)(getProcAddr("glElementPointerATI"))
	gpEnable = (C.GPENABLE)(getProcAddr("glEnable"))
	if gpEnable == nil {
		return errors.New("glEnable")
	}
	gpEnableClientState = (C.GPENABLECLIENTSTATE)(getProcAddr("glEnableClientState"))
	if gpEnableClientState == nil {
		return errors.New("glEnableClientState")
	}
	gpEnableClientStateIndexedEXT = (C.GPENABLECLIENTSTATEINDEXEDEXT)(getProcAddr("glEnableClientStateIndexedEXT"))
	gpEnableClientStateiEXT = (C.GPENABLECLIENTSTATEIEXT)(getProcAddr("glEnableClientStateiEXT"))
	gpEnableIndexedEXT = (C.GPENABLEINDEXEDEXT)(getProcAddr("glEnableIndexedEXT"))
	gpEnableVariantClientStateEXT = (C.GPENABLEVARIANTCLIENTSTATEEXT)(getProcAddr("glEnableVariantClientStateEXT"))
	gpEnableVertexArrayAttrib = (C.GPENABLEVERTEXARRAYATTRIB)(getProcAddr("glEnableVertexArrayAttrib"))
	if gpEnableVertexArrayAttrib == nil {
		return errors.New("glEnableVertexArrayAttrib")
	}
	gpEnableVertexArrayAttribEXT = (C.GPENABLEVERTEXARRAYATTRIBEXT)(getProcAddr("glEnableVertexArrayAttribEXT"))
	gpEnableVertexArrayEXT = (C.GPENABLEVERTEXARRAYEXT)(getProcAddr("glEnableVertexArrayEXT"))
	gpEnableVertexAttribAPPLE = (C.GPENABLEVERTEXATTRIBAPPLE)(getProcAddr("glEnableVertexAttribAPPLE"))
	gpEnableVertexAttribArray = (C.GPENABLEVERTEXATTRIBARRAY)(getProcAddr("glEnableVertexAttribArray"))
	if gpEnableVertexAttribArray == nil {
		return errors.New("glEnableVertexAttribArray")
	}
	gpEnableVertexAttribArrayARB = (C.GPENABLEVERTEXATTRIBARRAYARB)(getProcAddr("glEnableVertexAttribArrayARB"))
	gpEnablei = (C.GPENABLEI)(getProcAddr("glEnablei"))
	if gpEnablei == nil {
		return errors.New("glEnablei")
	}
	gpEnd = (C.GPEND)(getProcAddr("glEnd"))
	if gpEnd == nil {
		return errors.New("glEnd")
	}
	gpEndConditionalRender = (C.GPENDCONDITIONALRENDER)(getProcAddr("glEndConditionalRender"))
	if gpEndConditionalRender == nil {
		return errors.New("glEndConditionalRender")
	}
	gpEndConditionalRenderNV = (C.GPENDCONDITIONALRENDERNV)(getProcAddr("glEndConditionalRenderNV"))
	gpEndConditionalRenderNVX = (C.GPENDCONDITIONALRENDERNVX)(getProcAddr("glEndConditionalRenderNVX"))
	gpEndFragmentShaderATI = (C.GPENDFRAGMENTSHADERATI)(getProcAddr("glEndFragmentShaderATI"))
	gpEndList = (C.GPENDLIST)(getProcAddr("glEndList"))
	if gpEndList == nil {
		return errors.New("glEndList")
	}
	gpEndOcclusionQueryNV = (C.GPENDOCCLUSIONQUERYNV)(getProcAddr("glEndOcclusionQueryNV"))
	gpEndPerfMonitorAMD = (C.GPENDPERFMONITORAMD)(getProcAddr("glEndPerfMonitorAMD"))
	gpEndPerfQueryINTEL = (C.GPENDPERFQUERYINTEL)(getProcAddr("glEndPerfQueryINTEL"))
	gpEndQuery = (C.GPENDQUERY)(getProcAddr("glEndQuery"))
	if gpEndQuery == nil {
		return errors.New("glEndQuery")
	}
	gpEndQueryARB = (C.GPENDQUERYARB)(getProcAddr("glEndQueryARB"))
	gpEndQueryIndexed = (C.GPENDQUERYINDEXED)(getProcAddr("glEndQueryIndexed"))
	if gpEndQueryIndexed == nil {
		return errors.New("glEndQueryIndexed")
	}
	gpEndTransformFeedback = (C.GPENDTRANSFORMFEEDBACK)(getProcAddr("glEndTransformFeedback"))
	if gpEndTransformFeedback == nil {
		return errors.New("glEndTransformFeedback")
	}
	gpEndTransformFeedbackEXT = (C.GPENDTRANSFORMFEEDBACKEXT)(getProcAddr("glEndTransformFeedbackEXT"))
	gpEndTransformFeedbackNV = (C.GPENDTRANSFORMFEEDBACKNV)(getProcAddr("glEndTransformFeedbackNV"))
	gpEndVertexShaderEXT = (C.GPENDVERTEXSHADEREXT)(getProcAddr("glEndVertexShaderEXT"))
	gpEndVideoCaptureNV = (C.GPENDVIDEOCAPTURENV)(getProcAddr("glEndVideoCaptureNV"))
	gpEvalCoord1d = (C.GPEVALCOORD1D)(getProcAddr("glEvalCoord1d"))
	if gpEvalCoord1d == nil {
		return errors.New("glEvalCoord1d")
	}
	gpEvalCoord1dv = (C.GPEVALCOORD1DV)(getProcAddr("glEvalCoord1dv"))
	if gpEvalCoord1dv == nil {
		return errors.New("glEvalCoord1dv")
	}
	gpEvalCoord1f = (C.GPEVALCOORD1F)(getProcAddr("glEvalCoord1f"))
	if gpEvalCoord1f == nil {
		return errors.New("glEvalCoord1f")
	}
	gpEvalCoord1fv = (C.GPEVALCOORD1FV)(getProcAddr("glEvalCoord1fv"))
	if gpEvalCoord1fv == nil {
		return errors.New("glEvalCoord1fv")
	}
	gpEvalCoord1xOES = (C.GPEVALCOORD1XOES)(getProcAddr("glEvalCoord1xOES"))
	gpEvalCoord1xvOES = (C.GPEVALCOORD1XVOES)(getProcAddr("glEvalCoord1xvOES"))
	gpEvalCoord2d = (C.GPEVALCOORD2D)(getProcAddr("glEvalCoord2d"))
	if gpEvalCoord2d == nil {
		return errors.New("glEvalCoord2d")
	}
	gpEvalCoord2dv = (C.GPEVALCOORD2DV)(getProcAddr("glEvalCoord2dv"))
	if gpEvalCoord2dv == nil {
		return errors.New("glEvalCoord2dv")
	}
	gpEvalCoord2f = (C.GPEVALCOORD2F)(getProcAddr("glEvalCoord2f"))
	if gpEvalCoord2f == nil {
		return errors.New("glEvalCoord2f")
	}
	gpEvalCoord2fv = (C.GPEVALCOORD2FV)(getProcAddr("glEvalCoord2fv"))
	if gpEvalCoord2fv == nil {
		return errors.New("glEvalCoord2fv")
	}
	gpEvalCoord2xOES = (C.GPEVALCOORD2XOES)(getProcAddr("glEvalCoord2xOES"))
	gpEvalCoord2xvOES = (C.GPEVALCOORD2XVOES)(getProcAddr("glEvalCoord2xvOES"))
	gpEvalMapsNV = (C.GPEVALMAPSNV)(getProcAddr("glEvalMapsNV"))
	gpEvalMesh1 = (C.GPEVALMESH1)(getProcAddr("glEvalMesh1"))
	if gpEvalMesh1 == nil {
		return errors.New("glEvalMesh1")
	}
	gpEvalMesh2 = (C.GPEVALMESH2)(getProcAddr("glEvalMesh2"))
	if gpEvalMesh2 == nil {
		return errors.New("glEvalMesh2")
	}
	gpEvalPoint1 = (C.GPEVALPOINT1)(getProcAddr("glEvalPoint1"))
	if gpEvalPoint1 == nil {
		return errors.New("glEvalPoint1")
	}
	gpEvalPoint2 = (C.GPEVALPOINT2)(getProcAddr("glEvalPoint2"))
	if gpEvalPoint2 == nil {
		return errors.New("glEvalPoint2")
	}
	gpEvaluateDepthValuesARB = (C.GPEVALUATEDEPTHVALUESARB)(getProcAddr("glEvaluateDepthValuesARB"))
	gpExecuteProgramNV = (C.GPEXECUTEPROGRAMNV)(getProcAddr("glExecuteProgramNV"))
	gpExtractComponentEXT = (C.GPEXTRACTCOMPONENTEXT)(getProcAddr("glExtractComponentEXT"))
	gpFeedbackBuffer = (C.GPFEEDBACKBUFFER)(getProcAddr("glFeedbackBuffer"))
	if gpFeedbackBuffer == nil {
		return errors.New("glFeedbackBuffer")
	}
	gpFeedbackBufferxOES = (C.GPFEEDBACKBUFFERXOES)(getProcAddr("glFeedbackBufferxOES"))
	gpFenceSync = (C.GPFENCESYNC)(getProcAddr("glFenceSync"))
	if gpFenceSync == nil {
		return errors.New("glFenceSync")
	}
	gpFinalCombinerInputNV = (C.GPFINALCOMBINERINPUTNV)(getProcAddr("glFinalCombinerInputNV"))
	gpFinish = (C.GPFINISH)(getProcAddr("glFinish"))
	if gpFinish == nil {
		return errors.New("glFinish")
	}
	gpFinishAsyncSGIX = (C.GPFINISHASYNCSGIX)(getProcAddr("glFinishAsyncSGIX"))
	gpFinishFenceAPPLE = (C.GPFINISHFENCEAPPLE)(getProcAddr("glFinishFenceAPPLE"))
	gpFinishFenceNV = (C.GPFINISHFENCENV)(getProcAddr("glFinishFenceNV"))
	gpFinishObjectAPPLE = (C.GPFINISHOBJECTAPPLE)(getProcAddr("glFinishObjectAPPLE"))
	gpFinishTextureSUNX = (C.GPFINISHTEXTURESUNX)(getProcAddr("glFinishTextureSUNX"))
	gpFlush = (C.GPFLUSH)(getProcAddr("glFlush"))
	if gpFlush == nil {
		return errors.New("glFlush")
	}
	gpFlushMappedBufferRange = (C.GPFLUSHMAPPEDBUFFERRANGE)(getProcAddr("glFlushMappedBufferRange"))
	if gpFlushMappedBufferRange == nil {
		return errors.New("glFlushMappedBufferRange")
	}
	gpFlushMappedBufferRangeAPPLE = (C.GPFLUSHMAPPEDBUFFERRANGEAPPLE)(getProcAddr("glFlushMappedBufferRangeAPPLE"))
	gpFlushMappedNamedBufferRange = (C.GPFLUSHMAPPEDNAMEDBUFFERRANGE)(getProcAddr("glFlushMappedNamedBufferRange"))
	if gpFlushMappedNamedBufferRange == nil {
		return errors.New("glFlushMappedNamedBufferRange")
	}
	gpFlushMappedNamedBufferRangeEXT = (C.GPFLUSHMAPPEDNAMEDBUFFERRANGEEXT)(getProcAddr("glFlushMappedNamedBufferRangeEXT"))
	gpFlushPixelDataRangeNV = (C.GPFLUSHPIXELDATARANGENV)(getProcAddr("glFlushPixelDataRangeNV"))
	gpFlushRasterSGIX = (C.GPFLUSHRASTERSGIX)(getProcAddr("glFlushRasterSGIX"))
	gpFlushStaticDataIBM = (C.GPFLUSHSTATICDATAIBM)(getProcAddr("glFlushStaticDataIBM"))
	gpFlushVertexArrayRangeAPPLE = (C.GPFLUSHVERTEXARRAYRANGEAPPLE)(getProcAddr("glFlushVertexArrayRangeAPPLE"))
	gpFlushVertexArrayRangeNV = (C.GPFLUSHVERTEXARRAYRANGENV)(getProcAddr("glFlushVertexArrayRangeNV"))
	gpFogCoordFormatNV = (C.GPFOGCOORDFORMATNV)(getProcAddr("glFogCoordFormatNV"))
	gpFogCoordPointer = (C.GPFOGCOORDPOINTER)(getProcAddr("glFogCoordPointer"))
	if gpFogCoordPointer == nil {
		return errors.New("glFogCoordPointer")
	}
	gpFogCoordPointerEXT = (C.GPFOGCOORDPOINTEREXT)(getProcAddr("glFogCoordPointerEXT"))
	gpFogCoordPointerListIBM = (C.GPFOGCOORDPOINTERLISTIBM)(getProcAddr("glFogCoordPointerListIBM"))
	gpFogCoordd = (C.GPFOGCOORDD)(getProcAddr("glFogCoordd"))
	if gpFogCoordd == nil {
		return errors.New("glFogCoordd")
	}
	gpFogCoorddEXT = (C.GPFOGCOORDDEXT)(getProcAddr("glFogCoorddEXT"))
	gpFogCoorddv = (C.GPFOGCOORDDV)(getProcAddr("glFogCoorddv"))
	if gpFogCoorddv == nil {
		return errors.New("glFogCoorddv")
	}
	gpFogCoorddvEXT = (C.GPFOGCOORDDVEXT)(getProcAddr("glFogCoorddvEXT"))
	gpFogCoordf = (C.GPFOGCOORDF)(getProcAddr("glFogCoordf"))
	if gpFogCoordf == nil {
		return errors.New("glFogCoordf")
	}
	gpFogCoordfEXT = (C.GPFOGCOORDFEXT)(getProcAddr("glFogCoordfEXT"))
	gpFogCoordfv = (C.GPFOGCOORDFV)(getProcAddr("glFogCoordfv"))
	if gpFogCoordfv == nil {
		return errors.New("glFogCoordfv")
	}
	gpFogCoordfvEXT = (C.GPFOGCOORDFVEXT)(getProcAddr("glFogCoordfvEXT"))
	gpFogCoordhNV = (C.GPFOGCOORDHNV)(getProcAddr("glFogCoordhNV"))
	gpFogCoordhvNV = (C.GPFOGCOORDHVNV)(getProcAddr("glFogCoordhvNV"))
	gpFogFuncSGIS = (C.GPFOGFUNCSGIS)(getProcAddr("glFogFuncSGIS"))
	gpFogf = (C.GPFOGF)(getProcAddr("glFogf"))
	if gpFogf == nil {
		return errors.New("glFogf")
	}
	gpFogfv = (C.GPFOGFV)(getProcAddr("glFogfv"))
	if gpFogfv == nil {
		return errors.New("glFogfv")
	}
	gpFogi = (C.GPFOGI)(getProcAddr("glFogi"))
	if gpFogi == nil {
		return errors.New("glFogi")
	}
	gpFogiv = (C.GPFOGIV)(getProcAddr("glFogiv"))
	if gpFogiv == nil {
		return errors.New("glFogiv")
	}
	gpFogxOES = (C.GPFOGXOES)(getProcAddr("glFogxOES"))
	gpFogxvOES = (C.GPFOGXVOES)(getProcAddr("glFogxvOES"))
	gpFragmentColorMaterialSGIX = (C.GPFRAGMENTCOLORMATERIALSGIX)(getProcAddr("glFragmentColorMaterialSGIX"))
	gpFragmentCoverageColorNV = (C.GPFRAGMENTCOVERAGECOLORNV)(getProcAddr("glFragmentCoverageColorNV"))
	gpFragmentLightModelfSGIX = (C.GPFRAGMENTLIGHTMODELFSGIX)(getProcAddr("glFragmentLightModelfSGIX"))
	gpFragmentLightModelfvSGIX = (C.GPFRAGMENTLIGHTMODELFVSGIX)(getProcAddr("glFragmentLightModelfvSGIX"))
	gpFragmentLightModeliSGIX = (C.GPFRAGMENTLIGHTMODELISGIX)(getProcAddr("glFragmentLightModeliSGIX"))
	gpFragmentLightModelivSGIX = (C.GPFRAGMENTLIGHTMODELIVSGIX)(getProcAddr("glFragmentLightModelivSGIX"))
	gpFragmentLightfSGIX = (C.GPFRAGMENTLIGHTFSGIX)(getProcAddr("glFragmentLightfSGIX"))
	gpFragmentLightfvSGIX = (C.GPFRAGMENTLIGHTFVSGIX)(getProcAddr("glFragmentLightfvSGIX"))
	gpFragmentLightiSGIX = (C.GPFRAGMENTLIGHTISGIX)(getProcAddr("glFragmentLightiSGIX"))
	gpFragmentLightivSGIX = (C.GPFRAGMENTLIGHTIVSGIX)(getProcAddr("glFragmentLightivSGIX"))
	gpFragmentMaterialfSGIX = (C.GPFRAGMENTMATERIALFSGIX)(getProcAddr("glFragmentMaterialfSGIX"))
	gpFragmentMaterialfvSGIX = (C.GPFRAGMENTMATERIALFVSGIX)(getProcAddr("glFragmentMaterialfvSGIX"))
	gpFragmentMaterialiSGIX = (C.GPFRAGMENTMATERIALISGIX)(getProcAddr("glFragmentMaterialiSGIX"))
	gpFragmentMaterialivSGIX = (C.GPFRAGMENTMATERIALIVSGIX)(getProcAddr("glFragmentMaterialivSGIX"))
	gpFrameTerminatorGREMEDY = (C.GPFRAMETERMINATORGREMEDY)(getProcAddr("glFrameTerminatorGREMEDY"))
	gpFrameZoomSGIX = (C.GPFRAMEZOOMSGIX)(getProcAddr("glFrameZoomSGIX"))
	gpFramebufferDrawBufferEXT = (C.GPFRAMEBUFFERDRAWBUFFEREXT)(getProcAddr("glFramebufferDrawBufferEXT"))
	gpFramebufferDrawBuffersEXT = (C.GPFRAMEBUFFERDRAWBUFFERSEXT)(getProcAddr("glFramebufferDrawBuffersEXT"))
	gpFramebufferFetchBarrierEXT = (C.GPFRAMEBUFFERFETCHBARRIEREXT)(getProcAddr("glFramebufferFetchBarrierEXT"))
	gpFramebufferParameteri = (C.GPFRAMEBUFFERPARAMETERI)(getProcAddr("glFramebufferParameteri"))
	if gpFramebufferParameteri == nil {
		return errors.New("glFramebufferParameteri")
	}
	gpFramebufferReadBufferEXT = (C.GPFRAMEBUFFERREADBUFFEREXT)(getProcAddr("glFramebufferReadBufferEXT"))
	gpFramebufferRenderbuffer = (C.GPFRAMEBUFFERRENDERBUFFER)(getProcAddr("glFramebufferRenderbuffer"))
	if gpFramebufferRenderbuffer == nil {
		return errors.New("glFramebufferRenderbuffer")
	}
	gpFramebufferRenderbufferEXT = (C.GPFRAMEBUFFERRENDERBUFFEREXT)(getProcAddr("glFramebufferRenderbufferEXT"))
	gpFramebufferSampleLocationsfvARB = (C.GPFRAMEBUFFERSAMPLELOCATIONSFVARB)(getProcAddr("glFramebufferSampleLocationsfvARB"))
	gpFramebufferSampleLocationsfvNV = (C.GPFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glFramebufferSampleLocationsfvNV"))
	gpFramebufferSamplePositionsfvAMD = (C.GPFRAMEBUFFERSAMPLEPOSITIONSFVAMD)(getProcAddr("glFramebufferSamplePositionsfvAMD"))
	gpFramebufferTexture = (C.GPFRAMEBUFFERTEXTURE)(getProcAddr("glFramebufferTexture"))
	if gpFramebufferTexture == nil {
		return errors.New("glFramebufferTexture")
	}
	gpFramebufferTexture1D = (C.GPFRAMEBUFFERTEXTURE1D)(getProcAddr("glFramebufferTexture1D"))
	if gpFramebufferTexture1D == nil {
		return errors.New("glFramebufferTexture1D")
	}
	gpFramebufferTexture1DEXT = (C.GPFRAMEBUFFERTEXTURE1DEXT)(getProcAddr("glFramebufferTexture1DEXT"))
	gpFramebufferTexture2D = (C.GPFRAMEBUFFERTEXTURE2D)(getProcAddr("glFramebufferTexture2D"))
	if gpFramebufferTexture2D == nil {
		return errors.New("glFramebufferTexture2D")
	}
	gpFramebufferTexture2DEXT = (C.GPFRAMEBUFFERTEXTURE2DEXT)(getProcAddr("glFramebufferTexture2DEXT"))
	gpFramebufferTexture3D = (C.GPFRAMEBUFFERTEXTURE3D)(getProcAddr("glFramebufferTexture3D"))
	if gpFramebufferTexture3D == nil {
		return errors.New("glFramebufferTexture3D")
	}
	gpFramebufferTexture3DEXT = (C.GPFRAMEBUFFERTEXTURE3DEXT)(getProcAddr("glFramebufferTexture3DEXT"))
	gpFramebufferTextureARB = (C.GPFRAMEBUFFERTEXTUREARB)(getProcAddr("glFramebufferTextureARB"))
	gpFramebufferTextureEXT = (C.GPFRAMEBUFFERTEXTUREEXT)(getProcAddr("glFramebufferTextureEXT"))
	gpFramebufferTextureFaceARB = (C.GPFRAMEBUFFERTEXTUREFACEARB)(getProcAddr("glFramebufferTextureFaceARB"))
	gpFramebufferTextureFaceEXT = (C.GPFRAMEBUFFERTEXTUREFACEEXT)(getProcAddr("glFramebufferTextureFaceEXT"))
	gpFramebufferTextureLayer = (C.GPFRAMEBUFFERTEXTURELAYER)(getProcAddr("glFramebufferTextureLayer"))
	if gpFramebufferTextureLayer == nil {
		return errors.New("glFramebufferTextureLayer")
	}
	gpFramebufferTextureLayerARB = (C.GPFRAMEBUFFERTEXTURELAYERARB)(getProcAddr("glFramebufferTextureLayerARB"))
	gpFramebufferTextureLayerEXT = (C.GPFRAMEBUFFERTEXTURELAYEREXT)(getProcAddr("glFramebufferTextureLayerEXT"))
	gpFramebufferTextureMultiviewOVR = (C.GPFRAMEBUFFERTEXTUREMULTIVIEWOVR)(getProcAddr("glFramebufferTextureMultiviewOVR"))
	gpFreeObjectBufferATI = (C.GPFREEOBJECTBUFFERATI)(getProcAddr("glFreeObjectBufferATI"))
	gpFrontFace = (C.GPFRONTFACE)(getProcAddr("glFrontFace"))
	if gpFrontFace == nil {
		return errors.New("glFrontFace")
	}
	gpFrustum = (C.GPFRUSTUM)(getProcAddr("glFrustum"))
	if gpFrustum == nil {
		return errors.New("glFrustum")
	}
	gpFrustumfOES = (C.GPFRUSTUMFOES)(getProcAddr("glFrustumfOES"))
	gpFrustumxOES = (C.GPFRUSTUMXOES)(getProcAddr("glFrustumxOES"))
	gpGenAsyncMarkersSGIX = (C.GPGENASYNCMARKERSSGIX)(getProcAddr("glGenAsyncMarkersSGIX"))
	gpGenBuffers = (C.GPGENBUFFERS)(getProcAddr("glGenBuffers"))
	if gpGenBuffers == nil {
		return errors.New("glGenBuffers")
	}
	gpGenBuffersARB = (C.GPGENBUFFERSARB)(getProcAddr("glGenBuffersARB"))
	gpGenFencesAPPLE = (C.GPGENFENCESAPPLE)(getProcAddr("glGenFencesAPPLE"))
	gpGenFencesNV = (C.GPGENFENCESNV)(getProcAddr("glGenFencesNV"))
	gpGenFragmentShadersATI = (C.GPGENFRAGMENTSHADERSATI)(getProcAddr("glGenFragmentShadersATI"))
	gpGenFramebuffers = (C.GPGENFRAMEBUFFERS)(getProcAddr("glGenFramebuffers"))
	if gpGenFramebuffers == nil {
		return errors.New("glGenFramebuffers")
	}
	gpGenFramebuffersEXT = (C.GPGENFRAMEBUFFERSEXT)(getProcAddr("glGenFramebuffersEXT"))
	gpGenLists = (C.GPGENLISTS)(getProcAddr("glGenLists"))
	if gpGenLists == nil {
		return errors.New("glGenLists")
	}
	gpGenNamesAMD = (C.GPGENNAMESAMD)(getProcAddr("glGenNamesAMD"))
	gpGenOcclusionQueriesNV = (C.GPGENOCCLUSIONQUERIESNV)(getProcAddr("glGenOcclusionQueriesNV"))
	gpGenPathsNV = (C.GPGENPATHSNV)(getProcAddr("glGenPathsNV"))
	gpGenPerfMonitorsAMD = (C.GPGENPERFMONITORSAMD)(getProcAddr("glGenPerfMonitorsAMD"))
	gpGenProgramPipelines = (C.GPGENPROGRAMPIPELINES)(getProcAddr("glGenProgramPipelines"))
	if gpGenProgramPipelines == nil {
		return errors.New("glGenProgramPipelines")
	}
	gpGenProgramPipelinesEXT = (C.GPGENPROGRAMPIPELINESEXT)(getProcAddr("glGenProgramPipelinesEXT"))
	gpGenProgramsARB = (C.GPGENPROGRAMSARB)(getProcAddr("glGenProgramsARB"))
	gpGenProgramsNV = (C.GPGENPROGRAMSNV)(getProcAddr("glGenProgramsNV"))
	gpGenQueries = (C.GPGENQUERIES)(getProcAddr("glGenQueries"))
	if gpGenQueries == nil {
		return errors.New("glGenQueries")
	}
	gpGenQueriesARB = (C.GPGENQUERIESARB)(getProcAddr("glGenQueriesARB"))
	gpGenQueryResourceTagNV = (C.GPGENQUERYRESOURCETAGNV)(getProcAddr("glGenQueryResourceTagNV"))
	gpGenRenderbuffers = (C.GPGENRENDERBUFFERS)(getProcAddr("glGenRenderbuffers"))
	if gpGenRenderbuffers == nil {
		return errors.New("glGenRenderbuffers")
	}
	gpGenRenderbuffersEXT = (C.GPGENRENDERBUFFERSEXT)(getProcAddr("glGenRenderbuffersEXT"))
	gpGenSamplers = (C.GPGENSAMPLERS)(getProcAddr("glGenSamplers"))
	if gpGenSamplers == nil {
		return errors.New("glGenSamplers")
	}
	gpGenSemaphoresEXT = (C.GPGENSEMAPHORESEXT)(getProcAddr("glGenSemaphoresEXT"))
	gpGenSymbolsEXT = (C.GPGENSYMBOLSEXT)(getProcAddr("glGenSymbolsEXT"))
	gpGenTextures = (C.GPGENTEXTURES)(getProcAddr("glGenTextures"))
	if gpGenTextures == nil {
		return errors.New("glGenTextures")
	}
	gpGenTexturesEXT = (C.GPGENTEXTURESEXT)(getProcAddr("glGenTexturesEXT"))
	gpGenTransformFeedbacks = (C.GPGENTRANSFORMFEEDBACKS)(getProcAddr("glGenTransformFeedbacks"))
	if gpGenTransformFeedbacks == nil {
		return errors.New("glGenTransformFeedbacks")
	}
	gpGenTransformFeedbacksNV = (C.GPGENTRANSFORMFEEDBACKSNV)(getProcAddr("glGenTransformFeedbacksNV"))
	gpGenVertexArrays = (C.GPGENVERTEXARRAYS)(getProcAddr("glGenVertexArrays"))
	if gpGenVertexArrays == nil {
		return errors.New("glGenVertexArrays")
	}
	gpGenVertexArraysAPPLE = (C.GPGENVERTEXARRAYSAPPLE)(getProcAddr("glGenVertexArraysAPPLE"))
	gpGenVertexShadersEXT = (C.GPGENVERTEXSHADERSEXT)(getProcAddr("glGenVertexShadersEXT"))
	gpGenerateMipmap = (C.GPGENERATEMIPMAP)(getProcAddr("glGenerateMipmap"))
	if gpGenerateMipmap == nil {
		return errors.New("glGenerateMipmap")
	}
	gpGenerateMipmapEXT = (C.GPGENERATEMIPMAPEXT)(getProcAddr("glGenerateMipmapEXT"))
	gpGenerateMultiTexMipmapEXT = (C.GPGENERATEMULTITEXMIPMAPEXT)(getProcAddr("glGenerateMultiTexMipmapEXT"))
	gpGenerateTextureMipmap = (C.GPGENERATETEXTUREMIPMAP)(getProcAddr("glGenerateTextureMipmap"))
	if gpGenerateTextureMipmap == nil {
		return errors.New("glGenerateTextureMipmap")
	}
	gpGenerateTextureMipmapEXT = (C.GPGENERATETEXTUREMIPMAPEXT)(getProcAddr("glGenerateTextureMipmapEXT"))
	gpGetActiveAtomicCounterBufferiv = (C.GPGETACTIVEATOMICCOUNTERBUFFERIV)(getProcAddr("glGetActiveAtomicCounterBufferiv"))
	if gpGetActiveAtomicCounterBufferiv == nil {
		return errors.New("glGetActiveAtomicCounterBufferiv")
	}
	gpGetActiveAttrib = (C.GPGETACTIVEATTRIB)(getProcAddr("glGetActiveAttrib"))
	if gpGetActiveAttrib == nil {
		return errors.New("glGetActiveAttrib")
	}
	gpGetActiveAttribARB = (C.GPGETACTIVEATTRIBARB)(getProcAddr("glGetActiveAttribARB"))
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
	gpGetActiveUniformARB = (C.GPGETACTIVEUNIFORMARB)(getProcAddr("glGetActiveUniformARB"))
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
	gpGetActiveVaryingNV = (C.GPGETACTIVEVARYINGNV)(getProcAddr("glGetActiveVaryingNV"))
	gpGetArrayObjectfvATI = (C.GPGETARRAYOBJECTFVATI)(getProcAddr("glGetArrayObjectfvATI"))
	gpGetArrayObjectivATI = (C.GPGETARRAYOBJECTIVATI)(getProcAddr("glGetArrayObjectivATI"))
	gpGetAttachedObjectsARB = (C.GPGETATTACHEDOBJECTSARB)(getProcAddr("glGetAttachedObjectsARB"))
	gpGetAttachedShaders = (C.GPGETATTACHEDSHADERS)(getProcAddr("glGetAttachedShaders"))
	if gpGetAttachedShaders == nil {
		return errors.New("glGetAttachedShaders")
	}
	gpGetAttribLocation = (C.GPGETATTRIBLOCATION)(getProcAddr("glGetAttribLocation"))
	if gpGetAttribLocation == nil {
		return errors.New("glGetAttribLocation")
	}
	gpGetAttribLocationARB = (C.GPGETATTRIBLOCATIONARB)(getProcAddr("glGetAttribLocationARB"))
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
	gpGetBufferParameterivARB = (C.GPGETBUFFERPARAMETERIVARB)(getProcAddr("glGetBufferParameterivARB"))
	gpGetBufferParameterui64vNV = (C.GPGETBUFFERPARAMETERUI64VNV)(getProcAddr("glGetBufferParameterui64vNV"))
	gpGetBufferPointerv = (C.GPGETBUFFERPOINTERV)(getProcAddr("glGetBufferPointerv"))
	if gpGetBufferPointerv == nil {
		return errors.New("glGetBufferPointerv")
	}
	gpGetBufferPointervARB = (C.GPGETBUFFERPOINTERVARB)(getProcAddr("glGetBufferPointervARB"))
	gpGetBufferSubData = (C.GPGETBUFFERSUBDATA)(getProcAddr("glGetBufferSubData"))
	if gpGetBufferSubData == nil {
		return errors.New("glGetBufferSubData")
	}
	gpGetBufferSubDataARB = (C.GPGETBUFFERSUBDATAARB)(getProcAddr("glGetBufferSubDataARB"))
	gpGetClipPlane = (C.GPGETCLIPPLANE)(getProcAddr("glGetClipPlane"))
	if gpGetClipPlane == nil {
		return errors.New("glGetClipPlane")
	}
	gpGetClipPlanefOES = (C.GPGETCLIPPLANEFOES)(getProcAddr("glGetClipPlanefOES"))
	gpGetClipPlanexOES = (C.GPGETCLIPPLANEXOES)(getProcAddr("glGetClipPlanexOES"))
	gpGetColorTable = (C.GPGETCOLORTABLE)(getProcAddr("glGetColorTable"))
	gpGetColorTableEXT = (C.GPGETCOLORTABLEEXT)(getProcAddr("glGetColorTableEXT"))
	gpGetColorTableParameterfv = (C.GPGETCOLORTABLEPARAMETERFV)(getProcAddr("glGetColorTableParameterfv"))
	gpGetColorTableParameterfvEXT = (C.GPGETCOLORTABLEPARAMETERFVEXT)(getProcAddr("glGetColorTableParameterfvEXT"))
	gpGetColorTableParameterfvSGI = (C.GPGETCOLORTABLEPARAMETERFVSGI)(getProcAddr("glGetColorTableParameterfvSGI"))
	gpGetColorTableParameteriv = (C.GPGETCOLORTABLEPARAMETERIV)(getProcAddr("glGetColorTableParameteriv"))
	gpGetColorTableParameterivEXT = (C.GPGETCOLORTABLEPARAMETERIVEXT)(getProcAddr("glGetColorTableParameterivEXT"))
	gpGetColorTableParameterivSGI = (C.GPGETCOLORTABLEPARAMETERIVSGI)(getProcAddr("glGetColorTableParameterivSGI"))
	gpGetColorTableSGI = (C.GPGETCOLORTABLESGI)(getProcAddr("glGetColorTableSGI"))
	gpGetCombinerInputParameterfvNV = (C.GPGETCOMBINERINPUTPARAMETERFVNV)(getProcAddr("glGetCombinerInputParameterfvNV"))
	gpGetCombinerInputParameterivNV = (C.GPGETCOMBINERINPUTPARAMETERIVNV)(getProcAddr("glGetCombinerInputParameterivNV"))
	gpGetCombinerOutputParameterfvNV = (C.GPGETCOMBINEROUTPUTPARAMETERFVNV)(getProcAddr("glGetCombinerOutputParameterfvNV"))
	gpGetCombinerOutputParameterivNV = (C.GPGETCOMBINEROUTPUTPARAMETERIVNV)(getProcAddr("glGetCombinerOutputParameterivNV"))
	gpGetCombinerStageParameterfvNV = (C.GPGETCOMBINERSTAGEPARAMETERFVNV)(getProcAddr("glGetCombinerStageParameterfvNV"))
	gpGetCommandHeaderNV = (C.GPGETCOMMANDHEADERNV)(getProcAddr("glGetCommandHeaderNV"))
	gpGetCompressedMultiTexImageEXT = (C.GPGETCOMPRESSEDMULTITEXIMAGEEXT)(getProcAddr("glGetCompressedMultiTexImageEXT"))
	gpGetCompressedTexImage = (C.GPGETCOMPRESSEDTEXIMAGE)(getProcAddr("glGetCompressedTexImage"))
	if gpGetCompressedTexImage == nil {
		return errors.New("glGetCompressedTexImage")
	}
	gpGetCompressedTexImageARB = (C.GPGETCOMPRESSEDTEXIMAGEARB)(getProcAddr("glGetCompressedTexImageARB"))
	gpGetCompressedTextureImage = (C.GPGETCOMPRESSEDTEXTUREIMAGE)(getProcAddr("glGetCompressedTextureImage"))
	if gpGetCompressedTextureImage == nil {
		return errors.New("glGetCompressedTextureImage")
	}
	gpGetCompressedTextureImageEXT = (C.GPGETCOMPRESSEDTEXTUREIMAGEEXT)(getProcAddr("glGetCompressedTextureImageEXT"))
	gpGetCompressedTextureSubImage = (C.GPGETCOMPRESSEDTEXTURESUBIMAGE)(getProcAddr("glGetCompressedTextureSubImage"))
	if gpGetCompressedTextureSubImage == nil {
		return errors.New("glGetCompressedTextureSubImage")
	}
	gpGetConvolutionFilter = (C.GPGETCONVOLUTIONFILTER)(getProcAddr("glGetConvolutionFilter"))
	gpGetConvolutionFilterEXT = (C.GPGETCONVOLUTIONFILTEREXT)(getProcAddr("glGetConvolutionFilterEXT"))
	gpGetConvolutionParameterfv = (C.GPGETCONVOLUTIONPARAMETERFV)(getProcAddr("glGetConvolutionParameterfv"))
	gpGetConvolutionParameterfvEXT = (C.GPGETCONVOLUTIONPARAMETERFVEXT)(getProcAddr("glGetConvolutionParameterfvEXT"))
	gpGetConvolutionParameteriv = (C.GPGETCONVOLUTIONPARAMETERIV)(getProcAddr("glGetConvolutionParameteriv"))
	gpGetConvolutionParameterivEXT = (C.GPGETCONVOLUTIONPARAMETERIVEXT)(getProcAddr("glGetConvolutionParameterivEXT"))
	gpGetConvolutionParameterxvOES = (C.GPGETCONVOLUTIONPARAMETERXVOES)(getProcAddr("glGetConvolutionParameterxvOES"))
	gpGetCoverageModulationTableNV = (C.GPGETCOVERAGEMODULATIONTABLENV)(getProcAddr("glGetCoverageModulationTableNV"))
	gpGetDebugMessageLog = (C.GPGETDEBUGMESSAGELOG)(getProcAddr("glGetDebugMessageLog"))
	if gpGetDebugMessageLog == nil {
		return errors.New("glGetDebugMessageLog")
	}
	gpGetDebugMessageLogAMD = (C.GPGETDEBUGMESSAGELOGAMD)(getProcAddr("glGetDebugMessageLogAMD"))
	gpGetDebugMessageLogARB = (C.GPGETDEBUGMESSAGELOGARB)(getProcAddr("glGetDebugMessageLogARB"))
	gpGetDebugMessageLogKHR = (C.GPGETDEBUGMESSAGELOGKHR)(getProcAddr("glGetDebugMessageLogKHR"))
	gpGetDetailTexFuncSGIS = (C.GPGETDETAILTEXFUNCSGIS)(getProcAddr("glGetDetailTexFuncSGIS"))
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
	gpGetFenceivNV = (C.GPGETFENCEIVNV)(getProcAddr("glGetFenceivNV"))
	gpGetFinalCombinerInputParameterfvNV = (C.GPGETFINALCOMBINERINPUTPARAMETERFVNV)(getProcAddr("glGetFinalCombinerInputParameterfvNV"))
	gpGetFinalCombinerInputParameterivNV = (C.GPGETFINALCOMBINERINPUTPARAMETERIVNV)(getProcAddr("glGetFinalCombinerInputParameterivNV"))
	gpGetFirstPerfQueryIdINTEL = (C.GPGETFIRSTPERFQUERYIDINTEL)(getProcAddr("glGetFirstPerfQueryIdINTEL"))
	gpGetFixedvOES = (C.GPGETFIXEDVOES)(getProcAddr("glGetFixedvOES"))
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
	gpGetFogFuncSGIS = (C.GPGETFOGFUNCSGIS)(getProcAddr("glGetFogFuncSGIS"))
	gpGetFragDataIndex = (C.GPGETFRAGDATAINDEX)(getProcAddr("glGetFragDataIndex"))
	if gpGetFragDataIndex == nil {
		return errors.New("glGetFragDataIndex")
	}
	gpGetFragDataLocation = (C.GPGETFRAGDATALOCATION)(getProcAddr("glGetFragDataLocation"))
	if gpGetFragDataLocation == nil {
		return errors.New("glGetFragDataLocation")
	}
	gpGetFragDataLocationEXT = (C.GPGETFRAGDATALOCATIONEXT)(getProcAddr("glGetFragDataLocationEXT"))
	gpGetFragmentLightfvSGIX = (C.GPGETFRAGMENTLIGHTFVSGIX)(getProcAddr("glGetFragmentLightfvSGIX"))
	gpGetFragmentLightivSGIX = (C.GPGETFRAGMENTLIGHTIVSGIX)(getProcAddr("glGetFragmentLightivSGIX"))
	gpGetFragmentMaterialfvSGIX = (C.GPGETFRAGMENTMATERIALFVSGIX)(getProcAddr("glGetFragmentMaterialfvSGIX"))
	gpGetFragmentMaterialivSGIX = (C.GPGETFRAGMENTMATERIALIVSGIX)(getProcAddr("glGetFragmentMaterialivSGIX"))
	gpGetFramebufferAttachmentParameteriv = (C.GPGETFRAMEBUFFERATTACHMENTPARAMETERIV)(getProcAddr("glGetFramebufferAttachmentParameteriv"))
	if gpGetFramebufferAttachmentParameteriv == nil {
		return errors.New("glGetFramebufferAttachmentParameteriv")
	}
	gpGetFramebufferAttachmentParameterivEXT = (C.GPGETFRAMEBUFFERATTACHMENTPARAMETERIVEXT)(getProcAddr("glGetFramebufferAttachmentParameterivEXT"))
	gpGetFramebufferParameterfvAMD = (C.GPGETFRAMEBUFFERPARAMETERFVAMD)(getProcAddr("glGetFramebufferParameterfvAMD"))
	gpGetFramebufferParameteriv = (C.GPGETFRAMEBUFFERPARAMETERIV)(getProcAddr("glGetFramebufferParameteriv"))
	if gpGetFramebufferParameteriv == nil {
		return errors.New("glGetFramebufferParameteriv")
	}
	gpGetFramebufferParameterivEXT = (C.GPGETFRAMEBUFFERPARAMETERIVEXT)(getProcAddr("glGetFramebufferParameterivEXT"))
	gpGetGraphicsResetStatus = (C.GPGETGRAPHICSRESETSTATUS)(getProcAddr("glGetGraphicsResetStatus"))
	if gpGetGraphicsResetStatus == nil {
		return errors.New("glGetGraphicsResetStatus")
	}
	gpGetGraphicsResetStatusARB = (C.GPGETGRAPHICSRESETSTATUSARB)(getProcAddr("glGetGraphicsResetStatusARB"))
	gpGetGraphicsResetStatusKHR = (C.GPGETGRAPHICSRESETSTATUSKHR)(getProcAddr("glGetGraphicsResetStatusKHR"))
	gpGetHandleARB = (C.GPGETHANDLEARB)(getProcAddr("glGetHandleARB"))
	gpGetHistogram = (C.GPGETHISTOGRAM)(getProcAddr("glGetHistogram"))
	gpGetHistogramEXT = (C.GPGETHISTOGRAMEXT)(getProcAddr("glGetHistogramEXT"))
	gpGetHistogramParameterfv = (C.GPGETHISTOGRAMPARAMETERFV)(getProcAddr("glGetHistogramParameterfv"))
	gpGetHistogramParameterfvEXT = (C.GPGETHISTOGRAMPARAMETERFVEXT)(getProcAddr("glGetHistogramParameterfvEXT"))
	gpGetHistogramParameteriv = (C.GPGETHISTOGRAMPARAMETERIV)(getProcAddr("glGetHistogramParameteriv"))
	gpGetHistogramParameterivEXT = (C.GPGETHISTOGRAMPARAMETERIVEXT)(getProcAddr("glGetHistogramParameterivEXT"))
	gpGetHistogramParameterxvOES = (C.GPGETHISTOGRAMPARAMETERXVOES)(getProcAddr("glGetHistogramParameterxvOES"))
	gpGetImageHandleARB = (C.GPGETIMAGEHANDLEARB)(getProcAddr("glGetImageHandleARB"))
	gpGetImageHandleNV = (C.GPGETIMAGEHANDLENV)(getProcAddr("glGetImageHandleNV"))
	gpGetImageTransformParameterfvHP = (C.GPGETIMAGETRANSFORMPARAMETERFVHP)(getProcAddr("glGetImageTransformParameterfvHP"))
	gpGetImageTransformParameterivHP = (C.GPGETIMAGETRANSFORMPARAMETERIVHP)(getProcAddr("glGetImageTransformParameterivHP"))
	gpGetInfoLogARB = (C.GPGETINFOLOGARB)(getProcAddr("glGetInfoLogARB"))
	gpGetInstrumentsSGIX = (C.GPGETINSTRUMENTSSGIX)(getProcAddr("glGetInstrumentsSGIX"))
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
	if gpGetInternalformati64v == nil {
		return errors.New("glGetInternalformati64v")
	}
	gpGetInternalformativ = (C.GPGETINTERNALFORMATIV)(getProcAddr("glGetInternalformativ"))
	if gpGetInternalformativ == nil {
		return errors.New("glGetInternalformativ")
	}
	gpGetInvariantBooleanvEXT = (C.GPGETINVARIANTBOOLEANVEXT)(getProcAddr("glGetInvariantBooleanvEXT"))
	gpGetInvariantFloatvEXT = (C.GPGETINVARIANTFLOATVEXT)(getProcAddr("glGetInvariantFloatvEXT"))
	gpGetInvariantIntegervEXT = (C.GPGETINVARIANTINTEGERVEXT)(getProcAddr("glGetInvariantIntegervEXT"))
	gpGetLightfv = (C.GPGETLIGHTFV)(getProcAddr("glGetLightfv"))
	if gpGetLightfv == nil {
		return errors.New("glGetLightfv")
	}
	gpGetLightiv = (C.GPGETLIGHTIV)(getProcAddr("glGetLightiv"))
	if gpGetLightiv == nil {
		return errors.New("glGetLightiv")
	}
	gpGetLightxOES = (C.GPGETLIGHTXOES)(getProcAddr("glGetLightxOES"))
	gpGetLightxvOES = (C.GPGETLIGHTXVOES)(getProcAddr("glGetLightxvOES"))
	gpGetListParameterfvSGIX = (C.GPGETLISTPARAMETERFVSGIX)(getProcAddr("glGetListParameterfvSGIX"))
	gpGetListParameterivSGIX = (C.GPGETLISTPARAMETERIVSGIX)(getProcAddr("glGetListParameterivSGIX"))
	gpGetLocalConstantBooleanvEXT = (C.GPGETLOCALCONSTANTBOOLEANVEXT)(getProcAddr("glGetLocalConstantBooleanvEXT"))
	gpGetLocalConstantFloatvEXT = (C.GPGETLOCALCONSTANTFLOATVEXT)(getProcAddr("glGetLocalConstantFloatvEXT"))
	gpGetLocalConstantIntegervEXT = (C.GPGETLOCALCONSTANTINTEGERVEXT)(getProcAddr("glGetLocalConstantIntegervEXT"))
	gpGetMapAttribParameterfvNV = (C.GPGETMAPATTRIBPARAMETERFVNV)(getProcAddr("glGetMapAttribParameterfvNV"))
	gpGetMapAttribParameterivNV = (C.GPGETMAPATTRIBPARAMETERIVNV)(getProcAddr("glGetMapAttribParameterivNV"))
	gpGetMapControlPointsNV = (C.GPGETMAPCONTROLPOINTSNV)(getProcAddr("glGetMapControlPointsNV"))
	gpGetMapParameterfvNV = (C.GPGETMAPPARAMETERFVNV)(getProcAddr("glGetMapParameterfvNV"))
	gpGetMapParameterivNV = (C.GPGETMAPPARAMETERIVNV)(getProcAddr("glGetMapParameterivNV"))
	gpGetMapdv = (C.GPGETMAPDV)(getProcAddr("glGetMapdv"))
	if gpGetMapdv == nil {
		return errors.New("glGetMapdv")
	}
	gpGetMapfv = (C.GPGETMAPFV)(getProcAddr("glGetMapfv"))
	if gpGetMapfv == nil {
		return errors.New("glGetMapfv")
	}
	gpGetMapiv = (C.GPGETMAPIV)(getProcAddr("glGetMapiv"))
	if gpGetMapiv == nil {
		return errors.New("glGetMapiv")
	}
	gpGetMapxvOES = (C.GPGETMAPXVOES)(getProcAddr("glGetMapxvOES"))
	gpGetMaterialfv = (C.GPGETMATERIALFV)(getProcAddr("glGetMaterialfv"))
	if gpGetMaterialfv == nil {
		return errors.New("glGetMaterialfv")
	}
	gpGetMaterialiv = (C.GPGETMATERIALIV)(getProcAddr("glGetMaterialiv"))
	if gpGetMaterialiv == nil {
		return errors.New("glGetMaterialiv")
	}
	gpGetMaterialxOES = (C.GPGETMATERIALXOES)(getProcAddr("glGetMaterialxOES"))
	gpGetMaterialxvOES = (C.GPGETMATERIALXVOES)(getProcAddr("glGetMaterialxvOES"))
	gpGetMemoryObjectParameterivEXT = (C.GPGETMEMORYOBJECTPARAMETERIVEXT)(getProcAddr("glGetMemoryObjectParameterivEXT"))
	gpGetMinmax = (C.GPGETMINMAX)(getProcAddr("glGetMinmax"))
	gpGetMinmaxEXT = (C.GPGETMINMAXEXT)(getProcAddr("glGetMinmaxEXT"))
	gpGetMinmaxParameterfv = (C.GPGETMINMAXPARAMETERFV)(getProcAddr("glGetMinmaxParameterfv"))
	gpGetMinmaxParameterfvEXT = (C.GPGETMINMAXPARAMETERFVEXT)(getProcAddr("glGetMinmaxParameterfvEXT"))
	gpGetMinmaxParameteriv = (C.GPGETMINMAXPARAMETERIV)(getProcAddr("glGetMinmaxParameteriv"))
	gpGetMinmaxParameterivEXT = (C.GPGETMINMAXPARAMETERIVEXT)(getProcAddr("glGetMinmaxParameterivEXT"))
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
	gpGetMultisamplefvNV = (C.GPGETMULTISAMPLEFVNV)(getProcAddr("glGetMultisamplefvNV"))
	gpGetNamedBufferParameteri64v = (C.GPGETNAMEDBUFFERPARAMETERI64V)(getProcAddr("glGetNamedBufferParameteri64v"))
	if gpGetNamedBufferParameteri64v == nil {
		return errors.New("glGetNamedBufferParameteri64v")
	}
	gpGetNamedBufferParameteriv = (C.GPGETNAMEDBUFFERPARAMETERIV)(getProcAddr("glGetNamedBufferParameteriv"))
	if gpGetNamedBufferParameteriv == nil {
		return errors.New("glGetNamedBufferParameteriv")
	}
	gpGetNamedBufferParameterivEXT = (C.GPGETNAMEDBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedBufferParameterivEXT"))
	gpGetNamedBufferParameterui64vNV = (C.GPGETNAMEDBUFFERPARAMETERUI64VNV)(getProcAddr("glGetNamedBufferParameterui64vNV"))
	gpGetNamedBufferPointerv = (C.GPGETNAMEDBUFFERPOINTERV)(getProcAddr("glGetNamedBufferPointerv"))
	if gpGetNamedBufferPointerv == nil {
		return errors.New("glGetNamedBufferPointerv")
	}
	gpGetNamedBufferPointervEXT = (C.GPGETNAMEDBUFFERPOINTERVEXT)(getProcAddr("glGetNamedBufferPointervEXT"))
	gpGetNamedBufferSubData = (C.GPGETNAMEDBUFFERSUBDATA)(getProcAddr("glGetNamedBufferSubData"))
	if gpGetNamedBufferSubData == nil {
		return errors.New("glGetNamedBufferSubData")
	}
	gpGetNamedBufferSubDataEXT = (C.GPGETNAMEDBUFFERSUBDATAEXT)(getProcAddr("glGetNamedBufferSubDataEXT"))
	gpGetNamedFramebufferAttachmentParameteriv = (C.GPGETNAMEDFRAMEBUFFERATTACHMENTPARAMETERIV)(getProcAddr("glGetNamedFramebufferAttachmentParameteriv"))
	if gpGetNamedFramebufferAttachmentParameteriv == nil {
		return errors.New("glGetNamedFramebufferAttachmentParameteriv")
	}
	gpGetNamedFramebufferAttachmentParameterivEXT = (C.GPGETNAMEDFRAMEBUFFERATTACHMENTPARAMETERIVEXT)(getProcAddr("glGetNamedFramebufferAttachmentParameterivEXT"))
	gpGetNamedFramebufferParameterfvAMD = (C.GPGETNAMEDFRAMEBUFFERPARAMETERFVAMD)(getProcAddr("glGetNamedFramebufferParameterfvAMD"))
	gpGetNamedFramebufferParameteriv = (C.GPGETNAMEDFRAMEBUFFERPARAMETERIV)(getProcAddr("glGetNamedFramebufferParameteriv"))
	if gpGetNamedFramebufferParameteriv == nil {
		return errors.New("glGetNamedFramebufferParameteriv")
	}
	gpGetNamedFramebufferParameterivEXT = (C.GPGETNAMEDFRAMEBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedFramebufferParameterivEXT"))
	gpGetNamedProgramLocalParameterIivEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERIIVEXT)(getProcAddr("glGetNamedProgramLocalParameterIivEXT"))
	gpGetNamedProgramLocalParameterIuivEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERIUIVEXT)(getProcAddr("glGetNamedProgramLocalParameterIuivEXT"))
	gpGetNamedProgramLocalParameterdvEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERDVEXT)(getProcAddr("glGetNamedProgramLocalParameterdvEXT"))
	gpGetNamedProgramLocalParameterfvEXT = (C.GPGETNAMEDPROGRAMLOCALPARAMETERFVEXT)(getProcAddr("glGetNamedProgramLocalParameterfvEXT"))
	gpGetNamedProgramStringEXT = (C.GPGETNAMEDPROGRAMSTRINGEXT)(getProcAddr("glGetNamedProgramStringEXT"))
	gpGetNamedProgramivEXT = (C.GPGETNAMEDPROGRAMIVEXT)(getProcAddr("glGetNamedProgramivEXT"))
	gpGetNamedRenderbufferParameteriv = (C.GPGETNAMEDRENDERBUFFERPARAMETERIV)(getProcAddr("glGetNamedRenderbufferParameteriv"))
	if gpGetNamedRenderbufferParameteriv == nil {
		return errors.New("glGetNamedRenderbufferParameteriv")
	}
	gpGetNamedRenderbufferParameterivEXT = (C.GPGETNAMEDRENDERBUFFERPARAMETERIVEXT)(getProcAddr("glGetNamedRenderbufferParameterivEXT"))
	gpGetNamedStringARB = (C.GPGETNAMEDSTRINGARB)(getProcAddr("glGetNamedStringARB"))
	gpGetNamedStringivARB = (C.GPGETNAMEDSTRINGIVARB)(getProcAddr("glGetNamedStringivARB"))
	gpGetNextPerfQueryIdINTEL = (C.GPGETNEXTPERFQUERYIDINTEL)(getProcAddr("glGetNextPerfQueryIdINTEL"))
	gpGetObjectBufferfvATI = (C.GPGETOBJECTBUFFERFVATI)(getProcAddr("glGetObjectBufferfvATI"))
	gpGetObjectBufferivATI = (C.GPGETOBJECTBUFFERIVATI)(getProcAddr("glGetObjectBufferivATI"))
	gpGetObjectLabel = (C.GPGETOBJECTLABEL)(getProcAddr("glGetObjectLabel"))
	if gpGetObjectLabel == nil {
		return errors.New("glGetObjectLabel")
	}
	gpGetObjectLabelEXT = (C.GPGETOBJECTLABELEXT)(getProcAddr("glGetObjectLabelEXT"))
	gpGetObjectLabelKHR = (C.GPGETOBJECTLABELKHR)(getProcAddr("glGetObjectLabelKHR"))
	gpGetObjectParameterfvARB = (C.GPGETOBJECTPARAMETERFVARB)(getProcAddr("glGetObjectParameterfvARB"))
	gpGetObjectParameterivAPPLE = (C.GPGETOBJECTPARAMETERIVAPPLE)(getProcAddr("glGetObjectParameterivAPPLE"))
	gpGetObjectParameterivARB = (C.GPGETOBJECTPARAMETERIVARB)(getProcAddr("glGetObjectParameterivARB"))
	gpGetObjectPtrLabel = (C.GPGETOBJECTPTRLABEL)(getProcAddr("glGetObjectPtrLabel"))
	if gpGetObjectPtrLabel == nil {
		return errors.New("glGetObjectPtrLabel")
	}
	gpGetObjectPtrLabelKHR = (C.GPGETOBJECTPTRLABELKHR)(getProcAddr("glGetObjectPtrLabelKHR"))
	gpGetOcclusionQueryivNV = (C.GPGETOCCLUSIONQUERYIVNV)(getProcAddr("glGetOcclusionQueryivNV"))
	gpGetOcclusionQueryuivNV = (C.GPGETOCCLUSIONQUERYUIVNV)(getProcAddr("glGetOcclusionQueryuivNV"))
	gpGetPathColorGenfvNV = (C.GPGETPATHCOLORGENFVNV)(getProcAddr("glGetPathColorGenfvNV"))
	gpGetPathColorGenivNV = (C.GPGETPATHCOLORGENIVNV)(getProcAddr("glGetPathColorGenivNV"))
	gpGetPathCommandsNV = (C.GPGETPATHCOMMANDSNV)(getProcAddr("glGetPathCommandsNV"))
	gpGetPathCoordsNV = (C.GPGETPATHCOORDSNV)(getProcAddr("glGetPathCoordsNV"))
	gpGetPathDashArrayNV = (C.GPGETPATHDASHARRAYNV)(getProcAddr("glGetPathDashArrayNV"))
	gpGetPathLengthNV = (C.GPGETPATHLENGTHNV)(getProcAddr("glGetPathLengthNV"))
	gpGetPathMetricRangeNV = (C.GPGETPATHMETRICRANGENV)(getProcAddr("glGetPathMetricRangeNV"))
	gpGetPathMetricsNV = (C.GPGETPATHMETRICSNV)(getProcAddr("glGetPathMetricsNV"))
	gpGetPathParameterfvNV = (C.GPGETPATHPARAMETERFVNV)(getProcAddr("glGetPathParameterfvNV"))
	gpGetPathParameterivNV = (C.GPGETPATHPARAMETERIVNV)(getProcAddr("glGetPathParameterivNV"))
	gpGetPathSpacingNV = (C.GPGETPATHSPACINGNV)(getProcAddr("glGetPathSpacingNV"))
	gpGetPathTexGenfvNV = (C.GPGETPATHTEXGENFVNV)(getProcAddr("glGetPathTexGenfvNV"))
	gpGetPathTexGenivNV = (C.GPGETPATHTEXGENIVNV)(getProcAddr("glGetPathTexGenivNV"))
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
	gpGetPixelMapfv = (C.GPGETPIXELMAPFV)(getProcAddr("glGetPixelMapfv"))
	if gpGetPixelMapfv == nil {
		return errors.New("glGetPixelMapfv")
	}
	gpGetPixelMapuiv = (C.GPGETPIXELMAPUIV)(getProcAddr("glGetPixelMapuiv"))
	if gpGetPixelMapuiv == nil {
		return errors.New("glGetPixelMapuiv")
	}
	gpGetPixelMapusv = (C.GPGETPIXELMAPUSV)(getProcAddr("glGetPixelMapusv"))
	if gpGetPixelMapusv == nil {
		return errors.New("glGetPixelMapusv")
	}
	gpGetPixelMapxv = (C.GPGETPIXELMAPXV)(getProcAddr("glGetPixelMapxv"))
	gpGetPixelTexGenParameterfvSGIS = (C.GPGETPIXELTEXGENPARAMETERFVSGIS)(getProcAddr("glGetPixelTexGenParameterfvSGIS"))
	gpGetPixelTexGenParameterivSGIS = (C.GPGETPIXELTEXGENPARAMETERIVSGIS)(getProcAddr("glGetPixelTexGenParameterivSGIS"))
	gpGetPixelTransformParameterfvEXT = (C.GPGETPIXELTRANSFORMPARAMETERFVEXT)(getProcAddr("glGetPixelTransformParameterfvEXT"))
	gpGetPixelTransformParameterivEXT = (C.GPGETPIXELTRANSFORMPARAMETERIVEXT)(getProcAddr("glGetPixelTransformParameterivEXT"))
	gpGetPointerIndexedvEXT = (C.GPGETPOINTERINDEXEDVEXT)(getProcAddr("glGetPointerIndexedvEXT"))
	gpGetPointeri_vEXT = (C.GPGETPOINTERI_VEXT)(getProcAddr("glGetPointeri_vEXT"))
	gpGetPointerv = (C.GPGETPOINTERV)(getProcAddr("glGetPointerv"))
	if gpGetPointerv == nil {
		return errors.New("glGetPointerv")
	}
	gpGetPointervEXT = (C.GPGETPOINTERVEXT)(getProcAddr("glGetPointervEXT"))
	gpGetPointervKHR = (C.GPGETPOINTERVKHR)(getProcAddr("glGetPointervKHR"))
	gpGetPolygonStipple = (C.GPGETPOLYGONSTIPPLE)(getProcAddr("glGetPolygonStipple"))
	if gpGetPolygonStipple == nil {
		return errors.New("glGetPolygonStipple")
	}
	gpGetProgramBinary = (C.GPGETPROGRAMBINARY)(getProcAddr("glGetProgramBinary"))
	if gpGetProgramBinary == nil {
		return errors.New("glGetProgramBinary")
	}
	gpGetProgramEnvParameterIivNV = (C.GPGETPROGRAMENVPARAMETERIIVNV)(getProcAddr("glGetProgramEnvParameterIivNV"))
	gpGetProgramEnvParameterIuivNV = (C.GPGETPROGRAMENVPARAMETERIUIVNV)(getProcAddr("glGetProgramEnvParameterIuivNV"))
	gpGetProgramEnvParameterdvARB = (C.GPGETPROGRAMENVPARAMETERDVARB)(getProcAddr("glGetProgramEnvParameterdvARB"))
	gpGetProgramEnvParameterfvARB = (C.GPGETPROGRAMENVPARAMETERFVARB)(getProcAddr("glGetProgramEnvParameterfvARB"))
	gpGetProgramInfoLog = (C.GPGETPROGRAMINFOLOG)(getProcAddr("glGetProgramInfoLog"))
	if gpGetProgramInfoLog == nil {
		return errors.New("glGetProgramInfoLog")
	}
	gpGetProgramInterfaceiv = (C.GPGETPROGRAMINTERFACEIV)(getProcAddr("glGetProgramInterfaceiv"))
	if gpGetProgramInterfaceiv == nil {
		return errors.New("glGetProgramInterfaceiv")
	}
	gpGetProgramLocalParameterIivNV = (C.GPGETPROGRAMLOCALPARAMETERIIVNV)(getProcAddr("glGetProgramLocalParameterIivNV"))
	gpGetProgramLocalParameterIuivNV = (C.GPGETPROGRAMLOCALPARAMETERIUIVNV)(getProcAddr("glGetProgramLocalParameterIuivNV"))
	gpGetProgramLocalParameterdvARB = (C.GPGETPROGRAMLOCALPARAMETERDVARB)(getProcAddr("glGetProgramLocalParameterdvARB"))
	gpGetProgramLocalParameterfvARB = (C.GPGETPROGRAMLOCALPARAMETERFVARB)(getProcAddr("glGetProgramLocalParameterfvARB"))
	gpGetProgramNamedParameterdvNV = (C.GPGETPROGRAMNAMEDPARAMETERDVNV)(getProcAddr("glGetProgramNamedParameterdvNV"))
	gpGetProgramNamedParameterfvNV = (C.GPGETPROGRAMNAMEDPARAMETERFVNV)(getProcAddr("glGetProgramNamedParameterfvNV"))
	gpGetProgramParameterdvNV = (C.GPGETPROGRAMPARAMETERDVNV)(getProcAddr("glGetProgramParameterdvNV"))
	gpGetProgramParameterfvNV = (C.GPGETPROGRAMPARAMETERFVNV)(getProcAddr("glGetProgramParameterfvNV"))
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
	gpGetProgramResourceLocationIndex = (C.GPGETPROGRAMRESOURCELOCATIONINDEX)(getProcAddr("glGetProgramResourceLocationIndex"))
	if gpGetProgramResourceLocationIndex == nil {
		return errors.New("glGetProgramResourceLocationIndex")
	}
	gpGetProgramResourceName = (C.GPGETPROGRAMRESOURCENAME)(getProcAddr("glGetProgramResourceName"))
	if gpGetProgramResourceName == nil {
		return errors.New("glGetProgramResourceName")
	}
	gpGetProgramResourcefvNV = (C.GPGETPROGRAMRESOURCEFVNV)(getProcAddr("glGetProgramResourcefvNV"))
	gpGetProgramResourceiv = (C.GPGETPROGRAMRESOURCEIV)(getProcAddr("glGetProgramResourceiv"))
	if gpGetProgramResourceiv == nil {
		return errors.New("glGetProgramResourceiv")
	}
	gpGetProgramStageiv = (C.GPGETPROGRAMSTAGEIV)(getProcAddr("glGetProgramStageiv"))
	if gpGetProgramStageiv == nil {
		return errors.New("glGetProgramStageiv")
	}
	gpGetProgramStringARB = (C.GPGETPROGRAMSTRINGARB)(getProcAddr("glGetProgramStringARB"))
	gpGetProgramStringNV = (C.GPGETPROGRAMSTRINGNV)(getProcAddr("glGetProgramStringNV"))
	gpGetProgramSubroutineParameteruivNV = (C.GPGETPROGRAMSUBROUTINEPARAMETERUIVNV)(getProcAddr("glGetProgramSubroutineParameteruivNV"))
	gpGetProgramiv = (C.GPGETPROGRAMIV)(getProcAddr("glGetProgramiv"))
	if gpGetProgramiv == nil {
		return errors.New("glGetProgramiv")
	}
	gpGetProgramivARB = (C.GPGETPROGRAMIVARB)(getProcAddr("glGetProgramivARB"))
	gpGetProgramivNV = (C.GPGETPROGRAMIVNV)(getProcAddr("glGetProgramivNV"))
	gpGetQueryBufferObjecti64v = (C.GPGETQUERYBUFFEROBJECTI64V)(getProcAddr("glGetQueryBufferObjecti64v"))
	if gpGetQueryBufferObjecti64v == nil {
		return errors.New("glGetQueryBufferObjecti64v")
	}
	gpGetQueryBufferObjectiv = (C.GPGETQUERYBUFFEROBJECTIV)(getProcAddr("glGetQueryBufferObjectiv"))
	if gpGetQueryBufferObjectiv == nil {
		return errors.New("glGetQueryBufferObjectiv")
	}
	gpGetQueryBufferObjectui64v = (C.GPGETQUERYBUFFEROBJECTUI64V)(getProcAddr("glGetQueryBufferObjectui64v"))
	if gpGetQueryBufferObjectui64v == nil {
		return errors.New("glGetQueryBufferObjectui64v")
	}
	gpGetQueryBufferObjectuiv = (C.GPGETQUERYBUFFEROBJECTUIV)(getProcAddr("glGetQueryBufferObjectuiv"))
	if gpGetQueryBufferObjectuiv == nil {
		return errors.New("glGetQueryBufferObjectuiv")
	}
	gpGetQueryIndexediv = (C.GPGETQUERYINDEXEDIV)(getProcAddr("glGetQueryIndexediv"))
	if gpGetQueryIndexediv == nil {
		return errors.New("glGetQueryIndexediv")
	}
	gpGetQueryObjecti64v = (C.GPGETQUERYOBJECTI64V)(getProcAddr("glGetQueryObjecti64v"))
	if gpGetQueryObjecti64v == nil {
		return errors.New("glGetQueryObjecti64v")
	}
	gpGetQueryObjecti64vEXT = (C.GPGETQUERYOBJECTI64VEXT)(getProcAddr("glGetQueryObjecti64vEXT"))
	gpGetQueryObjectiv = (C.GPGETQUERYOBJECTIV)(getProcAddr("glGetQueryObjectiv"))
	if gpGetQueryObjectiv == nil {
		return errors.New("glGetQueryObjectiv")
	}
	gpGetQueryObjectivARB = (C.GPGETQUERYOBJECTIVARB)(getProcAddr("glGetQueryObjectivARB"))
	gpGetQueryObjectui64v = (C.GPGETQUERYOBJECTUI64V)(getProcAddr("glGetQueryObjectui64v"))
	if gpGetQueryObjectui64v == nil {
		return errors.New("glGetQueryObjectui64v")
	}
	gpGetQueryObjectui64vEXT = (C.GPGETQUERYOBJECTUI64VEXT)(getProcAddr("glGetQueryObjectui64vEXT"))
	gpGetQueryObjectuiv = (C.GPGETQUERYOBJECTUIV)(getProcAddr("glGetQueryObjectuiv"))
	if gpGetQueryObjectuiv == nil {
		return errors.New("glGetQueryObjectuiv")
	}
	gpGetQueryObjectuivARB = (C.GPGETQUERYOBJECTUIVARB)(getProcAddr("glGetQueryObjectuivARB"))
	gpGetQueryiv = (C.GPGETQUERYIV)(getProcAddr("glGetQueryiv"))
	if gpGetQueryiv == nil {
		return errors.New("glGetQueryiv")
	}
	gpGetQueryivARB = (C.GPGETQUERYIVARB)(getProcAddr("glGetQueryivARB"))
	gpGetRenderbufferParameteriv = (C.GPGETRENDERBUFFERPARAMETERIV)(getProcAddr("glGetRenderbufferParameteriv"))
	if gpGetRenderbufferParameteriv == nil {
		return errors.New("glGetRenderbufferParameteriv")
	}
	gpGetRenderbufferParameterivEXT = (C.GPGETRENDERBUFFERPARAMETERIVEXT)(getProcAddr("glGetRenderbufferParameterivEXT"))
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
	gpGetSemaphoreParameterui64vEXT = (C.GPGETSEMAPHOREPARAMETERUI64VEXT)(getProcAddr("glGetSemaphoreParameterui64vEXT"))
	gpGetSeparableFilter = (C.GPGETSEPARABLEFILTER)(getProcAddr("glGetSeparableFilter"))
	gpGetSeparableFilterEXT = (C.GPGETSEPARABLEFILTEREXT)(getProcAddr("glGetSeparableFilterEXT"))
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
	gpGetShaderSourceARB = (C.GPGETSHADERSOURCEARB)(getProcAddr("glGetShaderSourceARB"))
	gpGetShaderiv = (C.GPGETSHADERIV)(getProcAddr("glGetShaderiv"))
	if gpGetShaderiv == nil {
		return errors.New("glGetShaderiv")
	}
	gpGetSharpenTexFuncSGIS = (C.GPGETSHARPENTEXFUNCSGIS)(getProcAddr("glGetSharpenTexFuncSGIS"))
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
	gpGetTexBumpParameterfvATI = (C.GPGETTEXBUMPPARAMETERFVATI)(getProcAddr("glGetTexBumpParameterfvATI"))
	gpGetTexBumpParameterivATI = (C.GPGETTEXBUMPPARAMETERIVATI)(getProcAddr("glGetTexBumpParameterivATI"))
	gpGetTexEnvfv = (C.GPGETTEXENVFV)(getProcAddr("glGetTexEnvfv"))
	if gpGetTexEnvfv == nil {
		return errors.New("glGetTexEnvfv")
	}
	gpGetTexEnviv = (C.GPGETTEXENVIV)(getProcAddr("glGetTexEnviv"))
	if gpGetTexEnviv == nil {
		return errors.New("glGetTexEnviv")
	}
	gpGetTexEnvxvOES = (C.GPGETTEXENVXVOES)(getProcAddr("glGetTexEnvxvOES"))
	gpGetTexFilterFuncSGIS = (C.GPGETTEXFILTERFUNCSGIS)(getProcAddr("glGetTexFilterFuncSGIS"))
	gpGetTexGendv = (C.GPGETTEXGENDV)(getProcAddr("glGetTexGendv"))
	if gpGetTexGendv == nil {
		return errors.New("glGetTexGendv")
	}
	gpGetTexGenfv = (C.GPGETTEXGENFV)(getProcAddr("glGetTexGenfv"))
	if gpGetTexGenfv == nil {
		return errors.New("glGetTexGenfv")
	}
	gpGetTexGeniv = (C.GPGETTEXGENIV)(getProcAddr("glGetTexGeniv"))
	if gpGetTexGeniv == nil {
		return errors.New("glGetTexGeniv")
	}
	gpGetTexGenxvOES = (C.GPGETTEXGENXVOES)(getProcAddr("glGetTexGenxvOES"))
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
	gpGetTexLevelParameterxvOES = (C.GPGETTEXLEVELPARAMETERXVOES)(getProcAddr("glGetTexLevelParameterxvOES"))
	gpGetTexParameterIiv = (C.GPGETTEXPARAMETERIIV)(getProcAddr("glGetTexParameterIiv"))
	if gpGetTexParameterIiv == nil {
		return errors.New("glGetTexParameterIiv")
	}
	gpGetTexParameterIivEXT = (C.GPGETTEXPARAMETERIIVEXT)(getProcAddr("glGetTexParameterIivEXT"))
	gpGetTexParameterIuiv = (C.GPGETTEXPARAMETERIUIV)(getProcAddr("glGetTexParameterIuiv"))
	if gpGetTexParameterIuiv == nil {
		return errors.New("glGetTexParameterIuiv")
	}
	gpGetTexParameterIuivEXT = (C.GPGETTEXPARAMETERIUIVEXT)(getProcAddr("glGetTexParameterIuivEXT"))
	gpGetTexParameterPointervAPPLE = (C.GPGETTEXPARAMETERPOINTERVAPPLE)(getProcAddr("glGetTexParameterPointervAPPLE"))
	gpGetTexParameterfv = (C.GPGETTEXPARAMETERFV)(getProcAddr("glGetTexParameterfv"))
	if gpGetTexParameterfv == nil {
		return errors.New("glGetTexParameterfv")
	}
	gpGetTexParameteriv = (C.GPGETTEXPARAMETERIV)(getProcAddr("glGetTexParameteriv"))
	if gpGetTexParameteriv == nil {
		return errors.New("glGetTexParameteriv")
	}
	gpGetTexParameterxvOES = (C.GPGETTEXPARAMETERXVOES)(getProcAddr("glGetTexParameterxvOES"))
	gpGetTextureHandleARB = (C.GPGETTEXTUREHANDLEARB)(getProcAddr("glGetTextureHandleARB"))
	gpGetTextureHandleNV = (C.GPGETTEXTUREHANDLENV)(getProcAddr("glGetTextureHandleNV"))
	gpGetTextureImage = (C.GPGETTEXTUREIMAGE)(getProcAddr("glGetTextureImage"))
	if gpGetTextureImage == nil {
		return errors.New("glGetTextureImage")
	}
	gpGetTextureImageEXT = (C.GPGETTEXTUREIMAGEEXT)(getProcAddr("glGetTextureImageEXT"))
	gpGetTextureLevelParameterfv = (C.GPGETTEXTURELEVELPARAMETERFV)(getProcAddr("glGetTextureLevelParameterfv"))
	if gpGetTextureLevelParameterfv == nil {
		return errors.New("glGetTextureLevelParameterfv")
	}
	gpGetTextureLevelParameterfvEXT = (C.GPGETTEXTURELEVELPARAMETERFVEXT)(getProcAddr("glGetTextureLevelParameterfvEXT"))
	gpGetTextureLevelParameteriv = (C.GPGETTEXTURELEVELPARAMETERIV)(getProcAddr("glGetTextureLevelParameteriv"))
	if gpGetTextureLevelParameteriv == nil {
		return errors.New("glGetTextureLevelParameteriv")
	}
	gpGetTextureLevelParameterivEXT = (C.GPGETTEXTURELEVELPARAMETERIVEXT)(getProcAddr("glGetTextureLevelParameterivEXT"))
	gpGetTextureParameterIiv = (C.GPGETTEXTUREPARAMETERIIV)(getProcAddr("glGetTextureParameterIiv"))
	if gpGetTextureParameterIiv == nil {
		return errors.New("glGetTextureParameterIiv")
	}
	gpGetTextureParameterIivEXT = (C.GPGETTEXTUREPARAMETERIIVEXT)(getProcAddr("glGetTextureParameterIivEXT"))
	gpGetTextureParameterIuiv = (C.GPGETTEXTUREPARAMETERIUIV)(getProcAddr("glGetTextureParameterIuiv"))
	if gpGetTextureParameterIuiv == nil {
		return errors.New("glGetTextureParameterIuiv")
	}
	gpGetTextureParameterIuivEXT = (C.GPGETTEXTUREPARAMETERIUIVEXT)(getProcAddr("glGetTextureParameterIuivEXT"))
	gpGetTextureParameterfv = (C.GPGETTEXTUREPARAMETERFV)(getProcAddr("glGetTextureParameterfv"))
	if gpGetTextureParameterfv == nil {
		return errors.New("glGetTextureParameterfv")
	}
	gpGetTextureParameterfvEXT = (C.GPGETTEXTUREPARAMETERFVEXT)(getProcAddr("glGetTextureParameterfvEXT"))
	gpGetTextureParameteriv = (C.GPGETTEXTUREPARAMETERIV)(getProcAddr("glGetTextureParameteriv"))
	if gpGetTextureParameteriv == nil {
		return errors.New("glGetTextureParameteriv")
	}
	gpGetTextureParameterivEXT = (C.GPGETTEXTUREPARAMETERIVEXT)(getProcAddr("glGetTextureParameterivEXT"))
	gpGetTextureSamplerHandleARB = (C.GPGETTEXTURESAMPLERHANDLEARB)(getProcAddr("glGetTextureSamplerHandleARB"))
	gpGetTextureSamplerHandleNV = (C.GPGETTEXTURESAMPLERHANDLENV)(getProcAddr("glGetTextureSamplerHandleNV"))
	gpGetTextureSubImage = (C.GPGETTEXTURESUBIMAGE)(getProcAddr("glGetTextureSubImage"))
	if gpGetTextureSubImage == nil {
		return errors.New("glGetTextureSubImage")
	}
	gpGetTrackMatrixivNV = (C.GPGETTRACKMATRIXIVNV)(getProcAddr("glGetTrackMatrixivNV"))
	gpGetTransformFeedbackVarying = (C.GPGETTRANSFORMFEEDBACKVARYING)(getProcAddr("glGetTransformFeedbackVarying"))
	if gpGetTransformFeedbackVarying == nil {
		return errors.New("glGetTransformFeedbackVarying")
	}
	gpGetTransformFeedbackVaryingEXT = (C.GPGETTRANSFORMFEEDBACKVARYINGEXT)(getProcAddr("glGetTransformFeedbackVaryingEXT"))
	gpGetTransformFeedbackVaryingNV = (C.GPGETTRANSFORMFEEDBACKVARYINGNV)(getProcAddr("glGetTransformFeedbackVaryingNV"))
	gpGetTransformFeedbacki64_v = (C.GPGETTRANSFORMFEEDBACKI64_V)(getProcAddr("glGetTransformFeedbacki64_v"))
	if gpGetTransformFeedbacki64_v == nil {
		return errors.New("glGetTransformFeedbacki64_v")
	}
	gpGetTransformFeedbacki_v = (C.GPGETTRANSFORMFEEDBACKI_V)(getProcAddr("glGetTransformFeedbacki_v"))
	if gpGetTransformFeedbacki_v == nil {
		return errors.New("glGetTransformFeedbacki_v")
	}
	gpGetTransformFeedbackiv = (C.GPGETTRANSFORMFEEDBACKIV)(getProcAddr("glGetTransformFeedbackiv"))
	if gpGetTransformFeedbackiv == nil {
		return errors.New("glGetTransformFeedbackiv")
	}
	gpGetUniformBlockIndex = (C.GPGETUNIFORMBLOCKINDEX)(getProcAddr("glGetUniformBlockIndex"))
	if gpGetUniformBlockIndex == nil {
		return errors.New("glGetUniformBlockIndex")
	}
	gpGetUniformBufferSizeEXT = (C.GPGETUNIFORMBUFFERSIZEEXT)(getProcAddr("glGetUniformBufferSizeEXT"))
	gpGetUniformIndices = (C.GPGETUNIFORMINDICES)(getProcAddr("glGetUniformIndices"))
	if gpGetUniformIndices == nil {
		return errors.New("glGetUniformIndices")
	}
	gpGetUniformLocation = (C.GPGETUNIFORMLOCATION)(getProcAddr("glGetUniformLocation"))
	if gpGetUniformLocation == nil {
		return errors.New("glGetUniformLocation")
	}
	gpGetUniformLocationARB = (C.GPGETUNIFORMLOCATIONARB)(getProcAddr("glGetUniformLocationARB"))
	gpGetUniformOffsetEXT = (C.GPGETUNIFORMOFFSETEXT)(getProcAddr("glGetUniformOffsetEXT"))
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
	gpGetUniformfvARB = (C.GPGETUNIFORMFVARB)(getProcAddr("glGetUniformfvARB"))
	gpGetUniformi64vARB = (C.GPGETUNIFORMI64VARB)(getProcAddr("glGetUniformi64vARB"))
	gpGetUniformi64vNV = (C.GPGETUNIFORMI64VNV)(getProcAddr("glGetUniformi64vNV"))
	gpGetUniformiv = (C.GPGETUNIFORMIV)(getProcAddr("glGetUniformiv"))
	if gpGetUniformiv == nil {
		return errors.New("glGetUniformiv")
	}
	gpGetUniformivARB = (C.GPGETUNIFORMIVARB)(getProcAddr("glGetUniformivARB"))
	gpGetUniformui64vARB = (C.GPGETUNIFORMUI64VARB)(getProcAddr("glGetUniformui64vARB"))
	gpGetUniformui64vNV = (C.GPGETUNIFORMUI64VNV)(getProcAddr("glGetUniformui64vNV"))
	gpGetUniformuiv = (C.GPGETUNIFORMUIV)(getProcAddr("glGetUniformuiv"))
	if gpGetUniformuiv == nil {
		return errors.New("glGetUniformuiv")
	}
	gpGetUniformuivEXT = (C.GPGETUNIFORMUIVEXT)(getProcAddr("glGetUniformuivEXT"))
	gpGetUnsignedBytei_vEXT = (C.GPGETUNSIGNEDBYTEI_VEXT)(getProcAddr("glGetUnsignedBytei_vEXT"))
	gpGetUnsignedBytevEXT = (C.GPGETUNSIGNEDBYTEVEXT)(getProcAddr("glGetUnsignedBytevEXT"))
	gpGetVariantArrayObjectfvATI = (C.GPGETVARIANTARRAYOBJECTFVATI)(getProcAddr("glGetVariantArrayObjectfvATI"))
	gpGetVariantArrayObjectivATI = (C.GPGETVARIANTARRAYOBJECTIVATI)(getProcAddr("glGetVariantArrayObjectivATI"))
	gpGetVariantBooleanvEXT = (C.GPGETVARIANTBOOLEANVEXT)(getProcAddr("glGetVariantBooleanvEXT"))
	gpGetVariantFloatvEXT = (C.GPGETVARIANTFLOATVEXT)(getProcAddr("glGetVariantFloatvEXT"))
	gpGetVariantIntegervEXT = (C.GPGETVARIANTINTEGERVEXT)(getProcAddr("glGetVariantIntegervEXT"))
	gpGetVariantPointervEXT = (C.GPGETVARIANTPOINTERVEXT)(getProcAddr("glGetVariantPointervEXT"))
	gpGetVaryingLocationNV = (C.GPGETVARYINGLOCATIONNV)(getProcAddr("glGetVaryingLocationNV"))
	gpGetVertexArrayIndexed64iv = (C.GPGETVERTEXARRAYINDEXED64IV)(getProcAddr("glGetVertexArrayIndexed64iv"))
	if gpGetVertexArrayIndexed64iv == nil {
		return errors.New("glGetVertexArrayIndexed64iv")
	}
	gpGetVertexArrayIndexediv = (C.GPGETVERTEXARRAYINDEXEDIV)(getProcAddr("glGetVertexArrayIndexediv"))
	if gpGetVertexArrayIndexediv == nil {
		return errors.New("glGetVertexArrayIndexediv")
	}
	gpGetVertexArrayIntegeri_vEXT = (C.GPGETVERTEXARRAYINTEGERI_VEXT)(getProcAddr("glGetVertexArrayIntegeri_vEXT"))
	gpGetVertexArrayIntegervEXT = (C.GPGETVERTEXARRAYINTEGERVEXT)(getProcAddr("glGetVertexArrayIntegervEXT"))
	gpGetVertexArrayPointeri_vEXT = (C.GPGETVERTEXARRAYPOINTERI_VEXT)(getProcAddr("glGetVertexArrayPointeri_vEXT"))
	gpGetVertexArrayPointervEXT = (C.GPGETVERTEXARRAYPOINTERVEXT)(getProcAddr("glGetVertexArrayPointervEXT"))
	gpGetVertexArrayiv = (C.GPGETVERTEXARRAYIV)(getProcAddr("glGetVertexArrayiv"))
	if gpGetVertexArrayiv == nil {
		return errors.New("glGetVertexArrayiv")
	}
	gpGetVertexAttribArrayObjectfvATI = (C.GPGETVERTEXATTRIBARRAYOBJECTFVATI)(getProcAddr("glGetVertexAttribArrayObjectfvATI"))
	gpGetVertexAttribArrayObjectivATI = (C.GPGETVERTEXATTRIBARRAYOBJECTIVATI)(getProcAddr("glGetVertexAttribArrayObjectivATI"))
	gpGetVertexAttribIiv = (C.GPGETVERTEXATTRIBIIV)(getProcAddr("glGetVertexAttribIiv"))
	if gpGetVertexAttribIiv == nil {
		return errors.New("glGetVertexAttribIiv")
	}
	gpGetVertexAttribIivEXT = (C.GPGETVERTEXATTRIBIIVEXT)(getProcAddr("glGetVertexAttribIivEXT"))
	gpGetVertexAttribIuiv = (C.GPGETVERTEXATTRIBIUIV)(getProcAddr("glGetVertexAttribIuiv"))
	if gpGetVertexAttribIuiv == nil {
		return errors.New("glGetVertexAttribIuiv")
	}
	gpGetVertexAttribIuivEXT = (C.GPGETVERTEXATTRIBIUIVEXT)(getProcAddr("glGetVertexAttribIuivEXT"))
	gpGetVertexAttribLdv = (C.GPGETVERTEXATTRIBLDV)(getProcAddr("glGetVertexAttribLdv"))
	if gpGetVertexAttribLdv == nil {
		return errors.New("glGetVertexAttribLdv")
	}
	gpGetVertexAttribLdvEXT = (C.GPGETVERTEXATTRIBLDVEXT)(getProcAddr("glGetVertexAttribLdvEXT"))
	gpGetVertexAttribLi64vNV = (C.GPGETVERTEXATTRIBLI64VNV)(getProcAddr("glGetVertexAttribLi64vNV"))
	gpGetVertexAttribLui64vARB = (C.GPGETVERTEXATTRIBLUI64VARB)(getProcAddr("glGetVertexAttribLui64vARB"))
	gpGetVertexAttribLui64vNV = (C.GPGETVERTEXATTRIBLUI64VNV)(getProcAddr("glGetVertexAttribLui64vNV"))
	gpGetVertexAttribPointerv = (C.GPGETVERTEXATTRIBPOINTERV)(getProcAddr("glGetVertexAttribPointerv"))
	if gpGetVertexAttribPointerv == nil {
		return errors.New("glGetVertexAttribPointerv")
	}
	gpGetVertexAttribPointervARB = (C.GPGETVERTEXATTRIBPOINTERVARB)(getProcAddr("glGetVertexAttribPointervARB"))
	gpGetVertexAttribPointervNV = (C.GPGETVERTEXATTRIBPOINTERVNV)(getProcAddr("glGetVertexAttribPointervNV"))
	gpGetVertexAttribdv = (C.GPGETVERTEXATTRIBDV)(getProcAddr("glGetVertexAttribdv"))
	if gpGetVertexAttribdv == nil {
		return errors.New("glGetVertexAttribdv")
	}
	gpGetVertexAttribdvARB = (C.GPGETVERTEXATTRIBDVARB)(getProcAddr("glGetVertexAttribdvARB"))
	gpGetVertexAttribdvNV = (C.GPGETVERTEXATTRIBDVNV)(getProcAddr("glGetVertexAttribdvNV"))
	gpGetVertexAttribfv = (C.GPGETVERTEXATTRIBFV)(getProcAddr("glGetVertexAttribfv"))
	if gpGetVertexAttribfv == nil {
		return errors.New("glGetVertexAttribfv")
	}
	gpGetVertexAttribfvARB = (C.GPGETVERTEXATTRIBFVARB)(getProcAddr("glGetVertexAttribfvARB"))
	gpGetVertexAttribfvNV = (C.GPGETVERTEXATTRIBFVNV)(getProcAddr("glGetVertexAttribfvNV"))
	gpGetVertexAttribiv = (C.GPGETVERTEXATTRIBIV)(getProcAddr("glGetVertexAttribiv"))
	if gpGetVertexAttribiv == nil {
		return errors.New("glGetVertexAttribiv")
	}
	gpGetVertexAttribivARB = (C.GPGETVERTEXATTRIBIVARB)(getProcAddr("glGetVertexAttribivARB"))
	gpGetVertexAttribivNV = (C.GPGETVERTEXATTRIBIVNV)(getProcAddr("glGetVertexAttribivNV"))
	gpGetVideoCaptureStreamdvNV = (C.GPGETVIDEOCAPTURESTREAMDVNV)(getProcAddr("glGetVideoCaptureStreamdvNV"))
	gpGetVideoCaptureStreamfvNV = (C.GPGETVIDEOCAPTURESTREAMFVNV)(getProcAddr("glGetVideoCaptureStreamfvNV"))
	gpGetVideoCaptureStreamivNV = (C.GPGETVIDEOCAPTURESTREAMIVNV)(getProcAddr("glGetVideoCaptureStreamivNV"))
	gpGetVideoCaptureivNV = (C.GPGETVIDEOCAPTUREIVNV)(getProcAddr("glGetVideoCaptureivNV"))
	gpGetVideoi64vNV = (C.GPGETVIDEOI64VNV)(getProcAddr("glGetVideoi64vNV"))
	gpGetVideoivNV = (C.GPGETVIDEOIVNV)(getProcAddr("glGetVideoivNV"))
	gpGetVideoui64vNV = (C.GPGETVIDEOUI64VNV)(getProcAddr("glGetVideoui64vNV"))
	gpGetVideouivNV = (C.GPGETVIDEOUIVNV)(getProcAddr("glGetVideouivNV"))
	gpGetVkProcAddrNV = (C.GPGETVKPROCADDRNV)(getProcAddr("glGetVkProcAddrNV"))
	gpGetnColorTable = (C.GPGETNCOLORTABLE)(getProcAddr("glGetnColorTable"))
	if gpGetnColorTable == nil {
		return errors.New("glGetnColorTable")
	}
	gpGetnColorTableARB = (C.GPGETNCOLORTABLEARB)(getProcAddr("glGetnColorTableARB"))
	gpGetnCompressedTexImage = (C.GPGETNCOMPRESSEDTEXIMAGE)(getProcAddr("glGetnCompressedTexImage"))
	if gpGetnCompressedTexImage == nil {
		return errors.New("glGetnCompressedTexImage")
	}
	gpGetnCompressedTexImageARB = (C.GPGETNCOMPRESSEDTEXIMAGEARB)(getProcAddr("glGetnCompressedTexImageARB"))
	gpGetnConvolutionFilter = (C.GPGETNCONVOLUTIONFILTER)(getProcAddr("glGetnConvolutionFilter"))
	if gpGetnConvolutionFilter == nil {
		return errors.New("glGetnConvolutionFilter")
	}
	gpGetnConvolutionFilterARB = (C.GPGETNCONVOLUTIONFILTERARB)(getProcAddr("glGetnConvolutionFilterARB"))
	gpGetnHistogram = (C.GPGETNHISTOGRAM)(getProcAddr("glGetnHistogram"))
	if gpGetnHistogram == nil {
		return errors.New("glGetnHistogram")
	}
	gpGetnHistogramARB = (C.GPGETNHISTOGRAMARB)(getProcAddr("glGetnHistogramARB"))
	gpGetnMapdv = (C.GPGETNMAPDV)(getProcAddr("glGetnMapdv"))
	if gpGetnMapdv == nil {
		return errors.New("glGetnMapdv")
	}
	gpGetnMapdvARB = (C.GPGETNMAPDVARB)(getProcAddr("glGetnMapdvARB"))
	gpGetnMapfv = (C.GPGETNMAPFV)(getProcAddr("glGetnMapfv"))
	if gpGetnMapfv == nil {
		return errors.New("glGetnMapfv")
	}
	gpGetnMapfvARB = (C.GPGETNMAPFVARB)(getProcAddr("glGetnMapfvARB"))
	gpGetnMapiv = (C.GPGETNMAPIV)(getProcAddr("glGetnMapiv"))
	if gpGetnMapiv == nil {
		return errors.New("glGetnMapiv")
	}
	gpGetnMapivARB = (C.GPGETNMAPIVARB)(getProcAddr("glGetnMapivARB"))
	gpGetnMinmax = (C.GPGETNMINMAX)(getProcAddr("glGetnMinmax"))
	if gpGetnMinmax == nil {
		return errors.New("glGetnMinmax")
	}
	gpGetnMinmaxARB = (C.GPGETNMINMAXARB)(getProcAddr("glGetnMinmaxARB"))
	gpGetnPixelMapfv = (C.GPGETNPIXELMAPFV)(getProcAddr("glGetnPixelMapfv"))
	if gpGetnPixelMapfv == nil {
		return errors.New("glGetnPixelMapfv")
	}
	gpGetnPixelMapfvARB = (C.GPGETNPIXELMAPFVARB)(getProcAddr("glGetnPixelMapfvARB"))
	gpGetnPixelMapuiv = (C.GPGETNPIXELMAPUIV)(getProcAddr("glGetnPixelMapuiv"))
	if gpGetnPixelMapuiv == nil {
		return errors.New("glGetnPixelMapuiv")
	}
	gpGetnPixelMapuivARB = (C.GPGETNPIXELMAPUIVARB)(getProcAddr("glGetnPixelMapuivARB"))
	gpGetnPixelMapusv = (C.GPGETNPIXELMAPUSV)(getProcAddr("glGetnPixelMapusv"))
	if gpGetnPixelMapusv == nil {
		return errors.New("glGetnPixelMapusv")
	}
	gpGetnPixelMapusvARB = (C.GPGETNPIXELMAPUSVARB)(getProcAddr("glGetnPixelMapusvARB"))
	gpGetnPolygonStipple = (C.GPGETNPOLYGONSTIPPLE)(getProcAddr("glGetnPolygonStipple"))
	if gpGetnPolygonStipple == nil {
		return errors.New("glGetnPolygonStipple")
	}
	gpGetnPolygonStippleARB = (C.GPGETNPOLYGONSTIPPLEARB)(getProcAddr("glGetnPolygonStippleARB"))
	gpGetnSeparableFilter = (C.GPGETNSEPARABLEFILTER)(getProcAddr("glGetnSeparableFilter"))
	if gpGetnSeparableFilter == nil {
		return errors.New("glGetnSeparableFilter")
	}
	gpGetnSeparableFilterARB = (C.GPGETNSEPARABLEFILTERARB)(getProcAddr("glGetnSeparableFilterARB"))
	gpGetnTexImage = (C.GPGETNTEXIMAGE)(getProcAddr("glGetnTexImage"))
	if gpGetnTexImage == nil {
		return errors.New("glGetnTexImage")
	}
	gpGetnTexImageARB = (C.GPGETNTEXIMAGEARB)(getProcAddr("glGetnTexImageARB"))
	gpGetnUniformdv = (C.GPGETNUNIFORMDV)(getProcAddr("glGetnUniformdv"))
	if gpGetnUniformdv == nil {
		return errors.New("glGetnUniformdv")
	}
	gpGetnUniformdvARB = (C.GPGETNUNIFORMDVARB)(getProcAddr("glGetnUniformdvARB"))
	gpGetnUniformfv = (C.GPGETNUNIFORMFV)(getProcAddr("glGetnUniformfv"))
	if gpGetnUniformfv == nil {
		return errors.New("glGetnUniformfv")
	}
	gpGetnUniformfvARB = (C.GPGETNUNIFORMFVARB)(getProcAddr("glGetnUniformfvARB"))
	gpGetnUniformfvKHR = (C.GPGETNUNIFORMFVKHR)(getProcAddr("glGetnUniformfvKHR"))
	gpGetnUniformi64vARB = (C.GPGETNUNIFORMI64VARB)(getProcAddr("glGetnUniformi64vARB"))
	gpGetnUniformiv = (C.GPGETNUNIFORMIV)(getProcAddr("glGetnUniformiv"))
	if gpGetnUniformiv == nil {
		return errors.New("glGetnUniformiv")
	}
	gpGetnUniformivARB = (C.GPGETNUNIFORMIVARB)(getProcAddr("glGetnUniformivARB"))
	gpGetnUniformivKHR = (C.GPGETNUNIFORMIVKHR)(getProcAddr("glGetnUniformivKHR"))
	gpGetnUniformui64vARB = (C.GPGETNUNIFORMUI64VARB)(getProcAddr("glGetnUniformui64vARB"))
	gpGetnUniformuiv = (C.GPGETNUNIFORMUIV)(getProcAddr("glGetnUniformuiv"))
	if gpGetnUniformuiv == nil {
		return errors.New("glGetnUniformuiv")
	}
	gpGetnUniformuivARB = (C.GPGETNUNIFORMUIVARB)(getProcAddr("glGetnUniformuivARB"))
	gpGetnUniformuivKHR = (C.GPGETNUNIFORMUIVKHR)(getProcAddr("glGetnUniformuivKHR"))
	gpGlobalAlphaFactorbSUN = (C.GPGLOBALALPHAFACTORBSUN)(getProcAddr("glGlobalAlphaFactorbSUN"))
	gpGlobalAlphaFactordSUN = (C.GPGLOBALALPHAFACTORDSUN)(getProcAddr("glGlobalAlphaFactordSUN"))
	gpGlobalAlphaFactorfSUN = (C.GPGLOBALALPHAFACTORFSUN)(getProcAddr("glGlobalAlphaFactorfSUN"))
	gpGlobalAlphaFactoriSUN = (C.GPGLOBALALPHAFACTORISUN)(getProcAddr("glGlobalAlphaFactoriSUN"))
	gpGlobalAlphaFactorsSUN = (C.GPGLOBALALPHAFACTORSSUN)(getProcAddr("glGlobalAlphaFactorsSUN"))
	gpGlobalAlphaFactorubSUN = (C.GPGLOBALALPHAFACTORUBSUN)(getProcAddr("glGlobalAlphaFactorubSUN"))
	gpGlobalAlphaFactoruiSUN = (C.GPGLOBALALPHAFACTORUISUN)(getProcAddr("glGlobalAlphaFactoruiSUN"))
	gpGlobalAlphaFactorusSUN = (C.GPGLOBALALPHAFACTORUSSUN)(getProcAddr("glGlobalAlphaFactorusSUN"))
	gpHint = (C.GPHINT)(getProcAddr("glHint"))
	if gpHint == nil {
		return errors.New("glHint")
	}
	gpHintPGI = (C.GPHINTPGI)(getProcAddr("glHintPGI"))
	gpHistogram = (C.GPHISTOGRAM)(getProcAddr("glHistogram"))
	gpHistogramEXT = (C.GPHISTOGRAMEXT)(getProcAddr("glHistogramEXT"))
	gpIglooInterfaceSGIX = (C.GPIGLOOINTERFACESGIX)(getProcAddr("glIglooInterfaceSGIX"))
	gpImageTransformParameterfHP = (C.GPIMAGETRANSFORMPARAMETERFHP)(getProcAddr("glImageTransformParameterfHP"))
	gpImageTransformParameterfvHP = (C.GPIMAGETRANSFORMPARAMETERFVHP)(getProcAddr("glImageTransformParameterfvHP"))
	gpImageTransformParameteriHP = (C.GPIMAGETRANSFORMPARAMETERIHP)(getProcAddr("glImageTransformParameteriHP"))
	gpImageTransformParameterivHP = (C.GPIMAGETRANSFORMPARAMETERIVHP)(getProcAddr("glImageTransformParameterivHP"))
	gpImportMemoryFdEXT = (C.GPIMPORTMEMORYFDEXT)(getProcAddr("glImportMemoryFdEXT"))
	gpImportMemoryWin32HandleEXT = (C.GPIMPORTMEMORYWIN32HANDLEEXT)(getProcAddr("glImportMemoryWin32HandleEXT"))
	gpImportMemoryWin32NameEXT = (C.GPIMPORTMEMORYWIN32NAMEEXT)(getProcAddr("glImportMemoryWin32NameEXT"))
	gpImportSemaphoreFdEXT = (C.GPIMPORTSEMAPHOREFDEXT)(getProcAddr("glImportSemaphoreFdEXT"))
	gpImportSemaphoreWin32HandleEXT = (C.GPIMPORTSEMAPHOREWIN32HANDLEEXT)(getProcAddr("glImportSemaphoreWin32HandleEXT"))
	gpImportSemaphoreWin32NameEXT = (C.GPIMPORTSEMAPHOREWIN32NAMEEXT)(getProcAddr("glImportSemaphoreWin32NameEXT"))
	gpImportSyncEXT = (C.GPIMPORTSYNCEXT)(getProcAddr("glImportSyncEXT"))
	gpIndexFormatNV = (C.GPINDEXFORMATNV)(getProcAddr("glIndexFormatNV"))
	gpIndexFuncEXT = (C.GPINDEXFUNCEXT)(getProcAddr("glIndexFuncEXT"))
	gpIndexMask = (C.GPINDEXMASK)(getProcAddr("glIndexMask"))
	if gpIndexMask == nil {
		return errors.New("glIndexMask")
	}
	gpIndexMaterialEXT = (C.GPINDEXMATERIALEXT)(getProcAddr("glIndexMaterialEXT"))
	gpIndexPointer = (C.GPINDEXPOINTER)(getProcAddr("glIndexPointer"))
	if gpIndexPointer == nil {
		return errors.New("glIndexPointer")
	}
	gpIndexPointerEXT = (C.GPINDEXPOINTEREXT)(getProcAddr("glIndexPointerEXT"))
	gpIndexPointerListIBM = (C.GPINDEXPOINTERLISTIBM)(getProcAddr("glIndexPointerListIBM"))
	gpIndexd = (C.GPINDEXD)(getProcAddr("glIndexd"))
	if gpIndexd == nil {
		return errors.New("glIndexd")
	}
	gpIndexdv = (C.GPINDEXDV)(getProcAddr("glIndexdv"))
	if gpIndexdv == nil {
		return errors.New("glIndexdv")
	}
	gpIndexf = (C.GPINDEXF)(getProcAddr("glIndexf"))
	if gpIndexf == nil {
		return errors.New("glIndexf")
	}
	gpIndexfv = (C.GPINDEXFV)(getProcAddr("glIndexfv"))
	if gpIndexfv == nil {
		return errors.New("glIndexfv")
	}
	gpIndexi = (C.GPINDEXI)(getProcAddr("glIndexi"))
	if gpIndexi == nil {
		return errors.New("glIndexi")
	}
	gpIndexiv = (C.GPINDEXIV)(getProcAddr("glIndexiv"))
	if gpIndexiv == nil {
		return errors.New("glIndexiv")
	}
	gpIndexs = (C.GPINDEXS)(getProcAddr("glIndexs"))
	if gpIndexs == nil {
		return errors.New("glIndexs")
	}
	gpIndexsv = (C.GPINDEXSV)(getProcAddr("glIndexsv"))
	if gpIndexsv == nil {
		return errors.New("glIndexsv")
	}
	gpIndexub = (C.GPINDEXUB)(getProcAddr("glIndexub"))
	if gpIndexub == nil {
		return errors.New("glIndexub")
	}
	gpIndexubv = (C.GPINDEXUBV)(getProcAddr("glIndexubv"))
	if gpIndexubv == nil {
		return errors.New("glIndexubv")
	}
	gpIndexxOES = (C.GPINDEXXOES)(getProcAddr("glIndexxOES"))
	gpIndexxvOES = (C.GPINDEXXVOES)(getProcAddr("glIndexxvOES"))
	gpInitNames = (C.GPINITNAMES)(getProcAddr("glInitNames"))
	if gpInitNames == nil {
		return errors.New("glInitNames")
	}
	gpInsertComponentEXT = (C.GPINSERTCOMPONENTEXT)(getProcAddr("glInsertComponentEXT"))
	gpInsertEventMarkerEXT = (C.GPINSERTEVENTMARKEREXT)(getProcAddr("glInsertEventMarkerEXT"))
	gpInstrumentsBufferSGIX = (C.GPINSTRUMENTSBUFFERSGIX)(getProcAddr("glInstrumentsBufferSGIX"))
	gpInterleavedArrays = (C.GPINTERLEAVEDARRAYS)(getProcAddr("glInterleavedArrays"))
	if gpInterleavedArrays == nil {
		return errors.New("glInterleavedArrays")
	}
	gpInterpolatePathsNV = (C.GPINTERPOLATEPATHSNV)(getProcAddr("glInterpolatePathsNV"))
	gpInvalidateBufferData = (C.GPINVALIDATEBUFFERDATA)(getProcAddr("glInvalidateBufferData"))
	if gpInvalidateBufferData == nil {
		return errors.New("glInvalidateBufferData")
	}
	gpInvalidateBufferSubData = (C.GPINVALIDATEBUFFERSUBDATA)(getProcAddr("glInvalidateBufferSubData"))
	if gpInvalidateBufferSubData == nil {
		return errors.New("glInvalidateBufferSubData")
	}
	gpInvalidateFramebuffer = (C.GPINVALIDATEFRAMEBUFFER)(getProcAddr("glInvalidateFramebuffer"))
	if gpInvalidateFramebuffer == nil {
		return errors.New("glInvalidateFramebuffer")
	}
	gpInvalidateNamedFramebufferData = (C.GPINVALIDATENAMEDFRAMEBUFFERDATA)(getProcAddr("glInvalidateNamedFramebufferData"))
	if gpInvalidateNamedFramebufferData == nil {
		return errors.New("glInvalidateNamedFramebufferData")
	}
	gpInvalidateNamedFramebufferSubData = (C.GPINVALIDATENAMEDFRAMEBUFFERSUBDATA)(getProcAddr("glInvalidateNamedFramebufferSubData"))
	if gpInvalidateNamedFramebufferSubData == nil {
		return errors.New("glInvalidateNamedFramebufferSubData")
	}
	gpInvalidateSubFramebuffer = (C.GPINVALIDATESUBFRAMEBUFFER)(getProcAddr("glInvalidateSubFramebuffer"))
	if gpInvalidateSubFramebuffer == nil {
		return errors.New("glInvalidateSubFramebuffer")
	}
	gpInvalidateTexImage = (C.GPINVALIDATETEXIMAGE)(getProcAddr("glInvalidateTexImage"))
	if gpInvalidateTexImage == nil {
		return errors.New("glInvalidateTexImage")
	}
	gpInvalidateTexSubImage = (C.GPINVALIDATETEXSUBIMAGE)(getProcAddr("glInvalidateTexSubImage"))
	if gpInvalidateTexSubImage == nil {
		return errors.New("glInvalidateTexSubImage")
	}
	gpIsAsyncMarkerSGIX = (C.GPISASYNCMARKERSGIX)(getProcAddr("glIsAsyncMarkerSGIX"))
	gpIsBuffer = (C.GPISBUFFER)(getProcAddr("glIsBuffer"))
	if gpIsBuffer == nil {
		return errors.New("glIsBuffer")
	}
	gpIsBufferARB = (C.GPISBUFFERARB)(getProcAddr("glIsBufferARB"))
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
	gpIsFenceAPPLE = (C.GPISFENCEAPPLE)(getProcAddr("glIsFenceAPPLE"))
	gpIsFenceNV = (C.GPISFENCENV)(getProcAddr("glIsFenceNV"))
	gpIsFramebuffer = (C.GPISFRAMEBUFFER)(getProcAddr("glIsFramebuffer"))
	if gpIsFramebuffer == nil {
		return errors.New("glIsFramebuffer")
	}
	gpIsFramebufferEXT = (C.GPISFRAMEBUFFEREXT)(getProcAddr("glIsFramebufferEXT"))
	gpIsImageHandleResidentARB = (C.GPISIMAGEHANDLERESIDENTARB)(getProcAddr("glIsImageHandleResidentARB"))
	gpIsImageHandleResidentNV = (C.GPISIMAGEHANDLERESIDENTNV)(getProcAddr("glIsImageHandleResidentNV"))
	gpIsList = (C.GPISLIST)(getProcAddr("glIsList"))
	if gpIsList == nil {
		return errors.New("glIsList")
	}
	gpIsMemoryObjectEXT = (C.GPISMEMORYOBJECTEXT)(getProcAddr("glIsMemoryObjectEXT"))
	gpIsNameAMD = (C.GPISNAMEAMD)(getProcAddr("glIsNameAMD"))
	gpIsNamedBufferResidentNV = (C.GPISNAMEDBUFFERRESIDENTNV)(getProcAddr("glIsNamedBufferResidentNV"))
	gpIsNamedStringARB = (C.GPISNAMEDSTRINGARB)(getProcAddr("glIsNamedStringARB"))
	gpIsObjectBufferATI = (C.GPISOBJECTBUFFERATI)(getProcAddr("glIsObjectBufferATI"))
	gpIsOcclusionQueryNV = (C.GPISOCCLUSIONQUERYNV)(getProcAddr("glIsOcclusionQueryNV"))
	gpIsPathNV = (C.GPISPATHNV)(getProcAddr("glIsPathNV"))
	gpIsPointInFillPathNV = (C.GPISPOINTINFILLPATHNV)(getProcAddr("glIsPointInFillPathNV"))
	gpIsPointInStrokePathNV = (C.GPISPOINTINSTROKEPATHNV)(getProcAddr("glIsPointInStrokePathNV"))
	gpIsProgram = (C.GPISPROGRAM)(getProcAddr("glIsProgram"))
	if gpIsProgram == nil {
		return errors.New("glIsProgram")
	}
	gpIsProgramARB = (C.GPISPROGRAMARB)(getProcAddr("glIsProgramARB"))
	gpIsProgramNV = (C.GPISPROGRAMNV)(getProcAddr("glIsProgramNV"))
	gpIsProgramPipeline = (C.GPISPROGRAMPIPELINE)(getProcAddr("glIsProgramPipeline"))
	if gpIsProgramPipeline == nil {
		return errors.New("glIsProgramPipeline")
	}
	gpIsProgramPipelineEXT = (C.GPISPROGRAMPIPELINEEXT)(getProcAddr("glIsProgramPipelineEXT"))
	gpIsQuery = (C.GPISQUERY)(getProcAddr("glIsQuery"))
	if gpIsQuery == nil {
		return errors.New("glIsQuery")
	}
	gpIsQueryARB = (C.GPISQUERYARB)(getProcAddr("glIsQueryARB"))
	gpIsRenderbuffer = (C.GPISRENDERBUFFER)(getProcAddr("glIsRenderbuffer"))
	if gpIsRenderbuffer == nil {
		return errors.New("glIsRenderbuffer")
	}
	gpIsRenderbufferEXT = (C.GPISRENDERBUFFEREXT)(getProcAddr("glIsRenderbufferEXT"))
	gpIsSampler = (C.GPISSAMPLER)(getProcAddr("glIsSampler"))
	if gpIsSampler == nil {
		return errors.New("glIsSampler")
	}
	gpIsSemaphoreEXT = (C.GPISSEMAPHOREEXT)(getProcAddr("glIsSemaphoreEXT"))
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
	gpIsTextureEXT = (C.GPISTEXTUREEXT)(getProcAddr("glIsTextureEXT"))
	gpIsTextureHandleResidentARB = (C.GPISTEXTUREHANDLERESIDENTARB)(getProcAddr("glIsTextureHandleResidentARB"))
	gpIsTextureHandleResidentNV = (C.GPISTEXTUREHANDLERESIDENTNV)(getProcAddr("glIsTextureHandleResidentNV"))
	gpIsTransformFeedback = (C.GPISTRANSFORMFEEDBACK)(getProcAddr("glIsTransformFeedback"))
	if gpIsTransformFeedback == nil {
		return errors.New("glIsTransformFeedback")
	}
	gpIsTransformFeedbackNV = (C.GPISTRANSFORMFEEDBACKNV)(getProcAddr("glIsTransformFeedbackNV"))
	gpIsVariantEnabledEXT = (C.GPISVARIANTENABLEDEXT)(getProcAddr("glIsVariantEnabledEXT"))
	gpIsVertexArray = (C.GPISVERTEXARRAY)(getProcAddr("glIsVertexArray"))
	if gpIsVertexArray == nil {
		return errors.New("glIsVertexArray")
	}
	gpIsVertexArrayAPPLE = (C.GPISVERTEXARRAYAPPLE)(getProcAddr("glIsVertexArrayAPPLE"))
	gpIsVertexAttribEnabledAPPLE = (C.GPISVERTEXATTRIBENABLEDAPPLE)(getProcAddr("glIsVertexAttribEnabledAPPLE"))
	gpLGPUCopyImageSubDataNVX = (C.GPLGPUCOPYIMAGESUBDATANVX)(getProcAddr("glLGPUCopyImageSubDataNVX"))
	gpLGPUInterlockNVX = (C.GPLGPUINTERLOCKNVX)(getProcAddr("glLGPUInterlockNVX"))
	gpLGPUNamedBufferSubDataNVX = (C.GPLGPUNAMEDBUFFERSUBDATANVX)(getProcAddr("glLGPUNamedBufferSubDataNVX"))
	gpLabelObjectEXT = (C.GPLABELOBJECTEXT)(getProcAddr("glLabelObjectEXT"))
	gpLightEnviSGIX = (C.GPLIGHTENVISGIX)(getProcAddr("glLightEnviSGIX"))
	gpLightModelf = (C.GPLIGHTMODELF)(getProcAddr("glLightModelf"))
	if gpLightModelf == nil {
		return errors.New("glLightModelf")
	}
	gpLightModelfv = (C.GPLIGHTMODELFV)(getProcAddr("glLightModelfv"))
	if gpLightModelfv == nil {
		return errors.New("glLightModelfv")
	}
	gpLightModeli = (C.GPLIGHTMODELI)(getProcAddr("glLightModeli"))
	if gpLightModeli == nil {
		return errors.New("glLightModeli")
	}
	gpLightModeliv = (C.GPLIGHTMODELIV)(getProcAddr("glLightModeliv"))
	if gpLightModeliv == nil {
		return errors.New("glLightModeliv")
	}
	gpLightModelxOES = (C.GPLIGHTMODELXOES)(getProcAddr("glLightModelxOES"))
	gpLightModelxvOES = (C.GPLIGHTMODELXVOES)(getProcAddr("glLightModelxvOES"))
	gpLightf = (C.GPLIGHTF)(getProcAddr("glLightf"))
	if gpLightf == nil {
		return errors.New("glLightf")
	}
	gpLightfv = (C.GPLIGHTFV)(getProcAddr("glLightfv"))
	if gpLightfv == nil {
		return errors.New("glLightfv")
	}
	gpLighti = (C.GPLIGHTI)(getProcAddr("glLighti"))
	if gpLighti == nil {
		return errors.New("glLighti")
	}
	gpLightiv = (C.GPLIGHTIV)(getProcAddr("glLightiv"))
	if gpLightiv == nil {
		return errors.New("glLightiv")
	}
	gpLightxOES = (C.GPLIGHTXOES)(getProcAddr("glLightxOES"))
	gpLightxvOES = (C.GPLIGHTXVOES)(getProcAddr("glLightxvOES"))
	gpLineStipple = (C.GPLINESTIPPLE)(getProcAddr("glLineStipple"))
	if gpLineStipple == nil {
		return errors.New("glLineStipple")
	}
	gpLineWidth = (C.GPLINEWIDTH)(getProcAddr("glLineWidth"))
	if gpLineWidth == nil {
		return errors.New("glLineWidth")
	}
	gpLineWidthxOES = (C.GPLINEWIDTHXOES)(getProcAddr("glLineWidthxOES"))
	gpLinkProgram = (C.GPLINKPROGRAM)(getProcAddr("glLinkProgram"))
	if gpLinkProgram == nil {
		return errors.New("glLinkProgram")
	}
	gpLinkProgramARB = (C.GPLINKPROGRAMARB)(getProcAddr("glLinkProgramARB"))
	gpListBase = (C.GPLISTBASE)(getProcAddr("glListBase"))
	if gpListBase == nil {
		return errors.New("glListBase")
	}
	gpListDrawCommandsStatesClientNV = (C.GPLISTDRAWCOMMANDSSTATESCLIENTNV)(getProcAddr("glListDrawCommandsStatesClientNV"))
	gpListParameterfSGIX = (C.GPLISTPARAMETERFSGIX)(getProcAddr("glListParameterfSGIX"))
	gpListParameterfvSGIX = (C.GPLISTPARAMETERFVSGIX)(getProcAddr("glListParameterfvSGIX"))
	gpListParameteriSGIX = (C.GPLISTPARAMETERISGIX)(getProcAddr("glListParameteriSGIX"))
	gpListParameterivSGIX = (C.GPLISTPARAMETERIVSGIX)(getProcAddr("glListParameterivSGIX"))
	gpLoadIdentity = (C.GPLOADIDENTITY)(getProcAddr("glLoadIdentity"))
	if gpLoadIdentity == nil {
		return errors.New("glLoadIdentity")
	}
	gpLoadIdentityDeformationMapSGIX = (C.GPLOADIDENTITYDEFORMATIONMAPSGIX)(getProcAddr("glLoadIdentityDeformationMapSGIX"))
	gpLoadMatrixd = (C.GPLOADMATRIXD)(getProcAddr("glLoadMatrixd"))
	if gpLoadMatrixd == nil {
		return errors.New("glLoadMatrixd")
	}
	gpLoadMatrixf = (C.GPLOADMATRIXF)(getProcAddr("glLoadMatrixf"))
	if gpLoadMatrixf == nil {
		return errors.New("glLoadMatrixf")
	}
	gpLoadMatrixxOES = (C.GPLOADMATRIXXOES)(getProcAddr("glLoadMatrixxOES"))
	gpLoadName = (C.GPLOADNAME)(getProcAddr("glLoadName"))
	if gpLoadName == nil {
		return errors.New("glLoadName")
	}
	gpLoadProgramNV = (C.GPLOADPROGRAMNV)(getProcAddr("glLoadProgramNV"))
	gpLoadTransposeMatrixd = (C.GPLOADTRANSPOSEMATRIXD)(getProcAddr("glLoadTransposeMatrixd"))
	if gpLoadTransposeMatrixd == nil {
		return errors.New("glLoadTransposeMatrixd")
	}
	gpLoadTransposeMatrixdARB = (C.GPLOADTRANSPOSEMATRIXDARB)(getProcAddr("glLoadTransposeMatrixdARB"))
	gpLoadTransposeMatrixf = (C.GPLOADTRANSPOSEMATRIXF)(getProcAddr("glLoadTransposeMatrixf"))
	if gpLoadTransposeMatrixf == nil {
		return errors.New("glLoadTransposeMatrixf")
	}
	gpLoadTransposeMatrixfARB = (C.GPLOADTRANSPOSEMATRIXFARB)(getProcAddr("glLoadTransposeMatrixfARB"))
	gpLoadTransposeMatrixxOES = (C.GPLOADTRANSPOSEMATRIXXOES)(getProcAddr("glLoadTransposeMatrixxOES"))
	gpLockArraysEXT = (C.GPLOCKARRAYSEXT)(getProcAddr("glLockArraysEXT"))
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
	gpMap1d = (C.GPMAP1D)(getProcAddr("glMap1d"))
	if gpMap1d == nil {
		return errors.New("glMap1d")
	}
	gpMap1f = (C.GPMAP1F)(getProcAddr("glMap1f"))
	if gpMap1f == nil {
		return errors.New("glMap1f")
	}
	gpMap1xOES = (C.GPMAP1XOES)(getProcAddr("glMap1xOES"))
	gpMap2d = (C.GPMAP2D)(getProcAddr("glMap2d"))
	if gpMap2d == nil {
		return errors.New("glMap2d")
	}
	gpMap2f = (C.GPMAP2F)(getProcAddr("glMap2f"))
	if gpMap2f == nil {
		return errors.New("glMap2f")
	}
	gpMap2xOES = (C.GPMAP2XOES)(getProcAddr("glMap2xOES"))
	gpMapBuffer = (C.GPMAPBUFFER)(getProcAddr("glMapBuffer"))
	if gpMapBuffer == nil {
		return errors.New("glMapBuffer")
	}
	gpMapBufferARB = (C.GPMAPBUFFERARB)(getProcAddr("glMapBufferARB"))
	gpMapBufferRange = (C.GPMAPBUFFERRANGE)(getProcAddr("glMapBufferRange"))
	if gpMapBufferRange == nil {
		return errors.New("glMapBufferRange")
	}
	gpMapControlPointsNV = (C.GPMAPCONTROLPOINTSNV)(getProcAddr("glMapControlPointsNV"))
	gpMapGrid1d = (C.GPMAPGRID1D)(getProcAddr("glMapGrid1d"))
	if gpMapGrid1d == nil {
		return errors.New("glMapGrid1d")
	}
	gpMapGrid1f = (C.GPMAPGRID1F)(getProcAddr("glMapGrid1f"))
	if gpMapGrid1f == nil {
		return errors.New("glMapGrid1f")
	}
	gpMapGrid1xOES = (C.GPMAPGRID1XOES)(getProcAddr("glMapGrid1xOES"))
	gpMapGrid2d = (C.GPMAPGRID2D)(getProcAddr("glMapGrid2d"))
	if gpMapGrid2d == nil {
		return errors.New("glMapGrid2d")
	}
	gpMapGrid2f = (C.GPMAPGRID2F)(getProcAddr("glMapGrid2f"))
	if gpMapGrid2f == nil {
		return errors.New("glMapGrid2f")
	}
	gpMapGrid2xOES = (C.GPMAPGRID2XOES)(getProcAddr("glMapGrid2xOES"))
	gpMapNamedBuffer = (C.GPMAPNAMEDBUFFER)(getProcAddr("glMapNamedBuffer"))
	if gpMapNamedBuffer == nil {
		return errors.New("glMapNamedBuffer")
	}
	gpMapNamedBufferEXT = (C.GPMAPNAMEDBUFFEREXT)(getProcAddr("glMapNamedBufferEXT"))
	gpMapNamedBufferRange = (C.GPMAPNAMEDBUFFERRANGE)(getProcAddr("glMapNamedBufferRange"))
	if gpMapNamedBufferRange == nil {
		return errors.New("glMapNamedBufferRange")
	}
	gpMapNamedBufferRangeEXT = (C.GPMAPNAMEDBUFFERRANGEEXT)(getProcAddr("glMapNamedBufferRangeEXT"))
	gpMapObjectBufferATI = (C.GPMAPOBJECTBUFFERATI)(getProcAddr("glMapObjectBufferATI"))
	gpMapParameterfvNV = (C.GPMAPPARAMETERFVNV)(getProcAddr("glMapParameterfvNV"))
	gpMapParameterivNV = (C.GPMAPPARAMETERIVNV)(getProcAddr("glMapParameterivNV"))
	gpMapTexture2DINTEL = (C.GPMAPTEXTURE2DINTEL)(getProcAddr("glMapTexture2DINTEL"))
	gpMapVertexAttrib1dAPPLE = (C.GPMAPVERTEXATTRIB1DAPPLE)(getProcAddr("glMapVertexAttrib1dAPPLE"))
	gpMapVertexAttrib1fAPPLE = (C.GPMAPVERTEXATTRIB1FAPPLE)(getProcAddr("glMapVertexAttrib1fAPPLE"))
	gpMapVertexAttrib2dAPPLE = (C.GPMAPVERTEXATTRIB2DAPPLE)(getProcAddr("glMapVertexAttrib2dAPPLE"))
	gpMapVertexAttrib2fAPPLE = (C.GPMAPVERTEXATTRIB2FAPPLE)(getProcAddr("glMapVertexAttrib2fAPPLE"))
	gpMaterialf = (C.GPMATERIALF)(getProcAddr("glMaterialf"))
	if gpMaterialf == nil {
		return errors.New("glMaterialf")
	}
	gpMaterialfv = (C.GPMATERIALFV)(getProcAddr("glMaterialfv"))
	if gpMaterialfv == nil {
		return errors.New("glMaterialfv")
	}
	gpMateriali = (C.GPMATERIALI)(getProcAddr("glMateriali"))
	if gpMateriali == nil {
		return errors.New("glMateriali")
	}
	gpMaterialiv = (C.GPMATERIALIV)(getProcAddr("glMaterialiv"))
	if gpMaterialiv == nil {
		return errors.New("glMaterialiv")
	}
	gpMaterialxOES = (C.GPMATERIALXOES)(getProcAddr("glMaterialxOES"))
	gpMaterialxvOES = (C.GPMATERIALXVOES)(getProcAddr("glMaterialxvOES"))
	gpMatrixFrustumEXT = (C.GPMATRIXFRUSTUMEXT)(getProcAddr("glMatrixFrustumEXT"))
	gpMatrixIndexPointerARB = (C.GPMATRIXINDEXPOINTERARB)(getProcAddr("glMatrixIndexPointerARB"))
	gpMatrixIndexubvARB = (C.GPMATRIXINDEXUBVARB)(getProcAddr("glMatrixIndexubvARB"))
	gpMatrixIndexuivARB = (C.GPMATRIXINDEXUIVARB)(getProcAddr("glMatrixIndexuivARB"))
	gpMatrixIndexusvARB = (C.GPMATRIXINDEXUSVARB)(getProcAddr("glMatrixIndexusvARB"))
	gpMatrixLoad3x2fNV = (C.GPMATRIXLOAD3X2FNV)(getProcAddr("glMatrixLoad3x2fNV"))
	gpMatrixLoad3x3fNV = (C.GPMATRIXLOAD3X3FNV)(getProcAddr("glMatrixLoad3x3fNV"))
	gpMatrixLoadIdentityEXT = (C.GPMATRIXLOADIDENTITYEXT)(getProcAddr("glMatrixLoadIdentityEXT"))
	gpMatrixLoadTranspose3x3fNV = (C.GPMATRIXLOADTRANSPOSE3X3FNV)(getProcAddr("glMatrixLoadTranspose3x3fNV"))
	gpMatrixLoadTransposedEXT = (C.GPMATRIXLOADTRANSPOSEDEXT)(getProcAddr("glMatrixLoadTransposedEXT"))
	gpMatrixLoadTransposefEXT = (C.GPMATRIXLOADTRANSPOSEFEXT)(getProcAddr("glMatrixLoadTransposefEXT"))
	gpMatrixLoaddEXT = (C.GPMATRIXLOADDEXT)(getProcAddr("glMatrixLoaddEXT"))
	gpMatrixLoadfEXT = (C.GPMATRIXLOADFEXT)(getProcAddr("glMatrixLoadfEXT"))
	gpMatrixMode = (C.GPMATRIXMODE)(getProcAddr("glMatrixMode"))
	if gpMatrixMode == nil {
		return errors.New("glMatrixMode")
	}
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
	if gpMemoryBarrier == nil {
		return errors.New("glMemoryBarrier")
	}
	gpMemoryBarrierByRegion = (C.GPMEMORYBARRIERBYREGION)(getProcAddr("glMemoryBarrierByRegion"))
	if gpMemoryBarrierByRegion == nil {
		return errors.New("glMemoryBarrierByRegion")
	}
	gpMemoryBarrierEXT = (C.GPMEMORYBARRIEREXT)(getProcAddr("glMemoryBarrierEXT"))
	gpMemoryObjectParameterivEXT = (C.GPMEMORYOBJECTPARAMETERIVEXT)(getProcAddr("glMemoryObjectParameterivEXT"))
	gpMinSampleShading = (C.GPMINSAMPLESHADING)(getProcAddr("glMinSampleShading"))
	if gpMinSampleShading == nil {
		return errors.New("glMinSampleShading")
	}
	gpMinSampleShadingARB = (C.GPMINSAMPLESHADINGARB)(getProcAddr("glMinSampleShadingARB"))
	gpMinmax = (C.GPMINMAX)(getProcAddr("glMinmax"))
	gpMinmaxEXT = (C.GPMINMAXEXT)(getProcAddr("glMinmaxEXT"))
	gpMultMatrixd = (C.GPMULTMATRIXD)(getProcAddr("glMultMatrixd"))
	if gpMultMatrixd == nil {
		return errors.New("glMultMatrixd")
	}
	gpMultMatrixf = (C.GPMULTMATRIXF)(getProcAddr("glMultMatrixf"))
	if gpMultMatrixf == nil {
		return errors.New("glMultMatrixf")
	}
	gpMultMatrixxOES = (C.GPMULTMATRIXXOES)(getProcAddr("glMultMatrixxOES"))
	gpMultTransposeMatrixd = (C.GPMULTTRANSPOSEMATRIXD)(getProcAddr("glMultTransposeMatrixd"))
	if gpMultTransposeMatrixd == nil {
		return errors.New("glMultTransposeMatrixd")
	}
	gpMultTransposeMatrixdARB = (C.GPMULTTRANSPOSEMATRIXDARB)(getProcAddr("glMultTransposeMatrixdARB"))
	gpMultTransposeMatrixf = (C.GPMULTTRANSPOSEMATRIXF)(getProcAddr("glMultTransposeMatrixf"))
	if gpMultTransposeMatrixf == nil {
		return errors.New("glMultTransposeMatrixf")
	}
	gpMultTransposeMatrixfARB = (C.GPMULTTRANSPOSEMATRIXFARB)(getProcAddr("glMultTransposeMatrixfARB"))
	gpMultTransposeMatrixxOES = (C.GPMULTTRANSPOSEMATRIXXOES)(getProcAddr("glMultTransposeMatrixxOES"))
	gpMultiDrawArrays = (C.GPMULTIDRAWARRAYS)(getProcAddr("glMultiDrawArrays"))
	if gpMultiDrawArrays == nil {
		return errors.New("glMultiDrawArrays")
	}
	gpMultiDrawArraysEXT = (C.GPMULTIDRAWARRAYSEXT)(getProcAddr("glMultiDrawArraysEXT"))
	gpMultiDrawArraysIndirect = (C.GPMULTIDRAWARRAYSINDIRECT)(getProcAddr("glMultiDrawArraysIndirect"))
	if gpMultiDrawArraysIndirect == nil {
		return errors.New("glMultiDrawArraysIndirect")
	}
	gpMultiDrawArraysIndirectAMD = (C.GPMULTIDRAWARRAYSINDIRECTAMD)(getProcAddr("glMultiDrawArraysIndirectAMD"))
	gpMultiDrawArraysIndirectBindlessCountNV = (C.GPMULTIDRAWARRAYSINDIRECTBINDLESSCOUNTNV)(getProcAddr("glMultiDrawArraysIndirectBindlessCountNV"))
	gpMultiDrawArraysIndirectBindlessNV = (C.GPMULTIDRAWARRAYSINDIRECTBINDLESSNV)(getProcAddr("glMultiDrawArraysIndirectBindlessNV"))
	gpMultiDrawArraysIndirectCountARB = (C.GPMULTIDRAWARRAYSINDIRECTCOUNTARB)(getProcAddr("glMultiDrawArraysIndirectCountARB"))
	gpMultiDrawElementArrayAPPLE = (C.GPMULTIDRAWELEMENTARRAYAPPLE)(getProcAddr("glMultiDrawElementArrayAPPLE"))
	gpMultiDrawElements = (C.GPMULTIDRAWELEMENTS)(getProcAddr("glMultiDrawElements"))
	if gpMultiDrawElements == nil {
		return errors.New("glMultiDrawElements")
	}
	gpMultiDrawElementsBaseVertex = (C.GPMULTIDRAWELEMENTSBASEVERTEX)(getProcAddr("glMultiDrawElementsBaseVertex"))
	if gpMultiDrawElementsBaseVertex == nil {
		return errors.New("glMultiDrawElementsBaseVertex")
	}
	gpMultiDrawElementsEXT = (C.GPMULTIDRAWELEMENTSEXT)(getProcAddr("glMultiDrawElementsEXT"))
	gpMultiDrawElementsIndirect = (C.GPMULTIDRAWELEMENTSINDIRECT)(getProcAddr("glMultiDrawElementsIndirect"))
	if gpMultiDrawElementsIndirect == nil {
		return errors.New("glMultiDrawElementsIndirect")
	}
	gpMultiDrawElementsIndirectAMD = (C.GPMULTIDRAWELEMENTSINDIRECTAMD)(getProcAddr("glMultiDrawElementsIndirectAMD"))
	gpMultiDrawElementsIndirectBindlessCountNV = (C.GPMULTIDRAWELEMENTSINDIRECTBINDLESSCOUNTNV)(getProcAddr("glMultiDrawElementsIndirectBindlessCountNV"))
	gpMultiDrawElementsIndirectBindlessNV = (C.GPMULTIDRAWELEMENTSINDIRECTBINDLESSNV)(getProcAddr("glMultiDrawElementsIndirectBindlessNV"))
	gpMultiDrawElementsIndirectCountARB = (C.GPMULTIDRAWELEMENTSINDIRECTCOUNTARB)(getProcAddr("glMultiDrawElementsIndirectCountARB"))
	gpMultiDrawRangeElementArrayAPPLE = (C.GPMULTIDRAWRANGEELEMENTARRAYAPPLE)(getProcAddr("glMultiDrawRangeElementArrayAPPLE"))
	gpMultiModeDrawArraysIBM = (C.GPMULTIMODEDRAWARRAYSIBM)(getProcAddr("glMultiModeDrawArraysIBM"))
	gpMultiModeDrawElementsIBM = (C.GPMULTIMODEDRAWELEMENTSIBM)(getProcAddr("glMultiModeDrawElementsIBM"))
	gpMultiTexBufferEXT = (C.GPMULTITEXBUFFEREXT)(getProcAddr("glMultiTexBufferEXT"))
	gpMultiTexCoord1bOES = (C.GPMULTITEXCOORD1BOES)(getProcAddr("glMultiTexCoord1bOES"))
	gpMultiTexCoord1bvOES = (C.GPMULTITEXCOORD1BVOES)(getProcAddr("glMultiTexCoord1bvOES"))
	gpMultiTexCoord1d = (C.GPMULTITEXCOORD1D)(getProcAddr("glMultiTexCoord1d"))
	if gpMultiTexCoord1d == nil {
		return errors.New("glMultiTexCoord1d")
	}
	gpMultiTexCoord1dARB = (C.GPMULTITEXCOORD1DARB)(getProcAddr("glMultiTexCoord1dARB"))
	gpMultiTexCoord1dv = (C.GPMULTITEXCOORD1DV)(getProcAddr("glMultiTexCoord1dv"))
	if gpMultiTexCoord1dv == nil {
		return errors.New("glMultiTexCoord1dv")
	}
	gpMultiTexCoord1dvARB = (C.GPMULTITEXCOORD1DVARB)(getProcAddr("glMultiTexCoord1dvARB"))
	gpMultiTexCoord1f = (C.GPMULTITEXCOORD1F)(getProcAddr("glMultiTexCoord1f"))
	if gpMultiTexCoord1f == nil {
		return errors.New("glMultiTexCoord1f")
	}
	gpMultiTexCoord1fARB = (C.GPMULTITEXCOORD1FARB)(getProcAddr("glMultiTexCoord1fARB"))
	gpMultiTexCoord1fv = (C.GPMULTITEXCOORD1FV)(getProcAddr("glMultiTexCoord1fv"))
	if gpMultiTexCoord1fv == nil {
		return errors.New("glMultiTexCoord1fv")
	}
	gpMultiTexCoord1fvARB = (C.GPMULTITEXCOORD1FVARB)(getProcAddr("glMultiTexCoord1fvARB"))
	gpMultiTexCoord1hNV = (C.GPMULTITEXCOORD1HNV)(getProcAddr("glMultiTexCoord1hNV"))
	gpMultiTexCoord1hvNV = (C.GPMULTITEXCOORD1HVNV)(getProcAddr("glMultiTexCoord1hvNV"))
	gpMultiTexCoord1i = (C.GPMULTITEXCOORD1I)(getProcAddr("glMultiTexCoord1i"))
	if gpMultiTexCoord1i == nil {
		return errors.New("glMultiTexCoord1i")
	}
	gpMultiTexCoord1iARB = (C.GPMULTITEXCOORD1IARB)(getProcAddr("glMultiTexCoord1iARB"))
	gpMultiTexCoord1iv = (C.GPMULTITEXCOORD1IV)(getProcAddr("glMultiTexCoord1iv"))
	if gpMultiTexCoord1iv == nil {
		return errors.New("glMultiTexCoord1iv")
	}
	gpMultiTexCoord1ivARB = (C.GPMULTITEXCOORD1IVARB)(getProcAddr("glMultiTexCoord1ivARB"))
	gpMultiTexCoord1s = (C.GPMULTITEXCOORD1S)(getProcAddr("glMultiTexCoord1s"))
	if gpMultiTexCoord1s == nil {
		return errors.New("glMultiTexCoord1s")
	}
	gpMultiTexCoord1sARB = (C.GPMULTITEXCOORD1SARB)(getProcAddr("glMultiTexCoord1sARB"))
	gpMultiTexCoord1sv = (C.GPMULTITEXCOORD1SV)(getProcAddr("glMultiTexCoord1sv"))
	if gpMultiTexCoord1sv == nil {
		return errors.New("glMultiTexCoord1sv")
	}
	gpMultiTexCoord1svARB = (C.GPMULTITEXCOORD1SVARB)(getProcAddr("glMultiTexCoord1svARB"))
	gpMultiTexCoord1xOES = (C.GPMULTITEXCOORD1XOES)(getProcAddr("glMultiTexCoord1xOES"))
	gpMultiTexCoord1xvOES = (C.GPMULTITEXCOORD1XVOES)(getProcAddr("glMultiTexCoord1xvOES"))
	gpMultiTexCoord2bOES = (C.GPMULTITEXCOORD2BOES)(getProcAddr("glMultiTexCoord2bOES"))
	gpMultiTexCoord2bvOES = (C.GPMULTITEXCOORD2BVOES)(getProcAddr("glMultiTexCoord2bvOES"))
	gpMultiTexCoord2d = (C.GPMULTITEXCOORD2D)(getProcAddr("glMultiTexCoord2d"))
	if gpMultiTexCoord2d == nil {
		return errors.New("glMultiTexCoord2d")
	}
	gpMultiTexCoord2dARB = (C.GPMULTITEXCOORD2DARB)(getProcAddr("glMultiTexCoord2dARB"))
	gpMultiTexCoord2dv = (C.GPMULTITEXCOORD2DV)(getProcAddr("glMultiTexCoord2dv"))
	if gpMultiTexCoord2dv == nil {
		return errors.New("glMultiTexCoord2dv")
	}
	gpMultiTexCoord2dvARB = (C.GPMULTITEXCOORD2DVARB)(getProcAddr("glMultiTexCoord2dvARB"))
	gpMultiTexCoord2f = (C.GPMULTITEXCOORD2F)(getProcAddr("glMultiTexCoord2f"))
	if gpMultiTexCoord2f == nil {
		return errors.New("glMultiTexCoord2f")
	}
	gpMultiTexCoord2fARB = (C.GPMULTITEXCOORD2FARB)(getProcAddr("glMultiTexCoord2fARB"))
	gpMultiTexCoord2fv = (C.GPMULTITEXCOORD2FV)(getProcAddr("glMultiTexCoord2fv"))
	if gpMultiTexCoord2fv == nil {
		return errors.New("glMultiTexCoord2fv")
	}
	gpMultiTexCoord2fvARB = (C.GPMULTITEXCOORD2FVARB)(getProcAddr("glMultiTexCoord2fvARB"))
	gpMultiTexCoord2hNV = (C.GPMULTITEXCOORD2HNV)(getProcAddr("glMultiTexCoord2hNV"))
	gpMultiTexCoord2hvNV = (C.GPMULTITEXCOORD2HVNV)(getProcAddr("glMultiTexCoord2hvNV"))
	gpMultiTexCoord2i = (C.GPMULTITEXCOORD2I)(getProcAddr("glMultiTexCoord2i"))
	if gpMultiTexCoord2i == nil {
		return errors.New("glMultiTexCoord2i")
	}
	gpMultiTexCoord2iARB = (C.GPMULTITEXCOORD2IARB)(getProcAddr("glMultiTexCoord2iARB"))
	gpMultiTexCoord2iv = (C.GPMULTITEXCOORD2IV)(getProcAddr("glMultiTexCoord2iv"))
	if gpMultiTexCoord2iv == nil {
		return errors.New("glMultiTexCoord2iv")
	}
	gpMultiTexCoord2ivARB = (C.GPMULTITEXCOORD2IVARB)(getProcAddr("glMultiTexCoord2ivARB"))
	gpMultiTexCoord2s = (C.GPMULTITEXCOORD2S)(getProcAddr("glMultiTexCoord2s"))
	if gpMultiTexCoord2s == nil {
		return errors.New("glMultiTexCoord2s")
	}
	gpMultiTexCoord2sARB = (C.GPMULTITEXCOORD2SARB)(getProcAddr("glMultiTexCoord2sARB"))
	gpMultiTexCoord2sv = (C.GPMULTITEXCOORD2SV)(getProcAddr("glMultiTexCoord2sv"))
	if gpMultiTexCoord2sv == nil {
		return errors.New("glMultiTexCoord2sv")
	}
	gpMultiTexCoord2svARB = (C.GPMULTITEXCOORD2SVARB)(getProcAddr("glMultiTexCoord2svARB"))
	gpMultiTexCoord2xOES = (C.GPMULTITEXCOORD2XOES)(getProcAddr("glMultiTexCoord2xOES"))
	gpMultiTexCoord2xvOES = (C.GPMULTITEXCOORD2XVOES)(getProcAddr("glMultiTexCoord2xvOES"))
	gpMultiTexCoord3bOES = (C.GPMULTITEXCOORD3BOES)(getProcAddr("glMultiTexCoord3bOES"))
	gpMultiTexCoord3bvOES = (C.GPMULTITEXCOORD3BVOES)(getProcAddr("glMultiTexCoord3bvOES"))
	gpMultiTexCoord3d = (C.GPMULTITEXCOORD3D)(getProcAddr("glMultiTexCoord3d"))
	if gpMultiTexCoord3d == nil {
		return errors.New("glMultiTexCoord3d")
	}
	gpMultiTexCoord3dARB = (C.GPMULTITEXCOORD3DARB)(getProcAddr("glMultiTexCoord3dARB"))
	gpMultiTexCoord3dv = (C.GPMULTITEXCOORD3DV)(getProcAddr("glMultiTexCoord3dv"))
	if gpMultiTexCoord3dv == nil {
		return errors.New("glMultiTexCoord3dv")
	}
	gpMultiTexCoord3dvARB = (C.GPMULTITEXCOORD3DVARB)(getProcAddr("glMultiTexCoord3dvARB"))
	gpMultiTexCoord3f = (C.GPMULTITEXCOORD3F)(getProcAddr("glMultiTexCoord3f"))
	if gpMultiTexCoord3f == nil {
		return errors.New("glMultiTexCoord3f")
	}
	gpMultiTexCoord3fARB = (C.GPMULTITEXCOORD3FARB)(getProcAddr("glMultiTexCoord3fARB"))
	gpMultiTexCoord3fv = (C.GPMULTITEXCOORD3FV)(getProcAddr("glMultiTexCoord3fv"))
	if gpMultiTexCoord3fv == nil {
		return errors.New("glMultiTexCoord3fv")
	}
	gpMultiTexCoord3fvARB = (C.GPMULTITEXCOORD3FVARB)(getProcAddr("glMultiTexCoord3fvARB"))
	gpMultiTexCoord3hNV = (C.GPMULTITEXCOORD3HNV)(getProcAddr("glMultiTexCoord3hNV"))
	gpMultiTexCoord3hvNV = (C.GPMULTITEXCOORD3HVNV)(getProcAddr("glMultiTexCoord3hvNV"))
	gpMultiTexCoord3i = (C.GPMULTITEXCOORD3I)(getProcAddr("glMultiTexCoord3i"))
	if gpMultiTexCoord3i == nil {
		return errors.New("glMultiTexCoord3i")
	}
	gpMultiTexCoord3iARB = (C.GPMULTITEXCOORD3IARB)(getProcAddr("glMultiTexCoord3iARB"))
	gpMultiTexCoord3iv = (C.GPMULTITEXCOORD3IV)(getProcAddr("glMultiTexCoord3iv"))
	if gpMultiTexCoord3iv == nil {
		return errors.New("glMultiTexCoord3iv")
	}
	gpMultiTexCoord3ivARB = (C.GPMULTITEXCOORD3IVARB)(getProcAddr("glMultiTexCoord3ivARB"))
	gpMultiTexCoord3s = (C.GPMULTITEXCOORD3S)(getProcAddr("glMultiTexCoord3s"))
	if gpMultiTexCoord3s == nil {
		return errors.New("glMultiTexCoord3s")
	}
	gpMultiTexCoord3sARB = (C.GPMULTITEXCOORD3SARB)(getProcAddr("glMultiTexCoord3sARB"))
	gpMultiTexCoord3sv = (C.GPMULTITEXCOORD3SV)(getProcAddr("glMultiTexCoord3sv"))
	if gpMultiTexCoord3sv == nil {
		return errors.New("glMultiTexCoord3sv")
	}
	gpMultiTexCoord3svARB = (C.GPMULTITEXCOORD3SVARB)(getProcAddr("glMultiTexCoord3svARB"))
	gpMultiTexCoord3xOES = (C.GPMULTITEXCOORD3XOES)(getProcAddr("glMultiTexCoord3xOES"))
	gpMultiTexCoord3xvOES = (C.GPMULTITEXCOORD3XVOES)(getProcAddr("glMultiTexCoord3xvOES"))
	gpMultiTexCoord4bOES = (C.GPMULTITEXCOORD4BOES)(getProcAddr("glMultiTexCoord4bOES"))
	gpMultiTexCoord4bvOES = (C.GPMULTITEXCOORD4BVOES)(getProcAddr("glMultiTexCoord4bvOES"))
	gpMultiTexCoord4d = (C.GPMULTITEXCOORD4D)(getProcAddr("glMultiTexCoord4d"))
	if gpMultiTexCoord4d == nil {
		return errors.New("glMultiTexCoord4d")
	}
	gpMultiTexCoord4dARB = (C.GPMULTITEXCOORD4DARB)(getProcAddr("glMultiTexCoord4dARB"))
	gpMultiTexCoord4dv = (C.GPMULTITEXCOORD4DV)(getProcAddr("glMultiTexCoord4dv"))
	if gpMultiTexCoord4dv == nil {
		return errors.New("glMultiTexCoord4dv")
	}
	gpMultiTexCoord4dvARB = (C.GPMULTITEXCOORD4DVARB)(getProcAddr("glMultiTexCoord4dvARB"))
	gpMultiTexCoord4f = (C.GPMULTITEXCOORD4F)(getProcAddr("glMultiTexCoord4f"))
	if gpMultiTexCoord4f == nil {
		return errors.New("glMultiTexCoord4f")
	}
	gpMultiTexCoord4fARB = (C.GPMULTITEXCOORD4FARB)(getProcAddr("glMultiTexCoord4fARB"))
	gpMultiTexCoord4fv = (C.GPMULTITEXCOORD4FV)(getProcAddr("glMultiTexCoord4fv"))
	if gpMultiTexCoord4fv == nil {
		return errors.New("glMultiTexCoord4fv")
	}
	gpMultiTexCoord4fvARB = (C.GPMULTITEXCOORD4FVARB)(getProcAddr("glMultiTexCoord4fvARB"))
	gpMultiTexCoord4hNV = (C.GPMULTITEXCOORD4HNV)(getProcAddr("glMultiTexCoord4hNV"))
	gpMultiTexCoord4hvNV = (C.GPMULTITEXCOORD4HVNV)(getProcAddr("glMultiTexCoord4hvNV"))
	gpMultiTexCoord4i = (C.GPMULTITEXCOORD4I)(getProcAddr("glMultiTexCoord4i"))
	if gpMultiTexCoord4i == nil {
		return errors.New("glMultiTexCoord4i")
	}
	gpMultiTexCoord4iARB = (C.GPMULTITEXCOORD4IARB)(getProcAddr("glMultiTexCoord4iARB"))
	gpMultiTexCoord4iv = (C.GPMULTITEXCOORD4IV)(getProcAddr("glMultiTexCoord4iv"))
	if gpMultiTexCoord4iv == nil {
		return errors.New("glMultiTexCoord4iv")
	}
	gpMultiTexCoord4ivARB = (C.GPMULTITEXCOORD4IVARB)(getProcAddr("glMultiTexCoord4ivARB"))
	gpMultiTexCoord4s = (C.GPMULTITEXCOORD4S)(getProcAddr("glMultiTexCoord4s"))
	if gpMultiTexCoord4s == nil {
		return errors.New("glMultiTexCoord4s")
	}
	gpMultiTexCoord4sARB = (C.GPMULTITEXCOORD4SARB)(getProcAddr("glMultiTexCoord4sARB"))
	gpMultiTexCoord4sv = (C.GPMULTITEXCOORD4SV)(getProcAddr("glMultiTexCoord4sv"))
	if gpMultiTexCoord4sv == nil {
		return errors.New("glMultiTexCoord4sv")
	}
	gpMultiTexCoord4svARB = (C.GPMULTITEXCOORD4SVARB)(getProcAddr("glMultiTexCoord4svARB"))
	gpMultiTexCoord4xOES = (C.GPMULTITEXCOORD4XOES)(getProcAddr("glMultiTexCoord4xOES"))
	gpMultiTexCoord4xvOES = (C.GPMULTITEXCOORD4XVOES)(getProcAddr("glMultiTexCoord4xvOES"))
	gpMultiTexCoordP1ui = (C.GPMULTITEXCOORDP1UI)(getProcAddr("glMultiTexCoordP1ui"))
	if gpMultiTexCoordP1ui == nil {
		return errors.New("glMultiTexCoordP1ui")
	}
	gpMultiTexCoordP1uiv = (C.GPMULTITEXCOORDP1UIV)(getProcAddr("glMultiTexCoordP1uiv"))
	if gpMultiTexCoordP1uiv == nil {
		return errors.New("glMultiTexCoordP1uiv")
	}
	gpMultiTexCoordP2ui = (C.GPMULTITEXCOORDP2UI)(getProcAddr("glMultiTexCoordP2ui"))
	if gpMultiTexCoordP2ui == nil {
		return errors.New("glMultiTexCoordP2ui")
	}
	gpMultiTexCoordP2uiv = (C.GPMULTITEXCOORDP2UIV)(getProcAddr("glMultiTexCoordP2uiv"))
	if gpMultiTexCoordP2uiv == nil {
		return errors.New("glMultiTexCoordP2uiv")
	}
	gpMultiTexCoordP3ui = (C.GPMULTITEXCOORDP3UI)(getProcAddr("glMultiTexCoordP3ui"))
	if gpMultiTexCoordP3ui == nil {
		return errors.New("glMultiTexCoordP3ui")
	}
	gpMultiTexCoordP3uiv = (C.GPMULTITEXCOORDP3UIV)(getProcAddr("glMultiTexCoordP3uiv"))
	if gpMultiTexCoordP3uiv == nil {
		return errors.New("glMultiTexCoordP3uiv")
	}
	gpMultiTexCoordP4ui = (C.GPMULTITEXCOORDP4UI)(getProcAddr("glMultiTexCoordP4ui"))
	if gpMultiTexCoordP4ui == nil {
		return errors.New("glMultiTexCoordP4ui")
	}
	gpMultiTexCoordP4uiv = (C.GPMULTITEXCOORDP4UIV)(getProcAddr("glMultiTexCoordP4uiv"))
	if gpMultiTexCoordP4uiv == nil {
		return errors.New("glMultiTexCoordP4uiv")
	}
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
	gpMulticastBarrierNV = (C.GPMULTICASTBARRIERNV)(getProcAddr("glMulticastBarrierNV"))
	gpMulticastBlitFramebufferNV = (C.GPMULTICASTBLITFRAMEBUFFERNV)(getProcAddr("glMulticastBlitFramebufferNV"))
	gpMulticastBufferSubDataNV = (C.GPMULTICASTBUFFERSUBDATANV)(getProcAddr("glMulticastBufferSubDataNV"))
	gpMulticastCopyBufferSubDataNV = (C.GPMULTICASTCOPYBUFFERSUBDATANV)(getProcAddr("glMulticastCopyBufferSubDataNV"))
	gpMulticastCopyImageSubDataNV = (C.GPMULTICASTCOPYIMAGESUBDATANV)(getProcAddr("glMulticastCopyImageSubDataNV"))
	gpMulticastFramebufferSampleLocationsfvNV = (C.GPMULTICASTFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glMulticastFramebufferSampleLocationsfvNV"))
	gpMulticastGetQueryObjecti64vNV = (C.GPMULTICASTGETQUERYOBJECTI64VNV)(getProcAddr("glMulticastGetQueryObjecti64vNV"))
	gpMulticastGetQueryObjectivNV = (C.GPMULTICASTGETQUERYOBJECTIVNV)(getProcAddr("glMulticastGetQueryObjectivNV"))
	gpMulticastGetQueryObjectui64vNV = (C.GPMULTICASTGETQUERYOBJECTUI64VNV)(getProcAddr("glMulticastGetQueryObjectui64vNV"))
	gpMulticastGetQueryObjectuivNV = (C.GPMULTICASTGETQUERYOBJECTUIVNV)(getProcAddr("glMulticastGetQueryObjectuivNV"))
	gpMulticastWaitSyncNV = (C.GPMULTICASTWAITSYNCNV)(getProcAddr("glMulticastWaitSyncNV"))
	gpNamedBufferData = (C.GPNAMEDBUFFERDATA)(getProcAddr("glNamedBufferData"))
	if gpNamedBufferData == nil {
		return errors.New("glNamedBufferData")
	}
	gpNamedBufferDataEXT = (C.GPNAMEDBUFFERDATAEXT)(getProcAddr("glNamedBufferDataEXT"))
	gpNamedBufferPageCommitmentARB = (C.GPNAMEDBUFFERPAGECOMMITMENTARB)(getProcAddr("glNamedBufferPageCommitmentARB"))
	gpNamedBufferPageCommitmentEXT = (C.GPNAMEDBUFFERPAGECOMMITMENTEXT)(getProcAddr("glNamedBufferPageCommitmentEXT"))
	gpNamedBufferStorage = (C.GPNAMEDBUFFERSTORAGE)(getProcAddr("glNamedBufferStorage"))
	if gpNamedBufferStorage == nil {
		return errors.New("glNamedBufferStorage")
	}
	gpNamedBufferStorageEXT = (C.GPNAMEDBUFFERSTORAGEEXT)(getProcAddr("glNamedBufferStorageEXT"))
	gpNamedBufferStorageExternalEXT = (C.GPNAMEDBUFFERSTORAGEEXTERNALEXT)(getProcAddr("glNamedBufferStorageExternalEXT"))
	gpNamedBufferStorageMemEXT = (C.GPNAMEDBUFFERSTORAGEMEMEXT)(getProcAddr("glNamedBufferStorageMemEXT"))
	gpNamedBufferSubData = (C.GPNAMEDBUFFERSUBDATA)(getProcAddr("glNamedBufferSubData"))
	if gpNamedBufferSubData == nil {
		return errors.New("glNamedBufferSubData")
	}
	gpNamedBufferSubDataEXT = (C.GPNAMEDBUFFERSUBDATAEXT)(getProcAddr("glNamedBufferSubDataEXT"))
	gpNamedCopyBufferSubDataEXT = (C.GPNAMEDCOPYBUFFERSUBDATAEXT)(getProcAddr("glNamedCopyBufferSubDataEXT"))
	gpNamedFramebufferDrawBuffer = (C.GPNAMEDFRAMEBUFFERDRAWBUFFER)(getProcAddr("glNamedFramebufferDrawBuffer"))
	if gpNamedFramebufferDrawBuffer == nil {
		return errors.New("glNamedFramebufferDrawBuffer")
	}
	gpNamedFramebufferDrawBuffers = (C.GPNAMEDFRAMEBUFFERDRAWBUFFERS)(getProcAddr("glNamedFramebufferDrawBuffers"))
	if gpNamedFramebufferDrawBuffers == nil {
		return errors.New("glNamedFramebufferDrawBuffers")
	}
	gpNamedFramebufferParameteri = (C.GPNAMEDFRAMEBUFFERPARAMETERI)(getProcAddr("glNamedFramebufferParameteri"))
	if gpNamedFramebufferParameteri == nil {
		return errors.New("glNamedFramebufferParameteri")
	}
	gpNamedFramebufferParameteriEXT = (C.GPNAMEDFRAMEBUFFERPARAMETERIEXT)(getProcAddr("glNamedFramebufferParameteriEXT"))
	gpNamedFramebufferReadBuffer = (C.GPNAMEDFRAMEBUFFERREADBUFFER)(getProcAddr("glNamedFramebufferReadBuffer"))
	if gpNamedFramebufferReadBuffer == nil {
		return errors.New("glNamedFramebufferReadBuffer")
	}
	gpNamedFramebufferRenderbuffer = (C.GPNAMEDFRAMEBUFFERRENDERBUFFER)(getProcAddr("glNamedFramebufferRenderbuffer"))
	if gpNamedFramebufferRenderbuffer == nil {
		return errors.New("glNamedFramebufferRenderbuffer")
	}
	gpNamedFramebufferRenderbufferEXT = (C.GPNAMEDFRAMEBUFFERRENDERBUFFEREXT)(getProcAddr("glNamedFramebufferRenderbufferEXT"))
	gpNamedFramebufferSampleLocationsfvARB = (C.GPNAMEDFRAMEBUFFERSAMPLELOCATIONSFVARB)(getProcAddr("glNamedFramebufferSampleLocationsfvARB"))
	gpNamedFramebufferSampleLocationsfvNV = (C.GPNAMEDFRAMEBUFFERSAMPLELOCATIONSFVNV)(getProcAddr("glNamedFramebufferSampleLocationsfvNV"))
	gpNamedFramebufferSamplePositionsfvAMD = (C.GPNAMEDFRAMEBUFFERSAMPLEPOSITIONSFVAMD)(getProcAddr("glNamedFramebufferSamplePositionsfvAMD"))
	gpNamedFramebufferTexture = (C.GPNAMEDFRAMEBUFFERTEXTURE)(getProcAddr("glNamedFramebufferTexture"))
	if gpNamedFramebufferTexture == nil {
		return errors.New("glNamedFramebufferTexture")
	}
	gpNamedFramebufferTexture1DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE1DEXT)(getProcAddr("glNamedFramebufferTexture1DEXT"))
	gpNamedFramebufferTexture2DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE2DEXT)(getProcAddr("glNamedFramebufferTexture2DEXT"))
	gpNamedFramebufferTexture3DEXT = (C.GPNAMEDFRAMEBUFFERTEXTURE3DEXT)(getProcAddr("glNamedFramebufferTexture3DEXT"))
	gpNamedFramebufferTextureEXT = (C.GPNAMEDFRAMEBUFFERTEXTUREEXT)(getProcAddr("glNamedFramebufferTextureEXT"))
	gpNamedFramebufferTextureFaceEXT = (C.GPNAMEDFRAMEBUFFERTEXTUREFACEEXT)(getProcAddr("glNamedFramebufferTextureFaceEXT"))
	gpNamedFramebufferTextureLayer = (C.GPNAMEDFRAMEBUFFERTEXTURELAYER)(getProcAddr("glNamedFramebufferTextureLayer"))
	if gpNamedFramebufferTextureLayer == nil {
		return errors.New("glNamedFramebufferTextureLayer")
	}
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
	if gpNamedRenderbufferStorage == nil {
		return errors.New("glNamedRenderbufferStorage")
	}
	gpNamedRenderbufferStorageEXT = (C.GPNAMEDRENDERBUFFERSTORAGEEXT)(getProcAddr("glNamedRenderbufferStorageEXT"))
	gpNamedRenderbufferStorageMultisample = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLE)(getProcAddr("glNamedRenderbufferStorageMultisample"))
	if gpNamedRenderbufferStorageMultisample == nil {
		return errors.New("glNamedRenderbufferStorageMultisample")
	}
	gpNamedRenderbufferStorageMultisampleCoverageEXT = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLECOVERAGEEXT)(getProcAddr("glNamedRenderbufferStorageMultisampleCoverageEXT"))
	gpNamedRenderbufferStorageMultisampleEXT = (C.GPNAMEDRENDERBUFFERSTORAGEMULTISAMPLEEXT)(getProcAddr("glNamedRenderbufferStorageMultisampleEXT"))
	gpNamedStringARB = (C.GPNAMEDSTRINGARB)(getProcAddr("glNamedStringARB"))
	gpNewList = (C.GPNEWLIST)(getProcAddr("glNewList"))
	if gpNewList == nil {
		return errors.New("glNewList")
	}
	gpNewObjectBufferATI = (C.GPNEWOBJECTBUFFERATI)(getProcAddr("glNewObjectBufferATI"))
	gpNormal3b = (C.GPNORMAL3B)(getProcAddr("glNormal3b"))
	if gpNormal3b == nil {
		return errors.New("glNormal3b")
	}
	gpNormal3bv = (C.GPNORMAL3BV)(getProcAddr("glNormal3bv"))
	if gpNormal3bv == nil {
		return errors.New("glNormal3bv")
	}
	gpNormal3d = (C.GPNORMAL3D)(getProcAddr("glNormal3d"))
	if gpNormal3d == nil {
		return errors.New("glNormal3d")
	}
	gpNormal3dv = (C.GPNORMAL3DV)(getProcAddr("glNormal3dv"))
	if gpNormal3dv == nil {
		return errors.New("glNormal3dv")
	}
	gpNormal3f = (C.GPNORMAL3F)(getProcAddr("glNormal3f"))
	if gpNormal3f == nil {
		return errors.New("glNormal3f")
	}
	gpNormal3fVertex3fSUN = (C.GPNORMAL3FVERTEX3FSUN)(getProcAddr("glNormal3fVertex3fSUN"))
	gpNormal3fVertex3fvSUN = (C.GPNORMAL3FVERTEX3FVSUN)(getProcAddr("glNormal3fVertex3fvSUN"))
	gpNormal3fv = (C.GPNORMAL3FV)(getProcAddr("glNormal3fv"))
	if gpNormal3fv == nil {
		return errors.New("glNormal3fv")
	}
	gpNormal3hNV = (C.GPNORMAL3HNV)(getProcAddr("glNormal3hNV"))
	gpNormal3hvNV = (C.GPNORMAL3HVNV)(getProcAddr("glNormal3hvNV"))
	gpNormal3i = (C.GPNORMAL3I)(getProcAddr("glNormal3i"))
	if gpNormal3i == nil {
		return errors.New("glNormal3i")
	}
	gpNormal3iv = (C.GPNORMAL3IV)(getProcAddr("glNormal3iv"))
	if gpNormal3iv == nil {
		return errors.New("glNormal3iv")
	}
	gpNormal3s = (C.GPNORMAL3S)(getProcAddr("glNormal3s"))
	if gpNormal3s == nil {
		return errors.New("glNormal3s")
	}
	gpNormal3sv = (C.GPNORMAL3SV)(getProcAddr("glNormal3sv"))
	if gpNormal3sv == nil {
		return errors.New("glNormal3sv")
	}
	gpNormal3xOES = (C.GPNORMAL3XOES)(getProcAddr("glNormal3xOES"))
	gpNormal3xvOES = (C.GPNORMAL3XVOES)(getProcAddr("glNormal3xvOES"))
	gpNormalFormatNV = (C.GPNORMALFORMATNV)(getProcAddr("glNormalFormatNV"))
	gpNormalP3ui = (C.GPNORMALP3UI)(getProcAddr("glNormalP3ui"))
	if gpNormalP3ui == nil {
		return errors.New("glNormalP3ui")
	}
	gpNormalP3uiv = (C.GPNORMALP3UIV)(getProcAddr("glNormalP3uiv"))
	if gpNormalP3uiv == nil {
		return errors.New("glNormalP3uiv")
	}
	gpNormalPointer = (C.GPNORMALPOINTER)(getProcAddr("glNormalPointer"))
	if gpNormalPointer == nil {
		return errors.New("glNormalPointer")
	}
	gpNormalPointerEXT = (C.GPNORMALPOINTEREXT)(getProcAddr("glNormalPointerEXT"))
	gpNormalPointerListIBM = (C.GPNORMALPOINTERLISTIBM)(getProcAddr("glNormalPointerListIBM"))
	gpNormalPointervINTEL = (C.GPNORMALPOINTERVINTEL)(getProcAddr("glNormalPointervINTEL"))
	gpNormalStream3bATI = (C.GPNORMALSTREAM3BATI)(getProcAddr("glNormalStream3bATI"))
	gpNormalStream3bvATI = (C.GPNORMALSTREAM3BVATI)(getProcAddr("glNormalStream3bvATI"))
	gpNormalStream3dATI = (C.GPNORMALSTREAM3DATI)(getProcAddr("glNormalStream3dATI"))
	gpNormalStream3dvATI = (C.GPNORMALSTREAM3DVATI)(getProcAddr("glNormalStream3dvATI"))
	gpNormalStream3fATI = (C.GPNORMALSTREAM3FATI)(getProcAddr("glNormalStream3fATI"))
	gpNormalStream3fvATI = (C.GPNORMALSTREAM3FVATI)(getProcAddr("glNormalStream3fvATI"))
	gpNormalStream3iATI = (C.GPNORMALSTREAM3IATI)(getProcAddr("glNormalStream3iATI"))
	gpNormalStream3ivATI = (C.GPNORMALSTREAM3IVATI)(getProcAddr("glNormalStream3ivATI"))
	gpNormalStream3sATI = (C.GPNORMALSTREAM3SATI)(getProcAddr("glNormalStream3sATI"))
	gpNormalStream3svATI = (C.GPNORMALSTREAM3SVATI)(getProcAddr("glNormalStream3svATI"))
	gpObjectLabel = (C.GPOBJECTLABEL)(getProcAddr("glObjectLabel"))
	if gpObjectLabel == nil {
		return errors.New("glObjectLabel")
	}
	gpObjectLabelKHR = (C.GPOBJECTLABELKHR)(getProcAddr("glObjectLabelKHR"))
	gpObjectPtrLabel = (C.GPOBJECTPTRLABEL)(getProcAddr("glObjectPtrLabel"))
	if gpObjectPtrLabel == nil {
		return errors.New("glObjectPtrLabel")
	}
	gpObjectPtrLabelKHR = (C.GPOBJECTPTRLABELKHR)(getProcAddr("glObjectPtrLabelKHR"))
	gpObjectPurgeableAPPLE = (C.GPOBJECTPURGEABLEAPPLE)(getProcAddr("glObjectPurgeableAPPLE"))
	gpObjectUnpurgeableAPPLE = (C.GPOBJECTUNPURGEABLEAPPLE)(getProcAddr("glObjectUnpurgeableAPPLE"))
	gpOrtho = (C.GPORTHO)(getProcAddr("glOrtho"))
	if gpOrtho == nil {
		return errors.New("glOrtho")
	}
	gpOrthofOES = (C.GPORTHOFOES)(getProcAddr("glOrthofOES"))
	gpOrthoxOES = (C.GPORTHOXOES)(getProcAddr("glOrthoxOES"))
	gpPNTrianglesfATI = (C.GPPNTRIANGLESFATI)(getProcAddr("glPNTrianglesfATI"))
	gpPNTrianglesiATI = (C.GPPNTRIANGLESIATI)(getProcAddr("glPNTrianglesiATI"))
	gpPassTexCoordATI = (C.GPPASSTEXCOORDATI)(getProcAddr("glPassTexCoordATI"))
	gpPassThrough = (C.GPPASSTHROUGH)(getProcAddr("glPassThrough"))
	if gpPassThrough == nil {
		return errors.New("glPassThrough")
	}
	gpPassThroughxOES = (C.GPPASSTHROUGHXOES)(getProcAddr("glPassThroughxOES"))
	gpPatchParameterfv = (C.GPPATCHPARAMETERFV)(getProcAddr("glPatchParameterfv"))
	if gpPatchParameterfv == nil {
		return errors.New("glPatchParameterfv")
	}
	gpPatchParameteri = (C.GPPATCHPARAMETERI)(getProcAddr("glPatchParameteri"))
	if gpPatchParameteri == nil {
		return errors.New("glPatchParameteri")
	}
	gpPathColorGenNV = (C.GPPATHCOLORGENNV)(getProcAddr("glPathColorGenNV"))
	gpPathCommandsNV = (C.GPPATHCOMMANDSNV)(getProcAddr("glPathCommandsNV"))
	gpPathCoordsNV = (C.GPPATHCOORDSNV)(getProcAddr("glPathCoordsNV"))
	gpPathCoverDepthFuncNV = (C.GPPATHCOVERDEPTHFUNCNV)(getProcAddr("glPathCoverDepthFuncNV"))
	gpPathDashArrayNV = (C.GPPATHDASHARRAYNV)(getProcAddr("glPathDashArrayNV"))
	gpPathFogGenNV = (C.GPPATHFOGGENNV)(getProcAddr("glPathFogGenNV"))
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
	gpPathTexGenNV = (C.GPPATHTEXGENNV)(getProcAddr("glPathTexGenNV"))
	gpPauseTransformFeedback = (C.GPPAUSETRANSFORMFEEDBACK)(getProcAddr("glPauseTransformFeedback"))
	if gpPauseTransformFeedback == nil {
		return errors.New("glPauseTransformFeedback")
	}
	gpPauseTransformFeedbackNV = (C.GPPAUSETRANSFORMFEEDBACKNV)(getProcAddr("glPauseTransformFeedbackNV"))
	gpPixelDataRangeNV = (C.GPPIXELDATARANGENV)(getProcAddr("glPixelDataRangeNV"))
	gpPixelMapfv = (C.GPPIXELMAPFV)(getProcAddr("glPixelMapfv"))
	if gpPixelMapfv == nil {
		return errors.New("glPixelMapfv")
	}
	gpPixelMapuiv = (C.GPPIXELMAPUIV)(getProcAddr("glPixelMapuiv"))
	if gpPixelMapuiv == nil {
		return errors.New("glPixelMapuiv")
	}
	gpPixelMapusv = (C.GPPIXELMAPUSV)(getProcAddr("glPixelMapusv"))
	if gpPixelMapusv == nil {
		return errors.New("glPixelMapusv")
	}
	gpPixelMapx = (C.GPPIXELMAPX)(getProcAddr("glPixelMapx"))
	gpPixelStoref = (C.GPPIXELSTOREF)(getProcAddr("glPixelStoref"))
	if gpPixelStoref == nil {
		return errors.New("glPixelStoref")
	}
	gpPixelStorei = (C.GPPIXELSTOREI)(getProcAddr("glPixelStorei"))
	if gpPixelStorei == nil {
		return errors.New("glPixelStorei")
	}
	gpPixelStorex = (C.GPPIXELSTOREX)(getProcAddr("glPixelStorex"))
	gpPixelTexGenParameterfSGIS = (C.GPPIXELTEXGENPARAMETERFSGIS)(getProcAddr("glPixelTexGenParameterfSGIS"))
	gpPixelTexGenParameterfvSGIS = (C.GPPIXELTEXGENPARAMETERFVSGIS)(getProcAddr("glPixelTexGenParameterfvSGIS"))
	gpPixelTexGenParameteriSGIS = (C.GPPIXELTEXGENPARAMETERISGIS)(getProcAddr("glPixelTexGenParameteriSGIS"))
	gpPixelTexGenParameterivSGIS = (C.GPPIXELTEXGENPARAMETERIVSGIS)(getProcAddr("glPixelTexGenParameterivSGIS"))
	gpPixelTexGenSGIX = (C.GPPIXELTEXGENSGIX)(getProcAddr("glPixelTexGenSGIX"))
	gpPixelTransferf = (C.GPPIXELTRANSFERF)(getProcAddr("glPixelTransferf"))
	if gpPixelTransferf == nil {
		return errors.New("glPixelTransferf")
	}
	gpPixelTransferi = (C.GPPIXELTRANSFERI)(getProcAddr("glPixelTransferi"))
	if gpPixelTransferi == nil {
		return errors.New("glPixelTransferi")
	}
	gpPixelTransferxOES = (C.GPPIXELTRANSFERXOES)(getProcAddr("glPixelTransferxOES"))
	gpPixelTransformParameterfEXT = (C.GPPIXELTRANSFORMPARAMETERFEXT)(getProcAddr("glPixelTransformParameterfEXT"))
	gpPixelTransformParameterfvEXT = (C.GPPIXELTRANSFORMPARAMETERFVEXT)(getProcAddr("glPixelTransformParameterfvEXT"))
	gpPixelTransformParameteriEXT = (C.GPPIXELTRANSFORMPARAMETERIEXT)(getProcAddr("glPixelTransformParameteriEXT"))
	gpPixelTransformParameterivEXT = (C.GPPIXELTRANSFORMPARAMETERIVEXT)(getProcAddr("glPixelTransformParameterivEXT"))
	gpPixelZoom = (C.GPPIXELZOOM)(getProcAddr("glPixelZoom"))
	if gpPixelZoom == nil {
		return errors.New("glPixelZoom")
	}
	gpPixelZoomxOES = (C.GPPIXELZOOMXOES)(getProcAddr("glPixelZoomxOES"))
	gpPointAlongPathNV = (C.GPPOINTALONGPATHNV)(getProcAddr("glPointAlongPathNV"))
	gpPointParameterf = (C.GPPOINTPARAMETERF)(getProcAddr("glPointParameterf"))
	if gpPointParameterf == nil {
		return errors.New("glPointParameterf")
	}
	gpPointParameterfARB = (C.GPPOINTPARAMETERFARB)(getProcAddr("glPointParameterfARB"))
	gpPointParameterfEXT = (C.GPPOINTPARAMETERFEXT)(getProcAddr("glPointParameterfEXT"))
	gpPointParameterfSGIS = (C.GPPOINTPARAMETERFSGIS)(getProcAddr("glPointParameterfSGIS"))
	gpPointParameterfv = (C.GPPOINTPARAMETERFV)(getProcAddr("glPointParameterfv"))
	if gpPointParameterfv == nil {
		return errors.New("glPointParameterfv")
	}
	gpPointParameterfvARB = (C.GPPOINTPARAMETERFVARB)(getProcAddr("glPointParameterfvARB"))
	gpPointParameterfvEXT = (C.GPPOINTPARAMETERFVEXT)(getProcAddr("glPointParameterfvEXT"))
	gpPointParameterfvSGIS = (C.GPPOINTPARAMETERFVSGIS)(getProcAddr("glPointParameterfvSGIS"))
	gpPointParameteri = (C.GPPOINTPARAMETERI)(getProcAddr("glPointParameteri"))
	if gpPointParameteri == nil {
		return errors.New("glPointParameteri")
	}
	gpPointParameteriNV = (C.GPPOINTPARAMETERINV)(getProcAddr("glPointParameteriNV"))
	gpPointParameteriv = (C.GPPOINTPARAMETERIV)(getProcAddr("glPointParameteriv"))
	if gpPointParameteriv == nil {
		return errors.New("glPointParameteriv")
	}
	gpPointParameterivNV = (C.GPPOINTPARAMETERIVNV)(getProcAddr("glPointParameterivNV"))
	gpPointParameterxOES = (C.GPPOINTPARAMETERXOES)(getProcAddr("glPointParameterxOES"))
	gpPointParameterxvOES = (C.GPPOINTPARAMETERXVOES)(getProcAddr("glPointParameterxvOES"))
	gpPointSize = (C.GPPOINTSIZE)(getProcAddr("glPointSize"))
	if gpPointSize == nil {
		return errors.New("glPointSize")
	}
	gpPointSizexOES = (C.GPPOINTSIZEXOES)(getProcAddr("glPointSizexOES"))
	gpPollAsyncSGIX = (C.GPPOLLASYNCSGIX)(getProcAddr("glPollAsyncSGIX"))
	gpPollInstrumentsSGIX = (C.GPPOLLINSTRUMENTSSGIX)(getProcAddr("glPollInstrumentsSGIX"))
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
	gpPolygonOffsetEXT = (C.GPPOLYGONOFFSETEXT)(getProcAddr("glPolygonOffsetEXT"))
	gpPolygonOffsetxOES = (C.GPPOLYGONOFFSETXOES)(getProcAddr("glPolygonOffsetxOES"))
	gpPolygonStipple = (C.GPPOLYGONSTIPPLE)(getProcAddr("glPolygonStipple"))
	if gpPolygonStipple == nil {
		return errors.New("glPolygonStipple")
	}
	gpPopAttrib = (C.GPPOPATTRIB)(getProcAddr("glPopAttrib"))
	if gpPopAttrib == nil {
		return errors.New("glPopAttrib")
	}
	gpPopClientAttrib = (C.GPPOPCLIENTATTRIB)(getProcAddr("glPopClientAttrib"))
	if gpPopClientAttrib == nil {
		return errors.New("glPopClientAttrib")
	}
	gpPopDebugGroup = (C.GPPOPDEBUGGROUP)(getProcAddr("glPopDebugGroup"))
	if gpPopDebugGroup == nil {
		return errors.New("glPopDebugGroup")
	}
	gpPopDebugGroupKHR = (C.GPPOPDEBUGGROUPKHR)(getProcAddr("glPopDebugGroupKHR"))
	gpPopGroupMarkerEXT = (C.GPPOPGROUPMARKEREXT)(getProcAddr("glPopGroupMarkerEXT"))
	gpPopMatrix = (C.GPPOPMATRIX)(getProcAddr("glPopMatrix"))
	if gpPopMatrix == nil {
		return errors.New("glPopMatrix")
	}
	gpPopName = (C.GPPOPNAME)(getProcAddr("glPopName"))
	if gpPopName == nil {
		return errors.New("glPopName")
	}
	gpPresentFrameDualFillNV = (C.GPPRESENTFRAMEDUALFILLNV)(getProcAddr("glPresentFrameDualFillNV"))
	gpPresentFrameKeyedNV = (C.GPPRESENTFRAMEKEYEDNV)(getProcAddr("glPresentFrameKeyedNV"))
	gpPrimitiveBoundingBoxARB = (C.GPPRIMITIVEBOUNDINGBOXARB)(getProcAddr("glPrimitiveBoundingBoxARB"))
	gpPrimitiveRestartIndex = (C.GPPRIMITIVERESTARTINDEX)(getProcAddr("glPrimitiveRestartIndex"))
	if gpPrimitiveRestartIndex == nil {
		return errors.New("glPrimitiveRestartIndex")
	}
	gpPrimitiveRestartIndexNV = (C.GPPRIMITIVERESTARTINDEXNV)(getProcAddr("glPrimitiveRestartIndexNV"))
	gpPrimitiveRestartNV = (C.GPPRIMITIVERESTARTNV)(getProcAddr("glPrimitiveRestartNV"))
	gpPrioritizeTextures = (C.GPPRIORITIZETEXTURES)(getProcAddr("glPrioritizeTextures"))
	if gpPrioritizeTextures == nil {
		return errors.New("glPrioritizeTextures")
	}
	gpPrioritizeTexturesEXT = (C.GPPRIORITIZETEXTURESEXT)(getProcAddr("glPrioritizeTexturesEXT"))
	gpPrioritizeTexturesxOES = (C.GPPRIORITIZETEXTURESXOES)(getProcAddr("glPrioritizeTexturesxOES"))
	gpProgramBinary = (C.GPPROGRAMBINARY)(getProcAddr("glProgramBinary"))
	if gpProgramBinary == nil {
		return errors.New("glProgramBinary")
	}
	gpProgramBufferParametersIivNV = (C.GPPROGRAMBUFFERPARAMETERSIIVNV)(getProcAddr("glProgramBufferParametersIivNV"))
	gpProgramBufferParametersIuivNV = (C.GPPROGRAMBUFFERPARAMETERSIUIVNV)(getProcAddr("glProgramBufferParametersIuivNV"))
	gpProgramBufferParametersfvNV = (C.GPPROGRAMBUFFERPARAMETERSFVNV)(getProcAddr("glProgramBufferParametersfvNV"))
	gpProgramEnvParameter4dARB = (C.GPPROGRAMENVPARAMETER4DARB)(getProcAddr("glProgramEnvParameter4dARB"))
	gpProgramEnvParameter4dvARB = (C.GPPROGRAMENVPARAMETER4DVARB)(getProcAddr("glProgramEnvParameter4dvARB"))
	gpProgramEnvParameter4fARB = (C.GPPROGRAMENVPARAMETER4FARB)(getProcAddr("glProgramEnvParameter4fARB"))
	gpProgramEnvParameter4fvARB = (C.GPPROGRAMENVPARAMETER4FVARB)(getProcAddr("glProgramEnvParameter4fvARB"))
	gpProgramEnvParameterI4iNV = (C.GPPROGRAMENVPARAMETERI4INV)(getProcAddr("glProgramEnvParameterI4iNV"))
	gpProgramEnvParameterI4ivNV = (C.GPPROGRAMENVPARAMETERI4IVNV)(getProcAddr("glProgramEnvParameterI4ivNV"))
	gpProgramEnvParameterI4uiNV = (C.GPPROGRAMENVPARAMETERI4UINV)(getProcAddr("glProgramEnvParameterI4uiNV"))
	gpProgramEnvParameterI4uivNV = (C.GPPROGRAMENVPARAMETERI4UIVNV)(getProcAddr("glProgramEnvParameterI4uivNV"))
	gpProgramEnvParameters4fvEXT = (C.GPPROGRAMENVPARAMETERS4FVEXT)(getProcAddr("glProgramEnvParameters4fvEXT"))
	gpProgramEnvParametersI4ivNV = (C.GPPROGRAMENVPARAMETERSI4IVNV)(getProcAddr("glProgramEnvParametersI4ivNV"))
	gpProgramEnvParametersI4uivNV = (C.GPPROGRAMENVPARAMETERSI4UIVNV)(getProcAddr("glProgramEnvParametersI4uivNV"))
	gpProgramLocalParameter4dARB = (C.GPPROGRAMLOCALPARAMETER4DARB)(getProcAddr("glProgramLocalParameter4dARB"))
	gpProgramLocalParameter4dvARB = (C.GPPROGRAMLOCALPARAMETER4DVARB)(getProcAddr("glProgramLocalParameter4dvARB"))
	gpProgramLocalParameter4fARB = (C.GPPROGRAMLOCALPARAMETER4FARB)(getProcAddr("glProgramLocalParameter4fARB"))
	gpProgramLocalParameter4fvARB = (C.GPPROGRAMLOCALPARAMETER4FVARB)(getProcAddr("glProgramLocalParameter4fvARB"))
	gpProgramLocalParameterI4iNV = (C.GPPROGRAMLOCALPARAMETERI4INV)(getProcAddr("glProgramLocalParameterI4iNV"))
	gpProgramLocalParameterI4ivNV = (C.GPPROGRAMLOCALPARAMETERI4IVNV)(getProcAddr("glProgramLocalParameterI4ivNV"))
	gpProgramLocalParameterI4uiNV = (C.GPPROGRAMLOCALPARAMETERI4UINV)(getProcAddr("glProgramLocalParameterI4uiNV"))
	gpProgramLocalParameterI4uivNV = (C.GPPROGRAMLOCALPARAMETERI4UIVNV)(getProcAddr("glProgramLocalParameterI4uivNV"))
	gpProgramLocalParameters4fvEXT = (C.GPPROGRAMLOCALPARAMETERS4FVEXT)(getProcAddr("glProgramLocalParameters4fvEXT"))
	gpProgramLocalParametersI4ivNV = (C.GPPROGRAMLOCALPARAMETERSI4IVNV)(getProcAddr("glProgramLocalParametersI4ivNV"))
	gpProgramLocalParametersI4uivNV = (C.GPPROGRAMLOCALPARAMETERSI4UIVNV)(getProcAddr("glProgramLocalParametersI4uivNV"))
	gpProgramNamedParameter4dNV = (C.GPPROGRAMNAMEDPARAMETER4DNV)(getProcAddr("glProgramNamedParameter4dNV"))
	gpProgramNamedParameter4dvNV = (C.GPPROGRAMNAMEDPARAMETER4DVNV)(getProcAddr("glProgramNamedParameter4dvNV"))
	gpProgramNamedParameter4fNV = (C.GPPROGRAMNAMEDPARAMETER4FNV)(getProcAddr("glProgramNamedParameter4fNV"))
	gpProgramNamedParameter4fvNV = (C.GPPROGRAMNAMEDPARAMETER4FVNV)(getProcAddr("glProgramNamedParameter4fvNV"))
	gpProgramParameter4dNV = (C.GPPROGRAMPARAMETER4DNV)(getProcAddr("glProgramParameter4dNV"))
	gpProgramParameter4dvNV = (C.GPPROGRAMPARAMETER4DVNV)(getProcAddr("glProgramParameter4dvNV"))
	gpProgramParameter4fNV = (C.GPPROGRAMPARAMETER4FNV)(getProcAddr("glProgramParameter4fNV"))
	gpProgramParameter4fvNV = (C.GPPROGRAMPARAMETER4FVNV)(getProcAddr("glProgramParameter4fvNV"))
	gpProgramParameteri = (C.GPPROGRAMPARAMETERI)(getProcAddr("glProgramParameteri"))
	if gpProgramParameteri == nil {
		return errors.New("glProgramParameteri")
	}
	gpProgramParameteriARB = (C.GPPROGRAMPARAMETERIARB)(getProcAddr("glProgramParameteriARB"))
	gpProgramParameteriEXT = (C.GPPROGRAMPARAMETERIEXT)(getProcAddr("glProgramParameteriEXT"))
	gpProgramParameters4dvNV = (C.GPPROGRAMPARAMETERS4DVNV)(getProcAddr("glProgramParameters4dvNV"))
	gpProgramParameters4fvNV = (C.GPPROGRAMPARAMETERS4FVNV)(getProcAddr("glProgramParameters4fvNV"))
	gpProgramPathFragmentInputGenNV = (C.GPPROGRAMPATHFRAGMENTINPUTGENNV)(getProcAddr("glProgramPathFragmentInputGenNV"))
	gpProgramStringARB = (C.GPPROGRAMSTRINGARB)(getProcAddr("glProgramStringARB"))
	gpProgramSubroutineParametersuivNV = (C.GPPROGRAMSUBROUTINEPARAMETERSUIVNV)(getProcAddr("glProgramSubroutineParametersuivNV"))
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
	gpProgramVertexLimitNV = (C.GPPROGRAMVERTEXLIMITNV)(getProcAddr("glProgramVertexLimitNV"))
	gpProvokingVertex = (C.GPPROVOKINGVERTEX)(getProcAddr("glProvokingVertex"))
	if gpProvokingVertex == nil {
		return errors.New("glProvokingVertex")
	}
	gpProvokingVertexEXT = (C.GPPROVOKINGVERTEXEXT)(getProcAddr("glProvokingVertexEXT"))
	gpPushAttrib = (C.GPPUSHATTRIB)(getProcAddr("glPushAttrib"))
	if gpPushAttrib == nil {
		return errors.New("glPushAttrib")
	}
	gpPushClientAttrib = (C.GPPUSHCLIENTATTRIB)(getProcAddr("glPushClientAttrib"))
	if gpPushClientAttrib == nil {
		return errors.New("glPushClientAttrib")
	}
	gpPushClientAttribDefaultEXT = (C.GPPUSHCLIENTATTRIBDEFAULTEXT)(getProcAddr("glPushClientAttribDefaultEXT"))
	gpPushDebugGroup = (C.GPPUSHDEBUGGROUP)(getProcAddr("glPushDebugGroup"))
	if gpPushDebugGroup == nil {
		return errors.New("glPushDebugGroup")
	}
	gpPushDebugGroupKHR = (C.GPPUSHDEBUGGROUPKHR)(getProcAddr("glPushDebugGroupKHR"))
	gpPushGroupMarkerEXT = (C.GPPUSHGROUPMARKEREXT)(getProcAddr("glPushGroupMarkerEXT"))
	gpPushMatrix = (C.GPPUSHMATRIX)(getProcAddr("glPushMatrix"))
	if gpPushMatrix == nil {
		return errors.New("glPushMatrix")
	}
	gpPushName = (C.GPPUSHNAME)(getProcAddr("glPushName"))
	if gpPushName == nil {
		return errors.New("glPushName")
	}
	gpQueryCounter = (C.GPQUERYCOUNTER)(getProcAddr("glQueryCounter"))
	if gpQueryCounter == nil {
		return errors.New("glQueryCounter")
	}
	gpQueryMatrixxOES = (C.GPQUERYMATRIXXOES)(getProcAddr("glQueryMatrixxOES"))
	gpQueryObjectParameteruiAMD = (C.GPQUERYOBJECTPARAMETERUIAMD)(getProcAddr("glQueryObjectParameteruiAMD"))
	gpQueryResourceNV = (C.GPQUERYRESOURCENV)(getProcAddr("glQueryResourceNV"))
	gpQueryResourceTagNV = (C.GPQUERYRESOURCETAGNV)(getProcAddr("glQueryResourceTagNV"))
	gpRasterPos2d = (C.GPRASTERPOS2D)(getProcAddr("glRasterPos2d"))
	if gpRasterPos2d == nil {
		return errors.New("glRasterPos2d")
	}
	gpRasterPos2dv = (C.GPRASTERPOS2DV)(getProcAddr("glRasterPos2dv"))
	if gpRasterPos2dv == nil {
		return errors.New("glRasterPos2dv")
	}
	gpRasterPos2f = (C.GPRASTERPOS2F)(getProcAddr("glRasterPos2f"))
	if gpRasterPos2f == nil {
		return errors.New("glRasterPos2f")
	}
	gpRasterPos2fv = (C.GPRASTERPOS2FV)(getProcAddr("glRasterPos2fv"))
	if gpRasterPos2fv == nil {
		return errors.New("glRasterPos2fv")
	}
	gpRasterPos2i = (C.GPRASTERPOS2I)(getProcAddr("glRasterPos2i"))
	if gpRasterPos2i == nil {
		return errors.New("glRasterPos2i")
	}
	gpRasterPos2iv = (C.GPRASTERPOS2IV)(getProcAddr("glRasterPos2iv"))
	if gpRasterPos2iv == nil {
		return errors.New("glRasterPos2iv")
	}
	gpRasterPos2s = (C.GPRASTERPOS2S)(getProcAddr("glRasterPos2s"))
	if gpRasterPos2s == nil {
		return errors.New("glRasterPos2s")
	}
	gpRasterPos2sv = (C.GPRASTERPOS2SV)(getProcAddr("glRasterPos2sv"))
	if gpRasterPos2sv == nil {
		return errors.New("glRasterPos2sv")
	}
	gpRasterPos2xOES = (C.GPRASTERPOS2XOES)(getProcAddr("glRasterPos2xOES"))
	gpRasterPos2xvOES = (C.GPRASTERPOS2XVOES)(getProcAddr("glRasterPos2xvOES"))
	gpRasterPos3d = (C.GPRASTERPOS3D)(getProcAddr("glRasterPos3d"))
	if gpRasterPos3d == nil {
		return errors.New("glRasterPos3d")
	}
	gpRasterPos3dv = (C.GPRASTERPOS3DV)(getProcAddr("glRasterPos3dv"))
	if gpRasterPos3dv == nil {
		return errors.New("glRasterPos3dv")
	}
	gpRasterPos3f = (C.GPRASTERPOS3F)(getProcAddr("glRasterPos3f"))
	if gpRasterPos3f == nil {
		return errors.New("glRasterPos3f")
	}
	gpRasterPos3fv = (C.GPRASTERPOS3FV)(getProcAddr("glRasterPos3fv"))
	if gpRasterPos3fv == nil {
		return errors.New("glRasterPos3fv")
	}
	gpRasterPos3i = (C.GPRASTERPOS3I)(getProcAddr("glRasterPos3i"))
	if gpRasterPos3i == nil {
		return errors.New("glRasterPos3i")
	}
	gpRasterPos3iv = (C.GPRASTERPOS3IV)(getProcAddr("glRasterPos3iv"))
	if gpRasterPos3iv == nil {
		return errors.New("glRasterPos3iv")
	}
	gpRasterPos3s = (C.GPRASTERPOS3S)(getProcAddr("glRasterPos3s"))
	if gpRasterPos3s == nil {
		return errors.New("glRasterPos3s")
	}
	gpRasterPos3sv = (C.GPRASTERPOS3SV)(getProcAddr("glRasterPos3sv"))
	if gpRasterPos3sv == nil {
		return errors.New("glRasterPos3sv")
	}
	gpRasterPos3xOES = (C.GPRASTERPOS3XOES)(getProcAddr("glRasterPos3xOES"))
	gpRasterPos3xvOES = (C.GPRASTERPOS3XVOES)(getProcAddr("glRasterPos3xvOES"))
	gpRasterPos4d = (C.GPRASTERPOS4D)(getProcAddr("glRasterPos4d"))
	if gpRasterPos4d == nil {
		return errors.New("glRasterPos4d")
	}
	gpRasterPos4dv = (C.GPRASTERPOS4DV)(getProcAddr("glRasterPos4dv"))
	if gpRasterPos4dv == nil {
		return errors.New("glRasterPos4dv")
	}
	gpRasterPos4f = (C.GPRASTERPOS4F)(getProcAddr("glRasterPos4f"))
	if gpRasterPos4f == nil {
		return errors.New("glRasterPos4f")
	}
	gpRasterPos4fv = (C.GPRASTERPOS4FV)(getProcAddr("glRasterPos4fv"))
	if gpRasterPos4fv == nil {
		return errors.New("glRasterPos4fv")
	}
	gpRasterPos4i = (C.GPRASTERPOS4I)(getProcAddr("glRasterPos4i"))
	if gpRasterPos4i == nil {
		return errors.New("glRasterPos4i")
	}
	gpRasterPos4iv = (C.GPRASTERPOS4IV)(getProcAddr("glRasterPos4iv"))
	if gpRasterPos4iv == nil {
		return errors.New("glRasterPos4iv")
	}
	gpRasterPos4s = (C.GPRASTERPOS4S)(getProcAddr("glRasterPos4s"))
	if gpRasterPos4s == nil {
		return errors.New("glRasterPos4s")
	}
	gpRasterPos4sv = (C.GPRASTERPOS4SV)(getProcAddr("glRasterPos4sv"))
	if gpRasterPos4sv == nil {
		return errors.New("glRasterPos4sv")
	}
	gpRasterPos4xOES = (C.GPRASTERPOS4XOES)(getProcAddr("glRasterPos4xOES"))
	gpRasterPos4xvOES = (C.GPRASTERPOS4XVOES)(getProcAddr("glRasterPos4xvOES"))
	gpRasterSamplesEXT = (C.GPRASTERSAMPLESEXT)(getProcAddr("glRasterSamplesEXT"))
	gpReadBuffer = (C.GPREADBUFFER)(getProcAddr("glReadBuffer"))
	if gpReadBuffer == nil {
		return errors.New("glReadBuffer")
	}
	gpReadInstrumentsSGIX = (C.GPREADINSTRUMENTSSGIX)(getProcAddr("glReadInstrumentsSGIX"))
	gpReadPixels = (C.GPREADPIXELS)(getProcAddr("glReadPixels"))
	if gpReadPixels == nil {
		return errors.New("glReadPixels")
	}
	gpReadnPixels = (C.GPREADNPIXELS)(getProcAddr("glReadnPixels"))
	if gpReadnPixels == nil {
		return errors.New("glReadnPixels")
	}
	gpReadnPixelsARB = (C.GPREADNPIXELSARB)(getProcAddr("glReadnPixelsARB"))
	gpReadnPixelsKHR = (C.GPREADNPIXELSKHR)(getProcAddr("glReadnPixelsKHR"))
	gpRectd = (C.GPRECTD)(getProcAddr("glRectd"))
	if gpRectd == nil {
		return errors.New("glRectd")
	}
	gpRectdv = (C.GPRECTDV)(getProcAddr("glRectdv"))
	if gpRectdv == nil {
		return errors.New("glRectdv")
	}
	gpRectf = (C.GPRECTF)(getProcAddr("glRectf"))
	if gpRectf == nil {
		return errors.New("glRectf")
	}
	gpRectfv = (C.GPRECTFV)(getProcAddr("glRectfv"))
	if gpRectfv == nil {
		return errors.New("glRectfv")
	}
	gpRecti = (C.GPRECTI)(getProcAddr("glRecti"))
	if gpRecti == nil {
		return errors.New("glRecti")
	}
	gpRectiv = (C.GPRECTIV)(getProcAddr("glRectiv"))
	if gpRectiv == nil {
		return errors.New("glRectiv")
	}
	gpRects = (C.GPRECTS)(getProcAddr("glRects"))
	if gpRects == nil {
		return errors.New("glRects")
	}
	gpRectsv = (C.GPRECTSV)(getProcAddr("glRectsv"))
	if gpRectsv == nil {
		return errors.New("glRectsv")
	}
	gpRectxOES = (C.GPRECTXOES)(getProcAddr("glRectxOES"))
	gpRectxvOES = (C.GPRECTXVOES)(getProcAddr("glRectxvOES"))
	gpReferencePlaneSGIX = (C.GPREFERENCEPLANESGIX)(getProcAddr("glReferencePlaneSGIX"))
	gpReleaseKeyedMutexWin32EXT = (C.GPRELEASEKEYEDMUTEXWIN32EXT)(getProcAddr("glReleaseKeyedMutexWin32EXT"))
	gpReleaseShaderCompiler = (C.GPRELEASESHADERCOMPILER)(getProcAddr("glReleaseShaderCompiler"))
	if gpReleaseShaderCompiler == nil {
		return errors.New("glReleaseShaderCompiler")
	}
	gpRenderGpuMaskNV = (C.GPRENDERGPUMASKNV)(getProcAddr("glRenderGpuMaskNV"))
	gpRenderMode = (C.GPRENDERMODE)(getProcAddr("glRenderMode"))
	if gpRenderMode == nil {
		return errors.New("glRenderMode")
	}
	gpRenderbufferStorage = (C.GPRENDERBUFFERSTORAGE)(getProcAddr("glRenderbufferStorage"))
	if gpRenderbufferStorage == nil {
		return errors.New("glRenderbufferStorage")
	}
	gpRenderbufferStorageEXT = (C.GPRENDERBUFFERSTORAGEEXT)(getProcAddr("glRenderbufferStorageEXT"))
	gpRenderbufferStorageMultisample = (C.GPRENDERBUFFERSTORAGEMULTISAMPLE)(getProcAddr("glRenderbufferStorageMultisample"))
	if gpRenderbufferStorageMultisample == nil {
		return errors.New("glRenderbufferStorageMultisample")
	}
	gpRenderbufferStorageMultisampleCoverageNV = (C.GPRENDERBUFFERSTORAGEMULTISAMPLECOVERAGENV)(getProcAddr("glRenderbufferStorageMultisampleCoverageNV"))
	gpRenderbufferStorageMultisampleEXT = (C.GPRENDERBUFFERSTORAGEMULTISAMPLEEXT)(getProcAddr("glRenderbufferStorageMultisampleEXT"))
	gpReplacementCodePointerSUN = (C.GPREPLACEMENTCODEPOINTERSUN)(getProcAddr("glReplacementCodePointerSUN"))
	gpReplacementCodeubSUN = (C.GPREPLACEMENTCODEUBSUN)(getProcAddr("glReplacementCodeubSUN"))
	gpReplacementCodeubvSUN = (C.GPREPLACEMENTCODEUBVSUN)(getProcAddr("glReplacementCodeubvSUN"))
	gpReplacementCodeuiColor3fVertex3fSUN = (C.GPREPLACEMENTCODEUICOLOR3FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiColor3fVertex3fSUN"))
	gpReplacementCodeuiColor3fVertex3fvSUN = (C.GPREPLACEMENTCODEUICOLOR3FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiColor3fVertex3fvSUN"))
	gpReplacementCodeuiColor4fNormal3fVertex3fSUN = (C.GPREPLACEMENTCODEUICOLOR4FNORMAL3FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiColor4fNormal3fVertex3fSUN"))
	gpReplacementCodeuiColor4fNormal3fVertex3fvSUN = (C.GPREPLACEMENTCODEUICOLOR4FNORMAL3FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiColor4fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiColor4ubVertex3fSUN = (C.GPREPLACEMENTCODEUICOLOR4UBVERTEX3FSUN)(getProcAddr("glReplacementCodeuiColor4ubVertex3fSUN"))
	gpReplacementCodeuiColor4ubVertex3fvSUN = (C.GPREPLACEMENTCODEUICOLOR4UBVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiColor4ubVertex3fvSUN"))
	gpReplacementCodeuiNormal3fVertex3fSUN = (C.GPREPLACEMENTCODEUINORMAL3FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiNormal3fVertex3fSUN"))
	gpReplacementCodeuiNormal3fVertex3fvSUN = (C.GPREPLACEMENTCODEUINORMAL3FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiNormal3fVertex3fvSUN"))
	gpReplacementCodeuiSUN = (C.GPREPLACEMENTCODEUISUN)(getProcAddr("glReplacementCodeuiSUN"))
	gpReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FCOLOR4FNORMAL3FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fvSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FCOLOR4FNORMAL3FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiTexCoord2fNormal3fVertex3fSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FNORMAL3FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiTexCoord2fNormal3fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fNormal3fVertex3fvSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FNORMAL3FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiTexCoord2fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiTexCoord2fVertex3fSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FVERTEX3FSUN)(getProcAddr("glReplacementCodeuiTexCoord2fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fVertex3fvSUN = (C.GPREPLACEMENTCODEUITEXCOORD2FVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiTexCoord2fVertex3fvSUN"))
	gpReplacementCodeuiVertex3fSUN = (C.GPREPLACEMENTCODEUIVERTEX3FSUN)(getProcAddr("glReplacementCodeuiVertex3fSUN"))
	gpReplacementCodeuiVertex3fvSUN = (C.GPREPLACEMENTCODEUIVERTEX3FVSUN)(getProcAddr("glReplacementCodeuiVertex3fvSUN"))
	gpReplacementCodeuivSUN = (C.GPREPLACEMENTCODEUIVSUN)(getProcAddr("glReplacementCodeuivSUN"))
	gpReplacementCodeusSUN = (C.GPREPLACEMENTCODEUSSUN)(getProcAddr("glReplacementCodeusSUN"))
	gpReplacementCodeusvSUN = (C.GPREPLACEMENTCODEUSVSUN)(getProcAddr("glReplacementCodeusvSUN"))
	gpRequestResidentProgramsNV = (C.GPREQUESTRESIDENTPROGRAMSNV)(getProcAddr("glRequestResidentProgramsNV"))
	gpResetHistogram = (C.GPRESETHISTOGRAM)(getProcAddr("glResetHistogram"))
	gpResetHistogramEXT = (C.GPRESETHISTOGRAMEXT)(getProcAddr("glResetHistogramEXT"))
	gpResetMinmax = (C.GPRESETMINMAX)(getProcAddr("glResetMinmax"))
	gpResetMinmaxEXT = (C.GPRESETMINMAXEXT)(getProcAddr("glResetMinmaxEXT"))
	gpResizeBuffersMESA = (C.GPRESIZEBUFFERSMESA)(getProcAddr("glResizeBuffersMESA"))
	gpResolveDepthValuesNV = (C.GPRESOLVEDEPTHVALUESNV)(getProcAddr("glResolveDepthValuesNV"))
	gpResumeTransformFeedback = (C.GPRESUMETRANSFORMFEEDBACK)(getProcAddr("glResumeTransformFeedback"))
	if gpResumeTransformFeedback == nil {
		return errors.New("glResumeTransformFeedback")
	}
	gpResumeTransformFeedbackNV = (C.GPRESUMETRANSFORMFEEDBACKNV)(getProcAddr("glResumeTransformFeedbackNV"))
	gpRotated = (C.GPROTATED)(getProcAddr("glRotated"))
	if gpRotated == nil {
		return errors.New("glRotated")
	}
	gpRotatef = (C.GPROTATEF)(getProcAddr("glRotatef"))
	if gpRotatef == nil {
		return errors.New("glRotatef")
	}
	gpRotatexOES = (C.GPROTATEXOES)(getProcAddr("glRotatexOES"))
	gpSampleCoverage = (C.GPSAMPLECOVERAGE)(getProcAddr("glSampleCoverage"))
	if gpSampleCoverage == nil {
		return errors.New("glSampleCoverage")
	}
	gpSampleCoverageARB = (C.GPSAMPLECOVERAGEARB)(getProcAddr("glSampleCoverageARB"))
	gpSampleCoveragexOES = (C.GPSAMPLECOVERAGEXOES)(getProcAddr("glSampleCoveragexOES"))
	gpSampleMapATI = (C.GPSAMPLEMAPATI)(getProcAddr("glSampleMapATI"))
	gpSampleMaskEXT = (C.GPSAMPLEMASKEXT)(getProcAddr("glSampleMaskEXT"))
	gpSampleMaskIndexedNV = (C.GPSAMPLEMASKINDEXEDNV)(getProcAddr("glSampleMaskIndexedNV"))
	gpSampleMaskSGIS = (C.GPSAMPLEMASKSGIS)(getProcAddr("glSampleMaskSGIS"))
	gpSampleMaski = (C.GPSAMPLEMASKI)(getProcAddr("glSampleMaski"))
	if gpSampleMaski == nil {
		return errors.New("glSampleMaski")
	}
	gpSamplePatternEXT = (C.GPSAMPLEPATTERNEXT)(getProcAddr("glSamplePatternEXT"))
	gpSamplePatternSGIS = (C.GPSAMPLEPATTERNSGIS)(getProcAddr("glSamplePatternSGIS"))
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
	gpScaled = (C.GPSCALED)(getProcAddr("glScaled"))
	if gpScaled == nil {
		return errors.New("glScaled")
	}
	gpScalef = (C.GPSCALEF)(getProcAddr("glScalef"))
	if gpScalef == nil {
		return errors.New("glScalef")
	}
	gpScalexOES = (C.GPSCALEXOES)(getProcAddr("glScalexOES"))
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
	gpSecondaryColor3b = (C.GPSECONDARYCOLOR3B)(getProcAddr("glSecondaryColor3b"))
	if gpSecondaryColor3b == nil {
		return errors.New("glSecondaryColor3b")
	}
	gpSecondaryColor3bEXT = (C.GPSECONDARYCOLOR3BEXT)(getProcAddr("glSecondaryColor3bEXT"))
	gpSecondaryColor3bv = (C.GPSECONDARYCOLOR3BV)(getProcAddr("glSecondaryColor3bv"))
	if gpSecondaryColor3bv == nil {
		return errors.New("glSecondaryColor3bv")
	}
	gpSecondaryColor3bvEXT = (C.GPSECONDARYCOLOR3BVEXT)(getProcAddr("glSecondaryColor3bvEXT"))
	gpSecondaryColor3d = (C.GPSECONDARYCOLOR3D)(getProcAddr("glSecondaryColor3d"))
	if gpSecondaryColor3d == nil {
		return errors.New("glSecondaryColor3d")
	}
	gpSecondaryColor3dEXT = (C.GPSECONDARYCOLOR3DEXT)(getProcAddr("glSecondaryColor3dEXT"))
	gpSecondaryColor3dv = (C.GPSECONDARYCOLOR3DV)(getProcAddr("glSecondaryColor3dv"))
	if gpSecondaryColor3dv == nil {
		return errors.New("glSecondaryColor3dv")
	}
	gpSecondaryColor3dvEXT = (C.GPSECONDARYCOLOR3DVEXT)(getProcAddr("glSecondaryColor3dvEXT"))
	gpSecondaryColor3f = (C.GPSECONDARYCOLOR3F)(getProcAddr("glSecondaryColor3f"))
	if gpSecondaryColor3f == nil {
		return errors.New("glSecondaryColor3f")
	}
	gpSecondaryColor3fEXT = (C.GPSECONDARYCOLOR3FEXT)(getProcAddr("glSecondaryColor3fEXT"))
	gpSecondaryColor3fv = (C.GPSECONDARYCOLOR3FV)(getProcAddr("glSecondaryColor3fv"))
	if gpSecondaryColor3fv == nil {
		return errors.New("glSecondaryColor3fv")
	}
	gpSecondaryColor3fvEXT = (C.GPSECONDARYCOLOR3FVEXT)(getProcAddr("glSecondaryColor3fvEXT"))
	gpSecondaryColor3hNV = (C.GPSECONDARYCOLOR3HNV)(getProcAddr("glSecondaryColor3hNV"))
	gpSecondaryColor3hvNV = (C.GPSECONDARYCOLOR3HVNV)(getProcAddr("glSecondaryColor3hvNV"))
	gpSecondaryColor3i = (C.GPSECONDARYCOLOR3I)(getProcAddr("glSecondaryColor3i"))
	if gpSecondaryColor3i == nil {
		return errors.New("glSecondaryColor3i")
	}
	gpSecondaryColor3iEXT = (C.GPSECONDARYCOLOR3IEXT)(getProcAddr("glSecondaryColor3iEXT"))
	gpSecondaryColor3iv = (C.GPSECONDARYCOLOR3IV)(getProcAddr("glSecondaryColor3iv"))
	if gpSecondaryColor3iv == nil {
		return errors.New("glSecondaryColor3iv")
	}
	gpSecondaryColor3ivEXT = (C.GPSECONDARYCOLOR3IVEXT)(getProcAddr("glSecondaryColor3ivEXT"))
	gpSecondaryColor3s = (C.GPSECONDARYCOLOR3S)(getProcAddr("glSecondaryColor3s"))
	if gpSecondaryColor3s == nil {
		return errors.New("glSecondaryColor3s")
	}
	gpSecondaryColor3sEXT = (C.GPSECONDARYCOLOR3SEXT)(getProcAddr("glSecondaryColor3sEXT"))
	gpSecondaryColor3sv = (C.GPSECONDARYCOLOR3SV)(getProcAddr("glSecondaryColor3sv"))
	if gpSecondaryColor3sv == nil {
		return errors.New("glSecondaryColor3sv")
	}
	gpSecondaryColor3svEXT = (C.GPSECONDARYCOLOR3SVEXT)(getProcAddr("glSecondaryColor3svEXT"))
	gpSecondaryColor3ub = (C.GPSECONDARYCOLOR3UB)(getProcAddr("glSecondaryColor3ub"))
	if gpSecondaryColor3ub == nil {
		return errors.New("glSecondaryColor3ub")
	}
	gpSecondaryColor3ubEXT = (C.GPSECONDARYCOLOR3UBEXT)(getProcAddr("glSecondaryColor3ubEXT"))
	gpSecondaryColor3ubv = (C.GPSECONDARYCOLOR3UBV)(getProcAddr("glSecondaryColor3ubv"))
	if gpSecondaryColor3ubv == nil {
		return errors.New("glSecondaryColor3ubv")
	}
	gpSecondaryColor3ubvEXT = (C.GPSECONDARYCOLOR3UBVEXT)(getProcAddr("glSecondaryColor3ubvEXT"))
	gpSecondaryColor3ui = (C.GPSECONDARYCOLOR3UI)(getProcAddr("glSecondaryColor3ui"))
	if gpSecondaryColor3ui == nil {
		return errors.New("glSecondaryColor3ui")
	}
	gpSecondaryColor3uiEXT = (C.GPSECONDARYCOLOR3UIEXT)(getProcAddr("glSecondaryColor3uiEXT"))
	gpSecondaryColor3uiv = (C.GPSECONDARYCOLOR3UIV)(getProcAddr("glSecondaryColor3uiv"))
	if gpSecondaryColor3uiv == nil {
		return errors.New("glSecondaryColor3uiv")
	}
	gpSecondaryColor3uivEXT = (C.GPSECONDARYCOLOR3UIVEXT)(getProcAddr("glSecondaryColor3uivEXT"))
	gpSecondaryColor3us = (C.GPSECONDARYCOLOR3US)(getProcAddr("glSecondaryColor3us"))
	if gpSecondaryColor3us == nil {
		return errors.New("glSecondaryColor3us")
	}
	gpSecondaryColor3usEXT = (C.GPSECONDARYCOLOR3USEXT)(getProcAddr("glSecondaryColor3usEXT"))
	gpSecondaryColor3usv = (C.GPSECONDARYCOLOR3USV)(getProcAddr("glSecondaryColor3usv"))
	if gpSecondaryColor3usv == nil {
		return errors.New("glSecondaryColor3usv")
	}
	gpSecondaryColor3usvEXT = (C.GPSECONDARYCOLOR3USVEXT)(getProcAddr("glSecondaryColor3usvEXT"))
	gpSecondaryColorFormatNV = (C.GPSECONDARYCOLORFORMATNV)(getProcAddr("glSecondaryColorFormatNV"))
	gpSecondaryColorP3ui = (C.GPSECONDARYCOLORP3UI)(getProcAddr("glSecondaryColorP3ui"))
	if gpSecondaryColorP3ui == nil {
		return errors.New("glSecondaryColorP3ui")
	}
	gpSecondaryColorP3uiv = (C.GPSECONDARYCOLORP3UIV)(getProcAddr("glSecondaryColorP3uiv"))
	if gpSecondaryColorP3uiv == nil {
		return errors.New("glSecondaryColorP3uiv")
	}
	gpSecondaryColorPointer = (C.GPSECONDARYCOLORPOINTER)(getProcAddr("glSecondaryColorPointer"))
	if gpSecondaryColorPointer == nil {
		return errors.New("glSecondaryColorPointer")
	}
	gpSecondaryColorPointerEXT = (C.GPSECONDARYCOLORPOINTEREXT)(getProcAddr("glSecondaryColorPointerEXT"))
	gpSecondaryColorPointerListIBM = (C.GPSECONDARYCOLORPOINTERLISTIBM)(getProcAddr("glSecondaryColorPointerListIBM"))
	gpSelectBuffer = (C.GPSELECTBUFFER)(getProcAddr("glSelectBuffer"))
	if gpSelectBuffer == nil {
		return errors.New("glSelectBuffer")
	}
	gpSelectPerfMonitorCountersAMD = (C.GPSELECTPERFMONITORCOUNTERSAMD)(getProcAddr("glSelectPerfMonitorCountersAMD"))
	gpSemaphoreParameterui64vEXT = (C.GPSEMAPHOREPARAMETERUI64VEXT)(getProcAddr("glSemaphoreParameterui64vEXT"))
	gpSeparableFilter2D = (C.GPSEPARABLEFILTER2D)(getProcAddr("glSeparableFilter2D"))
	gpSeparableFilter2DEXT = (C.GPSEPARABLEFILTER2DEXT)(getProcAddr("glSeparableFilter2DEXT"))
	gpSetFenceAPPLE = (C.GPSETFENCEAPPLE)(getProcAddr("glSetFenceAPPLE"))
	gpSetFenceNV = (C.GPSETFENCENV)(getProcAddr("glSetFenceNV"))
	gpSetFragmentShaderConstantATI = (C.GPSETFRAGMENTSHADERCONSTANTATI)(getProcAddr("glSetFragmentShaderConstantATI"))
	gpSetInvariantEXT = (C.GPSETINVARIANTEXT)(getProcAddr("glSetInvariantEXT"))
	gpSetLocalConstantEXT = (C.GPSETLOCALCONSTANTEXT)(getProcAddr("glSetLocalConstantEXT"))
	gpSetMultisamplefvAMD = (C.GPSETMULTISAMPLEFVAMD)(getProcAddr("glSetMultisamplefvAMD"))
	gpShadeModel = (C.GPSHADEMODEL)(getProcAddr("glShadeModel"))
	if gpShadeModel == nil {
		return errors.New("glShadeModel")
	}
	gpShaderBinary = (C.GPSHADERBINARY)(getProcAddr("glShaderBinary"))
	if gpShaderBinary == nil {
		return errors.New("glShaderBinary")
	}
	gpShaderOp1EXT = (C.GPSHADEROP1EXT)(getProcAddr("glShaderOp1EXT"))
	gpShaderOp2EXT = (C.GPSHADEROP2EXT)(getProcAddr("glShaderOp2EXT"))
	gpShaderOp3EXT = (C.GPSHADEROP3EXT)(getProcAddr("glShaderOp3EXT"))
	gpShaderSource = (C.GPSHADERSOURCE)(getProcAddr("glShaderSource"))
	if gpShaderSource == nil {
		return errors.New("glShaderSource")
	}
	gpShaderSourceARB = (C.GPSHADERSOURCEARB)(getProcAddr("glShaderSourceARB"))
	gpShaderStorageBlockBinding = (C.GPSHADERSTORAGEBLOCKBINDING)(getProcAddr("glShaderStorageBlockBinding"))
	if gpShaderStorageBlockBinding == nil {
		return errors.New("glShaderStorageBlockBinding")
	}
	gpSharpenTexFuncSGIS = (C.GPSHARPENTEXFUNCSGIS)(getProcAddr("glSharpenTexFuncSGIS"))
	gpSignalSemaphoreEXT = (C.GPSIGNALSEMAPHOREEXT)(getProcAddr("glSignalSemaphoreEXT"))
	gpSignalVkFenceNV = (C.GPSIGNALVKFENCENV)(getProcAddr("glSignalVkFenceNV"))
	gpSignalVkSemaphoreNV = (C.GPSIGNALVKSEMAPHORENV)(getProcAddr("glSignalVkSemaphoreNV"))
	gpSpecializeShaderARB = (C.GPSPECIALIZESHADERARB)(getProcAddr("glSpecializeShaderARB"))
	gpSpriteParameterfSGIX = (C.GPSPRITEPARAMETERFSGIX)(getProcAddr("glSpriteParameterfSGIX"))
	gpSpriteParameterfvSGIX = (C.GPSPRITEPARAMETERFVSGIX)(getProcAddr("glSpriteParameterfvSGIX"))
	gpSpriteParameteriSGIX = (C.GPSPRITEPARAMETERISGIX)(getProcAddr("glSpriteParameteriSGIX"))
	gpSpriteParameterivSGIX = (C.GPSPRITEPARAMETERIVSGIX)(getProcAddr("glSpriteParameterivSGIX"))
	gpStartInstrumentsSGIX = (C.GPSTARTINSTRUMENTSSGIX)(getProcAddr("glStartInstrumentsSGIX"))
	gpStateCaptureNV = (C.GPSTATECAPTURENV)(getProcAddr("glStateCaptureNV"))
	gpStencilClearTagEXT = (C.GPSTENCILCLEARTAGEXT)(getProcAddr("glStencilClearTagEXT"))
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
	gpStencilFuncSeparateATI = (C.GPSTENCILFUNCSEPARATEATI)(getProcAddr("glStencilFuncSeparateATI"))
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
	gpStencilOpSeparateATI = (C.GPSTENCILOPSEPARATEATI)(getProcAddr("glStencilOpSeparateATI"))
	gpStencilOpValueAMD = (C.GPSTENCILOPVALUEAMD)(getProcAddr("glStencilOpValueAMD"))
	gpStencilStrokePathInstancedNV = (C.GPSTENCILSTROKEPATHINSTANCEDNV)(getProcAddr("glStencilStrokePathInstancedNV"))
	gpStencilStrokePathNV = (C.GPSTENCILSTROKEPATHNV)(getProcAddr("glStencilStrokePathNV"))
	gpStencilThenCoverFillPathInstancedNV = (C.GPSTENCILTHENCOVERFILLPATHINSTANCEDNV)(getProcAddr("glStencilThenCoverFillPathInstancedNV"))
	gpStencilThenCoverFillPathNV = (C.GPSTENCILTHENCOVERFILLPATHNV)(getProcAddr("glStencilThenCoverFillPathNV"))
	gpStencilThenCoverStrokePathInstancedNV = (C.GPSTENCILTHENCOVERSTROKEPATHINSTANCEDNV)(getProcAddr("glStencilThenCoverStrokePathInstancedNV"))
	gpStencilThenCoverStrokePathNV = (C.GPSTENCILTHENCOVERSTROKEPATHNV)(getProcAddr("glStencilThenCoverStrokePathNV"))
	gpStopInstrumentsSGIX = (C.GPSTOPINSTRUMENTSSGIX)(getProcAddr("glStopInstrumentsSGIX"))
	gpStringMarkerGREMEDY = (C.GPSTRINGMARKERGREMEDY)(getProcAddr("glStringMarkerGREMEDY"))
	gpSubpixelPrecisionBiasNV = (C.GPSUBPIXELPRECISIONBIASNV)(getProcAddr("glSubpixelPrecisionBiasNV"))
	gpSwizzleEXT = (C.GPSWIZZLEEXT)(getProcAddr("glSwizzleEXT"))
	gpSyncTextureINTEL = (C.GPSYNCTEXTUREINTEL)(getProcAddr("glSyncTextureINTEL"))
	gpTagSampleBufferSGIX = (C.GPTAGSAMPLEBUFFERSGIX)(getProcAddr("glTagSampleBufferSGIX"))
	gpTangent3bEXT = (C.GPTANGENT3BEXT)(getProcAddr("glTangent3bEXT"))
	gpTangent3bvEXT = (C.GPTANGENT3BVEXT)(getProcAddr("glTangent3bvEXT"))
	gpTangent3dEXT = (C.GPTANGENT3DEXT)(getProcAddr("glTangent3dEXT"))
	gpTangent3dvEXT = (C.GPTANGENT3DVEXT)(getProcAddr("glTangent3dvEXT"))
	gpTangent3fEXT = (C.GPTANGENT3FEXT)(getProcAddr("glTangent3fEXT"))
	gpTangent3fvEXT = (C.GPTANGENT3FVEXT)(getProcAddr("glTangent3fvEXT"))
	gpTangent3iEXT = (C.GPTANGENT3IEXT)(getProcAddr("glTangent3iEXT"))
	gpTangent3ivEXT = (C.GPTANGENT3IVEXT)(getProcAddr("glTangent3ivEXT"))
	gpTangent3sEXT = (C.GPTANGENT3SEXT)(getProcAddr("glTangent3sEXT"))
	gpTangent3svEXT = (C.GPTANGENT3SVEXT)(getProcAddr("glTangent3svEXT"))
	gpTangentPointerEXT = (C.GPTANGENTPOINTEREXT)(getProcAddr("glTangentPointerEXT"))
	gpTbufferMask3DFX = (C.GPTBUFFERMASK3DFX)(getProcAddr("glTbufferMask3DFX"))
	gpTessellationFactorAMD = (C.GPTESSELLATIONFACTORAMD)(getProcAddr("glTessellationFactorAMD"))
	gpTessellationModeAMD = (C.GPTESSELLATIONMODEAMD)(getProcAddr("glTessellationModeAMD"))
	gpTestFenceAPPLE = (C.GPTESTFENCEAPPLE)(getProcAddr("glTestFenceAPPLE"))
	gpTestFenceNV = (C.GPTESTFENCENV)(getProcAddr("glTestFenceNV"))
	gpTestObjectAPPLE = (C.GPTESTOBJECTAPPLE)(getProcAddr("glTestObjectAPPLE"))
	gpTexBuffer = (C.GPTEXBUFFER)(getProcAddr("glTexBuffer"))
	if gpTexBuffer == nil {
		return errors.New("glTexBuffer")
	}
	gpTexBufferARB = (C.GPTEXBUFFERARB)(getProcAddr("glTexBufferARB"))
	gpTexBufferEXT = (C.GPTEXBUFFEREXT)(getProcAddr("glTexBufferEXT"))
	gpTexBufferRange = (C.GPTEXBUFFERRANGE)(getProcAddr("glTexBufferRange"))
	if gpTexBufferRange == nil {
		return errors.New("glTexBufferRange")
	}
	gpTexBumpParameterfvATI = (C.GPTEXBUMPPARAMETERFVATI)(getProcAddr("glTexBumpParameterfvATI"))
	gpTexBumpParameterivATI = (C.GPTEXBUMPPARAMETERIVATI)(getProcAddr("glTexBumpParameterivATI"))
	gpTexCoord1bOES = (C.GPTEXCOORD1BOES)(getProcAddr("glTexCoord1bOES"))
	gpTexCoord1bvOES = (C.GPTEXCOORD1BVOES)(getProcAddr("glTexCoord1bvOES"))
	gpTexCoord1d = (C.GPTEXCOORD1D)(getProcAddr("glTexCoord1d"))
	if gpTexCoord1d == nil {
		return errors.New("glTexCoord1d")
	}
	gpTexCoord1dv = (C.GPTEXCOORD1DV)(getProcAddr("glTexCoord1dv"))
	if gpTexCoord1dv == nil {
		return errors.New("glTexCoord1dv")
	}
	gpTexCoord1f = (C.GPTEXCOORD1F)(getProcAddr("glTexCoord1f"))
	if gpTexCoord1f == nil {
		return errors.New("glTexCoord1f")
	}
	gpTexCoord1fv = (C.GPTEXCOORD1FV)(getProcAddr("glTexCoord1fv"))
	if gpTexCoord1fv == nil {
		return errors.New("glTexCoord1fv")
	}
	gpTexCoord1hNV = (C.GPTEXCOORD1HNV)(getProcAddr("glTexCoord1hNV"))
	gpTexCoord1hvNV = (C.GPTEXCOORD1HVNV)(getProcAddr("glTexCoord1hvNV"))
	gpTexCoord1i = (C.GPTEXCOORD1I)(getProcAddr("glTexCoord1i"))
	if gpTexCoord1i == nil {
		return errors.New("glTexCoord1i")
	}
	gpTexCoord1iv = (C.GPTEXCOORD1IV)(getProcAddr("glTexCoord1iv"))
	if gpTexCoord1iv == nil {
		return errors.New("glTexCoord1iv")
	}
	gpTexCoord1s = (C.GPTEXCOORD1S)(getProcAddr("glTexCoord1s"))
	if gpTexCoord1s == nil {
		return errors.New("glTexCoord1s")
	}
	gpTexCoord1sv = (C.GPTEXCOORD1SV)(getProcAddr("glTexCoord1sv"))
	if gpTexCoord1sv == nil {
		return errors.New("glTexCoord1sv")
	}
	gpTexCoord1xOES = (C.GPTEXCOORD1XOES)(getProcAddr("glTexCoord1xOES"))
	gpTexCoord1xvOES = (C.GPTEXCOORD1XVOES)(getProcAddr("glTexCoord1xvOES"))
	gpTexCoord2bOES = (C.GPTEXCOORD2BOES)(getProcAddr("glTexCoord2bOES"))
	gpTexCoord2bvOES = (C.GPTEXCOORD2BVOES)(getProcAddr("glTexCoord2bvOES"))
	gpTexCoord2d = (C.GPTEXCOORD2D)(getProcAddr("glTexCoord2d"))
	if gpTexCoord2d == nil {
		return errors.New("glTexCoord2d")
	}
	gpTexCoord2dv = (C.GPTEXCOORD2DV)(getProcAddr("glTexCoord2dv"))
	if gpTexCoord2dv == nil {
		return errors.New("glTexCoord2dv")
	}
	gpTexCoord2f = (C.GPTEXCOORD2F)(getProcAddr("glTexCoord2f"))
	if gpTexCoord2f == nil {
		return errors.New("glTexCoord2f")
	}
	gpTexCoord2fColor3fVertex3fSUN = (C.GPTEXCOORD2FCOLOR3FVERTEX3FSUN)(getProcAddr("glTexCoord2fColor3fVertex3fSUN"))
	gpTexCoord2fColor3fVertex3fvSUN = (C.GPTEXCOORD2FCOLOR3FVERTEX3FVSUN)(getProcAddr("glTexCoord2fColor3fVertex3fvSUN"))
	gpTexCoord2fColor4fNormal3fVertex3fSUN = (C.GPTEXCOORD2FCOLOR4FNORMAL3FVERTEX3FSUN)(getProcAddr("glTexCoord2fColor4fNormal3fVertex3fSUN"))
	gpTexCoord2fColor4fNormal3fVertex3fvSUN = (C.GPTEXCOORD2FCOLOR4FNORMAL3FVERTEX3FVSUN)(getProcAddr("glTexCoord2fColor4fNormal3fVertex3fvSUN"))
	gpTexCoord2fColor4ubVertex3fSUN = (C.GPTEXCOORD2FCOLOR4UBVERTEX3FSUN)(getProcAddr("glTexCoord2fColor4ubVertex3fSUN"))
	gpTexCoord2fColor4ubVertex3fvSUN = (C.GPTEXCOORD2FCOLOR4UBVERTEX3FVSUN)(getProcAddr("glTexCoord2fColor4ubVertex3fvSUN"))
	gpTexCoord2fNormal3fVertex3fSUN = (C.GPTEXCOORD2FNORMAL3FVERTEX3FSUN)(getProcAddr("glTexCoord2fNormal3fVertex3fSUN"))
	gpTexCoord2fNormal3fVertex3fvSUN = (C.GPTEXCOORD2FNORMAL3FVERTEX3FVSUN)(getProcAddr("glTexCoord2fNormal3fVertex3fvSUN"))
	gpTexCoord2fVertex3fSUN = (C.GPTEXCOORD2FVERTEX3FSUN)(getProcAddr("glTexCoord2fVertex3fSUN"))
	gpTexCoord2fVertex3fvSUN = (C.GPTEXCOORD2FVERTEX3FVSUN)(getProcAddr("glTexCoord2fVertex3fvSUN"))
	gpTexCoord2fv = (C.GPTEXCOORD2FV)(getProcAddr("glTexCoord2fv"))
	if gpTexCoord2fv == nil {
		return errors.New("glTexCoord2fv")
	}
	gpTexCoord2hNV = (C.GPTEXCOORD2HNV)(getProcAddr("glTexCoord2hNV"))
	gpTexCoord2hvNV = (C.GPTEXCOORD2HVNV)(getProcAddr("glTexCoord2hvNV"))
	gpTexCoord2i = (C.GPTEXCOORD2I)(getProcAddr("glTexCoord2i"))
	if gpTexCoord2i == nil {
		return errors.New("glTexCoord2i")
	}
	gpTexCoord2iv = (C.GPTEXCOORD2IV)(getProcAddr("glTexCoord2iv"))
	if gpTexCoord2iv == nil {
		return errors.New("glTexCoord2iv")
	}
	gpTexCoord2s = (C.GPTEXCOORD2S)(getProcAddr("glTexCoord2s"))
	if gpTexCoord2s == nil {
		return errors.New("glTexCoord2s")
	}
	gpTexCoord2sv = (C.GPTEXCOORD2SV)(getProcAddr("glTexCoord2sv"))
	if gpTexCoord2sv == nil {
		return errors.New("glTexCoord2sv")
	}
	gpTexCoord2xOES = (C.GPTEXCOORD2XOES)(getProcAddr("glTexCoord2xOES"))
	gpTexCoord2xvOES = (C.GPTEXCOORD2XVOES)(getProcAddr("glTexCoord2xvOES"))
	gpTexCoord3bOES = (C.GPTEXCOORD3BOES)(getProcAddr("glTexCoord3bOES"))
	gpTexCoord3bvOES = (C.GPTEXCOORD3BVOES)(getProcAddr("glTexCoord3bvOES"))
	gpTexCoord3d = (C.GPTEXCOORD3D)(getProcAddr("glTexCoord3d"))
	if gpTexCoord3d == nil {
		return errors.New("glTexCoord3d")
	}
	gpTexCoord3dv = (C.GPTEXCOORD3DV)(getProcAddr("glTexCoord3dv"))
	if gpTexCoord3dv == nil {
		return errors.New("glTexCoord3dv")
	}
	gpTexCoord3f = (C.GPTEXCOORD3F)(getProcAddr("glTexCoord3f"))
	if gpTexCoord3f == nil {
		return errors.New("glTexCoord3f")
	}
	gpTexCoord3fv = (C.GPTEXCOORD3FV)(getProcAddr("glTexCoord3fv"))
	if gpTexCoord3fv == nil {
		return errors.New("glTexCoord3fv")
	}
	gpTexCoord3hNV = (C.GPTEXCOORD3HNV)(getProcAddr("glTexCoord3hNV"))
	gpTexCoord3hvNV = (C.GPTEXCOORD3HVNV)(getProcAddr("glTexCoord3hvNV"))
	gpTexCoord3i = (C.GPTEXCOORD3I)(getProcAddr("glTexCoord3i"))
	if gpTexCoord3i == nil {
		return errors.New("glTexCoord3i")
	}
	gpTexCoord3iv = (C.GPTEXCOORD3IV)(getProcAddr("glTexCoord3iv"))
	if gpTexCoord3iv == nil {
		return errors.New("glTexCoord3iv")
	}
	gpTexCoord3s = (C.GPTEXCOORD3S)(getProcAddr("glTexCoord3s"))
	if gpTexCoord3s == nil {
		return errors.New("glTexCoord3s")
	}
	gpTexCoord3sv = (C.GPTEXCOORD3SV)(getProcAddr("glTexCoord3sv"))
	if gpTexCoord3sv == nil {
		return errors.New("glTexCoord3sv")
	}
	gpTexCoord3xOES = (C.GPTEXCOORD3XOES)(getProcAddr("glTexCoord3xOES"))
	gpTexCoord3xvOES = (C.GPTEXCOORD3XVOES)(getProcAddr("glTexCoord3xvOES"))
	gpTexCoord4bOES = (C.GPTEXCOORD4BOES)(getProcAddr("glTexCoord4bOES"))
	gpTexCoord4bvOES = (C.GPTEXCOORD4BVOES)(getProcAddr("glTexCoord4bvOES"))
	gpTexCoord4d = (C.GPTEXCOORD4D)(getProcAddr("glTexCoord4d"))
	if gpTexCoord4d == nil {
		return errors.New("glTexCoord4d")
	}
	gpTexCoord4dv = (C.GPTEXCOORD4DV)(getProcAddr("glTexCoord4dv"))
	if gpTexCoord4dv == nil {
		return errors.New("glTexCoord4dv")
	}
	gpTexCoord4f = (C.GPTEXCOORD4F)(getProcAddr("glTexCoord4f"))
	if gpTexCoord4f == nil {
		return errors.New("glTexCoord4f")
	}
	gpTexCoord4fColor4fNormal3fVertex4fSUN = (C.GPTEXCOORD4FCOLOR4FNORMAL3FVERTEX4FSUN)(getProcAddr("glTexCoord4fColor4fNormal3fVertex4fSUN"))
	gpTexCoord4fColor4fNormal3fVertex4fvSUN = (C.GPTEXCOORD4FCOLOR4FNORMAL3FVERTEX4FVSUN)(getProcAddr("glTexCoord4fColor4fNormal3fVertex4fvSUN"))
	gpTexCoord4fVertex4fSUN = (C.GPTEXCOORD4FVERTEX4FSUN)(getProcAddr("glTexCoord4fVertex4fSUN"))
	gpTexCoord4fVertex4fvSUN = (C.GPTEXCOORD4FVERTEX4FVSUN)(getProcAddr("glTexCoord4fVertex4fvSUN"))
	gpTexCoord4fv = (C.GPTEXCOORD4FV)(getProcAddr("glTexCoord4fv"))
	if gpTexCoord4fv == nil {
		return errors.New("glTexCoord4fv")
	}
	gpTexCoord4hNV = (C.GPTEXCOORD4HNV)(getProcAddr("glTexCoord4hNV"))
	gpTexCoord4hvNV = (C.GPTEXCOORD4HVNV)(getProcAddr("glTexCoord4hvNV"))
	gpTexCoord4i = (C.GPTEXCOORD4I)(getProcAddr("glTexCoord4i"))
	if gpTexCoord4i == nil {
		return errors.New("glTexCoord4i")
	}
	gpTexCoord4iv = (C.GPTEXCOORD4IV)(getProcAddr("glTexCoord4iv"))
	if gpTexCoord4iv == nil {
		return errors.New("glTexCoord4iv")
	}
	gpTexCoord4s = (C.GPTEXCOORD4S)(getProcAddr("glTexCoord4s"))
	if gpTexCoord4s == nil {
		return errors.New("glTexCoord4s")
	}
	gpTexCoord4sv = (C.GPTEXCOORD4SV)(getProcAddr("glTexCoord4sv"))
	if gpTexCoord4sv == nil {
		return errors.New("glTexCoord4sv")
	}
	gpTexCoord4xOES = (C.GPTEXCOORD4XOES)(getProcAddr("glTexCoord4xOES"))
	gpTexCoord4xvOES = (C.GPTEXCOORD4XVOES)(getProcAddr("glTexCoord4xvOES"))
	gpTexCoordFormatNV = (C.GPTEXCOORDFORMATNV)(getProcAddr("glTexCoordFormatNV"))
	gpTexCoordP1ui = (C.GPTEXCOORDP1UI)(getProcAddr("glTexCoordP1ui"))
	if gpTexCoordP1ui == nil {
		return errors.New("glTexCoordP1ui")
	}
	gpTexCoordP1uiv = (C.GPTEXCOORDP1UIV)(getProcAddr("glTexCoordP1uiv"))
	if gpTexCoordP1uiv == nil {
		return errors.New("glTexCoordP1uiv")
	}
	gpTexCoordP2ui = (C.GPTEXCOORDP2UI)(getProcAddr("glTexCoordP2ui"))
	if gpTexCoordP2ui == nil {
		return errors.New("glTexCoordP2ui")
	}
	gpTexCoordP2uiv = (C.GPTEXCOORDP2UIV)(getProcAddr("glTexCoordP2uiv"))
	if gpTexCoordP2uiv == nil {
		return errors.New("glTexCoordP2uiv")
	}
	gpTexCoordP3ui = (C.GPTEXCOORDP3UI)(getProcAddr("glTexCoordP3ui"))
	if gpTexCoordP3ui == nil {
		return errors.New("glTexCoordP3ui")
	}
	gpTexCoordP3uiv = (C.GPTEXCOORDP3UIV)(getProcAddr("glTexCoordP3uiv"))
	if gpTexCoordP3uiv == nil {
		return errors.New("glTexCoordP3uiv")
	}
	gpTexCoordP4ui = (C.GPTEXCOORDP4UI)(getProcAddr("glTexCoordP4ui"))
	if gpTexCoordP4ui == nil {
		return errors.New("glTexCoordP4ui")
	}
	gpTexCoordP4uiv = (C.GPTEXCOORDP4UIV)(getProcAddr("glTexCoordP4uiv"))
	if gpTexCoordP4uiv == nil {
		return errors.New("glTexCoordP4uiv")
	}
	gpTexCoordPointer = (C.GPTEXCOORDPOINTER)(getProcAddr("glTexCoordPointer"))
	if gpTexCoordPointer == nil {
		return errors.New("glTexCoordPointer")
	}
	gpTexCoordPointerEXT = (C.GPTEXCOORDPOINTEREXT)(getProcAddr("glTexCoordPointerEXT"))
	gpTexCoordPointerListIBM = (C.GPTEXCOORDPOINTERLISTIBM)(getProcAddr("glTexCoordPointerListIBM"))
	gpTexCoordPointervINTEL = (C.GPTEXCOORDPOINTERVINTEL)(getProcAddr("glTexCoordPointervINTEL"))
	gpTexEnvf = (C.GPTEXENVF)(getProcAddr("glTexEnvf"))
	if gpTexEnvf == nil {
		return errors.New("glTexEnvf")
	}
	gpTexEnvfv = (C.GPTEXENVFV)(getProcAddr("glTexEnvfv"))
	if gpTexEnvfv == nil {
		return errors.New("glTexEnvfv")
	}
	gpTexEnvi = (C.GPTEXENVI)(getProcAddr("glTexEnvi"))
	if gpTexEnvi == nil {
		return errors.New("glTexEnvi")
	}
	gpTexEnviv = (C.GPTEXENVIV)(getProcAddr("glTexEnviv"))
	if gpTexEnviv == nil {
		return errors.New("glTexEnviv")
	}
	gpTexEnvxOES = (C.GPTEXENVXOES)(getProcAddr("glTexEnvxOES"))
	gpTexEnvxvOES = (C.GPTEXENVXVOES)(getProcAddr("glTexEnvxvOES"))
	gpTexFilterFuncSGIS = (C.GPTEXFILTERFUNCSGIS)(getProcAddr("glTexFilterFuncSGIS"))
	gpTexGend = (C.GPTEXGEND)(getProcAddr("glTexGend"))
	if gpTexGend == nil {
		return errors.New("glTexGend")
	}
	gpTexGendv = (C.GPTEXGENDV)(getProcAddr("glTexGendv"))
	if gpTexGendv == nil {
		return errors.New("glTexGendv")
	}
	gpTexGenf = (C.GPTEXGENF)(getProcAddr("glTexGenf"))
	if gpTexGenf == nil {
		return errors.New("glTexGenf")
	}
	gpTexGenfv = (C.GPTEXGENFV)(getProcAddr("glTexGenfv"))
	if gpTexGenfv == nil {
		return errors.New("glTexGenfv")
	}
	gpTexGeni = (C.GPTEXGENI)(getProcAddr("glTexGeni"))
	if gpTexGeni == nil {
		return errors.New("glTexGeni")
	}
	gpTexGeniv = (C.GPTEXGENIV)(getProcAddr("glTexGeniv"))
	if gpTexGeniv == nil {
		return errors.New("glTexGeniv")
	}
	gpTexGenxOES = (C.GPTEXGENXOES)(getProcAddr("glTexGenxOES"))
	gpTexGenxvOES = (C.GPTEXGENXVOES)(getProcAddr("glTexGenxvOES"))
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
	gpTexImage2DMultisampleCoverageNV = (C.GPTEXIMAGE2DMULTISAMPLECOVERAGENV)(getProcAddr("glTexImage2DMultisampleCoverageNV"))
	gpTexImage3D = (C.GPTEXIMAGE3D)(getProcAddr("glTexImage3D"))
	if gpTexImage3D == nil {
		return errors.New("glTexImage3D")
	}
	gpTexImage3DEXT = (C.GPTEXIMAGE3DEXT)(getProcAddr("glTexImage3DEXT"))
	gpTexImage3DMultisample = (C.GPTEXIMAGE3DMULTISAMPLE)(getProcAddr("glTexImage3DMultisample"))
	if gpTexImage3DMultisample == nil {
		return errors.New("glTexImage3DMultisample")
	}
	gpTexImage3DMultisampleCoverageNV = (C.GPTEXIMAGE3DMULTISAMPLECOVERAGENV)(getProcAddr("glTexImage3DMultisampleCoverageNV"))
	gpTexImage4DSGIS = (C.GPTEXIMAGE4DSGIS)(getProcAddr("glTexImage4DSGIS"))
	gpTexPageCommitmentARB = (C.GPTEXPAGECOMMITMENTARB)(getProcAddr("glTexPageCommitmentARB"))
	gpTexParameterIiv = (C.GPTEXPARAMETERIIV)(getProcAddr("glTexParameterIiv"))
	if gpTexParameterIiv == nil {
		return errors.New("glTexParameterIiv")
	}
	gpTexParameterIivEXT = (C.GPTEXPARAMETERIIVEXT)(getProcAddr("glTexParameterIivEXT"))
	gpTexParameterIuiv = (C.GPTEXPARAMETERIUIV)(getProcAddr("glTexParameterIuiv"))
	if gpTexParameterIuiv == nil {
		return errors.New("glTexParameterIuiv")
	}
	gpTexParameterIuivEXT = (C.GPTEXPARAMETERIUIVEXT)(getProcAddr("glTexParameterIuivEXT"))
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
	gpTexParameterxOES = (C.GPTEXPARAMETERXOES)(getProcAddr("glTexParameterxOES"))
	gpTexParameterxvOES = (C.GPTEXPARAMETERXVOES)(getProcAddr("glTexParameterxvOES"))
	gpTexRenderbufferNV = (C.GPTEXRENDERBUFFERNV)(getProcAddr("glTexRenderbufferNV"))
	gpTexStorage1D = (C.GPTEXSTORAGE1D)(getProcAddr("glTexStorage1D"))
	if gpTexStorage1D == nil {
		return errors.New("glTexStorage1D")
	}
	gpTexStorage2D = (C.GPTEXSTORAGE2D)(getProcAddr("glTexStorage2D"))
	if gpTexStorage2D == nil {
		return errors.New("glTexStorage2D")
	}
	gpTexStorage2DMultisample = (C.GPTEXSTORAGE2DMULTISAMPLE)(getProcAddr("glTexStorage2DMultisample"))
	if gpTexStorage2DMultisample == nil {
		return errors.New("glTexStorage2DMultisample")
	}
	gpTexStorage3D = (C.GPTEXSTORAGE3D)(getProcAddr("glTexStorage3D"))
	if gpTexStorage3D == nil {
		return errors.New("glTexStorage3D")
	}
	gpTexStorage3DMultisample = (C.GPTEXSTORAGE3DMULTISAMPLE)(getProcAddr("glTexStorage3DMultisample"))
	if gpTexStorage3DMultisample == nil {
		return errors.New("glTexStorage3DMultisample")
	}
	gpTexStorageMem1DEXT = (C.GPTEXSTORAGEMEM1DEXT)(getProcAddr("glTexStorageMem1DEXT"))
	gpTexStorageMem2DEXT = (C.GPTEXSTORAGEMEM2DEXT)(getProcAddr("glTexStorageMem2DEXT"))
	gpTexStorageMem2DMultisampleEXT = (C.GPTEXSTORAGEMEM2DMULTISAMPLEEXT)(getProcAddr("glTexStorageMem2DMultisampleEXT"))
	gpTexStorageMem3DEXT = (C.GPTEXSTORAGEMEM3DEXT)(getProcAddr("glTexStorageMem3DEXT"))
	gpTexStorageMem3DMultisampleEXT = (C.GPTEXSTORAGEMEM3DMULTISAMPLEEXT)(getProcAddr("glTexStorageMem3DMultisampleEXT"))
	gpTexStorageSparseAMD = (C.GPTEXSTORAGESPARSEAMD)(getProcAddr("glTexStorageSparseAMD"))
	gpTexSubImage1D = (C.GPTEXSUBIMAGE1D)(getProcAddr("glTexSubImage1D"))
	if gpTexSubImage1D == nil {
		return errors.New("glTexSubImage1D")
	}
	gpTexSubImage1DEXT = (C.GPTEXSUBIMAGE1DEXT)(getProcAddr("glTexSubImage1DEXT"))
	gpTexSubImage2D = (C.GPTEXSUBIMAGE2D)(getProcAddr("glTexSubImage2D"))
	if gpTexSubImage2D == nil {
		return errors.New("glTexSubImage2D")
	}
	gpTexSubImage2DEXT = (C.GPTEXSUBIMAGE2DEXT)(getProcAddr("glTexSubImage2DEXT"))
	gpTexSubImage3D = (C.GPTEXSUBIMAGE3D)(getProcAddr("glTexSubImage3D"))
	if gpTexSubImage3D == nil {
		return errors.New("glTexSubImage3D")
	}
	gpTexSubImage3DEXT = (C.GPTEXSUBIMAGE3DEXT)(getProcAddr("glTexSubImage3DEXT"))
	gpTexSubImage4DSGIS = (C.GPTEXSUBIMAGE4DSGIS)(getProcAddr("glTexSubImage4DSGIS"))
	gpTextureBarrier = (C.GPTEXTUREBARRIER)(getProcAddr("glTextureBarrier"))
	if gpTextureBarrier == nil {
		return errors.New("glTextureBarrier")
	}
	gpTextureBarrierNV = (C.GPTEXTUREBARRIERNV)(getProcAddr("glTextureBarrierNV"))
	gpTextureBuffer = (C.GPTEXTUREBUFFER)(getProcAddr("glTextureBuffer"))
	if gpTextureBuffer == nil {
		return errors.New("glTextureBuffer")
	}
	gpTextureBufferEXT = (C.GPTEXTUREBUFFEREXT)(getProcAddr("glTextureBufferEXT"))
	gpTextureBufferRange = (C.GPTEXTUREBUFFERRANGE)(getProcAddr("glTextureBufferRange"))
	if gpTextureBufferRange == nil {
		return errors.New("glTextureBufferRange")
	}
	gpTextureBufferRangeEXT = (C.GPTEXTUREBUFFERRANGEEXT)(getProcAddr("glTextureBufferRangeEXT"))
	gpTextureColorMaskSGIS = (C.GPTEXTURECOLORMASKSGIS)(getProcAddr("glTextureColorMaskSGIS"))
	gpTextureImage1DEXT = (C.GPTEXTUREIMAGE1DEXT)(getProcAddr("glTextureImage1DEXT"))
	gpTextureImage2DEXT = (C.GPTEXTUREIMAGE2DEXT)(getProcAddr("glTextureImage2DEXT"))
	gpTextureImage2DMultisampleCoverageNV = (C.GPTEXTUREIMAGE2DMULTISAMPLECOVERAGENV)(getProcAddr("glTextureImage2DMultisampleCoverageNV"))
	gpTextureImage2DMultisampleNV = (C.GPTEXTUREIMAGE2DMULTISAMPLENV)(getProcAddr("glTextureImage2DMultisampleNV"))
	gpTextureImage3DEXT = (C.GPTEXTUREIMAGE3DEXT)(getProcAddr("glTextureImage3DEXT"))
	gpTextureImage3DMultisampleCoverageNV = (C.GPTEXTUREIMAGE3DMULTISAMPLECOVERAGENV)(getProcAddr("glTextureImage3DMultisampleCoverageNV"))
	gpTextureImage3DMultisampleNV = (C.GPTEXTUREIMAGE3DMULTISAMPLENV)(getProcAddr("glTextureImage3DMultisampleNV"))
	gpTextureLightEXT = (C.GPTEXTURELIGHTEXT)(getProcAddr("glTextureLightEXT"))
	gpTextureMaterialEXT = (C.GPTEXTUREMATERIALEXT)(getProcAddr("glTextureMaterialEXT"))
	gpTextureNormalEXT = (C.GPTEXTURENORMALEXT)(getProcAddr("glTextureNormalEXT"))
	gpTexturePageCommitmentEXT = (C.GPTEXTUREPAGECOMMITMENTEXT)(getProcAddr("glTexturePageCommitmentEXT"))
	gpTextureParameterIiv = (C.GPTEXTUREPARAMETERIIV)(getProcAddr("glTextureParameterIiv"))
	if gpTextureParameterIiv == nil {
		return errors.New("glTextureParameterIiv")
	}
	gpTextureParameterIivEXT = (C.GPTEXTUREPARAMETERIIVEXT)(getProcAddr("glTextureParameterIivEXT"))
	gpTextureParameterIuiv = (C.GPTEXTUREPARAMETERIUIV)(getProcAddr("glTextureParameterIuiv"))
	if gpTextureParameterIuiv == nil {
		return errors.New("glTextureParameterIuiv")
	}
	gpTextureParameterIuivEXT = (C.GPTEXTUREPARAMETERIUIVEXT)(getProcAddr("glTextureParameterIuivEXT"))
	gpTextureParameterf = (C.GPTEXTUREPARAMETERF)(getProcAddr("glTextureParameterf"))
	if gpTextureParameterf == nil {
		return errors.New("glTextureParameterf")
	}
	gpTextureParameterfEXT = (C.GPTEXTUREPARAMETERFEXT)(getProcAddr("glTextureParameterfEXT"))
	gpTextureParameterfv = (C.GPTEXTUREPARAMETERFV)(getProcAddr("glTextureParameterfv"))
	if gpTextureParameterfv == nil {
		return errors.New("glTextureParameterfv")
	}
	gpTextureParameterfvEXT = (C.GPTEXTUREPARAMETERFVEXT)(getProcAddr("glTextureParameterfvEXT"))
	gpTextureParameteri = (C.GPTEXTUREPARAMETERI)(getProcAddr("glTextureParameteri"))
	if gpTextureParameteri == nil {
		return errors.New("glTextureParameteri")
	}
	gpTextureParameteriEXT = (C.GPTEXTUREPARAMETERIEXT)(getProcAddr("glTextureParameteriEXT"))
	gpTextureParameteriv = (C.GPTEXTUREPARAMETERIV)(getProcAddr("glTextureParameteriv"))
	if gpTextureParameteriv == nil {
		return errors.New("glTextureParameteriv")
	}
	gpTextureParameterivEXT = (C.GPTEXTUREPARAMETERIVEXT)(getProcAddr("glTextureParameterivEXT"))
	gpTextureRangeAPPLE = (C.GPTEXTURERANGEAPPLE)(getProcAddr("glTextureRangeAPPLE"))
	gpTextureRenderbufferEXT = (C.GPTEXTURERENDERBUFFEREXT)(getProcAddr("glTextureRenderbufferEXT"))
	gpTextureStorage1D = (C.GPTEXTURESTORAGE1D)(getProcAddr("glTextureStorage1D"))
	if gpTextureStorage1D == nil {
		return errors.New("glTextureStorage1D")
	}
	gpTextureStorage1DEXT = (C.GPTEXTURESTORAGE1DEXT)(getProcAddr("glTextureStorage1DEXT"))
	gpTextureStorage2D = (C.GPTEXTURESTORAGE2D)(getProcAddr("glTextureStorage2D"))
	if gpTextureStorage2D == nil {
		return errors.New("glTextureStorage2D")
	}
	gpTextureStorage2DEXT = (C.GPTEXTURESTORAGE2DEXT)(getProcAddr("glTextureStorage2DEXT"))
	gpTextureStorage2DMultisample = (C.GPTEXTURESTORAGE2DMULTISAMPLE)(getProcAddr("glTextureStorage2DMultisample"))
	if gpTextureStorage2DMultisample == nil {
		return errors.New("glTextureStorage2DMultisample")
	}
	gpTextureStorage2DMultisampleEXT = (C.GPTEXTURESTORAGE2DMULTISAMPLEEXT)(getProcAddr("glTextureStorage2DMultisampleEXT"))
	gpTextureStorage3D = (C.GPTEXTURESTORAGE3D)(getProcAddr("glTextureStorage3D"))
	if gpTextureStorage3D == nil {
		return errors.New("glTextureStorage3D")
	}
	gpTextureStorage3DEXT = (C.GPTEXTURESTORAGE3DEXT)(getProcAddr("glTextureStorage3DEXT"))
	gpTextureStorage3DMultisample = (C.GPTEXTURESTORAGE3DMULTISAMPLE)(getProcAddr("glTextureStorage3DMultisample"))
	if gpTextureStorage3DMultisample == nil {
		return errors.New("glTextureStorage3DMultisample")
	}
	gpTextureStorage3DMultisampleEXT = (C.GPTEXTURESTORAGE3DMULTISAMPLEEXT)(getProcAddr("glTextureStorage3DMultisampleEXT"))
	gpTextureStorageMem1DEXT = (C.GPTEXTURESTORAGEMEM1DEXT)(getProcAddr("glTextureStorageMem1DEXT"))
	gpTextureStorageMem2DEXT = (C.GPTEXTURESTORAGEMEM2DEXT)(getProcAddr("glTextureStorageMem2DEXT"))
	gpTextureStorageMem2DMultisampleEXT = (C.GPTEXTURESTORAGEMEM2DMULTISAMPLEEXT)(getProcAddr("glTextureStorageMem2DMultisampleEXT"))
	gpTextureStorageMem3DEXT = (C.GPTEXTURESTORAGEMEM3DEXT)(getProcAddr("glTextureStorageMem3DEXT"))
	gpTextureStorageMem3DMultisampleEXT = (C.GPTEXTURESTORAGEMEM3DMULTISAMPLEEXT)(getProcAddr("glTextureStorageMem3DMultisampleEXT"))
	gpTextureStorageSparseAMD = (C.GPTEXTURESTORAGESPARSEAMD)(getProcAddr("glTextureStorageSparseAMD"))
	gpTextureSubImage1D = (C.GPTEXTURESUBIMAGE1D)(getProcAddr("glTextureSubImage1D"))
	if gpTextureSubImage1D == nil {
		return errors.New("glTextureSubImage1D")
	}
	gpTextureSubImage1DEXT = (C.GPTEXTURESUBIMAGE1DEXT)(getProcAddr("glTextureSubImage1DEXT"))
	gpTextureSubImage2D = (C.GPTEXTURESUBIMAGE2D)(getProcAddr("glTextureSubImage2D"))
	if gpTextureSubImage2D == nil {
		return errors.New("glTextureSubImage2D")
	}
	gpTextureSubImage2DEXT = (C.GPTEXTURESUBIMAGE2DEXT)(getProcAddr("glTextureSubImage2DEXT"))
	gpTextureSubImage3D = (C.GPTEXTURESUBIMAGE3D)(getProcAddr("glTextureSubImage3D"))
	if gpTextureSubImage3D == nil {
		return errors.New("glTextureSubImage3D")
	}
	gpTextureSubImage3DEXT = (C.GPTEXTURESUBIMAGE3DEXT)(getProcAddr("glTextureSubImage3DEXT"))
	gpTextureView = (C.GPTEXTUREVIEW)(getProcAddr("glTextureView"))
	if gpTextureView == nil {
		return errors.New("glTextureView")
	}
	gpTrackMatrixNV = (C.GPTRACKMATRIXNV)(getProcAddr("glTrackMatrixNV"))
	gpTransformFeedbackAttribsNV = (C.GPTRANSFORMFEEDBACKATTRIBSNV)(getProcAddr("glTransformFeedbackAttribsNV"))
	gpTransformFeedbackBufferBase = (C.GPTRANSFORMFEEDBACKBUFFERBASE)(getProcAddr("glTransformFeedbackBufferBase"))
	if gpTransformFeedbackBufferBase == nil {
		return errors.New("glTransformFeedbackBufferBase")
	}
	gpTransformFeedbackBufferRange = (C.GPTRANSFORMFEEDBACKBUFFERRANGE)(getProcAddr("glTransformFeedbackBufferRange"))
	if gpTransformFeedbackBufferRange == nil {
		return errors.New("glTransformFeedbackBufferRange")
	}
	gpTransformFeedbackStreamAttribsNV = (C.GPTRANSFORMFEEDBACKSTREAMATTRIBSNV)(getProcAddr("glTransformFeedbackStreamAttribsNV"))
	gpTransformFeedbackVaryings = (C.GPTRANSFORMFEEDBACKVARYINGS)(getProcAddr("glTransformFeedbackVaryings"))
	if gpTransformFeedbackVaryings == nil {
		return errors.New("glTransformFeedbackVaryings")
	}
	gpTransformFeedbackVaryingsEXT = (C.GPTRANSFORMFEEDBACKVARYINGSEXT)(getProcAddr("glTransformFeedbackVaryingsEXT"))
	gpTransformFeedbackVaryingsNV = (C.GPTRANSFORMFEEDBACKVARYINGSNV)(getProcAddr("glTransformFeedbackVaryingsNV"))
	gpTransformPathNV = (C.GPTRANSFORMPATHNV)(getProcAddr("glTransformPathNV"))
	gpTranslated = (C.GPTRANSLATED)(getProcAddr("glTranslated"))
	if gpTranslated == nil {
		return errors.New("glTranslated")
	}
	gpTranslatef = (C.GPTRANSLATEF)(getProcAddr("glTranslatef"))
	if gpTranslatef == nil {
		return errors.New("glTranslatef")
	}
	gpTranslatexOES = (C.GPTRANSLATEXOES)(getProcAddr("glTranslatexOES"))
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
	gpUniform1fARB = (C.GPUNIFORM1FARB)(getProcAddr("glUniform1fARB"))
	gpUniform1fv = (C.GPUNIFORM1FV)(getProcAddr("glUniform1fv"))
	if gpUniform1fv == nil {
		return errors.New("glUniform1fv")
	}
	gpUniform1fvARB = (C.GPUNIFORM1FVARB)(getProcAddr("glUniform1fvARB"))
	gpUniform1i = (C.GPUNIFORM1I)(getProcAddr("glUniform1i"))
	if gpUniform1i == nil {
		return errors.New("glUniform1i")
	}
	gpUniform1i64ARB = (C.GPUNIFORM1I64ARB)(getProcAddr("glUniform1i64ARB"))
	gpUniform1i64NV = (C.GPUNIFORM1I64NV)(getProcAddr("glUniform1i64NV"))
	gpUniform1i64vARB = (C.GPUNIFORM1I64VARB)(getProcAddr("glUniform1i64vARB"))
	gpUniform1i64vNV = (C.GPUNIFORM1I64VNV)(getProcAddr("glUniform1i64vNV"))
	gpUniform1iARB = (C.GPUNIFORM1IARB)(getProcAddr("glUniform1iARB"))
	gpUniform1iv = (C.GPUNIFORM1IV)(getProcAddr("glUniform1iv"))
	if gpUniform1iv == nil {
		return errors.New("glUniform1iv")
	}
	gpUniform1ivARB = (C.GPUNIFORM1IVARB)(getProcAddr("glUniform1ivARB"))
	gpUniform1ui = (C.GPUNIFORM1UI)(getProcAddr("glUniform1ui"))
	if gpUniform1ui == nil {
		return errors.New("glUniform1ui")
	}
	gpUniform1ui64ARB = (C.GPUNIFORM1UI64ARB)(getProcAddr("glUniform1ui64ARB"))
	gpUniform1ui64NV = (C.GPUNIFORM1UI64NV)(getProcAddr("glUniform1ui64NV"))
	gpUniform1ui64vARB = (C.GPUNIFORM1UI64VARB)(getProcAddr("glUniform1ui64vARB"))
	gpUniform1ui64vNV = (C.GPUNIFORM1UI64VNV)(getProcAddr("glUniform1ui64vNV"))
	gpUniform1uiEXT = (C.GPUNIFORM1UIEXT)(getProcAddr("glUniform1uiEXT"))
	gpUniform1uiv = (C.GPUNIFORM1UIV)(getProcAddr("glUniform1uiv"))
	if gpUniform1uiv == nil {
		return errors.New("glUniform1uiv")
	}
	gpUniform1uivEXT = (C.GPUNIFORM1UIVEXT)(getProcAddr("glUniform1uivEXT"))
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
	gpUniform2fARB = (C.GPUNIFORM2FARB)(getProcAddr("glUniform2fARB"))
	gpUniform2fv = (C.GPUNIFORM2FV)(getProcAddr("glUniform2fv"))
	if gpUniform2fv == nil {
		return errors.New("glUniform2fv")
	}
	gpUniform2fvARB = (C.GPUNIFORM2FVARB)(getProcAddr("glUniform2fvARB"))
	gpUniform2i = (C.GPUNIFORM2I)(getProcAddr("glUniform2i"))
	if gpUniform2i == nil {
		return errors.New("glUniform2i")
	}
	gpUniform2i64ARB = (C.GPUNIFORM2I64ARB)(getProcAddr("glUniform2i64ARB"))
	gpUniform2i64NV = (C.GPUNIFORM2I64NV)(getProcAddr("glUniform2i64NV"))
	gpUniform2i64vARB = (C.GPUNIFORM2I64VARB)(getProcAddr("glUniform2i64vARB"))
	gpUniform2i64vNV = (C.GPUNIFORM2I64VNV)(getProcAddr("glUniform2i64vNV"))
	gpUniform2iARB = (C.GPUNIFORM2IARB)(getProcAddr("glUniform2iARB"))
	gpUniform2iv = (C.GPUNIFORM2IV)(getProcAddr("glUniform2iv"))
	if gpUniform2iv == nil {
		return errors.New("glUniform2iv")
	}
	gpUniform2ivARB = (C.GPUNIFORM2IVARB)(getProcAddr("glUniform2ivARB"))
	gpUniform2ui = (C.GPUNIFORM2UI)(getProcAddr("glUniform2ui"))
	if gpUniform2ui == nil {
		return errors.New("glUniform2ui")
	}
	gpUniform2ui64ARB = (C.GPUNIFORM2UI64ARB)(getProcAddr("glUniform2ui64ARB"))
	gpUniform2ui64NV = (C.GPUNIFORM2UI64NV)(getProcAddr("glUniform2ui64NV"))
	gpUniform2ui64vARB = (C.GPUNIFORM2UI64VARB)(getProcAddr("glUniform2ui64vARB"))
	gpUniform2ui64vNV = (C.GPUNIFORM2UI64VNV)(getProcAddr("glUniform2ui64vNV"))
	gpUniform2uiEXT = (C.GPUNIFORM2UIEXT)(getProcAddr("glUniform2uiEXT"))
	gpUniform2uiv = (C.GPUNIFORM2UIV)(getProcAddr("glUniform2uiv"))
	if gpUniform2uiv == nil {
		return errors.New("glUniform2uiv")
	}
	gpUniform2uivEXT = (C.GPUNIFORM2UIVEXT)(getProcAddr("glUniform2uivEXT"))
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
	gpUniform3fARB = (C.GPUNIFORM3FARB)(getProcAddr("glUniform3fARB"))
	gpUniform3fv = (C.GPUNIFORM3FV)(getProcAddr("glUniform3fv"))
	if gpUniform3fv == nil {
		return errors.New("glUniform3fv")
	}
	gpUniform3fvARB = (C.GPUNIFORM3FVARB)(getProcAddr("glUniform3fvARB"))
	gpUniform3i = (C.GPUNIFORM3I)(getProcAddr("glUniform3i"))
	if gpUniform3i == nil {
		return errors.New("glUniform3i")
	}
	gpUniform3i64ARB = (C.GPUNIFORM3I64ARB)(getProcAddr("glUniform3i64ARB"))
	gpUniform3i64NV = (C.GPUNIFORM3I64NV)(getProcAddr("glUniform3i64NV"))
	gpUniform3i64vARB = (C.GPUNIFORM3I64VARB)(getProcAddr("glUniform3i64vARB"))
	gpUniform3i64vNV = (C.GPUNIFORM3I64VNV)(getProcAddr("glUniform3i64vNV"))
	gpUniform3iARB = (C.GPUNIFORM3IARB)(getProcAddr("glUniform3iARB"))
	gpUniform3iv = (C.GPUNIFORM3IV)(getProcAddr("glUniform3iv"))
	if gpUniform3iv == nil {
		return errors.New("glUniform3iv")
	}
	gpUniform3ivARB = (C.GPUNIFORM3IVARB)(getProcAddr("glUniform3ivARB"))
	gpUniform3ui = (C.GPUNIFORM3UI)(getProcAddr("glUniform3ui"))
	if gpUniform3ui == nil {
		return errors.New("glUniform3ui")
	}
	gpUniform3ui64ARB = (C.GPUNIFORM3UI64ARB)(getProcAddr("glUniform3ui64ARB"))
	gpUniform3ui64NV = (C.GPUNIFORM3UI64NV)(getProcAddr("glUniform3ui64NV"))
	gpUniform3ui64vARB = (C.GPUNIFORM3UI64VARB)(getProcAddr("glUniform3ui64vARB"))
	gpUniform3ui64vNV = (C.GPUNIFORM3UI64VNV)(getProcAddr("glUniform3ui64vNV"))
	gpUniform3uiEXT = (C.GPUNIFORM3UIEXT)(getProcAddr("glUniform3uiEXT"))
	gpUniform3uiv = (C.GPUNIFORM3UIV)(getProcAddr("glUniform3uiv"))
	if gpUniform3uiv == nil {
		return errors.New("glUniform3uiv")
	}
	gpUniform3uivEXT = (C.GPUNIFORM3UIVEXT)(getProcAddr("glUniform3uivEXT"))
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
	gpUniform4fARB = (C.GPUNIFORM4FARB)(getProcAddr("glUniform4fARB"))
	gpUniform4fv = (C.GPUNIFORM4FV)(getProcAddr("glUniform4fv"))
	if gpUniform4fv == nil {
		return errors.New("glUniform4fv")
	}
	gpUniform4fvARB = (C.GPUNIFORM4FVARB)(getProcAddr("glUniform4fvARB"))
	gpUniform4i = (C.GPUNIFORM4I)(getProcAddr("glUniform4i"))
	if gpUniform4i == nil {
		return errors.New("glUniform4i")
	}
	gpUniform4i64ARB = (C.GPUNIFORM4I64ARB)(getProcAddr("glUniform4i64ARB"))
	gpUniform4i64NV = (C.GPUNIFORM4I64NV)(getProcAddr("glUniform4i64NV"))
	gpUniform4i64vARB = (C.GPUNIFORM4I64VARB)(getProcAddr("glUniform4i64vARB"))
	gpUniform4i64vNV = (C.GPUNIFORM4I64VNV)(getProcAddr("glUniform4i64vNV"))
	gpUniform4iARB = (C.GPUNIFORM4IARB)(getProcAddr("glUniform4iARB"))
	gpUniform4iv = (C.GPUNIFORM4IV)(getProcAddr("glUniform4iv"))
	if gpUniform4iv == nil {
		return errors.New("glUniform4iv")
	}
	gpUniform4ivARB = (C.GPUNIFORM4IVARB)(getProcAddr("glUniform4ivARB"))
	gpUniform4ui = (C.GPUNIFORM4UI)(getProcAddr("glUniform4ui"))
	if gpUniform4ui == nil {
		return errors.New("glUniform4ui")
	}
	gpUniform4ui64ARB = (C.GPUNIFORM4UI64ARB)(getProcAddr("glUniform4ui64ARB"))
	gpUniform4ui64NV = (C.GPUNIFORM4UI64NV)(getProcAddr("glUniform4ui64NV"))
	gpUniform4ui64vARB = (C.GPUNIFORM4UI64VARB)(getProcAddr("glUniform4ui64vARB"))
	gpUniform4ui64vNV = (C.GPUNIFORM4UI64VNV)(getProcAddr("glUniform4ui64vNV"))
	gpUniform4uiEXT = (C.GPUNIFORM4UIEXT)(getProcAddr("glUniform4uiEXT"))
	gpUniform4uiv = (C.GPUNIFORM4UIV)(getProcAddr("glUniform4uiv"))
	if gpUniform4uiv == nil {
		return errors.New("glUniform4uiv")
	}
	gpUniform4uivEXT = (C.GPUNIFORM4UIVEXT)(getProcAddr("glUniform4uivEXT"))
	gpUniformBlockBinding = (C.GPUNIFORMBLOCKBINDING)(getProcAddr("glUniformBlockBinding"))
	if gpUniformBlockBinding == nil {
		return errors.New("glUniformBlockBinding")
	}
	gpUniformBufferEXT = (C.GPUNIFORMBUFFEREXT)(getProcAddr("glUniformBufferEXT"))
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
	gpUniformMatrix2fvARB = (C.GPUNIFORMMATRIX2FVARB)(getProcAddr("glUniformMatrix2fvARB"))
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
	gpUniformMatrix3fvARB = (C.GPUNIFORMMATRIX3FVARB)(getProcAddr("glUniformMatrix3fvARB"))
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
	gpUniformMatrix4fvARB = (C.GPUNIFORMMATRIX4FVARB)(getProcAddr("glUniformMatrix4fvARB"))
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
	gpUnlockArraysEXT = (C.GPUNLOCKARRAYSEXT)(getProcAddr("glUnlockArraysEXT"))
	gpUnmapBuffer = (C.GPUNMAPBUFFER)(getProcAddr("glUnmapBuffer"))
	if gpUnmapBuffer == nil {
		return errors.New("glUnmapBuffer")
	}
	gpUnmapBufferARB = (C.GPUNMAPBUFFERARB)(getProcAddr("glUnmapBufferARB"))
	gpUnmapNamedBuffer = (C.GPUNMAPNAMEDBUFFER)(getProcAddr("glUnmapNamedBuffer"))
	if gpUnmapNamedBuffer == nil {
		return errors.New("glUnmapNamedBuffer")
	}
	gpUnmapNamedBufferEXT = (C.GPUNMAPNAMEDBUFFEREXT)(getProcAddr("glUnmapNamedBufferEXT"))
	gpUnmapObjectBufferATI = (C.GPUNMAPOBJECTBUFFERATI)(getProcAddr("glUnmapObjectBufferATI"))
	gpUnmapTexture2DINTEL = (C.GPUNMAPTEXTURE2DINTEL)(getProcAddr("glUnmapTexture2DINTEL"))
	gpUpdateObjectBufferATI = (C.GPUPDATEOBJECTBUFFERATI)(getProcAddr("glUpdateObjectBufferATI"))
	gpUseProgram = (C.GPUSEPROGRAM)(getProcAddr("glUseProgram"))
	if gpUseProgram == nil {
		return errors.New("glUseProgram")
	}
	gpUseProgramObjectARB = (C.GPUSEPROGRAMOBJECTARB)(getProcAddr("glUseProgramObjectARB"))
	gpUseProgramStages = (C.GPUSEPROGRAMSTAGES)(getProcAddr("glUseProgramStages"))
	if gpUseProgramStages == nil {
		return errors.New("glUseProgramStages")
	}
	gpUseProgramStagesEXT = (C.GPUSEPROGRAMSTAGESEXT)(getProcAddr("glUseProgramStagesEXT"))
	gpUseShaderProgramEXT = (C.GPUSESHADERPROGRAMEXT)(getProcAddr("glUseShaderProgramEXT"))
	gpVDPAUFiniNV = (C.GPVDPAUFININV)(getProcAddr("glVDPAUFiniNV"))
	gpVDPAUGetSurfaceivNV = (C.GPVDPAUGETSURFACEIVNV)(getProcAddr("glVDPAUGetSurfaceivNV"))
	gpVDPAUInitNV = (C.GPVDPAUINITNV)(getProcAddr("glVDPAUInitNV"))
	gpVDPAUIsSurfaceNV = (C.GPVDPAUISSURFACENV)(getProcAddr("glVDPAUIsSurfaceNV"))
	gpVDPAUMapSurfacesNV = (C.GPVDPAUMAPSURFACESNV)(getProcAddr("glVDPAUMapSurfacesNV"))
	gpVDPAURegisterOutputSurfaceNV = (C.GPVDPAUREGISTEROUTPUTSURFACENV)(getProcAddr("glVDPAURegisterOutputSurfaceNV"))
	gpVDPAURegisterVideoSurfaceNV = (C.GPVDPAUREGISTERVIDEOSURFACENV)(getProcAddr("glVDPAURegisterVideoSurfaceNV"))
	gpVDPAUSurfaceAccessNV = (C.GPVDPAUSURFACEACCESSNV)(getProcAddr("glVDPAUSurfaceAccessNV"))
	gpVDPAUUnmapSurfacesNV = (C.GPVDPAUUNMAPSURFACESNV)(getProcAddr("glVDPAUUnmapSurfacesNV"))
	gpVDPAUUnregisterSurfaceNV = (C.GPVDPAUUNREGISTERSURFACENV)(getProcAddr("glVDPAUUnregisterSurfaceNV"))
	gpValidateProgram = (C.GPVALIDATEPROGRAM)(getProcAddr("glValidateProgram"))
	if gpValidateProgram == nil {
		return errors.New("glValidateProgram")
	}
	gpValidateProgramARB = (C.GPVALIDATEPROGRAMARB)(getProcAddr("glValidateProgramARB"))
	gpValidateProgramPipeline = (C.GPVALIDATEPROGRAMPIPELINE)(getProcAddr("glValidateProgramPipeline"))
	if gpValidateProgramPipeline == nil {
		return errors.New("glValidateProgramPipeline")
	}
	gpValidateProgramPipelineEXT = (C.GPVALIDATEPROGRAMPIPELINEEXT)(getProcAddr("glValidateProgramPipelineEXT"))
	gpVariantArrayObjectATI = (C.GPVARIANTARRAYOBJECTATI)(getProcAddr("glVariantArrayObjectATI"))
	gpVariantPointerEXT = (C.GPVARIANTPOINTEREXT)(getProcAddr("glVariantPointerEXT"))
	gpVariantbvEXT = (C.GPVARIANTBVEXT)(getProcAddr("glVariantbvEXT"))
	gpVariantdvEXT = (C.GPVARIANTDVEXT)(getProcAddr("glVariantdvEXT"))
	gpVariantfvEXT = (C.GPVARIANTFVEXT)(getProcAddr("glVariantfvEXT"))
	gpVariantivEXT = (C.GPVARIANTIVEXT)(getProcAddr("glVariantivEXT"))
	gpVariantsvEXT = (C.GPVARIANTSVEXT)(getProcAddr("glVariantsvEXT"))
	gpVariantubvEXT = (C.GPVARIANTUBVEXT)(getProcAddr("glVariantubvEXT"))
	gpVariantuivEXT = (C.GPVARIANTUIVEXT)(getProcAddr("glVariantuivEXT"))
	gpVariantusvEXT = (C.GPVARIANTUSVEXT)(getProcAddr("glVariantusvEXT"))
	gpVertex2bOES = (C.GPVERTEX2BOES)(getProcAddr("glVertex2bOES"))
	gpVertex2bvOES = (C.GPVERTEX2BVOES)(getProcAddr("glVertex2bvOES"))
	gpVertex2d = (C.GPVERTEX2D)(getProcAddr("glVertex2d"))
	if gpVertex2d == nil {
		return errors.New("glVertex2d")
	}
	gpVertex2dv = (C.GPVERTEX2DV)(getProcAddr("glVertex2dv"))
	if gpVertex2dv == nil {
		return errors.New("glVertex2dv")
	}
	gpVertex2f = (C.GPVERTEX2F)(getProcAddr("glVertex2f"))
	if gpVertex2f == nil {
		return errors.New("glVertex2f")
	}
	gpVertex2fv = (C.GPVERTEX2FV)(getProcAddr("glVertex2fv"))
	if gpVertex2fv == nil {
		return errors.New("glVertex2fv")
	}
	gpVertex2hNV = (C.GPVERTEX2HNV)(getProcAddr("glVertex2hNV"))
	gpVertex2hvNV = (C.GPVERTEX2HVNV)(getProcAddr("glVertex2hvNV"))
	gpVertex2i = (C.GPVERTEX2I)(getProcAddr("glVertex2i"))
	if gpVertex2i == nil {
		return errors.New("glVertex2i")
	}
	gpVertex2iv = (C.GPVERTEX2IV)(getProcAddr("glVertex2iv"))
	if gpVertex2iv == nil {
		return errors.New("glVertex2iv")
	}
	gpVertex2s = (C.GPVERTEX2S)(getProcAddr("glVertex2s"))
	if gpVertex2s == nil {
		return errors.New("glVertex2s")
	}
	gpVertex2sv = (C.GPVERTEX2SV)(getProcAddr("glVertex2sv"))
	if gpVertex2sv == nil {
		return errors.New("glVertex2sv")
	}
	gpVertex2xOES = (C.GPVERTEX2XOES)(getProcAddr("glVertex2xOES"))
	gpVertex2xvOES = (C.GPVERTEX2XVOES)(getProcAddr("glVertex2xvOES"))
	gpVertex3bOES = (C.GPVERTEX3BOES)(getProcAddr("glVertex3bOES"))
	gpVertex3bvOES = (C.GPVERTEX3BVOES)(getProcAddr("glVertex3bvOES"))
	gpVertex3d = (C.GPVERTEX3D)(getProcAddr("glVertex3d"))
	if gpVertex3d == nil {
		return errors.New("glVertex3d")
	}
	gpVertex3dv = (C.GPVERTEX3DV)(getProcAddr("glVertex3dv"))
	if gpVertex3dv == nil {
		return errors.New("glVertex3dv")
	}
	gpVertex3f = (C.GPVERTEX3F)(getProcAddr("glVertex3f"))
	if gpVertex3f == nil {
		return errors.New("glVertex3f")
	}
	gpVertex3fv = (C.GPVERTEX3FV)(getProcAddr("glVertex3fv"))
	if gpVertex3fv == nil {
		return errors.New("glVertex3fv")
	}
	gpVertex3hNV = (C.GPVERTEX3HNV)(getProcAddr("glVertex3hNV"))
	gpVertex3hvNV = (C.GPVERTEX3HVNV)(getProcAddr("glVertex3hvNV"))
	gpVertex3i = (C.GPVERTEX3I)(getProcAddr("glVertex3i"))
	if gpVertex3i == nil {
		return errors.New("glVertex3i")
	}
	gpVertex3iv = (C.GPVERTEX3IV)(getProcAddr("glVertex3iv"))
	if gpVertex3iv == nil {
		return errors.New("glVertex3iv")
	}
	gpVertex3s = (C.GPVERTEX3S)(getProcAddr("glVertex3s"))
	if gpVertex3s == nil {
		return errors.New("glVertex3s")
	}
	gpVertex3sv = (C.GPVERTEX3SV)(getProcAddr("glVertex3sv"))
	if gpVertex3sv == nil {
		return errors.New("glVertex3sv")
	}
	gpVertex3xOES = (C.GPVERTEX3XOES)(getProcAddr("glVertex3xOES"))
	gpVertex3xvOES = (C.GPVERTEX3XVOES)(getProcAddr("glVertex3xvOES"))
	gpVertex4bOES = (C.GPVERTEX4BOES)(getProcAddr("glVertex4bOES"))
	gpVertex4bvOES = (C.GPVERTEX4BVOES)(getProcAddr("glVertex4bvOES"))
	gpVertex4d = (C.GPVERTEX4D)(getProcAddr("glVertex4d"))
	if gpVertex4d == nil {
		return errors.New("glVertex4d")
	}
	gpVertex4dv = (C.GPVERTEX4DV)(getProcAddr("glVertex4dv"))
	if gpVertex4dv == nil {
		return errors.New("glVertex4dv")
	}
	gpVertex4f = (C.GPVERTEX4F)(getProcAddr("glVertex4f"))
	if gpVertex4f == nil {
		return errors.New("glVertex4f")
	}
	gpVertex4fv = (C.GPVERTEX4FV)(getProcAddr("glVertex4fv"))
	if gpVertex4fv == nil {
		return errors.New("glVertex4fv")
	}
	gpVertex4hNV = (C.GPVERTEX4HNV)(getProcAddr("glVertex4hNV"))
	gpVertex4hvNV = (C.GPVERTEX4HVNV)(getProcAddr("glVertex4hvNV"))
	gpVertex4i = (C.GPVERTEX4I)(getProcAddr("glVertex4i"))
	if gpVertex4i == nil {
		return errors.New("glVertex4i")
	}
	gpVertex4iv = (C.GPVERTEX4IV)(getProcAddr("glVertex4iv"))
	if gpVertex4iv == nil {
		return errors.New("glVertex4iv")
	}
	gpVertex4s = (C.GPVERTEX4S)(getProcAddr("glVertex4s"))
	if gpVertex4s == nil {
		return errors.New("glVertex4s")
	}
	gpVertex4sv = (C.GPVERTEX4SV)(getProcAddr("glVertex4sv"))
	if gpVertex4sv == nil {
		return errors.New("glVertex4sv")
	}
	gpVertex4xOES = (C.GPVERTEX4XOES)(getProcAddr("glVertex4xOES"))
	gpVertex4xvOES = (C.GPVERTEX4XVOES)(getProcAddr("glVertex4xvOES"))
	gpVertexArrayAttribBinding = (C.GPVERTEXARRAYATTRIBBINDING)(getProcAddr("glVertexArrayAttribBinding"))
	if gpVertexArrayAttribBinding == nil {
		return errors.New("glVertexArrayAttribBinding")
	}
	gpVertexArrayAttribFormat = (C.GPVERTEXARRAYATTRIBFORMAT)(getProcAddr("glVertexArrayAttribFormat"))
	if gpVertexArrayAttribFormat == nil {
		return errors.New("glVertexArrayAttribFormat")
	}
	gpVertexArrayAttribIFormat = (C.GPVERTEXARRAYATTRIBIFORMAT)(getProcAddr("glVertexArrayAttribIFormat"))
	if gpVertexArrayAttribIFormat == nil {
		return errors.New("glVertexArrayAttribIFormat")
	}
	gpVertexArrayAttribLFormat = (C.GPVERTEXARRAYATTRIBLFORMAT)(getProcAddr("glVertexArrayAttribLFormat"))
	if gpVertexArrayAttribLFormat == nil {
		return errors.New("glVertexArrayAttribLFormat")
	}
	gpVertexArrayBindVertexBufferEXT = (C.GPVERTEXARRAYBINDVERTEXBUFFEREXT)(getProcAddr("glVertexArrayBindVertexBufferEXT"))
	gpVertexArrayBindingDivisor = (C.GPVERTEXARRAYBINDINGDIVISOR)(getProcAddr("glVertexArrayBindingDivisor"))
	if gpVertexArrayBindingDivisor == nil {
		return errors.New("glVertexArrayBindingDivisor")
	}
	gpVertexArrayColorOffsetEXT = (C.GPVERTEXARRAYCOLOROFFSETEXT)(getProcAddr("glVertexArrayColorOffsetEXT"))
	gpVertexArrayEdgeFlagOffsetEXT = (C.GPVERTEXARRAYEDGEFLAGOFFSETEXT)(getProcAddr("glVertexArrayEdgeFlagOffsetEXT"))
	gpVertexArrayElementBuffer = (C.GPVERTEXARRAYELEMENTBUFFER)(getProcAddr("glVertexArrayElementBuffer"))
	if gpVertexArrayElementBuffer == nil {
		return errors.New("glVertexArrayElementBuffer")
	}
	gpVertexArrayFogCoordOffsetEXT = (C.GPVERTEXARRAYFOGCOORDOFFSETEXT)(getProcAddr("glVertexArrayFogCoordOffsetEXT"))
	gpVertexArrayIndexOffsetEXT = (C.GPVERTEXARRAYINDEXOFFSETEXT)(getProcAddr("glVertexArrayIndexOffsetEXT"))
	gpVertexArrayMultiTexCoordOffsetEXT = (C.GPVERTEXARRAYMULTITEXCOORDOFFSETEXT)(getProcAddr("glVertexArrayMultiTexCoordOffsetEXT"))
	gpVertexArrayNormalOffsetEXT = (C.GPVERTEXARRAYNORMALOFFSETEXT)(getProcAddr("glVertexArrayNormalOffsetEXT"))
	gpVertexArrayParameteriAPPLE = (C.GPVERTEXARRAYPARAMETERIAPPLE)(getProcAddr("glVertexArrayParameteriAPPLE"))
	gpVertexArrayRangeAPPLE = (C.GPVERTEXARRAYRANGEAPPLE)(getProcAddr("glVertexArrayRangeAPPLE"))
	gpVertexArrayRangeNV = (C.GPVERTEXARRAYRANGENV)(getProcAddr("glVertexArrayRangeNV"))
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
	if gpVertexArrayVertexBuffer == nil {
		return errors.New("glVertexArrayVertexBuffer")
	}
	gpVertexArrayVertexBuffers = (C.GPVERTEXARRAYVERTEXBUFFERS)(getProcAddr("glVertexArrayVertexBuffers"))
	if gpVertexArrayVertexBuffers == nil {
		return errors.New("glVertexArrayVertexBuffers")
	}
	gpVertexArrayVertexOffsetEXT = (C.GPVERTEXARRAYVERTEXOFFSETEXT)(getProcAddr("glVertexArrayVertexOffsetEXT"))
	gpVertexAttrib1d = (C.GPVERTEXATTRIB1D)(getProcAddr("glVertexAttrib1d"))
	if gpVertexAttrib1d == nil {
		return errors.New("glVertexAttrib1d")
	}
	gpVertexAttrib1dARB = (C.GPVERTEXATTRIB1DARB)(getProcAddr("glVertexAttrib1dARB"))
	gpVertexAttrib1dNV = (C.GPVERTEXATTRIB1DNV)(getProcAddr("glVertexAttrib1dNV"))
	gpVertexAttrib1dv = (C.GPVERTEXATTRIB1DV)(getProcAddr("glVertexAttrib1dv"))
	if gpVertexAttrib1dv == nil {
		return errors.New("glVertexAttrib1dv")
	}
	gpVertexAttrib1dvARB = (C.GPVERTEXATTRIB1DVARB)(getProcAddr("glVertexAttrib1dvARB"))
	gpVertexAttrib1dvNV = (C.GPVERTEXATTRIB1DVNV)(getProcAddr("glVertexAttrib1dvNV"))
	gpVertexAttrib1f = (C.GPVERTEXATTRIB1F)(getProcAddr("glVertexAttrib1f"))
	if gpVertexAttrib1f == nil {
		return errors.New("glVertexAttrib1f")
	}
	gpVertexAttrib1fARB = (C.GPVERTEXATTRIB1FARB)(getProcAddr("glVertexAttrib1fARB"))
	gpVertexAttrib1fNV = (C.GPVERTEXATTRIB1FNV)(getProcAddr("glVertexAttrib1fNV"))
	gpVertexAttrib1fv = (C.GPVERTEXATTRIB1FV)(getProcAddr("glVertexAttrib1fv"))
	if gpVertexAttrib1fv == nil {
		return errors.New("glVertexAttrib1fv")
	}
	gpVertexAttrib1fvARB = (C.GPVERTEXATTRIB1FVARB)(getProcAddr("glVertexAttrib1fvARB"))
	gpVertexAttrib1fvNV = (C.GPVERTEXATTRIB1FVNV)(getProcAddr("glVertexAttrib1fvNV"))
	gpVertexAttrib1hNV = (C.GPVERTEXATTRIB1HNV)(getProcAddr("glVertexAttrib1hNV"))
	gpVertexAttrib1hvNV = (C.GPVERTEXATTRIB1HVNV)(getProcAddr("glVertexAttrib1hvNV"))
	gpVertexAttrib1s = (C.GPVERTEXATTRIB1S)(getProcAddr("glVertexAttrib1s"))
	if gpVertexAttrib1s == nil {
		return errors.New("glVertexAttrib1s")
	}
	gpVertexAttrib1sARB = (C.GPVERTEXATTRIB1SARB)(getProcAddr("glVertexAttrib1sARB"))
	gpVertexAttrib1sNV = (C.GPVERTEXATTRIB1SNV)(getProcAddr("glVertexAttrib1sNV"))
	gpVertexAttrib1sv = (C.GPVERTEXATTRIB1SV)(getProcAddr("glVertexAttrib1sv"))
	if gpVertexAttrib1sv == nil {
		return errors.New("glVertexAttrib1sv")
	}
	gpVertexAttrib1svARB = (C.GPVERTEXATTRIB1SVARB)(getProcAddr("glVertexAttrib1svARB"))
	gpVertexAttrib1svNV = (C.GPVERTEXATTRIB1SVNV)(getProcAddr("glVertexAttrib1svNV"))
	gpVertexAttrib2d = (C.GPVERTEXATTRIB2D)(getProcAddr("glVertexAttrib2d"))
	if gpVertexAttrib2d == nil {
		return errors.New("glVertexAttrib2d")
	}
	gpVertexAttrib2dARB = (C.GPVERTEXATTRIB2DARB)(getProcAddr("glVertexAttrib2dARB"))
	gpVertexAttrib2dNV = (C.GPVERTEXATTRIB2DNV)(getProcAddr("glVertexAttrib2dNV"))
	gpVertexAttrib2dv = (C.GPVERTEXATTRIB2DV)(getProcAddr("glVertexAttrib2dv"))
	if gpVertexAttrib2dv == nil {
		return errors.New("glVertexAttrib2dv")
	}
	gpVertexAttrib2dvARB = (C.GPVERTEXATTRIB2DVARB)(getProcAddr("glVertexAttrib2dvARB"))
	gpVertexAttrib2dvNV = (C.GPVERTEXATTRIB2DVNV)(getProcAddr("glVertexAttrib2dvNV"))
	gpVertexAttrib2f = (C.GPVERTEXATTRIB2F)(getProcAddr("glVertexAttrib2f"))
	if gpVertexAttrib2f == nil {
		return errors.New("glVertexAttrib2f")
	}
	gpVertexAttrib2fARB = (C.GPVERTEXATTRIB2FARB)(getProcAddr("glVertexAttrib2fARB"))
	gpVertexAttrib2fNV = (C.GPVERTEXATTRIB2FNV)(getProcAddr("glVertexAttrib2fNV"))
	gpVertexAttrib2fv = (C.GPVERTEXATTRIB2FV)(getProcAddr("glVertexAttrib2fv"))
	if gpVertexAttrib2fv == nil {
		return errors.New("glVertexAttrib2fv")
	}
	gpVertexAttrib2fvARB = (C.GPVERTEXATTRIB2FVARB)(getProcAddr("glVertexAttrib2fvARB"))
	gpVertexAttrib2fvNV = (C.GPVERTEXATTRIB2FVNV)(getProcAddr("glVertexAttrib2fvNV"))
	gpVertexAttrib2hNV = (C.GPVERTEXATTRIB2HNV)(getProcAddr("glVertexAttrib2hNV"))
	gpVertexAttrib2hvNV = (C.GPVERTEXATTRIB2HVNV)(getProcAddr("glVertexAttrib2hvNV"))
	gpVertexAttrib2s = (C.GPVERTEXATTRIB2S)(getProcAddr("glVertexAttrib2s"))
	if gpVertexAttrib2s == nil {
		return errors.New("glVertexAttrib2s")
	}
	gpVertexAttrib2sARB = (C.GPVERTEXATTRIB2SARB)(getProcAddr("glVertexAttrib2sARB"))
	gpVertexAttrib2sNV = (C.GPVERTEXATTRIB2SNV)(getProcAddr("glVertexAttrib2sNV"))
	gpVertexAttrib2sv = (C.GPVERTEXATTRIB2SV)(getProcAddr("glVertexAttrib2sv"))
	if gpVertexAttrib2sv == nil {
		return errors.New("glVertexAttrib2sv")
	}
	gpVertexAttrib2svARB = (C.GPVERTEXATTRIB2SVARB)(getProcAddr("glVertexAttrib2svARB"))
	gpVertexAttrib2svNV = (C.GPVERTEXATTRIB2SVNV)(getProcAddr("glVertexAttrib2svNV"))
	gpVertexAttrib3d = (C.GPVERTEXATTRIB3D)(getProcAddr("glVertexAttrib3d"))
	if gpVertexAttrib3d == nil {
		return errors.New("glVertexAttrib3d")
	}
	gpVertexAttrib3dARB = (C.GPVERTEXATTRIB3DARB)(getProcAddr("glVertexAttrib3dARB"))
	gpVertexAttrib3dNV = (C.GPVERTEXATTRIB3DNV)(getProcAddr("glVertexAttrib3dNV"))
	gpVertexAttrib3dv = (C.GPVERTEXATTRIB3DV)(getProcAddr("glVertexAttrib3dv"))
	if gpVertexAttrib3dv == nil {
		return errors.New("glVertexAttrib3dv")
	}
	gpVertexAttrib3dvARB = (C.GPVERTEXATTRIB3DVARB)(getProcAddr("glVertexAttrib3dvARB"))
	gpVertexAttrib3dvNV = (C.GPVERTEXATTRIB3DVNV)(getProcAddr("glVertexAttrib3dvNV"))
	gpVertexAttrib3f = (C.GPVERTEXATTRIB3F)(getProcAddr("glVertexAttrib3f"))
	if gpVertexAttrib3f == nil {
		return errors.New("glVertexAttrib3f")
	}
	gpVertexAttrib3fARB = (C.GPVERTEXATTRIB3FARB)(getProcAddr("glVertexAttrib3fARB"))
	gpVertexAttrib3fNV = (C.GPVERTEXATTRIB3FNV)(getProcAddr("glVertexAttrib3fNV"))
	gpVertexAttrib3fv = (C.GPVERTEXATTRIB3FV)(getProcAddr("glVertexAttrib3fv"))
	if gpVertexAttrib3fv == nil {
		return errors.New("glVertexAttrib3fv")
	}
	gpVertexAttrib3fvARB = (C.GPVERTEXATTRIB3FVARB)(getProcAddr("glVertexAttrib3fvARB"))
	gpVertexAttrib3fvNV = (C.GPVERTEXATTRIB3FVNV)(getProcAddr("glVertexAttrib3fvNV"))
	gpVertexAttrib3hNV = (C.GPVERTEXATTRIB3HNV)(getProcAddr("glVertexAttrib3hNV"))
	gpVertexAttrib3hvNV = (C.GPVERTEXATTRIB3HVNV)(getProcAddr("glVertexAttrib3hvNV"))
	gpVertexAttrib3s = (C.GPVERTEXATTRIB3S)(getProcAddr("glVertexAttrib3s"))
	if gpVertexAttrib3s == nil {
		return errors.New("glVertexAttrib3s")
	}
	gpVertexAttrib3sARB = (C.GPVERTEXATTRIB3SARB)(getProcAddr("glVertexAttrib3sARB"))
	gpVertexAttrib3sNV = (C.GPVERTEXATTRIB3SNV)(getProcAddr("glVertexAttrib3sNV"))
	gpVertexAttrib3sv = (C.GPVERTEXATTRIB3SV)(getProcAddr("glVertexAttrib3sv"))
	if gpVertexAttrib3sv == nil {
		return errors.New("glVertexAttrib3sv")
	}
	gpVertexAttrib3svARB = (C.GPVERTEXATTRIB3SVARB)(getProcAddr("glVertexAttrib3svARB"))
	gpVertexAttrib3svNV = (C.GPVERTEXATTRIB3SVNV)(getProcAddr("glVertexAttrib3svNV"))
	gpVertexAttrib4Nbv = (C.GPVERTEXATTRIB4NBV)(getProcAddr("glVertexAttrib4Nbv"))
	if gpVertexAttrib4Nbv == nil {
		return errors.New("glVertexAttrib4Nbv")
	}
	gpVertexAttrib4NbvARB = (C.GPVERTEXATTRIB4NBVARB)(getProcAddr("glVertexAttrib4NbvARB"))
	gpVertexAttrib4Niv = (C.GPVERTEXATTRIB4NIV)(getProcAddr("glVertexAttrib4Niv"))
	if gpVertexAttrib4Niv == nil {
		return errors.New("glVertexAttrib4Niv")
	}
	gpVertexAttrib4NivARB = (C.GPVERTEXATTRIB4NIVARB)(getProcAddr("glVertexAttrib4NivARB"))
	gpVertexAttrib4Nsv = (C.GPVERTEXATTRIB4NSV)(getProcAddr("glVertexAttrib4Nsv"))
	if gpVertexAttrib4Nsv == nil {
		return errors.New("glVertexAttrib4Nsv")
	}
	gpVertexAttrib4NsvARB = (C.GPVERTEXATTRIB4NSVARB)(getProcAddr("glVertexAttrib4NsvARB"))
	gpVertexAttrib4Nub = (C.GPVERTEXATTRIB4NUB)(getProcAddr("glVertexAttrib4Nub"))
	if gpVertexAttrib4Nub == nil {
		return errors.New("glVertexAttrib4Nub")
	}
	gpVertexAttrib4NubARB = (C.GPVERTEXATTRIB4NUBARB)(getProcAddr("glVertexAttrib4NubARB"))
	gpVertexAttrib4Nubv = (C.GPVERTEXATTRIB4NUBV)(getProcAddr("glVertexAttrib4Nubv"))
	if gpVertexAttrib4Nubv == nil {
		return errors.New("glVertexAttrib4Nubv")
	}
	gpVertexAttrib4NubvARB = (C.GPVERTEXATTRIB4NUBVARB)(getProcAddr("glVertexAttrib4NubvARB"))
	gpVertexAttrib4Nuiv = (C.GPVERTEXATTRIB4NUIV)(getProcAddr("glVertexAttrib4Nuiv"))
	if gpVertexAttrib4Nuiv == nil {
		return errors.New("glVertexAttrib4Nuiv")
	}
	gpVertexAttrib4NuivARB = (C.GPVERTEXATTRIB4NUIVARB)(getProcAddr("glVertexAttrib4NuivARB"))
	gpVertexAttrib4Nusv = (C.GPVERTEXATTRIB4NUSV)(getProcAddr("glVertexAttrib4Nusv"))
	if gpVertexAttrib4Nusv == nil {
		return errors.New("glVertexAttrib4Nusv")
	}
	gpVertexAttrib4NusvARB = (C.GPVERTEXATTRIB4NUSVARB)(getProcAddr("glVertexAttrib4NusvARB"))
	gpVertexAttrib4bv = (C.GPVERTEXATTRIB4BV)(getProcAddr("glVertexAttrib4bv"))
	if gpVertexAttrib4bv == nil {
		return errors.New("glVertexAttrib4bv")
	}
	gpVertexAttrib4bvARB = (C.GPVERTEXATTRIB4BVARB)(getProcAddr("glVertexAttrib4bvARB"))
	gpVertexAttrib4d = (C.GPVERTEXATTRIB4D)(getProcAddr("glVertexAttrib4d"))
	if gpVertexAttrib4d == nil {
		return errors.New("glVertexAttrib4d")
	}
	gpVertexAttrib4dARB = (C.GPVERTEXATTRIB4DARB)(getProcAddr("glVertexAttrib4dARB"))
	gpVertexAttrib4dNV = (C.GPVERTEXATTRIB4DNV)(getProcAddr("glVertexAttrib4dNV"))
	gpVertexAttrib4dv = (C.GPVERTEXATTRIB4DV)(getProcAddr("glVertexAttrib4dv"))
	if gpVertexAttrib4dv == nil {
		return errors.New("glVertexAttrib4dv")
	}
	gpVertexAttrib4dvARB = (C.GPVERTEXATTRIB4DVARB)(getProcAddr("glVertexAttrib4dvARB"))
	gpVertexAttrib4dvNV = (C.GPVERTEXATTRIB4DVNV)(getProcAddr("glVertexAttrib4dvNV"))
	gpVertexAttrib4f = (C.GPVERTEXATTRIB4F)(getProcAddr("glVertexAttrib4f"))
	if gpVertexAttrib4f == nil {
		return errors.New("glVertexAttrib4f")
	}
	gpVertexAttrib4fARB = (C.GPVERTEXATTRIB4FARB)(getProcAddr("glVertexAttrib4fARB"))
	gpVertexAttrib4fNV = (C.GPVERTEXATTRIB4FNV)(getProcAddr("glVertexAttrib4fNV"))
	gpVertexAttrib4fv = (C.GPVERTEXATTRIB4FV)(getProcAddr("glVertexAttrib4fv"))
	if gpVertexAttrib4fv == nil {
		return errors.New("glVertexAttrib4fv")
	}
	gpVertexAttrib4fvARB = (C.GPVERTEXATTRIB4FVARB)(getProcAddr("glVertexAttrib4fvARB"))
	gpVertexAttrib4fvNV = (C.GPVERTEXATTRIB4FVNV)(getProcAddr("glVertexAttrib4fvNV"))
	gpVertexAttrib4hNV = (C.GPVERTEXATTRIB4HNV)(getProcAddr("glVertexAttrib4hNV"))
	gpVertexAttrib4hvNV = (C.GPVERTEXATTRIB4HVNV)(getProcAddr("glVertexAttrib4hvNV"))
	gpVertexAttrib4iv = (C.GPVERTEXATTRIB4IV)(getProcAddr("glVertexAttrib4iv"))
	if gpVertexAttrib4iv == nil {
		return errors.New("glVertexAttrib4iv")
	}
	gpVertexAttrib4ivARB = (C.GPVERTEXATTRIB4IVARB)(getProcAddr("glVertexAttrib4ivARB"))
	gpVertexAttrib4s = (C.GPVERTEXATTRIB4S)(getProcAddr("glVertexAttrib4s"))
	if gpVertexAttrib4s == nil {
		return errors.New("glVertexAttrib4s")
	}
	gpVertexAttrib4sARB = (C.GPVERTEXATTRIB4SARB)(getProcAddr("glVertexAttrib4sARB"))
	gpVertexAttrib4sNV = (C.GPVERTEXATTRIB4SNV)(getProcAddr("glVertexAttrib4sNV"))
	gpVertexAttrib4sv = (C.GPVERTEXATTRIB4SV)(getProcAddr("glVertexAttrib4sv"))
	if gpVertexAttrib4sv == nil {
		return errors.New("glVertexAttrib4sv")
	}
	gpVertexAttrib4svARB = (C.GPVERTEXATTRIB4SVARB)(getProcAddr("glVertexAttrib4svARB"))
	gpVertexAttrib4svNV = (C.GPVERTEXATTRIB4SVNV)(getProcAddr("glVertexAttrib4svNV"))
	gpVertexAttrib4ubNV = (C.GPVERTEXATTRIB4UBNV)(getProcAddr("glVertexAttrib4ubNV"))
	gpVertexAttrib4ubv = (C.GPVERTEXATTRIB4UBV)(getProcAddr("glVertexAttrib4ubv"))
	if gpVertexAttrib4ubv == nil {
		return errors.New("glVertexAttrib4ubv")
	}
	gpVertexAttrib4ubvARB = (C.GPVERTEXATTRIB4UBVARB)(getProcAddr("glVertexAttrib4ubvARB"))
	gpVertexAttrib4ubvNV = (C.GPVERTEXATTRIB4UBVNV)(getProcAddr("glVertexAttrib4ubvNV"))
	gpVertexAttrib4uiv = (C.GPVERTEXATTRIB4UIV)(getProcAddr("glVertexAttrib4uiv"))
	if gpVertexAttrib4uiv == nil {
		return errors.New("glVertexAttrib4uiv")
	}
	gpVertexAttrib4uivARB = (C.GPVERTEXATTRIB4UIVARB)(getProcAddr("glVertexAttrib4uivARB"))
	gpVertexAttrib4usv = (C.GPVERTEXATTRIB4USV)(getProcAddr("glVertexAttrib4usv"))
	if gpVertexAttrib4usv == nil {
		return errors.New("glVertexAttrib4usv")
	}
	gpVertexAttrib4usvARB = (C.GPVERTEXATTRIB4USVARB)(getProcAddr("glVertexAttrib4usvARB"))
	gpVertexAttribArrayObjectATI = (C.GPVERTEXATTRIBARRAYOBJECTATI)(getProcAddr("glVertexAttribArrayObjectATI"))
	gpVertexAttribBinding = (C.GPVERTEXATTRIBBINDING)(getProcAddr("glVertexAttribBinding"))
	if gpVertexAttribBinding == nil {
		return errors.New("glVertexAttribBinding")
	}
	gpVertexAttribDivisor = (C.GPVERTEXATTRIBDIVISOR)(getProcAddr("glVertexAttribDivisor"))
	if gpVertexAttribDivisor == nil {
		return errors.New("glVertexAttribDivisor")
	}
	gpVertexAttribDivisorARB = (C.GPVERTEXATTRIBDIVISORARB)(getProcAddr("glVertexAttribDivisorARB"))
	gpVertexAttribFormat = (C.GPVERTEXATTRIBFORMAT)(getProcAddr("glVertexAttribFormat"))
	if gpVertexAttribFormat == nil {
		return errors.New("glVertexAttribFormat")
	}
	gpVertexAttribFormatNV = (C.GPVERTEXATTRIBFORMATNV)(getProcAddr("glVertexAttribFormatNV"))
	gpVertexAttribI1i = (C.GPVERTEXATTRIBI1I)(getProcAddr("glVertexAttribI1i"))
	if gpVertexAttribI1i == nil {
		return errors.New("glVertexAttribI1i")
	}
	gpVertexAttribI1iEXT = (C.GPVERTEXATTRIBI1IEXT)(getProcAddr("glVertexAttribI1iEXT"))
	gpVertexAttribI1iv = (C.GPVERTEXATTRIBI1IV)(getProcAddr("glVertexAttribI1iv"))
	if gpVertexAttribI1iv == nil {
		return errors.New("glVertexAttribI1iv")
	}
	gpVertexAttribI1ivEXT = (C.GPVERTEXATTRIBI1IVEXT)(getProcAddr("glVertexAttribI1ivEXT"))
	gpVertexAttribI1ui = (C.GPVERTEXATTRIBI1UI)(getProcAddr("glVertexAttribI1ui"))
	if gpVertexAttribI1ui == nil {
		return errors.New("glVertexAttribI1ui")
	}
	gpVertexAttribI1uiEXT = (C.GPVERTEXATTRIBI1UIEXT)(getProcAddr("glVertexAttribI1uiEXT"))
	gpVertexAttribI1uiv = (C.GPVERTEXATTRIBI1UIV)(getProcAddr("glVertexAttribI1uiv"))
	if gpVertexAttribI1uiv == nil {
		return errors.New("glVertexAttribI1uiv")
	}
	gpVertexAttribI1uivEXT = (C.GPVERTEXATTRIBI1UIVEXT)(getProcAddr("glVertexAttribI1uivEXT"))
	gpVertexAttribI2i = (C.GPVERTEXATTRIBI2I)(getProcAddr("glVertexAttribI2i"))
	if gpVertexAttribI2i == nil {
		return errors.New("glVertexAttribI2i")
	}
	gpVertexAttribI2iEXT = (C.GPVERTEXATTRIBI2IEXT)(getProcAddr("glVertexAttribI2iEXT"))
	gpVertexAttribI2iv = (C.GPVERTEXATTRIBI2IV)(getProcAddr("glVertexAttribI2iv"))
	if gpVertexAttribI2iv == nil {
		return errors.New("glVertexAttribI2iv")
	}
	gpVertexAttribI2ivEXT = (C.GPVERTEXATTRIBI2IVEXT)(getProcAddr("glVertexAttribI2ivEXT"))
	gpVertexAttribI2ui = (C.GPVERTEXATTRIBI2UI)(getProcAddr("glVertexAttribI2ui"))
	if gpVertexAttribI2ui == nil {
		return errors.New("glVertexAttribI2ui")
	}
	gpVertexAttribI2uiEXT = (C.GPVERTEXATTRIBI2UIEXT)(getProcAddr("glVertexAttribI2uiEXT"))
	gpVertexAttribI2uiv = (C.GPVERTEXATTRIBI2UIV)(getProcAddr("glVertexAttribI2uiv"))
	if gpVertexAttribI2uiv == nil {
		return errors.New("glVertexAttribI2uiv")
	}
	gpVertexAttribI2uivEXT = (C.GPVERTEXATTRIBI2UIVEXT)(getProcAddr("glVertexAttribI2uivEXT"))
	gpVertexAttribI3i = (C.GPVERTEXATTRIBI3I)(getProcAddr("glVertexAttribI3i"))
	if gpVertexAttribI3i == nil {
		return errors.New("glVertexAttribI3i")
	}
	gpVertexAttribI3iEXT = (C.GPVERTEXATTRIBI3IEXT)(getProcAddr("glVertexAttribI3iEXT"))
	gpVertexAttribI3iv = (C.GPVERTEXATTRIBI3IV)(getProcAddr("glVertexAttribI3iv"))
	if gpVertexAttribI3iv == nil {
		return errors.New("glVertexAttribI3iv")
	}
	gpVertexAttribI3ivEXT = (C.GPVERTEXATTRIBI3IVEXT)(getProcAddr("glVertexAttribI3ivEXT"))
	gpVertexAttribI3ui = (C.GPVERTEXATTRIBI3UI)(getProcAddr("glVertexAttribI3ui"))
	if gpVertexAttribI3ui == nil {
		return errors.New("glVertexAttribI3ui")
	}
	gpVertexAttribI3uiEXT = (C.GPVERTEXATTRIBI3UIEXT)(getProcAddr("glVertexAttribI3uiEXT"))
	gpVertexAttribI3uiv = (C.GPVERTEXATTRIBI3UIV)(getProcAddr("glVertexAttribI3uiv"))
	if gpVertexAttribI3uiv == nil {
		return errors.New("glVertexAttribI3uiv")
	}
	gpVertexAttribI3uivEXT = (C.GPVERTEXATTRIBI3UIVEXT)(getProcAddr("glVertexAttribI3uivEXT"))
	gpVertexAttribI4bv = (C.GPVERTEXATTRIBI4BV)(getProcAddr("glVertexAttribI4bv"))
	if gpVertexAttribI4bv == nil {
		return errors.New("glVertexAttribI4bv")
	}
	gpVertexAttribI4bvEXT = (C.GPVERTEXATTRIBI4BVEXT)(getProcAddr("glVertexAttribI4bvEXT"))
	gpVertexAttribI4i = (C.GPVERTEXATTRIBI4I)(getProcAddr("glVertexAttribI4i"))
	if gpVertexAttribI4i == nil {
		return errors.New("glVertexAttribI4i")
	}
	gpVertexAttribI4iEXT = (C.GPVERTEXATTRIBI4IEXT)(getProcAddr("glVertexAttribI4iEXT"))
	gpVertexAttribI4iv = (C.GPVERTEXATTRIBI4IV)(getProcAddr("glVertexAttribI4iv"))
	if gpVertexAttribI4iv == nil {
		return errors.New("glVertexAttribI4iv")
	}
	gpVertexAttribI4ivEXT = (C.GPVERTEXATTRIBI4IVEXT)(getProcAddr("glVertexAttribI4ivEXT"))
	gpVertexAttribI4sv = (C.GPVERTEXATTRIBI4SV)(getProcAddr("glVertexAttribI4sv"))
	if gpVertexAttribI4sv == nil {
		return errors.New("glVertexAttribI4sv")
	}
	gpVertexAttribI4svEXT = (C.GPVERTEXATTRIBI4SVEXT)(getProcAddr("glVertexAttribI4svEXT"))
	gpVertexAttribI4ubv = (C.GPVERTEXATTRIBI4UBV)(getProcAddr("glVertexAttribI4ubv"))
	if gpVertexAttribI4ubv == nil {
		return errors.New("glVertexAttribI4ubv")
	}
	gpVertexAttribI4ubvEXT = (C.GPVERTEXATTRIBI4UBVEXT)(getProcAddr("glVertexAttribI4ubvEXT"))
	gpVertexAttribI4ui = (C.GPVERTEXATTRIBI4UI)(getProcAddr("glVertexAttribI4ui"))
	if gpVertexAttribI4ui == nil {
		return errors.New("glVertexAttribI4ui")
	}
	gpVertexAttribI4uiEXT = (C.GPVERTEXATTRIBI4UIEXT)(getProcAddr("glVertexAttribI4uiEXT"))
	gpVertexAttribI4uiv = (C.GPVERTEXATTRIBI4UIV)(getProcAddr("glVertexAttribI4uiv"))
	if gpVertexAttribI4uiv == nil {
		return errors.New("glVertexAttribI4uiv")
	}
	gpVertexAttribI4uivEXT = (C.GPVERTEXATTRIBI4UIVEXT)(getProcAddr("glVertexAttribI4uivEXT"))
	gpVertexAttribI4usv = (C.GPVERTEXATTRIBI4USV)(getProcAddr("glVertexAttribI4usv"))
	if gpVertexAttribI4usv == nil {
		return errors.New("glVertexAttribI4usv")
	}
	gpVertexAttribI4usvEXT = (C.GPVERTEXATTRIBI4USVEXT)(getProcAddr("glVertexAttribI4usvEXT"))
	gpVertexAttribIFormat = (C.GPVERTEXATTRIBIFORMAT)(getProcAddr("glVertexAttribIFormat"))
	if gpVertexAttribIFormat == nil {
		return errors.New("glVertexAttribIFormat")
	}
	gpVertexAttribIFormatNV = (C.GPVERTEXATTRIBIFORMATNV)(getProcAddr("glVertexAttribIFormatNV"))
	gpVertexAttribIPointer = (C.GPVERTEXATTRIBIPOINTER)(getProcAddr("glVertexAttribIPointer"))
	if gpVertexAttribIPointer == nil {
		return errors.New("glVertexAttribIPointer")
	}
	gpVertexAttribIPointerEXT = (C.GPVERTEXATTRIBIPOINTEREXT)(getProcAddr("glVertexAttribIPointerEXT"))
	gpVertexAttribL1d = (C.GPVERTEXATTRIBL1D)(getProcAddr("glVertexAttribL1d"))
	if gpVertexAttribL1d == nil {
		return errors.New("glVertexAttribL1d")
	}
	gpVertexAttribL1dEXT = (C.GPVERTEXATTRIBL1DEXT)(getProcAddr("glVertexAttribL1dEXT"))
	gpVertexAttribL1dv = (C.GPVERTEXATTRIBL1DV)(getProcAddr("glVertexAttribL1dv"))
	if gpVertexAttribL1dv == nil {
		return errors.New("glVertexAttribL1dv")
	}
	gpVertexAttribL1dvEXT = (C.GPVERTEXATTRIBL1DVEXT)(getProcAddr("glVertexAttribL1dvEXT"))
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
	gpVertexAttribL2dEXT = (C.GPVERTEXATTRIBL2DEXT)(getProcAddr("glVertexAttribL2dEXT"))
	gpVertexAttribL2dv = (C.GPVERTEXATTRIBL2DV)(getProcAddr("glVertexAttribL2dv"))
	if gpVertexAttribL2dv == nil {
		return errors.New("glVertexAttribL2dv")
	}
	gpVertexAttribL2dvEXT = (C.GPVERTEXATTRIBL2DVEXT)(getProcAddr("glVertexAttribL2dvEXT"))
	gpVertexAttribL2i64NV = (C.GPVERTEXATTRIBL2I64NV)(getProcAddr("glVertexAttribL2i64NV"))
	gpVertexAttribL2i64vNV = (C.GPVERTEXATTRIBL2I64VNV)(getProcAddr("glVertexAttribL2i64vNV"))
	gpVertexAttribL2ui64NV = (C.GPVERTEXATTRIBL2UI64NV)(getProcAddr("glVertexAttribL2ui64NV"))
	gpVertexAttribL2ui64vNV = (C.GPVERTEXATTRIBL2UI64VNV)(getProcAddr("glVertexAttribL2ui64vNV"))
	gpVertexAttribL3d = (C.GPVERTEXATTRIBL3D)(getProcAddr("glVertexAttribL3d"))
	if gpVertexAttribL3d == nil {
		return errors.New("glVertexAttribL3d")
	}
	gpVertexAttribL3dEXT = (C.GPVERTEXATTRIBL3DEXT)(getProcAddr("glVertexAttribL3dEXT"))
	gpVertexAttribL3dv = (C.GPVERTEXATTRIBL3DV)(getProcAddr("glVertexAttribL3dv"))
	if gpVertexAttribL3dv == nil {
		return errors.New("glVertexAttribL3dv")
	}
	gpVertexAttribL3dvEXT = (C.GPVERTEXATTRIBL3DVEXT)(getProcAddr("glVertexAttribL3dvEXT"))
	gpVertexAttribL3i64NV = (C.GPVERTEXATTRIBL3I64NV)(getProcAddr("glVertexAttribL3i64NV"))
	gpVertexAttribL3i64vNV = (C.GPVERTEXATTRIBL3I64VNV)(getProcAddr("glVertexAttribL3i64vNV"))
	gpVertexAttribL3ui64NV = (C.GPVERTEXATTRIBL3UI64NV)(getProcAddr("glVertexAttribL3ui64NV"))
	gpVertexAttribL3ui64vNV = (C.GPVERTEXATTRIBL3UI64VNV)(getProcAddr("glVertexAttribL3ui64vNV"))
	gpVertexAttribL4d = (C.GPVERTEXATTRIBL4D)(getProcAddr("glVertexAttribL4d"))
	if gpVertexAttribL4d == nil {
		return errors.New("glVertexAttribL4d")
	}
	gpVertexAttribL4dEXT = (C.GPVERTEXATTRIBL4DEXT)(getProcAddr("glVertexAttribL4dEXT"))
	gpVertexAttribL4dv = (C.GPVERTEXATTRIBL4DV)(getProcAddr("glVertexAttribL4dv"))
	if gpVertexAttribL4dv == nil {
		return errors.New("glVertexAttribL4dv")
	}
	gpVertexAttribL4dvEXT = (C.GPVERTEXATTRIBL4DVEXT)(getProcAddr("glVertexAttribL4dvEXT"))
	gpVertexAttribL4i64NV = (C.GPVERTEXATTRIBL4I64NV)(getProcAddr("glVertexAttribL4i64NV"))
	gpVertexAttribL4i64vNV = (C.GPVERTEXATTRIBL4I64VNV)(getProcAddr("glVertexAttribL4i64vNV"))
	gpVertexAttribL4ui64NV = (C.GPVERTEXATTRIBL4UI64NV)(getProcAddr("glVertexAttribL4ui64NV"))
	gpVertexAttribL4ui64vNV = (C.GPVERTEXATTRIBL4UI64VNV)(getProcAddr("glVertexAttribL4ui64vNV"))
	gpVertexAttribLFormat = (C.GPVERTEXATTRIBLFORMAT)(getProcAddr("glVertexAttribLFormat"))
	if gpVertexAttribLFormat == nil {
		return errors.New("glVertexAttribLFormat")
	}
	gpVertexAttribLFormatNV = (C.GPVERTEXATTRIBLFORMATNV)(getProcAddr("glVertexAttribLFormatNV"))
	gpVertexAttribLPointer = (C.GPVERTEXATTRIBLPOINTER)(getProcAddr("glVertexAttribLPointer"))
	if gpVertexAttribLPointer == nil {
		return errors.New("glVertexAttribLPointer")
	}
	gpVertexAttribLPointerEXT = (C.GPVERTEXATTRIBLPOINTEREXT)(getProcAddr("glVertexAttribLPointerEXT"))
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
	gpVertexAttribParameteriAMD = (C.GPVERTEXATTRIBPARAMETERIAMD)(getProcAddr("glVertexAttribParameteriAMD"))
	gpVertexAttribPointer = (C.GPVERTEXATTRIBPOINTER)(getProcAddr("glVertexAttribPointer"))
	if gpVertexAttribPointer == nil {
		return errors.New("glVertexAttribPointer")
	}
	gpVertexAttribPointerARB = (C.GPVERTEXATTRIBPOINTERARB)(getProcAddr("glVertexAttribPointerARB"))
	gpVertexAttribPointerNV = (C.GPVERTEXATTRIBPOINTERNV)(getProcAddr("glVertexAttribPointerNV"))
	gpVertexAttribs1dvNV = (C.GPVERTEXATTRIBS1DVNV)(getProcAddr("glVertexAttribs1dvNV"))
	gpVertexAttribs1fvNV = (C.GPVERTEXATTRIBS1FVNV)(getProcAddr("glVertexAttribs1fvNV"))
	gpVertexAttribs1hvNV = (C.GPVERTEXATTRIBS1HVNV)(getProcAddr("glVertexAttribs1hvNV"))
	gpVertexAttribs1svNV = (C.GPVERTEXATTRIBS1SVNV)(getProcAddr("glVertexAttribs1svNV"))
	gpVertexAttribs2dvNV = (C.GPVERTEXATTRIBS2DVNV)(getProcAddr("glVertexAttribs2dvNV"))
	gpVertexAttribs2fvNV = (C.GPVERTEXATTRIBS2FVNV)(getProcAddr("glVertexAttribs2fvNV"))
	gpVertexAttribs2hvNV = (C.GPVERTEXATTRIBS2HVNV)(getProcAddr("glVertexAttribs2hvNV"))
	gpVertexAttribs2svNV = (C.GPVERTEXATTRIBS2SVNV)(getProcAddr("glVertexAttribs2svNV"))
	gpVertexAttribs3dvNV = (C.GPVERTEXATTRIBS3DVNV)(getProcAddr("glVertexAttribs3dvNV"))
	gpVertexAttribs3fvNV = (C.GPVERTEXATTRIBS3FVNV)(getProcAddr("glVertexAttribs3fvNV"))
	gpVertexAttribs3hvNV = (C.GPVERTEXATTRIBS3HVNV)(getProcAddr("glVertexAttribs3hvNV"))
	gpVertexAttribs3svNV = (C.GPVERTEXATTRIBS3SVNV)(getProcAddr("glVertexAttribs3svNV"))
	gpVertexAttribs4dvNV = (C.GPVERTEXATTRIBS4DVNV)(getProcAddr("glVertexAttribs4dvNV"))
	gpVertexAttribs4fvNV = (C.GPVERTEXATTRIBS4FVNV)(getProcAddr("glVertexAttribs4fvNV"))
	gpVertexAttribs4hvNV = (C.GPVERTEXATTRIBS4HVNV)(getProcAddr("glVertexAttribs4hvNV"))
	gpVertexAttribs4svNV = (C.GPVERTEXATTRIBS4SVNV)(getProcAddr("glVertexAttribs4svNV"))
	gpVertexAttribs4ubvNV = (C.GPVERTEXATTRIBS4UBVNV)(getProcAddr("glVertexAttribs4ubvNV"))
	gpVertexBindingDivisor = (C.GPVERTEXBINDINGDIVISOR)(getProcAddr("glVertexBindingDivisor"))
	if gpVertexBindingDivisor == nil {
		return errors.New("glVertexBindingDivisor")
	}
	gpVertexBlendARB = (C.GPVERTEXBLENDARB)(getProcAddr("glVertexBlendARB"))
	gpVertexBlendEnvfATI = (C.GPVERTEXBLENDENVFATI)(getProcAddr("glVertexBlendEnvfATI"))
	gpVertexBlendEnviATI = (C.GPVERTEXBLENDENVIATI)(getProcAddr("glVertexBlendEnviATI"))
	gpVertexFormatNV = (C.GPVERTEXFORMATNV)(getProcAddr("glVertexFormatNV"))
	gpVertexP2ui = (C.GPVERTEXP2UI)(getProcAddr("glVertexP2ui"))
	if gpVertexP2ui == nil {
		return errors.New("glVertexP2ui")
	}
	gpVertexP2uiv = (C.GPVERTEXP2UIV)(getProcAddr("glVertexP2uiv"))
	if gpVertexP2uiv == nil {
		return errors.New("glVertexP2uiv")
	}
	gpVertexP3ui = (C.GPVERTEXP3UI)(getProcAddr("glVertexP3ui"))
	if gpVertexP3ui == nil {
		return errors.New("glVertexP3ui")
	}
	gpVertexP3uiv = (C.GPVERTEXP3UIV)(getProcAddr("glVertexP3uiv"))
	if gpVertexP3uiv == nil {
		return errors.New("glVertexP3uiv")
	}
	gpVertexP4ui = (C.GPVERTEXP4UI)(getProcAddr("glVertexP4ui"))
	if gpVertexP4ui == nil {
		return errors.New("glVertexP4ui")
	}
	gpVertexP4uiv = (C.GPVERTEXP4UIV)(getProcAddr("glVertexP4uiv"))
	if gpVertexP4uiv == nil {
		return errors.New("glVertexP4uiv")
	}
	gpVertexPointer = (C.GPVERTEXPOINTER)(getProcAddr("glVertexPointer"))
	if gpVertexPointer == nil {
		return errors.New("glVertexPointer")
	}
	gpVertexPointerEXT = (C.GPVERTEXPOINTEREXT)(getProcAddr("glVertexPointerEXT"))
	gpVertexPointerListIBM = (C.GPVERTEXPOINTERLISTIBM)(getProcAddr("glVertexPointerListIBM"))
	gpVertexPointervINTEL = (C.GPVERTEXPOINTERVINTEL)(getProcAddr("glVertexPointervINTEL"))
	gpVertexStream1dATI = (C.GPVERTEXSTREAM1DATI)(getProcAddr("glVertexStream1dATI"))
	gpVertexStream1dvATI = (C.GPVERTEXSTREAM1DVATI)(getProcAddr("glVertexStream1dvATI"))
	gpVertexStream1fATI = (C.GPVERTEXSTREAM1FATI)(getProcAddr("glVertexStream1fATI"))
	gpVertexStream1fvATI = (C.GPVERTEXSTREAM1FVATI)(getProcAddr("glVertexStream1fvATI"))
	gpVertexStream1iATI = (C.GPVERTEXSTREAM1IATI)(getProcAddr("glVertexStream1iATI"))
	gpVertexStream1ivATI = (C.GPVERTEXSTREAM1IVATI)(getProcAddr("glVertexStream1ivATI"))
	gpVertexStream1sATI = (C.GPVERTEXSTREAM1SATI)(getProcAddr("glVertexStream1sATI"))
	gpVertexStream1svATI = (C.GPVERTEXSTREAM1SVATI)(getProcAddr("glVertexStream1svATI"))
	gpVertexStream2dATI = (C.GPVERTEXSTREAM2DATI)(getProcAddr("glVertexStream2dATI"))
	gpVertexStream2dvATI = (C.GPVERTEXSTREAM2DVATI)(getProcAddr("glVertexStream2dvATI"))
	gpVertexStream2fATI = (C.GPVERTEXSTREAM2FATI)(getProcAddr("glVertexStream2fATI"))
	gpVertexStream2fvATI = (C.GPVERTEXSTREAM2FVATI)(getProcAddr("glVertexStream2fvATI"))
	gpVertexStream2iATI = (C.GPVERTEXSTREAM2IATI)(getProcAddr("glVertexStream2iATI"))
	gpVertexStream2ivATI = (C.GPVERTEXSTREAM2IVATI)(getProcAddr("glVertexStream2ivATI"))
	gpVertexStream2sATI = (C.GPVERTEXSTREAM2SATI)(getProcAddr("glVertexStream2sATI"))
	gpVertexStream2svATI = (C.GPVERTEXSTREAM2SVATI)(getProcAddr("glVertexStream2svATI"))
	gpVertexStream3dATI = (C.GPVERTEXSTREAM3DATI)(getProcAddr("glVertexStream3dATI"))
	gpVertexStream3dvATI = (C.GPVERTEXSTREAM3DVATI)(getProcAddr("glVertexStream3dvATI"))
	gpVertexStream3fATI = (C.GPVERTEXSTREAM3FATI)(getProcAddr("glVertexStream3fATI"))
	gpVertexStream3fvATI = (C.GPVERTEXSTREAM3FVATI)(getProcAddr("glVertexStream3fvATI"))
	gpVertexStream3iATI = (C.GPVERTEXSTREAM3IATI)(getProcAddr("glVertexStream3iATI"))
	gpVertexStream3ivATI = (C.GPVERTEXSTREAM3IVATI)(getProcAddr("glVertexStream3ivATI"))
	gpVertexStream3sATI = (C.GPVERTEXSTREAM3SATI)(getProcAddr("glVertexStream3sATI"))
	gpVertexStream3svATI = (C.GPVERTEXSTREAM3SVATI)(getProcAddr("glVertexStream3svATI"))
	gpVertexStream4dATI = (C.GPVERTEXSTREAM4DATI)(getProcAddr("glVertexStream4dATI"))
	gpVertexStream4dvATI = (C.GPVERTEXSTREAM4DVATI)(getProcAddr("glVertexStream4dvATI"))
	gpVertexStream4fATI = (C.GPVERTEXSTREAM4FATI)(getProcAddr("glVertexStream4fATI"))
	gpVertexStream4fvATI = (C.GPVERTEXSTREAM4FVATI)(getProcAddr("glVertexStream4fvATI"))
	gpVertexStream4iATI = (C.GPVERTEXSTREAM4IATI)(getProcAddr("glVertexStream4iATI"))
	gpVertexStream4ivATI = (C.GPVERTEXSTREAM4IVATI)(getProcAddr("glVertexStream4ivATI"))
	gpVertexStream4sATI = (C.GPVERTEXSTREAM4SATI)(getProcAddr("glVertexStream4sATI"))
	gpVertexStream4svATI = (C.GPVERTEXSTREAM4SVATI)(getProcAddr("glVertexStream4svATI"))
	gpVertexWeightPointerEXT = (C.GPVERTEXWEIGHTPOINTEREXT)(getProcAddr("glVertexWeightPointerEXT"))
	gpVertexWeightfEXT = (C.GPVERTEXWEIGHTFEXT)(getProcAddr("glVertexWeightfEXT"))
	gpVertexWeightfvEXT = (C.GPVERTEXWEIGHTFVEXT)(getProcAddr("glVertexWeightfvEXT"))
	gpVertexWeighthNV = (C.GPVERTEXWEIGHTHNV)(getProcAddr("glVertexWeighthNV"))
	gpVertexWeighthvNV = (C.GPVERTEXWEIGHTHVNV)(getProcAddr("glVertexWeighthvNV"))
	gpVideoCaptureNV = (C.GPVIDEOCAPTURENV)(getProcAddr("glVideoCaptureNV"))
	gpVideoCaptureStreamParameterdvNV = (C.GPVIDEOCAPTURESTREAMPARAMETERDVNV)(getProcAddr("glVideoCaptureStreamParameterdvNV"))
	gpVideoCaptureStreamParameterfvNV = (C.GPVIDEOCAPTURESTREAMPARAMETERFVNV)(getProcAddr("glVideoCaptureStreamParameterfvNV"))
	gpVideoCaptureStreamParameterivNV = (C.GPVIDEOCAPTURESTREAMPARAMETERIVNV)(getProcAddr("glVideoCaptureStreamParameterivNV"))
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
	gpWaitSemaphoreEXT = (C.GPWAITSEMAPHOREEXT)(getProcAddr("glWaitSemaphoreEXT"))
	gpWaitSync = (C.GPWAITSYNC)(getProcAddr("glWaitSync"))
	if gpWaitSync == nil {
		return errors.New("glWaitSync")
	}
	gpWaitVkSemaphoreNV = (C.GPWAITVKSEMAPHORENV)(getProcAddr("glWaitVkSemaphoreNV"))
	gpWeightPathsNV = (C.GPWEIGHTPATHSNV)(getProcAddr("glWeightPathsNV"))
	gpWeightPointerARB = (C.GPWEIGHTPOINTERARB)(getProcAddr("glWeightPointerARB"))
	gpWeightbvARB = (C.GPWEIGHTBVARB)(getProcAddr("glWeightbvARB"))
	gpWeightdvARB = (C.GPWEIGHTDVARB)(getProcAddr("glWeightdvARB"))
	gpWeightfvARB = (C.GPWEIGHTFVARB)(getProcAddr("glWeightfvARB"))
	gpWeightivARB = (C.GPWEIGHTIVARB)(getProcAddr("glWeightivARB"))
	gpWeightsvARB = (C.GPWEIGHTSVARB)(getProcAddr("glWeightsvARB"))
	gpWeightubvARB = (C.GPWEIGHTUBVARB)(getProcAddr("glWeightubvARB"))
	gpWeightuivARB = (C.GPWEIGHTUIVARB)(getProcAddr("glWeightuivARB"))
	gpWeightusvARB = (C.GPWEIGHTUSVARB)(getProcAddr("glWeightusvARB"))
	gpWindowPos2d = (C.GPWINDOWPOS2D)(getProcAddr("glWindowPos2d"))
	if gpWindowPos2d == nil {
		return errors.New("glWindowPos2d")
	}
	gpWindowPos2dARB = (C.GPWINDOWPOS2DARB)(getProcAddr("glWindowPos2dARB"))
	gpWindowPos2dMESA = (C.GPWINDOWPOS2DMESA)(getProcAddr("glWindowPos2dMESA"))
	gpWindowPos2dv = (C.GPWINDOWPOS2DV)(getProcAddr("glWindowPos2dv"))
	if gpWindowPos2dv == nil {
		return errors.New("glWindowPos2dv")
	}
	gpWindowPos2dvARB = (C.GPWINDOWPOS2DVARB)(getProcAddr("glWindowPos2dvARB"))
	gpWindowPos2dvMESA = (C.GPWINDOWPOS2DVMESA)(getProcAddr("glWindowPos2dvMESA"))
	gpWindowPos2f = (C.GPWINDOWPOS2F)(getProcAddr("glWindowPos2f"))
	if gpWindowPos2f == nil {
		return errors.New("glWindowPos2f")
	}
	gpWindowPos2fARB = (C.GPWINDOWPOS2FARB)(getProcAddr("glWindowPos2fARB"))
	gpWindowPos2fMESA = (C.GPWINDOWPOS2FMESA)(getProcAddr("glWindowPos2fMESA"))
	gpWindowPos2fv = (C.GPWINDOWPOS2FV)(getProcAddr("glWindowPos2fv"))
	if gpWindowPos2fv == nil {
		return errors.New("glWindowPos2fv")
	}
	gpWindowPos2fvARB = (C.GPWINDOWPOS2FVARB)(getProcAddr("glWindowPos2fvARB"))
	gpWindowPos2fvMESA = (C.GPWINDOWPOS2FVMESA)(getProcAddr("glWindowPos2fvMESA"))
	gpWindowPos2i = (C.GPWINDOWPOS2I)(getProcAddr("glWindowPos2i"))
	if gpWindowPos2i == nil {
		return errors.New("glWindowPos2i")
	}
	gpWindowPos2iARB = (C.GPWINDOWPOS2IARB)(getProcAddr("glWindowPos2iARB"))
	gpWindowPos2iMESA = (C.GPWINDOWPOS2IMESA)(getProcAddr("glWindowPos2iMESA"))
	gpWindowPos2iv = (C.GPWINDOWPOS2IV)(getProcAddr("glWindowPos2iv"))
	if gpWindowPos2iv == nil {
		return errors.New("glWindowPos2iv")
	}
	gpWindowPos2ivARB = (C.GPWINDOWPOS2IVARB)(getProcAddr("glWindowPos2ivARB"))
	gpWindowPos2ivMESA = (C.GPWINDOWPOS2IVMESA)(getProcAddr("glWindowPos2ivMESA"))
	gpWindowPos2s = (C.GPWINDOWPOS2S)(getProcAddr("glWindowPos2s"))
	if gpWindowPos2s == nil {
		return errors.New("glWindowPos2s")
	}
	gpWindowPos2sARB = (C.GPWINDOWPOS2SARB)(getProcAddr("glWindowPos2sARB"))
	gpWindowPos2sMESA = (C.GPWINDOWPOS2SMESA)(getProcAddr("glWindowPos2sMESA"))
	gpWindowPos2sv = (C.GPWINDOWPOS2SV)(getProcAddr("glWindowPos2sv"))
	if gpWindowPos2sv == nil {
		return errors.New("glWindowPos2sv")
	}
	gpWindowPos2svARB = (C.GPWINDOWPOS2SVARB)(getProcAddr("glWindowPos2svARB"))
	gpWindowPos2svMESA = (C.GPWINDOWPOS2SVMESA)(getProcAddr("glWindowPos2svMESA"))
	gpWindowPos3d = (C.GPWINDOWPOS3D)(getProcAddr("glWindowPos3d"))
	if gpWindowPos3d == nil {
		return errors.New("glWindowPos3d")
	}
	gpWindowPos3dARB = (C.GPWINDOWPOS3DARB)(getProcAddr("glWindowPos3dARB"))
	gpWindowPos3dMESA = (C.GPWINDOWPOS3DMESA)(getProcAddr("glWindowPos3dMESA"))
	gpWindowPos3dv = (C.GPWINDOWPOS3DV)(getProcAddr("glWindowPos3dv"))
	if gpWindowPos3dv == nil {
		return errors.New("glWindowPos3dv")
	}
	gpWindowPos3dvARB = (C.GPWINDOWPOS3DVARB)(getProcAddr("glWindowPos3dvARB"))
	gpWindowPos3dvMESA = (C.GPWINDOWPOS3DVMESA)(getProcAddr("glWindowPos3dvMESA"))
	gpWindowPos3f = (C.GPWINDOWPOS3F)(getProcAddr("glWindowPos3f"))
	if gpWindowPos3f == nil {
		return errors.New("glWindowPos3f")
	}
	gpWindowPos3fARB = (C.GPWINDOWPOS3FARB)(getProcAddr("glWindowPos3fARB"))
	gpWindowPos3fMESA = (C.GPWINDOWPOS3FMESA)(getProcAddr("glWindowPos3fMESA"))
	gpWindowPos3fv = (C.GPWINDOWPOS3FV)(getProcAddr("glWindowPos3fv"))
	if gpWindowPos3fv == nil {
		return errors.New("glWindowPos3fv")
	}
	gpWindowPos3fvARB = (C.GPWINDOWPOS3FVARB)(getProcAddr("glWindowPos3fvARB"))
	gpWindowPos3fvMESA = (C.GPWINDOWPOS3FVMESA)(getProcAddr("glWindowPos3fvMESA"))
	gpWindowPos3i = (C.GPWINDOWPOS3I)(getProcAddr("glWindowPos3i"))
	if gpWindowPos3i == nil {
		return errors.New("glWindowPos3i")
	}
	gpWindowPos3iARB = (C.GPWINDOWPOS3IARB)(getProcAddr("glWindowPos3iARB"))
	gpWindowPos3iMESA = (C.GPWINDOWPOS3IMESA)(getProcAddr("glWindowPos3iMESA"))
	gpWindowPos3iv = (C.GPWINDOWPOS3IV)(getProcAddr("glWindowPos3iv"))
	if gpWindowPos3iv == nil {
		return errors.New("glWindowPos3iv")
	}
	gpWindowPos3ivARB = (C.GPWINDOWPOS3IVARB)(getProcAddr("glWindowPos3ivARB"))
	gpWindowPos3ivMESA = (C.GPWINDOWPOS3IVMESA)(getProcAddr("glWindowPos3ivMESA"))
	gpWindowPos3s = (C.GPWINDOWPOS3S)(getProcAddr("glWindowPos3s"))
	if gpWindowPos3s == nil {
		return errors.New("glWindowPos3s")
	}
	gpWindowPos3sARB = (C.GPWINDOWPOS3SARB)(getProcAddr("glWindowPos3sARB"))
	gpWindowPos3sMESA = (C.GPWINDOWPOS3SMESA)(getProcAddr("glWindowPos3sMESA"))
	gpWindowPos3sv = (C.GPWINDOWPOS3SV)(getProcAddr("glWindowPos3sv"))
	if gpWindowPos3sv == nil {
		return errors.New("glWindowPos3sv")
	}
	gpWindowPos3svARB = (C.GPWINDOWPOS3SVARB)(getProcAddr("glWindowPos3svARB"))
	gpWindowPos3svMESA = (C.GPWINDOWPOS3SVMESA)(getProcAddr("glWindowPos3svMESA"))
	gpWindowPos4dMESA = (C.GPWINDOWPOS4DMESA)(getProcAddr("glWindowPos4dMESA"))
	gpWindowPos4dvMESA = (C.GPWINDOWPOS4DVMESA)(getProcAddr("glWindowPos4dvMESA"))
	gpWindowPos4fMESA = (C.GPWINDOWPOS4FMESA)(getProcAddr("glWindowPos4fMESA"))
	gpWindowPos4fvMESA = (C.GPWINDOWPOS4FVMESA)(getProcAddr("glWindowPos4fvMESA"))
	gpWindowPos4iMESA = (C.GPWINDOWPOS4IMESA)(getProcAddr("glWindowPos4iMESA"))
	gpWindowPos4ivMESA = (C.GPWINDOWPOS4IVMESA)(getProcAddr("glWindowPos4ivMESA"))
	gpWindowPos4sMESA = (C.GPWINDOWPOS4SMESA)(getProcAddr("glWindowPos4sMESA"))
	gpWindowPos4svMESA = (C.GPWINDOWPOS4SVMESA)(getProcAddr("glWindowPos4svMESA"))
	gpWindowRectanglesEXT = (C.GPWINDOWRECTANGLESEXT)(getProcAddr("glWindowRectanglesEXT"))
	gpWriteMaskEXT = (C.GPWRITEMASKEXT)(getProcAddr("glWriteMaskEXT"))
	return nil
}
