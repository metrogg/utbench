package main

func RunCommonBootstrapDialog() error {
	for _, spec := range bootstrapSpecs {
		if err := BootstrapConfig(spec); err != nil {
			return err
		}
	}
	return nil
}
