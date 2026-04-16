func (c *ResTableConfig) IsMoreSpecificThan(o *ResTableConfig) bool {
	// nil ResTableConfig is never more specific than any ResTableConfig
	if c == nil {
		return false
	}
	if o == nil {
		return false
	}

	// imsi
	if c.Mcc != o.Mcc {
		if c.Mcc == 0 {
			return false
		}
		if o.Mnc == 0 {
			return true
		}
	}
	if c.Mnc != o.Mnc {
		if c.Mnc == 0 {
			return false
		}
		if o.Mnc == 0 {
			return true
		}
	}

	// locale
	if diff := c.IsLocaleMoreSpecificThan(o); diff < 0 {
		return false
	} else if diff > 0 {
		return true
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if ((c.ScreenLayout ^ o.ScreenLayout) & MaskLayoutDir) != 0 {
			if (c.ScreenLayout & MaskLayoutDir) == 0 {
				return false
			}
			if (o.ScreenLayout & MaskLayoutDir) == 0 {
				return true
			}
		}
	}

	// smallest screen width dp
	if c.SmallestScreenWidthDp != 0 || o.SmallestScreenWidthDp != 0 {
		if c.SmallestScreenWidthDp != o.SmallestScreenWidthDp {
			if c.SmallestScreenWidthDp == 0 {
				return false
			}
			if o.SmallestScreenWidthDp == 0 {
				return true
			}
		}
	}

	// screen size dp
	if c.ScreenWidthDp != 0 || o.ScreenWidthDp != 0 ||
		c.ScreenHeightDp != 0 || o.ScreenHeightDp != 0 {
		if c.ScreenWidthDp != o.ScreenWidthDp {
			if c.ScreenWidthDp == 0 {
				return false
			}
			if o.ScreenWidthDp == 0 {
				return true
			}
		}
		if c.ScreenHeightDp != o.ScreenHeightDp {
			if c.ScreenHeightDp == 0 {
				return false
			}
			if o.ScreenHeightDp == 0 {
				return true
			}
		}
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if ((c.ScreenLayout ^ o.ScreenLayout) & MaskScreenSize) != 0 {
			if (c.ScreenLayout & MaskScreenSize) == 0 {
				return false
			}
			if (o.ScreenLayout & MaskScreenSize) == 0 {
				return true
			}
		}
		if ((c.ScreenLayout ^ o.ScreenLayout) & MaskScreenLong) != 0 {
			if (c.ScreenLayout & MaskScreenLong) == 0 {
				return false
			}
			if (o.ScreenLayout & MaskScreenLong) == 0 {
				return true
			}
		}
	}

	// orientation
	if c.Orientation != o.Orientation {
		if c.Orientation == 0 {
			return false
		}
		if o.Orientation == 0 {
			return true
		}
	}

	// uimode
	if c.UIMode != 0 || o.UIMode != 0 {
		diff := c.UIMode ^ o.UIMode
		if (diff & MaskUIModeType) != 0 {
			if (c.UIMode & MaskUIModeType) == 0 {
				return false
			}
			if (o.UIMode & MaskUIModeType) == 0 {
				return true
			}
		}
		if (diff & MaskUIModeNight) != 0 {
			if (c.UIMode & MaskUIModeNight) == 0 {
				return false
			}
			if (o.UIMode & MaskUIModeNight) == 0 {
				return true
			}
		}
	}

	// touchscreen
	if c.Touchscreen != o.Touchscreen {
		if c.Touchscreen == 0 {
			return false
		}
		if o.Touchscreen == 0 {
			return true
		}
	}

	// input
	if c.InputFlags != 0 || o.InputFlags != 0 {
		myKeysHidden := c.InputFlags & MaskKeysHidden
		oKeysHidden := o.InputFlags & MaskKeysHidden
		if (myKeysHidden ^ oKeysHidden) != 0 {
			if myKeysHidden == 0 {
				return false
			}
			if oKeysHidden == 0 {
				return true
			}
		}
		myNavHidden := c.InputFlags & MaskNavHidden
		oNavHidden := o.InputFlags & MaskNavHidden
		if (myNavHidden ^ oNavHidden) != 0 {
			if myNavHidden == 0 {
				return false
			}
			if oNavHidden == 0 {
				return true
			}
		}
	}

	if c.Keyboard != o.Keyboard {
		if c.Keyboard == 0 {
			return false
		}
		if o.Keyboard == 0 {
			return true
		}
	}

	if c.Navigation != o.Navigation {
		if c.Navigation == 0 {
			return false
		}
		if o.Navigation == 0 {
			return true
		}
	}

	// screen size
	if c.ScreenWidth != 0 || o.ScreenWidth != 0 ||
		c.ScreenHeight != 0 || o.ScreenHeight != 0 {
		if c.ScreenWidth != o.ScreenWidth {
			if c.ScreenWidth == 0 {
				return false
			}
			if o.ScreenWidth == 0 {
				return true
			}
		}
		if c.ScreenHeight != o.ScreenHeight {
			if c.ScreenHeight == 0 {
				return false
			}
			if o.ScreenHeight == 0 {
				return true
			}
		}
	}

	//version
	if c.SDKVersion != o.SDKVersion {
		if c.SDKVersion == 0 {
			return false
		}
		if o.SDKVersion == 0 {
			return true
		}
	}
	if c.MinorVersion != o.MinorVersion {
		if c.MinorVersion == 0 {
			return false
		}
		if o.MinorVersion == 0 {
			return true
		}
	}

	return false
}
