# Table of Contents:
- [Auth logic](#auth-logic)
- [C4 User Service component diagram](#c4-user-service-component-diagram)
- [C4 User Service code diagram](#c4-user-service-code-diagram)

### Auth logic:
During registration and login, the user receives a JWT access token and a JWT refresh token. 
When the access token expires (15 minutes), the user can send a request to `/refresh-token` to obtain a new pair 
of tokens using the refresh token. The refresh token is valid for one week. If the user does not use the application 
for a period exceeding the refresh token's validity, the application will require them to log in again.

### C4 User Service component diagram:
![System Architecture](./Component_CleverVillageSystem_UserService.svg)

### C4 User Service code diagram:
![System Architecture](./Code_CleverVillageSystem_UserService.svg)
