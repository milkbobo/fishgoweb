package queue

import (
	"errors"
	. "github.com/milkbobo/fishgoweb/app/log"
	. "github.com/milkbobo/fishgoweb/util"
	"sync"
	"sync/atomic"
	"unsafe"
)

type memoryQueueChannel struct {
	channel *InfiniteChannel
}

type memoryQueueStore struct {
	router        map[string]map[string]*memoryQueueChannel
	routerPointer *map[string][]*memoryQueueChannel
	mutex         sync.Mutex
	waitgroup     *sync.WaitGroup
	exitChan      chan bool
}

func newMemoryQueue(log Log, config QueueConfig) (queueStoreInterface, error) {
	result := &memoryQueueStore{}
	result.router = map[string]map[string]*memoryQueueChannel{}
	result.routerPointer = &map[string][]*memoryQueueChannel{}
	result.waitgroup = &sync.WaitGroup{}
	result.exitChan = make(chan bool, 16)
	return result, nil
}

func (this *memoryQueueStore) getRouter() map[string][]*memoryQueueChannel {
	router := *(*map[string][]*memoryQueueChannel)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&this.routerPointer)),
	))
	return router
}

func (this *memoryQueueStore) setRouter(topicId string, queueName string) (*memoryQueueChannel, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	queues, isExist := this.router[topicId]
	if isExist == false {
		queues = map[string]*memoryQueueChannel{}
		this.router[topicId] = queues
	}
	_, isExist = queues[queueName]
	if isExist == true {
		return nil, errors.New("has exist queue " + queueName)
	}
	result := &memoryQueueChannel{
		channel: NewInfiniteChannel(),
	}
	queues[queueName] = result

	newRouter := map[string][]*memoryQueueChannel{}
	for topicId, single := range this.router {
		singleChannel := []*memoryQueueChannel{}
		for _, single2 := range single {
			singleChannel = append(singleChannel, single2)
		}
		newRouter[topicId] = singleChannel
	}
	atomic.StorePointer(
		(*unsafe.Pointer)(unsafe.Pointer(&this.routerPointer)),
		unsafe.Pointer(&newRouter),
	)
	return result, nil
}

func (this *memoryQueueStore) Produce(topicId string, data []byte) error {
	router := this.getRouter()
	queues, isExist := router[topicId]
	if isExist == false {
		return errors.New("dos not exist topicId " + topicId)
	}
	for _, queue := range queues {
		this.waitgroup.Add(1)
		queue.channel.Write(data)
	}
	return nil
}

func (this *memoryQueueStore) Consume(topicId string, queueName string, poolSize int, listener queueStoreListener) error {
	queue, err := this.setRouter(topicId, queueName)
	if err != nil {
		return err
	}
	channel := queue.channel
	for i := 0; i < poolSize; i++ {
		go func() {
			for single := range channel.Read() {
				listener(single.([]byte))
				this.waitgroup.Done()
			}
		}()
	}
	return nil
}

func (this *memoryQueueStore) Run() error {
	this.waitgroup.Wait()
	this.exitChan <- true
	return nil
}

func (this *memoryQueueStore) Close() {
	<-this.exitChan
}
