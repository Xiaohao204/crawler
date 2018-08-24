package engine


//Scheduler:        &scheduler.QueuedScheduler{},
//WorkerCount:      100,
//ItemChan:         itemChan,
//RequestProcessor: engine.Worker
type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(Request) (ParseResult, error)


type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	//case r := <-s.requestChan:
	//requestQ = append(requestQ, r)
	//case w := <-s.workerChan:
	//workerQ = append(workerQ, w)
	//case activeWorker <- activeRequest:
	//workerQ = workerQ[1:]
	//requestQ = requestQ[1:]
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(),
			out, e.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicate(r.Url) {
			continue
		}
		//s.requestChan <- r
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			go func(i Item) {
				e.ItemChan <- i
			}(item)
		}

		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

func (e *ConcurrentEngine) createWorker(
	in chan Request,
	out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			//s.workerChan <- in
			ready.WorkerReady(in)
			request := <-in
			//engine.Worker
			result, err := e.RequestProcessor(
				request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}
