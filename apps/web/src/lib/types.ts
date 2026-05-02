// proto 타입의 최소 미러 (web 에선 ts-proto 결과 대신 가벼운 인터페이스로 사용).
export type TripStatus = 'TRIP_STATUS_UNSPECIFIED' | 'TRIP_STATUS_PLANNING' | 'TRIP_STATUS_ONGOING' | 'TRIP_STATUS_COMPLETED' | 'TRIP_STATUS_ARCHIVED'

export interface Trip {
  tripId: string
  ownerId: string
  title: string
  description?: string
  destination?: string
  countryCode?: string
  startDate?: string
  endDate?: string
  totalBudget?: number
  budgetCurrency?: string
  status?: TripStatus
  createdAt?: string
  updatedAt?: string
}

export interface ListTripsResponse {
  trips: Trip[]
}

export interface Day {
  dayId: string
  tripId: string
  dayDate: string
  dayNumber: number
  region?: string
  weatherForecast?: string
  notes?: string
}

export interface ListDaysResponse {
  days: Day[]
}

export interface Schedule {
  scheduleId: string
  dayId: string
  startTime?: string
  endTime?: string
  title: string
  category?: string
  cost?: number
  costCurrency?: string
  orderIdx?: number
  description?: string
}

export interface ListSchedulesResponse {
  schedules: Schedule[]
}
