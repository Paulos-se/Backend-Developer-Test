package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)
const(
	
	DB_USER ="DATABASE_USER"
    DB_PASSWORD="YOUR_PASSWORD_HERE"
    DB_NAME="DATABASE_NAME"

    
)

func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db,err := sql.Open("postgres", dbinfo)

    checkErr(err)
          
    return db
    
}

type Spot struct{
	Distance float32 `json:"distance"`
	Id string `json:"spotId"`
    SpotName string `json:"name"`
	SpotWebsite string `json:"website"`
	SpotCoordinates string `json:"coordinates"`
  	SpotDescription string `json:"description"`
  	SpotRating float64 `json:"rating"`
}

type JsonResponse struct {
    Type    string `json:"code"`
    Data    []Spot `json:"spots"`
    Message string `json:"message"`
}


func main() {

    // Init the mux router
    router := mux.NewRouter()


    // Get spots endpoint
    router.HandleFunc("/spots/", GetSpots).Methods("GET")

    // serve the app
    fmt.Println("Listening on port 3050")
    log.Fatal(http.ListenAndServe(":3050", router))
}

func printMessage(message string){
    fmt.Println(" ")
    fmt.Println(message)
    fmt.Println(" ")
}

func checkErr(err error){
    if err!=nil{
	panic(err)
    }
}

func GetSpots(w http.ResponseWriter, r *http.Request){
    db:=setupDB()

	var spots[]Spot
	q:=r.URL.Query()

	long:=q.Get("longitude")
	lat:=q.Get("latitude")
	rad:=q.Get("radius")
	shape:=q.Get("type")

	var response = JsonResponse{}

	longitude, err := strconv.ParseFloat(long, 32); 
	if err != nil {
    w.WriteHeader(http.StatusBadRequest) 
	response =JsonResponse{Type:"400", Data:spots, Message: "Incorrect param Longitude"} 
	json.NewEncoder(w).Encode(response)
	return
  } 
    latitude, err := strconv.ParseFloat(lat, 64); 
	
if err != nil {
    w.WriteHeader(http.StatusBadRequest)  
	response =JsonResponse{Type:"400", Data:spots, Message: "Incorrect param Latitude"}
	json.NewEncoder(w).Encode(response)
	return
  } 
	radius, err := strconv.ParseFloat(rad, 32); 
	
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)  
	response =JsonResponse{Type:"400", Data:spots, Message: "Incorrect param Radius"}
	json.NewEncoder(w).Encode(response)
	return
  } 
	
	var sqlquery=``;

    printMessage("Getting spots...")

	if strings.ToLower(shape)=="square" {

    sqlquery=`SELECT id,name,(CASE WHEN website IS NULL THEN
   ''
    ELSE
    website
   END),coordinates,(CASE WHEN description IS NULL THEN
   ''
    ELSE
    description
   END),rating,ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography) FROM "MY_TABLE" 
  WHERE 
   ST_Within((coordinates::geometry),ST_Envelope(ST_GeomFromText(CONCAT('LINESTRING	(', 
      $1 - ($3 * ((1 / ((2 * PI() / 360) * 6378.137)) / 1000))/ COS($2 * (PI() / 180)), ' ', 
      $2 - ($3 * ((1 / ((2 * PI() / 360) * 6378.137)) / 1000)), ',', 
      $1 + ($3 * ((1 / ((2 * PI() / 360) * 6378.137)) / 1000))/ COS($2 * (PI() / 180)), ' ', 
      $2 + ($3 * ((1 / ((2 * PI()/ 360) * 6378.137)) / 1000)), ')'), 4326)))
	ORDER BY 
	(CASE WHEN ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography)>50 then ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography) 
	ELSE -1 * rating END)`;

    } else if strings.ToLower(shape)=="circle"{

    sqlquery=`SELECT id,name,(CASE WHEN website IS NULL THEN
   ''
    ELSE
    website
	END),coordinates,(CASE WHEN description IS NULL THEN
   ''
    ELSE
    description
	END),rating,ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography) 
	FROM "MY_TABLE"
	WHERE ST_DWithin(coordinates, Geography (ST_MakePoint($1,$2)),$3) ORDER BY(CASE WHEN ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography)>=50 THEN ST_Distance(ST_MakePoint($1,$2)::geography,ST_AsText(coordinates)::geography) else -1 * rating END)`;

    }else{

	w.WriteHeader(http.StatusBadRequest)  
	response =JsonResponse{Type:"400", Data:spots, Message: "Incorrect param Type"}
	json.NewEncoder(w).Encode(response)
	return
	}
	
		
	

    rows, err :=db.Query(sqlquery,longitude,latitude,radius)

	if err != nil {
    w.WriteHeader(http.StatusNotFound)  
	response =JsonResponse{Type:"404", Data:spots, Message: "Not Found"}
	json.NewEncoder(w).Encode(response)
	return
  } 

    

    

    for rows.Next(){
		var distance float32
		var id string
        var spotName string
		var spotWebsite string
		var spotCoordinates string
		var spotDescription string 
		var spotRating float64 

        err=rows.Scan(&id, &spotName, &spotWebsite, &spotCoordinates, &spotDescription, &spotRating,&distance)
	
	
        

        spots=append(spots,Spot{Id:id,SpotName:spotName, SpotWebsite: spotWebsite, SpotCoordinates:spotCoordinates ,SpotDescription: spotDescription, SpotRating: spotRating,Distance:distance})
	}
	
	response =JsonResponse{Type:"200", Data:spots, Message: "Success"}
  

        json.NewEncoder(w).Encode(response)


}


