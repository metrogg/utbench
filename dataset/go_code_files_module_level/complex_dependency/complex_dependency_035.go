func (d *InstanceDiff) Same(d2 *InstanceDiff) (bool, string) {
	// we can safely compare the pointers without a lock
	switch {
	case d == nil && d2 == nil:
		return true, ""
	case d == nil || d2 == nil:
		return false, "one nil"
	case d == d2:
		return true, ""
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// If we're going from requiring new to NOT requiring new, then we have
	// to see if all required news were computed. If so, it is allowed since
	// computed may also mean "same value and therefore not new".
	oldNew := d.requiresNew()
	newNew := d2.RequiresNew()
	if oldNew && !newNew {
		oldNew = false

		// This section builds a list of ignorable attributes for requiresNew
		// by removing off any elements of collections going to zero elements.
		// For collections going to zero, they may not exist at all in the
		// new diff (and hence RequiresNew == false).
		ignoreAttrs := make(map[string]struct{})
		for k, diffOld := range d.Attributes {
			if !strings.HasSuffix(k, ".%") && !strings.HasSuffix(k, ".#") {
				continue
			}

			// This case is in here as a protection measure. The bug that this
			// code originally fixed (GH-11349) didn't have to deal with computed
			// so I'm not 100% sure what the correct behavior is. Best to leave
			// the old behavior.
			if diffOld.NewComputed {
				continue
			}

			// We're looking for the case a map goes to exactly 0.
			if diffOld.New != "0" {
				continue
			}

			// Found it! Ignore all of these. The prefix here is stripping
			// off the "%" so it is just "k."
			prefix := k[:len(k)-1]
			for k2, _ := range d.Attributes {
				if strings.HasPrefix(k2, prefix) {
					ignoreAttrs[k2] = struct{}{}
				}
			}
		}

		for k, rd := range d.Attributes {
			if _, ok := ignoreAttrs[k]; ok {
				continue
			}

			// If the field is requires new and NOT computed, then what
			// we have is a diff mismatch for sure. We set that the old
			// diff does REQUIRE a ForceNew.
			if rd != nil && rd.RequiresNew && !rd.NewComputed {
				oldNew = true
				break
			}
		}
	}

	if oldNew != newNew {
		return false, fmt.Sprintf(
			"diff RequiresNew; old: %t, new: %t", oldNew, newNew)
	}

	// Verify that destroy matches. The second boolean here allows us to
	// have mismatching Destroy if we're moving from RequiresNew true
	// to false above. Therefore, the second boolean will only pass if
	// we're moving from Destroy: true to false as well.
	if d.Destroy != d2.GetDestroy() && d.requiresNew() == oldNew {
		return false, fmt.Sprintf(
			"diff: Destroy; old: %t, new: %t", d.Destroy, d2.GetDestroy())
	}

	// Go through the old diff and make sure the new diff has all the
	// same attributes. To start, build up the check map to be all the keys.
	checkOld := make(map[string]struct{})
	checkNew := make(map[string]struct{})
	for k, _ := range d.Attributes {
		checkOld[k] = struct{}{}
	}
	for k, _ := range d2.CopyAttributes() {
		checkNew[k] = struct{}{}
	}

	// Make an ordered list so we are sure the approximated hashes are left
	// to process at the end of the loop
	keys := make([]string, 0, len(d.Attributes))
	for k, _ := range d.Attributes {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	for _, k := range keys {
		diffOld := d.Attributes[k]

		if _, ok := checkOld[k]; !ok {
			// We're not checking this key for whatever reason (see where
			// check is modified).
			continue
		}

		// Remove this key since we'll never hit it again
		delete(checkOld, k)
		delete(checkNew, k)

		_, ok := d2.GetAttribute(k)
		if !ok {
			// If there's no new attribute, and the old diff expected the attribute
			// to be removed, that's just fine.
			if diffOld.NewRemoved {
				continue
			}

			// If the last diff was a computed value then the absense of
			// that value is allowed since it may mean the value ended up
			// being the same.
			if diffOld.NewComputed {
				ok = true
			}

			// No exact match, but maybe this is a set containing computed
			// values. So check if there is an approximate hash in the key
			// and if so, try to match the key.
			if strings.Contains(k, "~") {
				parts := strings.Split(k, ".")
				parts2 := append([]string(nil), parts...)

				re := regexp.MustCompile(`^~\d+$`)
				for i, part := range parts {
					if re.MatchString(part) {
						// we're going to consider this the base of a
						// computed hash, and remove all longer matching fields
						ok = true

						parts2[i] = `\d+`
						parts2 = parts2[:i+1]
						break
					}
				}

				re, err := regexp.Compile("^" + strings.Join(parts2, `\.`))
				if err != nil {
					return false, fmt.Sprintf("regexp failed to compile; err: %#v", err)
				}

				for k2, _ := range checkNew {
					if re.MatchString(k2) {
						delete(checkNew, k2)
					}
				}
			}

			// This is a little tricky, but when a diff contains a computed
			// list, set, or map that can only be interpolated after the apply
			// command has created the dependent resources, it could turn out
			// that the result is actually the same as the existing state which
			// would remove the key from the diff.
			if diffOld.NewComputed && (strings.HasSuffix(k, ".#") || strings.HasSuffix(k, ".%")) {
				ok = true
			}

			// Similarly, in a RequiresNew scenario, a list that shows up in the plan
			// diff can disappear from the apply diff, which is calculated from an
			// empty state.
			if d.requiresNew() && (strings.HasSuffix(k, ".#") || strings.HasSuffix(k, ".%")) {
				ok = true
			}

			if !ok {
				return false, fmt.Sprintf("attribute mismatch: %s", k)
			}
		}

		// search for the suffix of the base of a [computed] map, list or set.
		match := multiVal.FindStringSubmatch(k)

		if diffOld.NewComputed && len(match) == 2 {
			matchLen := len(match[1])

			// This is a computed list, set, or map, so remove any keys with
			// this prefix from the check list.
			kprefix := k[:len(k)-matchLen]
			for k2, _ := range checkOld {
				if strings.HasPrefix(k2, kprefix) {
					delete(checkOld, k2)
				}
			}
			for k2, _ := range checkNew {
				if strings.HasPrefix(k2, kprefix) {
					delete(checkNew, k2)
				}
			}
		}

		// We don't compare the values because we can't currently actually
		// guarantee to generate the same value two two diffs created from
		// the same state+config: we have some pesky interpolation functions
		// that do not behave as pure functions (uuid, timestamp) and so they
		// can be different each time a diff is produced.
		// FIXME: Re-organize our config handling so that we don't re-evaluate
		// expressions when we produce a second comparison diff during
		// apply (for EvalCompareDiff).
	}

	// Check for leftover attributes
	if len(checkNew) > 0 {
		extras := make([]string, 0, len(checkNew))
		for attr, _ := range checkNew {
			extras = append(extras, attr)
		}
		return false,
			fmt.Sprintf("extra attributes: %s", strings.Join(extras, ", "))
	}

	return true, ""
}
