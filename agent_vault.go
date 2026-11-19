package infisical

import (
	"errors"
	"strings"

	api "github.com/infisical/go-sdk/packages/api/agentvault"
)

// Enums

type AgentVaultTrafficPolicy = api.TrafficPolicy
type AgentVaultSessionStatus = api.SessionStatus
type AgentVaultSessionScope = api.SessionScope
type AgentVaultMemberType = api.MemberType
type AgentVaultDelegationScope = api.DelegationScope
type AgentVaultOrderDirection = api.OrderDirection
type AgentVaultServiceType = api.ServiceType
type AgentVaultSubstitutionSurface = api.SubstitutionSurface
type AgentVaultMCPTransport = api.MCPTransport

const (
	AgentVaultTrafficPolicyAnyHost      = api.TrafficPolicyAnyHost
	AgentVaultTrafficPolicyServiceHosts = api.TrafficPolicyServiceHosts

	AgentVaultSessionStatusActive  = api.SessionStatusActive
	AgentVaultSessionStatusRevoked = api.SessionStatusRevoked
	AgentVaultSessionStatusExpired = api.SessionStatusExpired

	AgentVaultSessionScopeMine = api.SessionScopeMine
	AgentVaultSessionScopeAll  = api.SessionScopeAll

	AgentVaultMemberTypeUser            = api.MemberTypeUser
	AgentVaultMemberTypeMachineIdentity = api.MemberTypeMachineIdentity
	AgentVaultMemberTypeGroup           = api.MemberTypeGroup

	AgentVaultDelegationScopeAllMembers = api.DelegationScopeAllMembers
	AgentVaultDelegationScopeSelected   = api.DelegationScopeSelected

	AgentVaultOrderDirectionAsc  = api.OrderDirectionAsc
	AgentVaultOrderDirectionDesc = api.OrderDirectionDesc

	AgentVaultServiceTypeHTTP     = api.ServiceTypeHTTP
	AgentVaultServiceTypeMCP      = api.ServiceTypeMCP
	AgentVaultServiceTypePostgres = api.ServiceTypePostgres
	AgentVaultServiceTypeMySQL    = api.ServiceTypeMySQL
	AgentVaultServiceTypeSSH      = api.ServiceTypeSSH

	AgentVaultSubstitutionSurfacePath   = api.SubstitutionSurfacePath
	AgentVaultSubstitutionSurfaceQuery  = api.SubstitutionSurfaceQuery
	AgentVaultSubstitutionSurfaceHeader = api.SubstitutionSurfaceHeader
	AgentVaultSubstitutionSurfaceBody   = api.SubstitutionSurfaceBody

	AgentVaultMCPTransportStreamableHTTP = api.MCPTransportStreamableHTTP
	AgentVaultMCPTransportSSE            = api.MCPTransportSSE
)

// Results

type AgentVaultProxy = api.Proxy
type AgentVaultProxyRef = api.ProxyRef
type AgentVaultProxyWithToken = api.ProxyWithToken
type AgentVaultListProxiesResult = api.ListProxiesResponse

type AgentVaultSession = api.Session
type AgentVaultCreatedSession = api.CreatedSession
type AgentVaultSessionActor = api.SessionActor
type AgentVaultSessionUser = api.SessionUser
type AgentVaultSessionAccessBundle = api.SessionAccessBundle
type AgentVaultSessionService = api.SessionService
type AgentVaultSessionServiceCredential = api.SessionServiceCredential
type AgentVaultSessionCustomHeader = api.SessionCustomHeader
type AgentVaultSessionSubstitution = api.SessionSubstitution
type AgentVaultServiceTLS = api.ServiceTLS
type AgentVaultListSessionsResult = api.ListSessionsResponse

type AgentVaultAccessBundle = api.AccessBundle
type AgentVaultAccessBundleRef = api.AccessBundleRef
type AgentVaultAccessBundleSummary = api.AccessBundleSummary
type AgentVaultAccessBundleWithServices = api.AccessBundleWithServices
type AgentVaultAccessBundleService = api.AccessBundleService
type AgentVaultServiceRef = api.ServiceRef
type AgentVaultListAccessBundlesResult = api.ListAccessBundlesResponse

type AgentVaultActorRef = api.ActorRef
type AgentVaultGrant = api.Grant
type AgentVaultGrantee = api.Grantee
type AgentVaultMember = api.Member
type AgentVaultListGranteesResult = api.ListGranteesResponse
type AgentVaultListMembersResult = api.ListMembersResponse
type AgentVaultGrantResult = api.GrantResponse
type AgentVaultRevokeGrantsResult = api.RevokeGrantsResponse

type AgentVaultDelegation = api.Delegation
type AgentVaultDelegationUser = api.DelegationUser
type AgentVaultDelegationGroup = api.DelegationGroup
type AgentVaultDelegationMachineIdentity = api.DelegationMachineIdentity
type AgentVaultListDelegationsResult = api.ListDelegationsResponse

// Inline services, sent with a session

