package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

type Task struct {
	UserId       string `json:"UserId"`
	Name         string `json:"Name"`
	Id           string `json:"ID"`
	Type         string `json:"Type"`
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
	sessionToken, sessionUserid, err := getSessionValues(w, r)
	if err != nil {
		//The error has been sent in func getSeessionValues
		return
	}
	log.Info(sessionToken)
	log.Info(sessionUserid)

	var Heretask Task
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
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	task := common_proto.Task{
		UserId:       sessionUserid,
		Name:         Heretask.Name,
		Type:         Heretask.Type,
		Image:        Heretask.Image,
		DataCenter:   Heretask.DataCenter,
		DataCenterId: Heretask.DataCenterId,
	}
	tcrq := taskmgr.CreateTaskRequest{
		UserId: sessionUserid,
		Task:   &task,
	}
	tcrp, err := dc.CreateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
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
	sessionToken, sessionUserid, err := getSessionValues(w, r)
	if err != nil {
		return
	}
	log.Info(sessionToken)
	log.Info(sessionUserid)

	var Heretask Task
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
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	task := common_proto.Task{
		UserId:       sessionUserid,
		Name:         Heretask.Name,
		Type:         Heretask.Type,
		Image:        Heretask.Image,
		DataCenter:   Heretask.DataCenter,
		DataCenterId: Heretask.DataCenterId,
	}
	tcrq := taskmgr.UpdateTaskRequest{
		UserId: sessionUserid,
		Task:   &task,
	}
	Err2, err := dc.UpdateTask(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		return
	}
	if Err2 == nil {
		log.Info("Task created successfully. \n")
	} else {
		log.Info("Fail to create task. \n")
	}

}

func ListTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Task Lists")
	sessionToken, sessionUserid, err := getSessionValues(w, r)
	if err != nil {
		return
	}
	log.Info(sessionToken)
	log.Info(sessionUserid)

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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
	if rsp, err := dc.TaskList(tokenContext, &taskmgr.ID{UserId: sessionUserid}); err != nil {
		log.Fatal(err.Error())
	} else {
		userTasks = append(userTasks, rsp.Tasks...)
		if len(userTasks) == 0 {
			log.Printf("no tasks belongs to %s", sessionUserid)
		} else {
			log.Println(len(userTasks), "tasks belongs to ", sessionUserid)
			jsonTaskList, _ := json.Marshal(userTasks)
			w.Write(jsonTaskList)
			for i := range userTasks {
				log.Println(userTasks[i])
				//w.Write([]byte(fmt.Sprintf("%s", userTasks[i])))
			}

		}
	}

}

func CancelTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Cancel Task")
	sessionToken, sessionUserid, err := getSessionValues(w, r)
	if err != nil {
		return
	}
	log.Info(sessionToken)
	log.Info(sessionUserid)

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
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if _, err := dc.CancelTask(tokenContext, &taskmgr.Request{UserId: sessionUserid, TaskId: NewRequest.TaskId}); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("CancelTask Ok")
	}
}

func PurgeTask(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Purge Task")
	sessionToken, sessionUserid, err := getSessionValues(w, r)
	if err != nil {
		return
	}
	log.Info(sessionToken)
	log.Info(sessionUserid)

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
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	dc := taskmgr.NewTaskMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	tokenContext, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if _, err := dc.PurgeTask(tokenContext, &taskmgr.Request{UserId: sessionUserid, TaskId: NewRequest.TaskId}); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("PurgeTask Ok")
	}
}
