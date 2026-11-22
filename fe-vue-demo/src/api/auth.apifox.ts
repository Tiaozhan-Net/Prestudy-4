/**
 * Apifox「JavaScript」请求代码原样版。
 *
 * 来源：Apifox 的“生成代码 → JavaScript”。
 * 请求部分保留了 var、Headers、requestOptions、fetch、then/catch 的原始结构；
 * 只把 raw 里写死的示例值换成了函数参数。
 *
 * 为了让代码能被调用，外面包了函数；但只 console.log，不 return。
 * login 保留 Apifox 原版，不带 credentials；register/profile/logout 为了课堂演示
 * 补上了 credentials: 'include'。实际项目请看走 /api 代理的 auth.axios.ts。
 */

const API_BASE_URL = 'http://127.0.0.1:8088/api'

interface LoginPayload {
  name: string
  password: string
}

export function login(name: string, password: string) {
  var myHeaders = new Headers()
  myHeaders.append('Content-Type', 'application/json')

  // 只改这里：把 Apifox 示例里的“小明 / abc123456”换成参数
  var raw = JSON.stringify({
    name,
    password,
  })

  var requestOptions = {
    method: 'POST',
    headers: myHeaders,
    body: raw,
    redirect: 'follow' as RequestRedirect,
  }

  fetch(`${API_BASE_URL}/auth/login`, requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.log('error', error))
}

export function register(name: string, password: string) {
  var myHeaders = new Headers()
  myHeaders.append('Content-Type', 'application/json')

  var raw = JSON.stringify({
    name,
    password,
  })

  var requestOptions = {
    method: 'POST',
    headers: myHeaders,
    body: raw,
    credentials: 'include' as RequestCredentials,
    redirect: 'follow' as RequestRedirect,
  }

  fetch(`${API_BASE_URL}/auth/register`, requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.log('error', error))
}

export function profile() {
  var requestOptions = {
    method: 'GET',
    credentials: 'include' as RequestCredentials,
    redirect: 'follow' as RequestRedirect,
  }

  fetch(`${API_BASE_URL}/auth/profile`, requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.log('error', error))
}

export function logout() {
  var requestOptions = {
    method: 'DELETE',
    credentials: 'include' as RequestCredentials,
    redirect: 'follow' as RequestRedirect,
  }

  fetch(`${API_BASE_URL}/auth/logout`, requestOptions)
    .then((response) => response.text())
    .then((result) => console.log(result))
    .catch((error) => console.log('error', error))
}

export const authApifox = {
  login,
  register,
  profile,
  logout,
}
