package lib

import (
	"encoding/json"
	"strings"

	"github.com/beego/beego"
	"github.com/beego/beego/validation"
)

//CreateValidationMap ranslates validation structure to map
// that can be easly presented in template
func CreateValidationMap(valid validation.Validation) map[string]map[string]string {
	v := make(map[string]map[string]string)
	/*
			{
				"email": {
					"Requrired" : "Can not be empty"
				},
				"password" :{

			  }
		  }
	*/
	for _, err := range valid.Errors {
		beego.Notice(err.Key, err.Message)
		k := strings.Split(err.Key, ".")
		var field, errorType string
		if len(k) > 1 {
			field = k[0]
			errorType = k[1]
		} else {
			field = err.Key
			errorType = " "
		}
		beego.Error(field)
		if _, ok := v[field]; !ok {
			v[field] = make(map[string]string)
		}
		v[field][errorType] = err.Message
	}
	return v

}

//Dump any structure as json string
func Dump(obj interface{}) {
	result, _ := json.MarshalIndent(obj, "", "\t")
	beego.Debug(string(result))
}

//CopyStruct serializes src and tries to deserialize it to dst
func CopyStruct(src interface{}, dst interface{}) error {
	jsonString, err := json.Marshal(&src)
	if err != nil {
		beego.Error("Unable to marshal object")
		return err
	}

	if err := json.Unmarshal([]byte(jsonString), &dst); err != nil {
		beego.Error("Unable to unmarshal object")
		return err
	}

	return nil
}
