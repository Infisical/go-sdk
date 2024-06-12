package api

import (
	"github.com/infisical/go-sdk/packages/models"
)

// List folders
type ListFoldersV1Request struct {
	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	Path        string `json:"path,omitempty"`
}

type ListFoldersV1Response struct {
	Folders []models.Folder `json:"folders"`
}

// Update folder
type UpdateFolderV1Request struct {
	folderID string `json:"-"`

	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
}

type UpdateFolderV1Response struct {
	Folder models.Folder `json:"folder"`
}

// Create folder
type CreateFolderV1Request struct {
	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
}

type CreateFolderV1Response struct {
	Folder models.Folder `json:"folder"`
}

// Delete folder
type DeleteFolderV1Request struct {
	folderID string `json:"-"`

	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
}

type DeleteFolderV1Response struct {
	Folder models.Folder `json:"folder"`
}
