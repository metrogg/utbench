func ImportStatusAsError(g *GlobalContext, s *keybase1.Status) error {
	if s == nil {
		return nil
	}
	switch s.Code {
	case SCOk:
		return nil
	case SCGeneric:
		return errors.New(s.Desc)
	case SCBadSession:
		return BadSessionError{s.Desc}
	case SCBadLoginPassword:
		return PassphraseError{s.Desc}
	case SCKeyBadGen:
		return KeyGenError{s.Desc}
	case SCAlreadyLoggedIn:
		return LoggedInError{}
	case SCCanceled:
		return CanceledError{s.Desc}
	case SCInputCanceled:
		return InputCanceledError{}
	case SCKeyNoSecret:
		return NoSecretKeyError{}
	case SCLoginRequired:
		return LoginRequiredError{s.Desc}
	case SCNoSession:
		return NoSessionError{}
	case SCKeyCorrupted:
		return KeyCorruptedError{s.Desc}
	case SCOffline:
		return OfflineError{}
	case SCKeyInUse:
		var fp *PGPFingerprint
		if len(s.Desc) > 0 {
			fp, _ = PGPFingerprintFromHex(s.Desc)
		}
		return KeyExistsError{fp}
	case SCKeyNotFound:
		return NoKeyError{s.Desc}
	case SCKeyNoEldest:
		return NoSigChainError{}
	case SCStreamExists:
		return StreamExistsError{}
	case SCBadInvitationCode:
		return BadInvitationCodeError{}
	case SCStreamNotFound:
		return StreamNotFoundError{}
	case SCStreamWrongKind:
		return StreamWrongKindError{}
	case SCStreamEOF:
		return io.EOF
	case SCSelfNotFound:
		return SelfNotFoundError{msg: s.Desc}
	case SCDeviceNotFound:
		return NoDeviceError{Reason: s.Desc}
	case SCDecryptionKeyNotFound:
		return NoDecryptionKeyError{Msg: s.Desc}
	case SCTimeout:
		return TimeoutError{}
	case SCDeviceMismatch:
		return ReceiverDeviceError{Msg: s.Desc}
	case SCBadKexPhrase:
		return InvalidKexPhraseError{}
	case SCReloginRequired:
		return ReloginRequiredError{}
	case SCDeviceRequired:
		return DeviceRequiredError{}
	case SCMissingResult:
		return IdentifyDidNotCompleteError{}
	case SCSibkeyAlreadyExists:
		return SibkeyAlreadyExistsError{}
	case SCSigCreationDisallowed:
		service := ""
		if len(s.Fields) > 0 && s.Fields[0].Key == "remote_service" {
			service = s.Fields[0].Value
		}
		return ServiceDoesNotSupportNewProofsError{Service: service}
	case SCNoUI:
		return NoUIError{Which: s.Desc}
	case SCNoUIDelegation:
		return UIDelegationUnavailableError{}
	case SCProfileNotPublic:
		return ProfileNotPublicError{msg: s.Desc}
	case SCIdentifyFailed:
		var assertion string
		if len(s.Fields) > 0 && s.Fields[0].Key == "assertion" {
			assertion = s.Fields[0].Value
		}
		return IdentifyFailedError{Assertion: assertion, Reason: s.Desc}
	case SCIdentifiesFailed:
		return IdentifiesFailedError{}
	case SCIdentifySummaryError:
		ret := IdentifySummaryError{}
		for _, pair := range s.Fields {
			if pair.Key == "username" {
				ret.username = NewNormalizedUsername(pair.Value)
			} else {
				// The other keys are expected to be "problem_%d".
				ret.problems = append(ret.problems, pair.Value)
			}
		}
		return ret
	case SCTrackingBroke:
		return TrackingBrokeError{}
	case SCResolutionFailed:
		var input string
		if len(s.Fields) > 0 && s.Fields[0].Key == "input" {
			input = s.Fields[0].Value
		}
		return ResolutionError{Msg: s.Desc, Input: input}
	case SCAccountReset:
		var e keybase1.UserVersion
		var r keybase1.Seqno
		seqnoFromString := func(s string) keybase1.Seqno {
			i, _ := strconv.Atoi(s)
			return keybase1.Seqno(i)
		}
		for _, field := range s.Fields {
			switch field.Key {
			case "e_uid":
				e.Uid, _ = keybase1.UIDFromString(field.Value)
			case "e_version":
				e.EldestSeqno = seqnoFromString(field.Value)
			case "r_version":
				r = seqnoFromString(field.Value)
			}
		}
		return NewAccountResetError(e, r)
	case SCKeyNoPGPEncryption:
		ret := NoPGPEncryptionKeyError{User: s.Desc}
		for _, field := range s.Fields {
			switch field.Key {
			case "HasKeybaseEncryptionKey":
				ret.HasKeybaseEncryptionKey = true
			}
		}
		return ret
	case SCKeyNoNaClEncryption:
		ret := NoNaClEncryptionKeyError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "Username":
				ret.Username = field.Value
			case "HasPGPKey":
				ret.HasPGPKey = true
			case "HasPUK":
				ret.HasPUK = true
			case "HasPaperKey":
				ret.HasPaperKey = true
			case "HasDeviceKey":
				ret.HasDeviceKey = true
			}
		}
		return ret
	case SCWrongCryptoFormat:
		ret := WrongCryptoFormatError{Operation: s.Desc}
		for _, field := range s.Fields {
			switch field.Key {
			case "wanted":
				ret.Wanted = CryptoMessageFormat(field.Value)
			case "received":
				ret.Received = CryptoMessageFormat(field.Value)
			}
		}
		return ret
	case SCKeySyncedPGPNotFound:
		return NoSyncedPGPKeyError{}
	case SCKeyNoMatchingGPG:
		ret := NoMatchingGPGKeysError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "fingerprints":
				ret.Fingerprints = strings.Split(field.Value, ",")
			case "has_active_device":
				ret.HasActiveDevice = true
			}
		}
		return ret
	case SCDevicePrevProvisioned:
		return DeviceAlreadyProvisionedError{}
	case SCDeviceProvisionViaDevice:
		return ProvisionViaDeviceRequiredError{}
	case SCDeviceNoProvision:
		return ProvisionUnavailableError{}
	case SCGPGUnavailable:
		return GPGUnavailableError{}
	case SCNotFound:
		return NotFoundError{Msg: s.Desc}
	case SCDeleted:
		return UserDeletedError{Msg: s.Desc}
	case SCDecryptionError:
		ret := DecryptionError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "Cause":
				ret.Cause = fmt.Errorf(field.Value)
			}
		}
		return ret
	case SCSigCannotVerify:
		ret := kbcrypto.VerificationError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "Cause":
				ret.Cause = fmt.Errorf(field.Value)
			}
		}
		return ret
	case SCKeyRevoked:
		return KeyRevokedError{msg: s.Desc}
	case SCDeviceNameInUse:
		return DeviceNameInUseError{}
	case SCDeviceBadName:
		return DeviceBadNameError{}
	case SCGenericAPIError:
		var code int
		for _, field := range s.Fields {
			switch field.Key {
			case "code":
				var err error
				code, err = strconv.Atoi(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing generic API error code: %s", err)
				}
			}
		}
		return &APIError{
			Msg:  s.Desc,
			Code: code,
		}
	case SCChatInternal:
		return ChatInternalError{}
	case SCChatStalePreviousState:
		return ChatStalePreviousStateError{}
	case SCChatEphemeralRetentionPolicyViolatedError:
		var maxAge gregor1.DurationSec
		for _, field := range s.Fields {
			switch field.Key {
			case "MaxAge":
				dur, err := time.ParseDuration(field.Value)
				if err == nil {
					maxAge = gregor1.ToDurationSec(dur)
				}
				break
			}
		}
		return ChatEphemeralRetentionPolicyViolatedError{maxAge}
	case SCChatConvExists:
		var convID chat1.ConversationID
		for _, field := range s.Fields {
			switch field.Key {
			case "ConvID":
				bs, err := chat1.MakeConvID(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing ChatConvExistsError")
				}
				convID = chat1.ConversationID(bs)
			}
		}
		return ChatConvExistsError{
			ConvID: convID,
		}
	case SCChatUnknownTLFID:
		var tlfID chat1.TLFID
		for _, field := range s.Fields {
			switch field.Key {
			case "TlfID":
				var err error
				tlfID, err = chat1.MakeTLFID(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing chat unknown TLF ID error")
				}
			}
		}
		return ChatUnknownTLFIDError{
			TlfID: tlfID,
		}
	case SCChatNotInConv:
		var uid gregor1.UID
		for _, field := range s.Fields {
			switch field.Key {
			case "UID":
				val, err := hex.DecodeString(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing chat not in conv UID")
				}
				uid = gregor1.UID(val)
			}
		}
		return ChatNotInConvError{
			UID: uid,
		}
	case SCChatNotInTeam:
		var uid gregor1.UID
		for _, field := range s.Fields {
			switch field.Key {
			case "UID":
				val, err := hex.DecodeString(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing chat not in conv UID")
				}
				uid = gregor1.UID(val)
			}
		}
		return ChatNotInTeamError{
			UID: uid,
		}
	case SCChatTLFFinalized:
		var tlfID chat1.TLFID
		for _, field := range s.Fields {
			switch field.Key {
			case "TlfID":
				var err error
				tlfID, err = chat1.MakeTLFID(field.Value)
				if err != nil && g != nil {
					g.Log.Warning("error parsing chat tlf finalized TLFID: %s", err.Error())
				}
			}
		}
		return ChatTLFFinalizedError{
			TlfID: tlfID,
		}
	case SCChatBadMsg:
		return ChatBadMsgError{Msg: s.Desc}
	case SCChatBroadcast:
		return ChatBroadcastError{Msg: s.Desc}
	case SCChatRateLimit:
		var rlimit chat1.RateLimit
		for _, field := range s.Fields {
			switch field.Key {
			case "RateLimit":
				var err error
				err = json.Unmarshal([]byte(field.Value), &rlimit)
				if err != nil && g != nil {
					g.Log.Warning("error parsing chat rate limit: %s", err.Error())
				}
			}
		}
		if rlimit.Name == "" && g != nil {
			g.Log.Warning("error rate limit information not found")
		}
		return ChatRateLimitError{
			RateLimit: rlimit,
			Msg:       s.Desc,
		}
	case SCChatAlreadySuperseded:
		return ChatAlreadySupersededError{Msg: s.Desc}
	case SCChatAlreadyDeleted:
		return ChatAlreadyDeletedError{Msg: s.Desc}
	case SCBadEmail:
		return BadEmailError{Msg: s.Desc}
	case SCExists:
		return ExistsError{Msg: s.Desc}
	case SCInvalidAddress:
		return InvalidAddressError{Msg: s.Desc}
	case SCChatCollision:
		return ChatCollisionError{}
	case SCChatMessageCollision:
		var headerHash string
		for _, field := range s.Fields {
			switch field.Key {
			case "HeaderHash":
				headerHash = field.Value
			}
		}
		return ChatMessageCollisionError{
			HeaderHash: headerHash,
		}
	case SCChatDuplicateMessage:
		var soutboxID string
		for _, field := range s.Fields {
			switch field.Key {
			case "OutboxID":
				soutboxID = field.Value
			}
		}
		boutboxID, _ := hex.DecodeString(soutboxID)
		return ChatDuplicateMessageError{
			OutboxID: chat1.OutboxID(boutboxID),
		}
	case SCChatClientError:
		return ChatClientError{Msg: s.Desc}
	case SCNeedSelfRekey:
		ret := NeedSelfRekeyError{Msg: s.Desc}
		for _, field := range s.Fields {
			switch field.Key {
			case "Tlf":
				ret.Tlf = field.Value
			}
		}
		return ret
	case SCNeedOtherRekey:
		ret := NeedOtherRekeyError{Msg: s.Desc}
		for _, field := range s.Fields {
			switch field.Key {
			case "Tlf":
				ret.Tlf = field.Value
			}
		}
		return ret
	case SCLoginStateTimeout:
		var e LoginStateTimeoutError
		for _, field := range s.Fields {
			switch field.Key {
			case "ActiveRequest":
				e.ActiveRequest = field.Value
			case "AttemptedRequest":
				e.AttemptedRequest = field.Value
			case "Duration":
				dur, err := time.ParseDuration(field.Value)
				if err == nil {
					e.Duration = dur
				}
			}
		}
		return e
	case SCRevokeCurrentDevice:
		return RevokeCurrentDeviceError{}
	case SCRevokeLastDevice:
		return RevokeLastDeviceError{}
	case SCRevokeLastDevicePGP:
		return RevokeLastDevicePGPError{}
	case SCTeamKeyMaskNotFound:
		e := KeyMaskNotFoundError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "App":
				e.App = keybase1.TeamApplication(field.IntValue())
			case "Gen":
				e.Gen = keybase1.PerTeamKeyGeneration(field.IntValue())
			}
		}
		return e
	case SCDeviceProvisionOffline:
		return ProvisionFailedOfflineError{}
	case SCGitInvalidRepoName:
		e := InvalidRepoNameError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "Name":
				e.Name = field.Value
			}
		}
		return e
	case SCGitRepoAlreadyExists:
		e := RepoAlreadyExistsError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "DesiredName":
				e.DesiredName = field.Value
			case "ExistingName":
				e.ExistingName = field.Value
			case "ExistingID":
				e.ExistingID = field.Value
			}
		}
		return e
	case SCGitRepoDoesntExist:
		e := RepoDoesntExistError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "Name":
				e.Name = field.Value
			}
		}
		return e
	case SCNoOp:
		return NoOpError{Desc: s.Desc}
	case SCNoSpaceOnDevice:
		return NoSpaceOnDeviceError{Desc: s.Desc}
	case SCTeamInviteBadToken:
		return TeamInviteBadTokenError{}
	case SCTeamInviteTokenReused:
		return TeamInviteTokenReusedError{}
	case SCTeamBadMembership:
		return TeamBadMembershipError{}
	case SCTeamProvisionalCanKey, SCTeamProvisionalCannotKey:
		e := TeamProvisionalError{}
		for _, field := range s.Fields {
			switch field.Key {
			case "IsPublic":
				if field.Value == "1" {
					e.IsPublic = true
				}
			case "PreResolveDisplayName":
				e.PreResolveDisplayName = field.Value
			}
		}
		if s.Code == SCTeamProvisionalCanKey {
			e.CanKey = true
		}
		return e
	case SCEphemeralPairwiseMACsMissingUIDs:
		uids := []keybase1.UID{}
		for _, field := range s.Fields {
			uids = append(uids, keybase1.UID(field.Value))
		}
		return NewEphemeralPairwiseMACsMissingUIDsError(uids)
	case SCMerkleClientError:
		e := MerkleClientError{m: s.Desc}
		for _, field := range s.Fields {
			if field.Key == "type" {
				i, err := strconv.Atoi(field.Value)
				if err != nil {
					g.Log.Warning("error parsing merkle error type: %s", err)
				} else {
					e.t = merkleClientErrorType(i)
				}
			}
		}
		return e
	case SCTeamFTLOutdated:
		return NewTeamFTLOutdatedError(s.Desc)
	case SCFeatureFlag:
		var feature Feature
		for _, field := range s.Fields {
			if field.Key == "feature" {
				feature = Feature(field.Value)
			}
		}
		return NewFeatureFlagError(s.Desc, feature)
	case SCNoPaperKeys:
		return NoPaperKeysError{}
	default:
		ase := AppStatusError{
			Code:   s.Code,
			Name:   s.Name,
			Desc:   s.Desc,
			Fields: make(map[string]string),
		}
		for _, f := range s.Fields {
			ase.Fields[f.Key] = f.Value
		}
		return ase
	}
}
