# 發文、留言、心情系統 API 文檔

## 📝 文章 API

### 1. 建立文章
**POST** `/posts`

**需要登入**: ✅

**Request Body**:
```json
{
  "title": "我的第一篇文章",
  "content": "這是文章的完整內容...",
  "summary": "這是文章摘要",
  "cover_image": "https://example.com/image.jpg",
  "status": "published",
  "visibility": "public",
  "tags": ["技術", "Go語言", "教學"]
}
```

**欄位說明**:
- `title` (string, required): 文章標題，1-255 字元
- `content` (string, required): 文章內容
- `summary` (string, optional): 文章摘要，最多 500 字元
- `cover_image` (string, optional): 封面圖片 URL
- `status` (string, optional): 狀態，可選 `draft`、`published`、`archived`，預設 `draft`
- `visibility` (string, optional): 可見性，可選 `public`、`private`、`friends`，預設 `public`
- `tags` (array, optional): 標籤陣列

**Response (201)**:
```json
{
  "message": "Post created successfully",
  "post": {
    "ID": 1,
    "CreatedAt": "2025-11-07T10:00:00Z",
    "UpdatedAt": "2025-11-07T10:00:00Z",
    "title": "我的第一篇文章",
    "content": "這是文章的完整內容...",
    "status": "published",
    "visibility": "public",
    "view_count": 0,
    "author": {
      "ID": 1,
      "nickname": "user123"
    },
    "tags": [
      {"ID": 1, "name": "技術", "color": "#3B82F6"},
      {"ID": 2, "name": "Go語言", "color": "#3B82F6"}
    ]
  }
}
```

---

### 2. 取得文章列表
**GET** `/posts`

**需要登入**: ❌（選擇性）

**Query Parameters**:
- `page` (int, optional): 頁碼，預設 1
- `page_size` (int, optional): 每頁數量，預設 10，最多 100
- `status` (string, optional): 篩選狀態
- `tag` (string, optional): 篩選標籤

**範例**: `/posts?page=1&page_size=10&tag=技術`

