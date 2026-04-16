func (cli *CLI) ParseFlags(args []string) (*config.Config, []string, bool, bool, bool, error) {
	var dry, once, isVersion bool

	c := config.DefaultConfig()

	if s := os.Getenv("CT_LOCAL_CONFIG"); s != "" {
		envConfig, err := config.Parse(s)
		if err != nil {
			return nil, nil, false, false, false, err
		}
		c = c.Merge(envConfig)
	}

	// configPaths stores the list of configuration paths on disk
	configPaths := make([]string, 0, 6)

	// Parse the flags and options
	flags := flag.NewFlagSet(version.Name, flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)
	flags.Usage = func() {}

	flags.Var((funcVar)(func(s string) error {
		configPaths = append(configPaths, s)
		return nil
	}), "config", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.Address = config.String(s)
		return nil
	}), "consul-addr", "")

	flags.Var((funcVar)(func(s string) error {
		a, err := config.ParseAuthConfig(s)
		if err != nil {
			return err
		}
		c.Consul.Auth = a
		return nil
	}), "consul-auth", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Consul.Retry.Enabled = config.Bool(b)
		return nil
	}), "consul-retry", "")

	flags.Var((funcIntVar)(func(i int) error {
		c.Consul.Retry.Attempts = config.Int(i)
		return nil
	}), "consul-retry-attempts", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Retry.Backoff = config.TimeDuration(d)
		return nil
	}), "consul-retry-backoff", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Retry.MaxBackoff = config.TimeDuration(d)
		return nil
	}), "consul-retry-max-backoff", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Consul.SSL.Enabled = config.Bool(b)
		return nil
	}), "consul-ssl", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.SSL.CaCert = config.String(s)
		return nil
	}), "consul-ssl-ca-cert", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.SSL.CaPath = config.String(s)
		return nil
	}), "consul-ssl-ca-path", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.SSL.Cert = config.String(s)
		return nil
	}), "consul-ssl-cert", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.SSL.Key = config.String(s)
		return nil
	}), "consul-ssl-key", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.SSL.ServerName = config.String(s)
		return nil
	}), "consul-ssl-server-name", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Consul.SSL.Verify = config.Bool(b)
		return nil
	}), "consul-ssl-verify", "")

	flags.Var((funcVar)(func(s string) error {
		c.Consul.Token = config.String(s)
		return nil
	}), "consul-token", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Transport.DialKeepAlive = config.TimeDuration(d)
		return nil
	}), "consul-transport-dial-keep-alive", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Transport.DialTimeout = config.TimeDuration(d)
		return nil
	}), "consul-transport-dial-timeout", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Consul.Transport.DisableKeepAlives = config.Bool(b)
		return nil
	}), "consul-transport-disable-keep-alives", "")

	flags.Var((funcIntVar)(func(i int) error {
		c.Consul.Transport.MaxIdleConnsPerHost = config.Int(i)
		return nil
	}), "consul-transport-max-idle-conns-per-host", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Transport.TLSHandshakeTimeout = config.TimeDuration(d)
		return nil
	}), "consul-transport-tls-handshake-timeout", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Dedup.Enabled = config.Bool(b)
		return nil
	}), "dedup", "")

	flags.BoolVar(&dry, "dry", false, "")

	flags.Var((funcVar)(func(s string) error {
		c.Exec.Enabled = config.Bool(true)
		c.Exec.Command = config.String(s)
		return nil
	}), "exec", "")

	flags.Var((funcVar)(func(s string) error {
		sig, err := signals.Parse(s)
		if err != nil {
			return err
		}
		c.Exec.KillSignal = config.Signal(sig)
		return nil
	}), "exec-kill-signal", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Exec.KillTimeout = config.TimeDuration(d)
		return nil
	}), "exec-kill-timeout", "")

	flags.Var((funcVar)(func(s string) error {
		sig, err := signals.Parse(s)
		if err != nil {
			return err
		}
		c.Exec.ReloadSignal = config.Signal(sig)
		return nil
	}), "exec-reload-signal", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Exec.Splay = config.TimeDuration(d)
		return nil
	}), "exec-splay", "")

	flags.Var((funcVar)(func(s string) error {
		sig, err := signals.Parse(s)
		if err != nil {
			return err
		}
		c.KillSignal = config.Signal(sig)
		return nil
	}), "kill-signal", "")

	flags.Var((funcVar)(func(s string) error {
		c.LogLevel = config.String(s)
		return nil
	}), "log-level", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.MaxStale = config.TimeDuration(d)
		return nil
	}), "max-stale", "")

	flags.BoolVar(&once, "once", false, "")

	flags.Var((funcVar)(func(s string) error {
		c.PidFile = config.String(s)
		return nil
	}), "pid-file", "")

	flags.Var((funcVar)(func(s string) error {
		sig, err := signals.Parse(s)
		if err != nil {
			return err
		}
		c.ReloadSignal = config.Signal(sig)
		return nil
	}), "reload-signal", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Consul.Retry.Backoff = config.TimeDuration(d)
		return nil
	}), "retry", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Syslog.Enabled = config.Bool(b)
		return nil
	}), "syslog", "")

	flags.Var((funcVar)(func(s string) error {
		c.Syslog.Facility = config.String(s)
		return nil
	}), "syslog-facility", "")

	flags.Var((funcVar)(func(s string) error {
		t, err := config.ParseTemplateConfig(s)
		if err != nil {
			return err
		}
		*c.Templates = append(*c.Templates, t)
		return nil
	}), "template", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.Address = config.String(s)
		return nil
	}), "vault-addr", "")

	flags.Var((funcDurationVar)(func(t time.Duration) error {
		c.Vault.Grace = config.TimeDuration(t)
		return nil
	}), "vault-grace", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.RenewToken = config.Bool(b)
		return nil
	}), "vault-renew-token", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.Retry.Enabled = config.Bool(b)
		return nil
	}), "vault-retry", "")

	flags.Var((funcIntVar)(func(i int) error {
		c.Vault.Retry.Attempts = config.Int(i)
		return nil
	}), "vault-retry-attempts", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Vault.Retry.Backoff = config.TimeDuration(d)
		return nil
	}), "vault-retry-backoff", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Vault.Retry.MaxBackoff = config.TimeDuration(d)
		return nil
	}), "vault-retry-max-backoff", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.SSL.Enabled = config.Bool(b)
		return nil
	}), "vault-ssl", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.SSL.CaCert = config.String(s)
		return nil
	}), "vault-ssl-ca-cert", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.SSL.CaPath = config.String(s)
		return nil
	}), "vault-ssl-ca-path", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.SSL.Cert = config.String(s)
		return nil
	}), "vault-ssl-cert", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.SSL.Key = config.String(s)
		return nil
	}), "vault-ssl-key", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.SSL.ServerName = config.String(s)
		return nil
	}), "vault-ssl-server-name", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.SSL.Verify = config.Bool(b)
		return nil
	}), "vault-ssl-verify", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Vault.Transport.DialKeepAlive = config.TimeDuration(d)
		return nil
	}), "vault-transport-dial-keep-alive", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Vault.Transport.DialTimeout = config.TimeDuration(d)
		return nil
	}), "vault-transport-dial-timeout", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.Transport.DisableKeepAlives = config.Bool(b)
		return nil
	}), "vault-transport-disable-keep-alives", "")

	flags.Var((funcIntVar)(func(i int) error {
		c.Vault.Transport.MaxIdleConnsPerHost = config.Int(i)
		return nil
	}), "vault-transport-max-idle-conns-per-host", "")

	flags.Var((funcDurationVar)(func(d time.Duration) error {
		c.Vault.Transport.TLSHandshakeTimeout = config.TimeDuration(d)
		return nil
	}), "vault-transport-tls-handshake-timeout", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.Token = config.String(s)
		return nil
	}), "vault-token", "")

	flags.Var((funcVar)(func(s string) error {
		c.Vault.VaultAgentTokenFile = config.String(s)
		return nil
	}), "vault-agent-token-file", "")

	flags.Var((funcBoolVar)(func(b bool) error {
		c.Vault.UnwrapToken = config.Bool(b)
		return nil
	}), "vault-unwrap-token", "")

	flags.Var((funcVar)(func(s string) error {
		w, err := config.ParseWaitConfig(s)
		if err != nil {
			return err
		}
		c.Wait = w
		return nil
	}), "wait", "")

	flags.BoolVar(&isVersion, "v", false, "")
	flags.BoolVar(&isVersion, "version", false, "")

	// If there was a parser error, stop
	if err := flags.Parse(args); err != nil {
		return nil, nil, false, false, false, err
	}

	// Error if extra arguments are present
	args = flags.Args()
	if len(args) > 0 {
		return nil, nil, false, false, false, fmt.Errorf("cli: extra args: %q", args)
	}

	return c, configPaths, once, dry, isVersion, nil
}
