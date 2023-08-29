package result

type Result struct {
	err error
	obj any
}

func Ok(obj any) Result {
	return Result{
		err: nil,
		obj: obj,
	}
}

func Err(err error) Result {
	return Result{
		err: err,
		obj: nil,
	}
}

func (r Result) Err() error {
	return r.err
}

func (r Result) Ok() any {
	return r.obj
}
