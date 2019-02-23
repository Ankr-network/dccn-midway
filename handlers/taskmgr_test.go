package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"
	"fmt"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)
type taskidstruct struct {
	task_id string `json:"task_id"`
}

func checkReplicastatus(t *testing.T, client *http.Client, target int32, taskid string, sessionTokenarray string) bool {
	t.Log("Checking the status of task ", taskid)
	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Add("Authorization", sessionTokenarray)
	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
		return false
	}
	defer respList.Body.Close()
	newbody := make([]*common_proto.Task, 0)
	bytebody, _ := ioutil.ReadAll(respList.Body)
	_ = json.Unmarshal(bytebody, &newbody)
	//t.Log(newbody)
	for i := range newbody {
		if newbody[i].Id == taskid {
			if newbody[i].Attributes.Replica == target {
				t.Log("Find the task and check the number of replica, it is the same!", taskid)
				return true
			}
			t.Log("Find the task and check the number of replica, it is not the same!", newbody[i].Attributes.Replica)
			return false
		}
	}
	t.Log("Did not find the taskid ", taskid)
	return false
}

func checkIDstatus(t *testing.T, client *http.Client, target common_proto.TaskStatus, taskid string, sessionTokenarray string) bool {
	t.Log("Checking the status of task ", taskid)
	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Add("Authorization", sessionTokenarray)
	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
		return false
	}
	defer respList.Body.Close()
	newbody := make([]*common_proto.Task, 0)
	bytebody, _ := ioutil.ReadAll(respList.Body)
	_ = json.Unmarshal(bytebody, &newbody)
	//t.Log(newbody)
	for i := range newbody {
		if newbody[i].Id == taskid {
			if newbody[i].Status == target {
				t.Log("Find the task and check the status of task, it is the same!", taskid)
				return true
			}
			t.Log("Find the task and check the status of task, it is not the same!", newbody[i].Status)
			return false
		}
	}
	t.Log("Did not find the taskid ", taskid)
	return false
}

