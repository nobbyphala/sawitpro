package repository

const (
	queryInserProfile = `
		INSERT INTO
			user_profile
			(full_name, phone_number, password, created_at, updated_at)
		VALUES
			($1, $2 ,$3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id`

	queryGetProfileById = `
		SELECT
			id, 
			full_name, 
			phone_number, 
			password
		FROM
			user_profile
		WHERE
			id = $1`

	queryUpdateProfileById = `
		UPDATE
			user_profile
		SET
			full_name = $1
			phone_name = $2
		WHERE
			id = $3
		`

	queryGetProfileByPhoneNumber = `
		SELECT
			id, 
			full_name, 
			phone_number, 
			password
		FROM
			user_profile
		WHERE
			phone_number = $1`

	queryIncreaseSuccessLoginCount = `
		UPDATE 
			user_profile
		SET 
			success_count = success_count + 1
		WHERE 
			id = $1`
)
