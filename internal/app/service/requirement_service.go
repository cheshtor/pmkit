package service

import (
	"pmkit/internal/app/dao"
	"pmkit/internal/app/model"
)

var requirementDao dao.RequirementDao
var contentDao dao.RequirementContentDao
var commentDao dao.RequirementCommentDao
var trackDao dao.RequirementTrackDao

type RequirementService struct {
}

func (s *RequirementService) SaveRequirement(requirement *model.Requirement) (*model.Requirement, error) {
	return nil, nil
}

func (s *RequirementService) EditRequirement(t *model.ThreadLocal, requirement *model.Requirement) (bool, error) {
	return true, nil
}

func (s *RequirementService) ChangeRequirementStatus(t *model.ThreadLocal, id int64, status model.RequirementStatus) (bool, error) {
	return true, nil
}

func (s *RequirementService) RemoveRequirement(t *model.ThreadLocal, id int64) (bool, error) {
	return true, nil
}

func (s *RequirementService) GetRequirementById(id int64) (*model.Requirement, error) {
	return nil, nil
}

func (s *RequirementService) SearchRequirement(condition *model.Requirement, pageNo, pageSize int64) (*model.Page, error) {
	return nil, nil
}

func (s *RequirementService) SaveRequirementTrack(track *model.RequirementTrack) (*model.RequirementTrack, error) {
	return nil, nil
}

func (s *RequirementService) GetAllRequirementTracks(requirementId int64) ([]*model.RequirementTrack, error) {
	return nil, nil
}

func (s *RequirementService) SaveRequirementComment(comment *model.RequirementComment) (*model.RequirementComment, error) {
	return nil, nil
}

func (s *RequirementService) RemoveRequirementComment(commentId int64) (bool, error) {
	return true, nil
}

func (s *RequirementService) GetAllRequirementComment(requirementId int64) ([]*model.RequirementComment, error) {
	return nil, nil
}
