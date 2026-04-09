func (u *EventType) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// AppLinkTeam : (apps) Linked app for team
		AppLinkTeam json.RawMessage `json:"app_link_team,omitempty"`
		// AppLinkUser : (apps) Linked app for member
		AppLinkUser json.RawMessage `json:"app_link_user,omitempty"`
		// AppUnlinkTeam : (apps) Unlinked app for team
		AppUnlinkTeam json.RawMessage `json:"app_unlink_team,omitempty"`
		// AppUnlinkUser : (apps) Unlinked app for member
		AppUnlinkUser json.RawMessage `json:"app_unlink_user,omitempty"`
		// FileAddComment : (comments) Added file comment
		FileAddComment json.RawMessage `json:"file_add_comment,omitempty"`
		// FileChangeCommentSubscription : (comments) Subscribed to or
		// unsubscribed from comment notifications for file
		FileChangeCommentSubscription json.RawMessage `json:"file_change_comment_subscription,omitempty"`
		// FileDeleteComment : (comments) Deleted file comment
		FileDeleteComment json.RawMessage `json:"file_delete_comment,omitempty"`
		// FileEditComment : (comments) Edited file comment
		FileEditComment json.RawMessage `json:"file_edit_comment,omitempty"`
		// FileLikeComment : (comments) Liked file comment (deprecated, no
		// longer logged)
		FileLikeComment json.RawMessage `json:"file_like_comment,omitempty"`
		// FileResolveComment : (comments) Resolved file comment
		FileResolveComment json.RawMessage `json:"file_resolve_comment,omitempty"`
		// FileUnlikeComment : (comments) Unliked file comment (deprecated, no
		// longer logged)
		FileUnlikeComment json.RawMessage `json:"file_unlike_comment,omitempty"`
		// FileUnresolveComment : (comments) Unresolved file comment
		FileUnresolveComment json.RawMessage `json:"file_unresolve_comment,omitempty"`
		// DeviceChangeIpDesktop : (devices) Changed IP address associated with
		// active desktop session
		DeviceChangeIpDesktop json.RawMessage `json:"device_change_ip_desktop,omitempty"`
		// DeviceChangeIpMobile : (devices) Changed IP address associated with
		// active mobile session
		DeviceChangeIpMobile json.RawMessage `json:"device_change_ip_mobile,omitempty"`
		// DeviceChangeIpWeb : (devices) Changed IP address associated with
		// active web session
		DeviceChangeIpWeb json.RawMessage `json:"device_change_ip_web,omitempty"`
		// DeviceDeleteOnUnlinkFail : (devices) Failed to delete all files from
		// unlinked device
		DeviceDeleteOnUnlinkFail json.RawMessage `json:"device_delete_on_unlink_fail,omitempty"`
		// DeviceDeleteOnUnlinkSuccess : (devices) Deleted all files from
		// unlinked device
		DeviceDeleteOnUnlinkSuccess json.RawMessage `json:"device_delete_on_unlink_success,omitempty"`
		// DeviceLinkFail : (devices) Failed to link device
		DeviceLinkFail json.RawMessage `json:"device_link_fail,omitempty"`
		// DeviceLinkSuccess : (devices) Linked device
		DeviceLinkSuccess json.RawMessage `json:"device_link_success,omitempty"`
		// DeviceManagementDisabled : (devices) Disabled device management
		// (deprecated, no longer logged)
		DeviceManagementDisabled json.RawMessage `json:"device_management_disabled,omitempty"`
		// DeviceManagementEnabled : (devices) Enabled device management
		// (deprecated, no longer logged)
		DeviceManagementEnabled json.RawMessage `json:"device_management_enabled,omitempty"`
		// DeviceUnlink : (devices) Disconnected device
		DeviceUnlink json.RawMessage `json:"device_unlink,omitempty"`
		// EmmRefreshAuthToken : (devices) Refreshed auth token used for setting
		// up enterprise mobility management
		EmmRefreshAuthToken json.RawMessage `json:"emm_refresh_auth_token,omitempty"`
		// AccountCaptureChangeAvailability : (domains) Granted/revoked option
		// to enable account capture on team domains
		AccountCaptureChangeAvailability json.RawMessage `json:"account_capture_change_availability,omitempty"`
		// AccountCaptureMigrateAccount : (domains) Account-captured user
		// migrated account to team
		AccountCaptureMigrateAccount json.RawMessage `json:"account_capture_migrate_account,omitempty"`
		// AccountCaptureNotificationEmailsSent : (domains) Sent proactive
		// account capture email to all unmanaged members
		AccountCaptureNotificationEmailsSent json.RawMessage `json:"account_capture_notification_emails_sent,omitempty"`
		// AccountCaptureRelinquishAccount : (domains) Account-captured user
		// changed account email to personal email
		AccountCaptureRelinquishAccount json.RawMessage `json:"account_capture_relinquish_account,omitempty"`
		// DisabledDomainInvites : (domains) Disabled domain invites
		// (deprecated, no longer logged)
		DisabledDomainInvites json.RawMessage `json:"disabled_domain_invites,omitempty"`
		// DomainInvitesApproveRequestToJoinTeam : (domains) Approved user's
		// request to join team
		DomainInvitesApproveRequestToJoinTeam json.RawMessage `json:"domain_invites_approve_request_to_join_team,omitempty"`
		// DomainInvitesDeclineRequestToJoinTeam : (domains) Declined user's
		// request to join team
		DomainInvitesDeclineRequestToJoinTeam json.RawMessage `json:"domain_invites_decline_request_to_join_team,omitempty"`
		// DomainInvitesEmailExistingUsers : (domains) Sent domain invites to
		// existing domain accounts (deprecated, no longer logged)
		DomainInvitesEmailExistingUsers json.RawMessage `json:"domain_invites_email_existing_users,omitempty"`
		// DomainInvitesRequestToJoinTeam : (domains) Requested to join team
		DomainInvitesRequestToJoinTeam json.RawMessage `json:"domain_invites_request_to_join_team,omitempty"`
		// DomainInvitesSetInviteNewUserPrefToNo : (domains) Disabled
		// "Automatically invite new users" (deprecated, no longer logged)
		DomainInvitesSetInviteNewUserPrefToNo json.RawMessage `json:"domain_invites_set_invite_new_user_pref_to_no,omitempty"`
		// DomainInvitesSetInviteNewUserPrefToYes : (domains) Enabled
		// "Automatically invite new users" (deprecated, no longer logged)
		DomainInvitesSetInviteNewUserPrefToYes json.RawMessage `json:"domain_invites_set_invite_new_user_pref_to_yes,omitempty"`
		// DomainVerificationAddDomainFail : (domains) Failed to verify team
		// domain
		DomainVerificationAddDomainFail json.RawMessage `json:"domain_verification_add_domain_fail,omitempty"`
		// DomainVerificationAddDomainSuccess : (domains) Verified team domain
		DomainVerificationAddDomainSuccess json.RawMessage `json:"domain_verification_add_domain_success,omitempty"`
		// DomainVerificationRemoveDomain : (domains) Removed domain from list
		// of verified team domains
		DomainVerificationRemoveDomain json.RawMessage `json:"domain_verification_remove_domain,omitempty"`
		// EnabledDomainInvites : (domains) Enabled domain invites (deprecated,
		// no longer logged)
		EnabledDomainInvites json.RawMessage `json:"enabled_domain_invites,omitempty"`
		// CreateFolder : (file_operations) Created folders (deprecated, no
		// longer logged)
		CreateFolder json.RawMessage `json:"create_folder,omitempty"`
		// FileAdd : (file_operations) Added files and/or folders
		FileAdd json.RawMessage `json:"file_add,omitempty"`
		// FileCopy : (file_operations) Copied files and/or folders
		FileCopy json.RawMessage `json:"file_copy,omitempty"`
		// FileDelete : (file_operations) Deleted files and/or folders
		FileDelete json.RawMessage `json:"file_delete,omitempty"`
		// FileDownload : (file_operations) Downloaded files and/or folders
		FileDownload json.RawMessage `json:"file_download,omitempty"`
		// FileEdit : (file_operations) Edited files
		FileEdit json.RawMessage `json:"file_edit,omitempty"`
		// FileGetCopyReference : (file_operations) Created copy reference to
		// file/folder
		FileGetCopyReference json.RawMessage `json:"file_get_copy_reference,omitempty"`
		// FileMove : (file_operations) Moved files and/or folders
		FileMove json.RawMessage `json:"file_move,omitempty"`
		// FilePermanentlyDelete : (file_operations) Permanently deleted files
		// and/or folders
		FilePermanentlyDelete json.RawMessage `json:"file_permanently_delete,omitempty"`
		// FilePreview : (file_operations) Previewed files and/or folders
		FilePreview json.RawMessage `json:"file_preview,omitempty"`
		// FileRename : (file_operations) Renamed files and/or folders
		FileRename json.RawMessage `json:"file_rename,omitempty"`
		// FileRestore : (file_operations) Restored deleted files and/or folders
		FileRestore json.RawMessage `json:"file_restore,omitempty"`
		// FileRevert : (file_operations) Reverted files to previous version
		FileRevert json.RawMessage `json:"file_revert,omitempty"`
		// FileRollbackChanges : (file_operations) Rolled back file actions
		FileRollbackChanges json.RawMessage `json:"file_rollback_changes,omitempty"`
		// FileSaveCopyReference : (file_operations) Saved file/folder using
		// copy reference
		FileSaveCopyReference json.RawMessage `json:"file_save_copy_reference,omitempty"`
		// FileRequestChange : (file_requests) Changed file request
		FileRequestChange json.RawMessage `json:"file_request_change,omitempty"`
		// FileRequestClose : (file_requests) Closed file request
		FileRequestClose json.RawMessage `json:"file_request_close,omitempty"`
		// FileRequestCreate : (file_requests) Created file request
		FileRequestCreate json.RawMessage `json:"file_request_create,omitempty"`
		// FileRequestReceiveFile : (file_requests) Received files for file
		// request
		FileRequestReceiveFile json.RawMessage `json:"file_request_receive_file,omitempty"`
		// GroupAddExternalId : (groups) Added external ID for group
		GroupAddExternalId json.RawMessage `json:"group_add_external_id,omitempty"`
		// GroupAddMember : (groups) Added team members to group
		GroupAddMember json.RawMessage `json:"group_add_member,omitempty"`
		// GroupChangeExternalId : (groups) Changed external ID for group
		GroupChangeExternalId json.RawMessage `json:"group_change_external_id,omitempty"`
		// GroupChangeManagementType : (groups) Changed group management type
		GroupChangeManagementType json.RawMessage `json:"group_change_management_type,omitempty"`
		// GroupChangeMemberRole : (groups) Changed manager permissions of group
		// member
		GroupChangeMemberRole json.RawMessage `json:"group_change_member_role,omitempty"`
		// GroupCreate : (groups) Created group
		GroupCreate json.RawMessage `json:"group_create,omitempty"`
		// GroupDelete : (groups) Deleted group
		GroupDelete json.RawMessage `json:"group_delete,omitempty"`
		// GroupDescriptionUpdated : (groups) Updated group (deprecated, no
		// longer logged)
		GroupDescriptionUpdated json.RawMessage `json:"group_description_updated,omitempty"`
		// GroupJoinPolicyUpdated : (groups) Updated group join policy
		// (deprecated, no longer logged)
		GroupJoinPolicyUpdated json.RawMessage `json:"group_join_policy_updated,omitempty"`
		// GroupMoved : (groups) Moved group (deprecated, no longer logged)
		GroupMoved json.RawMessage `json:"group_moved,omitempty"`
		// GroupRemoveExternalId : (groups) Removed external ID for group
		GroupRemoveExternalId json.RawMessage `json:"group_remove_external_id,omitempty"`
		// GroupRemoveMember : (groups) Removed team members from group
		GroupRemoveMember json.RawMessage `json:"group_remove_member,omitempty"`
		// GroupRename : (groups) Renamed group
		GroupRename json.RawMessage `json:"group_rename,omitempty"`
		// EmmError : (logins) Failed to sign in via EMM (deprecated, replaced
		// by 'Failed to sign in')
		EmmError json.RawMessage `json:"emm_error,omitempty"`
		// LoginFail : (logins) Failed to sign in
		LoginFail json.RawMessage `json:"login_fail,omitempty"`
		// LoginSuccess : (logins) Signed in
		LoginSuccess json.RawMessage `json:"login_success,omitempty"`
		// Logout : (logins) Signed out
		Logout json.RawMessage `json:"logout,omitempty"`
		// ResellerSupportSessionEnd : (logins) Ended reseller support session
		ResellerSupportSessionEnd json.RawMessage `json:"reseller_support_session_end,omitempty"`
		// ResellerSupportSessionStart : (logins) Started reseller support
		// session
		ResellerSupportSessionStart json.RawMessage `json:"reseller_support_session_start,omitempty"`
		// SignInAsSessionEnd : (logins) Ended admin sign-in-as session
		SignInAsSessionEnd json.RawMessage `json:"sign_in_as_session_end,omitempty"`
		// SignInAsSessionStart : (logins) Started admin sign-in-as session
		SignInAsSessionStart json.RawMessage `json:"sign_in_as_session_start,omitempty"`
		// SsoError : (logins) Failed to sign in via SSO (deprecated, replaced
		// by 'Failed to sign in')
		SsoError json.RawMessage `json:"sso_error,omitempty"`
		// MemberAddName : (members) Added team member name
		MemberAddName json.RawMessage `json:"member_add_name,omitempty"`
		// MemberChangeAdminRole : (members) Changed team member admin role
		MemberChangeAdminRole json.RawMessage `json:"member_change_admin_role,omitempty"`
		// MemberChangeEmail : (members) Changed team member email
		MemberChangeEmail json.RawMessage `json:"member_change_email,omitempty"`
		// MemberChangeMembershipType : (members) Changed membership type
		// (limited/full) of member (deprecated, no longer logged)
		MemberChangeMembershipType json.RawMessage `json:"member_change_membership_type,omitempty"`
		// MemberChangeName : (members) Changed team member name
		MemberChangeName json.RawMessage `json:"member_change_name,omitempty"`
		// MemberChangeStatus : (members) Changed member status (invited,
		// joined, suspended, etc.)
		MemberChangeStatus json.RawMessage `json:"member_change_status,omitempty"`
		// MemberDeleteManualContacts : (members) Cleared manually added
		// contacts
		MemberDeleteManualContacts json.RawMessage `json:"member_delete_manual_contacts,omitempty"`
		// MemberPermanentlyDeleteAccountContents : (members) Permanently
		// deleted contents of deleted team member account
		MemberPermanentlyDeleteAccountContents json.RawMessage `json:"member_permanently_delete_account_contents,omitempty"`
		// MemberSpaceLimitsAddCustomQuota : (members) Set custom member space
		// limit
		MemberSpaceLimitsAddCustomQuota json.RawMessage `json:"member_space_limits_add_custom_quota,omitempty"`
		// MemberSpaceLimitsChangeCustomQuota : (members) Changed custom member
		// space limit
		MemberSpaceLimitsChangeCustomQuota json.RawMessage `json:"member_space_limits_change_custom_quota,omitempty"`
		// MemberSpaceLimitsChangeStatus : (members) Changed space limit status
		MemberSpaceLimitsChangeStatus json.RawMessage `json:"member_space_limits_change_status,omitempty"`
		// MemberSpaceLimitsRemoveCustomQuota : (members) Removed custom member
		// space limit
		MemberSpaceLimitsRemoveCustomQuota json.RawMessage `json:"member_space_limits_remove_custom_quota,omitempty"`
		// MemberSuggest : (members) Suggested person to add to team
		MemberSuggest json.RawMessage `json:"member_suggest,omitempty"`
		// MemberTransferAccountContents : (members) Transferred contents of
		// deleted member account to another member
		MemberTransferAccountContents json.RawMessage `json:"member_transfer_account_contents,omitempty"`
		// SecondaryMailsPolicyChanged : (members) Secondary mails policy
		// changed
		SecondaryMailsPolicyChanged json.RawMessage `json:"secondary_mails_policy_changed,omitempty"`
		// PaperContentAddMember : (paper) Added team member to Paper doc/folder
		PaperContentAddMember json.RawMessage `json:"paper_content_add_member,omitempty"`
		// PaperContentAddToFolder : (paper) Added Paper doc/folder to folder
		PaperContentAddToFolder json.RawMessage `json:"paper_content_add_to_folder,omitempty"`
		// PaperContentArchive : (paper) Archived Paper doc/folder
		PaperContentArchive json.RawMessage `json:"paper_content_archive,omitempty"`
		// PaperContentCreate : (paper) Created Paper doc/folder
		PaperContentCreate json.RawMessage `json:"paper_content_create,omitempty"`
		// PaperContentPermanentlyDelete : (paper) Permanently deleted Paper
		// doc/folder
		PaperContentPermanentlyDelete json.RawMessage `json:"paper_content_permanently_delete,omitempty"`
		// PaperContentRemoveFromFolder : (paper) Removed Paper doc/folder from
		// folder
		PaperContentRemoveFromFolder json.RawMessage `json:"paper_content_remove_from_folder,omitempty"`
		// PaperContentRemoveMember : (paper) Removed team member from Paper
		// doc/folder
		PaperContentRemoveMember json.RawMessage `json:"paper_content_remove_member,omitempty"`
		// PaperContentRename : (paper) Renamed Paper doc/folder
		PaperContentRename json.RawMessage `json:"paper_content_rename,omitempty"`
		// PaperContentRestore : (paper) Restored archived Paper doc/folder
		PaperContentRestore json.RawMessage `json:"paper_content_restore,omitempty"`
		// PaperDocAddComment : (paper) Added Paper doc comment
		PaperDocAddComment json.RawMessage `json:"paper_doc_add_comment,omitempty"`
		// PaperDocChangeMemberRole : (paper) Changed team member permissions
		// for Paper doc
		PaperDocChangeMemberRole json.RawMessage `json:"paper_doc_change_member_role,omitempty"`
		// PaperDocChangeSharingPolicy : (paper) Changed sharing setting for
		// Paper doc
		PaperDocChangeSharingPolicy json.RawMessage `json:"paper_doc_change_sharing_policy,omitempty"`
		// PaperDocChangeSubscription : (paper) Followed/unfollowed Paper doc
		PaperDocChangeSubscription json.RawMessage `json:"paper_doc_change_subscription,omitempty"`
		// PaperDocDeleted : (paper) Archived Paper doc (deprecated, no longer
		// logged)
		PaperDocDeleted json.RawMessage `json:"paper_doc_deleted,omitempty"`
		// PaperDocDeleteComment : (paper) Deleted Paper doc comment
		PaperDocDeleteComment json.RawMessage `json:"paper_doc_delete_comment,omitempty"`
		// PaperDocDownload : (paper) Downloaded Paper doc in specific format
		PaperDocDownload json.RawMessage `json:"paper_doc_download,omitempty"`
		// PaperDocEdit : (paper) Edited Paper doc
		PaperDocEdit json.RawMessage `json:"paper_doc_edit,omitempty"`
		// PaperDocEditComment : (paper) Edited Paper doc comment
		PaperDocEditComment json.RawMessage `json:"paper_doc_edit_comment,omitempty"`
		// PaperDocFollowed : (paper) Followed Paper doc (deprecated, replaced
		// by 'Followed/unfollowed Paper doc')
		PaperDocFollowed json.RawMessage `json:"paper_doc_followed,omitempty"`
		// PaperDocMention : (paper) Mentioned team member in Paper doc
		PaperDocMention json.RawMessage `json:"paper_doc_mention,omitempty"`
		// PaperDocOwnershipChanged : (paper) Transferred ownership of Paper doc
		PaperDocOwnershipChanged json.RawMessage `json:"paper_doc_ownership_changed,omitempty"`
		// PaperDocRequestAccess : (paper) Requested access to Paper doc
		PaperDocRequestAccess json.RawMessage `json:"paper_doc_request_access,omitempty"`
		// PaperDocResolveComment : (paper) Resolved Paper doc comment
		PaperDocResolveComment json.RawMessage `json:"paper_doc_resolve_comment,omitempty"`
		// PaperDocRevert : (paper) Restored Paper doc to previous version
		PaperDocRevert json.RawMessage `json:"paper_doc_revert,omitempty"`
		// PaperDocSlackShare : (paper) Shared Paper doc via Slack
		PaperDocSlackShare json.RawMessage `json:"paper_doc_slack_share,omitempty"`
		// PaperDocTeamInvite : (paper) Shared Paper doc with team member
		// (deprecated, no longer logged)
		PaperDocTeamInvite json.RawMessage `json:"paper_doc_team_invite,omitempty"`
		// PaperDocTrashed : (paper) Deleted Paper doc
		PaperDocTrashed json.RawMessage `json:"paper_doc_trashed,omitempty"`
		// PaperDocUnresolveComment : (paper) Unresolved Paper doc comment
		PaperDocUnresolveComment json.RawMessage `json:"paper_doc_unresolve_comment,omitempty"`
		// PaperDocUntrashed : (paper) Restored Paper doc
		PaperDocUntrashed json.RawMessage `json:"paper_doc_untrashed,omitempty"`
		// PaperDocView : (paper) Viewed Paper doc
		PaperDocView json.RawMessage `json:"paper_doc_view,omitempty"`
		// PaperExternalViewAllow : (paper) Changed Paper external sharing
		// setting to anyone (deprecated, no longer logged)
		PaperExternalViewAllow json.RawMessage `json:"paper_external_view_allow,omitempty"`
		// PaperExternalViewDefaultTeam : (paper) Changed Paper external sharing
		// setting to default team (deprecated, no longer logged)
		PaperExternalViewDefaultTeam json.RawMessage `json:"paper_external_view_default_team,omitempty"`
		// PaperExternalViewForbid : (paper) Changed Paper external sharing
		// setting to team-only (deprecated, no longer logged)
		PaperExternalViewForbid json.RawMessage `json:"paper_external_view_forbid,omitempty"`
		// PaperFolderChangeSubscription : (paper) Followed/unfollowed Paper
		// folder
		PaperFolderChangeSubscription json.RawMessage `json:"paper_folder_change_subscription,omitempty"`
		// PaperFolderDeleted : (paper) Archived Paper folder (deprecated, no
		// longer logged)
		PaperFolderDeleted json.RawMessage `json:"paper_folder_deleted,omitempty"`
		// PaperFolderFollowed : (paper) Followed Paper folder (deprecated,
		// replaced by 'Followed/unfollowed Paper folder')
		PaperFolderFollowed json.RawMessage `json:"paper_folder_followed,omitempty"`
		// PaperFolderTeamInvite : (paper) Shared Paper folder with member
		// (deprecated, no longer logged)
		PaperFolderTeamInvite json.RawMessage `json:"paper_folder_team_invite,omitempty"`
		// PasswordChange : (passwords) Changed password
		PasswordChange json.RawMessage `json:"password_change,omitempty"`
		// PasswordReset : (passwords) Reset password
		PasswordReset json.RawMessage `json:"password_reset,omitempty"`
		// PasswordResetAll : (passwords) Reset all team member passwords
		PasswordResetAll json.RawMessage `json:"password_reset_all,omitempty"`
		// EmmCreateExceptionsReport : (reports) Created EMM-excluded users
		// report
		EmmCreateExceptionsReport json.RawMessage `json:"emm_create_exceptions_report,omitempty"`
		// EmmCreateUsageReport : (reports) Created EMM mobile app usage report
		EmmCreateUsageReport json.RawMessage `json:"emm_create_usage_report,omitempty"`
		// ExportMembersReport : (reports) Created member data report
		ExportMembersReport json.RawMessage `json:"export_members_report,omitempty"`
		// PaperAdminExportStart : (reports) Exported all team Paper docs
		PaperAdminExportStart json.RawMessage `json:"paper_admin_export_start,omitempty"`
		// SmartSyncCreateAdminPrivilegeReport : (reports) Created Smart Sync
		// non-admin devices report
		SmartSyncCreateAdminPrivilegeReport json.RawMessage `json:"smart_sync_create_admin_privilege_report,omitempty"`
		// TeamActivityCreateReport : (reports) Created team activity report
		TeamActivityCreateReport json.RawMessage `json:"team_activity_create_report,omitempty"`
		// CollectionShare : (sharing) Shared album
		CollectionShare json.RawMessage `json:"collection_share,omitempty"`
		// NoteAclInviteOnly : (sharing) Changed Paper doc to invite-only
		// (deprecated, no longer logged)
		NoteAclInviteOnly json.RawMessage `json:"note_acl_invite_only,omitempty"`
		// NoteAclLink : (sharing) Changed Paper doc to link-accessible
		// (deprecated, no longer logged)
		NoteAclLink json.RawMessage `json:"note_acl_link,omitempty"`
		// NoteAclTeamLink : (sharing) Changed Paper doc to link-accessible for
		// team (deprecated, no longer logged)
		NoteAclTeamLink json.RawMessage `json:"note_acl_team_link,omitempty"`
		// NoteShared : (sharing) Shared Paper doc (deprecated, no longer
		// logged)
		NoteShared json.RawMessage `json:"note_shared,omitempty"`
		// NoteShareReceive : (sharing) Shared received Paper doc (deprecated,
		// no longer logged)
		NoteShareReceive json.RawMessage `json:"note_share_receive,omitempty"`
		// OpenNoteShared : (sharing) Opened shared Paper doc (deprecated, no
		// longer logged)
		OpenNoteShared json.RawMessage `json:"open_note_shared,omitempty"`
		// SfAddGroup : (sharing) Added team to shared folder (deprecated, no
		// longer logged)
		SfAddGroup json.RawMessage `json:"sf_add_group,omitempty"`
		// SfAllowNonMembersToViewSharedLinks : (sharing) Allowed
		// non-collaborators to view links to files in shared folder
		// (deprecated, no longer logged)
		SfAllowNonMembersToViewSharedLinks json.RawMessage `json:"sf_allow_non_members_to_view_shared_links,omitempty"`
		// SfExternalInviteWarn : (sharing) Set team members to see warning
		// before sharing folders outside team (deprecated, no longer logged)
		SfExternalInviteWarn json.RawMessage `json:"sf_external_invite_warn,omitempty"`
		// SfFbInvite : (sharing) Invited Facebook users to shared folder
		// (deprecated, no longer logged)
		SfFbInvite json.RawMessage `json:"sf_fb_invite,omitempty"`
		// SfFbInviteChangeRole : (sharing) Changed Facebook user's role in
		// shared folder (deprecated, no longer logged)
		SfFbInviteChangeRole json.RawMessage `json:"sf_fb_invite_change_role,omitempty"`
		// SfFbUninvite : (sharing) Uninvited Facebook user from shared folder
		// (deprecated, no longer logged)
		SfFbUninvite json.RawMessage `json:"sf_fb_uninvite,omitempty"`
		// SfInviteGroup : (sharing) Invited group to shared folder (deprecated,
		// no longer logged)
		SfInviteGroup json.RawMessage `json:"sf_invite_group,omitempty"`
		// SfTeamGrantAccess : (sharing) Granted access to shared folder
		// (deprecated, no longer logged)
		SfTeamGrantAccess json.RawMessage `json:"sf_team_grant_access,omitempty"`
		// SfTeamInvite : (sharing) Invited team members to shared folder
		// (deprecated, replaced by 'Invited user to Dropbox and added them to
		// shared file/folder')
		SfTeamInvite json.RawMessage `json:"sf_team_invite,omitempty"`
		// SfTeamInviteChangeRole : (sharing) Changed team member's role in
		// shared folder (deprecated, no longer logged)
		SfTeamInviteChangeRole json.RawMessage `json:"sf_team_invite_change_role,omitempty"`
		// SfTeamJoin : (sharing) Joined team member's shared folder
		// (deprecated, no longer logged)
		SfTeamJoin json.RawMessage `json:"sf_team_join,omitempty"`
		// SfTeamJoinFromOobLink : (sharing) Joined team member's shared folder
		// from link (deprecated, no longer logged)
		SfTeamJoinFromOobLink json.RawMessage `json:"sf_team_join_from_oob_link,omitempty"`
		// SfTeamUninvite : (sharing) Unshared folder with team member
		// (deprecated, replaced by 'Removed invitee from shared file/folder
		// before invite was accepted')
		SfTeamUninvite json.RawMessage `json:"sf_team_uninvite,omitempty"`
		// SharedContentAddInvitees : (sharing) Invited user to Dropbox and
		// added them to shared file/folder
		SharedContentAddInvitees json.RawMessage `json:"shared_content_add_invitees,omitempty"`
		// SharedContentAddLinkExpiry : (sharing) Added expiration date to link
		// for shared file/folder
		SharedContentAddLinkExpiry json.RawMessage `json:"shared_content_add_link_expiry,omitempty"`
		// SharedContentAddLinkPassword : (sharing) Added password to link for
		// shared file/folder
		SharedContentAddLinkPassword json.RawMessage `json:"shared_content_add_link_password,omitempty"`
		// SharedContentAddMember : (sharing) Added users and/or groups to
		// shared file/folder
		SharedContentAddMember json.RawMessage `json:"shared_content_add_member,omitempty"`
		// SharedContentChangeDownloadsPolicy : (sharing) Changed whether
		// members can download shared file/folder
		SharedContentChangeDownloadsPolicy json.RawMessage `json:"shared_content_change_downloads_policy,omitempty"`
		// SharedContentChangeInviteeRole : (sharing) Changed access type of
		// invitee to shared file/folder before invite was accepted
		SharedContentChangeInviteeRole json.RawMessage `json:"shared_content_change_invitee_role,omitempty"`
		// SharedContentChangeLinkAudience : (sharing) Changed link audience of
		// shared file/folder
		SharedContentChangeLinkAudience json.RawMessage `json:"shared_content_change_link_audience,omitempty"`
		// SharedContentChangeLinkExpiry : (sharing) Changed link expiration of
		// shared file/folder
		SharedContentChangeLinkExpiry json.RawMessage `json:"shared_content_change_link_expiry,omitempty"`
		// SharedContentChangeLinkPassword : (sharing) Changed link password of
		// shared file/folder
		SharedContentChangeLinkPassword json.RawMessage `json:"shared_content_change_link_password,omitempty"`
		// SharedContentChangeMemberRole : (sharing) Changed access type of
		// shared file/folder member
		SharedContentChangeMemberRole json.RawMessage `json:"shared_content_change_member_role,omitempty"`
		// SharedContentChangeViewerInfoPolicy : (sharing) Changed whether
		// members can see who viewed shared file/folder
		SharedContentChangeViewerInfoPolicy json.RawMessage `json:"shared_content_change_viewer_info_policy,omitempty"`
		// SharedContentClaimInvitation : (sharing) Acquired membership of
		// shared file/folder by accepting invite
		SharedContentClaimInvitation json.RawMessage `json:"shared_content_claim_invitation,omitempty"`
		// SharedContentCopy : (sharing) Copied shared file/folder to own
		// Dropbox
		SharedContentCopy json.RawMessage `json:"shared_content_copy,omitempty"`
		// SharedContentDownload : (sharing) Downloaded shared file/folder
		SharedContentDownload json.RawMessage `json:"shared_content_download,omitempty"`
		// SharedContentRelinquishMembership : (sharing) Left shared file/folder
		SharedContentRelinquishMembership json.RawMessage `json:"shared_content_relinquish_membership,omitempty"`
		// SharedContentRemoveInvitees : (sharing) Removed invitee from shared
		// file/folder before invite was accepted
		SharedContentRemoveInvitees json.RawMessage `json:"shared_content_remove_invitees,omitempty"`
		// SharedContentRemoveLinkExpiry : (sharing) Removed link expiration
		// date of shared file/folder
		SharedContentRemoveLinkExpiry json.RawMessage `json:"shared_content_remove_link_expiry,omitempty"`
		// SharedContentRemoveLinkPassword : (sharing) Removed link password of
		// shared file/folder
		SharedContentRemoveLinkPassword json.RawMessage `json:"shared_content_remove_link_password,omitempty"`
		// SharedContentRemoveMember : (sharing) Removed user/group from shared
		// file/folder
		SharedContentRemoveMember json.RawMessage `json:"shared_content_remove_member,omitempty"`
		// SharedContentRequestAccess : (sharing) Requested access to shared
		// file/folder
		SharedContentRequestAccess json.RawMessage `json:"shared_content_request_access,omitempty"`
		// SharedContentUnshare : (sharing) Unshared file/folder by clearing
		// membership and turning off link
		SharedContentUnshare json.RawMessage `json:"shared_content_unshare,omitempty"`
		// SharedContentView : (sharing) Previewed shared file/folder
		SharedContentView json.RawMessage `json:"shared_content_view,omitempty"`
		// SharedFolderChangeLinkPolicy : (sharing) Changed who can access
		// shared folder via link
		SharedFolderChangeLinkPolicy json.RawMessage `json:"shared_folder_change_link_policy,omitempty"`
		// SharedFolderChangeMembersInheritancePolicy : (sharing) Changed
		// whether shared folder inherits members from parent folder
		SharedFolderChangeMembersInheritancePolicy json.RawMessage `json:"shared_folder_change_members_inheritance_policy,omitempty"`
		// SharedFolderChangeMembersManagementPolicy : (sharing) Changed who can
		// add/remove members of shared folder
		SharedFolderChangeMembersManagementPolicy json.RawMessage `json:"shared_folder_change_members_management_policy,omitempty"`
		// SharedFolderChangeMembersPolicy : (sharing) Changed who can become
		// member of shared folder
		SharedFolderChangeMembersPolicy json.RawMessage `json:"shared_folder_change_members_policy,omitempty"`
		// SharedFolderCreate : (sharing) Created shared folder
		SharedFolderCreate json.RawMessage `json:"shared_folder_create,omitempty"`
		// SharedFolderDeclineInvitation : (sharing) Declined team member's
		// invite to shared folder
		SharedFolderDeclineInvitation json.RawMessage `json:"shared_folder_decline_invitation,omitempty"`
		// SharedFolderMount : (sharing) Added shared folder to own Dropbox
		SharedFolderMount json.RawMessage `json:"shared_folder_mount,omitempty"`
		// SharedFolderNest : (sharing) Changed parent of shared folder
		SharedFolderNest json.RawMessage `json:"shared_folder_nest,omitempty"`
		// SharedFolderTransferOwnership : (sharing) Transferred ownership of
		// shared folder to another member
		SharedFolderTransferOwnership json.RawMessage `json:"shared_folder_transfer_ownership,omitempty"`
		// SharedFolderUnmount : (sharing) Deleted shared folder from Dropbox
		SharedFolderUnmount json.RawMessage `json:"shared_folder_unmount,omitempty"`
		// SharedLinkAddExpiry : (sharing) Added shared link expiration date
		SharedLinkAddExpiry json.RawMessage `json:"shared_link_add_expiry,omitempty"`
		// SharedLinkChangeExpiry : (sharing) Changed shared link expiration
		// date
		SharedLinkChangeExpiry json.RawMessage `json:"shared_link_change_expiry,omitempty"`
		// SharedLinkChangeVisibility : (sharing) Changed visibility of shared
		// link
		SharedLinkChangeVisibility json.RawMessage `json:"shared_link_change_visibility,omitempty"`
		// SharedLinkCopy : (sharing) Added file/folder to Dropbox from shared
		// link
		SharedLinkCopy json.RawMessage `json:"shared_link_copy,omitempty"`
		// SharedLinkCreate : (sharing) Created shared link
		SharedLinkCreate json.RawMessage `json:"shared_link_create,omitempty"`
		// SharedLinkDisable : (sharing) Removed shared link
		SharedLinkDisable json.RawMessage `json:"shared_link_disable,omitempty"`
		// SharedLinkDownload : (sharing) Downloaded file/folder from shared
		// link
		SharedLinkDownload json.RawMessage `json:"shared_link_download,omitempty"`
		// SharedLinkRemoveExpiry : (sharing) Removed shared link expiration
		// date
		SharedLinkRemoveExpiry json.RawMessage `json:"shared_link_remove_expiry,omitempty"`
		// SharedLinkShare : (sharing) Added members as audience of shared link
		SharedLinkShare json.RawMessage `json:"shared_link_share,omitempty"`
		// SharedLinkView : (sharing) Opened shared link
		SharedLinkView json.RawMessage `json:"shared_link_view,omitempty"`
		// SharedNoteOpened : (sharing) Opened shared Paper doc (deprecated, no
		// longer logged)
		SharedNoteOpened json.RawMessage `json:"shared_note_opened,omitempty"`
		// ShmodelGroupShare : (sharing) Shared link with group (deprecated, no
		// longer logged)
		ShmodelGroupShare json.RawMessage `json:"shmodel_group_share,omitempty"`
		// ShowcaseAccessGranted : (showcase) Granted access to showcase
		ShowcaseAccessGranted json.RawMessage `json:"showcase_access_granted,omitempty"`
		// ShowcaseAddMember : (showcase) Added member to showcase
		ShowcaseAddMember json.RawMessage `json:"showcase_add_member,omitempty"`
		// ShowcaseArchived : (showcase) Archived showcase
		ShowcaseArchived json.RawMessage `json:"showcase_archived,omitempty"`
		// ShowcaseCreated : (showcase) Created showcase
		ShowcaseCreated json.RawMessage `json:"showcase_created,omitempty"`
		// ShowcaseDeleteComment : (showcase) Deleted showcase comment
		ShowcaseDeleteComment json.RawMessage `json:"showcase_delete_comment,omitempty"`
		// ShowcaseEdited : (showcase) Edited showcase
		ShowcaseEdited json.RawMessage `json:"showcase_edited,omitempty"`
		// ShowcaseEditComment : (showcase) Edited showcase comment
		ShowcaseEditComment json.RawMessage `json:"showcase_edit_comment,omitempty"`
		// ShowcaseFileAdded : (showcase) Added file to showcase
		ShowcaseFileAdded json.RawMessage `json:"showcase_file_added,omitempty"`
		// ShowcaseFileDownload : (showcase) Downloaded file from showcase
		ShowcaseFileDownload json.RawMessage `json:"showcase_file_download,omitempty"`
		// ShowcaseFileRemoved : (showcase) Removed file from showcase
		ShowcaseFileRemoved json.RawMessage `json:"showcase_file_removed,omitempty"`
		// ShowcaseFileView : (showcase) Viewed file in showcase
		ShowcaseFileView json.RawMessage `json:"showcase_file_view,omitempty"`
		// ShowcasePermanentlyDeleted : (showcase) Permanently deleted showcase
		ShowcasePermanentlyDeleted json.RawMessage `json:"showcase_permanently_deleted,omitempty"`
		// ShowcasePostComment : (showcase) Added showcase comment
		ShowcasePostComment json.RawMessage `json:"showcase_post_comment,omitempty"`
		// ShowcaseRemoveMember : (showcase) Removed member from showcase
		ShowcaseRemoveMember json.RawMessage `json:"showcase_remove_member,omitempty"`
		// ShowcaseRenamed : (showcase) Renamed showcase
		ShowcaseRenamed json.RawMessage `json:"showcase_renamed,omitempty"`
		// ShowcaseRequestAccess : (showcase) Requested access to showcase
		ShowcaseRequestAccess json.RawMessage `json:"showcase_request_access,omitempty"`
		// ShowcaseResolveComment : (showcase) Resolved showcase comment
		ShowcaseResolveComment json.RawMessage `json:"showcase_resolve_comment,omitempty"`
		// ShowcaseRestored : (showcase) Unarchived showcase
		ShowcaseRestored json.RawMessage `json:"showcase_restored,omitempty"`
		// ShowcaseTrashed : (showcase) Deleted showcase
		ShowcaseTrashed json.RawMessage `json:"showcase_trashed,omitempty"`
		// ShowcaseTrashedDeprecated : (showcase) Deleted showcase (old version)
		// (deprecated, replaced by 'Deleted showcase')
		ShowcaseTrashedDeprecated json.RawMessage `json:"showcase_trashed_deprecated,omitempty"`
		// ShowcaseUnresolveComment : (showcase) Unresolved showcase comment
		ShowcaseUnresolveComment json.RawMessage `json:"showcase_unresolve_comment,omitempty"`
		// ShowcaseUntrashed : (showcase) Restored showcase
		ShowcaseUntrashed json.RawMessage `json:"showcase_untrashed,omitempty"`
		// ShowcaseUntrashedDeprecated : (showcase) Restored showcase (old
		// version) (deprecated, replaced by 'Restored showcase')
		ShowcaseUntrashedDeprecated json.RawMessage `json:"showcase_untrashed_deprecated,omitempty"`
		// ShowcaseView : (showcase) Viewed showcase
		ShowcaseView json.RawMessage `json:"showcase_view,omitempty"`
		// SsoAddCert : (sso) Added X.509 certificate for SSO
		SsoAddCert json.RawMessage `json:"sso_add_cert,omitempty"`
		// SsoAddLoginUrl : (sso) Added sign-in URL for SSO
		SsoAddLoginUrl json.RawMessage `json:"sso_add_login_url,omitempty"`
		// SsoAddLogoutUrl : (sso) Added sign-out URL for SSO
		SsoAddLogoutUrl json.RawMessage `json:"sso_add_logout_url,omitempty"`
		// SsoChangeCert : (sso) Changed X.509 certificate for SSO
		SsoChangeCert json.RawMessage `json:"sso_change_cert,omitempty"`
		// SsoChangeLoginUrl : (sso) Changed sign-in URL for SSO
		SsoChangeLoginUrl json.RawMessage `json:"sso_change_login_url,omitempty"`
		// SsoChangeLogoutUrl : (sso) Changed sign-out URL for SSO
		SsoChangeLogoutUrl json.RawMessage `json:"sso_change_logout_url,omitempty"`
		// SsoChangeSamlIdentityMode : (sso) Changed SAML identity mode for SSO
		SsoChangeSamlIdentityMode json.RawMessage `json:"sso_change_saml_identity_mode,omitempty"`
		// SsoRemoveCert : (sso) Removed X.509 certificate for SSO
		SsoRemoveCert json.RawMessage `json:"sso_remove_cert,omitempty"`
		// SsoRemoveLoginUrl : (sso) Removed sign-in URL for SSO
		SsoRemoveLoginUrl json.RawMessage `json:"sso_remove_login_url,omitempty"`
		// SsoRemoveLogoutUrl : (sso) Removed sign-out URL for SSO
		SsoRemoveLogoutUrl json.RawMessage `json:"sso_remove_logout_url,omitempty"`
		// TeamFolderChangeStatus : (team_folders) Changed archival status of
		// team folder
		TeamFolderChangeStatus json.RawMessage `json:"team_folder_change_status,omitempty"`
		// TeamFolderCreate : (team_folders) Created team folder in active
		// status
		TeamFolderCreate json.RawMessage `json:"team_folder_create,omitempty"`
		// TeamFolderDowngrade : (team_folders) Downgraded team folder to
		// regular shared folder
		TeamFolderDowngrade json.RawMessage `json:"team_folder_downgrade,omitempty"`
		// TeamFolderPermanentlyDelete : (team_folders) Permanently deleted
		// archived team folder
		TeamFolderPermanentlyDelete json.RawMessage `json:"team_folder_permanently_delete,omitempty"`
		// TeamFolderRename : (team_folders) Renamed active/archived team folder
		TeamFolderRename json.RawMessage `json:"team_folder_rename,omitempty"`
		// TeamSelectiveSyncSettingsChanged : (team_folders) Changed sync
		// default
		TeamSelectiveSyncSettingsChanged json.RawMessage `json:"team_selective_sync_settings_changed,omitempty"`
		// AccountCaptureChangePolicy : (team_policies) Changed account capture
		// setting on team domain
		AccountCaptureChangePolicy json.RawMessage `json:"account_capture_change_policy,omitempty"`
		// AllowDownloadDisabled : (team_policies) Disabled downloads
		// (deprecated, no longer logged)
		AllowDownloadDisabled json.RawMessage `json:"allow_download_disabled,omitempty"`
		// AllowDownloadEnabled : (team_policies) Enabled downloads (deprecated,
		// no longer logged)
		AllowDownloadEnabled json.RawMessage `json:"allow_download_enabled,omitempty"`
		// CameraUploadsPolicyChanged : (team_policies) Changed camera uploads
		// setting for team
		CameraUploadsPolicyChanged json.RawMessage `json:"camera_uploads_policy_changed,omitempty"`
		// DataPlacementRestrictionChangePolicy : (team_policies) Set
		// restrictions on data center locations where team data resides
		DataPlacementRestrictionChangePolicy json.RawMessage `json:"data_placement_restriction_change_policy,omitempty"`
		// DataPlacementRestrictionSatisfyPolicy : (team_policies) Completed
		// restrictions on data center locations where team data resides
		DataPlacementRestrictionSatisfyPolicy json.RawMessage `json:"data_placement_restriction_satisfy_policy,omitempty"`
		// DeviceApprovalsChangeDesktopPolicy : (team_policies) Set/removed
		// limit on number of computers member can link to team Dropbox account
		DeviceApprovalsChangeDesktopPolicy json.RawMessage `json:"device_approvals_change_desktop_policy,omitempty"`
		// DeviceApprovalsChangeMobilePolicy : (team_policies) Set/removed limit
		// on number of mobile devices member can link to team Dropbox account
		DeviceApprovalsChangeMobilePolicy json.RawMessage `json:"device_approvals_change_mobile_policy,omitempty"`
		// DeviceApprovalsChangeOverageAction : (team_policies) Changed device
		// approvals setting when member is over limit
		DeviceApprovalsChangeOverageAction json.RawMessage `json:"device_approvals_change_overage_action,omitempty"`
		// DeviceApprovalsChangeUnlinkAction : (team_policies) Changed device
		// approvals setting when member unlinks approved device
		DeviceApprovalsChangeUnlinkAction json.RawMessage `json:"device_approvals_change_unlink_action,omitempty"`
		// DirectoryRestrictionsAddMembers : (team_policies) Added members to
		// directory restrictions list
		DirectoryRestrictionsAddMembers json.RawMessage `json:"directory_restrictions_add_members,omitempty"`
		// DirectoryRestrictionsRemoveMembers : (team_policies) Removed members
		// from directory restrictions list
		DirectoryRestrictionsRemoveMembers json.RawMessage `json:"directory_restrictions_remove_members,omitempty"`
		// EmmAddException : (team_policies) Added members to EMM exception list
		EmmAddException json.RawMessage `json:"emm_add_exception,omitempty"`
		// EmmChangePolicy : (team_policies) Enabled/disabled enterprise
		// mobility management for members
		EmmChangePolicy json.RawMessage `json:"emm_change_policy,omitempty"`
		// EmmRemoveException : (team_policies) Removed members from EMM
		// exception list
		EmmRemoveException json.RawMessage `json:"emm_remove_exception,omitempty"`
		// ExtendedVersionHistoryChangePolicy : (team_policies) Accepted/opted
		// out of extended version history
		ExtendedVersionHistoryChangePolicy json.RawMessage `json:"extended_version_history_change_policy,omitempty"`
		// FileCommentsChangePolicy : (team_policies) Enabled/disabled
		// commenting on team files
		FileCommentsChangePolicy json.RawMessage `json:"file_comments_change_policy,omitempty"`
		// FileRequestsChangePolicy : (team_policies) Enabled/disabled file
		// requests
		FileRequestsChangePolicy json.RawMessage `json:"file_requests_change_policy,omitempty"`
		// FileRequestsEmailsEnabled : (team_policies) Enabled file request
		// emails for everyone (deprecated, no longer logged)
		FileRequestsEmailsEnabled json.RawMessage `json:"file_requests_emails_enabled,omitempty"`
		// FileRequestsEmailsRestrictedToTeamOnly : (team_policies) Enabled file
		// request emails for team (deprecated, no longer logged)
		FileRequestsEmailsRestrictedToTeamOnly json.RawMessage `json:"file_requests_emails_restricted_to_team_only,omitempty"`
		// GoogleSsoChangePolicy : (team_policies) Enabled/disabled Google
		// single sign-on for team
		GoogleSsoChangePolicy json.RawMessage `json:"google_sso_change_policy,omitempty"`
		// GroupUserManagementChangePolicy : (team_policies) Changed who can
		// create groups
		GroupUserManagementChangePolicy json.RawMessage `json:"group_user_management_change_policy,omitempty"`
		// MemberRequestsChangePolicy : (team_policies) Changed whether users
		// can find team when not invited
		MemberRequestsChangePolicy json.RawMessage `json:"member_requests_change_policy,omitempty"`
		// MemberSpaceLimitsAddException : (team_policies) Added members to
		// member space limit exception list
		MemberSpaceLimitsAddException json.RawMessage `json:"member_space_limits_add_exception,omitempty"`
		// MemberSpaceLimitsChangeCapsTypePolicy : (team_policies) Changed
		// member space limit type for team
		MemberSpaceLimitsChangeCapsTypePolicy json.RawMessage `json:"member_space_limits_change_caps_type_policy,omitempty"`
		// MemberSpaceLimitsChangePolicy : (team_policies) Changed team default
		// member space limit
		MemberSpaceLimitsChangePolicy json.RawMessage `json:"member_space_limits_change_policy,omitempty"`
		// MemberSpaceLimitsRemoveException : (team_policies) Removed members
		// from member space limit exception list
		MemberSpaceLimitsRemoveException json.RawMessage `json:"member_space_limits_remove_exception,omitempty"`
		// MemberSuggestionsChangePolicy : (team_policies) Enabled/disabled
		// option for team members to suggest people to add to team
		MemberSuggestionsChangePolicy json.RawMessage `json:"member_suggestions_change_policy,omitempty"`
		// MicrosoftOfficeAddinChangePolicy : (team_policies) Enabled/disabled
		// Microsoft Office add-in
		MicrosoftOfficeAddinChangePolicy json.RawMessage `json:"microsoft_office_addin_change_policy,omitempty"`
		// NetworkControlChangePolicy : (team_policies) Enabled/disabled network
		// control
		NetworkControlChangePolicy json.RawMessage `json:"network_control_change_policy,omitempty"`
		// PaperChangeDeploymentPolicy : (team_policies) Changed whether Dropbox
		// Paper, when enabled, is deployed to all members or to specific
		// members
		PaperChangeDeploymentPolicy json.RawMessage `json:"paper_change_deployment_policy,omitempty"`
		// PaperChangeMemberLinkPolicy : (team_policies) Changed whether
		// non-members can view Paper docs with link (deprecated, no longer
		// logged)
		PaperChangeMemberLinkPolicy json.RawMessage `json:"paper_change_member_link_policy,omitempty"`
		// PaperChangeMemberPolicy : (team_policies) Changed whether members can
		// share Paper docs outside team, and if docs are accessible only by
		// team members or anyone by default
		PaperChangeMemberPolicy json.RawMessage `json:"paper_change_member_policy,omitempty"`
		// PaperChangePolicy : (team_policies) Enabled/disabled Dropbox Paper
		// for team
		PaperChangePolicy json.RawMessage `json:"paper_change_policy,omitempty"`
		// PaperEnabledUsersGroupAddition : (team_policies) Added users to
		// Paper-enabled users list
		PaperEnabledUsersGroupAddition json.RawMessage `json:"paper_enabled_users_group_addition,omitempty"`
		// PaperEnabledUsersGroupRemoval : (team_policies) Removed users from
		// Paper-enabled users list
		PaperEnabledUsersGroupRemoval json.RawMessage `json:"paper_enabled_users_group_removal,omitempty"`
		// PermanentDeleteChangePolicy : (team_policies) Enabled/disabled
		// ability of team members to permanently delete content
		PermanentDeleteChangePolicy json.RawMessage `json:"permanent_delete_change_policy,omitempty"`
		// SharingChangeFolderJoinPolicy : (team_policies) Changed whether team
		// members can join shared folders owned outside team
		SharingChangeFolderJoinPolicy json.RawMessage `json:"sharing_change_folder_join_policy,omitempty"`
		// SharingChangeLinkPolicy : (team_policies) Changed whether members can
		// share links outside team, and if links are accessible only by team
		// members or anyone by default
		SharingChangeLinkPolicy json.RawMessage `json:"sharing_change_link_policy,omitempty"`
		// SharingChangeMemberPolicy : (team_policies) Changed whether members
		// can share files/folders outside team
		SharingChangeMemberPolicy json.RawMessage `json:"sharing_change_member_policy,omitempty"`
		// ShowcaseChangeDownloadPolicy : (team_policies) Enabled/disabled
		// downloading files from Dropbox Showcase for team
		ShowcaseChangeDownloadPolicy json.RawMessage `json:"showcase_change_download_policy,omitempty"`
		// ShowcaseChangeEnabledPolicy : (team_policies) Enabled/disabled
		// Dropbox Showcase for team
		ShowcaseChangeEnabledPolicy json.RawMessage `json:"showcase_change_enabled_policy,omitempty"`
		// ShowcaseChangeExternalSharingPolicy : (team_policies)
		// Enabled/disabled sharing Dropbox Showcase externally for team
		ShowcaseChangeExternalSharingPolicy json.RawMessage `json:"showcase_change_external_sharing_policy,omitempty"`
		// SmartSyncChangePolicy : (team_policies) Changed default Smart Sync
		// setting for team members
		SmartSyncChangePolicy json.RawMessage `json:"smart_sync_change_policy,omitempty"`
		// SmartSyncNotOptOut : (team_policies) Opted team into Smart Sync
		SmartSyncNotOptOut json.RawMessage `json:"smart_sync_not_opt_out,omitempty"`
		// SmartSyncOptOut : (team_policies) Opted team out of Smart Sync
		SmartSyncOptOut json.RawMessage `json:"smart_sync_opt_out,omitempty"`
		// SsoChangePolicy : (team_policies) Changed single sign-on setting for
		// team
		SsoChangePolicy json.RawMessage `json:"sso_change_policy,omitempty"`
		// TeamSelectiveSyncPolicyChanged : (team_policies) Enabled/disabled
		// Team Selective Sync for team
		TeamSelectiveSyncPolicyChanged json.RawMessage `json:"team_selective_sync_policy_changed,omitempty"`
		// TfaChangePolicy : (team_policies) Changed two-step verification
		// setting for team
		TfaChangePolicy json.RawMessage `json:"tfa_change_policy,omitempty"`
		// TwoAccountChangePolicy : (team_policies) Enabled/disabled option for
		// members to link personal Dropbox account and team account to same
		// computer
		TwoAccountChangePolicy json.RawMessage `json:"two_account_change_policy,omitempty"`
		// ViewerInfoPolicyChanged : (team_policies) Changed team policy for
		// viewer info
		ViewerInfoPolicyChanged json.RawMessage `json:"viewer_info_policy_changed,omitempty"`
		// WebSessionsChangeFixedLengthPolicy : (team_policies) Changed how long
		// members can stay signed in to Dropbox.com
		WebSessionsChangeFixedLengthPolicy json.RawMessage `json:"web_sessions_change_fixed_length_policy,omitempty"`
		// WebSessionsChangeIdleLengthPolicy : (team_policies) Changed how long
		// team members can be idle while signed in to Dropbox.com
		WebSessionsChangeIdleLengthPolicy json.RawMessage `json:"web_sessions_change_idle_length_policy,omitempty"`
		// TeamMergeFrom : (team_profile) Merged another team into this team
		TeamMergeFrom json.RawMessage `json:"team_merge_from,omitempty"`
		// TeamMergeTo : (team_profile) Merged this team into another team
		TeamMergeTo json.RawMessage `json:"team_merge_to,omitempty"`
		// TeamProfileAddLogo : (team_profile) Added team logo to display on
		// shared link headers
		TeamProfileAddLogo json.RawMessage `json:"team_profile_add_logo,omitempty"`
		// TeamProfileChangeDefaultLanguage : (team_profile) Changed default
		// language for team
		TeamProfileChangeDefaultLanguage json.RawMessage `json:"team_profile_change_default_language,omitempty"`
		// TeamProfileChangeLogo : (team_profile) Changed team logo displayed on
		// shared link headers
		TeamProfileChangeLogo json.RawMessage `json:"team_profile_change_logo,omitempty"`
		// TeamProfileChangeName : (team_profile) Changed team name
		TeamProfileChangeName json.RawMessage `json:"team_profile_change_name,omitempty"`
		// TeamProfileRemoveLogo : (team_profile) Removed team logo displayed on
		// shared link headers
		TeamProfileRemoveLogo json.RawMessage `json:"team_profile_remove_logo,omitempty"`
		// TfaAddBackupPhone : (tfa) Added backup phone for two-step
		// verification
		TfaAddBackupPhone json.RawMessage `json:"tfa_add_backup_phone,omitempty"`
		// TfaAddSecurityKey : (tfa) Added security key for two-step
		// verification
		TfaAddSecurityKey json.RawMessage `json:"tfa_add_security_key,omitempty"`
		// TfaChangeBackupPhone : (tfa) Changed backup phone for two-step
		// verification
		TfaChangeBackupPhone json.RawMessage `json:"tfa_change_backup_phone,omitempty"`
		// TfaChangeStatus : (tfa) Enabled/disabled/changed two-step
		// verification setting
		TfaChangeStatus json.RawMessage `json:"tfa_change_status,omitempty"`
		// TfaRemoveBackupPhone : (tfa) Removed backup phone for two-step
		// verification
		TfaRemoveBackupPhone json.RawMessage `json:"tfa_remove_backup_phone,omitempty"`
		// TfaRemoveSecurityKey : (tfa) Removed security key for two-step
		// verification
		TfaRemoveSecurityKey json.RawMessage `json:"tfa_remove_security_key,omitempty"`
		// TfaReset : (tfa) Reset two-step verification for team member
		TfaReset json.RawMessage `json:"tfa_reset,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "app_link_team":
		err = json.Unmarshal(body, &u.AppLinkTeam)

		if err != nil {
			return err
		}
	case "app_link_user":
		err = json.Unmarshal(body, &u.AppLinkUser)

		if err != nil {
			return err
		}
	case "app_unlink_team":
		err = json.Unmarshal(body, &u.AppUnlinkTeam)

		if err != nil {
			return err
		}
	case "app_unlink_user":
		err = json.Unmarshal(body, &u.AppUnlinkUser)

		if err != nil {
			return err
		}
	case "file_add_comment":
		err = json.Unmarshal(body, &u.FileAddComment)

		if err != nil {
			return err
		}
	case "file_change_comment_subscription":
		err = json.Unmarshal(body, &u.FileChangeCommentSubscription)

		if err != nil {
			return err
		}
	case "file_delete_comment":
		err = json.Unmarshal(body, &u.FileDeleteComment)

		if err != nil {
			return err
		}
	case "file_edit_comment":
		err = json.Unmarshal(body, &u.FileEditComment)

		if err != nil {
			return err
		}
	case "file_like_comment":
		err = json.Unmarshal(body, &u.FileLikeComment)

		if err != nil {
			return err
		}
	case "file_resolve_comment":
		err = json.Unmarshal(body, &u.FileResolveComment)

		if err != nil {
			return err
		}
	case "file_unlike_comment":
		err = json.Unmarshal(body, &u.FileUnlikeComment)

		if err != nil {
			return err
		}
	case "file_unresolve_comment":
		err = json.Unmarshal(body, &u.FileUnresolveComment)

		if err != nil {
			return err
		}
	case "device_change_ip_desktop":
		err = json.Unmarshal(body, &u.DeviceChangeIpDesktop)

		if err != nil {
			return err
		}
	case "device_change_ip_mobile":
		err = json.Unmarshal(body, &u.DeviceChangeIpMobile)

		if err != nil {
			return err
		}
	case "device_change_ip_web":
		err = json.Unmarshal(body, &u.DeviceChangeIpWeb)

		if err != nil {
			return err
		}
	case "device_delete_on_unlink_fail":
		err = json.Unmarshal(body, &u.DeviceDeleteOnUnlinkFail)

		if err != nil {
			return err
		}
	case "device_delete_on_unlink_success":
		err = json.Unmarshal(body, &u.DeviceDeleteOnUnlinkSuccess)

		if err != nil {
			return err
		}
	case "device_link_fail":
		err = json.Unmarshal(body, &u.DeviceLinkFail)

		if err != nil {
			return err
		}
	case "device_link_success":
		err = json.Unmarshal(body, &u.DeviceLinkSuccess)

		if err != nil {
			return err
		}
	case "device_management_disabled":
		err = json.Unmarshal(body, &u.DeviceManagementDisabled)

		if err != nil {
			return err
		}
	case "device_management_enabled":
		err = json.Unmarshal(body, &u.DeviceManagementEnabled)

		if err != nil {
			return err
		}
	case "device_unlink":
		err = json.Unmarshal(body, &u.DeviceUnlink)

		if err != nil {
			return err
		}
	case "emm_refresh_auth_token":
		err = json.Unmarshal(body, &u.EmmRefreshAuthToken)

		if err != nil {
			return err
		}
	case "account_capture_change_availability":
		err = json.Unmarshal(body, &u.AccountCaptureChangeAvailability)

		if err != nil {
			return err
		}
	case "account_capture_migrate_account":
		err = json.Unmarshal(body, &u.AccountCaptureMigrateAccount)

		if err != nil {
			return err
		}
	case "account_capture_notification_emails_sent":
		err = json.Unmarshal(body, &u.AccountCaptureNotificationEmailsSent)

		if err != nil {
			return err
		}
	case "account_capture_relinquish_account":
		err = json.Unmarshal(body, &u.AccountCaptureRelinquishAccount)

		if err != nil {
			return err
		}
	case "disabled_domain_invites":
		err = json.Unmarshal(body, &u.DisabledDomainInvites)

		if err != nil {
			return err
		}
	case "domain_invites_approve_request_to_join_team":
		err = json.Unmarshal(body, &u.DomainInvitesApproveRequestToJoinTeam)

		if err != nil {
			return err
		}
	case "domain_invites_decline_request_to_join_team":
		err = json.Unmarshal(body, &u.DomainInvitesDeclineRequestToJoinTeam)

		if err != nil {
			return err
		}
	case "domain_invites_email_existing_users":
		err = json.Unmarshal(body, &u.DomainInvitesEmailExistingUsers)

		if err != nil {
			return err
		}
	case "domain_invites_request_to_join_team":
		err = json.Unmarshal(body, &u.DomainInvitesRequestToJoinTeam)

		if err != nil {
			return err
		}
	case "domain_invites_set_invite_new_user_pref_to_no":
		err = json.Unmarshal(body, &u.DomainInvitesSetInviteNewUserPrefToNo)

		if err != nil {
			return err
		}
	case "domain_invites_set_invite_new_user_pref_to_yes":
		err = json.Unmarshal(body, &u.DomainInvitesSetInviteNewUserPrefToYes)

		if err != nil {
			return err
		}
	case "domain_verification_add_domain_fail":
		err = json.Unmarshal(body, &u.DomainVerificationAddDomainFail)

		if err != nil {
			return err
		}
	case "domain_verification_add_domain_success":
		err = json.Unmarshal(body, &u.DomainVerificationAddDomainSuccess)

		if err != nil {
			return err
		}
	case "domain_verification_remove_domain":
		err = json.Unmarshal(body, &u.DomainVerificationRemoveDomain)

		if err != nil {
			return err
		}
	case "enabled_domain_invites":
		err = json.Unmarshal(body, &u.EnabledDomainInvites)

		if err != nil {
			return err
		}
	case "create_folder":
		err = json.Unmarshal(body, &u.CreateFolder)

		if err != nil {
			return err
		}
	case "file_add":
		err = json.Unmarshal(body, &u.FileAdd)

		if err != nil {
			return err
		}
	case "file_copy":
		err = json.Unmarshal(body, &u.FileCopy)

		if err != nil {
			return err
		}
	case "file_delete":
		err = json.Unmarshal(body, &u.FileDelete)

		if err != nil {
			return err
		}
	case "file_download":
		err = json.Unmarshal(body, &u.FileDownload)

		if err != nil {
			return err
		}
	case "file_edit":
		err = json.Unmarshal(body, &u.FileEdit)

		if err != nil {
			return err
		}
	case "file_get_copy_reference":
		err = json.Unmarshal(body, &u.FileGetCopyReference)

		if err != nil {
			return err
		}
	case "file_move":
		err = json.Unmarshal(body, &u.FileMove)

		if err != nil {
			return err
		}
	case "file_permanently_delete":
		err = json.Unmarshal(body, &u.FilePermanentlyDelete)

		if err != nil {
			return err
		}
	case "file_preview":
		err = json.Unmarshal(body, &u.FilePreview)

		if err != nil {
			return err
		}
	case "file_rename":
		err = json.Unmarshal(body, &u.FileRename)

		if err != nil {
			return err
		}
	case "file_restore":
		err = json.Unmarshal(body, &u.FileRestore)

		if err != nil {
			return err
		}
	case "file_revert":
		err = json.Unmarshal(body, &u.FileRevert)

		if err != nil {
			return err
		}
	case "file_rollback_changes":
		err = json.Unmarshal(body, &u.FileRollbackChanges)

		if err != nil {
			return err
		}
	case "file_save_copy_reference":
		err = json.Unmarshal(body, &u.FileSaveCopyReference)

		if err != nil {
			return err
		}
	case "file_request_change":
		err = json.Unmarshal(body, &u.FileRequestChange)

		if err != nil {
			return err
		}
	case "file_request_close":
		err = json.Unmarshal(body, &u.FileRequestClose)

		if err != nil {
			return err
		}
	case "file_request_create":
		err = json.Unmarshal(body, &u.FileRequestCreate)

		if err != nil {
			return err
		}
	case "file_request_receive_file":
		err = json.Unmarshal(body, &u.FileRequestReceiveFile)

		if err != nil {
			return err
		}
	case "group_add_external_id":
		err = json.Unmarshal(body, &u.GroupAddExternalId)

		if err != nil {
			return err
		}
	case "group_add_member":
		err = json.Unmarshal(body, &u.GroupAddMember)

		if err != nil {
			return err
		}
	case "group_change_external_id":
		err = json.Unmarshal(body, &u.GroupChangeExternalId)

		if err != nil {
			return err
		}
	case "group_change_management_type":
		err = json.Unmarshal(body, &u.GroupChangeManagementType)

		if err != nil {
			return err
		}
	case "group_change_member_role":
		err = json.Unmarshal(body, &u.GroupChangeMemberRole)

		if err != nil {
			return err
		}
	case "group_create":
		err = json.Unmarshal(body, &u.GroupCreate)

		if err != nil {
			return err
		}
	case "group_delete":
		err = json.Unmarshal(body, &u.GroupDelete)

		if err != nil {
			return err
		}
	case "group_description_updated":
		err = json.Unmarshal(body, &u.GroupDescriptionUpdated)

		if err != nil {
			return err
		}
	case "group_join_policy_updated":
		err = json.Unmarshal(body, &u.GroupJoinPolicyUpdated)

		if err != nil {
			return err
		}
	case "group_moved":
		err = json.Unmarshal(body, &u.GroupMoved)

		if err != nil {
			return err
		}
	case "group_remove_external_id":
		err = json.Unmarshal(body, &u.GroupRemoveExternalId)

		if err != nil {
			return err
		}
	case "group_remove_member":
		err = json.Unmarshal(body, &u.GroupRemoveMember)

		if err != nil {
			return err
		}
	case "group_rename":
		err = json.Unmarshal(body, &u.GroupRename)

		if err != nil {
			return err
		}
	case "emm_error":
		err = json.Unmarshal(body, &u.EmmError)

		if err != nil {
			return err
		}
	case "login_fail":
		err = json.Unmarshal(body, &u.LoginFail)

		if err != nil {
			return err
		}
	case "login_success":
		err = json.Unmarshal(body, &u.LoginSuccess)

		if err != nil {
			return err
		}
	case "logout":
		err = json.Unmarshal(body, &u.Logout)

		if err != nil {
			return err
		}
	case "reseller_support_session_end":
		err = json.Unmarshal(body, &u.ResellerSupportSessionEnd)

		if err != nil {
			return err
		}
	case "reseller_support_session_start":
		err = json.Unmarshal(body, &u.ResellerSupportSessionStart)

		if err != nil {
			return err
		}
	case "sign_in_as_session_end":
		err = json.Unmarshal(body, &u.SignInAsSessionEnd)

		if err != nil {
			return err
		}
	case "sign_in_as_session_start":
		err = json.Unmarshal(body, &u.SignInAsSessionStart)

		if err != nil {
			return err
		}
	case "sso_error":
		err = json.Unmarshal(body, &u.SsoError)

		if err != nil {
			return err
		}
	case "member_add_name":
		err = json.Unmarshal(body, &u.MemberAddName)

		if err != nil {
			return err
		}
	case "member_change_admin_role":
		err = json.Unmarshal(body, &u.MemberChangeAdminRole)

		if err != nil {
			return err
		}
	case "member_change_email":
		err = json.Unmarshal(body, &u.MemberChangeEmail)

		if err != nil {
			return err
		}
	case "member_change_membership_type":
		err = json.Unmarshal(body, &u.MemberChangeMembershipType)

		if err != nil {
			return err
		}
	case "member_change_name":
		err = json.Unmarshal(body, &u.MemberChangeName)

		if err != nil {
			return err
		}
	case "member_change_status":
		err = json.Unmarshal(body, &u.MemberChangeStatus)

		if err != nil {
			return err
		}
	case "member_delete_manual_contacts":
		err = json.Unmarshal(body, &u.MemberDeleteManualContacts)

		if err != nil {
			return err
		}
	case "member_permanently_delete_account_contents":
		err = json.Unmarshal(body, &u.MemberPermanentlyDeleteAccountContents)

		if err != nil {
			return err
		}
	case "member_space_limits_add_custom_quota":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsAddCustomQuota)

		if err != nil {
			return err
		}
	case "member_space_limits_change_custom_quota":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeCustomQuota)

		if err != nil {
			return err
		}
	case "member_space_limits_change_status":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeStatus)

		if err != nil {
			return err
		}
	case "member_space_limits_remove_custom_quota":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsRemoveCustomQuota)

		if err != nil {
			return err
		}
	case "member_suggest":
		err = json.Unmarshal(body, &u.MemberSuggest)

		if err != nil {
			return err
		}
	case "member_transfer_account_contents":
		err = json.Unmarshal(body, &u.MemberTransferAccountContents)

		if err != nil {
			return err
		}
	case "secondary_mails_policy_changed":
		err = json.Unmarshal(body, &u.SecondaryMailsPolicyChanged)

		if err != nil {
			return err
		}
	case "paper_content_add_member":
		err = json.Unmarshal(body, &u.PaperContentAddMember)

		if err != nil {
			return err
		}
	case "paper_content_add_to_folder":
		err = json.Unmarshal(body, &u.PaperContentAddToFolder)

		if err != nil {
			return err
		}
	case "paper_content_archive":
		err = json.Unmarshal(body, &u.PaperContentArchive)

		if err != nil {
			return err
		}
	case "paper_content_create":
		err = json.Unmarshal(body, &u.PaperContentCreate)

		if err != nil {
			return err
		}
	case "paper_content_permanently_delete":
		err = json.Unmarshal(body, &u.PaperContentPermanentlyDelete)

		if err != nil {
			return err
		}
	case "paper_content_remove_from_folder":
		err = json.Unmarshal(body, &u.PaperContentRemoveFromFolder)

		if err != nil {
			return err
		}
	case "paper_content_remove_member":
		err = json.Unmarshal(body, &u.PaperContentRemoveMember)

		if err != nil {
			return err
		}
	case "paper_content_rename":
		err = json.Unmarshal(body, &u.PaperContentRename)

		if err != nil {
			return err
		}
	case "paper_content_restore":
		err = json.Unmarshal(body, &u.PaperContentRestore)

		if err != nil {
			return err
		}
	case "paper_doc_add_comment":
		err = json.Unmarshal(body, &u.PaperDocAddComment)

		if err != nil {
			return err
		}
	case "paper_doc_change_member_role":
		err = json.Unmarshal(body, &u.PaperDocChangeMemberRole)

		if err != nil {
			return err
		}
	case "paper_doc_change_sharing_policy":
		err = json.Unmarshal(body, &u.PaperDocChangeSharingPolicy)

		if err != nil {
			return err
		}
	case "paper_doc_change_subscription":
		err = json.Unmarshal(body, &u.PaperDocChangeSubscription)

		if err != nil {
			return err
		}
	case "paper_doc_deleted":
		err = json.Unmarshal(body, &u.PaperDocDeleted)

		if err != nil {
			return err
		}
	case "paper_doc_delete_comment":
		err = json.Unmarshal(body, &u.PaperDocDeleteComment)

		if err != nil {
			return err
		}
	case "paper_doc_download":
		err = json.Unmarshal(body, &u.PaperDocDownload)

		if err != nil {
			return err
		}
	case "paper_doc_edit":
		err = json.Unmarshal(body, &u.PaperDocEdit)

		if err != nil {
			return err
		}
	case "paper_doc_edit_comment":
		err = json.Unmarshal(body, &u.PaperDocEditComment)

		if err != nil {
			return err
		}
	case "paper_doc_followed":
		err = json.Unmarshal(body, &u.PaperDocFollowed)

		if err != nil {
			return err
		}
	case "paper_doc_mention":
		err = json.Unmarshal(body, &u.PaperDocMention)

		if err != nil {
			return err
		}
	case "paper_doc_ownership_changed":
		err = json.Unmarshal(body, &u.PaperDocOwnershipChanged)

		if err != nil {
			return err
		}
	case "paper_doc_request_access":
		err = json.Unmarshal(body, &u.PaperDocRequestAccess)

		if err != nil {
			return err
		}
	case "paper_doc_resolve_comment":
		err = json.Unmarshal(body, &u.PaperDocResolveComment)

		if err != nil {
			return err
		}
	case "paper_doc_revert":
		err = json.Unmarshal(body, &u.PaperDocRevert)

		if err != nil {
			return err
		}
	case "paper_doc_slack_share":
		err = json.Unmarshal(body, &u.PaperDocSlackShare)

		if err != nil {
			return err
		}
	case "paper_doc_team_invite":
		err = json.Unmarshal(body, &u.PaperDocTeamInvite)

		if err != nil {
			return err
		}
	case "paper_doc_trashed":
		err = json.Unmarshal(body, &u.PaperDocTrashed)

		if err != nil {
			return err
		}
	case "paper_doc_unresolve_comment":
		err = json.Unmarshal(body, &u.PaperDocUnresolveComment)

		if err != nil {
			return err
		}
	case "paper_doc_untrashed":
		err = json.Unmarshal(body, &u.PaperDocUntrashed)

		if err != nil {
			return err
		}
	case "paper_doc_view":
		err = json.Unmarshal(body, &u.PaperDocView)

		if err != nil {
			return err
		}
	case "paper_external_view_allow":
		err = json.Unmarshal(body, &u.PaperExternalViewAllow)

		if err != nil {
			return err
		}
	case "paper_external_view_default_team":
		err = json.Unmarshal(body, &u.PaperExternalViewDefaultTeam)

		if err != nil {
			return err
		}
	case "paper_external_view_forbid":
		err = json.Unmarshal(body, &u.PaperExternalViewForbid)

		if err != nil {
			return err
		}
	case "paper_folder_change_subscription":
		err = json.Unmarshal(body, &u.PaperFolderChangeSubscription)

		if err != nil {
			return err
		}
	case "paper_folder_deleted":
		err = json.Unmarshal(body, &u.PaperFolderDeleted)

		if err != nil {
			return err
		}
	case "paper_folder_followed":
		err = json.Unmarshal(body, &u.PaperFolderFollowed)

		if err != nil {
			return err
		}
	case "paper_folder_team_invite":
		err = json.Unmarshal(body, &u.PaperFolderTeamInvite)

		if err != nil {
			return err
		}
	case "password_change":
		err = json.Unmarshal(body, &u.PasswordChange)

		if err != nil {
			return err
		}
	case "password_reset":
		err = json.Unmarshal(body, &u.PasswordReset)

		if err != nil {
			return err
		}
	case "password_reset_all":
		err = json.Unmarshal(body, &u.PasswordResetAll)

		if err != nil {
			return err
		}
	case "emm_create_exceptions_report":
		err = json.Unmarshal(body, &u.EmmCreateExceptionsReport)

		if err != nil {
			return err
		}
	case "emm_create_usage_report":
		err = json.Unmarshal(body, &u.EmmCreateUsageReport)

		if err != nil {
			return err
		}
	case "export_members_report":
		err = json.Unmarshal(body, &u.ExportMembersReport)

		if err != nil {
			return err
		}
	case "paper_admin_export_start":
		err = json.Unmarshal(body, &u.PaperAdminExportStart)

		if err != nil {
			return err
		}
	case "smart_sync_create_admin_privilege_report":
		err = json.Unmarshal(body, &u.SmartSyncCreateAdminPrivilegeReport)

		if err != nil {
			return err
		}
	case "team_activity_create_report":
		err = json.Unmarshal(body, &u.TeamActivityCreateReport)

		if err != nil {
			return err
		}
	case "collection_share":
		err = json.Unmarshal(body, &u.CollectionShare)

		if err != nil {
			return err
		}
	case "note_acl_invite_only":
		err = json.Unmarshal(body, &u.NoteAclInviteOnly)

		if err != nil {
			return err
		}
	case "note_acl_link":
		err = json.Unmarshal(body, &u.NoteAclLink)

		if err != nil {
			return err
		}
	case "note_acl_team_link":
		err = json.Unmarshal(body, &u.NoteAclTeamLink)

		if err != nil {
			return err
		}
	case "note_shared":
		err = json.Unmarshal(body, &u.NoteShared)

		if err != nil {
			return err
		}
	case "note_share_receive":
		err = json.Unmarshal(body, &u.NoteShareReceive)

		if err != nil {
			return err
		}
	case "open_note_shared":
		err = json.Unmarshal(body, &u.OpenNoteShared)

		if err != nil {
			return err
		}
	case "sf_add_group":
		err = json.Unmarshal(body, &u.SfAddGroup)

		if err != nil {
			return err
		}
	case "sf_allow_non_members_to_view_shared_links":
		err = json.Unmarshal(body, &u.SfAllowNonMembersToViewSharedLinks)

		if err != nil {
			return err
		}
	case "sf_external_invite_warn":
		err = json.Unmarshal(body, &u.SfExternalInviteWarn)

		if err != nil {
			return err
		}
	case "sf_fb_invite":
		err = json.Unmarshal(body, &u.SfFbInvite)

		if err != nil {
			return err
		}
	case "sf_fb_invite_change_role":
		err = json.Unmarshal(body, &u.SfFbInviteChangeRole)

		if err != nil {
			return err
		}
	case "sf_fb_uninvite":
		err = json.Unmarshal(body, &u.SfFbUninvite)

		if err != nil {
			return err
		}
	case "sf_invite_group":
		err = json.Unmarshal(body, &u.SfInviteGroup)

		if err != nil {
			return err
		}
	case "sf_team_grant_access":
		err = json.Unmarshal(body, &u.SfTeamGrantAccess)

		if err != nil {
			return err
		}
	case "sf_team_invite":
		err = json.Unmarshal(body, &u.SfTeamInvite)

		if err != nil {
			return err
		}
	case "sf_team_invite_change_role":
		err = json.Unmarshal(body, &u.SfTeamInviteChangeRole)

		if err != nil {
			return err
		}
	case "sf_team_join":
		err = json.Unmarshal(body, &u.SfTeamJoin)

		if err != nil {
			return err
		}
	case "sf_team_join_from_oob_link":
		err = json.Unmarshal(body, &u.SfTeamJoinFromOobLink)

		if err != nil {
			return err
		}
	case "sf_team_uninvite":
		err = json.Unmarshal(body, &u.SfTeamUninvite)

		if err != nil {
			return err
		}
	case "shared_content_add_invitees":
		err = json.Unmarshal(body, &u.SharedContentAddInvitees)

		if err != nil {
			return err
		}
	case "shared_content_add_link_expiry":
		err = json.Unmarshal(body, &u.SharedContentAddLinkExpiry)

		if err != nil {
			return err
		}
	case "shared_content_add_link_password":
		err = json.Unmarshal(body, &u.SharedContentAddLinkPassword)

		if err != nil {
			return err
		}
	case "shared_content_add_member":
		err = json.Unmarshal(body, &u.SharedContentAddMember)

		if err != nil {
			return err
		}
	case "shared_content_change_downloads_policy":
		err = json.Unmarshal(body, &u.SharedContentChangeDownloadsPolicy)

		if err != nil {
			return err
		}
	case "shared_content_change_invitee_role":
		err = json.Unmarshal(body, &u.SharedContentChangeInviteeRole)

		if err != nil {
			return err
		}
	case "shared_content_change_link_audience":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkAudience)

		if err != nil {
			return err
		}
	case "shared_content_change_link_expiry":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkExpiry)

		if err != nil {
			return err
		}
	case "shared_content_change_link_password":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkPassword)

		if err != nil {
			return err
		}
	case "shared_content_change_member_role":
		err = json.Unmarshal(body, &u.SharedContentChangeMemberRole)

		if err != nil {
			return err
		}
	case "shared_content_change_viewer_info_policy":
		err = json.Unmarshal(body, &u.SharedContentChangeViewerInfoPolicy)

		if err != nil {
			return err
		}
	case "shared_content_claim_invitation":
		err = json.Unmarshal(body, &u.SharedContentClaimInvitation)

		if err != nil {
			return err
		}
	case "shared_content_copy":
		err = json.Unmarshal(body, &u.SharedContentCopy)

		if err != nil {
			return err
		}
	case "shared_content_download":
		err = json.Unmarshal(body, &u.SharedContentDownload)

		if err != nil {
			return err
		}
	case "shared_content_relinquish_membership":
		err = json.Unmarshal(body, &u.SharedContentRelinquishMembership)

		if err != nil {
			return err
		}
	case "shared_content_remove_invitees":
		err = json.Unmarshal(body, &u.SharedContentRemoveInvitees)

		if err != nil {
			return err
		}
	case "shared_content_remove_link_expiry":
		err = json.Unmarshal(body, &u.SharedContentRemoveLinkExpiry)

		if err != nil {
			return err
		}
	case "shared_content_remove_link_password":
		err = json.Unmarshal(body, &u.SharedContentRemoveLinkPassword)

		if err != nil {
			return err
		}
	case "shared_content_remove_member":
		err = json.Unmarshal(body, &u.SharedContentRemoveMember)

		if err != nil {
			return err
		}
	case "shared_content_request_access":
		err = json.Unmarshal(body, &u.SharedContentRequestAccess)

		if err != nil {
			return err
		}
	case "shared_content_unshare":
		err = json.Unmarshal(body, &u.SharedContentUnshare)

		if err != nil {
			return err
		}
	case "shared_content_view":
		err = json.Unmarshal(body, &u.SharedContentView)

		if err != nil {
			return err
		}
	case "shared_folder_change_link_policy":
		err = json.Unmarshal(body, &u.SharedFolderChangeLinkPolicy)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_inheritance_policy":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersInheritancePolicy)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_management_policy":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersManagementPolicy)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_policy":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersPolicy)

		if err != nil {
			return err
		}
	case "shared_folder_create":
		err = json.Unmarshal(body, &u.SharedFolderCreate)

		if err != nil {
			return err
		}
	case "shared_folder_decline_invitation":
		err = json.Unmarshal(body, &u.SharedFolderDeclineInvitation)

		if err != nil {
			return err
		}
	case "shared_folder_mount":
		err = json.Unmarshal(body, &u.SharedFolderMount)

		if err != nil {
			return err
		}
	case "shared_folder_nest":
		err = json.Unmarshal(body, &u.SharedFolderNest)

		if err != nil {
			return err
		}
	case "shared_folder_transfer_ownership":
		err = json.Unmarshal(body, &u.SharedFolderTransferOwnership)

		if err != nil {
			return err
		}
	case "shared_folder_unmount":
		err = json.Unmarshal(body, &u.SharedFolderUnmount)

		if err != nil {
			return err
		}
	case "shared_link_add_expiry":
		err = json.Unmarshal(body, &u.SharedLinkAddExpiry)

		if err != nil {
			return err
		}
	case "shared_link_change_expiry":
		err = json.Unmarshal(body, &u.SharedLinkChangeExpiry)

		if err != nil {
			return err
		}
	case "shared_link_change_visibility":
		err = json.Unmarshal(body, &u.SharedLinkChangeVisibility)

		if err != nil {
			return err
		}
	case "shared_link_copy":
		err = json.Unmarshal(body, &u.SharedLinkCopy)

		if err != nil {
			return err
		}
	case "shared_link_create":
		err = json.Unmarshal(body, &u.SharedLinkCreate)

		if err != nil {
			return err
		}
	case "shared_link_disable":
		err = json.Unmarshal(body, &u.SharedLinkDisable)

		if err != nil {
			return err
		}
	case "shared_link_download":
		err = json.Unmarshal(body, &u.SharedLinkDownload)

		if err != nil {
			return err
		}
	case "shared_link_remove_expiry":
		err = json.Unmarshal(body, &u.SharedLinkRemoveExpiry)

		if err != nil {
			return err
		}
	case "shared_link_share":
		err = json.Unmarshal(body, &u.SharedLinkShare)

		if err != nil {
			return err
		}
	case "shared_link_view":
		err = json.Unmarshal(body, &u.SharedLinkView)

		if err != nil {
			return err
		}
	case "shared_note_opened":
		err = json.Unmarshal(body, &u.SharedNoteOpened)

		if err != nil {
			return err
		}
	case "shmodel_group_share":
		err = json.Unmarshal(body, &u.ShmodelGroupShare)

		if err != nil {
			return err
		}
	case "showcase_access_granted":
		err = json.Unmarshal(body, &u.ShowcaseAccessGranted)

		if err != nil {
			return err
		}
	case "showcase_add_member":
		err = json.Unmarshal(body, &u.ShowcaseAddMember)

		if err != nil {
			return err
		}
	case "showcase_archived":
		err = json.Unmarshal(body, &u.ShowcaseArchived)

		if err != nil {
			return err
		}
	case "showcase_created":
		err = json.Unmarshal(body, &u.ShowcaseCreated)

		if err != nil {
			return err
		}
	case "showcase_delete_comment":
		err = json.Unmarshal(body, &u.ShowcaseDeleteComment)

		if err != nil {
			return err
		}
	case "showcase_edited":
		err = json.Unmarshal(body, &u.ShowcaseEdited)

		if err != nil {
			return err
		}
	case "showcase_edit_comment":
		err = json.Unmarshal(body, &u.ShowcaseEditComment)

		if err != nil {
			return err
		}
	case "showcase_file_added":
		err = json.Unmarshal(body, &u.ShowcaseFileAdded)

		if err != nil {
			return err
		}
	case "showcase_file_download":
		err = json.Unmarshal(body, &u.ShowcaseFileDownload)

		if err != nil {
			return err
		}
	case "showcase_file_removed":
		err = json.Unmarshal(body, &u.ShowcaseFileRemoved)

		if err != nil {
			return err
		}
	case "showcase_file_view":
		err = json.Unmarshal(body, &u.ShowcaseFileView)

		if err != nil {
			return err
		}
	case "showcase_permanently_deleted":
		err = json.Unmarshal(body, &u.ShowcasePermanentlyDeleted)

		if err != nil {
			return err
		}
	case "showcase_post_comment":
		err = json.Unmarshal(body, &u.ShowcasePostComment)

		if err != nil {
			return err
		}
	case "showcase_remove_member":
		err = json.Unmarshal(body, &u.ShowcaseRemoveMember)

		if err != nil {
			return err
		}
	case "showcase_renamed":
		err = json.Unmarshal(body, &u.ShowcaseRenamed)

		if err != nil {
			return err
		}
	case "showcase_request_access":
		err = json.Unmarshal(body, &u.ShowcaseRequestAccess)

		if err != nil {
			return err
		}
	case "showcase_resolve_comment":
		err = json.Unmarshal(body, &u.ShowcaseResolveComment)

		if err != nil {
			return err
		}
	case "showcase_restored":
		err = json.Unmarshal(body, &u.ShowcaseRestored)

		if err != nil {
			return err
		}
	case "showcase_trashed":
		err = json.Unmarshal(body, &u.ShowcaseTrashed)

		if err != nil {
			return err
		}
	case "showcase_trashed_deprecated":
		err = json.Unmarshal(body, &u.ShowcaseTrashedDeprecated)

		if err != nil {
			return err
		}
	case "showcase_unresolve_comment":
		err = json.Unmarshal(body, &u.ShowcaseUnresolveComment)

		if err != nil {
			return err
		}
	case "showcase_untrashed":
		err = json.Unmarshal(body, &u.ShowcaseUntrashed)

		if err != nil {
			return err
		}
	case "showcase_untrashed_deprecated":
		err = json.Unmarshal(body, &u.ShowcaseUntrashedDeprecated)

		if err != nil {
			return err
		}
	case "showcase_view":
		err = json.Unmarshal(body, &u.ShowcaseView)

		if err != nil {
			return err
		}
	case "sso_add_cert":
		err = json.Unmarshal(body, &u.SsoAddCert)

		if err != nil {
			return err
		}
	case "sso_add_login_url":
		err = json.Unmarshal(body, &u.SsoAddLoginUrl)

		if err != nil {
			return err
		}
	case "sso_add_logout_url":
		err = json.Unmarshal(body, &u.SsoAddLogoutUrl)

		if err != nil {
			return err
		}
	case "sso_change_cert":
		err = json.Unmarshal(body, &u.SsoChangeCert)

		if err != nil {
			return err
		}
	case "sso_change_login_url":
		err = json.Unmarshal(body, &u.SsoChangeLoginUrl)

		if err != nil {
			return err
		}
	case "sso_change_logout_url":
		err = json.Unmarshal(body, &u.SsoChangeLogoutUrl)

		if err != nil {
			return err
		}
	case "sso_change_saml_identity_mode":
		err = json.Unmarshal(body, &u.SsoChangeSamlIdentityMode)

		if err != nil {
			return err
		}
	case "sso_remove_cert":
		err = json.Unmarshal(body, &u.SsoRemoveCert)

		if err != nil {
			return err
		}
	case "sso_remove_login_url":
		err = json.Unmarshal(body, &u.SsoRemoveLoginUrl)

		if err != nil {
			return err
		}
	case "sso_remove_logout_url":
		err = json.Unmarshal(body, &u.SsoRemoveLogoutUrl)

		if err != nil {
			return err
		}
	case "team_folder_change_status":
		err = json.Unmarshal(body, &u.TeamFolderChangeStatus)

		if err != nil {
			return err
		}
	case "team_folder_create":
		err = json.Unmarshal(body, &u.TeamFolderCreate)

		if err != nil {
			return err
		}
	case "team_folder_downgrade":
		err = json.Unmarshal(body, &u.TeamFolderDowngrade)

		if err != nil {
			return err
		}
	case "team_folder_permanently_delete":
		err = json.Unmarshal(body, &u.TeamFolderPermanentlyDelete)

		if err != nil {
			return err
		}
	case "team_folder_rename":
		err = json.Unmarshal(body, &u.TeamFolderRename)

		if err != nil {
			return err
		}
	case "team_selective_sync_settings_changed":
		err = json.Unmarshal(body, &u.TeamSelectiveSyncSettingsChanged)

		if err != nil {
			return err
		}
	case "account_capture_change_policy":
		err = json.Unmarshal(body, &u.AccountCaptureChangePolicy)

		if err != nil {
			return err
		}
	case "allow_download_disabled":
		err = json.Unmarshal(body, &u.AllowDownloadDisabled)

		if err != nil {
			return err
		}
	case "allow_download_enabled":
		err = json.Unmarshal(body, &u.AllowDownloadEnabled)

		if err != nil {
			return err
		}
	case "camera_uploads_policy_changed":
		err = json.Unmarshal(body, &u.CameraUploadsPolicyChanged)

		if err != nil {
			return err
		}
	case "data_placement_restriction_change_policy":
		err = json.Unmarshal(body, &u.DataPlacementRestrictionChangePolicy)

		if err != nil {
			return err
		}
	case "data_placement_restriction_satisfy_policy":
		err = json.Unmarshal(body, &u.DataPlacementRestrictionSatisfyPolicy)

		if err != nil {
			return err
		}
	case "device_approvals_change_desktop_policy":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeDesktopPolicy)

		if err != nil {
			return err
		}
	case "device_approvals_change_mobile_policy":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeMobilePolicy)

		if err != nil {
			return err
		}
	case "device_approvals_change_overage_action":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeOverageAction)

		if err != nil {
			return err
		}
	case "device_approvals_change_unlink_action":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeUnlinkAction)

		if err != nil {
			return err
		}
	case "directory_restrictions_add_members":
		err = json.Unmarshal(body, &u.DirectoryRestrictionsAddMembers)

		if err != nil {
			return err
		}
	case "directory_restrictions_remove_members":
		err = json.Unmarshal(body, &u.DirectoryRestrictionsRemoveMembers)

		if err != nil {
			return err
		}
	case "emm_add_exception":
		err = json.Unmarshal(body, &u.EmmAddException)

		if err != nil {
			return err
		}
	case "emm_change_policy":
		err = json.Unmarshal(body, &u.EmmChangePolicy)

		if err != nil {
			return err
		}
	case "emm_remove_exception":
		err = json.Unmarshal(body, &u.EmmRemoveException)

		if err != nil {
			return err
		}
	case "extended_version_history_change_policy":
		err = json.Unmarshal(body, &u.ExtendedVersionHistoryChangePolicy)

		if err != nil {
			return err
		}
	case "file_comments_change_policy":
		err = json.Unmarshal(body, &u.FileCommentsChangePolicy)

		if err != nil {
			return err
		}
	case "file_requests_change_policy":
		err = json.Unmarshal(body, &u.FileRequestsChangePolicy)

		if err != nil {
			return err
		}
	case "file_requests_emails_enabled":
		err = json.Unmarshal(body, &u.FileRequestsEmailsEnabled)

		if err != nil {
			return err
		}
	case "file_requests_emails_restricted_to_team_only":
		err = json.Unmarshal(body, &u.FileRequestsEmailsRestrictedToTeamOnly)

		if err != nil {
			return err
		}
	case "google_sso_change_policy":
		err = json.Unmarshal(body, &u.GoogleSsoChangePolicy)

		if err != nil {
			return err
		}
	case "group_user_management_change_policy":
		err = json.Unmarshal(body, &u.GroupUserManagementChangePolicy)

		if err != nil {
			return err
		}
	case "member_requests_change_policy":
		err = json.Unmarshal(body, &u.MemberRequestsChangePolicy)

		if err != nil {
			return err
		}
	case "member_space_limits_add_exception":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsAddException)

		if err != nil {
			return err
		}
	case "member_space_limits_change_caps_type_policy":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeCapsTypePolicy)

		if err != nil {
			return err
		}
	case "member_space_limits_change_policy":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangePolicy)

		if err != nil {
			return err
		}
	case "member_space_limits_remove_exception":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsRemoveException)

		if err != nil {
			return err
		}
	case "member_suggestions_change_policy":
		err = json.Unmarshal(body, &u.MemberSuggestionsChangePolicy)

		if err != nil {
			return err
		}
	case "microsoft_office_addin_change_policy":
		err = json.Unmarshal(body, &u.MicrosoftOfficeAddinChangePolicy)

		if err != nil {
			return err
		}
	case "network_control_change_policy":
		err = json.Unmarshal(body, &u.NetworkControlChangePolicy)

		if err != nil {
			return err
		}
	case "paper_change_deployment_policy":
		err = json.Unmarshal(body, &u.PaperChangeDeploymentPolicy)

		if err != nil {
			return err
		}
	case "paper_change_member_link_policy":
		err = json.Unmarshal(body, &u.PaperChangeMemberLinkPolicy)

		if err != nil {
			return err
		}
	case "paper_change_member_policy":
		err = json.Unmarshal(body, &u.PaperChangeMemberPolicy)

		if err != nil {
			return err
		}
	case "paper_change_policy":
		err = json.Unmarshal(body, &u.PaperChangePolicy)

		if err != nil {
			return err
		}
	case "paper_enabled_users_group_addition":
		err = json.Unmarshal(body, &u.PaperEnabledUsersGroupAddition)

		if err != nil {
			return err
		}
	case "paper_enabled_users_group_removal":
		err = json.Unmarshal(body, &u.PaperEnabledUsersGroupRemoval)

		if err != nil {
			return err
		}
	case "permanent_delete_change_policy":
		err = json.Unmarshal(body, &u.PermanentDeleteChangePolicy)

		if err != nil {
			return err
		}
	case "sharing_change_folder_join_policy":
		err = json.Unmarshal(body, &u.SharingChangeFolderJoinPolicy)

		if err != nil {
			return err
		}
	case "sharing_change_link_policy":
		err = json.Unmarshal(body, &u.SharingChangeLinkPolicy)

		if err != nil {
			return err
		}
	case "sharing_change_member_policy":
		err = json.Unmarshal(body, &u.SharingChangeMemberPolicy)

		if err != nil {
			return err
		}
	case "showcase_change_download_policy":
		err = json.Unmarshal(body, &u.ShowcaseChangeDownloadPolicy)

		if err != nil {
			return err
		}
	case "showcase_change_enabled_policy":
		err = json.Unmarshal(body, &u.ShowcaseChangeEnabledPolicy)

		if err != nil {
			return err
		}
	case "showcase_change_external_sharing_policy":
		err = json.Unmarshal(body, &u.ShowcaseChangeExternalSharingPolicy)

		if err != nil {
			return err
		}
	case "smart_sync_change_policy":
		err = json.Unmarshal(body, &u.SmartSyncChangePolicy)

		if err != nil {
			return err
		}
	case "smart_sync_not_opt_out":
		err = json.Unmarshal(body, &u.SmartSyncNotOptOut)

		if err != nil {
			return err
		}
	case "smart_sync_opt_out":
		err = json.Unmarshal(body, &u.SmartSyncOptOut)

		if err != nil {
			return err
		}
	case "sso_change_policy":
		err = json.Unmarshal(body, &u.SsoChangePolicy)

		if err != nil {
			return err
		}
	case "team_selective_sync_policy_changed":
		err = json.Unmarshal(body, &u.TeamSelectiveSyncPolicyChanged)

		if err != nil {
			return err
		}
	case "tfa_change_policy":
		err = json.Unmarshal(body, &u.TfaChangePolicy)

		if err != nil {
			return err
		}
	case "two_account_change_policy":
		err = json.Unmarshal(body, &u.TwoAccountChangePolicy)

		if err != nil {
			return err
		}
	case "viewer_info_policy_changed":
		err = json.Unmarshal(body, &u.ViewerInfoPolicyChanged)

		if err != nil {
			return err
		}
	case "web_sessions_change_fixed_length_policy":
		err = json.Unmarshal(body, &u.WebSessionsChangeFixedLengthPolicy)

		if err != nil {
			return err
		}
	case "web_sessions_change_idle_length_policy":
		err = json.Unmarshal(body, &u.WebSessionsChangeIdleLengthPolicy)

		if err != nil {
			return err
		}
	case "team_merge_from":
		err = json.Unmarshal(body, &u.TeamMergeFrom)

		if err != nil {
			return err
		}
	case "team_merge_to":
		err = json.Unmarshal(body, &u.TeamMergeTo)

		if err != nil {
			return err
		}
	case "team_profile_add_logo":
		err = json.Unmarshal(body, &u.TeamProfileAddLogo)

		if err != nil {
			return err
		}
	case "team_profile_change_default_language":
		err = json.Unmarshal(body, &u.TeamProfileChangeDefaultLanguage)

		if err != nil {
			return err
		}
	case "team_profile_change_logo":
		err = json.Unmarshal(body, &u.TeamProfileChangeLogo)

		if err != nil {
			return err
		}
	case "team_profile_change_name":
		err = json.Unmarshal(body, &u.TeamProfileChangeName)

		if err != nil {
			return err
		}
	case "team_profile_remove_logo":
		err = json.Unmarshal(body, &u.TeamProfileRemoveLogo)

		if err != nil {
			return err
		}
	case "tfa_add_backup_phone":
		err = json.Unmarshal(body, &u.TfaAddBackupPhone)

		if err != nil {
			return err
		}
	case "tfa_add_security_key":
		err = json.Unmarshal(body, &u.TfaAddSecurityKey)

		if err != nil {
			return err
		}
	case "tfa_change_backup_phone":
		err = json.Unmarshal(body, &u.TfaChangeBackupPhone)

		if err != nil {
			return err
		}
	case "tfa_change_status":
		err = json.Unmarshal(body, &u.TfaChangeStatus)

		if err != nil {
			return err
		}
	case "tfa_remove_backup_phone":
		err = json.Unmarshal(body, &u.TfaRemoveBackupPhone)

		if err != nil {
			return err
		}
	case "tfa_remove_security_key":
		err = json.Unmarshal(body, &u.TfaRemoveSecurityKey)

		if err != nil {
			return err
		}
	case "tfa_reset":
		err = json.Unmarshal(body, &u.TfaReset)

		if err != nil {
			return err
		}
	}
	return nil
}
