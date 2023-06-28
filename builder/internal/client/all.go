package client

func Bundle() error {
	err := CurseForge()

	if err != nil {
		return err
	}

	err = Modrinth()

	if err != nil {
		return err
	}

	return nil
}
