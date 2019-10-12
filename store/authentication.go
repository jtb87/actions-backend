package store

import (
	"backend/entities"

	log "github.com/sirupsen/logrus"
)

// UserGetter gets a user from the database
func (st *DbStore) UserGetter(id int) (entities.Profile, error) {
	user := entities.Profile{}
	err := st.DB.Get(&user, "select id,username,password from users where id = $1", id)
	if err != nil {
		log.Warning(err)
		return entities.Profile{}, err
	}
	return user, nil
}

// AuthorizeToken checks the token and returns the corresponding profile
func (st *DbStore) AuthorizeToken(token string) (p entities.Profile, err error) {
	row := st.DB.QueryRowx("SELECT id, username, password, token_hash from profile where token_hash = $1", token)
	err = row.Scan(&p.ID, &p.Username, &p.Password, &p.Token)
	if err != nil {
		log.Println(err)
	}
	log.Infof("Token authorized for profile with id:%v", p.ID)
	return
}

// ProfileAuthentication gets a user from the database
func (st *DbStore) ProfileAuthentication(p *entities.Profile) (err error) {
	row := st.DB.QueryRowx("SELECT id, token_hash from profile WHERE username=$1 AND password=$2", p.Username, p.Password)
	err = row.Scan(&p.ID, &p.Token)
	if err != nil {
		log.Warning(err)
		return
	}
	log.Infof("Succesfully authenticated profile id: %v", p.ID)
	return
}
