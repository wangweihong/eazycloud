package iharbor

type ConfigRequest struct {
}

type ConfigResponse struct {
	AuthMode struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"auth_mode"`
	CountPerProject struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"count_per_project"`
	EmailFrom struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"email_from"`
	EmailHost struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"email_host"`
	EmailIdentity struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"email_identity"`
	EmailInsecure struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"email_insecure"`
	EmailPort struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"email_port"`
	EmailSsl struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"email_ssl"`
	EmailUsername struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"email_username"`
	HTTPAuthproxyEndpoint struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"http_authproxy_endpoint"`
	HTTPAuthproxyServerCertificate struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"http_authproxy_server_certificate"`
	HTTPAuthproxySkipSearch struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"http_authproxy_skip_search"`
	HTTPAuthproxyTokenreviewEndpoint struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"http_authproxy_tokenreview_endpoint"`
	HTTPAuthproxyVerifyCert struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"http_authproxy_verify_cert"`
	LdapBaseDn struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_base_dn"`
	LdapFilter struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_filter"`
	LdapGroupAdminDn struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_group_admin_dn"`
	LdapGroupAttributeName struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_group_attribute_name"`
	LdapGroupBaseDn struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_group_base_dn"`
	LdapGroupMembershipAttribute struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_group_membership_attribute"`
	LdapGroupSearchFilter struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_group_search_filter"`
	LdapGroupSearchScope struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"ldap_group_search_scope"`
	LdapScope struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"ldap_scope"`
	LdapSearchDn struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_search_dn"`
	LdapTimeout struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"ldap_timeout"`
	LdapUID struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_uid"`
	LdapURL struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"ldap_url"`
	LdapVerifyCert struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"ldap_verify_cert"`
	NotificationEnable struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"notification_enable"`
	OidcClientID struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"oidc_client_id"`
	OidcEndpoint struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"oidc_endpoint"`
	OidcGroupsClaim struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"oidc_groups_claim"`
	OidcName struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"oidc_name"`
	OidcScope struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"oidc_scope"`
	OidcVerifyCert struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"oidc_verify_cert"`
	ProjectCreationRestriction struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"project_creation_restriction"`
	QuotaPerProjectEnable struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"quota_per_project_enable"`
	ReadOnly struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"read_only"`
	RobotTokenDuration struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"robot_token_duration"`
	ScanAllPolicy struct {
		Value    interface{} `json:"value"`
		Editable bool        `json:"editable"`
	} `json:"scan_all_policy"`
	SelfRegistration struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"self_registration"`
	StoragePerProject struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"storage_per_project"`
	TokenExpiration struct {
		Value    int  `json:"value"`
		Editable bool `json:"editable"`
	} `json:"token_expiration"`
	UaaClientID struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"uaa_client_id"`
	UaaClientSecret struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"uaa_client_secret"`
	UaaEndpoint struct {
		Value    string `json:"value"`
		Editable bool   `json:"editable"`
	} `json:"uaa_endpoint"`
	UaaVerifyCert struct {
		Value    bool `json:"value"`
		Editable bool `json:"editable"`
	} `json:"uaa_verify_cert"`
}
