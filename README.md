## Background

I built An API and wrote some sql queries for the purpose of locating spots within a certain area. The intention here is to mimic the building of a real world backend service with GO language, psql and PostGIS (a spatial database extender for PostgreSQL).

## Step 1 - Creating tables and Seeding

To create a database, a table and seed the table please query `spots.sql` in the `db` folder. The before_update.txt in the root directory shows the output of the query.

## Step 2 - Quering the database

To update the database website coloumn with only the domain name please execute `1-update-the-table.sql` located in the `db` folder. The after_update.txt shows the updated table.

To count how many spots contain the same domain please query `2-count-same-domain.sql`. The output is in count-same-domain.txt.

To query spots which have a domain with a count greater than 1 execute `3-more-than-1.sql`. Just like the previous queries, the output is in more-than-1-domain.txt.

## Step 3 - Creating Endpoint and Fetching data

The `index.go` is your entry point. Download all go-gettable dependencies and then you should set up the connection to your database with the user name usually "postgres", password and database name.

The Endpoint should receive 4 query parameters and find all the spots within a given radius from the table "MY_TABLE" using the requested parameters.

- `Latitude`
- `Longitude`
- `Radius (in meters)`
- `Type (circle or square)`

Results are ordered by distance and if distance between two spots is smaller than 50m, then they are ordered by rating. Then endpoint returns an array of objects including all fields in the data set.

Example if you start the server on local port 3050 with `go run index.go` in the terminal and make a request with these query params `/spots/?latitude=51.51207609999999&longitude=-0.143967&radius=100&type=circle` it returns all the spots within 100metres radius of the latitude and longitude point requested (type is circle here):

```http
{
	"code": "200",
	"spots": [
		{
			"distance": 43.33627,
			"spotId": "d21e893d-112f-4ee9-a56c-66805665d3d9",
			"name": "Browns",
			"website": "browns-restaurants.co.uk",
			"coordinates": "0101000020E6100000F3C3AD275F75C2BFC466367E97C14940",
			"description": "Permanently closed.",
			"rating": 7.882712276466712
		},
		{
			"distance": 31.600357,
			"spotId": "2f3945f9-321c-4c9c-84ba-4a191aaf4a99",
			"name": "Lucknow 49",
			"website": "lucknowldn.com",
			"coordinates": "0101000020E61000001D700278B06AC2BF967840D994C14940",
			"description": "Lucknow 49 aims to be a local staple in the heart of Mayfair, providing an unexpected respite from the area’s formal business lunches. The kitchen will serve Awadhi cuisine, native to the Indian city Lucknow.",
			"rating": 5.725400423471108
		},
		{
			"distance": 0.0001826,
			"spotId": "00028503-4e87-48db-964d-a942a098406e",
			"name": "CELINE London New Bond Street Store",
			"website": "stores.celine.com",
			"coordinates": "0101000020E610000082035ABA826DC2BF1648ABB58BC14940",
			"description": "",
			"rating": 2.800566895178669
		},
		{
			"distance": 28.19163,
			"spotId": "695b9b17-3f1d-4566-a5ae-bf1cc02cd6a2",
			"name": "Pendulum of Mayfair",
			"website": "pendulumofmayfair.co.uk",
			"coordinates": "0101000020E610000024873E0E396CC2BF5AA846F993C14940",
			"description": "A superb antique clocks shop based in the heart of London. A family business that sells and restores clocks and antique furniture.",
			"rating": 2.5777556989500283
		},
		{
			"distance": 66.32787,
			"spotId": "853a2138-b300-4881-8a90-422f4425ae0c",
			"name": "Mayfair Brew House",
			"website": "",
			"coordinates": "0101000020E61000008470BB86CF8CC2BF41518F238CC14940",
			"description": "",
			"rating": 6.766181308639254
		},
		{
			"distance": 79.864586,
			"spotId": "b489951e-b72c-49eb-8feb-fd3312f54ea8",
			"name": "Gallery 1&2 Coffee Bar",
			"website": "",
			"coordinates": "0101000020E61000006A8BC635994CC2BF3B01A83D80C14940",
			"description": "Sotheby's",
			"rating": 3.6169075357763703
		},
		{
			"distance": 80.873276,
			"spotId": "ff3d4668-f5e4-468e-9414-ff4c832ac4cb",
			"name": "Umu",
			"website": "umurestaurant.com",
			"coordinates": "0101000020E6100000D455CBE7267BC2BF0BA2A47675C14940",
			"description": "2 Michelin-starred Kyoto influenced restaurant in Mayfair. On a mission to spread the Ikejime revolution throughout the UK.",
			"rating": 3.6574910683602724
		}
	],
	"message": "Success"
}
```