**Response (200)**:
```json
{
  "posts": [
    {
      "ID": 1,
      "title": "我的第一篇文章",
      "summary": "這是文章摘要",
      "author": {
        "ID": 1,
        "nickname": "user123"
      },
      "tags": [...],
      "view_count": 42,
      "created_at": "2025-11-07T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

### 3. 取得單一文章
**GET** `/posts/:id`

**需要登入**: ❌（選擇性，私人文章需要）

**Response (200)**:
```json
{
  "post": {
    "ID": 1,
    "title": "我的第一篇文章",
    "content": "完整內容...",
    "author": {...},
    "tags": [...],
    "comments": [...],
    "reactions": [...],
    "view_count": 43
  }
}
```

---

### 4. 更新文章
**PUT** `/posts/:id`

**需要登入**: ✅（僅作者）

**Request Body**: 與建立文章相同，所有欄位都是選擇性的

---

### 5. 刪除文章
**DELETE** `/posts/:id`

**需要登入**: ✅（僅作者或管理員）

**Response (200)**:
```json
{
  "message": "Post deleted successfully"
}
```

---

## 💬 留言 API

### 1. 建立留言
**POST** `/posts/:id/comments`

**需要登入**: ✅

**Request Body**:
```json
{
  "content": "這是我的留言內容",
  "parent_id": null
}
```

**欄位說明**:
- `content` (string, required): 留言內容
- `parent_id` (uint, optional): 父留言 ID（用於回覆）

**Response (201)**:
```json
{
  "message": "Comment created successfully",
  "comment": {
    "ID": 1,
    "post_id": 1,
    "content": "這是我的留言內容",
    "author": {
      "ID": 1,
      "nickname": "user123"
    },
    "parent_id": null,
    "created_at": "2025-11-07T10:00:00Z"
  }
}
```

---

### 2. 取得文章的所有留言
**GET** `/posts/:id/comments`

**需要登入**: ❌

**Response (200)**:
```json
{
  "comments": [
    {
      "ID": 1,
      "content": "這是留言",
      "author": {...},
      "replies": [
        {
          "ID": 2,
          "content": "這是回覆",
          "author": {...},
          "parent_id": 1
        }
      ],
      "reactions": [...],
      "is_edited": false
    }
  ]
}
```

---

### 3. 更新留言
**PUT** `/comments/:comment_id`

**需要登入**: ✅（僅作者）

**Request Body**:
```json
{
  "content": "更新後的留言內容"
}
```

---

### 4. 刪除留言（軟刪除）
**DELETE** `/comments/:comment_id`

**需要登入**: ✅（僅作者或管理員）

**說明**: 留言會被標記為已刪除，內容改為 "[此留言已刪除]"，但結構保留（為了維持對話串的完整性）

---

## ❤️ 反應/心情 API

### 1. 對文章新增反應
**POST** `/posts/:id/reactions`

**需要登入**: ✅

**Request Body**:
```json
{
  "type": "like"
}
```

**反應類型**:
- `like` - 讚 👍
- `love` - 愛心 ❤️
- `haha` - 哈哈 😆
- `wow` - 驚訝 😮
- `sad` - 難過 😢
- `angry` - 生氣 😠
- `care` - 關心 🤗

**行為說明**:
- 第一次：新增反應
- 相同類型：取消反應
- 不同類型：更新反應

**Response (201/200)**:
```json
{
  "message": "Reaction added",
  "reaction": {
    "ID": 1,
    "user_id": 1,
    "post_id": 1,
    "type": "like"
  }
}
```

---

### 2. 對留言新增反應
**POST** `/comments/:comment_id/reactions`

**需要登入**: ✅

**Request Body**: 與文章反應相同

---

### 3. 取得文章的反應統計
**GET** `/posts/:id/reactions`

**需要登入**: ❌（選擇性）

**Response (200)**:
```json
{
  "reactions": [
    {"type": "like", "count": 15},
    {"type": "love", "count": 8},
    {"type": "haha", "count": 3}
  ],
  "user_reaction": {
    "ID": 1,
    "type": "like"
  }
}
```

**說明**: `user_reaction` 僅在登入時返回，顯示當前使用者的反應

---

### 4. 取得留言的反應統計
**GET** `/comments/:comment_id/reactions`

**需要登入**: ❌（選擇性）

**Response**: 與文章反應統計相同

---

## 🔒 權限說明

### 文章權限
- **public**: 所有人可見
- **private**: 僅作者可見
- **friends**: 僅好友可見（需實作好友系統）

### 操作權限
- 建立：需登入
- 閱讀：依文章可見性
- 更新：僅作者
- 刪除：作者或管理員

---

## 📊 資料庫關係

```
User (使用者)
  ├─ has many Posts (文章)
  ├─ has many Comments (留言)
  └─ has many Reactions (反應)

Post (文章)
  ├─ belongs to User (作者)
  ├─ has many Tags (標籤) - 多對多
  ├─ has many Comments (留言)
  └─ has many Reactions (反應)

Comment (留言)
  ├─ belongs to Post (文章)
  ├─ belongs to User (作者)
  ├─ has many Replies (回覆) - 自關聯
  └─ has many Reactions (反應)

Reaction (反應)
  ├─ belongs to User
  └─ belongs to Post OR Comment (擇一)

Tag (標籤)
  └─ has many Posts - 多對多
```

---

## 🚀 使用範例

### 完整發文流程
```bash
# 1. 登入
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  -c cookies.txt

# 2. 建立文章
curl -X POST http://localhost/api/posts \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "title":"Go 語言入門",
    "content":"這是一篇關於 Go 的教學...",
    "status":"published",
    "tags":["Go","教學"]
  }'

# 3. 對文章按讚
curl -X POST http://localhost/api/posts/1/reactions \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"type":"like"}'

# 4. 留言
curl -X POST http://localhost/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"content":"寫得很棒！"}'
```

---

## ⚠️ 錯誤處理

所有 API 在發生錯誤時會返回適當的 HTTP 狀態碼和錯誤訊息：

```json
{
  "error": "錯誤描述"
}
```

常見狀態碼：
- `400` - 請求格式錯誤
- `401` - 未授權（需要登入）
- `403` - 禁止訪問（權限不足）
- `404` - 資源不存在
- `500` - 伺服器錯誤
