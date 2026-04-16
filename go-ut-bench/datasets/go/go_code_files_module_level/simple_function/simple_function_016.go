func (c *container) createArgs(cmds []string) []string {
	adHoc := (len(cmds) > 0)
	args := []string{}
	// AddHost
	for _, addHost := range c.AddHost() {
		args = append(args, "--add-host", addHost)
	}
	// BlkioWeight
	if c.BlkioWeight > 0 {
		args = append(args, "--blkio-weight", strconv.Itoa(c.BlkioWeight))
	}
	// BlkioWeightDevice
	for _, blkioWeightDevice := range c.BlkioWeightDevice() {
		args = append(args, "--blkio-weight-device", blkioWeightDevice)
	}
	// CapAdd
	for _, capAdd := range c.CapAdd() {
		args = append(args, "--cap-add", capAdd)
	}
	// CapDrop
	for _, capDrop := range c.CapDrop() {
		args = append(args, "--cap-drop", capDrop)
	}
	// CgroupParent
	if len(c.CgroupParent()) > 0 {
		args = append(args, "--cgroup-parent", c.CgroupParent())
	}
	// Cidfile
	if len(c.Cidfile()) > 0 {
		args = append(args, "--cidfile", c.Cidfile())
	}
	// CPUPeriod
	if c.CPUPeriod > 0 {
		args = append(args, "--cpu-period", strconv.Itoa(c.CPUPeriod))
	}
	// CPUQuota
	if c.CPUQuota > 0 {
		args = append(args, "--cpu-quota", strconv.Itoa(c.CPUQuota))
	}
	// CPU set
	if c.CPUset > 0 {
		args = append(args, "--cpuset", strconv.Itoa(c.CPUset))
	}
	// CPU shares
	if c.CPUShares > 0 {
		args = append(args, "--cpu-shares", strconv.Itoa(c.CPUShares))
	}
	// Device
	for _, device := range c.Device() {
		args = append(args, "--device", device)
	}
	// DeviceReadBps
	for _, deviceReadBps := range c.DeviceReadBps() {
		args = append(args, "--device-read-bps", deviceReadBps)
	}
	// DeviceReadIops
	for _, deviceReadIops := range c.DeviceReadIops() {
		args = append(args, "--device-read-iops", deviceReadIops)
	}
	// DeviceWriteBps
	for _, deviceWriteBps := range c.DeviceWriteBps() {
		args = append(args, "--device-write-bps", deviceWriteBps)
	}
	// DeviceWriteIops
	for _, deviceWriteIops := range c.DeviceWriteIops() {
		args = append(args, "--device-write-iops", deviceWriteIops)
	}
	// DNS
	for _, dns := range c.DNS() {
		args = append(args, "--dns", dns)
	}

	// DNSOpt
	for _, dnsOpt := range c.DNSOpt() {
		args = append(args, "--dns-opt", dnsOpt)
	}
	// DNS Search
	for _, dnsSearch := range c.DNSSearch() {
		args = append(args, "--dns-search", dnsSearch)
	}
	// Entrypoint
	if len(c.Entrypoint()) > 0 {
		args = append(args, "--entrypoint", c.Entrypoint())
	}
	// Env
	for _, env := range c.Env() {
		args = append(args, "--env", env)
	}
	// Env file
	for _, envFile := range c.EnvFile() {
		args = append(args, "--env-file", envFile)
	}
	// Expose
	for _, expose := range c.Expose() {
		args = append(args, "--expose", expose)
	}
	// GroupAdd
	for _, groupAdd := range c.GroupAdd() {
		args = append(args, "--group-add", groupAdd)
	}
	// Health Cmd
	if len(c.HealthCmd()) > 0 {
		args = append(args, "--health-cmd", c.HealthCmd())
	}
	// Health Interval
	if len(c.HealthInterval()) > 0 {
		args = append(args, "--health-interval", c.HealthInterval())
	}
	// Health Retries
	if c.HealthRetries > 0 {
		args = append(args, "--health-retries", strconv.Itoa(c.HealthRetries))
	} else if c.HealthcheckParams().Retries > 0 {
		args = append(args, "--health-retries", strconv.Itoa(c.HealthcheckParams().Retries))
	}
	// Health Timeout
	if len(c.HealthTimeout()) > 0 {
		args = append(args, "--health-timeout", c.HealthTimeout())
	}
	// Host
	if len(c.Hostname()) > 0 {
		args = append(args, "--hostname", c.Hostname())
	}
	// Init
	if c.Init {
		args = append(args, "--init")
	}
	// Interactive
	if c.Stdin_Open || c.Interactive {
		args = append(args, "--interactive")
	}
	// Ip
	if !adHoc {
		if len(c.Ip()) > 0 {
			args = append(args, "--ip", c.Ip())
		}
	}
	// Ip6
	if !adHoc {
		if len(c.Ip6()) > 0 {
			args = append(args, "--ip6", c.Ip6())
		}
	}
	// IPC
	if len(c.IPC()) > 0 {
		ipcContainer := containerReference(c.IPC())
		if len(ipcContainer) > 0 {
			if includes(allowed, ipcContainer) {
				args = append(args, "--ipc", "container:"+cfg.Container(ipcContainer).ActualName(false))
			}
		} else {
			args = append(args, "--ipc", c.IPC())
		}
	}
	// Isolation
	if len(c.Isolation()) > 0 {
		args = append(args, "--isolation", c.Isolation())
	}
	// KernelMemory
	if len(c.KernelMemory()) > 0 {
		args = append(args, "--kernel-memory", c.KernelMemory())
	}
	// Label
	for _, label := range c.Label() {
		args = append(args, "--label", label)
	}
	// LabelFile
	for _, labelFile := range c.LabelFile() {
		args = append(args, "--label-file", labelFile)
	}
	// Link
	for _, link := range c.Link() {
		linkParts := strings.Split(link, ":")
		linkName := linkParts[0]
		if includes(allowed, linkName) {
			linkParts[0] = cfg.Container(linkName).ActualName(false)
			args = append(args, "--link", strings.Join(linkParts, ":"))
		}
	}
	// External Links
	for _, externalLink := range c.ExternalLinks() {
		args = append(args, "--link", externalLink)
	}
	// LogDriver
	if len(c.LogDriver()) > 0 {
		args = append(args, "--log-driver", c.LogDriver())
	}
	// LogOpt
	for _, opt := range c.LogOpt() {
		args = append(args, "--log-opt", opt)
	}
	// LxcConf
	for _, lxcConf := range c.LxcConf() {
		args = append(args, "--lxc-conf", lxcConf)
	}
	// Mac address
	if len(c.MacAddress()) > 0 {
		args = append(args, "--mac-address", c.MacAddress())
	}
	// Memory
	if len(c.Memory()) > 0 {
		args = append(args, "--memory", c.Memory())
	}
	// MemoryReservation
	if len(c.MemoryReservation()) > 0 {
		args = append(args, "--memory-reservation", c.MemoryReservation())
	}
	// MemorySwap
	if len(c.MemorySwap()) > 0 {
		args = append(args, "--memory-swap", c.MemorySwap())
	}
	// MemorySwappiness
	if c.MemorySwappiness.Defined {
		args = append(args, "--memory-swappiness", strconv.Itoa(c.MemorySwappiness.Value))
	}
	// Net
	netParam := c.ActualNet()
	if len(netParam) > 0 && netParam != netBridge {
		args = append(args, "--net", netParam)
	}
	// NetAlias
	for _, netAlias := range c.NetAlias() {
		args = append(args, "--net-alias", netAlias)
	}
	// NoHealthcheck
	if c.NoHealthcheck || c.HealthcheckParams().Disable {
		args = append(args, "--no-healthcheck")
	}
	// OomKillDisable
	if c.OomKillDisable {
		args = append(args, "--oom-kill-disable")
	}
	// OomScoreAdj
	if len(c.OomScoreAdj()) > 0 {
		args = append(args, "--oom-score-adj", c.OomScoreAdj())
	}
	// PID
	if len(c.Pid()) > 0 {
		args = append(args, "--pid", c.Pid())
	}
	// Privileged
	if c.Privileged {
		args = append(args, "--privileged")
	}
	// Publish
	if !adHoc {
		for _, port := range c.Publish() {
			args = append(args, "--publish", port)
		}
	}
	// PublishAll
	if !adHoc {
		if c.PublishAll {
			args = append(args, "--publish-all")
		}
	}
	// ReadOnly
	if c.ReadOnly || c.Read_Only {
		args = append(args, "--read-only")
	}
	// Restart
	if len(c.Restart()) > 0 {
		args = append(args, "--restart", c.Restart())
	}
	// Rm
	if adHoc || c.RawRm {
		args = append(args, "--rm")
	}
	// SecurityOpt
	for _, securityOpt := range c.SecurityOpt() {
		args = append(args, "--security-opt", securityOpt)
	}
	// Share SSH socket
	if c.ShareSshSocket {
		sock_path := os.Getenv("SSH_AUTH_SOCK")
		args = append(args, "--volume", sock_path+":/ssh-socket")
		args = append(args, "--env", "SSH_AUTH_SOCK=/ssh-socket")
	}
	// ShmSize
	if len(c.ShmSize()) > 0 {
		args = append(args, "--shm-size", c.ShmSize())
	}
	// SigProxy
	if c.SigProxy.Falsy() {
		args = append(args, "--sig-proxy=false")
	}
	// StopSignal
	if len(c.StopSignal()) > 0 {
		args = append(args, "--stop-signal", c.StopSignal())
	}
	// StopTimeout
	if len(c.StopTimeout()) > 0 {
		args = append(args, "--stop-timeout", c.StopTimeout())
	}
	// Tmpfs
	for _, tmpfs := range c.Tmpfs() {
		args = append(args, "--tmpfs", tmpfs)
	}
	// Tty
	if c.Tty {
		args = append(args, "--tty")
	}
	// Ulimit
	for _, ulimit := range c.Ulimit() {
		args = append(args, "--ulimit", ulimit)
	}
	// User
	if len(c.User()) > 0 {
		args = append(args, "--user", c.User())
	}
	// Userns
	if len(c.Userns()) > 0 {
		args = append(args, "--userns", c.Userns())
	}
	// Uts
	if len(c.Uts()) > 0 {
		args = append(args, "--uts", c.Uts())
	}
	// Volumes
	for _, volume := range c.Volume() {
		volumeArgs := []string{"--volume"}
		am := cfg.AcceleratedMount(volume)
		if accelerationEnabled() && am != nil {
			am.Run()
			volumeArgs = append(volumeArgs, am.VolumeArg())
		} else {
			volumeArgs = append(volumeArgs, actualVolumeArg(volume))
		}
		args = append(args, volumeArgs...)
	}
	// VolumeDriver
	if len(c.VolumeDriver()) > 0 {
		args = append(args, "--volume-driver", c.VolumeDriver())
	}
	// VolumesFrom
	for _, volumesFrom := range c.VolumesFrom() {
		volumesFromParts := strings.Split(volumesFrom, ":")
		volumesFromName := volumesFromParts[0]
		if includes(allowed, volumesFromName) {
			volumesFromParts[0] = cfg.Container(volumesFromName).ActualName(false)
			args = append(args, "--volumes-from", strings.Join(volumesFromParts, ":"))
		}
	}
	// Workdir
	if len(c.Workdir()) > 0 {
		args = append(args, "--workdir", c.Workdir())
	}
	// Name
	args = append(args, "--name", c.ActualName(adHoc))
	// Image
	args = append(args, c.Image())
	// Command
	if len(cmds) > 0 {
		args = append(args, cmds...)
	} else {
		args = append(args, c.Cmd()...)
	}
	return args
}
