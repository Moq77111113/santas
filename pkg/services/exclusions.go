package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
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

	assignments := make(map[*ent.Member]*ent.Member)
	for _, mm := range mmx {
		var santa *ent.Member

		for {
			santa = santas[r.Intn(len(santas))]

			if !isExcluded(santa, mm.ExcludedMembers) && !santasAssignments[santa] {
				assignments[mm.Member] = santa
				santasAssignments[santa] = true
				break
			}
		}
	}

	return assignments, nil

}

// GenerateSecretSanta generates the Secret Santa pairings
func GenerateSecretSanta(group []GroupExclusions) ([]SecretSantaPairing, error) {
	if len(group) < 2 {
		return nil, errors.New("not enough participants")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 1. Sort members by the number of exclusions in ascending order
	sort.Slice(group, func(i, j int) bool {
		return len(group[i].ExcludedMembers) < len(group[j].ExcludedMembers)
	})

	// Create a list of all members
	members := make([]*ent.Member, len(group))
	for i, g := range group {
		members[i] = g.Member
	}

	available := make([]*ent.Member, len(members))
	copy(available, members)

	var pairings []SecretSantaPairing
	// 2. Assign each member to a Secret Santa, respecting exclusions
	for _, g := range group {
		// 3. Pick a random receiver from available members, excluding invalid ones
		validReceivers := filterAvailableReceivers(g.Member, g.ExcludedMembers, available)

		if len(validReceivers) == 0 {
			return nil, fmt.Errorf("no valid receivers left for %s", g.Member.Name)
		}

		receiver := validReceivers[r.Intn(len(validReceivers))]

		// 4. Add pairing
		pairings = append(pairings, SecretSantaPairing{Giver: g.Member, Receiver: receiver})

		// Remove the assigned receiver from the available list
		available = removeMember(available, receiver)
	}

	return pairings, nil

}

// Helper function to check if a member is excluded and filter available receivers
func filterAvailableReceivers(giver *ent.Member, excluded []*ent.Member, available []*ent.Member) []*ent.Member {
	validReceivers := []*ent.Member{}
	for _, candidate := range available {
		if giver.ID != candidate.ID && !isExcluded(candidate, excluded) {
			validReceivers = append(validReceivers, candidate)
		}
	}
	return validReceivers
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

// Helper function to remove a member from a slice
func removeMember(members []*ent.Member, member *ent.Member) []*ent.Member {
	for i, m := range members {
		if m.ID == member.ID {
			return append(members[:i], members[i+1:]...)
		}
	}
	return members
}