type AgentVaultInlineService = api.InlineService
type AgentVaultInlineHTTPService = api.InlineHTTPService
type AgentVaultInlineMCPService = api.InlineMCPService
type AgentVaultInlineDatabaseService = api.InlineDatabaseService
type AgentVaultInlineSSHService = api.InlineSSHService
type AgentVaultInlineTLS = api.InlineTLS
type AgentVaultInlineCustomHeader = api.InlineCustomHeader
type AgentVaultInlineSubstitution = api.InlineSubstitution
type AgentVaultHTTPCredential = api.HTTPCredential
type AgentVaultSSHCredential = api.SSHCredential
type AgentVaultBearerCredential = api.BearerCredential
type AgentVaultBasicCredential = api.BasicCredential
type AgentVaultPassthroughCredential = api.PassthroughCredential
type AgentVaultSSHKeyCredential = api.SSHKeyCredential

// StringPtr returns a pointer to the given string value, for options where a pointer to an empty string means
// something different from leaving the value out.
func StringPtr(v string) *string { return &v }

// Options

// AgentVaultProjectOptions names the Agent Vault project a delegation belongs to.
type AgentVaultProjectOptions struct {
	// The ID of the Agent Vault project. Required.
	ProjectID string
}

// AgentVaultOnBehalfOf names the user to act on behalf of, by email address or by ID, not both.
type AgentVaultOnBehalfOf struct {
	UserEmail string
	UserID    string
}

type AgentVaultListProxiesOptions struct {
	// The ID of the Agent Vault project to list proxies from. Required.
	ProjectID string
	// Match proxies by name or description.
	Search string
	// What to sort proxies by: "name" or "createdAt". Defaults to "createdAt".
	OrderBy string
	// Which way to sort. Defaults to descending.
	OrderDirection AgentVaultOrderDirection
	// The maximum number of proxies to return (1-100). Defaults to 20.
	Limit int
	// How many proxies to skip. Defaults to 0.
	Offset int
}

type AgentVaultCreateProxyOptions struct {
	// The ID of the Agent Vault project to register the proxy in. Required.
	ProjectID string
	// The name of the proxy (a slug, at most 64 characters).
	Name string
	// Which hosts an agent may reach through this proxy: any host, or only the hosts that the session's services
	// cover. Defaults to AgentVaultTrafficPolicyAnyHost.
	TrafficPolicy AgentVaultTrafficPolicy
	// Hosts that stay reachable under the service-hosts traffic policy even though none of the session's services
	// cover them.
	AllowedHosts []string
	// How often, in seconds, the proxy refreshes its sessions and settings (10-300). Defaults to 60.
	PollInterval int
	// What the proxy is for, up to 256 characters, without line breaks or other control characters.
	Description string
}

// AgentVaultUpdateProxyOptions changes only the values you set.
type AgentVaultUpdateProxyOptions struct {
	// The new name of the proxy (a slug, at most 64 characters).
	Name *string
	// Which hosts an agent may reach through this proxy.
	TrafficPolicy AgentVaultTrafficPolicy
	// Hosts that stay reachable under the service-hosts traffic policy. Leave it nil to keep them, or set an empty,
	// non-nil list to remove every allowed host.
	AllowedHosts []string
	// How often, in seconds, the proxy refreshes its sessions and settings (10-300).
	PollInterval int
	// The new description, up to 256 characters, without line breaks or other control characters. Point it at an
	// empty string to remove the description.
	Description *string
}

type AgentVaultListSessionsOptions struct {
	// The ID of the Agent Vault project to list sessions from. Required.
	ProjectID string
	// Whose sessions to list: your own, or everyone's (admins only). Defaults to AgentVaultSessionScopeMine.
	Scope AgentVaultSessionScope
	// Only return sessions with any of these statuses.
	Statuses []AgentVaultSessionStatus
	// Match sessions by label, actor name, actor email, the name or email of the user a session was started on behalf
	// of, or access bundle name.
	Search string
	// Only return sessions whose metadata matches every one of these key-value pairs exactly, such as
	// {"runId": "8812"}. Up to 20 pairs.
	Metadata map[string]string
	// The maximum number of sessions to return (1-100). Defaults to 20.
	Limit int
	// How many sessions to skip. Defaults to 0.
	Offset int
}

type AgentVaultCreateSessionOptions struct {
	// The ID of the Agent Vault project to create the session in. Required.
	ProjectID string
	// The name of the access bundle the session carries. Set AccessBundle, Services, or both.
	AccessBundle string
	// Services the session carries itself, instead of or beside its access bundle's: up to 20. Only Agent Vault admins
	// can send them. Each has the same fields as a service in an access bundle, except that its name is optional for
	// HTTP and MCP services (database and SSH services need one), its values can't use variables, and it can't use
	// OAuth or member-provided values.
	Services []AgentVaultInlineService
	// How long the session lasts: a duration such as "30m", "8h" or "7d" (at least 1m), or "never". Defaults to "7d",
	// or to "1h" for a session on behalf of a user.
	TTL string
	// A short description of what the session is for, such as "nightly deploy". Up to 128 characters, without line
	// breaks, tabs, or invisible characters.
	Label string
	// Key-value pairs to tag the session with, such as {"runId": "8812"}. Up to 20 keys. A key is 1 to 255 characters
	// and a value up to 1020, without line breaks, tabs, or invisible characters. Keys that match once trimmed are
	// rejected.
	Metadata map[string]string
	// The user to start the session on behalf of. Leave it nil for a session of your own.
	OnBehalfOf *AgentVaultOnBehalfOf
}

