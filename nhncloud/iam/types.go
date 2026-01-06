package iam

type Organization struct {
	ID          string `json:"orgId"`
	Name        string `json:"orgName"`
	Description string `json:"description,omitempty"`
	Status      string `json:"orgStatus"`
	CreatedAt   string `json:"createdDateTime,omitempty"`
	UpdatedAt   string `json:"modifiedDateTime,omitempty"`
}

type Project struct {
	ID             string `json:"projectId"`
	Name           string `json:"projectName"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"orgId"`
	Status         string `json:"projectStatusCode"`
	CreatedAt      string `json:"createdDateTime,omitempty"`
	UpdatedAt      string `json:"modifiedDateTime,omitempty"`
}

type Member struct {
	ID        string   `json:"uuid"`
	Email     string   `json:"emailAddress"`
	Name      string   `json:"memberName"`
	Status    string   `json:"memberStatus"`
	Roles     []string `json:"roles,omitempty"`
	CreatedAt string   `json:"createdDateTime,omitempty"`
	UpdatedAt string   `json:"modifiedDateTime,omitempty"`
}

type ListOrganizationsOutput struct {
	Organizations []Organization `json:"orgList"`
}

type GetOrganizationOutput struct {
	Organization Organization `json:"org"`
}

type ListProjectsOutput struct {
	Projects []Project `json:"projectList"`
}

type GetProjectOutput struct {
	Project Project `json:"project"`
}

type ListMembersOutput struct {
	Members []Member `json:"memberList"`
}

type GetMemberOutput struct {
	Member Member `json:"member"`
}

type InviteMemberInput struct {
	Email string   `json:"emailAddress"`
	Roles []string `json:"roles,omitempty"`
}

type InviteMemberOutput struct {
	MemberID string `json:"uuid"`
}
