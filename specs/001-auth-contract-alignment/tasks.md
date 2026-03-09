# Tasks: Auth contract alignment and token lifecycle stabilization

## Contracts
- [ ] Audit `contracts/proto/auth/v1/auth.proto`
- [ ] Replace/remove UUID-incompatible `user_id int64` fields
- [ ] Regenerate protobuf code
- [ ] Verify generated client/server interfaces compile against authsvc and BFF

## Auth service
- [ ] Implement `Login(ctx, *LoginRequest) (*LoginResponse, error)`
- [ ] Implement `Refresh(ctx, *RefreshRequest) (*RefreshResponse, error)`
- [ ] Implement `LogoutAll(ctx, *LogoutAllRequest) (*emptypb.Empty, error)`
- [ ] Ensure `Register`, `Logout`, `ValidateAccess` match final contract
- [ ] Populate `AccessExpiresAt`, `RefreshExpiresAt`, `SessionId`
- [ ] Decide whether authsvc returns refresh token directly or leaves final transport policy to BFF
- [ ] Replace raw transport errors with mapped auth/domain errors

## BFF
- [ ] Add missing auth client methods: Login, Refresh, LogoutAll
- [ ] Add routes: POST /login, POST /refresh, POST /logout-all
- [ ] Implement matching HTTP handlers
- [ ] Return access token in JSON body only
- [ ] Set refresh token in HttpOnly cookie
- [ ] Clear refresh cookie on logout/logout-all

## Config
- [ ] Normalize DB SSL env key usage across services and compose files
- [ ] Normalize service address env key usage across services and compose files
- [ ] Update `.env.example`

## Tests
- [ ] Add authsvc tests for token metadata population
- [ ] Add handler/client tests for full auth lifecycle
- [ ] Add regression test proving declared RPCs are implemented