func ClientLogin(t *testing.T, client *http.Client) string {

	var jsonStrlogin = []byte(`{"username":"testuser","email":"testuser28@mailinator.com", "password":"111111nn"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))
	reqlogin.Header.Set("X-Custom-Header", "myvalue")
	reqlogin.Header.Set("Content-Type", "application/json")

	resplogin, err := client.Do(reqlogin)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	sessionToken, _ := ioutil.ReadAll(resplogin.Body)
	t.Log("Sessiontoday:", string(sessionToken))
	sessionTokenarray := "bearer " + string(sessionToken)

	defer resplogin.Body.Close()

	t.Log("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	return sessionTokenarray
}

func TestCreateTask(t *testing.T) {

	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	t.Log("URL for Create Task:>", urlCreate)
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"Name": "updateuser",
    "Type": "1",
    "Image": "nginx:1.12",
	"Replica": 3,
	"DataCenterName": "datacenter01"
}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	var m map[string]interface{}
    err = json.Unmarshal([]byte(body), &m)
    if err != nil {
        t.Error("err")
    }
    mm := make(map[string]string)
    for k, v := range m {
        mm[k] = fmt.Sprint(v)
	}
	sbody := mm["task_id"]
	//body, _ := ioutil.ReadAll(respCreate.Body)
	//t.Log("Create Task Body:", string(body))
	//var Heretask taskidstruct
	//err = json.NewDecoder(respCreate.Body).Decode(&Heretask)
	//t.Log("Heretask:>", Heretask)
	//time.Sleep(time.Millisecond * 3000)

	if !checkIDstatus(t, client, common_proto.TaskStatus_RUNNING, sbody, sessionTokenarray) {
		t.Error("Tasks established faliure!")
	}
	if !checkReplicastatus(t, client, 3, sbody, sessionTokenarray) {
		t.Error("Tasks established faliure!")
	}
	//t.Log("Create Task Successfull!")
	pb := &Request{TaskID: sbody}
	jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err := http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Add("Authorization", sessionTokenarray)

	respPurge, err := client.Do(reqPurge)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respPurge.Body.Close()

	if respPurge.Status != "200 OK" {
		t.Error("Purge Status Error! Cannot login")
	}
}

func TestCreateTaskDouble(t *testing.T) {

	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	t.Log("URL for Create Task:>", urlCreate)
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"Name": "updateuser",
    "Type": "1",
    "Image": "nginx:1.12",
	"Replica": 3,
	"DataCenterName": "datacenter01"
}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	var m map[string]interface{}
    err = json.Unmarshal([]byte(body), &m)
    if err != nil {
        t.Error("err")
    }
    mm := make(map[string]string)
    for k, v := range m {
        mm[k] = fmt.Sprint(v)
	}
	sbody := mm["task_id"]

	time.Sleep(time.Millisecond * 3000)

	if !checkIDstatus(t, client, common_proto.TaskStatus_RUNNING, sbody, sessionTokenarray) {
		t.Error("Tasks established faliure!")
	}
	if !checkReplicastatus(t, client, 4, sbody, sessionTokenarray) {
		t.Error("Tasks established faliure!")
	}

	reqCreateDouble, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreateDouble.Header.Add("Authorization", sessionTokenarray)

	respCreateDouble, err := client.Do(reqCreateDouble)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreateDouble.Body.Close()

	t.Log("Create Task response Status:", respCreateDouble.Status)
	if respCreateDouble.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	time.Sleep(time.Millisecond * 3000)

	if !checkIDstatus(t, client, common_proto.TaskStatus_RUNNING, sbody, sessionTokenarray) {
		t.Error("Tasks has been overwrited")
	}
	if !checkReplicastatus(t, client, 4, sbody, sessionTokenarray) {
		t.Error("Tasks has been changed!")
	}

	//t.Log("Create Task Successfull!")
	pb := &Request{TaskID: sbody}
	jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err := http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Add("Authorization", sessionTokenarray)

	respPurge, err := client.Do(reqPurge)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respPurge.Body.Close()

	if respPurge.Status != "200 OK" {
		t.Error("Purge Status Error! Cannot login")
	}
}

func TestCreateTaskBADCred(t *testing.T) { //One cannot create task without login
	t.Log("URL for Create Task:>", urlCreate)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrCreate = []byte(`{"Name": "updateuser",
    "Type": "1",
    "Image": "nginx:1.12",
	"Replica": 3,
	"DataCenterName": "datacenter01"
}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Set("X-Custom-Header", "myvalue")
	reqCreate.Header.Set("Content-Type", "application/json")

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create response Status:", respCreate.Status)
	if respCreate.Status == "200 OK" {
		t.Error("Error! One should not create task without login!")
	}

}

func TestUpdateTask(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	t.Log("URL for Update Task:>", urlUpdate)

	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"Name": "updateuser",
    "Type": "1",
    "Image": "nginx:1.12",
	"Replica": 3,
	"DataCenterName": "datacenter01"
}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)
	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	var m map[string]interface{}
    err = json.Unmarshal([]byte(body), &m)
    if err != nil {
        t.Error("err")
    }
    mm := make(map[string]string)
    for k, v := range m {
        mm[k] = fmt.Sprint(v)
	}
	sbody := mm["task_id"]
	//sbody := string(body)
	checkIDstatus(t, client, common_proto.TaskStatus_RUNNING, sbody, sessionTokenarray)
	//payload := strings.NewReader("{\"id\":\"85e9d737-a8c1-4d02-86b4-81844f322a19\",\"user_id\":\"dba2456a-fe2a-4220-acbf-eafccf5b8af8\",\"name\":\"TestforPurgetask\",\"type\":\"web\",\"image\":\"nginx:1.12\",\"data_center\":\"Datacenter\",\"data_center_id\":\"10\",\"status\":7}")
	//time.Sleep(time.Millisecond*10000)
	jsonStrUpdate := Task{
		Name:         "xiaoWang",
		ID:            sbody,
		Image:        "nginx:1.12",
		Type:			"2",
		Replica:      1,
		DataCenterName:   "aslkdfjas",
	}
	jsonUpdateList, _ := json.Marshal(jsonStrUpdate)
	reqUpdate, err := http.NewRequest("POST", urlUpdate, bytes.NewBuffer(jsonUpdateList))
	reqUpdate.Header.Add("Authorization", sessionTokenarray)

	respUpdate, err := client.Do(reqUpdate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respUpdate.Body.Close()

	t.Log("Update response Status:", respUpdate.Status)
	if respUpdate.Status != "200 OK" {
		t.Error("Update Status Error! Cannot update the task")
	}
}

