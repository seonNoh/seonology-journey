// Keycloak SSO 통합 (PKCE).
//
// 세션 유지 정책:
//  - init 에서 `check-sso` 로 realm 세션이 살아있으면 silent 로그인
//  - onTokenExpired 콜백을 등록해 만료 직전/직후 자동 갱신
//  - 추가로 60 초마다 백그라운드에서 updateToken 을 호출해 realm
//    idle timeout 타이머를 갱신한다. 브라우저 탭이 장시간 열려있어도
//    refresh token 만 살아있으면 재로그인 없이 세션을 유지한다.
//  - updateToken 버퍼는 70 초. access token 이 기본 5 분일 때
//    유저 조작과 네트워크 왕복이 겹쳐도 만료되지 않는 여유를 확보.
//
// 실질적인 최대 세션 길이는 Keycloak realm 의 ssoSessionMaxLifespan /
// ssoSessionIdleTimeout 에 의해 결정된다. 앱에서는 refresh 호출을
// 게으르지 않게 하는 것까지만 책임진다.
import Keycloak from 'keycloak-js'

const config = {
  url: import.meta.env.VITE_KEYCLOAK_URL ?? 'https://auth.seonology.com',
  realm: import.meta.env.VITE_KEYCLOAK_REALM ?? 'seonology-journey',
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID ?? 'journey-web',
}

export const keycloak = new Keycloak(config)

// Buffer used when asking Keycloak for a new access token. If the token
// expires within this many seconds, updateToken refreshes. 70s keeps us
// comfortably ahead of network latency and clock skew.
const REFRESH_BUFFER_SECONDS = 70

// How often we poll updateToken while the tab is open. Short enough that
// realm idle timeouts (default 30 min) don't bite during active use.
const BACKGROUND_REFRESH_MS = 60_000

let initPromise: Promise<boolean> | null = null
let backgroundTimer: ReturnType<typeof setInterval> | null = null

export function initKeycloak(): Promise<boolean> {
  if (!initPromise) {
    initPromise = keycloak
      .init({
        onLoad: 'check-sso',
        pkceMethod: 'S256',
        silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
        // iframe session check disabled — we rely on updateToken instead.
        checkLoginIframe: false,
      })
      .then((authed) => {
        if (authed) setupAutoRefresh()
        return authed
      })
  }
  return initPromise
}

// Wire automatic refresh triggers. Safe to call multiple times; we clear
// any previous timer before starting a new one.
function setupAutoRefresh() {
  keycloak.onTokenExpired = () => {
    keycloak.updateToken(REFRESH_BUFFER_SECONDS).catch(() => {
      // Refresh token rejected — force re-login rather than silently
      // leaving the UI with an expired access token.
      login()
    })
  }

  if (backgroundTimer) clearInterval(backgroundTimer)
  backgroundTimer = setInterval(() => {
    if (!keycloak.authenticated) return
    keycloak.updateToken(REFRESH_BUFFER_SECONDS).catch(() => {
      // Swallow transient network errors. onTokenExpired will re-try
      // at actual expiry if the issue persists.
    })
  }, BACKGROUND_REFRESH_MS)

  // Proactively refresh when the user brings the tab back from sleep.
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible' && keycloak.authenticated) {
      keycloak.updateToken(REFRESH_BUFFER_SECONDS).catch(() => {})
    }
  })
}

export function login() {
  return keycloak.login({ redirectUri: window.location.href })
}

export function logout() {
  if (backgroundTimer) {
    clearInterval(backgroundTimer)
    backgroundTimer = null
  }
  return keycloak.logout({ redirectUri: window.location.origin })
}

// 토큰이 만료 버퍼 내에 있으면 새로 받음. 기존 호출부 호환성 유지.
export async function getToken(): Promise<string | undefined> {
  if (!keycloak.authenticated) return undefined
  try {
    await keycloak.updateToken(REFRESH_BUFFER_SECONDS)
  } catch {
    return undefined
  }
  return keycloak.token
}
