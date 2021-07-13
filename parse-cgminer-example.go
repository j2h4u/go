package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/savaki/jq"
)

//
// types
//

type jsonField struct {
	ExtractOperation, RawValue string
	SanitizedValue             float64
	ErrorDescription           string
}

type jsonFieldCollection map[string]jsonField

//
// functions
//

func isNum(s string) bool {
	// always true for int and float
	dotFound := false
	for _, v := range s {
		if v == '.' {
			if dotFound {
				return false
			}
			dotFound = true
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func unquoteString(s string) string {
	var unquotedString string
	var err error

	unquotedString, err = strconv.Unquote(s)

	if err == nil {
		return unquotedString
	} else {
		return s // as is
	}

}

func (j jsonFieldCollection) Populate(jsonRaw []byte) {
	var err error
	var hash string
	var thisField jsonField

	if j == nil {
		panic("j is nil")
	}

	for hash, thisField = range j {
		op, _ := jq.Parse(thisField.ExtractOperation) // create an Operation (!!! no error checking at all)
		value, _ := op.Apply(jsonRaw)                 // apply an Operation (!!! no error checking at all)

		thisField.RawValue = string(value)
		valueAsUnquotedString := unquoteString(thisField.RawValue)

		switch {
		case valueAsUnquotedString == "":
			thisField.SanitizedValue = 0
			thisField.ErrorDescription = "<empty string>"
		case valueAsUnquotedString == "null":
			thisField.SanitizedValue = 0
			thisField.ErrorDescription = "<null>"
		case isNum(valueAsUnquotedString):
			thisField.SanitizedValue, err = strconv.ParseFloat(valueAsUnquotedString, 64)
			if err != nil {
				thisField.SanitizedValue = 0
				thisField.ErrorDescription = "<cannot parse to float>"
			}
		default:
			thisField.SanitizedValue = 0
			thisField.ErrorDescription = "<NaN>"
		}

		// populate the map (passed by reference, of course)
		j[hash] = jsonField{
			RawValue:         thisField.RawValue,
			SanitizedValue:   thisField.SanitizedValue,
			ErrorDescription: thisField.ErrorDescription,
		}
	}
}

func (j jsonFieldCollection) PrettyPrint() {
	if j == nil {
		panic("j is nil")
	}

	// let's sort hash keys
	keys := make([]string, 0, len(j))
	for k := range j {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// and print them orderly
	for _, hash := range keys {
		thisField := j[hash]
		fmt.Printf("%-15s = %7.2f, raw value %s", hash, thisField.SanitizedValue, thisField.RawValue)
		if len(thisField.ErrorDescription) != 0 {
			fmt.Printf(", parsing error %s", thisField.ErrorDescription)
		}
		fmt.Printf("\n")
	}
}

//
// main
//

func main() {
	// an input json
	jsonRaw := []byte(`{"STATUS":[{"STATUS":"S","When":1626178788,"Code":70,"Msg":"BMMiner stats","Description":"bmminer 1.0.0"}],"STATS":[{"BMMiner":"2.0.0 rwglr","Miner":"30.0.1.3","CompileTime":"Fri Jul 9 18:31:39 UTC 2021","Type":"Antminer S9 Hiveon"},{"STATS":0,"ID":"BC50","Elapsed":4271,"Calls":0,"Wait":0,"Max":0,"Min":99999999,"GHS 5s":"5904.657","GHS av":5847.82,"miner_count":2,"frequency":"412","fan_num":2,"fan1":0,"fan2":0,"fan3":0,"fan4":0,"fan5":4080,"fan6":4920,"fan7":0,"fan8":0,"temp_num":2,"temp1":0,"temp2":0,"temp3":0,"temp4":0,"temp5":0,"temp6":53,"temp7":52,"temp8":0,"temp9":0,"temp10":0,"temp11":0,"temp12":0,"temp13":0,"temp14":0,"temp15":0,"temp16":0,"temp2_1":0,"temp2_2":0,"temp2_3":0,"temp2_4":0,"temp2_5":0,"temp2_6":68,"temp2_7":67,"temp2_8":0,"temp2_9":0,"temp2_10":0,"temp2_11":0,"temp2_12":0,"temp2_13":0,"temp2_14":0,"temp2_15":0,"temp2_16":0,"temp3_1":0,"temp3_2":0,"temp3_3":0,"temp3_4":0,"temp3_5":0,"temp3_6":0,"temp3_7":0,"temp3_8":0,"temp3_9":0,"temp3_10":0,"temp3_11":0,"temp3_12":0,"temp3_13":0,"temp3_14":0,"temp3_15":0,"temp3_16":0,"freq_avg1":0,"freq_avg2":0,"freq_avg3":0,"freq_avg4":0,"freq_avg5":0,"freq_avg6":412,"freq_avg7":411.46,"freq_avg8":0,"freq_avg9":0,"freq_avg10":0,"freq_avg11":0,"freq_avg12":0,"freq_avg13":0,"freq_avg14":0,"freq_avg15":0,"freq_avg16":0,"total_rateideal":5914.09,"total_freqavg":411.73,"total_acn":126,"total_rate":5904.65,"chain_rateideal1":0,"chain_rateideal2":0,"chain_rateideal3":0,"chain_rateideal4":0,"chain_rateideal5":0,"chain_rateideal6":2958.98,"chain_rateideal7":2955.1,"chain_rateideal8":0,"chain_rateideal9":0,"chain_rateideal10":0,"chain_rateideal11":0,"chain_rateideal12":0,"chain_rateideal13":0,"chain_rateideal14":0,"chain_rateideal15":0,"chain_rateideal16":0,"temp_max":68,"Device Hardware%":0.0001,"no_matching_work":7,"chain_acn1":0,"chain_acn2":0,"chain_acn3":0,"chain_acn4":0,"chain_acn5":0,"chain_acn6":63,"chain_acn7":63,"chain_acn8":0,"chain_acn9":0,"chain_acn10":0,"chain_acn11":0,"chain_acn12":0,"chain_acn13":0,"chain_acn14":0,"chain_acn15":0,"chain_acn16":0,"chain_acs1":"","chain_acs2":"","chain_acs3":"","chain_acs4":"","chain_acs5":"","chain_acs6":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo","chain_acs7":" oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo oooooooo ooooooo","chain_acs8":"","chain_acs9":"","chain_acs10":"","chain_acs11":"","chain_acs12":"","chain_acs13":"","chain_acs14":"","chain_acs15":"","chain_acs16":"","chain_hw1":0,"chain_hw2":0,"chain_hw3":0,"chain_hw4":0,"chain_hw5":0,"chain_hw6":6,"chain_hw7":1,"chain_hw8":0,"chain_hw9":0,"chain_hw10":0,"chain_hw11":0,"chain_hw12":0,"chain_hw13":0,"chain_hw14":0,"chain_hw15":0,"chain_hw16":0,"chain_rate1":"","chain_rate2":"","chain_rate3":"","chain_rate4":"","chain_rate5":"","chain_rate6":"2948.61","chain_rate7":"2956.04","chain_rate8":"","chain_rate9":"","chain_rate10":"","chain_rate11":"","chain_rate12":"","chain_rate13":"","chain_rate14":"","chain_rate15":"","chain_rate16":"","chain_xtime6":"{X5=1,X22=1,X43=1}","chain_xtime7":"{X54=1}","chain_offside_6":"0","chain_offside_7":"0","chain_opencore_6":"0","chain_opencore_7":"0","miner_version":"30.0.1.3","chain_power1":0,"chain_power2":0,"chain_power3":0,"chain_power4":0,"chain_power5":0,"chain_power6":234.84,"chain_power7":221.94,"chain_power8":0,"chain_power9":0,"chain_power10":0,"chain_power11":0,"chain_power12":0,"chain_power13":0,"chain_power14":0,"chain_power15":0,"chain_power16":0,"chain_power":"456.78 (AB)"}],"id":1}`)

	// prepare a map--init and then populate with extract operators
	jsonFieldCollection := make(jsonFieldCollection)
	jsonFieldCollection["GHS 5s"] = jsonField{ExtractOperation: ".STATS.[1].GHS 5s"}           // "5904.657"
	jsonFieldCollection["GHS av"] = jsonField{ExtractOperation: ".STATS.[1].GHS av"}           // 5847.82
	jsonFieldCollection["chain power"] = jsonField{ExtractOperation: ".STATS.[1].chain_power"} // "456.78 (AB)"
	jsonFieldCollection["chain_rate1"] = jsonField{ExtractOperation: ".STATS.[1].chain_rate1"} // ""

	jsonFieldCollection.Populate(jsonRaw)
	jsonFieldCollection.PrettyPrint()
}
