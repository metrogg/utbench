func (session *Session) bytes2Value(col *core.Column, fieldValue *reflect.Value, data []byte) error {
	if structConvert, ok := fieldValue.Addr().Interface().(core.Conversion); ok {
		return structConvert.FromDB(data)
	}

	if structConvert, ok := fieldValue.Interface().(core.Conversion); ok {
		return structConvert.FromDB(data)
	}

	var v interface{}
	key := col.Name
	fieldType := fieldValue.Type()

	switch fieldType.Kind() {
	case reflect.Complex64, reflect.Complex128:
		x := reflect.New(fieldType)
		if len(data) > 0 {
			err := DefaultJSONHandler.Unmarshal(data, x.Interface())
			if err != nil {
				session.engine.logger.Error(err)
				return err
			}
			fieldValue.Set(x.Elem())
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		v = data
		t := fieldType.Elem()
		k := t.Kind()
		if col.SQLType.IsText() {
			x := reflect.New(fieldType)
			if len(data) > 0 {
				err := DefaultJSONHandler.Unmarshal(data, x.Interface())
				if err != nil {
					session.engine.logger.Error(err)
					return err
				}
				fieldValue.Set(x.Elem())
			}
		} else if col.SQLType.IsBlob() {
			if k == reflect.Uint8 {
				fieldValue.Set(reflect.ValueOf(v))
			} else {
				x := reflect.New(fieldType)
				if len(data) > 0 {
					err := DefaultJSONHandler.Unmarshal(data, x.Interface())
					if err != nil {
						session.engine.logger.Error(err)
						return err
					}
					fieldValue.Set(x.Elem())
				}
			}
		} else {
			return ErrUnSupportedType
		}
	case reflect.String:
		fieldValue.SetString(string(data))
	case reflect.Bool:
		v, err := asBool(data)
		if err != nil {
			return fmt.Errorf("arg %v as bool: %s", key, err.Error())
		}
		fieldValue.Set(reflect.ValueOf(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sdata := string(data)
		var x int64
		var err error
		// for mysql, when use bit, it returned \x01
		if col.SQLType.Name == core.Bit &&
			session.engine.dialect.DBType() == core.MYSQL { // !nashtsai! TODO dialect needs to provide conversion interface API
			if len(data) == 1 {
				x = int64(data[0])
			} else {
				x = 0
			}
		} else if strings.HasPrefix(sdata, "0x") {
			x, err = strconv.ParseInt(sdata, 16, 64)
		} else if strings.HasPrefix(sdata, "0") {
			x, err = strconv.ParseInt(sdata, 8, 64)
		} else if strings.EqualFold(sdata, "true") {
			x = 1
		} else if strings.EqualFold(sdata, "false") {
			x = 0
		} else {
			x, err = strconv.ParseInt(sdata, 10, 64)
		}
		if err != nil {
			return fmt.Errorf("arg %v as int: %s", key, err.Error())
		}
		fieldValue.SetInt(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return fmt.Errorf("arg %v as float64: %s", key, err.Error())
		}
		fieldValue.SetFloat(x)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		x, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("arg %v as int: %s", key, err.Error())
		}
		fieldValue.SetUint(x)
	//Currently only support Time type
	case reflect.Struct:
		// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
		if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
			if err := nulVal.Scan(data); err != nil {
				return fmt.Errorf("sql.Scan(%v) failed: %s ", data, err.Error())
			}
		} else {
			if fieldType.ConvertibleTo(core.TimeType) {
				x, err := session.byte2Time(col, data)
				if err != nil {
					return err
				}
				v = x
				fieldValue.Set(reflect.ValueOf(v).Convert(fieldType))
			} else if session.statement.UseCascade {
				table, err := session.engine.autoMapType(*fieldValue)
				if err != nil {
					return err
				}

				// TODO: current only support 1 primary key
				if len(table.PrimaryKeys) > 1 {
					return errors.New("unsupported composited primary key cascade")
				}

				var pk = make(core.PK, len(table.PrimaryKeys))
				rawValueType := table.ColumnType(table.PKColumns()[0].FieldName)
				pk[0], err = str2PK(string(data), rawValueType)
				if err != nil {
					return err
				}

				if !isPKZero(pk) {
					// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
					// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
					// property to be fetched lazily
					structInter := reflect.New(fieldValue.Type())
					has, err := session.ID(pk).NoCascade().get(structInter.Interface())
					if err != nil {
						return err
					}
					if has {
						v = structInter.Elem().Interface()
						fieldValue.Set(reflect.ValueOf(v))
					} else {
						return errors.New("cascade obj is not exist")
					}
				}
			}
		}
	case reflect.Ptr:
		// !nashtsai! TODO merge duplicated codes above
		//typeStr := fieldType.String()
		switch fieldType.Elem().Kind() {
		// case "*string":
		case core.StringType.Kind():
			x := string(data)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*bool":
		case core.BoolType.Kind():
			d := string(data)
			v, err := strconv.ParseBool(d)
			if err != nil {
				return fmt.Errorf("arg %v as bool: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&v).Convert(fieldType))
		// case "*complex64":
		case core.Complex64Type.Kind():
			var x complex64
			if len(data) > 0 {
				err := DefaultJSONHandler.Unmarshal(data, &x)
				if err != nil {
					session.engine.logger.Error(err)
					return err
				}
				fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
			}
		// case "*complex128":
		case core.Complex128Type.Kind():
			var x complex128
			if len(data) > 0 {
				err := DefaultJSONHandler.Unmarshal(data, &x)
				if err != nil {
					session.engine.logger.Error(err)
					return err
				}
				fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
			}
		// case "*float64":
		case core.Float64Type.Kind():
			x, err := strconv.ParseFloat(string(data), 64)
			if err != nil {
				return fmt.Errorf("arg %v as float64: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*float32":
		case core.Float32Type.Kind():
			var x float32
			x1, err := strconv.ParseFloat(string(data), 32)
			if err != nil {
				return fmt.Errorf("arg %v as float32: %s", key, err.Error())
			}
			x = float32(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint64":
		case core.Uint64Type.Kind():
			var x uint64
			x, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint":
		case core.UintType.Kind():
			var x uint
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint32":
		case core.Uint32Type.Kind():
			var x uint32
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint32(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint8":
		case core.Uint8Type.Kind():
			var x uint8
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint8(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint16":
		case core.Uint16Type.Kind():
			var x uint16
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint16(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int64":
		case core.Int64Type.Kind():
			sdata := string(data)
			var x int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == core.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int64(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x, err = strconv.ParseInt(sdata, 16, 64)
			} else if strings.HasPrefix(sdata, "0") {
				x, err = strconv.ParseInt(sdata, 8, 64)
			} else {
				x, err = strconv.ParseInt(sdata, 10, 64)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int":
		case core.IntType.Kind():
			sdata := string(data)
			var x int
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == core.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int32":
		case core.Int32Type.Kind():
			sdata := string(data)
			var x int32
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == core.Bit &&
				session.engine.dialect.DBType() == core.MYSQL {
				if len(data) == 1 {
					x = int32(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int32(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int32(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int32(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int8":
		case core.Int8Type.Kind():
			sdata := string(data)
			var x int8
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == core.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int8(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int8(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int8(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int8(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int16":
		case core.Int16Type.Kind():
			sdata := string(data)
			var x int16
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == core.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int16(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int16(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int16(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int16(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*SomeStruct":
		case reflect.Struct:
			switch fieldType {
			// case "*.time.Time":
			case core.PtrTimeType:
				x, err := session.byte2Time(col, data)
				if err != nil {
					return err
				}
				v = x
				fieldValue.Set(reflect.ValueOf(&x))
			default:
				if session.statement.UseCascade {
					structInter := reflect.New(fieldType.Elem())
					table, err := session.engine.autoMapType(structInter.Elem())
					if err != nil {
						return err
					}

					if len(table.PrimaryKeys) > 1 {
						return errors.New("unsupported composited primary key cascade")
					}

					var pk = make(core.PK, len(table.PrimaryKeys))
					rawValueType := table.ColumnType(table.PKColumns()[0].FieldName)
					pk[0], err = str2PK(string(data), rawValueType)
					if err != nil {
						return err
					}

					if !isPKZero(pk) {
						// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
						// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
						// property to be fetched lazily
						has, err := session.ID(pk).NoCascade().get(structInter.Interface())
						if err != nil {
							return err
						}
						if has {
							v = structInter.Interface()
							fieldValue.Set(reflect.ValueOf(v))
						} else {
							return errors.New("cascade obj is not exist")
						}
					}
				} else {
					return fmt.Errorf("unsupported struct type in Scan: %s", fieldValue.Type().String())
				}
			}
		default:
			return fmt.Errorf("unsupported type in Scan: %s", fieldValue.Type().String())
		}
	default:
		return fmt.Errorf("unsupported type in Scan: %s", fieldValue.Type().String())
	}

	return nil
}
