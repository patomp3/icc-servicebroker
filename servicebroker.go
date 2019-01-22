package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

// ServiceRequest for mapping request of service
type ServiceRequest struct {
	RequestTransID string `json:"request_trans_id"`
	TvsCustomerID  int    `json:"tvs_customer_id"`
	TvsReferenceID string `json:"tvs_reference_id"`
	OrderType      string `json:"order_type"`
	OrderMessage   string `json:"order_message"`
	ByChannel      string `json:"by_channel"`
	ByUser         string `json:"by_user"`
}

// ServiceResponse for mapping response of service
type ServiceResponse struct {
	OrderTransID     string `json:"order_trans_id"`
	RequestTransID   string `json:"request_trans_id"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

// ProcessService to redirect message to each queue by config
func ProcessService(req ServiceRequest) ServiceResponse {
	var myReturn ServiceResponse

	// Insert Log Transaction
	var outID int64
	var oRs driver.Rows
	bResult := ExecuteStoreProcedure("QED", "begin PK_NPF_BROKERSERV.ProcessServiceTrans(:1,:2,:3,:4,:5,:6,:7,:8,:9); end;",
		req.TvsCustomerID, req.TvsReferenceID, req.OrderType, req.OrderMessage, req.ByChannel,
		req.ByUser, req.RequestTransID, sql.Out{Dest: &outID}, sql.Out{Dest: &oRs})

	if !bResult && outID <= 0 {
		//fmt.Printf("outId=%d", outID)
		// error insert transaction
		myReturn.ErrorCode = strconv.FormatInt(800, 10)
		myReturn.ErrorDescription = viper.GetString("errorcode" + myReturn.ErrorCode)
		myReturn.OrderTransID = strconv.FormatInt(outID, 10)
		myReturn.RequestTransID = req.RequestTransID
		return myReturn
	}

	// get result from cursor to process transaction flow
	var orderTransID int64
	var orderID int64
	var flowID int64
	var flowType string
	var flowName string
	if bResult && oRs != nil {
		values := make([]driver.Value, len(oRs.Columns()))
		// ok
		for oRs.Next(values) == nil {
			orderTransID = values[0].(int64)
			orderID = values[1].(int64)
			flowID = values[2].(int64)
			flowType = values[3].(string)
			flowName = values[4].(string)
			//fmt.Println(orderID)
			fmt.Printf("%d %d %d %s %s\n", orderTransID, orderID, flowID, flowType, flowName)

			//Process Order
			queueURL := viper.GetString(flowName + ".queueurl")
			queueName := viper.GetString(flowName + ".queuename")
			fmt.Printf("%s %s\n", queueName, queueURL)

			q := Send{queueURL, queueName}
			ch := q.Connect()
			q.SendMessage(ch, req.ByUser, strconv.FormatInt(orderTransID, 10), "application/json", req.OrderMessage)
			q.Close()

		}
	}

	myReturn.ErrorCode = strconv.FormatInt(0, 10)
	myReturn.ErrorDescription = viper.GetString("errorcode" + myReturn.ErrorCode)
	myReturn.OrderTransID = strconv.FormatInt(outID, 10)
	myReturn.RequestTransID = req.RequestTransID
	return myReturn
}

func submitToQueue() {

}
