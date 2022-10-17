\c spots

SELECT website,count(website) AS webcount 
FROM "MY_TABLE" GROUP BY website 
HAVING count(website)>1;