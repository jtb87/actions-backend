package store

import (
	"backend/entities"

	log "github.com/sirupsen/logrus"
)

// GetListOfActions gets an action from the database.
func (st *DbStore) GetListOfActions(ProfileID int) (actions []entities.Action, err error) {

	rows, err := st.DB.Queryx("SELECT id, profile_id, subject, description, category, action_date, planned_date, created_at, updated_at from action")
	if err != nil {
		log.Error(err)
		return
	}
	var act entities.Action
	for rows.Next() {
		err = rows.Scan(&act.ID, &act.ProfileID, &act.Subject, &act.Description, &act.Category, &act.ActionDate, &act.PlannedDate, &act.CreatedAt, &act.UpdateAt)
		if err != nil {
			log.Error(err)
			return
		}
		actions = append(actions, act)
	}

	return
}

// GetActionByID gets an action from the database.
func (st *DbStore) GetActionByID(id int) (act entities.Action, err error) {
	row := st.DB.QueryRowx("SELECT id, profile_id, subject, description, category, action_date, planned_date, created_at, updated_at from action where id = $1", id)
	err = row.Scan(&act.ID, &act.ProfileID, &act.Subject, &act.Description, &act.Category, &act.ActionDate, &act.PlannedDate, &act.CreatedAt, &act.UpdateAt)
	if err != nil {
		log.Error(err)
	}
	return
}

// CreateAction creates a new action entry
func (st *DbStore) CreateAction(a *entities.Action) (err error) {
	stmt, err := st.DB.Preparex("INSERT INTO action(subject, description, category, action_date, planned_date, profile_id) VALUES($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		log.Warn(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(a.Subject, a.Description, a.Category, a.ActionDate, a.PlannedDate, a.ProfileID).Scan(&a.ID)
	if err != nil {
		log.Warn(err)
	}
	log.Infof("Action created with id:%v", a.ID)
	return
}

// UpdateAction creates a new action entry
func (st *DbStore) UpdateAction(a *entities.Action) (err error) {
	stmt, err := st.DB.Preparex("UPDATE action SET subject = $1 , description = $2, category = $3 , action_date = $4, planned_date = $5 where id = $6 RETURNING updated_at")
	if err != nil {
		log.Warn(err)
		return
	}
	err = stmt.QueryRow(a.Subject, a.Description, a.Category, a.ActionDate, a.PlannedDate, a.ID).Scan(&a.UpdateAt)
	if err != nil {
		log.Warn(err)
		return
	}
	return
}

// DeleteAction deletes an action from the database.
func (st *DbStore) DeleteAction(id int) (err error) {
	stmt, err := st.DB.Preparex("DELETE FROM action WHERE id = $1")
	if err != nil {
		log.Error(err)
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Error(err)
	}
	log.Infof("Deleted acion with id %v", id)
	return
}
