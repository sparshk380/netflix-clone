package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10/orm"
)

type Nodes struct {
	tableName string `pg:"nodes"`

	NodeID           uint64 `json:"node_id" pg:"node_id,pk,type:serial"`
	NodeName         string `json:"node_name" pg:"node_name,unique"`
	AccountID        uint64 `json:"account_id" pg:"account_id"`
	ChainName        string `json:"chain_name" pg:"chain_name"`
	ChainType        string `json:"chain_type" pg:"chain_type"`
	Status           string `json:"status" pg:"status"`
	PublicIP         string `json:"public_ip" pg:"public_ip"`
	RPCPort          string `json:"rpc_port" pg:"rpc_port"`
	RPCWebSocketPort string `json:"rpc_websocket_port" pg:"rpc_websocket_port"`
	RestHTTPPort     string `json:"rest_http_port" pg:"rest_http_port"`
	Location         string `json:"location" pg:"location"`

	CreatedAt time.Time `json:"-" pg:"created_at,type:timestamp with time zone"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,type:timestamp with time zone"`
}

func (db *Client) CreateNodesSchema() error {

	err := db.Model((*Nodes)(nil)).CreateTable(
		&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true})
	if err != nil {
		return ErrTableCreationFailed
	}

	return nil
}

func (u *Nodes) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	return ctx, nil
}

func (u *Nodes) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.UpdatedAt = time.Now().UTC()

	return ctx, nil
}
