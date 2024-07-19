package common

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

func StringToUint64(str string) (uint64, error) {
	uintVar, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logrus.Error("failed to convert string to uint. Err: ", err.Error())
		return 0, ErrInvalidID
	}
	return uintVar, nil
}

func ConvertInterfaceToUint64(accountIDInterface interface{}) (uint64, error) {
	aidUint64, ok := accountIDInterface.(uint64)
	if !ok {
		logrus.Errorf("failed to convert accountid interface to uint64. accountIDInterface: %v", accountIDInterface)
		return 0, ErrUnauthorized
	}

	return aidUint64, nil
}
