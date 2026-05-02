import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { AppShell } from './components/AppShell'
import { HomePage } from './pages/HomePage'
import { TripListPage } from './pages/TripListPage'
import { TripDetailPage } from './pages/TripDetailPage'
import { DayDetailPage } from './pages/DayDetailPage'

export function App() {
  return (
    <BrowserRouter>
      <AppShell>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/trips" element={<TripListPage />} />
          <Route path="/trips/:tripId" element={<TripDetailPage />} />
          <Route path="/days/:dayId" element={<DayDetailPage />} />
        </Routes>
      </AppShell>
    </BrowserRouter>
  )
}
