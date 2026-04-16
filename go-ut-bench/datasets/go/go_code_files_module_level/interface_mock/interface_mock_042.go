func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.ClusterNetwork)(nil), (*network.ClusterNetwork)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterNetwork_To_network_ClusterNetwork(a.(*v1.ClusterNetwork), b.(*network.ClusterNetwork), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.ClusterNetwork)(nil), (*v1.ClusterNetwork)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_ClusterNetwork_To_v1_ClusterNetwork(a.(*network.ClusterNetwork), b.(*v1.ClusterNetwork), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterNetworkEntry)(nil), (*network.ClusterNetworkEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterNetworkEntry_To_network_ClusterNetworkEntry(a.(*v1.ClusterNetworkEntry), b.(*network.ClusterNetworkEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.ClusterNetworkEntry)(nil), (*v1.ClusterNetworkEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_ClusterNetworkEntry_To_v1_ClusterNetworkEntry(a.(*network.ClusterNetworkEntry), b.(*v1.ClusterNetworkEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterNetworkList)(nil), (*network.ClusterNetworkList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterNetworkList_To_network_ClusterNetworkList(a.(*v1.ClusterNetworkList), b.(*network.ClusterNetworkList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.ClusterNetworkList)(nil), (*v1.ClusterNetworkList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_ClusterNetworkList_To_v1_ClusterNetworkList(a.(*network.ClusterNetworkList), b.(*v1.ClusterNetworkList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EgressNetworkPolicy)(nil), (*network.EgressNetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EgressNetworkPolicy_To_network_EgressNetworkPolicy(a.(*v1.EgressNetworkPolicy), b.(*network.EgressNetworkPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.EgressNetworkPolicy)(nil), (*v1.EgressNetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_EgressNetworkPolicy_To_v1_EgressNetworkPolicy(a.(*network.EgressNetworkPolicy), b.(*v1.EgressNetworkPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EgressNetworkPolicyList)(nil), (*network.EgressNetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EgressNetworkPolicyList_To_network_EgressNetworkPolicyList(a.(*v1.EgressNetworkPolicyList), b.(*network.EgressNetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.EgressNetworkPolicyList)(nil), (*v1.EgressNetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_EgressNetworkPolicyList_To_v1_EgressNetworkPolicyList(a.(*network.EgressNetworkPolicyList), b.(*v1.EgressNetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EgressNetworkPolicyPeer)(nil), (*network.EgressNetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EgressNetworkPolicyPeer_To_network_EgressNetworkPolicyPeer(a.(*v1.EgressNetworkPolicyPeer), b.(*network.EgressNetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.EgressNetworkPolicyPeer)(nil), (*v1.EgressNetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_EgressNetworkPolicyPeer_To_v1_EgressNetworkPolicyPeer(a.(*network.EgressNetworkPolicyPeer), b.(*v1.EgressNetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EgressNetworkPolicyRule)(nil), (*network.EgressNetworkPolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EgressNetworkPolicyRule_To_network_EgressNetworkPolicyRule(a.(*v1.EgressNetworkPolicyRule), b.(*network.EgressNetworkPolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.EgressNetworkPolicyRule)(nil), (*v1.EgressNetworkPolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_EgressNetworkPolicyRule_To_v1_EgressNetworkPolicyRule(a.(*network.EgressNetworkPolicyRule), b.(*v1.EgressNetworkPolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EgressNetworkPolicySpec)(nil), (*network.EgressNetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EgressNetworkPolicySpec_To_network_EgressNetworkPolicySpec(a.(*v1.EgressNetworkPolicySpec), b.(*network.EgressNetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.EgressNetworkPolicySpec)(nil), (*v1.EgressNetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_EgressNetworkPolicySpec_To_v1_EgressNetworkPolicySpec(a.(*network.EgressNetworkPolicySpec), b.(*v1.EgressNetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HostSubnet)(nil), (*network.HostSubnet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HostSubnet_To_network_HostSubnet(a.(*v1.HostSubnet), b.(*network.HostSubnet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.HostSubnet)(nil), (*v1.HostSubnet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_HostSubnet_To_v1_HostSubnet(a.(*network.HostSubnet), b.(*v1.HostSubnet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HostSubnetList)(nil), (*network.HostSubnetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HostSubnetList_To_network_HostSubnetList(a.(*v1.HostSubnetList), b.(*network.HostSubnetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.HostSubnetList)(nil), (*v1.HostSubnetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_HostSubnetList_To_v1_HostSubnetList(a.(*network.HostSubnetList), b.(*v1.HostSubnetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NetNamespace)(nil), (*network.NetNamespace)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NetNamespace_To_network_NetNamespace(a.(*v1.NetNamespace), b.(*network.NetNamespace), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.NetNamespace)(nil), (*v1.NetNamespace)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_NetNamespace_To_v1_NetNamespace(a.(*network.NetNamespace), b.(*v1.NetNamespace), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NetNamespaceList)(nil), (*network.NetNamespaceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NetNamespaceList_To_network_NetNamespaceList(a.(*v1.NetNamespaceList), b.(*network.NetNamespaceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*network.NetNamespaceList)(nil), (*v1.NetNamespaceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_network_NetNamespaceList_To_v1_NetNamespaceList(a.(*network.NetNamespaceList), b.(*v1.NetNamespaceList), scope)
	}); err != nil {
		return err
	}
	return nil
}
