# Spec: Auth contract alignment and token lifecycle stabilization

## Context
The repository contains a BFF (`internal`), an auth gRPC service (`services/authsvc`), and a client service.
The authentication flow is only partially implemented and the public/internal contracts are inconsistent.

## Problem
The current auth implementation is not production-like because:
- gRPC contract and handler implementation do not fully match;
- BFF exposes only part of the auth flow;
- token metadata is not fully returned;
- refresh token transport is inconsistent;
- config keys are inconsistent across services.

## Goal
Make authentication flow internally consistent, testable, and safe enough to serve as the baseline for further production hardening.

## In scope
1. Align auth protobuf contract with actual domain model.
2. Fully implement the auth lifecycle:
   - register
   - login
   - refresh
   - logout current session
   - logout all sessions
   - validate access
3. Ensure token metadata is returned consistently.
4. Use one clear refresh token transport strategy.
5. Normalize config/env keys required by auth and client services.
6. Add automated tests covering the full auth lifecycle.

## Out of scope
- AWS deployment changes
- observability stack
- rate limiting
- mTLS between services
- password reset, email verification, MFA

## Functional requirements

### FR-1 Contract consistency
Auth gRPC service MUST implement every RPC declared in `contracts/proto/auth/v1/auth.proto`.
No declared RPC may remain unimplemented at runtime.

### FR-2 User ID consistency
User identifiers MUST use one consistent type across contracts and services.
Because domain user IDs are UUIDs, protobuf fields carrying user IDs MUST be string-based or removed if unnecessary.

### FR-3 TokenPair completeness
Successful register/login/refresh responses MUST include:
- access_token
- access_expires_at
- refresh_token or cookie-only transport decision
- refresh_expires_at
- session_id
Any returned metadata MUST be populated with real values.

### FR-4 Refresh token transport
The system MUST use one clear strategy for refresh token transport.
Baseline decision for production-like mode:
- refresh token is stored in HttpOnly cookie
- access token is returned in JSON body
- refresh token is not duplicated in JSON body

### FR-5 BFF auth surface
BFF MUST expose working HTTP endpoints for:
- POST /register
- POST /login
- POST /refresh
- POST /logout
- POST /logout-all

### FR-6 Error mapping
Known auth errors MUST be mapped to stable transport-level responses:
- invalid credentials
- invalid token
- expired/revoked session
- user already exists
Transport handlers MUST not return raw ambiguous errors where a domain-specific error exists.

### FR-7 Config consistency
All services MUST use the same env key names for DB SSL mode and service addresses.
Example: one key for DB SSL mode, one naming convention for service bind/listen addresses.

## Acceptance criteria
1. All RPCs in auth.proto are implemented and callable.
2. BFF can execute full register/login/refresh/logout/logout-all flow.
3. Token response metadata is fully populated.
4. Refresh token is not returned in both cookie and JSON.
5. Contract uses UUID-compatible user ID representation consistently.
6. Config/env naming is normalized across BFF, authsvc, clientsvc, and compose files.
7. Automated tests cover happy path and key failure cases.