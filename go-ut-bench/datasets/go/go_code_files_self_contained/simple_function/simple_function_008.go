package main

func ShutdownAll() error {
	errs := []error{}
	for root, _ := range driversByRoot {
		if err := ShutdownDriver(root); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return ErrBadDriverShutdown
	}
	return nil
}
