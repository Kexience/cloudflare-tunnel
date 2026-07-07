# AGENTS.md — 前端项目

## 项目概述

Cloudflare Tunnel 管理系统的前端部分，基于 Svelte 5 + TypeScript + Tailwind CSS 构建，提供用户认证和隧道管理的可视化界面。

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Svelte | 5.x (Runes 语法) | UI 框架 |
| TypeScript | 6.x | 类型系统 |
| Vite | 8.x | 构建工具 |
| Tailwind CSS | 4.x | 样式框架 |
| Axios | 1.x | HTTP 客户端 |
| svelte-routing | 2.x | 客户端路由 |
| Bun | - | 包管理器 (bun.lock) |

## 开发命令

```bash
# 安装依赖
bun install

# 启动开发服务器
bun run dev

# 构建生产版本
bun run build

# 预览生产构建
bun run preview

# 类型检查 + Svelte 检查
bun run check
```

## 目录结构

```
frontend/
├── src/
│   ├── main.ts              # 应用入口，挂载 App 组件
│   ├── App.svelte           # 根组件，路由配置
│   ├── app.css              # 全局样式 (Tailwind 导入)
│   ├── lib/
│   │   ├── api.ts           # Axios 实例，请求/响应拦截器
│   │   └── auth/
│   │       ├── index.ts     # 模块导出
│   │       ├── types.ts     # 类型定义 (User, AuthState, ApiResponse 等)
│   │       ├── api.ts       # 认证 API 调用 (login, register, getCurrentUser)
│   │       ├── store.ts     # Svelte store (authStore, isAuthenticated 等)
│   │       └── ProtectedRoute.svelte  # 路由守卫组件
│   └── pages/
│       ├── Login.svelte     # 登录页
│       ├── Register.svelte  # 注册页
│       └── Dashboard.svelte # 仪表盘页 (隧道管理主页)
├── public/                  # 静态资源
├── index.html               # HTML 入口
├── vite.config.ts           # Vite 配置 (含 API 代理)
├── svelte.config.js         # Svelte 配置
├── tsconfig.json            # TS 配置 (项目引用)
├── tsconfig.app.json        # 应用 TS 配置
├── tsconfig.node.json       # Node 环境 TS 配置
└── .env.example             # 环境变量示例
```

## 架构要点

### Svelte 5 Runes 语法

项目使用 Svelte 5 的 Runes 响应式语法，**不要使用旧版 Svelte 4 语法**：

```svelte
<!-- ✅ 正确：Svelte 5 Runes -->
<script lang="ts">
  let count = $state(0)
  let doubled = $derived(count * 2)
  let { title }: { title: string } = $props()
</script>

<!-- ❌ 错误：Svelte 4 旧语法 -->
<script lang="ts">
  let count = 0           // 不要这样做
  $: doubled = count * 2  // 不要这样做
  export let title        // 不要这样做
</script>
```

### 认证流程

1. **Token 存储**: JWT token 存储在 `localStorage`
2. **请求拦截**: `src/lib/api.ts` 的 Axios 拦截器自动在请求头添加 `Authorization: Bearer <token>`
3. **响应拦截**: 401 响应自动清除 token 并跳转登录页
4. **路由守卫**: `ProtectedRoute.svelte` 组件验证 token 有效性，无效则重定向到 `/login`
5. **状态管理**: `authStore` 使用 Svelte writable store 管理认证状态

### API 代理配置

开发模式下，Vite 将 `/api` 请求代理到后端：

```
前端请求: /api/v1/user/login
      ↓ Vite proxy (rewrite 去掉 /api 前缀)
后端接收: /v1/user/login → http://localhost:8083
```

### 路由结构

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | Loading | 自动检测登录状态并重定向 |
| `/login` | Login | 登录页 |
| `/register` | Register | 注册页 |
| `/dashboard` | Dashboard | 仪表盘 (需认证) |

### API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/user/register` | 用户注册 |
| POST | `/v1/user/login` | 用户登录，返回 token + user |
| GET | `/v1/user/me` | 获取当前用户信息 |

### 类型定义

```typescript
interface ApiResponse<T> {
  code: number      // 0 表示成功
  message: string
  data: T | null
}

interface User {
  id: number
  nickname: string
  username: string
  email: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `VITE_API_BASE_URL` | `/api` | API 基础路径 |

## 编码规范

1. **组件文件**: 使用 PascalCase 命名 (如 `Login.svelte`, `ProtectedRoute.svelte`)
2. **TypeScript**: 严格模式，禁止使用 `any`，所有类型必须明确声明
3. **样式**: 使用 Tailwind CSS 工具类，不写自定义 CSS (除非必要)
4. **状态管理**: 使用 Svelte store (`writable`/`derived`)，不引入外部状态库
5. **API 调用**: 统一通过 `src/lib/api.ts` 的 Axios 实例，不要直接使用 fetch
6. **路由**: 使用 `svelte-routing` 的 `navigate()` 函数进行编程式导航
7. **错误处理**: API 错误统一通过 `authStore.error` 展示，UI 使用中文提示

## 当前开发状态

已完成:
- [x] 项目基础搭建 (Vite + Svelte 5 + TypeScript + Tailwind)
- [x] 用户认证系统 (登录/注册/JWT)
- [x] 路由守卫和认证状态管理
- [x] 登录页、注册页、仪表盘基础 UI

待开发:
- [ ] 隧道 CRUD 管理界面
- [ ] 隧道状态实时监控
- [ ] 流量统计图表
- [ ] 日志查看功能
- [ ] 用户个人设置

## 注意事项

- 根目录 `AGENTS.md` 中写的前端技术栈为 "React 18+" 是**过时的**，实际使用的是 **Svelte 5**
- 包管理器使用 **bun**，不是 npm/yarn/pnpm
- Svelte 5 的 `$state`, `$derived`, `$props`, `$effect` 是编译器魔法，不是普通 JS 函数
- 构建产物输出到 `dist/` 目录
