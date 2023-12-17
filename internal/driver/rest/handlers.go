package rest

import (
	"encoding/json"
	"io"
	cloud_storage "moneyman-writer-go/internal/adapter/google/cloud-storage"
	"moneyman-writer-go/internal/core"
	"moneyman-writer-go/internal/model"
	x "moneyman-writer-go/internal/utils/logger"
	"net/http"
)

type RestController struct {
	svc        *core.Service
	downloader core.ObjectDownloader
}

func NewRestController(svc *core.Service, d core.ObjectDownloader) *RestController {
	return &RestController{
		svc:        svc,
		downloader: d,
	}
}

func (c *RestController) HandleGcsTransactionsUploadedEvent(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		x.Logger().Errorw("failed to read request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := cloud_storage.GcsEvent{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		x.Logger().Errorw("failed to parse request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.svc.SaveNewTransactionsFromObjectFile(r.Context(), c.downloader, &model.TransactionsFileUploadedEvent{
		Bucket:      event.Bucket,
		Name:        event.Name,
		TimeCreated: event.TimeCreated,
	})
	if err != nil {
		x.Logger().Errorw("failed to save new transactions", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
