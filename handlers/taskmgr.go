package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"

	"github.com/shopspring/decimal"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"github.com/Ankr-network/dccn-midway/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

const (
	BittrexADDRESS = "https://api.bittrex.com"
)
var (
	urlUSDT  = fmt.Sprintf("%s/api/v1.1/public/getticker?market=USDT-BTC", BittrexADDRESS)
	urlBTC = fmt.Sprintf("%s/api/v1.1/public/getticker?market=BTC-ANKR", BittrexADDRESS)
)

type Ticker struct {
	Bid  decimal.Decimal `json:"Bid"`
	Ask  decimal.Decimal `json:"Ask"`
	Last decimal.Decimal `json:"Last"`
}

type Bitraxbody struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Result Ticker `json:"result"`
}


type Task struct {
	Name         string `json:"Name"`
	ID           string `json:"ID"`
	Type         string `json:"Type"`
	Image        string `json:"Image"`
	Replica      int32  `json:"Replica"`
	hidden 		 string   `json:"hidden"`
	CreationDate uint64 `json:"CreationDate"`
	DataCenterName   string `json:"DataCenterName"`
	LastModifiedDate uint64 `json:"LastModifiedDate"`
	Schedule string `json:"Schedule"`
}

type Request struct {
	TaskID string `json:"TaskID"`
}

type TaskOverviewSendback struct {
	ClusterCount int32 `json:"cluster_count"`
	EnvironmentCount int32 `json:"environment_count"`
	RegionCount int32 `json:"region_count"`
	TotalTaskCount int32 `json:"total_task_count"`
	HealthTaskCount int32 `json:"health_task_count"`
}

type NetworkInfoSendback struct {
	UserCount int32 `json:"user_count"`
	HostCount int32 `json:"host_count"`
	EnvironmentCount int32 `json:"environment_count"`
	ContainerCount int32 `json:"container_count"`
	Traffic string`json:"traffic"`
}

