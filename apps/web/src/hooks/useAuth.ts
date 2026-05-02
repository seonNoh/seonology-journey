import { useEffect, useState } from 'react'
import { initKeycloak, keycloak, login, logout } from '../lib/keycloak'

export interface AuthState {
  ready: boolean
  authenticated: boolean
  username?: string
}

export function useAuth(): AuthState & { login: () => void; logout: () => void } {
  const [state, setState] = useState<AuthState>({ ready: false, authenticated: false })

  useEffect(() => {
    let cancelled = false
    initKeycloak()
      .then((authd) => {
        if (cancelled) return
        setState({
          ready: true,
          authenticated: authd,
          username:
            keycloak.tokenParsed && typeof keycloak.tokenParsed === 'object'
              ? ((keycloak.tokenParsed as Record<string, unknown>).preferred_username as string | undefined)
              : undefined,
        })
      })
      .catch(() => {
        if (cancelled) return
        setState({ ready: true, authenticated: false })
      })
    return () => {
      cancelled = true
    }
  }, [])

  return { ...state, login, logout }
}
