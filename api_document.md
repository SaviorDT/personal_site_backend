# Personal Site Backend API Documentation

## Authentication APIs

### POST /auth/register
**Description**: Register a new user account

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "username"
}
```

**Request Body Schema**:
- `email` (string, required): User's email address (must be valid email format)
- `password` (string, required): User's password (minimum 8 characters)
- `nickname` (string, required): User's display name

**Success Response (200)**:
```json
{
  "message": "User registered successfully",
  "user_id": 1
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
  ```
- `500 Internal Server Error`: Server error during registration
  ```json
  {
    "error": "Failed to create user"
  }
  ```

---

### POST /auth/login
**Description**: Login with email and password. Sets an `auth_token` HTTP-only cookie for the session.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Request Body Schema**:
- `email` (string, required): User's email address (must be valid email format)
- `password` (string, required): User's password (minimum 8 characters)

**Success Response (200)**:
```json
{
  "user_id": 1,
  "message": "Login successful",
  "role": "user",
  "nickname": "username",
}
```

**Response Schema**:
- `user_id` (uint): User's unique identifier
- `message` (string): Success message
- `role` (string): User's role (e.g., "user", "admin")
- `nickname` (string): User's display name

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'loginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
  ```
- `401 Unauthorized`: Invalid credentials
  ```json
  {
    "error": "Invalid email or password"
  }
  ```
- `500 Internal Server Error`: Server error during login
  ```json
  {
    "error": "Failed to generate token"
  }
  ```

---

### POST /auth/logout
**Description**: Logout user and clear authentication cookie. Removes the `auth_token` cookie.

**Success Response (200)**:
```json
{
  "message": "Logged out successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Server error during logout
  ```json
  {
    "error": "Failed to logout"
  }
  ```

---

### POST /auth/change-password
**Description**: Change user's password (requires login first)

**Request Body**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Request Body Schema**:
- `old_password` (string, required): Current password (minimum 8 characters)
- `new_password` (string, required): New password (minimum 8 characters)

**Success Response (200)**:
```json
{
  "message": "Password changed successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'changePasswordRequest.OldPassword' Error:Field validation for 'OldPassword' failed on the 'required' tag"
  }
  ```
- `401 Unauthorized`: Missing or invalid token, or incorrect old password
  ```json
  {
    "error": "Unauthorized"
  }
  ```
  ```json
  {
    "error": "Old password is incorrect"
  }
  ```
- `403 Forbidden`: Password change not allowed for this account type
  ```json
  {
    "error": "Password change only allowed for password-based accounts"
  }
  ```
- `500 Internal Server Error`: Server error during password change
  ```json
  {
    "error": "Failed to update password"
  }
  ```

### GET /auth/login-github
Description: Start GitHub OAuth login flow. Optionally accept a `redirect` query param to indicate where the browser should be redirected after a successful login.

Query Parameters:
- `redirect` (string, optional): A full URL to redirect to on success. This value is preserved via OAuth state and used in the callback.

Response:
- 302 Redirect to GitHub authorization URL

Notes:
- The server encodes a nonce and the `redirect` value into the OAuth `state` parameter.

### GET /auth/login-github-callback
Description: OAuth callback endpoint for GitHub. Exchanges the authorization code for a token, creates or finds the user, sets the `auth_token` HTTP-only cookie, and then redirects back to the provided `redirect` URL if present. If no `redirect` is provided, returns JSON.

Reads:
- `state` (from GitHub): contains an encoded object with fields `n` (nonce) and `r` (redirect URL).
- `redirect` (optional query): used only if not present in state.

On success with redirect present:
- 302 Redirect to the `redirect` URL with the following query parameters appended:
  - `login=success`
  - `user_id` (number)
  - `message` (string): "GitHub login successful"
  - `role` (string)
  - `nickname` (string)

On success without redirect:
- 200 JSON:
```json
{
  "message": "GitHub login successful",
  "user_id": 1,
  "role": "user",
  "nickname": "username"
}
```

Error Responses:
- `400 Bad Request`: Missing `code`
- `401 Unauthorized`: Code exchange failed
- `500 Internal Server Error`: OAuth not configured, DB or token errors

---

### GET /auth/login-google
Description: Start Google OAuth login flow. Optionally accept a `redirect` query param to indicate where the browser should be redirected after a successful login.

Query Parameters:
- `redirect` (string, optional): A full URL to redirect to on success. This value is preserved via OAuth state and used in the callback.

Response:
- 302 Redirect to Google authorization URL

Notes:
- The server encodes a nonce and the `redirect` value into the OAuth `state` parameter.

### GET /auth/login-google-callback
Description: OAuth callback endpoint for Google. Exchanges the authorization code, creates or finds the user, sets the `auth_token` HTTP-only cookie, and then redirects back to the provided `redirect` URL if present. If no `redirect` is provided, returns JSON.

Reads:
- `state` (from Google): contains an encoded object with fields `n` (nonce) and `r` (redirect URL).
- `redirect` (optional query): used only if not present in state.

On success with redirect present:
- 302 Redirect to the `redirect` URL with the following query parameters appended:
  - `login=success`
  - `user_id` (number)
  - `message` (string): "Google login successful"
  - `role` (string)
  - `nickname` (string)

On success without redirect:
- 200 JSON similar to GitHub callback.

Error Responses:
- `400 Bad Request`: Missing `code`
- `401 Unauthorized`: Code exchange failed
- `500 Internal Server Error`: OAuth not configured, DB or token errors

---

## Authentication Flow

1. **Register**: Create a new account using `/auth/register`. Or you don't need to do that if you use OAuth.
2. **Login**: Authenticate using `/auth/login` and you don't need to manage any thing about session. Or use `/auth/login-{3rd-platform}` to use OAuth login.
3. **Access Protected Resources**: token will saved in http only cookie
4. **Change Password**: Use `/auth/change-password` with valid authentication

---


## Error Handling

All endpoints return appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (validation errors)
- `401`: Unauthorized (authentication required or failed)
- `403`: Forbidden (insufficient permissions)
- `500`: Internal Server Error

## Notes

- All passwords must be at least 8 characters long
- Email addresses must be in valid email format
- Password changes are only allowed for accounts created with email/password (not OAuth accounts)