func (c *containerLXC) startCommon() (string, error) {
	// Load the go-lxc struct
	err := c.initLXC(true)
	if err != nil {
		return "", errors.Wrap(err, "Load go-lxc struct")
	}

	// Check that we're not already running
	if c.IsRunning() {
		return "", fmt.Errorf("The container is already running")
	}

	// Sanity checks for devices
	for name, m := range c.expandedDevices {
		switch m["type"] {
		case "disk":
			// When we want to attach a storage volume created via
			// the storage api m["source"] only contains the name of
			// the storage volume, not the path where it is mounted.
			// So do only check for the existence of m["source"]
			// when m["pool"] is empty.
			if m["pool"] == "" && m["source"] != "" && !shared.IsTrue(m["optional"]) && !shared.PathExists(shared.HostPath(m["source"])) {
				return "", fmt.Errorf("Missing source '%s' for disk '%s'", m["source"], name)
			}
		case "nic":
			if m["parent"] != "" && !shared.PathExists(fmt.Sprintf("/sys/class/net/%s", m["parent"])) {
				return "", fmt.Errorf("Missing parent '%s' for nic '%s'", m["parent"], name)
			}
		case "unix-char", "unix-block":
			srcPath, exist := m["source"]
			if !exist {
				srcPath = m["path"]
			}

			if srcPath != "" && m["required"] != "" && !shared.IsTrue(m["required"]) {
				err = deviceInotifyAddClosestLivingAncestor(c.state, filepath.Dir(srcPath))
				if err != nil {
					logger.Errorf("Failed to add \"%s\" to inotify targets", srcPath)
					return "", fmt.Errorf("Failed to setup inotify watch for '%s': %v", srcPath, err)
				}
			} else if srcPath != "" && m["major"] == "" && m["minor"] == "" && !shared.PathExists(srcPath) {
				return "", fmt.Errorf("Missing source '%s' for device '%s'", srcPath, name)
			}
		}
	}

	// Load any required kernel modules
	kernelModules := c.expandedConfig["linux.kernel_modules"]
	if kernelModules != "" {
		for _, module := range strings.Split(kernelModules, ",") {
			module = strings.TrimPrefix(module, " ")
			err := util.LoadModule(module)
			if err != nil {
				return "", fmt.Errorf("Failed to load kernel module '%s': %s", module, err)
			}
		}
	}

	var ourStart bool
	newSize, ok := c.LocalConfig()["volatile.apply_quota"]
	if ok {
		err := c.initStorage()
		if err != nil {
			return "", errors.Wrap(err, "Initialize storage")
		}

		size, err := shared.ParseByteSizeString(newSize)
		if err != nil {
			return "", err
		}
		err = c.storage.StorageEntitySetQuota(storagePoolVolumeTypeContainer, size, c)
		if err != nil {
			return "", errors.Wrap(err, "Set storage quota")
		}

		// Remove the volatile key from the DB
		err = c.state.Cluster.ContainerConfigRemove(c.id, "volatile.apply_quota")
		if err != nil {
			return "", errors.Wrap(err, "Remove volatile.apply_quota config key")
		}

		// Remove the volatile key from the in-memory configs
		delete(c.localConfig, "volatile.apply_quota")
		delete(c.expandedConfig, "volatile.apply_quota")
	}

	/* Deal with idmap changes */
	nextIdmap, err := c.NextIdmap()
	if err != nil {
		return "", errors.Wrap(err, "Set ID map")
	}

	diskIdmap, err := c.DiskIdmap()
	if err != nil {
		return "", errors.Wrap(err, "Set last ID map")
	}

	if !nextIdmap.Equals(diskIdmap) && !(diskIdmap == nil && c.state.OS.Shiftfs) {
		if shared.IsTrue(c.expandedConfig["security.protection.shift"]) {
			return "", fmt.Errorf("Container is protected against filesystem shifting")
		}

		logger.Debugf("Container idmap changed, remapping")
		c.updateProgress("Remapping container filesystem")

		ourStart, err = c.StorageStart()
		if err != nil {
			return "", errors.Wrap(err, "Storage start")
		}

		if diskIdmap != nil {
			if c.Storage().GetStorageType() == storageTypeZfs {
				err = diskIdmap.UnshiftRootfs(c.RootfsPath(), zfsIdmapSetSkipper)
			} else {
				err = diskIdmap.UnshiftRootfs(c.RootfsPath(), nil)
			}
			if err != nil {
				if ourStart {
					c.StorageStop()
				}
				return "", err
			}
		}

		if nextIdmap != nil && !c.state.OS.Shiftfs {
			if c.Storage().GetStorageType() == storageTypeZfs {
				err = nextIdmap.ShiftRootfs(c.RootfsPath(), zfsIdmapSetSkipper)
			} else {
				err = nextIdmap.ShiftRootfs(c.RootfsPath(), nil)
			}
			if err != nil {
				if ourStart {
					c.StorageStop()
				}
				return "", err
			}
		}

		jsonDiskIdmap := "[]"
		if nextIdmap != nil && !c.state.OS.Shiftfs {
			idmapBytes, err := json.Marshal(nextIdmap.Idmap)
			if err != nil {
				return "", err
			}
			jsonDiskIdmap = string(idmapBytes)
		}

		err = c.ConfigKeySet("volatile.last_state.idmap", jsonDiskIdmap)
		if err != nil {
			return "", errors.Wrapf(err, "Set volatile.last_state.idmap config key on container %q (id %d)", c.name, c.id)
		}

		c.updateProgress("")
	}

	var idmapBytes []byte
	if nextIdmap == nil {
		idmapBytes = []byte("[]")
	} else {
		idmapBytes, err = json.Marshal(nextIdmap.Idmap)
		if err != nil {
			return "", err
		}
	}

	if c.localConfig["volatile.idmap.current"] != string(idmapBytes) {
		err = c.ConfigKeySet("volatile.idmap.current", string(idmapBytes))
		if err != nil {
			return "", errors.Wrapf(err, "Set volatile.idmap.current config key on container %q (id %d)", c.name, c.id)
		}
	}

	// Generate the Seccomp profile
	if err := SeccompCreateProfile(c); err != nil {
		return "", err
	}

	// Cleanup any existing leftover devices
	c.removeUnixDevices()
	c.removeDiskDevices()
	c.removeNetworkFilters()
	c.removeProxyDevices()

	var usbs []usbDevice
	var sriov []string
	diskDevices := map[string]types.Device{}

	// Create the devices
	for _, k := range c.expandedDevices.DeviceNames() {
		m := c.expandedDevices[k]
		if shared.StringInSlice(m["type"], []string{"unix-char", "unix-block"}) {
			// Unix device
			paths, err := c.createUnixDevice(fmt.Sprintf("unix.%s", k), m, true)
			if err != nil {
				// Deal with device hotplug
				if m["required"] == "" || shared.IsTrue(m["required"]) {
					return "", err
				}

				srcPath := m["source"]
				if srcPath == "" {
					srcPath = m["path"]
				}
				srcPath = shared.HostPath(srcPath)

				err = deviceInotifyAddClosestLivingAncestor(c.state, srcPath)
				if err != nil {
					logger.Errorf("Failed to add \"%s\" to inotify targets", srcPath)
					return "", err
				}
				continue
			}
			devPath := paths[0]
			if c.isCurrentlyPrivileged() && !c.state.OS.RunningInUserNS && c.state.OS.CGroupDevicesController {
				// Add the new device cgroup rule
				dType, dMajor, dMinor, err := deviceGetAttributes(devPath)
				if err != nil {
					if m["required"] == "" || shared.IsTrue(m["required"]) {
						return "", err
					}
				} else {
					err = lxcSetConfigItem(c.c, "lxc.cgroup.devices.allow", fmt.Sprintf("%s %d:%d rwm", dType, dMajor, dMinor))
					if err != nil {
						return "", fmt.Errorf("Failed to add cgroup rule for device")
					}
				}
			}
		} else if m["type"] == "usb" {
			if usbs == nil {
				usbs, err = deviceLoadUsb()
				if err != nil {
					return "", err
				}
			}

			for _, usb := range usbs {
				if (m["vendorid"] != "" && usb.vendor != m["vendorid"]) || (m["productid"] != "" && usb.product != m["productid"]) {
					continue
				}

				err := c.setupUnixDevice(fmt.Sprintf("unix.%s", k), m, usb.major, usb.minor, usb.path, shared.IsTrue(m["required"]), false)
				if err != nil {
					return "", err
				}
			}
		} else if m["type"] == "gpu" {
			allGpus := deviceWantsAllGPUs(m)
			gpus, nvidiaDevices, err := deviceLoadGpu(allGpus)
			if err != nil {
				return "", err
			}

			sawNvidia := false
			found := false
			for _, gpu := range gpus {
				if (m["vendorid"] != "" && gpu.vendorID != m["vendorid"]) ||
					(m["pci"] != "" && gpu.pci != m["pci"]) ||
					(m["productid"] != "" && gpu.productID != m["productid"]) ||
					(m["id"] != "" && gpu.id != m["id"]) {
					continue
				}

				found = true

				err := c.setupUnixDevice(fmt.Sprintf("unix.%s", k), m, gpu.major, gpu.minor, gpu.path, true, false)
				if err != nil {
					return "", err
				}

				if !gpu.isNvidia {
					continue
				}

				if gpu.nvidia.path != "" {
					err = c.setupUnixDevice(fmt.Sprintf("unix.%s", k), m, gpu.nvidia.major, gpu.nvidia.minor, gpu.nvidia.path, true, false)
					if err != nil {
						return "", err
					}
				} else if !allGpus {
					errMsg := fmt.Errorf("Failed to detect correct \"/dev/nvidia\" path")
					logger.Errorf("%s", errMsg)
					return "", errMsg
				}

				sawNvidia = true
			}

			if sawNvidia {
				for _, gpu := range nvidiaDevices {
					if shared.IsTrue(c.expandedConfig["nvidia.runtime"]) {
						if !gpu.isCard {
							continue
						}
					}
					err := c.setupUnixDevice(fmt.Sprintf("unix.%s", k), m, gpu.major, gpu.minor, gpu.path, true, false)
					if err != nil {
						return "", err
					}
				}
			}

			if !found {
				msg := "Failed to detect requested GPU device"
				logger.Error(msg)
				return "", fmt.Errorf(msg)
			}
		} else if m["type"] == "disk" {
			if m["path"] != "/" {
				diskDevices[k] = m
			}
		} else if m["type"] == "nic" || m["type"] == "infiniband" {
			var err error
			var infiniband map[string]IBF
			if m["type"] == "infiniband" {
				infiniband, err = deviceLoadInfiniband()
				if err != nil {
					return "", err
				}
			}

			networkKeyPrefix := "lxc.net"
			if !util.RuntimeLiblxcVersionAtLeast(2, 1, 0) {
				networkKeyPrefix = "lxc.network"
			}

			m, err = c.fillNetworkDevice(k, m)
			if err != nil {
				return "", err
			}

			networkidx := -1
			reserved := []string{}
			// Record nictype == physical devices since those won't
			// be available for nictype == sriov.
			for _, dName := range c.expandedDevices.DeviceNames() {
				m := c.expandedDevices[dName]
				if m["type"] != "nic" && m["type"] != "infiniband" {
					continue
				}

				if m["nictype"] != "physical" {
					continue
				}

				reserved = append(reserved, m["parent"])
			}

			for _, dName := range c.expandedDevices.DeviceNames() {
				m := c.expandedDevices[dName]
				if m["type"] != "nic" && m["type"] != "infiniband" {
					continue
				}
				networkidx++

				if shared.StringInSlice(dName, sriov) {
					continue
				} else {
					sriov = append(sriov, dName)
				}

				if m["nictype"] != "sriov" {
					continue
				}

				m, err = c.fillSriovNetworkDevice(dName, m, reserved)
				if err != nil {
					return "", err
				}

				// Make sure that no one called dibs.
				reserved = append(reserved, m["host_name"])

				val := c.c.ConfigItem(fmt.Sprintf("%s.%d.type", networkKeyPrefix, networkidx))
				if len(val) == 0 || val[0] != "phys" {
					return "", fmt.Errorf("Network index corresponds to false network")
				}

				// Fill in correct name right now
				err = lxcSetConfigItem(c.c, fmt.Sprintf("%s.%d.link", networkKeyPrefix, networkidx), m["host_name"])
				if err != nil {
					return "", err
				}

				if m["type"] == "infiniband" {
					key := m["host_name"]
					ifDev, ok := infiniband[key]
					if !ok {
						return "", fmt.Errorf("Specified infiniband device \"%s\" not found", key)
					}

					err := c.addInfinibandDevices(dName, &ifDev, false)
					if err != nil {
						return "", err
					}
				}
			}

			if m["type"] == "infiniband" && m["nictype"] == "physical" {
				key := m["parent"]
				ifDev, ok := infiniband[key]
				if !ok {
					return "", fmt.Errorf("Specified infiniband device \"%s\" not found", key)
				}

				err := c.addInfinibandDevices(k, &ifDev, false)
				if err != nil {
					return "", err
				}
			}

			if m["nictype"] == "bridged" && shared.IsTrue(m["security.mac_filtering"]) {
				// Read device name from config
				vethName := ""
				for i := 0; i < len(c.c.ConfigItem(networkKeyPrefix)); i++ {
					val := c.c.ConfigItem(fmt.Sprintf("%s.%d.hwaddr", networkKeyPrefix, i))
					if len(val) == 0 || val[0] != m["hwaddr"] {
						continue
					}

					val = c.c.ConfigItem(fmt.Sprintf("%s.%d.link", networkKeyPrefix, i))
					if len(val) == 0 || val[0] != m["parent"] {
						continue
					}

					val = c.c.ConfigItem(fmt.Sprintf("%s.%d.veth.pair", networkKeyPrefix, i))
					if len(val) == 0 {
						continue
					}

					vethName = val[0]
					break
				}

				if vethName == "" {
					return "", fmt.Errorf("Failed to find device name for mac_filtering")
				}

				err = c.createNetworkFilter(vethName, m["parent"], m["hwaddr"])
				if err != nil {
					return "", err
				}
			}

			// Create VLAN devices
			if shared.StringInSlice(m["nictype"], []string{"macvlan", "physical"}) && m["vlan"] != "" {
				device := networkGetHostDevice(m["parent"], m["vlan"])
				if !shared.PathExists(fmt.Sprintf("/sys/class/net/%s", device)) {
					_, err := shared.RunCommand("ip", "link", "add", "link", m["parent"], "name", device, "up", "type", "vlan", "id", m["vlan"])
					if err != nil {
						return "", err
					}

					// Attempt to disable IPv6 router advertisement acceptance
					networkSysctlSet(fmt.Sprintf("ipv6/conf/%s/accept_ra", device), "0")
				}
			}
		}
	}

	err = c.addDiskDevices(diskDevices, func(name string, d types.Device) error {
		_, err := c.createDiskDevice(name, d)
		return err
	})
	if err != nil {
		return "", err
	}

	// Create any missing directory
	err = os.MkdirAll(c.LogPath(), 0700)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(c.DevicesPath(), 0711)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(c.ShmountsPath(), 0711)
	if err != nil {
		return "", err
	}

	// Rotate the log file
	logfile := c.LogFilePath()
	if shared.PathExists(logfile) {
		os.Remove(logfile + ".old")
		err := os.Rename(logfile, logfile+".old")
		if err != nil {
			return "", err
		}
	}

	// Storage is guaranteed to be mountable now.
	ourStart, err = c.StorageStart()
	if err != nil {
		return "", err
	}

	// Generate the LXC config
	configPath := filepath.Join(c.LogPath(), "lxc.conf")
	err = c.c.SaveConfigFile(configPath)
	if err != nil {
		os.Remove(configPath)
		return "", err
	}

	// Undo liblxc modifying container directory ownership
	err = os.Chown(c.Path(), 0, 0)
	if err != nil {
		if ourStart {
			c.StorageStop()
		}
		return "", err
	}

	// Set right permission to allow traversal
	var mode os.FileMode
	if c.isCurrentlyPrivileged() {
		mode = 0700
	} else {
		mode = 0711
	}

	err = os.Chmod(c.Path(), mode)
	if err != nil {
		if ourStart {
			c.StorageStop()
		}
		return "", err
	}

	// Update the backup.yaml file
	err = writeBackupFile(c)
	if err != nil {
		if ourStart {
			c.StorageStop()
		}
		return "", err
	}

	if !c.IsStateful() && shared.PathExists(c.StatePath()) {
		os.RemoveAll(c.StatePath())
	}

	_, err = c.StorageStop()
	if err != nil {
		return "", err
	}

	// Update time container was last started
	err = c.state.Cluster.ContainerLastUsedUpdate(c.id, time.Now().UTC())
	if err != nil {
		return "", fmt.Errorf("Error updating last used: %v", err)
	}

	// Unmount any previously mounted shiftfs
	syscall.Unmount(c.RootfsPath(), syscall.MNT_DETACH)

	return configPath, nil
}
