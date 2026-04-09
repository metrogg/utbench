func (c *Client) convertRaw(term string, raw []interface{}) interface{} {
	// The things you do to get proper types.
	switch term {
	case "bu":
		o, err := bitfinex.NewBalanceInfoFromRaw(raw)
		if err != nil {
			return err
		}
		bu := bitfinex.BalanceUpdate(*o)
		return &bu
	case "ps":
		o, err := bitfinex.NewPositionSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "pn":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pn := bitfinex.PositionNew(*o)
		return &pn
	case "pu":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pu := bitfinex.PositionUpdate(*o)
		return &pu
	case "pc":
		o, err := bitfinex.NewPositionFromRaw(raw)
		if err != nil {
			return err
		}
		pc := bitfinex.PositionCancel(*o)
		return &pc
	case "ws":
		o, err := bitfinex.NewWalletSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "wu":
		o, err := bitfinex.NewWalletFromRaw(raw)
		if err != nil {
			return err
		}
		wu := bitfinex.WalletUpdate(*o)
		return &wu
	case "os":
		o, err := bitfinex.NewOrderSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "on":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		on := bitfinex.OrderNew(*o)
		return &on
	case "ou":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		ou := bitfinex.OrderUpdate(*o)
		return &ou
	case "oc":
		o, err := bitfinex.NewOrderFromRaw(raw)
		if err != nil {
			return err
		}
		oc := bitfinex.OrderCancel(*o)
		return &oc
	case "hts":
		o, err := bitfinex.NewTradeExecutionUpdateSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		hts := bitfinex.HistoricalTradeSnapshot(*o)
		return &hts
	case "te":
		o, err := bitfinex.NewTradeExecutionFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "tu":
		tu, err := bitfinex.NewTradeExecutionUpdateFromRaw(raw)
		if err != nil {
			return err
		}
		return tu
	case "fte":
		o, err := bitfinex.NewFundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		fte := bitfinex.FundingTradeExecution(*o)
		return &fte
	case "ftu":
		o, err := bitfinex.NewFundingTradeFromRaw(raw)
		if err != nil {
			return err
		}
		ftu := bitfinex.FundingTradeUpdate(*o)
		return &ftu
	case "hfts":
		o, err := bitfinex.NewFundingTradeSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		nfts := bitfinex.HistoricalFundingTradeSnapshot(*o)
		return &nfts
	case "n":
		o, err := bitfinex.NewNotificationFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fos":
		o, err := bitfinex.NewFundingOfferSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fon":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		fon := bitfinex.FundingOfferNew(*o)
		return &fon
	case "fou":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		fou := bitfinex.FundingOfferUpdate(*o)
		return &fou
	case "foc":
		o, err := bitfinex.NewOfferFromRaw(raw)
		if err != nil {
			return err
		}
		foc := bitfinex.FundingOfferCancel(*o)
		return &foc
	case "fiu":
		o, err := bitfinex.NewFundingInfoFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcs":
		o, err := bitfinex.NewFundingCreditSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fcn":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcn := bitfinex.FundingCreditNew(*o)
		return &fcn
	case "fcu":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcu := bitfinex.FundingCreditUpdate(*o)
		return &fcu
	case "fcc":
		o, err := bitfinex.NewCreditFromRaw(raw)
		if err != nil {
			return err
		}
		fcc := bitfinex.FundingCreditCancel(*o)
		return &fcc
	case "fls":
		o, err := bitfinex.NewFundingLoanSnapshotFromRaw(raw)
		if err != nil {
			return err
		}
		return o
	case "fln":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		fln := bitfinex.FundingLoanNew(*o)
		return &fln
	case "flu":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		flu := bitfinex.FundingLoanUpdate(*o)
		return &flu
	case "flc":
		o, err := bitfinex.NewLoanFromRaw(raw)
		if err != nil {
			return err
		}
		flc := bitfinex.FundingLoanCancel(*o)
		return &flc
	//case "uac":
	case "hb":
		return &bitfinex.Heartbeat{}
	case "ats":
		// TODO: Is not in documentation, so figure out what it is.
		return nil
	case "oc-req":
		// TODO
		return nil
	case "on-req":
		// TODO
		return nil
	case "mis": // Should not be sent anymore as of 2017-04-01
		return nil
	case "miu":
		o, err := bitfinex.NewMarginInfoFromRaw(raw)
		if err != nil {
			return err
		}
		// return a strongly typed reference, rather than dereference a generic interface
		// too bad golang doesn't inherit an interface's underlying type when creating a reference to the interface
		if base, ok := o.(*bitfinex.MarginInfoBase); ok {
			return base
		}
		if update, ok := o.(*bitfinex.MarginInfoUpdate); ok {
			return update
		}
		return o // better than nothing
	default:
		c.log.Warningf("unhandled channel data, term: %s", term)
	}

	return fmt.Errorf("term %q not recognized", term)
}
