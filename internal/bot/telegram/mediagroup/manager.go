package mediagroup

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/archiver"
)

type MediaGroupManager struct {
	groups     map[string]*MediaGroup
	mutex      sync.RWMutex
	timeout    time.Duration
	onComplete func(*MediaGroup)
}

func NewMediaGroupManager(timeout time.Duration, onComplete func(*MediaGroup)) *MediaGroupManager {
	return &MediaGroupManager{
		groups:     make(map[string]*MediaGroup),
		timeout:    timeout,
		onComplete: onComplete,
	}
}

func (mgm *MediaGroupManager) AddFile(message *tgbotapi.Message, fileData archiver.FileData) {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	mediaGroupID := message.MediaGroupID
	if mediaGroupID == "" {
		// not part of a media group, handle immediately
		group := &MediaGroup{
			ID:     "",
			ChatID: message.Chat.ID,
			Files: []MediaGroupFile{
				{Message: message, FileData: fileData},
			},
			Created: time.Now(),
		}
		go mgm.onComplete(group)
		return
	}

	// check if group already exists
	group, exists := mgm.groups[mediaGroupID]
	if !exists {
		group = &MediaGroup{
			ID:      mediaGroupID,
			ChatID:  message.Chat.ID,
			Files:   []MediaGroupFile{},
			Created: time.Now(),
		}
		mgm.groups[mediaGroupID] = group
	}

	group.Files = append(group.Files, MediaGroupFile{
		Message:  message,
		FileData: fileData,
	})

	// reset timer
	if group.Timer != nil {
		group.Timer.Stop()
	}

	group.Timer = time.AfterFunc(mgm.timeout, func() {
		mgm.completeGroup(mediaGroupID)
	})
}

func (mgm *MediaGroupManager) completeGroup(mediaGroupID string) {
	mgm.mutex.Lock()
	group, exists := mgm.groups[mediaGroupID]
	if !exists {
		mgm.mutex.Unlock()
		return
	}
	delete(mgm.groups, mediaGroupID)
	mgm.mutex.Unlock()

	if group.Timer != nil {
		group.Timer.Stop()
	}

	mgm.onComplete(group)
}

func (mgm *MediaGroupManager) GetPendingGroups() int {
	mgm.mutex.RLock()
	defer mgm.mutex.RUnlock()
	return len(mgm.groups)
}
