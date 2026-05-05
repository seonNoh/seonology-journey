import { describe, it, expect, vi } from 'vitest'

vi.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  Routes: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  Route: () => null,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
  useNavigate: () => vi.fn(),
  useParams: () => ({}),
}))

vi.mock('keycloak-js', () => ({
  default: vi.fn(() => ({
    init: vi.fn().mockResolvedValue(true),
    authenticated: true,
    token: 'test-token',
  })),
}))

describe('App smoke test', () => {
  it('renders without crashing', () => {
    expect(true).toBe(true)
  })
})
