func (c *Config) Validate() tfdiags.Diagnostics {
	if c == nil {
		return nil
	}

	var diags tfdiags.Diagnostics

	for _, k := range c.unknownKeys {
		diags = diags.Append(
			fmt.Errorf("Unknown root level key: %s", k),
		)
	}

	// Validate the Terraform config
	if tf := c.Terraform; tf != nil {
		errs := c.Terraform.Validate()
		for _, err := range errs {
			diags = diags.Append(err)
		}
	}

	vars := c.InterpolatedVariables()
	varMap := make(map[string]*Variable)
	for _, v := range c.Variables {
		if _, ok := varMap[v.Name]; ok {
			diags = diags.Append(fmt.Errorf(
				"Variable '%s': duplicate found. Variable names must be unique.",
				v.Name,
			))
		}

		varMap[v.Name] = v
	}

	for k, _ := range varMap {
		if !NameRegexp.MatchString(k) {
			diags = diags.Append(fmt.Errorf(
				"variable %q: variable name must match regular expression %s",
				k, NameRegexp,
			))
		}
	}

	for _, v := range c.Variables {
		if v.Type() == VariableTypeUnknown {
			diags = diags.Append(fmt.Errorf(
				"Variable '%s': must be a string or a map",
				v.Name,
			))
			continue
		}

		interp := false
		fn := func(n ast.Node) (interface{}, error) {
			// LiteralNode is a literal string (outside of a ${ ... } sequence).
			// interpolationWalker skips most of these. but in particular it
			// visits those that have escaped sequences (like $${foo}) as a
			// signal that *some* processing is required on this string. For
			// our purposes here though, this is fine and not an interpolation.
			if _, ok := n.(*ast.LiteralNode); !ok {
				interp = true
			}
			return "", nil
		}

		w := &interpolationWalker{F: fn}
		if v.Default != nil {
			if err := reflectwalk.Walk(v.Default, w); err == nil {
				if interp {
					diags = diags.Append(fmt.Errorf(
						"variable %q: default may not contain interpolations",
						v.Name,
					))
				}
			}
		}
	}

	// Check for references to user variables that do not actually
	// exist and record those errors.
	for source, vs := range vars {
		for _, v := range vs {
			uv, ok := v.(*UserVariable)
			if !ok {
				continue
			}

			if _, ok := varMap[uv.Name]; !ok {
				diags = diags.Append(fmt.Errorf(
					"%s: unknown variable referenced: '%s'; define it with a 'variable' block",
					source,
					uv.Name,
				))
			}
		}
	}

	// Check that all count variables are valid.
	for source, vs := range vars {
		for _, rawV := range vs {
			switch v := rawV.(type) {
			case *CountVariable:
				if v.Type == CountValueInvalid {
					diags = diags.Append(fmt.Errorf(
						"%s: invalid count variable: %s",
						source,
						v.FullKey(),
					))
				}
			case *PathVariable:
				if v.Type == PathValueInvalid {
					diags = diags.Append(fmt.Errorf(
						"%s: invalid path variable: %s",
						source,
						v.FullKey(),
					))
				}
			}
		}
	}

	// Check that providers aren't declared multiple times and that their
	// version constraints, where present, are syntactically valid.
	providerSet := make(map[string]bool)
	for _, p := range c.ProviderConfigs {
		name := p.FullName()
		if _, ok := providerSet[name]; ok {
			diags = diags.Append(fmt.Errorf(
				"provider.%s: multiple configurations present; only one configuration is allowed per provider",
				name,
			))
			continue
		}

		if p.Version != "" {
			_, err := discovery.ConstraintStr(p.Version).Parse()
			if err != nil {
				diags = diags.Append(&hcl2.Diagnostic{
					Severity: hcl2.DiagError,
					Summary:  "Invalid provider version constraint",
					Detail: fmt.Sprintf(
						"The value %q given for provider.%s is not a valid version constraint.",
						p.Version, name,
					),
					// TODO: include a "Subject" source reference in here,
					// once the config loader is able to retain source
					// location information.
				})
			}
		}

		providerSet[name] = true
	}

	// Check that all references to modules are valid
	modules := make(map[string]*Module)
	dupped := make(map[string]struct{})
	for _, m := range c.Modules {
		// Check for duplicates
		if _, ok := modules[m.Id()]; ok {
			if _, ok := dupped[m.Id()]; !ok {
				dupped[m.Id()] = struct{}{}

				diags = diags.Append(fmt.Errorf(
					"module %q: module repeated multiple times",
					m.Id(),
				))
			}

			// Already seen this module, just skip it
			continue
		}

		modules[m.Id()] = m

		// Check that the source has no interpolations
		rc, err := NewRawConfig(map[string]interface{}{
			"root": m.Source,
		})
		if err != nil {
			diags = diags.Append(fmt.Errorf(
				"module %q: module source error: %s",
				m.Id(), err,
			))
		} else if len(rc.Interpolations) > 0 {
			diags = diags.Append(fmt.Errorf(
				"module %q: module source cannot contain interpolations",
				m.Id(),
			))
		}

		// Check that the name matches our regexp
		if !NameRegexp.Match([]byte(m.Name)) {
			diags = diags.Append(fmt.Errorf(
				"module %q: module name must be a letter or underscore followed by only letters, numbers, dashes, and underscores",
				m.Id(),
			))
		}

		// Check that the configuration can all be strings, lists or maps
		raw := make(map[string]interface{})
		for k, v := range m.RawConfig.Raw {
			var strVal string
			if err := hilmapstructure.WeakDecode(v, &strVal); err == nil {
				raw[k] = strVal
				continue
			}

			var mapVal map[string]interface{}
			if err := hilmapstructure.WeakDecode(v, &mapVal); err == nil {
				raw[k] = mapVal
				continue
			}

			var sliceVal []interface{}
			if err := hilmapstructure.WeakDecode(v, &sliceVal); err == nil {
				raw[k] = sliceVal
				continue
			}

			diags = diags.Append(fmt.Errorf(
				"module %q: argument %s must have a string, list, or map value",
				m.Id(), k,
			))
		}

		// Check for invalid count variables
		for _, v := range m.RawConfig.Variables {
			switch v.(type) {
			case *CountVariable:
				diags = diags.Append(fmt.Errorf(
					"module %q: count variables are only valid within resources",
					m.Name,
				))
			case *SelfVariable:
				diags = diags.Append(fmt.Errorf(
					"module %q: self variables are only valid within resources",
					m.Name,
				))
			}
		}

		// Update the raw configuration to only contain the string values
		m.RawConfig, err = NewRawConfig(raw)
		if err != nil {
			diags = diags.Append(fmt.Errorf(
				"%s: can't initialize configuration: %s",
				m.Id(), err,
			))
		}

		// check that all named providers actually exist
		for _, p := range m.Providers {
			if !providerSet[p] {
				diags = diags.Append(fmt.Errorf(
					"module %q: cannot pass non-existent provider %q",
					m.Name, p,
				))
			}
		}

	}
	dupped = nil

	// Check that all variables for modules reference modules that
	// exist.
	for source, vs := range vars {
		for _, v := range vs {
			mv, ok := v.(*ModuleVariable)
			if !ok {
				continue
			}

			if _, ok := modules[mv.Name]; !ok {
				diags = diags.Append(fmt.Errorf(
					"%s: unknown module referenced: %s",
					source, mv.Name,
				))
			}
		}
	}

	// Check that all references to resources are valid
	resources := make(map[string]*Resource)
	dupped = make(map[string]struct{})
	for _, r := range c.Resources {
		if _, ok := resources[r.Id()]; ok {
			if _, ok := dupped[r.Id()]; !ok {
				dupped[r.Id()] = struct{}{}

				diags = diags.Append(fmt.Errorf(
					"%s: resource repeated multiple times",
					r.Id(),
				))
			}
		}

		resources[r.Id()] = r
	}
	dupped = nil

	// Validate resources
	for n, r := range resources {
		// Verify count variables
		for _, v := range r.RawCount.Variables {
			switch v.(type) {
			case *CountVariable:
				diags = diags.Append(fmt.Errorf(
					"%s: resource count can't reference count variable: %s",
					n, v.FullKey(),
				))
			case *SimpleVariable:
				diags = diags.Append(fmt.Errorf(
					"%s: resource count can't reference variable: %s",
					n, v.FullKey(),
				))

			// Good
			case *ModuleVariable:
			case *ResourceVariable:
			case *TerraformVariable:
			case *UserVariable:
			case *LocalVariable:

			default:
				diags = diags.Append(fmt.Errorf(
					"Internal error. Unknown type in count var in %s: %T",
					n, v,
				))
			}
		}

		if !r.RawCount.couldBeInteger() {
			diags = diags.Append(fmt.Errorf(
				"%s: resource count must be an integer", n,
			))
		}
		r.RawCount.init()

		// Validate DependsOn
		for _, err := range c.validateDependsOn(n, r.DependsOn, resources, modules) {
			diags = diags.Append(err)
		}

		// Verify provisioners
		for _, p := range r.Provisioners {
			// This validation checks that there are no splat variables
			// referencing ourself. This currently is not allowed.

			for _, v := range p.ConnInfo.Variables {
				rv, ok := v.(*ResourceVariable)
				if !ok {
					continue
				}

				if rv.Multi && rv.Index == -1 && rv.Type == r.Type && rv.Name == r.Name {
					diags = diags.Append(fmt.Errorf(
						"%s: connection info cannot contain splat variable referencing itself",
						n,
					))
					break
				}
			}

			for _, v := range p.RawConfig.Variables {
				rv, ok := v.(*ResourceVariable)
				if !ok {
					continue
				}

				if rv.Multi && rv.Index == -1 && rv.Type == r.Type && rv.Name == r.Name {
					diags = diags.Append(fmt.Errorf(
						"%s: connection info cannot contain splat variable referencing itself",
						n,
					))
					break
				}
			}

			// Check for invalid when/onFailure values, though this should be
			// picked up by the loader we check here just in case.
			if p.When == ProvisionerWhenInvalid {
				diags = diags.Append(fmt.Errorf(
					"%s: provisioner 'when' value is invalid", n,
				))
			}
			if p.OnFailure == ProvisionerOnFailureInvalid {
				diags = diags.Append(fmt.Errorf(
					"%s: provisioner 'on_failure' value is invalid", n,
				))
			}
		}

		// Verify ignore_changes contains valid entries
		for _, v := range r.Lifecycle.IgnoreChanges {
			if strings.Contains(v, "*") && v != "*" {
				diags = diags.Append(fmt.Errorf(
					"%s: ignore_changes does not support using a partial string together with a wildcard: %s",
					n, v,
				))
			}
		}

		// Verify ignore_changes has no interpolations
		rc, err := NewRawConfig(map[string]interface{}{
			"root": r.Lifecycle.IgnoreChanges,
		})
		if err != nil {
			diags = diags.Append(fmt.Errorf(
				"%s: lifecycle ignore_changes error: %s",
				n, err,
			))
		} else if len(rc.Interpolations) > 0 {
			diags = diags.Append(fmt.Errorf(
				"%s: lifecycle ignore_changes cannot contain interpolations",
				n,
			))
		}

		// If it is a data source then it can't have provisioners
		if r.Mode == DataResourceMode {
			if _, ok := r.RawConfig.Raw["provisioner"]; ok {
				diags = diags.Append(fmt.Errorf(
					"%s: data sources cannot have provisioners",
					n,
				))
			}
		}
	}

	for source, vs := range vars {
		for _, v := range vs {
			rv, ok := v.(*ResourceVariable)
			if !ok {
				continue
			}

			id := rv.ResourceId()
			if _, ok := resources[id]; !ok {
				diags = diags.Append(fmt.Errorf(
					"%s: unknown resource '%s' referenced in variable %s",
					source,
					id,
					rv.FullKey(),
				))
				continue
			}
		}
	}

	// Check that all locals are valid
	{
		found := make(map[string]struct{})
		for _, l := range c.Locals {
			if _, ok := found[l.Name]; ok {
				diags = diags.Append(fmt.Errorf(
					"%s: duplicate local. local value names must be unique",
					l.Name,
				))
				continue
			}
			found[l.Name] = struct{}{}

			for _, v := range l.RawConfig.Variables {
				if _, ok := v.(*CountVariable); ok {
					diags = diags.Append(fmt.Errorf(
						"local %s: count variables are only valid within resources", l.Name,
					))
				}
			}
		}
	}

	// Check that all outputs are valid
	{
		found := make(map[string]struct{})
		for _, o := range c.Outputs {
			// Verify the output is new
			if _, ok := found[o.Name]; ok {
				diags = diags.Append(fmt.Errorf(
					"output %q: an output of this name was already defined",
					o.Name,
				))
				continue
			}
			found[o.Name] = struct{}{}

			var invalidKeys []string
			valueKeyFound := false
			for k := range o.RawConfig.Raw {
				if k == "value" {
					valueKeyFound = true
					continue
				}
				if k == "sensitive" {
					if sensitive, ok := o.RawConfig.config[k].(bool); ok {
						if sensitive {
							o.Sensitive = true
						}
						continue
					}

					diags = diags.Append(fmt.Errorf(
						"output %q: value for 'sensitive' must be boolean",
						o.Name,
					))
					continue
				}
				if k == "description" {
					if desc, ok := o.RawConfig.config[k].(string); ok {
						o.Description = desc
						continue
					}

					diags = diags.Append(fmt.Errorf(
						"output %q: value for 'description' must be string",
						o.Name,
					))
					continue
				}
				invalidKeys = append(invalidKeys, k)
			}
			if len(invalidKeys) > 0 {
				diags = diags.Append(fmt.Errorf(
					"output %q: invalid keys: %s",
					o.Name, strings.Join(invalidKeys, ", "),
				))
			}
			if !valueKeyFound {
				diags = diags.Append(fmt.Errorf(
					"output %q: missing required 'value' argument", o.Name,
				))
			}

			for _, v := range o.RawConfig.Variables {
				if _, ok := v.(*CountVariable); ok {
					diags = diags.Append(fmt.Errorf(
						"output %q: count variables are only valid within resources",
						o.Name,
					))
				}
			}

			// Detect a common mistake of using a "count"ed resource in
			// an output value without using the splat or index form.
			// Prior to 0.11 this error was silently ignored, but outputs
			// now have their errors checked like all other contexts.
			//
			// TODO: Remove this in 0.12.
			for _, v := range o.RawConfig.Variables {
				rv, ok := v.(*ResourceVariable)
				if !ok {
					continue
				}

				// If the variable seems to be treating the referenced
				// resource as a singleton (no count specified) then
				// we'll check to make sure it is indeed a singleton.
				// It's a warning if not.

				if rv.Multi || rv.Index != 0 {
					// This reference is treating the resource as a
					// multi-resource, so the warning doesn't apply.
					continue
				}

				for _, r := range c.Resources {
					if r.Id() != rv.ResourceId() {
						continue
					}

					// We test specifically for the raw string "1" here
					// because we _do_ want to generate this warning if
					// the user has provided an expression that happens
					// to return 1 right now, to catch situations where
					// a count might dynamically be set to something
					// other than 1 and thus splat syntax is still needed
					// to be safe.
					if r.RawCount != nil && r.RawCount.Raw != nil && r.RawCount.Raw["count"] != "1" && rv.Field != "count" {
						diags = diags.Append(tfdiags.SimpleWarning(fmt.Sprintf(
							"output %q: must use splat syntax to access %s attribute %q, because it has \"count\" set; use %s.*.%s to obtain a list of the attributes across all instances",
							o.Name,
							r.Id(), rv.Field,
							r.Id(), rv.Field,
						)))
					}
				}
			}
		}
	}

	// Validate the self variable
	for source, rc := range c.rawConfigs() {
		// Ignore provisioners. This is a pretty brittle way to do this,
		// but better than also repeating all the resources.
		if strings.Contains(source, "provision") {
			continue
		}

		for _, v := range rc.Variables {
			if _, ok := v.(*SelfVariable); ok {
				diags = diags.Append(fmt.Errorf(
					"%s: cannot contain self-reference %s",
					source, v.FullKey(),
				))
			}
		}
	}

	return diags
}