// AgentVaultUpdateSessionOptions changes only the values you set.
type AgentVaultUpdateSessionOptions struct {
	// The new label, up to 128 characters, without line breaks, tabs, or invisible characters. Point it at an empty
	// string to remove the label.
	Label *string
	// The session's new metadata, which replaces all of it. Leave it nil to keep it, or set an empty, non-nil map to
	// remove it. Up to 20 keys. A key is 1 to 255 characters and a value up to 1020, without line breaks, tabs, or
	// invisible characters. Keys that match once trimmed are rejected.
	Metadata map[string]string
}

type AgentVaultListAccessBundlesOptions struct {
	// The ID of the Agent Vault project to list access bundles from. Required.
	ProjectID string
	// Match access bundles by name or description.
	Search string
	// true lists only the access bundles with services that still need your personal values, and false lists the
	// rest. Defaults to both.
	RequiresSetup *bool
	// true lists only the access bundles granted to you, directly or through a group, and false lists the ones you
	// reach only as an admin. Defaults to both. Can't be combined with OnBehalfOf.
	Granted *bool
	// As a machine identity whose delegation covers the user, list the access bundles that a session on the user's
	// behalf could use: those both the user and your machine identity can use. RequiresSetup and
	// ServicesRequiringSetup then refer to the user's personal values.
	OnBehalfOf *AgentVaultOnBehalfOf
	// What to sort access bundles by: "name", "serviceCount" or "createdAt". Defaults to "createdAt".
	OrderBy string
	// Which way to sort. Defaults to descending.
	OrderDirection AgentVaultOrderDirection
	// The maximum number of access bundles to return (1-100). Defaults to 20.
	Limit int
	// How many access bundles to skip. Defaults to 0.
	Offset int
}

type AgentVaultCreateAccessBundleOptions struct {
	// The ID of the Agent Vault project to create the access bundle in. Required.
	ProjectID string
	// The name of the access bundle (a slug, at most 64 characters).
	Name string
	// What the access bundle is for, up to 256 characters, without line breaks or other control characters.
	Description string
}

// AgentVaultUpdateAccessBundleOptions changes only the values you set.
type AgentVaultUpdateAccessBundleOptions struct {
	// The new name of the access bundle (a slug, at most 64 characters).
	Name *string
	// The new description, up to 256 characters, without line breaks or other control characters. Point it at an
	// empty string to remove the description.
	Description *string
}

type AgentVaultListGranteesOptions struct {
	// List only users, only machine identities, or only groups. Defaults to all three.
	ActorType AgentVaultMemberType
	// Match by name, username or email address.
	Search string
	// The maximum number to return (1-100). Defaults to 20.
	Limit int
	// How many to skip. Defaults to 0.
	Offset int
}

// AgentVaultChangeGrantsOptions names who to grant an access bundle to, or revoke it from. Name at least one, and at
// most 100 in total.
type AgentVaultChangeGrantsOptions struct {
	UserIDs            []string
	MachineIdentityIDs []string
	GroupIDs           []string
}

type AgentVaultListDelegationsOptions struct {
	// The ID of the Agent Vault project to list delegations from. Required.
	ProjectID string
	// Only return the delegations of these machine identities, up to 100.
	MachineIdentityIDs []string
	// The maximum number of delegations to return (1-100). Defaults to 100.
	Limit int
	// How many delegations to skip. Defaults to 0.
	Offset int
}

type AgentVaultSetDelegationOptions struct {
	// The ID of the Agent Vault project the delegation belongs to. Required.
	ProjectID string
	// AgentVaultDelegationScopeAllMembers for every Agent Vault member, admins included, or
	// AgentVaultDelegationScopeSelected for only the users and groups you name.
	Scope AgentVaultDelegationScope
	// With the selected scope, the users the machine identity can act on behalf of, up to 100. Each must be a member
	// of Agent Vault. For a user who's in Agent Vault only through a group, name the group instead.
	UserIDs []string
	// With the selected scope, the groups whose members the machine identity can act on behalf of, up to 100. Each
	// must be a member of Agent Vault.
	GroupIDs []string
}

// AgentVaultChangeDelegationSubjectsOptions names users and groups to add to a delegation, or take out of it.
type AgentVaultChangeDelegationSubjectsOptions struct {
	// The ID of the Agent Vault project the delegation belongs to. Required.
	ProjectID string
	// Up to 100 users. Users the change wouldn't affect are skipped.
	UserIDs []string
	// Up to 100 groups. Groups the change wouldn't affect are skipped.
	GroupIDs []string
}

type AgentVaultListAvailableDelegationSubjectsOptions struct {
	// The ID of the Agent Vault project the delegation belongs to. Required.
	ProjectID string
	// List only users or only groups. Defaults to both.
	ActorType AgentVaultMemberType
	// Match users by name, username or email address, and groups by name.
	Search string
	// The maximum number of users and groups to return (1-100). Defaults to 20.
	Limit int
	// How many users and groups to skip. Defaults to 0.
	Offset int
}

// Interfaces

// AgentVaultInterface manages Infisical Agent Vault.
type AgentVaultInterface interface {
	Proxies() AgentVaultProxiesInterface
	Sessions() AgentVaultSessionsInterface
	AccessBundles() AgentVaultAccessBundlesInterface
	Delegations() AgentVaultDelegationsInterface
}

