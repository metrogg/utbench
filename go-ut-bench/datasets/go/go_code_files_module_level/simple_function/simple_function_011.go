func NewFactory(optionalClientConfig clientcmd.ClientConfig) *Factory {
	mapper := kubectl.ShortcutExpander{RESTMapper: api.RESTMapper}

	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.SetNormalizeFunc(util.WarnWordSepNormalizeFunc) // Warn for "_" flags

	clientConfig := optionalClientConfig
	if optionalClientConfig == nil {
		clientConfig = DefaultClientConfig(flags)
	}

	clients := NewClientCache(clientConfig)

	return &Factory{
		clients: clients,
		flags:   flags,

		Object: func() (meta.RESTMapper, runtime.ObjectTyper) {
			cfg, err := clientConfig.ClientConfig()
			CheckErr(err)
			cmdApiVersion := unversioned.GroupVersion{}
			if cfg.GroupVersion != nil {
				cmdApiVersion = *cfg.GroupVersion
			}

			return kubectl.OutputVersionMapper{RESTMapper: mapper, OutputVersions: []unversioned.GroupVersion{cmdApiVersion}}, api.Scheme
		},
		Client: func() (*client.Client, error) {
			return clients.ClientForVersion(nil)
		},
		ClientConfig: func() (*client.Config, error) {
			return clients.ClientConfigForVersion(nil)
		},
		ClientForMapping: func(mapping *meta.RESTMapping) (resource.RESTClient, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			if err != nil {
				return nil, err
			}
			switch mapping.GroupVersionKind.Group {
			case api.GroupName:
				return client.RESTClient, nil
			case autoscaling.GroupName:
				return client.AutoscalingClient.RESTClient, nil
			case batch.GroupName:
				return client.BatchClient.RESTClient, nil
			case extensions.GroupName:
				return client.ExtensionsClient.RESTClient, nil
			}
			return nil, fmt.Errorf("unable to get RESTClient for resource '%s'", mapping.Resource)
		},
		Describer: func(mapping *meta.RESTMapping) (kubectl.Describer, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			if err != nil {
				return nil, err
			}
			if describer, ok := kubectl.DescriberFor(mapping.GroupVersionKind.GroupKind(), client); ok {
				return describer, nil
			}
			return nil, fmt.Errorf("no description has been implemented for %q", mapping.GroupVersionKind.Kind)
		},
		Decoder: func(toInternal bool) runtime.Decoder {
			if toInternal {
				return api.Codecs.UniversalDecoder()
			}
			return api.Codecs.UniversalDeserializer()
		},
		JSONEncoder: func() runtime.Encoder {
			return api.Codecs.LegacyCodec(registered.EnabledVersions()...)
		},
		Printer: func(mapping *meta.RESTMapping, noHeaders, withNamespace bool, wide bool, showAll bool, showLabels bool, absoluteTimestamps bool, columnLabels []string) (kubectl.ResourcePrinter, error) {
			return kubectl.NewHumanReadablePrinter(noHeaders, withNamespace, wide, showAll, showLabels, absoluteTimestamps, columnLabels), nil
		},
		PodSelectorForObject: func(object runtime.Object) (string, error) {
			// TODO: replace with a swagger schema based approach (identify pod selector via schema introspection)
			switch t := object.(type) {
			case *api.ReplicationController:
				return kubectl.MakeLabels(t.Spec.Selector), nil
			case *api.Pod:
				if len(t.Labels) == 0 {
					return "", fmt.Errorf("the pod has no labels and cannot be exposed")
				}
				return kubectl.MakeLabels(t.Labels), nil
			case *api.Service:
				if t.Spec.Selector == nil {
					return "", fmt.Errorf("the service has no pod selector set")
				}
				return kubectl.MakeLabels(t.Spec.Selector), nil
			case *extensions.Deployment:
				selector, err := unversioned.LabelSelectorAsSelector(t.Spec.Selector)
				if err != nil {
					return "", fmt.Errorf("invalid label selector: %v", err)
				}
				return selector.String(), nil
			case *extensions.ReplicaSet:
				selector, err := unversioned.LabelSelectorAsSelector(t.Spec.Selector)
				if err != nil {
					return "", fmt.Errorf("failed to convert label selector to selector: %v", err)
				}
				return selector.String(), nil
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return "", err
				}
				return "", fmt.Errorf("cannot extract pod selector from %v", gvk)
			}
		},
		MapBasedSelectorForObject: func(object runtime.Object) (string, error) {
			// TODO: replace with a swagger schema based approach (identify pod selector via schema introspection)
			switch t := object.(type) {
			case *api.ReplicationController:
				return kubectl.MakeLabels(t.Spec.Selector), nil
			case *api.Pod:
				if len(t.Labels) == 0 {
					return "", fmt.Errorf("the pod has no labels and cannot be exposed")
				}
				return kubectl.MakeLabels(t.Labels), nil
			case *api.Service:
				if t.Spec.Selector == nil {
					return "", fmt.Errorf("the service has no pod selector set")
				}
				return kubectl.MakeLabels(t.Spec.Selector), nil
			case *extensions.Deployment:
				// TODO(madhusudancs): Make this smarter by admitting MatchExpressions with Equals
				// operator, DoubleEquals operator and In operator with only one element in the set.
				if len(t.Spec.Selector.MatchExpressions) > 0 {
					return "", fmt.Errorf("couldn't convert expressions - \"%+v\" to map-based selector format")
				}
				return kubectl.MakeLabels(t.Spec.Selector.MatchLabels), nil
			case *extensions.ReplicaSet:
				// TODO(madhusudancs): Make this smarter by admitting MatchExpressions with Equals
				// operator, DoubleEquals operator and In operator with only one element in the set.
				if len(t.Spec.Selector.MatchExpressions) > 0 {
					return "", fmt.Errorf("couldn't convert expressions - \"%+v\" to map-based selector format")
				}
				return kubectl.MakeLabels(t.Spec.Selector.MatchLabels), nil
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return "", err
				}
				return "", fmt.Errorf("cannot extract pod selector from %v", gvk)
			}
		},
		PortsForObject: func(object runtime.Object) ([]string, error) {
			// TODO: replace with a swagger schema based approach (identify pod selector via schema introspection)
			switch t := object.(type) {
			case *api.ReplicationController:
				return getPorts(t.Spec.Template.Spec), nil
			case *api.Pod:
				return getPorts(t.Spec), nil
			case *api.Service:
				return getServicePorts(t.Spec), nil
			case *extensions.Deployment:
				return getPorts(t.Spec.Template.Spec), nil
			case *extensions.ReplicaSet:
				return getPorts(t.Spec.Template.Spec), nil
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("cannot extract ports from %v", gvk)
			}
		},
		LabelsForObject: func(object runtime.Object) (map[string]string, error) {
			return meta.NewAccessor().Labels(object)
		},
		LogsForObject: func(object, options runtime.Object) (*client.Request, error) {
			c, err := clients.ClientForVersion(nil)
			if err != nil {
				return nil, err
			}

			switch t := object.(type) {
			case *api.Pod:
				opts, ok := options.(*api.PodLogOptions)
				if !ok {
					return nil, errors.New("provided options object is not a PodLogOptions")
				}
				return c.Pods(t.Namespace).GetLogs(t.Name, opts), nil
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("cannot get the logs from %v", gvk)
			}
		},
		PauseObject: func(object runtime.Object) (bool, error) {
			c, err := clients.ClientForVersion(nil)
			if err != nil {
				return false, err
			}

			switch t := object.(type) {
			case *extensions.Deployment:
				if t.Spec.Paused {
					return true, nil
				}
				t.Spec.Paused = true
				_, err := c.Extensions().Deployments(t.Namespace).Update(t)
				return false, err
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return false, err
				}
				return false, fmt.Errorf("cannot pause %v", gvk)
			}
		},
		ResumeObject: func(object runtime.Object) (bool, error) {
			c, err := clients.ClientForVersion(nil)
			if err != nil {
				return false, err
			}

			switch t := object.(type) {
			case *extensions.Deployment:
				if !t.Spec.Paused {
					return true, nil
				}
				t.Spec.Paused = false
				_, err := c.Extensions().Deployments(t.Namespace).Update(t)
				return false, err
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return false, err
				}
				return false, fmt.Errorf("cannot resume %v", gvk)
			}
		},
		Scaler: func(mapping *meta.RESTMapping) (kubectl.Scaler, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			if err != nil {
				return nil, err
			}
			return kubectl.ScalerFor(mapping.GroupVersionKind.GroupKind(), client)
		},
		Reaper: func(mapping *meta.RESTMapping) (kubectl.Reaper, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			if err != nil {
				return nil, err
			}
			return kubectl.ReaperFor(mapping.GroupVersionKind.GroupKind(), client)
		},
		HistoryViewer: func(mapping *meta.RESTMapping) (kubectl.HistoryViewer, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			clientset := clientset.FromUnversionedClient(client)
			if err != nil {
				return nil, err
			}
			return kubectl.HistoryViewerFor(mapping.GroupVersionKind.GroupKind(), clientset)
		},
		Rollbacker: func(mapping *meta.RESTMapping) (kubectl.Rollbacker, error) {
			mappingVersion := mapping.GroupVersionKind.GroupVersion()
			client, err := clients.ClientForVersion(&mappingVersion)
			if err != nil {
				return nil, err
			}
			return kubectl.RollbackerFor(mapping.GroupVersionKind.GroupKind(), client)
		},
		Validator: func(validate bool, cacheDir string) (validation.Schema, error) {
			if validate {
				client, err := clients.ClientForVersion(nil)
				if err != nil {
					return nil, err
				}
				dir := cacheDir
				if len(dir) > 0 {
					version, err := client.ServerVersion()
					if err != nil {
						return nil, err
					}
					dir = path.Join(cacheDir, version.String())
				}
				return &clientSwaggerSchema{
					c:        client,
					cacheDir: dir,
					mapper:   api.RESTMapper,
				}, nil
			}
			return validation.NullSchema{}, nil
		},
		SwaggerSchema: func(gvk unversioned.GroupVersionKind) (*swagger.ApiDeclaration, error) {
			version := gvk.GroupVersion()
			client, err := clients.ClientForVersion(&version)
			if err != nil {
				return nil, err
			}
			return client.Discovery().SwaggerSchema(version)
		},
		DefaultNamespace: func() (string, bool, error) {
			return clientConfig.Namespace()
		},
		Generators: func(cmdName string) map[string]kubectl.Generator {
			return DefaultGenerators(cmdName)
		},
		CanBeExposed: func(kind unversioned.GroupKind) error {
			switch kind {
			case api.Kind("ReplicationController"), api.Kind("Service"), api.Kind("Pod"), extensions.Kind("Deployment"), extensions.Kind("ReplicaSet"):
				// nothing to do here
			default:
				return fmt.Errorf("cannot expose a %s", kind)
			}
			return nil
		},
		CanBeAutoscaled: func(kind unversioned.GroupKind) error {
			switch kind {
			case api.Kind("ReplicationController"), extensions.Kind("Deployment"):
				// nothing to do here
			default:
				return fmt.Errorf("cannot autoscale a %v", kind)
			}
			return nil
		},
		AttachablePodForObject: func(object runtime.Object) (*api.Pod, error) {
			client, err := clients.ClientForVersion(nil)
			if err != nil {
				return nil, err
			}
			switch t := object.(type) {
			case *api.ReplicationController:
				selector := labels.SelectorFromSet(t.Spec.Selector)
				return GetFirstPod(client, t.Namespace, selector)
			case *extensions.Deployment:
				selector, err := unversioned.LabelSelectorAsSelector(t.Spec.Selector)
				if err != nil {
					return nil, fmt.Errorf("invalid label selector: %v", err)
				}
				return GetFirstPod(client, t.Namespace, selector)
			case *extensions.Job:
				selector, err := unversioned.LabelSelectorAsSelector(t.Spec.Selector)
				if err != nil {
					return nil, fmt.Errorf("invalid label selector: %v", err)
				}
				return GetFirstPod(client, t.Namespace, selector)
			case *api.Pod:
				return t, nil
			default:
				gvk, err := api.Scheme.ObjectKind(object)
				if err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("cannot attach to %v: not implemented", gvk)
			}
		},
		EditorEnvs: func() []string {
			return []string{"KUBE_EDITOR", "EDITOR"}
		},
	}
}