func TestUpdateTaskwithNotask(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	t.Log("URL for Update Task:>", urlUpdate)

	sessionTokenarray := ClientLogin(t, client)

	var jsonStrUpdate = []byte(`{"UserId": "123",
	"Name": "xiaowu",
	"Id": "120",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	}`)
	reqUpdate, err := http.NewRequest("POST", urlUpdate, bytes.NewBuffer(jsonStrUpdate))
	reqUpdate.Header.Add("Authorization", sessionTokenarray)

	respUpdate, err := client.Do(reqUpdate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respUpdate.Body.Close()

	t.Log("Update response Status:", respUpdate.Status)
	if respUpdate.Status == "200 OK" {
		t.Error("Update Status Error! One should not update the task when he do not has this task")
	}

	time.Sleep(500)

}

func TestUpdateTaskBADCred(t *testing.T) { //One cannot Update task without login
	t.Log("URL for Update Task:>", urlUpdate)
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

	t.Log("Signin response Status:", respUpdate.Status)
	if respUpdate.Status == "200 OK" {
		t.Error("Error! One should not update task without login!")
	}

}

func TestListTask(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	t.Log("URL for List Task:>", urlList)
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Add("Authorization", sessionTokenarray)

	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respList.Body.Close()

	t.Log("List response Status:", respList.Status)
	if respList.Status != "200 OK" {
		t.Error("Error! Cannot get the List of tasks!")
	}
	body, _ := ioutil.ReadAll(respList.Body)
	t.Log("List response Body:", string(body))
}

func TestListTaskBADCred(t *testing.T) { //One cannot create task without login
	t.Log("URL for List Task:>", urlList)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Set("X-Custom-Header", "myvalue")
	reqList.Header.Set("Content-Type", "application/json")

	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respList.Body.Close()

	t.Log("List response Status:", respList.Status)
	if respList.Status == "200 OK" {
		t.Error("Error! One should not list task without login!")
	}

}

func TestPurgeTask(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	t.Log("URL for Create Task:>", urlCreate)
	sessionTokenarray := ClientLogin(t, client)
	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "testforpurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1
	}`)
	
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskID: sbody}
	jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err := http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Add("Authorization", sessionTokenarray)

	respPurge, err := client.Do(reqPurge)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respPurge.Body.Close()

	if respPurge.Status != "200 OK" {
		t.Error("Purge Status Error! Cannot login")
	}

	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Add("Authorization", sessionTokenarray)
	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respList.Body.Close()

	newbody, _ := ioutil.ReadAll(respList.Body)
	t.Log("List Task Body:", string(newbody))
	if strings.Contains(string(newbody), sbody) {
		t.Error("Purge Status Error! Purge Faliure.")
	}
}

func TestPurgeTaskDoublePurge(t *testing.T) { //No man ever purges in the same tasks twice. --H.K1tahara
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("Create Status Error! Cannot login")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskID: sbody}
	jsonStrPurge, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqPurge, err := http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge.Header.Add("Authorization", sessionTokenarray)

	respPurge, err := client.Do(reqPurge)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respPurge.Body.Close()

	if respPurge.Status != "200 OK" {
		t.Error("Purge Status Error! Cannot login")
	}
	reqPurge1, err := http.NewRequest("POST", urlPurge, bytes.NewBuffer(jsonStrPurge))
	reqPurge1.Header.Add("Authorization", sessionTokenarray)
	respPurge1, err := client.Do(reqPurge1)
	if err != nil {
		t.Error("Purge Status Error! Purge Unsuccessful")
		//panic(err)
		//t.Error("err")
	}

	defer respPurge1.Body.Close()

	if respPurge1.Status == "200 OK" {
		t.Error("Purge Status Error! Purge Unsuccessful")
	}

}

func TestCancelTask(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("login Status Error!")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskID: sbody}
	jsonStrCancel, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqCancel, err := http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel.Header.Add("Authorization", sessionTokenarray)

	respCancel, err := client.Do(reqCancel)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCancel.Body.Close()

	if respCancel.Status != "200 OK" {
		t.Error("Cancel Status Error! Cannot login")
	}

	var jsonStrList = []byte(`{}`)
	reqList, err := http.NewRequest("POST", urlList, bytes.NewBuffer(jsonStrList))
	reqList.Header.Add("Authorization", sessionTokenarray)

	respList, err := client.Do(reqList)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respList.Body.Close()
	newbody := make([]*common_proto.Task, 0)
	bytebody, _ := ioutil.ReadAll(respList.Body)
	_ = json.Unmarshal(bytebody, &newbody)
	for i := range newbody {
		if newbody[i].Id == sbody {
			if !checkIDstatus(t, client, common_proto.TaskStatus_CANCELLED, sbody, sessionTokenarray) {
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
	sessionTokenarray := ClientLogin(t, client)

	var jsonStrCreate = []byte(`{"UserId": "123",
	"Name": "TestforPurgetask",
	"Id": "12",
    "Type": "web",
    "Image": "nginx:1.12",
	"Replica": 1,
	"DataCenter": "Datacenter",
	"DataCenterId": "10"}`)
	reqCreate, err := http.NewRequest("POST", urlCreate, bytes.NewBuffer(jsonStrCreate))
	reqCreate.Header.Add("Authorization", sessionTokenarray)

	respCreate, err := client.Do(reqCreate)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCreate.Body.Close()

	t.Log("Create Task response Status:", respCreate.Status)
	if respCreate.Status != "200 OK" {
		t.Error("Create Status Error! Cannot login")
	}
	body, _ := ioutil.ReadAll(respCreate.Body)
	t.Log("Create Task Body:", string(body))
	sbody := string(body)

	pb := &Request{TaskID: sbody}
	jsonStrCancel, err := json.Marshal(pb)
	if err != nil {
		t.Error("could not marshal JSON")
	}
	reqCancel, err := http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel.Header.Add("Authorization", sessionTokenarray)

	respCancel, err := client.Do(reqCancel)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respCancel.Body.Close()

	if respCancel.Status != "200 OK" {
		t.Error("Cancel Status Error! Cannot login")
	}

	reqCancel1, err := http.NewRequest("POST", urlCancel, bytes.NewBuffer(jsonStrCancel))
	reqCancel1.Header.Add("Authorization", sessionTokenarray)

	respCancel1, err := client.Do(reqCancel1)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	//defer respCancel1.Body.Close()

	if respCancel1.Status == "200 OK" {
		t.Error("Cancel Status Error! Cancel should not be successful")
	}

}
