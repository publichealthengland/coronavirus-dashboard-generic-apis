package metric_search

const query = `
SELECT metric,
       MAX(metric_name)           AS metric_name,
       MAX(title)                 AS category,
	   MAX(deprecated) AS deprecated,
       CASE
           WHEN MAX(tag) NOTNULL THEN JSONB_AGG(tag ORDER BY tag)
           ELSE '[]'::JSONB
       END AS tags,
	   MAX(docs.last_modified)::DATE::TEXT AS doc_last_modified
FROM covid19.metric_reference AS mr
	LEFT OUTER JOIN covid19.page 		AS pg ON mr.category = pg.id
	LEFT OUTER JOIN covid19.metric_tag  AS mt ON mr.metric = mt.metric_id
	LEFT OUTER JOIN covid19.tag         AS tg ON mt.tag_id = tg.id
    LEFT JOIN (
        SELECT metric_id, MAX(last_modified) AS last_modified
        FROM covid19.metric_asset AS ma
          LEFT JOIN covid19.metric_asset_to_metric AS matm ON ma.id = matm.asset_id
        GROUP BY matm.metric_id
    ) AS docs ON docs.metric_id = mr.metric
WHERE mr.released IS TRUE %s
GROUP BY metric
ORDER BY metric`

const metricTemplate = "%s"

const searchFilter = ` AND (
	   mr.metric ~* ('.*' || REGEXP_REPLACE($%d, '\s+', '', 'g') || '.*')
	OR mr.metric_name ILIKE ('%%' || REGEXP_REPLACE($%d, '\s+', '%%', 'g') || '%%')
  )`

const searchExactFilter = " AND mr.metric ILIKE $%d"

const categoryFilter = " AND LOWER(pg.title) = $%d"

const tagsFilter = ` AND mt.metric_id IN (
    SELECT metric_id
    FROM covid19.tag AS tg
        LEFT JOIN covid19.metric_tag AS mt ON mt.tag_id = tg.id
    WHERE tag = ANY ($%d::VARCHAR[])
  )`
