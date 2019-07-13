package mysql

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	pubsubmodel "github.com/ortuman/jackal/model/pubsub"
)

func (s *Storage) InsertOrUpdatePubSubNode(node *pubsubmodel.Node) error {
	return s.inTransaction(func(tx *sql.Tx) error {
		// if not existing, insert new node
		_, err := sq.Insert("pubsub_nodes").
			Columns("host", "name", "updated_at", "created_at").
			Suffix("ON DUPLICATE KEY UPDATE updated_at = NOW()").
			Values(node.Host, node.Name, nowExpr, nowExpr).
			RunWith(tx).Exec()
		if err != nil {
			return err
		}

		// fetch identifier
		var nodeIdentifier string

		err = sq.Select("id").
			From("pubsub_nodes").
			Where(sq.And{sq.Eq{"host": node.Host}, sq.Eq{"name": node.Name}}).
			RunWith(tx).QueryRow().Scan(&nodeIdentifier)
		if err != nil {
			return err
		}
		// delete previous node options
		_, err = sq.Delete("pubsub_node_options").
			Where(sq.Eq{"node_id": nodeIdentifier}).
			RunWith(tx).Exec()
		if err != nil {
			return err
		}
		// insert new option set
		for name, value := range node.Options.Map() {
			_, err = sq.Insert("pubsub_node_options").
				Columns("node_id", "name", "value", "updated_at", "created_at").
				Values(nodeIdentifier, name, value, nowExpr, nowExpr).
				RunWith(tx).Exec()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetPubSubNode(host, name string) (*pubsubmodel.Node, error) {
	rows, err := sq.Select("name", "value").
		From("pubsub_node_options").
		Where("node_id = (SELECT id FROM pubsub_nodes WHERE host = ? AND name = ?)", host, name).
		RunWith(s.db).Query()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	optMap, err := s.scanNodeOptionsMap(rows)
	if err != nil {
		return nil, err
	}
	opts, err := pubsubmodel.NewOptionsFromMap(optMap)
	if err != nil {
		return nil, err
	}
	return &pubsubmodel.Node{
		Host:    host,
		Name:    name,
		Options: *opts,
	}, nil
}

func (s *Storage) InsertNodeItem(item *pubsubmodel.Item, host, name string, maxNodeItems int) error {
	// TODO(ortuman): implement me!
	return errors.New("unimplemented method")
}

func (s *Storage) GetNodeItems(host, name string) ([]pubsubmodel.Item, error) {
	// TODO(ortuman): implement me!
	return nil, errors.New("unimplemented method")
}

func (s *Storage) InsertPubSubNodeItem(item *pubsubmodel.Item, host, name string, maxNodeItems int) error {
	// TODO(ortuman): implement me!
	return errors.New("unimplemented method")
}

func (s *Storage) GetPubSubNodeItems(host, name string) ([]pubsubmodel.Item, error) {
	// TODO(ortuman): implement me!
	return nil, errors.New("unimplemented method")
}

func (s *Storage) InsertPubSubNodeAffiliation(affiliatiaon *pubsubmodel.Affiliation, host, name string) error {
	// TODO(ortuman): implement me!
	return errors.New("unimplemented method")
}

func (s *Storage) GetPubSubNodeAffiliation(host, name string) ([]pubsubmodel.Affiliation, error) {
	// TODO(ortuman): implement me!
	return nil, errors.New("unimplemented method")
}

func (s *Storage) scanNodeOptionsMap(scanner rowsScanner) (map[string]string, error) {
	var optMap = make(map[string]string)
	for scanner.Next() {
		var opt, value string
		if err := scanner.Scan(&opt, &value); err != nil {
			return nil, err
		}
		optMap[opt] = value
	}
	return optMap, nil
}
