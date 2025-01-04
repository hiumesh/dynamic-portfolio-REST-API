package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryAttachment interface {
	Create(userId string, parentTable string, parentId uint, attachment *schemas.Attachment) error
	CreateMany(userId string, parentTable string, parentId uint, attachments *schemas.Attachments) error
	UpdateOrCreate(userId string, parentTable string, parentId uint, attachments *schemas.Attachments) error
}

type repositoryAttachment struct {
	db *gorm.DB
}

func (r *repositoryAttachment) Create(userId string, parentTable string, parentId uint, attachment *schemas.Attachment) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed to parse user id")
	}

	var att = models.Attachment{
		UserId:      userUUID,
		ParentTable: parentTable,
		ParentId:    parentId,
		FileName:    attachment.FileName,
		FileUrl:     attachment.FileUrl,
		FileType:    attachment.FileType,
		FileSize:    attachment.FileSize,
	}

	return r.db.Create(att).Error
}

func (r *repositoryAttachment) CreateMany(userId string, parentTable string, parentId uint, attachments *schemas.Attachments) (*models.Attachments, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	var atts = make(models.Attachments, 0)

	for _, attachment := range *attachments {
		atts = append(atts, models.Attachment{
			UserId:      userUUID,
			ParentTable: parentTable,
			ParentId:    parentId,
			FileName:    attachment.FileName,
			FileUrl:     attachment.FileUrl,
			FileType:    attachment.FileType,
			FileSize:    attachment.FileSize,
		})
	}

	if err := r.db.Create(&atts).Error; err != nil {
		return nil, err
	}

	return &atts, nil
}

func (r *repositoryAttachment) UpdateOrCreate(userId string, parentTable string, parentId uint, attachments *schemas.Attachments) (interface{}, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	var oldAttachments models.Attachments

	if err := r.db.Where("user_id = ? and parent_table = ? and parent_id = ?", userUUID, parentTable, parentId).Find(&oldAttachments).Error; err != nil {
		return nil, err
	}

	var oldAttachmentsMap = make(map[string]models.Attachment)
	var alreadyExistsAttachmentsMap = make(map[string]models.Attachment)

	for _, oldAttachment := range oldAttachments {
		oldAttachmentsMap[oldAttachment.FileUrl] = oldAttachment
	}

	logrus.Info(oldAttachmentsMap)

	var newAttachments models.Attachments

	for _, attachment := range *attachments {
		if _, ok := oldAttachmentsMap[attachment.FileUrl]; ok {
			alreadyExistsAttachmentsMap[attachment.FileUrl] = oldAttachmentsMap[attachment.FileUrl]
			delete(oldAttachmentsMap, attachment.FileUrl)
		} else {
			newAttachments = append(newAttachments, models.Attachment{
				UserId:      userUUID,
				ParentTable: parentTable,
				ParentId:    parentId,
				FileName:    attachment.FileName,
				FileUrl:     attachment.FileUrl,
				FileType:    attachment.FileType,
				FileSize:    attachment.FileSize,
			})
		}
	}

	if len(oldAttachmentsMap) > 0 {
		var oldAttachmentIds []uint
		for _, oldAttachment := range oldAttachmentsMap {
			oldAttachmentIds = append(oldAttachmentIds, oldAttachment.ID)
		}

		if err := r.db.Unscoped().Delete(&models.Attachment{}, "id in (?)", oldAttachmentIds).Error; err != nil {
			return nil, err
		}
	}

	if len(newAttachments) > 0 {
		if err := r.db.Create(&newAttachments).Error; err != nil {
			return nil, err
		}
	}

	var newAttachmentsMap = make(map[string]models.Attachment)

	for _, newAttachment := range newAttachments {
		newAttachmentsMap[newAttachment.FileUrl] = newAttachment
	}

	var responseAttachments models.Attachments

	for _, attachment := range *attachments {
		if _, ok := newAttachmentsMap[attachment.FileUrl]; ok {
			responseAttachments = append(responseAttachments, newAttachmentsMap[attachment.FileUrl])
		}
		if _, ok := alreadyExistsAttachmentsMap[attachment.FileUrl]; ok {
			responseAttachments = append(responseAttachments, alreadyExistsAttachmentsMap[attachment.FileUrl])
		}
	}

	return &responseAttachments, nil
}

func (r *repositoryAttachment) DeleteMany(userId string, parentTable string, parentId string) error {
	if err := r.db.Unscoped().Where("user_id = ? and parent_table = ? and parent_id = ?", userId, parentTable, parentId).Delete(&models.Attachment{}).Error; err != nil {
		return err
	}

	return nil
}

func NewAttachmentRepository(db *gorm.DB) *repositoryAttachment {
	return &repositoryAttachment{
		db: db,
	}
}
