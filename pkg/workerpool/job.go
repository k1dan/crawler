package workerpool

import "context"

type (
	JobID   string
	JobType string

	ExecutionFunc func(ctx context.Context, args any) (any, error)

	JobDescription struct {
		ID   JobID
		Type JobType
	}

	Result struct {
		Value       any
		Err         error
		Description JobDescription
	}

	Job struct {
		Description JobDescription
		ExecFunc    ExecutionFunc
		Args        any
	}
)

// execute func executes ExecFunc and return Result
func (j Job) execute(ctx context.Context) Result {
	executionResult, err := j.ExecFunc(ctx, j.Args)
	if err != nil {
		return Result{
			Err:         err,
			Description: j.Description,
		}
	}

	return Result{
		Value:       executionResult,
		Description: j.Description,
	}
}
