func NewSchema(in interface{}, context *compiler.Context) (*Schema, error) {
	errors := make([]error, 0)
	x := &Schema{}
	m, ok := compiler.UnpackMap(in)
	if !ok {
		message := fmt.Sprintf("has unexpected value: %+v (%T)", in, in)
		errors = append(errors, compiler.NewError(context, message))
	} else {
		allowedKeys := []string{"additionalProperties", "allOf", "anyOf", "default", "deprecated", "description", "discriminator", "enum", "example", "exclusiveMaximum", "exclusiveMinimum", "externalDocs", "format", "items", "maxItems", "maxLength", "maxProperties", "maximum", "minItems", "minLength", "minProperties", "minimum", "multipleOf", "not", "nullable", "oneOf", "pattern", "properties", "readOnly", "required", "title", "type", "uniqueItems", "writeOnly", "xml"}
		allowedPatterns := []*regexp.Regexp{pattern1}
		invalidKeys := compiler.InvalidKeysInMap(m, allowedKeys, allowedPatterns)
		if len(invalidKeys) > 0 {
			message := fmt.Sprintf("has invalid %s: %+v", compiler.PluralProperties(len(invalidKeys)), strings.Join(invalidKeys, ", "))
			errors = append(errors, compiler.NewError(context, message))
		}
		// bool nullable = 1;
		v1 := compiler.MapValueForKey(m, "nullable")
		if v1 != nil {
			x.Nullable, ok = v1.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for nullable: %+v (%T)", v1, v1)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// Discriminator discriminator = 2;
		v2 := compiler.MapValueForKey(m, "discriminator")
		if v2 != nil {
			var err error
			x.Discriminator, err = NewDiscriminator(v2, compiler.NewContext("discriminator", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// bool read_only = 3;
		v3 := compiler.MapValueForKey(m, "readOnly")
		if v3 != nil {
			x.ReadOnly, ok = v3.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for readOnly: %+v (%T)", v3, v3)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// bool write_only = 4;
		v4 := compiler.MapValueForKey(m, "writeOnly")
		if v4 != nil {
			x.WriteOnly, ok = v4.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for writeOnly: %+v (%T)", v4, v4)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// Xml xml = 5;
		v5 := compiler.MapValueForKey(m, "xml")
		if v5 != nil {
			var err error
			x.Xml, err = NewXml(v5, compiler.NewContext("xml", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// ExternalDocs external_docs = 6;
		v6 := compiler.MapValueForKey(m, "externalDocs")
		if v6 != nil {
			var err error
			x.ExternalDocs, err = NewExternalDocs(v6, compiler.NewContext("externalDocs", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// Any example = 7;
		v7 := compiler.MapValueForKey(m, "example")
		if v7 != nil {
			var err error
			x.Example, err = NewAny(v7, compiler.NewContext("example", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// bool deprecated = 8;
		v8 := compiler.MapValueForKey(m, "deprecated")
		if v8 != nil {
			x.Deprecated, ok = v8.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for deprecated: %+v (%T)", v8, v8)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// string title = 9;
		v9 := compiler.MapValueForKey(m, "title")
		if v9 != nil {
			x.Title, ok = v9.(string)
			if !ok {
				message := fmt.Sprintf("has unexpected value for title: %+v (%T)", v9, v9)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// float multiple_of = 10;
		v10 := compiler.MapValueForKey(m, "multipleOf")
		if v10 != nil {
			switch v10 := v10.(type) {
			case float64:
				x.MultipleOf = v10
			case float32:
				x.MultipleOf = float64(v10)
			case uint64:
				x.MultipleOf = float64(v10)
			case uint32:
				x.MultipleOf = float64(v10)
			case int64:
				x.MultipleOf = float64(v10)
			case int32:
				x.MultipleOf = float64(v10)
			case int:
				x.MultipleOf = float64(v10)
			default:
				message := fmt.Sprintf("has unexpected value for multipleOf: %+v (%T)", v10, v10)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// float maximum = 11;
		v11 := compiler.MapValueForKey(m, "maximum")
		if v11 != nil {
			switch v11 := v11.(type) {
			case float64:
				x.Maximum = v11
			case float32:
				x.Maximum = float64(v11)
			case uint64:
				x.Maximum = float64(v11)
			case uint32:
				x.Maximum = float64(v11)
			case int64:
				x.Maximum = float64(v11)
			case int32:
				x.Maximum = float64(v11)
			case int:
				x.Maximum = float64(v11)
			default:
				message := fmt.Sprintf("has unexpected value for maximum: %+v (%T)", v11, v11)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// bool exclusive_maximum = 12;
		v12 := compiler.MapValueForKey(m, "exclusiveMaximum")
		if v12 != nil {
			x.ExclusiveMaximum, ok = v12.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for exclusiveMaximum: %+v (%T)", v12, v12)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// float minimum = 13;
		v13 := compiler.MapValueForKey(m, "minimum")
		if v13 != nil {
			switch v13 := v13.(type) {
			case float64:
				x.Minimum = v13
			case float32:
				x.Minimum = float64(v13)
			case uint64:
				x.Minimum = float64(v13)
			case uint32:
				x.Minimum = float64(v13)
			case int64:
				x.Minimum = float64(v13)
			case int32:
				x.Minimum = float64(v13)
			case int:
				x.Minimum = float64(v13)
			default:
				message := fmt.Sprintf("has unexpected value for minimum: %+v (%T)", v13, v13)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// bool exclusive_minimum = 14;
		v14 := compiler.MapValueForKey(m, "exclusiveMinimum")
		if v14 != nil {
			x.ExclusiveMinimum, ok = v14.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for exclusiveMinimum: %+v (%T)", v14, v14)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 max_length = 15;
		v15 := compiler.MapValueForKey(m, "maxLength")
		if v15 != nil {
			t, ok := v15.(int)
			if ok {
				x.MaxLength = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for maxLength: %+v (%T)", v15, v15)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 min_length = 16;
		v16 := compiler.MapValueForKey(m, "minLength")
		if v16 != nil {
			t, ok := v16.(int)
			if ok {
				x.MinLength = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for minLength: %+v (%T)", v16, v16)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// string pattern = 17;
		v17 := compiler.MapValueForKey(m, "pattern")
		if v17 != nil {
			x.Pattern, ok = v17.(string)
			if !ok {
				message := fmt.Sprintf("has unexpected value for pattern: %+v (%T)", v17, v17)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 max_items = 18;
		v18 := compiler.MapValueForKey(m, "maxItems")
		if v18 != nil {
			t, ok := v18.(int)
			if ok {
				x.MaxItems = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for maxItems: %+v (%T)", v18, v18)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 min_items = 19;
		v19 := compiler.MapValueForKey(m, "minItems")
		if v19 != nil {
			t, ok := v19.(int)
			if ok {
				x.MinItems = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for minItems: %+v (%T)", v19, v19)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// bool unique_items = 20;
		v20 := compiler.MapValueForKey(m, "uniqueItems")
		if v20 != nil {
			x.UniqueItems, ok = v20.(bool)
			if !ok {
				message := fmt.Sprintf("has unexpected value for uniqueItems: %+v (%T)", v20, v20)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 max_properties = 21;
		v21 := compiler.MapValueForKey(m, "maxProperties")
		if v21 != nil {
			t, ok := v21.(int)
			if ok {
				x.MaxProperties = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for maxProperties: %+v (%T)", v21, v21)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 min_properties = 22;
		v22 := compiler.MapValueForKey(m, "minProperties")
		if v22 != nil {
			t, ok := v22.(int)
			if ok {
				x.MinProperties = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for minProperties: %+v (%T)", v22, v22)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// repeated string required = 23;
		v23 := compiler.MapValueForKey(m, "required")
		if v23 != nil {
			v, ok := v23.([]interface{})
			if ok {
				x.Required = compiler.ConvertInterfaceArrayToStringArray(v)
			} else {
				message := fmt.Sprintf("has unexpected value for required: %+v (%T)", v23, v23)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// repeated Any enum = 24;
		v24 := compiler.MapValueForKey(m, "enum")
		if v24 != nil {
			// repeated Any
			x.Enum = make([]*Any, 0)
			a, ok := v24.([]interface{})
			if ok {
				for _, item := range a {
					y, err := NewAny(item, compiler.NewContext("enum", context))
					if err != nil {
						errors = append(errors, err)
					}
					x.Enum = append(x.Enum, y)
				}
			}
		}
		// string type = 25;
		v25 := compiler.MapValueForKey(m, "type")
		if v25 != nil {
			x.Type, ok = v25.(string)
			if !ok {
				message := fmt.Sprintf("has unexpected value for type: %+v (%T)", v25, v25)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// repeated SchemaOrReference all_of = 26;
		v26 := compiler.MapValueForKey(m, "allOf")
		if v26 != nil {
			// repeated SchemaOrReference
			x.AllOf = make([]*SchemaOrReference, 0)
			a, ok := v26.([]interface{})
			if ok {
				for _, item := range a {
					y, err := NewSchemaOrReference(item, compiler.NewContext("allOf", context))
					if err != nil {
						errors = append(errors, err)
					}
					x.AllOf = append(x.AllOf, y)
				}
			}
		}
		// repeated SchemaOrReference one_of = 27;
		v27 := compiler.MapValueForKey(m, "oneOf")
		if v27 != nil {
			// repeated SchemaOrReference
			x.OneOf = make([]*SchemaOrReference, 0)
			a, ok := v27.([]interface{})
			if ok {
				for _, item := range a {
					y, err := NewSchemaOrReference(item, compiler.NewContext("oneOf", context))
					if err != nil {
						errors = append(errors, err)
					}
					x.OneOf = append(x.OneOf, y)
				}
			}
		}
		// repeated SchemaOrReference any_of = 28;
		v28 := compiler.MapValueForKey(m, "anyOf")
		if v28 != nil {
			// repeated SchemaOrReference
			x.AnyOf = make([]*SchemaOrReference, 0)
			a, ok := v28.([]interface{})
			if ok {
				for _, item := range a {
					y, err := NewSchemaOrReference(item, compiler.NewContext("anyOf", context))
					if err != nil {
						errors = append(errors, err)
					}
					x.AnyOf = append(x.AnyOf, y)
				}
			}
		}
		// Schema not = 29;
		v29 := compiler.MapValueForKey(m, "not")
		if v29 != nil {
			var err error
			x.Not, err = NewSchema(v29, compiler.NewContext("not", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// ItemsItem items = 30;
		v30 := compiler.MapValueForKey(m, "items")
		if v30 != nil {
			var err error
			x.Items, err = NewItemsItem(v30, compiler.NewContext("items", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// Properties properties = 31;
		v31 := compiler.MapValueForKey(m, "properties")
		if v31 != nil {
			var err error
			x.Properties, err = NewProperties(v31, compiler.NewContext("properties", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// AdditionalPropertiesItem additional_properties = 32;
		v32 := compiler.MapValueForKey(m, "additionalProperties")
		if v32 != nil {
			var err error
			x.AdditionalProperties, err = NewAdditionalPropertiesItem(v32, compiler.NewContext("additionalProperties", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// DefaultType default = 33;
		v33 := compiler.MapValueForKey(m, "default")
		if v33 != nil {
			var err error
			x.Default, err = NewDefaultType(v33, compiler.NewContext("default", context))
			if err != nil {
				errors = append(errors, err)
			}
		}
		// string description = 34;
		v34 := compiler.MapValueForKey(m, "description")
		if v34 != nil {
			x.Description, ok = v34.(string)
			if !ok {
				message := fmt.Sprintf("has unexpected value for description: %+v (%T)", v34, v34)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// string format = 35;
		v35 := compiler.MapValueForKey(m, "format")
		if v35 != nil {
			x.Format, ok = v35.(string)
			if !ok {
				message := fmt.Sprintf("has unexpected value for format: %+v (%T)", v35, v35)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// repeated NamedAny specification_extension = 36;
		// MAP: Any ^x-
		x.SpecificationExtension = make([]*NamedAny, 0)
		for _, item := range m {
			k, ok := compiler.StringValue(item.Key)
			if ok {
				v := item.Value
				if strings.HasPrefix(k, "x-") {
					pair := &NamedAny{}
					pair.Name = k
					result := &Any{}
					handled, resultFromExt, err := compiler.HandleExtension(context, v, k)
					if handled {
						if err != nil {
							errors = append(errors, err)
						} else {
							bytes, _ := yaml.Marshal(v)
							result.Yaml = string(bytes)
							result.Value = resultFromExt
							pair.Value = result
						}
					} else {
						pair.Value, err = NewAny(v, compiler.NewContext(k, context))
						if err != nil {
							errors = append(errors, err)
						}
					}
					x.SpecificationExtension = append(x.SpecificationExtension, pair)
				}
			}
		}
	}
	return x, compiler.NewErrorGroupOrNil(errors)
}
