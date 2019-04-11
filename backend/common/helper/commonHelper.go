package helper

// analogue of ternary operator for get
func IfGet(statement bool, a, b interface{}) interface{} {
	if statement {
		 return a
	}
	return b
}

// analogue of ternary operator for func exec 
func IfSet(statement bool, fnA, fnB interface{}){
	if statement {
		fnA.(func())()
  } else {
	fnB.(func())()
  }
}



