package websocket

import "sync"

type hubManager struct {
	hubs map[int]*Hub
}

var (
	once sync.Once
	hubM *hubManager
)

func NewHubManager() *hubManager {
	once.Do(func() {
		hubM = &hubManager{
			hubs: make(map[int]*Hub),
		}
	})
	return hubM
}

func (hm hubManager) GetHub(roomId int) *Hub {
	_, ok := hm.hubs[roomId]
	if !ok {
		hm.hubs[roomId] = NewHub()
		go hm.hubs[roomId].Run()
	}

	return hm.hubs[roomId]
}

func (hm hubManager) RemoveHub(roomId int) {
	delete(hm.hubs, roomId)
}
