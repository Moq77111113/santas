package services

import (
	"context"

	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/ent/exclusion"
)

type (
	ExclusionRepo struct {
		orm *ent.Client
	}

	AddExclusion struct {
		GroupId   int
		MemberId  int
		ExcludeId int
	}

	MemberExclusions struct {
		Member          *ent.Member   `json:"member"`
		ExcludedMembers []*ent.Member `json:"excludedMembers"`
	}
)

func NewExclusionRepo(orm *ent.Client) *ExclusionRepo {
	return &ExclusionRepo{orm: orm}
}

func (s *ExclusionRepo) AddExclusion(ctx context.Context, payload AddExclusion) (*ent.Exclusion, error) {
	return s.orm.Exclusion.
		Create().
		SetGroupID(payload.GroupId).
		SetMemberID(payload.MemberId).
		SetExcludeID(payload.ExcludeId).
		Save(ctx)
}

func (s *ExclusionRepo) GetExclusions(ctx context.Context, groupId int) ([]*MemberExclusions, error) {
	exc, err := s.orm.Exclusion.
		Query().
		Where(exclusion.GroupIDEQ(groupId)).
		WithMember().
		WithExclude().
		All(ctx)

	if err != nil {
		return nil, err
	}

	return toMemberExclusions(exc), nil

}

func (s *ExclusionRepo) RemoveExclusion(ctx context.Context, id int) error {
	return s.orm.Exclusion.DeleteOneID(id).Exec(ctx)
}

func toMemberExclusions(mms []*ent.Exclusion) []*MemberExclusions {
	members := make(map[int]*ent.Member)
	excludedMembers := make([]*ent.Member, 0, len(mms))
	for _, mm := range mms {
		if _, ok := members[mm.MemberID]; !ok {
			members[mm.MemberID] = mm.Edges.Member
		}
		excludedMembers = append(excludedMembers, mm.Edges.Exclude)
	}

	res := make([]*MemberExclusions, 0, len(members))
	for _, member := range members {
		res = append(res, &MemberExclusions{
			Member:          member,
			ExcludedMembers: excludedMembers,
		})
	}

	return res
}
