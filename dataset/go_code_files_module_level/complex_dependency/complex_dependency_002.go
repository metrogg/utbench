func InitWithProcAddrFunc(getProcAddr func(name string) unsafe.Pointer) error {
	gpAccum = uintptr(getProcAddr("glAccum"))
	if gpAccum == 0 {
		return errors.New("glAccum")
	}
	gpAccumxOES = uintptr(getProcAddr("glAccumxOES"))
	gpAcquireKeyedMutexWin32EXT = uintptr(getProcAddr("glAcquireKeyedMutexWin32EXT"))
	gpActiveProgramEXT = uintptr(getProcAddr("glActiveProgramEXT"))
	gpActiveShaderProgram = uintptr(getProcAddr("glActiveShaderProgram"))
	gpActiveShaderProgramEXT = uintptr(getProcAddr("glActiveShaderProgramEXT"))
	gpActiveStencilFaceEXT = uintptr(getProcAddr("glActiveStencilFaceEXT"))
	gpActiveTexture = uintptr(getProcAddr("glActiveTexture"))
	if gpActiveTexture == 0 {
		return errors.New("glActiveTexture")
	}
	gpActiveTextureARB = uintptr(getProcAddr("glActiveTextureARB"))
	gpActiveVaryingNV = uintptr(getProcAddr("glActiveVaryingNV"))
	gpAlphaFragmentOp1ATI = uintptr(getProcAddr("glAlphaFragmentOp1ATI"))
	gpAlphaFragmentOp2ATI = uintptr(getProcAddr("glAlphaFragmentOp2ATI"))
	gpAlphaFragmentOp3ATI = uintptr(getProcAddr("glAlphaFragmentOp3ATI"))
	gpAlphaFunc = uintptr(getProcAddr("glAlphaFunc"))
	if gpAlphaFunc == 0 {
		return errors.New("glAlphaFunc")
	}
	gpAlphaFuncxOES = uintptr(getProcAddr("glAlphaFuncxOES"))
	gpAlphaToCoverageDitherControlNV = uintptr(getProcAddr("glAlphaToCoverageDitherControlNV"))
	gpApplyFramebufferAttachmentCMAAINTEL = uintptr(getProcAddr("glApplyFramebufferAttachmentCMAAINTEL"))
	gpApplyTextureEXT = uintptr(getProcAddr("glApplyTextureEXT"))
	gpAreProgramsResidentNV = uintptr(getProcAddr("glAreProgramsResidentNV"))
	gpAreTexturesResident = uintptr(getProcAddr("glAreTexturesResident"))
	if gpAreTexturesResident == 0 {
		return errors.New("glAreTexturesResident")
	}
	gpAreTexturesResidentEXT = uintptr(getProcAddr("glAreTexturesResidentEXT"))
	gpArrayElement = uintptr(getProcAddr("glArrayElement"))
	if gpArrayElement == 0 {
		return errors.New("glArrayElement")
	}
	gpArrayElementEXT = uintptr(getProcAddr("glArrayElementEXT"))
	gpArrayObjectATI = uintptr(getProcAddr("glArrayObjectATI"))
	gpAsyncMarkerSGIX = uintptr(getProcAddr("glAsyncMarkerSGIX"))
	gpAttachObjectARB = uintptr(getProcAddr("glAttachObjectARB"))
	gpAttachShader = uintptr(getProcAddr("glAttachShader"))
	if gpAttachShader == 0 {
		return errors.New("glAttachShader")
	}
	gpBegin = uintptr(getProcAddr("glBegin"))
	if gpBegin == 0 {
		return errors.New("glBegin")
	}
	gpBeginConditionalRenderNV = uintptr(getProcAddr("glBeginConditionalRenderNV"))
	gpBeginConditionalRenderNVX = uintptr(getProcAddr("glBeginConditionalRenderNVX"))
	gpBeginFragmentShaderATI = uintptr(getProcAddr("glBeginFragmentShaderATI"))
	gpBeginOcclusionQueryNV = uintptr(getProcAddr("glBeginOcclusionQueryNV"))
	gpBeginPerfMonitorAMD = uintptr(getProcAddr("glBeginPerfMonitorAMD"))
	gpBeginPerfQueryINTEL = uintptr(getProcAddr("glBeginPerfQueryINTEL"))
	gpBeginQuery = uintptr(getProcAddr("glBeginQuery"))
	if gpBeginQuery == 0 {
		return errors.New("glBeginQuery")
	}
	gpBeginQueryARB = uintptr(getProcAddr("glBeginQueryARB"))
	gpBeginQueryIndexed = uintptr(getProcAddr("glBeginQueryIndexed"))
	gpBeginTransformFeedbackEXT = uintptr(getProcAddr("glBeginTransformFeedbackEXT"))
	gpBeginTransformFeedbackNV = uintptr(getProcAddr("glBeginTransformFeedbackNV"))
	gpBeginVertexShaderEXT = uintptr(getProcAddr("glBeginVertexShaderEXT"))
	gpBeginVideoCaptureNV = uintptr(getProcAddr("glBeginVideoCaptureNV"))
	gpBindAttribLocation = uintptr(getProcAddr("glBindAttribLocation"))
	if gpBindAttribLocation == 0 {
		return errors.New("glBindAttribLocation")
	}
	gpBindAttribLocationARB = uintptr(getProcAddr("glBindAttribLocationARB"))
	gpBindBuffer = uintptr(getProcAddr("glBindBuffer"))
	if gpBindBuffer == 0 {
		return errors.New("glBindBuffer")
	}
	gpBindBufferARB = uintptr(getProcAddr("glBindBufferARB"))
	gpBindBufferBase = uintptr(getProcAddr("glBindBufferBase"))
	gpBindBufferBaseEXT = uintptr(getProcAddr("glBindBufferBaseEXT"))
	gpBindBufferBaseNV = uintptr(getProcAddr("glBindBufferBaseNV"))
	gpBindBufferOffsetEXT = uintptr(getProcAddr("glBindBufferOffsetEXT"))
	gpBindBufferOffsetNV = uintptr(getProcAddr("glBindBufferOffsetNV"))
	gpBindBufferRange = uintptr(getProcAddr("glBindBufferRange"))
	gpBindBufferRangeEXT = uintptr(getProcAddr("glBindBufferRangeEXT"))
	gpBindBufferRangeNV = uintptr(getProcAddr("glBindBufferRangeNV"))
	gpBindBuffersBase = uintptr(getProcAddr("glBindBuffersBase"))
	gpBindBuffersRange = uintptr(getProcAddr("glBindBuffersRange"))
	gpBindFragDataLocationEXT = uintptr(getProcAddr("glBindFragDataLocationEXT"))
	gpBindFragDataLocationIndexed = uintptr(getProcAddr("glBindFragDataLocationIndexed"))
	gpBindFragmentShaderATI = uintptr(getProcAddr("glBindFragmentShaderATI"))
	gpBindFramebuffer = uintptr(getProcAddr("glBindFramebuffer"))
	gpBindFramebufferEXT = uintptr(getProcAddr("glBindFramebufferEXT"))
	gpBindImageTexture = uintptr(getProcAddr("glBindImageTexture"))
	gpBindImageTextureEXT = uintptr(getProcAddr("glBindImageTextureEXT"))
	gpBindImageTextures = uintptr(getProcAddr("glBindImageTextures"))
	gpBindLightParameterEXT = uintptr(getProcAddr("glBindLightParameterEXT"))
	gpBindMaterialParameterEXT = uintptr(getProcAddr("glBindMaterialParameterEXT"))
	gpBindMultiTextureEXT = uintptr(getProcAddr("glBindMultiTextureEXT"))
	gpBindParameterEXT = uintptr(getProcAddr("glBindParameterEXT"))
	gpBindProgramARB = uintptr(getProcAddr("glBindProgramARB"))
	gpBindProgramNV = uintptr(getProcAddr("glBindProgramNV"))
	gpBindProgramPipeline = uintptr(getProcAddr("glBindProgramPipeline"))
	gpBindProgramPipelineEXT = uintptr(getProcAddr("glBindProgramPipelineEXT"))
	gpBindRenderbuffer = uintptr(getProcAddr("glBindRenderbuffer"))
	gpBindRenderbufferEXT = uintptr(getProcAddr("glBindRenderbufferEXT"))
	gpBindSampler = uintptr(getProcAddr("glBindSampler"))
	gpBindSamplers = uintptr(getProcAddr("glBindSamplers"))
	gpBindTexGenParameterEXT = uintptr(getProcAddr("glBindTexGenParameterEXT"))
	gpBindTexture = uintptr(getProcAddr("glBindTexture"))
	if gpBindTexture == 0 {
		return errors.New("glBindTexture")
	}
	gpBindTextureEXT = uintptr(getProcAddr("glBindTextureEXT"))
	gpBindTextureUnit = uintptr(getProcAddr("glBindTextureUnit"))
	gpBindTextureUnitParameterEXT = uintptr(getProcAddr("glBindTextureUnitParameterEXT"))
	gpBindTextures = uintptr(getProcAddr("glBindTextures"))
	gpBindTransformFeedback = uintptr(getProcAddr("glBindTransformFeedback"))
	gpBindTransformFeedbackNV = uintptr(getProcAddr("glBindTransformFeedbackNV"))
	gpBindVertexArray = uintptr(getProcAddr("glBindVertexArray"))
	gpBindVertexArrayAPPLE = uintptr(getProcAddr("glBindVertexArrayAPPLE"))
	gpBindVertexBuffer = uintptr(getProcAddr("glBindVertexBuffer"))
	gpBindVertexBuffers = uintptr(getProcAddr("glBindVertexBuffers"))
	gpBindVertexShaderEXT = uintptr(getProcAddr("glBindVertexShaderEXT"))
	gpBindVideoCaptureStreamBufferNV = uintptr(getProcAddr("glBindVideoCaptureStreamBufferNV"))
	gpBindVideoCaptureStreamTextureNV = uintptr(getProcAddr("glBindVideoCaptureStreamTextureNV"))
	gpBinormal3bEXT = uintptr(getProcAddr("glBinormal3bEXT"))
	gpBinormal3bvEXT = uintptr(getProcAddr("glBinormal3bvEXT"))
	gpBinormal3dEXT = uintptr(getProcAddr("glBinormal3dEXT"))
	gpBinormal3dvEXT = uintptr(getProcAddr("glBinormal3dvEXT"))
	gpBinormal3fEXT = uintptr(getProcAddr("glBinormal3fEXT"))
	gpBinormal3fvEXT = uintptr(getProcAddr("glBinormal3fvEXT"))
	gpBinormal3iEXT = uintptr(getProcAddr("glBinormal3iEXT"))
	gpBinormal3ivEXT = uintptr(getProcAddr("glBinormal3ivEXT"))
	gpBinormal3sEXT = uintptr(getProcAddr("glBinormal3sEXT"))
	gpBinormal3svEXT = uintptr(getProcAddr("glBinormal3svEXT"))
	gpBinormalPointerEXT = uintptr(getProcAddr("glBinormalPointerEXT"))
	gpBitmap = uintptr(getProcAddr("glBitmap"))
	if gpBitmap == 0 {
		return errors.New("glBitmap")
	}
	gpBitmapxOES = uintptr(getProcAddr("glBitmapxOES"))
	gpBlendBarrierKHR = uintptr(getProcAddr("glBlendBarrierKHR"))
	gpBlendBarrierNV = uintptr(getProcAddr("glBlendBarrierNV"))
	gpBlendColor = uintptr(getProcAddr("glBlendColor"))
	if gpBlendColor == 0 {
		return errors.New("glBlendColor")
	}
	gpBlendColorEXT = uintptr(getProcAddr("glBlendColorEXT"))
	gpBlendColorxOES = uintptr(getProcAddr("glBlendColorxOES"))
	gpBlendEquation = uintptr(getProcAddr("glBlendEquation"))
	if gpBlendEquation == 0 {
		return errors.New("glBlendEquation")
	}
	gpBlendEquationEXT = uintptr(getProcAddr("glBlendEquationEXT"))
	gpBlendEquationIndexedAMD = uintptr(getProcAddr("glBlendEquationIndexedAMD"))
	gpBlendEquationSeparate = uintptr(getProcAddr("glBlendEquationSeparate"))
	if gpBlendEquationSeparate == 0 {
		return errors.New("glBlendEquationSeparate")
	}
	gpBlendEquationSeparateEXT = uintptr(getProcAddr("glBlendEquationSeparateEXT"))
	gpBlendEquationSeparateIndexedAMD = uintptr(getProcAddr("glBlendEquationSeparateIndexedAMD"))
	gpBlendEquationSeparateiARB = uintptr(getProcAddr("glBlendEquationSeparateiARB"))
	gpBlendEquationiARB = uintptr(getProcAddr("glBlendEquationiARB"))
	gpBlendFunc = uintptr(getProcAddr("glBlendFunc"))
	if gpBlendFunc == 0 {
		return errors.New("glBlendFunc")
	}
	gpBlendFuncIndexedAMD = uintptr(getProcAddr("glBlendFuncIndexedAMD"))
	gpBlendFuncSeparate = uintptr(getProcAddr("glBlendFuncSeparate"))
	if gpBlendFuncSeparate == 0 {
		return errors.New("glBlendFuncSeparate")
	}
	gpBlendFuncSeparateEXT = uintptr(getProcAddr("glBlendFuncSeparateEXT"))
	gpBlendFuncSeparateINGR = uintptr(getProcAddr("glBlendFuncSeparateINGR"))
	gpBlendFuncSeparateIndexedAMD = uintptr(getProcAddr("glBlendFuncSeparateIndexedAMD"))
	gpBlendFuncSeparateiARB = uintptr(getProcAddr("glBlendFuncSeparateiARB"))
	gpBlendFunciARB = uintptr(getProcAddr("glBlendFunciARB"))
	gpBlendParameteriNV = uintptr(getProcAddr("glBlendParameteriNV"))
	gpBlitFramebuffer = uintptr(getProcAddr("glBlitFramebuffer"))
	gpBlitFramebufferEXT = uintptr(getProcAddr("glBlitFramebufferEXT"))
	gpBlitNamedFramebuffer = uintptr(getProcAddr("glBlitNamedFramebuffer"))
	gpBufferAddressRangeNV = uintptr(getProcAddr("glBufferAddressRangeNV"))
	gpBufferData = uintptr(getProcAddr("glBufferData"))
	if gpBufferData == 0 {
		return errors.New("glBufferData")
	}
	gpBufferDataARB = uintptr(getProcAddr("glBufferDataARB"))
	gpBufferPageCommitmentARB = uintptr(getProcAddr("glBufferPageCommitmentARB"))
	gpBufferParameteriAPPLE = uintptr(getProcAddr("glBufferParameteriAPPLE"))
	gpBufferStorage = uintptr(getProcAddr("glBufferStorage"))
	gpBufferStorageExternalEXT = uintptr(getProcAddr("glBufferStorageExternalEXT"))
	gpBufferStorageMemEXT = uintptr(getProcAddr("glBufferStorageMemEXT"))
	gpBufferSubData = uintptr(getProcAddr("glBufferSubData"))
	if gpBufferSubData == 0 {
		return errors.New("glBufferSubData")
	}
	gpBufferSubDataARB = uintptr(getProcAddr("glBufferSubDataARB"))
	gpCallCommandListNV = uintptr(getProcAddr("glCallCommandListNV"))
	gpCallList = uintptr(getProcAddr("glCallList"))
	if gpCallList == 0 {
		return errors.New("glCallList")
	}
	gpCallLists = uintptr(getProcAddr("glCallLists"))
	if gpCallLists == 0 {
		return errors.New("glCallLists")
	}
	gpCheckFramebufferStatus = uintptr(getProcAddr("glCheckFramebufferStatus"))
	gpCheckFramebufferStatusEXT = uintptr(getProcAddr("glCheckFramebufferStatusEXT"))
	gpCheckNamedFramebufferStatus = uintptr(getProcAddr("glCheckNamedFramebufferStatus"))
	gpCheckNamedFramebufferStatusEXT = uintptr(getProcAddr("glCheckNamedFramebufferStatusEXT"))
	gpClampColorARB = uintptr(getProcAddr("glClampColorARB"))
	gpClear = uintptr(getProcAddr("glClear"))
	if gpClear == 0 {
		return errors.New("glClear")
	}
	gpClearAccum = uintptr(getProcAddr("glClearAccum"))
	if gpClearAccum == 0 {
		return errors.New("glClearAccum")
	}
	gpClearAccumxOES = uintptr(getProcAddr("glClearAccumxOES"))
	gpClearBufferData = uintptr(getProcAddr("glClearBufferData"))
	gpClearBufferSubData = uintptr(getProcAddr("glClearBufferSubData"))
	gpClearColor = uintptr(getProcAddr("glClearColor"))
	if gpClearColor == 0 {
		return errors.New("glClearColor")
	}
	gpClearColorIiEXT = uintptr(getProcAddr("glClearColorIiEXT"))
	gpClearColorIuiEXT = uintptr(getProcAddr("glClearColorIuiEXT"))
	gpClearColorxOES = uintptr(getProcAddr("glClearColorxOES"))
	gpClearDepth = uintptr(getProcAddr("glClearDepth"))
	if gpClearDepth == 0 {
		return errors.New("glClearDepth")
	}
	gpClearDepthdNV = uintptr(getProcAddr("glClearDepthdNV"))
	gpClearDepthf = uintptr(getProcAddr("glClearDepthf"))
	gpClearDepthfOES = uintptr(getProcAddr("glClearDepthfOES"))
	gpClearDepthxOES = uintptr(getProcAddr("glClearDepthxOES"))
	gpClearIndex = uintptr(getProcAddr("glClearIndex"))
	if gpClearIndex == 0 {
		return errors.New("glClearIndex")
	}
	gpClearNamedBufferData = uintptr(getProcAddr("glClearNamedBufferData"))
	gpClearNamedBufferDataEXT = uintptr(getProcAddr("glClearNamedBufferDataEXT"))
	gpClearNamedBufferSubData = uintptr(getProcAddr("glClearNamedBufferSubData"))
	gpClearNamedBufferSubDataEXT = uintptr(getProcAddr("glClearNamedBufferSubDataEXT"))
	gpClearNamedFramebufferfi = uintptr(getProcAddr("glClearNamedFramebufferfi"))
	gpClearNamedFramebufferfv = uintptr(getProcAddr("glClearNamedFramebufferfv"))
	gpClearNamedFramebufferiv = uintptr(getProcAddr("glClearNamedFramebufferiv"))
	gpClearNamedFramebufferuiv = uintptr(getProcAddr("glClearNamedFramebufferuiv"))
	gpClearStencil = uintptr(getProcAddr("glClearStencil"))
	if gpClearStencil == 0 {
		return errors.New("glClearStencil")
	}
	gpClearTexImage = uintptr(getProcAddr("glClearTexImage"))
	gpClearTexSubImage = uintptr(getProcAddr("glClearTexSubImage"))
	gpClientActiveTexture = uintptr(getProcAddr("glClientActiveTexture"))
	if gpClientActiveTexture == 0 {
		return errors.New("glClientActiveTexture")
	}
	gpClientActiveTextureARB = uintptr(getProcAddr("glClientActiveTextureARB"))
	gpClientActiveVertexStreamATI = uintptr(getProcAddr("glClientActiveVertexStreamATI"))
	gpClientAttribDefaultEXT = uintptr(getProcAddr("glClientAttribDefaultEXT"))
	gpClientWaitSync = uintptr(getProcAddr("glClientWaitSync"))
	gpClipControl = uintptr(getProcAddr("glClipControl"))
	gpClipPlane = uintptr(getProcAddr("glClipPlane"))
	if gpClipPlane == 0 {
		return errors.New("glClipPlane")
	}
	gpClipPlanefOES = uintptr(getProcAddr("glClipPlanefOES"))
	gpClipPlanexOES = uintptr(getProcAddr("glClipPlanexOES"))
	gpColor3b = uintptr(getProcAddr("glColor3b"))
	if gpColor3b == 0 {
		return errors.New("glColor3b")
	}
	gpColor3bv = uintptr(getProcAddr("glColor3bv"))
	if gpColor3bv == 0 {
		return errors.New("glColor3bv")
	}
	gpColor3d = uintptr(getProcAddr("glColor3d"))
	if gpColor3d == 0 {
		return errors.New("glColor3d")
	}
	gpColor3dv = uintptr(getProcAddr("glColor3dv"))
	if gpColor3dv == 0 {
		return errors.New("glColor3dv")
	}
	gpColor3f = uintptr(getProcAddr("glColor3f"))
	if gpColor3f == 0 {
		return errors.New("glColor3f")
	}
	gpColor3fVertex3fSUN = uintptr(getProcAddr("glColor3fVertex3fSUN"))
	gpColor3fVertex3fvSUN = uintptr(getProcAddr("glColor3fVertex3fvSUN"))
	gpColor3fv = uintptr(getProcAddr("glColor3fv"))
	if gpColor3fv == 0 {
		return errors.New("glColor3fv")
	}
	gpColor3hNV = uintptr(getProcAddr("glColor3hNV"))
	gpColor3hvNV = uintptr(getProcAddr("glColor3hvNV"))
	gpColor3i = uintptr(getProcAddr("glColor3i"))
	if gpColor3i == 0 {
		return errors.New("glColor3i")
	}
	gpColor3iv = uintptr(getProcAddr("glColor3iv"))
	if gpColor3iv == 0 {
		return errors.New("glColor3iv")
	}
	gpColor3s = uintptr(getProcAddr("glColor3s"))
	if gpColor3s == 0 {
		return errors.New("glColor3s")
	}
	gpColor3sv = uintptr(getProcAddr("glColor3sv"))
	if gpColor3sv == 0 {
		return errors.New("glColor3sv")
	}
	gpColor3ub = uintptr(getProcAddr("glColor3ub"))
	if gpColor3ub == 0 {
		return errors.New("glColor3ub")
	}
	gpColor3ubv = uintptr(getProcAddr("glColor3ubv"))
	if gpColor3ubv == 0 {
		return errors.New("glColor3ubv")
	}
	gpColor3ui = uintptr(getProcAddr("glColor3ui"))
	if gpColor3ui == 0 {
		return errors.New("glColor3ui")
	}
	gpColor3uiv = uintptr(getProcAddr("glColor3uiv"))
	if gpColor3uiv == 0 {
		return errors.New("glColor3uiv")
	}
	gpColor3us = uintptr(getProcAddr("glColor3us"))
	if gpColor3us == 0 {
		return errors.New("glColor3us")
	}
	gpColor3usv = uintptr(getProcAddr("glColor3usv"))
	if gpColor3usv == 0 {
		return errors.New("glColor3usv")
	}
	gpColor3xOES = uintptr(getProcAddr("glColor3xOES"))
	gpColor3xvOES = uintptr(getProcAddr("glColor3xvOES"))
	gpColor4b = uintptr(getProcAddr("glColor4b"))
	if gpColor4b == 0 {
		return errors.New("glColor4b")
	}
	gpColor4bv = uintptr(getProcAddr("glColor4bv"))
	if gpColor4bv == 0 {
		return errors.New("glColor4bv")
	}
	gpColor4d = uintptr(getProcAddr("glColor4d"))
	if gpColor4d == 0 {
		return errors.New("glColor4d")
	}
	gpColor4dv = uintptr(getProcAddr("glColor4dv"))
	if gpColor4dv == 0 {
		return errors.New("glColor4dv")
	}
	gpColor4f = uintptr(getProcAddr("glColor4f"))
	if gpColor4f == 0 {
		return errors.New("glColor4f")
	}
	gpColor4fNormal3fVertex3fSUN = uintptr(getProcAddr("glColor4fNormal3fVertex3fSUN"))
	gpColor4fNormal3fVertex3fvSUN = uintptr(getProcAddr("glColor4fNormal3fVertex3fvSUN"))
	gpColor4fv = uintptr(getProcAddr("glColor4fv"))
	if gpColor4fv == 0 {
		return errors.New("glColor4fv")
	}
	gpColor4hNV = uintptr(getProcAddr("glColor4hNV"))
	gpColor4hvNV = uintptr(getProcAddr("glColor4hvNV"))
	gpColor4i = uintptr(getProcAddr("glColor4i"))
	if gpColor4i == 0 {
		return errors.New("glColor4i")
	}
	gpColor4iv = uintptr(getProcAddr("glColor4iv"))
	if gpColor4iv == 0 {
		return errors.New("glColor4iv")
	}
	gpColor4s = uintptr(getProcAddr("glColor4s"))
	if gpColor4s == 0 {
		return errors.New("glColor4s")
	}
	gpColor4sv = uintptr(getProcAddr("glColor4sv"))
	if gpColor4sv == 0 {
		return errors.New("glColor4sv")
	}
	gpColor4ub = uintptr(getProcAddr("glColor4ub"))
	if gpColor4ub == 0 {
		return errors.New("glColor4ub")
	}
	gpColor4ubVertex2fSUN = uintptr(getProcAddr("glColor4ubVertex2fSUN"))
	gpColor4ubVertex2fvSUN = uintptr(getProcAddr("glColor4ubVertex2fvSUN"))
	gpColor4ubVertex3fSUN = uintptr(getProcAddr("glColor4ubVertex3fSUN"))
	gpColor4ubVertex3fvSUN = uintptr(getProcAddr("glColor4ubVertex3fvSUN"))
	gpColor4ubv = uintptr(getProcAddr("glColor4ubv"))
	if gpColor4ubv == 0 {
		return errors.New("glColor4ubv")
	}
	gpColor4ui = uintptr(getProcAddr("glColor4ui"))
	if gpColor4ui == 0 {
		return errors.New("glColor4ui")
	}
	gpColor4uiv = uintptr(getProcAddr("glColor4uiv"))
	if gpColor4uiv == 0 {
		return errors.New("glColor4uiv")
	}
	gpColor4us = uintptr(getProcAddr("glColor4us"))
	if gpColor4us == 0 {
		return errors.New("glColor4us")
	}
	gpColor4usv = uintptr(getProcAddr("glColor4usv"))
	if gpColor4usv == 0 {
		return errors.New("glColor4usv")
	}
	gpColor4xOES = uintptr(getProcAddr("glColor4xOES"))
	gpColor4xvOES = uintptr(getProcAddr("glColor4xvOES"))
	gpColorFormatNV = uintptr(getProcAddr("glColorFormatNV"))
	gpColorFragmentOp1ATI = uintptr(getProcAddr("glColorFragmentOp1ATI"))
	gpColorFragmentOp2ATI = uintptr(getProcAddr("glColorFragmentOp2ATI"))
	gpColorFragmentOp3ATI = uintptr(getProcAddr("glColorFragmentOp3ATI"))
	gpColorMask = uintptr(getProcAddr("glColorMask"))
	if gpColorMask == 0 {
		return errors.New("glColorMask")
	}
	gpColorMaskIndexedEXT = uintptr(getProcAddr("glColorMaskIndexedEXT"))
	gpColorMaterial = uintptr(getProcAddr("glColorMaterial"))
	if gpColorMaterial == 0 {
		return errors.New("glColorMaterial")
	}
	gpColorPointer = uintptr(getProcAddr("glColorPointer"))
	if gpColorPointer == 0 {
		return errors.New("glColorPointer")
	}
	gpColorPointerEXT = uintptr(getProcAddr("glColorPointerEXT"))
	gpColorPointerListIBM = uintptr(getProcAddr("glColorPointerListIBM"))
	gpColorPointervINTEL = uintptr(getProcAddr("glColorPointervINTEL"))
	gpColorSubTableEXT = uintptr(getProcAddr("glColorSubTableEXT"))
	gpColorTableEXT = uintptr(getProcAddr("glColorTableEXT"))
	gpColorTableParameterfvSGI = uintptr(getProcAddr("glColorTableParameterfvSGI"))
	gpColorTableParameterivSGI = uintptr(getProcAddr("glColorTableParameterivSGI"))
	gpColorTableSGI = uintptr(getProcAddr("glColorTableSGI"))
	gpCombinerInputNV = uintptr(getProcAddr("glCombinerInputNV"))
	gpCombinerOutputNV = uintptr(getProcAddr("glCombinerOutputNV"))
	gpCombinerParameterfNV = uintptr(getProcAddr("glCombinerParameterfNV"))
	gpCombinerParameterfvNV = uintptr(getProcAddr("glCombinerParameterfvNV"))
	gpCombinerParameteriNV = uintptr(getProcAddr("glCombinerParameteriNV"))
	gpCombinerParameterivNV = uintptr(getProcAddr("glCombinerParameterivNV"))
	gpCombinerStageParameterfvNV = uintptr(getProcAddr("glCombinerStageParameterfvNV"))
	gpCommandListSegmentsNV = uintptr(getProcAddr("glCommandListSegmentsNV"))
	gpCompileCommandListNV = uintptr(getProcAddr("glCompileCommandListNV"))
	gpCompileShader = uintptr(getProcAddr("glCompileShader"))
	if gpCompileShader == 0 {
		return errors.New("glCompileShader")
	}
	gpCompileShaderARB = uintptr(getProcAddr("glCompileShaderARB"))
	gpCompileShaderIncludeARB = uintptr(getProcAddr("glCompileShaderIncludeARB"))
	gpCompressedMultiTexImage1DEXT = uintptr(getProcAddr("glCompressedMultiTexImage1DEXT"))
	gpCompressedMultiTexImage2DEXT = uintptr(getProcAddr("glCompressedMultiTexImage2DEXT"))
	gpCompressedMultiTexImage3DEXT = uintptr(getProcAddr("glCompressedMultiTexImage3DEXT"))
	gpCompressedMultiTexSubImage1DEXT = uintptr(getProcAddr("glCompressedMultiTexSubImage1DEXT"))
	gpCompressedMultiTexSubImage2DEXT = uintptr(getProcAddr("glCompressedMultiTexSubImage2DEXT"))
	gpCompressedMultiTexSubImage3DEXT = uintptr(getProcAddr("glCompressedMultiTexSubImage3DEXT"))
	gpCompressedTexImage1D = uintptr(getProcAddr("glCompressedTexImage1D"))
	if gpCompressedTexImage1D == 0 {
		return errors.New("glCompressedTexImage1D")
	}
	gpCompressedTexImage1DARB = uintptr(getProcAddr("glCompressedTexImage1DARB"))
	gpCompressedTexImage2D = uintptr(getProcAddr("glCompressedTexImage2D"))
	if gpCompressedTexImage2D == 0 {
		return errors.New("glCompressedTexImage2D")
	}
	gpCompressedTexImage2DARB = uintptr(getProcAddr("glCompressedTexImage2DARB"))
	gpCompressedTexImage3D = uintptr(getProcAddr("glCompressedTexImage3D"))
	if gpCompressedTexImage3D == 0 {
		return errors.New("glCompressedTexImage3D")
	}
	gpCompressedTexImage3DARB = uintptr(getProcAddr("glCompressedTexImage3DARB"))
	gpCompressedTexSubImage1D = uintptr(getProcAddr("glCompressedTexSubImage1D"))
	if gpCompressedTexSubImage1D == 0 {
		return errors.New("glCompressedTexSubImage1D")
	}
	gpCompressedTexSubImage1DARB = uintptr(getProcAddr("glCompressedTexSubImage1DARB"))
	gpCompressedTexSubImage2D = uintptr(getProcAddr("glCompressedTexSubImage2D"))
	if gpCompressedTexSubImage2D == 0 {
		return errors.New("glCompressedTexSubImage2D")
	}
	gpCompressedTexSubImage2DARB = uintptr(getProcAddr("glCompressedTexSubImage2DARB"))
	gpCompressedTexSubImage3D = uintptr(getProcAddr("glCompressedTexSubImage3D"))
	if gpCompressedTexSubImage3D == 0 {
		return errors.New("glCompressedTexSubImage3D")
	}
	gpCompressedTexSubImage3DARB = uintptr(getProcAddr("glCompressedTexSubImage3DARB"))
	gpCompressedTextureImage1DEXT = uintptr(getProcAddr("glCompressedTextureImage1DEXT"))
	gpCompressedTextureImage2DEXT = uintptr(getProcAddr("glCompressedTextureImage2DEXT"))
	gpCompressedTextureImage3DEXT = uintptr(getProcAddr("glCompressedTextureImage3DEXT"))
	gpCompressedTextureSubImage1D = uintptr(getProcAddr("glCompressedTextureSubImage1D"))
	gpCompressedTextureSubImage1DEXT = uintptr(getProcAddr("glCompressedTextureSubImage1DEXT"))
	gpCompressedTextureSubImage2D = uintptr(getProcAddr("glCompressedTextureSubImage2D"))
	gpCompressedTextureSubImage2DEXT = uintptr(getProcAddr("glCompressedTextureSubImage2DEXT"))
	gpCompressedTextureSubImage3D = uintptr(getProcAddr("glCompressedTextureSubImage3D"))
	gpCompressedTextureSubImage3DEXT = uintptr(getProcAddr("glCompressedTextureSubImage3DEXT"))
	gpConservativeRasterParameterfNV = uintptr(getProcAddr("glConservativeRasterParameterfNV"))
	gpConservativeRasterParameteriNV = uintptr(getProcAddr("glConservativeRasterParameteriNV"))
	gpConvolutionFilter1DEXT = uintptr(getProcAddr("glConvolutionFilter1DEXT"))
	gpConvolutionFilter2DEXT = uintptr(getProcAddr("glConvolutionFilter2DEXT"))
	gpConvolutionParameterfEXT = uintptr(getProcAddr("glConvolutionParameterfEXT"))
	gpConvolutionParameterfvEXT = uintptr(getProcAddr("glConvolutionParameterfvEXT"))
	gpConvolutionParameteriEXT = uintptr(getProcAddr("glConvolutionParameteriEXT"))
	gpConvolutionParameterivEXT = uintptr(getProcAddr("glConvolutionParameterivEXT"))
	gpConvolutionParameterxOES = uintptr(getProcAddr("glConvolutionParameterxOES"))
	gpConvolutionParameterxvOES = uintptr(getProcAddr("glConvolutionParameterxvOES"))
	gpCopyBufferSubData = uintptr(getProcAddr("glCopyBufferSubData"))
	gpCopyColorSubTableEXT = uintptr(getProcAddr("glCopyColorSubTableEXT"))
	gpCopyColorTableSGI = uintptr(getProcAddr("glCopyColorTableSGI"))
	gpCopyConvolutionFilter1DEXT = uintptr(getProcAddr("glCopyConvolutionFilter1DEXT"))
	gpCopyConvolutionFilter2DEXT = uintptr(getProcAddr("glCopyConvolutionFilter2DEXT"))
	gpCopyImageSubData = uintptr(getProcAddr("glCopyImageSubData"))
	gpCopyImageSubDataNV = uintptr(getProcAddr("glCopyImageSubDataNV"))
	gpCopyMultiTexImage1DEXT = uintptr(getProcAddr("glCopyMultiTexImage1DEXT"))
	gpCopyMultiTexImage2DEXT = uintptr(getProcAddr("glCopyMultiTexImage2DEXT"))
	gpCopyMultiTexSubImage1DEXT = uintptr(getProcAddr("glCopyMultiTexSubImage1DEXT"))
	gpCopyMultiTexSubImage2DEXT = uintptr(getProcAddr("glCopyMultiTexSubImage2DEXT"))
	gpCopyMultiTexSubImage3DEXT = uintptr(getProcAddr("glCopyMultiTexSubImage3DEXT"))
	gpCopyNamedBufferSubData = uintptr(getProcAddr("glCopyNamedBufferSubData"))
	gpCopyPathNV = uintptr(getProcAddr("glCopyPathNV"))
	gpCopyPixels = uintptr(getProcAddr("glCopyPixels"))
	if gpCopyPixels == 0 {
		return errors.New("glCopyPixels")
	}
	gpCopyTexImage1D = uintptr(getProcAddr("glCopyTexImage1D"))
	if gpCopyTexImage1D == 0 {
		return errors.New("glCopyTexImage1D")
	}
	gpCopyTexImage1DEXT = uintptr(getProcAddr("glCopyTexImage1DEXT"))
	gpCopyTexImage2D = uintptr(getProcAddr("glCopyTexImage2D"))
	if gpCopyTexImage2D == 0 {
		return errors.New("glCopyTexImage2D")
	}
	gpCopyTexImage2DEXT = uintptr(getProcAddr("glCopyTexImage2DEXT"))
	gpCopyTexSubImage1D = uintptr(getProcAddr("glCopyTexSubImage1D"))
	if gpCopyTexSubImage1D == 0 {
		return errors.New("glCopyTexSubImage1D")
	}
	gpCopyTexSubImage1DEXT = uintptr(getProcAddr("glCopyTexSubImage1DEXT"))
	gpCopyTexSubImage2D = uintptr(getProcAddr("glCopyTexSubImage2D"))
	if gpCopyTexSubImage2D == 0 {
		return errors.New("glCopyTexSubImage2D")
	}
	gpCopyTexSubImage2DEXT = uintptr(getProcAddr("glCopyTexSubImage2DEXT"))
	gpCopyTexSubImage3D = uintptr(getProcAddr("glCopyTexSubImage3D"))
	if gpCopyTexSubImage3D == 0 {
		return errors.New("glCopyTexSubImage3D")
	}
	gpCopyTexSubImage3DEXT = uintptr(getProcAddr("glCopyTexSubImage3DEXT"))
	gpCopyTextureImage1DEXT = uintptr(getProcAddr("glCopyTextureImage1DEXT"))
	gpCopyTextureImage2DEXT = uintptr(getProcAddr("glCopyTextureImage2DEXT"))
	gpCopyTextureSubImage1D = uintptr(getProcAddr("glCopyTextureSubImage1D"))
	gpCopyTextureSubImage1DEXT = uintptr(getProcAddr("glCopyTextureSubImage1DEXT"))
	gpCopyTextureSubImage2D = uintptr(getProcAddr("glCopyTextureSubImage2D"))
	gpCopyTextureSubImage2DEXT = uintptr(getProcAddr("glCopyTextureSubImage2DEXT"))
	gpCopyTextureSubImage3D = uintptr(getProcAddr("glCopyTextureSubImage3D"))
	gpCopyTextureSubImage3DEXT = uintptr(getProcAddr("glCopyTextureSubImage3DEXT"))
	gpCoverFillPathInstancedNV = uintptr(getProcAddr("glCoverFillPathInstancedNV"))
	gpCoverFillPathNV = uintptr(getProcAddr("glCoverFillPathNV"))
	gpCoverStrokePathInstancedNV = uintptr(getProcAddr("glCoverStrokePathInstancedNV"))
	gpCoverStrokePathNV = uintptr(getProcAddr("glCoverStrokePathNV"))
	gpCoverageModulationNV = uintptr(getProcAddr("glCoverageModulationNV"))
	gpCoverageModulationTableNV = uintptr(getProcAddr("glCoverageModulationTableNV"))
	gpCreateBuffers = uintptr(getProcAddr("glCreateBuffers"))
	gpCreateCommandListsNV = uintptr(getProcAddr("glCreateCommandListsNV"))
	gpCreateFramebuffers = uintptr(getProcAddr("glCreateFramebuffers"))
	gpCreateMemoryObjectsEXT = uintptr(getProcAddr("glCreateMemoryObjectsEXT"))
	gpCreatePerfQueryINTEL = uintptr(getProcAddr("glCreatePerfQueryINTEL"))
	gpCreateProgram = uintptr(getProcAddr("glCreateProgram"))
	if gpCreateProgram == 0 {
		return errors.New("glCreateProgram")
	}
	gpCreateProgramObjectARB = uintptr(getProcAddr("glCreateProgramObjectARB"))
	gpCreateProgramPipelines = uintptr(getProcAddr("glCreateProgramPipelines"))
	gpCreateQueries = uintptr(getProcAddr("glCreateQueries"))
	gpCreateRenderbuffers = uintptr(getProcAddr("glCreateRenderbuffers"))
	gpCreateSamplers = uintptr(getProcAddr("glCreateSamplers"))
	gpCreateShader = uintptr(getProcAddr("glCreateShader"))
	if gpCreateShader == 0 {
		return errors.New("glCreateShader")
	}
	gpCreateShaderObjectARB = uintptr(getProcAddr("glCreateShaderObjectARB"))
	gpCreateShaderProgramEXT = uintptr(getProcAddr("glCreateShaderProgramEXT"))
	gpCreateShaderProgramv = uintptr(getProcAddr("glCreateShaderProgramv"))
	gpCreateShaderProgramvEXT = uintptr(getProcAddr("glCreateShaderProgramvEXT"))
	gpCreateStatesNV = uintptr(getProcAddr("glCreateStatesNV"))
	gpCreateSyncFromCLeventARB = uintptr(getProcAddr("glCreateSyncFromCLeventARB"))
	gpCreateTextures = uintptr(getProcAddr("glCreateTextures"))
	gpCreateTransformFeedbacks = uintptr(getProcAddr("glCreateTransformFeedbacks"))
	gpCreateVertexArrays = uintptr(getProcAddr("glCreateVertexArrays"))
	gpCullFace = uintptr(getProcAddr("glCullFace"))
	if gpCullFace == 0 {
		return errors.New("glCullFace")
	}
	gpCullParameterdvEXT = uintptr(getProcAddr("glCullParameterdvEXT"))
	gpCullParameterfvEXT = uintptr(getProcAddr("glCullParameterfvEXT"))
	gpCurrentPaletteMatrixARB = uintptr(getProcAddr("glCurrentPaletteMatrixARB"))
	gpDebugMessageCallback = uintptr(getProcAddr("glDebugMessageCallback"))
	gpDebugMessageCallbackAMD = uintptr(getProcAddr("glDebugMessageCallbackAMD"))
	gpDebugMessageCallbackARB = uintptr(getProcAddr("glDebugMessageCallbackARB"))
	gpDebugMessageCallbackKHR = uintptr(getProcAddr("glDebugMessageCallbackKHR"))
	gpDebugMessageControl = uintptr(getProcAddr("glDebugMessageControl"))
	gpDebugMessageControlARB = uintptr(getProcAddr("glDebugMessageControlARB"))
	gpDebugMessageControlKHR = uintptr(getProcAddr("glDebugMessageControlKHR"))
	gpDebugMessageEnableAMD = uintptr(getProcAddr("glDebugMessageEnableAMD"))
	gpDebugMessageInsert = uintptr(getProcAddr("glDebugMessageInsert"))
	gpDebugMessageInsertAMD = uintptr(getProcAddr("glDebugMessageInsertAMD"))
	gpDebugMessageInsertARB = uintptr(getProcAddr("glDebugMessageInsertARB"))
	gpDebugMessageInsertKHR = uintptr(getProcAddr("glDebugMessageInsertKHR"))
	gpDeformSGIX = uintptr(getProcAddr("glDeformSGIX"))
	gpDeformationMap3dSGIX = uintptr(getProcAddr("glDeformationMap3dSGIX"))
	gpDeformationMap3fSGIX = uintptr(getProcAddr("glDeformationMap3fSGIX"))
	gpDeleteAsyncMarkersSGIX = uintptr(getProcAddr("glDeleteAsyncMarkersSGIX"))
	gpDeleteBuffers = uintptr(getProcAddr("glDeleteBuffers"))
	if gpDeleteBuffers == 0 {
		return errors.New("glDeleteBuffers")
	}
	gpDeleteBuffersARB = uintptr(getProcAddr("glDeleteBuffersARB"))
	gpDeleteCommandListsNV = uintptr(getProcAddr("glDeleteCommandListsNV"))
	gpDeleteFencesAPPLE = uintptr(getProcAddr("glDeleteFencesAPPLE"))
	gpDeleteFencesNV = uintptr(getProcAddr("glDeleteFencesNV"))
	gpDeleteFragmentShaderATI = uintptr(getProcAddr("glDeleteFragmentShaderATI"))
	gpDeleteFramebuffers = uintptr(getProcAddr("glDeleteFramebuffers"))
	gpDeleteFramebuffersEXT = uintptr(getProcAddr("glDeleteFramebuffersEXT"))
	gpDeleteLists = uintptr(getProcAddr("glDeleteLists"))
	if gpDeleteLists == 0 {
		return errors.New("glDeleteLists")
	}
	gpDeleteMemoryObjectsEXT = uintptr(getProcAddr("glDeleteMemoryObjectsEXT"))
	gpDeleteNamedStringARB = uintptr(getProcAddr("glDeleteNamedStringARB"))
	gpDeleteNamesAMD = uintptr(getProcAddr("glDeleteNamesAMD"))
	gpDeleteObjectARB = uintptr(getProcAddr("glDeleteObjectARB"))
	gpDeleteOcclusionQueriesNV = uintptr(getProcAddr("glDeleteOcclusionQueriesNV"))
	gpDeletePathsNV = uintptr(getProcAddr("glDeletePathsNV"))
	gpDeletePerfMonitorsAMD = uintptr(getProcAddr("glDeletePerfMonitorsAMD"))
	gpDeletePerfQueryINTEL = uintptr(getProcAddr("glDeletePerfQueryINTEL"))
	gpDeleteProgram = uintptr(getProcAddr("glDeleteProgram"))
	if gpDeleteProgram == 0 {
		return errors.New("glDeleteProgram")
	}
	gpDeleteProgramPipelines = uintptr(getProcAddr("glDeleteProgramPipelines"))
	gpDeleteProgramPipelinesEXT = uintptr(getProcAddr("glDeleteProgramPipelinesEXT"))
	gpDeleteProgramsARB = uintptr(getProcAddr("glDeleteProgramsARB"))
	gpDeleteProgramsNV = uintptr(getProcAddr("glDeleteProgramsNV"))
	gpDeleteQueries = uintptr(getProcAddr("glDeleteQueries"))
	if gpDeleteQueries == 0 {
		return errors.New("glDeleteQueries")
	}
	gpDeleteQueriesARB = uintptr(getProcAddr("glDeleteQueriesARB"))
	gpDeleteQueryResourceTagNV = uintptr(getProcAddr("glDeleteQueryResourceTagNV"))
	gpDeleteRenderbuffers = uintptr(getProcAddr("glDeleteRenderbuffers"))
	gpDeleteRenderbuffersEXT = uintptr(getProcAddr("glDeleteRenderbuffersEXT"))
	gpDeleteSamplers = uintptr(getProcAddr("glDeleteSamplers"))
	gpDeleteSemaphoresEXT = uintptr(getProcAddr("glDeleteSemaphoresEXT"))
	gpDeleteShader = uintptr(getProcAddr("glDeleteShader"))
	if gpDeleteShader == 0 {
		return errors.New("glDeleteShader")
	}
	gpDeleteStatesNV = uintptr(getProcAddr("glDeleteStatesNV"))
	gpDeleteSync = uintptr(getProcAddr("glDeleteSync"))
	gpDeleteTextures = uintptr(getProcAddr("glDeleteTextures"))
	if gpDeleteTextures == 0 {
		return errors.New("glDeleteTextures")
	}
	gpDeleteTexturesEXT = uintptr(getProcAddr("glDeleteTexturesEXT"))
	gpDeleteTransformFeedbacks = uintptr(getProcAddr("glDeleteTransformFeedbacks"))
	gpDeleteTransformFeedbacksNV = uintptr(getProcAddr("glDeleteTransformFeedbacksNV"))
	gpDeleteVertexArrays = uintptr(getProcAddr("glDeleteVertexArrays"))
	gpDeleteVertexArraysAPPLE = uintptr(getProcAddr("glDeleteVertexArraysAPPLE"))
	gpDeleteVertexShaderEXT = uintptr(getProcAddr("glDeleteVertexShaderEXT"))
	gpDepthBoundsEXT = uintptr(getProcAddr("glDepthBoundsEXT"))
	gpDepthBoundsdNV = uintptr(getProcAddr("glDepthBoundsdNV"))
	gpDepthFunc = uintptr(getProcAddr("glDepthFunc"))
	if gpDepthFunc == 0 {
		return errors.New("glDepthFunc")
	}
	gpDepthMask = uintptr(getProcAddr("glDepthMask"))
	if gpDepthMask == 0 {
		return errors.New("glDepthMask")
	}
	gpDepthRange = uintptr(getProcAddr("glDepthRange"))
	if gpDepthRange == 0 {
		return errors.New("glDepthRange")
	}
	gpDepthRangeArrayv = uintptr(getProcAddr("glDepthRangeArrayv"))
	gpDepthRangeIndexed = uintptr(getProcAddr("glDepthRangeIndexed"))
	gpDepthRangedNV = uintptr(getProcAddr("glDepthRangedNV"))
	gpDepthRangef = uintptr(getProcAddr("glDepthRangef"))
	gpDepthRangefOES = uintptr(getProcAddr("glDepthRangefOES"))
	gpDepthRangexOES = uintptr(getProcAddr("glDepthRangexOES"))
	gpDetachObjectARB = uintptr(getProcAddr("glDetachObjectARB"))
	gpDetachShader = uintptr(getProcAddr("glDetachShader"))
	if gpDetachShader == 0 {
		return errors.New("glDetachShader")
	}
	gpDetailTexFuncSGIS = uintptr(getProcAddr("glDetailTexFuncSGIS"))
	gpDisable = uintptr(getProcAddr("glDisable"))
	if gpDisable == 0 {
		return errors.New("glDisable")
	}
	gpDisableClientState = uintptr(getProcAddr("glDisableClientState"))
	if gpDisableClientState == 0 {
		return errors.New("glDisableClientState")
	}
	gpDisableClientStateIndexedEXT = uintptr(getProcAddr("glDisableClientStateIndexedEXT"))
	gpDisableClientStateiEXT = uintptr(getProcAddr("glDisableClientStateiEXT"))
	gpDisableIndexedEXT = uintptr(getProcAddr("glDisableIndexedEXT"))
	gpDisableVariantClientStateEXT = uintptr(getProcAddr("glDisableVariantClientStateEXT"))
	gpDisableVertexArrayAttrib = uintptr(getProcAddr("glDisableVertexArrayAttrib"))
	gpDisableVertexArrayAttribEXT = uintptr(getProcAddr("glDisableVertexArrayAttribEXT"))
	gpDisableVertexArrayEXT = uintptr(getProcAddr("glDisableVertexArrayEXT"))
	gpDisableVertexAttribAPPLE = uintptr(getProcAddr("glDisableVertexAttribAPPLE"))
	gpDisableVertexAttribArray = uintptr(getProcAddr("glDisableVertexAttribArray"))
	if gpDisableVertexAttribArray == 0 {
		return errors.New("glDisableVertexAttribArray")
	}
	gpDisableVertexAttribArrayARB = uintptr(getProcAddr("glDisableVertexAttribArrayARB"))
	gpDispatchCompute = uintptr(getProcAddr("glDispatchCompute"))
	gpDispatchComputeGroupSizeARB = uintptr(getProcAddr("glDispatchComputeGroupSizeARB"))
	gpDispatchComputeIndirect = uintptr(getProcAddr("glDispatchComputeIndirect"))
	gpDrawArrays = uintptr(getProcAddr("glDrawArrays"))
	if gpDrawArrays == 0 {
		return errors.New("glDrawArrays")
	}
	gpDrawArraysEXT = uintptr(getProcAddr("glDrawArraysEXT"))
	gpDrawArraysIndirect = uintptr(getProcAddr("glDrawArraysIndirect"))
	gpDrawArraysInstancedARB = uintptr(getProcAddr("glDrawArraysInstancedARB"))
	gpDrawArraysInstancedBaseInstance = uintptr(getProcAddr("glDrawArraysInstancedBaseInstance"))
	gpDrawArraysInstancedEXT = uintptr(getProcAddr("glDrawArraysInstancedEXT"))
	gpDrawBuffer = uintptr(getProcAddr("glDrawBuffer"))
	if gpDrawBuffer == 0 {
		return errors.New("glDrawBuffer")
	}
	gpDrawBuffers = uintptr(getProcAddr("glDrawBuffers"))
	if gpDrawBuffers == 0 {
		return errors.New("glDrawBuffers")
	}
	gpDrawBuffersARB = uintptr(getProcAddr("glDrawBuffersARB"))
	gpDrawBuffersATI = uintptr(getProcAddr("glDrawBuffersATI"))
	gpDrawCommandsAddressNV = uintptr(getProcAddr("glDrawCommandsAddressNV"))
	gpDrawCommandsNV = uintptr(getProcAddr("glDrawCommandsNV"))
	gpDrawCommandsStatesAddressNV = uintptr(getProcAddr("glDrawCommandsStatesAddressNV"))
	gpDrawCommandsStatesNV = uintptr(getProcAddr("glDrawCommandsStatesNV"))
	gpDrawElementArrayAPPLE = uintptr(getProcAddr("glDrawElementArrayAPPLE"))
	gpDrawElementArrayATI = uintptr(getProcAddr("glDrawElementArrayATI"))
	gpDrawElements = uintptr(getProcAddr("glDrawElements"))
	if gpDrawElements == 0 {
		return errors.New("glDrawElements")
	}
	gpDrawElementsBaseVertex = uintptr(getProcAddr("glDrawElementsBaseVertex"))
	gpDrawElementsIndirect = uintptr(getProcAddr("glDrawElementsIndirect"))
	gpDrawElementsInstancedARB = uintptr(getProcAddr("glDrawElementsInstancedARB"))
	gpDrawElementsInstancedBaseInstance = uintptr(getProcAddr("glDrawElementsInstancedBaseInstance"))
	gpDrawElementsInstancedBaseVertex = uintptr(getProcAddr("glDrawElementsInstancedBaseVertex"))
	gpDrawElementsInstancedBaseVertexBaseInstance = uintptr(getProcAddr("glDrawElementsInstancedBaseVertexBaseInstance"))
	gpDrawElementsInstancedEXT = uintptr(getProcAddr("glDrawElementsInstancedEXT"))
	gpDrawMeshArraysSUN = uintptr(getProcAddr("glDrawMeshArraysSUN"))
	gpDrawPixels = uintptr(getProcAddr("glDrawPixels"))
	if gpDrawPixels == 0 {
		return errors.New("glDrawPixels")
	}
	gpDrawRangeElementArrayAPPLE = uintptr(getProcAddr("glDrawRangeElementArrayAPPLE"))
	gpDrawRangeElementArrayATI = uintptr(getProcAddr("glDrawRangeElementArrayATI"))
	gpDrawRangeElements = uintptr(getProcAddr("glDrawRangeElements"))
	if gpDrawRangeElements == 0 {
		return errors.New("glDrawRangeElements")
	}
	gpDrawRangeElementsBaseVertex = uintptr(getProcAddr("glDrawRangeElementsBaseVertex"))
	gpDrawRangeElementsEXT = uintptr(getProcAddr("glDrawRangeElementsEXT"))
	gpDrawTextureNV = uintptr(getProcAddr("glDrawTextureNV"))
	gpDrawTransformFeedback = uintptr(getProcAddr("glDrawTransformFeedback"))
	gpDrawTransformFeedbackInstanced = uintptr(getProcAddr("glDrawTransformFeedbackInstanced"))
	gpDrawTransformFeedbackNV = uintptr(getProcAddr("glDrawTransformFeedbackNV"))
	gpDrawTransformFeedbackStream = uintptr(getProcAddr("glDrawTransformFeedbackStream"))
	gpDrawTransformFeedbackStreamInstanced = uintptr(getProcAddr("glDrawTransformFeedbackStreamInstanced"))
	gpDrawVkImageNV = uintptr(getProcAddr("glDrawVkImageNV"))
	gpEGLImageTargetTexStorageEXT = uintptr(getProcAddr("glEGLImageTargetTexStorageEXT"))
	gpEGLImageTargetTextureStorageEXT = uintptr(getProcAddr("glEGLImageTargetTextureStorageEXT"))
	gpEdgeFlag = uintptr(getProcAddr("glEdgeFlag"))
	if gpEdgeFlag == 0 {
		return errors.New("glEdgeFlag")
	}
	gpEdgeFlagFormatNV = uintptr(getProcAddr("glEdgeFlagFormatNV"))
	gpEdgeFlagPointer = uintptr(getProcAddr("glEdgeFlagPointer"))
	if gpEdgeFlagPointer == 0 {
		return errors.New("glEdgeFlagPointer")
	}
	gpEdgeFlagPointerEXT = uintptr(getProcAddr("glEdgeFlagPointerEXT"))
	gpEdgeFlagPointerListIBM = uintptr(getProcAddr("glEdgeFlagPointerListIBM"))
	gpEdgeFlagv = uintptr(getProcAddr("glEdgeFlagv"))
	if gpEdgeFlagv == 0 {
		return errors.New("glEdgeFlagv")
	}
	gpElementPointerAPPLE = uintptr(getProcAddr("glElementPointerAPPLE"))
	gpElementPointerATI = uintptr(getProcAddr("glElementPointerATI"))
	gpEnable = uintptr(getProcAddr("glEnable"))
	if gpEnable == 0 {
		return errors.New("glEnable")
	}
	gpEnableClientState = uintptr(getProcAddr("glEnableClientState"))
	if gpEnableClientState == 0 {
		return errors.New("glEnableClientState")
	}
	gpEnableClientStateIndexedEXT = uintptr(getProcAddr("glEnableClientStateIndexedEXT"))
	gpEnableClientStateiEXT = uintptr(getProcAddr("glEnableClientStateiEXT"))
	gpEnableIndexedEXT = uintptr(getProcAddr("glEnableIndexedEXT"))
	gpEnableVariantClientStateEXT = uintptr(getProcAddr("glEnableVariantClientStateEXT"))
	gpEnableVertexArrayAttrib = uintptr(getProcAddr("glEnableVertexArrayAttrib"))
	gpEnableVertexArrayAttribEXT = uintptr(getProcAddr("glEnableVertexArrayAttribEXT"))
	gpEnableVertexArrayEXT = uintptr(getProcAddr("glEnableVertexArrayEXT"))
	gpEnableVertexAttribAPPLE = uintptr(getProcAddr("glEnableVertexAttribAPPLE"))
	gpEnableVertexAttribArray = uintptr(getProcAddr("glEnableVertexAttribArray"))
	if gpEnableVertexAttribArray == 0 {
		return errors.New("glEnableVertexAttribArray")
	}
	gpEnableVertexAttribArrayARB = uintptr(getProcAddr("glEnableVertexAttribArrayARB"))
	gpEnd = uintptr(getProcAddr("glEnd"))
	if gpEnd == 0 {
		return errors.New("glEnd")
	}
	gpEndConditionalRenderNV = uintptr(getProcAddr("glEndConditionalRenderNV"))
	gpEndConditionalRenderNVX = uintptr(getProcAddr("glEndConditionalRenderNVX"))
	gpEndFragmentShaderATI = uintptr(getProcAddr("glEndFragmentShaderATI"))
	gpEndList = uintptr(getProcAddr("glEndList"))
	if gpEndList == 0 {
		return errors.New("glEndList")
	}
	gpEndOcclusionQueryNV = uintptr(getProcAddr("glEndOcclusionQueryNV"))
	gpEndPerfMonitorAMD = uintptr(getProcAddr("glEndPerfMonitorAMD"))
	gpEndPerfQueryINTEL = uintptr(getProcAddr("glEndPerfQueryINTEL"))
	gpEndQuery = uintptr(getProcAddr("glEndQuery"))
	if gpEndQuery == 0 {
		return errors.New("glEndQuery")
	}
	gpEndQueryARB = uintptr(getProcAddr("glEndQueryARB"))
	gpEndQueryIndexed = uintptr(getProcAddr("glEndQueryIndexed"))
	gpEndTransformFeedbackEXT = uintptr(getProcAddr("glEndTransformFeedbackEXT"))
	gpEndTransformFeedbackNV = uintptr(getProcAddr("glEndTransformFeedbackNV"))
	gpEndVertexShaderEXT = uintptr(getProcAddr("glEndVertexShaderEXT"))
	gpEndVideoCaptureNV = uintptr(getProcAddr("glEndVideoCaptureNV"))
	gpEvalCoord1d = uintptr(getProcAddr("glEvalCoord1d"))
	if gpEvalCoord1d == 0 {
		return errors.New("glEvalCoord1d")
	}
	gpEvalCoord1dv = uintptr(getProcAddr("glEvalCoord1dv"))
	if gpEvalCoord1dv == 0 {
		return errors.New("glEvalCoord1dv")
	}
	gpEvalCoord1f = uintptr(getProcAddr("glEvalCoord1f"))
	if gpEvalCoord1f == 0 {
		return errors.New("glEvalCoord1f")
	}
	gpEvalCoord1fv = uintptr(getProcAddr("glEvalCoord1fv"))
	if gpEvalCoord1fv == 0 {
		return errors.New("glEvalCoord1fv")
	}
	gpEvalCoord1xOES = uintptr(getProcAddr("glEvalCoord1xOES"))
	gpEvalCoord1xvOES = uintptr(getProcAddr("glEvalCoord1xvOES"))
	gpEvalCoord2d = uintptr(getProcAddr("glEvalCoord2d"))
	if gpEvalCoord2d == 0 {
		return errors.New("glEvalCoord2d")
	}
	gpEvalCoord2dv = uintptr(getProcAddr("glEvalCoord2dv"))
	if gpEvalCoord2dv == 0 {
		return errors.New("glEvalCoord2dv")
	}
	gpEvalCoord2f = uintptr(getProcAddr("glEvalCoord2f"))
	if gpEvalCoord2f == 0 {
		return errors.New("glEvalCoord2f")
	}
	gpEvalCoord2fv = uintptr(getProcAddr("glEvalCoord2fv"))
	if gpEvalCoord2fv == 0 {
		return errors.New("glEvalCoord2fv")
	}
	gpEvalCoord2xOES = uintptr(getProcAddr("glEvalCoord2xOES"))
	gpEvalCoord2xvOES = uintptr(getProcAddr("glEvalCoord2xvOES"))
	gpEvalMapsNV = uintptr(getProcAddr("glEvalMapsNV"))
	gpEvalMesh1 = uintptr(getProcAddr("glEvalMesh1"))
	if gpEvalMesh1 == 0 {
		return errors.New("glEvalMesh1")
	}
	gpEvalMesh2 = uintptr(getProcAddr("glEvalMesh2"))
	if gpEvalMesh2 == 0 {
		return errors.New("glEvalMesh2")
	}
	gpEvalPoint1 = uintptr(getProcAddr("glEvalPoint1"))
	if gpEvalPoint1 == 0 {
		return errors.New("glEvalPoint1")
	}
	gpEvalPoint2 = uintptr(getProcAddr("glEvalPoint2"))
	if gpEvalPoint2 == 0 {
		return errors.New("glEvalPoint2")
	}
	gpEvaluateDepthValuesARB = uintptr(getProcAddr("glEvaluateDepthValuesARB"))
	gpExecuteProgramNV = uintptr(getProcAddr("glExecuteProgramNV"))
	gpExtractComponentEXT = uintptr(getProcAddr("glExtractComponentEXT"))
	gpFeedbackBuffer = uintptr(getProcAddr("glFeedbackBuffer"))
	if gpFeedbackBuffer == 0 {
		return errors.New("glFeedbackBuffer")
	}
	gpFeedbackBufferxOES = uintptr(getProcAddr("glFeedbackBufferxOES"))
	gpFenceSync = uintptr(getProcAddr("glFenceSync"))
	gpFinalCombinerInputNV = uintptr(getProcAddr("glFinalCombinerInputNV"))
	gpFinish = uintptr(getProcAddr("glFinish"))
	if gpFinish == 0 {
		return errors.New("glFinish")
	}
	gpFinishAsyncSGIX = uintptr(getProcAddr("glFinishAsyncSGIX"))
	gpFinishFenceAPPLE = uintptr(getProcAddr("glFinishFenceAPPLE"))
	gpFinishFenceNV = uintptr(getProcAddr("glFinishFenceNV"))
	gpFinishObjectAPPLE = uintptr(getProcAddr("glFinishObjectAPPLE"))
	gpFinishTextureSUNX = uintptr(getProcAddr("glFinishTextureSUNX"))
	gpFlush = uintptr(getProcAddr("glFlush"))
	if gpFlush == 0 {
		return errors.New("glFlush")
	}
	gpFlushMappedBufferRange = uintptr(getProcAddr("glFlushMappedBufferRange"))
	gpFlushMappedBufferRangeAPPLE = uintptr(getProcAddr("glFlushMappedBufferRangeAPPLE"))
	gpFlushMappedNamedBufferRange = uintptr(getProcAddr("glFlushMappedNamedBufferRange"))
	gpFlushMappedNamedBufferRangeEXT = uintptr(getProcAddr("glFlushMappedNamedBufferRangeEXT"))
	gpFlushPixelDataRangeNV = uintptr(getProcAddr("glFlushPixelDataRangeNV"))
	gpFlushRasterSGIX = uintptr(getProcAddr("glFlushRasterSGIX"))
	gpFlushStaticDataIBM = uintptr(getProcAddr("glFlushStaticDataIBM"))
	gpFlushVertexArrayRangeAPPLE = uintptr(getProcAddr("glFlushVertexArrayRangeAPPLE"))
	gpFlushVertexArrayRangeNV = uintptr(getProcAddr("glFlushVertexArrayRangeNV"))
	gpFogCoordFormatNV = uintptr(getProcAddr("glFogCoordFormatNV"))
	gpFogCoordPointer = uintptr(getProcAddr("glFogCoordPointer"))
	if gpFogCoordPointer == 0 {
		return errors.New("glFogCoordPointer")
	}
	gpFogCoordPointerEXT = uintptr(getProcAddr("glFogCoordPointerEXT"))
	gpFogCoordPointerListIBM = uintptr(getProcAddr("glFogCoordPointerListIBM"))
	gpFogCoordd = uintptr(getProcAddr("glFogCoordd"))
	if gpFogCoordd == 0 {
		return errors.New("glFogCoordd")
	}
	gpFogCoorddEXT = uintptr(getProcAddr("glFogCoorddEXT"))
	gpFogCoorddv = uintptr(getProcAddr("glFogCoorddv"))
	if gpFogCoorddv == 0 {
		return errors.New("glFogCoorddv")
	}
	gpFogCoorddvEXT = uintptr(getProcAddr("glFogCoorddvEXT"))
	gpFogCoordf = uintptr(getProcAddr("glFogCoordf"))
	if gpFogCoordf == 0 {
		return errors.New("glFogCoordf")
	}
	gpFogCoordfEXT = uintptr(getProcAddr("glFogCoordfEXT"))
	gpFogCoordfv = uintptr(getProcAddr("glFogCoordfv"))
	if gpFogCoordfv == 0 {
		return errors.New("glFogCoordfv")
	}
	gpFogCoordfvEXT = uintptr(getProcAddr("glFogCoordfvEXT"))
	gpFogCoordhNV = uintptr(getProcAddr("glFogCoordhNV"))
	gpFogCoordhvNV = uintptr(getProcAddr("glFogCoordhvNV"))
	gpFogFuncSGIS = uintptr(getProcAddr("glFogFuncSGIS"))
	gpFogf = uintptr(getProcAddr("glFogf"))
	if gpFogf == 0 {
		return errors.New("glFogf")
	}
	gpFogfv = uintptr(getProcAddr("glFogfv"))
	if gpFogfv == 0 {
		return errors.New("glFogfv")
	}
	gpFogi = uintptr(getProcAddr("glFogi"))
	if gpFogi == 0 {
		return errors.New("glFogi")
	}
	gpFogiv = uintptr(getProcAddr("glFogiv"))
	if gpFogiv == 0 {
		return errors.New("glFogiv")
	}
	gpFogxOES = uintptr(getProcAddr("glFogxOES"))
	gpFogxvOES = uintptr(getProcAddr("glFogxvOES"))
	gpFragmentColorMaterialSGIX = uintptr(getProcAddr("glFragmentColorMaterialSGIX"))
	gpFragmentCoverageColorNV = uintptr(getProcAddr("glFragmentCoverageColorNV"))
	gpFragmentLightModelfSGIX = uintptr(getProcAddr("glFragmentLightModelfSGIX"))
	gpFragmentLightModelfvSGIX = uintptr(getProcAddr("glFragmentLightModelfvSGIX"))
	gpFragmentLightModeliSGIX = uintptr(getProcAddr("glFragmentLightModeliSGIX"))
	gpFragmentLightModelivSGIX = uintptr(getProcAddr("glFragmentLightModelivSGIX"))
	gpFragmentLightfSGIX = uintptr(getProcAddr("glFragmentLightfSGIX"))
	gpFragmentLightfvSGIX = uintptr(getProcAddr("glFragmentLightfvSGIX"))
	gpFragmentLightiSGIX = uintptr(getProcAddr("glFragmentLightiSGIX"))
	gpFragmentLightivSGIX = uintptr(getProcAddr("glFragmentLightivSGIX"))
	gpFragmentMaterialfSGIX = uintptr(getProcAddr("glFragmentMaterialfSGIX"))
	gpFragmentMaterialfvSGIX = uintptr(getProcAddr("glFragmentMaterialfvSGIX"))
	gpFragmentMaterialiSGIX = uintptr(getProcAddr("glFragmentMaterialiSGIX"))
	gpFragmentMaterialivSGIX = uintptr(getProcAddr("glFragmentMaterialivSGIX"))
	gpFrameTerminatorGREMEDY = uintptr(getProcAddr("glFrameTerminatorGREMEDY"))
	gpFrameZoomSGIX = uintptr(getProcAddr("glFrameZoomSGIX"))
	gpFramebufferDrawBufferEXT = uintptr(getProcAddr("glFramebufferDrawBufferEXT"))
	gpFramebufferDrawBuffersEXT = uintptr(getProcAddr("glFramebufferDrawBuffersEXT"))
	gpFramebufferFetchBarrierEXT = uintptr(getProcAddr("glFramebufferFetchBarrierEXT"))
	gpFramebufferParameteri = uintptr(getProcAddr("glFramebufferParameteri"))
	gpFramebufferReadBufferEXT = uintptr(getProcAddr("glFramebufferReadBufferEXT"))
	gpFramebufferRenderbuffer = uintptr(getProcAddr("glFramebufferRenderbuffer"))
	gpFramebufferRenderbufferEXT = uintptr(getProcAddr("glFramebufferRenderbufferEXT"))
	gpFramebufferSampleLocationsfvARB = uintptr(getProcAddr("glFramebufferSampleLocationsfvARB"))
	gpFramebufferSampleLocationsfvNV = uintptr(getProcAddr("glFramebufferSampleLocationsfvNV"))
	gpFramebufferSamplePositionsfvAMD = uintptr(getProcAddr("glFramebufferSamplePositionsfvAMD"))
	gpFramebufferTexture1D = uintptr(getProcAddr("glFramebufferTexture1D"))
	gpFramebufferTexture1DEXT = uintptr(getProcAddr("glFramebufferTexture1DEXT"))
	gpFramebufferTexture2D = uintptr(getProcAddr("glFramebufferTexture2D"))
	gpFramebufferTexture2DEXT = uintptr(getProcAddr("glFramebufferTexture2DEXT"))
	gpFramebufferTexture3D = uintptr(getProcAddr("glFramebufferTexture3D"))
	gpFramebufferTexture3DEXT = uintptr(getProcAddr("glFramebufferTexture3DEXT"))
	gpFramebufferTextureARB = uintptr(getProcAddr("glFramebufferTextureARB"))
	gpFramebufferTextureEXT = uintptr(getProcAddr("glFramebufferTextureEXT"))
	gpFramebufferTextureFaceARB = uintptr(getProcAddr("glFramebufferTextureFaceARB"))
	gpFramebufferTextureFaceEXT = uintptr(getProcAddr("glFramebufferTextureFaceEXT"))
	gpFramebufferTextureLayer = uintptr(getProcAddr("glFramebufferTextureLayer"))
	gpFramebufferTextureLayerARB = uintptr(getProcAddr("glFramebufferTextureLayerARB"))
	gpFramebufferTextureLayerEXT = uintptr(getProcAddr("glFramebufferTextureLayerEXT"))
	gpFramebufferTextureMultiviewOVR = uintptr(getProcAddr("glFramebufferTextureMultiviewOVR"))
	gpFreeObjectBufferATI = uintptr(getProcAddr("glFreeObjectBufferATI"))
	gpFrontFace = uintptr(getProcAddr("glFrontFace"))
	if gpFrontFace == 0 {
		return errors.New("glFrontFace")
	}
	gpFrustum = uintptr(getProcAddr("glFrustum"))
	if gpFrustum == 0 {
		return errors.New("glFrustum")
	}
	gpFrustumfOES = uintptr(getProcAddr("glFrustumfOES"))
	gpFrustumxOES = uintptr(getProcAddr("glFrustumxOES"))
	gpGenAsyncMarkersSGIX = uintptr(getProcAddr("glGenAsyncMarkersSGIX"))
	gpGenBuffers = uintptr(getProcAddr("glGenBuffers"))
	if gpGenBuffers == 0 {
		return errors.New("glGenBuffers")
	}
	gpGenBuffersARB = uintptr(getProcAddr("glGenBuffersARB"))
	gpGenFencesAPPLE = uintptr(getProcAddr("glGenFencesAPPLE"))
	gpGenFencesNV = uintptr(getProcAddr("glGenFencesNV"))
	gpGenFragmentShadersATI = uintptr(getProcAddr("glGenFragmentShadersATI"))
	gpGenFramebuffers = uintptr(getProcAddr("glGenFramebuffers"))
	gpGenFramebuffersEXT = uintptr(getProcAddr("glGenFramebuffersEXT"))
	gpGenLists = uintptr(getProcAddr("glGenLists"))
	if gpGenLists == 0 {
		return errors.New("glGenLists")
	}
	gpGenNamesAMD = uintptr(getProcAddr("glGenNamesAMD"))
	gpGenOcclusionQueriesNV = uintptr(getProcAddr("glGenOcclusionQueriesNV"))
	gpGenPathsNV = uintptr(getProcAddr("glGenPathsNV"))
	gpGenPerfMonitorsAMD = uintptr(getProcAddr("glGenPerfMonitorsAMD"))
	gpGenProgramPipelines = uintptr(getProcAddr("glGenProgramPipelines"))
	gpGenProgramPipelinesEXT = uintptr(getProcAddr("glGenProgramPipelinesEXT"))
	gpGenProgramsARB = uintptr(getProcAddr("glGenProgramsARB"))
	gpGenProgramsNV = uintptr(getProcAddr("glGenProgramsNV"))
	gpGenQueries = uintptr(getProcAddr("glGenQueries"))
	if gpGenQueries == 0 {
		return errors.New("glGenQueries")
	}
	gpGenQueriesARB = uintptr(getProcAddr("glGenQueriesARB"))
	gpGenQueryResourceTagNV = uintptr(getProcAddr("glGenQueryResourceTagNV"))
	gpGenRenderbuffers = uintptr(getProcAddr("glGenRenderbuffers"))
	gpGenRenderbuffersEXT = uintptr(getProcAddr("glGenRenderbuffersEXT"))
	gpGenSamplers = uintptr(getProcAddr("glGenSamplers"))
	gpGenSemaphoresEXT = uintptr(getProcAddr("glGenSemaphoresEXT"))
	gpGenSymbolsEXT = uintptr(getProcAddr("glGenSymbolsEXT"))
	gpGenTextures = uintptr(getProcAddr("glGenTextures"))
	if gpGenTextures == 0 {
		return errors.New("glGenTextures")
	}
	gpGenTexturesEXT = uintptr(getProcAddr("glGenTexturesEXT"))
	gpGenTransformFeedbacks = uintptr(getProcAddr("glGenTransformFeedbacks"))
	gpGenTransformFeedbacksNV = uintptr(getProcAddr("glGenTransformFeedbacksNV"))
	gpGenVertexArrays = uintptr(getProcAddr("glGenVertexArrays"))
	gpGenVertexArraysAPPLE = uintptr(getProcAddr("glGenVertexArraysAPPLE"))
	gpGenVertexShadersEXT = uintptr(getProcAddr("glGenVertexShadersEXT"))
	gpGenerateMipmap = uintptr(getProcAddr("glGenerateMipmap"))
	gpGenerateMipmapEXT = uintptr(getProcAddr("glGenerateMipmapEXT"))
	gpGenerateMultiTexMipmapEXT = uintptr(getProcAddr("glGenerateMultiTexMipmapEXT"))
	gpGenerateTextureMipmap = uintptr(getProcAddr("glGenerateTextureMipmap"))
	gpGenerateTextureMipmapEXT = uintptr(getProcAddr("glGenerateTextureMipmapEXT"))
	gpGetActiveAtomicCounterBufferiv = uintptr(getProcAddr("glGetActiveAtomicCounterBufferiv"))
	gpGetActiveAttrib = uintptr(getProcAddr("glGetActiveAttrib"))
	if gpGetActiveAttrib == 0 {
		return errors.New("glGetActiveAttrib")
	}
	gpGetActiveAttribARB = uintptr(getProcAddr("glGetActiveAttribARB"))
	gpGetActiveSubroutineName = uintptr(getProcAddr("glGetActiveSubroutineName"))
	gpGetActiveSubroutineUniformName = uintptr(getProcAddr("glGetActiveSubroutineUniformName"))
	gpGetActiveSubroutineUniformiv = uintptr(getProcAddr("glGetActiveSubroutineUniformiv"))
	gpGetActiveUniform = uintptr(getProcAddr("glGetActiveUniform"))
	if gpGetActiveUniform == 0 {
		return errors.New("glGetActiveUniform")
	}
	gpGetActiveUniformARB = uintptr(getProcAddr("glGetActiveUniformARB"))
	gpGetActiveUniformBlockName = uintptr(getProcAddr("glGetActiveUniformBlockName"))
	gpGetActiveUniformBlockiv = uintptr(getProcAddr("glGetActiveUniformBlockiv"))
	gpGetActiveUniformName = uintptr(getProcAddr("glGetActiveUniformName"))
	gpGetActiveUniformsiv = uintptr(getProcAddr("glGetActiveUniformsiv"))
	gpGetActiveVaryingNV = uintptr(getProcAddr("glGetActiveVaryingNV"))
	gpGetArrayObjectfvATI = uintptr(getProcAddr("glGetArrayObjectfvATI"))
	gpGetArrayObjectivATI = uintptr(getProcAddr("glGetArrayObjectivATI"))
	gpGetAttachedObjectsARB = uintptr(getProcAddr("glGetAttachedObjectsARB"))
	gpGetAttachedShaders = uintptr(getProcAddr("glGetAttachedShaders"))
	if gpGetAttachedShaders == 0 {
		return errors.New("glGetAttachedShaders")
	}
	gpGetAttribLocation = uintptr(getProcAddr("glGetAttribLocation"))
	if gpGetAttribLocation == 0 {
		return errors.New("glGetAttribLocation")
	}
	gpGetAttribLocationARB = uintptr(getProcAddr("glGetAttribLocationARB"))
	gpGetBooleanIndexedvEXT = uintptr(getProcAddr("glGetBooleanIndexedvEXT"))
	gpGetBooleanv = uintptr(getProcAddr("glGetBooleanv"))
	if gpGetBooleanv == 0 {
		return errors.New("glGetBooleanv")
	}
	gpGetBufferParameteriv = uintptr(getProcAddr("glGetBufferParameteriv"))
	if gpGetBufferParameteriv == 0 {
		return errors.New("glGetBufferParameteriv")
	}
	gpGetBufferParameterivARB = uintptr(getProcAddr("glGetBufferParameterivARB"))
	gpGetBufferParameterui64vNV = uintptr(getProcAddr("glGetBufferParameterui64vNV"))
	gpGetBufferPointerv = uintptr(getProcAddr("glGetBufferPointerv"))
	if gpGetBufferPointerv == 0 {
		return errors.New("glGetBufferPointerv")
	}
	gpGetBufferPointervARB = uintptr(getProcAddr("glGetBufferPointervARB"))
	gpGetBufferSubData = uintptr(getProcAddr("glGetBufferSubData"))
	if gpGetBufferSubData == 0 {
		return errors.New("glGetBufferSubData")
	}
	gpGetBufferSubDataARB = uintptr(getProcAddr("glGetBufferSubDataARB"))
	gpGetClipPlane = uintptr(getProcAddr("glGetClipPlane"))
	if gpGetClipPlane == 0 {
		return errors.New("glGetClipPlane")
	}
	gpGetClipPlanefOES = uintptr(getProcAddr("glGetClipPlanefOES"))
	gpGetClipPlanexOES = uintptr(getProcAddr("glGetClipPlanexOES"))
	gpGetColorTableEXT = uintptr(getProcAddr("glGetColorTableEXT"))
	gpGetColorTableParameterfvEXT = uintptr(getProcAddr("glGetColorTableParameterfvEXT"))
	gpGetColorTableParameterfvSGI = uintptr(getProcAddr("glGetColorTableParameterfvSGI"))
	gpGetColorTableParameterivEXT = uintptr(getProcAddr("glGetColorTableParameterivEXT"))
	gpGetColorTableParameterivSGI = uintptr(getProcAddr("glGetColorTableParameterivSGI"))
	gpGetColorTableSGI = uintptr(getProcAddr("glGetColorTableSGI"))
	gpGetCombinerInputParameterfvNV = uintptr(getProcAddr("glGetCombinerInputParameterfvNV"))
	gpGetCombinerInputParameterivNV = uintptr(getProcAddr("glGetCombinerInputParameterivNV"))
	gpGetCombinerOutputParameterfvNV = uintptr(getProcAddr("glGetCombinerOutputParameterfvNV"))
	gpGetCombinerOutputParameterivNV = uintptr(getProcAddr("glGetCombinerOutputParameterivNV"))
	gpGetCombinerStageParameterfvNV = uintptr(getProcAddr("glGetCombinerStageParameterfvNV"))
	gpGetCommandHeaderNV = uintptr(getProcAddr("glGetCommandHeaderNV"))
	gpGetCompressedMultiTexImageEXT = uintptr(getProcAddr("glGetCompressedMultiTexImageEXT"))
	gpGetCompressedTexImage = uintptr(getProcAddr("glGetCompressedTexImage"))
	if gpGetCompressedTexImage == 0 {
		return errors.New("glGetCompressedTexImage")
	}
	gpGetCompressedTexImageARB = uintptr(getProcAddr("glGetCompressedTexImageARB"))
	gpGetCompressedTextureImage = uintptr(getProcAddr("glGetCompressedTextureImage"))
	gpGetCompressedTextureImageEXT = uintptr(getProcAddr("glGetCompressedTextureImageEXT"))
	gpGetCompressedTextureSubImage = uintptr(getProcAddr("glGetCompressedTextureSubImage"))
	gpGetConvolutionFilterEXT = uintptr(getProcAddr("glGetConvolutionFilterEXT"))
	gpGetConvolutionParameterfvEXT = uintptr(getProcAddr("glGetConvolutionParameterfvEXT"))
	gpGetConvolutionParameterivEXT = uintptr(getProcAddr("glGetConvolutionParameterivEXT"))
	gpGetConvolutionParameterxvOES = uintptr(getProcAddr("glGetConvolutionParameterxvOES"))
	gpGetCoverageModulationTableNV = uintptr(getProcAddr("glGetCoverageModulationTableNV"))
	gpGetDebugMessageLog = uintptr(getProcAddr("glGetDebugMessageLog"))
	gpGetDebugMessageLogAMD = uintptr(getProcAddr("glGetDebugMessageLogAMD"))
	gpGetDebugMessageLogARB = uintptr(getProcAddr("glGetDebugMessageLogARB"))
	gpGetDebugMessageLogKHR = uintptr(getProcAddr("glGetDebugMessageLogKHR"))
	gpGetDetailTexFuncSGIS = uintptr(getProcAddr("glGetDetailTexFuncSGIS"))
	gpGetDoubleIndexedvEXT = uintptr(getProcAddr("glGetDoubleIndexedvEXT"))
	gpGetDoublei_v = uintptr(getProcAddr("glGetDoublei_v"))
	gpGetDoublei_vEXT = uintptr(getProcAddr("glGetDoublei_vEXT"))
	gpGetDoublev = uintptr(getProcAddr("glGetDoublev"))
	if gpGetDoublev == 0 {
		return errors.New("glGetDoublev")
	}
	gpGetError = uintptr(getProcAddr("glGetError"))
	if gpGetError == 0 {
		return errors.New("glGetError")
	}
	gpGetFenceivNV = uintptr(getProcAddr("glGetFenceivNV"))
	gpGetFinalCombinerInputParameterfvNV = uintptr(getProcAddr("glGetFinalCombinerInputParameterfvNV"))
	gpGetFinalCombinerInputParameterivNV = uintptr(getProcAddr("glGetFinalCombinerInputParameterivNV"))
	gpGetFirstPerfQueryIdINTEL = uintptr(getProcAddr("glGetFirstPerfQueryIdINTEL"))
	gpGetFixedvOES = uintptr(getProcAddr("glGetFixedvOES"))
	gpGetFloatIndexedvEXT = uintptr(getProcAddr("glGetFloatIndexedvEXT"))
	gpGetFloati_v = uintptr(getProcAddr("glGetFloati_v"))
	gpGetFloati_vEXT = uintptr(getProcAddr("glGetFloati_vEXT"))
	gpGetFloatv = uintptr(getProcAddr("glGetFloatv"))
	if gpGetFloatv == 0 {
		return errors.New("glGetFloatv")
	}
	gpGetFogFuncSGIS = uintptr(getProcAddr("glGetFogFuncSGIS"))
	gpGetFragDataIndex = uintptr(getProcAddr("glGetFragDataIndex"))
	gpGetFragDataLocationEXT = uintptr(getProcAddr("glGetFragDataLocationEXT"))
	gpGetFragmentLightfvSGIX = uintptr(getProcAddr("glGetFragmentLightfvSGIX"))
	gpGetFragmentLightivSGIX = uintptr(getProcAddr("glGetFragmentLightivSGIX"))
	gpGetFragmentMaterialfvSGIX = uintptr(getProcAddr("glGetFragmentMaterialfvSGIX"))
	gpGetFragmentMaterialivSGIX = uintptr(getProcAddr("glGetFragmentMaterialivSGIX"))
	gpGetFramebufferAttachmentParameteriv = uintptr(getProcAddr("glGetFramebufferAttachmentParameteriv"))
	gpGetFramebufferAttachmentParameterivEXT = uintptr(getProcAddr("glGetFramebufferAttachmentParameterivEXT"))
	gpGetFramebufferParameterfvAMD = uintptr(getProcAddr("glGetFramebufferParameterfvAMD"))
	gpGetFramebufferParameteriv = uintptr(getProcAddr("glGetFramebufferParameteriv"))
	gpGetFramebufferParameterivEXT = uintptr(getProcAddr("glGetFramebufferParameterivEXT"))
	gpGetGraphicsResetStatus = uintptr(getProcAddr("glGetGraphicsResetStatus"))
	gpGetGraphicsResetStatusARB = uintptr(getProcAddr("glGetGraphicsResetStatusARB"))
	gpGetGraphicsResetStatusKHR = uintptr(getProcAddr("glGetGraphicsResetStatusKHR"))
	gpGetHandleARB = uintptr(getProcAddr("glGetHandleARB"))
	gpGetHistogramEXT = uintptr(getProcAddr("glGetHistogramEXT"))
	gpGetHistogramParameterfvEXT = uintptr(getProcAddr("glGetHistogramParameterfvEXT"))
	gpGetHistogramParameterivEXT = uintptr(getProcAddr("glGetHistogramParameterivEXT"))
	gpGetHistogramParameterxvOES = uintptr(getProcAddr("glGetHistogramParameterxvOES"))
	gpGetImageHandleARB = uintptr(getProcAddr("glGetImageHandleARB"))
	gpGetImageHandleNV = uintptr(getProcAddr("glGetImageHandleNV"))
	gpGetImageTransformParameterfvHP = uintptr(getProcAddr("glGetImageTransformParameterfvHP"))
	gpGetImageTransformParameterivHP = uintptr(getProcAddr("glGetImageTransformParameterivHP"))
	gpGetInfoLogARB = uintptr(getProcAddr("glGetInfoLogARB"))
	gpGetInstrumentsSGIX = uintptr(getProcAddr("glGetInstrumentsSGIX"))
	gpGetInteger64v = uintptr(getProcAddr("glGetInteger64v"))
	gpGetIntegerIndexedvEXT = uintptr(getProcAddr("glGetIntegerIndexedvEXT"))
	gpGetIntegeri_v = uintptr(getProcAddr("glGetIntegeri_v"))
	gpGetIntegerui64i_vNV = uintptr(getProcAddr("glGetIntegerui64i_vNV"))
	gpGetIntegerui64vNV = uintptr(getProcAddr("glGetIntegerui64vNV"))
	gpGetIntegerv = uintptr(getProcAddr("glGetIntegerv"))
	if gpGetIntegerv == 0 {
		return errors.New("glGetIntegerv")
	}
	gpGetInternalformatSampleivNV = uintptr(getProcAddr("glGetInternalformatSampleivNV"))
	gpGetInternalformati64v = uintptr(getProcAddr("glGetInternalformati64v"))
	gpGetInternalformativ = uintptr(getProcAddr("glGetInternalformativ"))
	gpGetInvariantBooleanvEXT = uintptr(getProcAddr("glGetInvariantBooleanvEXT"))
	gpGetInvariantFloatvEXT = uintptr(getProcAddr("glGetInvariantFloatvEXT"))
	gpGetInvariantIntegervEXT = uintptr(getProcAddr("glGetInvariantIntegervEXT"))
	gpGetLightfv = uintptr(getProcAddr("glGetLightfv"))
	if gpGetLightfv == 0 {
		return errors.New("glGetLightfv")
	}
	gpGetLightiv = uintptr(getProcAddr("glGetLightiv"))
	if gpGetLightiv == 0 {
		return errors.New("glGetLightiv")
	}
	gpGetLightxOES = uintptr(getProcAddr("glGetLightxOES"))
	gpGetLightxvOES = uintptr(getProcAddr("glGetLightxvOES"))
	gpGetListParameterfvSGIX = uintptr(getProcAddr("glGetListParameterfvSGIX"))
	gpGetListParameterivSGIX = uintptr(getProcAddr("glGetListParameterivSGIX"))
	gpGetLocalConstantBooleanvEXT = uintptr(getProcAddr("glGetLocalConstantBooleanvEXT"))
	gpGetLocalConstantFloatvEXT = uintptr(getProcAddr("glGetLocalConstantFloatvEXT"))
	gpGetLocalConstantIntegervEXT = uintptr(getProcAddr("glGetLocalConstantIntegervEXT"))
	gpGetMapAttribParameterfvNV = uintptr(getProcAddr("glGetMapAttribParameterfvNV"))
	gpGetMapAttribParameterivNV = uintptr(getProcAddr("glGetMapAttribParameterivNV"))
	gpGetMapControlPointsNV = uintptr(getProcAddr("glGetMapControlPointsNV"))
	gpGetMapParameterfvNV = uintptr(getProcAddr("glGetMapParameterfvNV"))
	gpGetMapParameterivNV = uintptr(getProcAddr("glGetMapParameterivNV"))
	gpGetMapdv = uintptr(getProcAddr("glGetMapdv"))
	if gpGetMapdv == 0 {
		return errors.New("glGetMapdv")
	}
	gpGetMapfv = uintptr(getProcAddr("glGetMapfv"))
	if gpGetMapfv == 0 {
		return errors.New("glGetMapfv")
	}
	gpGetMapiv = uintptr(getProcAddr("glGetMapiv"))
	if gpGetMapiv == 0 {
		return errors.New("glGetMapiv")
	}
	gpGetMapxvOES = uintptr(getProcAddr("glGetMapxvOES"))
	gpGetMaterialfv = uintptr(getProcAddr("glGetMaterialfv"))
	if gpGetMaterialfv == 0 {
		return errors.New("glGetMaterialfv")
	}
	gpGetMaterialiv = uintptr(getProcAddr("glGetMaterialiv"))
	if gpGetMaterialiv == 0 {
		return errors.New("glGetMaterialiv")
	}
	gpGetMaterialxOES = uintptr(getProcAddr("glGetMaterialxOES"))
	gpGetMaterialxvOES = uintptr(getProcAddr("glGetMaterialxvOES"))
	gpGetMemoryObjectParameterivEXT = uintptr(getProcAddr("glGetMemoryObjectParameterivEXT"))
	gpGetMinmaxEXT = uintptr(getProcAddr("glGetMinmaxEXT"))
	gpGetMinmaxParameterfvEXT = uintptr(getProcAddr("glGetMinmaxParameterfvEXT"))
	gpGetMinmaxParameterivEXT = uintptr(getProcAddr("glGetMinmaxParameterivEXT"))
	gpGetMultiTexEnvfvEXT = uintptr(getProcAddr("glGetMultiTexEnvfvEXT"))
	gpGetMultiTexEnvivEXT = uintptr(getProcAddr("glGetMultiTexEnvivEXT"))
	gpGetMultiTexGendvEXT = uintptr(getProcAddr("glGetMultiTexGendvEXT"))
	gpGetMultiTexGenfvEXT = uintptr(getProcAddr("glGetMultiTexGenfvEXT"))
	gpGetMultiTexGenivEXT = uintptr(getProcAddr("glGetMultiTexGenivEXT"))
	gpGetMultiTexImageEXT = uintptr(getProcAddr("glGetMultiTexImageEXT"))
	gpGetMultiTexLevelParameterfvEXT = uintptr(getProcAddr("glGetMultiTexLevelParameterfvEXT"))
	gpGetMultiTexLevelParameterivEXT = uintptr(getProcAddr("glGetMultiTexLevelParameterivEXT"))
	gpGetMultiTexParameterIivEXT = uintptr(getProcAddr("glGetMultiTexParameterIivEXT"))
	gpGetMultiTexParameterIuivEXT = uintptr(getProcAddr("glGetMultiTexParameterIuivEXT"))
	gpGetMultiTexParameterfvEXT = uintptr(getProcAddr("glGetMultiTexParameterfvEXT"))
	gpGetMultiTexParameterivEXT = uintptr(getProcAddr("glGetMultiTexParameterivEXT"))
	gpGetMultisamplefv = uintptr(getProcAddr("glGetMultisamplefv"))
	gpGetMultisamplefvNV = uintptr(getProcAddr("glGetMultisamplefvNV"))
	gpGetNamedBufferParameteri64v = uintptr(getProcAddr("glGetNamedBufferParameteri64v"))
	gpGetNamedBufferParameteriv = uintptr(getProcAddr("glGetNamedBufferParameteriv"))
	gpGetNamedBufferParameterivEXT = uintptr(getProcAddr("glGetNamedBufferParameterivEXT"))
	gpGetNamedBufferParameterui64vNV = uintptr(getProcAddr("glGetNamedBufferParameterui64vNV"))
	gpGetNamedBufferPointerv = uintptr(getProcAddr("glGetNamedBufferPointerv"))
	gpGetNamedBufferPointervEXT = uintptr(getProcAddr("glGetNamedBufferPointervEXT"))
	gpGetNamedBufferSubData = uintptr(getProcAddr("glGetNamedBufferSubData"))
	gpGetNamedBufferSubDataEXT = uintptr(getProcAddr("glGetNamedBufferSubDataEXT"))
	gpGetNamedFramebufferAttachmentParameteriv = uintptr(getProcAddr("glGetNamedFramebufferAttachmentParameteriv"))
	gpGetNamedFramebufferAttachmentParameterivEXT = uintptr(getProcAddr("glGetNamedFramebufferAttachmentParameterivEXT"))
	gpGetNamedFramebufferParameterfvAMD = uintptr(getProcAddr("glGetNamedFramebufferParameterfvAMD"))
	gpGetNamedFramebufferParameteriv = uintptr(getProcAddr("glGetNamedFramebufferParameteriv"))
	gpGetNamedFramebufferParameterivEXT = uintptr(getProcAddr("glGetNamedFramebufferParameterivEXT"))
	gpGetNamedProgramLocalParameterIivEXT = uintptr(getProcAddr("glGetNamedProgramLocalParameterIivEXT"))
	gpGetNamedProgramLocalParameterIuivEXT = uintptr(getProcAddr("glGetNamedProgramLocalParameterIuivEXT"))
	gpGetNamedProgramLocalParameterdvEXT = uintptr(getProcAddr("glGetNamedProgramLocalParameterdvEXT"))
	gpGetNamedProgramLocalParameterfvEXT = uintptr(getProcAddr("glGetNamedProgramLocalParameterfvEXT"))
	gpGetNamedProgramStringEXT = uintptr(getProcAddr("glGetNamedProgramStringEXT"))
	gpGetNamedProgramivEXT = uintptr(getProcAddr("glGetNamedProgramivEXT"))
	gpGetNamedRenderbufferParameteriv = uintptr(getProcAddr("glGetNamedRenderbufferParameteriv"))
	gpGetNamedRenderbufferParameterivEXT = uintptr(getProcAddr("glGetNamedRenderbufferParameterivEXT"))
	gpGetNamedStringARB = uintptr(getProcAddr("glGetNamedStringARB"))
	gpGetNamedStringivARB = uintptr(getProcAddr("glGetNamedStringivARB"))
	gpGetNextPerfQueryIdINTEL = uintptr(getProcAddr("glGetNextPerfQueryIdINTEL"))
	gpGetObjectBufferfvATI = uintptr(getProcAddr("glGetObjectBufferfvATI"))
	gpGetObjectBufferivATI = uintptr(getProcAddr("glGetObjectBufferivATI"))
	gpGetObjectLabel = uintptr(getProcAddr("glGetObjectLabel"))
	gpGetObjectLabelEXT = uintptr(getProcAddr("glGetObjectLabelEXT"))
	gpGetObjectLabelKHR = uintptr(getProcAddr("glGetObjectLabelKHR"))
	gpGetObjectParameterfvARB = uintptr(getProcAddr("glGetObjectParameterfvARB"))
	gpGetObjectParameterivAPPLE = uintptr(getProcAddr("glGetObjectParameterivAPPLE"))
	gpGetObjectParameterivARB = uintptr(getProcAddr("glGetObjectParameterivARB"))
	gpGetObjectPtrLabel = uintptr(getProcAddr("glGetObjectPtrLabel"))
	gpGetObjectPtrLabelKHR = uintptr(getProcAddr("glGetObjectPtrLabelKHR"))
	gpGetOcclusionQueryivNV = uintptr(getProcAddr("glGetOcclusionQueryivNV"))
	gpGetOcclusionQueryuivNV = uintptr(getProcAddr("glGetOcclusionQueryuivNV"))
	gpGetPathCommandsNV = uintptr(getProcAddr("glGetPathCommandsNV"))
	gpGetPathCoordsNV = uintptr(getProcAddr("glGetPathCoordsNV"))
	gpGetPathDashArrayNV = uintptr(getProcAddr("glGetPathDashArrayNV"))
	gpGetPathLengthNV = uintptr(getProcAddr("glGetPathLengthNV"))
	gpGetPathMetricRangeNV = uintptr(getProcAddr("glGetPathMetricRangeNV"))
	gpGetPathMetricsNV = uintptr(getProcAddr("glGetPathMetricsNV"))
	gpGetPathParameterfvNV = uintptr(getProcAddr("glGetPathParameterfvNV"))
	gpGetPathParameterivNV = uintptr(getProcAddr("glGetPathParameterivNV"))
	gpGetPathSpacingNV = uintptr(getProcAddr("glGetPathSpacingNV"))
	gpGetPerfCounterInfoINTEL = uintptr(getProcAddr("glGetPerfCounterInfoINTEL"))
	gpGetPerfMonitorCounterDataAMD = uintptr(getProcAddr("glGetPerfMonitorCounterDataAMD"))
	gpGetPerfMonitorCounterInfoAMD = uintptr(getProcAddr("glGetPerfMonitorCounterInfoAMD"))
	gpGetPerfMonitorCounterStringAMD = uintptr(getProcAddr("glGetPerfMonitorCounterStringAMD"))
	gpGetPerfMonitorCountersAMD = uintptr(getProcAddr("glGetPerfMonitorCountersAMD"))
	gpGetPerfMonitorGroupStringAMD = uintptr(getProcAddr("glGetPerfMonitorGroupStringAMD"))
	gpGetPerfMonitorGroupsAMD = uintptr(getProcAddr("glGetPerfMonitorGroupsAMD"))
	gpGetPerfQueryDataINTEL = uintptr(getProcAddr("glGetPerfQueryDataINTEL"))
	gpGetPerfQueryIdByNameINTEL = uintptr(getProcAddr("glGetPerfQueryIdByNameINTEL"))
	gpGetPerfQueryInfoINTEL = uintptr(getProcAddr("glGetPerfQueryInfoINTEL"))
	gpGetPixelMapfv = uintptr(getProcAddr("glGetPixelMapfv"))
	if gpGetPixelMapfv == 0 {
		return errors.New("glGetPixelMapfv")
	}
	gpGetPixelMapuiv = uintptr(getProcAddr("glGetPixelMapuiv"))
	if gpGetPixelMapuiv == 0 {
		return errors.New("glGetPixelMapuiv")
	}
	gpGetPixelMapusv = uintptr(getProcAddr("glGetPixelMapusv"))
	if gpGetPixelMapusv == 0 {
		return errors.New("glGetPixelMapusv")
	}
	gpGetPixelMapxv = uintptr(getProcAddr("glGetPixelMapxv"))
	gpGetPixelTexGenParameterfvSGIS = uintptr(getProcAddr("glGetPixelTexGenParameterfvSGIS"))
	gpGetPixelTexGenParameterivSGIS = uintptr(getProcAddr("glGetPixelTexGenParameterivSGIS"))
	gpGetPixelTransformParameterfvEXT = uintptr(getProcAddr("glGetPixelTransformParameterfvEXT"))
	gpGetPixelTransformParameterivEXT = uintptr(getProcAddr("glGetPixelTransformParameterivEXT"))
	gpGetPointerIndexedvEXT = uintptr(getProcAddr("glGetPointerIndexedvEXT"))
	gpGetPointeri_vEXT = uintptr(getProcAddr("glGetPointeri_vEXT"))
	gpGetPointerv = uintptr(getProcAddr("glGetPointerv"))
	if gpGetPointerv == 0 {
		return errors.New("glGetPointerv")
	}
	gpGetPointervEXT = uintptr(getProcAddr("glGetPointervEXT"))
	gpGetPointervKHR = uintptr(getProcAddr("glGetPointervKHR"))
	gpGetPolygonStipple = uintptr(getProcAddr("glGetPolygonStipple"))
	if gpGetPolygonStipple == 0 {
		return errors.New("glGetPolygonStipple")
	}
	gpGetProgramBinary = uintptr(getProcAddr("glGetProgramBinary"))
	gpGetProgramEnvParameterIivNV = uintptr(getProcAddr("glGetProgramEnvParameterIivNV"))
	gpGetProgramEnvParameterIuivNV = uintptr(getProcAddr("glGetProgramEnvParameterIuivNV"))
	gpGetProgramEnvParameterdvARB = uintptr(getProcAddr("glGetProgramEnvParameterdvARB"))
	gpGetProgramEnvParameterfvARB = uintptr(getProcAddr("glGetProgramEnvParameterfvARB"))
	gpGetProgramInfoLog = uintptr(getProcAddr("glGetProgramInfoLog"))
	if gpGetProgramInfoLog == 0 {
		return errors.New("glGetProgramInfoLog")
	}
	gpGetProgramInterfaceiv = uintptr(getProcAddr("glGetProgramInterfaceiv"))
	gpGetProgramLocalParameterIivNV = uintptr(getProcAddr("glGetProgramLocalParameterIivNV"))
	gpGetProgramLocalParameterIuivNV = uintptr(getProcAddr("glGetProgramLocalParameterIuivNV"))
	gpGetProgramLocalParameterdvARB = uintptr(getProcAddr("glGetProgramLocalParameterdvARB"))
	gpGetProgramLocalParameterfvARB = uintptr(getProcAddr("glGetProgramLocalParameterfvARB"))
	gpGetProgramNamedParameterdvNV = uintptr(getProcAddr("glGetProgramNamedParameterdvNV"))
	gpGetProgramNamedParameterfvNV = uintptr(getProcAddr("glGetProgramNamedParameterfvNV"))
	gpGetProgramParameterdvNV = uintptr(getProcAddr("glGetProgramParameterdvNV"))
	gpGetProgramParameterfvNV = uintptr(getProcAddr("glGetProgramParameterfvNV"))
	gpGetProgramPipelineInfoLog = uintptr(getProcAddr("glGetProgramPipelineInfoLog"))
	gpGetProgramPipelineInfoLogEXT = uintptr(getProcAddr("glGetProgramPipelineInfoLogEXT"))
	gpGetProgramPipelineiv = uintptr(getProcAddr("glGetProgramPipelineiv"))
	gpGetProgramPipelineivEXT = uintptr(getProcAddr("glGetProgramPipelineivEXT"))
	gpGetProgramResourceIndex = uintptr(getProcAddr("glGetProgramResourceIndex"))
	gpGetProgramResourceLocation = uintptr(getProcAddr("glGetProgramResourceLocation"))
	gpGetProgramResourceLocationIndex = uintptr(getProcAddr("glGetProgramResourceLocationIndex"))
	gpGetProgramResourceName = uintptr(getProcAddr("glGetProgramResourceName"))
	gpGetProgramResourcefvNV = uintptr(getProcAddr("glGetProgramResourcefvNV"))
	gpGetProgramResourceiv = uintptr(getProcAddr("glGetProgramResourceiv"))
	gpGetProgramStageiv = uintptr(getProcAddr("glGetProgramStageiv"))
	gpGetProgramStringARB = uintptr(getProcAddr("glGetProgramStringARB"))
	gpGetProgramStringNV = uintptr(getProcAddr("glGetProgramStringNV"))
	gpGetProgramSubroutineParameteruivNV = uintptr(getProcAddr("glGetProgramSubroutineParameteruivNV"))
	gpGetProgramiv = uintptr(getProcAddr("glGetProgramiv"))
	if gpGetProgramiv == 0 {
		return errors.New("glGetProgramiv")
	}
	gpGetProgramivARB = uintptr(getProcAddr("glGetProgramivARB"))
	gpGetProgramivNV = uintptr(getProcAddr("glGetProgramivNV"))
	gpGetQueryBufferObjecti64v = uintptr(getProcAddr("glGetQueryBufferObjecti64v"))
	gpGetQueryBufferObjectiv = uintptr(getProcAddr("glGetQueryBufferObjectiv"))
	gpGetQueryBufferObjectui64v = uintptr(getProcAddr("glGetQueryBufferObjectui64v"))
	gpGetQueryBufferObjectuiv = uintptr(getProcAddr("glGetQueryBufferObjectuiv"))
	gpGetQueryIndexediv = uintptr(getProcAddr("glGetQueryIndexediv"))
	gpGetQueryObjecti64v = uintptr(getProcAddr("glGetQueryObjecti64v"))
	gpGetQueryObjecti64vEXT = uintptr(getProcAddr("glGetQueryObjecti64vEXT"))
	gpGetQueryObjectiv = uintptr(getProcAddr("glGetQueryObjectiv"))
	if gpGetQueryObjectiv == 0 {
		return errors.New("glGetQueryObjectiv")
	}
	gpGetQueryObjectivARB = uintptr(getProcAddr("glGetQueryObjectivARB"))
	gpGetQueryObjectui64v = uintptr(getProcAddr("glGetQueryObjectui64v"))
	gpGetQueryObjectui64vEXT = uintptr(getProcAddr("glGetQueryObjectui64vEXT"))
	gpGetQueryObjectuiv = uintptr(getProcAddr("glGetQueryObjectuiv"))
	if gpGetQueryObjectuiv == 0 {
		return errors.New("glGetQueryObjectuiv")
	}
	gpGetQueryObjectuivARB = uintptr(getProcAddr("glGetQueryObjectuivARB"))
	gpGetQueryiv = uintptr(getProcAddr("glGetQueryiv"))
	if gpGetQueryiv == 0 {
		return errors.New("glGetQueryiv")
	}
	gpGetQueryivARB = uintptr(getProcAddr("glGetQueryivARB"))
	gpGetRenderbufferParameteriv = uintptr(getProcAddr("glGetRenderbufferParameteriv"))
	gpGetRenderbufferParameterivEXT = uintptr(getProcAddr("glGetRenderbufferParameterivEXT"))
	gpGetSamplerParameterIiv = uintptr(getProcAddr("glGetSamplerParameterIiv"))
	gpGetSamplerParameterIuiv = uintptr(getProcAddr("glGetSamplerParameterIuiv"))
	gpGetSamplerParameterfv = uintptr(getProcAddr("glGetSamplerParameterfv"))
	gpGetSamplerParameteriv = uintptr(getProcAddr("glGetSamplerParameteriv"))
	gpGetSemaphoreParameterui64vEXT = uintptr(getProcAddr("glGetSemaphoreParameterui64vEXT"))
	gpGetSeparableFilterEXT = uintptr(getProcAddr("glGetSeparableFilterEXT"))
	gpGetShaderInfoLog = uintptr(getProcAddr("glGetShaderInfoLog"))
	if gpGetShaderInfoLog == 0 {
		return errors.New("glGetShaderInfoLog")
	}
	gpGetShaderPrecisionFormat = uintptr(getProcAddr("glGetShaderPrecisionFormat"))
	gpGetShaderSource = uintptr(getProcAddr("glGetShaderSource"))
	if gpGetShaderSource == 0 {
		return errors.New("glGetShaderSource")
	}
	gpGetShaderSourceARB = uintptr(getProcAddr("glGetShaderSourceARB"))
	gpGetShaderiv = uintptr(getProcAddr("glGetShaderiv"))
	if gpGetShaderiv == 0 {
		return errors.New("glGetShaderiv")
	}
	gpGetSharpenTexFuncSGIS = uintptr(getProcAddr("glGetSharpenTexFuncSGIS"))
	gpGetStageIndexNV = uintptr(getProcAddr("glGetStageIndexNV"))
	gpGetString = uintptr(getProcAddr("glGetString"))
	if gpGetString == 0 {
		return errors.New("glGetString")
	}
	gpGetSubroutineIndex = uintptr(getProcAddr("glGetSubroutineIndex"))
	gpGetSubroutineUniformLocation = uintptr(getProcAddr("glGetSubroutineUniformLocation"))
	gpGetSynciv = uintptr(getProcAddr("glGetSynciv"))
	gpGetTexBumpParameterfvATI = uintptr(getProcAddr("glGetTexBumpParameterfvATI"))
	gpGetTexBumpParameterivATI = uintptr(getProcAddr("glGetTexBumpParameterivATI"))
	gpGetTexEnvfv = uintptr(getProcAddr("glGetTexEnvfv"))
	if gpGetTexEnvfv == 0 {
		return errors.New("glGetTexEnvfv")
	}
	gpGetTexEnviv = uintptr(getProcAddr("glGetTexEnviv"))
	if gpGetTexEnviv == 0 {
		return errors.New("glGetTexEnviv")
	}
	gpGetTexEnvxvOES = uintptr(getProcAddr("glGetTexEnvxvOES"))
	gpGetTexFilterFuncSGIS = uintptr(getProcAddr("glGetTexFilterFuncSGIS"))
	gpGetTexGendv = uintptr(getProcAddr("glGetTexGendv"))
	if gpGetTexGendv == 0 {
		return errors.New("glGetTexGendv")
	}
	gpGetTexGenfv = uintptr(getProcAddr("glGetTexGenfv"))
	if gpGetTexGenfv == 0 {
		return errors.New("glGetTexGenfv")
	}
	gpGetTexGeniv = uintptr(getProcAddr("glGetTexGeniv"))
	if gpGetTexGeniv == 0 {
		return errors.New("glGetTexGeniv")
	}
	gpGetTexGenxvOES = uintptr(getProcAddr("glGetTexGenxvOES"))
	gpGetTexImage = uintptr(getProcAddr("glGetTexImage"))
	if gpGetTexImage == 0 {
		return errors.New("glGetTexImage")
	}
	gpGetTexLevelParameterfv = uintptr(getProcAddr("glGetTexLevelParameterfv"))
	if gpGetTexLevelParameterfv == 0 {
		return errors.New("glGetTexLevelParameterfv")
	}
	gpGetTexLevelParameteriv = uintptr(getProcAddr("glGetTexLevelParameteriv"))
	if gpGetTexLevelParameteriv == 0 {
		return errors.New("glGetTexLevelParameteriv")
	}
	gpGetTexLevelParameterxvOES = uintptr(getProcAddr("glGetTexLevelParameterxvOES"))
	gpGetTexParameterIivEXT = uintptr(getProcAddr("glGetTexParameterIivEXT"))
	gpGetTexParameterIuivEXT = uintptr(getProcAddr("glGetTexParameterIuivEXT"))
	gpGetTexParameterPointervAPPLE = uintptr(getProcAddr("glGetTexParameterPointervAPPLE"))
	gpGetTexParameterfv = uintptr(getProcAddr("glGetTexParameterfv"))
	if gpGetTexParameterfv == 0 {
		return errors.New("glGetTexParameterfv")
	}
	gpGetTexParameteriv = uintptr(getProcAddr("glGetTexParameteriv"))
	if gpGetTexParameteriv == 0 {
		return errors.New("glGetTexParameteriv")
	}
	gpGetTexParameterxvOES = uintptr(getProcAddr("glGetTexParameterxvOES"))
	gpGetTextureHandleARB = uintptr(getProcAddr("glGetTextureHandleARB"))
	gpGetTextureHandleNV = uintptr(getProcAddr("glGetTextureHandleNV"))
	gpGetTextureImage = uintptr(getProcAddr("glGetTextureImage"))
	gpGetTextureImageEXT = uintptr(getProcAddr("glGetTextureImageEXT"))
	gpGetTextureLevelParameterfv = uintptr(getProcAddr("glGetTextureLevelParameterfv"))
	gpGetTextureLevelParameterfvEXT = uintptr(getProcAddr("glGetTextureLevelParameterfvEXT"))
	gpGetTextureLevelParameteriv = uintptr(getProcAddr("glGetTextureLevelParameteriv"))
	gpGetTextureLevelParameterivEXT = uintptr(getProcAddr("glGetTextureLevelParameterivEXT"))
	gpGetTextureParameterIiv = uintptr(getProcAddr("glGetTextureParameterIiv"))
	gpGetTextureParameterIivEXT = uintptr(getProcAddr("glGetTextureParameterIivEXT"))
	gpGetTextureParameterIuiv = uintptr(getProcAddr("glGetTextureParameterIuiv"))
	gpGetTextureParameterIuivEXT = uintptr(getProcAddr("glGetTextureParameterIuivEXT"))
	gpGetTextureParameterfv = uintptr(getProcAddr("glGetTextureParameterfv"))
	gpGetTextureParameterfvEXT = uintptr(getProcAddr("glGetTextureParameterfvEXT"))
	gpGetTextureParameteriv = uintptr(getProcAddr("glGetTextureParameteriv"))
	gpGetTextureParameterivEXT = uintptr(getProcAddr("glGetTextureParameterivEXT"))
	gpGetTextureSamplerHandleARB = uintptr(getProcAddr("glGetTextureSamplerHandleARB"))
	gpGetTextureSamplerHandleNV = uintptr(getProcAddr("glGetTextureSamplerHandleNV"))
	gpGetTextureSubImage = uintptr(getProcAddr("glGetTextureSubImage"))
	gpGetTrackMatrixivNV = uintptr(getProcAddr("glGetTrackMatrixivNV"))
	gpGetTransformFeedbackVaryingEXT = uintptr(getProcAddr("glGetTransformFeedbackVaryingEXT"))
	gpGetTransformFeedbackVaryingNV = uintptr(getProcAddr("glGetTransformFeedbackVaryingNV"))
	gpGetTransformFeedbacki64_v = uintptr(getProcAddr("glGetTransformFeedbacki64_v"))
	gpGetTransformFeedbacki_v = uintptr(getProcAddr("glGetTransformFeedbacki_v"))
	gpGetTransformFeedbackiv = uintptr(getProcAddr("glGetTransformFeedbackiv"))
	gpGetUniformBlockIndex = uintptr(getProcAddr("glGetUniformBlockIndex"))
	gpGetUniformBufferSizeEXT = uintptr(getProcAddr("glGetUniformBufferSizeEXT"))
	gpGetUniformIndices = uintptr(getProcAddr("glGetUniformIndices"))
	gpGetUniformLocation = uintptr(getProcAddr("glGetUniformLocation"))
	if gpGetUniformLocation == 0 {
		return errors.New("glGetUniformLocation")
	}
	gpGetUniformLocationARB = uintptr(getProcAddr("glGetUniformLocationARB"))
	gpGetUniformOffsetEXT = uintptr(getProcAddr("glGetUniformOffsetEXT"))
	gpGetUniformSubroutineuiv = uintptr(getProcAddr("glGetUniformSubroutineuiv"))
	gpGetUniformdv = uintptr(getProcAddr("glGetUniformdv"))
	gpGetUniformfv = uintptr(getProcAddr("glGetUniformfv"))
	if gpGetUniformfv == 0 {
		return errors.New("glGetUniformfv")
	}
	gpGetUniformfvARB = uintptr(getProcAddr("glGetUniformfvARB"))
	gpGetUniformi64vARB = uintptr(getProcAddr("glGetUniformi64vARB"))
	gpGetUniformi64vNV = uintptr(getProcAddr("glGetUniformi64vNV"))
	gpGetUniformiv = uintptr(getProcAddr("glGetUniformiv"))
	if gpGetUniformiv == 0 {
		return errors.New("glGetUniformiv")
	}
	gpGetUniformivARB = uintptr(getProcAddr("glGetUniformivARB"))
	gpGetUniformui64vARB = uintptr(getProcAddr("glGetUniformui64vARB"))
	gpGetUniformui64vNV = uintptr(getProcAddr("glGetUniformui64vNV"))
	gpGetUniformuivEXT = uintptr(getProcAddr("glGetUniformuivEXT"))
	gpGetUnsignedBytei_vEXT = uintptr(getProcAddr("glGetUnsignedBytei_vEXT"))
	gpGetUnsignedBytevEXT = uintptr(getProcAddr("glGetUnsignedBytevEXT"))
	gpGetVariantArrayObjectfvATI = uintptr(getProcAddr("glGetVariantArrayObjectfvATI"))
	gpGetVariantArrayObjectivATI = uintptr(getProcAddr("glGetVariantArrayObjectivATI"))
	gpGetVariantBooleanvEXT = uintptr(getProcAddr("glGetVariantBooleanvEXT"))
	gpGetVariantFloatvEXT = uintptr(getProcAddr("glGetVariantFloatvEXT"))
	gpGetVariantIntegervEXT = uintptr(getProcAddr("glGetVariantIntegervEXT"))
	gpGetVariantPointervEXT = uintptr(getProcAddr("glGetVariantPointervEXT"))
	gpGetVaryingLocationNV = uintptr(getProcAddr("glGetVaryingLocationNV"))
	gpGetVertexArrayIndexed64iv = uintptr(getProcAddr("glGetVertexArrayIndexed64iv"))
	gpGetVertexArrayIndexediv = uintptr(getProcAddr("glGetVertexArrayIndexediv"))
	gpGetVertexArrayIntegeri_vEXT = uintptr(getProcAddr("glGetVertexArrayIntegeri_vEXT"))
	gpGetVertexArrayIntegervEXT = uintptr(getProcAddr("glGetVertexArrayIntegervEXT"))
	gpGetVertexArrayPointeri_vEXT = uintptr(getProcAddr("glGetVertexArrayPointeri_vEXT"))
	gpGetVertexArrayPointervEXT = uintptr(getProcAddr("glGetVertexArrayPointervEXT"))
	gpGetVertexArrayiv = uintptr(getProcAddr("glGetVertexArrayiv"))
	gpGetVertexAttribArrayObjectfvATI = uintptr(getProcAddr("glGetVertexAttribArrayObjectfvATI"))
	gpGetVertexAttribArrayObjectivATI = uintptr(getProcAddr("glGetVertexAttribArrayObjectivATI"))
	gpGetVertexAttribIivEXT = uintptr(getProcAddr("glGetVertexAttribIivEXT"))
	gpGetVertexAttribIuivEXT = uintptr(getProcAddr("glGetVertexAttribIuivEXT"))
	gpGetVertexAttribLdv = uintptr(getProcAddr("glGetVertexAttribLdv"))
	gpGetVertexAttribLdvEXT = uintptr(getProcAddr("glGetVertexAttribLdvEXT"))
	gpGetVertexAttribLi64vNV = uintptr(getProcAddr("glGetVertexAttribLi64vNV"))
	gpGetVertexAttribLui64vARB = uintptr(getProcAddr("glGetVertexAttribLui64vARB"))
	gpGetVertexAttribLui64vNV = uintptr(getProcAddr("glGetVertexAttribLui64vNV"))
	gpGetVertexAttribPointerv = uintptr(getProcAddr("glGetVertexAttribPointerv"))
	if gpGetVertexAttribPointerv == 0 {
		return errors.New("glGetVertexAttribPointerv")
	}
	gpGetVertexAttribPointervARB = uintptr(getProcAddr("glGetVertexAttribPointervARB"))
	gpGetVertexAttribPointervNV = uintptr(getProcAddr("glGetVertexAttribPointervNV"))
	gpGetVertexAttribdv = uintptr(getProcAddr("glGetVertexAttribdv"))
	if gpGetVertexAttribdv == 0 {
		return errors.New("glGetVertexAttribdv")
	}
	gpGetVertexAttribdvARB = uintptr(getProcAddr("glGetVertexAttribdvARB"))
	gpGetVertexAttribdvNV = uintptr(getProcAddr("glGetVertexAttribdvNV"))
	gpGetVertexAttribfv = uintptr(getProcAddr("glGetVertexAttribfv"))
	if gpGetVertexAttribfv == 0 {
		return errors.New("glGetVertexAttribfv")
	}
	gpGetVertexAttribfvARB = uintptr(getProcAddr("glGetVertexAttribfvARB"))
	gpGetVertexAttribfvNV = uintptr(getProcAddr("glGetVertexAttribfvNV"))
	gpGetVertexAttribiv = uintptr(getProcAddr("glGetVertexAttribiv"))
	if gpGetVertexAttribiv == 0 {
		return errors.New("glGetVertexAttribiv")
	}
	gpGetVertexAttribivARB = uintptr(getProcAddr("glGetVertexAttribivARB"))
	gpGetVertexAttribivNV = uintptr(getProcAddr("glGetVertexAttribivNV"))
	gpGetVideoCaptureStreamdvNV = uintptr(getProcAddr("glGetVideoCaptureStreamdvNV"))
	gpGetVideoCaptureStreamfvNV = uintptr(getProcAddr("glGetVideoCaptureStreamfvNV"))
	gpGetVideoCaptureStreamivNV = uintptr(getProcAddr("glGetVideoCaptureStreamivNV"))
	gpGetVideoCaptureivNV = uintptr(getProcAddr("glGetVideoCaptureivNV"))
	gpGetVideoi64vNV = uintptr(getProcAddr("glGetVideoi64vNV"))
	gpGetVideoivNV = uintptr(getProcAddr("glGetVideoivNV"))
	gpGetVideoui64vNV = uintptr(getProcAddr("glGetVideoui64vNV"))
	gpGetVideouivNV = uintptr(getProcAddr("glGetVideouivNV"))
	gpGetVkProcAddrNV = uintptr(getProcAddr("glGetVkProcAddrNV"))
	gpGetnCompressedTexImageARB = uintptr(getProcAddr("glGetnCompressedTexImageARB"))
	gpGetnTexImageARB = uintptr(getProcAddr("glGetnTexImageARB"))
	gpGetnUniformdvARB = uintptr(getProcAddr("glGetnUniformdvARB"))
	gpGetnUniformfv = uintptr(getProcAddr("glGetnUniformfv"))
	gpGetnUniformfvARB = uintptr(getProcAddr("glGetnUniformfvARB"))
	gpGetnUniformfvKHR = uintptr(getProcAddr("glGetnUniformfvKHR"))
	gpGetnUniformi64vARB = uintptr(getProcAddr("glGetnUniformi64vARB"))
	gpGetnUniformiv = uintptr(getProcAddr("glGetnUniformiv"))
	gpGetnUniformivARB = uintptr(getProcAddr("glGetnUniformivARB"))
	gpGetnUniformivKHR = uintptr(getProcAddr("glGetnUniformivKHR"))
	gpGetnUniformui64vARB = uintptr(getProcAddr("glGetnUniformui64vARB"))
	gpGetnUniformuiv = uintptr(getProcAddr("glGetnUniformuiv"))
	gpGetnUniformuivARB = uintptr(getProcAddr("glGetnUniformuivARB"))
	gpGetnUniformuivKHR = uintptr(getProcAddr("glGetnUniformuivKHR"))
	gpGlobalAlphaFactorbSUN = uintptr(getProcAddr("glGlobalAlphaFactorbSUN"))
	gpGlobalAlphaFactordSUN = uintptr(getProcAddr("glGlobalAlphaFactordSUN"))
	gpGlobalAlphaFactorfSUN = uintptr(getProcAddr("glGlobalAlphaFactorfSUN"))
	gpGlobalAlphaFactoriSUN = uintptr(getProcAddr("glGlobalAlphaFactoriSUN"))
	gpGlobalAlphaFactorsSUN = uintptr(getProcAddr("glGlobalAlphaFactorsSUN"))
	gpGlobalAlphaFactorubSUN = uintptr(getProcAddr("glGlobalAlphaFactorubSUN"))
	gpGlobalAlphaFactoruiSUN = uintptr(getProcAddr("glGlobalAlphaFactoruiSUN"))
	gpGlobalAlphaFactorusSUN = uintptr(getProcAddr("glGlobalAlphaFactorusSUN"))
	gpHint = uintptr(getProcAddr("glHint"))
	if gpHint == 0 {
		return errors.New("glHint")
	}
	gpHintPGI = uintptr(getProcAddr("glHintPGI"))
	gpHistogramEXT = uintptr(getProcAddr("glHistogramEXT"))
	gpIglooInterfaceSGIX = uintptr(getProcAddr("glIglooInterfaceSGIX"))
	gpImageTransformParameterfHP = uintptr(getProcAddr("glImageTransformParameterfHP"))
	gpImageTransformParameterfvHP = uintptr(getProcAddr("glImageTransformParameterfvHP"))
	gpImageTransformParameteriHP = uintptr(getProcAddr("glImageTransformParameteriHP"))
	gpImageTransformParameterivHP = uintptr(getProcAddr("glImageTransformParameterivHP"))
	gpImportMemoryFdEXT = uintptr(getProcAddr("glImportMemoryFdEXT"))
	gpImportMemoryWin32HandleEXT = uintptr(getProcAddr("glImportMemoryWin32HandleEXT"))
	gpImportMemoryWin32NameEXT = uintptr(getProcAddr("glImportMemoryWin32NameEXT"))
	gpImportSemaphoreFdEXT = uintptr(getProcAddr("glImportSemaphoreFdEXT"))
	gpImportSemaphoreWin32HandleEXT = uintptr(getProcAddr("glImportSemaphoreWin32HandleEXT"))
	gpImportSemaphoreWin32NameEXT = uintptr(getProcAddr("glImportSemaphoreWin32NameEXT"))
	gpImportSyncEXT = uintptr(getProcAddr("glImportSyncEXT"))
	gpIndexFormatNV = uintptr(getProcAddr("glIndexFormatNV"))
	gpIndexFuncEXT = uintptr(getProcAddr("glIndexFuncEXT"))
	gpIndexMask = uintptr(getProcAddr("glIndexMask"))
	if gpIndexMask == 0 {
		return errors.New("glIndexMask")
	}
	gpIndexMaterialEXT = uintptr(getProcAddr("glIndexMaterialEXT"))
	gpIndexPointer = uintptr(getProcAddr("glIndexPointer"))
	if gpIndexPointer == 0 {
		return errors.New("glIndexPointer")
	}
	gpIndexPointerEXT = uintptr(getProcAddr("glIndexPointerEXT"))
	gpIndexPointerListIBM = uintptr(getProcAddr("glIndexPointerListIBM"))
	gpIndexd = uintptr(getProcAddr("glIndexd"))
	if gpIndexd == 0 {
		return errors.New("glIndexd")
	}
	gpIndexdv = uintptr(getProcAddr("glIndexdv"))
	if gpIndexdv == 0 {
		return errors.New("glIndexdv")
	}
	gpIndexf = uintptr(getProcAddr("glIndexf"))
	if gpIndexf == 0 {
		return errors.New("glIndexf")
	}
	gpIndexfv = uintptr(getProcAddr("glIndexfv"))
	if gpIndexfv == 0 {
		return errors.New("glIndexfv")
	}
	gpIndexi = uintptr(getProcAddr("glIndexi"))
	if gpIndexi == 0 {
		return errors.New("glIndexi")
	}
	gpIndexiv = uintptr(getProcAddr("glIndexiv"))
	if gpIndexiv == 0 {
		return errors.New("glIndexiv")
	}
	gpIndexs = uintptr(getProcAddr("glIndexs"))
	if gpIndexs == 0 {
		return errors.New("glIndexs")
	}
	gpIndexsv = uintptr(getProcAddr("glIndexsv"))
	if gpIndexsv == 0 {
		return errors.New("glIndexsv")
	}
	gpIndexub = uintptr(getProcAddr("glIndexub"))
	if gpIndexub == 0 {
		return errors.New("glIndexub")
	}
	gpIndexubv = uintptr(getProcAddr("glIndexubv"))
	if gpIndexubv == 0 {
		return errors.New("glIndexubv")
	}
	gpIndexxOES = uintptr(getProcAddr("glIndexxOES"))
	gpIndexxvOES = uintptr(getProcAddr("glIndexxvOES"))
	gpInitNames = uintptr(getProcAddr("glInitNames"))
	if gpInitNames == 0 {
		return errors.New("glInitNames")
	}
	gpInsertComponentEXT = uintptr(getProcAddr("glInsertComponentEXT"))
	gpInsertEventMarkerEXT = uintptr(getProcAddr("glInsertEventMarkerEXT"))
	gpInstrumentsBufferSGIX = uintptr(getProcAddr("glInstrumentsBufferSGIX"))
	gpInterleavedArrays = uintptr(getProcAddr("glInterleavedArrays"))
	if gpInterleavedArrays == 0 {
		return errors.New("glInterleavedArrays")
	}
	gpInterpolatePathsNV = uintptr(getProcAddr("glInterpolatePathsNV"))
	gpInvalidateBufferData = uintptr(getProcAddr("glInvalidateBufferData"))
	gpInvalidateBufferSubData = uintptr(getProcAddr("glInvalidateBufferSubData"))
	gpInvalidateFramebuffer = uintptr(getProcAddr("glInvalidateFramebuffer"))
	gpInvalidateNamedFramebufferData = uintptr(getProcAddr("glInvalidateNamedFramebufferData"))
	gpInvalidateNamedFramebufferSubData = uintptr(getProcAddr("glInvalidateNamedFramebufferSubData"))
	gpInvalidateSubFramebuffer = uintptr(getProcAddr("glInvalidateSubFramebuffer"))
	gpInvalidateTexImage = uintptr(getProcAddr("glInvalidateTexImage"))
	gpInvalidateTexSubImage = uintptr(getProcAddr("glInvalidateTexSubImage"))
	gpIsAsyncMarkerSGIX = uintptr(getProcAddr("glIsAsyncMarkerSGIX"))
	gpIsBuffer = uintptr(getProcAddr("glIsBuffer"))
	if gpIsBuffer == 0 {
		return errors.New("glIsBuffer")
	}
	gpIsBufferARB = uintptr(getProcAddr("glIsBufferARB"))
	gpIsBufferResidentNV = uintptr(getProcAddr("glIsBufferResidentNV"))
	gpIsCommandListNV = uintptr(getProcAddr("glIsCommandListNV"))
	gpIsEnabled = uintptr(getProcAddr("glIsEnabled"))
	if gpIsEnabled == 0 {
		return errors.New("glIsEnabled")
	}
	gpIsEnabledIndexedEXT = uintptr(getProcAddr("glIsEnabledIndexedEXT"))
	gpIsFenceAPPLE = uintptr(getProcAddr("glIsFenceAPPLE"))
	gpIsFenceNV = uintptr(getProcAddr("glIsFenceNV"))
	gpIsFramebuffer = uintptr(getProcAddr("glIsFramebuffer"))
	gpIsFramebufferEXT = uintptr(getProcAddr("glIsFramebufferEXT"))
	gpIsImageHandleResidentARB = uintptr(getProcAddr("glIsImageHandleResidentARB"))
	gpIsImageHandleResidentNV = uintptr(getProcAddr("glIsImageHandleResidentNV"))
	gpIsList = uintptr(getProcAddr("glIsList"))
	if gpIsList == 0 {
		return errors.New("glIsList")
	}
	gpIsMemoryObjectEXT = uintptr(getProcAddr("glIsMemoryObjectEXT"))
	gpIsNameAMD = uintptr(getProcAddr("glIsNameAMD"))
	gpIsNamedBufferResidentNV = uintptr(getProcAddr("glIsNamedBufferResidentNV"))
	gpIsNamedStringARB = uintptr(getProcAddr("glIsNamedStringARB"))
	gpIsObjectBufferATI = uintptr(getProcAddr("glIsObjectBufferATI"))
	gpIsOcclusionQueryNV = uintptr(getProcAddr("glIsOcclusionQueryNV"))
	gpIsPathNV = uintptr(getProcAddr("glIsPathNV"))
	gpIsPointInFillPathNV = uintptr(getProcAddr("glIsPointInFillPathNV"))
	gpIsPointInStrokePathNV = uintptr(getProcAddr("glIsPointInStrokePathNV"))
	gpIsProgram = uintptr(getProcAddr("glIsProgram"))
	if gpIsProgram == 0 {
		return errors.New("glIsProgram")
	}
	gpIsProgramARB = uintptr(getProcAddr("glIsProgramARB"))
	gpIsProgramNV = uintptr(getProcAddr("glIsProgramNV"))
	gpIsProgramPipeline = uintptr(getProcAddr("glIsProgramPipeline"))
	gpIsProgramPipelineEXT = uintptr(getProcAddr("glIsProgramPipelineEXT"))
	gpIsQuery = uintptr(getProcAddr("glIsQuery"))
	if gpIsQuery == 0 {
		return errors.New("glIsQuery")
	}
	gpIsQueryARB = uintptr(getProcAddr("glIsQueryARB"))
	gpIsRenderbuffer = uintptr(getProcAddr("glIsRenderbuffer"))
	gpIsRenderbufferEXT = uintptr(getProcAddr("glIsRenderbufferEXT"))
	gpIsSampler = uintptr(getProcAddr("glIsSampler"))
	gpIsSemaphoreEXT = uintptr(getProcAddr("glIsSemaphoreEXT"))
	gpIsShader = uintptr(getProcAddr("glIsShader"))
	if gpIsShader == 0 {
		return errors.New("glIsShader")
	}
	gpIsStateNV = uintptr(getProcAddr("glIsStateNV"))
	gpIsSync = uintptr(getProcAddr("glIsSync"))
	gpIsTexture = uintptr(getProcAddr("glIsTexture"))
	if gpIsTexture == 0 {
		return errors.New("glIsTexture")
	}
	gpIsTextureEXT = uintptr(getProcAddr("glIsTextureEXT"))
	gpIsTextureHandleResidentARB = uintptr(getProcAddr("glIsTextureHandleResidentARB"))
	gpIsTextureHandleResidentNV = uintptr(getProcAddr("glIsTextureHandleResidentNV"))
	gpIsTransformFeedback = uintptr(getProcAddr("glIsTransformFeedback"))
	gpIsTransformFeedbackNV = uintptr(getProcAddr("glIsTransformFeedbackNV"))
	gpIsVariantEnabledEXT = uintptr(getProcAddr("glIsVariantEnabledEXT"))
	gpIsVertexArray = uintptr(getProcAddr("glIsVertexArray"))
	gpIsVertexArrayAPPLE = uintptr(getProcAddr("glIsVertexArrayAPPLE"))
	gpIsVertexAttribEnabledAPPLE = uintptr(getProcAddr("glIsVertexAttribEnabledAPPLE"))
	gpLGPUCopyImageSubDataNVX = uintptr(getProcAddr("glLGPUCopyImageSubDataNVX"))
	gpLGPUInterlockNVX = uintptr(getProcAddr("glLGPUInterlockNVX"))
	gpLGPUNamedBufferSubDataNVX = uintptr(getProcAddr("glLGPUNamedBufferSubDataNVX"))
	gpLabelObjectEXT = uintptr(getProcAddr("glLabelObjectEXT"))
	gpLightEnviSGIX = uintptr(getProcAddr("glLightEnviSGIX"))
	gpLightModelf = uintptr(getProcAddr("glLightModelf"))
	if gpLightModelf == 0 {
		return errors.New("glLightModelf")
	}
	gpLightModelfv = uintptr(getProcAddr("glLightModelfv"))
	if gpLightModelfv == 0 {
		return errors.New("glLightModelfv")
	}
	gpLightModeli = uintptr(getProcAddr("glLightModeli"))
	if gpLightModeli == 0 {
		return errors.New("glLightModeli")
	}
	gpLightModeliv = uintptr(getProcAddr("glLightModeliv"))
	if gpLightModeliv == 0 {
		return errors.New("glLightModeliv")
	}
	gpLightModelxOES = uintptr(getProcAddr("glLightModelxOES"))
	gpLightModelxvOES = uintptr(getProcAddr("glLightModelxvOES"))
	gpLightf = uintptr(getProcAddr("glLightf"))
	if gpLightf == 0 {
		return errors.New("glLightf")
	}
	gpLightfv = uintptr(getProcAddr("glLightfv"))
	if gpLightfv == 0 {
		return errors.New("glLightfv")
	}
	gpLighti = uintptr(getProcAddr("glLighti"))
	if gpLighti == 0 {
		return errors.New("glLighti")
	}
	gpLightiv = uintptr(getProcAddr("glLightiv"))
	if gpLightiv == 0 {
		return errors.New("glLightiv")
	}
	gpLightxOES = uintptr(getProcAddr("glLightxOES"))
	gpLightxvOES = uintptr(getProcAddr("glLightxvOES"))
	gpLineStipple = uintptr(getProcAddr("glLineStipple"))
	if gpLineStipple == 0 {
		return errors.New("glLineStipple")
	}
	gpLineWidth = uintptr(getProcAddr("glLineWidth"))
	if gpLineWidth == 0 {
		return errors.New("glLineWidth")
	}
	gpLineWidthxOES = uintptr(getProcAddr("glLineWidthxOES"))
	gpLinkProgram = uintptr(getProcAddr("glLinkProgram"))
	if gpLinkProgram == 0 {
		return errors.New("glLinkProgram")
	}
	gpLinkProgramARB = uintptr(getProcAddr("glLinkProgramARB"))
	gpListBase = uintptr(getProcAddr("glListBase"))
	if gpListBase == 0 {
		return errors.New("glListBase")
	}
	gpListDrawCommandsStatesClientNV = uintptr(getProcAddr("glListDrawCommandsStatesClientNV"))
	gpListParameterfSGIX = uintptr(getProcAddr("glListParameterfSGIX"))
	gpListParameterfvSGIX = uintptr(getProcAddr("glListParameterfvSGIX"))
	gpListParameteriSGIX = uintptr(getProcAddr("glListParameteriSGIX"))
	gpListParameterivSGIX = uintptr(getProcAddr("glListParameterivSGIX"))
	gpLoadIdentity = uintptr(getProcAddr("glLoadIdentity"))
	if gpLoadIdentity == 0 {
		return errors.New("glLoadIdentity")
	}
	gpLoadIdentityDeformationMapSGIX = uintptr(getProcAddr("glLoadIdentityDeformationMapSGIX"))
	gpLoadMatrixd = uintptr(getProcAddr("glLoadMatrixd"))
	if gpLoadMatrixd == 0 {
		return errors.New("glLoadMatrixd")
	}
	gpLoadMatrixf = uintptr(getProcAddr("glLoadMatrixf"))
	if gpLoadMatrixf == 0 {
		return errors.New("glLoadMatrixf")
	}
	gpLoadMatrixxOES = uintptr(getProcAddr("glLoadMatrixxOES"))
	gpLoadName = uintptr(getProcAddr("glLoadName"))
	if gpLoadName == 0 {
		return errors.New("glLoadName")
	}
	gpLoadProgramNV = uintptr(getProcAddr("glLoadProgramNV"))
	gpLoadTransposeMatrixd = uintptr(getProcAddr("glLoadTransposeMatrixd"))
	if gpLoadTransposeMatrixd == 0 {
		return errors.New("glLoadTransposeMatrixd")
	}
	gpLoadTransposeMatrixdARB = uintptr(getProcAddr("glLoadTransposeMatrixdARB"))
	gpLoadTransposeMatrixf = uintptr(getProcAddr("glLoadTransposeMatrixf"))
	if gpLoadTransposeMatrixf == 0 {
		return errors.New("glLoadTransposeMatrixf")
	}
	gpLoadTransposeMatrixfARB = uintptr(getProcAddr("glLoadTransposeMatrixfARB"))
	gpLoadTransposeMatrixxOES = uintptr(getProcAddr("glLoadTransposeMatrixxOES"))
	gpLockArraysEXT = uintptr(getProcAddr("glLockArraysEXT"))
	gpLogicOp = uintptr(getProcAddr("glLogicOp"))
	if gpLogicOp == 0 {
		return errors.New("glLogicOp")
	}
	gpMakeBufferNonResidentNV = uintptr(getProcAddr("glMakeBufferNonResidentNV"))
	gpMakeBufferResidentNV = uintptr(getProcAddr("glMakeBufferResidentNV"))
	gpMakeImageHandleNonResidentARB = uintptr(getProcAddr("glMakeImageHandleNonResidentARB"))
	gpMakeImageHandleNonResidentNV = uintptr(getProcAddr("glMakeImageHandleNonResidentNV"))
	gpMakeImageHandleResidentARB = uintptr(getProcAddr("glMakeImageHandleResidentARB"))
	gpMakeImageHandleResidentNV = uintptr(getProcAddr("glMakeImageHandleResidentNV"))
	gpMakeNamedBufferNonResidentNV = uintptr(getProcAddr("glMakeNamedBufferNonResidentNV"))
	gpMakeNamedBufferResidentNV = uintptr(getProcAddr("glMakeNamedBufferResidentNV"))
	gpMakeTextureHandleNonResidentARB = uintptr(getProcAddr("glMakeTextureHandleNonResidentARB"))
	gpMakeTextureHandleNonResidentNV = uintptr(getProcAddr("glMakeTextureHandleNonResidentNV"))
	gpMakeTextureHandleResidentARB = uintptr(getProcAddr("glMakeTextureHandleResidentARB"))
	gpMakeTextureHandleResidentNV = uintptr(getProcAddr("glMakeTextureHandleResidentNV"))
	gpMap1d = uintptr(getProcAddr("glMap1d"))
	if gpMap1d == 0 {
		return errors.New("glMap1d")
	}
	gpMap1f = uintptr(getProcAddr("glMap1f"))
	if gpMap1f == 0 {
		return errors.New("glMap1f")
	}
	gpMap1xOES = uintptr(getProcAddr("glMap1xOES"))
	gpMap2d = uintptr(getProcAddr("glMap2d"))
	if gpMap2d == 0 {
		return errors.New("glMap2d")
	}
	gpMap2f = uintptr(getProcAddr("glMap2f"))
	if gpMap2f == 0 {
		return errors.New("glMap2f")
	}
	gpMap2xOES = uintptr(getProcAddr("glMap2xOES"))
	gpMapBuffer = uintptr(getProcAddr("glMapBuffer"))
	if gpMapBuffer == 0 {
		return errors.New("glMapBuffer")
	}
	gpMapBufferARB = uintptr(getProcAddr("glMapBufferARB"))
	gpMapBufferRange = uintptr(getProcAddr("glMapBufferRange"))
	gpMapControlPointsNV = uintptr(getProcAddr("glMapControlPointsNV"))
	gpMapGrid1d = uintptr(getProcAddr("glMapGrid1d"))
	if gpMapGrid1d == 0 {
		return errors.New("glMapGrid1d")
	}
	gpMapGrid1f = uintptr(getProcAddr("glMapGrid1f"))
	if gpMapGrid1f == 0 {
		return errors.New("glMapGrid1f")
	}
	gpMapGrid1xOES = uintptr(getProcAddr("glMapGrid1xOES"))
	gpMapGrid2d = uintptr(getProcAddr("glMapGrid2d"))
	if gpMapGrid2d == 0 {
		return errors.New("glMapGrid2d")
	}
	gpMapGrid2f = uintptr(getProcAddr("glMapGrid2f"))
	if gpMapGrid2f == 0 {
		return errors.New("glMapGrid2f")
	}
	gpMapGrid2xOES = uintptr(getProcAddr("glMapGrid2xOES"))
	gpMapNamedBuffer = uintptr(getProcAddr("glMapNamedBuffer"))
	gpMapNamedBufferEXT = uintptr(getProcAddr("glMapNamedBufferEXT"))
	gpMapNamedBufferRange = uintptr(getProcAddr("glMapNamedBufferRange"))
	gpMapNamedBufferRangeEXT = uintptr(getProcAddr("glMapNamedBufferRangeEXT"))
	gpMapObjectBufferATI = uintptr(getProcAddr("glMapObjectBufferATI"))
	gpMapParameterfvNV = uintptr(getProcAddr("glMapParameterfvNV"))
	gpMapParameterivNV = uintptr(getProcAddr("glMapParameterivNV"))
	gpMapTexture2DINTEL = uintptr(getProcAddr("glMapTexture2DINTEL"))
	gpMapVertexAttrib1dAPPLE = uintptr(getProcAddr("glMapVertexAttrib1dAPPLE"))
	gpMapVertexAttrib1fAPPLE = uintptr(getProcAddr("glMapVertexAttrib1fAPPLE"))
	gpMapVertexAttrib2dAPPLE = uintptr(getProcAddr("glMapVertexAttrib2dAPPLE"))
	gpMapVertexAttrib2fAPPLE = uintptr(getProcAddr("glMapVertexAttrib2fAPPLE"))
	gpMaterialf = uintptr(getProcAddr("glMaterialf"))
	if gpMaterialf == 0 {
		return errors.New("glMaterialf")
	}
	gpMaterialfv = uintptr(getProcAddr("glMaterialfv"))
	if gpMaterialfv == 0 {
		return errors.New("glMaterialfv")
	}
	gpMateriali = uintptr(getProcAddr("glMateriali"))
	if gpMateriali == 0 {
		return errors.New("glMateriali")
	}
	gpMaterialiv = uintptr(getProcAddr("glMaterialiv"))
	if gpMaterialiv == 0 {
		return errors.New("glMaterialiv")
	}
	gpMaterialxOES = uintptr(getProcAddr("glMaterialxOES"))
	gpMaterialxvOES = uintptr(getProcAddr("glMaterialxvOES"))
	gpMatrixFrustumEXT = uintptr(getProcAddr("glMatrixFrustumEXT"))
	gpMatrixIndexPointerARB = uintptr(getProcAddr("glMatrixIndexPointerARB"))
	gpMatrixIndexubvARB = uintptr(getProcAddr("glMatrixIndexubvARB"))
	gpMatrixIndexuivARB = uintptr(getProcAddr("glMatrixIndexuivARB"))
	gpMatrixIndexusvARB = uintptr(getProcAddr("glMatrixIndexusvARB"))
	gpMatrixLoad3x2fNV = uintptr(getProcAddr("glMatrixLoad3x2fNV"))
	gpMatrixLoad3x3fNV = uintptr(getProcAddr("glMatrixLoad3x3fNV"))
	gpMatrixLoadIdentityEXT = uintptr(getProcAddr("glMatrixLoadIdentityEXT"))
	gpMatrixLoadTranspose3x3fNV = uintptr(getProcAddr("glMatrixLoadTranspose3x3fNV"))
	gpMatrixLoadTransposedEXT = uintptr(getProcAddr("glMatrixLoadTransposedEXT"))
	gpMatrixLoadTransposefEXT = uintptr(getProcAddr("glMatrixLoadTransposefEXT"))
	gpMatrixLoaddEXT = uintptr(getProcAddr("glMatrixLoaddEXT"))
	gpMatrixLoadfEXT = uintptr(getProcAddr("glMatrixLoadfEXT"))
	gpMatrixMode = uintptr(getProcAddr("glMatrixMode"))
	if gpMatrixMode == 0 {
		return errors.New("glMatrixMode")
	}
	gpMatrixMult3x2fNV = uintptr(getProcAddr("glMatrixMult3x2fNV"))
	gpMatrixMult3x3fNV = uintptr(getProcAddr("glMatrixMult3x3fNV"))
	gpMatrixMultTranspose3x3fNV = uintptr(getProcAddr("glMatrixMultTranspose3x3fNV"))
	gpMatrixMultTransposedEXT = uintptr(getProcAddr("glMatrixMultTransposedEXT"))
	gpMatrixMultTransposefEXT = uintptr(getProcAddr("glMatrixMultTransposefEXT"))
	gpMatrixMultdEXT = uintptr(getProcAddr("glMatrixMultdEXT"))
	gpMatrixMultfEXT = uintptr(getProcAddr("glMatrixMultfEXT"))
	gpMatrixOrthoEXT = uintptr(getProcAddr("glMatrixOrthoEXT"))
	gpMatrixPopEXT = uintptr(getProcAddr("glMatrixPopEXT"))
	gpMatrixPushEXT = uintptr(getProcAddr("glMatrixPushEXT"))
	gpMatrixRotatedEXT = uintptr(getProcAddr("glMatrixRotatedEXT"))
	gpMatrixRotatefEXT = uintptr(getProcAddr("glMatrixRotatefEXT"))
	gpMatrixScaledEXT = uintptr(getProcAddr("glMatrixScaledEXT"))
	gpMatrixScalefEXT = uintptr(getProcAddr("glMatrixScalefEXT"))
	gpMatrixTranslatedEXT = uintptr(getProcAddr("glMatrixTranslatedEXT"))
	gpMatrixTranslatefEXT = uintptr(getProcAddr("glMatrixTranslatefEXT"))
	gpMaxShaderCompilerThreadsARB = uintptr(getProcAddr("glMaxShaderCompilerThreadsARB"))
	gpMaxShaderCompilerThreadsKHR = uintptr(getProcAddr("glMaxShaderCompilerThreadsKHR"))
	gpMemoryBarrier = uintptr(getProcAddr("glMemoryBarrier"))
	gpMemoryBarrierByRegion = uintptr(getProcAddr("glMemoryBarrierByRegion"))
	gpMemoryBarrierEXT = uintptr(getProcAddr("glMemoryBarrierEXT"))
	gpMemoryObjectParameterivEXT = uintptr(getProcAddr("glMemoryObjectParameterivEXT"))
	gpMinSampleShadingARB = uintptr(getProcAddr("glMinSampleShadingARB"))
	gpMinmaxEXT = uintptr(getProcAddr("glMinmaxEXT"))
	gpMultMatrixd = uintptr(getProcAddr("glMultMatrixd"))
	if gpMultMatrixd == 0 {
		return errors.New("glMultMatrixd")
	}
	gpMultMatrixf = uintptr(getProcAddr("glMultMatrixf"))
	if gpMultMatrixf == 0 {
		return errors.New("glMultMatrixf")
	}
	gpMultMatrixxOES = uintptr(getProcAddr("glMultMatrixxOES"))
	gpMultTransposeMatrixd = uintptr(getProcAddr("glMultTransposeMatrixd"))
	if gpMultTransposeMatrixd == 0 {
		return errors.New("glMultTransposeMatrixd")
	}
	gpMultTransposeMatrixdARB = uintptr(getProcAddr("glMultTransposeMatrixdARB"))
	gpMultTransposeMatrixf = uintptr(getProcAddr("glMultTransposeMatrixf"))
	if gpMultTransposeMatrixf == 0 {
		return errors.New("glMultTransposeMatrixf")
	}
	gpMultTransposeMatrixfARB = uintptr(getProcAddr("glMultTransposeMatrixfARB"))
	gpMultTransposeMatrixxOES = uintptr(getProcAddr("glMultTransposeMatrixxOES"))
	gpMultiDrawArrays = uintptr(getProcAddr("glMultiDrawArrays"))
	if gpMultiDrawArrays == 0 {
		return errors.New("glMultiDrawArrays")
	}
	gpMultiDrawArraysEXT = uintptr(getProcAddr("glMultiDrawArraysEXT"))
	gpMultiDrawArraysIndirect = uintptr(getProcAddr("glMultiDrawArraysIndirect"))
	gpMultiDrawArraysIndirectAMD = uintptr(getProcAddr("glMultiDrawArraysIndirectAMD"))
	gpMultiDrawArraysIndirectBindlessCountNV = uintptr(getProcAddr("glMultiDrawArraysIndirectBindlessCountNV"))
	gpMultiDrawArraysIndirectBindlessNV = uintptr(getProcAddr("glMultiDrawArraysIndirectBindlessNV"))
	gpMultiDrawArraysIndirectCountARB = uintptr(getProcAddr("glMultiDrawArraysIndirectCountARB"))
	gpMultiDrawElementArrayAPPLE = uintptr(getProcAddr("glMultiDrawElementArrayAPPLE"))
	gpMultiDrawElements = uintptr(getProcAddr("glMultiDrawElements"))
	if gpMultiDrawElements == 0 {
		return errors.New("glMultiDrawElements")
	}
	gpMultiDrawElementsBaseVertex = uintptr(getProcAddr("glMultiDrawElementsBaseVertex"))
	gpMultiDrawElementsEXT = uintptr(getProcAddr("glMultiDrawElementsEXT"))
	gpMultiDrawElementsIndirect = uintptr(getProcAddr("glMultiDrawElementsIndirect"))
	gpMultiDrawElementsIndirectAMD = uintptr(getProcAddr("glMultiDrawElementsIndirectAMD"))
	gpMultiDrawElementsIndirectBindlessCountNV = uintptr(getProcAddr("glMultiDrawElementsIndirectBindlessCountNV"))
	gpMultiDrawElementsIndirectBindlessNV = uintptr(getProcAddr("glMultiDrawElementsIndirectBindlessNV"))
	gpMultiDrawElementsIndirectCountARB = uintptr(getProcAddr("glMultiDrawElementsIndirectCountARB"))
	gpMultiDrawRangeElementArrayAPPLE = uintptr(getProcAddr("glMultiDrawRangeElementArrayAPPLE"))
	gpMultiModeDrawArraysIBM = uintptr(getProcAddr("glMultiModeDrawArraysIBM"))
	gpMultiModeDrawElementsIBM = uintptr(getProcAddr("glMultiModeDrawElementsIBM"))
	gpMultiTexBufferEXT = uintptr(getProcAddr("glMultiTexBufferEXT"))
	gpMultiTexCoord1bOES = uintptr(getProcAddr("glMultiTexCoord1bOES"))
	gpMultiTexCoord1bvOES = uintptr(getProcAddr("glMultiTexCoord1bvOES"))
	gpMultiTexCoord1d = uintptr(getProcAddr("glMultiTexCoord1d"))
	if gpMultiTexCoord1d == 0 {
		return errors.New("glMultiTexCoord1d")
	}
	gpMultiTexCoord1dARB = uintptr(getProcAddr("glMultiTexCoord1dARB"))
	gpMultiTexCoord1dv = uintptr(getProcAddr("glMultiTexCoord1dv"))
	if gpMultiTexCoord1dv == 0 {
		return errors.New("glMultiTexCoord1dv")
	}
	gpMultiTexCoord1dvARB = uintptr(getProcAddr("glMultiTexCoord1dvARB"))
	gpMultiTexCoord1f = uintptr(getProcAddr("glMultiTexCoord1f"))
	if gpMultiTexCoord1f == 0 {
		return errors.New("glMultiTexCoord1f")
	}
	gpMultiTexCoord1fARB = uintptr(getProcAddr("glMultiTexCoord1fARB"))
	gpMultiTexCoord1fv = uintptr(getProcAddr("glMultiTexCoord1fv"))
	if gpMultiTexCoord1fv == 0 {
		return errors.New("glMultiTexCoord1fv")
	}
	gpMultiTexCoord1fvARB = uintptr(getProcAddr("glMultiTexCoord1fvARB"))
	gpMultiTexCoord1hNV = uintptr(getProcAddr("glMultiTexCoord1hNV"))
	gpMultiTexCoord1hvNV = uintptr(getProcAddr("glMultiTexCoord1hvNV"))
	gpMultiTexCoord1i = uintptr(getProcAddr("glMultiTexCoord1i"))
	if gpMultiTexCoord1i == 0 {
		return errors.New("glMultiTexCoord1i")
	}
	gpMultiTexCoord1iARB = uintptr(getProcAddr("glMultiTexCoord1iARB"))
	gpMultiTexCoord1iv = uintptr(getProcAddr("glMultiTexCoord1iv"))
	if gpMultiTexCoord1iv == 0 {
		return errors.New("glMultiTexCoord1iv")
	}
	gpMultiTexCoord1ivARB = uintptr(getProcAddr("glMultiTexCoord1ivARB"))
	gpMultiTexCoord1s = uintptr(getProcAddr("glMultiTexCoord1s"))
	if gpMultiTexCoord1s == 0 {
		return errors.New("glMultiTexCoord1s")
	}
	gpMultiTexCoord1sARB = uintptr(getProcAddr("glMultiTexCoord1sARB"))
	gpMultiTexCoord1sv = uintptr(getProcAddr("glMultiTexCoord1sv"))
	if gpMultiTexCoord1sv == 0 {
		return errors.New("glMultiTexCoord1sv")
	}
	gpMultiTexCoord1svARB = uintptr(getProcAddr("glMultiTexCoord1svARB"))
	gpMultiTexCoord1xOES = uintptr(getProcAddr("glMultiTexCoord1xOES"))
	gpMultiTexCoord1xvOES = uintptr(getProcAddr("glMultiTexCoord1xvOES"))
	gpMultiTexCoord2bOES = uintptr(getProcAddr("glMultiTexCoord2bOES"))
	gpMultiTexCoord2bvOES = uintptr(getProcAddr("glMultiTexCoord2bvOES"))
	gpMultiTexCoord2d = uintptr(getProcAddr("glMultiTexCoord2d"))
	if gpMultiTexCoord2d == 0 {
		return errors.New("glMultiTexCoord2d")
	}
	gpMultiTexCoord2dARB = uintptr(getProcAddr("glMultiTexCoord2dARB"))
	gpMultiTexCoord2dv = uintptr(getProcAddr("glMultiTexCoord2dv"))
	if gpMultiTexCoord2dv == 0 {
		return errors.New("glMultiTexCoord2dv")
	}
	gpMultiTexCoord2dvARB = uintptr(getProcAddr("glMultiTexCoord2dvARB"))
	gpMultiTexCoord2f = uintptr(getProcAddr("glMultiTexCoord2f"))
	if gpMultiTexCoord2f == 0 {
		return errors.New("glMultiTexCoord2f")
	}
	gpMultiTexCoord2fARB = uintptr(getProcAddr("glMultiTexCoord2fARB"))
	gpMultiTexCoord2fv = uintptr(getProcAddr("glMultiTexCoord2fv"))
	if gpMultiTexCoord2fv == 0 {
		return errors.New("glMultiTexCoord2fv")
	}
	gpMultiTexCoord2fvARB = uintptr(getProcAddr("glMultiTexCoord2fvARB"))
	gpMultiTexCoord2hNV = uintptr(getProcAddr("glMultiTexCoord2hNV"))
	gpMultiTexCoord2hvNV = uintptr(getProcAddr("glMultiTexCoord2hvNV"))
	gpMultiTexCoord2i = uintptr(getProcAddr("glMultiTexCoord2i"))
	if gpMultiTexCoord2i == 0 {
		return errors.New("glMultiTexCoord2i")
	}
	gpMultiTexCoord2iARB = uintptr(getProcAddr("glMultiTexCoord2iARB"))
	gpMultiTexCoord2iv = uintptr(getProcAddr("glMultiTexCoord2iv"))
	if gpMultiTexCoord2iv == 0 {
		return errors.New("glMultiTexCoord2iv")
	}
	gpMultiTexCoord2ivARB = uintptr(getProcAddr("glMultiTexCoord2ivARB"))
	gpMultiTexCoord2s = uintptr(getProcAddr("glMultiTexCoord2s"))
	if gpMultiTexCoord2s == 0 {
		return errors.New("glMultiTexCoord2s")
	}
	gpMultiTexCoord2sARB = uintptr(getProcAddr("glMultiTexCoord2sARB"))
	gpMultiTexCoord2sv = uintptr(getProcAddr("glMultiTexCoord2sv"))
	if gpMultiTexCoord2sv == 0 {
		return errors.New("glMultiTexCoord2sv")
	}
	gpMultiTexCoord2svARB = uintptr(getProcAddr("glMultiTexCoord2svARB"))
	gpMultiTexCoord2xOES = uintptr(getProcAddr("glMultiTexCoord2xOES"))
	gpMultiTexCoord2xvOES = uintptr(getProcAddr("glMultiTexCoord2xvOES"))
	gpMultiTexCoord3bOES = uintptr(getProcAddr("glMultiTexCoord3bOES"))
	gpMultiTexCoord3bvOES = uintptr(getProcAddr("glMultiTexCoord3bvOES"))
	gpMultiTexCoord3d = uintptr(getProcAddr("glMultiTexCoord3d"))
	if gpMultiTexCoord3d == 0 {
		return errors.New("glMultiTexCoord3d")
	}
	gpMultiTexCoord3dARB = uintptr(getProcAddr("glMultiTexCoord3dARB"))
	gpMultiTexCoord3dv = uintptr(getProcAddr("glMultiTexCoord3dv"))
	if gpMultiTexCoord3dv == 0 {
		return errors.New("glMultiTexCoord3dv")
	}
	gpMultiTexCoord3dvARB = uintptr(getProcAddr("glMultiTexCoord3dvARB"))
	gpMultiTexCoord3f = uintptr(getProcAddr("glMultiTexCoord3f"))
	if gpMultiTexCoord3f == 0 {
		return errors.New("glMultiTexCoord3f")
	}
	gpMultiTexCoord3fARB = uintptr(getProcAddr("glMultiTexCoord3fARB"))
	gpMultiTexCoord3fv = uintptr(getProcAddr("glMultiTexCoord3fv"))
	if gpMultiTexCoord3fv == 0 {
		return errors.New("glMultiTexCoord3fv")
	}
	gpMultiTexCoord3fvARB = uintptr(getProcAddr("glMultiTexCoord3fvARB"))
	gpMultiTexCoord3hNV = uintptr(getProcAddr("glMultiTexCoord3hNV"))
	gpMultiTexCoord3hvNV = uintptr(getProcAddr("glMultiTexCoord3hvNV"))
	gpMultiTexCoord3i = uintptr(getProcAddr("glMultiTexCoord3i"))
	if gpMultiTexCoord3i == 0 {
		return errors.New("glMultiTexCoord3i")
	}
	gpMultiTexCoord3iARB = uintptr(getProcAddr("glMultiTexCoord3iARB"))
	gpMultiTexCoord3iv = uintptr(getProcAddr("glMultiTexCoord3iv"))
	if gpMultiTexCoord3iv == 0 {
		return errors.New("glMultiTexCoord3iv")
	}
	gpMultiTexCoord3ivARB = uintptr(getProcAddr("glMultiTexCoord3ivARB"))
	gpMultiTexCoord3s = uintptr(getProcAddr("glMultiTexCoord3s"))
	if gpMultiTexCoord3s == 0 {
		return errors.New("glMultiTexCoord3s")
	}
	gpMultiTexCoord3sARB = uintptr(getProcAddr("glMultiTexCoord3sARB"))
	gpMultiTexCoord3sv = uintptr(getProcAddr("glMultiTexCoord3sv"))
	if gpMultiTexCoord3sv == 0 {
		return errors.New("glMultiTexCoord3sv")
	}
	gpMultiTexCoord3svARB = uintptr(getProcAddr("glMultiTexCoord3svARB"))
	gpMultiTexCoord3xOES = uintptr(getProcAddr("glMultiTexCoord3xOES"))
	gpMultiTexCoord3xvOES = uintptr(getProcAddr("glMultiTexCoord3xvOES"))
	gpMultiTexCoord4bOES = uintptr(getProcAddr("glMultiTexCoord4bOES"))
	gpMultiTexCoord4bvOES = uintptr(getProcAddr("glMultiTexCoord4bvOES"))
	gpMultiTexCoord4d = uintptr(getProcAddr("glMultiTexCoord4d"))
	if gpMultiTexCoord4d == 0 {
		return errors.New("glMultiTexCoord4d")
	}
	gpMultiTexCoord4dARB = uintptr(getProcAddr("glMultiTexCoord4dARB"))
	gpMultiTexCoord4dv = uintptr(getProcAddr("glMultiTexCoord4dv"))
	if gpMultiTexCoord4dv == 0 {
		return errors.New("glMultiTexCoord4dv")
	}
	gpMultiTexCoord4dvARB = uintptr(getProcAddr("glMultiTexCoord4dvARB"))
	gpMultiTexCoord4f = uintptr(getProcAddr("glMultiTexCoord4f"))
	if gpMultiTexCoord4f == 0 {
		return errors.New("glMultiTexCoord4f")
	}
	gpMultiTexCoord4fARB = uintptr(getProcAddr("glMultiTexCoord4fARB"))
	gpMultiTexCoord4fv = uintptr(getProcAddr("glMultiTexCoord4fv"))
	if gpMultiTexCoord4fv == 0 {
		return errors.New("glMultiTexCoord4fv")
	}
	gpMultiTexCoord4fvARB = uintptr(getProcAddr("glMultiTexCoord4fvARB"))
	gpMultiTexCoord4hNV = uintptr(getProcAddr("glMultiTexCoord4hNV"))
	gpMultiTexCoord4hvNV = uintptr(getProcAddr("glMultiTexCoord4hvNV"))
	gpMultiTexCoord4i = uintptr(getProcAddr("glMultiTexCoord4i"))
	if gpMultiTexCoord4i == 0 {
		return errors.New("glMultiTexCoord4i")
	}
	gpMultiTexCoord4iARB = uintptr(getProcAddr("glMultiTexCoord4iARB"))
	gpMultiTexCoord4iv = uintptr(getProcAddr("glMultiTexCoord4iv"))
	if gpMultiTexCoord4iv == 0 {
		return errors.New("glMultiTexCoord4iv")
	}
	gpMultiTexCoord4ivARB = uintptr(getProcAddr("glMultiTexCoord4ivARB"))
	gpMultiTexCoord4s = uintptr(getProcAddr("glMultiTexCoord4s"))
	if gpMultiTexCoord4s == 0 {
		return errors.New("glMultiTexCoord4s")
	}
	gpMultiTexCoord4sARB = uintptr(getProcAddr("glMultiTexCoord4sARB"))
	gpMultiTexCoord4sv = uintptr(getProcAddr("glMultiTexCoord4sv"))
	if gpMultiTexCoord4sv == 0 {
		return errors.New("glMultiTexCoord4sv")
	}
	gpMultiTexCoord4svARB = uintptr(getProcAddr("glMultiTexCoord4svARB"))
	gpMultiTexCoord4xOES = uintptr(getProcAddr("glMultiTexCoord4xOES"))
	gpMultiTexCoord4xvOES = uintptr(getProcAddr("glMultiTexCoord4xvOES"))
	gpMultiTexCoordPointerEXT = uintptr(getProcAddr("glMultiTexCoordPointerEXT"))
	gpMultiTexEnvfEXT = uintptr(getProcAddr("glMultiTexEnvfEXT"))
	gpMultiTexEnvfvEXT = uintptr(getProcAddr("glMultiTexEnvfvEXT"))
	gpMultiTexEnviEXT = uintptr(getProcAddr("glMultiTexEnviEXT"))
	gpMultiTexEnvivEXT = uintptr(getProcAddr("glMultiTexEnvivEXT"))
	gpMultiTexGendEXT = uintptr(getProcAddr("glMultiTexGendEXT"))
	gpMultiTexGendvEXT = uintptr(getProcAddr("glMultiTexGendvEXT"))
	gpMultiTexGenfEXT = uintptr(getProcAddr("glMultiTexGenfEXT"))
	gpMultiTexGenfvEXT = uintptr(getProcAddr("glMultiTexGenfvEXT"))
	gpMultiTexGeniEXT = uintptr(getProcAddr("glMultiTexGeniEXT"))
	gpMultiTexGenivEXT = uintptr(getProcAddr("glMultiTexGenivEXT"))
	gpMultiTexImage1DEXT = uintptr(getProcAddr("glMultiTexImage1DEXT"))
	gpMultiTexImage2DEXT = uintptr(getProcAddr("glMultiTexImage2DEXT"))
	gpMultiTexImage3DEXT = uintptr(getProcAddr("glMultiTexImage3DEXT"))
	gpMultiTexParameterIivEXT = uintptr(getProcAddr("glMultiTexParameterIivEXT"))
	gpMultiTexParameterIuivEXT = uintptr(getProcAddr("glMultiTexParameterIuivEXT"))
	gpMultiTexParameterfEXT = uintptr(getProcAddr("glMultiTexParameterfEXT"))
	gpMultiTexParameterfvEXT = uintptr(getProcAddr("glMultiTexParameterfvEXT"))
	gpMultiTexParameteriEXT = uintptr(getProcAddr("glMultiTexParameteriEXT"))
	gpMultiTexParameterivEXT = uintptr(getProcAddr("glMultiTexParameterivEXT"))
	gpMultiTexRenderbufferEXT = uintptr(getProcAddr("glMultiTexRenderbufferEXT"))
	gpMultiTexSubImage1DEXT = uintptr(getProcAddr("glMultiTexSubImage1DEXT"))
	gpMultiTexSubImage2DEXT = uintptr(getProcAddr("glMultiTexSubImage2DEXT"))
	gpMultiTexSubImage3DEXT = uintptr(getProcAddr("glMultiTexSubImage3DEXT"))
	gpMulticastBarrierNV = uintptr(getProcAddr("glMulticastBarrierNV"))
	gpMulticastBlitFramebufferNV = uintptr(getProcAddr("glMulticastBlitFramebufferNV"))
	gpMulticastBufferSubDataNV = uintptr(getProcAddr("glMulticastBufferSubDataNV"))
	gpMulticastCopyBufferSubDataNV = uintptr(getProcAddr("glMulticastCopyBufferSubDataNV"))
	gpMulticastCopyImageSubDataNV = uintptr(getProcAddr("glMulticastCopyImageSubDataNV"))
	gpMulticastFramebufferSampleLocationsfvNV = uintptr(getProcAddr("glMulticastFramebufferSampleLocationsfvNV"))
	gpMulticastGetQueryObjecti64vNV = uintptr(getProcAddr("glMulticastGetQueryObjecti64vNV"))
	gpMulticastGetQueryObjectivNV = uintptr(getProcAddr("glMulticastGetQueryObjectivNV"))
	gpMulticastGetQueryObjectui64vNV = uintptr(getProcAddr("glMulticastGetQueryObjectui64vNV"))
	gpMulticastGetQueryObjectuivNV = uintptr(getProcAddr("glMulticastGetQueryObjectuivNV"))
	gpMulticastWaitSyncNV = uintptr(getProcAddr("glMulticastWaitSyncNV"))
	gpNamedBufferData = uintptr(getProcAddr("glNamedBufferData"))
	gpNamedBufferDataEXT = uintptr(getProcAddr("glNamedBufferDataEXT"))
	gpNamedBufferPageCommitmentARB = uintptr(getProcAddr("glNamedBufferPageCommitmentARB"))
	gpNamedBufferPageCommitmentEXT = uintptr(getProcAddr("glNamedBufferPageCommitmentEXT"))
	gpNamedBufferStorage = uintptr(getProcAddr("glNamedBufferStorage"))
	gpNamedBufferStorageEXT = uintptr(getProcAddr("glNamedBufferStorageEXT"))
	gpNamedBufferStorageExternalEXT = uintptr(getProcAddr("glNamedBufferStorageExternalEXT"))
	gpNamedBufferStorageMemEXT = uintptr(getProcAddr("glNamedBufferStorageMemEXT"))
	gpNamedBufferSubData = uintptr(getProcAddr("glNamedBufferSubData"))
	gpNamedBufferSubDataEXT = uintptr(getProcAddr("glNamedBufferSubDataEXT"))
	gpNamedCopyBufferSubDataEXT = uintptr(getProcAddr("glNamedCopyBufferSubDataEXT"))
	gpNamedFramebufferDrawBuffer = uintptr(getProcAddr("glNamedFramebufferDrawBuffer"))
	gpNamedFramebufferDrawBuffers = uintptr(getProcAddr("glNamedFramebufferDrawBuffers"))
	gpNamedFramebufferParameteri = uintptr(getProcAddr("glNamedFramebufferParameteri"))
	gpNamedFramebufferParameteriEXT = uintptr(getProcAddr("glNamedFramebufferParameteriEXT"))
	gpNamedFramebufferReadBuffer = uintptr(getProcAddr("glNamedFramebufferReadBuffer"))
	gpNamedFramebufferRenderbuffer = uintptr(getProcAddr("glNamedFramebufferRenderbuffer"))
	gpNamedFramebufferRenderbufferEXT = uintptr(getProcAddr("glNamedFramebufferRenderbufferEXT"))
	gpNamedFramebufferSampleLocationsfvARB = uintptr(getProcAddr("glNamedFramebufferSampleLocationsfvARB"))
	gpNamedFramebufferSampleLocationsfvNV = uintptr(getProcAddr("glNamedFramebufferSampleLocationsfvNV"))
	gpNamedFramebufferSamplePositionsfvAMD = uintptr(getProcAddr("glNamedFramebufferSamplePositionsfvAMD"))
	gpNamedFramebufferTexture = uintptr(getProcAddr("glNamedFramebufferTexture"))
	gpNamedFramebufferTexture1DEXT = uintptr(getProcAddr("glNamedFramebufferTexture1DEXT"))
	gpNamedFramebufferTexture2DEXT = uintptr(getProcAddr("glNamedFramebufferTexture2DEXT"))
	gpNamedFramebufferTexture3DEXT = uintptr(getProcAddr("glNamedFramebufferTexture3DEXT"))
	gpNamedFramebufferTextureEXT = uintptr(getProcAddr("glNamedFramebufferTextureEXT"))
	gpNamedFramebufferTextureFaceEXT = uintptr(getProcAddr("glNamedFramebufferTextureFaceEXT"))
	gpNamedFramebufferTextureLayer = uintptr(getProcAddr("glNamedFramebufferTextureLayer"))
	gpNamedFramebufferTextureLayerEXT = uintptr(getProcAddr("glNamedFramebufferTextureLayerEXT"))
	gpNamedProgramLocalParameter4dEXT = uintptr(getProcAddr("glNamedProgramLocalParameter4dEXT"))
	gpNamedProgramLocalParameter4dvEXT = uintptr(getProcAddr("glNamedProgramLocalParameter4dvEXT"))
	gpNamedProgramLocalParameter4fEXT = uintptr(getProcAddr("glNamedProgramLocalParameter4fEXT"))
	gpNamedProgramLocalParameter4fvEXT = uintptr(getProcAddr("glNamedProgramLocalParameter4fvEXT"))
	gpNamedProgramLocalParameterI4iEXT = uintptr(getProcAddr("glNamedProgramLocalParameterI4iEXT"))
	gpNamedProgramLocalParameterI4ivEXT = uintptr(getProcAddr("glNamedProgramLocalParameterI4ivEXT"))
	gpNamedProgramLocalParameterI4uiEXT = uintptr(getProcAddr("glNamedProgramLocalParameterI4uiEXT"))
	gpNamedProgramLocalParameterI4uivEXT = uintptr(getProcAddr("glNamedProgramLocalParameterI4uivEXT"))
	gpNamedProgramLocalParameters4fvEXT = uintptr(getProcAddr("glNamedProgramLocalParameters4fvEXT"))
	gpNamedProgramLocalParametersI4ivEXT = uintptr(getProcAddr("glNamedProgramLocalParametersI4ivEXT"))
	gpNamedProgramLocalParametersI4uivEXT = uintptr(getProcAddr("glNamedProgramLocalParametersI4uivEXT"))
	gpNamedProgramStringEXT = uintptr(getProcAddr("glNamedProgramStringEXT"))
	gpNamedRenderbufferStorage = uintptr(getProcAddr("glNamedRenderbufferStorage"))
	gpNamedRenderbufferStorageEXT = uintptr(getProcAddr("glNamedRenderbufferStorageEXT"))
	gpNamedRenderbufferStorageMultisample = uintptr(getProcAddr("glNamedRenderbufferStorageMultisample"))
	gpNamedRenderbufferStorageMultisampleCoverageEXT = uintptr(getProcAddr("glNamedRenderbufferStorageMultisampleCoverageEXT"))
	gpNamedRenderbufferStorageMultisampleEXT = uintptr(getProcAddr("glNamedRenderbufferStorageMultisampleEXT"))
	gpNamedStringARB = uintptr(getProcAddr("glNamedStringARB"))
	gpNewList = uintptr(getProcAddr("glNewList"))
	if gpNewList == 0 {
		return errors.New("glNewList")
	}
	gpNewObjectBufferATI = uintptr(getProcAddr("glNewObjectBufferATI"))
	gpNormal3b = uintptr(getProcAddr("glNormal3b"))
	if gpNormal3b == 0 {
		return errors.New("glNormal3b")
	}
	gpNormal3bv = uintptr(getProcAddr("glNormal3bv"))
	if gpNormal3bv == 0 {
		return errors.New("glNormal3bv")
	}
	gpNormal3d = uintptr(getProcAddr("glNormal3d"))
	if gpNormal3d == 0 {
		return errors.New("glNormal3d")
	}
	gpNormal3dv = uintptr(getProcAddr("glNormal3dv"))
	if gpNormal3dv == 0 {
		return errors.New("glNormal3dv")
	}
	gpNormal3f = uintptr(getProcAddr("glNormal3f"))
	if gpNormal3f == 0 {
		return errors.New("glNormal3f")
	}
	gpNormal3fVertex3fSUN = uintptr(getProcAddr("glNormal3fVertex3fSUN"))
	gpNormal3fVertex3fvSUN = uintptr(getProcAddr("glNormal3fVertex3fvSUN"))
	gpNormal3fv = uintptr(getProcAddr("glNormal3fv"))
	if gpNormal3fv == 0 {
		return errors.New("glNormal3fv")
	}
	gpNormal3hNV = uintptr(getProcAddr("glNormal3hNV"))
	gpNormal3hvNV = uintptr(getProcAddr("glNormal3hvNV"))
	gpNormal3i = uintptr(getProcAddr("glNormal3i"))
	if gpNormal3i == 0 {
		return errors.New("glNormal3i")
	}
	gpNormal3iv = uintptr(getProcAddr("glNormal3iv"))
	if gpNormal3iv == 0 {
		return errors.New("glNormal3iv")
	}
	gpNormal3s = uintptr(getProcAddr("glNormal3s"))
	if gpNormal3s == 0 {
		return errors.New("glNormal3s")
	}
	gpNormal3sv = uintptr(getProcAddr("glNormal3sv"))
	if gpNormal3sv == 0 {
		return errors.New("glNormal3sv")
	}
	gpNormal3xOES = uintptr(getProcAddr("glNormal3xOES"))
	gpNormal3xvOES = uintptr(getProcAddr("glNormal3xvOES"))
	gpNormalFormatNV = uintptr(getProcAddr("glNormalFormatNV"))
	gpNormalPointer = uintptr(getProcAddr("glNormalPointer"))
	if gpNormalPointer == 0 {
		return errors.New("glNormalPointer")
	}
	gpNormalPointerEXT = uintptr(getProcAddr("glNormalPointerEXT"))
	gpNormalPointerListIBM = uintptr(getProcAddr("glNormalPointerListIBM"))
	gpNormalPointervINTEL = uintptr(getProcAddr("glNormalPointervINTEL"))
	gpNormalStream3bATI = uintptr(getProcAddr("glNormalStream3bATI"))
	gpNormalStream3bvATI = uintptr(getProcAddr("glNormalStream3bvATI"))
	gpNormalStream3dATI = uintptr(getProcAddr("glNormalStream3dATI"))
	gpNormalStream3dvATI = uintptr(getProcAddr("glNormalStream3dvATI"))
	gpNormalStream3fATI = uintptr(getProcAddr("glNormalStream3fATI"))
	gpNormalStream3fvATI = uintptr(getProcAddr("glNormalStream3fvATI"))
	gpNormalStream3iATI = uintptr(getProcAddr("glNormalStream3iATI"))
	gpNormalStream3ivATI = uintptr(getProcAddr("glNormalStream3ivATI"))
	gpNormalStream3sATI = uintptr(getProcAddr("glNormalStream3sATI"))
	gpNormalStream3svATI = uintptr(getProcAddr("glNormalStream3svATI"))
	gpObjectLabel = uintptr(getProcAddr("glObjectLabel"))
	gpObjectLabelKHR = uintptr(getProcAddr("glObjectLabelKHR"))
	gpObjectPtrLabel = uintptr(getProcAddr("glObjectPtrLabel"))
	gpObjectPtrLabelKHR = uintptr(getProcAddr("glObjectPtrLabelKHR"))
	gpObjectPurgeableAPPLE = uintptr(getProcAddr("glObjectPurgeableAPPLE"))
	gpObjectUnpurgeableAPPLE = uintptr(getProcAddr("glObjectUnpurgeableAPPLE"))
	gpOrtho = uintptr(getProcAddr("glOrtho"))
	if gpOrtho == 0 {
		return errors.New("glOrtho")
	}
	gpOrthofOES = uintptr(getProcAddr("glOrthofOES"))
	gpOrthoxOES = uintptr(getProcAddr("glOrthoxOES"))
	gpPNTrianglesfATI = uintptr(getProcAddr("glPNTrianglesfATI"))
	gpPNTrianglesiATI = uintptr(getProcAddr("glPNTrianglesiATI"))
	gpPassTexCoordATI = uintptr(getProcAddr("glPassTexCoordATI"))
	gpPassThrough = uintptr(getProcAddr("glPassThrough"))
	if gpPassThrough == 0 {
		return errors.New("glPassThrough")
	}
	gpPassThroughxOES = uintptr(getProcAddr("glPassThroughxOES"))
	gpPatchParameterfv = uintptr(getProcAddr("glPatchParameterfv"))
	gpPatchParameteri = uintptr(getProcAddr("glPatchParameteri"))
	gpPathCommandsNV = uintptr(getProcAddr("glPathCommandsNV"))
	gpPathCoordsNV = uintptr(getProcAddr("glPathCoordsNV"))
	gpPathCoverDepthFuncNV = uintptr(getProcAddr("glPathCoverDepthFuncNV"))
	gpPathDashArrayNV = uintptr(getProcAddr("glPathDashArrayNV"))
	gpPathGlyphIndexArrayNV = uintptr(getProcAddr("glPathGlyphIndexArrayNV"))
	gpPathGlyphIndexRangeNV = uintptr(getProcAddr("glPathGlyphIndexRangeNV"))
	gpPathGlyphRangeNV = uintptr(getProcAddr("glPathGlyphRangeNV"))
	gpPathGlyphsNV = uintptr(getProcAddr("glPathGlyphsNV"))
	gpPathMemoryGlyphIndexArrayNV = uintptr(getProcAddr("glPathMemoryGlyphIndexArrayNV"))
	gpPathParameterfNV = uintptr(getProcAddr("glPathParameterfNV"))
	gpPathParameterfvNV = uintptr(getProcAddr("glPathParameterfvNV"))
	gpPathParameteriNV = uintptr(getProcAddr("glPathParameteriNV"))
	gpPathParameterivNV = uintptr(getProcAddr("glPathParameterivNV"))
	gpPathStencilDepthOffsetNV = uintptr(getProcAddr("glPathStencilDepthOffsetNV"))
	gpPathStencilFuncNV = uintptr(getProcAddr("glPathStencilFuncNV"))
	gpPathStringNV = uintptr(getProcAddr("glPathStringNV"))
	gpPathSubCommandsNV = uintptr(getProcAddr("glPathSubCommandsNV"))
	gpPathSubCoordsNV = uintptr(getProcAddr("glPathSubCoordsNV"))
	gpPauseTransformFeedback = uintptr(getProcAddr("glPauseTransformFeedback"))
	gpPauseTransformFeedbackNV = uintptr(getProcAddr("glPauseTransformFeedbackNV"))
	gpPixelDataRangeNV = uintptr(getProcAddr("glPixelDataRangeNV"))
	gpPixelMapfv = uintptr(getProcAddr("glPixelMapfv"))
	if gpPixelMapfv == 0 {
		return errors.New("glPixelMapfv")
	}
	gpPixelMapuiv = uintptr(getProcAddr("glPixelMapuiv"))
	if gpPixelMapuiv == 0 {
		return errors.New("glPixelMapuiv")
	}
	gpPixelMapusv = uintptr(getProcAddr("glPixelMapusv"))
	if gpPixelMapusv == 0 {
		return errors.New("glPixelMapusv")
	}
	gpPixelMapx = uintptr(getProcAddr("glPixelMapx"))
	gpPixelStoref = uintptr(getProcAddr("glPixelStoref"))
	if gpPixelStoref == 0 {
		return errors.New("glPixelStoref")
	}
	gpPixelStorei = uintptr(getProcAddr("glPixelStorei"))
	if gpPixelStorei == 0 {
		return errors.New("glPixelStorei")
	}
	gpPixelStorex = uintptr(getProcAddr("glPixelStorex"))
	gpPixelTexGenParameterfSGIS = uintptr(getProcAddr("glPixelTexGenParameterfSGIS"))
	gpPixelTexGenParameterfvSGIS = uintptr(getProcAddr("glPixelTexGenParameterfvSGIS"))
	gpPixelTexGenParameteriSGIS = uintptr(getProcAddr("glPixelTexGenParameteriSGIS"))
	gpPixelTexGenParameterivSGIS = uintptr(getProcAddr("glPixelTexGenParameterivSGIS"))
	gpPixelTexGenSGIX = uintptr(getProcAddr("glPixelTexGenSGIX"))
	gpPixelTransferf = uintptr(getProcAddr("glPixelTransferf"))
	if gpPixelTransferf == 0 {
		return errors.New("glPixelTransferf")
	}
	gpPixelTransferi = uintptr(getProcAddr("glPixelTransferi"))
	if gpPixelTransferi == 0 {
		return errors.New("glPixelTransferi")
	}
	gpPixelTransferxOES = uintptr(getProcAddr("glPixelTransferxOES"))
	gpPixelTransformParameterfEXT = uintptr(getProcAddr("glPixelTransformParameterfEXT"))
	gpPixelTransformParameterfvEXT = uintptr(getProcAddr("glPixelTransformParameterfvEXT"))
	gpPixelTransformParameteriEXT = uintptr(getProcAddr("glPixelTransformParameteriEXT"))
	gpPixelTransformParameterivEXT = uintptr(getProcAddr("glPixelTransformParameterivEXT"))
	gpPixelZoom = uintptr(getProcAddr("glPixelZoom"))
	if gpPixelZoom == 0 {
		return errors.New("glPixelZoom")
	}
	gpPixelZoomxOES = uintptr(getProcAddr("glPixelZoomxOES"))
	gpPointAlongPathNV = uintptr(getProcAddr("glPointAlongPathNV"))
	gpPointParameterf = uintptr(getProcAddr("glPointParameterf"))
	if gpPointParameterf == 0 {
		return errors.New("glPointParameterf")
	}
	gpPointParameterfARB = uintptr(getProcAddr("glPointParameterfARB"))
	gpPointParameterfEXT = uintptr(getProcAddr("glPointParameterfEXT"))
	gpPointParameterfSGIS = uintptr(getProcAddr("glPointParameterfSGIS"))
	gpPointParameterfv = uintptr(getProcAddr("glPointParameterfv"))
	if gpPointParameterfv == 0 {
		return errors.New("glPointParameterfv")
	}
	gpPointParameterfvARB = uintptr(getProcAddr("glPointParameterfvARB"))
	gpPointParameterfvEXT = uintptr(getProcAddr("glPointParameterfvEXT"))
	gpPointParameterfvSGIS = uintptr(getProcAddr("glPointParameterfvSGIS"))
	gpPointParameteri = uintptr(getProcAddr("glPointParameteri"))
	if gpPointParameteri == 0 {
		return errors.New("glPointParameteri")
	}
	gpPointParameteriNV = uintptr(getProcAddr("glPointParameteriNV"))
	gpPointParameteriv = uintptr(getProcAddr("glPointParameteriv"))
	if gpPointParameteriv == 0 {
		return errors.New("glPointParameteriv")
	}
	gpPointParameterivNV = uintptr(getProcAddr("glPointParameterivNV"))
	gpPointParameterxOES = uintptr(getProcAddr("glPointParameterxOES"))
	gpPointParameterxvOES = uintptr(getProcAddr("glPointParameterxvOES"))
	gpPointSize = uintptr(getProcAddr("glPointSize"))
	if gpPointSize == 0 {
		return errors.New("glPointSize")
	}
	gpPointSizexOES = uintptr(getProcAddr("glPointSizexOES"))
	gpPollAsyncSGIX = uintptr(getProcAddr("glPollAsyncSGIX"))
	gpPollInstrumentsSGIX = uintptr(getProcAddr("glPollInstrumentsSGIX"))
	gpPolygonMode = uintptr(getProcAddr("glPolygonMode"))
	if gpPolygonMode == 0 {
		return errors.New("glPolygonMode")
	}
	gpPolygonOffset = uintptr(getProcAddr("glPolygonOffset"))
	if gpPolygonOffset == 0 {
		return errors.New("glPolygonOffset")
	}
	gpPolygonOffsetClamp = uintptr(getProcAddr("glPolygonOffsetClamp"))
	gpPolygonOffsetClampEXT = uintptr(getProcAddr("glPolygonOffsetClampEXT"))
	gpPolygonOffsetEXT = uintptr(getProcAddr("glPolygonOffsetEXT"))
	gpPolygonOffsetxOES = uintptr(getProcAddr("glPolygonOffsetxOES"))
	gpPolygonStipple = uintptr(getProcAddr("glPolygonStipple"))
	if gpPolygonStipple == 0 {
		return errors.New("glPolygonStipple")
	}
	gpPopAttrib = uintptr(getProcAddr("glPopAttrib"))
	if gpPopAttrib == 0 {
		return errors.New("glPopAttrib")
	}
	gpPopClientAttrib = uintptr(getProcAddr("glPopClientAttrib"))
	if gpPopClientAttrib == 0 {
		return errors.New("glPopClientAttrib")
	}
	gpPopDebugGroup = uintptr(getProcAddr("glPopDebugGroup"))
	gpPopDebugGroupKHR = uintptr(getProcAddr("glPopDebugGroupKHR"))
	gpPopGroupMarkerEXT = uintptr(getProcAddr("glPopGroupMarkerEXT"))
	gpPopMatrix = uintptr(getProcAddr("glPopMatrix"))
	if gpPopMatrix == 0 {
		return errors.New("glPopMatrix")
	}
	gpPopName = uintptr(getProcAddr("glPopName"))
	if gpPopName == 0 {
		return errors.New("glPopName")
	}
	gpPresentFrameDualFillNV = uintptr(getProcAddr("glPresentFrameDualFillNV"))
	gpPresentFrameKeyedNV = uintptr(getProcAddr("glPresentFrameKeyedNV"))
	gpPrimitiveBoundingBoxARB = uintptr(getProcAddr("glPrimitiveBoundingBoxARB"))
	gpPrimitiveRestartIndexNV = uintptr(getProcAddr("glPrimitiveRestartIndexNV"))
	gpPrimitiveRestartNV = uintptr(getProcAddr("glPrimitiveRestartNV"))
	gpPrioritizeTextures = uintptr(getProcAddr("glPrioritizeTextures"))
	if gpPrioritizeTextures == 0 {
		return errors.New("glPrioritizeTextures")
	}
	gpPrioritizeTexturesEXT = uintptr(getProcAddr("glPrioritizeTexturesEXT"))
	gpPrioritizeTexturesxOES = uintptr(getProcAddr("glPrioritizeTexturesxOES"))
	gpProgramBinary = uintptr(getProcAddr("glProgramBinary"))
	gpProgramBufferParametersIivNV = uintptr(getProcAddr("glProgramBufferParametersIivNV"))
	gpProgramBufferParametersIuivNV = uintptr(getProcAddr("glProgramBufferParametersIuivNV"))
	gpProgramBufferParametersfvNV = uintptr(getProcAddr("glProgramBufferParametersfvNV"))
	gpProgramEnvParameter4dARB = uintptr(getProcAddr("glProgramEnvParameter4dARB"))
	gpProgramEnvParameter4dvARB = uintptr(getProcAddr("glProgramEnvParameter4dvARB"))
	gpProgramEnvParameter4fARB = uintptr(getProcAddr("glProgramEnvParameter4fARB"))
	gpProgramEnvParameter4fvARB = uintptr(getProcAddr("glProgramEnvParameter4fvARB"))
	gpProgramEnvParameterI4iNV = uintptr(getProcAddr("glProgramEnvParameterI4iNV"))
	gpProgramEnvParameterI4ivNV = uintptr(getProcAddr("glProgramEnvParameterI4ivNV"))
	gpProgramEnvParameterI4uiNV = uintptr(getProcAddr("glProgramEnvParameterI4uiNV"))
	gpProgramEnvParameterI4uivNV = uintptr(getProcAddr("glProgramEnvParameterI4uivNV"))
	gpProgramEnvParameters4fvEXT = uintptr(getProcAddr("glProgramEnvParameters4fvEXT"))
	gpProgramEnvParametersI4ivNV = uintptr(getProcAddr("glProgramEnvParametersI4ivNV"))
	gpProgramEnvParametersI4uivNV = uintptr(getProcAddr("glProgramEnvParametersI4uivNV"))
	gpProgramLocalParameter4dARB = uintptr(getProcAddr("glProgramLocalParameter4dARB"))
	gpProgramLocalParameter4dvARB = uintptr(getProcAddr("glProgramLocalParameter4dvARB"))
	gpProgramLocalParameter4fARB = uintptr(getProcAddr("glProgramLocalParameter4fARB"))
	gpProgramLocalParameter4fvARB = uintptr(getProcAddr("glProgramLocalParameter4fvARB"))
	gpProgramLocalParameterI4iNV = uintptr(getProcAddr("glProgramLocalParameterI4iNV"))
	gpProgramLocalParameterI4ivNV = uintptr(getProcAddr("glProgramLocalParameterI4ivNV"))
	gpProgramLocalParameterI4uiNV = uintptr(getProcAddr("glProgramLocalParameterI4uiNV"))
	gpProgramLocalParameterI4uivNV = uintptr(getProcAddr("glProgramLocalParameterI4uivNV"))
	gpProgramLocalParameters4fvEXT = uintptr(getProcAddr("glProgramLocalParameters4fvEXT"))
	gpProgramLocalParametersI4ivNV = uintptr(getProcAddr("glProgramLocalParametersI4ivNV"))
	gpProgramLocalParametersI4uivNV = uintptr(getProcAddr("glProgramLocalParametersI4uivNV"))
	gpProgramNamedParameter4dNV = uintptr(getProcAddr("glProgramNamedParameter4dNV"))
	gpProgramNamedParameter4dvNV = uintptr(getProcAddr("glProgramNamedParameter4dvNV"))
	gpProgramNamedParameter4fNV = uintptr(getProcAddr("glProgramNamedParameter4fNV"))
	gpProgramNamedParameter4fvNV = uintptr(getProcAddr("glProgramNamedParameter4fvNV"))
	gpProgramParameter4dNV = uintptr(getProcAddr("glProgramParameter4dNV"))
	gpProgramParameter4dvNV = uintptr(getProcAddr("glProgramParameter4dvNV"))
	gpProgramParameter4fNV = uintptr(getProcAddr("glProgramParameter4fNV"))
	gpProgramParameter4fvNV = uintptr(getProcAddr("glProgramParameter4fvNV"))
	gpProgramParameteri = uintptr(getProcAddr("glProgramParameteri"))
	gpProgramParameteriARB = uintptr(getProcAddr("glProgramParameteriARB"))
	gpProgramParameteriEXT = uintptr(getProcAddr("glProgramParameteriEXT"))
	gpProgramParameters4dvNV = uintptr(getProcAddr("glProgramParameters4dvNV"))
	gpProgramParameters4fvNV = uintptr(getProcAddr("glProgramParameters4fvNV"))
	gpProgramPathFragmentInputGenNV = uintptr(getProcAddr("glProgramPathFragmentInputGenNV"))
	gpProgramStringARB = uintptr(getProcAddr("glProgramStringARB"))
	gpProgramSubroutineParametersuivNV = uintptr(getProcAddr("glProgramSubroutineParametersuivNV"))
	gpProgramUniform1d = uintptr(getProcAddr("glProgramUniform1d"))
	gpProgramUniform1dEXT = uintptr(getProcAddr("glProgramUniform1dEXT"))
	gpProgramUniform1dv = uintptr(getProcAddr("glProgramUniform1dv"))
	gpProgramUniform1dvEXT = uintptr(getProcAddr("glProgramUniform1dvEXT"))
	gpProgramUniform1f = uintptr(getProcAddr("glProgramUniform1f"))
	gpProgramUniform1fEXT = uintptr(getProcAddr("glProgramUniform1fEXT"))
	gpProgramUniform1fv = uintptr(getProcAddr("glProgramUniform1fv"))
	gpProgramUniform1fvEXT = uintptr(getProcAddr("glProgramUniform1fvEXT"))
	gpProgramUniform1i = uintptr(getProcAddr("glProgramUniform1i"))
	gpProgramUniform1i64ARB = uintptr(getProcAddr("glProgramUniform1i64ARB"))
	gpProgramUniform1i64NV = uintptr(getProcAddr("glProgramUniform1i64NV"))
	gpProgramUniform1i64vARB = uintptr(getProcAddr("glProgramUniform1i64vARB"))
	gpProgramUniform1i64vNV = uintptr(getProcAddr("glProgramUniform1i64vNV"))
	gpProgramUniform1iEXT = uintptr(getProcAddr("glProgramUniform1iEXT"))
	gpProgramUniform1iv = uintptr(getProcAddr("glProgramUniform1iv"))
	gpProgramUniform1ivEXT = uintptr(getProcAddr("glProgramUniform1ivEXT"))
	gpProgramUniform1ui = uintptr(getProcAddr("glProgramUniform1ui"))
	gpProgramUniform1ui64ARB = uintptr(getProcAddr("glProgramUniform1ui64ARB"))
	gpProgramUniform1ui64NV = uintptr(getProcAddr("glProgramUniform1ui64NV"))
	gpProgramUniform1ui64vARB = uintptr(getProcAddr("glProgramUniform1ui64vARB"))
	gpProgramUniform1ui64vNV = uintptr(getProcAddr("glProgramUniform1ui64vNV"))
	gpProgramUniform1uiEXT = uintptr(getProcAddr("glProgramUniform1uiEXT"))
	gpProgramUniform1uiv = uintptr(getProcAddr("glProgramUniform1uiv"))
	gpProgramUniform1uivEXT = uintptr(getProcAddr("glProgramUniform1uivEXT"))
	gpProgramUniform2d = uintptr(getProcAddr("glProgramUniform2d"))
	gpProgramUniform2dEXT = uintptr(getProcAddr("glProgramUniform2dEXT"))
	gpProgramUniform2dv = uintptr(getProcAddr("glProgramUniform2dv"))
	gpProgramUniform2dvEXT = uintptr(getProcAddr("glProgramUniform2dvEXT"))
	gpProgramUniform2f = uintptr(getProcAddr("glProgramUniform2f"))
	gpProgramUniform2fEXT = uintptr(getProcAddr("glProgramUniform2fEXT"))
	gpProgramUniform2fv = uintptr(getProcAddr("glProgramUniform2fv"))
	gpProgramUniform2fvEXT = uintptr(getProcAddr("glProgramUniform2fvEXT"))
	gpProgramUniform2i = uintptr(getProcAddr("glProgramUniform2i"))
	gpProgramUniform2i64ARB = uintptr(getProcAddr("glProgramUniform2i64ARB"))
	gpProgramUniform2i64NV = uintptr(getProcAddr("glProgramUniform2i64NV"))
	gpProgramUniform2i64vARB = uintptr(getProcAddr("glProgramUniform2i64vARB"))
	gpProgramUniform2i64vNV = uintptr(getProcAddr("glProgramUniform2i64vNV"))
	gpProgramUniform2iEXT = uintptr(getProcAddr("glProgramUniform2iEXT"))
	gpProgramUniform2iv = uintptr(getProcAddr("glProgramUniform2iv"))
	gpProgramUniform2ivEXT = uintptr(getProcAddr("glProgramUniform2ivEXT"))
	gpProgramUniform2ui = uintptr(getProcAddr("glProgramUniform2ui"))
	gpProgramUniform2ui64ARB = uintptr(getProcAddr("glProgramUniform2ui64ARB"))
	gpProgramUniform2ui64NV = uintptr(getProcAddr("glProgramUniform2ui64NV"))
	gpProgramUniform2ui64vARB = uintptr(getProcAddr("glProgramUniform2ui64vARB"))
	gpProgramUniform2ui64vNV = uintptr(getProcAddr("glProgramUniform2ui64vNV"))
	gpProgramUniform2uiEXT = uintptr(getProcAddr("glProgramUniform2uiEXT"))
	gpProgramUniform2uiv = uintptr(getProcAddr("glProgramUniform2uiv"))
	gpProgramUniform2uivEXT = uintptr(getProcAddr("glProgramUniform2uivEXT"))
	gpProgramUniform3d = uintptr(getProcAddr("glProgramUniform3d"))
	gpProgramUniform3dEXT = uintptr(getProcAddr("glProgramUniform3dEXT"))
	gpProgramUniform3dv = uintptr(getProcAddr("glProgramUniform3dv"))
	gpProgramUniform3dvEXT = uintptr(getProcAddr("glProgramUniform3dvEXT"))
	gpProgramUniform3f = uintptr(getProcAddr("glProgramUniform3f"))
	gpProgramUniform3fEXT = uintptr(getProcAddr("glProgramUniform3fEXT"))
	gpProgramUniform3fv = uintptr(getProcAddr("glProgramUniform3fv"))
	gpProgramUniform3fvEXT = uintptr(getProcAddr("glProgramUniform3fvEXT"))
	gpProgramUniform3i = uintptr(getProcAddr("glProgramUniform3i"))
	gpProgramUniform3i64ARB = uintptr(getProcAddr("glProgramUniform3i64ARB"))
	gpProgramUniform3i64NV = uintptr(getProcAddr("glProgramUniform3i64NV"))
	gpProgramUniform3i64vARB = uintptr(getProcAddr("glProgramUniform3i64vARB"))
	gpProgramUniform3i64vNV = uintptr(getProcAddr("glProgramUniform3i64vNV"))
	gpProgramUniform3iEXT = uintptr(getProcAddr("glProgramUniform3iEXT"))
	gpProgramUniform3iv = uintptr(getProcAddr("glProgramUniform3iv"))
	gpProgramUniform3ivEXT = uintptr(getProcAddr("glProgramUniform3ivEXT"))
	gpProgramUniform3ui = uintptr(getProcAddr("glProgramUniform3ui"))
	gpProgramUniform3ui64ARB = uintptr(getProcAddr("glProgramUniform3ui64ARB"))
	gpProgramUniform3ui64NV = uintptr(getProcAddr("glProgramUniform3ui64NV"))
	gpProgramUniform3ui64vARB = uintptr(getProcAddr("glProgramUniform3ui64vARB"))
	gpProgramUniform3ui64vNV = uintptr(getProcAddr("glProgramUniform3ui64vNV"))
	gpProgramUniform3uiEXT = uintptr(getProcAddr("glProgramUniform3uiEXT"))
	gpProgramUniform3uiv = uintptr(getProcAddr("glProgramUniform3uiv"))
	gpProgramUniform3uivEXT = uintptr(getProcAddr("glProgramUniform3uivEXT"))
	gpProgramUniform4d = uintptr(getProcAddr("glProgramUniform4d"))
	gpProgramUniform4dEXT = uintptr(getProcAddr("glProgramUniform4dEXT"))
	gpProgramUniform4dv = uintptr(getProcAddr("glProgramUniform4dv"))
	gpProgramUniform4dvEXT = uintptr(getProcAddr("glProgramUniform4dvEXT"))
	gpProgramUniform4f = uintptr(getProcAddr("glProgramUniform4f"))
	gpProgramUniform4fEXT = uintptr(getProcAddr("glProgramUniform4fEXT"))
	gpProgramUniform4fv = uintptr(getProcAddr("glProgramUniform4fv"))
	gpProgramUniform4fvEXT = uintptr(getProcAddr("glProgramUniform4fvEXT"))
	gpProgramUniform4i = uintptr(getProcAddr("glProgramUniform4i"))
	gpProgramUniform4i64ARB = uintptr(getProcAddr("glProgramUniform4i64ARB"))
	gpProgramUniform4i64NV = uintptr(getProcAddr("glProgramUniform4i64NV"))
	gpProgramUniform4i64vARB = uintptr(getProcAddr("glProgramUniform4i64vARB"))
	gpProgramUniform4i64vNV = uintptr(getProcAddr("glProgramUniform4i64vNV"))
	gpProgramUniform4iEXT = uintptr(getProcAddr("glProgramUniform4iEXT"))
	gpProgramUniform4iv = uintptr(getProcAddr("glProgramUniform4iv"))
	gpProgramUniform4ivEXT = uintptr(getProcAddr("glProgramUniform4ivEXT"))
	gpProgramUniform4ui = uintptr(getProcAddr("glProgramUniform4ui"))
	gpProgramUniform4ui64ARB = uintptr(getProcAddr("glProgramUniform4ui64ARB"))
	gpProgramUniform4ui64NV = uintptr(getProcAddr("glProgramUniform4ui64NV"))
	gpProgramUniform4ui64vARB = uintptr(getProcAddr("glProgramUniform4ui64vARB"))
	gpProgramUniform4ui64vNV = uintptr(getProcAddr("glProgramUniform4ui64vNV"))
	gpProgramUniform4uiEXT = uintptr(getProcAddr("glProgramUniform4uiEXT"))
	gpProgramUniform4uiv = uintptr(getProcAddr("glProgramUniform4uiv"))
	gpProgramUniform4uivEXT = uintptr(getProcAddr("glProgramUniform4uivEXT"))
	gpProgramUniformHandleui64ARB = uintptr(getProcAddr("glProgramUniformHandleui64ARB"))
	gpProgramUniformHandleui64NV = uintptr(getProcAddr("glProgramUniformHandleui64NV"))
	gpProgramUniformHandleui64vARB = uintptr(getProcAddr("glProgramUniformHandleui64vARB"))
	gpProgramUniformHandleui64vNV = uintptr(getProcAddr("glProgramUniformHandleui64vNV"))
	gpProgramUniformMatrix2dv = uintptr(getProcAddr("glProgramUniformMatrix2dv"))
	gpProgramUniformMatrix2dvEXT = uintptr(getProcAddr("glProgramUniformMatrix2dvEXT"))
	gpProgramUniformMatrix2fv = uintptr(getProcAddr("glProgramUniformMatrix2fv"))
	gpProgramUniformMatrix2fvEXT = uintptr(getProcAddr("glProgramUniformMatrix2fvEXT"))
	gpProgramUniformMatrix2x3dv = uintptr(getProcAddr("glProgramUniformMatrix2x3dv"))
	gpProgramUniformMatrix2x3dvEXT = uintptr(getProcAddr("glProgramUniformMatrix2x3dvEXT"))
	gpProgramUniformMatrix2x3fv = uintptr(getProcAddr("glProgramUniformMatrix2x3fv"))
	gpProgramUniformMatrix2x3fvEXT = uintptr(getProcAddr("glProgramUniformMatrix2x3fvEXT"))
	gpProgramUniformMatrix2x4dv = uintptr(getProcAddr("glProgramUniformMatrix2x4dv"))
	gpProgramUniformMatrix2x4dvEXT = uintptr(getProcAddr("glProgramUniformMatrix2x4dvEXT"))
	gpProgramUniformMatrix2x4fv = uintptr(getProcAddr("glProgramUniformMatrix2x4fv"))
	gpProgramUniformMatrix2x4fvEXT = uintptr(getProcAddr("glProgramUniformMatrix2x4fvEXT"))
	gpProgramUniformMatrix3dv = uintptr(getProcAddr("glProgramUniformMatrix3dv"))
	gpProgramUniformMatrix3dvEXT = uintptr(getProcAddr("glProgramUniformMatrix3dvEXT"))
	gpProgramUniformMatrix3fv = uintptr(getProcAddr("glProgramUniformMatrix3fv"))
	gpProgramUniformMatrix3fvEXT = uintptr(getProcAddr("glProgramUniformMatrix3fvEXT"))
	gpProgramUniformMatrix3x2dv = uintptr(getProcAddr("glProgramUniformMatrix3x2dv"))
	gpProgramUniformMatrix3x2dvEXT = uintptr(getProcAddr("glProgramUniformMatrix3x2dvEXT"))
	gpProgramUniformMatrix3x2fv = uintptr(getProcAddr("glProgramUniformMatrix3x2fv"))
	gpProgramUniformMatrix3x2fvEXT = uintptr(getProcAddr("glProgramUniformMatrix3x2fvEXT"))
	gpProgramUniformMatrix3x4dv = uintptr(getProcAddr("glProgramUniformMatrix3x4dv"))
	gpProgramUniformMatrix3x4dvEXT = uintptr(getProcAddr("glProgramUniformMatrix3x4dvEXT"))
	gpProgramUniformMatrix3x4fv = uintptr(getProcAddr("glProgramUniformMatrix3x4fv"))
	gpProgramUniformMatrix3x4fvEXT = uintptr(getProcAddr("glProgramUniformMatrix3x4fvEXT"))
	gpProgramUniformMatrix4dv = uintptr(getProcAddr("glProgramUniformMatrix4dv"))
	gpProgramUniformMatrix4dvEXT = uintptr(getProcAddr("glProgramUniformMatrix4dvEXT"))
	gpProgramUniformMatrix4fv = uintptr(getProcAddr("glProgramUniformMatrix4fv"))
	gpProgramUniformMatrix4fvEXT = uintptr(getProcAddr("glProgramUniformMatrix4fvEXT"))
	gpProgramUniformMatrix4x2dv = uintptr(getProcAddr("glProgramUniformMatrix4x2dv"))
	gpProgramUniformMatrix4x2dvEXT = uintptr(getProcAddr("glProgramUniformMatrix4x2dvEXT"))
	gpProgramUniformMatrix4x2fv = uintptr(getProcAddr("glProgramUniformMatrix4x2fv"))
	gpProgramUniformMatrix4x2fvEXT = uintptr(getProcAddr("glProgramUniformMatrix4x2fvEXT"))
	gpProgramUniformMatrix4x3dv = uintptr(getProcAddr("glProgramUniformMatrix4x3dv"))
	gpProgramUniformMatrix4x3dvEXT = uintptr(getProcAddr("glProgramUniformMatrix4x3dvEXT"))
	gpProgramUniformMatrix4x3fv = uintptr(getProcAddr("glProgramUniformMatrix4x3fv"))
	gpProgramUniformMatrix4x3fvEXT = uintptr(getProcAddr("glProgramUniformMatrix4x3fvEXT"))
	gpProgramUniformui64NV = uintptr(getProcAddr("glProgramUniformui64NV"))
	gpProgramUniformui64vNV = uintptr(getProcAddr("glProgramUniformui64vNV"))
	gpProgramVertexLimitNV = uintptr(getProcAddr("glProgramVertexLimitNV"))
	gpProvokingVertex = uintptr(getProcAddr("glProvokingVertex"))
	gpProvokingVertexEXT = uintptr(getProcAddr("glProvokingVertexEXT"))
	gpPushAttrib = uintptr(getProcAddr("glPushAttrib"))
	if gpPushAttrib == 0 {
		return errors.New("glPushAttrib")
	}
	gpPushClientAttrib = uintptr(getProcAddr("glPushClientAttrib"))
	if gpPushClientAttrib == 0 {
		return errors.New("glPushClientAttrib")
	}
	gpPushClientAttribDefaultEXT = uintptr(getProcAddr("glPushClientAttribDefaultEXT"))
	gpPushDebugGroup = uintptr(getProcAddr("glPushDebugGroup"))
	gpPushDebugGroupKHR = uintptr(getProcAddr("glPushDebugGroupKHR"))
	gpPushGroupMarkerEXT = uintptr(getProcAddr("glPushGroupMarkerEXT"))
	gpPushMatrix = uintptr(getProcAddr("glPushMatrix"))
	if gpPushMatrix == 0 {
		return errors.New("glPushMatrix")
	}
	gpPushName = uintptr(getProcAddr("glPushName"))
	if gpPushName == 0 {
		return errors.New("glPushName")
	}
	gpQueryCounter = uintptr(getProcAddr("glQueryCounter"))
	gpQueryMatrixxOES = uintptr(getProcAddr("glQueryMatrixxOES"))
	gpQueryObjectParameteruiAMD = uintptr(getProcAddr("glQueryObjectParameteruiAMD"))
	gpQueryResourceNV = uintptr(getProcAddr("glQueryResourceNV"))
	gpQueryResourceTagNV = uintptr(getProcAddr("glQueryResourceTagNV"))
	gpRasterPos2d = uintptr(getProcAddr("glRasterPos2d"))
	if gpRasterPos2d == 0 {
		return errors.New("glRasterPos2d")
	}
	gpRasterPos2dv = uintptr(getProcAddr("glRasterPos2dv"))
	if gpRasterPos2dv == 0 {
		return errors.New("glRasterPos2dv")
	}
	gpRasterPos2f = uintptr(getProcAddr("glRasterPos2f"))
	if gpRasterPos2f == 0 {
		return errors.New("glRasterPos2f")
	}
	gpRasterPos2fv = uintptr(getProcAddr("glRasterPos2fv"))
	if gpRasterPos2fv == 0 {
		return errors.New("glRasterPos2fv")
	}
	gpRasterPos2i = uintptr(getProcAddr("glRasterPos2i"))
	if gpRasterPos2i == 0 {
		return errors.New("glRasterPos2i")
	}
	gpRasterPos2iv = uintptr(getProcAddr("glRasterPos2iv"))
	if gpRasterPos2iv == 0 {
		return errors.New("glRasterPos2iv")
	}
	gpRasterPos2s = uintptr(getProcAddr("glRasterPos2s"))
	if gpRasterPos2s == 0 {
		return errors.New("glRasterPos2s")
	}
	gpRasterPos2sv = uintptr(getProcAddr("glRasterPos2sv"))
	if gpRasterPos2sv == 0 {
		return errors.New("glRasterPos2sv")
	}
	gpRasterPos2xOES = uintptr(getProcAddr("glRasterPos2xOES"))
	gpRasterPos2xvOES = uintptr(getProcAddr("glRasterPos2xvOES"))
	gpRasterPos3d = uintptr(getProcAddr("glRasterPos3d"))
	if gpRasterPos3d == 0 {
		return errors.New("glRasterPos3d")
	}
	gpRasterPos3dv = uintptr(getProcAddr("glRasterPos3dv"))
	if gpRasterPos3dv == 0 {
		return errors.New("glRasterPos3dv")
	}
	gpRasterPos3f = uintptr(getProcAddr("glRasterPos3f"))
	if gpRasterPos3f == 0 {
		return errors.New("glRasterPos3f")
	}
	gpRasterPos3fv = uintptr(getProcAddr("glRasterPos3fv"))
	if gpRasterPos3fv == 0 {
		return errors.New("glRasterPos3fv")
	}
	gpRasterPos3i = uintptr(getProcAddr("glRasterPos3i"))
	if gpRasterPos3i == 0 {
		return errors.New("glRasterPos3i")
	}
	gpRasterPos3iv = uintptr(getProcAddr("glRasterPos3iv"))
	if gpRasterPos3iv == 0 {
		return errors.New("glRasterPos3iv")
	}
	gpRasterPos3s = uintptr(getProcAddr("glRasterPos3s"))
	if gpRasterPos3s == 0 {
		return errors.New("glRasterPos3s")
	}
	gpRasterPos3sv = uintptr(getProcAddr("glRasterPos3sv"))
	if gpRasterPos3sv == 0 {
		return errors.New("glRasterPos3sv")
	}
	gpRasterPos3xOES = uintptr(getProcAddr("glRasterPos3xOES"))
	gpRasterPos3xvOES = uintptr(getProcAddr("glRasterPos3xvOES"))
	gpRasterPos4d = uintptr(getProcAddr("glRasterPos4d"))
	if gpRasterPos4d == 0 {
		return errors.New("glRasterPos4d")
	}
	gpRasterPos4dv = uintptr(getProcAddr("glRasterPos4dv"))
	if gpRasterPos4dv == 0 {
		return errors.New("glRasterPos4dv")
	}
	gpRasterPos4f = uintptr(getProcAddr("glRasterPos4f"))
	if gpRasterPos4f == 0 {
		return errors.New("glRasterPos4f")
	}
	gpRasterPos4fv = uintptr(getProcAddr("glRasterPos4fv"))
	if gpRasterPos4fv == 0 {
		return errors.New("glRasterPos4fv")
	}
	gpRasterPos4i = uintptr(getProcAddr("glRasterPos4i"))
	if gpRasterPos4i == 0 {
		return errors.New("glRasterPos4i")
	}
	gpRasterPos4iv = uintptr(getProcAddr("glRasterPos4iv"))
	if gpRasterPos4iv == 0 {
		return errors.New("glRasterPos4iv")
	}
	gpRasterPos4s = uintptr(getProcAddr("glRasterPos4s"))
	if gpRasterPos4s == 0 {
		return errors.New("glRasterPos4s")
	}
	gpRasterPos4sv = uintptr(getProcAddr("glRasterPos4sv"))
	if gpRasterPos4sv == 0 {
		return errors.New("glRasterPos4sv")
	}
	gpRasterPos4xOES = uintptr(getProcAddr("glRasterPos4xOES"))
	gpRasterPos4xvOES = uintptr(getProcAddr("glRasterPos4xvOES"))
	gpRasterSamplesEXT = uintptr(getProcAddr("glRasterSamplesEXT"))
	gpReadBuffer = uintptr(getProcAddr("glReadBuffer"))
	if gpReadBuffer == 0 {
		return errors.New("glReadBuffer")
	}
	gpReadInstrumentsSGIX = uintptr(getProcAddr("glReadInstrumentsSGIX"))
	gpReadPixels = uintptr(getProcAddr("glReadPixels"))
	if gpReadPixels == 0 {
		return errors.New("glReadPixels")
	}
	gpReadnPixels = uintptr(getProcAddr("glReadnPixels"))
	gpReadnPixelsARB = uintptr(getProcAddr("glReadnPixelsARB"))
	gpReadnPixelsKHR = uintptr(getProcAddr("glReadnPixelsKHR"))
	gpRectd = uintptr(getProcAddr("glRectd"))
	if gpRectd == 0 {
		return errors.New("glRectd")
	}
	gpRectdv = uintptr(getProcAddr("glRectdv"))
	if gpRectdv == 0 {
		return errors.New("glRectdv")
	}
	gpRectf = uintptr(getProcAddr("glRectf"))
	if gpRectf == 0 {
		return errors.New("glRectf")
	}
	gpRectfv = uintptr(getProcAddr("glRectfv"))
	if gpRectfv == 0 {
		return errors.New("glRectfv")
	}
	gpRecti = uintptr(getProcAddr("glRecti"))
	if gpRecti == 0 {
		return errors.New("glRecti")
	}
	gpRectiv = uintptr(getProcAddr("glRectiv"))
	if gpRectiv == 0 {
		return errors.New("glRectiv")
	}
	gpRects = uintptr(getProcAddr("glRects"))
	if gpRects == 0 {
		return errors.New("glRects")
	}
	gpRectsv = uintptr(getProcAddr("glRectsv"))
	if gpRectsv == 0 {
		return errors.New("glRectsv")
	}
	gpRectxOES = uintptr(getProcAddr("glRectxOES"))
	gpRectxvOES = uintptr(getProcAddr("glRectxvOES"))
	gpReferencePlaneSGIX = uintptr(getProcAddr("glReferencePlaneSGIX"))
	gpReleaseKeyedMutexWin32EXT = uintptr(getProcAddr("glReleaseKeyedMutexWin32EXT"))
	gpReleaseShaderCompiler = uintptr(getProcAddr("glReleaseShaderCompiler"))
	gpRenderGpuMaskNV = uintptr(getProcAddr("glRenderGpuMaskNV"))
	gpRenderMode = uintptr(getProcAddr("glRenderMode"))
	if gpRenderMode == 0 {
		return errors.New("glRenderMode")
	}
	gpRenderbufferStorage = uintptr(getProcAddr("glRenderbufferStorage"))
	gpRenderbufferStorageEXT = uintptr(getProcAddr("glRenderbufferStorageEXT"))
	gpRenderbufferStorageMultisample = uintptr(getProcAddr("glRenderbufferStorageMultisample"))
	gpRenderbufferStorageMultisampleCoverageNV = uintptr(getProcAddr("glRenderbufferStorageMultisampleCoverageNV"))
	gpRenderbufferStorageMultisampleEXT = uintptr(getProcAddr("glRenderbufferStorageMultisampleEXT"))
	gpReplacementCodePointerSUN = uintptr(getProcAddr("glReplacementCodePointerSUN"))
	gpReplacementCodeubSUN = uintptr(getProcAddr("glReplacementCodeubSUN"))
	gpReplacementCodeubvSUN = uintptr(getProcAddr("glReplacementCodeubvSUN"))
	gpReplacementCodeuiColor3fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiColor3fVertex3fSUN"))
	gpReplacementCodeuiColor3fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiColor3fVertex3fvSUN"))
	gpReplacementCodeuiColor4fNormal3fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiColor4fNormal3fVertex3fSUN"))
	gpReplacementCodeuiColor4fNormal3fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiColor4fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiColor4ubVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiColor4ubVertex3fSUN"))
	gpReplacementCodeuiColor4ubVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiColor4ubVertex3fvSUN"))
	gpReplacementCodeuiNormal3fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiNormal3fVertex3fSUN"))
	gpReplacementCodeuiNormal3fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiNormal3fVertex3fvSUN"))
	gpReplacementCodeuiSUN = uintptr(getProcAddr("glReplacementCodeuiSUN"))
	gpReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fColor4fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiTexCoord2fNormal3fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fNormal3fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fNormal3fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fNormal3fVertex3fvSUN"))
	gpReplacementCodeuiTexCoord2fVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fVertex3fSUN"))
	gpReplacementCodeuiTexCoord2fVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiTexCoord2fVertex3fvSUN"))
	gpReplacementCodeuiVertex3fSUN = uintptr(getProcAddr("glReplacementCodeuiVertex3fSUN"))
	gpReplacementCodeuiVertex3fvSUN = uintptr(getProcAddr("glReplacementCodeuiVertex3fvSUN"))
	gpReplacementCodeuivSUN = uintptr(getProcAddr("glReplacementCodeuivSUN"))
	gpReplacementCodeusSUN = uintptr(getProcAddr("glReplacementCodeusSUN"))
	gpReplacementCodeusvSUN = uintptr(getProcAddr("glReplacementCodeusvSUN"))
	gpRequestResidentProgramsNV = uintptr(getProcAddr("glRequestResidentProgramsNV"))
	gpResetHistogramEXT = uintptr(getProcAddr("glResetHistogramEXT"))
	gpResetMinmaxEXT = uintptr(getProcAddr("glResetMinmaxEXT"))
	gpResizeBuffersMESA = uintptr(getProcAddr("glResizeBuffersMESA"))
	gpResolveDepthValuesNV = uintptr(getProcAddr("glResolveDepthValuesNV"))
	gpResumeTransformFeedback = uintptr(getProcAddr("glResumeTransformFeedback"))
	gpResumeTransformFeedbackNV = uintptr(getProcAddr("glResumeTransformFeedbackNV"))
	gpRotated = uintptr(getProcAddr("glRotated"))
	if gpRotated == 0 {
		return errors.New("glRotated")
	}
	gpRotatef = uintptr(getProcAddr("glRotatef"))
	if gpRotatef == 0 {
		return errors.New("glRotatef")
	}
	gpRotatexOES = uintptr(getProcAddr("glRotatexOES"))
	gpSampleCoverage = uintptr(getProcAddr("glSampleCoverage"))
	if gpSampleCoverage == 0 {
		return errors.New("glSampleCoverage")
	}
	gpSampleCoverageARB = uintptr(getProcAddr("glSampleCoverageARB"))
	gpSampleCoveragexOES = uintptr(getProcAddr("glSampleCoveragexOES"))
	gpSampleMapATI = uintptr(getProcAddr("glSampleMapATI"))
	gpSampleMaskEXT = uintptr(getProcAddr("glSampleMaskEXT"))
	gpSampleMaskIndexedNV = uintptr(getProcAddr("glSampleMaskIndexedNV"))
	gpSampleMaskSGIS = uintptr(getProcAddr("glSampleMaskSGIS"))
	gpSampleMaski = uintptr(getProcAddr("glSampleMaski"))
	gpSamplePatternEXT = uintptr(getProcAddr("glSamplePatternEXT"))
	gpSamplePatternSGIS = uintptr(getProcAddr("glSamplePatternSGIS"))
	gpSamplerParameterIiv = uintptr(getProcAddr("glSamplerParameterIiv"))
	gpSamplerParameterIuiv = uintptr(getProcAddr("glSamplerParameterIuiv"))
	gpSamplerParameterf = uintptr(getProcAddr("glSamplerParameterf"))
	gpSamplerParameterfv = uintptr(getProcAddr("glSamplerParameterfv"))
	gpSamplerParameteri = uintptr(getProcAddr("glSamplerParameteri"))
	gpSamplerParameteriv = uintptr(getProcAddr("glSamplerParameteriv"))
	gpScaled = uintptr(getProcAddr("glScaled"))
	if gpScaled == 0 {
		return errors.New("glScaled")
	}
	gpScalef = uintptr(getProcAddr("glScalef"))
	if gpScalef == 0 {
		return errors.New("glScalef")
	}
	gpScalexOES = uintptr(getProcAddr("glScalexOES"))
	gpScissor = uintptr(getProcAddr("glScissor"))
	if gpScissor == 0 {
		return errors.New("glScissor")
	}
	gpScissorArrayv = uintptr(getProcAddr("glScissorArrayv"))
	gpScissorIndexed = uintptr(getProcAddr("glScissorIndexed"))
	gpScissorIndexedv = uintptr(getProcAddr("glScissorIndexedv"))
	gpSecondaryColor3b = uintptr(getProcAddr("glSecondaryColor3b"))
	if gpSecondaryColor3b == 0 {
		return errors.New("glSecondaryColor3b")
	}
	gpSecondaryColor3bEXT = uintptr(getProcAddr("glSecondaryColor3bEXT"))
	gpSecondaryColor3bv = uintptr(getProcAddr("glSecondaryColor3bv"))
	if gpSecondaryColor3bv == 0 {
		return errors.New("glSecondaryColor3bv")
	}
	gpSecondaryColor3bvEXT = uintptr(getProcAddr("glSecondaryColor3bvEXT"))
	gpSecondaryColor3d = uintptr(getProcAddr("glSecondaryColor3d"))
	if gpSecondaryColor3d == 0 {
		return errors.New("glSecondaryColor3d")
	}
	gpSecondaryColor3dEXT = uintptr(getProcAddr("glSecondaryColor3dEXT"))
	gpSecondaryColor3dv = uintptr(getProcAddr("glSecondaryColor3dv"))
	if gpSecondaryColor3dv == 0 {
		return errors.New("glSecondaryColor3dv")
	}
	gpSecondaryColor3dvEXT = uintptr(getProcAddr("glSecondaryColor3dvEXT"))
	gpSecondaryColor3f = uintptr(getProcAddr("glSecondaryColor3f"))
	if gpSecondaryColor3f == 0 {
		return errors.New("glSecondaryColor3f")
	}
	gpSecondaryColor3fEXT = uintptr(getProcAddr("glSecondaryColor3fEXT"))
	gpSecondaryColor3fv = uintptr(getProcAddr("glSecondaryColor3fv"))
	if gpSecondaryColor3fv == 0 {
		return errors.New("glSecondaryColor3fv")
	}
	gpSecondaryColor3fvEXT = uintptr(getProcAddr("glSecondaryColor3fvEXT"))
	gpSecondaryColor3hNV = uintptr(getProcAddr("glSecondaryColor3hNV"))
	gpSecondaryColor3hvNV = uintptr(getProcAddr("glSecondaryColor3hvNV"))
	gpSecondaryColor3i = uintptr(getProcAddr("glSecondaryColor3i"))
	if gpSecondaryColor3i == 0 {
		return errors.New("glSecondaryColor3i")
	}
	gpSecondaryColor3iEXT = uintptr(getProcAddr("glSecondaryColor3iEXT"))
	gpSecondaryColor3iv = uintptr(getProcAddr("glSecondaryColor3iv"))
	if gpSecondaryColor3iv == 0 {
		return errors.New("glSecondaryColor3iv")
	}
	gpSecondaryColor3ivEXT = uintptr(getProcAddr("glSecondaryColor3ivEXT"))
	gpSecondaryColor3s = uintptr(getProcAddr("glSecondaryColor3s"))
	if gpSecondaryColor3s == 0 {
		return errors.New("glSecondaryColor3s")
	}
	gpSecondaryColor3sEXT = uintptr(getProcAddr("glSecondaryColor3sEXT"))
	gpSecondaryColor3sv = uintptr(getProcAddr("glSecondaryColor3sv"))
	if gpSecondaryColor3sv == 0 {
		return errors.New("glSecondaryColor3sv")
	}
	gpSecondaryColor3svEXT = uintptr(getProcAddr("glSecondaryColor3svEXT"))
	gpSecondaryColor3ub = uintptr(getProcAddr("glSecondaryColor3ub"))
	if gpSecondaryColor3ub == 0 {
		return errors.New("glSecondaryColor3ub")
	}
	gpSecondaryColor3ubEXT = uintptr(getProcAddr("glSecondaryColor3ubEXT"))
	gpSecondaryColor3ubv = uintptr(getProcAddr("glSecondaryColor3ubv"))
	if gpSecondaryColor3ubv == 0 {
		return errors.New("glSecondaryColor3ubv")
	}
	gpSecondaryColor3ubvEXT = uintptr(getProcAddr("glSecondaryColor3ubvEXT"))
	gpSecondaryColor3ui = uintptr(getProcAddr("glSecondaryColor3ui"))
	if gpSecondaryColor3ui == 0 {
		return errors.New("glSecondaryColor3ui")
	}
	gpSecondaryColor3uiEXT = uintptr(getProcAddr("glSecondaryColor3uiEXT"))
	gpSecondaryColor3uiv = uintptr(getProcAddr("glSecondaryColor3uiv"))
	if gpSecondaryColor3uiv == 0 {
		return errors.New("glSecondaryColor3uiv")
	}
	gpSecondaryColor3uivEXT = uintptr(getProcAddr("glSecondaryColor3uivEXT"))
	gpSecondaryColor3us = uintptr(getProcAddr("glSecondaryColor3us"))
	if gpSecondaryColor3us == 0 {
		return errors.New("glSecondaryColor3us")
	}
	gpSecondaryColor3usEXT = uintptr(getProcAddr("glSecondaryColor3usEXT"))
	gpSecondaryColor3usv = uintptr(getProcAddr("glSecondaryColor3usv"))
	if gpSecondaryColor3usv == 0 {
		return errors.New("glSecondaryColor3usv")
	}
	gpSecondaryColor3usvEXT = uintptr(getProcAddr("glSecondaryColor3usvEXT"))
	gpSecondaryColorFormatNV = uintptr(getProcAddr("glSecondaryColorFormatNV"))
	gpSecondaryColorPointer = uintptr(getProcAddr("glSecondaryColorPointer"))
	if gpSecondaryColorPointer == 0 {
		return errors.New("glSecondaryColorPointer")
	}
	gpSecondaryColorPointerEXT = uintptr(getProcAddr("glSecondaryColorPointerEXT"))
	gpSecondaryColorPointerListIBM = uintptr(getProcAddr("glSecondaryColorPointerListIBM"))
	gpSelectBuffer = uintptr(getProcAddr("glSelectBuffer"))
	if gpSelectBuffer == 0 {
		return errors.New("glSelectBuffer")
	}
	gpSelectPerfMonitorCountersAMD = uintptr(getProcAddr("glSelectPerfMonitorCountersAMD"))
	gpSemaphoreParameterui64vEXT = uintptr(getProcAddr("glSemaphoreParameterui64vEXT"))
	gpSeparableFilter2DEXT = uintptr(getProcAddr("glSeparableFilter2DEXT"))
	gpSetFenceAPPLE = uintptr(getProcAddr("glSetFenceAPPLE"))
	gpSetFenceNV = uintptr(getProcAddr("glSetFenceNV"))
	gpSetFragmentShaderConstantATI = uintptr(getProcAddr("glSetFragmentShaderConstantATI"))
	gpSetInvariantEXT = uintptr(getProcAddr("glSetInvariantEXT"))
	gpSetLocalConstantEXT = uintptr(getProcAddr("glSetLocalConstantEXT"))
	gpSetMultisamplefvAMD = uintptr(getProcAddr("glSetMultisamplefvAMD"))
	gpShadeModel = uintptr(getProcAddr("glShadeModel"))
	if gpShadeModel == 0 {
		return errors.New("glShadeModel")
	}
	gpShaderBinary = uintptr(getProcAddr("glShaderBinary"))
	gpShaderOp1EXT = uintptr(getProcAddr("glShaderOp1EXT"))
	gpShaderOp2EXT = uintptr(getProcAddr("glShaderOp2EXT"))
	gpShaderOp3EXT = uintptr(getProcAddr("glShaderOp3EXT"))
	gpShaderSource = uintptr(getProcAddr("glShaderSource"))
	if gpShaderSource == 0 {
		return errors.New("glShaderSource")
	}
	gpShaderSourceARB = uintptr(getProcAddr("glShaderSourceARB"))
	gpShaderStorageBlockBinding = uintptr(getProcAddr("glShaderStorageBlockBinding"))
	gpSharpenTexFuncSGIS = uintptr(getProcAddr("glSharpenTexFuncSGIS"))
	gpSignalSemaphoreEXT = uintptr(getProcAddr("glSignalSemaphoreEXT"))
	gpSignalVkFenceNV = uintptr(getProcAddr("glSignalVkFenceNV"))
	gpSignalVkSemaphoreNV = uintptr(getProcAddr("glSignalVkSemaphoreNV"))
	gpSpecializeShaderARB = uintptr(getProcAddr("glSpecializeShaderARB"))
	gpSpriteParameterfSGIX = uintptr(getProcAddr("glSpriteParameterfSGIX"))
	gpSpriteParameterfvSGIX = uintptr(getProcAddr("glSpriteParameterfvSGIX"))
	gpSpriteParameteriSGIX = uintptr(getProcAddr("glSpriteParameteriSGIX"))
	gpSpriteParameterivSGIX = uintptr(getProcAddr("glSpriteParameterivSGIX"))
	gpStartInstrumentsSGIX = uintptr(getProcAddr("glStartInstrumentsSGIX"))
	gpStateCaptureNV = uintptr(getProcAddr("glStateCaptureNV"))
	gpStencilClearTagEXT = uintptr(getProcAddr("glStencilClearTagEXT"))
	gpStencilFillPathInstancedNV = uintptr(getProcAddr("glStencilFillPathInstancedNV"))
	gpStencilFillPathNV = uintptr(getProcAddr("glStencilFillPathNV"))
	gpStencilFunc = uintptr(getProcAddr("glStencilFunc"))
	if gpStencilFunc == 0 {
		return errors.New("glStencilFunc")
	}
	gpStencilFuncSeparate = uintptr(getProcAddr("glStencilFuncSeparate"))
	if gpStencilFuncSeparate == 0 {
		return errors.New("glStencilFuncSeparate")
	}
	gpStencilFuncSeparateATI = uintptr(getProcAddr("glStencilFuncSeparateATI"))
	gpStencilMask = uintptr(getProcAddr("glStencilMask"))
	if gpStencilMask == 0 {
		return errors.New("glStencilMask")
	}
	gpStencilMaskSeparate = uintptr(getProcAddr("glStencilMaskSeparate"))
	if gpStencilMaskSeparate == 0 {
		return errors.New("glStencilMaskSeparate")
	}
	gpStencilOp = uintptr(getProcAddr("glStencilOp"))
	if gpStencilOp == 0 {
		return errors.New("glStencilOp")
	}
	gpStencilOpSeparate = uintptr(getProcAddr("glStencilOpSeparate"))
	if gpStencilOpSeparate == 0 {
		return errors.New("glStencilOpSeparate")
	}
	gpStencilOpSeparateATI = uintptr(getProcAddr("glStencilOpSeparateATI"))
	gpStencilOpValueAMD = uintptr(getProcAddr("glStencilOpValueAMD"))
	gpStencilStrokePathInstancedNV = uintptr(getProcAddr("glStencilStrokePathInstancedNV"))
	gpStencilStrokePathNV = uintptr(getProcAddr("glStencilStrokePathNV"))
	gpStencilThenCoverFillPathInstancedNV = uintptr(getProcAddr("glStencilThenCoverFillPathInstancedNV"))
	gpStencilThenCoverFillPathNV = uintptr(getProcAddr("glStencilThenCoverFillPathNV"))
	gpStencilThenCoverStrokePathInstancedNV = uintptr(getProcAddr("glStencilThenCoverStrokePathInstancedNV"))
	gpStencilThenCoverStrokePathNV = uintptr(getProcAddr("glStencilThenCoverStrokePathNV"))
	gpStopInstrumentsSGIX = uintptr(getProcAddr("glStopInstrumentsSGIX"))
	gpStringMarkerGREMEDY = uintptr(getProcAddr("glStringMarkerGREMEDY"))
	gpSubpixelPrecisionBiasNV = uintptr(getProcAddr("glSubpixelPrecisionBiasNV"))
	gpSwizzleEXT = uintptr(getProcAddr("glSwizzleEXT"))
	gpSyncTextureINTEL = uintptr(getProcAddr("glSyncTextureINTEL"))
	gpTagSampleBufferSGIX = uintptr(getProcAddr("glTagSampleBufferSGIX"))
	gpTangent3bEXT = uintptr(getProcAddr("glTangent3bEXT"))
	gpTangent3bvEXT = uintptr(getProcAddr("glTangent3bvEXT"))
	gpTangent3dEXT = uintptr(getProcAddr("glTangent3dEXT"))
	gpTangent3dvEXT = uintptr(getProcAddr("glTangent3dvEXT"))
	gpTangent3fEXT = uintptr(getProcAddr("glTangent3fEXT"))
	gpTangent3fvEXT = uintptr(getProcAddr("glTangent3fvEXT"))
	gpTangent3iEXT = uintptr(getProcAddr("glTangent3iEXT"))
	gpTangent3ivEXT = uintptr(getProcAddr("glTangent3ivEXT"))
	gpTangent3sEXT = uintptr(getProcAddr("glTangent3sEXT"))
	gpTangent3svEXT = uintptr(getProcAddr("glTangent3svEXT"))
	gpTangentPointerEXT = uintptr(getProcAddr("glTangentPointerEXT"))
	gpTbufferMask3DFX = uintptr(getProcAddr("glTbufferMask3DFX"))
	gpTessellationFactorAMD = uintptr(getProcAddr("glTessellationFactorAMD"))
	gpTessellationModeAMD = uintptr(getProcAddr("glTessellationModeAMD"))
	gpTestFenceAPPLE = uintptr(getProcAddr("glTestFenceAPPLE"))
	gpTestFenceNV = uintptr(getProcAddr("glTestFenceNV"))
	gpTestObjectAPPLE = uintptr(getProcAddr("glTestObjectAPPLE"))
	gpTexBufferARB = uintptr(getProcAddr("glTexBufferARB"))
	gpTexBufferEXT = uintptr(getProcAddr("glTexBufferEXT"))
	gpTexBufferRange = uintptr(getProcAddr("glTexBufferRange"))
	gpTexBumpParameterfvATI = uintptr(getProcAddr("glTexBumpParameterfvATI"))
	gpTexBumpParameterivATI = uintptr(getProcAddr("glTexBumpParameterivATI"))
	gpTexCoord1bOES = uintptr(getProcAddr("glTexCoord1bOES"))
	gpTexCoord1bvOES = uintptr(getProcAddr("glTexCoord1bvOES"))
	gpTexCoord1d = uintptr(getProcAddr("glTexCoord1d"))
	if gpTexCoord1d == 0 {
		return errors.New("glTexCoord1d")
	}
	gpTexCoord1dv = uintptr(getProcAddr("glTexCoord1dv"))
	if gpTexCoord1dv == 0 {
		return errors.New("glTexCoord1dv")
	}
	gpTexCoord1f = uintptr(getProcAddr("glTexCoord1f"))
	if gpTexCoord1f == 0 {
		return errors.New("glTexCoord1f")
	}
	gpTexCoord1fv = uintptr(getProcAddr("glTexCoord1fv"))
	if gpTexCoord1fv == 0 {
		return errors.New("glTexCoord1fv")
	}
	gpTexCoord1hNV = uintptr(getProcAddr("glTexCoord1hNV"))
	gpTexCoord1hvNV = uintptr(getProcAddr("glTexCoord1hvNV"))
	gpTexCoord1i = uintptr(getProcAddr("glTexCoord1i"))
	if gpTexCoord1i == 0 {
		return errors.New("glTexCoord1i")
	}
	gpTexCoord1iv = uintptr(getProcAddr("glTexCoord1iv"))
	if gpTexCoord1iv == 0 {
		return errors.New("glTexCoord1iv")
	}
	gpTexCoord1s = uintptr(getProcAddr("glTexCoord1s"))
	if gpTexCoord1s == 0 {
		return errors.New("glTexCoord1s")
	}
	gpTexCoord1sv = uintptr(getProcAddr("glTexCoord1sv"))
	if gpTexCoord1sv == 0 {
		return errors.New("glTexCoord1sv")
	}
	gpTexCoord1xOES = uintptr(getProcAddr("glTexCoord1xOES"))
	gpTexCoord1xvOES = uintptr(getProcAddr("glTexCoord1xvOES"))
	gpTexCoord2bOES = uintptr(getProcAddr("glTexCoord2bOES"))
	gpTexCoord2bvOES = uintptr(getProcAddr("glTexCoord2bvOES"))
	gpTexCoord2d = uintptr(getProcAddr("glTexCoord2d"))
	if gpTexCoord2d == 0 {
		return errors.New("glTexCoord2d")
	}
	gpTexCoord2dv = uintptr(getProcAddr("glTexCoord2dv"))
	if gpTexCoord2dv == 0 {
		return errors.New("glTexCoord2dv")
	}
	gpTexCoord2f = uintptr(getProcAddr("glTexCoord2f"))
	if gpTexCoord2f == 0 {
		return errors.New("glTexCoord2f")
	}
	gpTexCoord2fColor3fVertex3fSUN = uintptr(getProcAddr("glTexCoord2fColor3fVertex3fSUN"))
	gpTexCoord2fColor3fVertex3fvSUN = uintptr(getProcAddr("glTexCoord2fColor3fVertex3fvSUN"))
	gpTexCoord2fColor4fNormal3fVertex3fSUN = uintptr(getProcAddr("glTexCoord2fColor4fNormal3fVertex3fSUN"))
	gpTexCoord2fColor4fNormal3fVertex3fvSUN = uintptr(getProcAddr("glTexCoord2fColor4fNormal3fVertex3fvSUN"))
	gpTexCoord2fColor4ubVertex3fSUN = uintptr(getProcAddr("glTexCoord2fColor4ubVertex3fSUN"))
	gpTexCoord2fColor4ubVertex3fvSUN = uintptr(getProcAddr("glTexCoord2fColor4ubVertex3fvSUN"))
	gpTexCoord2fNormal3fVertex3fSUN = uintptr(getProcAddr("glTexCoord2fNormal3fVertex3fSUN"))
	gpTexCoord2fNormal3fVertex3fvSUN = uintptr(getProcAddr("glTexCoord2fNormal3fVertex3fvSUN"))
	gpTexCoord2fVertex3fSUN = uintptr(getProcAddr("glTexCoord2fVertex3fSUN"))
	gpTexCoord2fVertex3fvSUN = uintptr(getProcAddr("glTexCoord2fVertex3fvSUN"))
	gpTexCoord2fv = uintptr(getProcAddr("glTexCoord2fv"))
	if gpTexCoord2fv == 0 {
		return errors.New("glTexCoord2fv")
	}
	gpTexCoord2hNV = uintptr(getProcAddr("glTexCoord2hNV"))
	gpTexCoord2hvNV = uintptr(getProcAddr("glTexCoord2hvNV"))
	gpTexCoord2i = uintptr(getProcAddr("glTexCoord2i"))
	if gpTexCoord2i == 0 {
		return errors.New("glTexCoord2i")
	}
	gpTexCoord2iv = uintptr(getProcAddr("glTexCoord2iv"))
	if gpTexCoord2iv == 0 {
		return errors.New("glTexCoord2iv")
	}
	gpTexCoord2s = uintptr(getProcAddr("glTexCoord2s"))
	if gpTexCoord2s == 0 {
		return errors.New("glTexCoord2s")
	}
	gpTexCoord2sv = uintptr(getProcAddr("glTexCoord2sv"))
	if gpTexCoord2sv == 0 {
		return errors.New("glTexCoord2sv")
	}
	gpTexCoord2xOES = uintptr(getProcAddr("glTexCoord2xOES"))
	gpTexCoord2xvOES = uintptr(getProcAddr("glTexCoord2xvOES"))
	gpTexCoord3bOES = uintptr(getProcAddr("glTexCoord3bOES"))
	gpTexCoord3bvOES = uintptr(getProcAddr("glTexCoord3bvOES"))
	gpTexCoord3d = uintptr(getProcAddr("glTexCoord3d"))
	if gpTexCoord3d == 0 {
		return errors.New("glTexCoord3d")
	}
	gpTexCoord3dv = uintptr(getProcAddr("glTexCoord3dv"))
	if gpTexCoord3dv == 0 {
		return errors.New("glTexCoord3dv")
	}
	gpTexCoord3f = uintptr(getProcAddr("glTexCoord3f"))
	if gpTexCoord3f == 0 {
		return errors.New("glTexCoord3f")
	}
	gpTexCoord3fv = uintptr(getProcAddr("glTexCoord3fv"))
	if gpTexCoord3fv == 0 {
		return errors.New("glTexCoord3fv")
	}
	gpTexCoord3hNV = uintptr(getProcAddr("glTexCoord3hNV"))
	gpTexCoord3hvNV = uintptr(getProcAddr("glTexCoord3hvNV"))
	gpTexCoord3i = uintptr(getProcAddr("glTexCoord3i"))
	if gpTexCoord3i == 0 {
		return errors.New("glTexCoord3i")
	}
	gpTexCoord3iv = uintptr(getProcAddr("glTexCoord3iv"))
	if gpTexCoord3iv == 0 {
		return errors.New("glTexCoord3iv")
	}
	gpTexCoord3s = uintptr(getProcAddr("glTexCoord3s"))
	if gpTexCoord3s == 0 {
		return errors.New("glTexCoord3s")
	}
	gpTexCoord3sv = uintptr(getProcAddr("glTexCoord3sv"))
	if gpTexCoord3sv == 0 {
		return errors.New("glTexCoord3sv")
	}
	gpTexCoord3xOES = uintptr(getProcAddr("glTexCoord3xOES"))
	gpTexCoord3xvOES = uintptr(getProcAddr("glTexCoord3xvOES"))
	gpTexCoord4bOES = uintptr(getProcAddr("glTexCoord4bOES"))
	gpTexCoord4bvOES = uintptr(getProcAddr("glTexCoord4bvOES"))
	gpTexCoord4d = uintptr(getProcAddr("glTexCoord4d"))
	if gpTexCoord4d == 0 {
		return errors.New("glTexCoord4d")
	}
	gpTexCoord4dv = uintptr(getProcAddr("glTexCoord4dv"))
	if gpTexCoord4dv == 0 {
		return errors.New("glTexCoord4dv")
	}
	gpTexCoord4f = uintptr(getProcAddr("glTexCoord4f"))
	if gpTexCoord4f == 0 {
		return errors.New("glTexCoord4f")
	}
	gpTexCoord4fColor4fNormal3fVertex4fSUN = uintptr(getProcAddr("glTexCoord4fColor4fNormal3fVertex4fSUN"))
	gpTexCoord4fColor4fNormal3fVertex4fvSUN = uintptr(getProcAddr("glTexCoord4fColor4fNormal3fVertex4fvSUN"))
	gpTexCoord4fVertex4fSUN = uintptr(getProcAddr("glTexCoord4fVertex4fSUN"))
	gpTexCoord4fVertex4fvSUN = uintptr(getProcAddr("glTexCoord4fVertex4fvSUN"))
	gpTexCoord4fv = uintptr(getProcAddr("glTexCoord4fv"))
	if gpTexCoord4fv == 0 {
		return errors.New("glTexCoord4fv")
	}
	gpTexCoord4hNV = uintptr(getProcAddr("glTexCoord4hNV"))
	gpTexCoord4hvNV = uintptr(getProcAddr("glTexCoord4hvNV"))
	gpTexCoord4i = uintptr(getProcAddr("glTexCoord4i"))
	if gpTexCoord4i == 0 {
		return errors.New("glTexCoord4i")
	}
	gpTexCoord4iv = uintptr(getProcAddr("glTexCoord4iv"))
	if gpTexCoord4iv == 0 {
		return errors.New("glTexCoord4iv")
	}
	gpTexCoord4s = uintptr(getProcAddr("glTexCoord4s"))
	if gpTexCoord4s == 0 {
		return errors.New("glTexCoord4s")
	}
	gpTexCoord4sv = uintptr(getProcAddr("glTexCoord4sv"))
	if gpTexCoord4sv == 0 {
		return errors.New("glTexCoord4sv")
	}
	gpTexCoord4xOES = uintptr(getProcAddr("glTexCoord4xOES"))
	gpTexCoord4xvOES = uintptr(getProcAddr("glTexCoord4xvOES"))
	gpTexCoordFormatNV = uintptr(getProcAddr("glTexCoordFormatNV"))
	gpTexCoordPointer = uintptr(getProcAddr("glTexCoordPointer"))
	if gpTexCoordPointer == 0 {
		return errors.New("glTexCoordPointer")
	}
	gpTexCoordPointerEXT = uintptr(getProcAddr("glTexCoordPointerEXT"))
	gpTexCoordPointerListIBM = uintptr(getProcAddr("glTexCoordPointerListIBM"))
	gpTexCoordPointervINTEL = uintptr(getProcAddr("glTexCoordPointervINTEL"))
	gpTexEnvf = uintptr(getProcAddr("glTexEnvf"))
	if gpTexEnvf == 0 {
		return errors.New("glTexEnvf")
	}
	gpTexEnvfv = uintptr(getProcAddr("glTexEnvfv"))
	if gpTexEnvfv == 0 {
		return errors.New("glTexEnvfv")
	}
	gpTexEnvi = uintptr(getProcAddr("glTexEnvi"))
	if gpTexEnvi == 0 {
		return errors.New("glTexEnvi")
	}
	gpTexEnviv = uintptr(getProcAddr("glTexEnviv"))
	if gpTexEnviv == 0 {
		return errors.New("glTexEnviv")
	}
	gpTexEnvxOES = uintptr(getProcAddr("glTexEnvxOES"))
	gpTexEnvxvOES = uintptr(getProcAddr("glTexEnvxvOES"))
	gpTexFilterFuncSGIS = uintptr(getProcAddr("glTexFilterFuncSGIS"))
	gpTexGend = uintptr(getProcAddr("glTexGend"))
	if gpTexGend == 0 {
		return errors.New("glTexGend")
	}
	gpTexGendv = uintptr(getProcAddr("glTexGendv"))
	if gpTexGendv == 0 {
		return errors.New("glTexGendv")
	}
	gpTexGenf = uintptr(getProcAddr("glTexGenf"))
	if gpTexGenf == 0 {
		return errors.New("glTexGenf")
	}
	gpTexGenfv = uintptr(getProcAddr("glTexGenfv"))
	if gpTexGenfv == 0 {
		return errors.New("glTexGenfv")
	}
	gpTexGeni = uintptr(getProcAddr("glTexGeni"))
	if gpTexGeni == 0 {
		return errors.New("glTexGeni")
	}
	gpTexGeniv = uintptr(getProcAddr("glTexGeniv"))
	if gpTexGeniv == 0 {
		return errors.New("glTexGeniv")
	}
	gpTexGenxOES = uintptr(getProcAddr("glTexGenxOES"))
	gpTexGenxvOES = uintptr(getProcAddr("glTexGenxvOES"))
	gpTexImage1D = uintptr(getProcAddr("glTexImage1D"))
	if gpTexImage1D == 0 {
		return errors.New("glTexImage1D")
	}
	gpTexImage2D = uintptr(getProcAddr("glTexImage2D"))
	if gpTexImage2D == 0 {
		return errors.New("glTexImage2D")
	}
	gpTexImage2DMultisample = uintptr(getProcAddr("glTexImage2DMultisample"))
	gpTexImage2DMultisampleCoverageNV = uintptr(getProcAddr("glTexImage2DMultisampleCoverageNV"))
	gpTexImage3D = uintptr(getProcAddr("glTexImage3D"))
	if gpTexImage3D == 0 {
		return errors.New("glTexImage3D")
	}
	gpTexImage3DEXT = uintptr(getProcAddr("glTexImage3DEXT"))
	gpTexImage3DMultisample = uintptr(getProcAddr("glTexImage3DMultisample"))
	gpTexImage3DMultisampleCoverageNV = uintptr(getProcAddr("glTexImage3DMultisampleCoverageNV"))
	gpTexImage4DSGIS = uintptr(getProcAddr("glTexImage4DSGIS"))
	gpTexPageCommitmentARB = uintptr(getProcAddr("glTexPageCommitmentARB"))
	gpTexParameterIivEXT = uintptr(getProcAddr("glTexParameterIivEXT"))
	gpTexParameterIuivEXT = uintptr(getProcAddr("glTexParameterIuivEXT"))
	gpTexParameterf = uintptr(getProcAddr("glTexParameterf"))
	if gpTexParameterf == 0 {
		return errors.New("glTexParameterf")
	}
	gpTexParameterfv = uintptr(getProcAddr("glTexParameterfv"))
	if gpTexParameterfv == 0 {
		return errors.New("glTexParameterfv")
	}
	gpTexParameteri = uintptr(getProcAddr("glTexParameteri"))
	if gpTexParameteri == 0 {
		return errors.New("glTexParameteri")
	}
	gpTexParameteriv = uintptr(getProcAddr("glTexParameteriv"))
	if gpTexParameteriv == 0 {
		return errors.New("glTexParameteriv")
	}
	gpTexParameterxOES = uintptr(getProcAddr("glTexParameterxOES"))
	gpTexParameterxvOES = uintptr(getProcAddr("glTexParameterxvOES"))
	gpTexRenderbufferNV = uintptr(getProcAddr("glTexRenderbufferNV"))
	gpTexStorage1D = uintptr(getProcAddr("glTexStorage1D"))
	gpTexStorage2D = uintptr(getProcAddr("glTexStorage2D"))
	gpTexStorage2DMultisample = uintptr(getProcAddr("glTexStorage2DMultisample"))
	gpTexStorage3D = uintptr(getProcAddr("glTexStorage3D"))
	gpTexStorage3DMultisample = uintptr(getProcAddr("glTexStorage3DMultisample"))
	gpTexStorageMem1DEXT = uintptr(getProcAddr("glTexStorageMem1DEXT"))
	gpTexStorageMem2DEXT = uintptr(getProcAddr("glTexStorageMem2DEXT"))
	gpTexStorageMem2DMultisampleEXT = uintptr(getProcAddr("glTexStorageMem2DMultisampleEXT"))
	gpTexStorageMem3DEXT = uintptr(getProcAddr("glTexStorageMem3DEXT"))
	gpTexStorageMem3DMultisampleEXT = uintptr(getProcAddr("glTexStorageMem3DMultisampleEXT"))
	gpTexStorageSparseAMD = uintptr(getProcAddr("glTexStorageSparseAMD"))
	gpTexSubImage1D = uintptr(getProcAddr("glTexSubImage1D"))
	if gpTexSubImage1D == 0 {
		return errors.New("glTexSubImage1D")
	}
	gpTexSubImage1DEXT = uintptr(getProcAddr("glTexSubImage1DEXT"))
	gpTexSubImage2D = uintptr(getProcAddr("glTexSubImage2D"))
	if gpTexSubImage2D == 0 {
		return errors.New("glTexSubImage2D")
	}
	gpTexSubImage2DEXT = uintptr(getProcAddr("glTexSubImage2DEXT"))
	gpTexSubImage3D = uintptr(getProcAddr("glTexSubImage3D"))
	if gpTexSubImage3D == 0 {
		return errors.New("glTexSubImage3D")
	}
	gpTexSubImage3DEXT = uintptr(getProcAddr("glTexSubImage3DEXT"))
	gpTexSubImage4DSGIS = uintptr(getProcAddr("glTexSubImage4DSGIS"))
	gpTextureBarrier = uintptr(getProcAddr("glTextureBarrier"))
	gpTextureBarrierNV = uintptr(getProcAddr("glTextureBarrierNV"))
	gpTextureBuffer = uintptr(getProcAddr("glTextureBuffer"))
	gpTextureBufferEXT = uintptr(getProcAddr("glTextureBufferEXT"))
	gpTextureBufferRange = uintptr(getProcAddr("glTextureBufferRange"))
	gpTextureBufferRangeEXT = uintptr(getProcAddr("glTextureBufferRangeEXT"))
	gpTextureColorMaskSGIS = uintptr(getProcAddr("glTextureColorMaskSGIS"))
	gpTextureImage1DEXT = uintptr(getProcAddr("glTextureImage1DEXT"))
	gpTextureImage2DEXT = uintptr(getProcAddr("glTextureImage2DEXT"))
	gpTextureImage2DMultisampleCoverageNV = uintptr(getProcAddr("glTextureImage2DMultisampleCoverageNV"))
	gpTextureImage2DMultisampleNV = uintptr(getProcAddr("glTextureImage2DMultisampleNV"))
	gpTextureImage3DEXT = uintptr(getProcAddr("glTextureImage3DEXT"))
	gpTextureImage3DMultisampleCoverageNV = uintptr(getProcAddr("glTextureImage3DMultisampleCoverageNV"))
	gpTextureImage3DMultisampleNV = uintptr(getProcAddr("glTextureImage3DMultisampleNV"))
	gpTextureLightEXT = uintptr(getProcAddr("glTextureLightEXT"))
	gpTextureMaterialEXT = uintptr(getProcAddr("glTextureMaterialEXT"))
	gpTextureNormalEXT = uintptr(getProcAddr("glTextureNormalEXT"))
	gpTexturePageCommitmentEXT = uintptr(getProcAddr("glTexturePageCommitmentEXT"))
	gpTextureParameterIiv = uintptr(getProcAddr("glTextureParameterIiv"))
	gpTextureParameterIivEXT = uintptr(getProcAddr("glTextureParameterIivEXT"))
	gpTextureParameterIuiv = uintptr(getProcAddr("glTextureParameterIuiv"))
	gpTextureParameterIuivEXT = uintptr(getProcAddr("glTextureParameterIuivEXT"))
	gpTextureParameterf = uintptr(getProcAddr("glTextureParameterf"))
	gpTextureParameterfEXT = uintptr(getProcAddr("glTextureParameterfEXT"))
	gpTextureParameterfv = uintptr(getProcAddr("glTextureParameterfv"))
	gpTextureParameterfvEXT = uintptr(getProcAddr("glTextureParameterfvEXT"))
	gpTextureParameteri = uintptr(getProcAddr("glTextureParameteri"))
	gpTextureParameteriEXT = uintptr(getProcAddr("glTextureParameteriEXT"))
	gpTextureParameteriv = uintptr(getProcAddr("glTextureParameteriv"))
	gpTextureParameterivEXT = uintptr(getProcAddr("glTextureParameterivEXT"))
	gpTextureRangeAPPLE = uintptr(getProcAddr("glTextureRangeAPPLE"))
	gpTextureRenderbufferEXT = uintptr(getProcAddr("glTextureRenderbufferEXT"))
	gpTextureStorage1D = uintptr(getProcAddr("glTextureStorage1D"))
	gpTextureStorage1DEXT = uintptr(getProcAddr("glTextureStorage1DEXT"))
	gpTextureStorage2D = uintptr(getProcAddr("glTextureStorage2D"))
	gpTextureStorage2DEXT = uintptr(getProcAddr("glTextureStorage2DEXT"))
	gpTextureStorage2DMultisample = uintptr(getProcAddr("glTextureStorage2DMultisample"))
	gpTextureStorage2DMultisampleEXT = uintptr(getProcAddr("glTextureStorage2DMultisampleEXT"))
	gpTextureStorage3D = uintptr(getProcAddr("glTextureStorage3D"))
	gpTextureStorage3DEXT = uintptr(getProcAddr("glTextureStorage3DEXT"))
	gpTextureStorage3DMultisample = uintptr(getProcAddr("glTextureStorage3DMultisample"))
	gpTextureStorage3DMultisampleEXT = uintptr(getProcAddr("glTextureStorage3DMultisampleEXT"))
	gpTextureStorageMem1DEXT = uintptr(getProcAddr("glTextureStorageMem1DEXT"))
	gpTextureStorageMem2DEXT = uintptr(getProcAddr("glTextureStorageMem2DEXT"))
	gpTextureStorageMem2DMultisampleEXT = uintptr(getProcAddr("glTextureStorageMem2DMultisampleEXT"))
	gpTextureStorageMem3DEXT = uintptr(getProcAddr("glTextureStorageMem3DEXT"))
	gpTextureStorageMem3DMultisampleEXT = uintptr(getProcAddr("glTextureStorageMem3DMultisampleEXT"))
	gpTextureStorageSparseAMD = uintptr(getProcAddr("glTextureStorageSparseAMD"))
	gpTextureSubImage1D = uintptr(getProcAddr("glTextureSubImage1D"))
	gpTextureSubImage1DEXT = uintptr(getProcAddr("glTextureSubImage1DEXT"))
	gpTextureSubImage2D = uintptr(getProcAddr("glTextureSubImage2D"))
	gpTextureSubImage2DEXT = uintptr(getProcAddr("glTextureSubImage2DEXT"))
	gpTextureSubImage3D = uintptr(getProcAddr("glTextureSubImage3D"))
	gpTextureSubImage3DEXT = uintptr(getProcAddr("glTextureSubImage3DEXT"))
	gpTextureView = uintptr(getProcAddr("glTextureView"))
	gpTrackMatrixNV = uintptr(getProcAddr("glTrackMatrixNV"))
	gpTransformFeedbackAttribsNV = uintptr(getProcAddr("glTransformFeedbackAttribsNV"))
	gpTransformFeedbackBufferBase = uintptr(getProcAddr("glTransformFeedbackBufferBase"))
	gpTransformFeedbackBufferRange = uintptr(getProcAddr("glTransformFeedbackBufferRange"))
	gpTransformFeedbackStreamAttribsNV = uintptr(getProcAddr("glTransformFeedbackStreamAttribsNV"))
	gpTransformFeedbackVaryingsEXT = uintptr(getProcAddr("glTransformFeedbackVaryingsEXT"))
	gpTransformFeedbackVaryingsNV = uintptr(getProcAddr("glTransformFeedbackVaryingsNV"))
	gpTransformPathNV = uintptr(getProcAddr("glTransformPathNV"))
	gpTranslated = uintptr(getProcAddr("glTranslated"))
	if gpTranslated == 0 {
		return errors.New("glTranslated")
	}
	gpTranslatef = uintptr(getProcAddr("glTranslatef"))
	if gpTranslatef == 0 {
		return errors.New("glTranslatef")
	}
	gpTranslatexOES = uintptr(getProcAddr("glTranslatexOES"))
	gpUniform1d = uintptr(getProcAddr("glUniform1d"))
	gpUniform1dv = uintptr(getProcAddr("glUniform1dv"))
	gpUniform1f = uintptr(getProcAddr("glUniform1f"))
	if gpUniform1f == 0 {
		return errors.New("glUniform1f")
	}
	gpUniform1fARB = uintptr(getProcAddr("glUniform1fARB"))
	gpUniform1fv = uintptr(getProcAddr("glUniform1fv"))
	if gpUniform1fv == 0 {
		return errors.New("glUniform1fv")
	}
	gpUniform1fvARB = uintptr(getProcAddr("glUniform1fvARB"))
	gpUniform1i = uintptr(getProcAddr("glUniform1i"))
	if gpUniform1i == 0 {
		return errors.New("glUniform1i")
	}
	gpUniform1i64ARB = uintptr(getProcAddr("glUniform1i64ARB"))
	gpUniform1i64NV = uintptr(getProcAddr("glUniform1i64NV"))
	gpUniform1i64vARB = uintptr(getProcAddr("glUniform1i64vARB"))
	gpUniform1i64vNV = uintptr(getProcAddr("glUniform1i64vNV"))
	gpUniform1iARB = uintptr(getProcAddr("glUniform1iARB"))
	gpUniform1iv = uintptr(getProcAddr("glUniform1iv"))
	if gpUniform1iv == 0 {
		return errors.New("glUniform1iv")
	}
	gpUniform1ivARB = uintptr(getProcAddr("glUniform1ivARB"))
	gpUniform1ui64ARB = uintptr(getProcAddr("glUniform1ui64ARB"))
	gpUniform1ui64NV = uintptr(getProcAddr("glUniform1ui64NV"))
	gpUniform1ui64vARB = uintptr(getProcAddr("glUniform1ui64vARB"))
	gpUniform1ui64vNV = uintptr(getProcAddr("glUniform1ui64vNV"))
	gpUniform1uiEXT = uintptr(getProcAddr("glUniform1uiEXT"))
	gpUniform1uivEXT = uintptr(getProcAddr("glUniform1uivEXT"))
	gpUniform2d = uintptr(getProcAddr("glUniform2d"))
	gpUniform2dv = uintptr(getProcAddr("glUniform2dv"))
	gpUniform2f = uintptr(getProcAddr("glUniform2f"))
	if gpUniform2f == 0 {
		return errors.New("glUniform2f")
	}
	gpUniform2fARB = uintptr(getProcAddr("glUniform2fARB"))
	gpUniform2fv = uintptr(getProcAddr("glUniform2fv"))
	if gpUniform2fv == 0 {
		return errors.New("glUniform2fv")
	}
	gpUniform2fvARB = uintptr(getProcAddr("glUniform2fvARB"))
	gpUniform2i = uintptr(getProcAddr("glUniform2i"))
	if gpUniform2i == 0 {
		return errors.New("glUniform2i")
	}
	gpUniform2i64ARB = uintptr(getProcAddr("glUniform2i64ARB"))
	gpUniform2i64NV = uintptr(getProcAddr("glUniform2i64NV"))
	gpUniform2i64vARB = uintptr(getProcAddr("glUniform2i64vARB"))
	gpUniform2i64vNV = uintptr(getProcAddr("glUniform2i64vNV"))
	gpUniform2iARB = uintptr(getProcAddr("glUniform2iARB"))
	gpUniform2iv = uintptr(getProcAddr("glUniform2iv"))
	if gpUniform2iv == 0 {
		return errors.New("glUniform2iv")
	}
	gpUniform2ivARB = uintptr(getProcAddr("glUniform2ivARB"))
	gpUniform2ui64ARB = uintptr(getProcAddr("glUniform2ui64ARB"))
	gpUniform2ui64NV = uintptr(getProcAddr("glUniform2ui64NV"))
	gpUniform2ui64vARB = uintptr(getProcAddr("glUniform2ui64vARB"))
	gpUniform2ui64vNV = uintptr(getProcAddr("glUniform2ui64vNV"))
	gpUniform2uiEXT = uintptr(getProcAddr("glUniform2uiEXT"))
	gpUniform2uivEXT = uintptr(getProcAddr("glUniform2uivEXT"))
	gpUniform3d = uintptr(getProcAddr("glUniform3d"))
	gpUniform3dv = uintptr(getProcAddr("glUniform3dv"))
	gpUniform3f = uintptr(getProcAddr("glUniform3f"))
	if gpUniform3f == 0 {
		return errors.New("glUniform3f")
	}
	gpUniform3fARB = uintptr(getProcAddr("glUniform3fARB"))
	gpUniform3fv = uintptr(getProcAddr("glUniform3fv"))
	if gpUniform3fv == 0 {
		return errors.New("glUniform3fv")
	}
	gpUniform3fvARB = uintptr(getProcAddr("glUniform3fvARB"))
	gpUniform3i = uintptr(getProcAddr("glUniform3i"))
	if gpUniform3i == 0 {
		return errors.New("glUniform3i")
	}
	gpUniform3i64ARB = uintptr(getProcAddr("glUniform3i64ARB"))
	gpUniform3i64NV = uintptr(getProcAddr("glUniform3i64NV"))
	gpUniform3i64vARB = uintptr(getProcAddr("glUniform3i64vARB"))
	gpUniform3i64vNV = uintptr(getProcAddr("glUniform3i64vNV"))
	gpUniform3iARB = uintptr(getProcAddr("glUniform3iARB"))
	gpUniform3iv = uintptr(getProcAddr("glUniform3iv"))
	if gpUniform3iv == 0 {
		return errors.New("glUniform3iv")
	}
	gpUniform3ivARB = uintptr(getProcAddr("glUniform3ivARB"))
	gpUniform3ui64ARB = uintptr(getProcAddr("glUniform3ui64ARB"))
	gpUniform3ui64NV = uintptr(getProcAddr("glUniform3ui64NV"))
	gpUniform3ui64vARB = uintptr(getProcAddr("glUniform3ui64vARB"))
	gpUniform3ui64vNV = uintptr(getProcAddr("glUniform3ui64vNV"))
	gpUniform3uiEXT = uintptr(getProcAddr("glUniform3uiEXT"))
	gpUniform3uivEXT = uintptr(getProcAddr("glUniform3uivEXT"))
	gpUniform4d = uintptr(getProcAddr("glUniform4d"))
	gpUniform4dv = uintptr(getProcAddr("glUniform4dv"))
	gpUniform4f = uintptr(getProcAddr("glUniform4f"))
	if gpUniform4f == 0 {
		return errors.New("glUniform4f")
	}
	gpUniform4fARB = uintptr(getProcAddr("glUniform4fARB"))
	gpUniform4fv = uintptr(getProcAddr("glUniform4fv"))
	if gpUniform4fv == 0 {
		return errors.New("glUniform4fv")
	}
	gpUniform4fvARB = uintptr(getProcAddr("glUniform4fvARB"))
	gpUniform4i = uintptr(getProcAddr("glUniform4i"))
	if gpUniform4i == 0 {
		return errors.New("glUniform4i")
	}
	gpUniform4i64ARB = uintptr(getProcAddr("glUniform4i64ARB"))
	gpUniform4i64NV = uintptr(getProcAddr("glUniform4i64NV"))
	gpUniform4i64vARB = uintptr(getProcAddr("glUniform4i64vARB"))
	gpUniform4i64vNV = uintptr(getProcAddr("glUniform4i64vNV"))
	gpUniform4iARB = uintptr(getProcAddr("glUniform4iARB"))
	gpUniform4iv = uintptr(getProcAddr("glUniform4iv"))
	if gpUniform4iv == 0 {
		return errors.New("glUniform4iv")
	}
	gpUniform4ivARB = uintptr(getProcAddr("glUniform4ivARB"))
	gpUniform4ui64ARB = uintptr(getProcAddr("glUniform4ui64ARB"))
	gpUniform4ui64NV = uintptr(getProcAddr("glUniform4ui64NV"))
	gpUniform4ui64vARB = uintptr(getProcAddr("glUniform4ui64vARB"))
	gpUniform4ui64vNV = uintptr(getProcAddr("glUniform4ui64vNV"))
	gpUniform4uiEXT = uintptr(getProcAddr("glUniform4uiEXT"))
	gpUniform4uivEXT = uintptr(getProcAddr("glUniform4uivEXT"))
	gpUniformBlockBinding = uintptr(getProcAddr("glUniformBlockBinding"))
	gpUniformBufferEXT = uintptr(getProcAddr("glUniformBufferEXT"))
	gpUniformHandleui64ARB = uintptr(getProcAddr("glUniformHandleui64ARB"))
	gpUniformHandleui64NV = uintptr(getProcAddr("glUniformHandleui64NV"))
	gpUniformHandleui64vARB = uintptr(getProcAddr("glUniformHandleui64vARB"))
	gpUniformHandleui64vNV = uintptr(getProcAddr("glUniformHandleui64vNV"))
	gpUniformMatrix2dv = uintptr(getProcAddr("glUniformMatrix2dv"))
	gpUniformMatrix2fv = uintptr(getProcAddr("glUniformMatrix2fv"))
	if gpUniformMatrix2fv == 0 {
		return errors.New("glUniformMatrix2fv")
	}
	gpUniformMatrix2fvARB = uintptr(getProcAddr("glUniformMatrix2fvARB"))
	gpUniformMatrix2x3dv = uintptr(getProcAddr("glUniformMatrix2x3dv"))
	gpUniformMatrix2x3fv = uintptr(getProcAddr("glUniformMatrix2x3fv"))
	if gpUniformMatrix2x3fv == 0 {
		return errors.New("glUniformMatrix2x3fv")
	}
	gpUniformMatrix2x4dv = uintptr(getProcAddr("glUniformMatrix2x4dv"))
	gpUniformMatrix2x4fv = uintptr(getProcAddr("glUniformMatrix2x4fv"))
	if gpUniformMatrix2x4fv == 0 {
		return errors.New("glUniformMatrix2x4fv")
	}
	gpUniformMatrix3dv = uintptr(getProcAddr("glUniformMatrix3dv"))
	gpUniformMatrix3fv = uintptr(getProcAddr("glUniformMatrix3fv"))
	if gpUniformMatrix3fv == 0 {
		return errors.New("glUniformMatrix3fv")
	}
	gpUniformMatrix3fvARB = uintptr(getProcAddr("glUniformMatrix3fvARB"))
	gpUniformMatrix3x2dv = uintptr(getProcAddr("glUniformMatrix3x2dv"))
	gpUniformMatrix3x2fv = uintptr(getProcAddr("glUniformMatrix3x2fv"))
	if gpUniformMatrix3x2fv == 0 {
		return errors.New("glUniformMatrix3x2fv")
	}
	gpUniformMatrix3x4dv = uintptr(getProcAddr("glUniformMatrix3x4dv"))
	gpUniformMatrix3x4fv = uintptr(getProcAddr("glUniformMatrix3x4fv"))
	if gpUniformMatrix3x4fv == 0 {
		return errors.New("glUniformMatrix3x4fv")
	}
	gpUniformMatrix4dv = uintptr(getProcAddr("glUniformMatrix4dv"))
	gpUniformMatrix4fv = uintptr(getProcAddr("glUniformMatrix4fv"))
	if gpUniformMatrix4fv == 0 {
		return errors.New("glUniformMatrix4fv")
	}
	gpUniformMatrix4fvARB = uintptr(getProcAddr("glUniformMatrix4fvARB"))
	gpUniformMatrix4x2dv = uintptr(getProcAddr("glUniformMatrix4x2dv"))
	gpUniformMatrix4x2fv = uintptr(getProcAddr("glUniformMatrix4x2fv"))
	if gpUniformMatrix4x2fv == 0 {
		return errors.New("glUniformMatrix4x2fv")
	}
	gpUniformMatrix4x3dv = uintptr(getProcAddr("glUniformMatrix4x3dv"))
	gpUniformMatrix4x3fv = uintptr(getProcAddr("glUniformMatrix4x3fv"))
	if gpUniformMatrix4x3fv == 0 {
		return errors.New("glUniformMatrix4x3fv")
	}
	gpUniformSubroutinesuiv = uintptr(getProcAddr("glUniformSubroutinesuiv"))
	gpUniformui64NV = uintptr(getProcAddr("glUniformui64NV"))
	gpUniformui64vNV = uintptr(getProcAddr("glUniformui64vNV"))
	gpUnlockArraysEXT = uintptr(getProcAddr("glUnlockArraysEXT"))
	gpUnmapBuffer = uintptr(getProcAddr("glUnmapBuffer"))
	if gpUnmapBuffer == 0 {
		return errors.New("glUnmapBuffer")
	}
	gpUnmapBufferARB = uintptr(getProcAddr("glUnmapBufferARB"))
	gpUnmapNamedBuffer = uintptr(getProcAddr("glUnmapNamedBuffer"))
	gpUnmapNamedBufferEXT = uintptr(getProcAddr("glUnmapNamedBufferEXT"))
	gpUnmapObjectBufferATI = uintptr(getProcAddr("glUnmapObjectBufferATI"))
	gpUnmapTexture2DINTEL = uintptr(getProcAddr("glUnmapTexture2DINTEL"))
	gpUpdateObjectBufferATI = uintptr(getProcAddr("glUpdateObjectBufferATI"))
	gpUseProgram = uintptr(getProcAddr("glUseProgram"))
	if gpUseProgram == 0 {
		return errors.New("glUseProgram")
	}
	gpUseProgramObjectARB = uintptr(getProcAddr("glUseProgramObjectARB"))
	gpUseProgramStages = uintptr(getProcAddr("glUseProgramStages"))
	gpUseProgramStagesEXT = uintptr(getProcAddr("glUseProgramStagesEXT"))
	gpUseShaderProgramEXT = uintptr(getProcAddr("glUseShaderProgramEXT"))
	gpVDPAUFiniNV = uintptr(getProcAddr("glVDPAUFiniNV"))
	gpVDPAUGetSurfaceivNV = uintptr(getProcAddr("glVDPAUGetSurfaceivNV"))
	gpVDPAUInitNV = uintptr(getProcAddr("glVDPAUInitNV"))
	gpVDPAUIsSurfaceNV = uintptr(getProcAddr("glVDPAUIsSurfaceNV"))
	gpVDPAUMapSurfacesNV = uintptr(getProcAddr("glVDPAUMapSurfacesNV"))
	gpVDPAURegisterOutputSurfaceNV = uintptr(getProcAddr("glVDPAURegisterOutputSurfaceNV"))
	gpVDPAURegisterVideoSurfaceNV = uintptr(getProcAddr("glVDPAURegisterVideoSurfaceNV"))
	gpVDPAUSurfaceAccessNV = uintptr(getProcAddr("glVDPAUSurfaceAccessNV"))
	gpVDPAUUnmapSurfacesNV = uintptr(getProcAddr("glVDPAUUnmapSurfacesNV"))
	gpVDPAUUnregisterSurfaceNV = uintptr(getProcAddr("glVDPAUUnregisterSurfaceNV"))
	gpValidateProgram = uintptr(getProcAddr("glValidateProgram"))
	if gpValidateProgram == 0 {
		return errors.New("glValidateProgram")
	}
	gpValidateProgramARB = uintptr(getProcAddr("glValidateProgramARB"))
	gpValidateProgramPipeline = uintptr(getProcAddr("glValidateProgramPipeline"))
	gpValidateProgramPipelineEXT = uintptr(getProcAddr("glValidateProgramPipelineEXT"))
	gpVariantArrayObjectATI = uintptr(getProcAddr("glVariantArrayObjectATI"))
	gpVariantPointerEXT = uintptr(getProcAddr("glVariantPointerEXT"))
	gpVariantbvEXT = uintptr(getProcAddr("glVariantbvEXT"))
	gpVariantdvEXT = uintptr(getProcAddr("glVariantdvEXT"))
	gpVariantfvEXT = uintptr(getProcAddr("glVariantfvEXT"))
	gpVariantivEXT = uintptr(getProcAddr("glVariantivEXT"))
	gpVariantsvEXT = uintptr(getProcAddr("glVariantsvEXT"))
	gpVariantubvEXT = uintptr(getProcAddr("glVariantubvEXT"))
	gpVariantuivEXT = uintptr(getProcAddr("glVariantuivEXT"))
	gpVariantusvEXT = uintptr(getProcAddr("glVariantusvEXT"))
	gpVertex2bOES = uintptr(getProcAddr("glVertex2bOES"))
	gpVertex2bvOES = uintptr(getProcAddr("glVertex2bvOES"))
	gpVertex2d = uintptr(getProcAddr("glVertex2d"))
	if gpVertex2d == 0 {
		return errors.New("glVertex2d")
	}
	gpVertex2dv = uintptr(getProcAddr("glVertex2dv"))
	if gpVertex2dv == 0 {
		return errors.New("glVertex2dv")
	}
	gpVertex2f = uintptr(getProcAddr("glVertex2f"))
	if gpVertex2f == 0 {
		return errors.New("glVertex2f")
	}
	gpVertex2fv = uintptr(getProcAddr("glVertex2fv"))
	if gpVertex2fv == 0 {
		return errors.New("glVertex2fv")
	}
	gpVertex2hNV = uintptr(getProcAddr("glVertex2hNV"))
	gpVertex2hvNV = uintptr(getProcAddr("glVertex2hvNV"))
	gpVertex2i = uintptr(getProcAddr("glVertex2i"))
	if gpVertex2i == 0 {
		return errors.New("glVertex2i")
	}
	gpVertex2iv = uintptr(getProcAddr("glVertex2iv"))
	if gpVertex2iv == 0 {
		return errors.New("glVertex2iv")
	}
	gpVertex2s = uintptr(getProcAddr("glVertex2s"))
	if gpVertex2s == 0 {
		return errors.New("glVertex2s")
	}
	gpVertex2sv = uintptr(getProcAddr("glVertex2sv"))
	if gpVertex2sv == 0 {
		return errors.New("glVertex2sv")
	}
	gpVertex2xOES = uintptr(getProcAddr("glVertex2xOES"))
	gpVertex2xvOES = uintptr(getProcAddr("glVertex2xvOES"))
	gpVertex3bOES = uintptr(getProcAddr("glVertex3bOES"))
	gpVertex3bvOES = uintptr(getProcAddr("glVertex3bvOES"))
	gpVertex3d = uintptr(getProcAddr("glVertex3d"))
	if gpVertex3d == 0 {
		return errors.New("glVertex3d")
	}
	gpVertex3dv = uintptr(getProcAddr("glVertex3dv"))
	if gpVertex3dv == 0 {
		return errors.New("glVertex3dv")
	}
	gpVertex3f = uintptr(getProcAddr("glVertex3f"))
	if gpVertex3f == 0 {
		return errors.New("glVertex3f")
	}
	gpVertex3fv = uintptr(getProcAddr("glVertex3fv"))
	if gpVertex3fv == 0 {
		return errors.New("glVertex3fv")
	}
	gpVertex3hNV = uintptr(getProcAddr("glVertex3hNV"))
	gpVertex3hvNV = uintptr(getProcAddr("glVertex3hvNV"))
	gpVertex3i = uintptr(getProcAddr("glVertex3i"))
	if gpVertex3i == 0 {
		return errors.New("glVertex3i")
	}
	gpVertex3iv = uintptr(getProcAddr("glVertex3iv"))
	if gpVertex3iv == 0 {
		return errors.New("glVertex3iv")
	}
	gpVertex3s = uintptr(getProcAddr("glVertex3s"))
	if gpVertex3s == 0 {
		return errors.New("glVertex3s")
	}
	gpVertex3sv = uintptr(getProcAddr("glVertex3sv"))
	if gpVertex3sv == 0 {
		return errors.New("glVertex3sv")
	}
	gpVertex3xOES = uintptr(getProcAddr("glVertex3xOES"))
	gpVertex3xvOES = uintptr(getProcAddr("glVertex3xvOES"))
	gpVertex4bOES = uintptr(getProcAddr("glVertex4bOES"))
	gpVertex4bvOES = uintptr(getProcAddr("glVertex4bvOES"))
	gpVertex4d = uintptr(getProcAddr("glVertex4d"))
	if gpVertex4d == 0 {
		return errors.New("glVertex4d")
	}
	gpVertex4dv = uintptr(getProcAddr("glVertex4dv"))
	if gpVertex4dv == 0 {
		return errors.New("glVertex4dv")
	}
	gpVertex4f = uintptr(getProcAddr("glVertex4f"))
	if gpVertex4f == 0 {
		return errors.New("glVertex4f")
	}
	gpVertex4fv = uintptr(getProcAddr("glVertex4fv"))
	if gpVertex4fv == 0 {
		return errors.New("glVertex4fv")
	}
	gpVertex4hNV = uintptr(getProcAddr("glVertex4hNV"))
	gpVertex4hvNV = uintptr(getProcAddr("glVertex4hvNV"))
	gpVertex4i = uintptr(getProcAddr("glVertex4i"))
	if gpVertex4i == 0 {
		return errors.New("glVertex4i")
	}
	gpVertex4iv = uintptr(getProcAddr("glVertex4iv"))
	if gpVertex4iv == 0 {
		return errors.New("glVertex4iv")
	}
	gpVertex4s = uintptr(getProcAddr("glVertex4s"))
	if gpVertex4s == 0 {
		return errors.New("glVertex4s")
	}
	gpVertex4sv = uintptr(getProcAddr("glVertex4sv"))
	if gpVertex4sv == 0 {
		return errors.New("glVertex4sv")
	}
	gpVertex4xOES = uintptr(getProcAddr("glVertex4xOES"))
	gpVertex4xvOES = uintptr(getProcAddr("glVertex4xvOES"))
	gpVertexArrayAttribBinding = uintptr(getProcAddr("glVertexArrayAttribBinding"))
	gpVertexArrayAttribFormat = uintptr(getProcAddr("glVertexArrayAttribFormat"))
	gpVertexArrayAttribIFormat = uintptr(getProcAddr("glVertexArrayAttribIFormat"))
	gpVertexArrayAttribLFormat = uintptr(getProcAddr("glVertexArrayAttribLFormat"))
	gpVertexArrayBindVertexBufferEXT = uintptr(getProcAddr("glVertexArrayBindVertexBufferEXT"))
	gpVertexArrayBindingDivisor = uintptr(getProcAddr("glVertexArrayBindingDivisor"))
	gpVertexArrayColorOffsetEXT = uintptr(getProcAddr("glVertexArrayColorOffsetEXT"))
	gpVertexArrayEdgeFlagOffsetEXT = uintptr(getProcAddr("glVertexArrayEdgeFlagOffsetEXT"))
	gpVertexArrayElementBuffer = uintptr(getProcAddr("glVertexArrayElementBuffer"))
	gpVertexArrayFogCoordOffsetEXT = uintptr(getProcAddr("glVertexArrayFogCoordOffsetEXT"))
	gpVertexArrayIndexOffsetEXT = uintptr(getProcAddr("glVertexArrayIndexOffsetEXT"))
	gpVertexArrayMultiTexCoordOffsetEXT = uintptr(getProcAddr("glVertexArrayMultiTexCoordOffsetEXT"))
	gpVertexArrayNormalOffsetEXT = uintptr(getProcAddr("glVertexArrayNormalOffsetEXT"))
	gpVertexArrayParameteriAPPLE = uintptr(getProcAddr("glVertexArrayParameteriAPPLE"))
	gpVertexArrayRangeAPPLE = uintptr(getProcAddr("glVertexArrayRangeAPPLE"))
	gpVertexArrayRangeNV = uintptr(getProcAddr("glVertexArrayRangeNV"))
	gpVertexArraySecondaryColorOffsetEXT = uintptr(getProcAddr("glVertexArraySecondaryColorOffsetEXT"))
	gpVertexArrayTexCoordOffsetEXT = uintptr(getProcAddr("glVertexArrayTexCoordOffsetEXT"))
	gpVertexArrayVertexAttribBindingEXT = uintptr(getProcAddr("glVertexArrayVertexAttribBindingEXT"))
	gpVertexArrayVertexAttribDivisorEXT = uintptr(getProcAddr("glVertexArrayVertexAttribDivisorEXT"))
	gpVertexArrayVertexAttribFormatEXT = uintptr(getProcAddr("glVertexArrayVertexAttribFormatEXT"))
	gpVertexArrayVertexAttribIFormatEXT = uintptr(getProcAddr("glVertexArrayVertexAttribIFormatEXT"))
	gpVertexArrayVertexAttribIOffsetEXT = uintptr(getProcAddr("glVertexArrayVertexAttribIOffsetEXT"))
	gpVertexArrayVertexAttribLFormatEXT = uintptr(getProcAddr("glVertexArrayVertexAttribLFormatEXT"))
	gpVertexArrayVertexAttribLOffsetEXT = uintptr(getProcAddr("glVertexArrayVertexAttribLOffsetEXT"))
	gpVertexArrayVertexAttribOffsetEXT = uintptr(getProcAddr("glVertexArrayVertexAttribOffsetEXT"))
	gpVertexArrayVertexBindingDivisorEXT = uintptr(getProcAddr("glVertexArrayVertexBindingDivisorEXT"))
	gpVertexArrayVertexBuffer = uintptr(getProcAddr("glVertexArrayVertexBuffer"))
	gpVertexArrayVertexBuffers = uintptr(getProcAddr("glVertexArrayVertexBuffers"))
	gpVertexArrayVertexOffsetEXT = uintptr(getProcAddr("glVertexArrayVertexOffsetEXT"))
	gpVertexAttrib1d = uintptr(getProcAddr("glVertexAttrib1d"))
	if gpVertexAttrib1d == 0 {
		return errors.New("glVertexAttrib1d")
	}
	gpVertexAttrib1dARB = uintptr(getProcAddr("glVertexAttrib1dARB"))
	gpVertexAttrib1dNV = uintptr(getProcAddr("glVertexAttrib1dNV"))
	gpVertexAttrib1dv = uintptr(getProcAddr("glVertexAttrib1dv"))
	if gpVertexAttrib1dv == 0 {
		return errors.New("glVertexAttrib1dv")
	}
	gpVertexAttrib1dvARB = uintptr(getProcAddr("glVertexAttrib1dvARB"))
	gpVertexAttrib1dvNV = uintptr(getProcAddr("glVertexAttrib1dvNV"))
	gpVertexAttrib1f = uintptr(getProcAddr("glVertexAttrib1f"))
	if gpVertexAttrib1f == 0 {
		return errors.New("glVertexAttrib1f")
	}
	gpVertexAttrib1fARB = uintptr(getProcAddr("glVertexAttrib1fARB"))
	gpVertexAttrib1fNV = uintptr(getProcAddr("glVertexAttrib1fNV"))
	gpVertexAttrib1fv = uintptr(getProcAddr("glVertexAttrib1fv"))
	if gpVertexAttrib1fv == 0 {
		return errors.New("glVertexAttrib1fv")
	}
	gpVertexAttrib1fvARB = uintptr(getProcAddr("glVertexAttrib1fvARB"))
	gpVertexAttrib1fvNV = uintptr(getProcAddr("glVertexAttrib1fvNV"))
	gpVertexAttrib1hNV = uintptr(getProcAddr("glVertexAttrib1hNV"))
	gpVertexAttrib1hvNV = uintptr(getProcAddr("glVertexAttrib1hvNV"))
	gpVertexAttrib1s = uintptr(getProcAddr("glVertexAttrib1s"))
	if gpVertexAttrib1s == 0 {
		return errors.New("glVertexAttrib1s")
	}
	gpVertexAttrib1sARB = uintptr(getProcAddr("glVertexAttrib1sARB"))
	gpVertexAttrib1sNV = uintptr(getProcAddr("glVertexAttrib1sNV"))
	gpVertexAttrib1sv = uintptr(getProcAddr("glVertexAttrib1sv"))
	if gpVertexAttrib1sv == 0 {
		return errors.New("glVertexAttrib1sv")
	}
	gpVertexAttrib1svARB = uintptr(getProcAddr("glVertexAttrib1svARB"))
	gpVertexAttrib1svNV = uintptr(getProcAddr("glVertexAttrib1svNV"))
	gpVertexAttrib2d = uintptr(getProcAddr("glVertexAttrib2d"))
	if gpVertexAttrib2d == 0 {
		return errors.New("glVertexAttrib2d")
	}
	gpVertexAttrib2dARB = uintptr(getProcAddr("glVertexAttrib2dARB"))
	gpVertexAttrib2dNV = uintptr(getProcAddr("glVertexAttrib2dNV"))
	gpVertexAttrib2dv = uintptr(getProcAddr("glVertexAttrib2dv"))
	if gpVertexAttrib2dv == 0 {
		return errors.New("glVertexAttrib2dv")
	}
	gpVertexAttrib2dvARB = uintptr(getProcAddr("glVertexAttrib2dvARB"))
	gpVertexAttrib2dvNV = uintptr(getProcAddr("glVertexAttrib2dvNV"))
	gpVertexAttrib2f = uintptr(getProcAddr("glVertexAttrib2f"))
	if gpVertexAttrib2f == 0 {
		return errors.New("glVertexAttrib2f")
	}
	gpVertexAttrib2fARB = uintptr(getProcAddr("glVertexAttrib2fARB"))
	gpVertexAttrib2fNV = uintptr(getProcAddr("glVertexAttrib2fNV"))
	gpVertexAttrib2fv = uintptr(getProcAddr("glVertexAttrib2fv"))
	if gpVertexAttrib2fv == 0 {
		return errors.New("glVertexAttrib2fv")
	}
	gpVertexAttrib2fvARB = uintptr(getProcAddr("glVertexAttrib2fvARB"))
	gpVertexAttrib2fvNV = uintptr(getProcAddr("glVertexAttrib2fvNV"))
	gpVertexAttrib2hNV = uintptr(getProcAddr("glVertexAttrib2hNV"))
	gpVertexAttrib2hvNV = uintptr(getProcAddr("glVertexAttrib2hvNV"))
	gpVertexAttrib2s = uintptr(getProcAddr("glVertexAttrib2s"))
	if gpVertexAttrib2s == 0 {
		return errors.New("glVertexAttrib2s")
	}
	gpVertexAttrib2sARB = uintptr(getProcAddr("glVertexAttrib2sARB"))
	gpVertexAttrib2sNV = uintptr(getProcAddr("glVertexAttrib2sNV"))
	gpVertexAttrib2sv = uintptr(getProcAddr("glVertexAttrib2sv"))
	if gpVertexAttrib2sv == 0 {
		return errors.New("glVertexAttrib2sv")
	}
	gpVertexAttrib2svARB = uintptr(getProcAddr("glVertexAttrib2svARB"))
	gpVertexAttrib2svNV = uintptr(getProcAddr("glVertexAttrib2svNV"))
	gpVertexAttrib3d = uintptr(getProcAddr("glVertexAttrib3d"))
	if gpVertexAttrib3d == 0 {
		return errors.New("glVertexAttrib3d")
	}
	gpVertexAttrib3dARB = uintptr(getProcAddr("glVertexAttrib3dARB"))
	gpVertexAttrib3dNV = uintptr(getProcAddr("glVertexAttrib3dNV"))
	gpVertexAttrib3dv = uintptr(getProcAddr("glVertexAttrib3dv"))
	if gpVertexAttrib3dv == 0 {
		return errors.New("glVertexAttrib3dv")
	}
	gpVertexAttrib3dvARB = uintptr(getProcAddr("glVertexAttrib3dvARB"))
	gpVertexAttrib3dvNV = uintptr(getProcAddr("glVertexAttrib3dvNV"))
	gpVertexAttrib3f = uintptr(getProcAddr("glVertexAttrib3f"))
	if gpVertexAttrib3f == 0 {
		return errors.New("glVertexAttrib3f")
	}
	gpVertexAttrib3fARB = uintptr(getProcAddr("glVertexAttrib3fARB"))
	gpVertexAttrib3fNV = uintptr(getProcAddr("glVertexAttrib3fNV"))
	gpVertexAttrib3fv = uintptr(getProcAddr("glVertexAttrib3fv"))
	if gpVertexAttrib3fv == 0 {
		return errors.New("glVertexAttrib3fv")
	}
	gpVertexAttrib3fvARB = uintptr(getProcAddr("glVertexAttrib3fvARB"))
	gpVertexAttrib3fvNV = uintptr(getProcAddr("glVertexAttrib3fvNV"))
	gpVertexAttrib3hNV = uintptr(getProcAddr("glVertexAttrib3hNV"))
	gpVertexAttrib3hvNV = uintptr(getProcAddr("glVertexAttrib3hvNV"))
	gpVertexAttrib3s = uintptr(getProcAddr("glVertexAttrib3s"))
	if gpVertexAttrib3s == 0 {
		return errors.New("glVertexAttrib3s")
	}
	gpVertexAttrib3sARB = uintptr(getProcAddr("glVertexAttrib3sARB"))
	gpVertexAttrib3sNV = uintptr(getProcAddr("glVertexAttrib3sNV"))
	gpVertexAttrib3sv = uintptr(getProcAddr("glVertexAttrib3sv"))
	if gpVertexAttrib3sv == 0 {
		return errors.New("glVertexAttrib3sv")
	}
	gpVertexAttrib3svARB = uintptr(getProcAddr("glVertexAttrib3svARB"))
	gpVertexAttrib3svNV = uintptr(getProcAddr("glVertexAttrib3svNV"))
	gpVertexAttrib4Nbv = uintptr(getProcAddr("glVertexAttrib4Nbv"))
	if gpVertexAttrib4Nbv == 0 {
		return errors.New("glVertexAttrib4Nbv")
	}
	gpVertexAttrib4NbvARB = uintptr(getProcAddr("glVertexAttrib4NbvARB"))
	gpVertexAttrib4Niv = uintptr(getProcAddr("glVertexAttrib4Niv"))
	if gpVertexAttrib4Niv == 0 {
		return errors.New("glVertexAttrib4Niv")
	}
	gpVertexAttrib4NivARB = uintptr(getProcAddr("glVertexAttrib4NivARB"))
	gpVertexAttrib4Nsv = uintptr(getProcAddr("glVertexAttrib4Nsv"))
	if gpVertexAttrib4Nsv == 0 {
		return errors.New("glVertexAttrib4Nsv")
	}
	gpVertexAttrib4NsvARB = uintptr(getProcAddr("glVertexAttrib4NsvARB"))
	gpVertexAttrib4Nub = uintptr(getProcAddr("glVertexAttrib4Nub"))
	if gpVertexAttrib4Nub == 0 {
		return errors.New("glVertexAttrib4Nub")
	}
	gpVertexAttrib4NubARB = uintptr(getProcAddr("glVertexAttrib4NubARB"))
	gpVertexAttrib4Nubv = uintptr(getProcAddr("glVertexAttrib4Nubv"))
	if gpVertexAttrib4Nubv == 0 {
		return errors.New("glVertexAttrib4Nubv")
	}
	gpVertexAttrib4NubvARB = uintptr(getProcAddr("glVertexAttrib4NubvARB"))
	gpVertexAttrib4Nuiv = uintptr(getProcAddr("glVertexAttrib4Nuiv"))
	if gpVertexAttrib4Nuiv == 0 {
		return errors.New("glVertexAttrib4Nuiv")
	}
	gpVertexAttrib4NuivARB = uintptr(getProcAddr("glVertexAttrib4NuivARB"))
	gpVertexAttrib4Nusv = uintptr(getProcAddr("glVertexAttrib4Nusv"))
	if gpVertexAttrib4Nusv == 0 {
		return errors.New("glVertexAttrib4Nusv")
	}
	gpVertexAttrib4NusvARB = uintptr(getProcAddr("glVertexAttrib4NusvARB"))
	gpVertexAttrib4bv = uintptr(getProcAddr("glVertexAttrib4bv"))
	if gpVertexAttrib4bv == 0 {
		return errors.New("glVertexAttrib4bv")
	}
	gpVertexAttrib4bvARB = uintptr(getProcAddr("glVertexAttrib4bvARB"))
	gpVertexAttrib4d = uintptr(getProcAddr("glVertexAttrib4d"))
	if gpVertexAttrib4d == 0 {
		return errors.New("glVertexAttrib4d")
	}
	gpVertexAttrib4dARB = uintptr(getProcAddr("glVertexAttrib4dARB"))
	gpVertexAttrib4dNV = uintptr(getProcAddr("glVertexAttrib4dNV"))
	gpVertexAttrib4dv = uintptr(getProcAddr("glVertexAttrib4dv"))
	if gpVertexAttrib4dv == 0 {
		return errors.New("glVertexAttrib4dv")
	}
	gpVertexAttrib4dvARB = uintptr(getProcAddr("glVertexAttrib4dvARB"))
	gpVertexAttrib4dvNV = uintptr(getProcAddr("glVertexAttrib4dvNV"))
	gpVertexAttrib4f = uintptr(getProcAddr("glVertexAttrib4f"))
	if gpVertexAttrib4f == 0 {
		return errors.New("glVertexAttrib4f")
	}
	gpVertexAttrib4fARB = uintptr(getProcAddr("glVertexAttrib4fARB"))
	gpVertexAttrib4fNV = uintptr(getProcAddr("glVertexAttrib4fNV"))
	gpVertexAttrib4fv = uintptr(getProcAddr("glVertexAttrib4fv"))
	if gpVertexAttrib4fv == 0 {
		return errors.New("glVertexAttrib4fv")
	}
	gpVertexAttrib4fvARB = uintptr(getProcAddr("glVertexAttrib4fvARB"))
	gpVertexAttrib4fvNV = uintptr(getProcAddr("glVertexAttrib4fvNV"))
	gpVertexAttrib4hNV = uintptr(getProcAddr("glVertexAttrib4hNV"))
	gpVertexAttrib4hvNV = uintptr(getProcAddr("glVertexAttrib4hvNV"))
	gpVertexAttrib4iv = uintptr(getProcAddr("glVertexAttrib4iv"))
	if gpVertexAttrib4iv == 0 {
		return errors.New("glVertexAttrib4iv")
	}
	gpVertexAttrib4ivARB = uintptr(getProcAddr("glVertexAttrib4ivARB"))
	gpVertexAttrib4s = uintptr(getProcAddr("glVertexAttrib4s"))
	if gpVertexAttrib4s == 0 {
		return errors.New("glVertexAttrib4s")
	}
	gpVertexAttrib4sARB = uintptr(getProcAddr("glVertexAttrib4sARB"))
	gpVertexAttrib4sNV = uintptr(getProcAddr("glVertexAttrib4sNV"))
	gpVertexAttrib4sv = uintptr(getProcAddr("glVertexAttrib4sv"))
	if gpVertexAttrib4sv == 0 {
		return errors.New("glVertexAttrib4sv")
	}
	gpVertexAttrib4svARB = uintptr(getProcAddr("glVertexAttrib4svARB"))
	gpVertexAttrib4svNV = uintptr(getProcAddr("glVertexAttrib4svNV"))
	gpVertexAttrib4ubNV = uintptr(getProcAddr("glVertexAttrib4ubNV"))
	gpVertexAttrib4ubv = uintptr(getProcAddr("glVertexAttrib4ubv"))
	if gpVertexAttrib4ubv == 0 {
		return errors.New("glVertexAttrib4ubv")
	}
	gpVertexAttrib4ubvARB = uintptr(getProcAddr("glVertexAttrib4ubvARB"))
	gpVertexAttrib4ubvNV = uintptr(getProcAddr("glVertexAttrib4ubvNV"))
	gpVertexAttrib4uiv = uintptr(getProcAddr("glVertexAttrib4uiv"))
	if gpVertexAttrib4uiv == 0 {
		return errors.New("glVertexAttrib4uiv")
	}
	gpVertexAttrib4uivARB = uintptr(getProcAddr("glVertexAttrib4uivARB"))
	gpVertexAttrib4usv = uintptr(getProcAddr("glVertexAttrib4usv"))
	if gpVertexAttrib4usv == 0 {
		return errors.New("glVertexAttrib4usv")
	}
	gpVertexAttrib4usvARB = uintptr(getProcAddr("glVertexAttrib4usvARB"))
	gpVertexAttribArrayObjectATI = uintptr(getProcAddr("glVertexAttribArrayObjectATI"))
	gpVertexAttribBinding = uintptr(getProcAddr("glVertexAttribBinding"))
	gpVertexAttribDivisorARB = uintptr(getProcAddr("glVertexAttribDivisorARB"))
	gpVertexAttribFormat = uintptr(getProcAddr("glVertexAttribFormat"))
	gpVertexAttribFormatNV = uintptr(getProcAddr("glVertexAttribFormatNV"))
	gpVertexAttribI1iEXT = uintptr(getProcAddr("glVertexAttribI1iEXT"))
	gpVertexAttribI1ivEXT = uintptr(getProcAddr("glVertexAttribI1ivEXT"))
	gpVertexAttribI1uiEXT = uintptr(getProcAddr("glVertexAttribI1uiEXT"))
	gpVertexAttribI1uivEXT = uintptr(getProcAddr("glVertexAttribI1uivEXT"))
	gpVertexAttribI2iEXT = uintptr(getProcAddr("glVertexAttribI2iEXT"))
	gpVertexAttribI2ivEXT = uintptr(getProcAddr("glVertexAttribI2ivEXT"))
	gpVertexAttribI2uiEXT = uintptr(getProcAddr("glVertexAttribI2uiEXT"))
	gpVertexAttribI2uivEXT = uintptr(getProcAddr("glVertexAttribI2uivEXT"))
	gpVertexAttribI3iEXT = uintptr(getProcAddr("glVertexAttribI3iEXT"))
	gpVertexAttribI3ivEXT = uintptr(getProcAddr("glVertexAttribI3ivEXT"))
	gpVertexAttribI3uiEXT = uintptr(getProcAddr("glVertexAttribI3uiEXT"))
	gpVertexAttribI3uivEXT = uintptr(getProcAddr("glVertexAttribI3uivEXT"))
	gpVertexAttribI4bvEXT = uintptr(getProcAddr("glVertexAttribI4bvEXT"))
	gpVertexAttribI4iEXT = uintptr(getProcAddr("glVertexAttribI4iEXT"))
	gpVertexAttribI4ivEXT = uintptr(getProcAddr("glVertexAttribI4ivEXT"))
	gpVertexAttribI4svEXT = uintptr(getProcAddr("glVertexAttribI4svEXT"))
	gpVertexAttribI4ubvEXT = uintptr(getProcAddr("glVertexAttribI4ubvEXT"))
	gpVertexAttribI4uiEXT = uintptr(getProcAddr("glVertexAttribI4uiEXT"))
	gpVertexAttribI4uivEXT = uintptr(getProcAddr("glVertexAttribI4uivEXT"))
	gpVertexAttribI4usvEXT = uintptr(getProcAddr("glVertexAttribI4usvEXT"))
	gpVertexAttribIFormat = uintptr(getProcAddr("glVertexAttribIFormat"))
	gpVertexAttribIFormatNV = uintptr(getProcAddr("glVertexAttribIFormatNV"))
	gpVertexAttribIPointerEXT = uintptr(getProcAddr("glVertexAttribIPointerEXT"))
	gpVertexAttribL1d = uintptr(getProcAddr("glVertexAttribL1d"))
	gpVertexAttribL1dEXT = uintptr(getProcAddr("glVertexAttribL1dEXT"))
	gpVertexAttribL1dv = uintptr(getProcAddr("glVertexAttribL1dv"))
	gpVertexAttribL1dvEXT = uintptr(getProcAddr("glVertexAttribL1dvEXT"))
	gpVertexAttribL1i64NV = uintptr(getProcAddr("glVertexAttribL1i64NV"))
	gpVertexAttribL1i64vNV = uintptr(getProcAddr("glVertexAttribL1i64vNV"))
	gpVertexAttribL1ui64ARB = uintptr(getProcAddr("glVertexAttribL1ui64ARB"))
	gpVertexAttribL1ui64NV = uintptr(getProcAddr("glVertexAttribL1ui64NV"))
	gpVertexAttribL1ui64vARB = uintptr(getProcAddr("glVertexAttribL1ui64vARB"))
	gpVertexAttribL1ui64vNV = uintptr(getProcAddr("glVertexAttribL1ui64vNV"))
	gpVertexAttribL2d = uintptr(getProcAddr("glVertexAttribL2d"))
	gpVertexAttribL2dEXT = uintptr(getProcAddr("glVertexAttribL2dEXT"))
	gpVertexAttribL2dv = uintptr(getProcAddr("glVertexAttribL2dv"))
	gpVertexAttribL2dvEXT = uintptr(getProcAddr("glVertexAttribL2dvEXT"))
	gpVertexAttribL2i64NV = uintptr(getProcAddr("glVertexAttribL2i64NV"))
	gpVertexAttribL2i64vNV = uintptr(getProcAddr("glVertexAttribL2i64vNV"))
	gpVertexAttribL2ui64NV = uintptr(getProcAddr("glVertexAttribL2ui64NV"))
	gpVertexAttribL2ui64vNV = uintptr(getProcAddr("glVertexAttribL2ui64vNV"))
	gpVertexAttribL3d = uintptr(getProcAddr("glVertexAttribL3d"))
	gpVertexAttribL3dEXT = uintptr(getProcAddr("glVertexAttribL3dEXT"))
	gpVertexAttribL3dv = uintptr(getProcAddr("glVertexAttribL3dv"))
	gpVertexAttribL3dvEXT = uintptr(getProcAddr("glVertexAttribL3dvEXT"))
	gpVertexAttribL3i64NV = uintptr(getProcAddr("glVertexAttribL3i64NV"))
	gpVertexAttribL3i64vNV = uintptr(getProcAddr("glVertexAttribL3i64vNV"))
	gpVertexAttribL3ui64NV = uintptr(getProcAddr("glVertexAttribL3ui64NV"))
	gpVertexAttribL3ui64vNV = uintptr(getProcAddr("glVertexAttribL3ui64vNV"))
	gpVertexAttribL4d = uintptr(getProcAddr("glVertexAttribL4d"))
	gpVertexAttribL4dEXT = uintptr(getProcAddr("glVertexAttribL4dEXT"))
	gpVertexAttribL4dv = uintptr(getProcAddr("glVertexAttribL4dv"))
	gpVertexAttribL4dvEXT = uintptr(getProcAddr("glVertexAttribL4dvEXT"))
	gpVertexAttribL4i64NV = uintptr(getProcAddr("glVertexAttribL4i64NV"))
	gpVertexAttribL4i64vNV = uintptr(getProcAddr("glVertexAttribL4i64vNV"))
	gpVertexAttribL4ui64NV = uintptr(getProcAddr("glVertexAttribL4ui64NV"))
	gpVertexAttribL4ui64vNV = uintptr(getProcAddr("glVertexAttribL4ui64vNV"))
	gpVertexAttribLFormat = uintptr(getProcAddr("glVertexAttribLFormat"))
	gpVertexAttribLFormatNV = uintptr(getProcAddr("glVertexAttribLFormatNV"))
	gpVertexAttribLPointer = uintptr(getProcAddr("glVertexAttribLPointer"))
	gpVertexAttribLPointerEXT = uintptr(getProcAddr("glVertexAttribLPointerEXT"))
	gpVertexAttribP1ui = uintptr(getProcAddr("glVertexAttribP1ui"))
	gpVertexAttribP1uiv = uintptr(getProcAddr("glVertexAttribP1uiv"))
	gpVertexAttribP2ui = uintptr(getProcAddr("glVertexAttribP2ui"))
	gpVertexAttribP2uiv = uintptr(getProcAddr("glVertexAttribP2uiv"))
	gpVertexAttribP3ui = uintptr(getProcAddr("glVertexAttribP3ui"))
	gpVertexAttribP3uiv = uintptr(getProcAddr("glVertexAttribP3uiv"))
	gpVertexAttribP4ui = uintptr(getProcAddr("glVertexAttribP4ui"))
	gpVertexAttribP4uiv = uintptr(getProcAddr("glVertexAttribP4uiv"))
	gpVertexAttribParameteriAMD = uintptr(getProcAddr("glVertexAttribParameteriAMD"))
	gpVertexAttribPointer = uintptr(getProcAddr("glVertexAttribPointer"))
	if gpVertexAttribPointer == 0 {
		return errors.New("glVertexAttribPointer")
	}
	gpVertexAttribPointerARB = uintptr(getProcAddr("glVertexAttribPointerARB"))
	gpVertexAttribPointerNV = uintptr(getProcAddr("glVertexAttribPointerNV"))
	gpVertexAttribs1dvNV = uintptr(getProcAddr("glVertexAttribs1dvNV"))
	gpVertexAttribs1fvNV = uintptr(getProcAddr("glVertexAttribs1fvNV"))
	gpVertexAttribs1hvNV = uintptr(getProcAddr("glVertexAttribs1hvNV"))
	gpVertexAttribs1svNV = uintptr(getProcAddr("glVertexAttribs1svNV"))
	gpVertexAttribs2dvNV = uintptr(getProcAddr("glVertexAttribs2dvNV"))
	gpVertexAttribs2fvNV = uintptr(getProcAddr("glVertexAttribs2fvNV"))
	gpVertexAttribs2hvNV = uintptr(getProcAddr("glVertexAttribs2hvNV"))
	gpVertexAttribs2svNV = uintptr(getProcAddr("glVertexAttribs2svNV"))
	gpVertexAttribs3dvNV = uintptr(getProcAddr("glVertexAttribs3dvNV"))
	gpVertexAttribs3fvNV = uintptr(getProcAddr("glVertexAttribs3fvNV"))
	gpVertexAttribs3hvNV = uintptr(getProcAddr("glVertexAttribs3hvNV"))
	gpVertexAttribs3svNV = uintptr(getProcAddr("glVertexAttribs3svNV"))
	gpVertexAttribs4dvNV = uintptr(getProcAddr("glVertexAttribs4dvNV"))
	gpVertexAttribs4fvNV = uintptr(getProcAddr("glVertexAttribs4fvNV"))
	gpVertexAttribs4hvNV = uintptr(getProcAddr("glVertexAttribs4hvNV"))
	gpVertexAttribs4svNV = uintptr(getProcAddr("glVertexAttribs4svNV"))
	gpVertexAttribs4ubvNV = uintptr(getProcAddr("glVertexAttribs4ubvNV"))
	gpVertexBindingDivisor = uintptr(getProcAddr("glVertexBindingDivisor"))
	gpVertexBlendARB = uintptr(getProcAddr("glVertexBlendARB"))
	gpVertexBlendEnvfATI = uintptr(getProcAddr("glVertexBlendEnvfATI"))
	gpVertexBlendEnviATI = uintptr(getProcAddr("glVertexBlendEnviATI"))
	gpVertexFormatNV = uintptr(getProcAddr("glVertexFormatNV"))
	gpVertexPointer = uintptr(getProcAddr("glVertexPointer"))
	if gpVertexPointer == 0 {
		return errors.New("glVertexPointer")
	}
	gpVertexPointerEXT = uintptr(getProcAddr("glVertexPointerEXT"))
	gpVertexPointerListIBM = uintptr(getProcAddr("glVertexPointerListIBM"))
	gpVertexPointervINTEL = uintptr(getProcAddr("glVertexPointervINTEL"))
	gpVertexStream1dATI = uintptr(getProcAddr("glVertexStream1dATI"))
	gpVertexStream1dvATI = uintptr(getProcAddr("glVertexStream1dvATI"))
	gpVertexStream1fATI = uintptr(getProcAddr("glVertexStream1fATI"))
	gpVertexStream1fvATI = uintptr(getProcAddr("glVertexStream1fvATI"))
	gpVertexStream1iATI = uintptr(getProcAddr("glVertexStream1iATI"))
	gpVertexStream1ivATI = uintptr(getProcAddr("glVertexStream1ivATI"))
	gpVertexStream1sATI = uintptr(getProcAddr("glVertexStream1sATI"))
	gpVertexStream1svATI = uintptr(getProcAddr("glVertexStream1svATI"))
	gpVertexStream2dATI = uintptr(getProcAddr("glVertexStream2dATI"))
	gpVertexStream2dvATI = uintptr(getProcAddr("glVertexStream2dvATI"))
	gpVertexStream2fATI = uintptr(getProcAddr("glVertexStream2fATI"))
	gpVertexStream2fvATI = uintptr(getProcAddr("glVertexStream2fvATI"))
	gpVertexStream2iATI = uintptr(getProcAddr("glVertexStream2iATI"))
	gpVertexStream2ivATI = uintptr(getProcAddr("glVertexStream2ivATI"))
	gpVertexStream2sATI = uintptr(getProcAddr("glVertexStream2sATI"))
	gpVertexStream2svATI = uintptr(getProcAddr("glVertexStream2svATI"))
	gpVertexStream3dATI = uintptr(getProcAddr("glVertexStream3dATI"))
	gpVertexStream3dvATI = uintptr(getProcAddr("glVertexStream3dvATI"))
	gpVertexStream3fATI = uintptr(getProcAddr("glVertexStream3fATI"))
	gpVertexStream3fvATI = uintptr(getProcAddr("glVertexStream3fvATI"))
	gpVertexStream3iATI = uintptr(getProcAddr("glVertexStream3iATI"))
	gpVertexStream3ivATI = uintptr(getProcAddr("glVertexStream3ivATI"))
	gpVertexStream3sATI = uintptr(getProcAddr("glVertexStream3sATI"))
	gpVertexStream3svATI = uintptr(getProcAddr("glVertexStream3svATI"))
	gpVertexStream4dATI = uintptr(getProcAddr("glVertexStream4dATI"))
	gpVertexStream4dvATI = uintptr(getProcAddr("glVertexStream4dvATI"))
	gpVertexStream4fATI = uintptr(getProcAddr("glVertexStream4fATI"))
	gpVertexStream4fvATI = uintptr(getProcAddr("glVertexStream4fvATI"))
	gpVertexStream4iATI = uintptr(getProcAddr("glVertexStream4iATI"))
	gpVertexStream4ivATI = uintptr(getProcAddr("glVertexStream4ivATI"))
	gpVertexStream4sATI = uintptr(getProcAddr("glVertexStream4sATI"))
	gpVertexStream4svATI = uintptr(getProcAddr("glVertexStream4svATI"))
	gpVertexWeightPointerEXT = uintptr(getProcAddr("glVertexWeightPointerEXT"))
	gpVertexWeightfEXT = uintptr(getProcAddr("glVertexWeightfEXT"))
	gpVertexWeightfvEXT = uintptr(getProcAddr("glVertexWeightfvEXT"))
	gpVertexWeighthNV = uintptr(getProcAddr("glVertexWeighthNV"))
	gpVertexWeighthvNV = uintptr(getProcAddr("glVertexWeighthvNV"))
	gpVideoCaptureNV = uintptr(getProcAddr("glVideoCaptureNV"))
	gpVideoCaptureStreamParameterdvNV = uintptr(getProcAddr("glVideoCaptureStreamParameterdvNV"))
	gpVideoCaptureStreamParameterfvNV = uintptr(getProcAddr("glVideoCaptureStreamParameterfvNV"))
	gpVideoCaptureStreamParameterivNV = uintptr(getProcAddr("glVideoCaptureStreamParameterivNV"))
	gpViewport = uintptr(getProcAddr("glViewport"))
	if gpViewport == 0 {
		return errors.New("glViewport")
	}
	gpViewportArrayv = uintptr(getProcAddr("glViewportArrayv"))
	gpViewportIndexedf = uintptr(getProcAddr("glViewportIndexedf"))
	gpViewportIndexedfv = uintptr(getProcAddr("glViewportIndexedfv"))
	gpViewportPositionWScaleNV = uintptr(getProcAddr("glViewportPositionWScaleNV"))
	gpViewportSwizzleNV = uintptr(getProcAddr("glViewportSwizzleNV"))
	gpWaitSemaphoreEXT = uintptr(getProcAddr("glWaitSemaphoreEXT"))
	gpWaitSync = uintptr(getProcAddr("glWaitSync"))
	gpWaitVkSemaphoreNV = uintptr(getProcAddr("glWaitVkSemaphoreNV"))
	gpWeightPathsNV = uintptr(getProcAddr("glWeightPathsNV"))
	gpWeightPointerARB = uintptr(getProcAddr("glWeightPointerARB"))
	gpWeightbvARB = uintptr(getProcAddr("glWeightbvARB"))
	gpWeightdvARB = uintptr(getProcAddr("glWeightdvARB"))
	gpWeightfvARB = uintptr(getProcAddr("glWeightfvARB"))
	gpWeightivARB = uintptr(getProcAddr("glWeightivARB"))
	gpWeightsvARB = uintptr(getProcAddr("glWeightsvARB"))
	gpWeightubvARB = uintptr(getProcAddr("glWeightubvARB"))
	gpWeightuivARB = uintptr(getProcAddr("glWeightuivARB"))
	gpWeightusvARB = uintptr(getProcAddr("glWeightusvARB"))
	gpWindowPos2d = uintptr(getProcAddr("glWindowPos2d"))
	if gpWindowPos2d == 0 {
		return errors.New("glWindowPos2d")
	}
	gpWindowPos2dARB = uintptr(getProcAddr("glWindowPos2dARB"))
	gpWindowPos2dMESA = uintptr(getProcAddr("glWindowPos2dMESA"))
	gpWindowPos2dv = uintptr(getProcAddr("glWindowPos2dv"))
	if gpWindowPos2dv == 0 {
		return errors.New("glWindowPos2dv")
	}
	gpWindowPos2dvARB = uintptr(getProcAddr("glWindowPos2dvARB"))
	gpWindowPos2dvMESA = uintptr(getProcAddr("glWindowPos2dvMESA"))
	gpWindowPos2f = uintptr(getProcAddr("glWindowPos2f"))
	if gpWindowPos2f == 0 {
		return errors.New("glWindowPos2f")
	}
	gpWindowPos2fARB = uintptr(getProcAddr("glWindowPos2fARB"))
	gpWindowPos2fMESA = uintptr(getProcAddr("glWindowPos2fMESA"))
	gpWindowPos2fv = uintptr(getProcAddr("glWindowPos2fv"))
	if gpWindowPos2fv == 0 {
		return errors.New("glWindowPos2fv")
	}
	gpWindowPos2fvARB = uintptr(getProcAddr("glWindowPos2fvARB"))
	gpWindowPos2fvMESA = uintptr(getProcAddr("glWindowPos2fvMESA"))
	gpWindowPos2i = uintptr(getProcAddr("glWindowPos2i"))
	if gpWindowPos2i == 0 {
		return errors.New("glWindowPos2i")
	}
	gpWindowPos2iARB = uintptr(getProcAddr("glWindowPos2iARB"))
	gpWindowPos2iMESA = uintptr(getProcAddr("glWindowPos2iMESA"))
	gpWindowPos2iv = uintptr(getProcAddr("glWindowPos2iv"))
	if gpWindowPos2iv == 0 {
		return errors.New("glWindowPos2iv")
	}
	gpWindowPos2ivARB = uintptr(getProcAddr("glWindowPos2ivARB"))
	gpWindowPos2ivMESA = uintptr(getProcAddr("glWindowPos2ivMESA"))
	gpWindowPos2s = uintptr(getProcAddr("glWindowPos2s"))
	if gpWindowPos2s == 0 {
		return errors.New("glWindowPos2s")
	}
	gpWindowPos2sARB = uintptr(getProcAddr("glWindowPos2sARB"))
	gpWindowPos2sMESA = uintptr(getProcAddr("glWindowPos2sMESA"))
	gpWindowPos2sv = uintptr(getProcAddr("glWindowPos2sv"))
	if gpWindowPos2sv == 0 {
		return errors.New("glWindowPos2sv")
	}
	gpWindowPos2svARB = uintptr(getProcAddr("glWindowPos2svARB"))
	gpWindowPos2svMESA = uintptr(getProcAddr("glWindowPos2svMESA"))
	gpWindowPos3d = uintptr(getProcAddr("glWindowPos3d"))
	if gpWindowPos3d == 0 {
		return errors.New("glWindowPos3d")
	}
	gpWindowPos3dARB = uintptr(getProcAddr("glWindowPos3dARB"))
	gpWindowPos3dMESA = uintptr(getProcAddr("glWindowPos3dMESA"))
	gpWindowPos3dv = uintptr(getProcAddr("glWindowPos3dv"))
	if gpWindowPos3dv == 0 {
		return errors.New("glWindowPos3dv")
	}
	gpWindowPos3dvARB = uintptr(getProcAddr("glWindowPos3dvARB"))
	gpWindowPos3dvMESA = uintptr(getProcAddr("glWindowPos3dvMESA"))
	gpWindowPos3f = uintptr(getProcAddr("glWindowPos3f"))
	if gpWindowPos3f == 0 {
		return errors.New("glWindowPos3f")
	}
	gpWindowPos3fARB = uintptr(getProcAddr("glWindowPos3fARB"))
	gpWindowPos3fMESA = uintptr(getProcAddr("glWindowPos3fMESA"))
	gpWindowPos3fv = uintptr(getProcAddr("glWindowPos3fv"))
	if gpWindowPos3fv == 0 {
		return errors.New("glWindowPos3fv")
	}
	gpWindowPos3fvARB = uintptr(getProcAddr("glWindowPos3fvARB"))
	gpWindowPos3fvMESA = uintptr(getProcAddr("glWindowPos3fvMESA"))
	gpWindowPos3i = uintptr(getProcAddr("glWindowPos3i"))
	if gpWindowPos3i == 0 {
		return errors.New("glWindowPos3i")
	}
	gpWindowPos3iARB = uintptr(getProcAddr("glWindowPos3iARB"))
	gpWindowPos3iMESA = uintptr(getProcAddr("glWindowPos3iMESA"))
	gpWindowPos3iv = uintptr(getProcAddr("glWindowPos3iv"))
	if gpWindowPos3iv == 0 {
		return errors.New("glWindowPos3iv")
	}
	gpWindowPos3ivARB = uintptr(getProcAddr("glWindowPos3ivARB"))
	gpWindowPos3ivMESA = uintptr(getProcAddr("glWindowPos3ivMESA"))
	gpWindowPos3s = uintptr(getProcAddr("glWindowPos3s"))
	if gpWindowPos3s == 0 {
		return errors.New("glWindowPos3s")
	}
	gpWindowPos3sARB = uintptr(getProcAddr("glWindowPos3sARB"))
	gpWindowPos3sMESA = uintptr(getProcAddr("glWindowPos3sMESA"))
	gpWindowPos3sv = uintptr(getProcAddr("glWindowPos3sv"))
	if gpWindowPos3sv == 0 {
		return errors.New("glWindowPos3sv")
	}
	gpWindowPos3svARB = uintptr(getProcAddr("glWindowPos3svARB"))
	gpWindowPos3svMESA = uintptr(getProcAddr("glWindowPos3svMESA"))
	gpWindowPos4dMESA = uintptr(getProcAddr("glWindowPos4dMESA"))
	gpWindowPos4dvMESA = uintptr(getProcAddr("glWindowPos4dvMESA"))
	gpWindowPos4fMESA = uintptr(getProcAddr("glWindowPos4fMESA"))
	gpWindowPos4fvMESA = uintptr(getProcAddr("glWindowPos4fvMESA"))
	gpWindowPos4iMESA = uintptr(getProcAddr("glWindowPos4iMESA"))
	gpWindowPos4ivMESA = uintptr(getProcAddr("glWindowPos4ivMESA"))
	gpWindowPos4sMESA = uintptr(getProcAddr("glWindowPos4sMESA"))
	gpWindowPos4svMESA = uintptr(getProcAddr("glWindowPos4svMESA"))
	gpWindowRectanglesEXT = uintptr(getProcAddr("glWindowRectanglesEXT"))
	gpWriteMaskEXT = uintptr(getProcAddr("glWriteMaskEXT"))
	return nil
}
