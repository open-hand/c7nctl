package model

type UserInfo struct {
	ID                   int    `json:"id"`
	OrganizationID       int    `json:"organizationId"`
	OrganizationName     string `json:"organizationName"`
	OrganizationCode     string `json:"organizationCode"`
	LoginName            string `json:"loginName"`
	Email                string `json:"email"`
	RealName             string `json:"realName"`
	Phone                string `json:"phone"`
	InternationalTelCode string `json:"internationalTelCode"`
	ImageURL             string `json:"imageUrl"`
	Language             string `json:"language"`
	TimeZone             string `json:"timeZone"`
	Locked               bool   `json:"locked"`
	Ldap                 bool   `json:"ldap"`
	Enabled              bool   `json:"enabled"`
	Admin                bool   `json:"admin"`
	ObjectVersionNumber  int    `json:"objectVersionNumber"`
}


type UserAttrInfo struct {

	IamUserId  int `json:"iamUserId"`
	GitlabUserId int  `json:"gitlabUserId"`
}
