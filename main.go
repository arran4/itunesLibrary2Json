package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"bytes"
	"encoding/json"
	"fmt"
)

func getValue(decoder *xml.Decoder) interface{} {
	for ;; {
		if nexttokenraw, err := decoder.Token(); err != nil {
			log.Printf("Error decoding token: %v %v", err, nexttokenraw)
			return err
		} else {
			//			log.Printf("next token: %v", nexttoken)
			nextthis := ""
			switch nexttoken := nexttokenraw.(type) {
			case xml.StartElement:
				nextthis = nexttoken.Name.Local
			//			log.Printf("Of type: %s", nextthis)
				switch nextthis {
				case "integer":
					var result int
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return result
					}
				case "date":
					var result string
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return result
					}
				case "string":
					var result string
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return result
					}
				case "true":
					var result string
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return true
					}
				case "false":
					var result string
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return false
					}
				case "array":
					var result []interface{}
					for ;; {
						element := getValue(decoder)
						if _, theend := element.(xml.EndElement); theend {
							break;
						}
						result = append(result, element)
					}
					return result
				case "data":
					var result string
					if err := decoder.DecodeElement(&result, &nexttoken); err != nil {
						log.Printf("Got error decoding int: %v", err)
						return err
					} else {
						return result
					}
				case "dict":
					return makeDict(decoder, nexttoken)
				default:
					log.Printf("Unknown type", nextthis)
				}
			case xml.EndElement:
				return nexttokenraw
			case xml.Comment:
//				log.Printf("Comment ignored: %v", nexttoken)
				continue
			case xml.CharData:
//				log.Printf("CharData ignored: %v", nexttoken)
				continue
			default:
//				log.Printf("Unknown token; %s", nexttoken)
				continue;
			}
		}
	}

}

func getKeyValue(decoder *xml.Decoder, token *xml.StartElement) (key string, value interface{}) {
	decoder.DecodeElement(&key, token)
	value = getValue(decoder)
	return
}

func makeDict(decoder *xml.Decoder, token xml.Token) (result map[string]interface{}) {
	if token == nil { return }
	switch token := token.(type) {
	case xml.StartElement:
	case xml.EndElement:
		return
	default:
		log.Printf("Unknown token type; %s", token)
	}
	result = map[string]interface{}{}
//	log.Printf("%s", this)
	for ;; {
		if nexttoken, err := decoder.Token(); err != nil {
			log.Printf("Error decoding token: %v %v", err, nexttoken)
			return
		} else {
//			log.Printf("next token: %v", nexttoken)
			nextthis := ""
			switch nexttoken := nexttoken.(type) {
			case xml.StartElement:
				nextthis = nexttoken.Name.Local
			case xml.EndElement:
				return
			case xml.Comment:
//				log.Printf("Comment ignored: %v", nexttoken)
				continue
			case xml.CharData:
//				log.Printf("CharData ignored: %v", nexttoken)
				continue
			default:
//				log.Printf("Unknown token; %s", nexttoken)
				continue
			}
//			log.Printf("Of type: %s", nextthis)
			switch nextthis {
			case "key":
				if se, ok := nexttoken.(xml.StartElement); ok {
					key, value := getKeyValue(decoder, &se)
					result[key] = value
				}
			case "dict":
				result[nextthis] = makeDict(decoder, nexttoken)
				default:
					log.Printf("Unknown type", nextthis)
			}
		}
	}
	return
}

func main() {
	filename := "iTunes Music Library.xml"
	if data, err := ioutil.ReadFile(filename); err != nil {
		log.Printf("Error opening file: %v", err)
		return
	} else if reader := bytes.NewReader(data); reader == nil {
		log.Printf("Error creating reader")
	} else if decoder := xml.NewDecoder(reader); decoder == nil {
		log.Printf("Error decoding file")
		return
	} else {
		var token xml.Token
		var err error
    	for ;; {
			if token, err = decoder.Token(); err != nil {
				if err.Error() != "EOF" {
					log.Printf("Error getting next token: %v", err)
				}
				return
			} else {
//				log.Printf("%v", token)
				switch token.(type) {
				case xml.StartElement:
					result := makeDict(decoder, token)
					if out, err := json.MarshalIndent(result, "  ", "  "); err != nil {
						log.Printf("Error marshalling map: %v", err)
					} else if err := ioutil.WriteFile("out.json", out, 0644); err != nil {
						log.Printf("Error writing file: %v", err)
					} else {
						fmt.Print("Done")
					}
//					return
				}
			}
		}
	}
}