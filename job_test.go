package owl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

type fakeJob struct {
}

func TestDebounce(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
	}()

	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))
}

func TestDebounceSleepShort(t *testing.T) {

	jobs := make(chan Job)
	amount := int64(10)

	go func() {
		jobs <- &fakeJob{}
		time.Sleep(time.Duration(11) * time.Millisecond)
		jobs <- &fakeJob{}

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
		jobs <- &fakeJob{}
		time.Sleep(time.Duration(110) * time.Millisecond)
		jobs <- &fakeJob{}
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
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		time.Sleep(time.Duration(110) * time.Millisecond)
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
		time.Sleep(time.Duration(11) * time.Millisecond)
		jobs <- &fakeJob{}
		jobs <- &fakeJob{}
	}()
	results := Debounce(jobs, amount)

	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.NotNil(t, <-results)
	assert.Equal(t, 0, len(results))


}