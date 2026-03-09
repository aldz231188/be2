# Plan: Auth contract alignment and token lifecycle stabilization

## Design decisions
1. Keep UUID as the canonical user ID type.
2. Update protobuf contract where it conflicts with UUID-based domain identity.
3. Treat authsvc gRPC contract as the source of truth for BFF auth client integration.
4. Use cookie-based refresh token transport in BFF for production-like mode.
5. Return only access token in JSON responses from BFF.
6. Keep current JWT + DB-backed session model.
7. Do not change AWS/infrastructure in this iteration.

## Technical changes

### Contracts
- Review `contracts/proto/auth/v1/auth.proto`
- Fix inconsistent `user_id` representation
- Regenerate protobuf stubs

### Auth service
- Implement Login
- Rename/align Refresh and LogoutAll handlers to actual gRPC interface
- Populate full TokenPair fields in app layer
- Ensure stable domain-to-transport error mapping

### BFF
- Implement HTTP handlers/routes for login, refresh, logout-all
- Extend auth gRPC client with missing methods
- Return access token payload only
- Set/clear refresh cookie consistently

### Config
- Normalize env names:
  - DB SSL mode
  - auth/client service address keys

### Tests
- Add/adjust tests for:
  - register success/conflict
  - login success/invalid credentials
  - refresh success/revoked session
  - logout current session
  - logout all sessions
  - validate access
