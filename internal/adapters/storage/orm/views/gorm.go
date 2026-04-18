package views

var Views = map[string]string{
	"vw_organizations_members": `
CREATE OR REPLACE VIEW vw_organizations_members AS
SELECT 
	om.* ,
	u.first_name ,
	u.last_name ,
	u.display_name ,
	u.email ,
	u.avatar 
FROM organization_members om
LEFT JOIN users u ON om.user_id = u.id 
`,
}