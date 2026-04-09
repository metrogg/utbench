func (t *teamSigchainPlayer) addInnerLink(
	prevState *TeamSigChainState, link *ChainLinkUnpacked, signer SignerX,
	isInflate bool) (
	res checkInnerLinkResult, err error) {

	if link.inner == nil {
		return res, NewStubbedError(link)
	}
	payload := *link.inner

	if !signer.signer.Uid.Exists() {
		return res, NewInvalidLink(link, "empty link signer")
	}

	// This may be superfluous.
	err = link.AssertInnerOuterMatch()
	if err != nil {
		return res, err
	}

	// completely ignore these fields
	_ = payload.ExpireIn
	_ = payload.SeqType

	if payload.Tag != "signature" {
		return res, NewInvalidLink(link, "unrecognized tag: '%s'", payload.Tag)
	}

	if payload.Body.Team == nil {
		return res, NewInvalidLink(link, "missing team section")
	}
	team := payload.Body.Team

	if len(team.ID) == 0 {
		return res, NewInvalidLink(link, "missing team id")
	}
	teamID, err := keybase1.TeamIDFromString(string(team.ID))
	if err != nil {
		return res, err
	}

	if teamID.IsPublic() != team.Public {
		return res, fmt.Errorf("link specified public:%v but ID is public:%v",
			team.Public, teamID.IsPublic())
	}

	if prevState != nil && !prevState.inner.Id.Equal(teamID) {
		return res, fmt.Errorf("wrong team id: %s != %s", teamID.String(), prevState.inner.Id.String())
	}

	if prevState != nil && prevState.IsImplicit() != team.Implicit {
		return res, fmt.Errorf("link specified implicit:%v but team was already implicit:%v",
			team.Implicit, prevState.IsImplicit())
	}

	if prevState != nil && prevState.IsPublic() != team.Public {
		return res, fmt.Errorf("link specified public:%v but team was already public:%v",
			team.Implicit, prevState.IsImplicit())
	}

	if team.Public && (!team.Implicit || (prevState != nil && !prevState.IsImplicit())) {
		return res, fmt.Errorf("public non-implicit teams are not supported")
	}

	err = t.checkSeqnoToAdd(prevState, link.Seqno(), isInflate)
	if err != nil {
		return res, err
	}
	// When isInflate then it is likely prevSeqno != prevState.GetLatestSeqno()
	prevSeqno := link.Seqno() - 1

	allowInflate := func(allow bool) error {
		if isInflate && !allow {
			return fmt.Errorf("inflating link type not supported: %v", payload.Body.Type)
		}
		return nil
	}
	allowInImplicitTeam := func(allow bool) error {
		if team.Implicit && !allow {
			return NewImplicitTeamOperationError(payload.Body.Type)
		}
		return nil
	}
	enforceFirstInChain := func(firstInChain bool) error {
		if firstInChain {
			if prevState != nil {
				return fmt.Errorf("link type '%s' unexpected at seqno:%v", payload.Body.Type, prevState.inner.LastSeqno+1)
			}
		} else {
			if prevState == nil {
				return fmt.Errorf("link type '%s' unexpected at beginning", payload.Body.Type)
			}
		}
		return nil
	}
	enforceGeneric := func(name string, rule Tristate, hasReal bool) error {
		switch rule {
		case TristateDisallow:
			if hasReal {
				return fmt.Errorf("sigchain link contains unexpected '%s'", name)
			}
		case TristateRequire:
			if !hasReal {
				return fmt.Errorf("sigchain link missing %s", name)
			}
		case TristateOptional:
		default:
			return fmt.Errorf("unsupported tristate (fault): %v", rule)
		}
		return nil
	}
	enforce := func(rules LinkRules) error {
		return libkb.PickFirstError(
			enforceGeneric("name", rules.Name, team.Name != nil),
			enforceGeneric("members", rules.Members, team.Members != nil),
			enforceGeneric("parent", rules.Parent, team.Parent != nil),
			enforceGeneric("subteam", rules.Subteam, team.Subteam != nil),
			enforceGeneric("per-team-key", rules.PerTeamKey, team.PerTeamKey != nil),
			enforceGeneric("admin", rules.Admin, team.Admin != nil),
			enforceGeneric("invites", rules.Invites, team.Invites != nil),
			enforceGeneric("completed-invites", rules.CompletedInvites, team.CompletedInvites != nil),
			enforceGeneric("settings", rules.Settings, team.Settings != nil),
			enforceGeneric("kbfs", rules.KBFS, team.KBFS != nil),
			enforceGeneric("box-summary-hash", rules.BoxSummaryHash, team.BoxSummaryHash != nil),
			allowInImplicitTeam(rules.AllowInImplicitTeam),
			allowInflate(rules.AllowInflate),
			enforceFirstInChain(rules.FirstInChain),
		)
	}

	checkAdmin := func(op string) (signerIsExplicitOwner bool, err error) {
		signerRole, err := prevState.GetUserRoleAtSeqno(signer.signer, prevSeqno)
		if err != nil {
			signerRole = keybase1.TeamRole_NONE
		}
		signerIsExplicitOwner = signerRole == keybase1.TeamRole_OWNER
		if signerRole.IsAdminOrAbove() || signer.implicitAdmin {
			return signerIsExplicitOwner, nil
		}
		return signerIsExplicitOwner, fmt.Errorf("link signer does not have permission to %s: %v is a %v", op, signer, signerRole)
	}

	checkExplicitWriter := func(op string) (err error) {
		signerRole, err := prevState.GetUserRoleAtSeqno(signer.signer, prevSeqno)
		if err != nil {
			signerRole = keybase1.TeamRole_NONE
		}
		if !signerRole.IsWriterOrAbove() {
			return fmt.Errorf("link signer does not have writer permission to %s: %v is a %v", op, signer, signerRole)
		}
		return nil
	}

	moveState := func() {
		// Move prevState to res.newState.
		// Re-use the object without deep copying. There must be no other live references into prevState.
		res.newState = *prevState
		prevState = nil
	}
	isHighLink := false

	switch libkb.LinkType(payload.Body.Type) {
	case libkb.LinkTypeTeamRoot:
		isHighLink = true
		err = enforce(LinkRules{
			Name:                TristateRequire,
			Members:             TristateRequire,
			PerTeamKey:          TristateRequire,
			BoxSummaryHash:      TristateOptional,
			Invites:             TristateOptional,
			Settings:            TristateOptional,
			AllowInImplicitTeam: true,
			FirstInChain:        true,
		})
		if err != nil {
			return res, err
		}

		// Check the team name
		teamName, err := keybase1.TeamNameFromString(string(*team.Name))
		if err != nil {
			return res, err
		}
		if !teamName.IsRootTeam() {
			return res, fmt.Errorf("root team has subteam name: %s", teamName)
		}

		// Whether this is an implicit team
		isImplicit := teamName.IsImplicit()
		if isImplicit != team.Implicit {
			return res, fmt.Errorf("link specified implicit:%v but name specified implicit:%v",
				team.Implicit, isImplicit)
		}

		// Check the team ID
		// assert that team_name = hash(team_id)
		// this is only true for root teams
		if !teamID.Equal(teamName.ToTeamID(team.Public)) {
			return res, fmt.Errorf("team id:%s does not match team name:%s", teamID, teamName)
		}
		if teamID.IsSubTeam() {
			return res, fmt.Errorf("malformed root team id")
		}

		roleUpdates, err := t.sanityCheckMembers(*team.Members, sanityCheckMembersOptions{
			requireOwners:       true,
			allowRemovals:       false,
			onlyOwnersOrReaders: isImplicit,
		})
		if err != nil {
			return res, err
		}

		perTeamKey, err := t.checkPerTeamKey(*link.source, *team.PerTeamKey, 1)
		if err != nil {
			return res, err
		}

		perTeamKeys := make(map[keybase1.PerTeamKeyGeneration]keybase1.PerTeamKey)
		perTeamKeys[keybase1.PerTeamKeyGeneration(1)] = perTeamKey

		res.newState = TeamSigChainState{
			inner: keybase1.TeamSigChainState{
				Reader:       t.reader,
				Id:           teamID,
				Implicit:     isImplicit,
				Public:       team.Public,
				RootAncestor: teamName.RootAncestorName(),
				NameDepth:    teamName.Depth(),
				NameLog: []keybase1.TeamNameLogPoint{{
					LastPart: teamName.LastPart(),
					Seqno:    1,
				}},
				LastSeqno:        1,
				LastLinkID:       link.LinkID().Export(),
				ParentID:         nil,
				UserLog:          make(map[keybase1.UserVersion][]keybase1.UserLogPoint),
				SubteamLog:       make(map[keybase1.TeamID][]keybase1.SubteamLogPoint),
				PerTeamKeys:      perTeamKeys,
				PerTeamKeyCTime:  keybase1.UnixTime(payload.Ctime),
				LinkIDs:          make(map[keybase1.Seqno]keybase1.LinkID),
				StubbedLinks:     make(map[keybase1.Seqno]bool),
				ActiveInvites:    make(map[keybase1.TeamInviteID]keybase1.TeamInvite),
				ObsoleteInvites:  make(map[keybase1.TeamInviteID]keybase1.TeamInvite),
				TlfLegacyUpgrade: make(map[keybase1.TeamApplication]keybase1.TeamLegacyTLFUpgradeChainInfo),
				MerkleRoots:      make(map[keybase1.Seqno]keybase1.MerkleRootV2),
			}}

		t.updateMembership(&res.newState, roleUpdates, payload.SignatureMetadata())

		if team.Invites != nil {
			if isImplicit {
				signerIsExplicitOwner := true
				additions, cancelations, err := t.sanityCheckInvites(signer.signer, signerIsExplicitOwner,
					*team.Invites, link.SigID(), sanityCheckInvitesOptions{
						isRootTeam:   true,
						implicitTeam: isImplicit,
					})
				if err != nil {
					return res, err
				}
				t.updateInvites(&res.newState, additions, cancelations)
			} else {
				return res, fmt.Errorf("invites not allowed in root link")
			}
		}

		// check that the signer is an owner
		if res.newState.getUserRole(signer.signer) != keybase1.TeamRole_OWNER {
			return res, fmt.Errorf("signer is not an owner: %v (%v)", signer, team.Members.Owners)
		}

		if settings := team.Settings; settings != nil {
			err = t.parseTeamSettings(settings, &res.newState)
			if err != nil {
				return res, err
			}
		}
	case libkb.LinkTypeChangeMembership:
		err = enforce(LinkRules{
			Members:             TristateRequire,
			PerTeamKey:          TristateOptional,
			Admin:               TristateOptional,
			CompletedInvites:    TristateOptional,
			BoxSummaryHash:      TristateOptional,
			AllowInImplicitTeam: true,
		})
		if err != nil {
			return res, err
		}

		// Check that the signer is at least an ADMIN or is an IMPLICIT ADMIN to have permission to make this link.
		var signerIsExplicitOwner bool
		signerIsExplicitOwner, err = checkAdmin("change membership")
		if err != nil {
			return res, err
		}

		roleUpdates, err := t.sanityCheckMembers(*team.Members, sanityCheckMembersOptions{
			disallowOwners:      prevState.IsSubteam(),
			allowRemovals:       true,
			onlyOwnersOrReaders: prevState.IsImplicit(),
		})
		if err != nil {
			return res, err
		}

		// Only owners can add more owners.
		if (len(roleUpdates[keybase1.TeamRole_OWNER]) > 0) && !signerIsExplicitOwner {
			return res, fmt.Errorf("non-owner cannot add owners")
		}

		// Only owners can remove owners.
		if t.roleUpdatesDemoteOwners(prevState, roleUpdates) && !signerIsExplicitOwner {
			return res, fmt.Errorf("non-owner cannot demote owners")
		}

		if prevState.IsImplicit() {
			// In implicit teams there are only 2 kinds of membership changes allowed:
			// 1. Resolve an invite. Adds 1 user and completes 1 invite.
			//    Though sometimes a new user is not added, due to a conflict.
			// 2. Accept a reset user. Adds 1 user and removes 1 user.
			//    Where the new one has the same UID and role as the old and a greater EldestSeqno.

			// Here's a case that is not straightforward:
			// There is an impteam alice,leland%2,bob@twitter.
			// Leland resets and then proves bob@twitter. On the team chain alice accepts Leland's reset.
			// So she removes leland%2, adds leland%3, and completes the bob@twitter invite.

			// Here's another:
			// There is an impteam leland#bob@twitter.
			// Leland proves bob@twitter. On the team chain leland completes the invite.
			// Now it's just leland.

			// Check that the invites being completed are all active.
			// For non-implicit teams we are more lenient, but here we need the counts to match up.
			invitees := make(map[keybase1.UID]bool)
			parsedCompletedInvites := make(map[keybase1.TeamInviteID]keybase1.UserVersion)
			for inviteID, invitee := range team.CompletedInvites {
				_, ok := prevState.inner.ActiveInvites[inviteID]
				if !ok {
					return res, NewImplicitTeamOperationError("completed invite %v but was not active",
						inviteID)
				}
				uv, err := keybase1.ParseUserVersion(invitee)
				if err != nil {
					return res, err
				}
				invitees[uv.Uid] = true
				parsedCompletedInvites[inviteID] = uv
			}
			nCompleted := len(team.CompletedInvites)

			// Check these two properties:
			// - Every removal must come with an addition of a successor. Ignore role.
			// - Every addition must either be paired with a removal, or resolve an invite. Ignore role.
			// This is a coarse check that ignores role changes.

			type removal struct {
				uv        keybase1.UserVersion
				satisfied bool
			}
			removals := make(map[keybase1.UID]removal)
			for _, uv := range roleUpdates[keybase1.TeamRole_NONE] {
				removals[uv.Uid] = removal{uv: uv}
			}
			var nCompletedExpected int
			additions := make(map[keybase1.UID]bool)
			// Every addition must either be paired with a removal or resolve an invite.
			for _, uv := range append(roleUpdates[keybase1.TeamRole_OWNER], roleUpdates[keybase1.TeamRole_READER]...) {
				removal, ok := removals[uv.Uid]
				if ok {
					if removal.uv.EldestSeqno >= uv.EldestSeqno {
						return res, NewImplicitTeamOperationError("replaced with older eldest seqno: %v -> %v",
							removal.uv.EldestSeqno, uv.EldestSeqno)
					}
					removal.satisfied = true
					removals[uv.Uid] = removal
					if invitees[uv.Uid] && uv.EldestSeqno > removal.uv.EldestSeqno {
						// If we are removing someone that is also a completed invite, then it must
						// be replacing a reset user with a new version. Expect an invite in this case.
						nCompletedExpected++
						additions[uv.Uid] = true
					}
				} else {
					// This is a new user, so must be a completed invite.
					nCompletedExpected++
					additions[uv.Uid] = true
				}
			}
			// All removals must have come with successor.
			for _, r := range removals {
				if !r.satisfied {
					return res, NewImplicitTeamOperationError("removal without addition for %v", r.uv)
				}
			}
			// Completed invites that do not bring in new members mean
			// SBS consolidations.
			for _, uv := range parsedCompletedInvites {
				_, ok := additions[uv.Uid]
				if !ok {
					if prevState.getUserRole(uv) == keybase1.TeamRole_NONE {
						return res, NewImplicitTeamOperationError("trying to moot invite but there is no member for %v", uv)
					}
					nCompleted--
				}
			}
			// The number of completed invites must match.
			if nCompletedExpected != nCompleted {
				return res, NewImplicitTeamOperationError("illegal membership change: %d != %d",
					nCompletedExpected, nCompleted)
			}
		}

		isHighLink, err = t.roleUpdateChangedHighSet(prevState, roleUpdates)
		if err != nil {
			return res, fmt.Errorf("could not determine if high user set changed")
		}

		moveState()
		t.updateMembership(&res.newState, roleUpdates, payload.SignatureMetadata())
		t.completeInvites(&res.newState, team.CompletedInvites)
		t.obsoleteInvites(&res.newState, roleUpdates, payload.SignatureMetadata())

		// Note: If someone was removed, the per-team-key should be rotated. This is not checked though.

		if team.PerTeamKey != nil {
			lastKey, err := res.newState.GetLatestPerTeamKey()
			if err != nil {
				return res, fmt.Errorf("getting previous per-team-key: %s", err)
			}
			newKey, err := t.checkPerTeamKey(*link.source, *team.PerTeamKey, lastKey.Gen+keybase1.PerTeamKeyGeneration(1))
			if err != nil {
				return res, err
			}
			res.newState.inner.PerTeamKeys[newKey.Gen] = newKey
			res.newState.inner.PerTeamKeyCTime = keybase1.UnixTime(payload.Ctime)
		}
	case libkb.LinkTypeRotateKey:
		err = enforce(LinkRules{
			PerTeamKey:          TristateRequire,
			Admin:               TristateOptional,
			BoxSummaryHash:      TristateOptional,
			AllowInImplicitTeam: true,
		})
		if err != nil {
			return res, err
		}

		// Check that the signer is at least a writer to have permission to make this link.
		if !signer.implicitAdmin {
			signerRole, err := prevState.GetUserRoleAtSeqno(signer.signer, prevSeqno)
			if err != nil {
				return res, err
			}
			switch signerRole {
			case keybase1.TeamRole_WRITER, keybase1.TeamRole_ADMIN, keybase1.TeamRole_OWNER:
				// ok
			default:
				return res, fmt.Errorf("link signer does not have permission to rotate key: %v is a %v", signer, signerRole)
			}
		}

		lastKey, err := prevState.GetLatestPerTeamKey()
		if err != nil {
			return res, fmt.Errorf("getting previous per-team-key: %s", err)
		}
		newKey, err := t.checkPerTeamKey(*link.source, *team.PerTeamKey, lastKey.Gen+keybase1.PerTeamKeyGeneration(1))
		if err != nil {
			return res, err
		}

		moveState()
		res.newState.inner.PerTeamKeys[newKey.Gen] = newKey
		res.newState.inner.PerTeamKeyCTime = keybase1.UnixTime(payload.Ctime)
	case libkb.LinkTypeLeave:
		err = enforce(LinkRules{ /* Just about everything is restricted. */ })
		if err != nil {
			return res, err
		}
		// Key rotation should never be allowed since FullVerify sometimes does not run on leave links.

		// Check that the signer is at least a reader.
		// Implicit admins cannot leave a subteam.
		signerRole, err := prevState.GetUserRoleAtSeqno(signer.signer, prevSeqno)
		if err != nil {
			return res, err
		}
		switch signerRole {
		case keybase1.TeamRole_READER, keybase1.TeamRole_WRITER, keybase1.TeamRole_ADMIN, keybase1.TeamRole_OWNER:
			// ok
		default:
			return res, fmt.Errorf("link signer does not have permission to leave: %v is a %v", signer, signerRole)
		}

		// The last owner of a team should not leave.
		// But that's really up to them and the server. We're just reading what has happened.

		moveState()
		res.newState.inform(signer.signer, keybase1.TeamRole_NONE, payload.SignatureMetadata())
	case libkb.LinkTypeNewSubteam:
		err = enforce(LinkRules{
			Subteam:      TristateRequire,
			Admin:        TristateOptional,
			AllowInflate: true,
		})
		if err != nil {
			return res, err
		}

		// Check the subteam ID
		subteamID, err := t.assertIsSubteamID(string(team.Subteam.ID))
		if err != nil {
			return res, err
		}

		// Check the subteam name
		subteamName, err := t.assertSubteamName(prevState, link.Seqno(), string(team.Subteam.Name))
		if err != nil {
			return res, err
		}

		_, err = checkAdmin("make subteam")
		if err != nil {
			return res, err
		}

		moveState()

		// informSubteam will take care of asserting that these links are inflated
		// in order for each subteam.
		err = res.newState.informSubteam(subteamID, subteamName, link.Seqno())
		if err != nil {
			return res, fmt.Errorf("adding new subteam: %v", err)
		}
	case libkb.LinkTypeSubteamHead:
		isHighLink = true

		err = enforce(LinkRules{
			Name:           TristateRequire,
			Members:        TristateRequire,
			Parent:         TristateRequire,
			PerTeamKey:     TristateRequire,
			Admin:          TristateOptional,
			Settings:       TristateOptional,
			BoxSummaryHash: TristateOptional,
			FirstInChain:   true,
		})
		if err != nil {
			return res, err
		}

		if team.Public {
			return res, fmt.Errorf("public subteams are not supported")
		}

		// Check the subteam ID
		if !teamID.IsSubTeam() {
			return res, fmt.Errorf("malformed subteam id")
		}

		// Check parent ID
		parentID, err := keybase1.TeamIDFromString(string(team.Parent.ID))
		if err != nil {
			return res, fmt.Errorf("invalid parent id: %v", err)
		}

		// Check the initial subteam name
		teamName, err := keybase1.TeamNameFromString(string(*team.Name))
		if err != nil {
			return res, err
		}
		if teamName.IsRootTeam() {
			return res, fmt.Errorf("subteam has root team name: %s", teamName)
		}
		if teamName.IsImplicit() || team.Implicit {
			return res, NewImplicitTeamOperationError(payload.Body.Type)
		}

		roleUpdates, err := t.sanityCheckMembers(*team.Members, sanityCheckMembersOptions{
			disallowOwners: true,
			allowRemovals:  false,
		})
		if err != nil {
			return res, err
		}

		perTeamKey, err := t.checkPerTeamKey(*link.source, *team.PerTeamKey, 1)
		if err != nil {
			return res, err
		}

		perTeamKeys := make(map[keybase1.PerTeamKeyGeneration]keybase1.PerTeamKey)
		perTeamKeys[keybase1.PerTeamKeyGeneration(1)] = perTeamKey

		res.newState = TeamSigChainState{
			inner: keybase1.TeamSigChainState{
				Reader:       t.reader,
				Id:           teamID,
				Implicit:     false,
				Public:       false,
				RootAncestor: teamName.RootAncestorName(),
				NameDepth:    teamName.Depth(),
				NameLog: []keybase1.TeamNameLogPoint{{
					LastPart: teamName.LastPart(),
					Seqno:    1,
				}},
				LastSeqno:       1,
				LastLinkID:      link.LinkID().Export(),
				ParentID:        &parentID,
				UserLog:         make(map[keybase1.UserVersion][]keybase1.UserLogPoint),
				SubteamLog:      make(map[keybase1.TeamID][]keybase1.SubteamLogPoint),
				PerTeamKeys:     perTeamKeys,
				PerTeamKeyCTime: keybase1.UnixTime(payload.Ctime),
				LinkIDs:         make(map[keybase1.Seqno]keybase1.LinkID),
				StubbedLinks:    make(map[keybase1.Seqno]bool),
				ActiveInvites:   make(map[keybase1.TeamInviteID]keybase1.TeamInvite),
				ObsoleteInvites: make(map[keybase1.TeamInviteID]keybase1.TeamInvite),
				MerkleRoots:     make(map[keybase1.Seqno]keybase1.MerkleRootV2),
			}}

		t.updateMembership(&res.newState, roleUpdates, payload.SignatureMetadata())
		if settings := team.Settings; settings != nil {
			err = t.parseTeamSettings(settings, &res.newState)
			if err != nil {
				return res, err
			}
		}
	case libkb.LinkTypeRenameSubteam:
		err = enforce(LinkRules{
			Subteam:      TristateRequire,
			Admin:        TristateOptional,
			AllowInflate: true,
		})
		if err != nil {
			return res, err
		}

		_, err = checkAdmin("rename subteam")
		if err != nil {
			return res, err
		}

		// Check the subteam ID
		subteamID, err := t.assertIsSubteamID(string(team.Subteam.ID))
		if err != nil {
			return res, err
		}

		// Check the subteam name
		subteamName, err := t.assertSubteamName(prevState, link.Seqno(), string(team.Subteam.Name))
		if err != nil {
			return res, err
		}

		moveState()

		// informSubteam will take care of asserting that these links are inflated
		// in order for each subteam.
		err = res.newState.informSubteam(subteamID, subteamName, link.Seqno())
		if err != nil {
			return res, fmt.Errorf("adding new subteam: %v", err)
		}
	case libkb.LinkTypeRenameUpPointer:
		err = enforce(LinkRules{
			Name:   TristateRequire,
			Parent: TristateRequire,
			Admin:  TristateOptional,
		})
		if err != nil {
			return res, err
		}

		// These links only occur in subteam.
		if !prevState.IsSubteam() {
			return res, fmt.Errorf("got %v in root team", payload.Body.Type)
		}

		// Sanity check that the parent doesn't claim to have changed.
		parentID, err := keybase1.TeamIDFromString(string(team.Parent.ID))
		if err != nil {
			return res, fmt.Errorf("invalid parent team id: %v", err)
		}
		if !parentID.Eq(*prevState.GetParentID()) {
			return res, fmt.Errorf("wrong parent team ID: %s != %s", parentID, prevState.GetParentID())
		}

		// Ideally we would assert that the name
		// But we may not have an up-to-date picture at this time of the parent's name.
		// So assert this:
		// - The root team name is the same.
		// - The depth of the new name is the same.
		newName, err := keybase1.TeamNameFromString(string(*team.Name))
		if err != nil {
			return res, fmt.Errorf("invalid team name '%s': %v", *team.Name, err)
		}
		if newName.IsRootTeam() {
			return res, fmt.Errorf("cannot rename to root team name: %v", newName.String())
		}
		if !newName.RootAncestorName().Eq(prevState.inner.RootAncestor) {
			return res, fmt.Errorf("rename cannot change root ancestor team name: %v -> %v", prevState.inner.RootAncestor, newName)
		}
		if newName.Depth() != prevState.inner.NameDepth {
			return res, fmt.Errorf("rename cannot change team nesting depth: %v -> %v", prevState.inner.NameDepth, newName)
		}

		moveState()

		res.newState.inner.NameLog = append(res.newState.inner.NameLog, keybase1.TeamNameLogPoint{
			LastPart: newName.LastPart(),
			Seqno:    link.Seqno(),
		})
	case libkb.LinkTypeDeleteSubteam:
		err = enforce(LinkRules{
			Subteam:      TristateRequire,
			Admin:        TristateOptional,
			AllowInflate: true,
		})
		if err != nil {
			return res, err
		}

		_, err = checkAdmin("delete subteam")
		if err != nil {
			return res, err
		}

		// Check the subteam ID
		subteamID, err := t.assertIsSubteamID(string(team.Subteam.ID))
		if err != nil {
			return res, err
		}

		// Check the subteam name
		_, err = t.assertSubteamName(prevState, link.Seqno(), string(team.Subteam.Name))
		if err != nil {
			return res, err
		}

		moveState()

		err = res.newState.informSubteamDelete(subteamID, link.Seqno())
		if err != nil {
			return res, fmt.Errorf("error deleting subteam: %v", err)
		}
	case libkb.LinkTypeInvite:
		err = enforce(LinkRules{
			Admin:               TristateOptional,
			Invites:             TristateRequire,
			AllowInImplicitTeam: true,
		})
		if err != nil {
			return res, err
		}

		signerIsExplicitOwner, err := checkAdmin("invite")
		if err != nil {
			return res, err
		}

		additions, cancelations, err := t.sanityCheckInvites(signer.signer, signerIsExplicitOwner,
			*team.Invites, link.SigID(), sanityCheckInvitesOptions{
				isRootTeam:   !prevState.IsSubteam(),
				implicitTeam: prevState.IsImplicit(),
			})
		if err != nil {
			return res, err
		}

		if prevState.IsImplicit() {
			// Check to see if the additions were previously members of the team
			checkImpteamInvites := func() error {
				addedUIDs := make(map[keybase1.UID]bool)
				for _, invites := range additions {
					for _, invite := range invites {
						cat, err := invite.Type.C()
						if err != nil {
							return err
						}
						if cat == keybase1.TeamInviteCategory_KEYBASE {
							uv, err := invite.KeybaseUserVersion()
							if err != nil {
								return err
							}
							addedUIDs[uv.Uid] = true
							_, err = prevState.GetLatestUVWithUID(uv.Uid)
							if err == nil {
								// Found crypto member in previous
								// state, we are good!
								continue
							}
							_, _, found := prevState.FindActiveKeybaseInvite(uv.Uid)
							if found {
								// Found PUKless member in previous
								// state, still fine!
								continue
							}
							// Neither crypto member nor PUKless member
							// found, we can't allow this addition.
							return fmt.Errorf("Not found previous version of user %s", uv.Uid)
						}
						return fmt.Errorf("invalid invite type in implicit team: %v", cat)
					}
				}

				var cancelledUVs []keybase1.UserVersion
				for _, inviteID := range cancelations {
					invite, found := prevState.FindActiveInviteByID(inviteID)
					if !found {
						// This is harmless and also we might be canceling
						// an obsolete invite.
						continue
					}
					inviteUv, err := invite.KeybaseUserVersion()
					if err != nil {
						return fmt.Errorf("cancelled invite is not valid keybase-type invite: %v", err)
					}
					cancelledUVs = append(cancelledUVs, inviteUv)
				}

				for _, uv := range cancelledUVs {
					if !addedUIDs[uv.Uid] {
						return fmt.Errorf("cancelling invite for %v without inviting back a new version", uv)
					}
				}
				return nil
			}
			if err := checkImpteamInvites(); err != nil {
				return res, NewImplicitTeamOperationError("Error in link %q: %v", payload.Body.Type, err)
			}
		}

		moveState()
		t.updateInvites(&res.newState, additions, cancelations)
	case libkb.LinkTypeSettings:
		err = enforce(LinkRules{
			Admin:    TristateOptional,
			Settings: TristateRequire,
			// Allow key rotation in settings link. Closing an open team
			// should rotate team key.
			PerTeamKey: TristateOptional,
			// At the moment the only team setting is banned in implicit teams.
			// But in the future there could be allowed settings that also use this link type.
			AllowInImplicitTeam: true,
		})
		if err != nil {
			return res, err
		}

		_, err = checkAdmin("change settings")
		if err != nil {
			return res, err
		}

		moveState()
		err = t.parseTeamSettings(team.Settings, &res.newState)
		if err != nil {
			return res, err
		}

		// When team is changed from open to closed, per-team-key should be rotated. But
		// this is not enforced.
		if team.PerTeamKey != nil {
			lastKey, err := res.newState.GetLatestPerTeamKey()
			if err != nil {
				return res, fmt.Errorf("getting previous per-team-key: %s", err)
			}
			newKey, err := t.checkPerTeamKey(*link.source, *team.PerTeamKey, lastKey.Gen+keybase1.PerTeamKeyGeneration(1))
			if err != nil {
				return res, err
			}
			res.newState.inner.PerTeamKeys[newKey.Gen] = newKey
			res.newState.inner.PerTeamKeyCTime = keybase1.UnixTime(payload.Ctime)
		}
	case libkb.LinkTypeDeleteRoot:
		return res, NewTeamDeletedError()
	case libkb.LinkTypeDeleteUpPointer:
		return res, NewTeamDeletedError()
	case libkb.LinkTypeKBFSSettings:
		err = enforce(LinkRules{
			Admin:               TristateOptional,
			KBFS:                TristateRequire,
			AllowInImplicitTeam: true,
		})
		if err != nil {
			return res, err
		}

		err = checkExplicitWriter("change KBFS settings")
		if err != nil {
			return res, err
		}

		moveState()
		err = t.parseKBFSTLFUpgrade(team.KBFS, &res.newState)
		if err != nil {
			return res, err
		}
	case "":
		return res, errors.New("empty body type")
	default:
		if link.outerLink.IgnoreIfUnsupported {
			moveState()
		} else {
			return res, fmt.Errorf("unsupported link type: %s", payload.Body.Type)
		}
	}

	if isHighLink {
		res.newState.inner.LastHighLinkID = link.LinkID().Export()
		res.newState.inner.LastHighSeqno = link.Seqno()
	}
	return res, nil
}
