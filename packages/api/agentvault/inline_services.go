package api

import "encoding/json"

// Inline services: services sent with a session, instead of added to an access bundle. Their values are sent exactly
// as given, since inline services can't use variables, so no value can contain text in double braces.
//
// Each type writes its own "type" field when it's sent, so callers can't send a service under the wrong type. The
// MarshalJSON methods marshal a copy under a local type with no methods, which keeps them from calling themselves.

// InlineService is a service sent with a session: an InlineHTTPService, an InlineMCPService, an
// InlinePostgresService, an InlineMySQLService or an InlineSSHService.
type InlineService interface {
	isInlineService()
}

// HTTPCredential is the credential of an InlineHTTPService: a BearerCredential, a BasicCredential or a
// PassthroughCredential.
type HTTPCredential interface {
	isHTTPCredential()
}

// SSHCredential is the credential of an InlineSSHService: a LoginCredential or an SSHKeyCredential.
type SSHCredential interface {
	isSSHCredential()
}

// BearerCredential is a token that the proxy writes to a header of each request.
type BearerCredential struct {
	// The token.
	Value string `json:"value"`
	// The header the token is written to. Defaults to "Authorization".
	HeaderName string `json:"headerName,omitempty"`
	// Written before the token, separated by one space. Defaults to "Bearer" for the Authorization header, and to
	// nothing for any other header. Point it at an empty string to send the token alone in the Authorization header.
	HeaderPrefix *string `json:"headerPrefix,omitempty"`
}

func (BearerCredential) isHTTPCredential() {}

func (c BearerCredential) MarshalJSON() ([]byte, error) {
	type fields BearerCredential
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{"bearer", fields(c)})
}

// BasicCredential is a username and password that the proxy sends to an HTTP service as HTTP basic authentication.
// Set a username, a password, or both.
type BasicCredential struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (BasicCredential) isHTTPCredential() {}

func (c BasicCredential) MarshalJSON() ([]byte, error) {
	type fields BasicCredential
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{"basic", fields(c)})
}

// LoginCredential is the account the proxy logs in to a database or SSH server as.
type LoginCredential struct {
	// Required.
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (LoginCredential) isSSHCredential() {}

func (c LoginCredential) MarshalJSON() ([]byte, error) {
	type fields LoginCredential
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{"basic", fields(c)})
}

// PassthroughCredential is no credential: the proxy lets requests to the service's hosts through without attaching
// one.
type PassthroughCredential struct{}

func (PassthroughCredential) isHTTPCredential() {}

func (PassthroughCredential) MarshalJSON() ([]byte, error) {
	return []byte(`{"type":"passthrough"}`), nil
}

// SSHKeyCredential is an SSH key the proxy logs in with, in OpenSSH or PEM format and without a passphrase.
type SSHKeyCredential struct {
	Username   string `json:"username"`
	PrivateKey string `json:"privateKey"`
	// An SSH certificate for the key: the contents of its -cert.pub file.
	Certificate string `json:"certificate,omitempty"`
}

func (SSHKeyCredential) isSSHCredential() {}

func (c SSHKeyCredential) MarshalJSON() ([]byte, error) {
	type fields SSHKeyCredential
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{"ssh-key", fields(c)})
}

// InlineCustomHeader is a header that the proxy attaches to each request, besides the credential.
type InlineCustomHeader struct {
	Name string `json:"name"`
	// Written before the value, separated by one space.
	Prefix string `json:"prefix,omitempty"`
	Value  string `json:"value"`
}

// InlineSubstitution is a placeholder that the proxy replaces with a real value before forwarding a request.
type InlineSubstitution struct {
	Placeholder string `json:"placeholder"`
	// The parts of each request where the proxy replaces the placeholder.
	Surfaces []SubstitutionSurface `json:"surfaces"`
	Value    string                `json:"value"`
}

// InlineHTTPService is an HTTP API sent with a session.
type InlineHTTPService struct {
	// A name for the service, in lowercase letters, numbers and hyphens. Without one, the service goes by its first
	// host.
	Name string `json:"name,omitempty"`
	// What the service is for, up to 256 characters. Agents in the session can read it, so never put a secret in it.
	Description string `json:"description,omitempty"`
	// The hosts the service covers, separated by commas, such as "api.github.com". A leading "*." matches exactly one
	// label.
	HostPattern string `json:"hostPattern"`
	// The HTTP methods the service allows, such as "GET". Leave it nil to allow every method. An empty, non-nil list
	// is sent as it is, rather than dropped, so it can't widen the service to every method.
	AllowedMethods []string `json:"allowedMethods,omitzero"`
	// The path prefixes the service allows, matched on whole segments. Leave it nil to allow every path. An empty,
	// non-nil list is sent as it is, rather than dropped, so it can't widen the service to every path.
	AllowedPathPrefixes []string       `json:"allowedPathPrefixes,omitzero"`
	Credential          HTTPCredential `json:"credential"`
	// Up to 20.
	CustomHeaders []InlineCustomHeader `json:"customHeaders,omitempty"`
	// Up to 20.
	Substitutions []InlineSubstitution `json:"substitutions,omitempty"`
}

func (InlineHTTPService) isInlineService() {}

