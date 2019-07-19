package master_side

func (*Job) StartSession(sch *Scheduler) error {
	address, err := sch.bcHandler.CreateComputationContract()
	if err != nil {
		return err
	}

	return nil
}
