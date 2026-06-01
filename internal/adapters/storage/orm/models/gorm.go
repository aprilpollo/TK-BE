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
		&SubTasksModel{},
		&SubTaskAssignModel{},
		&TaskAssignModel{},
		&TaskCommentModel{},
		&TaskCommentFileModel{},
		&TaskAttachmentModel{},
		&TaskStatusModel{},
		&TaskPriorityModel{},
	}
}
