package api

import "time"

// Infisical Agent Vault: proxies, sessions, access bundles and delegations.

// TrafficPolicy is which hosts an agent may reach through an Agent Vault proxy.
type TrafficPolicy string

const (
	// TrafficPolicyAnyHost lets agents reach any host.
	TrafficPolicyAnyHost TrafficPolicy = "any-host"
	// TrafficPolicyServiceHosts lets agents reach only the hosts that the session's services cover, plus the proxy's
	// allowed hosts.
	TrafficPolicyServiceHosts TrafficPolicy = "service-hosts"

	// TrafficPolicyServiceHostsOldName is the old name of service-hosts. Infisical stores service-hosts under this name
	// and still accepts it, while a version from before the rename accepts and returns only this name.
	TrafficPolicyServiceHostsOldName TrafficPolicy = "bundle-hosts"
)

// SessionStatus is the status of an Agent Vault session.
type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusRevoked SessionStatus = "revoked"
	SessionStatusExpired SessionStatus = "expired"
)

// SessionScope is whose Agent Vault sessions to list.
type SessionScope string

const (
	SessionScopeMine SessionScope = "mine"
	SessionScopeAll  SessionScope = "all"
)

// MemberType is a kind of Agent Vault member: a user, a machine identity or a group.
type MemberType string

const (
	MemberTypeUser            MemberType = "user"
	MemberTypeMachineIdentity MemberType = "machineIdentity"
	MemberTypeGroup           MemberType = "group"
)

// DelegationScope is which users a machine identity can create Agent Vault sessions on behalf of.
type DelegationScope string

const (
	// DelegationScopeAllMembers covers every Agent Vault member, admins included.
	DelegationScopeAllMembers DelegationScope = "all-members"
	// DelegationScopeSelected covers only the users and groups the delegation names.
	DelegationScopeSelected DelegationScope = "selected"
)

// OrderDirection is which way a list is sorted.
type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "asc"
	OrderDirectionDesc OrderDirection = "desc"
)

// ServiceType is a kind of service Agent Vault serves. Each returns only its own fields.
type ServiceType string

const (
	ServiceTypeHTTP     ServiceType = "http"
	ServiceTypeMCP      ServiceType = "mcp"
	ServiceTypePostgres ServiceType = "postgres"
	ServiceTypeMySQL    ServiceType = "mysql"
	ServiceTypeSSH      ServiceType = "ssh"
)

// SubstitutionSurface is a part of a request where the proxy replaces a substitution's placeholder with its value.
type SubstitutionSurface string

const (
	SubstitutionSurfacePath   SubstitutionSurface = "path"
	SubstitutionSurfaceQuery  SubstitutionSurface = "query"
	SubstitutionSurfaceHeader SubstitutionSurface = "header"
	SubstitutionSurfaceBody   SubstitutionSurface = "body"
)

// MCPTransport is how an agent talks to an MCP server.
type MCPTransport string

const (
	MCPTransportStreamableHTTP MCPTransport = "streamable-http"
	// MCPTransportSSE is the legacy HTTP+SSE transport.
	MCPTransportSSE MCPTransport = "sse"
)

// Proxies

// Proxy is an Agent Vault proxy. Only admins see its settings: TrafficPolicy, AllowedHosts, PollInterval and
// CreatedAt.
type Proxy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// What the proxy is for, or nil if it has no description.
	Description *string `json:"description"`
	// When the proxy last checked in, or nil if it never has.
	Heartbeat *time.Time `json:"heartbeat"`
	// Whether the proxy has checked in recently enough to be considered up.
	IsHealthy bool `json:"isHealthy"`
	// The SHA-256 fingerprint of the proxy's certificate authority.
	RootCaFingerprint *string `json:"rootCaFingerprint"`
	// When the proxy's certificate authority expires.
	RootCaExpiresAt *time.Time    `json:"rootCaExpiresAt"`
	TrafficPolicy   TrafficPolicy `json:"trafficPolicy,omitempty"`
	// Hosts that stay reachable under the service-hosts traffic policy, as one comma-separated string.
	AllowedHosts *string `json:"allowedHosts,omitempty"`
	// How often, in seconds, the proxy refreshes its sessions and settings.
	PollInterval int        `json:"pollInterval,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
}

// ProxyRef is an Agent Vault proxy, by ID and name.
type ProxyRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProxyWithToken is an Agent Vault proxy with its proxy token.
type ProxyWithToken struct {
	Proxy Proxy `json:"proxy"`
	// The proxy token. It's only returned once.
	ProxyToken string `json:"proxyToken"`
}

type ListProxiesResponse struct {
	Proxies    []Proxy `json:"proxies"`
	TotalCount int     `json:"totalCount"`
}

type SingleProxyResponse struct {
	Proxy Proxy `json:"proxy"`
}

type DeleteProxyResponse struct {
	Proxy ProxyRef `json:"proxy"`
}

// Sessions

// SessionActor is the user or machine identity an Agent Vault session belongs to. Type says which: a user has
// Username, Email, FirstName and LastName, and a machine identity has Name.
type SessionActor struct {
	Type MemberType `json:"type"`
	// The ID of the user or machine identity, or nil once it's deleted.
	ID *string `json:"id"`
	// The name of the machine identity. Once it's deleted, the name recorded when the session was created.
	Name string `json:"name,omitempty"`
	// The username of the user. Once the user is deleted, the email address recorded when the session was created.
	Username  string  `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	// The last name of the user, or nil once the user is deleted.
	LastName *string `json:"lastName,omitempty"`
}

