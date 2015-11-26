package main

import (
"fmt"
"github.com/julienschmidt/httprouter"
"net/http"
"encoding/json"
"strconv"
)

var keyValueMap map[int]string

type KeyResponse struct
{
    Key int `json: key`
    Value string `json: value`
}


type AllKeysResponse struct
{
	Values []KeyResponse
}

func main(){
                    
                    fmt.Println("=========================")
                    keyValueMap = make(map[int]string)

                    mux := httprouter.New()
                    mux.GET("/keys", getAllKeys)
                    mux.GET("/keys/:key_id", getValue)
                    mux.PUT("/keys/:key_id/:value", setKeyValue)

                    server := http.Server{
                        Addr:        "0.0.0.0:3001",
                        Handler: mux,
                    }
           
                    server.ListenAndServe()
            }




func setKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

        keyString:=p.ByName("key_id")
        //Convert to int
        key, err := strconv.Atoi(keyString)
        if err!=nil{
			fmt.Println("Error in conversion")
        }
        value:=p.ByName("value")

        fmt.Println("the obtained key , value pair:",key,value)

        //adding the value to the map
        keyValueMap[key]=value
         //Set the response
        rw.WriteHeader(http.StatusOK)

}


func getValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	keyString:=p.ByName("key_id")
	 key, err := strconv.Atoi(keyString)
        if err!=nil{
			fmt.Println("Error in conversion")
        }

        if _, ok := keyValueMap[key]; !ok {
                   //err_noKey:=errors.New("Supplied key is not found in the system. Please check.")
                    fmt.Println("Supplied key is not found in the system. Please check.")

               }
               value:=keyValueMap[key]



	resp:= KeyResponse{}
	resp.Key=key
	resp.Value=value

	//marshalling into a json

           respJson, err4 := json.Marshal(resp)
           if err4!=nil{
            fmt.Print("Error occcured in marshalling")
        }

        rw.Header().Set("Content-Type","application/json")
        rw.WriteHeader(http.StatusOK)
        fmt.Fprintf(rw, "%s", respJson)

	}


func getAllKeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var Values []KeyResponse

	for key,value := range keyValueMap {

		var temp KeyResponse
		temp.Key=key
		temp.Value=value
		Values=append(Values,temp)
	}

	//Create 
	AllKeys:=AllKeysResponse{Values}

	//marshalling into a json

           respJson, err := json.Marshal(AllKeys)
           if err!=nil{
            fmt.Print("Error occcured in marshalling")
        }

        rw.Header().Set("Content-Type","application/json")
        rw.WriteHeader(http.StatusOK)
        fmt.Fprintf(rw, "%s", respJson)

}