func (t *TranslatorRules) merge(tNew *TranslatorRules) {

	t.Plural = stringMerge(t.Plural, tNew.Plural)

	if tNew.PluralRuleFunc != nil {
		t.PluralRuleFunc = tNew.PluralRuleFunc
	}

	t.Direction = stringMerge(t.Direction, tNew.Direction)

	t.Numbers.Symbols.Decimal = stringMerge(t.Numbers.Symbols.Decimal, tNew.Numbers.Symbols.Decimal)
	t.Numbers.Symbols.Group = stringMerge(t.Numbers.Symbols.Group, tNew.Numbers.Symbols.Group)
	t.Numbers.Symbols.Negative = stringMerge(t.Numbers.Symbols.Negative, tNew.Numbers.Symbols.Negative)
	t.Numbers.Symbols.Percent = stringMerge(t.Numbers.Symbols.Percent, tNew.Numbers.Symbols.Percent)
	t.Numbers.Symbols.Permille = stringMerge(t.Numbers.Symbols.Permille, tNew.Numbers.Symbols.Permille)
	t.Numbers.Symbols.Permille = stringMerge(t.Numbers.Symbols.Permille, tNew.Numbers.Symbols.Permille)
	t.Numbers.Symbols.Permille = stringMerge(t.Numbers.Symbols.Permille, tNew.Numbers.Symbols.Permille)
	t.Numbers.Formats.Decimal = stringMerge(t.Numbers.Formats.Decimal, tNew.Numbers.Formats.Decimal)
	t.Numbers.Formats.Currency = stringMerge(t.Numbers.Formats.Currency, tNew.Numbers.Formats.Currency)
	t.Numbers.Formats.Percent = stringMerge(t.Numbers.Formats.Percent, tNew.Numbers.Formats.Percent)

	for i, c := range tNew.Currencies {
		if t.Currencies == nil {
			t.Currencies = tNew.Currencies
		} else if _, ok := t.Currencies[i]; !ok {
			t.Currencies[i] = c
		} else {
			tmp := t.Currencies[i]
			tmp.Symbol = stringMerge(tmp.Symbol, c.Symbol)
			t.Currencies[i] = tmp
		}
	}

	t.DateTime.TimeSeparator = stringMerge(t.DateTime.TimeSeparator, tNew.DateTime.TimeSeparator)

	t.DateTime.Formats.Date.Full = stringMerge(t.DateTime.Formats.Date.Full, tNew.DateTime.Formats.Date.Full)
	t.DateTime.Formats.Date.Long = stringMerge(t.DateTime.Formats.Date.Long, tNew.DateTime.Formats.Date.Long)
	t.DateTime.Formats.Date.Medium = stringMerge(t.DateTime.Formats.Date.Medium, tNew.DateTime.Formats.Date.Medium)
	t.DateTime.Formats.Date.Short = stringMerge(t.DateTime.Formats.Date.Short, tNew.DateTime.Formats.Date.Short)
	t.DateTime.Formats.Time.Full = stringMerge(t.DateTime.Formats.Time.Full, tNew.DateTime.Formats.Time.Full)
	t.DateTime.Formats.Time.Long = stringMerge(t.DateTime.Formats.Time.Long, tNew.DateTime.Formats.Time.Long)
	t.DateTime.Formats.Time.Medium = stringMerge(t.DateTime.Formats.Time.Medium, tNew.DateTime.Formats.Time.Medium)
	t.DateTime.Formats.Time.Short = stringMerge(t.DateTime.Formats.Time.Short, tNew.DateTime.Formats.Time.Short)
	t.DateTime.Formats.DateTime.Full = stringMerge(t.DateTime.Formats.DateTime.Full, tNew.DateTime.Formats.DateTime.Full)
	t.DateTime.Formats.DateTime.Long = stringMerge(t.DateTime.Formats.DateTime.Long, tNew.DateTime.Formats.DateTime.Long)
	t.DateTime.Formats.DateTime.Medium = stringMerge(t.DateTime.Formats.DateTime.Medium, tNew.DateTime.Formats.DateTime.Medium)
	t.DateTime.Formats.DateTime.Short = stringMerge(t.DateTime.Formats.DateTime.Short, tNew.DateTime.Formats.DateTime.Short)

	t.DateTime.FormatNames.Months.Abbreviated.Month1 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month1, tNew.DateTime.FormatNames.Months.Abbreviated.Month1)
	t.DateTime.FormatNames.Months.Abbreviated.Month2 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month2, tNew.DateTime.FormatNames.Months.Abbreviated.Month2)
	t.DateTime.FormatNames.Months.Abbreviated.Month3 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month3, tNew.DateTime.FormatNames.Months.Abbreviated.Month3)
	t.DateTime.FormatNames.Months.Abbreviated.Month4 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month4, tNew.DateTime.FormatNames.Months.Abbreviated.Month4)
	t.DateTime.FormatNames.Months.Abbreviated.Month5 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month5, tNew.DateTime.FormatNames.Months.Abbreviated.Month5)
	t.DateTime.FormatNames.Months.Abbreviated.Month6 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month6, tNew.DateTime.FormatNames.Months.Abbreviated.Month6)
	t.DateTime.FormatNames.Months.Abbreviated.Month7 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month7, tNew.DateTime.FormatNames.Months.Abbreviated.Month7)
	t.DateTime.FormatNames.Months.Abbreviated.Month8 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month8, tNew.DateTime.FormatNames.Months.Abbreviated.Month8)
	t.DateTime.FormatNames.Months.Abbreviated.Month9 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month9, tNew.DateTime.FormatNames.Months.Abbreviated.Month9)
	t.DateTime.FormatNames.Months.Abbreviated.Month10 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month10, tNew.DateTime.FormatNames.Months.Abbreviated.Month10)
	t.DateTime.FormatNames.Months.Abbreviated.Month11 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month11, tNew.DateTime.FormatNames.Months.Abbreviated.Month11)
	t.DateTime.FormatNames.Months.Abbreviated.Month12 = stringMerge(t.DateTime.FormatNames.Months.Abbreviated.Month12, tNew.DateTime.FormatNames.Months.Abbreviated.Month12)

	t.DateTime.FormatNames.Months.Narrow.Month1 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month1, tNew.DateTime.FormatNames.Months.Narrow.Month1)
	t.DateTime.FormatNames.Months.Narrow.Month2 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month2, tNew.DateTime.FormatNames.Months.Narrow.Month2)
	t.DateTime.FormatNames.Months.Narrow.Month3 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month3, tNew.DateTime.FormatNames.Months.Narrow.Month3)
	t.DateTime.FormatNames.Months.Narrow.Month4 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month4, tNew.DateTime.FormatNames.Months.Narrow.Month4)
	t.DateTime.FormatNames.Months.Narrow.Month5 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month5, tNew.DateTime.FormatNames.Months.Narrow.Month5)
	t.DateTime.FormatNames.Months.Narrow.Month6 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month6, tNew.DateTime.FormatNames.Months.Narrow.Month6)
	t.DateTime.FormatNames.Months.Narrow.Month7 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month7, tNew.DateTime.FormatNames.Months.Narrow.Month7)
	t.DateTime.FormatNames.Months.Narrow.Month8 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month8, tNew.DateTime.FormatNames.Months.Narrow.Month8)
	t.DateTime.FormatNames.Months.Narrow.Month9 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month9, tNew.DateTime.FormatNames.Months.Narrow.Month9)
	t.DateTime.FormatNames.Months.Narrow.Month10 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month10, tNew.DateTime.FormatNames.Months.Narrow.Month10)
	t.DateTime.FormatNames.Months.Narrow.Month11 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month11, tNew.DateTime.FormatNames.Months.Narrow.Month11)
	t.DateTime.FormatNames.Months.Narrow.Month12 = stringMerge(t.DateTime.FormatNames.Months.Narrow.Month12, tNew.DateTime.FormatNames.Months.Narrow.Month12)

	t.DateTime.FormatNames.Months.Wide.Month1 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month1, tNew.DateTime.FormatNames.Months.Wide.Month1)
	t.DateTime.FormatNames.Months.Wide.Month2 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month2, tNew.DateTime.FormatNames.Months.Wide.Month2)
	t.DateTime.FormatNames.Months.Wide.Month3 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month3, tNew.DateTime.FormatNames.Months.Wide.Month3)
	t.DateTime.FormatNames.Months.Wide.Month4 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month4, tNew.DateTime.FormatNames.Months.Wide.Month4)
	t.DateTime.FormatNames.Months.Wide.Month5 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month5, tNew.DateTime.FormatNames.Months.Wide.Month5)
	t.DateTime.FormatNames.Months.Wide.Month6 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month6, tNew.DateTime.FormatNames.Months.Wide.Month6)
	t.DateTime.FormatNames.Months.Wide.Month7 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month7, tNew.DateTime.FormatNames.Months.Wide.Month7)
	t.DateTime.FormatNames.Months.Wide.Month8 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month8, tNew.DateTime.FormatNames.Months.Wide.Month8)
	t.DateTime.FormatNames.Months.Wide.Month9 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month9, tNew.DateTime.FormatNames.Months.Wide.Month9)
	t.DateTime.FormatNames.Months.Wide.Month10 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month10, tNew.DateTime.FormatNames.Months.Wide.Month10)
	t.DateTime.FormatNames.Months.Wide.Month11 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month11, tNew.DateTime.FormatNames.Months.Wide.Month11)
	t.DateTime.FormatNames.Months.Wide.Month12 = stringMerge(t.DateTime.FormatNames.Months.Wide.Month12, tNew.DateTime.FormatNames.Months.Wide.Month12)

	t.DateTime.FormatNames.Days.Abbreviated.Sun = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Sun, tNew.DateTime.FormatNames.Days.Abbreviated.Sun)
	t.DateTime.FormatNames.Days.Abbreviated.Mon = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Mon, tNew.DateTime.FormatNames.Days.Abbreviated.Mon)
	t.DateTime.FormatNames.Days.Abbreviated.Tue = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Tue, tNew.DateTime.FormatNames.Days.Abbreviated.Tue)
	t.DateTime.FormatNames.Days.Abbreviated.Wed = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Wed, tNew.DateTime.FormatNames.Days.Abbreviated.Wed)
	t.DateTime.FormatNames.Days.Abbreviated.Thu = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Thu, tNew.DateTime.FormatNames.Days.Abbreviated.Thu)
	t.DateTime.FormatNames.Days.Abbreviated.Fri = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Fri, tNew.DateTime.FormatNames.Days.Abbreviated.Fri)
	t.DateTime.FormatNames.Days.Abbreviated.Sat = stringMerge(t.DateTime.FormatNames.Days.Abbreviated.Sat, tNew.DateTime.FormatNames.Days.Abbreviated.Sat)

	t.DateTime.FormatNames.Days.Narrow.Sun = stringMerge(t.DateTime.FormatNames.Days.Narrow.Sun, tNew.DateTime.FormatNames.Days.Narrow.Sun)
	t.DateTime.FormatNames.Days.Narrow.Mon = stringMerge(t.DateTime.FormatNames.Days.Narrow.Mon, tNew.DateTime.FormatNames.Days.Narrow.Mon)
	t.DateTime.FormatNames.Days.Narrow.Tue = stringMerge(t.DateTime.FormatNames.Days.Narrow.Tue, tNew.DateTime.FormatNames.Days.Narrow.Tue)
	t.DateTime.FormatNames.Days.Narrow.Wed = stringMerge(t.DateTime.FormatNames.Days.Narrow.Wed, tNew.DateTime.FormatNames.Days.Narrow.Wed)
	t.DateTime.FormatNames.Days.Narrow.Thu = stringMerge(t.DateTime.FormatNames.Days.Narrow.Thu, tNew.DateTime.FormatNames.Days.Narrow.Thu)
	t.DateTime.FormatNames.Days.Narrow.Fri = stringMerge(t.DateTime.FormatNames.Days.Narrow.Fri, tNew.DateTime.FormatNames.Days.Narrow.Fri)
	t.DateTime.FormatNames.Days.Narrow.Sat = stringMerge(t.DateTime.FormatNames.Days.Narrow.Sat, tNew.DateTime.FormatNames.Days.Narrow.Sat)

	t.DateTime.FormatNames.Days.Short.Sun = stringMerge(t.DateTime.FormatNames.Days.Short.Sun, tNew.DateTime.FormatNames.Days.Short.Sun)
	t.DateTime.FormatNames.Days.Short.Mon = stringMerge(t.DateTime.FormatNames.Days.Short.Mon, tNew.DateTime.FormatNames.Days.Short.Mon)
	t.DateTime.FormatNames.Days.Short.Tue = stringMerge(t.DateTime.FormatNames.Days.Short.Tue, tNew.DateTime.FormatNames.Days.Short.Tue)
	t.DateTime.FormatNames.Days.Short.Wed = stringMerge(t.DateTime.FormatNames.Days.Short.Wed, tNew.DateTime.FormatNames.Days.Short.Wed)
	t.DateTime.FormatNames.Days.Short.Thu = stringMerge(t.DateTime.FormatNames.Days.Short.Thu, tNew.DateTime.FormatNames.Days.Short.Thu)
	t.DateTime.FormatNames.Days.Short.Fri = stringMerge(t.DateTime.FormatNames.Days.Short.Fri, tNew.DateTime.FormatNames.Days.Short.Fri)
	t.DateTime.FormatNames.Days.Short.Sat = stringMerge(t.DateTime.FormatNames.Days.Short.Sat, tNew.DateTime.FormatNames.Days.Short.Sat)

	t.DateTime.FormatNames.Days.Wide.Sun = stringMerge(t.DateTime.FormatNames.Days.Wide.Sun, tNew.DateTime.FormatNames.Days.Wide.Sun)
	t.DateTime.FormatNames.Days.Wide.Mon = stringMerge(t.DateTime.FormatNames.Days.Wide.Mon, tNew.DateTime.FormatNames.Days.Wide.Mon)
	t.DateTime.FormatNames.Days.Wide.Tue = stringMerge(t.DateTime.FormatNames.Days.Wide.Tue, tNew.DateTime.FormatNames.Days.Wide.Tue)
	t.DateTime.FormatNames.Days.Wide.Wed = stringMerge(t.DateTime.FormatNames.Days.Wide.Wed, tNew.DateTime.FormatNames.Days.Wide.Wed)
	t.DateTime.FormatNames.Days.Wide.Thu = stringMerge(t.DateTime.FormatNames.Days.Wide.Thu, tNew.DateTime.FormatNames.Days.Wide.Thu)
	t.DateTime.FormatNames.Days.Wide.Fri = stringMerge(t.DateTime.FormatNames.Days.Wide.Fri, tNew.DateTime.FormatNames.Days.Wide.Fri)
	t.DateTime.FormatNames.Days.Wide.Sat = stringMerge(t.DateTime.FormatNames.Days.Wide.Sat, tNew.DateTime.FormatNames.Days.Wide.Sat)

	t.DateTime.FormatNames.Periods.Abbreviated.AM = stringMerge(t.DateTime.FormatNames.Periods.Abbreviated.AM, tNew.DateTime.FormatNames.Periods.Abbreviated.AM)
	t.DateTime.FormatNames.Periods.Abbreviated.PM = stringMerge(t.DateTime.FormatNames.Periods.Abbreviated.PM, tNew.DateTime.FormatNames.Periods.Abbreviated.PM)
	t.DateTime.FormatNames.Periods.Narrow.AM = stringMerge(t.DateTime.FormatNames.Periods.Narrow.AM, tNew.DateTime.FormatNames.Periods.Narrow.AM)
	t.DateTime.FormatNames.Periods.Narrow.PM = stringMerge(t.DateTime.FormatNames.Periods.Narrow.PM, tNew.DateTime.FormatNames.Periods.Narrow.PM)
	t.DateTime.FormatNames.Periods.Wide.AM = stringMerge(t.DateTime.FormatNames.Periods.Wide.AM, tNew.DateTime.FormatNames.Periods.Wide.AM)
	t.DateTime.FormatNames.Periods.Wide.PM = stringMerge(t.DateTime.FormatNames.Periods.Wide.PM, tNew.DateTime.FormatNames.Periods.Wide.PM)
}
