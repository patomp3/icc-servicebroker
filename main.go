package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type config struct {
	port string
}

/*func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}*/

func main() {
	fmt.Println("ICC Service Broker Started...")

	var env string
	var cfg config

	// For no assign parameter env. using default to Test
	if len(os.Args) > 1 {
		env = os.Args[1]
	} else {
		env = "test"
	}

	// Load configuration
	viper.SetConfigName("app")    // no need to include file extension
	viper.AddConfigPath("config") // set the path of your config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found..." + err.Error())
	} else {
		if strings.ToLower(env) == "production" {
			cfg.port = viper.GetString("production.port")

		} else {
			cfg.port = viper.GetString("development.port")

		}
		fmt.Println("Loading Configuration...")
		fmt.Println("Env = " + env)
		fmt.Println("Server Port=" + cfg.port)
	}

	router := mux.NewRouter()
	router.HandleFunc("/servicebroker", serviceBroker).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+cfg.port, router))

	/*
		// Call Sent to Queue
		q := QueueService{"amqp://admin:admin@172.19.218.104:5672/", "myqueue"}
		ch := q.Connect()
		q.SendMessage(ch, "text/plain", "Test Message 1")
		q.SendMessage(ch, "text/plain", "Test Message 1222")
		q.SendMessage(ch, "text/plain", "Test Message 333")
		q.Close()

		r := QueueService{"amqp://admin:admin@172.19.218.104:5672/", "myqueue2"}
		ch2 := r.Connect()
		r.SendMessage(ch2, "text/plain", "Test Message 2")
		r.Close()
	*/

	//Test call smsservice
	/*var info *DBInfo
	base := "ICC"
	info = GetDBInfo(base)
	fmt.Println("Get DB info:" + "ICC")
	fmt.Println("user =" + info.user)
	fmt.Println("password =" + info.password)
	fmt.Println("dsnurl =" + info.dsnURL)

	fmt.Println("Test Start Execute Query")
	rows, err := ExecuteSQL("QED", "update tmp_aaa set status = 'Y' where cust='patom1'")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("row effects:%d\n", rows)*/

	/*fmt.Println("Test Select statment")
	rows, err := SelectSQL("QED", "select cust, status from tmp_aaa")
	// close database connection after this main function finished
	defer rows.Close()
	if err != nil {

	} else {
		var cust string
		var status string
		for rows.Next() {
			rows.Scan(&cust, &status)
			fmt.Printf("%s, %s\n", cust, status)
		}
	}*/
	//End test call smsservice

	/*fmt.Println("Execute Store no return")
	ExecuteStoreProcedure("QED", "begin PK_IBS_TMP.ProcessExample(:1,:2); end;", "patom", "A")*/

	/*var result int
	fmt.Println("Execute Store return result")
	ExecuteStoreProcedure("QED", "begin PK_IBS_TMP.Calculate(:1,:2,:3); end;", "5", "6", sql.Out{Dest: &result})
	fmt.Println(result)*/

	/*var result driver.Rows
	fmt.Println("Execute Store return cursor")
	bResult := ExecuteStoreProcedure("QED", "begin PK_IBS_TMP.GetReson(:1,:2); end;", "131", sql.Out{Dest: &result})
	if bResult && result != nil {
		values := make([]driver.Value, len(result.Columns()))
		// ok
		for result.Next(values) == nil {
			fmt.Printf("%d %d %s\n", values[0], values[1], values[2])
		}
	}*/

	/*var outID int
	fmt.Println("Execute store return out param")
	bResult := ExecuteStoreProcedure("QED", "begin PK_IBS_IVRLITE.InsertIncommingLog(:1,:2,:3); end;", "12345", "test", sql.Out{Dest: &outID})
	if bResult {
		fmt.Printf("Out Id=%d", outID)
	}*/
	//router := mux.NewRouter()
	//router.HandleFunc("/servicebroker", serviceBroker).Methods("POST")

	//log.Fatal(http.ListenAndServe(":"+cfg.port, router))
}

func serviceBroker(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//Read Json Request
	var req ServiceRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	log.Println("Request incoming...")
	log.Println(req)

	//call recon api
	var res ServiceResponse
	res = ProcessService(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
