package practise

import (
	"fmt"
	"os"
)

func ErrorTest() {
	file, err := open("/fake-dir/fake-file");
	if err != nil {
		fmt.Println("open file error!", err)
	}else{
		fmt.Println("open file success!", file)
	}
}

func open(path string) (file *os.File,err error){
	if file,err = os.Open(path);err != nil {
		fmt.Println("open file error!", err)
	}
	return
}

const limit = 0.000001
type ErrNegativeSqrt float64

func Sqrt(x float64) (float64, error) {
	if x == 0{
		return 0, nil
	}
	if x > 0 {
		z := 1.0
		for ;z*z - x > limit || z*z - x < -limit;z -= (z*z - x) / (2*z) {}
		return z, nil
	}
	return 0,ErrNegativeSqrt(x)
}

func (e ErrNegativeSqrt) Error() string {
	//return fmt.Sprint(e)//这样会死循环,因为Sprint内部会判断参数类型，如果实现error接口的话，会继续调用Error
	return fmt.Sprint("cannot Sqrt negative number:",float64(e))
}