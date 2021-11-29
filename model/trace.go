package model

import (
    guuid "github.com/google/uuid"
    "os"
    "time"
)

type Trace struct {
    Id string `json:"id"`
    AccountId string `json:"account_id"`
    WebTxnId string `json:"web_txn_id"`
    RavenpodTxnId string `json:"ravenpod_txn_id"`
    BlockchainTxnId string `json:"blockchain_txn_id"`
    InvocationId string `json:"invocation_id"`
    Hostname string `json:"hostname"`
    Channel string `json:"channel"`
    IsTransactionStart bool `json:"is_transaction_start"`
    SequenceNumber int `json:"sequence_number"`
    NestLevel int `json:"nest_level"`
    ModuleName string `json:"module_name"`
    FunctionName string `json:"function_name"`
    Args string `json:"args"`
    ReturnedResult string `json:"returned_result"`
    EventType int `json:"event_type"`
    EventData string `json:"event_data"`
    EventTime string `json:"event_time"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    TransientMap string `json:"transient_map"`
    Collection string `json:"collection"`
}

type KeyTrace struct {
    Id string `json:"id"`
    AccountId string `json:"account_id"`
    WebTxnId string `json:"web_txn_id"`
    RavenpodTxnId string `json:"ravenpod_txn_id"`
    BlockchainTxnId string `json:"blockchain_txn_id"`
    MspId string `json:"msp_id"`
    Hostname string `json:"hostname"`
    ChannelName string `json:"channel_name"`
    KeyContent string `json:"key_content"`
    ValueContent string `json:"value_content"`
    OperationType int `json:"operation_type"`
    TimeTaken int64 `json:"time_taken"`
    ReadAt string `json:"read_at"`
    WrittenAt string `json:"written_at"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
    Collection string `json:"collection"`
    PurgedAt string `json:"purged_at"`
    TrackedHashLabel string `json:"tracked_hash_label"`
    TrackedHashValue string `json:"tracked_hash_value"`
}

const (
	EVENT_TYPE_ENTRY = 1
	EVENT_TYPE_EXIT  = 2
    OPERATION_TYPE_READ = 1000
    OPERATION_TYPE_WRITE = 2000
    OPERATION_TYPE_DELETE = 3000
)

func NewKeyTraceRecord(accountId string,
    webTxnId string,
    ravenpodTxnId string,
    blockchainTxnId string,
    mspId string,
    channelName string,
    collection string,
    keyContent string,
    valueContent string,
    operationType int,
    trackedHashLabel string,
    trackedHashValue string,
    timeTaken int64) KeyTrace {

    currentTime := time.Now()
    format := "2006-01-02 15:04:05.000"
    hostname, _ := os.Hostname()
    
    readAt := ""
    writtenAt := ""
    purgedAt := ""

    switch operationType {
    case OPERATION_TYPE_READ:
        readAt = currentTime.Format(format)
    case OPERATION_TYPE_WRITE:
        writtenAt = currentTime.Format(format)
    case OPERATION_TYPE_DELETE:
        purgedAt = currentTime.Format(format)
    }
    
    keyTraceRecord := KeyTrace{
        Id: guuid.New().String(),
        AccountId: accountId,
        WebTxnId: webTxnId,
        RavenpodTxnId: ravenpodTxnId,
        BlockchainTxnId: blockchainTxnId,
        MspId: mspId,
        Hostname: hostname,
        ChannelName: channelName,
        Collection: collection,
        KeyContent: keyContent,
        ValueContent: valueContent,
        OperationType: operationType,
        TimeTaken: timeTaken,
        ReadAt: readAt,
        WrittenAt: writtenAt,
        PurgedAt: purgedAt,
        TrackedHashLabel: trackedHashLabel,
        TrackedHashValue: trackedHashValue,
        CreatedAt: currentTime.Format(format)}

    return keyTraceRecord

}



func NewTraceRecord(accountId string,
    webTxnId string,
    ravenpodTxnId string,
    blockchainTxnId string,
    invocationId string,
    channel string,
    isTransactionStart bool,
    sequenceNumber int, 
    nestLevel int, 
    moduleName string, 
    functionName string, 
    args string, 
    transientMap string, 
    collection string, 
    returnedResult string, 
    eventType int, 
    eventData string) Trace {

    currentTime := time.Now()
    format := "2006-01-02 15:04:05.000"
    hostname, _ := os.Hostname()

	traceRecord := Trace{
        Id: guuid.New().String(),
        AccountId: accountId,
        WebTxnId: webTxnId,
        RavenpodTxnId: ravenpodTxnId,
        BlockchainTxnId: blockchainTxnId,
        InvocationId: invocationId,
        Hostname: hostname,
        Channel: channel,
        IsTransactionStart: isTransactionStart,
        SequenceNumber: sequenceNumber,
        NestLevel: nestLevel,
        ModuleName: moduleName,
        FunctionName: functionName,
        Args: args,
        ReturnedResult: returnedResult,
        EventType: eventType,
        EventData: eventData,
        EventTime: currentTime.Format(format),
        CreatedAt: currentTime.Format(format),
        UpdatedAt: "",
        TransientMap: transientMap,
        Collection: collection}

    return traceRecord

}
