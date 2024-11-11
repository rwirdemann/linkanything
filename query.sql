-- name: GetLinks :many
select Id, Title, uri, created, tags, draft from links
order by created desc limit $1 offset $2;