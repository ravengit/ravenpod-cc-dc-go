package instrumentation

import (
	"encoding/json"
	"github.com/ravengit/ravenpod-cc-dc-go/datapublisher"
	"github.com/ravengit/ravenpod-cc-dc-go/model"
	rpruntime "github.com/ravengit/ravenpod-cc-dc-go/runtime"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	FABRIC_LEDGER_MODULE_NAME = "FABRIC LEDGER"
	CHANNEL_STATE_DATA        = "CHANNEL STATE DATA"
)

func purifyTransientMapToJSONString(tMap map[string][]byte) string {
	mapClone := make(map[string][]byte)
	for k, v := range tMap {
		if !strings.HasPrefix(k, "rp_") {
			mapClone[k] = v
		}
	}
	tMapInJsonStr := ""
	if len(mapClone) > 0 {
		j, _ := json.Marshal(mapClone)
		tMapInJsonStr = string(j)
	}
	return tMapInJsonStr
}

func InstrumentStateDeletion(blockchainTxnId string, invocationId string, collection string, key string, transientMap map[string][]byte, eventType int, entryTime time.Time) {
	// Get method name
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	splits := strings.Split(frame.Function, ".")
	methodName := splits[len(splits)-1]
	log.Println("[RAVENPOD] Capturing trace. methodNmae, eventType", methodName, eventType)

	// Get collection and key
	if len(collection) == 0 {
		collection = CHANNEL_STATE_DATA
	}

	hasRavenpodData := transientMap["rp_webTxnId"]
	if len(hasRavenpodData) > 0 {
		dataPublisher := datapublisher.GetDataPublisher()
		webTxnId := string(transientMap["rp_webTxnId"])
		ravenpodTxnId := string(transientMap["rp_ravenpodTxnId"])
		accountId := string(transientMap["rp_accountId"])
		channel := string(transientMap["rp_channel"])
		nestLevel, _ := strconv.Atoi(string(transientMap["rp_nestLevel"]))
		sequenceNumber, _ := strconv.Atoi(string(transientMap["rp_sequenceNumber"]))
		mspId := os.Getenv("CORE_PEER_LOCALMSPID")
		traceRecord := model.NewTraceRecord(
			accountId,
			webTxnId,
			ravenpodTxnId,
			blockchainTxnId,
			invocationId,
			channel,
			false,
			sequenceNumber,
			nestLevel,
			FABRIC_LEDGER_MODULE_NAME,
			methodName,
			key,
			purifyTransientMapToJSONString(transientMap),
			collection,
			"",
			eventType,
			"")
		if eventType == model.EVENT_TYPE_ENTRY {
			dataPublisher.PushRecord(traceRecord, accountId)
			nestLevel++
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
		} else {
			timeTaken := time.Now().Sub(entryTime).Milliseconds()
			nestLevel--
			traceRecord.NestLevel = nestLevel
			dataPublisher.PushRecord(traceRecord, accountId)
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
			keyTraceRecord := model.NewKeyTraceRecord(
				accountId,
				webTxnId,
				ravenpodTxnId,
				blockchainTxnId,
				mspId,
				channel,
				collection,
				key,
				"",
				model.OPERATION_TYPE_DELETE,
				"",
				"",
				timeTaken)
			dataPublisher.PushRecord(keyTraceRecord, accountId)
		}
	} else {
		log.Println("[RAVENPOD] Ravenpod context data not found. Did you enable Ravenpod data collector in the web app?")
		return
	}

}

