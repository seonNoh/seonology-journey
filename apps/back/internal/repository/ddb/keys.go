package ddb

import "fmt"

// Key partition/sort key builders for single-table and multi-table patterns.

func UserKey(userID string) string {
	return fmt.Sprintf("USER#%s", userID)
}

func TripKey(tripID string) string {
	return fmt.Sprintf("TRIP#%s", tripID)
}

func DayKey(dayID string) string {
	return fmt.Sprintf("DAY#%s", dayID)
}

func ScheduleKey(scheduleID string) string {
	return fmt.Sprintf("SCHEDULE#%s", scheduleID)
}

func MealKey(mealType string) string {
	return fmt.Sprintf("MEAL#%s", mealType)
}

func MediaKey(takenAt, mediaID string) string {
	return fmt.Sprintf("MEDIA#%s#%s", takenAt, mediaID)
}

func ExpenseKey(expenseID string) string {
	return fmt.Sprintf("EXPENSE#%s", expenseID)
}

func NoteKey(noteID string) string {
	return fmt.Sprintf("NOTE#%s", noteID)
}

func CompanionKey(userID string) string {
	return fmt.Sprintf("COMPANION#%s", userID)
}

func ChecklistKey(checklistID string) string {
	return fmt.Sprintf("CHECKLIST#%s", checklistID)
}

func ReservationKey(reservationID string) string {
	return fmt.Sprintf("RESERVATION#%s", reservationID)
}

func TagKey(tagID string) string {
	return fmt.Sprintf("TAG#%s", tagID)
}

func TemplateKey(templateID string) string {
	return fmt.Sprintf("TEMPLATE#%s", templateID)
}

func FavoriteKey(placeID string) string {
	return fmt.Sprintf("FAVORITE#%s", placeID)
}

func ShareKey(code string) string {
	return fmt.Sprintf("SHARE#%s", code)
}

// AccommodationSK is a fixed sort key (one per day).
const AccommodationSK = "ACCOMMODATION"
