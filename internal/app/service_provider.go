package app

import (
	"context"
	"log"

	api "github.com/ukrainskykirill/chat-server/internal/api/chats"
	"github.com/ukrainskykirill/chat-server/internal/client/db"
	"github.com/ukrainskykirill/chat-server/internal/client/db/pg"
	"github.com/ukrainskykirill/chat-server/internal/client/db/transaction"
	"github.com/ukrainskykirill/chat-server/internal/closer"
	"github.com/ukrainskykirill/chat-server/internal/config"
	"github.com/ukrainskykirill/chat-server/internal/repository"
	chatRepo "github.com/ukrainskykirill/chat-server/internal/repository/chats"
	msgRepo "github.com/ukrainskykirill/chat-server/internal/repository/messages"
	"github.com/ukrainskykirill/chat-server/internal/service"
	chatServ "github.com/ukrainskykirill/chat-server/internal/service/chats"
	msgServ "github.com/ukrainskykirill/chat-server/internal/service/messages"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

	chatRepo repository.ChatsRepository
	msgRepo  repository.MessagesRepository

	chatServ service.ChatsService
	msgServ  service.MessagesService

	api *api.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PGConfig() config.PGConfig {
	if sp.pgConfig == nil {
		cfg, err := config.NewDBConfig()
		if err != nil {
			log.Fatalf("Error loading config: %s", err.Error())
		}

		sp.pgConfig = cfg
	}

	return sp.pgConfig
}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Error loading config: %s", err.Error())
		}
		sp.grpcConfig = cfg
	}
	return sp.grpcConfig
}

func (sp *serviceProvider) DBClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		cl, err := pg.New(ctx, sp.PGConfig().URL())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		closer.Add(cl.Close)

		sp.dbClient = cl
	}

	return sp.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (sp *serviceProvider) ChatRepo(ctx context.Context) repository.ChatsRepository {
	if sp.chatRepo == nil {
		sp.chatRepo = chatRepo.NewChatRepository(sp.DBClient(ctx))
	}

	return sp.chatRepo
}

func (sp *serviceProvider) MsgRepo(ctx context.Context) repository.MessagesRepository {
	if sp.msgRepo == nil {
		sp.msgRepo = msgRepo.NewMessageRepository(sp.DBClient(ctx))
	}

	return sp.msgRepo
}

func (sp *serviceProvider) ChatService(ctx context.Context) service.ChatsService {
	if sp.chatServ == nil {
		sp.chatServ = chatServ.NewServ(sp.ChatRepo(ctx))
	}

	return sp.chatServ
}

func (sp *serviceProvider) MsgService(ctx context.Context) service.MessagesService {
	if sp.msgServ == nil {
		sp.msgServ = msgServ.NewServ(sp.MsgRepo(ctx), sp.ChatRepo(ctx))
	}

	return sp.msgServ
}

func (sp *serviceProvider) ChatsAPI(ctx context.Context) *api.Implementation {
	if sp.api == nil {
		sp.api = api.NewChatImplementation(sp.ChatService(ctx), sp.MsgService(ctx))
	}

	return sp.api
}