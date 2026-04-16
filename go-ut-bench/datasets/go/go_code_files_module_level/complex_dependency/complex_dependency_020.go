func newClient(apiurl string, apikey string, secret string, async bool, verifyssl bool) *CloudStackClient {
	jar, _ := cookiejar.New(nil)
	cs := &CloudStackClient{
		client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: !verifyssl},
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: time.Duration(60 * time.Second),
		},
		baseURL: apiurl,
		apiKey:  apikey,
		secret:  secret,
		async:   async,
		options: []OptionFunc{},
		timeout: 300,
	}
	cs.APIDiscovery = NewAPIDiscoveryService(cs)
	cs.Account = NewAccountService(cs)
	cs.Address = NewAddressService(cs)
	cs.AffinityGroup = NewAffinityGroupService(cs)
	cs.Alert = NewAlertService(cs)
	cs.Asyncjob = NewAsyncjobService(cs)
	cs.Authentication = NewAuthenticationService(cs)
	cs.AutoScale = NewAutoScaleService(cs)
	cs.Baremetal = NewBaremetalService(cs)
	cs.BigSwitchBCF = NewBigSwitchBCFService(cs)
	cs.BrocadeVCS = NewBrocadeVCSService(cs)
	cs.Certificate = NewCertificateService(cs)
	cs.CloudIdentifier = NewCloudIdentifierService(cs)
	cs.Cluster = NewClusterService(cs)
	cs.Configuration = NewConfigurationService(cs)
	cs.Custom = NewCustomService(cs)
	cs.DiskOffering = NewDiskOfferingService(cs)
	cs.Domain = NewDomainService(cs)
	cs.Event = NewEventService(cs)
	cs.Firewall = NewFirewallService(cs)
	cs.GuestOS = NewGuestOSService(cs)
	cs.Host = NewHostService(cs)
	cs.Hypervisor = NewHypervisorService(cs)
	cs.ISO = NewISOService(cs)
	cs.ImageStore = NewImageStoreService(cs)
	cs.InternalLB = NewInternalLBService(cs)
	cs.LDAP = NewLDAPService(cs)
	cs.Limit = NewLimitService(cs)
	cs.LoadBalancer = NewLoadBalancerService(cs)
	cs.NAT = NewNATService(cs)
	cs.NetworkACL = NewNetworkACLService(cs)
	cs.NetworkDevice = NewNetworkDeviceService(cs)
	cs.NetworkOffering = NewNetworkOfferingService(cs)
	cs.Network = NewNetworkService(cs)
	cs.Nic = NewNicService(cs)
	cs.NiciraNVP = NewNiciraNVPService(cs)
	cs.NuageVSP = NewNuageVSPService(cs)
	cs.OutofbandManagement = NewOutofbandManagementService(cs)
	cs.OvsElement = NewOvsElementService(cs)
	cs.Pod = NewPodService(cs)
	cs.Pool = NewPoolService(cs)
	cs.PortableIP = NewPortableIPService(cs)
	cs.Project = NewProjectService(cs)
	cs.Quota = NewQuotaService(cs)
	cs.Region = NewRegionService(cs)
	cs.Resourcemetadata = NewResourcemetadataService(cs)
	cs.Resourcetags = NewResourcetagsService(cs)
	cs.Role = NewRoleService(cs)
	cs.Router = NewRouterService(cs)
	cs.SSH = NewSSHService(cs)
	cs.SecurityGroup = NewSecurityGroupService(cs)
	cs.ServiceOffering = NewServiceOfferingService(cs)
	cs.Snapshot = NewSnapshotService(cs)
	cs.StoragePool = NewStoragePoolService(cs)
	cs.StratosphereSSP = NewStratosphereSSPService(cs)
	cs.Swift = NewSwiftService(cs)
	cs.SystemCapacity = NewSystemCapacityService(cs)
	cs.SystemVM = NewSystemVMService(cs)
	cs.Template = NewTemplateService(cs)
	cs.UCS = NewUCSService(cs)
	cs.Usage = NewUsageService(cs)
	cs.User = NewUserService(cs)
	cs.VLAN = NewVLANService(cs)
	cs.VMGroup = NewVMGroupService(cs)
	cs.VPC = NewVPCService(cs)
	cs.VPN = NewVPNService(cs)
	cs.VirtualMachine = NewVirtualMachineService(cs)
	cs.Volume = NewVolumeService(cs)
	cs.Zone = NewZoneService(cs)
	return cs
}
