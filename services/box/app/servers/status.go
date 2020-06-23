package servers

import "telegram_boxes/services/box/protobuf/services/core/protobuf"

var (
	OK         = protobuf.Status_OK
	RECOVERING = protobuf.Status_Recovering
	FATAL      = protobuf.Status_Fatal
)