// AgentVaultProxiesInterface manages Agent Vault proxies.
type AgentVaultProxiesInterface interface {
	// List lists the project's Agent Vault proxies. A proxy's settings are only set for admins.
	List(options AgentVaultListProxiesOptions) (AgentVaultListProxiesResult, error)
	// Create registers an Agent Vault proxy and issues its proxy token. The token is only returned once.
	Create(options AgentVaultCreateProxyOptions) (AgentVaultProxyWithToken, error)
	// Update updates an Agent Vault proxy's name, description or settings. Only the values you set are changed.
	Update(proxyID string, options AgentVaultUpdateProxyOptions) (AgentVaultProxy, error)
	// Delete deletes an Agent Vault proxy, and returns its ID and name.
	Delete(proxyID string) (AgentVaultProxyRef, error)
	// RotateToken issues a new proxy token. The proxy's earlier tokens stop working and it gets a new certificate
	// authority. The token is only returned once.
	RotateToken(proxyID string) (AgentVaultProxyWithToken, error)
	// Revoke revokes an Agent Vault proxy's token. Rotate its token to bring it back.
	Revoke(proxyID string) (AgentVaultProxy, error)
}

// AgentVaultSessionsInterface manages Agent Vault sessions.
type AgentVaultSessionsInterface interface {
	// List lists Agent Vault sessions.
	List(options AgentVaultListSessionsOptions) (AgentVaultListSessionsResult, error)
	// Get gets an Agent Vault session.
	Get(sessionID string) (AgentVaultSession, error)
	// Create creates an Agent Vault session with an access bundle you can reach, services of its own, or both, for
	// yourself or, as a machine identity with a delegation, on behalf of a user.
	//
	// Only Agent Vault admins can send services with a session. On behalf of a user, that's your machine identity:
	// the user doesn't need the Admin role, and the services don't use the user's personal values.
	//
	// To start a session on behalf of a user, you must be a machine identity whose delegation covers the user. The
	// proxy then uses the user's personal values, and the session can use only an access bundle that both the user
	// and your machine identity can use. The session still belongs to your machine identity, so the user can't see it
	// unless they're an Agent Vault admin.
	//
	// The created session includes its token, which is only returned once. Its services leave out their tokens,
	// passwords and values.
	Create(options AgentVaultCreateSessionOptions) (AgentVaultCreatedSession, error)
	// Update updates the label or metadata of an Agent Vault session. Set Label, Metadata, or both; a value you leave
	// out stays as it is. This also works after the session expires or is revoked, but not once the session's owner,
	// or the user it was started on behalf of, is deleted.
	Update(sessionID string, options AgentVaultUpdateSessionOptions) (AgentVaultSession, error)
	// Revoke revokes an Agent Vault session. Revoking an already revoked session changes nothing. The response leaves
	// out the session's access bundle and services, so they're nil; call Get for them.
	Revoke(sessionID string) (AgentVaultSession, error)
}

// AgentVaultAccessBundlesInterface manages Agent Vault access bundles and who they're granted to. Every method except
// List and Get requires the Admin role in Agent Vault.
type AgentVaultAccessBundlesInterface interface {
	// List lists the Agent Vault access bundles you can reach. Admins reach every access bundle.
	List(options AgentVaultListAccessBundlesOptions) (AgentVaultListAccessBundlesResult, error)
	// Get gets an Agent Vault access bundle you can reach, with its services. A service's VariableReferences is only
	// set for admins.
	Get(accessBundleID string) (AgentVaultAccessBundleWithServices, error)
	// Create creates an Agent Vault access bundle. It's granted to you, unless you're in Agent Vault only through a
	// group.
	Create(options AgentVaultCreateAccessBundleOptions) (AgentVaultAccessBundle, error)
	// Update updates an Agent Vault access bundle's name or description. Only the values you set are changed.
	Update(accessBundleID string, options AgentVaultUpdateAccessBundleOptions) (AgentVaultAccessBundle, error)
	// Delete deletes an Agent Vault access bundle.
	Delete(accessBundleID string) (AgentVaultAccessBundle, error)
	// ListGrantees lists the users, machine identities and groups that an Agent Vault access bundle is granted to.
	ListGrantees(accessBundleID string, options AgentVaultListGranteesOptions) (AgentVaultListGranteesResult, error)
	// ListAvailableGrantees lists the Agent Vault members that an access bundle isn't granted to yet, so you can grant
	// it to them.
	ListAvailableGrantees(accessBundleID string, options AgentVaultListGranteesOptions) (AgentVaultListMembersResult, error)
	// Grant grants an Agent Vault access bundle to users, machine identities or groups. Each must be a member of Agent
	// Vault. It returns the new grants, and the grantees that already had the access bundle in Skipped.
	Grant(accessBundleID string, options AgentVaultChangeGrantsOptions) (AgentVaultGrantResult, error)
	// Revoke revokes an Agent Vault access bundle from users, machine identities or groups. It returns the revoked
	// grants, and the grantees that didn't have the access bundle in Skipped.
	Revoke(accessBundleID string, options AgentVaultChangeGrantsOptions) (AgentVaultRevokeGrantsResult, error)
}

