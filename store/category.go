package store

import (
	"backend/entities"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// SliceScanInterface interface
type SliceScanInterface interface {
	SliceScan() ([]interface{}, error)
}

// DetectData detects whether or not there are values for the given column index
func DetectData(slci SliceScanInterface, index int) bool {
	sl, err := slci.SliceScan()
	switch {
	case err != nil:
		return false
	case sl[index] == nil:
		return false
	default:
		return true
	}
}

// GetListOfCategories gets an cation from the database.
func (st *DbStore) GetListOfCategories(ProfileID int) (categories []entities.Category, err error) {
	qry := `
	SELECT 
	 cat.id, 
	 cat.name, 
	 cat.interval, 
	 cat.interval_type, 
	 cat.created_at, 
	 cat.updated_at, 
	 la.id, 
	 la.subject, 
	 la.action_date 
	FROM category as cat left join action as la on la.id  = (
		SELECT ac.id
		FROM action as ac
		WHERE 
			ac.category_id = cat.id AND ac.action_date is not null
		order by ac.action_date desc, id desc
		LIMIT 1
	)`

	rows, err := st.DB.Queryx(qry)
	if err != nil {
		log.Error(err)
		return
	}
	for rows.Next() {
		var cat entities.Category
		if DetectData(rows, 6) {
			var lastAction entities.Action
			err = rows.Scan(
				&cat.ID,
				&cat.Name,
				&cat.Interval,
				&cat.IntervalType,
				&cat.CreatedAt,
				&cat.UpdatedAt,
				&lastAction.ID,
				&lastAction.Subject,
				&lastAction.ActionDate)
			if err != nil {
				log.Error(err)
				return
			}
			cat.LastAction = &lastAction
		} else {
			var nullvalue *interface{}
			err = rows.Scan(&cat.ID, &cat.Name, &cat.Interval, &cat.IntervalType, &cat.CreatedAt, &cat.UpdatedAt, &nullvalue, &nullvalue, &nullvalue)
			if err != nil {
				log.Error(err)
				return
			}
		}
		cat.CalcDaysSinceLastAction()
		categories = append(categories, cat)
	}

	return
}

// GetCategory gets a category from the database.
func (st *DbStore) GetCategory(id int) (cat entities.Category, err error) {
	qry := `	
	SELECT 
	cat.id, 
	cat.name, 
	cat.interval, 
	cat.interval_type, 
	cat.created_at, 
	cat.updated_at, 
	la.id, 
	la.subject, 
	la.action_date 
   FROM category as cat left join action as la on la.id  = 
   	(
	   SELECT ac.id
	   FROM action as ac
	   WHERE 
		   ac.category_id = cat.id AND ac.action_date is not null
	   order by ac.action_date desc, id desc
	   LIMIT 1
	) 
	WHERE cat.id = $1
	`
	// Using rows so the row is not closed after scanned
	rows, err := st.DB.Queryx(qry, id)
	if err != nil {
		log.Error(err)
		return
	}
	for rows.Next() {
		if DetectData(rows, 6) {
			var lastAction entities.Action
			err = rows.Scan(
				&cat.ID,
				&cat.Name,
				&cat.Interval,
				&cat.IntervalType,
				&cat.CreatedAt,
				&cat.UpdatedAt,
				&lastAction.ID,
				&lastAction.Subject,
				&lastAction.ActionDate)
			if err != nil {
				log.Error("something is up")
				log.Error(err)
				return
			}
			cat.LastAction = &lastAction
		} else {
			var nullvalue *interface{}
			err = rows.Scan(&cat.ID, &cat.Name, &cat.Interval, &cat.IntervalType, &cat.CreatedAt, &cat.UpdatedAt, &nullvalue, &nullvalue, &nullvalue)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
	if cat.ID == 0 {
		err = fmt.Errorf("category with id='%v' not found", id)
		return
	}
	cat.CalcDaysSinceLastAction()
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
