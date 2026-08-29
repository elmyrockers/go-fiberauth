# go-xvelope

[![Go Reference](https://pkg.go.dev/badge/github.com/elmyrockers/go-xvelope.svg)](https://pkg.go.dev/github.com/elmyrockers/go-xvelope)
[![Go Version](https://img.shields.io/badge/go1.27+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

<div align="center">
    <img src="img/logo.svg" width="900" />
</div><br>

Secure, stateless auth library for Go. AES-256-GCM encrypted cookies and opaque tokens inspired by ASP.NET Core Identity.

## Why xvelope? 

Most authentication libraries either force you into heavy server-side session stores (like Redis) or rely on standard JWTs that expose user claims in plain text to the client. 

**xvelope** takes a different approach by bringing the robust architecture of **ASP.NET Core Identity** to Go, focusing strictly on statelessness and data privacy.

### Core Concepts

* **Stateless by Design:** Your server doesn't need to track session state in memory or a database cache. The entire authentication context travels securely with the HTTP request.
* **Sealed, Unreadable Cookies:** Instead of base64-encoded JWTs, `xvelope` uses **AES-256-GCM Envelope Encryption**. The client receives a mathematically sealed cookie. If they try to decode it, they see nothing. If they tamper with it, the GCM auth tag instantly rejects it.
* **Opaque Bearer Tokens:** For API and mobile clients, `xvelope` issues opaque access and refresh tokens. Zero Personally Identifiable Information (PII) or internal role data is exposed to the frontend.
* **Familiar Identity Engine:** Built around the tried-and-true concepts of `UserManager`, `SignInManager`, and `Claims` from ASP.NET Core, giving you an enterprise-grade API tailored for Go.

## The Architecture

When a user logs in, `xvelope` processes the identity:
1. Gathers the user's `Claims` (ID, Email, Roles).
2. Serializes and compresses the payload.
3. Encrypts the payload using an AES-256-GCM master key.
4. Wraps the cipher in an `HttpOnly`, `Secure`, `SameSite=Strict` cookie or an Opaque Token pair.
5. On the next request, the middleware instantly decrypts and reconstructs the Go context—with zero database lookups.