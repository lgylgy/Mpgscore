package main

import (
	"crypto/tls"
	"net"

	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Controller struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Connect(mongoAddr, database, collection string) error {
	dialInfo, err := mgo.ParseURL(mongoAddr)
	if err != nil {
		return err
	}
	//Below part is similar to above.
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	c.session = session
	c.collection = session.DB(database).C(collection)

	index := mgo.Index{
		Key:        []string{"id"},
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

func (c *Controller) ListPlayers() ([]*Player, error) {
	var result []*Player
	if err := c.collection.Find(nil).Sort("name").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Controller) ListTeamPlayers(team string) ([]*Player, error) {
	var result []*Player
	if err := c.collection.Find(bson.D{{Name: "team", Value: team}}).Sort("name").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Controller) AddPlayer(player *Player) (*Player, error) {
	player.ID = xid.New().String()
	if err := c.collection.Insert(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (c *Controller) GetPlayer(id string) (*Player, error) {
	player := &Player{}
	if err := c.collection.Find(bson.D{{Name: "id", Value: id}}).One(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (c *Controller) UpdatePlayer(player *Player) (*Player, error) {
	if err := c.collection.Update(bson.D{{Name: "id", Value: player.ID}}, player); err != nil {
		return nil, err
	}
	return player, nil
}
