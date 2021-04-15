package xml

import (
	"fmt"
	"strconv"
	"strings"
)

//
//
// XmlGetJobAttrMap
//
//
func XmlGetJobAttrs(mv interface{}) ([]string, int) {
	attrMetaMap := make(map[string]interface{})
	var xmlAttrs []string

	XmlTraverse(-1, &attrMetaMap, mv)

	maxDepth := 0
	for _, value := range attrMetaMap {
		valueMap := value.(map[string]interface{})
		if valueMap["depth"].(int) > maxDepth {
			maxDepth = valueMap["depth"].(int)
		}
	}

	for _, value := range attrMetaMap {
		valueMap := value.(map[string]interface{})
		if valueMap["depth"].(int) == maxDepth {
			key := fmt.Sprintf("%v", valueMap["key"])
			xmlAttrs = append(xmlAttrs, key)
		}
	}

	return xmlAttrs, maxDepth
}

//
//
// XmlTraverse
//
//
func XmlTraverse(depth int, attrMetaMap *map[string]interface{}, mv interface{}) bool {
	if _, ok := mv.(map[string]interface{}); ok {
		for key, value := range mv.(map[string]interface{}) {
			if XmlTraverse(depth+1, attrMetaMap, value) {
				metaMap := make(map[string]interface{})
				metaMap["depth"] = depth + 1
				metaMap["key"] = key
				(*attrMetaMap)[key+strconv.Itoa(depth+1)] = metaMap
			}
		}
	} else if _, ok := mv.([]interface{}); ok {
		for _, value := range mv.([]interface{}) {
			XmlTraverse(depth, attrMetaMap, value)
		}
	} else if _, ok := mv.(string); ok {
		return true
	}

	return true
}

//
//
// XmlGetJobMap
//
//
func XmlGetJobMap(mv interface{}, jobIndex string, depth int) map[string]interface{} {
	jobMap := make(map[string]interface{})

	XmlJobCrawl(-1, depth, jobIndex, &jobMap, mv)

	return jobMap
}

//
//
// XmlJobCrawl
//
//
func XmlJobCrawl(depth int, jobDepth int, jobIndex string, jobMap *map[string]interface{}, mv interface{}) bool {
	if _, ok := mv.(map[string]interface{}); ok {
		for _, value := range mv.(map[string]interface{}) {
			XmlJobCrawl(depth+1, jobDepth, jobIndex, jobMap, value)
		}
	} else if _, ok := mv.([]interface{}); ok {
		for _, value := range mv.([]interface{}) {
			if depth == jobDepth-1 {
				if _, ok := value.(map[string]interface{}); ok {
					//jobUuid := uuid.New().String()
					valueMsi := value.(map[string]interface{})
					SanitizeMsi(&valueMsi)
					key := fmt.Sprintf("%v", valueMsi[jobIndex])
					//(*jobMap)[jobUuid] = valueMsi
					(*jobMap)[key] = valueMsi
				}
			}
			XmlJobCrawl(depth+1, jobDepth, jobIndex, jobMap, value)
		}
	} else if _, ok := mv.(string); ok {
		return true
	}

	return true
}

//
//
// SanitizeMsi
//
//
func SanitizeMsi(msi *map[string]interface{}) {
	for key, value := range *msi {
		if strings.ContainsAny(value.(string), "%") {
			value = strings.ReplaceAll(value.(string), "%", "&percnt;")
		}
		if strings.ContainsAny(value.(string), "\"") {
			value = strings.ReplaceAll(value.(string), "\"", "&quot;")
		}
		(*msi)[key] = value
	}
}
