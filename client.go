package main

import (
"fmt"
"hash/fnv"
"container/ring"
"strconv"
"net/http"
"io/ioutil"
"encoding/json"
)

type KeyValuePair struct
{
	Key int `json: key`
	Value string `json: value`
}

type AllKeysResponse struct
{
	Values []KeyValuePair
}


func main(){
		//3 10  7 works    7 6 10

	//Saras- 5,0,8
	//Rad: 5, 8, 11 
	var caches [13]string
	server1:="http://localhost:3000"
	server2:="http://localhost:3001"
	server3:="http://localhost:3002"

	//Calling server hashing
	serverIndex1:=serverHash(server1)
	serverIndex2:=serverHash(server2)
	serverIndex3:=serverHash(server3)

	caches[serverIndex1]=server1
	caches[serverIndex2]=server2
	caches[serverIndex3]=server3


	keyInput:=[]int{1,2,3,4,5,6,7,8,9,10}
	valueInput:=[]string{"a","b","c","d","e","f","g","h","i","j"}


	r := ring.New(len(caches))


	//for loop for putting the values into the ring
	for i := 0; i < r.Len(); i++ {
		r.Value = caches[i]
		r = r.Next()
	}

	p:=r


	fmt.Println("Sharding and putting the key value pairs:")

	//Get an index in the ring
	for i:=0;i<len(keyInput);i++{

		var hashValue uint32
		hashValue = hash(strconv.Itoa(keyInput[i]))

		a := int(hashValue)

		//get the index to be inserted into
		index:=a%(len(caches))
	//fmt.Println("the value of index:",index)

		for j:=0;j<index;j++{
			r= r.Next()
		}

	//if empty move fwd
		for r.Value=="" {
			r= r.Next()
		}
		fmt.Println("Key: ",keyInput[i]," Value: ",valueInput[i],"Server chosen:",r.Value)

	//Forming the url to send
		serverName,_:=(r.Value).(string) 
		url:=serverName+"/keys/"+strconv.Itoa(keyInput[i])+"/"+valueInput[i]
	//fmt.Println(url)
		req1, errReqC := http.NewRequest("PUT", url, nil)
		if errReqC!=nil{
			errMsg:="Request creation error"
			fmt.Println(errMsg)
         //errorCheck(errMsg,rw)
			return
		}

		client := &http.Client{}
		resp, errClient := client.Do(req1)
		if errClient != nil {
			errMsg:="Request creation error.Check server side."
			fmt.Println(errMsg)
        //errorCheck(errMsg,rw)
			return
		}
		defer resp.Body.Close()


	//Reset to the beginning of the circular array 
		r=p
	}





//GETTING the values
	fmt.Println(" ")
	fmt.Println("Getting the key Value pairs: ")

	for i:=0;i<len(keyInput);i++{
		var hashValue uint32
		hashValue = hash(strconv.Itoa(keyInput[i]))

		a := int(hashValue)

		//get the index to be inserted into
		index:=a%(len(caches))
	//fmt.Println("the value of index:",index)

		for j:=0;j<index;j++{
			r= r.Next()
		}

	//if empty move fwd
		for r.Value=="" {
			r= r.Next()
		}

	//Forming the url to send
		serverName,_:=(r.Value).(string) 
		url:=serverName+"/keys/"+strconv.Itoa(keyInput[i])
	//fmt.Println(url)
		resp, err := http.Get(url);
		if err != nil {
			fmt.Println("Get the key error.Check the server side.")
                      
		}

		defer resp.Body.Close()
		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			fmt.Println("Get request error")

		}
		var result KeyValuePair 
                //Unmarshall the response into a json
		err2:=json.Unmarshal(body,&result)
		if err2 != nil {
			fmt.Println("Get request error")
		}

		fmt.Println(result.Key,"=>",result.Value,"    from "+ serverName)
	//Reset to the beginning of the circular array 
		r=p
	}



	fmt.Println("")
	//get all values
	url1:= server1+"/keys/"
	fmt.Println("Printing all the key value pairs from server: "+server1)
	PrintAllServerKeys(url1)
	fmt.Println("")

	url2:= server2+"/keys/"
	fmt.Println("Printing all the key value pairs from server: "+server2)
	PrintAllServerKeys(url2)
	fmt.Println("")

	url3:= server3+"/keys/"
	fmt.Println("Printing all the key value pairs from server: "+server3)
	PrintAllServerKeys(url3)
	fmt.Println("")
	
}



func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}


func serverHash(url string) int {
	hashValue:=hash(url)
	index:=(hashValue*54321)%13
	return int(index)

}

func PrintAllServerKeys(url string) {
	resp, err := http.Get(url);
	if err != nil {
		fmt.Println("Get all keys request error.Check server side")
                        //err_get:=errors.New("Get error")
	}

	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Get all keys request error")

	}
	var result1 AllKeysResponse 
                //Unmarshall the response into a json
	err2:=json.Unmarshal(body,&result1)
	if err2 != nil {
		fmt.Println("Get all keys unmarshall error")
	}


	for l:=0;l<len(result1.Values);l++{
		temp:=result1.Values[l]
		fmt.Println(temp.Key,"=>",temp.Value)
	}
}

