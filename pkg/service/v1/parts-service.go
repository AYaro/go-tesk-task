package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	v1 "github.com/AYaro/go-test-task/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type partServiceServer struct {
	db *sql.DB
}

func NewPartServiceServer(db *sql.DB) v1.PartServiceServer {
	return &partServiceServer{db: db}
}

// checkAPI проверяет совпадает ли версися апи клиента и сервиса
func (s *partServiceServer) checkAPI(api string) error {
	if len(api) > 0 && apiVersion != api {
		return status.Errorf(codes.Unimplemented,
			"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
	}
	return nil
}

func (s *partServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database: "+err.Error())
	}
	return c, nil
}

// Создание новой детали
func (s *partServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// проверка на совпадание версии API клиента и сервера
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var id int
	err := s.db.QueryRow("INSERT INTO part (manufacturer_id , vendor_code) VALUES($1 , $2) RETURNING id", req.Part.ManufacturerId, req.Part.VendorCode).Scan(&id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert data into Part: "+err.Error())
	}
	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  int64(id),
	}, nil
}

func (s *partServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	var createdAt time.Time
	part := v1.Part{}
	err := s.db.QueryRow("SELECT id, manufacturer_id, vendor_code, created_at FROM part WHERE id=$1 AND deleted_at IS NULL", req.Id).Scan(&part.Id, &part.ManufacturerId, &part.VendorCode, &createdAt)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select fields from Part: "+err.Error())
	}

	part.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		return nil, status.Error(codes.Unknown, "createdAt field has invalid format: "+err.Error())
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		Part: &part,
	}, nil

}

func (s *partServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	res, err := s.db.Exec("UPDATE Part SET manufacturer_id=$1, vendor_code=$2 WHERE ID=$3 AND deleted_at IS NULL",
		req.Part.ManufacturerId, req.Part.VendorCode, req.Part.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update Part: "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Part with ID='%d' is not found",
			req.Part.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
	}, nil
}

func (s *partServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	res, err := s.db.Exec("UPDATE part SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1 AND deleted_at IS NULL", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete Part: "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Part with ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: rows,
	}, nil
}

func (s *partServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT id, manufacturer_id, vendor_code, created_at FROM part WHERE deleted_at IS NULL")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from Part: "+err.Error())
	}
	defer rows.Close()

	var createdAt time.Time
	list := []*v1.Part{}
	for rows.Next() {
		part := &v1.Part{}
		if err := rows.Scan(&part.Id, &part.ManufacturerId, &part.VendorCode, &createdAt); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve values from Part: "+err.Error())
		}
		part.CreatedAt, err = ptypes.TimestampProto(createdAt)
		if err != nil {
			return nil, status.Error(codes.Unknown, "createdAt field has invalid format: "+err.Error())
		}
		list = append(list, part)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from Part: "+err.Error())
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		Parts: list,
	}, nil
}
