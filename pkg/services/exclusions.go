package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/moq77111113/chmoly-santas/ent"
	"github.com/moq77111113/chmoly-santas/ent/exclusion"
	"github.com/moq77111113/chmoly-santas/ent/group"
	"github.com/moq77111113/chmoly-santas/ent/member"
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

	MemberExcludedBy struct {
		Member     *ent.Member
		ExcludedBy []*ent.Member
	}

	SecretSantaPairing struct {
		Giver    *ent.Member `json:"giver"`
		Receiver *ent.Member `json:"receiver"`
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

func (s *ExclusionRepo) MemberWithExclusions(ctx context.Context, groupId, memberId int) (*GroupExclusions, error) {

	mm, err := s.orm.Member.Query().Where(member.IDEQ(memberId)).Only(ctx)

	if err != nil {
		return nil, err
	}

	exc, err := s.orm.Exclusion.
		Query().
		Where(
			exclusion.GroupIDEQ(groupId),
			exclusion.MemberIDEQ(memberId),
		).
		WithExclude().
		All(ctx)

	if err != nil {
		return nil, err
	}

	excludedMembers := make([]*ent.Member, 0, len(exc))
	for _, e := range exc {
		excludedMembers = append(excludedMembers, e.Edges.Exclude)
	}

	return &GroupExclusions{
		Member:          mm,
		ExcludedMembers: excludedMembers,
	}, nil
}

func (s *ExclusionRepo) MembersExcludedBy(ctx context.Context, groupId, memberId int) (*MemberExcludedBy, error) {

	mm, err := s.orm.Member.Query().Where(member.IDEQ(memberId)).Only(ctx)

	if err != nil {
		return nil, err
	}

	exc, err := s.orm.Exclusion.
		Query().
		Where(
			exclusion.GroupIDEQ(groupId),
			exclusion.ExcludeIDEQ(memberId),
		).
		WithMember().
		All(ctx)

	if err != nil {
		return nil, err
	}

	excludedBy := make([]*ent.Member, 0, len(exc))

	for _, e := range exc {
		excludedBy = append(excludedBy, e.Edges.Member)
	}

	return &MemberExcludedBy{
		Member:     mm,
		ExcludedBy: excludedBy,
	}, nil
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

func (s *ExclusionRepo) GenerateSanta(ctx context.Context, groupId int) ([]*SecretSantaPairing, error) {

	// Fetch members and their exclusions
	mmx, err := s.MembersWithExclusions(ctx, groupId)
	if err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	santasAssignments := make(map[*ent.Member]bool)
	// Initialize santas slice with the correct length
	santas := make([]*ent.Member, len(mmx))

	// Collect members into the santas slice
	for i, mm := range mmx {
		santas[i] = mm.Member
	}

	assignments := make(map[*ent.Member]*ent.Member)

	// Iterate over each member and find a valid Santa for them.
	for _, mm := range mmx {
		var santa *ent.Member
		maxAttempts := len(santas) * 2 // Prevent potential infinite loops
		attempts := 0
		for {
			// If the number of attempts exceeds the maximum allowed, return an error.

			if attempts >= maxAttempts {
				return nil, fmt.Errorf("unable to find a valid santa for %v", mm.Member)
			}

			// Randomly select a Santa from the santas slice.
			santa = santas[r.Intn(len(santas))]

			// Check if the selected Santa is valid:
			// - The Santa is not the member themselves.
			// - The Santa is not in the member's exclusion list.
			// - The Santa has not already been assigned.
			if santa.ID != mm.Member.ID && !isExcluded(santa, mm.ExcludedMembers) && !santasAssignments[santa] {
				assignments[mm.Member] = santa
				santasAssignments[santa] = true
				break
			}
			attempts++
		}
	}

	var pairings []*SecretSantaPairing
	for giver, receiver := range assignments {
		pairings = append(pairings, &SecretSantaPairing{Giver: giver, Receiver: receiver})
	}

	return pairings, nil
}

// Helper function to check if a member is excluded
func isExcluded(member *ent.Member, excluded []*ent.Member) bool {
	for _, ex := range excluded {
		if ex.ID == member.ID {
			return true
		}
	}
	return false
}
