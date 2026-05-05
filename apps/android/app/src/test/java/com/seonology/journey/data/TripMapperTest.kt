package com.seonology.journey.data

import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertNotNull
import org.junit.jupiter.api.Test

class TripMapperTest {

    @Test
    fun `Trip data class defaults are correct`() {
        val trip = Trip(tripId = "t1", title = "Tokyo Trip")
        assertEquals("t1", trip.tripId)
        assertEquals("Tokyo Trip", trip.title)
        assertEquals(null, trip.destination)
        assertEquals(null, trip.status)
    }

    @Test
    fun `CreateTripRequest holds all fields`() {
        val req = CreateTripRequest(
            title = "Hokkaido",
            destination = "Sapporo",
            startDate = "2025-03-01",
            endDate = "2025-03-07",
            budgetCurrency = "JPY",
        )
        assertNotNull(req)
        assertEquals("JPY", req.budgetCurrency)
    }

    @Test
    fun `Day data class dayNumber defaults to 0`() {
        val day = Day(dayId = "d1", tripId = "t1", dayDate = "2025-03-01")
        assertEquals(0, day.dayNumber)
    }
}
