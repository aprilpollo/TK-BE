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
		&ProjectModel{},
		&ProjectStatusModel{},
		&ProjectNotificationSettingModel{},
		&TasksModel{},
		&TaskAssignModel{},
		&TaskCommentModel{},
		&TaskAttachmentModel{},
		&TaskStatusModel{},
		&TaskPriorityModel{},
	}
}