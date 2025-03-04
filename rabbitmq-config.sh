#!/bin/bash
# Execute with: docker exec -it crm-be_rabbitmq_1 /bin/bash -c "$(cat rabbitmq-config.sh)"

# Exchanges
rabbitmqadmin declare exchange name=auth.events type=topic
rabbitmqadmin declare exchange name=user.events type=topic
rabbitmqadmin declare exchange name=feedback.events type=topic
rabbitmqadmin declare exchange name=reward.events type=topic
rabbitmqadmin declare exchange name=inventory.events type=topic
rabbitmqadmin declare exchange name=article.events type=topic
rabbitmqadmin declare exchange name=notification.events type=topic

# Queues
rabbitmqadmin declare queue name=user.auth.events
rabbitmqadmin declare queue name=feedback.point.events
rabbitmqadmin declare queue name=reward.claim.events
rabbitmqadmin declare queue name=inventory.stock.events
rabbitmqadmin declare queue name=notification.email.events
rabbitmqadmin declare queue name=notification.push.events

# Bindings
rabbitmqadmin declare binding source=auth.events destination=user.auth.events routing_key="auth.login.*"
rabbitmqadmin declare binding source=feedback.events destination=feedback.point.events routing_key="feedback.created"
rabbitmqadmin declare binding source=reward.events destination=reward.claim.events routing_key="reward.claimed"
rabbitmqadmin declare binding source=reward.events destination=inventory.stock.events routing_key="reward.claimed"
rabbitmqadmin declare binding source=auth.events destination=notification.email.events routing_key="auth.reset.*"
rabbitmqadmin declare binding source=article.events destination=notification.push.events routing_key="article.created"