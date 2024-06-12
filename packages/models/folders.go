package models

import "time"

type Folder struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAT  time.Time `json:"updatedAt"`
	EnvId      string    `json:"envId"`
	ParentID   string    `json:"parentId"`
	IsReserved bool      `json:"isReserved"`
}
