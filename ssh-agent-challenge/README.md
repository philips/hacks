# ssh-agent OAUTH 2.0 bridge

Nearly everyone who does server infrastructure has an ssh key. This is a strong
public-key based identity. Not only that but it is often the only identity that
is available when people remote into a host via ssh thanks to "ssh-agent
forwarding".

This tool acts a bridge so that people can use their ssh-identity to back
another identity, initially an identity backed by dex, an OAUTH 2.0 / OIDC
provider.

## Manual

### List Public Keys

ssh-agent-tool list-keys

### Challenge Sign

ssh-agent-tool sign

### Challenge Verify