type USDTANKR struct {
	Price decimal.Decimal `json:"price"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Create Tasks")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c

	var Heretask Task
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&Heretask)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Info("Something went wrong! in Createtask", err)
		return
	}
	log.Info(Heretask)

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})
	if err != nil {
		log.Info("Cannot convert Hidden into Bool!")
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return 
	}
	attribute := common_proto.TaskAttributes{
		Replica:      int32(Heretask.Replica),
	}

	task := common_proto.Task{
		Id:       	  Heretask.ID,
		Name:         Heretask.Name,
		Attributes:   &attribute,
		DataCenterName: Heretask.DataCenterName,
	}
	switch Heretask.Type {
	case "0":
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}
	case "1":
		task.Type = common_proto.TaskType_JOB
		task.TypeData = &common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image: Heretask.Image}}
	case "2":
		task.Type = common_proto.TaskType_CRONJOB
		task.TypeData = &common_proto.Task_TypeCronJob{TypeCronJob: &common_proto.TaskTypeCronJob{Image: Heretask.Image, Schedule: Heretask.Schedule}}
	default:
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}

	}
	/*switch u := task.Type.(type) {
		case *common_proto.Task_TypeDeployment: // u.Number contains the number.
		case *common_proto.Task_TypeDeployment: // u.Name contains the string.
		}*/
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tcrq := taskmgr.CreateTaskRequest{
		Task:   &task,
	}
	tcrp, err := dc.CreateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	JsonUser, err := json.Marshal(tcrp)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong in Marshall Request! %s\n", err)
		return
	}
	w.Write(JsonUser)
	log.Info("Task created successfully. \n", tcrp.TaskId)
	//w.Write([]byte(fmt.Sprintf("%s", tcrp.TaskId)))

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Update Tasks")
	sessionToken, err := util.SessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

	var Heretask Task
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&Heretask)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	log.Info(Heretask)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	if err != nil {
		log.Info("Cannot convert Hidden into Bool!")
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return 
	}
	attribute := common_proto.TaskAttributes{
		Replica:      int32(Heretask.Replica),
	}

	task := common_proto.Task{
		Id:       	  Heretask.ID,
		Name:         Heretask.Name,
		Attributes:   &attribute,
		DataCenterName: Heretask.DataCenterName,
	}
	switch Heretask.Type {
	case "0":
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}
	case "1":
		task.Type = common_proto.TaskType_JOB
		task.TypeData = &common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image: Heretask.Image}}
	case "2":
		task.Type = common_proto.TaskType_CRONJOB
		task.TypeData = &common_proto.Task_TypeCronJob{TypeCronJob: &common_proto.TaskTypeCronJob{Image: Heretask.Image, Schedule: Heretask.Schedule}}
	default:
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}
	}
	tcrq := taskmgr.UpdateTaskRequest{
		Task:   &task,
	}
	_, err = dc.UpdateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	log.Info("Task updated successfully. \n")
	w.Write([]byte("Task updated successfully."))
}

func ListTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Lists")
	sessionToken, err := util.SessionTokenValue(w, r)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)
	var NewRequest Request
	//if r.Body == {} {
	//	NewRequest = Request{}
	//} else {
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&NewRequest)
	//}
	log.Info("xiaouw")
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, "The body is empty! Please request with a vaild object", http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
		//NewRequest = Request{}
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	userTasks := make([]*common_proto.Task, 0)
	rsp, err := dc.TaskList(tokenContext, &taskmgr.TaskListRequest{TaskFilter: &taskmgr.TaskFilter{TaskId: NewRequest.TaskID}})
	if err != nil {
		log.Info(err)
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	userTasks = append(userTasks, rsp.Tasks...)
	if len(userTasks) == 0 {
		log.Printf("no tasks belongs to You!")
	//	return
	}
	log.Println(len(userTasks), "tasks belongs to You!")
	jsonTaskList, _ := json.Marshal(userTasks)
	w.Write(jsonTaskList)
	for i := range userTasks {
		log.Println(userTasks[i])
	}
}

func CancelTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Cancel Task")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c

	var NewRequest Request
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&NewRequest)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if _, err := dc.CancelTask(tokenContext, 
		&taskmgr.TaskID{TaskId: NewRequest.TaskID}); err != nil {
		log.Println(err.Error())
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	log.Println("CancelTask Ok")
	w.Write([]byte("CancelTask Ok"))
}

func PurgeTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Purge Task")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	var NewRequest Request
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&NewRequest)
	log.Info(NewRequest)
	log.Info("Above it the new Request")
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if _, err := dc.PurgeTask(tokenContext, 
		&taskmgr.TaskID{TaskId: NewRequest.TaskID}); err != nil {
		log.Println(err.Error())
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	log.Println("PurgeTask Ok")
	w.Write([]byte("PurgeTask Ok"))
}
/*
func TaskDetail(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Detail")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c
	var NewRequest Request
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&NewRequest)
	log.Info(NewRequest)
	log.Info("Above it the new Request")
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Something went wrong! %s\n", err)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if tcrp, err := dc.TaskDetail(tokenContext, 
		&taskmgr.Request{UserId: "", TaskId: NewRequest.TaskID}); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Println("TaskDetail Ok")
		w.Write([]byte(fmt.Sprintf("%s", tcrp.Task)))
	}
}
*/
func DataCenterList(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Datacenter Lists")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
	}

	defer conn.Close()
	dc := dcmgr.NewDCAPIClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	Datacenters := make([]*common_proto.DataCenter, 0)
	rsp, err := dc.DataCenterList(tokenContext, &common_proto.Empty{})
	if err != nil {
		log.Info(err.Error())
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	Datacenters = append(Datacenters, rsp.DcList...)
	if len(Datacenters) == 0 {
		log.Printf("no datacenter is running now")
		return
	}
	log.Println(len(Datacenters), "datacenters is running now")
	jsonDcList, _ := json.Marshal(Datacenters)
	w.Write(jsonDcList)
	for i := range Datacenters {
		log.Println(Datacenters[i])
	}
}


func TaskOverview(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Overview")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	rsp, err := dc.TaskOverview(tokenContext, &common_proto.Empty{})
	if err != nil {
		log.Info(err.Error())
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	Sendback := TaskOverviewSendback{
		ClusterCount:	rsp.ClusterCount,
		EnvironmentCount: rsp.EnvironmentCount,
		RegionCount:	rsp.RegionCount,
		TotalTaskCount: rsp.TotalTaskCount,
		HealthTaskCount: rsp.HealthTaskCount,
	}
	//int32 environment_count = 2;
	//int32 region_count = 3;
	//int32 total_task_count = 4;
	//int32 health_task_count = 5;
	jsonDcList, err := json.Marshal(Sendback)
	if err != nil {
		log.Info("Marshal Error", err)
		http.Error(w, util.ParseError(err), http.StatusNotFound)
	}
	w.Write(jsonDcList)
}

func NetworkInfo(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Network Information")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
	}

	defer conn.Close()
	dc := dcmgr.NewDCAPIClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	rsp, err := dc.NetworkInfo(tokenContext, &common_proto.Empty{})
	if err != nil {
		log.Info(err.Error())
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	Sendback := NetworkInfoSendback{
		UserCount:	rsp.UserCount,
		HostCount:	rsp.HostCount,
		EnvironmentCount: rsp.EnvironmentCount,
		ContainerCount: rsp.ContainerCount,
		
	}
	switch rsp.Traffic {
	case 0:
		Sendback.Traffic = "N/A"
	case 1:
		Sendback.Traffic = "LIGHT"
	case 2:
		Sendback.Traffic = "MEDIUM"
	case 3:
		Sendback.Traffic = "HEAVY"
	default:
		Sendback.Traffic = "N/A"
	}
	jsonDcList, err := json.Marshal(Sendback)
	if err != nil {
		log.Info("Marshal Error", err)
		http.Error(w, util.ParseError(err), http.StatusNotFound)
	}
	w.Write(jsonDcList)
}

func TaskLeaderBoard(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task LeaderBoard")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	LeaderBoard := make([]*taskmgr.TaskLeaderBoardDetail, 0)
	rsp, err := dc.TaskLeaderBoard(tokenContext, &common_proto.Empty{})
	if err != nil {
		log.Info(err.Error())
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	LeaderBoard = append(LeaderBoard, rsp.List...)
	if len(LeaderBoard) == 0 {
		log.Printf("no LeaderBoard is running now")
		return
	}
	log.Println(len(LeaderBoard), "leaderboard is running now")
	jsonLeaderBoard, _ := json.Marshal(LeaderBoard)
	w.Write(jsonLeaderBoard)
	for i := range LeaderBoard {
		log.Println(LeaderBoard[i])
	}
}


func DCLeaderBoard(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("DataCenter LeaderBoard")
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
	}

	defer conn.Close()
	dc := dcmgr.NewDCAPIClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	LeaderBoard := make([]*dcmgr.DataCenterLeaderBoardDetail, 0)
	rsp, err := dc.DataCenterLeaderBoard(tokenContext, &common_proto.Empty{})
	if err != nil {
		log.Info(err.Error())
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	LeaderBoard = append(LeaderBoard, rsp.List...)
	if len(LeaderBoard) == 0 {
		log.Printf("no LeaderBoard is running now")
		return
	}
	log.Println(len(LeaderBoard), "leaderboard is running now")
	jsonLeaderBoard, _ := json.Marshal(LeaderBoard)
	w.Write(jsonLeaderBoard)
	for i := range LeaderBoard {
		log.Println(LeaderBoard[i])
	}
}

func AnkrPrice(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("AnkrPrice")
	var jsonStrList = []byte(`{}`)
	reqUSDT, err := http.NewRequest("GET", urlUSDT, bytes.NewBuffer(jsonStrList))
	log.Info(reqUSDT)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	reqBTC, err := http.NewRequest("GET", urlBTC, bytes.NewBuffer(jsonStrList))
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	num := 0
	client := &http.Client{}
	respUSDT, err := client.Do(reqUSDT)
	log.Info(string(respUSDT))
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	//defer respUSDT.Body.Close()
	btcusdt, _ := ioutil.ReadAll(respUSDT.Body)
	
	var usdtbody Bitraxbody
	err = json.Unmarshal(btcusdt, &usdtbody)
	usdt := usdtbody.Result.Last
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	if usdtbody.Success == false {
		num = 0
		for usdtbody.Success == false {
			log.Info("xiaowag")
			respUSDT, err = client.Do(reqUSDT)
			btcusdt, _ = ioutil.ReadAll(respUSDT.Body)
			err = json.Unmarshal(btcusdt, &usdtbody)
			num = num + 1
			if num > 30 {
				http.Error(w, "There is something wrong with the Bitrex API, please try again", http.StatusBadRequest)
				return
			}
		}	
	}
	respBTC, err := client.Do(reqBTC)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	defer respBTC.Body.Close()
	btcankr, _ := ioutil.ReadAll(respBTC.Body)
	var btcbody Bitraxbody
	_ = json.Unmarshal(btcankr, &btcbody)
	btc := btcbody.Result.Last

	if btcbody.Success == false {
		num = 0
		for btcbody.Success == false {
			log.Info("xiaowu")
			respBTC, err = client.Do(reqBTC)
			btcankr, _ = ioutil.ReadAll(respBTC.Body)
			err = json.Unmarshal(btcankr, &btcbody)
			num = num + 1
			if num > 30 {
				http.Error(w, "There is something wrong with the Bitrex API, please try again", http.StatusBadRequest)
				return
			}
		}
		
	}
	log.Info(btc)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
	}
	jsonPrice := USDTANKR{
		Price: btc.Mul(usdt),
	}
	OutputPrice, err := json.Marshal(jsonPrice)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong in Marshall Request! %s\n", err)
		return
	}
	w.Write(OutputPrice)
}
