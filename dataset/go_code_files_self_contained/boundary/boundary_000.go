package main

import (
	"fmt"
	"os"
)

func applySymlinks() error {
	fmt.Println("Removing bad symlinks")
	var err error
	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Frameworks") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/Frameworks"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Libraries") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/Libraries"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Resources") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/Resources"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/ThrustShell Framework") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/ThrustShell Framework"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Libraries") == false {
		if err = os.Remove(getAppDirectory() + "/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Libraries"); err != nil {
			return err
		}
	}

	fmt.Println("Applying Symlinks")

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/A",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Frameworks") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Frameworks",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Frameworks"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Libraries") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Libraries",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Libraries"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Resources") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Resources",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Resources"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/ThrustShell Framework") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/ThrustShell Framework",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/ThrustShell Framework"); err != nil {
			return err
		}
	}

	if PathNotExist(getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Libraries") == true {
		if err = os.Symlink(
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/A/Libraries/Libraries",
			getAppDirectory()+"/Contents/Frameworks/ThrustShell Framework.framework/Versions/Current/Libraries"); err != nil {
			return err
		}
	}

	return nil
}
