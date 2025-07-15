---
description: When using Remote environments, each user session connects to a remote SSH server acting as a proxy.
toc_depth: 5
---

# Remote environment

When using Remote environments, each user session establishes a connection to a remote SSH server, effectively acting as an SSH proxy. This allows Bifröst to serve as a gateway or jump host to other systems.

## Configuration {: #configuration}

<<property("type", "Environment Type", default="remote", required=True)>>
Has to be set to `remote` to enable the remote environment.

<<property("loginAllowed", "bool", template_context="../context/authorization.md", default=True)>>
Has to be true (after being evaluated) that the user is allowed to use this environment.

<<property("host", "string", template_context="../context/authorization.md", required=True)>>
The hostname or IP address of the remote SSH server to connect to.

<<property("port", "uint16", template_context="../context/authorization.md", default=22)>>
The port number of the remote SSH server.

<<property("user", "string", template_context="../context/authorization.md", default="{{.authorization.user.name}}")>>
The username to use when connecting to the remote SSH server.

<<property("password", "string", template_context="../context/authorization.md")>>
The password to use for authentication with the remote SSH server.

!!! warning
     Storing passwords in configuration files is not recommended for production use. Consider using SSH key-based authentication instead when it becomes available.

<<property("privateKey", "string", template_context="../context/authorization.md")>>
Path to private key file for SSH key-based authentication with the remote server. If provided, this will be used instead of password authentication.

!!! note
     This feature is a planned enhancement. Currently, only password authentication is supported.

<<property("banner", "string", template_context="../context/authorization.md", default="")>>
Will be displayed to the user upon connection to the remote environment.

<<property("portForwardingAllowed", "bool", template_context="../context/authorization.md", default=True)>>
Reserved for future use. Port forwarding is not currently implemented in the remote environment.

## Examples {: #examples}

1. Basic remote connection:
   ```yaml
   type: remote
   host: "remote-server.example.com"
   user: "{{.authorization.user.name}}"
   password: "secret-password"
   ```

2. Connect to specific port with custom user mapping:
   ```yaml
   type: remote
   host: "10.0.1.100"
   port: 2222
   user: "admin"
   password: "{{env `REMOTE_PASSWORD`}}"
   ```

3. With login restrictions:
   ```yaml
   type: remote
   host: "prod-server.internal"
   user: "{{.authorization.user.name}}"
   password: "{{env `PROD_PASSWORD`}}"
   loginAllowed: |
     {{ .authorization.user.groups | firstMatching `{{.name | eq "production-access"}}` }}
   ```