func InstrumentStateSetter(blockchainTxnId string, invocationId string, collection string, key string, value []byte, transientMap map[string][]byte, eventType int, entryTime time.Time) {
	// Get method name
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	splits := strings.Split(frame.Function, ".")
	methodName := splits[len(splits)-1]
	log.Println("[RAVENPOD] Capturing trace. methodNmae, eventType", methodName, eventType)

	// Get collection, key and value
	valueInStr := string(FilterKey(value))
	if len(collection) == 0 {
		collection = CHANNEL_STATE_DATA
	}

	args := struct {
		Key   string
		Value string
	}{
		Key:   key,
		Value: valueInStr,
	}
	argsBuffer, _ := json.Marshal(args)

	hasRavenpodData := transientMap["rp_webTxnId"]
	if len(hasRavenpodData) > 0 {
		dataPublisher := datapublisher.GetDataPublisher()
		webTxnId := string(transientMap["rp_webTxnId"])
		ravenpodTxnId := string(transientMap["rp_ravenpodTxnId"])
		accountId := string(transientMap["rp_accountId"])
		channel := string(transientMap["rp_channel"])
		nestLevel, _ := strconv.Atoi(string(transientMap["rp_nestLevel"]))
		sequenceNumber, _ := strconv.Atoi(string(transientMap["rp_sequenceNumber"]))
		mspId := os.Getenv("CORE_PEER_LOCALMSPID")
		traceRecord := model.NewTraceRecord(
			accountId,
			webTxnId,
			ravenpodTxnId,
			blockchainTxnId,
			invocationId,
			channel,
			false,
			sequenceNumber,
			nestLevel,
			FABRIC_LEDGER_MODULE_NAME,
			methodName,
			string(argsBuffer),
			string(FilterKey([]byte(purifyTransientMapToJSONString(transientMap)))),
			collection,
			"",
			eventType,
			"")
		if eventType == model.EVENT_TYPE_ENTRY {
			dataPublisher.PushRecord(traceRecord, accountId)
			nestLevel++
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
		} else {
			timeTaken := time.Now().Sub(entryTime).Milliseconds()
			nestLevel--
			traceRecord.NestLevel = nestLevel
			traceRecord.Args = ""
			dataPublisher.PushRecord(traceRecord, accountId)
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
			hashTrackingFunc := rpruntime.GetRuntimeOptions().HashTrackingFunc
			var label, hash string = "", ""
			if hashTrackingFunc != nil {
				label, hash = hashTrackingFunc(collection, key, value)
			}
			keyTraceRecord := model.NewKeyTraceRecord(
				accountId,
				webTxnId,
				ravenpodTxnId,
				blockchainTxnId,
				mspId,
				channel,
				collection,
				key,
				valueInStr,
				model.OPERATION_TYPE_WRITE,
				label,
				hash,
				timeTaken)
			dataPublisher.PushRecord(keyTraceRecord, accountId)
		}
	} else {
		log.Println("[RAVENPOD] Ravenpod context data not found. Did you enable Ravenpod data collector in the web app?")
		return
	}

}

func InstrumentStateGetter(blockchainTxnId string, invocationId string, collection string, startKey string, endKey string, returnedValue []byte, transientMap map[string][]byte, eventType int, entryTime time.Time) {
	// Get method name
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	splits := strings.Split(frame.Function, ".")
	methodName := splits[len(splits)-1]
	if eventType == model.EVENT_TYPE_EXIT {
		methodName = splits[len(splits)-2]
	}
	log.Println("[RAVENPOD] Capturing trace. methodNmae, eventType", methodName, eventType)

	// Get collection and key
	if len(collection) == 0 {
		collection = CHANNEL_STATE_DATA
	}
	key := startKey
	if len(endKey) > 0 {
		key += ":" + endKey
	}

	hasRavenpodData := transientMap["rp_webTxnId"]
	if len(hasRavenpodData) > 0 {
		dataPublisher := datapublisher.GetDataPublisher()
		webTxnId := string(transientMap["rp_webTxnId"])
		ravenpodTxnId := string(transientMap["rp_ravenpodTxnId"])
		accountId := string(transientMap["rp_accountId"])
		channel := string(transientMap["rp_channel"])
		nestLevel, _ := strconv.Atoi(string(transientMap["rp_nestLevel"]))
		sequenceNumber, _ := strconv.Atoi(string(transientMap["rp_sequenceNumber"]))
		mspId := os.Getenv("CORE_PEER_LOCALMSPID")
		traceRecord := model.NewTraceRecord(
			accountId,
			webTxnId,
			ravenpodTxnId,
			blockchainTxnId,
			invocationId,
			channel,
			false,
			sequenceNumber,
			nestLevel,
			FABRIC_LEDGER_MODULE_NAME,
			methodName,
			key,
			string(FilterKey([]byte(purifyTransientMapToJSONString(transientMap)))),
			collection,
			string(FilterKey(returnedValue)),
			eventType,
			"")
		if eventType == model.EVENT_TYPE_ENTRY {
			dataPublisher.PushRecord(traceRecord, accountId)
			nestLevel++
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
		} else {
			timeTaken := time.Now().Sub(entryTime).Milliseconds()
			nestLevel--
			traceRecord.NestLevel = nestLevel
			dataPublisher.PushRecord(traceRecord, accountId)
			sequenceNumber++
			transientMap["rp_sequenceNumber"] = []byte(strconv.Itoa(sequenceNumber))
			transientMap["rp_nestLevel"] = []byte(strconv.Itoa(nestLevel))
			hashTrackingFunc := rpruntime.GetRuntimeOptions().HashTrackingFunc
			var label, hash string = "", ""
			if hashTrackingFunc != nil {
				label, hash = hashTrackingFunc(collection, key, returnedValue)
			}
			keyTraceRecord := model.NewKeyTraceRecord(
				accountId,
				webTxnId,
				ravenpodTxnId,
				blockchainTxnId,
				mspId,
				channel,
				collection,
				key,
				string(FilterKey(returnedValue)),
				model.OPERATION_TYPE_READ,
				label,
				hash,
				timeTaken)
			dataPublisher.PushRecord(keyTraceRecord, accountId)
		}
	} else {
		log.Println("[RAVENPOD] Ravenpod context data not found. Did you enable Ravenpod data collector in the web app?")
		return
	}

}
