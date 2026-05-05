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
  location?: GeoPoint
}

export type MealType =
  | 'MEAL_TYPE_UNSPECIFIED'
  | 'MEAL_TYPE_BREAKFAST'
  | 'MEAL_TYPE_LUNCH'
  | 'MEAL_TYPE_DINNER'

export type MealSource =
  | 'MEAL_SOURCE_UNSPECIFIED'
  | 'MEAL_SOURCE_HOTEL'
  | 'MEAL_SOURCE_LOCAL'
  | 'MEAL_SOURCE_CONVENIENCE'
  | 'MEAL_SOURCE_SKIP'

export interface Meal {
  dayId: string
  mealType: MealType
  source?: MealSource
  restaurantName?: string
  menu?: string
  cost?: Money
  rating?: number
  review?: string
}

export interface ListMealsResponse {
  meals: Meal[]
}

export interface Accommodation {
  dayId: string
  name: string
  checkInTime?: string
  checkOutTime?: string
  cost?: Money
  amenities?: string
  address?: string
}

export type ExpenseCategory =
  | 'EXPENSE_CATEGORY_UNSPECIFIED'
  | 'EXPENSE_CATEGORY_TRANSPORT'
  | 'EXPENSE_CATEGORY_FOOD'
  | 'EXPENSE_CATEGORY_LODGING'
  | 'EXPENSE_CATEGORY_ACTIVITY'
  | 'EXPENSE_CATEGORY_SHOPPING'
  | 'EXPENSE_CATEGORY_OTHER'

export type PaymentMethod =
  | 'PAYMENT_METHOD_UNSPECIFIED'
  | 'PAYMENT_METHOD_CASH'
  | 'PAYMENT_METHOD_CARD'
  | 'PAYMENT_METHOD_TRANSFER'

export interface Expense {
  id: string
  tripId: string
  dayId?: string
  category: ExpenseCategory
  amount: Money
  paymentMethod?: PaymentMethod
  description?: string
  spentAt?: string
}

export interface ListExpensesResponse {
  expenses: Expense[]
}

export interface ExpenseSummary {
  tripId: string
  grandTotal?: Money
  byCategory: { category: ExpenseCategory; total: Money }[]
  byDay: { date: string; total: Money }[]
}

export interface CreateExpenseInput {
  dayId?: string
  category: ExpenseCategory
  amount: Money
  paymentMethod?: PaymentMethod
  description?: string
  spentAt?: string
}

export type ChecklistCategory =
  | 'CHECKLIST_CATEGORY_UNSPECIFIED'
  | 'CHECKLIST_CATEGORY_PACKING'
  | 'CHECKLIST_CATEGORY_TODO'
  | 'CHECKLIST_CATEGORY_BOOKING'

export interface ChecklistItem {
  id: string
  tripId: string
  category: ChecklistCategory
  item: string
  isChecked?: boolean
}

export interface ListChecklistResponse {
  items: ChecklistItem[]
}

export interface Note {
  id: string
  tripId: string
  dayId?: string
  content: string
  mood?: string
}

export interface ListNotesResponse {
  notes: Note[]
}

export type ReservationType =
  | 'RESERVATION_TYPE_UNSPECIFIED'
  | 'RESERVATION_TYPE_FLIGHT'
  | 'RESERVATION_TYPE_HOTEL'
  | 'RESERVATION_TYPE_ACTIVITY'
  | 'RESERVATION_TYPE_RESTAURANT'
  | 'RESERVATION_TYPE_TRANSPORT'

export interface Reservation {
  id: string
  tripId: string
  type: ReservationType
  vendor?: string
  confirmNumber?: string
  reservedAt?: string
  cost?: Money
  attachmentS3Key?: string
  notes?: string
}

export interface ListReservationsResponse {
  reservations: Reservation[]
}

export interface CreateReservationInput {
  type: ReservationType
  vendor?: string
  confirmNumber?: string
  reservedAt?: string
  cost?: Money
  notes?: string
}

export interface Tag {
  id: string
  userId?: string
  name: string
  color?: string
}

export interface ListTagsResponse {
  tags: Tag[]
}

export type CompanionRole =
  | 'COMPANION_ROLE_UNSPECIFIED'
  | 'COMPANION_ROLE_OWNER'
  | 'COMPANION_ROLE_EDITOR'
  | 'COMPANION_ROLE_VIEWER'

export interface Companion {
  tripId: string
  memberId: string
  displayName?: string
  avatarUrl?: string
  role: CompanionRole
  invitedAt?: string
}

export interface ListCompanionsResponse {
  companions: Companion[]
}

export interface Media {
  id: string
  tripId: string
  dayId?: string
  scheduleId?: string
  s3Key: string
  thumbnailS3Key?: string
  mimeType?: string
  size?: number | string
  takenAt?: string
  caption?: string
}

export interface ListMediaResponse {
  items: Media[]
}

export interface GetUploadUrlResponse {
  uploadUrl: string
  s3Key: string
  expiresAt?: string
  mediaId: string
}

export interface GetMediaUrlResponse {
  url: string
  expiresAt?: string
}

export interface Share {
  code: string
  tripId: string
  permission: CompanionRole
  expiresAt?: string
}
