# auth 请求代码对照

| 文件 | 来源 | 写法 | 是否被页面使用 |
| --- | --- | --- | --- |
| `auth.ts` | 项目当前使用 | Axios 实例 + 拦截器 + `request<T>`，走 `/api` 代理 | 是 |
| `auth.apifox.ts` | Apifox 生成代码 | `fetch` + `then/catch`，直接请求 `http://127.0.0.1:8088/api`；login 保留原版不带 credentials，其余接口为演示补了 `credentials: 'include'` | 否 |
| `auth.axios.ts` | 手写对照版 | Axios + `async/await` + `try/catch`，`baseURL` 为 `/api`，走 Vite 代理；只 `console.log`，没有 `return` / `new Promise` / `Promise.reject` | 否 |

`auth.apifox.ts` 和 `auth.axios.ts` 只用于课堂对照，不参与页面逻辑；两份都只把结果 `console.log` 出来，不返回数据。

两个文件都实现了相同的四个接口：

- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/profile`
- `DELETE /auth/logout`

## 主要差异

- Apifox 版是“复制下来就能看”的请求代码，保留了 `var`、`Headers`、`requestOptions`、`fetch`、`.then()`、`.catch()`；它使用完整的 `http://127.0.0.1:8088/api` 地址，直连后端，不经过 Vite 代理。login 不带 `credentials`，直接跑只能看请求结果，不能保存跨域登录 Cookie。
- Axios 版把每个请求都包在 `try/catch` 里，成功 `console.log(response.data)`，失败 `console.error(...)`，不写 `return`、`return Promise`、`new Promise` 或 `Promise.reject`。
- Axios 版 `baseURL` 是 `/api`，请求会走 Vite 代理，浏览器看到的是同源的 5173；`axiosConfig` 里仍保留 `withCredentials: true`，以后改成直连后端或在生产环境跨域时仍然需要它。
- 如果要把 Apifox 版拿去实际登录，需要自己补 `credentials: 'include'`；这就是“原始代码”和“能用于项目”的差别。
