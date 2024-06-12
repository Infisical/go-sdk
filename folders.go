package infisical

import (
	api "github.com/infisical/go-sdk/packages/api/folders"
	"github.com/infisical/go-sdk/packages/models"
)

type ListFoldersOptions = api.ListFoldersV1Request
type UpdateFolderOptions = api.UpdateFolderV1Request
type CreateFolderOptions = api.CreateFolderV1Request
type DeleteFolderOptions = api.DeleteFolderV1Request

type FoldersInterface interface {
	List(options ListFoldersOptions) ([]models.Folder, error)
	Update(options UpdateFolderOptions) (models.Folder, error)
	Create(options CreateFolderOptions) (models.Folder, error)
	Delete(options DeleteFolderOptions) (models.Folder, error)
}

type Folders struct {
	client *InfisicalClient
}

func (s *Folders) List(options ListFoldersOptions) ([]models.Folder, error) {
	res, err := api.CallListFoldersV1(s.client.httpClient, options)

	if err != nil {
		return nil, err
	}

	folders := append([]models.Folder(nil), res.Folders...) // Clone main folders slice, we will modify this if imports are enabled
	return folders, nil
}

func (s *Folders) Update(options UpdateFolderOptions) (models.Folder, error) {
	res, err := api.CallUpdateFolderV1(s.client.httpClient, options)

	if err != nil {
		return models.Folder{}, err
	}

	return res.Folder, nil
}

func (s *Folders) Create(options CreateFolderOptions) (models.Folder, error) {
	res, err := api.CallCreateFolderV1(s.client.httpClient, options)

	if err != nil {
		return models.Folder{}, err
	}

	return res.Folder, nil
}

func (s *Folders) Delete(options DeleteFolderOptions) (models.Folder, error) {
	res, err := api.CallDeleteFolderV1(s.client.httpClient, options)

	if err != nil {
		return models.Folder{}, err
	}

	return res.Folder, nil
}

func NewFolders(client *InfisicalClient) FoldersInterface {
	return &Folders{client: client}
}
