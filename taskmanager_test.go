
package main
 
import (
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"net/http/cookiejar"
	"encoding/json"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)
 
const (
	urllogin = "http://localhost:8000/login"
	urlCreate = "http://localhost:8000/create"
	urlUpdate = "http://localhost:8000/update"
	urlPurge = "http://localhost:8000/purge"
	urlCancel = "http://localhost:8000/delete"
	urllist = "http://localhost:8000/list"
)
 
/*type Person struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Age int `json:"age"`
}*/


func TestCreateTask(t *testing.T) {
	
	fmt.Println("URL for login:>", urllogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	
	fmt.Println("URL for Create Task:>", urlCreate)
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforCreatetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK"{
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	fmt.Println("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskId: sbody}
    jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err :=http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Set("X-Custom-Header", "myvalue")
    reqPurge.Header.Set("Content-Type", "application/json")

	respPurge, err := client.Do(reqPurge)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respPurge.Body.Close()
	
	if respPurge.Status != "200 OK"{
		t.Error("Purge Status Error! Cannot login")
	}
}

func TestCreateTaskBADCred(t *testing.T) { //One cannot create task without login
	fmt.Println("URL for Create Task:>", urlCreate)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "xiaowu",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter01",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Signin response Status:", respCreate.Status)
	if respCreate.Status == "200 OK"{
		t.Error("Error! One should not create task without login!")
	}

	
}


func TestUpdateTask(t *testing.T) {
	fmt.Println("URL for login:>", urllogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	fmt.Println("URL for Update Task:>", urlUpdate)
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrUpdate = []byte(`{"UserId": "123",
	"Name": "xiaowu",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "aslkdfjas",
	"DataCenterId": "10"}`)
	reqUpdate, err := http.NewRequest("POST", urlUpdate, bytes.NewBuffer(jsonStrUpdate))
    reqUpdate.Header.Set("X-Custom-Header", "myvalue")
    reqUpdate.Header.Set("Content-Type", "application/json")
	
	respUpdate, err := client.Do(reqUpdate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respUpdate.Body.Close()

	fmt.Println("Update response Status:", respUpdate.Status)
	if respUpdate.Status != "200 OK"{
		t.Error("Update Status Error! Cannot update the task")
	}

}

func TestUpdateTaskBADCred(t *testing.T) { //One cannot Update task without login
	fmt.Println("URL for Update Task:>", urlUpdate)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrUpdate = []byte(`{"UserId": "123",
	"Name": "xiaowu",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "aslkdfjas",
	"DataCenterId": "10"}`)
	reqUpdate, err := http.NewRequest("POST", urlUpdate, bytes.NewBuffer(jsonStrUpdate))
    reqUpdate.Header.Set("X-Custom-Header", "myvalue")
    reqUpdate.Header.Set("Content-Type", "application/json")
	
	respUpdate, err := client.Do(reqUpdate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respUpdate.Body.Close()

	fmt.Println("Signin response Status:", respUpdate.Status)
	if respUpdate.Status == "200 OK"{
		t.Error("Error! One should not update task without login!")
	}

	
}




func TestListTask(t *testing.T) {
	fmt.Println("URL for login:>", urllogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	
	fmt.Println("URL for List Task:>", urllist)
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urllist, bytes.NewBuffer(jsonStrList))
    reqList.Header.Set("X-Custom-Header", "myvalue")
    reqList.Header.Set("Content-Type", "application/json")
	
	respList, err := client.Do(reqList)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respList.Body.Close()

	fmt.Println("List response Status:", respList.Status)
	if respList.Status != "200 OK"{
		t.Error("Error! Cannot get the List of tasks!")
	}
	body, _ := ioutil.ReadAll(respList.Body)
    fmt.Println("List response Body:", string(body))
}

func TestListTaskBADCred(t *testing.T) { //One cannot create task without login
	fmt.Println("URL for List Task:>", urllist)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urllist, bytes.NewBuffer(jsonStrList))
    reqList.Header.Set("X-Custom-Header", "myvalue")
    reqList.Header.Set("Content-Type", "application/json")
	
	respList, err := client.Do(reqList)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respList.Body.Close()

	fmt.Println("List response Status:", respList.Status)
	if respList.Status == "200 OK"{
		t.Error("Error! One should not list task without login!")
	}

	
}

func TestPurgeTask(t *testing.T) {
	fmt.Println("URL for login:>", urllogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	
	fmt.Println("URL for Create Task:>", urlCreate)
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK"{
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	fmt.Println("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskId: sbody}
    jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err :=http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Set("X-Custom-Header", "myvalue")
    reqPurge.Header.Set("Content-Type", "application/json")

	respPurge, err := client.Do(reqPurge)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respPurge.Body.Close()
	
	if respPurge.Status != "200 OK"{
		t.Error("Purge Status Error! Cannot login")
	}

	var jsonStrList = []byte(`{}`)
	reqList, err :=http.NewRequest("POST", urllist, bytes.NewBuffer(jsonStrList))
	reqList.Header.Set("X-Custom-Header", "myvalue")
	reqList.Header.Set("Content-Type", "application/json")
	respList, err := client.Do(reqList)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respList.Body.Close()

	newbody, _ := ioutil.ReadAll(respList.Body)
	if bytes.ContainsAny(newbody, sbody) {
		t.Error("Purge Status Error! Purge Faliure.")
	}
}

func TestPurgeTaskDoublePurge(t *testing.T) { //No man ever purges in the same tasks twice. --H.K1tahara
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK"{
		t.Error("Create Status Error! Cannot login")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	fmt.Println("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskId: sbody}
    jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err :=http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Set("X-Custom-Header", "myvalue")
    reqPurge.Header.Set("Content-Type", "application/json")

	respPurge, err := client.Do(reqPurge)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respPurge.Body.Close()
	
	if respPurge.Status != "200 OK"{
		t.Error("Purge Status Error! Cannot login")
	}
	respPurge, err = client.Do(reqPurge)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respPurge.Body.Close()
	
	if respPurge.Status == "200 OK"{
		t.Error("Purge Status Error! Purge Unsuccessful")
	}

}



func TestCancelTask(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK"{
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	fmt.Println("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskId: sbody}
    jsonStrCancel, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqCancel, err :=http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel.Header.Set("X-Custom-Header", "myvalue")
    reqCancel.Header.Set("Content-Type", "application/json")

	respCancel, err := client.Do(reqCancel)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respCancel.Body.Close()
	
	if respCancel.Status != "200 OK"{
		t.Error("Cancel Status Error! Cannot login")
	}

	var jsonStrList = []byte(`{}`)
	reqList, err :=http.NewRequest("POST", urllist, bytes.NewBuffer(jsonStrList))
	reqList.Header.Set("X-Custom-Header", "myvalue")
	reqList.Header.Set("Content-Type", "application/json")
	respList, err := client.Do(reqList)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respList.Body.Close()
	newbody := make([]*common_proto.Task, 0)
	bytebody, _ := ioutil.ReadAll(respList.Body)
	_ = json.Unmarshal(bytebody, &newbody)
	for i := range newbody{
		if newbody[i].Id == sbody{
			if newbody[i].Status != common_proto.TaskStatus_CANCEL {
				t.Error("Error! Fail to cancel the task")
			}
		}
	}

}

func TestCancelTaskDoubleCancel(t *testing.T) { //No man ever cancels in the same tasks twice, too. --H.K1tahara
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	
	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
    reqCreate.Header.Set("X-Custom-Header", "myvalue")
    reqCreate.Header.Set("Content-Type", "application/json")
	
	respCreate, err := client.Do(reqCreate)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer respCreate.Body.Close()

	fmt.Println("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK"{
		t.Error("Create Status Error! Cannot login")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	fmt.Println("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskId: sbody}
    jsonStrCancel, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqCancel, err :=http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel.Header.Set("X-Custom-Header", "myvalue")
    reqCancel.Header.Set("Content-Type", "application/json")

	respCancel, err := client.Do(reqCancel)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respCancel.Body.Close()
	
	if respCancel.Status != "200 OK"{
		t.Error("Cancel Status Error! Cannot login")
	}

	reqCancel1, err :=http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel1.Header.Set("X-Custom-Header", "myvalue")
    reqCancel1.Header.Set("Content-Type", "application/json")

	respCancel1, err := client.Do(reqCancel1)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
	defer respCancel1.Body.Close()
	
	if respCancel1.Status == "200 OK"{
		t.Error("Cancel Status Error! Cancel Unsuccessful")
	}

}
