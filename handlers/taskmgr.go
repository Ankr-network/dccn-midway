package handlers

import (
	"context"
	"encoding/json"
	"errors"
	//"fmt"
	"net/http"
	"time"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"github.com/Ankr-network/dccn-midway/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

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
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}
	case "2":
		task.Type = common_proto.TaskType_JOB
		task.TypeData = &common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image: Heretask.Image}}
	case "3":
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
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.TypeData = &common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: Heretask.Image}}
	case "2":
		task.Type = common_proto.TaskType_JOB
		task.TypeData = &common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image: Heretask.Image}}
	case "3":
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
		return
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
	jsonDcList, err := json.Marshal(rsp)
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
	jsonDcList, err := json.Marshal(rsp)
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