package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"errors"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

type Task struct {
	UserId       string `json:"UserId"`
	Name         string `json:"Name"`
	Id           string `json:"ID"`
	Type         string`json:"Type"`
	Image        string `json:"Image"`
	Replica      int32  `json:"Replica"`
	DataCenter   string `json:"DataCenter"`
	DataCenterId string `json:"DataCenterId"`
}

type Request struct {
	TaskId string `json:"TaskId"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Create Tasks")
	c, err := sessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c

	Heretask := common_proto.Task{}
	// Get the JSON body and decode into credentials
	err = json.NewDecoder(r.Body).Decode(&Heretask)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Info("Something went wrong! ", err)
		return
	}

	log.Info(Heretask)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tcrq := taskmgr.CreateTaskRequest{
		UserId: "sessionUserid",
		Task:   &Heretask,
	}
	tcrp, err := dc.CreateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if tcrp.Error == nil {
		log.Info("Task created successfully. \n", tcrp.TaskId)
		w.Write([]byte(fmt.Sprintf("%s", tcrp.TaskId)))
	} else {
		log.Info("Fail to create task. \n", tcrp.Error)
	}

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Update Tasks")
	sessionToken, err:= sessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

	Heretask := common_proto.Task{}
	// Get the JSON body and decode into credentials
	err1 := json.NewDecoder(r.Body).Decode(&Heretask)
	if err1 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Something went wrong! %s\n", err1)
		return
	}

	log.Info(Heretask)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	/*task := common_proto.Task{
		UserId:       "sessionUserid",
		Name:         Heretask.Name,
		Type:         Heretask.Type,
		Image:        Heretask.Image,
		DataCenter:   Heretask.DataCenter,
		DataCenterId: Heretask.DataCenterId,
	}*/
	tcrq := taskmgr.UpdateTaskRequest{
		UserId: "sessionUserid",
		Task:   &Heretask,
	}
	_, err = dc.UpdateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info("Task updated successfully. \n")
}

func ListTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Lists")
	sessionToken, err := sessionTokenValue(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

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

	userTasks := make([]*common_proto.Task, 0)
	rsp, err := dc.TaskList(tokenContext, &taskmgr.ID{UserId: "sessionUserid"});
	if  err != nil {
		log.Info(err)
		w.WriteHeader(http.StatusBadRequest)
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
	c, err := sessionTokenValue(w, r)
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
	err1 := json.NewDecoder(r.Body).Decode(&NewRequest)
	if err1 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Something went wrong! %s\n", err1)
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

	if _, err := dc.CancelTask(tokenContext, &taskmgr.Request{UserId: "sessionUserid", TaskId: NewRequest.TaskId}); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Println("CancelTask Ok")
	}
}

func PurgeTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Purge Task")
	c, err := sessionTokenValue(w, r)
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
	err1 := json.NewDecoder(r.Body).Decode(&NewRequest)
	log.Info(NewRequest)
	log.Info("Above it the new Request")
	if err1 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Something went wrong! %s\n", err1)
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

	if _, err := dc.PurgeTask(tokenContext, &taskmgr.Request{UserId: "", TaskId: NewRequest.TaskId}); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Println("PurgeTask Ok")
	}
}

func TaskDetail(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Detail")
	c, err := sessionTokenValue(w, r)
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

	if tcrp, err := dc.TaskDetail(tokenContext, &taskmgr.Request{UserId: "", TaskId: NewRequest.TaskId}); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Println("TaskDetail Ok")
		w.Write([]byte(fmt.Sprintf("%s", tcrp.Task)))
	}
}

func DataCenterList(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Datacenter Lists")
	c, err := sessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
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
	rsp, err := dc.DataCenterList(tokenContext, &dcmgr.DataCenterListRequest{UserId: ""});
	if err != nil {
		log.Info(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
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