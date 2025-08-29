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

## Storage APIs
**Description**:
```
All users can access its own storage and users who have not logged in share a storage.

Don't upload any thing you don't want to share whit admins here. Admins could see all your files.
```

### POST /storage/folder/*folder_path
**Description**: Create a new folder in the storage system

**Path Parameters**:
- `folder_path` (string, required): The folder path to create (supports nested paths, e.g., `/documents/2024/reports`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Directory created successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to create directory
  ```json
  {
    "error": "Failed to create directory"
  }
  ```

**Example**:
```bash
POST /storage/folder/documents/2024/reports
```

---

### GET /storage/folder/*folder_path
**Description**: List contents of a folder

**Path Parameters**:
- `folder_path` (string, required): The folder path to list (supports nested paths, e.g., `/documents/2024`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
[
    {
        "is_dir": true,
        "name": "test",
        "size": 4096
    },
    {
        "is_dir": false,
        "name": "test.txt",
        "size": 3
    }
]
```

**Response Schema**:
  - `name` (string): File or folder name
  - `is_dir` (bool): Whether it is a folder
  - `size` (number): File or folder size in bytes

**Error Responses**:
- `500 Internal Server Error`: Failed to list folder contents
  ```json
  {
    "error": "Failed to list folder contents"
  }
  ```

**Example**:
```bash
GET /storage/folder/documents/2024
```

---

### PATCH /storage/folder/*folder_path
**Description**: Update a folder

**Path Parameters**:
- `folder_path` (string, required): The current folder path to update

**Request Body**:
```json
{
  "path": "new folder path"
}
```

**Request Body Schema**:
- `path` (string, optional): The new folder path.


**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Folder updated successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to update folder
  ```json
  {
    "error": "Failed to update folder"
  }
  ```
- `500 Internal Server Error`: Failed to rename folder
  ```json
  {
    "error": "Failed to move folder"
  }
  ```

**Example**:
```bash
PATCH /storage/folder/documents/old_name
```

---

### DELETE /storage/folder/*folder_path
**Description**: Delete a folder and all its contents

**Path Parameters**:
- `folder_path` (string, required): The folder path to delete

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Folder deleted successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to delete folder
  ```json
  {
    "error": "Failed to delete folder"
  }
  ```

**Example**:
```bash
DELETE /storage/folder/documents/temp
```

---

### GET /storage/file/*file_path
**Description**: Download/retrieve a file from storage

**Path Parameters**:
- `file_path` (string, required): The file path to retrieve (supports nested paths, e.g., `/documents/2024/report.pdf`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
- Returns the file content with appropriate Content-Type header
- File is served directly for download or display

**Error Responses**:
- `400 Bad Request`: Invalid file path
  ```json
  {
    "error": "Cannot get file"
  }
  ```
- `404 Not Found`: File not found (served by web server)

**Example**:
```bash
GET /storage/file/documents/2024/report.pdf
```

---

### POST /storage/file/*file_path
**Description**: Upload a file to storage with chunked upload support for large files

**Path Parameters**:
- `file_path` (string, required): The destination file path (supports nested paths, e.g., `/documents/2024/report.pdf`)

**Form Data (multipart/form-data)**:
- `file_id` (string, required): Unique identifier for the file (used for chunked uploads)
- `chunk_index` (integer, required): Index of current chunk (0-based)
- `total_chunks` (integer, required): Total number of chunks for this file
- `chunk_data` (file, required): The file chunk data

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification
- `Content-Type`: multipart/form-data

**Success Response (200)**:
For non-final chunks:
```json
{
  "message": "Chunk uploaded successfully"
}
```

For final chunk (file complete):
```json
{
  "message": "File uploaded successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request parameters
  ```json
  {
    "error": "Cannot upload file"
  }
  ```
- `400 Bad Request`: Missing chunk data
  ```json
  {
    "error": "Missing chunk_data"
  }
  ```
- `500 Internal Server Error`: Upload failed
  ```json
  {
    "error": "Failed to save file"
  }
  ```

**Chunked Upload Process**:
1. Split large files into chunks (recommended: 1-10MB per chunk)
2. Upload each chunk with the same `file_id` and sequential `chunk_index`
3. Server automatically merges chunks when the final chunk is received
4. Temporary chunks are stored in `tmp/upload_chunks/{file_id}/` during upload

**Example**:
```bash
# Upload chunk 0 of 3
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=0
total_chunks=3
chunk_data=<binary_chunk_0>

# Upload chunk 1 of 3
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=1
total_chunks=3
chunk_data=<binary_chunk_1>

# Upload final chunk 2 of 3 (triggers merge)
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=2
total_chunks=3
chunk_data=<binary_chunk_2>
```

---

### PATCH /storage/file/*file_path
**Description**: Update/move a file to a new location

**Path Parameters**:
- `file_path` (string, required): The current file path to update

**Request Body**:
```json
{
  "path": "/new/location/filename.ext"
}
```

**Request Body Schema**:
- `path` (string, optional): New file path (relative to storage root)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "File updated successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request payload
  ```json
  {
    "error": "Invalid request payload"
  }
  ```
- `400 Bad Request`: Invalid new file path
  ```json
  {
    "error": "Invalid new file path"
  }
  ```
- `500 Internal Server Error`: Failed to update file
  ```json
  {
    "error": "Failed to update file"
  }
  ```
- `500 Internal Server Error`: Failed to move file
  ```json
  {
    "error": "Failed to move file"
  }
  ```

**Example**:
```bash
PATCH /storage/file/documents/old_name.txt
Content-Type: application/json

{
  "path": "/documents/renamed_file.txt"
}
```

---

### DELETE /storage/file/*file_path
**Description**: Delete a file from storage

**Path Parameters**:
- `file_path` (string, required): The file path to delete

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "File deleted successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid file path
  ```json
  {
    "error": "Cannot delete file"
  }
  ```
- `500 Internal Server Error`: Failed to delete file
  ```json
  {
    "error": "Failed to delete file"
  }
  ```

**Example**:
```bash
DELETE /storage/file/documents/unwanted_file.txt
```

---

## Storage Notes

- All folder and file paths support nested directory structures
- Authentication is optional for storage operations.
- User have their own storage if they logged in and they share a storage with other users if they did not log in.
- Folder and file names are case-sensitive
- The `*folder_path` and `*file_path` parameters capture the entire path after `/folder/` or `/file/`

## Battle Cat APIs

### GET /battle-cat/levels
**Description**: Filter and list Battle Cat levels by stage and up to 3 enemies. Returns an array of level collections. Internally, results may include matches for 3-enemy, 2-enemy, and 1-enemy combinations.

**Query Parameters**:
- `stage` (string, required, max length 3): Stage identifier to filter on
- `enemy` (string, required, repeated, max 3): Enemy names; pass 1 to 3 values
  - Examples: `?enemy=dog&enemy=snake`

**Success Response (200)**:
Returns an array of collections. Each collection contains the requested enemies echo and a list of levels.
```json
[
  {
    "enemies": ["enemyA", "enemyB", "enemyC"],
    "levels": [
      { "level": "001", "name": "Stage 1", "hp": 1200, "enemies": "enemyA, enemyX" },
      { "level": "002", "name": "Stage 2", "hp": 1500, "enemies": "enemyB, enemyY" }
    ]
  }
]
```

**Response Schema**:
- Array of objects with:
  - `enemies` (string[]): The enemies from the query (echoed)
  - `levels` (array): Matched levels for a particular enemy combination
    - `level` (string): Level code/id
    - `name` (string): Level name
    - `hp` (number): Level HP
    - `enemies` (string): Original enemies string from DB for that level

**Error Responses**:
- `401 Unauthorized`: Invalid query parameters (binding/validation failed)
  ```json
  { "error": "Key: 'FilterLevelsRequest.Enemies' Error:Field validation for 'Enemies' failed on the 'max' tag" }
  ```

**Examples**:
```bash
# Single enemy
GET /battle-cat/levels?stage=E01&enemy=dog

# Two enemies (repeat the query parameter)
GET /battle-cat/levels?stage=E01&enemy=dog&enemy=snake

# Three enemies
GET /battle-cat/levels?stage=E01&enemy=dog&enemy=snake&enemy=boss
```

Notes:
- The endpoint may return multiple collections representing matches for three-enemy, pairwise, and single-enemy filters.

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