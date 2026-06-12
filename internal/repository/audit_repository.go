package repository

import "database/sql"

type AuditRepository interface {
	LogEventTx(
		tx *sql.Tx,
		hashedNIM string,
		eventType string,
	) error
}

type auditRepository struct{}

func NewAuditRepository() AuditRepository {
	return &auditRepository{}
}

func (r *auditRepository) LogEventTx(
	tx *sql.Tx,
	hashedNIM string,
	eventType string,
) error {

	_, err := tx.Exec(`
		INSERT INTO voter_events
		(
			hashed_nim,
			event_type
		)
		VALUES
		(
			$1,
			$2
		)
	`,
		hashedNIM,
		eventType,
	)

	return err
}
