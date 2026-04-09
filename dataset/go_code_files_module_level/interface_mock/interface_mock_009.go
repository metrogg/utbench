func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.Image)(nil), (*image.Image)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Image_To_image_Image(a.(*v1.Image), b.(*image.Image), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.Image)(nil), (*v1.Image)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_Image_To_v1_Image(a.(*image.Image), b.(*v1.Image), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageBlobReferences)(nil), (*image.ImageBlobReferences)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageBlobReferences_To_image_ImageBlobReferences(a.(*v1.ImageBlobReferences), b.(*image.ImageBlobReferences), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageBlobReferences)(nil), (*v1.ImageBlobReferences)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageBlobReferences_To_v1_ImageBlobReferences(a.(*image.ImageBlobReferences), b.(*v1.ImageBlobReferences), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageImportSpec)(nil), (*image.ImageImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageImportSpec_To_image_ImageImportSpec(a.(*v1.ImageImportSpec), b.(*image.ImageImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageImportSpec)(nil), (*v1.ImageImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageImportSpec_To_v1_ImageImportSpec(a.(*image.ImageImportSpec), b.(*v1.ImageImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageImportStatus)(nil), (*image.ImageImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageImportStatus_To_image_ImageImportStatus(a.(*v1.ImageImportStatus), b.(*image.ImageImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageImportStatus)(nil), (*v1.ImageImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageImportStatus_To_v1_ImageImportStatus(a.(*image.ImageImportStatus), b.(*v1.ImageImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageLayer)(nil), (*image.ImageLayer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageLayer_To_image_ImageLayer(a.(*v1.ImageLayer), b.(*image.ImageLayer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageLayer)(nil), (*v1.ImageLayer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageLayer_To_v1_ImageLayer(a.(*image.ImageLayer), b.(*v1.ImageLayer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageLayerData)(nil), (*image.ImageLayerData)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageLayerData_To_image_ImageLayerData(a.(*v1.ImageLayerData), b.(*image.ImageLayerData), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageLayerData)(nil), (*v1.ImageLayerData)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageLayerData_To_v1_ImageLayerData(a.(*image.ImageLayerData), b.(*v1.ImageLayerData), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageList)(nil), (*image.ImageList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageList_To_image_ImageList(a.(*v1.ImageList), b.(*image.ImageList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageList)(nil), (*v1.ImageList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageList_To_v1_ImageList(a.(*image.ImageList), b.(*v1.ImageList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageLookupPolicy)(nil), (*image.ImageLookupPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageLookupPolicy_To_image_ImageLookupPolicy(a.(*v1.ImageLookupPolicy), b.(*image.ImageLookupPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageLookupPolicy)(nil), (*v1.ImageLookupPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageLookupPolicy_To_v1_ImageLookupPolicy(a.(*image.ImageLookupPolicy), b.(*v1.ImageLookupPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageSignature)(nil), (*image.ImageSignature)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageSignature_To_image_ImageSignature(a.(*v1.ImageSignature), b.(*image.ImageSignature), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageSignature)(nil), (*v1.ImageSignature)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageSignature_To_v1_ImageSignature(a.(*image.ImageSignature), b.(*v1.ImageSignature), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStream)(nil), (*image.ImageStream)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStream_To_image_ImageStream(a.(*v1.ImageStream), b.(*image.ImageStream), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStream)(nil), (*v1.ImageStream)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStream_To_v1_ImageStream(a.(*image.ImageStream), b.(*v1.ImageStream), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamImage)(nil), (*image.ImageStreamImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamImage_To_image_ImageStreamImage(a.(*v1.ImageStreamImage), b.(*image.ImageStreamImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamImage)(nil), (*v1.ImageStreamImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamImage_To_v1_ImageStreamImage(a.(*image.ImageStreamImage), b.(*v1.ImageStreamImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamImport)(nil), (*image.ImageStreamImport)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamImport_To_image_ImageStreamImport(a.(*v1.ImageStreamImport), b.(*image.ImageStreamImport), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamImport)(nil), (*v1.ImageStreamImport)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamImport_To_v1_ImageStreamImport(a.(*image.ImageStreamImport), b.(*v1.ImageStreamImport), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamImportSpec)(nil), (*image.ImageStreamImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamImportSpec_To_image_ImageStreamImportSpec(a.(*v1.ImageStreamImportSpec), b.(*image.ImageStreamImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamImportSpec)(nil), (*v1.ImageStreamImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamImportSpec_To_v1_ImageStreamImportSpec(a.(*image.ImageStreamImportSpec), b.(*v1.ImageStreamImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamImportStatus)(nil), (*image.ImageStreamImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamImportStatus_To_image_ImageStreamImportStatus(a.(*v1.ImageStreamImportStatus), b.(*image.ImageStreamImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamImportStatus)(nil), (*v1.ImageStreamImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamImportStatus_To_v1_ImageStreamImportStatus(a.(*image.ImageStreamImportStatus), b.(*v1.ImageStreamImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamLayers)(nil), (*image.ImageStreamLayers)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamLayers_To_image_ImageStreamLayers(a.(*v1.ImageStreamLayers), b.(*image.ImageStreamLayers), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamLayers)(nil), (*v1.ImageStreamLayers)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamLayers_To_v1_ImageStreamLayers(a.(*image.ImageStreamLayers), b.(*v1.ImageStreamLayers), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamList)(nil), (*image.ImageStreamList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamList_To_image_ImageStreamList(a.(*v1.ImageStreamList), b.(*image.ImageStreamList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamList)(nil), (*v1.ImageStreamList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamList_To_v1_ImageStreamList(a.(*image.ImageStreamList), b.(*v1.ImageStreamList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamMapping)(nil), (*image.ImageStreamMapping)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamMapping_To_image_ImageStreamMapping(a.(*v1.ImageStreamMapping), b.(*image.ImageStreamMapping), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamMapping)(nil), (*v1.ImageStreamMapping)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamMapping_To_v1_ImageStreamMapping(a.(*image.ImageStreamMapping), b.(*v1.ImageStreamMapping), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamSpec)(nil), (*image.ImageStreamSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamSpec_To_image_ImageStreamSpec(a.(*v1.ImageStreamSpec), b.(*image.ImageStreamSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamSpec)(nil), (*v1.ImageStreamSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamSpec_To_v1_ImageStreamSpec(a.(*image.ImageStreamSpec), b.(*v1.ImageStreamSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamStatus)(nil), (*image.ImageStreamStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamStatus_To_image_ImageStreamStatus(a.(*v1.ImageStreamStatus), b.(*image.ImageStreamStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamStatus)(nil), (*v1.ImageStreamStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamStatus_To_v1_ImageStreamStatus(a.(*image.ImageStreamStatus), b.(*v1.ImageStreamStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamTag)(nil), (*image.ImageStreamTag)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamTag_To_image_ImageStreamTag(a.(*v1.ImageStreamTag), b.(*image.ImageStreamTag), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamTag)(nil), (*v1.ImageStreamTag)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamTag_To_v1_ImageStreamTag(a.(*image.ImageStreamTag), b.(*v1.ImageStreamTag), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageStreamTagList)(nil), (*image.ImageStreamTagList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamTagList_To_image_ImageStreamTagList(a.(*v1.ImageStreamTagList), b.(*image.ImageStreamTagList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.ImageStreamTagList)(nil), (*v1.ImageStreamTagList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamTagList_To_v1_ImageStreamTagList(a.(*image.ImageStreamTagList), b.(*v1.ImageStreamTagList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RepositoryImportSpec)(nil), (*image.RepositoryImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RepositoryImportSpec_To_image_RepositoryImportSpec(a.(*v1.RepositoryImportSpec), b.(*image.RepositoryImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.RepositoryImportSpec)(nil), (*v1.RepositoryImportSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_RepositoryImportSpec_To_v1_RepositoryImportSpec(a.(*image.RepositoryImportSpec), b.(*v1.RepositoryImportSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RepositoryImportStatus)(nil), (*image.RepositoryImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RepositoryImportStatus_To_image_RepositoryImportStatus(a.(*v1.RepositoryImportStatus), b.(*image.RepositoryImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.RepositoryImportStatus)(nil), (*v1.RepositoryImportStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_RepositoryImportStatus_To_v1_RepositoryImportStatus(a.(*image.RepositoryImportStatus), b.(*v1.RepositoryImportStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SignatureCondition)(nil), (*image.SignatureCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SignatureCondition_To_image_SignatureCondition(a.(*v1.SignatureCondition), b.(*image.SignatureCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.SignatureCondition)(nil), (*v1.SignatureCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_SignatureCondition_To_v1_SignatureCondition(a.(*image.SignatureCondition), b.(*v1.SignatureCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SignatureGenericEntity)(nil), (*image.SignatureGenericEntity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SignatureGenericEntity_To_image_SignatureGenericEntity(a.(*v1.SignatureGenericEntity), b.(*image.SignatureGenericEntity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.SignatureGenericEntity)(nil), (*v1.SignatureGenericEntity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_SignatureGenericEntity_To_v1_SignatureGenericEntity(a.(*image.SignatureGenericEntity), b.(*v1.SignatureGenericEntity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SignatureIssuer)(nil), (*image.SignatureIssuer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SignatureIssuer_To_image_SignatureIssuer(a.(*v1.SignatureIssuer), b.(*image.SignatureIssuer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.SignatureIssuer)(nil), (*v1.SignatureIssuer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_SignatureIssuer_To_v1_SignatureIssuer(a.(*image.SignatureIssuer), b.(*v1.SignatureIssuer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SignatureSubject)(nil), (*image.SignatureSubject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SignatureSubject_To_image_SignatureSubject(a.(*v1.SignatureSubject), b.(*image.SignatureSubject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.SignatureSubject)(nil), (*v1.SignatureSubject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_SignatureSubject_To_v1_SignatureSubject(a.(*image.SignatureSubject), b.(*v1.SignatureSubject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagEvent)(nil), (*image.TagEvent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagEvent_To_image_TagEvent(a.(*v1.TagEvent), b.(*image.TagEvent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.TagEvent)(nil), (*v1.TagEvent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_TagEvent_To_v1_TagEvent(a.(*image.TagEvent), b.(*v1.TagEvent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagEventCondition)(nil), (*image.TagEventCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagEventCondition_To_image_TagEventCondition(a.(*v1.TagEventCondition), b.(*image.TagEventCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.TagEventCondition)(nil), (*v1.TagEventCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_TagEventCondition_To_v1_TagEventCondition(a.(*image.TagEventCondition), b.(*v1.TagEventCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagImportPolicy)(nil), (*image.TagImportPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagImportPolicy_To_image_TagImportPolicy(a.(*v1.TagImportPolicy), b.(*image.TagImportPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.TagImportPolicy)(nil), (*v1.TagImportPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_TagImportPolicy_To_v1_TagImportPolicy(a.(*image.TagImportPolicy), b.(*v1.TagImportPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagReference)(nil), (*image.TagReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagReference_To_image_TagReference(a.(*v1.TagReference), b.(*image.TagReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.TagReference)(nil), (*v1.TagReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_TagReference_To_v1_TagReference(a.(*image.TagReference), b.(*v1.TagReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagReferencePolicy)(nil), (*image.TagReferencePolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagReferencePolicy_To_image_TagReferencePolicy(a.(*v1.TagReferencePolicy), b.(*image.TagReferencePolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*image.TagReferencePolicy)(nil), (*v1.TagReferencePolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_TagReferencePolicy_To_v1_TagReferencePolicy(a.(*image.TagReferencePolicy), b.(*v1.TagReferencePolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*image.ImageStreamSpec)(nil), (*v1.ImageStreamSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamSpec_To_v1_ImageStreamSpec(a.(*image.ImageStreamSpec), b.(*v1.ImageStreamSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*image.ImageStreamStatus)(nil), (*v1.ImageStreamStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_ImageStreamStatus_To_v1_ImageStreamStatus(a.(*image.ImageStreamStatus), b.(*v1.ImageStreamStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*image.Image)(nil), (*v1.Image)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_image_Image_To_v1_Image(a.(*image.Image), b.(*v1.Image), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ImageStreamSpec)(nil), (*image.ImageStreamSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamSpec_To_image_ImageStreamSpec(a.(*v1.ImageStreamSpec), b.(*image.ImageStreamSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ImageStreamStatus)(nil), (*image.ImageStreamStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageStreamStatus_To_image_ImageStreamStatus(a.(*v1.ImageStreamStatus), b.(*image.ImageStreamStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.Image)(nil), (*image.Image)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Image_To_image_Image(a.(*v1.Image), b.(*image.Image), scope)
	}); err != nil {
		return err
	}
	return nil
}
