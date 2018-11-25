package holiday

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DayInfo struct {
	Code int `json:code`
	Data int `json:data`
	//Day  string
}

type HolidayClient struct {
	Window []TimeWindow
}

type TimeWindow2 struct {
	Begin string `yaml:"begin"`
	End   string `yaml:"end"`
}

type TimeWindow struct {
	Begin time.Time
	End   time.Time
}

var Holidays = make(map[string]DayInfo)

const BaseUrl = "http://api.goseek.cn/Tools/holiday?date="
const layout = "15:04:05"

func NewHolidayClient() *HolidayClient {
	return &HolidayClient{}
}

func NewHolidayClientWithTimeWindow(t []TimeWindow2) (*HolidayClient, error) {

	var client *HolidayClient = &HolidayClient{}
	//var err error

	for _, v := range t {
		begin, err := time.Parse(layout, v.Begin) //定义一个小时:分钟维度的time
		if err != nil {
			return nil, fmt.Errorf("begin time format error, must be HH:MM:SS ")
		}
		end, err := time.Parse(layout, v.End)
		if err != nil {
			return nil, fmt.Errorf("end time format error, must be HH:MM:SS ")
		}

		if begin.After(end) {
			return nil, fmt.Errorf("beginTime should not > endTime ")
		}

		fmt.Println("begin,end: ", begin, end)
		//var a  TimeWindow
		client.Window = append(client.Window, TimeWindow{Begin: begin, End: end})
	}
	return client, nil
}

//节假日前一天上线窗口判断，通过限制的窗口，去查询是否在窗口中，如在明天是节假日，且时间在上线窗口中，则返回true，其他情况返回false。
func (hc *HolidayClient) BeforeHolidayWindow(nowTime time.Time) (bool, error) {
	//判断今天是否节假日，如果是则结束
	if res, err := HolidayInfo(nowTime); !res && err == nil {
		return false, fmt.Errorf("today is holiday ，permission denied")
	}
	//const layout = "15:04:05"
	//判断当前时间是否在指定时间窗中.如果是，继续判断是否为节假日前一天

	//把时间戳转为layout格式string后，再转回layout的time，只保留 HH:MM 进行对比
	bb, _ := time.Parse(layout, nowTime.Format(layout))

	for i := range hc.Window {
		fmt.Println("bb,end,begin:", bb, hc.Window[i].End, hc.Window[i].Begin)
		if (bb.Before(hc.Window[i].End) && bb.After(hc.Window[i].Begin)) || bb.Equal(hc.Window[i].End) || bb.Equal(hc.Window[i].Begin) {
			fmt.Println("nowtime is in window ")
			tomorrow := nowTime.AddDate(0, 0, 1)

			//判断明天是否节假日，如果是，则不允许上线！
			k, err1 := HolidayInfo(tomorrow)
			if !k && err1 == nil {
				fmt.Println(" nowtime is in the online window, but tomorrow is holiday,permission denied !")
				return true, nil
			}

		} else {
			fmt.Println(bb, "nowtime is not in window ")

		}
	}
	return false, nil

}

// 参数为Time格式， 工作日返回true，节假日返回false！
func HolidayInfo(d time.Time) (bool, error) {
	//格式化为指定类型
	day := d.Format("20060102")
	key, err1 := Holidays[day]
	if err1 && key.Data == 0 {
		return true, nil
	}

	url := BaseUrl + day
	response, err2 := http.Get(url)
	if err2 != nil {
		return false, err2
	}
	defer response.Body.Close()

	var body []byte
	body, _ = ioutil.ReadAll(response.Body)
	info := &DayInfo{}
	err := json.Unmarshal(body, info)
	if err != nil {
		fmt.Println("json unmarshal failed! ")
		return false, err
	}

	Holidays[day] = *info

	if info.Data == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// 参数为Time格式， 工作日返回true，节假日返回false！
func NowDayInfo() (bool, error) {
	return HolidayInfo(time.Now())
}
