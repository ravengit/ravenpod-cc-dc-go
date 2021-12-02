package model

import (
    guuid "github.com/google/uuid"
    "os"
    "time"
)
    
type Trace struct {
    TableName string `json:"tableName"`
    Id string `json:"id"`
    AccountId string `json:"accountId"`
    WebTxnId string `json:"webTxnId"`
    RavenpodTxnId string `json:"ravenpodTxnId"`
    BlockchainTxnId string `json:"blockchainTxnId"`
    InvocationId string `json:"invocationId"`
    Hostname string `json:"hostname"`
    Channel string `json:"channel"`
    IsTransactionStart bool `json:"isTransactionStart"`
    SequenceNumber int `json:"sequenceNumber"`
    NestLevel int `json:"nestLevel"`
    ModuleName string `json:"moduleName"`
    FunctionName string `json:"functionName"`
    Args string `json:"args"`
    ReturnedResult string `json:"returnedResult"`
    EventType int `json:"eventType"`
    EventData string `json:"eventData"`
    EventTime string `json:"eventTime"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
    TransientMap string `json:"transientMap"`
    Collection string `json:"collection"`
}

type KeyTrace struct {
    TableName string `json:"tableName"`
    Id string `json:"id"`
    AccountId string `json:"accountId"`
    WebTxnId string `json:"webTxnId"`
    RavenpodTxnId string `json:"ravenpodTxnId"`
    BlockchainTxnId string `json:"blockchainTxnId"`
    MspId string `json:"mspId"`
    Hostname string `json:"hostname"`
    ChannelName string `json:"channelName"`
    KeyContent string `json:"keyContent"`
    ValueContent string `json:"valueContent"`
    OperationType int `json:"operationType"`
    TimeTaken int64 `json:"timeTaken"`
    ReadAt string `json:"readAt"`
    WrittenAt string `json:"writtenAt"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
    Collection string `json:"collection"`
    PurgedAt string `json:"purgedAt"`
    TrackedHashLabel string `json:"trackedHashLabel"`
    TrackedHashValue string `json:"trackedHashValue"`
}

const (
	EVENT_TYPE_ENTRY = 1
	EVENT_TYPE_EXIT  = 2
    OPERATION_TYPE_READ = 1000
    OPERATION_TYPE_WRITE = 2000
    OPERATION_TYPE_DELETE = 3000
    TRACE_TABLE_NAME = "traces"
    KEY_TRACE_TABLE_NAME = "key_traces"
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
        TableName: KEY_TRACE_TABLE_NAME,
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
        TableName: TRACE_TABLE_NAME,
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
