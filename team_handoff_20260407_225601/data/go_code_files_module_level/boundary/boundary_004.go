func (u *EventDetails) UnmarshalJSON(body []byte) error {
	type wrap struct {
		dropbox.Tagged
		// AppLinkTeamDetails : has no documentation (yet)
		AppLinkTeamDetails json.RawMessage `json:"app_link_team_details,omitempty"`
		// AppLinkUserDetails : has no documentation (yet)
		AppLinkUserDetails json.RawMessage `json:"app_link_user_details,omitempty"`
		// AppUnlinkTeamDetails : has no documentation (yet)
		AppUnlinkTeamDetails json.RawMessage `json:"app_unlink_team_details,omitempty"`
		// AppUnlinkUserDetails : has no documentation (yet)
		AppUnlinkUserDetails json.RawMessage `json:"app_unlink_user_details,omitempty"`
		// FileAddCommentDetails : has no documentation (yet)
		FileAddCommentDetails json.RawMessage `json:"file_add_comment_details,omitempty"`
		// FileChangeCommentSubscriptionDetails : has no documentation (yet)
		FileChangeCommentSubscriptionDetails json.RawMessage `json:"file_change_comment_subscription_details,omitempty"`
		// FileDeleteCommentDetails : has no documentation (yet)
		FileDeleteCommentDetails json.RawMessage `json:"file_delete_comment_details,omitempty"`
		// FileEditCommentDetails : has no documentation (yet)
		FileEditCommentDetails json.RawMessage `json:"file_edit_comment_details,omitempty"`
		// FileLikeCommentDetails : has no documentation (yet)
		FileLikeCommentDetails json.RawMessage `json:"file_like_comment_details,omitempty"`
		// FileResolveCommentDetails : has no documentation (yet)
		FileResolveCommentDetails json.RawMessage `json:"file_resolve_comment_details,omitempty"`
		// FileUnlikeCommentDetails : has no documentation (yet)
		FileUnlikeCommentDetails json.RawMessage `json:"file_unlike_comment_details,omitempty"`
		// FileUnresolveCommentDetails : has no documentation (yet)
		FileUnresolveCommentDetails json.RawMessage `json:"file_unresolve_comment_details,omitempty"`
		// DeviceChangeIpDesktopDetails : has no documentation (yet)
		DeviceChangeIpDesktopDetails json.RawMessage `json:"device_change_ip_desktop_details,omitempty"`
		// DeviceChangeIpMobileDetails : has no documentation (yet)
		DeviceChangeIpMobileDetails json.RawMessage `json:"device_change_ip_mobile_details,omitempty"`
		// DeviceChangeIpWebDetails : has no documentation (yet)
		DeviceChangeIpWebDetails json.RawMessage `json:"device_change_ip_web_details,omitempty"`
		// DeviceDeleteOnUnlinkFailDetails : has no documentation (yet)
		DeviceDeleteOnUnlinkFailDetails json.RawMessage `json:"device_delete_on_unlink_fail_details,omitempty"`
		// DeviceDeleteOnUnlinkSuccessDetails : has no documentation (yet)
		DeviceDeleteOnUnlinkSuccessDetails json.RawMessage `json:"device_delete_on_unlink_success_details,omitempty"`
		// DeviceLinkFailDetails : has no documentation (yet)
		DeviceLinkFailDetails json.RawMessage `json:"device_link_fail_details,omitempty"`
		// DeviceLinkSuccessDetails : has no documentation (yet)
		DeviceLinkSuccessDetails json.RawMessage `json:"device_link_success_details,omitempty"`
		// DeviceManagementDisabledDetails : has no documentation (yet)
		DeviceManagementDisabledDetails json.RawMessage `json:"device_management_disabled_details,omitempty"`
		// DeviceManagementEnabledDetails : has no documentation (yet)
		DeviceManagementEnabledDetails json.RawMessage `json:"device_management_enabled_details,omitempty"`
		// DeviceUnlinkDetails : has no documentation (yet)
		DeviceUnlinkDetails json.RawMessage `json:"device_unlink_details,omitempty"`
		// EmmRefreshAuthTokenDetails : has no documentation (yet)
		EmmRefreshAuthTokenDetails json.RawMessage `json:"emm_refresh_auth_token_details,omitempty"`
		// AccountCaptureChangeAvailabilityDetails : has no documentation (yet)
		AccountCaptureChangeAvailabilityDetails json.RawMessage `json:"account_capture_change_availability_details,omitempty"`
		// AccountCaptureMigrateAccountDetails : has no documentation (yet)
		AccountCaptureMigrateAccountDetails json.RawMessage `json:"account_capture_migrate_account_details,omitempty"`
		// AccountCaptureNotificationEmailsSentDetails : has no documentation
		// (yet)
		AccountCaptureNotificationEmailsSentDetails json.RawMessage `json:"account_capture_notification_emails_sent_details,omitempty"`
		// AccountCaptureRelinquishAccountDetails : has no documentation (yet)
		AccountCaptureRelinquishAccountDetails json.RawMessage `json:"account_capture_relinquish_account_details,omitempty"`
		// DisabledDomainInvitesDetails : has no documentation (yet)
		DisabledDomainInvitesDetails json.RawMessage `json:"disabled_domain_invites_details,omitempty"`
		// DomainInvitesApproveRequestToJoinTeamDetails : has no documentation
		// (yet)
		DomainInvitesApproveRequestToJoinTeamDetails json.RawMessage `json:"domain_invites_approve_request_to_join_team_details,omitempty"`
		// DomainInvitesDeclineRequestToJoinTeamDetails : has no documentation
		// (yet)
		DomainInvitesDeclineRequestToJoinTeamDetails json.RawMessage `json:"domain_invites_decline_request_to_join_team_details,omitempty"`
		// DomainInvitesEmailExistingUsersDetails : has no documentation (yet)
		DomainInvitesEmailExistingUsersDetails json.RawMessage `json:"domain_invites_email_existing_users_details,omitempty"`
		// DomainInvitesRequestToJoinTeamDetails : has no documentation (yet)
		DomainInvitesRequestToJoinTeamDetails json.RawMessage `json:"domain_invites_request_to_join_team_details,omitempty"`
		// DomainInvitesSetInviteNewUserPrefToNoDetails : has no documentation
		// (yet)
		DomainInvitesSetInviteNewUserPrefToNoDetails json.RawMessage `json:"domain_invites_set_invite_new_user_pref_to_no_details,omitempty"`
		// DomainInvitesSetInviteNewUserPrefToYesDetails : has no documentation
		// (yet)
		DomainInvitesSetInviteNewUserPrefToYesDetails json.RawMessage `json:"domain_invites_set_invite_new_user_pref_to_yes_details,omitempty"`
		// DomainVerificationAddDomainFailDetails : has no documentation (yet)
		DomainVerificationAddDomainFailDetails json.RawMessage `json:"domain_verification_add_domain_fail_details,omitempty"`
		// DomainVerificationAddDomainSuccessDetails : has no documentation
		// (yet)
		DomainVerificationAddDomainSuccessDetails json.RawMessage `json:"domain_verification_add_domain_success_details,omitempty"`
		// DomainVerificationRemoveDomainDetails : has no documentation (yet)
		DomainVerificationRemoveDomainDetails json.RawMessage `json:"domain_verification_remove_domain_details,omitempty"`
		// EnabledDomainInvitesDetails : has no documentation (yet)
		EnabledDomainInvitesDetails json.RawMessage `json:"enabled_domain_invites_details,omitempty"`
		// CreateFolderDetails : has no documentation (yet)
		CreateFolderDetails json.RawMessage `json:"create_folder_details,omitempty"`
		// FileAddDetails : has no documentation (yet)
		FileAddDetails json.RawMessage `json:"file_add_details,omitempty"`
		// FileCopyDetails : has no documentation (yet)
		FileCopyDetails json.RawMessage `json:"file_copy_details,omitempty"`
		// FileDeleteDetails : has no documentation (yet)
		FileDeleteDetails json.RawMessage `json:"file_delete_details,omitempty"`
		// FileDownloadDetails : has no documentation (yet)
		FileDownloadDetails json.RawMessage `json:"file_download_details,omitempty"`
		// FileEditDetails : has no documentation (yet)
		FileEditDetails json.RawMessage `json:"file_edit_details,omitempty"`
		// FileGetCopyReferenceDetails : has no documentation (yet)
		FileGetCopyReferenceDetails json.RawMessage `json:"file_get_copy_reference_details,omitempty"`
		// FileMoveDetails : has no documentation (yet)
		FileMoveDetails json.RawMessage `json:"file_move_details,omitempty"`
		// FilePermanentlyDeleteDetails : has no documentation (yet)
		FilePermanentlyDeleteDetails json.RawMessage `json:"file_permanently_delete_details,omitempty"`
		// FilePreviewDetails : has no documentation (yet)
		FilePreviewDetails json.RawMessage `json:"file_preview_details,omitempty"`
		// FileRenameDetails : has no documentation (yet)
		FileRenameDetails json.RawMessage `json:"file_rename_details,omitempty"`
		// FileRestoreDetails : has no documentation (yet)
		FileRestoreDetails json.RawMessage `json:"file_restore_details,omitempty"`
		// FileRevertDetails : has no documentation (yet)
		FileRevertDetails json.RawMessage `json:"file_revert_details,omitempty"`
		// FileRollbackChangesDetails : has no documentation (yet)
		FileRollbackChangesDetails json.RawMessage `json:"file_rollback_changes_details,omitempty"`
		// FileSaveCopyReferenceDetails : has no documentation (yet)
		FileSaveCopyReferenceDetails json.RawMessage `json:"file_save_copy_reference_details,omitempty"`
		// FileRequestChangeDetails : has no documentation (yet)
		FileRequestChangeDetails json.RawMessage `json:"file_request_change_details,omitempty"`
		// FileRequestCloseDetails : has no documentation (yet)
		FileRequestCloseDetails json.RawMessage `json:"file_request_close_details,omitempty"`
		// FileRequestCreateDetails : has no documentation (yet)
		FileRequestCreateDetails json.RawMessage `json:"file_request_create_details,omitempty"`
		// FileRequestReceiveFileDetails : has no documentation (yet)
		FileRequestReceiveFileDetails json.RawMessage `json:"file_request_receive_file_details,omitempty"`
		// GroupAddExternalIdDetails : has no documentation (yet)
		GroupAddExternalIdDetails json.RawMessage `json:"group_add_external_id_details,omitempty"`
		// GroupAddMemberDetails : has no documentation (yet)
		GroupAddMemberDetails json.RawMessage `json:"group_add_member_details,omitempty"`
		// GroupChangeExternalIdDetails : has no documentation (yet)
		GroupChangeExternalIdDetails json.RawMessage `json:"group_change_external_id_details,omitempty"`
		// GroupChangeManagementTypeDetails : has no documentation (yet)
		GroupChangeManagementTypeDetails json.RawMessage `json:"group_change_management_type_details,omitempty"`
		// GroupChangeMemberRoleDetails : has no documentation (yet)
		GroupChangeMemberRoleDetails json.RawMessage `json:"group_change_member_role_details,omitempty"`
		// GroupCreateDetails : has no documentation (yet)
		GroupCreateDetails json.RawMessage `json:"group_create_details,omitempty"`
		// GroupDeleteDetails : has no documentation (yet)
		GroupDeleteDetails json.RawMessage `json:"group_delete_details,omitempty"`
		// GroupDescriptionUpdatedDetails : has no documentation (yet)
		GroupDescriptionUpdatedDetails json.RawMessage `json:"group_description_updated_details,omitempty"`
		// GroupJoinPolicyUpdatedDetails : has no documentation (yet)
		GroupJoinPolicyUpdatedDetails json.RawMessage `json:"group_join_policy_updated_details,omitempty"`
		// GroupMovedDetails : has no documentation (yet)
		GroupMovedDetails json.RawMessage `json:"group_moved_details,omitempty"`
		// GroupRemoveExternalIdDetails : has no documentation (yet)
		GroupRemoveExternalIdDetails json.RawMessage `json:"group_remove_external_id_details,omitempty"`
		// GroupRemoveMemberDetails : has no documentation (yet)
		GroupRemoveMemberDetails json.RawMessage `json:"group_remove_member_details,omitempty"`
		// GroupRenameDetails : has no documentation (yet)
		GroupRenameDetails json.RawMessage `json:"group_rename_details,omitempty"`
		// EmmErrorDetails : has no documentation (yet)
		EmmErrorDetails json.RawMessage `json:"emm_error_details,omitempty"`
		// LoginFailDetails : has no documentation (yet)
		LoginFailDetails json.RawMessage `json:"login_fail_details,omitempty"`
		// LoginSuccessDetails : has no documentation (yet)
		LoginSuccessDetails json.RawMessage `json:"login_success_details,omitempty"`
		// LogoutDetails : has no documentation (yet)
		LogoutDetails json.RawMessage `json:"logout_details,omitempty"`
		// ResellerSupportSessionEndDetails : has no documentation (yet)
		ResellerSupportSessionEndDetails json.RawMessage `json:"reseller_support_session_end_details,omitempty"`
		// ResellerSupportSessionStartDetails : has no documentation (yet)
		ResellerSupportSessionStartDetails json.RawMessage `json:"reseller_support_session_start_details,omitempty"`
		// SignInAsSessionEndDetails : has no documentation (yet)
		SignInAsSessionEndDetails json.RawMessage `json:"sign_in_as_session_end_details,omitempty"`
		// SignInAsSessionStartDetails : has no documentation (yet)
		SignInAsSessionStartDetails json.RawMessage `json:"sign_in_as_session_start_details,omitempty"`
		// SsoErrorDetails : has no documentation (yet)
		SsoErrorDetails json.RawMessage `json:"sso_error_details,omitempty"`
		// MemberAddNameDetails : has no documentation (yet)
		MemberAddNameDetails json.RawMessage `json:"member_add_name_details,omitempty"`
		// MemberChangeAdminRoleDetails : has no documentation (yet)
		MemberChangeAdminRoleDetails json.RawMessage `json:"member_change_admin_role_details,omitempty"`
		// MemberChangeEmailDetails : has no documentation (yet)
		MemberChangeEmailDetails json.RawMessage `json:"member_change_email_details,omitempty"`
		// MemberChangeMembershipTypeDetails : has no documentation (yet)
		MemberChangeMembershipTypeDetails json.RawMessage `json:"member_change_membership_type_details,omitempty"`
		// MemberChangeNameDetails : has no documentation (yet)
		MemberChangeNameDetails json.RawMessage `json:"member_change_name_details,omitempty"`
		// MemberChangeStatusDetails : has no documentation (yet)
		MemberChangeStatusDetails json.RawMessage `json:"member_change_status_details,omitempty"`
		// MemberDeleteManualContactsDetails : has no documentation (yet)
		MemberDeleteManualContactsDetails json.RawMessage `json:"member_delete_manual_contacts_details,omitempty"`
		// MemberPermanentlyDeleteAccountContentsDetails : has no documentation
		// (yet)
		MemberPermanentlyDeleteAccountContentsDetails json.RawMessage `json:"member_permanently_delete_account_contents_details,omitempty"`
		// MemberSpaceLimitsAddCustomQuotaDetails : has no documentation (yet)
		MemberSpaceLimitsAddCustomQuotaDetails json.RawMessage `json:"member_space_limits_add_custom_quota_details,omitempty"`
		// MemberSpaceLimitsChangeCustomQuotaDetails : has no documentation
		// (yet)
		MemberSpaceLimitsChangeCustomQuotaDetails json.RawMessage `json:"member_space_limits_change_custom_quota_details,omitempty"`
		// MemberSpaceLimitsChangeStatusDetails : has no documentation (yet)
		MemberSpaceLimitsChangeStatusDetails json.RawMessage `json:"member_space_limits_change_status_details,omitempty"`
		// MemberSpaceLimitsRemoveCustomQuotaDetails : has no documentation
		// (yet)
		MemberSpaceLimitsRemoveCustomQuotaDetails json.RawMessage `json:"member_space_limits_remove_custom_quota_details,omitempty"`
		// MemberSuggestDetails : has no documentation (yet)
		MemberSuggestDetails json.RawMessage `json:"member_suggest_details,omitempty"`
		// MemberTransferAccountContentsDetails : has no documentation (yet)
		MemberTransferAccountContentsDetails json.RawMessage `json:"member_transfer_account_contents_details,omitempty"`
		// SecondaryMailsPolicyChangedDetails : has no documentation (yet)
		SecondaryMailsPolicyChangedDetails json.RawMessage `json:"secondary_mails_policy_changed_details,omitempty"`
		// PaperContentAddMemberDetails : has no documentation (yet)
		PaperContentAddMemberDetails json.RawMessage `json:"paper_content_add_member_details,omitempty"`
		// PaperContentAddToFolderDetails : has no documentation (yet)
		PaperContentAddToFolderDetails json.RawMessage `json:"paper_content_add_to_folder_details,omitempty"`
		// PaperContentArchiveDetails : has no documentation (yet)
		PaperContentArchiveDetails json.RawMessage `json:"paper_content_archive_details,omitempty"`
		// PaperContentCreateDetails : has no documentation (yet)
		PaperContentCreateDetails json.RawMessage `json:"paper_content_create_details,omitempty"`
		// PaperContentPermanentlyDeleteDetails : has no documentation (yet)
		PaperContentPermanentlyDeleteDetails json.RawMessage `json:"paper_content_permanently_delete_details,omitempty"`
		// PaperContentRemoveFromFolderDetails : has no documentation (yet)
		PaperContentRemoveFromFolderDetails json.RawMessage `json:"paper_content_remove_from_folder_details,omitempty"`
		// PaperContentRemoveMemberDetails : has no documentation (yet)
		PaperContentRemoveMemberDetails json.RawMessage `json:"paper_content_remove_member_details,omitempty"`
		// PaperContentRenameDetails : has no documentation (yet)
		PaperContentRenameDetails json.RawMessage `json:"paper_content_rename_details,omitempty"`
		// PaperContentRestoreDetails : has no documentation (yet)
		PaperContentRestoreDetails json.RawMessage `json:"paper_content_restore_details,omitempty"`
		// PaperDocAddCommentDetails : has no documentation (yet)
		PaperDocAddCommentDetails json.RawMessage `json:"paper_doc_add_comment_details,omitempty"`
		// PaperDocChangeMemberRoleDetails : has no documentation (yet)
		PaperDocChangeMemberRoleDetails json.RawMessage `json:"paper_doc_change_member_role_details,omitempty"`
		// PaperDocChangeSharingPolicyDetails : has no documentation (yet)
		PaperDocChangeSharingPolicyDetails json.RawMessage `json:"paper_doc_change_sharing_policy_details,omitempty"`
		// PaperDocChangeSubscriptionDetails : has no documentation (yet)
		PaperDocChangeSubscriptionDetails json.RawMessage `json:"paper_doc_change_subscription_details,omitempty"`
		// PaperDocDeletedDetails : has no documentation (yet)
		PaperDocDeletedDetails json.RawMessage `json:"paper_doc_deleted_details,omitempty"`
		// PaperDocDeleteCommentDetails : has no documentation (yet)
		PaperDocDeleteCommentDetails json.RawMessage `json:"paper_doc_delete_comment_details,omitempty"`
		// PaperDocDownloadDetails : has no documentation (yet)
		PaperDocDownloadDetails json.RawMessage `json:"paper_doc_download_details,omitempty"`
		// PaperDocEditDetails : has no documentation (yet)
		PaperDocEditDetails json.RawMessage `json:"paper_doc_edit_details,omitempty"`
		// PaperDocEditCommentDetails : has no documentation (yet)
		PaperDocEditCommentDetails json.RawMessage `json:"paper_doc_edit_comment_details,omitempty"`
		// PaperDocFollowedDetails : has no documentation (yet)
		PaperDocFollowedDetails json.RawMessage `json:"paper_doc_followed_details,omitempty"`
		// PaperDocMentionDetails : has no documentation (yet)
		PaperDocMentionDetails json.RawMessage `json:"paper_doc_mention_details,omitempty"`
		// PaperDocOwnershipChangedDetails : has no documentation (yet)
		PaperDocOwnershipChangedDetails json.RawMessage `json:"paper_doc_ownership_changed_details,omitempty"`
		// PaperDocRequestAccessDetails : has no documentation (yet)
		PaperDocRequestAccessDetails json.RawMessage `json:"paper_doc_request_access_details,omitempty"`
		// PaperDocResolveCommentDetails : has no documentation (yet)
		PaperDocResolveCommentDetails json.RawMessage `json:"paper_doc_resolve_comment_details,omitempty"`
		// PaperDocRevertDetails : has no documentation (yet)
		PaperDocRevertDetails json.RawMessage `json:"paper_doc_revert_details,omitempty"`
		// PaperDocSlackShareDetails : has no documentation (yet)
		PaperDocSlackShareDetails json.RawMessage `json:"paper_doc_slack_share_details,omitempty"`
		// PaperDocTeamInviteDetails : has no documentation (yet)
		PaperDocTeamInviteDetails json.RawMessage `json:"paper_doc_team_invite_details,omitempty"`
		// PaperDocTrashedDetails : has no documentation (yet)
		PaperDocTrashedDetails json.RawMessage `json:"paper_doc_trashed_details,omitempty"`
		// PaperDocUnresolveCommentDetails : has no documentation (yet)
		PaperDocUnresolveCommentDetails json.RawMessage `json:"paper_doc_unresolve_comment_details,omitempty"`
		// PaperDocUntrashedDetails : has no documentation (yet)
		PaperDocUntrashedDetails json.RawMessage `json:"paper_doc_untrashed_details,omitempty"`
		// PaperDocViewDetails : has no documentation (yet)
		PaperDocViewDetails json.RawMessage `json:"paper_doc_view_details,omitempty"`
		// PaperExternalViewAllowDetails : has no documentation (yet)
		PaperExternalViewAllowDetails json.RawMessage `json:"paper_external_view_allow_details,omitempty"`
		// PaperExternalViewDefaultTeamDetails : has no documentation (yet)
		PaperExternalViewDefaultTeamDetails json.RawMessage `json:"paper_external_view_default_team_details,omitempty"`
		// PaperExternalViewForbidDetails : has no documentation (yet)
		PaperExternalViewForbidDetails json.RawMessage `json:"paper_external_view_forbid_details,omitempty"`
		// PaperFolderChangeSubscriptionDetails : has no documentation (yet)
		PaperFolderChangeSubscriptionDetails json.RawMessage `json:"paper_folder_change_subscription_details,omitempty"`
		// PaperFolderDeletedDetails : has no documentation (yet)
		PaperFolderDeletedDetails json.RawMessage `json:"paper_folder_deleted_details,omitempty"`
		// PaperFolderFollowedDetails : has no documentation (yet)
		PaperFolderFollowedDetails json.RawMessage `json:"paper_folder_followed_details,omitempty"`
		// PaperFolderTeamInviteDetails : has no documentation (yet)
		PaperFolderTeamInviteDetails json.RawMessage `json:"paper_folder_team_invite_details,omitempty"`
		// PasswordChangeDetails : has no documentation (yet)
		PasswordChangeDetails json.RawMessage `json:"password_change_details,omitempty"`
		// PasswordResetDetails : has no documentation (yet)
		PasswordResetDetails json.RawMessage `json:"password_reset_details,omitempty"`
		// PasswordResetAllDetails : has no documentation (yet)
		PasswordResetAllDetails json.RawMessage `json:"password_reset_all_details,omitempty"`
		// EmmCreateExceptionsReportDetails : has no documentation (yet)
		EmmCreateExceptionsReportDetails json.RawMessage `json:"emm_create_exceptions_report_details,omitempty"`
		// EmmCreateUsageReportDetails : has no documentation (yet)
		EmmCreateUsageReportDetails json.RawMessage `json:"emm_create_usage_report_details,omitempty"`
		// ExportMembersReportDetails : has no documentation (yet)
		ExportMembersReportDetails json.RawMessage `json:"export_members_report_details,omitempty"`
		// PaperAdminExportStartDetails : has no documentation (yet)
		PaperAdminExportStartDetails json.RawMessage `json:"paper_admin_export_start_details,omitempty"`
		// SmartSyncCreateAdminPrivilegeReportDetails : has no documentation
		// (yet)
		SmartSyncCreateAdminPrivilegeReportDetails json.RawMessage `json:"smart_sync_create_admin_privilege_report_details,omitempty"`
		// TeamActivityCreateReportDetails : has no documentation (yet)
		TeamActivityCreateReportDetails json.RawMessage `json:"team_activity_create_report_details,omitempty"`
		// CollectionShareDetails : has no documentation (yet)
		CollectionShareDetails json.RawMessage `json:"collection_share_details,omitempty"`
		// NoteAclInviteOnlyDetails : has no documentation (yet)
		NoteAclInviteOnlyDetails json.RawMessage `json:"note_acl_invite_only_details,omitempty"`
		// NoteAclLinkDetails : has no documentation (yet)
		NoteAclLinkDetails json.RawMessage `json:"note_acl_link_details,omitempty"`
		// NoteAclTeamLinkDetails : has no documentation (yet)
		NoteAclTeamLinkDetails json.RawMessage `json:"note_acl_team_link_details,omitempty"`
		// NoteSharedDetails : has no documentation (yet)
		NoteSharedDetails json.RawMessage `json:"note_shared_details,omitempty"`
		// NoteShareReceiveDetails : has no documentation (yet)
		NoteShareReceiveDetails json.RawMessage `json:"note_share_receive_details,omitempty"`
		// OpenNoteSharedDetails : has no documentation (yet)
		OpenNoteSharedDetails json.RawMessage `json:"open_note_shared_details,omitempty"`
		// SfAddGroupDetails : has no documentation (yet)
		SfAddGroupDetails json.RawMessage `json:"sf_add_group_details,omitempty"`
		// SfAllowNonMembersToViewSharedLinksDetails : has no documentation
		// (yet)
		SfAllowNonMembersToViewSharedLinksDetails json.RawMessage `json:"sf_allow_non_members_to_view_shared_links_details,omitempty"`
		// SfExternalInviteWarnDetails : has no documentation (yet)
		SfExternalInviteWarnDetails json.RawMessage `json:"sf_external_invite_warn_details,omitempty"`
		// SfFbInviteDetails : has no documentation (yet)
		SfFbInviteDetails json.RawMessage `json:"sf_fb_invite_details,omitempty"`
		// SfFbInviteChangeRoleDetails : has no documentation (yet)
		SfFbInviteChangeRoleDetails json.RawMessage `json:"sf_fb_invite_change_role_details,omitempty"`
		// SfFbUninviteDetails : has no documentation (yet)
		SfFbUninviteDetails json.RawMessage `json:"sf_fb_uninvite_details,omitempty"`
		// SfInviteGroupDetails : has no documentation (yet)
		SfInviteGroupDetails json.RawMessage `json:"sf_invite_group_details,omitempty"`
		// SfTeamGrantAccessDetails : has no documentation (yet)
		SfTeamGrantAccessDetails json.RawMessage `json:"sf_team_grant_access_details,omitempty"`
		// SfTeamInviteDetails : has no documentation (yet)
		SfTeamInviteDetails json.RawMessage `json:"sf_team_invite_details,omitempty"`
		// SfTeamInviteChangeRoleDetails : has no documentation (yet)
		SfTeamInviteChangeRoleDetails json.RawMessage `json:"sf_team_invite_change_role_details,omitempty"`
		// SfTeamJoinDetails : has no documentation (yet)
		SfTeamJoinDetails json.RawMessage `json:"sf_team_join_details,omitempty"`
		// SfTeamJoinFromOobLinkDetails : has no documentation (yet)
		SfTeamJoinFromOobLinkDetails json.RawMessage `json:"sf_team_join_from_oob_link_details,omitempty"`
		// SfTeamUninviteDetails : has no documentation (yet)
		SfTeamUninviteDetails json.RawMessage `json:"sf_team_uninvite_details,omitempty"`
		// SharedContentAddInviteesDetails : has no documentation (yet)
		SharedContentAddInviteesDetails json.RawMessage `json:"shared_content_add_invitees_details,omitempty"`
		// SharedContentAddLinkExpiryDetails : has no documentation (yet)
		SharedContentAddLinkExpiryDetails json.RawMessage `json:"shared_content_add_link_expiry_details,omitempty"`
		// SharedContentAddLinkPasswordDetails : has no documentation (yet)
		SharedContentAddLinkPasswordDetails json.RawMessage `json:"shared_content_add_link_password_details,omitempty"`
		// SharedContentAddMemberDetails : has no documentation (yet)
		SharedContentAddMemberDetails json.RawMessage `json:"shared_content_add_member_details,omitempty"`
		// SharedContentChangeDownloadsPolicyDetails : has no documentation
		// (yet)
		SharedContentChangeDownloadsPolicyDetails json.RawMessage `json:"shared_content_change_downloads_policy_details,omitempty"`
		// SharedContentChangeInviteeRoleDetails : has no documentation (yet)
		SharedContentChangeInviteeRoleDetails json.RawMessage `json:"shared_content_change_invitee_role_details,omitempty"`
		// SharedContentChangeLinkAudienceDetails : has no documentation (yet)
		SharedContentChangeLinkAudienceDetails json.RawMessage `json:"shared_content_change_link_audience_details,omitempty"`
		// SharedContentChangeLinkExpiryDetails : has no documentation (yet)
		SharedContentChangeLinkExpiryDetails json.RawMessage `json:"shared_content_change_link_expiry_details,omitempty"`
		// SharedContentChangeLinkPasswordDetails : has no documentation (yet)
		SharedContentChangeLinkPasswordDetails json.RawMessage `json:"shared_content_change_link_password_details,omitempty"`
		// SharedContentChangeMemberRoleDetails : has no documentation (yet)
		SharedContentChangeMemberRoleDetails json.RawMessage `json:"shared_content_change_member_role_details,omitempty"`
		// SharedContentChangeViewerInfoPolicyDetails : has no documentation
		// (yet)
		SharedContentChangeViewerInfoPolicyDetails json.RawMessage `json:"shared_content_change_viewer_info_policy_details,omitempty"`
		// SharedContentClaimInvitationDetails : has no documentation (yet)
		SharedContentClaimInvitationDetails json.RawMessage `json:"shared_content_claim_invitation_details,omitempty"`
		// SharedContentCopyDetails : has no documentation (yet)
		SharedContentCopyDetails json.RawMessage `json:"shared_content_copy_details,omitempty"`
		// SharedContentDownloadDetails : has no documentation (yet)
		SharedContentDownloadDetails json.RawMessage `json:"shared_content_download_details,omitempty"`
		// SharedContentRelinquishMembershipDetails : has no documentation (yet)
		SharedContentRelinquishMembershipDetails json.RawMessage `json:"shared_content_relinquish_membership_details,omitempty"`
		// SharedContentRemoveInviteesDetails : has no documentation (yet)
		SharedContentRemoveInviteesDetails json.RawMessage `json:"shared_content_remove_invitees_details,omitempty"`
		// SharedContentRemoveLinkExpiryDetails : has no documentation (yet)
		SharedContentRemoveLinkExpiryDetails json.RawMessage `json:"shared_content_remove_link_expiry_details,omitempty"`
		// SharedContentRemoveLinkPasswordDetails : has no documentation (yet)
		SharedContentRemoveLinkPasswordDetails json.RawMessage `json:"shared_content_remove_link_password_details,omitempty"`
		// SharedContentRemoveMemberDetails : has no documentation (yet)
		SharedContentRemoveMemberDetails json.RawMessage `json:"shared_content_remove_member_details,omitempty"`
		// SharedContentRequestAccessDetails : has no documentation (yet)
		SharedContentRequestAccessDetails json.RawMessage `json:"shared_content_request_access_details,omitempty"`
		// SharedContentUnshareDetails : has no documentation (yet)
		SharedContentUnshareDetails json.RawMessage `json:"shared_content_unshare_details,omitempty"`
		// SharedContentViewDetails : has no documentation (yet)
		SharedContentViewDetails json.RawMessage `json:"shared_content_view_details,omitempty"`
		// SharedFolderChangeLinkPolicyDetails : has no documentation (yet)
		SharedFolderChangeLinkPolicyDetails json.RawMessage `json:"shared_folder_change_link_policy_details,omitempty"`
		// SharedFolderChangeMembersInheritancePolicyDetails : has no
		// documentation (yet)
		SharedFolderChangeMembersInheritancePolicyDetails json.RawMessage `json:"shared_folder_change_members_inheritance_policy_details,omitempty"`
		// SharedFolderChangeMembersManagementPolicyDetails : has no
		// documentation (yet)
		SharedFolderChangeMembersManagementPolicyDetails json.RawMessage `json:"shared_folder_change_members_management_policy_details,omitempty"`
		// SharedFolderChangeMembersPolicyDetails : has no documentation (yet)
		SharedFolderChangeMembersPolicyDetails json.RawMessage `json:"shared_folder_change_members_policy_details,omitempty"`
		// SharedFolderCreateDetails : has no documentation (yet)
		SharedFolderCreateDetails json.RawMessage `json:"shared_folder_create_details,omitempty"`
		// SharedFolderDeclineInvitationDetails : has no documentation (yet)
		SharedFolderDeclineInvitationDetails json.RawMessage `json:"shared_folder_decline_invitation_details,omitempty"`
		// SharedFolderMountDetails : has no documentation (yet)
		SharedFolderMountDetails json.RawMessage `json:"shared_folder_mount_details,omitempty"`
		// SharedFolderNestDetails : has no documentation (yet)
		SharedFolderNestDetails json.RawMessage `json:"shared_folder_nest_details,omitempty"`
		// SharedFolderTransferOwnershipDetails : has no documentation (yet)
		SharedFolderTransferOwnershipDetails json.RawMessage `json:"shared_folder_transfer_ownership_details,omitempty"`
		// SharedFolderUnmountDetails : has no documentation (yet)
		SharedFolderUnmountDetails json.RawMessage `json:"shared_folder_unmount_details,omitempty"`
		// SharedLinkAddExpiryDetails : has no documentation (yet)
		SharedLinkAddExpiryDetails json.RawMessage `json:"shared_link_add_expiry_details,omitempty"`
		// SharedLinkChangeExpiryDetails : has no documentation (yet)
		SharedLinkChangeExpiryDetails json.RawMessage `json:"shared_link_change_expiry_details,omitempty"`
		// SharedLinkChangeVisibilityDetails : has no documentation (yet)
		SharedLinkChangeVisibilityDetails json.RawMessage `json:"shared_link_change_visibility_details,omitempty"`
		// SharedLinkCopyDetails : has no documentation (yet)
		SharedLinkCopyDetails json.RawMessage `json:"shared_link_copy_details,omitempty"`
		// SharedLinkCreateDetails : has no documentation (yet)
		SharedLinkCreateDetails json.RawMessage `json:"shared_link_create_details,omitempty"`
		// SharedLinkDisableDetails : has no documentation (yet)
		SharedLinkDisableDetails json.RawMessage `json:"shared_link_disable_details,omitempty"`
		// SharedLinkDownloadDetails : has no documentation (yet)
		SharedLinkDownloadDetails json.RawMessage `json:"shared_link_download_details,omitempty"`
		// SharedLinkRemoveExpiryDetails : has no documentation (yet)
		SharedLinkRemoveExpiryDetails json.RawMessage `json:"shared_link_remove_expiry_details,omitempty"`
		// SharedLinkShareDetails : has no documentation (yet)
		SharedLinkShareDetails json.RawMessage `json:"shared_link_share_details,omitempty"`
		// SharedLinkViewDetails : has no documentation (yet)
		SharedLinkViewDetails json.RawMessage `json:"shared_link_view_details,omitempty"`
		// SharedNoteOpenedDetails : has no documentation (yet)
		SharedNoteOpenedDetails json.RawMessage `json:"shared_note_opened_details,omitempty"`
		// ShmodelGroupShareDetails : has no documentation (yet)
		ShmodelGroupShareDetails json.RawMessage `json:"shmodel_group_share_details,omitempty"`
		// ShowcaseAccessGrantedDetails : has no documentation (yet)
		ShowcaseAccessGrantedDetails json.RawMessage `json:"showcase_access_granted_details,omitempty"`
		// ShowcaseAddMemberDetails : has no documentation (yet)
		ShowcaseAddMemberDetails json.RawMessage `json:"showcase_add_member_details,omitempty"`
		// ShowcaseArchivedDetails : has no documentation (yet)
		ShowcaseArchivedDetails json.RawMessage `json:"showcase_archived_details,omitempty"`
		// ShowcaseCreatedDetails : has no documentation (yet)
		ShowcaseCreatedDetails json.RawMessage `json:"showcase_created_details,omitempty"`
		// ShowcaseDeleteCommentDetails : has no documentation (yet)
		ShowcaseDeleteCommentDetails json.RawMessage `json:"showcase_delete_comment_details,omitempty"`
		// ShowcaseEditedDetails : has no documentation (yet)
		ShowcaseEditedDetails json.RawMessage `json:"showcase_edited_details,omitempty"`
		// ShowcaseEditCommentDetails : has no documentation (yet)
		ShowcaseEditCommentDetails json.RawMessage `json:"showcase_edit_comment_details,omitempty"`
		// ShowcaseFileAddedDetails : has no documentation (yet)
		ShowcaseFileAddedDetails json.RawMessage `json:"showcase_file_added_details,omitempty"`
		// ShowcaseFileDownloadDetails : has no documentation (yet)
		ShowcaseFileDownloadDetails json.RawMessage `json:"showcase_file_download_details,omitempty"`
		// ShowcaseFileRemovedDetails : has no documentation (yet)
		ShowcaseFileRemovedDetails json.RawMessage `json:"showcase_file_removed_details,omitempty"`
		// ShowcaseFileViewDetails : has no documentation (yet)
		ShowcaseFileViewDetails json.RawMessage `json:"showcase_file_view_details,omitempty"`
		// ShowcasePermanentlyDeletedDetails : has no documentation (yet)
		ShowcasePermanentlyDeletedDetails json.RawMessage `json:"showcase_permanently_deleted_details,omitempty"`
		// ShowcasePostCommentDetails : has no documentation (yet)
		ShowcasePostCommentDetails json.RawMessage `json:"showcase_post_comment_details,omitempty"`
		// ShowcaseRemoveMemberDetails : has no documentation (yet)
		ShowcaseRemoveMemberDetails json.RawMessage `json:"showcase_remove_member_details,omitempty"`
		// ShowcaseRenamedDetails : has no documentation (yet)
		ShowcaseRenamedDetails json.RawMessage `json:"showcase_renamed_details,omitempty"`
		// ShowcaseRequestAccessDetails : has no documentation (yet)
		ShowcaseRequestAccessDetails json.RawMessage `json:"showcase_request_access_details,omitempty"`
		// ShowcaseResolveCommentDetails : has no documentation (yet)
		ShowcaseResolveCommentDetails json.RawMessage `json:"showcase_resolve_comment_details,omitempty"`
		// ShowcaseRestoredDetails : has no documentation (yet)
		ShowcaseRestoredDetails json.RawMessage `json:"showcase_restored_details,omitempty"`
		// ShowcaseTrashedDetails : has no documentation (yet)
		ShowcaseTrashedDetails json.RawMessage `json:"showcase_trashed_details,omitempty"`
		// ShowcaseTrashedDeprecatedDetails : has no documentation (yet)
		ShowcaseTrashedDeprecatedDetails json.RawMessage `json:"showcase_trashed_deprecated_details,omitempty"`
		// ShowcaseUnresolveCommentDetails : has no documentation (yet)
		ShowcaseUnresolveCommentDetails json.RawMessage `json:"showcase_unresolve_comment_details,omitempty"`
		// ShowcaseUntrashedDetails : has no documentation (yet)
		ShowcaseUntrashedDetails json.RawMessage `json:"showcase_untrashed_details,omitempty"`
		// ShowcaseUntrashedDeprecatedDetails : has no documentation (yet)
		ShowcaseUntrashedDeprecatedDetails json.RawMessage `json:"showcase_untrashed_deprecated_details,omitempty"`
		// ShowcaseViewDetails : has no documentation (yet)
		ShowcaseViewDetails json.RawMessage `json:"showcase_view_details,omitempty"`
		// SsoAddCertDetails : has no documentation (yet)
		SsoAddCertDetails json.RawMessage `json:"sso_add_cert_details,omitempty"`
		// SsoAddLoginUrlDetails : has no documentation (yet)
		SsoAddLoginUrlDetails json.RawMessage `json:"sso_add_login_url_details,omitempty"`
		// SsoAddLogoutUrlDetails : has no documentation (yet)
		SsoAddLogoutUrlDetails json.RawMessage `json:"sso_add_logout_url_details,omitempty"`
		// SsoChangeCertDetails : has no documentation (yet)
		SsoChangeCertDetails json.RawMessage `json:"sso_change_cert_details,omitempty"`
		// SsoChangeLoginUrlDetails : has no documentation (yet)
		SsoChangeLoginUrlDetails json.RawMessage `json:"sso_change_login_url_details,omitempty"`
		// SsoChangeLogoutUrlDetails : has no documentation (yet)
		SsoChangeLogoutUrlDetails json.RawMessage `json:"sso_change_logout_url_details,omitempty"`
		// SsoChangeSamlIdentityModeDetails : has no documentation (yet)
		SsoChangeSamlIdentityModeDetails json.RawMessage `json:"sso_change_saml_identity_mode_details,omitempty"`
		// SsoRemoveCertDetails : has no documentation (yet)
		SsoRemoveCertDetails json.RawMessage `json:"sso_remove_cert_details,omitempty"`
		// SsoRemoveLoginUrlDetails : has no documentation (yet)
		SsoRemoveLoginUrlDetails json.RawMessage `json:"sso_remove_login_url_details,omitempty"`
		// SsoRemoveLogoutUrlDetails : has no documentation (yet)
		SsoRemoveLogoutUrlDetails json.RawMessage `json:"sso_remove_logout_url_details,omitempty"`
		// TeamFolderChangeStatusDetails : has no documentation (yet)
		TeamFolderChangeStatusDetails json.RawMessage `json:"team_folder_change_status_details,omitempty"`
		// TeamFolderCreateDetails : has no documentation (yet)
		TeamFolderCreateDetails json.RawMessage `json:"team_folder_create_details,omitempty"`
		// TeamFolderDowngradeDetails : has no documentation (yet)
		TeamFolderDowngradeDetails json.RawMessage `json:"team_folder_downgrade_details,omitempty"`
		// TeamFolderPermanentlyDeleteDetails : has no documentation (yet)
		TeamFolderPermanentlyDeleteDetails json.RawMessage `json:"team_folder_permanently_delete_details,omitempty"`
		// TeamFolderRenameDetails : has no documentation (yet)
		TeamFolderRenameDetails json.RawMessage `json:"team_folder_rename_details,omitempty"`
		// TeamSelectiveSyncSettingsChangedDetails : has no documentation (yet)
		TeamSelectiveSyncSettingsChangedDetails json.RawMessage `json:"team_selective_sync_settings_changed_details,omitempty"`
		// AccountCaptureChangePolicyDetails : has no documentation (yet)
		AccountCaptureChangePolicyDetails json.RawMessage `json:"account_capture_change_policy_details,omitempty"`
		// AllowDownloadDisabledDetails : has no documentation (yet)
		AllowDownloadDisabledDetails json.RawMessage `json:"allow_download_disabled_details,omitempty"`
		// AllowDownloadEnabledDetails : has no documentation (yet)
		AllowDownloadEnabledDetails json.RawMessage `json:"allow_download_enabled_details,omitempty"`
		// CameraUploadsPolicyChangedDetails : has no documentation (yet)
		CameraUploadsPolicyChangedDetails json.RawMessage `json:"camera_uploads_policy_changed_details,omitempty"`
		// DataPlacementRestrictionChangePolicyDetails : has no documentation
		// (yet)
		DataPlacementRestrictionChangePolicyDetails json.RawMessage `json:"data_placement_restriction_change_policy_details,omitempty"`
		// DataPlacementRestrictionSatisfyPolicyDetails : has no documentation
		// (yet)
		DataPlacementRestrictionSatisfyPolicyDetails json.RawMessage `json:"data_placement_restriction_satisfy_policy_details,omitempty"`
		// DeviceApprovalsChangeDesktopPolicyDetails : has no documentation
		// (yet)
		DeviceApprovalsChangeDesktopPolicyDetails json.RawMessage `json:"device_approvals_change_desktop_policy_details,omitempty"`
		// DeviceApprovalsChangeMobilePolicyDetails : has no documentation (yet)
		DeviceApprovalsChangeMobilePolicyDetails json.RawMessage `json:"device_approvals_change_mobile_policy_details,omitempty"`
		// DeviceApprovalsChangeOverageActionDetails : has no documentation
		// (yet)
		DeviceApprovalsChangeOverageActionDetails json.RawMessage `json:"device_approvals_change_overage_action_details,omitempty"`
		// DeviceApprovalsChangeUnlinkActionDetails : has no documentation (yet)
		DeviceApprovalsChangeUnlinkActionDetails json.RawMessage `json:"device_approvals_change_unlink_action_details,omitempty"`
		// DirectoryRestrictionsAddMembersDetails : has no documentation (yet)
		DirectoryRestrictionsAddMembersDetails json.RawMessage `json:"directory_restrictions_add_members_details,omitempty"`
		// DirectoryRestrictionsRemoveMembersDetails : has no documentation
		// (yet)
		DirectoryRestrictionsRemoveMembersDetails json.RawMessage `json:"directory_restrictions_remove_members_details,omitempty"`
		// EmmAddExceptionDetails : has no documentation (yet)
		EmmAddExceptionDetails json.RawMessage `json:"emm_add_exception_details,omitempty"`
		// EmmChangePolicyDetails : has no documentation (yet)
		EmmChangePolicyDetails json.RawMessage `json:"emm_change_policy_details,omitempty"`
		// EmmRemoveExceptionDetails : has no documentation (yet)
		EmmRemoveExceptionDetails json.RawMessage `json:"emm_remove_exception_details,omitempty"`
		// ExtendedVersionHistoryChangePolicyDetails : has no documentation
		// (yet)
		ExtendedVersionHistoryChangePolicyDetails json.RawMessage `json:"extended_version_history_change_policy_details,omitempty"`
		// FileCommentsChangePolicyDetails : has no documentation (yet)
		FileCommentsChangePolicyDetails json.RawMessage `json:"file_comments_change_policy_details,omitempty"`
		// FileRequestsChangePolicyDetails : has no documentation (yet)
		FileRequestsChangePolicyDetails json.RawMessage `json:"file_requests_change_policy_details,omitempty"`
		// FileRequestsEmailsEnabledDetails : has no documentation (yet)
		FileRequestsEmailsEnabledDetails json.RawMessage `json:"file_requests_emails_enabled_details,omitempty"`
		// FileRequestsEmailsRestrictedToTeamOnlyDetails : has no documentation
		// (yet)
		FileRequestsEmailsRestrictedToTeamOnlyDetails json.RawMessage `json:"file_requests_emails_restricted_to_team_only_details,omitempty"`
		// GoogleSsoChangePolicyDetails : has no documentation (yet)
		GoogleSsoChangePolicyDetails json.RawMessage `json:"google_sso_change_policy_details,omitempty"`
		// GroupUserManagementChangePolicyDetails : has no documentation (yet)
		GroupUserManagementChangePolicyDetails json.RawMessage `json:"group_user_management_change_policy_details,omitempty"`
		// MemberRequestsChangePolicyDetails : has no documentation (yet)
		MemberRequestsChangePolicyDetails json.RawMessage `json:"member_requests_change_policy_details,omitempty"`
		// MemberSpaceLimitsAddExceptionDetails : has no documentation (yet)
		MemberSpaceLimitsAddExceptionDetails json.RawMessage `json:"member_space_limits_add_exception_details,omitempty"`
		// MemberSpaceLimitsChangeCapsTypePolicyDetails : has no documentation
		// (yet)
		MemberSpaceLimitsChangeCapsTypePolicyDetails json.RawMessage `json:"member_space_limits_change_caps_type_policy_details,omitempty"`
		// MemberSpaceLimitsChangePolicyDetails : has no documentation (yet)
		MemberSpaceLimitsChangePolicyDetails json.RawMessage `json:"member_space_limits_change_policy_details,omitempty"`
		// MemberSpaceLimitsRemoveExceptionDetails : has no documentation (yet)
		MemberSpaceLimitsRemoveExceptionDetails json.RawMessage `json:"member_space_limits_remove_exception_details,omitempty"`
		// MemberSuggestionsChangePolicyDetails : has no documentation (yet)
		MemberSuggestionsChangePolicyDetails json.RawMessage `json:"member_suggestions_change_policy_details,omitempty"`
		// MicrosoftOfficeAddinChangePolicyDetails : has no documentation (yet)
		MicrosoftOfficeAddinChangePolicyDetails json.RawMessage `json:"microsoft_office_addin_change_policy_details,omitempty"`
		// NetworkControlChangePolicyDetails : has no documentation (yet)
		NetworkControlChangePolicyDetails json.RawMessage `json:"network_control_change_policy_details,omitempty"`
		// PaperChangeDeploymentPolicyDetails : has no documentation (yet)
		PaperChangeDeploymentPolicyDetails json.RawMessage `json:"paper_change_deployment_policy_details,omitempty"`
		// PaperChangeMemberLinkPolicyDetails : has no documentation (yet)
		PaperChangeMemberLinkPolicyDetails json.RawMessage `json:"paper_change_member_link_policy_details,omitempty"`
		// PaperChangeMemberPolicyDetails : has no documentation (yet)
		PaperChangeMemberPolicyDetails json.RawMessage `json:"paper_change_member_policy_details,omitempty"`
		// PaperChangePolicyDetails : has no documentation (yet)
		PaperChangePolicyDetails json.RawMessage `json:"paper_change_policy_details,omitempty"`
		// PaperEnabledUsersGroupAdditionDetails : has no documentation (yet)
		PaperEnabledUsersGroupAdditionDetails json.RawMessage `json:"paper_enabled_users_group_addition_details,omitempty"`
		// PaperEnabledUsersGroupRemovalDetails : has no documentation (yet)
		PaperEnabledUsersGroupRemovalDetails json.RawMessage `json:"paper_enabled_users_group_removal_details,omitempty"`
		// PermanentDeleteChangePolicyDetails : has no documentation (yet)
		PermanentDeleteChangePolicyDetails json.RawMessage `json:"permanent_delete_change_policy_details,omitempty"`
		// SharingChangeFolderJoinPolicyDetails : has no documentation (yet)
		SharingChangeFolderJoinPolicyDetails json.RawMessage `json:"sharing_change_folder_join_policy_details,omitempty"`
		// SharingChangeLinkPolicyDetails : has no documentation (yet)
		SharingChangeLinkPolicyDetails json.RawMessage `json:"sharing_change_link_policy_details,omitempty"`
		// SharingChangeMemberPolicyDetails : has no documentation (yet)
		SharingChangeMemberPolicyDetails json.RawMessage `json:"sharing_change_member_policy_details,omitempty"`
		// ShowcaseChangeDownloadPolicyDetails : has no documentation (yet)
		ShowcaseChangeDownloadPolicyDetails json.RawMessage `json:"showcase_change_download_policy_details,omitempty"`
		// ShowcaseChangeEnabledPolicyDetails : has no documentation (yet)
		ShowcaseChangeEnabledPolicyDetails json.RawMessage `json:"showcase_change_enabled_policy_details,omitempty"`
		// ShowcaseChangeExternalSharingPolicyDetails : has no documentation
		// (yet)
		ShowcaseChangeExternalSharingPolicyDetails json.RawMessage `json:"showcase_change_external_sharing_policy_details,omitempty"`
		// SmartSyncChangePolicyDetails : has no documentation (yet)
		SmartSyncChangePolicyDetails json.RawMessage `json:"smart_sync_change_policy_details,omitempty"`
		// SmartSyncNotOptOutDetails : has no documentation (yet)
		SmartSyncNotOptOutDetails json.RawMessage `json:"smart_sync_not_opt_out_details,omitempty"`
		// SmartSyncOptOutDetails : has no documentation (yet)
		SmartSyncOptOutDetails json.RawMessage `json:"smart_sync_opt_out_details,omitempty"`
		// SsoChangePolicyDetails : has no documentation (yet)
		SsoChangePolicyDetails json.RawMessage `json:"sso_change_policy_details,omitempty"`
		// TeamSelectiveSyncPolicyChangedDetails : has no documentation (yet)
		TeamSelectiveSyncPolicyChangedDetails json.RawMessage `json:"team_selective_sync_policy_changed_details,omitempty"`
		// TfaChangePolicyDetails : has no documentation (yet)
		TfaChangePolicyDetails json.RawMessage `json:"tfa_change_policy_details,omitempty"`
		// TwoAccountChangePolicyDetails : has no documentation (yet)
		TwoAccountChangePolicyDetails json.RawMessage `json:"two_account_change_policy_details,omitempty"`
		// ViewerInfoPolicyChangedDetails : has no documentation (yet)
		ViewerInfoPolicyChangedDetails json.RawMessage `json:"viewer_info_policy_changed_details,omitempty"`
		// WebSessionsChangeFixedLengthPolicyDetails : has no documentation
		// (yet)
		WebSessionsChangeFixedLengthPolicyDetails json.RawMessage `json:"web_sessions_change_fixed_length_policy_details,omitempty"`
		// WebSessionsChangeIdleLengthPolicyDetails : has no documentation (yet)
		WebSessionsChangeIdleLengthPolicyDetails json.RawMessage `json:"web_sessions_change_idle_length_policy_details,omitempty"`
		// TeamMergeFromDetails : has no documentation (yet)
		TeamMergeFromDetails json.RawMessage `json:"team_merge_from_details,omitempty"`
		// TeamMergeToDetails : has no documentation (yet)
		TeamMergeToDetails json.RawMessage `json:"team_merge_to_details,omitempty"`
		// TeamProfileAddLogoDetails : has no documentation (yet)
		TeamProfileAddLogoDetails json.RawMessage `json:"team_profile_add_logo_details,omitempty"`
		// TeamProfileChangeDefaultLanguageDetails : has no documentation (yet)
		TeamProfileChangeDefaultLanguageDetails json.RawMessage `json:"team_profile_change_default_language_details,omitempty"`
		// TeamProfileChangeLogoDetails : has no documentation (yet)
		TeamProfileChangeLogoDetails json.RawMessage `json:"team_profile_change_logo_details,omitempty"`
		// TeamProfileChangeNameDetails : has no documentation (yet)
		TeamProfileChangeNameDetails json.RawMessage `json:"team_profile_change_name_details,omitempty"`
		// TeamProfileRemoveLogoDetails : has no documentation (yet)
		TeamProfileRemoveLogoDetails json.RawMessage `json:"team_profile_remove_logo_details,omitempty"`
		// TfaAddBackupPhoneDetails : has no documentation (yet)
		TfaAddBackupPhoneDetails json.RawMessage `json:"tfa_add_backup_phone_details,omitempty"`
		// TfaAddSecurityKeyDetails : has no documentation (yet)
		TfaAddSecurityKeyDetails json.RawMessage `json:"tfa_add_security_key_details,omitempty"`
		// TfaChangeBackupPhoneDetails : has no documentation (yet)
		TfaChangeBackupPhoneDetails json.RawMessage `json:"tfa_change_backup_phone_details,omitempty"`
		// TfaChangeStatusDetails : has no documentation (yet)
		TfaChangeStatusDetails json.RawMessage `json:"tfa_change_status_details,omitempty"`
		// TfaRemoveBackupPhoneDetails : has no documentation (yet)
		TfaRemoveBackupPhoneDetails json.RawMessage `json:"tfa_remove_backup_phone_details,omitempty"`
		// TfaRemoveSecurityKeyDetails : has no documentation (yet)
		TfaRemoveSecurityKeyDetails json.RawMessage `json:"tfa_remove_security_key_details,omitempty"`
		// TfaResetDetails : has no documentation (yet)
		TfaResetDetails json.RawMessage `json:"tfa_reset_details,omitempty"`
		// MissingDetails : Hints that this event was returned with missing
		// details due to an internal error.
		MissingDetails json.RawMessage `json:"missing_details,omitempty"`
	}
	var w wrap
	var err error
	if err = json.Unmarshal(body, &w); err != nil {
		return err
	}
	u.Tag = w.Tag
	switch u.Tag {
	case "app_link_team_details":
		err = json.Unmarshal(body, &u.AppLinkTeamDetails)

		if err != nil {
			return err
		}
	case "app_link_user_details":
		err = json.Unmarshal(body, &u.AppLinkUserDetails)

		if err != nil {
			return err
		}
	case "app_unlink_team_details":
		err = json.Unmarshal(body, &u.AppUnlinkTeamDetails)

		if err != nil {
			return err
		}
	case "app_unlink_user_details":
		err = json.Unmarshal(body, &u.AppUnlinkUserDetails)

		if err != nil {
			return err
		}
	case "file_add_comment_details":
		err = json.Unmarshal(body, &u.FileAddCommentDetails)

		if err != nil {
			return err
		}
	case "file_change_comment_subscription_details":
		err = json.Unmarshal(body, &u.FileChangeCommentSubscriptionDetails)

		if err != nil {
			return err
		}
	case "file_delete_comment_details":
		err = json.Unmarshal(body, &u.FileDeleteCommentDetails)

		if err != nil {
			return err
		}
	case "file_edit_comment_details":
		err = json.Unmarshal(body, &u.FileEditCommentDetails)

		if err != nil {
			return err
		}
	case "file_like_comment_details":
		err = json.Unmarshal(body, &u.FileLikeCommentDetails)

		if err != nil {
			return err
		}
	case "file_resolve_comment_details":
		err = json.Unmarshal(body, &u.FileResolveCommentDetails)

		if err != nil {
			return err
		}
	case "file_unlike_comment_details":
		err = json.Unmarshal(body, &u.FileUnlikeCommentDetails)

		if err != nil {
			return err
		}
	case "file_unresolve_comment_details":
		err = json.Unmarshal(body, &u.FileUnresolveCommentDetails)

		if err != nil {
			return err
		}
	case "device_change_ip_desktop_details":
		err = json.Unmarshal(body, &u.DeviceChangeIpDesktopDetails)

		if err != nil {
			return err
		}
	case "device_change_ip_mobile_details":
		err = json.Unmarshal(body, &u.DeviceChangeIpMobileDetails)

		if err != nil {
			return err
		}
	case "device_change_ip_web_details":
		err = json.Unmarshal(body, &u.DeviceChangeIpWebDetails)

		if err != nil {
			return err
		}
	case "device_delete_on_unlink_fail_details":
		err = json.Unmarshal(body, &u.DeviceDeleteOnUnlinkFailDetails)

		if err != nil {
			return err
		}
	case "device_delete_on_unlink_success_details":
		err = json.Unmarshal(body, &u.DeviceDeleteOnUnlinkSuccessDetails)

		if err != nil {
			return err
		}
	case "device_link_fail_details":
		err = json.Unmarshal(body, &u.DeviceLinkFailDetails)

		if err != nil {
			return err
		}
	case "device_link_success_details":
		err = json.Unmarshal(body, &u.DeviceLinkSuccessDetails)

		if err != nil {
			return err
		}
	case "device_management_disabled_details":
		err = json.Unmarshal(body, &u.DeviceManagementDisabledDetails)

		if err != nil {
			return err
		}
	case "device_management_enabled_details":
		err = json.Unmarshal(body, &u.DeviceManagementEnabledDetails)

		if err != nil {
			return err
		}
	case "device_unlink_details":
		err = json.Unmarshal(body, &u.DeviceUnlinkDetails)

		if err != nil {
			return err
		}
	case "emm_refresh_auth_token_details":
		err = json.Unmarshal(body, &u.EmmRefreshAuthTokenDetails)

		if err != nil {
			return err
		}
	case "account_capture_change_availability_details":
		err = json.Unmarshal(body, &u.AccountCaptureChangeAvailabilityDetails)

		if err != nil {
			return err
		}
	case "account_capture_migrate_account_details":
		err = json.Unmarshal(body, &u.AccountCaptureMigrateAccountDetails)

		if err != nil {
			return err
		}
	case "account_capture_notification_emails_sent_details":
		err = json.Unmarshal(body, &u.AccountCaptureNotificationEmailsSentDetails)

		if err != nil {
			return err
		}
	case "account_capture_relinquish_account_details":
		err = json.Unmarshal(body, &u.AccountCaptureRelinquishAccountDetails)

		if err != nil {
			return err
		}
	case "disabled_domain_invites_details":
		err = json.Unmarshal(body, &u.DisabledDomainInvitesDetails)

		if err != nil {
			return err
		}
	case "domain_invites_approve_request_to_join_team_details":
		err = json.Unmarshal(body, &u.DomainInvitesApproveRequestToJoinTeamDetails)

		if err != nil {
			return err
		}
	case "domain_invites_decline_request_to_join_team_details":
		err = json.Unmarshal(body, &u.DomainInvitesDeclineRequestToJoinTeamDetails)

		if err != nil {
			return err
		}
	case "domain_invites_email_existing_users_details":
		err = json.Unmarshal(body, &u.DomainInvitesEmailExistingUsersDetails)

		if err != nil {
			return err
		}
	case "domain_invites_request_to_join_team_details":
		err = json.Unmarshal(body, &u.DomainInvitesRequestToJoinTeamDetails)

		if err != nil {
			return err
		}
	case "domain_invites_set_invite_new_user_pref_to_no_details":
		err = json.Unmarshal(body, &u.DomainInvitesSetInviteNewUserPrefToNoDetails)

		if err != nil {
			return err
		}
	case "domain_invites_set_invite_new_user_pref_to_yes_details":
		err = json.Unmarshal(body, &u.DomainInvitesSetInviteNewUserPrefToYesDetails)

		if err != nil {
			return err
		}
	case "domain_verification_add_domain_fail_details":
		err = json.Unmarshal(body, &u.DomainVerificationAddDomainFailDetails)

		if err != nil {
			return err
		}
	case "domain_verification_add_domain_success_details":
		err = json.Unmarshal(body, &u.DomainVerificationAddDomainSuccessDetails)

		if err != nil {
			return err
		}
	case "domain_verification_remove_domain_details":
		err = json.Unmarshal(body, &u.DomainVerificationRemoveDomainDetails)

		if err != nil {
			return err
		}
	case "enabled_domain_invites_details":
		err = json.Unmarshal(body, &u.EnabledDomainInvitesDetails)

		if err != nil {
			return err
		}
	case "create_folder_details":
		err = json.Unmarshal(body, &u.CreateFolderDetails)

		if err != nil {
			return err
		}
	case "file_add_details":
		err = json.Unmarshal(body, &u.FileAddDetails)

		if err != nil {
			return err
		}
	case "file_copy_details":
		err = json.Unmarshal(body, &u.FileCopyDetails)

		if err != nil {
			return err
		}
	case "file_delete_details":
		err = json.Unmarshal(body, &u.FileDeleteDetails)

		if err != nil {
			return err
		}
	case "file_download_details":
		err = json.Unmarshal(body, &u.FileDownloadDetails)

		if err != nil {
			return err
		}
	case "file_edit_details":
		err = json.Unmarshal(body, &u.FileEditDetails)

		if err != nil {
			return err
		}
	case "file_get_copy_reference_details":
		err = json.Unmarshal(body, &u.FileGetCopyReferenceDetails)

		if err != nil {
			return err
		}
	case "file_move_details":
		err = json.Unmarshal(body, &u.FileMoveDetails)

		if err != nil {
			return err
		}
	case "file_permanently_delete_details":
		err = json.Unmarshal(body, &u.FilePermanentlyDeleteDetails)

		if err != nil {
			return err
		}
	case "file_preview_details":
		err = json.Unmarshal(body, &u.FilePreviewDetails)

		if err != nil {
			return err
		}
	case "file_rename_details":
		err = json.Unmarshal(body, &u.FileRenameDetails)

		if err != nil {
			return err
		}
	case "file_restore_details":
		err = json.Unmarshal(body, &u.FileRestoreDetails)

		if err != nil {
			return err
		}
	case "file_revert_details":
		err = json.Unmarshal(body, &u.FileRevertDetails)

		if err != nil {
			return err
		}
	case "file_rollback_changes_details":
		err = json.Unmarshal(body, &u.FileRollbackChangesDetails)

		if err != nil {
			return err
		}
	case "file_save_copy_reference_details":
		err = json.Unmarshal(body, &u.FileSaveCopyReferenceDetails)

		if err != nil {
			return err
		}
	case "file_request_change_details":
		err = json.Unmarshal(body, &u.FileRequestChangeDetails)

		if err != nil {
			return err
		}
	case "file_request_close_details":
		err = json.Unmarshal(body, &u.FileRequestCloseDetails)

		if err != nil {
			return err
		}
	case "file_request_create_details":
		err = json.Unmarshal(body, &u.FileRequestCreateDetails)

		if err != nil {
			return err
		}
	case "file_request_receive_file_details":
		err = json.Unmarshal(body, &u.FileRequestReceiveFileDetails)

		if err != nil {
			return err
		}
	case "group_add_external_id_details":
		err = json.Unmarshal(body, &u.GroupAddExternalIdDetails)

		if err != nil {
			return err
		}
	case "group_add_member_details":
		err = json.Unmarshal(body, &u.GroupAddMemberDetails)

		if err != nil {
			return err
		}
	case "group_change_external_id_details":
		err = json.Unmarshal(body, &u.GroupChangeExternalIdDetails)

		if err != nil {
			return err
		}
	case "group_change_management_type_details":
		err = json.Unmarshal(body, &u.GroupChangeManagementTypeDetails)

		if err != nil {
			return err
		}
	case "group_change_member_role_details":
		err = json.Unmarshal(body, &u.GroupChangeMemberRoleDetails)

		if err != nil {
			return err
		}
	case "group_create_details":
		err = json.Unmarshal(body, &u.GroupCreateDetails)

		if err != nil {
			return err
		}
	case "group_delete_details":
		err = json.Unmarshal(body, &u.GroupDeleteDetails)

		if err != nil {
			return err
		}
	case "group_description_updated_details":
		err = json.Unmarshal(body, &u.GroupDescriptionUpdatedDetails)

		if err != nil {
			return err
		}
	case "group_join_policy_updated_details":
		err = json.Unmarshal(body, &u.GroupJoinPolicyUpdatedDetails)

		if err != nil {
			return err
		}
	case "group_moved_details":
		err = json.Unmarshal(body, &u.GroupMovedDetails)

		if err != nil {
			return err
		}
	case "group_remove_external_id_details":
		err = json.Unmarshal(body, &u.GroupRemoveExternalIdDetails)

		if err != nil {
			return err
		}
	case "group_remove_member_details":
		err = json.Unmarshal(body, &u.GroupRemoveMemberDetails)

		if err != nil {
			return err
		}
	case "group_rename_details":
		err = json.Unmarshal(body, &u.GroupRenameDetails)

		if err != nil {
			return err
		}
	case "emm_error_details":
		err = json.Unmarshal(body, &u.EmmErrorDetails)

		if err != nil {
			return err
		}
	case "login_fail_details":
		err = json.Unmarshal(body, &u.LoginFailDetails)

		if err != nil {
			return err
		}
	case "login_success_details":
		err = json.Unmarshal(body, &u.LoginSuccessDetails)

		if err != nil {
			return err
		}
	case "logout_details":
		err = json.Unmarshal(body, &u.LogoutDetails)

		if err != nil {
			return err
		}
	case "reseller_support_session_end_details":
		err = json.Unmarshal(body, &u.ResellerSupportSessionEndDetails)

		if err != nil {
			return err
		}
	case "reseller_support_session_start_details":
		err = json.Unmarshal(body, &u.ResellerSupportSessionStartDetails)

		if err != nil {
			return err
		}
	case "sign_in_as_session_end_details":
		err = json.Unmarshal(body, &u.SignInAsSessionEndDetails)

		if err != nil {
			return err
		}
	case "sign_in_as_session_start_details":
		err = json.Unmarshal(body, &u.SignInAsSessionStartDetails)

		if err != nil {
			return err
		}
	case "sso_error_details":
		err = json.Unmarshal(body, &u.SsoErrorDetails)

		if err != nil {
			return err
		}
	case "member_add_name_details":
		err = json.Unmarshal(body, &u.MemberAddNameDetails)

		if err != nil {
			return err
		}
	case "member_change_admin_role_details":
		err = json.Unmarshal(body, &u.MemberChangeAdminRoleDetails)

		if err != nil {
			return err
		}
	case "member_change_email_details":
		err = json.Unmarshal(body, &u.MemberChangeEmailDetails)

		if err != nil {
			return err
		}
	case "member_change_membership_type_details":
		err = json.Unmarshal(body, &u.MemberChangeMembershipTypeDetails)

		if err != nil {
			return err
		}
	case "member_change_name_details":
		err = json.Unmarshal(body, &u.MemberChangeNameDetails)

		if err != nil {
			return err
		}
	case "member_change_status_details":
		err = json.Unmarshal(body, &u.MemberChangeStatusDetails)

		if err != nil {
			return err
		}
	case "member_delete_manual_contacts_details":
		err = json.Unmarshal(body, &u.MemberDeleteManualContactsDetails)

		if err != nil {
			return err
		}
	case "member_permanently_delete_account_contents_details":
		err = json.Unmarshal(body, &u.MemberPermanentlyDeleteAccountContentsDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_add_custom_quota_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsAddCustomQuotaDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_change_custom_quota_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeCustomQuotaDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_change_status_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeStatusDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_remove_custom_quota_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsRemoveCustomQuotaDetails)

		if err != nil {
			return err
		}
	case "member_suggest_details":
		err = json.Unmarshal(body, &u.MemberSuggestDetails)

		if err != nil {
			return err
		}
	case "member_transfer_account_contents_details":
		err = json.Unmarshal(body, &u.MemberTransferAccountContentsDetails)

		if err != nil {
			return err
		}
	case "secondary_mails_policy_changed_details":
		err = json.Unmarshal(body, &u.SecondaryMailsPolicyChangedDetails)

		if err != nil {
			return err
		}
	case "paper_content_add_member_details":
		err = json.Unmarshal(body, &u.PaperContentAddMemberDetails)

		if err != nil {
			return err
		}
	case "paper_content_add_to_folder_details":
		err = json.Unmarshal(body, &u.PaperContentAddToFolderDetails)

		if err != nil {
			return err
		}
	case "paper_content_archive_details":
		err = json.Unmarshal(body, &u.PaperContentArchiveDetails)

		if err != nil {
			return err
		}
	case "paper_content_create_details":
		err = json.Unmarshal(body, &u.PaperContentCreateDetails)

		if err != nil {
			return err
		}
	case "paper_content_permanently_delete_details":
		err = json.Unmarshal(body, &u.PaperContentPermanentlyDeleteDetails)

		if err != nil {
			return err
		}
	case "paper_content_remove_from_folder_details":
		err = json.Unmarshal(body, &u.PaperContentRemoveFromFolderDetails)

		if err != nil {
			return err
		}
	case "paper_content_remove_member_details":
		err = json.Unmarshal(body, &u.PaperContentRemoveMemberDetails)

		if err != nil {
			return err
		}
	case "paper_content_rename_details":
		err = json.Unmarshal(body, &u.PaperContentRenameDetails)

		if err != nil {
			return err
		}
	case "paper_content_restore_details":
		err = json.Unmarshal(body, &u.PaperContentRestoreDetails)

		if err != nil {
			return err
		}
	case "paper_doc_add_comment_details":
		err = json.Unmarshal(body, &u.PaperDocAddCommentDetails)

		if err != nil {
			return err
		}
	case "paper_doc_change_member_role_details":
		err = json.Unmarshal(body, &u.PaperDocChangeMemberRoleDetails)

		if err != nil {
			return err
		}
	case "paper_doc_change_sharing_policy_details":
		err = json.Unmarshal(body, &u.PaperDocChangeSharingPolicyDetails)

		if err != nil {
			return err
		}
	case "paper_doc_change_subscription_details":
		err = json.Unmarshal(body, &u.PaperDocChangeSubscriptionDetails)

		if err != nil {
			return err
		}
	case "paper_doc_deleted_details":
		err = json.Unmarshal(body, &u.PaperDocDeletedDetails)

		if err != nil {
			return err
		}
	case "paper_doc_delete_comment_details":
		err = json.Unmarshal(body, &u.PaperDocDeleteCommentDetails)

		if err != nil {
			return err
		}
	case "paper_doc_download_details":
		err = json.Unmarshal(body, &u.PaperDocDownloadDetails)

		if err != nil {
			return err
		}
	case "paper_doc_edit_details":
		err = json.Unmarshal(body, &u.PaperDocEditDetails)

		if err != nil {
			return err
		}
	case "paper_doc_edit_comment_details":
		err = json.Unmarshal(body, &u.PaperDocEditCommentDetails)

		if err != nil {
			return err
		}
	case "paper_doc_followed_details":
		err = json.Unmarshal(body, &u.PaperDocFollowedDetails)

		if err != nil {
			return err
		}
	case "paper_doc_mention_details":
		err = json.Unmarshal(body, &u.PaperDocMentionDetails)

		if err != nil {
			return err
		}
	case "paper_doc_ownership_changed_details":
		err = json.Unmarshal(body, &u.PaperDocOwnershipChangedDetails)

		if err != nil {
			return err
		}
	case "paper_doc_request_access_details":
		err = json.Unmarshal(body, &u.PaperDocRequestAccessDetails)

		if err != nil {
			return err
		}
	case "paper_doc_resolve_comment_details":
		err = json.Unmarshal(body, &u.PaperDocResolveCommentDetails)

		if err != nil {
			return err
		}
	case "paper_doc_revert_details":
		err = json.Unmarshal(body, &u.PaperDocRevertDetails)

		if err != nil {
			return err
		}
	case "paper_doc_slack_share_details":
		err = json.Unmarshal(body, &u.PaperDocSlackShareDetails)

		if err != nil {
			return err
		}
	case "paper_doc_team_invite_details":
		err = json.Unmarshal(body, &u.PaperDocTeamInviteDetails)

		if err != nil {
			return err
		}
	case "paper_doc_trashed_details":
		err = json.Unmarshal(body, &u.PaperDocTrashedDetails)

		if err != nil {
			return err
		}
	case "paper_doc_unresolve_comment_details":
		err = json.Unmarshal(body, &u.PaperDocUnresolveCommentDetails)

		if err != nil {
			return err
		}
	case "paper_doc_untrashed_details":
		err = json.Unmarshal(body, &u.PaperDocUntrashedDetails)

		if err != nil {
			return err
		}
	case "paper_doc_view_details":
		err = json.Unmarshal(body, &u.PaperDocViewDetails)

		if err != nil {
			return err
		}
	case "paper_external_view_allow_details":
		err = json.Unmarshal(body, &u.PaperExternalViewAllowDetails)

		if err != nil {
			return err
		}
	case "paper_external_view_default_team_details":
		err = json.Unmarshal(body, &u.PaperExternalViewDefaultTeamDetails)

		if err != nil {
			return err
		}
	case "paper_external_view_forbid_details":
		err = json.Unmarshal(body, &u.PaperExternalViewForbidDetails)

		if err != nil {
			return err
		}
	case "paper_folder_change_subscription_details":
		err = json.Unmarshal(body, &u.PaperFolderChangeSubscriptionDetails)

		if err != nil {
			return err
		}
	case "paper_folder_deleted_details":
		err = json.Unmarshal(body, &u.PaperFolderDeletedDetails)

		if err != nil {
			return err
		}
	case "paper_folder_followed_details":
		err = json.Unmarshal(body, &u.PaperFolderFollowedDetails)

		if err != nil {
			return err
		}
	case "paper_folder_team_invite_details":
		err = json.Unmarshal(body, &u.PaperFolderTeamInviteDetails)

		if err != nil {
			return err
		}
	case "password_change_details":
		err = json.Unmarshal(body, &u.PasswordChangeDetails)

		if err != nil {
			return err
		}
	case "password_reset_details":
		err = json.Unmarshal(body, &u.PasswordResetDetails)

		if err != nil {
			return err
		}
	case "password_reset_all_details":
		err = json.Unmarshal(body, &u.PasswordResetAllDetails)

		if err != nil {
			return err
		}
	case "emm_create_exceptions_report_details":
		err = json.Unmarshal(body, &u.EmmCreateExceptionsReportDetails)

		if err != nil {
			return err
		}
	case "emm_create_usage_report_details":
		err = json.Unmarshal(body, &u.EmmCreateUsageReportDetails)

		if err != nil {
			return err
		}
	case "export_members_report_details":
		err = json.Unmarshal(body, &u.ExportMembersReportDetails)

		if err != nil {
			return err
		}
	case "paper_admin_export_start_details":
		err = json.Unmarshal(body, &u.PaperAdminExportStartDetails)

		if err != nil {
			return err
		}
	case "smart_sync_create_admin_privilege_report_details":
		err = json.Unmarshal(body, &u.SmartSyncCreateAdminPrivilegeReportDetails)

		if err != nil {
			return err
		}
	case "team_activity_create_report_details":
		err = json.Unmarshal(body, &u.TeamActivityCreateReportDetails)

		if err != nil {
			return err
		}
	case "collection_share_details":
		err = json.Unmarshal(body, &u.CollectionShareDetails)

		if err != nil {
			return err
		}
	case "note_acl_invite_only_details":
		err = json.Unmarshal(body, &u.NoteAclInviteOnlyDetails)

		if err != nil {
			return err
		}
	case "note_acl_link_details":
		err = json.Unmarshal(body, &u.NoteAclLinkDetails)

		if err != nil {
			return err
		}
	case "note_acl_team_link_details":
		err = json.Unmarshal(body, &u.NoteAclTeamLinkDetails)

		if err != nil {
			return err
		}
	case "note_shared_details":
		err = json.Unmarshal(body, &u.NoteSharedDetails)

		if err != nil {
			return err
		}
	case "note_share_receive_details":
		err = json.Unmarshal(body, &u.NoteShareReceiveDetails)

		if err != nil {
			return err
		}
	case "open_note_shared_details":
		err = json.Unmarshal(body, &u.OpenNoteSharedDetails)

		if err != nil {
			return err
		}
	case "sf_add_group_details":
		err = json.Unmarshal(body, &u.SfAddGroupDetails)

		if err != nil {
			return err
		}
	case "sf_allow_non_members_to_view_shared_links_details":
		err = json.Unmarshal(body, &u.SfAllowNonMembersToViewSharedLinksDetails)

		if err != nil {
			return err
		}
	case "sf_external_invite_warn_details":
		err = json.Unmarshal(body, &u.SfExternalInviteWarnDetails)

		if err != nil {
			return err
		}
	case "sf_fb_invite_details":
		err = json.Unmarshal(body, &u.SfFbInviteDetails)

		if err != nil {
			return err
		}
	case "sf_fb_invite_change_role_details":
		err = json.Unmarshal(body, &u.SfFbInviteChangeRoleDetails)

		if err != nil {
			return err
		}
	case "sf_fb_uninvite_details":
		err = json.Unmarshal(body, &u.SfFbUninviteDetails)

		if err != nil {
			return err
		}
	case "sf_invite_group_details":
		err = json.Unmarshal(body, &u.SfInviteGroupDetails)

		if err != nil {
			return err
		}
	case "sf_team_grant_access_details":
		err = json.Unmarshal(body, &u.SfTeamGrantAccessDetails)

		if err != nil {
			return err
		}
	case "sf_team_invite_details":
		err = json.Unmarshal(body, &u.SfTeamInviteDetails)

		if err != nil {
			return err
		}
	case "sf_team_invite_change_role_details":
		err = json.Unmarshal(body, &u.SfTeamInviteChangeRoleDetails)

		if err != nil {
			return err
		}
	case "sf_team_join_details":
		err = json.Unmarshal(body, &u.SfTeamJoinDetails)

		if err != nil {
			return err
		}
	case "sf_team_join_from_oob_link_details":
		err = json.Unmarshal(body, &u.SfTeamJoinFromOobLinkDetails)

		if err != nil {
			return err
		}
	case "sf_team_uninvite_details":
		err = json.Unmarshal(body, &u.SfTeamUninviteDetails)

		if err != nil {
			return err
		}
	case "shared_content_add_invitees_details":
		err = json.Unmarshal(body, &u.SharedContentAddInviteesDetails)

		if err != nil {
			return err
		}
	case "shared_content_add_link_expiry_details":
		err = json.Unmarshal(body, &u.SharedContentAddLinkExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_content_add_link_password_details":
		err = json.Unmarshal(body, &u.SharedContentAddLinkPasswordDetails)

		if err != nil {
			return err
		}
	case "shared_content_add_member_details":
		err = json.Unmarshal(body, &u.SharedContentAddMemberDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_downloads_policy_details":
		err = json.Unmarshal(body, &u.SharedContentChangeDownloadsPolicyDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_invitee_role_details":
		err = json.Unmarshal(body, &u.SharedContentChangeInviteeRoleDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_link_audience_details":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkAudienceDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_link_expiry_details":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_link_password_details":
		err = json.Unmarshal(body, &u.SharedContentChangeLinkPasswordDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_member_role_details":
		err = json.Unmarshal(body, &u.SharedContentChangeMemberRoleDetails)

		if err != nil {
			return err
		}
	case "shared_content_change_viewer_info_policy_details":
		err = json.Unmarshal(body, &u.SharedContentChangeViewerInfoPolicyDetails)

		if err != nil {
			return err
		}
	case "shared_content_claim_invitation_details":
		err = json.Unmarshal(body, &u.SharedContentClaimInvitationDetails)

		if err != nil {
			return err
		}
	case "shared_content_copy_details":
		err = json.Unmarshal(body, &u.SharedContentCopyDetails)

		if err != nil {
			return err
		}
	case "shared_content_download_details":
		err = json.Unmarshal(body, &u.SharedContentDownloadDetails)

		if err != nil {
			return err
		}
	case "shared_content_relinquish_membership_details":
		err = json.Unmarshal(body, &u.SharedContentRelinquishMembershipDetails)

		if err != nil {
			return err
		}
	case "shared_content_remove_invitees_details":
		err = json.Unmarshal(body, &u.SharedContentRemoveInviteesDetails)

		if err != nil {
			return err
		}
	case "shared_content_remove_link_expiry_details":
		err = json.Unmarshal(body, &u.SharedContentRemoveLinkExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_content_remove_link_password_details":
		err = json.Unmarshal(body, &u.SharedContentRemoveLinkPasswordDetails)

		if err != nil {
			return err
		}
	case "shared_content_remove_member_details":
		err = json.Unmarshal(body, &u.SharedContentRemoveMemberDetails)

		if err != nil {
			return err
		}
	case "shared_content_request_access_details":
		err = json.Unmarshal(body, &u.SharedContentRequestAccessDetails)

		if err != nil {
			return err
		}
	case "shared_content_unshare_details":
		err = json.Unmarshal(body, &u.SharedContentUnshareDetails)

		if err != nil {
			return err
		}
	case "shared_content_view_details":
		err = json.Unmarshal(body, &u.SharedContentViewDetails)

		if err != nil {
			return err
		}
	case "shared_folder_change_link_policy_details":
		err = json.Unmarshal(body, &u.SharedFolderChangeLinkPolicyDetails)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_inheritance_policy_details":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersInheritancePolicyDetails)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_management_policy_details":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersManagementPolicyDetails)

		if err != nil {
			return err
		}
	case "shared_folder_change_members_policy_details":
		err = json.Unmarshal(body, &u.SharedFolderChangeMembersPolicyDetails)

		if err != nil {
			return err
		}
	case "shared_folder_create_details":
		err = json.Unmarshal(body, &u.SharedFolderCreateDetails)

		if err != nil {
			return err
		}
	case "shared_folder_decline_invitation_details":
		err = json.Unmarshal(body, &u.SharedFolderDeclineInvitationDetails)

		if err != nil {
			return err
		}
	case "shared_folder_mount_details":
		err = json.Unmarshal(body, &u.SharedFolderMountDetails)

		if err != nil {
			return err
		}
	case "shared_folder_nest_details":
		err = json.Unmarshal(body, &u.SharedFolderNestDetails)

		if err != nil {
			return err
		}
	case "shared_folder_transfer_ownership_details":
		err = json.Unmarshal(body, &u.SharedFolderTransferOwnershipDetails)

		if err != nil {
			return err
		}
	case "shared_folder_unmount_details":
		err = json.Unmarshal(body, &u.SharedFolderUnmountDetails)

		if err != nil {
			return err
		}
	case "shared_link_add_expiry_details":
		err = json.Unmarshal(body, &u.SharedLinkAddExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_link_change_expiry_details":
		err = json.Unmarshal(body, &u.SharedLinkChangeExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_link_change_visibility_details":
		err = json.Unmarshal(body, &u.SharedLinkChangeVisibilityDetails)

		if err != nil {
			return err
		}
	case "shared_link_copy_details":
		err = json.Unmarshal(body, &u.SharedLinkCopyDetails)

		if err != nil {
			return err
		}
	case "shared_link_create_details":
		err = json.Unmarshal(body, &u.SharedLinkCreateDetails)

		if err != nil {
			return err
		}
	case "shared_link_disable_details":
		err = json.Unmarshal(body, &u.SharedLinkDisableDetails)

		if err != nil {
			return err
		}
	case "shared_link_download_details":
		err = json.Unmarshal(body, &u.SharedLinkDownloadDetails)

		if err != nil {
			return err
		}
	case "shared_link_remove_expiry_details":
		err = json.Unmarshal(body, &u.SharedLinkRemoveExpiryDetails)

		if err != nil {
			return err
		}
	case "shared_link_share_details":
		err = json.Unmarshal(body, &u.SharedLinkShareDetails)

		if err != nil {
			return err
		}
	case "shared_link_view_details":
		err = json.Unmarshal(body, &u.SharedLinkViewDetails)

		if err != nil {
			return err
		}
	case "shared_note_opened_details":
		err = json.Unmarshal(body, &u.SharedNoteOpenedDetails)

		if err != nil {
			return err
		}
	case "shmodel_group_share_details":
		err = json.Unmarshal(body, &u.ShmodelGroupShareDetails)

		if err != nil {
			return err
		}
	case "showcase_access_granted_details":
		err = json.Unmarshal(body, &u.ShowcaseAccessGrantedDetails)

		if err != nil {
			return err
		}
	case "showcase_add_member_details":
		err = json.Unmarshal(body, &u.ShowcaseAddMemberDetails)

		if err != nil {
			return err
		}
	case "showcase_archived_details":
		err = json.Unmarshal(body, &u.ShowcaseArchivedDetails)

		if err != nil {
			return err
		}
	case "showcase_created_details":
		err = json.Unmarshal(body, &u.ShowcaseCreatedDetails)

		if err != nil {
			return err
		}
	case "showcase_delete_comment_details":
		err = json.Unmarshal(body, &u.ShowcaseDeleteCommentDetails)

		if err != nil {
			return err
		}
	case "showcase_edited_details":
		err = json.Unmarshal(body, &u.ShowcaseEditedDetails)

		if err != nil {
			return err
		}
	case "showcase_edit_comment_details":
		err = json.Unmarshal(body, &u.ShowcaseEditCommentDetails)

		if err != nil {
			return err
		}
	case "showcase_file_added_details":
		err = json.Unmarshal(body, &u.ShowcaseFileAddedDetails)

		if err != nil {
			return err
		}
	case "showcase_file_download_details":
		err = json.Unmarshal(body, &u.ShowcaseFileDownloadDetails)

		if err != nil {
			return err
		}
	case "showcase_file_removed_details":
		err = json.Unmarshal(body, &u.ShowcaseFileRemovedDetails)

		if err != nil {
			return err
		}
	case "showcase_file_view_details":
		err = json.Unmarshal(body, &u.ShowcaseFileViewDetails)

		if err != nil {
			return err
		}
	case "showcase_permanently_deleted_details":
		err = json.Unmarshal(body, &u.ShowcasePermanentlyDeletedDetails)

		if err != nil {
			return err
		}
	case "showcase_post_comment_details":
		err = json.Unmarshal(body, &u.ShowcasePostCommentDetails)

		if err != nil {
			return err
		}
	case "showcase_remove_member_details":
		err = json.Unmarshal(body, &u.ShowcaseRemoveMemberDetails)

		if err != nil {
			return err
		}
	case "showcase_renamed_details":
		err = json.Unmarshal(body, &u.ShowcaseRenamedDetails)

		if err != nil {
			return err
		}
	case "showcase_request_access_details":
		err = json.Unmarshal(body, &u.ShowcaseRequestAccessDetails)

		if err != nil {
			return err
		}
	case "showcase_resolve_comment_details":
		err = json.Unmarshal(body, &u.ShowcaseResolveCommentDetails)

		if err != nil {
			return err
		}
	case "showcase_restored_details":
		err = json.Unmarshal(body, &u.ShowcaseRestoredDetails)

		if err != nil {
			return err
		}
	case "showcase_trashed_details":
		err = json.Unmarshal(body, &u.ShowcaseTrashedDetails)

		if err != nil {
			return err
		}
	case "showcase_trashed_deprecated_details":
		err = json.Unmarshal(body, &u.ShowcaseTrashedDeprecatedDetails)

		if err != nil {
			return err
		}
	case "showcase_unresolve_comment_details":
		err = json.Unmarshal(body, &u.ShowcaseUnresolveCommentDetails)

		if err != nil {
			return err
		}
	case "showcase_untrashed_details":
		err = json.Unmarshal(body, &u.ShowcaseUntrashedDetails)

		if err != nil {
			return err
		}
	case "showcase_untrashed_deprecated_details":
		err = json.Unmarshal(body, &u.ShowcaseUntrashedDeprecatedDetails)

		if err != nil {
			return err
		}
	case "showcase_view_details":
		err = json.Unmarshal(body, &u.ShowcaseViewDetails)

		if err != nil {
			return err
		}
	case "sso_add_cert_details":
		err = json.Unmarshal(body, &u.SsoAddCertDetails)

		if err != nil {
			return err
		}
	case "sso_add_login_url_details":
		err = json.Unmarshal(body, &u.SsoAddLoginUrlDetails)

		if err != nil {
			return err
		}
	case "sso_add_logout_url_details":
		err = json.Unmarshal(body, &u.SsoAddLogoutUrlDetails)

		if err != nil {
			return err
		}
	case "sso_change_cert_details":
		err = json.Unmarshal(body, &u.SsoChangeCertDetails)

		if err != nil {
			return err
		}
	case "sso_change_login_url_details":
		err = json.Unmarshal(body, &u.SsoChangeLoginUrlDetails)

		if err != nil {
			return err
		}
	case "sso_change_logout_url_details":
		err = json.Unmarshal(body, &u.SsoChangeLogoutUrlDetails)

		if err != nil {
			return err
		}
	case "sso_change_saml_identity_mode_details":
		err = json.Unmarshal(body, &u.SsoChangeSamlIdentityModeDetails)

		if err != nil {
			return err
		}
	case "sso_remove_cert_details":
		err = json.Unmarshal(body, &u.SsoRemoveCertDetails)

		if err != nil {
			return err
		}
	case "sso_remove_login_url_details":
		err = json.Unmarshal(body, &u.SsoRemoveLoginUrlDetails)

		if err != nil {
			return err
		}
	case "sso_remove_logout_url_details":
		err = json.Unmarshal(body, &u.SsoRemoveLogoutUrlDetails)

		if err != nil {
			return err
		}
	case "team_folder_change_status_details":
		err = json.Unmarshal(body, &u.TeamFolderChangeStatusDetails)

		if err != nil {
			return err
		}
	case "team_folder_create_details":
		err = json.Unmarshal(body, &u.TeamFolderCreateDetails)

		if err != nil {
			return err
		}
	case "team_folder_downgrade_details":
		err = json.Unmarshal(body, &u.TeamFolderDowngradeDetails)

		if err != nil {
			return err
		}
	case "team_folder_permanently_delete_details":
		err = json.Unmarshal(body, &u.TeamFolderPermanentlyDeleteDetails)

		if err != nil {
			return err
		}
	case "team_folder_rename_details":
		err = json.Unmarshal(body, &u.TeamFolderRenameDetails)

		if err != nil {
			return err
		}
	case "team_selective_sync_settings_changed_details":
		err = json.Unmarshal(body, &u.TeamSelectiveSyncSettingsChangedDetails)

		if err != nil {
			return err
		}
	case "account_capture_change_policy_details":
		err = json.Unmarshal(body, &u.AccountCaptureChangePolicyDetails)

		if err != nil {
			return err
		}
	case "allow_download_disabled_details":
		err = json.Unmarshal(body, &u.AllowDownloadDisabledDetails)

		if err != nil {
			return err
		}
	case "allow_download_enabled_details":
		err = json.Unmarshal(body, &u.AllowDownloadEnabledDetails)

		if err != nil {
			return err
		}
	case "camera_uploads_policy_changed_details":
		err = json.Unmarshal(body, &u.CameraUploadsPolicyChangedDetails)

		if err != nil {
			return err
		}
	case "data_placement_restriction_change_policy_details":
		err = json.Unmarshal(body, &u.DataPlacementRestrictionChangePolicyDetails)

		if err != nil {
			return err
		}
	case "data_placement_restriction_satisfy_policy_details":
		err = json.Unmarshal(body, &u.DataPlacementRestrictionSatisfyPolicyDetails)

		if err != nil {
			return err
		}
	case "device_approvals_change_desktop_policy_details":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeDesktopPolicyDetails)

		if err != nil {
			return err
		}
	case "device_approvals_change_mobile_policy_details":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeMobilePolicyDetails)

		if err != nil {
			return err
		}
	case "device_approvals_change_overage_action_details":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeOverageActionDetails)

		if err != nil {
			return err
		}
	case "device_approvals_change_unlink_action_details":
		err = json.Unmarshal(body, &u.DeviceApprovalsChangeUnlinkActionDetails)

		if err != nil {
			return err
		}
	case "directory_restrictions_add_members_details":
		err = json.Unmarshal(body, &u.DirectoryRestrictionsAddMembersDetails)

		if err != nil {
			return err
		}
	case "directory_restrictions_remove_members_details":
		err = json.Unmarshal(body, &u.DirectoryRestrictionsRemoveMembersDetails)

		if err != nil {
			return err
		}
	case "emm_add_exception_details":
		err = json.Unmarshal(body, &u.EmmAddExceptionDetails)

		if err != nil {
			return err
		}
	case "emm_change_policy_details":
		err = json.Unmarshal(body, &u.EmmChangePolicyDetails)

		if err != nil {
			return err
		}
	case "emm_remove_exception_details":
		err = json.Unmarshal(body, &u.EmmRemoveExceptionDetails)

		if err != nil {
			return err
		}
	case "extended_version_history_change_policy_details":
		err = json.Unmarshal(body, &u.ExtendedVersionHistoryChangePolicyDetails)

		if err != nil {
			return err
		}
	case "file_comments_change_policy_details":
		err = json.Unmarshal(body, &u.FileCommentsChangePolicyDetails)

		if err != nil {
			return err
		}
	case "file_requests_change_policy_details":
		err = json.Unmarshal(body, &u.FileRequestsChangePolicyDetails)

		if err != nil {
			return err
		}
	case "file_requests_emails_enabled_details":
		err = json.Unmarshal(body, &u.FileRequestsEmailsEnabledDetails)

		if err != nil {
			return err
		}
	case "file_requests_emails_restricted_to_team_only_details":
		err = json.Unmarshal(body, &u.FileRequestsEmailsRestrictedToTeamOnlyDetails)

		if err != nil {
			return err
		}
	case "google_sso_change_policy_details":
		err = json.Unmarshal(body, &u.GoogleSsoChangePolicyDetails)

		if err != nil {
			return err
		}
	case "group_user_management_change_policy_details":
		err = json.Unmarshal(body, &u.GroupUserManagementChangePolicyDetails)

		if err != nil {
			return err
		}
	case "member_requests_change_policy_details":
		err = json.Unmarshal(body, &u.MemberRequestsChangePolicyDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_add_exception_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsAddExceptionDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_change_caps_type_policy_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangeCapsTypePolicyDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_change_policy_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsChangePolicyDetails)

		if err != nil {
			return err
		}
	case "member_space_limits_remove_exception_details":
		err = json.Unmarshal(body, &u.MemberSpaceLimitsRemoveExceptionDetails)

		if err != nil {
			return err
		}
	case "member_suggestions_change_policy_details":
		err = json.Unmarshal(body, &u.MemberSuggestionsChangePolicyDetails)

		if err != nil {
			return err
		}
	case "microsoft_office_addin_change_policy_details":
		err = json.Unmarshal(body, &u.MicrosoftOfficeAddinChangePolicyDetails)

		if err != nil {
			return err
		}
	case "network_control_change_policy_details":
		err = json.Unmarshal(body, &u.NetworkControlChangePolicyDetails)

		if err != nil {
			return err
		}
	case "paper_change_deployment_policy_details":
		err = json.Unmarshal(body, &u.PaperChangeDeploymentPolicyDetails)

		if err != nil {
			return err
		}
	case "paper_change_member_link_policy_details":
		err = json.Unmarshal(body, &u.PaperChangeMemberLinkPolicyDetails)

		if err != nil {
			return err
		}
	case "paper_change_member_policy_details":
		err = json.Unmarshal(body, &u.PaperChangeMemberPolicyDetails)

		if err != nil {
			return err
		}
	case "paper_change_policy_details":
		err = json.Unmarshal(body, &u.PaperChangePolicyDetails)

		if err != nil {
			return err
		}
	case "paper_enabled_users_group_addition_details":
		err = json.Unmarshal(body, &u.PaperEnabledUsersGroupAdditionDetails)

		if err != nil {
			return err
		}
	case "paper_enabled_users_group_removal_details":
		err = json.Unmarshal(body, &u.PaperEnabledUsersGroupRemovalDetails)

		if err != nil {
			return err
		}
	case "permanent_delete_change_policy_details":
		err = json.Unmarshal(body, &u.PermanentDeleteChangePolicyDetails)

		if err != nil {
			return err
		}
	case "sharing_change_folder_join_policy_details":
		err = json.Unmarshal(body, &u.SharingChangeFolderJoinPolicyDetails)

		if err != nil {
			return err
		}
	case "sharing_change_link_policy_details":
		err = json.Unmarshal(body, &u.SharingChangeLinkPolicyDetails)

		if err != nil {
			return err
		}
	case "sharing_change_member_policy_details":
		err = json.Unmarshal(body, &u.SharingChangeMemberPolicyDetails)

		if err != nil {
			return err
		}
	case "showcase_change_download_policy_details":
		err = json.Unmarshal(body, &u.ShowcaseChangeDownloadPolicyDetails)

		if err != nil {
			return err
		}
	case "showcase_change_enabled_policy_details":
		err = json.Unmarshal(body, &u.ShowcaseChangeEnabledPolicyDetails)

		if err != nil {
			return err
		}
	case "showcase_change_external_sharing_policy_details":
		err = json.Unmarshal(body, &u.ShowcaseChangeExternalSharingPolicyDetails)

		if err != nil {
			return err
		}
	case "smart_sync_change_policy_details":
		err = json.Unmarshal(body, &u.SmartSyncChangePolicyDetails)

		if err != nil {
			return err
		}
	case "smart_sync_not_opt_out_details":
		err = json.Unmarshal(body, &u.SmartSyncNotOptOutDetails)

		if err != nil {
			return err
		}
	case "smart_sync_opt_out_details":
		err = json.Unmarshal(body, &u.SmartSyncOptOutDetails)

		if err != nil {
			return err
		}
	case "sso_change_policy_details":
		err = json.Unmarshal(body, &u.SsoChangePolicyDetails)

		if err != nil {
			return err
		}
	case "team_selective_sync_policy_changed_details":
		err = json.Unmarshal(body, &u.TeamSelectiveSyncPolicyChangedDetails)

		if err != nil {
			return err
		}
	case "tfa_change_policy_details":
		err = json.Unmarshal(body, &u.TfaChangePolicyDetails)

		if err != nil {
			return err
		}
	case "two_account_change_policy_details":
		err = json.Unmarshal(body, &u.TwoAccountChangePolicyDetails)

		if err != nil {
			return err
		}
	case "viewer_info_policy_changed_details":
		err = json.Unmarshal(body, &u.ViewerInfoPolicyChangedDetails)

		if err != nil {
			return err
		}
	case "web_sessions_change_fixed_length_policy_details":
		err = json.Unmarshal(body, &u.WebSessionsChangeFixedLengthPolicyDetails)

		if err != nil {
			return err
		}
	case "web_sessions_change_idle_length_policy_details":
		err = json.Unmarshal(body, &u.WebSessionsChangeIdleLengthPolicyDetails)

		if err != nil {
			return err
		}
	case "team_merge_from_details":
		err = json.Unmarshal(body, &u.TeamMergeFromDetails)

		if err != nil {
			return err
		}
	case "team_merge_to_details":
		err = json.Unmarshal(body, &u.TeamMergeToDetails)

		if err != nil {
			return err
		}
	case "team_profile_add_logo_details":
		err = json.Unmarshal(body, &u.TeamProfileAddLogoDetails)

		if err != nil {
			return err
		}
	case "team_profile_change_default_language_details":
		err = json.Unmarshal(body, &u.TeamProfileChangeDefaultLanguageDetails)

		if err != nil {
			return err
		}
	case "team_profile_change_logo_details":
		err = json.Unmarshal(body, &u.TeamProfileChangeLogoDetails)

		if err != nil {
			return err
		}
	case "team_profile_change_name_details":
		err = json.Unmarshal(body, &u.TeamProfileChangeNameDetails)

		if err != nil {
			return err
		}
	case "team_profile_remove_logo_details":
		err = json.Unmarshal(body, &u.TeamProfileRemoveLogoDetails)

		if err != nil {
			return err
		}
	case "tfa_add_backup_phone_details":
		err = json.Unmarshal(body, &u.TfaAddBackupPhoneDetails)

		if err != nil {
			return err
		}
	case "tfa_add_security_key_details":
		err = json.Unmarshal(body, &u.TfaAddSecurityKeyDetails)

		if err != nil {
			return err
		}
	case "tfa_change_backup_phone_details":
		err = json.Unmarshal(body, &u.TfaChangeBackupPhoneDetails)

		if err != nil {
			return err
		}
	case "tfa_change_status_details":
		err = json.Unmarshal(body, &u.TfaChangeStatusDetails)

		if err != nil {
			return err
		}
	case "tfa_remove_backup_phone_details":
		err = json.Unmarshal(body, &u.TfaRemoveBackupPhoneDetails)

		if err != nil {
			return err
		}
	case "tfa_remove_security_key_details":
		err = json.Unmarshal(body, &u.TfaRemoveSecurityKeyDetails)

		if err != nil {
			return err
		}
	case "tfa_reset_details":
		err = json.Unmarshal(body, &u.TfaResetDetails)

		if err != nil {
			return err
		}
	case "missing_details":
		err = json.Unmarshal(body, &u.MissingDetails)

		if err != nil {
			return err
		}
	}
	return nil
}
