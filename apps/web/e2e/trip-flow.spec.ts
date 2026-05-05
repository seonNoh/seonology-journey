import { test, expect } from '@playwright/test'

test.describe('Trip flow E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Stub auth – in real CI this would go through Keycloak login
    await page.goto('/')
  })

  test('homepage renders', async ({ page }) => {
    await expect(page.locator('body')).toBeVisible()
  })

  test('trip list page loads', async ({ page }) => {
    await page.goto('/trips')
    await expect(page).toHaveURL(/\/trips/)
  })

  test('skip-to-content link is present', async ({ page }) => {
    const skipLink = page.locator('a.skip-link')
    await expect(skipLink).toHaveAttribute('href', '#main-content')
  })

  test('create trip → add schedule → upload photo → share flow', async ({ page }) => {
    // This is a skeleton E2E scenario; requires running backend or MSW mocks
    await page.goto('/trips')

    // Trip creation (would need mocked API)
    // const createBtn = page.getByRole('button', { name: /새 여행|New Trip/i })
    // await createBtn.click()

    // For now just verify navigation works
    await page.goto('/trips/test-trip-id/share')
    await expect(page.locator('h1')).toBeVisible()
  })

  test('nearby page renders map picker', async ({ page }) => {
    await page.goto('/trips/test-trip-id/nearby')
    await expect(page.locator('h1')).toContainText(/주변 검색|Nearby/)
  })

  test('transit page has input fields', async ({ page }) => {
    await page.goto('/trips/test-trip-id/transit')
    await expect(page.locator('input')).toHaveCount(3) // origin, dest, time
  })
})
