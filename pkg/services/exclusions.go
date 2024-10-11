package services

import (
	"context"
	"math/rand"
	"time"

	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/ent/exclusion"
	"github.com/moq77111113/chmoly-santas/ent/group"
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

	GroupExclusions struct {
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

func (s *ExclusionRepo) MembersWithExclusions(ctx context.Context, groupId int) ([]*GroupExclusions, error) {
	mms, err := s.orm.Group.Query().Where(group.IDEQ(groupId)).QueryMembers().All(ctx)
	if err != nil {
		return nil, err
	}

	if len(mms) == 0 {
		return nil, nil
	}

	exc, err := s.orm.Exclusion.
		Query().
		Where(exclusion.GroupIDEQ(groupId)).
		WithMember().
		WithExclude().
		All(ctx)

	if err != nil {
		return nil, err
	}

	mmsWithExclusions := make([]*GroupExclusions, 0, len(mms))
	for _, mm := range mms {
		excludedMembers := make([]*ent.Member, 0, len(exc))
		for _, e := range exc {
			if e.MemberID == mm.ID {
				excludedMembers = append(excludedMembers, e.Edges.Exclude)
			}
		}
		mmsWithExclusions = append(mmsWithExclusions, &GroupExclusions{
			Member:          mm,
			ExcludedMembers: excludedMembers,
		})
	}

	return mmsWithExclusions, nil

}

func (s *ExclusionRepo) RemoveExclusion(ctx context.Context, groupId, memberId, excludeId int) error {
	_, err := s.orm.Exclusion.
		Delete().
		Where(
			exclusion.GroupIDEQ(groupId),
			exclusion.MemberIDEQ(memberId),
			exclusion.ExcludeIDEQ(excludeId),
		).
		Exec(ctx)
	return err
}

func (s *ExclusionRepo) GenerateSanta(ctx context.Context, groupId int) (map[*ent.Member]*ent.Member, error) {
	mmx, err := s.MembersWithExclusions(ctx, groupId)
	if err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	santasAssignments := make(map[*ent.Member]bool)
	santas := make([]*ent.Member, 0, len(mmx))

	for i, mm := range mmx {
		santas[i] = mm.Member
	}

	assignements := make(map[*ent.Member]*ent.Member)
	for _, mm := range mmx {
		var santa *ent.Member

		for {
			santa = santas[r.Intn(len(santas))]

			if !isExcluded(santa, mm.ExcludedMembers) && !santasAssignments[santa] {
				assignements[mm.Member] = santa
				santasAssignments[santa] = true
				break
			}
		}
	}

	return assignements, nil

}

func isExcluded(santa *ent.Member, excludedMembers []*ent.Member) bool {
	for _, excluded := range excludedMembers {
		if santa.ID == excluded.ID { // Assuming Member has an ID field to compare
			return true
		}
	}
	return false
}
