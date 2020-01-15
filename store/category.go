package store

import (
	"backend/entities"

	log "github.com/sirupsen/logrus"
)

// GetListOfCategories gets an cation from the database.
func (st *DbStore) GetListOfCategories(ProfileID int) (categories []entities.Category, err error) {
	rows, err := st.DB.Queryx("SELECT id, name, interval, interval_type, created_at, updated_at from category")
	if err != nil {
		log.Error(err)
		return
	}
	var cat entities.Category
	for rows.Next() {
		err = rows.Scan(&cat.ID, &cat.Name, &cat.Interval, &cat.IntervalType, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			log.Error(err)
			return
		}
		categories = append(categories, cat)
	}

	return
}

// GetCategory gets a category from the database.
func (st *DbStore) GetCategory(id int) (cat entities.Category, err error) {
	row := st.DB.QueryRowx("SELECT id, name, interval, interval_type, created_at, updated_at from category where id = $1", id)
	err = row.Scan(&cat.ID, &cat.Name, &cat.Interval, &cat.IntervalType, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		log.Error(err)
	}
	return
}

// CreateCategory creates a new action entry
func (st *DbStore) CreateCategory(c *entities.Category) (err error) {
	stmt, err := st.DB.Preparex("INSERT INTO category(name, interval, interval_type) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		log.Warn(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.Name, c.Interval, c.IntervalType).Scan(&c.ID)
	if err != nil {
		log.Warn(err)
	}
	log.Infof("Category created with id:%v", c.ID)
	return
}

// DeleteCategory deletes an action from the database.
func (st *DbStore) DeleteCategory(id int) (err error) {
	stmt, err := st.DB.Preparex("DELETE FROM category WHERE id = $1")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Deleted category with id %v", id)
	return
}
