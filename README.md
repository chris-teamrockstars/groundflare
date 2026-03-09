
# Groundflare Cloudflare Pentest Proxy


## Concept

This Pentest Proxy can be used to transparently pentest
websites that are behind cloudflare's web application firewall.

If your pentest tools support a SOCKS proxy, you can point them
to this (locally or remotely running) proxy instance.


## Features

This is currently work in progress:

- [ ] Intercept NameResolver calls and store a JSON file-based database
- [ ] Resolve handler connects via TLS and looks for subject names to discover IPs
- [ ] If not available, check DNS resolvers for A/AAAA/CNAME entries
- [ ] If not available, check DNS dumpster for A/AAAA IP entries
- [ ] Intercept SOCKS Connect IPv4 calls to connect instead to Origin Server
- [ ] Intercept SOCKS Connect IPv6 calls to connect instead to Origin Server

