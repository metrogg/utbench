package main

func Bootstrap() error {
	if ExecutableNotExist() == true {
		var err error
		if err = prepareExecutable(); err != nil {
			return err
		}
		if err = prepareInfoPropertiesListTemplate(); err != nil {
			return err
		}

		return writeInfoPropertiesList()
	}

	return nil
}
