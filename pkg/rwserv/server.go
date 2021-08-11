package rwserv

import (
	"encoding/json"
	"github.com/rwlist/engine/pkg/auth"
	"github.com/rwlist/engine/pkg/domain"
	"github.com/rwlist/engine/pkg/jsonrpc"
)

type Server struct {
	dbms domain.DBMS
}

func NewServer(dbms domain.DBMS) *Server {
	return &Server{
		dbms: dbms,
	}
}

func (s *Server) Handle(request *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	var (
		method = request.Method
		params = request.Params
		user   = &auth.User{IsAdmin: true} // TODO: use tokens for auth
	)

	var (
		res jsonrpc.Result
		err error
	)

	switch method {
	case "databases.getAll":
		res, err = s.DatabasesGetAll(user)

	case "databases.create":
		var req CreateDatabaseRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.DatabasesCreate(user, &req)

	case "databases.drop":
		var req DropDatabaseRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.DatabasesDrop(user, &req)

	case "lists.getAll":
		var req GetAllListsRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.ListsGetAll(user, &req)

	case "lists.create":
		var req CreateListRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.ListsCreate(user, &req)

	case "lists.drop":
		var req DropListRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.ListsDrop(user, &req)

	case "list.insertMany":
		var req InsertManyRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.ListInsertMany(user, &req)

	case "list.readRange":
		var req ReadRangeRequest
		if err = json.Unmarshal(params, &req); err != nil {
			break
		}
		res, err = s.ListReadRange(user, &req)

	default:
		return nil, &jsonrpc.MethodNotFound
	}

	if err != nil {
		return nil, &jsonrpc.Error{
			Message: err.Error(), // TODO: replace with "internal error"
		}
	}

	return res, nil
}

func (s *Server) DatabasesGetAll(user *auth.User) (*AllDatabasesResponse, error) {
	dbs, err := s.dbms.AllDatabases(user)
	if err != nil {
		return nil, err
	}

	infos := []domain.DatabaseInfo{}
	for _, db := range dbs {
		info, err := db.Info()
		if err != nil {
			return nil, err
		}

		infos = append(infos, *info)
	}

	return &AllDatabasesResponse{
		Databases: infos,
	}, nil
}

func (s *Server) DatabasesCreate(user *auth.User, req *CreateDatabaseRequest) (*CreateDatabaseResponse, error) {
	db, err := s.dbms.CreateDatabase(user, req.Database)
	if err != nil {
		return nil, err
	}

	info, err := db.Info()
	if err != nil {
		return nil, err
	}

	info2 := CreateDatabaseResponse(*info)
	return &info2, nil
}

func (s *Server) DatabasesDrop(user *auth.User, req *DropDatabaseRequest) (*DropDatabaseResponse, error) {
	err := s.dbms.DropDatabase(user, req.Database)
	if err != nil {
		return nil, err
	}

	return &DropDatabaseResponse{}, nil
}

func (s *Server) ListsGetAll(user *auth.User, req *GetAllListsRequest) (*GetAllListsResponse, error) {
	db, err := s.dbms.Database(user, req.Database)
	if err != nil {
		return nil, err
	}

	lists, err := db.AllLists()
	if err != nil {
		return nil, err
	}

	if lists == nil {
		lists = []domain.ListInfo{}
	}

	return &GetAllListsResponse{
		Lists: lists,
	}, nil
}

func (s *Server) ListsCreate(user *auth.User, req *CreateListRequest) (*CreateListResponse, error) {
	db, err := s.dbms.Database(user, req.Database)
	if err != nil {
		return nil, err
	}

	info, err := db.CreateList(&domain.CreateListRequest{
		ListName: req.ListName,
		Engine:   req.Engine,
	})
	if err != nil {
		return nil, err
	}

	info2 := CreateListResponse(*info)
	return &info2, nil
}

func (s *Server) ListsDrop(user *auth.User, req *DropListRequest) (*DropListResponse, error) {
	db, err := s.dbms.Database(user, req.Database)
	if err != nil {
		return nil, err
	}

	err = db.DropList(req.ListName)
	if err != nil {
		return nil, err
	}

	return &DropListResponse{}, nil
}

func (s *Server) ListInsertMany(user *auth.User, req *InsertManyRequest) (*InsertManyResponse, error) {
	db, err := s.dbms.Database(user, req.Database)
	if err != nil {
		return nil, err
	}

	err = db.OpenList(req.ListName, func(list domain.List) error {
		for _, e := range req.Entries {
			err := list.Insert(domain.Entity(e))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &InsertManyResponse{}, nil
}

func (s *Server) ListReadRange(user *auth.User, req *ReadRangeRequest) (*ReadRangeResponse, error) {
	db, err := s.dbms.Database(user, req.Database)
	if err != nil {
		return nil, err
	}

	var entities []domain.Entity
	err = db.OpenList(req.ListName, func(list domain.List) error {
		entities, err = list.ReadRange(req.Offset, req.Limit)
		return err
	})
	if err != nil {
		return nil, err
	}

	res := []json.RawMessage{}
	for _, e := range entities {
		res = append(res, json.RawMessage(e))
	}

	return &ReadRangeResponse{
		Entries: res,
	}, nil
}
