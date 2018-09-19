package db

import (
	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mpgscore/api"
)

type Controller struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Connect(mongoAddr, database, collection string) error {
	session, err := mgo.Dial(mongoAddr)
	if err != nil {
		return err
	}
	c.session = session
	c.collection = session.DB(database).C(collection)

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.collection.EnsureIndex(index)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) Close() {
	c.session.Close()
}

func (c *Controller) ListPlayers() ([]*api.Player, error) {
	var result []*api.Player
	if err := c.collection.Find(nil).Sort("name").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Controller) ListTeamPlayers(team string) ([]*api.Player, error) {
	var result []*api.Player
	if err := c.collection.Find(bson.D{{Name: "team", Value: team}}).Sort("name").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Controller) AddPlayer(player *api.Player) (*api.Player, error) {
	player.ID = xid.New().String()
	if err := c.collection.Insert(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (c *Controller) GetPlayer(id string) (*api.Player, error) {
	player := &api.Player{}
	if err := c.collection.Find(bson.D{{Name: "id", Value: id}}).One(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (c *Controller) UpdatePlayer(player *api.Player) (*api.Player, error) {
	if err := c.collection.Update(bson.D{{Name: "id", Value: player.ID}}, player); err != nil {
		return nil, err
	}
	return player, nil
}
