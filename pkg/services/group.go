package services

import (
	"context"

	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/ent/group"
	"github.com/moq77111113/chmoly-santas/ent/member"
)

type (
	GroupRepo struct {
		orm *ent.Client
	}
)

func NewGroupRepo(orm *ent.Client) *GroupRepo {
	return &GroupRepo{orm: orm}
}

func (s *GroupRepo) CreateGroup(ctx context.Context, name string) (*ent.Group, error) {

	gr, err := s.orm.Group.Create().SetName(name).Save(ctx)

	return gr, err
}

func (s *GroupRepo) Get(ctx context.Context, id int) (*ent.Group, error) {
	return s.orm.Group.Get(ctx, id)
}

func (s *GroupRepo) GetMembers(ctx context.Context, id int) ([]*ent.Member, error) {
	return s.orm.Group.Query().Where(group.IDEQ(id)).WithMembers().QueryMembers().All(ctx)
}

func (s *GroupRepo) Remove(ctx context.Context, name string) error {

	gr, err := s.orm.Group.Query().Where(group.NameEqualFold(name)).Only(ctx)
	if err != nil {
		return err
	}

	err = s.orm.Group.DeleteOne(gr).Exec(ctx)

	return err
}

func (s *GroupRepo) AddMember(ctx context.Context, id int, memberName string) (*ent.Member, error) {

	gr, err := s.orm.Group.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mm, err := s.orm.Member.Query().Where(member.NameEqualFold(memberName)).Only(ctx)

	if err != nil {
		mm, err = s.orm.Member.Create().SetName(memberName).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	_, err = gr.Update().AddMembers(mm).Save(ctx)

	return mm, err
}

func (s *GroupRepo) RemoveMember(ctx context.Context, id, memberId int) (*ent.Member, error) {

	mm, err := s.orm.Group.Query().Where(group.IDEQ(id)).WithMembers().QueryMembers().Where(member.IDEQ(memberId)).Only(ctx)

	if err != nil {
		return nil, err
	}

	mm, err = mm.Update().RemoveGroupIDs(id).Save(ctx)

	if err != nil {
		return nil, err
	}

	return mm, nil
}
