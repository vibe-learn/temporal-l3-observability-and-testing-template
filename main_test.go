package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

// TestTaskQueueDefault — pure unit test. No Temporal required.
func TestTaskQueueDefault(t *testing.T) {
	if os.Getenv("TEMPORAL_TASK_QUEUE") == "" && TaskQueue() != "lesson-l3_observability_and_testing-tq" {
		t.Errorf("TaskQueue() default = %q, want lesson-l3_observability_and_testing-tq", TaskQueue())
	}
}

// TestPayoutWorkflow — юнит-тест на TestWorkflowEnvironment: активности
// замоканы через env.OnActivity(...).Return(...), проверяем, что workflow
// доходит до конца без ошибки. Бежит в CI без Temporal-сервера.
//
// По мере реализации добавляй проверки порядка вызовов, ретраев, компенсаций
// и т.п. согласно README.md (см. env.AssertExpectations, env.OnActivity().Once()).
func TestPayoutWorkflow(t *testing.T) {
	ts := &testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	env.OnActivity(Payout, mock.Anything).Return("ok", nil)
	env.OnActivity(ReversePayout, mock.Anything).Return("ok", nil)
	env.ExecuteWorkflow(PayoutWorkflow)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

// TestIntegration — требует запущенный Temporal (docker compose up -d).
// SKIPPED по умолчанию; включается через TEMPORAL_INTEGRATION=1.
func TestIntegration(t *testing.T) {
	if os.Getenv("TEMPORAL_INTEGRATION") == "" {
		t.Skip("set TEMPORAL_INTEGRATION=1 and run `docker compose up -d` to enable")
	}
	// TODO: client.Dial(TemporalAddress()), запусти воркер на TaskQueue(),
	// стартани PayoutWorkflow через client.ExecuteWorkflow и проверь результат
	// поведения урока «Observability и тестирование».
}
