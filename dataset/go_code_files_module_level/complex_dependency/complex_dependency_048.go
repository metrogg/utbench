func (uw *UnitWriter) appSystemdUnit(pa *preparedApp, binPath string, opts []*unit.UnitOption) []*unit.UnitOption {
	if uw.err != nil {
		return nil
	}

	flavor, systemdVersion, err := GetFlavor(uw.p)
	if err != nil {
		uw.err = errwrap.Wrap(errors.New("unable to determine stage1 flavor"), err)
		return nil
	}

	ra := pa.app
	app := ra.App
	appName := ra.Name
	imgName := uw.p.AppNameToImageName(ra.Name)

	podAbsRoot, err := filepath.Abs(uw.p.Root)
	if err != nil {
		uw.err = err
		return nil
	}

	var supplementaryGroups []string
	for _, g := range app.SupplementaryGIDs {
		supplementaryGroups = append(supplementaryGroups, strconv.Itoa(g))
	}

	// Write env file
	if err := common.WriteEnvFile(common.ComposeEnviron(pa.env), &uw.p.UidRange, EnvFilePath(uw.p.Root, pa.app.Name)); err != nil {
		uw.err = errwrap.Wrapf("unable to write environment file", err)
		return nil
	}

	execStart := append([]string{binPath}, app.Exec[1:]...)
	execStartString := quoteExec(execStart)
	opts = append(opts,
		unit.NewUnitOption("Service", "ExecStart", execStartString),
		unit.NewUnitOption("Service", "RootDirectory", common.RelAppRootfsPath(appName)),
		unit.NewUnitOption("Service", "WorkingDirectory", app.WorkingDirectory),
		unit.NewUnitOption("Service", "EnvironmentFile", RelEnvFilePath(appName)),
		unit.NewUnitOption("Service", "User", strconv.Itoa(int(pa.uid))),
		unit.NewUnitOption("Service", "Group", strconv.Itoa(int(pa.gid))),
		unit.NewUnitOption("Service", "PermissionsStartOnly", "true"),
		unit.NewUnitOption("Unit", "Requires", InstantiatedPrepareAppUnitName(ra.Name)),
		unit.NewUnitOption("Unit", "After", InstantiatedPrepareAppUnitName(ra.Name)),
	)

	if len(supplementaryGroups) > 0 {
		opts = appendOptionsList(opts, "Service", "SupplementaryGroups", "", supplementaryGroups...)
	}

	if !uw.p.InsecureOptions.DisableCapabilities {
		opts = append(opts, unit.NewUnitOption("Service", "CapabilityBoundingSet", strings.Join(pa.capabilities, " ")))
	}

	// Apply seccomp isolator, if any and not opt-ing out;
	// see https://www.freedesktop.org/software/systemd/man/systemd.exec.html#SystemCallFilter=
	if pa.seccomp != nil {
		opts, err = seccompUnitOptions(opts, pa.seccomp)
		if err != nil {
			uw.err = errwrap.Wrapf("unable to apply seccomp options", err)
			return nil
		}
	}
	opts = append(opts, unit.NewUnitOption("Service", "NoNewPrivileges", strconv.FormatBool(pa.noNewPrivileges)))

	if ra.ReadOnlyRootFS {
		for _, m := range pa.mounts {
			mntPath, err := EvaluateSymlinksInsideApp(podAbsRoot, m.Mount.Path)
			if err != nil {
				uw.err = err
				return nil
			}

			if !m.ReadOnly {
				rwDir := filepath.Join(common.RelAppRootfsPath(ra.Name), mntPath)
				opts = appendOptionsList(opts, "Service", "ReadWriteDirectories", "", rwDir)
			}
		}
		opts = appendOptionsList(opts, "Service", "ReadOnlyDirectories", "", common.RelAppRootfsPath(ra.Name))
	}

	// Unless we have --insecure-options=paths, then do some path protections:
	//
	// * prevent access to sensitive kernel tunables
	// * Run the app in a separate mount namespace
	//
	if !uw.p.InsecureOptions.DisablePaths {
		// Systemd 231+ has InaccessiblePaths
		// older versions only have InaccessibleDirectories
		// Paths prepended with "-" are ignored if they don't exist.
		if systemdVersion >= 231 {
			opts = appendOptionsList(opts, "Service", "InaccessiblePaths", "-", pa.relAppPaths(pa.hiddenPaths)...)
			opts = appendOptionsList(opts, "Service", "InaccessiblePaths", "-", pa.relAppPaths(pa.hiddenDirs)...)
			opts = appendOptionsList(opts, "Service", "ReadOnlyPaths", "-", pa.relAppPaths(pa.roPaths)...)
		} else {
			opts = appendOptionsList(opts, "Service", "InaccessibleDirectories", "-", pa.relAppPaths(pa.hiddenDirs)...)
			opts = appendOptionsList(opts, "Service", "ReadOnlyDirectories", "-", pa.relAppPaths(pa.roPaths)...)
		}

		if systemdVersion >= 233 {
			// ProtectKernelTunables is introduced in systemd-v232 but didn't work
			// until v233 due to a systemd bug, see
			// https://github.com/systemd/systemd/pull/4594
			// However, from v233, setting ProtectKernelTunables + RootDirectory causes
			// MountAPIVFS to be enabled unconditionally, which we don't want.
			//
			// opts = append(opts, unit.NewUnitOption("Service", "ProtectKernelTunables", "true"))

			// MountAPIVFS is introduced in systemd-233. Don't let systemd mount /sys:
			// it is mounted by prepare-app (tested by TestVolumeSysfs)
			opts = append(opts, unit.NewUnitOption("Service", "MountAPIVFS", "false"))
		}

		// MountFlags=shared creates a new mount namespace and (as unintuitive
		// as it might seem) makes sure the mount is slave+shared.
		opts = append(opts, unit.NewUnitOption("Service", "MountFlags", "shared"))
	}

	// Generate default device policy for the app, as well as the list of allowed devices.
	// For kvm flavor, devices are VM-specific and restricting them is not strictly needed.
	if !uw.p.InsecureOptions.DisablePaths && flavor != "kvm" {
		opts = append(opts, unit.NewUnitOption("Service", "DevicePolicy", "closed"))
		deviceAllows, err := generateDeviceAllows(common.Stage1RootfsPath(podAbsRoot), appName, app.MountPoints, pa.mounts, &uw.p.UidRange)
		if err != nil {
			uw.err = err
			return nil
		}
		for _, dev := range deviceAllows {
			opts = append(opts, unit.NewUnitOption("Service", "DeviceAllow", dev))
		}
	}

	for _, eh := range app.EventHandlers {
		var typ string
		switch eh.Name {
		case "pre-start":
			typ = "ExecStartPre"
		case "post-stop":
			typ = "ExecStopPost"
		default:
			uw.err = fmt.Errorf("unrecognized eventHandler: %v", eh.Name)
			return nil
		}
		exec := quoteExec(eh.Exec)
		opts = append(opts, unit.NewUnitOption("Service", typ, exec))
	}

	// Resource isolators
	if pa.resources.MemoryLimit != nil {
		opts = append(opts, unit.NewUnitOption("Service", "MemoryLimit", strconv.FormatUint(*pa.resources.MemoryLimit, 10)))
	}
	if pa.resources.CPUQuota != nil {
		quota := strconv.FormatUint(*pa.resources.CPUQuota, 10) + "%"
		opts = append(opts, unit.NewUnitOption("Service", "CPUQuota", quota))
	}
	if pa.resources.LinuxCPUShares != nil {
		opts = append(opts, unit.NewUnitOption("Service", "CPUShares", strconv.FormatUint(*pa.resources.LinuxCPUShares, 10)))
	}
	if pa.resources.LinuxOOMScoreAdjust != nil {
		opts = append(opts, unit.NewUnitOption("Service", "OOMScoreAdjust", strconv.Itoa(*pa.resources.LinuxOOMScoreAdjust)))
	}

	var saPorts []types.Port
	for _, p := range ra.App.Ports {
		if p.SocketActivated {
			saPorts = append(saPorts, p)
		}
	}

	if len(saPorts) > 0 {
		sockopts := []*unit.UnitOption{
			unit.NewUnitOption("Unit", "Description", fmt.Sprintf("Application=%v Image=%v %s", appName, imgName, "socket-activated ports")),
			unit.NewUnitOption("Unit", "DefaultDependencies", "false"),
			unit.NewUnitOption("Socket", "BindIPv6Only", "both"),
			unit.NewUnitOption("Socket", "Service", ServiceUnitName(appName)),
		}

		for _, sap := range saPorts {
			var proto string
			switch sap.Protocol {
			case "tcp":
				proto = "ListenStream"
			case "udp":
				proto = "ListenDatagram"
			default:
				uw.err = fmt.Errorf("unrecognized protocol: %v", sap.Protocol)
				return nil
			}
			// We find the host port for the pod's port and use that in the
			// socket unit file.
			// This is so because systemd inside the pod will match based on
			// the socket port number, and since the socket was created on the
			// host, it will have the host port number.
			port := findHostPort(*uw.p.Manifest, sap.Name)
			if port == 0 {
				log.Printf("warning: no --port option for socket-activated port %q, assuming port %d as specified in the manifest", sap.Name, sap.Port)
				port = sap.Port
			}
			sockopts = append(sockopts, unit.NewUnitOption("Socket", proto, fmt.Sprintf("%v", port)))
		}

		file, err := os.OpenFile(SocketUnitPath(uw.p.Root, appName), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			uw.err = errwrap.Wrap(errors.New("failed to create socket file"), err)
			return nil
		}
		defer file.Close()

		if _, err = io.Copy(file, unit.Serialize(sockopts)); err != nil {
			uw.err = errwrap.Wrap(errors.New("failed to write socket unit file"), err)
			return nil
		}

		if err = os.Symlink(path.Join("..", SocketUnitName(appName)), SocketWantPath(uw.p.Root, appName)); err != nil {
			uw.err = errwrap.Wrap(errors.New("failed to link socket want"), err)
			return nil
		}

		opts = append(opts, unit.NewUnitOption("Unit", "Requires", SocketUnitName(appName)))
	}
	return opts
}
