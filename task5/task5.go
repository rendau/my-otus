package task5

/*
	структура для хранения состояния воркера
*/
type workerSt struct {
	// канал для передачи горутине задании
	input       chan func() error
	inputClosed bool
}

/*
	структура для результата выполнения (горутиной) каждой задачи
*/
type resultSt struct {
	worker *workerSt
	err    error
	closed bool
}

/*
	Возвращает количество успешно выполненных задании и количество ощибок.
	Количество фактических ошибок может быть больше чем errLimit,
	так как на мемент остановки еще работающие горутины могут вернуть ошибку.
	Анолочно и с количеством успешно выполненных задании.
*/
func Fun1(jobs []func() error, workersLen, errLimit int) (int, int) {
	jobsLen := len(jobs)

	if jobsLen == 0 || workersLen == 0 {
		return 0, 0
	}

	if jobsLen < workersLen {
		workersLen = jobsLen
	}

	// Канал для получения результатов
	output := make(chan resultSt)
	defer close(output)

	workers := make([]*workerSt, workersLen)

	for i := 0; i < workersLen; i++ {
		workers[i] = &workerSt{
			input: make(chan func() error),
		}
		// создаем горутины
		go workerRoutine(workers[i], output)
	}

	var jobI int
	var successCnt int
	var errCnt int
	var closedCnt int
	var worker *workerSt

	// первая партия задании
	for _, worker = range workers {
		worker.input <- jobs[jobI]
		jobI++
	}

	for res := range output {
		if res.closed {
			closedCnt++
			if closedCnt >= workersLen {
				break
			}
			continue
		}

		if res.err != nil {
			errCnt++

			if errLimit > 0 && errCnt >= errLimit {
				stopWorkers(workers)
				continue
			}
		} else {
			successCnt++
		}

		if jobI >= jobsLen {
			stopWorkers(workers)
			continue
		}

		if !res.worker.inputClosed {
			res.worker.input <- jobs[jobI]
			jobI++
		}
	}

	return successCnt, errCnt
}

func stopWorkers(workers []*workerSt) {
	for _, worker := range workers {
		if !worker.inputClosed {
			close(worker.input)
			worker.inputClosed = true
		}
	}
}

/*
	Функция горутины. Горутина будет: читать задачи из input, и записывать результат в output
*/
func workerRoutine(ctx *workerSt, output chan<- resultSt) {
	result := resultSt{
		worker: ctx,
	}

	// defer для защиты от паники в задании
	defer func() {
		result.err = nil
		result.closed = true
		output <- result
	}()

	for job := range ctx.input {
		result.err = job()
		output <- result
	}
}
