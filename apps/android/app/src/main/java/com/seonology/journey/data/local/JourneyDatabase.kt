package com.seonology.journey.data.local

import androidx.room.Dao
import androidx.room.Database
import androidx.room.Entity
import androidx.room.Index
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.PrimaryKey
import androidx.room.Query
import androidx.room.RoomDatabase
import kotlinx.coroutines.flow.Flow

// --- Entities ---
//
// All tables carry updatedAt (ms since epoch) for last-write-wins sync and
// isDirty to mark pending local changes. Foreign-key columns (tripId,
// dayId) are explicitly indexed so per-parent queries become O(log N)
// rather than full-table scans.

@Entity(
    tableName = "trips",
    indices = [Index("updatedAt"), Index("isDirty")],
)
data class TripEntity(
    @PrimaryKey val tripId: String,
    val title: String,
    val description: String? = null,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
    val status: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

@Entity(
    tableName = "days",
    indices = [Index("tripId"), Index("updatedAt"), Index("isDirty")],
)
data class DayEntity(
    @PrimaryKey val dayId: String,
    val tripId: String,
    val dayDate: String,
    val dayNumber: Int = 0,
    val region: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

@Entity(
    tableName = "schedules",
    indices = [Index("dayId"), Index("updatedAt"), Index("isDirty")],
)
data class ScheduleEntity(
    @PrimaryKey val scheduleId: String,
    val dayId: String,
    val title: String,
    val category: String? = null,
    val startTime: String? = null,
    val endTime: String? = null,
    val locationLat: Double? = null,
    val locationLng: Double? = null,
    val notes: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

@Entity(
    tableName = "expenses",
    indices = [Index("tripId"), Index("dayId"), Index("updatedAt"), Index("isDirty")],
)
data class ExpenseEntity(
    @PrimaryKey val expenseId: String,
    val tripId: String,
    val dayId: String? = null,
    val category: String,
    val amount: Double,
    val currency: String,
    val description: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

@Entity(
    tableName = "notes",
    indices = [Index("tripId"), Index("updatedAt"), Index("isDirty")],
)
data class NoteEntity(
    @PrimaryKey val noteId: String,
    val tripId: String,
    val title: String,
    val content: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

@Entity(
    tableName = "media",
    indices = [Index("tripId"), Index("dayId"), Index("uploadedAt")],
)
data class MediaEntity(
    @PrimaryKey val mediaId: String,
    val tripId: String,
    val dayId: String? = null,
    val s3Key: String,
    val mimeType: String? = null,
    val thumbnailUrl: String? = null,
    val uploadedAt: Long = System.currentTimeMillis(),
    val isDirty: Boolean = false,
)

/** Bookkeeping for incremental pull. One row per logical resource. */
@Entity(tableName = "sync_state")
data class SyncStateEntity(
    @PrimaryKey val resource: String,
    val lastSyncedAt: Long,
)

// --- DAOs ---

@Dao
interface TripDao {
    @Query("SELECT * FROM trips ORDER BY updatedAt DESC")
    fun observeAll(): Flow<List<TripEntity>>

    @Query("SELECT * FROM trips WHERE tripId = :id LIMIT 1")
    suspend fun get(id: String): TripEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(trip: TripEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(trips: List<TripEntity>)

    @Query("DELETE FROM trips WHERE tripId = :id")
    suspend fun delete(id: String)

    /** Pending local changes awaiting push to server. */
    @Query("SELECT * FROM trips WHERE isDirty = 1")
    suspend fun pending(): List<TripEntity>

    @Query("UPDATE trips SET isDirty = 0 WHERE tripId = :id")
    suspend fun clearDirty(id: String)
}

@Dao
interface DayDao {
    @Query("SELECT * FROM days WHERE tripId = :tripId ORDER BY dayNumber")
    fun observeByTrip(tripId: String): Flow<List<DayEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(days: List<DayEntity>)

    @Query("DELETE FROM days WHERE tripId = :tripId")
    suspend fun deleteByTrip(tripId: String)
}

@Dao
interface ScheduleDao {
    @Query("SELECT * FROM schedules WHERE dayId = :dayId ORDER BY startTime")
    fun observeByDay(dayId: String): Flow<List<ScheduleEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(schedule: ScheduleEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(schedules: List<ScheduleEntity>)

    @Query("SELECT * FROM schedules WHERE isDirty = 1")
    suspend fun pending(): List<ScheduleEntity>

    @Query("UPDATE schedules SET isDirty = 0 WHERE scheduleId = :id")
    suspend fun clearDirty(id: String)
}

@Dao
interface ExpenseDao {
    @Query("SELECT * FROM expenses WHERE tripId = :tripId ORDER BY updatedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<ExpenseEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(expense: ExpenseEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(expenses: List<ExpenseEntity>)

    @Query("SELECT * FROM expenses WHERE isDirty = 1")
    suspend fun pending(): List<ExpenseEntity>

    @Query("UPDATE expenses SET isDirty = 0 WHERE expenseId = :id")
    suspend fun clearDirty(id: String)
}

@Dao
interface NoteDao {
    @Query("SELECT * FROM notes WHERE tripId = :tripId ORDER BY updatedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<NoteEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(note: NoteEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(notes: List<NoteEntity>)

    @Query("SELECT * FROM notes WHERE isDirty = 1")
    suspend fun pending(): List<NoteEntity>

    @Query("UPDATE notes SET isDirty = 0 WHERE noteId = :id")
    suspend fun clearDirty(id: String)
}

@Dao
interface MediaDao {
    @Query("SELECT * FROM media WHERE tripId = :tripId ORDER BY uploadedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<MediaEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(media: MediaEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(items: List<MediaEntity>)
}

@Dao
interface SyncStateDao {
    @Query("SELECT lastSyncedAt FROM sync_state WHERE resource = :resource")
    suspend fun lastSyncedAt(resource: String): Long?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun put(state: SyncStateEntity)
}

// --- Database ---

@Database(
    entities = [
        TripEntity::class,
        DayEntity::class,
        ScheduleEntity::class,
        ExpenseEntity::class,
        NoteEntity::class,
        MediaEntity::class,
        SyncStateEntity::class,
    ],
    version = 2,
    exportSchema = false,
)
abstract class JourneyDatabase : RoomDatabase() {
    abstract fun tripDao(): TripDao
    abstract fun dayDao(): DayDao
    abstract fun scheduleDao(): ScheduleDao
    abstract fun expenseDao(): ExpenseDao
    abstract fun noteDao(): NoteDao
    abstract fun mediaDao(): MediaDao
    abstract fun syncStateDao(): SyncStateDao
}
