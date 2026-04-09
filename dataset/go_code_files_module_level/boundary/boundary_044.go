func preProcess(cldrVar *cldr.CLDR) {

	for _, l := range cldrVar.Locales() {

		fmt.Println("Pre Processing:", l)

		split := strings.SplitN(l, "_", 2)
		baseLocale := split[0]
		// inheritedLocale := baseLocale

		// // one of the inherited english locales
		// // http://cldr.unicode.org/development/development-process/design-proposals/english-inheritance
		// if l == "en_001" || l == "en_GB" {
		// 	inheritedLocale = l
		// }

		trans := &translator{
			Locale:     l,
			BaseLocale: baseLocale,
			// InheritedLocale: inheritedLocale,
		}

		// if is a base locale
		if len(split) == 1 {
			baseTranslators[baseLocale] = trans
		}

		// baseTranslators[l] = trans
		// baseTranslators[baseLocale] = trans // allowing for unofficial fallback if none exists
		translators[l] = trans

		// get number, currency and datetime symbols

		// number values
		ldml := cldrVar.RawLDML(l)

		// some just have no data...
		if ldml.Numbers != nil {

			if len(ldml.Numbers.Symbols) > 0 {

				symbol := ldml.Numbers.Symbols[0]

				// Try to get the default numbering system instead of the first one
				systems := ldml.Numbers.DefaultNumberingSystem
				// There shouldn't really be more than one DefaultNumberingSystem
				if len(systems) > 0 {
					if dns := systems[0].Data(); dns != "" {
						for k := range ldml.Numbers.Symbols {
							if ldml.Numbers.Symbols[k].NumberSystem == dns {
								symbol = ldml.Numbers.Symbols[k]
								break
							}
						}
					}
				}

				if len(symbol.Decimal) > 0 {
					trans.Decimal = symbol.Decimal[0].Data()
				}
				if len(symbol.Group) > 0 {
					trans.Group = symbol.Group[0].Data()
				}
				if len(symbol.MinusSign) > 0 {
					trans.Minus = symbol.MinusSign[0].Data()
				}
				if len(symbol.PercentSign) > 0 {
					trans.Percent = symbol.PercentSign[0].Data()
				}
				if len(symbol.PerMille) > 0 {
					trans.PerMille = symbol.PerMille[0].Data()
				}

				if len(symbol.TimeSeparator) > 0 {
					trans.TimeSeparator = symbol.TimeSeparator[0].Data()
				}

				if len(symbol.Infinity) > 0 {
					trans.Infinity = symbol.Infinity[0].Data()
				}
			}

			if ldml.Numbers.Currencies != nil {

				for _, currency := range ldml.Numbers.Currencies.Currency {

					if len(strings.TrimSpace(currency.Type)) == 0 {
						continue
					}

					globalCurrenciesMap[currency.Type] = struct{}{}
				}
			}

			if len(ldml.Numbers.DecimalFormats) > 0 && len(ldml.Numbers.DecimalFormats[0].DecimalFormatLength) > 0 {

				for _, dfl := range ldml.Numbers.DecimalFormats[0].DecimalFormatLength {
					if len(dfl.Type) == 0 {
						trans.DecimalNumberFormat = dfl.DecimalFormat[0].Pattern[0].Data()
						break
					}
				}
			}

			if len(ldml.Numbers.PercentFormats) > 0 && len(ldml.Numbers.PercentFormats[0].PercentFormatLength) > 0 {

				for _, dfl := range ldml.Numbers.PercentFormats[0].PercentFormatLength {
					if len(dfl.Type) == 0 {
						trans.PercentNumberFormat = dfl.PercentFormat[0].Pattern[0].Data()
						break
					}
				}
			}

			if len(ldml.Numbers.CurrencyFormats) > 0 && len(ldml.Numbers.CurrencyFormats[0].CurrencyFormatLength) > 0 {

				if len(ldml.Numbers.CurrencyFormats[0].CurrencyFormatLength[0].CurrencyFormat) > 1 {

					split := strings.SplitN(ldml.Numbers.CurrencyFormats[0].CurrencyFormatLength[0].CurrencyFormat[1].Pattern[0].Data(), ";", 2)

					trans.CurrencyNumberFormat = split[0]

					if len(split) > 1 && len(split[1]) > 0 {
						trans.NegativeCurrencyNumberFormat = split[1]
					} else {
						trans.NegativeCurrencyNumberFormat = trans.CurrencyNumberFormat
					}
				} else {
					trans.CurrencyNumberFormat = ldml.Numbers.CurrencyFormats[0].CurrencyFormatLength[0].CurrencyFormat[0].Pattern[0].Data()
					trans.NegativeCurrencyNumberFormat = trans.CurrencyNumberFormat
				}
			}
		}

		if ldml.Dates != nil {

			if ldml.Dates.TimeZoneNames != nil {

				for _, zone := range ldml.Dates.TimeZoneNames.Metazone {

					for _, short := range zone.Short {

						if len(short.Standard) > 0 {
							za, ok := timezones[zone.Type]
							if !ok {
								za = new(zoneAbbrev)
								timezones[zone.Type] = za
							}
							za.standard = short.Standard[0].Data()
						}

						if len(short.Daylight) > 0 {
							za, ok := timezones[zone.Type]
							if !ok {
								za = new(zoneAbbrev)
								timezones[zone.Type] = za
							}
							za.daylight = short.Daylight[0].Data()
						}
					}

					for _, long := range zone.Long {

						if trans.timezones == nil {
							trans.timezones = make(map[string]*zoneAbbrev)
						}

						if len(long.Standard) > 0 {
							za, ok := trans.timezones[zone.Type]
							if !ok {
								za = new(zoneAbbrev)
								trans.timezones[zone.Type] = za
							}
							za.standard = long.Standard[0].Data()
						}

						za, ok := trans.timezones[zone.Type]
						if !ok {
							za = new(zoneAbbrev)
							trans.timezones[zone.Type] = za
						}

						if len(long.Daylight) > 0 {
							za.daylight = long.Daylight[0].Data()
						} else {
							za.daylight = za.standard
						}
					}
				}
			}

			if ldml.Dates.Calendars != nil {

				var calendar *cldr.Calendar

				for _, cal := range ldml.Dates.Calendars.Calendar {
					if cal.Type == "gregorian" {
						calendar = cal
					}
				}

				if calendar != nil {

					if calendar.DateFormats != nil {

						for _, datefmt := range calendar.DateFormats.DateFormatLength {

							switch datefmt.Type {
							case "full":
								trans.FmtDateFull = datefmt.DateFormat[0].Pattern[0].Data()

							case "long":
								trans.FmtDateLong = datefmt.DateFormat[0].Pattern[0].Data()

							case "medium":
								trans.FmtDateMedium = datefmt.DateFormat[0].Pattern[0].Data()

							case "short":
								trans.FmtDateShort = datefmt.DateFormat[0].Pattern[0].Data()
							}
						}
					}

					if calendar.TimeFormats != nil {

						for _, datefmt := range calendar.TimeFormats.TimeFormatLength {

							switch datefmt.Type {
							case "full":
								trans.FmtTimeFull = datefmt.TimeFormat[0].Pattern[0].Data()
							case "long":
								trans.FmtTimeLong = datefmt.TimeFormat[0].Pattern[0].Data()
							case "medium":
								trans.FmtTimeMedium = datefmt.TimeFormat[0].Pattern[0].Data()
							case "short":
								trans.FmtTimeShort = datefmt.TimeFormat[0].Pattern[0].Data()
							}
						}
					}

					if calendar.Months != nil {

						// month context starts at 'format', but there is also has 'stand-alone'
						// I'm making the decision to use the 'stand-alone' if, and only if,
						// the value does not exist in the 'format' month context
						var abbrSet, narrSet, wideSet bool

						for _, monthctx := range calendar.Months.MonthContext {

							for _, months := range monthctx.MonthWidth {

								var monthData []string

								for _, m := range months.Month {

									if len(m.Data()) == 0 {
										continue
									}

									switch m.Type {
									case "1":
										monthData = append(monthData, m.Data())
									case "2":
										monthData = append(monthData, m.Data())
									case "3":
										monthData = append(monthData, m.Data())
									case "4":
										monthData = append(monthData, m.Data())
									case "5":
										monthData = append(monthData, m.Data())
									case "6":
										monthData = append(monthData, m.Data())
									case "7":
										monthData = append(monthData, m.Data())
									case "8":
										monthData = append(monthData, m.Data())
									case "9":
										monthData = append(monthData, m.Data())
									case "10":
										monthData = append(monthData, m.Data())
									case "11":
										monthData = append(monthData, m.Data())
									case "12":
										monthData = append(monthData, m.Data())
									}
								}

								if len(monthData) > 0 {

									// making array indexes line up with month values
									// so I'll have an extra empty value, it's way faster
									// than a switch over all type values...
									monthData = append(make([]string, 1, len(monthData)+1), monthData...)

									switch months.Type {
									case "abbreviated":
										if !abbrSet {
											abbrSet = true
											trans.FmtMonthsAbbreviated = fmt.Sprintf("%#v", monthData)
										}
									case "narrow":
										if !narrSet {
											narrSet = true
											trans.FmtMonthsNarrow = fmt.Sprintf("%#v", monthData)
										}
									case "wide":
										if !wideSet {
											wideSet = true
											trans.FmtMonthsWide = fmt.Sprintf("%#v", monthData)
										}
									}
								}
							}
						}
					}

					if calendar.Days != nil {

						// day context starts at 'format', but there is also has 'stand-alone'
						// I'm making the decision to use the 'stand-alone' if, and only if,
						// the value does not exist in the 'format' day context
						var abbrSet, narrSet, shortSet, wideSet bool

						for _, dayctx := range calendar.Days.DayContext {

							for _, days := range dayctx.DayWidth {

								var dayData []string

								for _, d := range days.Day {

									switch d.Type {
									case "sun":
										dayData = append(dayData, d.Data())
									case "mon":
										dayData = append(dayData, d.Data())
									case "tue":
										dayData = append(dayData, d.Data())
									case "wed":
										dayData = append(dayData, d.Data())
									case "thu":
										dayData = append(dayData, d.Data())
									case "fri":
										dayData = append(dayData, d.Data())
									case "sat":
										dayData = append(dayData, d.Data())
									}
								}

								if len(dayData) > 0 {
									switch days.Type {
									case "abbreviated":
										if !abbrSet {
											abbrSet = true
											trans.FmtDaysAbbreviated = fmt.Sprintf("%#v", dayData)
										}
									case "narrow":
										if !narrSet {
											narrSet = true
											trans.FmtDaysNarrow = fmt.Sprintf("%#v", dayData)
										}
									case "short":
										if !shortSet {
											shortSet = true
											trans.FmtDaysShort = fmt.Sprintf("%#v", dayData)
										}
									case "wide":
										if !wideSet {
											wideSet = true
											trans.FmtDaysWide = fmt.Sprintf("%#v", dayData)
										}
									}
								}
							}
						}
					}

					if calendar.DayPeriods != nil {

						// day periods context starts at 'format', but there is also has 'stand-alone'
						// I'm making the decision to use the 'stand-alone' if, and only if,
						// the value does not exist in the 'format' day period context
						var abbrSet, narrSet, shortSet, wideSet bool

						for _, ctx := range calendar.DayPeriods.DayPeriodContext {

							for _, width := range ctx.DayPeriodWidth {

								// [0] = AM
								// [0] = PM
								ampm := make([]string, 2, 2)

								for _, d := range width.DayPeriod {

									if d.Type == "am" {
										ampm[0] = d.Data()
										continue
									}

									if d.Type == "pm" {
										ampm[1] = d.Data()
									}
								}

								switch width.Type {
								case "abbreviated":
									if !abbrSet {
										abbrSet = true
										trans.FmtPeriodsAbbreviated = fmt.Sprintf("%#v", ampm)
									}
								case "narrow":
									if !narrSet {
										narrSet = true
										trans.FmtPeriodsNarrow = fmt.Sprintf("%#v", ampm)
									}
								case "short":
									if !shortSet {
										shortSet = true
										trans.FmtPeriodsShort = fmt.Sprintf("%#v", ampm)
									}
								case "wide":
									if !wideSet {
										wideSet = true
										trans.FmtPeriodsWide = fmt.Sprintf("%#v", ampm)
									}
								}
							}
						}
					}

					if calendar.Eras != nil {

						// [0] = BC
						// [0] = AD
						abbrev := make([]string, 2, 2)
						narr := make([]string, 2, 2)
						wide := make([]string, 2, 2)

						if calendar.Eras.EraAbbr != nil {

							if len(calendar.Eras.EraAbbr.Era) == 4 {
								abbrev[0] = calendar.Eras.EraAbbr.Era[0].Data()
								abbrev[1] = calendar.Eras.EraAbbr.Era[2].Data()
							} else if len(calendar.Eras.EraAbbr.Era) == 2 {
								abbrev[0] = calendar.Eras.EraAbbr.Era[0].Data()
								abbrev[1] = calendar.Eras.EraAbbr.Era[1].Data()
							}
						}

						if calendar.Eras.EraNarrow != nil {

							if len(calendar.Eras.EraNarrow.Era) == 4 {
								narr[0] = calendar.Eras.EraNarrow.Era[0].Data()
								narr[1] = calendar.Eras.EraNarrow.Era[2].Data()
							} else if len(calendar.Eras.EraNarrow.Era) == 2 {
								narr[0] = calendar.Eras.EraNarrow.Era[0].Data()
								narr[1] = calendar.Eras.EraNarrow.Era[1].Data()
							}
						}

						if calendar.Eras.EraNames != nil {

							if len(calendar.Eras.EraNames.Era) == 4 {
								wide[0] = calendar.Eras.EraNames.Era[0].Data()
								wide[1] = calendar.Eras.EraNames.Era[2].Data()
							} else if len(calendar.Eras.EraNames.Era) == 2 {
								wide[0] = calendar.Eras.EraNames.Era[0].Data()
								wide[1] = calendar.Eras.EraNames.Era[1].Data()
							}
						}

						trans.FmtErasAbbreviated = fmt.Sprintf("%#v", abbrev)
						trans.FmtErasNarrow = fmt.Sprintf("%#v", narr)
						trans.FmtErasWide = fmt.Sprintf("%#v", wide)
					}
				}
			}
		}
	}

	for k := range globalCurrenciesMap {
		globalCurrencies = append(globalCurrencies, k)
	}

	sort.Strings(globalCurrencies)

	for i, loc := range globalCurrencies {
		globCurrencyIdxMap[loc] = i
	}
}