// AgentVaultDelegationsInterface manages Agent Vault delegations, which let machine identities create sessions on
// behalf of users. Every method requires the Admin role in Agent Vault, except that a machine identity can get its
// own delegation.
type AgentVaultDelegationsInterface interface {
	// List lists the machine identities that can create Agent Vault sessions on behalf of users. A listed delegation
	// has no MissingAccessBundles; call Get for them.
	List(options AgentVaultListDelegationsOptions) (AgentVaultListDelegationsResult, error)
	// Get gets which users a machine identity can create Agent Vault sessions on behalf of, with the access bundles
	// the machine identity can't create sessions with. If the machine identity has no delegation, it returns an
	// *errors.APIError whose StatusCode is 404.
	Get(machineIdentityID string, options AgentVaultProjectOptions) (AgentVaultDelegation, error)
	// Set sets which users a machine identity can create Agent Vault sessions on behalf of, replacing any delegation
	// it has. The machine identity must be a member of Agent Vault.
	Set(machineIdentityID string, options AgentVaultSetDelegationOptions) (AgentVaultDelegation, error)
	// Delete stops a machine identity from creating Agent Vault sessions on behalf of users. The proxy refuses the
	// sessions it already created on their behalf until a delegation covers those users again, so revoke the sessions
	// to end them for good. The deleted delegation has no MissingAccessBundles.
	Delete(machineIdentityID string, options AgentVaultProjectOptions) (AgentVaultDelegation, error)
	// AddSubjects adds users and groups to a machine identity's delegation, keeping the ones it already names. If the
	// machine identity has no delegation, this creates one with the selected scope. A delegation to all members can't
	// take users or groups.
	AddSubjects(machineIdentityID string, options AgentVaultChangeDelegationSubjectsOptions) (AgentVaultDelegation, error)
	// RemoveSubjects takes users and groups out of a machine identity's delegation. The proxy refuses the sessions it
	// created on behalf of the users you take out until the delegation covers them again, so revoke the sessions to
	// end them for good. Taking out the last user or group removes the delegation, and then it returns nil.
	RemoveSubjects(machineIdentityID string, options AgentVaultChangeDelegationSubjectsOptions) (*AgentVaultDelegation, error)
	// ListAvailableSubjects lists the Agent Vault users and groups that can still be added to a machine identity's
	// delegation. The machine identity must be a member of Agent Vault.
	ListAvailableSubjects(machineIdentityID string, options AgentVaultListAvailableDelegationSubjectsOptions) (AgentVaultListMembersResult, error)
}

// Implementation

type AgentVault struct {
	proxies       *AgentVaultProxies
	sessions      *AgentVaultSessions
	accessBundles *AgentVaultAccessBundles
	delegations   *AgentVaultDelegations
}

type AgentVaultProxies struct{ client *InfisicalClient }
type AgentVaultSessions struct{ client *InfisicalClient }
type AgentVaultAccessBundles struct{ client *InfisicalClient }
type AgentVaultDelegations struct{ client *InfisicalClient }

func NewAgentVault(client *InfisicalClient) AgentVaultInterface {
	return &AgentVault{
		proxies:       &AgentVaultProxies{client: client},
		sessions:      &AgentVaultSessions{client: client},
		accessBundles: &AgentVaultAccessBundles{client: client},
		delegations:   &AgentVaultDelegations{client: client},
	}
}

func (a *AgentVault) Proxies() AgentVaultProxiesInterface             { return a.proxies }
func (a *AgentVault) Sessions() AgentVaultSessionsInterface           { return a.sessions }
func (a *AgentVault) AccessBundles() AgentVaultAccessBundlesInterface { return a.accessBundles }
func (a *AgentVault) Delegations() AgentVaultDelegationsInterface     { return a.delegations }

// requireProject refuses a request without a project. An organization can have several Agent Vault projects, so a
// method that lists or creates something at project level, or addresses a delegation, has to name one.
func requireProject(projectID string) error {
	if projectID == "" {
		return errors.New("agent vault: ProjectID is required")
	}
	return nil
}

func withDefault[T comparable](value T, fallback T) T {
	var zero T
	if value == zero {
		return fallback
	}
	return value
}

// joinHosts joins host patterns into the one comma-separated string the API takes. The API rejects an empty string,
// so it reports false for one.
func joinHosts(allowedHosts []string) (string, bool) {
	hosts := strings.Join(allowedHosts, ",")
	return hosts, strings.TrimSpace(hosts) != ""
}

// toWireTrafficPolicy sends service-hosts by its old name, which every Infisical version accepts.
func toWireTrafficPolicy(trafficPolicy AgentVaultTrafficPolicy) AgentVaultTrafficPolicy {
	if trafficPolicy == api.TrafficPolicyServiceHosts {
		return api.TrafficPolicyServiceHostsOldName
	}
	return trafficPolicy
}

// viewProxy reports the old name of service-hosts, which an Infisical version from before the rename returns, as
// service-hosts.
func viewProxy(proxy AgentVaultProxy) AgentVaultProxy {
	if proxy.TrafficPolicy == api.TrafficPolicyServiceHostsOldName {
		proxy.TrafficPolicy = api.TrafficPolicyServiceHosts
	}
	return proxy
}

// orEmpty gives a session from an Infisical version without inline services or metadata an empty list and map of
// them, rather than none.
func orEmpty(services []AgentVaultSessionService, metadata map[string]string) ([]AgentVaultSessionService, map[string]string) {
	if services == nil {
		services = []AgentVaultSessionService{}
	}
	if metadata == nil {
		metadata = map[string]string{}
	}
	return services, metadata
}

