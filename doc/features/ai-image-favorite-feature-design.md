# 📌 AI 卡片网站 - 图像收藏功能技术方案

> **文档状态**：草案  
> **作者**：@hzycoder  
> **最后更新**：2026-03-31  
> **关联需求**：用户希望收藏自己生成的 AI 图像，便于后续查看或分享。

---

## 1️⃣ 功能目标

- ✅ 已登录用户可收藏/取消收藏自己生成的图像  
- ✅ 支持在个人中心查看所有收藏图像  
- ✅ 收藏数据持久化，支持高并发读写  
- ❌ 不支持收藏他人图像（初期版本）

---

## 2️⃣ 核心场景

| 场景 | 描述 |
|------|------|
| 收藏图像 | 用户点击“❤️收藏”按钮，将当前图像加入收藏夹 |
| 取消收藏 | 再次点击取消收藏 |
| 查看收藏列表 | 在 `/favorites` 页面按时间倒序展示 |

---

## 3️⃣ 技术选型

| 模块 | 选型 | 理由 |
|------|------|------|
| 后端语言 | Go 1.22+ | 高性能、简洁并发模型 |
| Web 框架 | Gin | 轻量、生态成熟 |
| ORM | GORM | 支持 PostgreSQL/MySQL，自动迁移 |
| 数据库 | PostgreSQL 15 | 原生 UUID、JSONB、强一致性 |
| 认证 | JWT | 无状态，适合前后端分离 |
| 部署 | Docker + Kubernetes | 与现有服务一致 |

---

## 4️⃣ 数据库设计

### 表结构

```sql
-- 收藏关系表（核心）
CREATE TABLE favorites (
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_id   UUID    NOT NULL REFERENCES generated_images(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, image_id)
);
```

> 💡 使用复合主键 `(user_id, image_id)` 天然去重，无需额外唯一索引。

---

## 5️⃣ API 接口设计

### `POST /api/favorites`
**收藏图像**

- **请求体**
  ```json
  { "image_id": "a1b2c3d4-..." }
  ```
- **权限**：需登录
- **响应**：`200 OK` 或 `404 Not Found`（图像不存在）

---

### `DELETE /api/favorites/{image_id}`
**取消收藏**

- **路径参数**：`image_id`
- **权限**：需登录
- **响应**：`200 OK` 或 `404`（未收藏）

---

### `GET /api/favorites`
**获取收藏列表**

- **响应示例**
  ```json
  {
    "data": [
      {
        "image_id": "a1b2...",
        "prompt": "夏日海滩，棕榈树",
        "image_url": "https://cdn/abc.jpg",
        "created_at": "2026-03-31T10:00:00Z"
      }
    ]
  }
  ```

---

## 6️⃣ 安全与校验

- 🔒 **必须验证**：`image_id` 是否属于当前用户（防止跨用户收藏）
- 🛡️ **防刷机制**：单用户每秒最多 3 次收藏操作（通过 Redis 限流）
- 🧪 **幂等性**：重复收藏同一图像应静默成功（使用 `ON CONFLICT DO NOTHING`）

---

## 7️⃣ 开发任务拆解（TODO）

| 任务 | 负责人 | 预估工时 | 状态 |
|------|--------|--------|------|
| 创建 `favorites` 表并迁移 | @dev1 | 1h | ⏳ |
| 实现收藏/取消 API | @dev1 | 3h | ⏳ |
| 添加收藏列表接口 | @dev1 | 2h | ⏳ |
| 前端对接收藏按钮 | @dev2 | 4h | ⏳ |
| 编写单元测试（Ginkgo） | @dev1 | 2h | ⏳ |
| 压测 & 限流配置 | @dev1 | 2h | ⏳ |

---

## 8️⃣ 验收标准

- [ ] 用户可成功收藏/取消自己的图像  
- [ ] 收藏列表按时间倒序展示  
- [ ] 并发 100 QPS 下无数据错乱  
- [ ] 未登录用户调用返回 `401`

---

## 9️⃣ 扩展性考虑（未来）

- 支持创建多个收藏夹（如“灵感”、“商用”）  
- 收藏图像支持添加标签  
- 提供“导出收藏”为 JSON 功能

---

> 📎 **附：相关文档链接**  
> - [用户认证方案](./auth-design.md)  
> - [图像生成 API 文档](./image-generation-api.md)
