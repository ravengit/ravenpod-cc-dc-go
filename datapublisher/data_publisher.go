package datapublisher

import (
	"sync"
	"log"
    "encoding/json"    
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"	
)

const FLUSH_TIMER = 15
const BATCH_SIZE = 10
const STREAM_NAME = "prodDataStream"

type AWSKinesis struct {
	stream          string
	region          string
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

type DataRecordSet struct {
    items             []*kinesis.PutRecordsRequestEntry
	mutex             *sync.Mutex
}

type DataPublisher struct {
	kinesisSession *kinesis.Kinesis
	drs DataRecordSet
}

var (
	doOnce sync.Once
	dataPublisher *DataPublisher // singleton
)

func InitDataPublisher(region string, dataPipelineAccessKey string, dataPipelineSecretAccessKey string) {
	doOnce.Do(func() {
		log.Println("[RAVENPOD] Data pipeline: ", region, dataPipelineAccessKey, dataPipelineSecretAccessKey)
		sess := session.New(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(dataPipelineAccessKey, dataPipelineSecretAccessKey, ""),
		})
		dataPublisher = &DataPublisher{
			kinesisSession: kinesis.New(sess),
			drs: DataRecordSet{mutex: &sync.Mutex{}},
		}
		ticker := time.NewTicker(FLUSH_TIMER * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
			   select {
				case <- ticker.C:
					log.Println("[RAVENPOD] Flush timer wakes up")
                    if len(dataPublisher.drs.items) > 0 {
                        log.Println("[RAVENPOD] No. of records in recordset", len(dataPublisher.drs.items))
                        dataPublisher.flushRecordSet(STREAM_NAME, &dataPublisher.drs.items)
                        log.Println("[RAVENPOD] Force flushing performed successfully")
                    } else {
                        log.Println("[RAVENPOD] No record in recordset")
                    }                
				case <- quit:
					ticker.Stop()
					return
				}
			}
		}()		
    })
}

func GetDataPublisher() *DataPublisher {
    return dataPublisher
}

func (dp *DataPublisher) PushRecord(record interface{}, partitionKey string) {
	dp.drs.mutex.Lock()
	defer dp.drs.mutex.Unlock()

    jsonBuffer, err := json.Marshal(&record)
    if err != nil {
        log.Println("Error marshalling the record data", err)
        return
	}

	item := &kinesis.PutRecordsRequestEntry{
		Data:         jsonBuffer,
		PartitionKey: aws.String(partitionKey),
	}

	dp.drs.items = append(dp.drs.items, item)

	// log.Println("Pipeline items: ", dp.trs.items)

	if len(dp.drs.items) >= BATCH_SIZE {
		dp.flushRecordSet(STREAM_NAME, &dp.drs.items)
	}

}

func (dp *DataPublisher) flushRecordSet(streamName string, items *[]*kinesis.PutRecordsRequestEntry) {
	_, err := dp.kinesisSession.PutRecords(&kinesis.PutRecordsInput{
		Records:    *items,
		StreamName: aws.String(streamName),
	})
	if err != nil {
		log.Fatal("Error pushing recordset to data pipeline", err)
	} else {
		log.Println("Recordset pushed to data pipeline.", streamName)
		*items = nil
	}
}