func viewSession(session AgentVaultSession) AgentVaultSession {
	session.Services, session.Metadata = orEmpty(session.Services, session.Metadata)
	return session
}

// checkOnBehalfOf checks that a user to act on behalf of is named by exactly one of email address and ID. A user named
// by neither is refused rather than dropped, so an empty lookup result can't quietly create a session for the caller,
// with the caller's own access.
func checkOnBehalfOf(onBehalfOf *AgentVaultOnBehalfOf) error {
	if onBehalfOf == nil {
		return nil
	}
	if onBehalfOf.UserEmail != "" && onBehalfOf.UserID != "" {
		return errors.New("agent vault: set either UserEmail or UserID in OnBehalfOf, not both")
	}
	if onBehalfOf.UserEmail == "" && onBehalfOf.UserID == "" {
		return errors.New("agent vault: OnBehalfOf names no user: set UserEmail or UserID")
	}
	return nil
}

// Proxies

func (p *AgentVaultProxies) List(options AgentVaultListProxiesOptions) (AgentVaultListProxiesResult, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultListProxiesResult{}, err
	}

	res, err := api.CallListProxiesV1(p.client.httpClient, api.ListProxiesRequest{
		ProjectID:      options.ProjectID,
		Search:         options.Search,
		OrderBy:        withDefault(options.OrderBy, "createdAt"),
		OrderDirection: withDefault(options.OrderDirection, api.OrderDirectionDesc),
		Limit:          withDefault(options.Limit, 20),
		Offset:         options.Offset,
	})
	if err != nil {
		return AgentVaultListProxiesResult{}, err
	}

	for i, proxy := range res.Proxies {
		res.Proxies[i] = viewProxy(proxy)
	}
	return res, nil
}

func (p *AgentVaultProxies) Create(options AgentVaultCreateProxyOptions) (AgentVaultProxyWithToken, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultProxyWithToken{}, err
	}

	request := api.CreateProxyRequest{
		ProjectID:     options.ProjectID,
		Name:          options.Name,
		TrafficPolicy: toWireTrafficPolicy(withDefault(options.TrafficPolicy, api.TrafficPolicyAnyHost)),
		PollInterval:  withDefault(options.PollInterval, 60),
		Description:   options.Description,
	}
	// An empty value means no allowed hosts, so it's left out.
	if hosts, ok := joinHosts(options.AllowedHosts); ok {
		request.AllowedHosts = hosts
	}

	res, err := api.CallCreateProxyV1(p.client.onceHTTPClient(), request)
	if err != nil {
		return AgentVaultProxyWithToken{}, err
	}

	res.Proxy = viewProxy(res.Proxy)
	return res, nil
}

func (p *AgentVaultProxies) Update(proxyID string, options AgentVaultUpdateProxyOptions) (AgentVaultProxy, error) {
	request := api.UpdateProxyRequest{
		Name:          options.Name,
		TrafficPolicy: toWireTrafficPolicy(options.TrafficPolicy),
		PollInterval:  options.PollInterval,
		Description:   options.Description,
	}
	if options.AllowedHosts != nil {
		// An empty list removes every allowed host, which the API takes as null.
		if hosts, ok := joinHosts(options.AllowedHosts); ok {
			request.AllowedHosts = api.NullableOf(&hosts)
		} else {
			request.AllowedHosts = api.NullableOf[string](nil)
		}
	}

	res, err := api.CallUpdateProxyV1(p.client.httpClient, proxyID, request)
	if err != nil {
		return AgentVaultProxy{}, err
	}
	return viewProxy(res.Proxy), nil
}

func (p *AgentVaultProxies) Delete(proxyID string) (AgentVaultProxyRef, error) {
	res, err := api.CallDeleteProxyV1(p.client.onceHTTPClient(), proxyID)
	if err != nil {
		return AgentVaultProxyRef{}, err
	}
	return res.Proxy, nil
}

func (p *AgentVaultProxies) RotateToken(proxyID string) (AgentVaultProxyWithToken, error) {
	res, err := api.CallRotateProxyTokenV1(p.client.onceHTTPClient(), proxyID)
	if err != nil {
		return AgentVaultProxyWithToken{}, err
	}
	res.Proxy = viewProxy(res.Proxy)
	return res, nil
}

func (p *AgentVaultProxies) Revoke(proxyID string) (AgentVaultProxy, error) {
	res, err := api.CallRevokeProxyV1(p.client.httpClient, proxyID)
	if err != nil {
		return AgentVaultProxy{}, err
	}
	return viewProxy(res.Proxy), nil
}

// Sessions

func (s *AgentVaultSessions) List(options AgentVaultListSessionsOptions) (AgentVaultListSessionsResult, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultListSessionsResult{}, err
	}

	res, err := api.CallListSessionsV1(s.client.httpClient, api.ListSessionsRequest{
		ProjectID: options.ProjectID,
		Scope:     withDefault(options.Scope, api.SessionScopeMine),
		Statuses:  options.Statuses,
		Search:    options.Search,
		Metadata:  options.Metadata,
		Limit:     withDefault(options.Limit, 20),
		Offset:    options.Offset,
	})
	if err != nil {
		return AgentVaultListSessionsResult{}, err
	}

	for i, session := range res.Sessions {
		res.Sessions[i] = viewSession(session)
	}
	return res, nil
}

