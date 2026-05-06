package com.seonology.journey.data

import com.squareup.moshi.JsonClass
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

// --- DTOs ---

@JsonClass(generateAdapter = true)
data class Trip(
    val id: String = "",
    val title: String = "",
    val description: String? = null,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
    val status: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class PageInfo(
    val nextCursor: String? = null,
    val hasMore: Boolean = false,
)

@JsonClass(generateAdapter = true)
data class ListTripsResponse(
    val trips: List<Trip> = emptyList(),
    val page: PageInfo? = null,
)

@JsonClass(generateAdapter = true)
data class CreateTripRequest(
    val title: String,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
)

@JsonClass(generateAdapter = true)
data class TripEnvelope(val trip: Trip)

@JsonClass(generateAdapter = true)
data class Day(
    val id: String = "",
    val tripId: String = "",
    val date: String = "",
    val dayNumber: Int = 0,
    val region: String? = null,
    val weather: String? = null,
    val dailySummary: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListDaysResponse(val days: List<Day> = emptyList())

@JsonClass(generateAdapter = true)
data class GeoPoint(val latitude: Double = 0.0, val longitude: Double = 0.0)

@JsonClass(generateAdapter = true)
data class GeocodePlace(
    val placeId: String = "",
    val name: String = "",
    val address: String = "",
    val location: GeoPoint? = null,
)

@JsonClass(generateAdapter = true)
data class GeocodeResponse(val places: List<GeocodePlace> = emptyList())

@JsonClass(generateAdapter = true)
data class Schedule(
    val id: String = "",
    val dayId: String = "",
    val order: Int = 0,
    val title: String = "",
    val startTime: String? = null,
    val endTime: String? = null,
    val placeName: String? = null,
    val location: GeoPoint? = null,
    val notes: String? = null,
    val region: String? = null,
    val cost: Money? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class CreateScheduleRequest(
    val title: String,
    val startTime: String? = null,
    val endTime: String? = null,
    val placeName: String? = null,
    val location: GeoPoint? = null,
    val notes: String? = null,
    val region: String? = null,
    val cost: Money? = null,
)

@JsonClass(generateAdapter = true)
data class ScheduleEnvelope(val schedule: Schedule = Schedule())

@JsonClass(generateAdapter = true)
data class ListSchedulesResponse(val schedules: List<Schedule> = emptyList())

@JsonClass(generateAdapter = true)
data class Money(val currency: String, val amount: Long)

@JsonClass(generateAdapter = true)
data class Expense(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val category: String = "",
    val amount: Money? = null,
    val description: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListExpensesResponse(
    val expenses: List<Expense> = emptyList(),
    val page: PageInfo? = null,
)

@JsonClass(generateAdapter = true)
data class Note(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val content: String = "",
    val mood: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListNotesResponse(val notes: List<Note> = emptyList())

@JsonClass(generateAdapter = true)
data class Meal(
    val dayId: String = "",
    val mealType: String = "",
    val source: String? = null,
    val restaurantName: String? = null,
    val menu: String? = null,
    val cost: Money? = null,
    val rating: Int? = null,
    val review: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListMealsResponse(val meals: List<Meal> = emptyList())

@JsonClass(generateAdapter = true)
data class Accommodation(
    val dayId: String = "",
    val name: String = "",
    val checkInTime: String? = null,
    val checkOutTime: String? = null,
    val cost: Money? = null,
    val amenities: String? = null,
    val address: String? = null,
)

@JsonClass(generateAdapter = true)
data class AccommodationEnvelope(val accommodation: Accommodation? = null)

@JsonClass(generateAdapter = true)
data class Media(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val s3Key: String = "",
    val thumbnailS3Key: String? = null,
    val mimeType: String? = null,
    val takenAt: String? = null,
    val caption: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListMediaResponse(
    val items: List<Media> = emptyList(),
    val page: PageInfo? = null,
)

// --- Checklist ---

@JsonClass(generateAdapter = true)
data class ChecklistItem(
    val id: String = "",
    val tripId: String = "",
    val category: String = "",
    val item: String = "",
    val isChecked: Boolean? = null,
)

@JsonClass(generateAdapter = true)
data class ListChecklistResponse(val items: List<ChecklistItem> = emptyList())

@JsonClass(generateAdapter = true)
data class ChecklistItemEnvelope(val item: ChecklistItem = ChecklistItem())

@JsonClass(generateAdapter = true)
data class CreateChecklistRequest(val item: String, val category: String)

@JsonClass(generateAdapter = true)
data class UpdateChecklistRequest(val isChecked: Boolean? = null, val item: String? = null)

// --- Reservations ---

@JsonClass(generateAdapter = true)
data class Reservation(
    val id: String = "",
    val tripId: String = "",
    val type: String = "",
    val vendor: String? = null,
    val confirmNumber: String? = null,
    val reservedAt: String? = null,
    val cost: Money? = null,
    val attachmentS3Key: String? = null,
    val notes: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListReservationsResponse(val reservations: List<Reservation> = emptyList())

@JsonClass(generateAdapter = true)
data class CreateReservationRequest(
    val type: String,
    val vendor: String? = null,
    val confirmNumber: String? = null,
    val reservedAt: String? = null,
    val cost: Money? = null,
    val notes: String? = null,
)

// --- Tags ---

@JsonClass(generateAdapter = true)
data class Tag(
    val id: String = "",
    val userId: String? = null,
    val name: String = "",
    val color: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListTagsResponse(val tags: List<Tag> = emptyList())

@JsonClass(generateAdapter = true)
data class CreateTagRequest(val name: String, val color: String? = null)

// --- Companions ---

@JsonClass(generateAdapter = true)
data class Companion(
    val tripId: String = "",
    val memberId: String = "",
    val displayName: String? = null,
    val avatarUrl: String? = null,
    val role: String = "",
    val invitedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListCompanionsResponse(val companions: List<Companion> = emptyList())

@JsonClass(generateAdapter = true)
data class AddCompanionRequest(val memberId: String, val role: String)

@JsonClass(generateAdapter = true)
data class UpdateCompanionRequest(val role: String)

// --- Share ---

@JsonClass(generateAdapter = true)
data class Share(
    val code: String = "",
    val tripId: String = "",
    val permission: String = "",
    val expiresAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListSharesResponse(val shares: List<Share> = emptyList())

@JsonClass(generateAdapter = true)
data class CreateShareRequest(val permission: String, val expiresInHours: Int? = null)

// --- Expenses CRUD ---

@JsonClass(generateAdapter = true)
data class ExpenseEnvelope(val expense: Expense = Expense())

@JsonClass(generateAdapter = true)
data class CreateExpenseRequest(
    val dayId: String? = null,
    val category: String,
    val amount: Money,
    val paymentMethod: String? = null,
    val description: String? = null,
    val spentAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ExpenseSummaryEntry(val category: String? = null, val date: String? = null, val total: Money? = null)

@JsonClass(generateAdapter = true)
data class ExpenseSummary(
    val tripId: String = "",
    val grandTotal: Money? = null,
    val byCategory: List<ExpenseSummaryEntry> = emptyList(),
    val byDay: List<ExpenseSummaryEntry> = emptyList(),
)

// --- Notes CRUD ---

@JsonClass(generateAdapter = true)
data class NoteEnvelope(val note: Note = Note())

@JsonClass(generateAdapter = true)
data class CreateNoteRequest(
    val dayId: String? = null,
    val content: String,
    val mood: String? = null,
)

// --- Meals upsert ---

@JsonClass(generateAdapter = true)
data class UpsertMealRequest(
    val mealType: String,
    val source: String? = null,
    val restaurantName: String? = null,
    val menu: String? = null,
    val cost: Money? = null,
    val rating: Int? = null,
    val review: String? = null,
)

@JsonClass(generateAdapter = true)
data class MealEnvelope(val meal: Meal = Meal())

// --- Accommodation upsert ---

@JsonClass(generateAdapter = true)
data class UpsertAccommodationRequest(
    val name: String,
    val checkInTime: String? = null,
    val checkOutTime: String? = null,
    val cost: Money? = null,
    val amenities: String? = null,
    val address: String? = null,
)

// --- Schedule reorder ---

@JsonClass(generateAdapter = true)
data class ReorderSchedulesRequest(val orderedIds: List<String>)

// --- Media URL/Upload ---

@JsonClass(generateAdapter = true)
data class GetMediaUrlResponse(val url: String = "", val expiresAt: String? = null)

@JsonClass(generateAdapter = true)
data class GetUploadUrlRequest(val filename: String, val mimeType: String, val size: Long)

@JsonClass(generateAdapter = true)
data class GetUploadUrlResponse(
    val uploadUrl: String = "",
    val s3Key: String = "",
    val expiresAt: String? = null,
    val mediaId: String = "",
)

@JsonClass(generateAdapter = true)
data class ConfirmMediaRequest(
    val mediaId: String,
    val s3Key: String,
    val dayId: String? = null,
    val scheduleId: String? = null,
    val caption: String? = null,
    val takenAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class MediaEnvelope(val media: Media = Media())

// --- Nearby / Transit ---

@JsonClass(generateAdapter = true)
data class NearbyPlace(
    val name: String = "",
    val address: String = "",
    val lat: Double = 0.0,
    val lng: Double = 0.0,
    val rating: Double? = null,
    val types: List<String>? = null,
)

@JsonClass(generateAdapter = true)
data class NearbyResponse(val results: List<NearbyPlace> = emptyList())

@JsonClass(generateAdapter = true)
data class TransitStep(val instruction: String = "", val distance: String = "")

@JsonClass(generateAdapter = true)
data class TransitRoute(
    val summary: String? = null,
    val duration: String? = null,
    val distance: String? = null,
    val steps: List<TransitStep> = emptyList(),
)

@JsonClass(generateAdapter = true)
data class TransitResponse(val routes: List<TransitRoute> = emptyList())

// --- API ---
//
// Endpoints exposed here mirror the REST gateway in apps/api. Only the
// methods actually used by the Android client are listed; adding more is
// a matter of appending signatures, no plumbing change needed.

interface JourneyApi {
    // Trips
    @GET("trips")
    suspend fun listTrips(
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListTripsResponse

    @POST("trips")
    suspend fun createTrip(@Body req: CreateTripRequest): TripEnvelope

    @GET("trips/{tripId}")
    suspend fun getTrip(@Path("tripId") tripId: String): TripEnvelope

    @PATCH("trips/{tripId}")
    suspend fun updateTrip(@Path("tripId") tripId: String, @Body req: CreateTripRequest): TripEnvelope

    @DELETE("trips/{tripId}")
    suspend fun deleteTrip(@Path("tripId") tripId: String)

    // Days
    @GET("trips/{tripId}/days")
    suspend fun listDays(@Path("tripId") tripId: String): ListDaysResponse

    // Schedules
    @GET("days/{dayId}/schedules")
    suspend fun listSchedules(@Path("dayId") dayId: String): ListSchedulesResponse

    @POST("days/{dayId}/schedules")
    suspend fun createSchedule(
        @Path("dayId") dayId: String,
        @Body body: CreateScheduleRequest,
    ): ScheduleEnvelope

    @DELETE("schedules/{scheduleId}")
    suspend fun deleteSchedule(@Path("scheduleId") scheduleId: String)

    // Meals
    @GET("days/{dayId}/meals")
    suspend fun listMeals(@Path("dayId") dayId: String): ListMealsResponse

    // Accommodation (단일 리소스)
    @GET("days/{dayId}/accommodation")
    suspend fun getAccommodation(@Path("dayId") dayId: String): AccommodationEnvelope

    // Expenses
    @GET("trips/{tripId}/expenses")
    suspend fun listExpenses(
        @Path("tripId") tripId: String,
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListExpensesResponse

    // Notes
    @GET("trips/{tripId}/notes")
    suspend fun listNotes(@Path("tripId") tripId: String): ListNotesResponse

    // Media
    @GET("trips/{tripId}/media")
    suspend fun listMedia(
        @Path("tripId") tripId: String,
        @Query("day_id") dayId: String? = null,
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListMediaResponse

    // External - geocoding via OpenStreetMap Nominatim (백엔드 프록시)
    @GET("external/geocode")
    suspend fun geocode(@Query("q") query: String): GeocodeResponse

    // External - nearby search
    @GET("external/nearby")
    suspend fun nearby(
        @Query("lat") lat: Double,
        @Query("lng") lng: Double,
        @Query("radius") radius: Int = 1000,
        @Query("type") type: String = "restaurant",
    ): NearbyResponse

    // External - transit search
    @GET("external/transit")
    suspend fun transit(
        @Query("origin_lat") originLat: Double,
        @Query("origin_lng") originLng: Double,
        @Query("dest_lat") destLat: Double,
        @Query("dest_lng") destLng: Double,
        @Query("departure_time") departureTime: String? = null,
    ): TransitResponse

    // Schedules - reorder
    @POST("days/{dayId}/schedules:reorder")
    suspend fun reorderSchedules(
        @Path("dayId") dayId: String,
        @Body body: ReorderSchedulesRequest,
    )

    // Meals - upsert / delete
    @PUT("days/{dayId}/meals")
    suspend fun upsertMeal(
        @Path("dayId") dayId: String,
        @Body body: UpsertMealRequest,
    ): MealEnvelope

    @DELETE("days/{dayId}/meals/{mealType}")
    suspend fun deleteMeal(
        @Path("dayId") dayId: String,
        @Path("mealType") mealType: String,
    )

    // Accommodation - upsert / delete
    @PUT("days/{dayId}/accommodation")
    suspend fun upsertAccommodation(
        @Path("dayId") dayId: String,
        @Body body: UpsertAccommodationRequest,
    ): AccommodationEnvelope

    @DELETE("days/{dayId}/accommodation")
    suspend fun deleteAccommodation(@Path("dayId") dayId: String)

    // Expenses - CRUD + summary
    @POST("trips/{tripId}/expenses")
    suspend fun createExpense(
        @Path("tripId") tripId: String,
        @Body body: CreateExpenseRequest,
    ): ExpenseEnvelope

    @DELETE("expenses/{expenseId}")
    suspend fun deleteExpense(@Path("expenseId") expenseId: String)

    @GET("trips/{tripId}/expense-summary")
    suspend fun expenseSummary(@Path("tripId") tripId: String): ExpenseSummary

    // Notes - CRUD
    @POST("trips/{tripId}/notes")
    suspend fun createNote(
        @Path("tripId") tripId: String,
        @Body body: CreateNoteRequest,
    ): NoteEnvelope

    @DELETE("notes/{noteId}")
    suspend fun deleteNote(@Path("noteId") noteId: String)

    // Checklist - CRUD
    @GET("trips/{tripId}/checklist")
    suspend fun listChecklist(@Path("tripId") tripId: String): ListChecklistResponse

    @POST("trips/{tripId}/checklist")
    suspend fun createChecklistItem(
        @Path("tripId") tripId: String,
        @Body body: CreateChecklistRequest,
    ): ChecklistItemEnvelope

    @PATCH("checklist/{itemId}")
    suspend fun updateChecklistItem(
        @Path("itemId") itemId: String,
        @Body body: UpdateChecklistRequest,
    ): ChecklistItemEnvelope

    @DELETE("checklist/{itemId}")
    suspend fun deleteChecklistItem(@Path("itemId") itemId: String)

    // Reservations - CRUD
    @GET("trips/{tripId}/reservations")
    suspend fun listReservations(@Path("tripId") tripId: String): ListReservationsResponse

    @POST("trips/{tripId}/reservations")
    suspend fun createReservation(
        @Path("tripId") tripId: String,
        @Body body: CreateReservationRequest,
    )

    @DELETE("reservations/{reservationId}")
    suspend fun deleteReservation(@Path("reservationId") reservationId: String)

    // Tags - global + per-trip
    @GET("tags")
    suspend fun listTags(): ListTagsResponse

    @POST("tags")
    suspend fun createTag(@Body body: CreateTagRequest)

    @DELETE("tags/{tagId}")
    suspend fun deleteTag(@Path("tagId") tagId: String)

    @GET("trips/{tripId}/tags")
    suspend fun listTripTags(@Path("tripId") tripId: String): ListTagsResponse

    @PUT("trips/{tripId}/tags/{tagId}")
    suspend fun attachTripTag(
        @Path("tripId") tripId: String,
        @Path("tagId") tagId: String,
    )

    @DELETE("trips/{tripId}/tags/{tagId}")
    suspend fun detachTripTag(
        @Path("tripId") tripId: String,
        @Path("tagId") tagId: String,
    )

    // Companions - CRUD
    @GET("trips/{tripId}/companions")
    suspend fun listCompanions(@Path("tripId") tripId: String): ListCompanionsResponse

    @POST("trips/{tripId}/companions")
    suspend fun addCompanion(
        @Path("tripId") tripId: String,
        @Body body: AddCompanionRequest,
    )

    @PATCH("trips/{tripId}/companions/{memberId}")
    suspend fun updateCompanion(
        @Path("tripId") tripId: String,
        @Path("memberId") memberId: String,
        @Body body: UpdateCompanionRequest,
    )

    @DELETE("trips/{tripId}/companions/{memberId}")
    suspend fun removeCompanion(
        @Path("tripId") tripId: String,
        @Path("memberId") memberId: String,
    )

    // Shares - CRUD
    @GET("trips/{tripId}/shares")
    suspend fun listShares(@Path("tripId") tripId: String): ListSharesResponse

    @POST("trips/{tripId}/shares")
    suspend fun createShare(
        @Path("tripId") tripId: String,
        @Body body: CreateShareRequest,
    ): Share

    @DELETE("shares/{code}")
    suspend fun deleteShare(@Path("code") code: String)

    // Media - URL / upload / delete
    @GET("media/{mediaId}/url")
    suspend fun getMediaUrl(
        @Path("mediaId") mediaId: String,
        @Query("thumbnail") thumbnail: Boolean? = null,
    ): GetMediaUrlResponse

    @POST("trips/{tripId}/media:upload-url")
    suspend fun getMediaUploadUrl(
        @Path("tripId") tripId: String,
        @Body body: GetUploadUrlRequest,
    ): GetUploadUrlResponse

    @POST("trips/{tripId}/media:confirm")
    suspend fun confirmMedia(
        @Path("tripId") tripId: String,
        @Body body: ConfirmMediaRequest,
    ): MediaEnvelope

    @DELETE("media/{mediaId}")
    suspend fun deleteMedia(@Path("mediaId") mediaId: String)
}
