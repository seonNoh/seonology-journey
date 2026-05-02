// Keycloak SSO 통합 (PKCE).
import Keycloak from 'keycloak-js'

const config = {
  url: import.meta.env.VITE_KEYCLOAK_URL ?? 'https://keycloak.seonology.com',
  realm: import.meta.env.VITE_KEYCLOAK_REALM ?? 'seonology',
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID ?? 'seonology-journey-web',
}

export const keycloak = new Keycloak(config)

// 한 번만 init.
let initPromise: Promise<boolean> | null = null

export function initKeycloak(): Promise<boolean> {
  if (!initPromise) {
    initPromise = keycloak.init({
      onLoad: 'check-sso',
      pkceMethod: 'S256',
      silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
      checkLoginIframe: false,
    })
  }
  return initPromise
}

export function login() {
  return keycloak.login({ redirectUri: window.location.href })
}

export function logout() {
  return keycloak.logout({ redirectUri: window.location.origin })
}

// 토큰을 자동 갱신해 만료 30초 전엔 새로 받음.
export async function getToken(): Promise<string | undefined> {
  if (!keycloak.authenticated) return undefined
  try {
    await keycloak.updateToken(30)
  } catch {
    return undefined
  }
  return keycloak.token
}
