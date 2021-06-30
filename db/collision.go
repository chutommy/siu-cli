package db

import "github.com/chutommy/siu/models"

// CheckCollision checks the unique of the motion.
func CheckCollision(m, exception models.Motion) (models.Motion, bool) {
	searches := []string{m.Name, m.URL, m.Shortcut}

	for _, search := range searches {
		if item, err := ReadOne(search); err == nil {
			if item != exception {
				return m, true
			}
		}
	}

	return models.Motion{}, false
}
