package com.seonology.journey.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.Entity
import androidx.room.PrimaryKey
import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow

// --- Entities ---

@Entity(tableName = "trips")
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
)

@Entity(tableName = "days")
data class DayEntity(
    @PrimaryKey val dayId: String,
    val tripId: String,
    val dayDate: String,
    val dayNumber: Int = 0,
    val region: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
)

@Entity(tableName = "schedules")
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
)

@Entity(tableName = "expenses")
data class ExpenseEntity(
    @PrimaryKey val expenseId: String,
    val tripId: String,
    val dayId: String? = null,
    val category: String,
    val amount: Double,
    val currency: String,
    val description: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
)

@Entity(tableName = "notes")
data class NoteEntity(
    @PrimaryKey val noteId: String,
    val tripId: String,
    val title: String,
    val content: String? = null,
    val updatedAt: Long = System.currentTimeMillis(),
)

@Entity(tableName = "media")
data class MediaEntity(
    @PrimaryKey val mediaId: String,
    val tripId: String,
    val dayId: String? = null,
    val s3Key: String,
    val mimeType: String? = null,
    val thumbnailUrl: String? = null,
    val uploadedAt: Long = System.currentTimeMillis(),
)

// --- DAOs ---

@Dao
interface TripDao {
    @Query("SELECT * FROM trips ORDER BY updatedAt DESC")
    fun observeAll(): Flow<List<TripEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(trip: TripEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(trips: List<TripEntity>)

    @Query("DELETE FROM trips WHERE tripId = :id")
    suspend fun delete(id: String)
}

@Dao
interface DayDao {
    @Query("SELECT * FROM days WHERE tripId = :tripId ORDER BY dayNumber")
    fun observeByTrip(tripId: String): Flow<List<DayEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(days: List<DayEntity>)
}

@Dao
interface ScheduleDao {
    @Query("SELECT * FROM schedules WHERE dayId = :dayId ORDER BY startTime")
    fun observeByDay(dayId: String): Flow<List<ScheduleEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(schedule: ScheduleEntity)
}

@Dao
interface ExpenseDao {
    @Query("SELECT * FROM expenses WHERE tripId = :tripId ORDER BY updatedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<ExpenseEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(expense: ExpenseEntity)
}

@Dao
interface NoteDao {
    @Query("SELECT * FROM notes WHERE tripId = :tripId ORDER BY updatedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<NoteEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(note: NoteEntity)
}

@Dao
interface MediaDao {
    @Query("SELECT * FROM media WHERE tripId = :tripId ORDER BY uploadedAt DESC")
    fun observeByTrip(tripId: String): Flow<List<MediaEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(media: MediaEntity)
}

// --- Database ---

@Database(
    entities = [TripEntity::class, DayEntity::class, ScheduleEntity::class, ExpenseEntity::class, NoteEntity::class, MediaEntity::class],
    version = 1,
    exportSchema = false,
)
abstract class JourneyDatabase : RoomDatabase() {
    abstract fun tripDao(): TripDao
    abstract fun dayDao(): DayDao
    abstract fun scheduleDao(): ScheduleDao
    abstract fun expenseDao(): ExpenseDao
    abstract fun noteDao(): NoteDao
    abstract fun mediaDao(): MediaDao
}
