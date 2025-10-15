package svcclient

import "time"

func (sc *SvcClient) UpdateName(name Name) error {
	if valid, err := name.IsValid(); !valid {
		return err
	}

	sc.name = name
	sc.updatedAt = time.Now().UTC()
	return nil
}

func (sc *SvcClient) UpdateDescription(description Description) error {
	if valid, err := description.IsValid(); !valid {
		return err
	}

	sc.description = description
	sc.updatedAt = time.Now().UTC()
	return nil
}

func (sc *SvcClient) RegenerateSecret(secretHash SecretHash) error {
	if valid, err := secretHash.IsValid(); !valid {
		return err
	}

	sc.secretHash = secretHash
	sc.updatedAt = time.Now().UTC()
	return nil
}

func (sc *SvcClient) Activate() error {
	sc.isActive = true
	sc.updatedAt = time.Now().UTC()
	return nil
}

func (sc *SvcClient) Deactivate() error {
	sc.isActive = false
	sc.updatedAt = time.Now().UTC()
	return nil
}
