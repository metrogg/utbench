package main

func Tokens() []string {
	tokens := make([]string, 0, len(BtcIndices)+len(DcrExchanges))
	var token string
	for token = range BtcIndices {
		tokens = append(tokens, token)
	}
	for token = range DcrExchanges {
		tokens = append(tokens, token)
	}
	return tokens
}
