package holiday

import (
	"fmt"
	"testing"
	"time"
)

// 参数为Time格式， 工作日返回true，节假日返回false！
func TestHolidayInfo(t *testing.T) {
	//fmt.Println(HolidayInfo(time.Now()))
	if res, err := HolidayInfo(time.Now()); err == nil && res {
		fmt.Println("今天是工作日", res)
	}

	if res, err := HolidayInfo(time.Now().AddDate(0, 0, 1)); err == nil && res {
		fmt.Println("明天是工作日", res)
	} else {
		fmt.Println("明天是节假日", res)

	}

}

func TestBeforeHolidayWindow(t *testing.T) {
	fmt.Println("TestBeforeHolidayWindow...")
	begin := "11:00:00"
	end := "13:00:00"
	sl := []TimeWindow2{}
	sl = append(sl, TimeWindow2{Begin: begin, End: end}, TimeWindow2{Begin: "14:00:00", End: "16:00:00"})
	window, _ := NewHolidayClientWithTimeWindow(sl)

	onLine := time.Now()
	res, err := window.BeforeHolidayWindow(onLine)
	if err != nil {
		fmt.Println("TestBeforeHolidayWindow ok !")
	}
	if res {
		fmt.Println("命中节假日前一天限制窗口！返回值为：", res)
	} else {
		fmt.Println("未命中节假日前一天限制窗口！返回值为：", res)
	}
}

// 参数为Time格式， 工作日返回true，节假日返回false！
func TestNowDayInfo(t *testing.T) {
	res, err := NowDayInfo()
	if err == nil {
		fmt.Println("NowDayInfo test...")
	}
	if res {
		fmt.Println("工作日", res)
	} else {
		fmt.Println("节假日", res)
	}

}
