package database

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// thx https://stackoverflow.com/questions/37889726/how-to-store-a-point-in-postgres-sql-database-using-gorm
func (c Coordinates) Value() (driver.Value, error) {
	out := []byte{'('}
	out = strconv.AppendFloat(out, c.Longitude, 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, c.Latitude, 'f', -1, 64)
	out = append(out, ')')
	return out, nil
}

func (c *Coordinates) Scan(src interface{}) (err error) {
	var data []byte
	switch src := src.(type) {
	case []byte:
		data = src
	case string:
		data = []byte(src)
	case nil:
		return nil
	default:
		return errors.New("(*Point).Scan: unsupported data type")
	}

	if len(data) == 0 {
		return nil
	}

	data = data[1 : len(data)-1] // drop the surrounding parentheses
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			if c.Longitude, err = strconv.ParseFloat(string(data[:i]), 64); err != nil {
				return err
			}
			if c.Latitude, err = strconv.ParseFloat(string(data[i+1:]), 64); err != nil {
				return err
			}
			break
		}
	}
	return nil
}
