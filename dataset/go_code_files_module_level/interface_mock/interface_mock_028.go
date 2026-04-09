func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSIDriver)(nil), (*storage.CSIDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSIDriver_To_storage_CSIDriver(a.(*v1beta1.CSIDriver), b.(*storage.CSIDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSIDriver)(nil), (*v1beta1.CSIDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSIDriver_To_v1beta1_CSIDriver(a.(*storage.CSIDriver), b.(*v1beta1.CSIDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSIDriverList)(nil), (*storage.CSIDriverList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSIDriverList_To_storage_CSIDriverList(a.(*v1beta1.CSIDriverList), b.(*storage.CSIDriverList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSIDriverList)(nil), (*v1beta1.CSIDriverList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSIDriverList_To_v1beta1_CSIDriverList(a.(*storage.CSIDriverList), b.(*v1beta1.CSIDriverList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSIDriverSpec)(nil), (*storage.CSIDriverSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSIDriverSpec_To_storage_CSIDriverSpec(a.(*v1beta1.CSIDriverSpec), b.(*storage.CSIDriverSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSIDriverSpec)(nil), (*v1beta1.CSIDriverSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSIDriverSpec_To_v1beta1_CSIDriverSpec(a.(*storage.CSIDriverSpec), b.(*v1beta1.CSIDriverSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSINode)(nil), (*storage.CSINode)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSINode_To_storage_CSINode(a.(*v1beta1.CSINode), b.(*storage.CSINode), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSINode)(nil), (*v1beta1.CSINode)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSINode_To_v1beta1_CSINode(a.(*storage.CSINode), b.(*v1beta1.CSINode), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSINodeDriver)(nil), (*storage.CSINodeDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSINodeDriver_To_storage_CSINodeDriver(a.(*v1beta1.CSINodeDriver), b.(*storage.CSINodeDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSINodeDriver)(nil), (*v1beta1.CSINodeDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSINodeDriver_To_v1beta1_CSINodeDriver(a.(*storage.CSINodeDriver), b.(*v1beta1.CSINodeDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSINodeList)(nil), (*storage.CSINodeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSINodeList_To_storage_CSINodeList(a.(*v1beta1.CSINodeList), b.(*storage.CSINodeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSINodeList)(nil), (*v1beta1.CSINodeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSINodeList_To_v1beta1_CSINodeList(a.(*storage.CSINodeList), b.(*v1beta1.CSINodeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.CSINodeSpec)(nil), (*storage.CSINodeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CSINodeSpec_To_storage_CSINodeSpec(a.(*v1beta1.CSINodeSpec), b.(*storage.CSINodeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.CSINodeSpec)(nil), (*v1beta1.CSINodeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_CSINodeSpec_To_v1beta1_CSINodeSpec(a.(*storage.CSINodeSpec), b.(*v1beta1.CSINodeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StorageClass)(nil), (*storage.StorageClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StorageClass_To_storage_StorageClass(a.(*v1beta1.StorageClass), b.(*storage.StorageClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.StorageClass)(nil), (*v1beta1.StorageClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_StorageClass_To_v1beta1_StorageClass(a.(*storage.StorageClass), b.(*v1beta1.StorageClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StorageClassList)(nil), (*storage.StorageClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StorageClassList_To_storage_StorageClassList(a.(*v1beta1.StorageClassList), b.(*storage.StorageClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.StorageClassList)(nil), (*v1beta1.StorageClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_StorageClassList_To_v1beta1_StorageClassList(a.(*storage.StorageClassList), b.(*v1beta1.StorageClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeAttachment)(nil), (*storage.VolumeAttachment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeAttachment_To_storage_VolumeAttachment(a.(*v1beta1.VolumeAttachment), b.(*storage.VolumeAttachment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeAttachment)(nil), (*v1beta1.VolumeAttachment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeAttachment_To_v1beta1_VolumeAttachment(a.(*storage.VolumeAttachment), b.(*v1beta1.VolumeAttachment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeAttachmentList)(nil), (*storage.VolumeAttachmentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeAttachmentList_To_storage_VolumeAttachmentList(a.(*v1beta1.VolumeAttachmentList), b.(*storage.VolumeAttachmentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeAttachmentList)(nil), (*v1beta1.VolumeAttachmentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeAttachmentList_To_v1beta1_VolumeAttachmentList(a.(*storage.VolumeAttachmentList), b.(*v1beta1.VolumeAttachmentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeAttachmentSource)(nil), (*storage.VolumeAttachmentSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeAttachmentSource_To_storage_VolumeAttachmentSource(a.(*v1beta1.VolumeAttachmentSource), b.(*storage.VolumeAttachmentSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeAttachmentSource)(nil), (*v1beta1.VolumeAttachmentSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeAttachmentSource_To_v1beta1_VolumeAttachmentSource(a.(*storage.VolumeAttachmentSource), b.(*v1beta1.VolumeAttachmentSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeAttachmentSpec)(nil), (*storage.VolumeAttachmentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeAttachmentSpec_To_storage_VolumeAttachmentSpec(a.(*v1beta1.VolumeAttachmentSpec), b.(*storage.VolumeAttachmentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeAttachmentSpec)(nil), (*v1beta1.VolumeAttachmentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeAttachmentSpec_To_v1beta1_VolumeAttachmentSpec(a.(*storage.VolumeAttachmentSpec), b.(*v1beta1.VolumeAttachmentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeAttachmentStatus)(nil), (*storage.VolumeAttachmentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeAttachmentStatus_To_storage_VolumeAttachmentStatus(a.(*v1beta1.VolumeAttachmentStatus), b.(*storage.VolumeAttachmentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeAttachmentStatus)(nil), (*v1beta1.VolumeAttachmentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeAttachmentStatus_To_v1beta1_VolumeAttachmentStatus(a.(*storage.VolumeAttachmentStatus), b.(*v1beta1.VolumeAttachmentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.VolumeError)(nil), (*storage.VolumeError)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_VolumeError_To_storage_VolumeError(a.(*v1beta1.VolumeError), b.(*storage.VolumeError), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*storage.VolumeError)(nil), (*v1beta1.VolumeError)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_storage_VolumeError_To_v1beta1_VolumeError(a.(*storage.VolumeError), b.(*v1beta1.VolumeError), scope)
	}); err != nil {
		return err
	}
	return nil
}
