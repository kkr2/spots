package repository

const (
	findTotalSpotsInRange = ` select count(s.id)
								from  spots s
								where st_intersects(ST_Transform( ST_MakeEnvelope(
								(ST_X(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))-($3)), 
								(ST_Y(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))-($3)),
								(ST_X(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))+($3)), 
								(ST_Y(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))+($3)), 
								2163
								), 4326)
								, s.coordinates)`

	findSpotsInRange = `
				select  s.id ,s."name" ,COALESCE(s.website , '') as website ,COALESCE(s.description , '') as description ,s.rating,ST_X(s.coordinates) AS "coordinates.longitude", ST_Y(s.coordinates) AS "coordinates.latitude" 
				from  spots s
				where st_intersects(ST_Transform( ST_MakeEnvelope(
				(ST_X(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))-($3)), 
				(ST_Y(ST_Transform(ST_SetSRID(ST_MakePoint($1 ,$2),4326),2163))-($3)),
				(ST_X(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))+($3)), 
				(ST_Y(ST_Transform(ST_SetSRID(ST_MakePoint($1, $2),4326),2163))+($3)), 
				2163
				), 4326)
				, s.coordinates)
				order by s.coordinates  <-> ST_SetSRID(ST_MakePoint($1, $2),4326)
				offset $4 
				limit $5
	`
)
