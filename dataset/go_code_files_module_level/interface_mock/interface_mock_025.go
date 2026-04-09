func (d *Daemon) EnableK8sWatcher(queueSize uint) error {
	if !k8s.IsEnabled() {
		log.Debug("Not enabling k8s event listener because k8s is not enabled")
		return nil
	}
	log.Info("Enabling k8s event listener")

	restConfig, err := k8s.CreateConfig()
	if err != nil {
		return fmt.Errorf("Unable to create rest configuration: %s", err)
	}

	apiextensionsclientset, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("Unable to create rest configuration for k8s CRD: %s", err)
	}

	err = cilium_v2.CreateCustomResourceDefinitions(apiextensionsclientset)
	if err != nil {
		return fmt.Errorf("Unable to create custom resource definition: %s", err)
	}
	d.k8sAPIGroups.addAPI(k8sAPIGroupCRD)

	ciliumNPClient := k8s.CiliumClient()

	serKNPs := serializer.NewFunctionQueue(queueSize)
	serSvcs := serializer.NewFunctionQueue(queueSize)
	serEps := serializer.NewFunctionQueue(queueSize)
	serIngresses := serializer.NewFunctionQueue(queueSize)
	serCNPs := serializer.NewFunctionQueue(queueSize)
	serPods := serializer.NewFunctionQueue(queueSize)
	serNodes := serializer.NewFunctionQueue(queueSize)
	serNamespaces := serializer.NewFunctionQueue(queueSize)

	_, policyController := informer.NewInformer(
		cache.NewListWatchFromClient(k8s.Client().NetworkingV1().RESTClient(),
			"networkpolicies", v1.NamespaceAll, fields.Everything()),
		&networkingv1.NetworkPolicy{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricKNP, metricCreate, valid, equal) }()
				if k8sNP := k8s.CopyObjToV1NetworkPolicy(obj); k8sNP != nil {
					valid = true
					serKNPs.Enqueue(func() error {
						err := d.addK8sNetworkPolicyV1(k8sNP)
						updateK8sEventMetric(metricKNP, metricCreate, err == nil)
						return nil
					}, serializer.NoRetry)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricKNP, metricUpdate, valid, equal) }()
				if oldK8sNP := k8s.CopyObjToV1NetworkPolicy(oldObj); oldK8sNP != nil {
					valid = true
					if newK8sNP := k8s.CopyObjToV1NetworkPolicy(newObj); newK8sNP != nil {
						if k8s.EqualV1NetworkPolicy(oldK8sNP, newK8sNP) {
							equal = true
							return
						}

						serKNPs.Enqueue(func() error {
							err := d.updateK8sNetworkPolicyV1(oldK8sNP, newK8sNP)
							updateK8sEventMetric(metricKNP, metricUpdate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricKNP, metricDelete, valid, equal) }()
				k8sNP := k8s.CopyObjToV1NetworkPolicy(obj)
				if k8sNP == nil {
					deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
					if !ok {
						return
					}
					// Delete was not observed by the watcher but is
					// removed from kube-apiserver. This is the last
					// known state and the object no longer exists.
					k8sNP = k8s.CopyObjToV1NetworkPolicy(deletedObj.Obj)
					if k8sNP == nil {
						return
					}
				}

				valid = true
				serKNPs.Enqueue(func() error {
					err := d.deleteK8sNetworkPolicyV1(k8sNP)
					updateK8sEventMetric(metricKNP, metricDelete, err == nil)
					return nil
				}, serializer.NoRetry)
			},
		},
		k8s.ConvertToNetworkPolicy,
	)
	d.blockWaitGroupToSyncResources(wait.NeverStop, policyController, k8sAPIGroupNetworkingV1Core)
	go policyController.Run(wait.NeverStop)

	d.k8sAPIGroups.addAPI(k8sAPIGroupNetworkingV1Core)

	_, svcController := informer.NewInformer(
		cache.NewListWatchFromClient(k8s.Client().CoreV1().RESTClient(),
			"services", v1.NamespaceAll, fields.Everything()),
		&v1.Service{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricService, metricCreate, valid, equal) }()
				if k8sSvc := k8s.CopyObjToV1Services(obj); k8sSvc != nil {
					valid = true
					serSvcs.Enqueue(func() error {
						err := d.addK8sServiceV1(k8sSvc)
						updateK8sEventMetric(metricService, metricCreate, err == nil)
						return nil
					}, serializer.NoRetry)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricService, metricUpdate, valid, equal) }()
				if oldk8sSvc := k8s.CopyObjToV1Services(oldObj); oldk8sSvc != nil {
					valid = true
					if newk8sSvc := k8s.CopyObjToV1Services(newObj); newk8sSvc != nil {
						if k8s.EqualV1Services(oldk8sSvc, newk8sSvc) {
							equal = true
							return
						}

						serSvcs.Enqueue(func() error {
							err := d.updateK8sServiceV1(oldk8sSvc, newk8sSvc)
							updateK8sEventMetric(metricService, metricUpdate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricService, metricDelete, valid, equal) }()
				k8sSvc := k8s.CopyObjToV1Services(obj)
				if k8sSvc == nil {
					deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
					if !ok {
						return
					}
					// Delete was not observed by the watcher but is
					// removed from kube-apiserver. This is the last
					// known state and the object no longer exists.
					k8sSvc = k8s.CopyObjToV1Services(deletedObj.Obj)
					if k8sSvc == nil {
						return
					}
				}

				valid = true
				serSvcs.Enqueue(func() error {
					err := d.deleteK8sServiceV1(k8sSvc)
					updateK8sEventMetric(metricService, metricDelete, err == nil)
					return nil
				}, serializer.NoRetry)
			},
		},
		k8s.ConvertToK8sService,
	)
	d.blockWaitGroupToSyncResources(wait.NeverStop, svcController, k8sAPIGroupServiceV1Core)
	go svcController.Run(wait.NeverStop)
	d.k8sAPIGroups.addAPI(k8sAPIGroupServiceV1Core)

	_, endpointController := informer.NewInformer(
		cache.NewListWatchFromClient(k8s.Client().CoreV1().RESTClient(),
			"endpoints", v1.NamespaceAll,
			fields.ParseSelectorOrDie(option.Config.K8sWatcherEndpointSelector),
		),
		&v1.Endpoints{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricEndpoint, metricCreate, valid, equal) }()
				if k8sEP := k8s.CopyObjToV1Endpoints(obj); k8sEP != nil {
					valid = true
					serEps.Enqueue(func() error {
						err := d.addK8sEndpointV1(k8sEP)
						updateK8sEventMetric(metricEndpoint, metricCreate, err == nil)
						return nil
					}, serializer.NoRetry)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricEndpoint, metricUpdate, valid, equal) }()
				if oldk8sEP := k8s.CopyObjToV1Endpoints(oldObj); oldk8sEP != nil {
					valid = true
					if newk8sEP := k8s.CopyObjToV1Endpoints(newObj); newk8sEP != nil {
						if k8s.EqualV1Endpoints(oldk8sEP, newk8sEP) {
							equal = true
							return
						}

						serEps.Enqueue(func() error {
							err := d.updateK8sEndpointV1(oldk8sEP, newk8sEP)
							updateK8sEventMetric(metricEndpoint, metricUpdate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricEndpoint, metricDelete, valid, equal) }()
				k8sEP := k8s.CopyObjToV1Endpoints(obj)
				if k8sEP == nil {
					deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
					if !ok {
						return
					}
					// Delete was not observed by the watcher but is
					// removed from kube-apiserver. This is the last
					// known state and the object no longer exists.
					k8sEP = k8s.CopyObjToV1Endpoints(deletedObj.Obj)
					if k8sEP == nil {
						return
					}
				}
				valid = true
				serEps.Enqueue(func() error {
					err := d.deleteK8sEndpointV1(k8sEP)
					updateK8sEventMetric(metricEndpoint, metricDelete, err == nil)
					return nil
				}, serializer.NoRetry)
			},
		},
		k8s.ConvertToK8sEndpoints,
	)
	d.blockWaitGroupToSyncResources(wait.NeverStop, endpointController, k8sAPIGroupEndpointV1Core)
	go endpointController.Run(wait.NeverStop)
	d.k8sAPIGroups.addAPI(k8sAPIGroupEndpointV1Core)

	if option.Config.IsLBEnabled() {
		_, ingressController := informer.NewInformer(
			cache.NewListWatchFromClient(k8s.Client().ExtensionsV1beta1().RESTClient(),
				"ingresses", v1.NamespaceAll, fields.Everything()),
			&v1beta1.Ingress{},
			0,
			cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					var valid, equal bool
					defer func() { d.k8sEventReceived(metricEndpoint, metricCreate, valid, equal) }()
					if k8sIngress := k8s.CopyObjToV1beta1Ingress(obj); k8sIngress != nil {
						valid = true
						serIngresses.Enqueue(func() error {
							err := d.addIngressV1beta1(k8sIngress)
							updateK8sEventMetric(metricIngress, metricCreate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				},
				UpdateFunc: func(oldObj, newObj interface{}) {
					var valid, equal bool
					defer func() { d.k8sEventReceived(metricEndpoint, metricUpdate, valid, equal) }()
					if oldk8sIngress := k8s.CopyObjToV1beta1Ingress(oldObj); oldk8sIngress != nil {
						valid = true
						if newk8sIngress := k8s.CopyObjToV1beta1Ingress(newObj); newk8sIngress != nil {
							if k8s.EqualV1beta1Ingress(oldk8sIngress, newk8sIngress) {
								equal = true
								return
							}

							serIngresses.Enqueue(func() error {
								err := d.updateIngressV1beta1(oldk8sIngress, newk8sIngress)
								updateK8sEventMetric(metricIngress, metricUpdate, err == nil)
								return nil
							}, serializer.NoRetry)
						}
					}
				},
				DeleteFunc: func(obj interface{}) {
					var valid, equal bool
					defer func() { d.k8sEventReceived(metricEndpoint, metricDelete, valid, equal) }()
					k8sIngress := k8s.CopyObjToV1beta1Ingress(obj)
					if k8sIngress == nil {
						deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
						if !ok {
							return
						}
						// Delete was not observed by the watcher but is
						// removed from kube-apiserver. This is the last
						// known state and the object no longer exists.
						k8sIngress = k8s.CopyObjToV1beta1Ingress(deletedObj.Obj)
						if k8sIngress == nil {
							return
						}
					}
					valid = true
					serEps.Enqueue(func() error {
						err := d.deleteIngressV1beta1(k8sIngress)
						updateK8sEventMetric(metricIngress, metricDelete, err == nil)
						return nil
					}, serializer.NoRetry)
				},
			},
			k8s.ConvertToIngress,
		)
		d.blockWaitGroupToSyncResources(wait.NeverStop, ingressController, k8sAPIGroupIngressV1Beta1)
		go ingressController.Run(wait.NeverStop)
		d.k8sAPIGroups.addAPI(k8sAPIGroupIngressV1Beta1)
	}

	var (
		cnpEventStore    cache.Store
		cnpConverterFunc informer.ConvertFunc
	)
	cnpStore := cache.NewStore(cache.DeletionHandlingMetaNamespaceKeyFunc)
	switch {
	case k8sversion.Capabilities().Patch:
		// k8s >= 1.13 does not require a store to update CNP status so
		// we don't even need to keep the status of a CNP with us.
		cnpConverterFunc = k8s.ConvertToCNP
	default:
		cnpEventStore = cnpStore
		cnpConverterFunc = k8s.ConvertToCNPWithStatus
	}

	ciliumV2Controller := informer.NewInformerWithStore(
		cache.NewListWatchFromClient(ciliumNPClient.CiliumV2().RESTClient(),
			"ciliumnetworkpolicies", v1.NamespaceAll, fields.Everything()),
		&cilium_v2.CiliumNetworkPolicy{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricCNP, metricCreate, valid, equal) }()
				if cnp := k8s.CopyObjToV2CNP(obj); cnp != nil {
					valid = true
					serCNPs.Enqueue(func() error {
						if cnp.RequiresDerivative() {
							return nil
						}
						err := d.addCiliumNetworkPolicyV2(ciliumNPClient, cnpEventStore, cnp)
						updateK8sEventMetric(metricCNP, metricCreate, err == nil)
						return nil
					}, serializer.NoRetry)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricCNP, metricUpdate, valid, equal) }()
				if oldCNP := k8s.CopyObjToV2CNP(oldObj); oldCNP != nil {
					valid = true
					if newCNP := k8s.CopyObjToV2CNP(newObj); newCNP != nil {
						if k8s.EqualV2CNP(oldCNP, newCNP) {
							equal = true
							return
						}

						serCNPs.Enqueue(func() error {
							if newCNP.RequiresDerivative() {
								return nil
							}

							err := d.updateCiliumNetworkPolicyV2(ciliumNPClient, cnpEventStore, oldCNP, newCNP)
							updateK8sEventMetric(metricCNP, metricUpdate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricCNP, metricDelete, valid, equal) }()
				cnp := k8s.CopyObjToV2CNP(obj)
				if cnp == nil {
					deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
					if !ok {
						return
					}
					// Delete was not observed by the watcher but is
					// removed from kube-apiserver. This is the last
					// known state and the object no longer exists.
					cnp = k8s.CopyObjToV2CNP(deletedObj.Obj)
					if cnp == nil {
						return
					}
				}
				valid = true
				serCNPs.Enqueue(func() error {
					err := d.deleteCiliumNetworkPolicyV2(cnp)
					updateK8sEventMetric(metricCNP, metricDelete, err == nil)
					return nil
				}, serializer.NoRetry)
			},
		},
		cnpConverterFunc,
		cnpStore,
	)
	d.blockWaitGroupToSyncResources(wait.NeverStop, ciliumV2Controller, k8sAPIGroupCiliumV2)
	go ciliumV2Controller.Run(wait.NeverStop)
	d.k8sAPIGroups.addAPI(k8sAPIGroupCiliumV2)

	asyncControllers := sync.WaitGroup{}
	asyncControllers.Add(1)
	go func() {
		var once sync.Once
		for {
			createPodController := func(fieldSelector fields.Selector) cache.Controller {
				_, podController := informer.NewInformer(
					cache.NewListWatchFromClient(k8s.Client().CoreV1().RESTClient(),
						"pods", v1.NamespaceAll, fieldSelector),
					&v1.Pod{},
					0,
					cache.ResourceEventHandlerFuncs{
						AddFunc: func(obj interface{}) {
							var valid, equal bool
							defer func() { d.k8sEventReceived(metricPod, metricCreate, valid, equal) }()
							if pod := k8s.CopyObjToV1Pod(obj); pod != nil {
								valid = true
								serPods.Enqueue(func() error {
									err := d.addK8sPodV1(pod)
									updateK8sEventMetric(metricPod, metricCreate, err == nil)
									return nil
								}, serializer.NoRetry)
							}
						},
						UpdateFunc: func(oldObj, newObj interface{}) {
							var valid, equal bool
							defer func() { d.k8sEventReceived(metricPod, metricUpdate, valid, equal) }()
							if oldPod := k8s.CopyObjToV1Pod(oldObj); oldPod != nil {
								valid = true
								if newPod := k8s.CopyObjToV1Pod(newObj); newPod != nil {
									if k8s.EqualV1Pod(oldPod, newPod) {
										equal = true
										return
									}

									serPods.Enqueue(func() error {
										err := d.updateK8sPodV1(oldPod, newPod)
										updateK8sEventMetric(metricPod, metricUpdate, err == nil)
										return nil
									}, serializer.NoRetry)
								}
							}
						},
						DeleteFunc: func(obj interface{}) {
							var valid, equal bool
							defer func() { d.k8sEventReceived(metricPod, metricDelete, valid, equal) }()
							if pod := k8s.CopyObjToV1Pod(obj); pod != nil {
								valid = true
								serPods.Enqueue(func() error {
									err := d.deleteK8sPodV1(pod)
									updateK8sEventMetric(metricPod, metricDelete, err == nil)
									return nil
								}, serializer.NoRetry)
							}
						},
					},
					k8s.ConvertToPod,
				)
				return podController
			}
			podController := createPodController(fields.Everything())

			isConnected := make(chan struct{})
			// once isConnected is closed, it will stop waiting on caches to be
			// synchronized.
			d.blockWaitGroupToSyncResources(isConnected, podController, k8sAPIGroupPodV1Core)
			once.Do(func() {
				asyncControllers.Done()
				d.k8sAPIGroups.addAPI(k8sAPIGroupPodV1Core)
			})
			go podController.Run(isConnected)

			if !option.Config.K8sEventHandover {
				return
			}

			// Replace pod controller by only receiving events from our own
			// node once we are connected to the kvstore.

			<-kvstore.Client().Connected()
			close(isConnected)

			log.WithField(logfields.Node, node.GetName()).Info("Connected to KVStore, watching for pod events on node")
			// Only watch for pod events for our node.
			podController = createPodController(fields.ParseSelectorOrDie("spec.nodeName=" + node.GetName()))
			isConnected = make(chan struct{})
			go podController.Run(isConnected)

			// Create a new pod controller when we are disconnected with the
			// kvstore
			<-kvstore.Client().Disconnected()
			close(isConnected)
			log.Info("Disconnected from KVStore, watching for pod events all nodes")
		}
	}()

	asyncControllers.Add(1)
	go func() {
		var once sync.Once
		for {
			_, nodeController := informer.NewInformer(
				cache.NewListWatchFromClient(k8s.Client().CoreV1().RESTClient(),
					"nodes", v1.NamespaceAll, fields.Everything()),
				&v1.Node{},
				0,
				cache.ResourceEventHandlerFuncs{
					AddFunc: func(obj interface{}) {
						var valid, equal bool
						defer func() { d.k8sEventReceived(metricNode, metricCreate, valid, equal) }()
						if Node := k8s.CopyObjToV1Node(obj); Node != nil {
							valid = true
							serNodes.Enqueue(func() error {
								err := d.addK8sNodeV1(Node)
								updateK8sEventMetric(metricNode, metricCreate, err == nil)
								return nil
							}, serializer.NoRetry)
						}
					},
					UpdateFunc: func(oldObj, newObj interface{}) {
						var valid, equal bool
						defer func() { d.k8sEventReceived(metricNode, metricUpdate, valid, equal) }()
						if oldNode := k8s.CopyObjToV1Node(oldObj); oldNode != nil {
							valid = true
							if newNode := k8s.CopyObjToV1Node(newObj); newNode != nil {
								if k8s.EqualV1Node(oldNode, newNode) {
									equal = true
									return
								}

								serNodes.Enqueue(func() error {
									err := d.updateK8sNodeV1(oldNode, newNode)
									updateK8sEventMetric(metricNode, metricUpdate, err == nil)
									return nil
								}, serializer.NoRetry)
							}
						}
					},
					DeleteFunc: func(obj interface{}) {
						var valid, equal bool
						defer func() { d.k8sEventReceived(metricNode, metricDelete, valid, equal) }()
						node := k8s.CopyObjToV1Node(obj)
						if node == nil {
							deletedObj, ok := obj.(cache.DeletedFinalStateUnknown)
							if !ok {
								return
							}
							// Delete was not observed by the watcher but is
							// removed from kube-apiserver. This is the last
							// known state and the object no longer exists.
							node = k8s.CopyObjToV1Node(deletedObj.Obj)
							if node == nil {
								return
							}
						}
						valid = true
						serNodes.Enqueue(func() error {
							err := d.deleteK8sNodeV1(node)
							updateK8sEventMetric(metricNode, metricDelete, err == nil)
							return nil
						}, serializer.NoRetry)
					},
				},
				k8s.ConvertToNode,
			)
			isConnected := make(chan struct{})
			// once isConnected is closed, it will stop waiting on caches to be
			// synchronized.
			d.blockWaitGroupToSyncResources(isConnected, nodeController, k8sAPIGroupNodeV1Core)

			once.Do(func() {
				// Signalize that we have put node controller in the wait group
				// to sync resources.
				asyncControllers.Done()
			})
			d.k8sAPIGroups.addAPI(k8sAPIGroupNodeV1Core)
			go nodeController.Run(isConnected)
			// TODO: do we really need to wait until etcd watcher signalizes
			// "listDone" or connected to it is sufficient?

			if !option.Config.K8sEventHandover {
				return
			}

			<-kvstore.Client().Connected()
			close(isConnected)

			log.Info("Connected to KVStore, stopping k8s node watcher")

			d.k8sAPIGroups.removeAPI(k8sAPIGroupNodeV1Core)
			// Create a new node controller when we are disconnected with the
			// kvstore
			<-kvstore.Client().Disconnected()

			log.Info("Disconnected from KVStore, restarting k8s node watcher")
		}
	}()

	_, namespaceController := informer.NewInformer(
		cache.NewListWatchFromClient(k8s.Client().CoreV1().RESTClient(),
			"namespaces", v1.NamespaceAll, fields.Everything()),
		&v1.Namespace{},
		0,
		cache.ResourceEventHandlerFuncs{
			// AddFunc does not matter since the endpoint will fetch
			// namespace labels when the endpoint is created
			// DelFunc does not matter since, when a namespace is deleted, all
			// pods belonging to that namespace are also deleted.
			UpdateFunc: func(oldObj, newObj interface{}) {
				var valid, equal bool
				defer func() { d.k8sEventReceived(metricNS, metricUpdate, valid, equal) }()
				if oldNS := k8s.CopyObjToV1Namespace(oldObj); oldNS != nil {
					valid = true
					if newNS := k8s.CopyObjToV1Namespace(newObj); newNS != nil {
						if k8s.EqualV1Namespace(oldNS, newNS) {
							equal = true
							return
						}

						serNamespaces.Enqueue(func() error {
							err := d.updateK8sV1Namespace(oldNS, newNS)
							updateK8sEventMetric(metricNS, metricUpdate, err == nil)
							return nil
						}, serializer.NoRetry)
					}
				}
			},
		},
		k8s.ConvertToNamespace,
	)

	go namespaceController.Run(wait.NeverStop)
	d.k8sAPIGroups.addAPI(k8sAPIGroupNamespaceV1Core)

	asyncControllers.Wait()

	return nil
}