// SessionUser is the user an Agent Vault session was started on behalf of.
type SessionUser struct {
	Type MemberType `json:"type"`
	// The ID of the user, or nil once the user is deleted.
	ID *string `json:"id"`
	// The username of the user. Once the user is deleted, the email address recorded when the session was created.
	Username  string  `json:"username"`
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	// The last name of the user, or nil once the user is deleted.
	LastName *string `json:"lastName"`
}

// SessionAccessBundle is the access bundle an Agent Vault session carries.
type SessionAccessBundle struct {
	// The ID of the access bundle, or nil once it's deleted.
	ID   *string `json:"id"`
	Name string  `json:"name"`
}

// ServiceTLS is how Agent Vault connects to a database service over TLS.
type ServiceTLS struct {
	Enabled bool `json:"enabled"`
	// Whether the server's certificate is checked.
	Verify bool `json:"verify"`
	// The CA certificate the server's certificate is checked against, or nil for the system's.
	CA *string `json:"ca"`
}

// SessionServiceCredential is the credential of a service sent with a session, without its token, password or values.
type SessionServiceCredential struct {
	// "bearer", "basic", "passthrough" or "ssh-key".
	Type string `json:"type"`
	// Set for a bearer credential.
	HeaderName string `json:"headerName,omitempty"`
	// Set for a bearer credential.
	HeaderPrefix string `json:"headerPrefix,omitempty"`
	// Set for an SSH key credential: whether it has an SSH certificate.
	HasCertificate bool `json:"hasCertificate,omitempty"`
}

// SessionCustomHeader is a header the proxy adds to a session's service, without its value.
type SessionCustomHeader struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

// SessionSubstitution is a placeholder the proxy replaces in a session's service, without its value.
type SessionSubstitution struct {
	Placeholder string                `json:"placeholder"`
	Surfaces    []SubstitutionSurface `json:"surfaces"`
}

// SessionService is a service sent with a session, as Infisical returns it: without its token, password or values.
// Type says which fields are set: HostPattern, AllowedMethods, AllowedPathPrefixes, CustomHeaders and Substitutions
// for "http"; ServerURL, Transport, HostPattern, AllowedTools and CustomHeaders for "mcp"; Host, Port, Database and
// TLS for "postgres" and "mysql"; and Host and Port for "ssh".
type SessionService struct {
	ID string `json:"id"`
	// The name the service was sent with, or its host if it was sent without one.
	Name string `json:"name"`
	// What the service is for, or nil if it was sent without a description.
	Description *string                  `json:"description,omitempty"`
	Type        ServiceType              `json:"type"`
	Credential  SessionServiceCredential `json:"credential"`

	HostPattern string `json:"hostPattern,omitempty"`
	// The HTTP methods the service allows, or nil for every method.
	AllowedMethods []string `json:"allowedMethods,omitempty"`
	// The path prefixes the service allows, or nil for every path.
	AllowedPathPrefixes []string              `json:"allowedPathPrefixes,omitempty"`
	CustomHeaders       []SessionCustomHeader `json:"customHeaders,omitempty"`
	Substitutions       []SessionSubstitution `json:"substitutions,omitempty"`

	ServerURL string       `json:"serverUrl,omitempty"`
	Transport MCPTransport `json:"transport,omitempty"`
	// The MCP tools that agents can call through the service, or nil for every tool. Also nil from an Infisical version
	// that doesn't return it.
	AllowedTools []string `json:"allowedTools,omitempty"`

	Host     string      `json:"host,omitempty"`
	Port     int         `json:"port,omitempty"`
	Database string      `json:"database,omitempty"`
	TLS      *ServiceTLS `json:"tls,omitempty"`
}

