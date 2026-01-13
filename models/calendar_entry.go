package models

import (
	"time"
)

type CalendarEntry struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      string    `json:"date"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"created_at"`
}

func GetCalendarEntries(userID int, startDate, endDate string) ([]CalendarEntry, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, date, checked, created_at
		FROM calendar_entries
		WHERE user_id = ? AND date BETWEEN ? AND ?
		ORDER BY date`,
		userID, startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CalendarEntry
	for rows.Next() {
		var entry CalendarEntry
		err := rows.Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Checked, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func CheckDay(userID int, date string) error {
	_, err := DB.Exec(`
		INSERT OR REPLACE INTO calendar_entries (user_id, date, checked)
		VALUES (?, ?, TRUE)`,
		userID, date)
	return err
}

func UncheckDay(userID int, date string) error {
	_, err := DB.Exec(`
		UPDATE calendar_entries SET checked = FALSE
		WHERE user_id = ? AND date = ?`,
		userID, date)
	return err
}

func GetCheckedDaysCount(userID int, startDate, endDate string) (int, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM calendar_entries
		WHERE user_id = ? AND checked = TRUE`,
		userID).Scan(&count)
	return count, err
}
