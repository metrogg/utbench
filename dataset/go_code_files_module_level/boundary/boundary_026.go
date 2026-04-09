func (s *Schema) validate(v interface{}) error {
	if s.Always != nil {
		if !*s.Always {
			return validationError("", "always fail")
		}
		return nil
	}

	if s.Ref != nil {
		if err := s.Ref.validate(v); err != nil {
			finishSchemaContext(err, s.Ref)
			var refURL string
			if s.URL == s.Ref.URL {
				refURL = s.Ref.Ptr
			} else {
				refURL = s.Ref.URL + s.Ref.Ptr
			}
			return validationError("$ref", "doesn't validate with %q", refURL).add(err)
		}

		// All other properties in a "$ref" object MUST be ignored
		return nil
	}

	if len(s.Types) > 0 {
		vType := jsonType(v)
		matched := false
		for _, t := range s.Types {
			if vType == t {
				matched = true
				break
			} else if t == "integer" && vType == "number" {
				if _, ok := new(big.Int).SetString(fmt.Sprint(v), 10); ok {
					matched = true
					break
				}
			}
		}
		if !matched {
			return validationError("type", "expected %s, but got %s", strings.Join(s.Types, " or "), vType)
		}
	}

	if len(s.Constant) > 0 {
		if !equals(v, s.Constant[0]) {
			switch jsonType(s.Constant[0]) {
			case "object", "array":
				return validationError("const", "const failed")
			default:
				return validationError("const", "value must be %#v", s.Constant[0])
			}
		}
	}

	if len(s.Enum) > 0 {
		matched := false
		for _, item := range s.Enum {
			if equals(v, item) {
				matched = true
				break
			}
		}
		if !matched {
			return validationError("enum", s.enumError)
		}
	}

	if s.format != nil && !s.format(v) {
		return validationError("format", "%q is not valid %q", v, s.Format)
	}

	if s.Not != nil && s.Not.validate(v) == nil {
		return validationError("not", "not failed")
	}

	for i, sch := range s.AllOf {
		if err := sch.validate(v); err != nil {
			return validationError("allOf/"+strconv.Itoa(i), "allOf failed").add(err)
		}
	}

	if len(s.AnyOf) > 0 {
		matched := false
		var causes []error
		for i, sch := range s.AnyOf {
			if err := sch.validate(v); err == nil {
				matched = true
				break
			} else {
				causes = append(causes, addContext("", strconv.Itoa(i), err))
			}
		}
		if !matched {
			return validationError("anyOf", "anyOf failed").add(causes...)
		}
	}

	if len(s.OneOf) > 0 {
		matched := -1
		var causes []error
		for i, sch := range s.OneOf {
			if err := sch.validate(v); err == nil {
				if matched == -1 {
					matched = i
				} else {
					return validationError("oneOf", "valid against schemas at indexes %d and %d", matched, i)
				}
			} else {
				causes = append(causes, addContext("", strconv.Itoa(i), err))
			}
		}
		if matched == -1 {
			return validationError("oneOf", "oneOf failed").add(causes...)
		}
	}

	if s.If != nil {
		if s.If.validate(v) == nil {
			if s.Then != nil {
				if err := s.Then.validate(v); err != nil {
					return validationError("then", "if-then failed").add(err)
				}
			}
		} else {
			if s.Else != nil {
				if err := s.Else.validate(v); err != nil {
					return validationError("else", "if-else failed").add(err)
				}
			}
		}
	}

	switch v := v.(type) {
	case map[string]interface{}:
		if s.MinProperties != -1 && len(v) < s.MinProperties {
			return validationError("minProperties", "minimum %d properties allowed, but found %d properties", s.MinProperties, len(v))
		}
		if s.MaxProperties != -1 && len(v) > s.MaxProperties {
			return validationError("maxProperties", "maximum %d properties allowed, but found %d properties", s.MaxProperties, len(v))
		}
		if len(s.Required) > 0 {
			var missing []string
			for _, pname := range s.Required {
				if _, ok := v[pname]; !ok {
					missing = append(missing, strconv.Quote(pname))
				}
			}
			if len(missing) > 0 {
				return validationError("required", "missing properties: %s", strings.Join(missing, ", "))
			}
		}

		var additionalProps map[string]struct{}
		if s.AdditionalProperties != nil {
			additionalProps = make(map[string]struct{}, len(v))
			for pname := range v {
				additionalProps[pname] = struct{}{}
			}
		}

		if len(s.Properties) > 0 {
			for pname, pschema := range s.Properties {
				if pvalue, ok := v[pname]; ok {
					delete(additionalProps, pname)
					if err := pschema.validate(pvalue); err != nil {
						return addContext(escape(pname), "properties/"+escape(pname), err)
					}
				}
			}
		}

		if s.PropertyNames != nil {
			for pname := range v {
				if err := s.PropertyNames.validate(pname); err != nil {
					return addContext(escape(pname), "propertyNames", err)
				}
			}
		}

		if s.RegexProperties {
			for pname := range v {
				if !isRegex(pname) {
					return validationError("", "patternProperty %q is not valid regex", pname)
				}
			}
		}
		for pattern, pschema := range s.PatternProperties {
			for pname, pvalue := range v {
				if pattern.MatchString(pname) {
					delete(additionalProps, pname)
					if err := pschema.validate(pvalue); err != nil {
						return addContext(escape(pname), "patternProperties/"+escape(pattern.String()), err)
					}
				}
			}
		}
		if s.AdditionalProperties != nil {
			if _, ok := s.AdditionalProperties.(bool); ok {
				if len(additionalProps) != 0 {
					pnames := make([]string, 0, len(additionalProps))
					for pname := range additionalProps {
						pnames = append(pnames, strconv.Quote(pname))
					}
					return validationError("additionalProperties", "additionalProperties %s not allowed", strings.Join(pnames, ", "))
				}
			} else {
				schema := s.AdditionalProperties.(*Schema)
				for pname := range additionalProps {
					if pvalue, ok := v[pname]; ok {
						if err := schema.validate(pvalue); err != nil {
							return addContext(escape(pname), "additionalProperties", err)
						}
					}
				}
			}
		}
		for dname, dvalue := range s.Dependencies {
			if _, ok := v[dname]; ok {
				switch dvalue := dvalue.(type) {
				case *Schema:
					if err := dvalue.validate(v); err != nil {
						return addContext("", "dependencies/"+escape(dname), err)
					}
				case []string:
					for i, pname := range dvalue {
						if _, ok := v[pname]; !ok {
							return validationError("dependencies/"+escape(dname)+"/"+strconv.Itoa(i), "property %q is required, if %q property exists", pname, dname)
						}
					}
				}
			}
		}

	case []interface{}:
		if s.MinItems != -1 && len(v) < s.MinItems {
			return validationError("minItems", "minimum %d items allowed, but found %d items", s.MinItems, len(v))
		}
		if s.MaxItems != -1 && len(v) > s.MaxItems {
			return validationError("maxItems", "maximum %d items allowed, but found %d items", s.MaxItems, len(v))
		}
		if s.UniqueItems {
			for i := 1; i < len(v); i++ {
				for j := 0; j < i; j++ {
					if equals(v[i], v[j]) {
						return validationError("uniqueItems", "items at index %d and %d are equal", j, i)
					}
				}
			}
		}
		switch items := s.Items.(type) {
		case *Schema:
			for i, item := range v {
				if err := items.validate(item); err != nil {
					return addContext(strconv.Itoa(i), "items", err)
				}
			}
		case []*Schema:
			if additionalItems, ok := s.AdditionalItems.(bool); ok {
				if !additionalItems && len(v) > len(items) {
					return validationError("additionalItems", "only %d items are allowed, but found %d items", len(items), len(v))
				}
			}
			for i, item := range v {
				if i < len(items) {
					if err := items[i].validate(item); err != nil {
						return addContext(strconv.Itoa(i), "items/"+strconv.Itoa(i), err)
					}
				} else if sch, ok := s.AdditionalItems.(*Schema); ok {
					if err := sch.validate(item); err != nil {
						return addContext(strconv.Itoa(i), "additionalItems", err)
					}
				} else {
					break
				}
			}
		}
		if s.Contains != nil {
			matched := false
			var causes []error
			for i, item := range v {
				if err := s.Contains.validate(item); err != nil {
					causes = append(causes, addContext(strconv.Itoa(i), "", err))
				} else {
					matched = true
					break
				}
			}
			if !matched {
				return validationError("contains", "contains failed").add(causes...)
			}
		}

	case string:
		if s.MinLength != -1 || s.MaxLength != -1 {
			length := utf8.RuneCount([]byte(v))
			if s.MinLength != -1 && length < s.MinLength {
				return validationError("minLength", "length must be >= %d, but got %d", s.MinLength, length)
			}
			if s.MaxLength != -1 && length > s.MaxLength {
				return validationError("maxLength", "length must be <= %d, but got %d", s.MaxLength, length)
			}
		}
		if s.Pattern != nil && !s.Pattern.MatchString(v) {
			return validationError("pattern", "does not match pattern %q", s.Pattern)
		}

		decoded := s.ContentEncoding == ""
		var content []byte
		if s.decoder != nil {
			b, err := s.decoder(v)
			if err != nil {
				return validationError("contentEncoding", "%q is not %s encoded", v, s.ContentEncoding)
			}
			content, decoded = b, true
		}
		if decoded && s.mediaType != nil {
			if s.decoder == nil {
				content = []byte(v)
			}
			if err := s.mediaType(content); err != nil {
				return validationError("contentMediaType", "value is not of mediatype %q", s.ContentMediaType)
			}
		}

	case json.Number, float64, int, int32, int64:
		num, _ := new(big.Float).SetString(fmt.Sprint(v))
		if s.Minimum != nil && num.Cmp(s.Minimum) < 0 {
			return validationError("minimum", "must be >= %v but found %v", s.Minimum, v)
		}
		if s.ExclusiveMinimum != nil && num.Cmp(s.ExclusiveMinimum) <= 0 {
			return validationError("exclusiveMinimum", "must be > %v but found %v", s.ExclusiveMinimum, v)
		}
		if s.Maximum != nil && num.Cmp(s.Maximum) > 0 {
			return validationError("maximum", "must be <= %v but found %v", s.Maximum, v)
		}
		if s.ExclusiveMaximum != nil && num.Cmp(s.ExclusiveMaximum) >= 0 {
			return validationError("exclusiveMaximum", "must be < %v but found %v", s.ExclusiveMaximum, v)
		}
		if s.MultipleOf != nil {
			if q := new(big.Float).Quo(num, s.MultipleOf); !q.IsInt() {
				return validationError("multipleOf", "%v not multipleOf %v", v, s.MultipleOf)
			}
		}
	}

	return nil
}
