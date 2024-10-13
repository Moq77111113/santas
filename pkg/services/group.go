package services

import (
	"context"

	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/ent/exclusion"
	"github.com/moq77111113/chmoly-santas/ent/group"
	"github.com/moq77111113/chmoly-santas/ent/member"
)

type (
	GroupRepo struct {
		orm *ent.Client
	}

	EnrichedWithOwner struct {
		Name  string      `json:"name"`
		Id    int         `json:"id"`
		Owner *ent.Member `json:"owner"`
	}
)

func NewGroupRepo(orm *ent.Client) *GroupRepo {
	return &GroupRepo{orm: orm}
}

func (s *GroupRepo) CreateGroup(ctx context.Context, o *ent.Member, name string) (*ent.Group, error) {

	gr, err := s.orm.Group.Create().SetName(name).SetOwner(o).Save(ctx)

	return gr, err
}

func (s *GroupRepo) List(ctx context.Context) ([]*EnrichedWithOwner, error) {
	grs, err := s.orm.Group.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	var enriched []*EnrichedWithOwner
	for _, g := range grs {
		en, err := enrichWithOwner(g)
		if err != nil {
			return nil, err
		}

		enriched = append(enriched, en)
	}

	return enriched, nil

}
func (s *GroupRepo) Get(ctx context.Context, id int) (*EnrichedWithOwner, error) {
	gr, err := s.orm.Group.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	en, err := enrichWithOwner(gr)
	if err != nil {
		return nil, err
	}

	return en, nil
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

func (s *GroupRepo) AddMember(ctx context.Context, id, member int) error {

	gr, err := s.orm.Group.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = gr.Update().AddMemberIDs(member).Save(ctx)

	return err

}

func (s *GroupRepo) RemoveMember(ctx context.Context, id, memberId int) (*ent.Member, error) {

	mm, err := s.orm.Group.Query().Where(group.IDEQ(id)).WithMembers().QueryMembers().Where(member.IDEQ(memberId)).Only(ctx)

	if err != nil {
		return nil, err
	}

	mm, err = mm.Update().RemoveGroupIDs(id).Save(ctx)

	s.orm.Exclusion.Delete().Where(exclusion.GroupID(id), exclusion.MemberID(memberId)).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mm, nil
}

func (s *GroupRepo) CreateMember(ctx context.Context, name string) (*ent.Member, error) {
	return s.orm.Member.Create().SetName(name).Save(ctx)
}

func enrichWithOwner(g *ent.Group) (*EnrichedWithOwner, error) {
	owner, err := g.QueryOwner().Only(context.Background())
	if err != nil {
		return nil, err
	}

	return &EnrichedWithOwner{
		Name:  g.Name,
		Id:    g.ID,
		Owner: owner,
	}, nil
}
