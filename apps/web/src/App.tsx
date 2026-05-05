import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { AppShell } from './components/AppShell'
import { HomePage } from './pages/HomePage'
import { TripListPage } from './pages/TripListPage'
import { TripDetailPage } from './pages/TripDetailPage'
import { DayDetailPage } from './pages/DayDetailPage'
import { ExpensesPage } from './pages/ExpensesPage'
import { NotesPage } from './pages/NotesPage'
import { ChecklistPage } from './pages/ChecklistPage'
import { ReservationsPage } from './pages/ReservationsPage'
import { TagsPage } from './pages/TagsPage'
import { CompanionsPage } from './pages/CompanionsPage'
import { MediaPage } from './pages/MediaPage'
import { SharePage } from './pages/SharePage'
import { NearbyPage, TransitPage } from './pages/ExternalPage'

export function App() {
  return (
    <BrowserRouter>
      <a href="#main-content" className="skip-link">
        본문으로 건너뛰기
      </a>
      <AppShell>
        <main id="main-content">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/trips" element={<TripListPage />} />
            <Route path="/trips/:tripId" element={<TripDetailPage />} />
            <Route path="/trips/:tripId/expenses" element={<ExpensesPage />} />
            <Route path="/trips/:tripId/notes" element={<NotesPage />} />
            <Route path="/trips/:tripId/checklist" element={<ChecklistPage />} />
            <Route path="/trips/:tripId/reservations" element={<ReservationsPage />} />
            <Route path="/trips/:tripId/tags" element={<TagsPage />} />
            <Route path="/trips/:tripId/companions" element={<CompanionsPage />} />
            <Route path="/trips/:tripId/media" element={<MediaPage />} />
            <Route path="/trips/:tripId/share" element={<SharePage />} />
            <Route path="/trips/:tripId/nearby" element={<NearbyPage />} />
            <Route path="/trips/:tripId/transit" element={<TransitPage />} />
            <Route path="/days/:dayId" element={<DayDetailPage />} />
          </Routes>
        </main>
      </AppShell>
    </BrowserRouter>
  )
}
