// proto/journey/v1 의 protojson 표현을 가벼운 ts 인터페이스로 미러.

export interface Money {
  currency: string
  amount: number | string // int64 는 protojson 에서 string 으로 직렬화될 수 있다.
}

export interface AuditTimestamps {
  createdAt?: string
  updatedAt?: string
}

export type TripStatus =
  | 'TRIP_STATUS_UNSPECIFIED'
  | 'TRIP_STATUS_PLANNING'
  | 'TRIP_STATUS_ONGOING'
  | 'TRIP_STATUS_COMPLETED'
  | 'TRIP_STATUS_ARCHIVED'

export interface Trip {
  id: string
  ownerId: string
  title: string
  description?: string
  startDate?: string
  endDate?: string
  status?: TripStatus
  coverImageUrl?: string
  totalBudget?: Money
  destination?: string
  countryCode?: string
  audit?: AuditTimestamps
}

export interface ListTripsResponse {
  trips: Trip[]
}

export interface Day {
  id: string
  tripId: string
  dayNumber: number
  date: string
  dayOfWeek?: string
  region?: string
  weather?: string
  dailySummary?: string
  audit?: AuditTimestamps
}

export interface ListDaysResponse {
  days: Day[]
}

export type ScheduleCategory = string
export type TransportType = string

export interface GeoPoint {
  latitude: number
  longitude: number
}

export interface Schedule {
  id: string
  dayId: string
  order: number
  startTime?: string
  endTime?: string
  title: string
  region?: string
  category?: ScheduleCategory
  transport?: TransportType
  transportDetail?: string
  cost?: Money
  placeName?: string
  location?: GeoPoint
  notes?: string
  isCompleted?: boolean
  audit?: AuditTimestamps
}

export interface ListSchedulesResponse {
  schedules: Schedule[]
}

export interface CreateTripInput {
  title: string
  description?: string
  startDate?: string
  endDate?: string
  destination?: string
  countryCode?: string
  totalBudget?: Money
}

export interface CreateScheduleInput {
  startTime?: string
  endTime?: string
  title: string
  region?: string
  notes?: string
  cost?: Money
  placeName?: string
}
