package utils

var (
	IdVerify               = Rules{"ID": {NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	DeptVerify             = Rules{"DeptSort": {NotEmpty()}, "DeptName": {NotEmpty()}, "ParentID": {NotEmpty()}}
	LoginVerify            = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityIds": {NotEmpty()}, "DeptId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty(), Le("100"), Gt("0")}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}, "DataScope": {NotEmpty()}, "Level": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
	InspectVerify          = Rules{"TaskId": {NotEmpty()}, "State": {NotEmpty()}}
	RecordIdVerify         = Rules{"RecordId": {NotEmpty()}}
	EmptyAppIdVerify       = Rules{"AppId": {NotEmpty()}}
	RecordSubmitVerify     = Rules{"AppId": {NotEmpty()}, "Form": {NotEmpty()}, "CurrentNode": {NotEmpty()}, "DeptId": {NotEmpty()}}
	AddApp                 = Rules{"Name": {NotEmpty()}, "Icon": {NotEmpty()}}
)