func (s *AgentVaultSessions) Get(sessionID string) (AgentVaultSession, error) {
	res, err := api.CallGetSessionV1(s.client.httpClient, sessionID)
	if err != nil {
		return AgentVaultSession{}, err
	}
	return viewSession(res.Session), nil
}

func (s *AgentVaultSessions) Create(options AgentVaultCreateSessionOptions) (AgentVaultCreatedSession, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultCreatedSession{}, err
	}
	if err := checkOnBehalfOf(options.OnBehalfOf); err != nil {
		return AgentVaultCreatedSession{}, err
	}

	var onBehalfOf *api.OnBehalfOfRequest
	if options.OnBehalfOf != nil {
		onBehalfOf = &api.OnBehalfOfRequest{UserEmail: options.OnBehalfOf.UserEmail, UserID: options.OnBehalfOf.UserID}
	}

	res, err := api.CallCreateSessionV1(s.client.onceHTTPClient(), api.CreateSessionRequest{
		ProjectID:    options.ProjectID,
		AccessBundle: options.AccessBundle,
		Services:     options.Services,
		TTL:          options.TTL,
		Label:        options.Label,
		Metadata:     options.Metadata,
		OnBehalfOf:   onBehalfOf,
	})
	if err != nil {
		return AgentVaultCreatedSession{}, err
	}

	session := res.Session
	session.Services, session.Metadata = orEmpty(session.Services, session.Metadata)
	return session, nil
}

func (s *AgentVaultSessions) Update(sessionID string, options AgentVaultUpdateSessionOptions) (AgentVaultSession, error) {
	res, err := api.CallUpdateSessionV1(s.client.httpClient, sessionID, api.UpdateSessionRequest{
		Label:    options.Label,
		Metadata: options.Metadata,
	})
	if err != nil {
		return AgentVaultSession{}, err
	}
	return viewSession(res.Session), nil
}

func (s *AgentVaultSessions) Revoke(sessionID string) (AgentVaultSession, error) {
	res, err := api.CallRevokeSessionV1(s.client.httpClient, sessionID)
	if err != nil {
		return AgentVaultSession{}, err
	}
	// Left as the API returns it: the response leaves out the access bundle and services, and an empty list would
	// read as a session that carries none.
	return res.Session, nil
}

// Access bundles

func (a *AgentVaultAccessBundles) List(options AgentVaultListAccessBundlesOptions) (AgentVaultListAccessBundlesResult, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultListAccessBundlesResult{}, err
	}
	if err := checkOnBehalfOf(options.OnBehalfOf); err != nil {
		return AgentVaultListAccessBundlesResult{}, err
	}

	request := api.ListAccessBundlesRequest{
		ProjectID:      options.ProjectID,
		Search:         options.Search,
		RequiresSetup:  options.RequiresSetup,
		Granted:        options.Granted,
		OrderBy:        withDefault(options.OrderBy, "createdAt"),
		OrderDirection: withDefault(options.OrderDirection, api.OrderDirectionDesc),
		Limit:          withDefault(options.Limit, 20),
		Offset:         options.Offset,
	}
	if options.OnBehalfOf != nil {
		request.OnBehalfOfUserEmail = options.OnBehalfOf.UserEmail
		request.OnBehalfOfUserID = options.OnBehalfOf.UserID
	}

	return api.CallListAccessBundlesV1(a.client.httpClient, request)
}

func (a *AgentVaultAccessBundles) Get(accessBundleID string) (AgentVaultAccessBundleWithServices, error) {
	res, err := api.CallGetAccessBundleV1(a.client.httpClient, accessBundleID)
	if err != nil {
		return AgentVaultAccessBundleWithServices{}, err
	}
	return res.AccessBundle, nil
}

func (a *AgentVaultAccessBundles) Create(options AgentVaultCreateAccessBundleOptions) (AgentVaultAccessBundle, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultAccessBundle{}, err
	}

	res, err := api.CallCreateAccessBundleV1(a.client.onceHTTPClient(), api.CreateAccessBundleRequest{
		ProjectID:   options.ProjectID,
		Name:        options.Name,
		Description: options.Description,
	})
	if err != nil {
		return AgentVaultAccessBundle{}, err
	}
	return res.AccessBundle, nil
}

func (a *AgentVaultAccessBundles) Update(accessBundleID string, options AgentVaultUpdateAccessBundleOptions) (AgentVaultAccessBundle, error) {
	res, err := api.CallUpdateAccessBundleV1(a.client.httpClient, accessBundleID, api.UpdateAccessBundleRequest{
		Name:        options.Name,
		Description: options.Description,
	})
	if err != nil {
		return AgentVaultAccessBundle{}, err
	}
	return res.AccessBundle, nil
}

func (a *AgentVaultAccessBundles) Delete(accessBundleID string) (AgentVaultAccessBundle, error) {
	res, err := api.CallDeleteAccessBundleV1(a.client.onceHTTPClient(), accessBundleID)
	if err != nil {
		return AgentVaultAccessBundle{}, err
	}
	return res.AccessBundle, nil
}

