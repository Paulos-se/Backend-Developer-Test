\c spots

SELECT website, count(website)
FROM "MY_TABLE"
GROUP BY website;