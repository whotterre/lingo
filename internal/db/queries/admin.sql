-- name: GetAdminForLogin :one
SELECT admin_id, email, password FROM admins
WHERE email = $1 
LIMIT 1;

-- name: CreateAdmin :one
INSERT INTO admins (first_name, last_name, email, password, profile_image_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAdminById :one
SELECT admin_id, first_name, last_name, email, profile_image_url, joined_at
FROM admins
WHERE admin_id = $1
LIMIT 1;

-- name: GetAdminByEmail :one
SELECT admin_id, first_name, last_name, email, profile_image_url, joined_at
FROM admins
WHERE email = $1
LIMIT 1;

-- name: UpdateAdmin :exec
UPDATE admins
SET first_name = $1,
    last_name = $2,
    email = $3,
    password = $4,
    profile_image_url = $5
WHERE admin_id = $6
RETURNING admin_id, first_name, last_name, email, profile_image_url, joined_at;

-- name: DeleteAdmin :exec
DELETE FROM admins
WHERE admin_id = $1
RETURNING admin_id, first_name, last_name, email, profile_image_url, joined_at;

-- name: GetAllAdmins :many
SELECT admin_id, first_name, last_name, email, profile_image_url, joined_at
FROM admins
ORDER BY joined_at DESC
LIMIT $1 OFFSET $2;

