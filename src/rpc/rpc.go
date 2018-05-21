package rpcdemo

import "errors"

type JsonServer struct {
}

type Args struct {
	X, Y int
}

func (j JsonServer) Div(args Args,result *float64) error{

	if args.Y == 0 {
		return errors.New("division by zero")
	}

	*result = float64(args.X) / float64(args.Y)
	return nil
}