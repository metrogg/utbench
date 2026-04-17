func (p *pluginManager) LoadPlugin(details *pluginDetails, emitter gomit.Emitter) (*loadedPlugin, serror.SnapError) {
	type result struct {
		lp  *loadedPlugin
		err serror.SnapError
	}
	resultChan := make(chan result)
	go func() {
		lPlugin := new(loadedPlugin)
		lPlugin.Details = details
		lPlugin.State = DetectedState

		var (
			ePlugin *plugin.ExecutablePlugin
			resp    plugin.Response
			err     error
		)

		if lPlugin.Details.Uri == nil {
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"path":   filepath.Base(lPlugin.Details.Exec[0]),
			}).Info("plugin load called")
			// We will create commands by appending the ExecPath to the actual command.
			// The ExecPath is a temporary location where the plugin/package will be
			// run from.
			commands := make([]string, len(lPlugin.Details.Exec))
			for i, e := range lPlugin.Details.Exec {
				commands[i] = filepath.Join(lPlugin.Details.ExecPath, e)
			}

			ePlugin, err = plugin.NewExecutablePlugin(
				p.GenerateArgs(int(log.GetLevel())).
					SetCertPath(details.CertPath).
					SetKeyPath(details.KeyPath).
					SetCACertPaths(details.CACertPaths).
					SetTLSEnabled(details.TLSEnabled),
				commands...)
			if err != nil {
				pmLogger.WithFields(log.Fields{
					"_block": "load-plugin",
					"error":  err.Error(),
				}).Error("load plugin error while creating executable plugin")
				resultChan <- result{nil, serror.New(err)}
				return
			}
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"path":   lPlugin.Details.Exec,
			}).Debug(fmt.Sprintf("plugin load timeout set to %ds", p.pluginLoadTimeout))
			resp, err = ePlugin.Run(time.Second * time.Duration(p.pluginLoadTimeout))
			if err != nil {
				pmLogger.WithFields(log.Fields{
					"_block": "load-plugin",
					"error":  err.Error(),
				}).Error("load plugin error when starting plugin")
				resultChan <- result{nil, serror.New(err)}
				return
			}

			ePlugin.SetName(resp.Meta.Name)

			key := fmt.Sprintf("%s"+core.Separator+"%s"+core.Separator+"%d", resp.Meta.Type.String(), resp.Meta.Name, resp.Meta.Version)
			if _, exists := p.loadedPlugins.table[key]; exists {
				resultChan <- result{nil, serror.New(ErrPluginAlreadyLoaded, map[string]interface{}{
					"plugin-name":    resp.Meta.Name,
					"plugin-version": resp.Meta.Version,
					"plugin-type":    resp.Type.String(),
				})}
			}
		} else {
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"uri":    lPlugin.Details.Uri.String(),
			}).Info("plugin load called")
			res, err := http.Get(lPlugin.Details.Uri.String())
			if err != nil {
				resultChan <- result{nil, serror.New(err)}
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				resultChan <- result{nil, serror.New(err)}
				return
			}
			err = json.Unmarshal(body, &resp)
			if err != nil {
				pmLogger.WithFields(log.Fields{
					"_block": "load-plugin",
					"error":  err.Error(),
				}).Error("error during json unmarshal")
			}
		}
		ap, err := newAvailablePlugin(resp, emitter, ePlugin, p.grpcSecurity)
		if err != nil {
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"error":  err.Error(),
			}).Error("load plugin error while creating available plugin")
			resultChan <- result{nil, serror.New(err)}
			return
		}

		if lPlugin.Details.Uri != nil {
			ap.SetIsRemote(true)
		}

		if resp.Meta.Unsecure {
			err = ap.client.Ping()
		} else {
			err = ap.client.SetKey()
		}

		if err != nil {
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"error":  err.Error(),
			}).Error("load plugin error while pinging the plugin")
			resultChan <- result{nil, serror.New(err)}
			return
		}

		// Get the ConfigPolicy and add it to the loaded plugin
		c, ok := ap.client.(plugin.Plugin)
		if !ok {
			resultChan <- result{nil, serror.New(errors.New("missing GetConfigPolicy function"))}
			return
		}
		cp, err := c.GetConfigPolicy()
		if err != nil {
			pmLogger.WithFields(log.Fields{
				"_block":         "load-plugin",
				"plugin-type":    "collector",
				"error":          err.Error(),
				"plugin-name":    ap.Name(),
				"plugin-version": ap.Version(),
				"plugin-id":      ap.ID(),
			}).Error("error in getting config policy")
			resultChan <- result{nil, serror.New(err)}
			return
		}

		lPlugin.ConfigPolicy = cp
		lPlugin.Meta = resp.Meta
		lPlugin.Type = resp.Type
		lPlugin.Token = resp.Token
		lPlugin.LoadedTime = time.Now()
		lPlugin.State = LoadedState

		if resp.Type == plugin.CollectorPluginType || resp.Type == plugin.StreamCollectorPluginType {
			cfgNode := p.pluginConfig.getPluginConfigDataNode(core.PluginType(resp.Type), resp.Meta.Name, resp.Meta.Version)

			if lPlugin.ConfigPolicy != nil {
				// Get plugin config defaults
				defaults := cdata.NewNode()
				for _, cpolicy := range lPlugin.ConfigPolicy.GetAll() {
					_, errs := cpolicy.AddDefaults(defaults.Table())
					if len(errs.Errors()) > 0 {
						for _, err := range errs.Errors() {
							pmLogger.WithFields(log.Fields{
								"_block":         "load-plugin",
								"plugin-type":    "collector",
								"plugin-name":    ap.Name(),
								"plugin-version": ap.Version(),
								"plugin-id":      ap.ID(),
							}).Error(err.Error())
						}
						resultChan <- result{nil, serror.New(errors.New("error getting default config"))}
						return

					}
				}

				// Update config policy with defaults
				cfgNode = cfgNode.ReverseMerge(defaults)
				cp, err = c.GetConfigPolicy()
				if err != nil {
					pmLogger.WithFields(log.Fields{
						"_block":         "load-plugin",
						"plugin-type":    "collector",
						"error":          err.Error(),
						"plugin-name":    ap.Name(),
						"plugin-version": ap.Version(),
						"plugin-id":      ap.ID(),
					}).Error("error in getting config policy")
					resultChan <- result{nil, serror.New(err)}
					return
				}
				lPlugin.ConfigPolicy = cp
			}

			colClient := ap.client.(client.PluginCollectorClient)
			if !ap.isRemote {
				defer ap.client.(client.PluginCollectorClient).Close()
			}

			cfg := plugin.ConfigType{
				ConfigDataNode: cfgNode,
			}

			metricTypes, err := colClient.GetMetricTypes(cfg)
			if err != nil {
				pmLogger.WithFields(log.Fields{
					"_block":         "load-plugin",
					"plugin-type":    resp.Type.String(),
					"error":          err.Error(),
					"plugin-name":    ap.Name(),
					"plugin-version": ap.Version(),
				}).Error("error in getting metric types")
				resultChan <- result{nil, serror.New(err)}
				return
			}

			// Add metric types to metric catalog
			for _, nmt := range metricTypes {
				// If the version is 0 default it to the plugin version
				// This honors the plugins explicit version but falls back
				// to the plugin version as default
				if nmt.Version() < 1 {
					// Since we have to override version we convert to a internal struct
					nmt = &metricType{
						namespace:          nmt.Namespace(),
						version:            resp.Meta.Version,
						lastAdvertisedTime: nmt.LastAdvertisedTime(),
						config:             nmt.Config(),
						data:               nmt.Data(),
						tags:               nmt.Tags(),
						description:        nmt.Description(),
						unit:               nmt.Unit(),
					}
				}
				// We quit and throw an error on bad metric versions (<1)
				// the is a safety catch otherwise the catalog will be corrupted
				if nmt.Version() < 1 {
					err := errors.New("Bad metric version from plugin")
					pmLogger.WithFields(log.Fields{
						"_block":           "load-plugin",
						"plugin-name":      resp.Meta.Name,
						"plugin-version":   resp.Meta.Version,
						"plugin-type":      resp.Meta.Type.String(),
						"plugin-path":      filepath.Base(lPlugin.Details.ExecPath),
						"metric-namespace": nmt.Namespace(),
						"metric-version":   nmt.Version(),
						"error":            err.Error(),
					}).Error("received metric with bad version")
					resultChan <- result{nil, serror.New(err)}
					return
				}

				//Add standard tags
				nmt = p.AddStandardAndWorkflowTags(nmt, nil)

				if err := p.metricCatalog.AddLoadedMetricType(lPlugin, nmt); err != nil {
					pmLogger.WithFields(log.Fields{
						"_block":           "load-plugin",
						"plugin-name":      resp.Meta.Name,
						"plugin-version":   resp.Meta.Version,
						"plugin-type":      resp.Meta.Type.String(),
						"plugin-path":      filepath.Base(lPlugin.Details.ExecPath),
						"metric-namespace": nmt.Namespace(),
						"metric-version":   nmt.Version(),
						"error":            err.Error(),
					}).Error("error adding loaded metric type")
					resultChan <- result{nil, serror.New(err)}
					return
				}
			}
		}

		if lPlugin.Details.Uri == nil {
			// Added so clients can adequately clean up connections
			ap.client.Kill("Retrieved necessary plugin info")
			err = ePlugin.Kill()
			if err != nil {
				pmLogger.WithFields(log.Fields{
					"_block": "load-plugin",
					"error":  err.Error(),
				}).Error("load plugin error while killing plugin executable plugin")
				resultChan <- result{nil, serror.New(err)}
				return
			}
		}

		if resp.State != plugin.PluginSuccess {
			e := fmt.Errorf("plugin loading did not succeed: %s\n", resp.ErrorMessage)
			pmLogger.WithFields(log.Fields{
				"_block":          "load-plugin",
				"error":           e,
				"plugin response": resp.ErrorMessage,
			}).Error("load plugin error")
			resultChan <- result{nil, serror.New(e)}
			return
		}

		aErr := p.loadedPlugins.add(lPlugin)
		if aErr != nil {
			pmLogger.WithFields(log.Fields{
				"_block": "load-plugin",
				"error":  aErr,
			}).Error("load plugin error while adding loaded plugin to load plugins collection")
			resultChan <- result{nil, aErr}
		}
		if ap.isRemote && aErr == nil {
			// monitor standalone plugins. Unload them from the plugin catalog and metrics list
			// when we detect they are no longer online.
			go func() {
				defer ap.client.(client.PluginCollectorClient).Close()
				for {
					time.Sleep(5 * time.Second)
					go ap.CheckHealth()
					if ap.failedHealthChecks > 3 {
						p.UnloadPlugin(lPlugin)
						return
					}
					if _, err := p.loadedPlugins.get(lPlugin.Key()); err != nil {
						// prevent leaking routine when plugin is unloaded normally
						return
					}

				}
			}()

		}
		resultChan <- result{lPlugin, nil}
		return
	}()

	select {
	case results := <-resultChan:
		return results.lp, results.err
	case <-time.After(time.Second * time.Duration(p.pluginLoadTimeout)):
		e := serror.New(errors.New("timed out waiting for plugin to load"))
		return nil, e
	}
}
