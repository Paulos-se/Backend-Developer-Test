\c spots

UPDATE "MY_TABLE" 
SET website=substring(website FROM '(?:.*://)?(?:www\.)?([^/?]*)');