func granteesRequest(options AgentVaultListGranteesOptions) api.ListGranteesRequest {
	return api.ListGranteesRequest{
		ActorType: options.ActorType,
		Search:    options.Search,
		Limit:     withDefault(options.Limit, 20),
		Offset:    options.Offset,
	}
}

func (a *AgentVaultAccessBundles) ListGrantees(accessBundleID string, options AgentVaultListGranteesOptions) (AgentVaultListGranteesResult, error) {
	return api.CallListAccessBundleGranteesV1(a.client.httpClient, accessBundleID, granteesRequest(options))
}

func (a *AgentVaultAccessBundles) ListAvailableGrantees(accessBundleID string, options AgentVaultListGranteesOptions) (AgentVaultListMembersResult, error) {
	return api.CallListAvailableAccessBundleGranteesV1(a.client.httpClient, accessBundleID, granteesRequest(options))
}

func grantsRequest(options AgentVaultChangeGrantsOptions) api.ChangeGrantsRequest {
	return api.ChangeGrantsRequest{
		UserIDs:            options.UserIDs,
		MachineIdentityIDs: options.MachineIdentityIDs,
		GroupIDs:           options.GroupIDs,
	}
}

func (a *AgentVaultAccessBundles) Grant(accessBundleID string, options AgentVaultChangeGrantsOptions) (AgentVaultGrantResult, error) {
	return api.CallGrantAccessBundleV1(a.client.onceHTTPClient(), accessBundleID, grantsRequest(options))
}

func (a *AgentVaultAccessBundles) Revoke(accessBundleID string, options AgentVaultChangeGrantsOptions) (AgentVaultRevokeGrantsResult, error) {
	return api.CallRevokeAccessBundleV1(a.client.onceHTTPClient(), accessBundleID, grantsRequest(options))
}

// Delegations

func (d *AgentVaultDelegations) List(options AgentVaultListDelegationsOptions) (AgentVaultListDelegationsResult, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultListDelegationsResult{}, err
	}

	return api.CallListDelegationsV1(d.client.httpClient, api.ListDelegationsRequest{
		ProjectID:          options.ProjectID,
		MachineIdentityIDs: options.MachineIdentityIDs,
		Limit:              withDefault(options.Limit, 100),
		Offset:             options.Offset,
	})
}

func (d *AgentVaultDelegations) Get(machineIdentityID string, options AgentVaultProjectOptions) (AgentVaultDelegation, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultDelegation{}, err
	}

	res, err := api.CallGetDelegationV1(d.client.httpClient, machineIdentityID, options.ProjectID)
	if err != nil {
		return AgentVaultDelegation{}, err
	}
	return res.Delegation, nil
}

func (d *AgentVaultDelegations) Set(machineIdentityID string, options AgentVaultSetDelegationOptions) (AgentVaultDelegation, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultDelegation{}, err
	}

	res, err := api.CallSetDelegationV1(d.client.httpClient, machineIdentityID, options.ProjectID, api.SetDelegationRequest{
		Scope:    options.Scope,
		UserIDs:  options.UserIDs,
		GroupIDs: options.GroupIDs,
	})
	if err != nil {
		return AgentVaultDelegation{}, err
	}
	return res.Delegation, nil
}

func (d *AgentVaultDelegations) Delete(machineIdentityID string, options AgentVaultProjectOptions) (AgentVaultDelegation, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultDelegation{}, err
	}

	res, err := api.CallDeleteDelegationV1(d.client.onceHTTPClient(), machineIdentityID, options.ProjectID)
	if err != nil {
		return AgentVaultDelegation{}, err
	}
	return res.Delegation, nil
}

func (d *AgentVaultDelegations) AddSubjects(machineIdentityID string, options AgentVaultChangeDelegationSubjectsOptions) (AgentVaultDelegation, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultDelegation{}, err
	}

	res, err := api.CallAddDelegationSubjectsV1(d.client.httpClient, machineIdentityID, options.ProjectID, api.ChangeDelegationSubjectsRequest{
		UserIDs:  options.UserIDs,
		GroupIDs: options.GroupIDs,
	})
	if err != nil {
		return AgentVaultDelegation{}, err
	}
	return res.Delegation, nil
}

func (d *AgentVaultDelegations) RemoveSubjects(machineIdentityID string, options AgentVaultChangeDelegationSubjectsOptions) (*AgentVaultDelegation, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return nil, err
	}

	res, err := api.CallRemoveDelegationSubjectsV1(d.client.onceHTTPClient(), machineIdentityID, options.ProjectID, api.ChangeDelegationSubjectsRequest{
		UserIDs:  options.UserIDs,
		GroupIDs: options.GroupIDs,
	})
	if err != nil {
		return nil, err
	}
	return res.Delegation, nil
}

func (d *AgentVaultDelegations) ListAvailableSubjects(machineIdentityID string, options AgentVaultListAvailableDelegationSubjectsOptions) (AgentVaultListMembersResult, error) {
	if err := requireProject(options.ProjectID); err != nil {
		return AgentVaultListMembersResult{}, err
	}

	return api.CallListAvailableDelegationSubjectsV1(d.client.httpClient, machineIdentityID, api.ListAvailableSubjectsRequest{
		ProjectID: options.ProjectID,
		ActorType: options.ActorType,
		Search:    options.Search,
		Limit:     withDefault(options.Limit, 20),
		Offset:    options.Offset,
	})
}
