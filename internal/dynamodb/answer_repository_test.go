package dynamodb

import (
	"fmt"
	"testing"

	"dochq.co.uk.answerservice/internal/domain"
	"github.com/go-test/deep"
)

func TestAnswerRepository(t *testing.T) {

	// Setup initial dataset.
	//
	answers := []*domain.Answer{
		{
			Key:   "name",
			Value: "John",
		},
		{
			Key:   "country",
			Value: "US",
		},
		{
			Key:   "city",
			Value: "NY",
		},
		{
			Key:   "address",
			Value: "street 1",
		},
	}

	// Test create operations.
	//
	for _, a := range answers {
		if err := testAnswerRepository.Create(a); err != nil {
			t.Errorf("unexpected err: %v", err)
			continue
		}
	}

	// Test get & update operations.
	//
	for _, a := range answers {
		foundAnswer, err := testAnswerRepository.Get(a.Key)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
			continue
		}
		if foundAnswer == nil {
			t.Errorf("answer not found by key: %v", foundAnswer.Key)
			continue
		}
		var (
			updatedAnswer = &domain.Answer{
				Key:   a.Key,
				Value: domain.AnswerValue(fmt.Sprintf("%v-new", a.Value)),
			}
		)
		if err := testAnswerRepository.Update(updatedAnswer); err != nil {
			t.Errorf("unexpected err: %v", err)
			continue
		}
		foundAnswer, err = testAnswerRepository.Get(a.Key)
		if err != nil {
			t.Errorf("unexpected err: %v", err)
			continue
		}
		if foundAnswer == nil {
			t.Errorf("updated answer not found by key: %v", foundAnswer.Key)
			continue
		}
		if diff := deep.Equal(*foundAnswer, *updatedAnswer); diff != nil {
			t.Error(diff)
			continue
		}
	}

	// Test delete operations.
	//
	for _, a := range answers {
		if err := testAnswerRepository.Delete(a.Key); err != nil {
			t.Errorf("unexpected err: %v", err)
			continue
		}
		foundAnswer, _ := testAnswerRepository.Get(a.Key)
		if foundAnswer != nil {
			t.Errorf("answer expected to be deleted: %v", a.Key)
			continue
		}
	}
}
