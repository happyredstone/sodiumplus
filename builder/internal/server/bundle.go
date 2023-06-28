package server

func Bundle() error {
	err := Tar()

	if err != nil {
		return err
	}

	err = Zip()

	if err != nil {
		return err
	}

	err = Rename()

	if err != nil {
		return err
	}

	return nil
}

func CleanBundle() error {
	err := Clean()

	if err != nil {
		return err
	}

	err = Bundle()

	if err != nil {
		return err
	}

	return nil
}