// Session is an Agent Vault session.
type Session struct {
	ID    string       `json:"id"`
	Label *string      `json:"label"`
	Actor SessionActor `json:"actor"`
	// The user the session was started on behalf of, or nil when the session is the actor's own.
	OnBehalfOf *SessionUser  `json:"onBehalfOf"`
	Status     SessionStatus `json:"status"`
	// When the session expires, or nil when it never does.
	ExpiresAt *time.Time `json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	// The access bundle the session carries, or nil if it carries only services of its own.
	AccessBundle *SessionAccessBundle `json:"accessBundle"`
	// The services sent with the session when it was created, without their secrets. Empty if it carries only an
	// access bundle.
	Services []SessionService `json:"services"`
	// The session's metadata as key-value pairs, such as {"runId": "8812"}. Empty if it has none.
	Metadata map[string]string `json:"metadata"`
}

// CreatedSession is an Agent Vault session as returned at creation.
type CreatedSession struct {
	ID string `json:"id"`
	// The Agent Vault project the session belongs to.
	ProjectID string  `json:"projectId"`
	Label     *string `json:"label"`
	// The session token. It's only returned once.
	Token string `json:"token"`
	// When the session expires, or nil when it never does.
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	// The access bundle the session carries, or nil if it carries only services of its own.
	AccessBundle *SessionAccessBundle `json:"accessBundle"`
	// The services sent with the session, without their secrets. Empty if it carries only an access bundle.
	Services []SessionService `json:"services"`
	// The session's metadata as key-value pairs. Empty if it has none.
	Metadata map[string]string `json:"metadata"`
}

type ListSessionsResponse struct {
	Sessions   []Session `json:"sessions"`
	TotalCount int       `json:"totalCount"`
}

type SingleSessionResponse struct {
	Session Session `json:"session"`
}

type CreateSessionResponse struct {
	Session CreatedSession `json:"session"`
}

// Access bundles

// ServiceRef is an Agent Vault service, by ID and name.
type ServiceRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AccessBundleRef is an Agent Vault access bundle, by ID and name.
type AccessBundleRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AccessBundle is an Agent Vault access bundle.
type AccessBundle struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// What the access bundle is for, or nil if it has no description.
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// AccessBundleSummary is an Agent Vault access bundle as listed, with what it holds.
type AccessBundleSummary struct {
	AccessBundle
	// How many services the access bundle holds.
	ServiceCount int `json:"serviceCount"`
	// How many users, machine identities and groups can reach the access bundle.
	MemberCount int `json:"memberCount"`
	// Every host pattern the access bundle's services cover.
	HostPatterns []string `json:"hostPatterns"`
	// The services that still need your personal values, or, when listed on behalf of a user, that user's.
	ServicesRequiringSetup []ServiceRef `json:"servicesRequiringSetup"`
}

// AccessBundleService is a service in an Agent Vault access bundle. The SDK doesn't manage services, so a service's
// nested settings, such as its credential, are left as the API returns them. Type says which of the other fields are
// set, as for SessionService.
type AccessBundleService struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// What the service is for, or nil if it has no description.
	Description *string        `json:"description,omitempty"`
	Type        ServiceType    `json:"type"`
	Credential  map[string]any `json:"credential"`
	// Only set for admins.
	VariableReferences []map[string]any `json:"variableReferences,omitempty"`
	// Your own values for the service, or nil when the admin sets the whole credential.
	PersonalValues map[string]any `json:"personalValues"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`

	HostPattern string `json:"hostPattern,omitempty"`
	// The HTTP methods the service allows, or nil for every method.
	AllowedMethods []string `json:"allowedMethods,omitempty"`
	// The path prefixes the service allows, or nil for every path.
	AllowedPathPrefixes []string         `json:"allowedPathPrefixes,omitempty"`
	CustomHeaders       []map[string]any `json:"customHeaders,omitempty"`
	Substitutions       []map[string]any `json:"substitutions,omitempty"`

	ServerURL string `json:"serverUrl,omitempty"`
	Transport string `json:"transport,omitempty"`
	// The MCP tools that agents can call through the service, by exact name, or nil for every tool.
	AllowedTools []string `json:"allowedTools,omitempty"`

	Host     string      `json:"host,omitempty"`
	Port     int         `json:"port,omitempty"`
	Database string      `json:"database,omitempty"`
	TLS      *ServiceTLS `json:"tls,omitempty"`
}

// AccessBundleWithServices is an Agent Vault access bundle with its services.
type AccessBundleWithServices struct {
	AccessBundle
	Services []AccessBundleService `json:"services"`
}

type ListAccessBundlesResponse struct {
	AccessBundles []AccessBundleSummary `json:"accessBundles"`
	TotalCount    int                   `json:"totalCount"`
}

type GetAccessBundleResponse struct {
	AccessBundle AccessBundleWithServices `json:"accessBundle"`
}

type SingleAccessBundleResponse struct {
	AccessBundle AccessBundle `json:"accessBundle"`
}

// Grants and members

// ActorRef is a user, machine identity or group, by type and ID.
type ActorRef struct {
	Type MemberType `json:"type"`
	ID   string     `json:"id"`
}

// Grant is a new grant of an access bundle.
type Grant struct {
	ActorRef
	GrantedAt time.Time `json:"grantedAt"`
}

// Grantee is a user, machine identity or group that an Agent Vault access bundle is granted to. Type says which: a
// user has Username, Email, FirstName, LastName and IsOrgMembershipPending, and a machine identity or group has Name.
type Grantee struct {
	Type      MemberType `json:"type"`
	ID        string     `json:"id"`
	Name      string     `json:"name,omitempty"`
	Username  string     `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	// Whether the user still has an unaccepted invitation to the organization.
	IsOrgMembershipPending bool `json:"isOrgMembershipPending,omitempty"`
	// When the access bundle was granted.
	GrantedAt time.Time `json:"grantedAt"`
}

