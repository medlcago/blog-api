COMPOSE=docker-compose -f docker-compose.yaml

.PHONY: up down logs purge

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

purge:
	$(COMPOSE) down --rmi all --volumes --remove-orphans