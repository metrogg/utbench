func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*APIEndpoint)(nil), (*kubeadm.APIEndpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_APIEndpoint_To_kubeadm_APIEndpoint(a.(*APIEndpoint), b.(*kubeadm.APIEndpoint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.APIEndpoint)(nil), (*APIEndpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_APIEndpoint_To_v1beta2_APIEndpoint(a.(*kubeadm.APIEndpoint), b.(*APIEndpoint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*APIServer)(nil), (*kubeadm.APIServer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_APIServer_To_kubeadm_APIServer(a.(*APIServer), b.(*kubeadm.APIServer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.APIServer)(nil), (*APIServer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_APIServer_To_v1beta2_APIServer(a.(*kubeadm.APIServer), b.(*APIServer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BootstrapToken)(nil), (*kubeadm.BootstrapToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_BootstrapToken_To_kubeadm_BootstrapToken(a.(*BootstrapToken), b.(*kubeadm.BootstrapToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.BootstrapToken)(nil), (*BootstrapToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_BootstrapToken_To_v1beta2_BootstrapToken(a.(*kubeadm.BootstrapToken), b.(*BootstrapToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BootstrapTokenDiscovery)(nil), (*kubeadm.BootstrapTokenDiscovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_BootstrapTokenDiscovery_To_kubeadm_BootstrapTokenDiscovery(a.(*BootstrapTokenDiscovery), b.(*kubeadm.BootstrapTokenDiscovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.BootstrapTokenDiscovery)(nil), (*BootstrapTokenDiscovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_BootstrapTokenDiscovery_To_v1beta2_BootstrapTokenDiscovery(a.(*kubeadm.BootstrapTokenDiscovery), b.(*BootstrapTokenDiscovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BootstrapTokenString)(nil), (*kubeadm.BootstrapTokenString)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_BootstrapTokenString_To_kubeadm_BootstrapTokenString(a.(*BootstrapTokenString), b.(*kubeadm.BootstrapTokenString), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.BootstrapTokenString)(nil), (*BootstrapTokenString)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_BootstrapTokenString_To_v1beta2_BootstrapTokenString(a.(*kubeadm.BootstrapTokenString), b.(*BootstrapTokenString), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterConfiguration)(nil), (*kubeadm.ClusterConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ClusterConfiguration_To_kubeadm_ClusterConfiguration(a.(*ClusterConfiguration), b.(*kubeadm.ClusterConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.ClusterConfiguration)(nil), (*ClusterConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_ClusterConfiguration_To_v1beta2_ClusterConfiguration(a.(*kubeadm.ClusterConfiguration), b.(*ClusterConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterStatus)(nil), (*kubeadm.ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ClusterStatus_To_kubeadm_ClusterStatus(a.(*ClusterStatus), b.(*kubeadm.ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.ClusterStatus)(nil), (*ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_ClusterStatus_To_v1beta2_ClusterStatus(a.(*kubeadm.ClusterStatus), b.(*ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ControlPlaneComponent)(nil), (*kubeadm.ControlPlaneComponent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ControlPlaneComponent_To_kubeadm_ControlPlaneComponent(a.(*ControlPlaneComponent), b.(*kubeadm.ControlPlaneComponent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.ControlPlaneComponent)(nil), (*ControlPlaneComponent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_ControlPlaneComponent_To_v1beta2_ControlPlaneComponent(a.(*kubeadm.ControlPlaneComponent), b.(*ControlPlaneComponent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DNS)(nil), (*kubeadm.DNS)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DNS_To_kubeadm_DNS(a.(*DNS), b.(*kubeadm.DNS), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.DNS)(nil), (*DNS)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_DNS_To_v1beta2_DNS(a.(*kubeadm.DNS), b.(*DNS), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Discovery)(nil), (*kubeadm.Discovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Discovery_To_kubeadm_Discovery(a.(*Discovery), b.(*kubeadm.Discovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.Discovery)(nil), (*Discovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_Discovery_To_v1beta2_Discovery(a.(*kubeadm.Discovery), b.(*Discovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Etcd)(nil), (*kubeadm.Etcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Etcd_To_kubeadm_Etcd(a.(*Etcd), b.(*kubeadm.Etcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.Etcd)(nil), (*Etcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_Etcd_To_v1beta2_Etcd(a.(*kubeadm.Etcd), b.(*Etcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExternalEtcd)(nil), (*kubeadm.ExternalEtcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ExternalEtcd_To_kubeadm_ExternalEtcd(a.(*ExternalEtcd), b.(*kubeadm.ExternalEtcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.ExternalEtcd)(nil), (*ExternalEtcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_ExternalEtcd_To_v1beta2_ExternalEtcd(a.(*kubeadm.ExternalEtcd), b.(*ExternalEtcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FileDiscovery)(nil), (*kubeadm.FileDiscovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_FileDiscovery_To_kubeadm_FileDiscovery(a.(*FileDiscovery), b.(*kubeadm.FileDiscovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.FileDiscovery)(nil), (*FileDiscovery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_FileDiscovery_To_v1beta2_FileDiscovery(a.(*kubeadm.FileDiscovery), b.(*FileDiscovery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HostPathMount)(nil), (*kubeadm.HostPathMount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_HostPathMount_To_kubeadm_HostPathMount(a.(*HostPathMount), b.(*kubeadm.HostPathMount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.HostPathMount)(nil), (*HostPathMount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_HostPathMount_To_v1beta2_HostPathMount(a.(*kubeadm.HostPathMount), b.(*HostPathMount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ImageMeta)(nil), (*kubeadm.ImageMeta)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ImageMeta_To_kubeadm_ImageMeta(a.(*ImageMeta), b.(*kubeadm.ImageMeta), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.ImageMeta)(nil), (*ImageMeta)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_ImageMeta_To_v1beta2_ImageMeta(a.(*kubeadm.ImageMeta), b.(*ImageMeta), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InitConfiguration)(nil), (*kubeadm.InitConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_InitConfiguration_To_kubeadm_InitConfiguration(a.(*InitConfiguration), b.(*kubeadm.InitConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.InitConfiguration)(nil), (*InitConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_InitConfiguration_To_v1beta2_InitConfiguration(a.(*kubeadm.InitConfiguration), b.(*InitConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JoinConfiguration)(nil), (*kubeadm.JoinConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_JoinConfiguration_To_kubeadm_JoinConfiguration(a.(*JoinConfiguration), b.(*kubeadm.JoinConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.JoinConfiguration)(nil), (*JoinConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_JoinConfiguration_To_v1beta2_JoinConfiguration(a.(*kubeadm.JoinConfiguration), b.(*JoinConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JoinControlPlane)(nil), (*kubeadm.JoinControlPlane)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_JoinControlPlane_To_kubeadm_JoinControlPlane(a.(*JoinControlPlane), b.(*kubeadm.JoinControlPlane), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.JoinControlPlane)(nil), (*JoinControlPlane)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_JoinControlPlane_To_v1beta2_JoinControlPlane(a.(*kubeadm.JoinControlPlane), b.(*JoinControlPlane), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LocalEtcd)(nil), (*kubeadm.LocalEtcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_LocalEtcd_To_kubeadm_LocalEtcd(a.(*LocalEtcd), b.(*kubeadm.LocalEtcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.LocalEtcd)(nil), (*LocalEtcd)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_LocalEtcd_To_v1beta2_LocalEtcd(a.(*kubeadm.LocalEtcd), b.(*LocalEtcd), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Networking)(nil), (*kubeadm.Networking)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Networking_To_kubeadm_Networking(a.(*Networking), b.(*kubeadm.Networking), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.Networking)(nil), (*Networking)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_Networking_To_v1beta2_Networking(a.(*kubeadm.Networking), b.(*Networking), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NodeRegistrationOptions)(nil), (*kubeadm.NodeRegistrationOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_NodeRegistrationOptions_To_kubeadm_NodeRegistrationOptions(a.(*NodeRegistrationOptions), b.(*kubeadm.NodeRegistrationOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kubeadm.NodeRegistrationOptions)(nil), (*NodeRegistrationOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kubeadm_NodeRegistrationOptions_To_v1beta2_NodeRegistrationOptions(a.(*kubeadm.NodeRegistrationOptions), b.(*NodeRegistrationOptions), scope)
	}); err != nil {
		return err
	}
	return nil
}
