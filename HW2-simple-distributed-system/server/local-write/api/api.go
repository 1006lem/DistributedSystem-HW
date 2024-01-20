package api

import (
	"context"
	"distributed-system/server/common"
	"distributed-system/server/env"
	"distributed-system/server/storage"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Handler struct {
	engine *gin.Engine
	server *http.Server
}

// New initializes Handler
func New() *Handler {
	return &Handler{
		engine: gin.Default(),
		server: &http.Server{},
	}
}

// InitRoutes initializes all routes for endpoint
func (h *Handler) InitRoutes() {
	// note
	h.engine.GET("/note", h.getAllNotesHandler)
	h.engine.GET("/note/:id", h.getNoteHandler)
	h.engine.POST("/note", h.postNoteHandler)
	h.engine.PUT("/note/:id", h.putNoteHandler)
	h.engine.PATCH("/note/:id", h.patchNoteHandler)
	h.engine.DELETE("/note/:id", h.deleteNoteHandler)

	// primary
	h.engine.GET("/primary/:id", h.getPrimaryHandler)

	// backup
	h.engine.POST("/backup", h.postBackupHandler)
	h.engine.PUT("/backup/:id", h.putBackupHandler)
	h.engine.PATCH("/backup/:id", h.patchBackupHandler)
	h.engine.DELETE("/backup/:id", h.deleteBackupHandler)
}

/*
// Run starts listening service in address
func (h *Handler) Run(wg *sync.WaitGroup) error {
	port := env.Config.ServicePort
	ip := "127.0.0.1"
	addr := ip + ":" + strconv.Itoa(port)

	log.Println("Starting local-write REST API Server at ", addr)
	go func() {
		defer wg.Done()
		err := h.engine.Run(addr)
		if err != nil {
			log.Fatalf("could not start API server in %s", addr)
		}
	}()
	select {}
}
*/

// Run starts listening service and handles graceful shutdown
func (h *Handler) Run(wg *sync.WaitGroup) {
	h.server.Addr = "127.0.0.1:" + strconv.Itoa(env.Config.ServicePort)
	h.server.Handler = h.engine

	go func() {
		defer wg.Done()
		log.Println("Starting local-write REST API Server at ", h.server.Addr)

		// PM이 아니라면, PM으로부터 memo를 다 가져옴 ([GET] /memo 요청 보냄)
		if !env.IsPM {
			PMAddress := env.ReplicaAddresses.Addresses[0]
			url := "http://" + PMAddress.IP + ":" + strconv.Itoa(PMAddress.Port) + "/note"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error sending GET request:", err)
				return
			}
			defer resp.Body.Close()

			// 응답 확인
			if resp.StatusCode == http.StatusOK {
				var notes common.Notes
				err := json.NewDecoder(resp.Body).Decode(&notes)
				if err != nil {
					fmt.Println("Error decoding JSON response:", err)
					return
				}
				fmt.Println("Initialize Data from PM ..")
				fmt.Printf("Response from %s:%d:\n%+v\n", PMAddress.IP, PMAddress.Port, notes)

				// update local storage
				err = storage.Manager.OpenAndWriteData(notes)
				if err != nil {
					return
				}
			} else {
				fmt.Printf("Error response from %s:%d: %s\n", PMAddress.IP, PMAddress.Port, resp.Status)
			}
		}

		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	h.handleShutdownSignal()
}

func (h *Handler) handleShutdownSignal() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v\n", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
