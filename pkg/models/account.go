package models

import (
	"github.com/go-pg/pg/v10/orm"
)

type Account struct {
	tableName string `pg:"accounts"`

	AccountID     uint64 `json:"account_id" pg:"account_id,type:serial"`
	PaymentWallet string `json:"payment_wallet" pg:"payment_wallet,pk"`
	Nonce         string `json:"nonce" pg:"nonce"`
}

func (db *Client) CreateAccountsSchema() error {

	err := db.Model((*Account)(nil)).CreateTable(
		&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true})
	if err != nil {
		return ErrTableCreationFailed
	}

	return nil
}
