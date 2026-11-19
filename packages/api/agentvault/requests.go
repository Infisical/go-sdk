package api

import "encoding/json"

// Requests, as sent to the API.

// Nullable is a request value that's left out while unset, sent as null when set to nil, and sent as itself
// otherwise. Its zero value is unset, so a field of this type is tagged omitzero.
type Nullable[T any] struct {
	set   bool
	value *T
}

// NullableOf returns a Nullable set to value, which may be nil.
func NullableOf[T any](value *T) Nullable[T] {
	return Nullable[T]{set: true, value: value}
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

type ListProxiesRequest struct {
	ProjectID      string
	Search         string
	OrderBy        string
	OrderDirection OrderDirection
	Limit          int
	Offset         int
}

type CreateProxyRequest struct {
	ProjectID     string        `json:"projectId"`
	Name          string        `json:"name"`
	TrafficPolicy TrafficPolicy `json:"trafficPolicy"`
	AllowedHosts  string        `json:"allowedHosts,omitempty"`
	PollInterval  int           `json:"pollInterval"`
	Description   string        `json:"description,omitempty"`
}

type UpdateProxyRequest struct {
	Name          *string       `json:"name,omitempty"`
	TrafficPolicy TrafficPolicy `json:"trafficPolicy,omitempty"`
	// Set to nil to remove every allowed host.
	AllowedHosts Nullable[string] `json:"allowedHosts,omitzero"`
	PollInterval int              `json:"pollInterval,omitempty"`
	Description  *string          `json:"description,omitempty"`
}

type ListSessionsRequest struct {
	ProjectID string
	Scope     SessionScope
	Statuses  []SessionStatus
	Search    string
	Metadata  map[string]string
	Limit     int
	Offset    int
}

// OnBehalfOfRequest names a user by email address or by ID, not both.
type OnBehalfOfRequest struct {
	UserEmail string `json:"userEmail,omitempty"`
	UserID    string `json:"userId,omitempty"`
}

type CreateSessionRequest struct {
	ProjectID string `json:"projectId"`
	// Left out without an access bundle, for a session that carries only services of its own.
	AccessBundle string          `json:"accessBundle,omitempty"`
	Services     []InlineService `json:"services,omitempty"`
	// Left out unless set, so the server picks the default: 7d, or 1h on behalf of a user.
	TTL        string             `json:"ttl,omitempty"`
	Label      string             `json:"label,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
	OnBehalfOf *OnBehalfOfRequest `json:"onBehalfOf,omitempty"`
}

type UpdateSessionRequest struct {
	// An empty string removes the label.
	Label *string `json:"label,omitempty"`
	// An empty, non-nil map removes the metadata, so it's sent rather than dropped.
	Metadata map[string]string `json:"metadata,omitzero"`
}

type ListAccessBundlesRequest struct {
	ProjectID           string
	Search              string
	RequiresSetup       *bool
	Granted             *bool
	OnBehalfOfUserEmail string
	OnBehalfOfUserID    string
	OrderBy             string
	OrderDirection      OrderDirection
	Limit               int
	Offset              int
}

type CreateAccessBundleRequest struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateAccessBundleRequest struct {
	Name *string `json:"name,omitempty"`
	// An empty string removes the description.
	Description *string `json:"description,omitempty"`
}

type ListGranteesRequest struct {
	ActorType MemberType
	Search    string
	Limit     int
	Offset    int
}

type ChangeGrantsRequest struct {
	UserIDs            []string `json:"userIds,omitempty"`
	MachineIdentityIDs []string `json:"machineIdentityIds,omitempty"`
	GroupIDs           []string `json:"groupIds,omitempty"`
}

type ListDelegationsRequest struct {
	ProjectID          string
	MachineIdentityIDs []string
	Limit              int
	Offset             int
}

type SetDelegationRequest struct {
	Scope    DelegationScope `json:"scope"`
	UserIDs  []string        `json:"userIds,omitempty"`
	GroupIDs []string        `json:"groupIds,omitempty"`
}

type ChangeDelegationSubjectsRequest struct {
	UserIDs  []string `json:"userIds,omitempty"`
	GroupIDs []string `json:"groupIds,omitempty"`
}

type ListAvailableSubjectsRequest struct {
	ProjectID string
	ActorType MemberType
	Search    string
	Limit     int
	Offset    int
}
