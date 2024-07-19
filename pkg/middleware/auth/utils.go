package auth

func convertInterfaceToString(walletAddressInterface interface{}) string {
	walletAddressString := walletAddressInterface.(string)

	return walletAddressString
}
