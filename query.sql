-- name: GetLinks :many
select distinct l.Id, l.title, l.uri, l.created, l.draft
from links l
         left join public.tags_links tl on l.id = tl.link_id
         left join public.tags tag on tag.id = tl.tag_id
where tag.name = ANY($1::varchar[]) or $1 is NULL
order by created desc
limit $2 offset $3;