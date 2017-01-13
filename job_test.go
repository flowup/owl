package owl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
	"log"
)

type FakeJob struct {
}

func NewFakeJob() *FakeJob {
	return &FakeJob{
	}
}

func (this*FakeJob) Start() JobResult {
	return JobResult{
		Output: "",
		Error:  nil,
	}
}

func (this *FakeJob) Stop() error {
	return nil
}

func TestDebounce(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
	}()

	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))
}

func TestDebounceSleepShort(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(11) * time.Millisecond)
		jobs <- &FakeJob{}

	}()
	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))

}

func TestDebounceSleepLong(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(110) * time.Millisecond)
		jobs <- &FakeJob{}
	}()
	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))
}

func TestDebounceMoreJobs(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(110) * time.Millisecond)
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(11) * time.Millisecond)
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
	}()
	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))

}

func TestScheduler(t *testing.T) {
	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(110) * time.Millisecond)
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
		time.Sleep(time.Duration(11) * time.Millisecond)
		jobs <- &FakeJob{}
		jobs <- &FakeJob{}
	}()
	debounced := Debounce(jobs, amount)

	results := Scheduler(debounced)

	for i := 0; i < 3; i++ {
		out := <-results
		fmt.Print(out.Output)
		if out.Error != nil {
			log.Fatal(out.Error)
		}
		assert.NotNil(t, out)

	}

	assert.Equal(t, 0, len(results))

}
