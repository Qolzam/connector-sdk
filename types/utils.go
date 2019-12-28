package types

//  Adapted to Golang from https://github.com/eclipse/mosquitto/blob/master/lib/util_topic.c#L138
//  LICENSE : https://github.com/eclipse/mosquitto/blob/master/LICENSE.txt
//
func checkTopicMatch(sub string, topic string) bool {
	if sub == "" || topic == "" {
		return false
	}

	if sub == topic {
		return true
	}

	var topicLength = len(topic)
	var subLength = len(sub)
	var position = 0
	var subIndex = 0
	var topicIndex = 0

	if (sub[subIndex] == '$' && topic[topicIndex] != '$') || (topic[topicIndex] == '$' && sub[subIndex] != '$') {
		return true
	}

	for subIndex < subLength && topicIndex < topicLength {
		if topic[topicIndex] == '+' || topic[topicIndex] == '#' {
			return false
		}

		if sub[subIndex] != topic[topicIndex] || topicIndex >= topicLength {
			// Check for wildcard matches
			if sub[subIndex] == '+' {
				// Check for bad "+foo" or "a/+foo" subscription
				if position > 0 && sub[subIndex-1] != '/' {
					return false
				}

				// Check for bad "foo+" or "foo+/a" subscription
				if subIndex+1 < subLength && sub[subIndex+1] != '/' {
					return false
				}

				position++
				subIndex++
				for topicIndex < topicLength && topic[topicIndex] != '/' {
					topicIndex++
				}

				if topicIndex >= topicLength && subIndex >= subLength {
					return true
				}
			} else if sub[subIndex] == '#' {
				// Check for bad "foo#" subscription
				if position > 0 && sub[subIndex-1] != '/' {
					return false
				}

				// Check for # not the final character of the sub, e.g. "#foo"
				if subIndex+1 < subLength {
					return false
				} else {
					return true
				}
			} else {
				// Check for e.g. foo/bar matching foo/+/#
				if topicIndex >= topicLength && position > 0 && sub[subIndex-1] == '+' && sub[subIndex] == '/' && sub[subIndex+1] == '#' {
					return true
				}

				// There is no match at this point, but is the sub invalid?
				for subIndex < subLength {
					if sub[subIndex] == '#' && subIndex+1 < subLength {
						return false
					}

					position++
					subIndex++
				}

				// Valid input, but no match
				return false
			}
		} else {
			// sub[spos] == topic[tpos]
			if topicIndex+1 > topicLength {

				// Check for e.g. foo matching foo/#
				if sub[subIndex+1] == '/' && sub[subIndex+2] == '#' && subIndex+3 >= subLength {
					return true
				}
			}

			position++
			subIndex++
			topicIndex++

			if subIndex >= subLength && topicIndex >= topicLength {
				return true
			} else if topicIndex >= topicLength && sub[subIndex] == '+' && subIndex+1 >= subLength {
				if position > 0 && sub[subIndex-1] != '/' {
					return false
				}

				position++
				subIndex++

				return true
			}
		}
	}

	if topicIndex < topicLength || subIndex < subLength {
		return false
	}

	return true
}
