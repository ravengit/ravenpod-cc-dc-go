package instrumentation

import (
	"encoding/json"
	"github.com/ravengit/ravenpod-cc-dc-go/runtime"
	"github.com/tidwall/gjson"
	"log"
)

func FilterKey(value []byte) []byte {

	runtimeOptions := runtime.GetRuntimeOptions()
	disableParam := runtimeOptions.DisableParam
	disableParamExceptions := runtimeOptions.DisableParamExceptions

	tMap := string(value)
	if !disableParam {
		return value
	}
	if len(disableParamExceptions) == 0 {
		log.Println("[RAVENPOD] No exceptions in parameter filtering. All parameters are filtered.")
		return []byte("")
	}
	if len(tMap) == 0 {
		return []byte("")
	}
	if !gjson.Valid(tMap) {
		log.Println("[RAVENPOD] Invalid JSON found in parameter filtering.", tMap)
		return []byte("")
	} else {
		log.Println("[RAVENPOD] Parameter filtering", tMap, disableParamExceptions)
		var newMap = make(map[string]interface{})
		values := gjson.GetMany(tMap, disableParamExceptions...)
		for i, s := range values {
			if s.Value() != nil {
				newMap[disableParamExceptions[i]] = s.Value()
			}
		}
		resultInBytes, _ := json.Marshal(newMap)
		log.Println("[RAVENPOD] Parameter filtering result", string(resultInBytes))
		return resultInBytes
	}
}
