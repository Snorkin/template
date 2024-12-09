package middleware

import (
	"example/config"
	"example/pkg/logger"
)

type MdwManager struct {
	cfg config.Config
	log logger.Logger
}

func NewMdwManager(
	cfg config.Config,
	log logger.Logger,
) MdwManager {
	return MdwManager{
		cfg: cfg,
		log: log,
	}
}