func (s InlineHTTPService) MarshalJSON() ([]byte, error) {
	type fields InlineHTTPService
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{string(ServiceTypeHTTP), fields(s)})
}

// InlineMCPService is an MCP server sent with a session.
type InlineMCPService struct {
	// A name for the service, in lowercase letters, numbers and hyphens. Without one, the service goes by its host.
	Name string `json:"name,omitempty"`
	// What the service is for, up to 256 characters. Agents in the session can read it, so never put a secret in it.
	Description string `json:"description,omitempty"`
	// The MCP server's URL, such as "https://mcp.example.com/mcp". The service covers its host.
	ServerURL string `json:"serverUrl"`
	// Defaults to MCPTransportStreamableHTTP.
	Transport  MCPTransport     `json:"transport"`
	Credential BearerCredential `json:"credential"`
	// The MCP tools that agents can call through the service, by exact, case-sensitive name: each 1 to 128 letters,
	// numbers, underscores, hyphens or dots. The proxy refuses a call to any other tool and leaves it out of the
	// server's tool list. Leave it nil to allow every tool. An empty, non-nil list is sent as it is, rather than
	// dropped, so it can't widen the service to every tool. Needs an Infisical version that accepts a tool list on a
	// session's services; an older one refuses the session.
	AllowedTools []string `json:"allowedTools,omitzero"`
	// Up to 20.
	CustomHeaders []InlineCustomHeader `json:"customHeaders,omitempty"`
}

func (InlineMCPService) isInlineService() {}

func (s InlineMCPService) MarshalJSON() ([]byte, error) {
	if s.Transport == "" {
		s.Transport = MCPTransportStreamableHTTP
	}
	type fields InlineMCPService
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{string(ServiceTypeMCP), fields(s)})
}

// InlineTLS is how the proxy connects to a database. Each field defaults to true, and CA to the system's certificate
// authorities.
type InlineTLS struct {
	Enabled *bool `json:"enabled,omitempty"`
	// Whether the proxy checks the database's certificate.
	Verify *bool `json:"verify,omitempty"`
	// PEM certificate(s) of the authority that signed the database's certificate, used only while Verify is on.
	CA string `json:"ca,omitempty"`
}

// InlinePostgresService is a PostgreSQL database sent with a session. The agent reaches it with
// `infisical agent-vault connect`.
type InlinePostgresService struct {
	// The name the agent connects with, in lowercase letters, numbers and hyphens, up to 63 characters.
	Name string `json:"name"`
	// What the service is for, up to 256 characters. Agents in the session can read it, so never put a secret in it.
	Description string `json:"description,omitempty"`
	// The database server's host name or IP address, as the proxy reaches it. No scheme, port or path.
	Host string `json:"host"`
	// Defaults to 5432.
	Port int `json:"port,omitempty"`
	// The database the proxy logs in to.
	Database string     `json:"database"`
	TLS      *InlineTLS `json:"tls,omitempty"`
	// The account the proxy logs in as.
	Credential LoginCredential `json:"credential"`
}

func (InlinePostgresService) isInlineService() {}

func (s InlinePostgresService) MarshalJSON() ([]byte, error) {
	type fields InlinePostgresService
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{string(ServiceTypePostgres), fields(s)})
}

// InlineMySQLService is a MySQL database sent with a session. The agent reaches it with
// `infisical agent-vault connect`.
type InlineMySQLService struct {
	// The name the agent connects with, in lowercase letters, numbers and hyphens, up to 63 characters.
	Name string `json:"name"`
	// What the service is for, up to 256 characters. Agents in the session can read it, so never put a secret in it.
	Description string `json:"description,omitempty"`
	// The database server's host name or IP address, as the proxy reaches it. No scheme, port or path.
	Host string `json:"host"`
	// Defaults to 3306.
	Port int `json:"port,omitempty"`
	// The database the proxy logs in to.
	Database string     `json:"database"`
	TLS      *InlineTLS `json:"tls,omitempty"`
	// The account the proxy logs in as.
	Credential LoginCredential `json:"credential"`
}

func (InlineMySQLService) isInlineService() {}

func (s InlineMySQLService) MarshalJSON() ([]byte, error) {
	type fields InlineMySQLService
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{string(ServiceTypeMySQL), fields(s)})
}

// InlineSSHService is an SSH server sent with a session. The agent reaches it with `infisical agent-vault connect`.
type InlineSSHService struct {
	// The name the agent connects with, in lowercase letters, numbers and hyphens, up to 63 characters.
	Name string `json:"name"`
	// What the service is for, up to 256 characters. Agents in the session can read it, so never put a secret in it.
	Description string `json:"description,omitempty"`
	// The SSH server's host name or IP address, as the proxy reaches it. No scheme, port or path.
	Host string `json:"host"`
	// Defaults to 22.
	Port int `json:"port,omitempty"`
	// The user and password, or the user and SSH key, the proxy logs in with.
	Credential SSHCredential `json:"credential"`
}

func (InlineSSHService) isInlineService() {}

func (s InlineSSHService) MarshalJSON() ([]byte, error) {
	type fields InlineSSHService
	return json.Marshal(struct {
		Type string `json:"type"`
		fields
	}{string(ServiceTypeSSH), fields(s)})
}
