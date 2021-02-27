package models

// Invitation describes a Invitation object.
type Invitation struct {
	ID                      *string `json:"id,omitempty"`
	InvitedUserDisplayName  *string `json:"invitedUserDisplayName,omitempty"`
	InvitedUserEmailAddress *string `json:"invitedUserEmailAddress,omitempty"`
	SendInvitationMessage   *bool   `json:"sendInvitationMessage,omitempty"`
	InviteRedirectURL       *string `json:"inviteRedirectUrl,omitempty"`
	InviteRedeemURL         *string `json:"inviteRedeemUrl,omitempty"`
	Status                  *string `json:"status,omitempty"`
	InvitedUserType         *string `json:"invitedUserType,omitempty"`

	InvitedUserMessageInfo *InvitedUserMessageInfo `json:"invitedUserMessageInfo,omitempty"`
	InvitedUser            *User                   `json:"invitedUser,omitempty"`
}

type InvitedUserMessageInfo struct {
	CCRecipients          *[]Recipient `json:"ccRecipients,omitempty"`
	CustomizedMessageBody *string      `json:"customizedMessageBody,omitempty"`
	MessageLanguage       *string      `json:"messageLanguage,omitempty"`
}

type Recipient struct {
	EmailAddress *EmailAddress `json:"emailAddress,omitempty"`
}

type EmailAddress struct {
	Address *string `json:"address,omitempty"`
	Name    *string `json:"name,omitempty"`
}
