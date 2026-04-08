package models

type TableNamer interface {
	TableName() string
}

type ModelList []TableNamer

func All() ModelList {
	return ModelList{
		&UserModel{},
		&OauthModel{},
		&OrganizationModel{},
		&OrganizationMemberModel{},
		&OrganizationMemberStatusModel{},
		&OrganizationMemberRoleModel{},
		&OrganizationMemberPagePermissionModel{},
		&TasksModel{},
		&TaskAssignModel{},
		&TaskCommentModel{},
		&TaskAttachmentModel{},
		&TaskStatusModel{},
		&TaskPriorityModel{},
	}
}