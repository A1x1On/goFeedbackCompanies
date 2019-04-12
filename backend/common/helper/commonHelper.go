package helper

// analogue of ternary operator for get
func GetIf(statement bool, a, b interface{}) interface{} {
	if statement {
		 return a
	}
	return b
}

// analogue of ternary operator for func exec 
func SetIf(statement bool, fnA, fnB func()){
	if statement {
		fnA()
  } else {
	fnB()
  }
}



