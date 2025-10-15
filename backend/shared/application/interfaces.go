// Package application contains use case interfaces and DTOspackage application

package application

import (
	"context"
)

// UseCase represents a use case interface
type UseCase[TRequest any, TResponse any] interface {
	Execute(ctx context.Context, request TRequest) (TResponse, error)
}

// CommandHandler handles commands (write operations)
type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

// QueryHandler handles queries (read operations)
type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

// EventHandler handles domain events
type EventHandler[TEvent any] interface {
	Handle(ctx context.Context, event TEvent) error
}

// Command represents a command (write operation)
type Command interface {
	GetCommandType() string
}

// Query represents a query (read operation)  
type Query interface {
	GetQueryType() string
}

// BaseCommand provides base implementation for commands
type BaseCommand struct {
	CommandType string `json:"command_type"`
}

func NewBaseCommand(commandType string) BaseCommand {
	return BaseCommand{
		CommandType: commandType,
	}
}

func (c BaseCommand) GetCommandType() string {
	return c.CommandType
}

// BaseQuery provides base implementation for queries
type BaseQuery struct {
	QueryType string `json:"query_type"`
}

func NewBaseQuery(queryType string) BaseQuery {
	return BaseQuery{
		QueryType: queryType,
	}
}

func (q BaseQuery) GetQueryType() string {
	return q.QueryType
}

// Result represents operation result with success/failure
type Result[T any] struct {
	Data    T      `json:"data,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// NewSuccessResult creates a successful result
func NewSuccessResult[T any](data T) Result[T] {
	return Result[T]{
		Data:    data,
		Success: true,
	}
}

// NewErrorResult creates a failed result
func NewErrorResult[T any](err string) Result[T] {
	var zero T
	return Result[T]{
		Data:    zero,
		Success: false,
		Error:   err,
	}
}