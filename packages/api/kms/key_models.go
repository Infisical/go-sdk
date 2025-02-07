package api

import "time"

type KmsKey struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Algorithm   string     `json:"algorithm"`
	ProjectID   string     `json:"projectId"`
	IsDisabled  bool       `json:"isDisabled"`
	IsReserved  bool       `json:"isReserved"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type KmsListKeysV1Request struct {
	ProjectID string `json:"projectId"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	OrderBy   string `json:"orderBy,omitempty"`
	OrderDir  string `json:"orderDir,omitempty"`
	Search    string `json:"search,omitempty"`
}

type ListKmsKeysV1Response struct {
	Keys []KmsKey `json:"keys"`
}

type KmsCreateKeyV1Request struct {
	Name        string `json:"name"`
	ProjectID   string `json:"projectId"`
	Description string `json:"description,omitempty"`
	Algorithm   string `json:"algorithm"`
}

type CreateKmsKeyV1Response struct {
	Key KmsKey `json:"key"`
}

type KmsUpdateKeyV1Request struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateKmsKeyV1Response struct {
	Key KmsKey `json:"key"`
}

type KmsDeleteKeyV1Request struct {
	ID string `json:"id"`
}

type DeleteKmsKeyV1Response struct {
	Key KmsKey `json:"key"`
}
