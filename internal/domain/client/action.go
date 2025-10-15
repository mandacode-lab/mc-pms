package client

import "time"

func (c *Client) UpdateName(name Name) error {
	if valid, err := name.IsValid(); !valid {
		return err
	}

	c.name = name
	c.updatedAt = time.Now().UTC()
	return nil
}

func (c *Client) UpdateDescription(description Description) error {
	if valid, err := description.IsValid(); !valid {
		return err
	}

	c.description = description
	c.updatedAt = time.Now().UTC()
	return nil
}

func (c *Client) RegenerateSecret(secretHash SecretHash) error {
	if valid, err := secretHash.IsValid(); !valid {
		return err
	}

	c.secretHash = secretHash
	c.updatedAt = time.Now().UTC()
	return nil
}

func (c *Client) Activate() error {
	c.isActive = true
	c.updatedAt = time.Now().UTC()
	return nil
}

func (c *Client) Deactivate() error {
	c.isActive = false
	c.updatedAt = time.Now().UTC()
	return nil
}
