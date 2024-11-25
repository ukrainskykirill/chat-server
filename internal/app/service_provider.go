package app

import (
	"context"
	"log"

	"github.com/ukrainskykirill/platform_common/pkg/closer"
	"github.com/ukrainskykirill/platform_common/pkg/db"
	"github.com/ukrainskykirill/platform_common/pkg/db/pg"
	"github.com/ukrainskykirill/platform_common/pkg/db/transaction"

	api "github.com/ukrainskykirill/chat-server/internal/api/chats"
	"github.com/ukrainskykirill/chat-server/internal/config"
	"github.com/ukrainskykirill/chat-server/internal/repository"
	chatRepo "github.com/ukrainskykirill/chat-server/internal/repository/chats"
	msgRepo "github.com/ukrainskykirill/chat-server/internal/repository/messages"
	"github.com/ukrainskykirill/chat-server/internal/service"
	chatServ "github.com/ukrainskykirill/chat-server/internal/service/chats"
	msgServ "github.com/ukrainskykirill/chat-server/internal/service/messages"
	streamServ "github.com/ukrainskykirill/chat-server/internal/service/stream"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

	chatRepo repository.ChatsRepository
	msgRepo  repository.MessagesRepository

	chatServ   service.ChatsService
	msgServ    service.MessagesService
	streamServ service.StreamService

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

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.DBClient(ctx).DB())
	}

	return sp.txManager
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

func (sp *serviceProvider) StreamService(ctx context.Context) service.StreamService {
	if sp.streamServ == nil {
		sp.streamServ = streamServ.NewStreamService()
	}

	return sp.streamServ
}

func (sp *serviceProvider) ChatService(ctx context.Context) service.ChatsService {
	if sp.chatServ == nil {
		sp.chatServ = chatServ.NewServ(
			sp.TxManager(ctx), sp.ChatRepo(ctx), sp.MsgRepo(ctx), sp.StreamService(ctx),
		)
	}

	return sp.chatServ
}

func (sp *serviceProvider) MsgService(ctx context.Context) service.MessagesService {
	if sp.msgServ == nil {
		sp.msgServ = msgServ.NewServ(sp.MsgRepo(ctx), sp.ChatRepo(ctx), sp.StreamService(ctx))
	}

	return sp.msgServ
}

func (sp *serviceProvider) ChatsAPI(ctx context.Context) *api.Implementation {
	if sp.api == nil {
		sp.api = api.NewChatImplementation(sp.ChatService(ctx), sp.MsgService(ctx))
	}

	return sp.api
}
