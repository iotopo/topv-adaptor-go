package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type ValueItem struct {
	Tag     string    `json:"tag"`
	Time    time.Time `json:"timestamp"` //"2017-07-14T02:40:09.999Z"
	Value   any       `json:"value"`
	Quality int       `json:"quality"`
}

type DataItem struct {
	Value any       `json:"value,omitempty"`
	Time  time.Time `json:"time,omitempty"`
}

type Result struct {
	Tag    string     `json:"tag"`
	Values []DataItem `json:"values"`
}
type HistoryResponse struct {
	Results []Result `json:"results"`
	Msg     string   `json:"msg,omitempty"`
	Code    string   `json:"code,omitempty"`
}

type TagPoint struct {
	Tag  string `json:"tag,omitempty"`
	Name string `json:"name,omitempty"`
}
type Device struct {
	ParentTag string    `json:"-"`
	Tag       string    `json:"tag,omitempty"`
	Name      string    `json:"name,omitempty"`
	Children  []*Device `json:"children,omitempty"`
	IsDevice  bool      `json:"isDevice,omitempty"`
}

func findLast(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		ProjectID string `json:"projectID,omitempty"`
		Tag       string `json:"tag,omitempty"`
		Device    bool   `json:"device,omitempty"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Fprintf(w, "ProjectID: %s, Tag: %s", reqBody.ProjectID, reqBody.Tag)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if reqBody.Device {
		// 查询设备标签下的所有测点
		json.NewEncoder(w).Encode([]ValueItem{
			{
				Tag:   reqBody.Tag,
				Time:  time.Now(),
				Value: "12.3",
			},
		})
	} else {
		// 查询单个测点
		json.NewEncoder(w).Encode(ValueItem{
			Tag:   reqBody.Tag,
			Time:  time.Now(),
			Value: "12.3",
		})
	}
}

func setValue(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		ProjectID string `json:"projectID,omitempty"`
		Tag       string `json:"tag,omitempty"`
		Value     string `json:"value"`
		Time      int64  `json:"time"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respData := struct {
		Code string `json:"code,omitempty"`
		Msg  string `json:"msg,omitempty"`
	}{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respData)
}

func queryHistory(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		ProjectID string   `json:"projectID,omitempty"`
		Tag       []string `json:"tag,omitempty"`
		Interval  string   `json:"interval,omitempty"`
		Start     string   `json:"start,omitempty"`
		End       string   `json:"end,omitempty"`
		Offset    int32    `json:"offset,omitempty"`
		Limit     int32    `json:"limit,omitempty"`
		Order     string   `json:"order,omitempty"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respData := HistoryResponse{
		Results: []Result{
			{Tag: reqBody.Tag[0], Values: []DataItem{{Value: "12.3", Time: time.Now()}}},
		},
	}
	json.NewEncoder(w).Encode(respData)
}

func queryPoins(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		ProjectID string `json:"projectID,omitempty"`
		ParentTag string `json:"parentTag,omitempty"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]TagPoint{
		{Tag: reqBody.ParentTag + ".a", Name: "a"},
		{Tag: reqBody.ParentTag + ".b", Name: "b"},
		{Tag: reqBody.ParentTag + ".c", Name: "c"},
	})
}

func queryDevices(w http.ResponseWriter, req *http.Request) {
	var reqBody struct {
		ProjectID string `json:"projectID,omitempty"`
	}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]Device{
		{Tag: "group1", Name: "group1", Children: []*Device{
			{Tag: "group1.dev1", Name: "dev1", IsDevice: true},
			{Tag: "group1.dev2", Name: "dev2", IsDevice: true},
			{Tag: "group1.dev3", Name: "dev3", IsDevice: true},
		}},
		{Tag: "group2", Name: "group2", Children: []*Device{
			{Tag: "group2.dev1", Name: "dev1", IsDevice: true},
			{Tag: "group2.dev2", Name: "dev2", IsDevice: true},
			{Tag: "group2.dev3", Name: "dev3", IsDevice: true},
		}},
		{Tag: "group3", Name: "group3", Children: []*Device{
			{Tag: "group3.dev1", Name: "dev1", IsDevice: true},
			{Tag: "group3.dev2", Name: "dev2", IsDevice: true},
			{Tag: "group3.dev3", Name: "dev3", IsDevice: true},
		}},
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go realPush(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/find_last", findLast)
	mux.HandleFunc("POST /api/query_history", queryHistory)
	mux.HandleFunc("POST /api/set_value", setValue)
	mux.HandleFunc("GET /api/query_points", queryPoins)
	mux.HandleFunc("GET /api/query_devices", queryDevices)

	slog.Info("run web server at :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