// Member is a user, machine identity or group that's a member of Agent Vault, with its Agent Vault role. Type says
// which: a user has Username, Email, FirstName, LastName and IsOrgMembershipPending, a machine identity has Name,
// IsManagedByAgentVault and OrgID, and a group has Name.
type Member struct {
	Type      MemberType `json:"type"`
	ID        string     `json:"id"`
	Name      string     `json:"name,omitempty"`
	Username  string     `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	// Whether the user still has an unaccepted invitation to the organization.
	IsOrgMembershipPending bool `json:"isOrgMembershipPending,omitempty"`
	// Whether Agent Vault created the machine identity.
	IsManagedByAgentVault bool `json:"isManagedByAgentVault,omitempty"`
	// The ID of the organization the machine identity belongs to.
	OrgID *string `json:"orgId,omitempty"`
	// The Agent Vault role, such as "admin" or "member".
	Role string `json:"role"`
	// Whether the member can currently reach Agent Vault.
	IsActive bool `json:"isActive"`
	// When the member was added to Agent Vault.
	AddedAt time.Time `json:"addedAt"`
}

type ListGranteesResponse struct {
	Members    []Grantee `json:"members"`
	TotalCount int       `json:"totalCount"`
}

type ListMembersResponse struct {
	Members    []Member `json:"members"`
	TotalCount int      `json:"totalCount"`
}

// GrantResponse is the result of granting an access bundle.
type GrantResponse struct {
	// The new grants.
	Members []Grant `json:"members"`
	// The grantees that already had the access bundle, left as they were.
	Skipped []ActorRef `json:"skipped"`
}

// RevokeGrantsResponse is the result of revoking an access bundle.
type RevokeGrantsResponse struct {
	// The revoked grants.
	Members []ActorRef `json:"members"`
	// The grantees that didn't have the access bundle, so nothing was revoked.
	Skipped []ActorRef `json:"skipped"`
}

// Delegations

// DelegationUser is a user in an Agent Vault delegation.
type DelegationUser struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	// When the user was added to the delegation.
	AddedAt time.Time `json:"addedAt"`
	// The access bundles this user can use and the machine identity can't. Only set by the methods that return them.
	MissingAccessBundles []AccessBundleRef `json:"missingAccessBundles,omitempty"`
}

// DelegationGroup is a group in an Agent Vault delegation. Its members are covered while they stay in the group.
type DelegationGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// When the group was added to the delegation.
	AddedAt time.Time `json:"addedAt"`
	// The access bundles this group can use and the machine identity can't. Only set by the methods that return them.
	MissingAccessBundles []AccessBundleRef `json:"missingAccessBundles,omitempty"`
}

// DelegationMachineIdentity is the machine identity a delegation belongs to.
type DelegationMachineIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Delegation is which users a machine identity can create Agent Vault sessions on behalf of.
//
// Get, Set, AddSubjects and RemoveSubjects also return the access bundles the machine identity can't create sessions
// over in MissingAccessBundles, and in each user's and group's MissingAccessBundles. A session on behalf of a user can
// use only an access bundle that both the user and the machine identity can use. The lists are empty when a machine
// identity that isn't an admin gets its own delegation.
type Delegation struct {
	MachineIdentity DelegationMachineIdentity `json:"machineIdentity"`
	Scope           DelegationScope           `json:"scope"`
	Users           []DelegationUser          `json:"users"`
	Groups          []DelegationGroup         `json:"groups"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	// The access bundles that users the delegation covers can use and the machine identity can't.
	MissingAccessBundles []AccessBundleRef `json:"missingAccessBundles,omitempty"`
}

type ListDelegationsResponse struct {
	Delegations []Delegation `json:"delegations"`
	TotalCount  int          `json:"totalCount"`
}

type SingleDelegationResponse struct {
	Delegation Delegation `json:"delegation"`
}

// RemoveDelegationSubjectsResponse holds the delegation after the change, or nil if it no longer names anyone.
type RemoveDelegationSubjectsResponse struct {
	Delegation *Delegation `json:"delegation"`
}
