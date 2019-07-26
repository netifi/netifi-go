package framing

type FrameType uint16

const (
	FrameTypeUndefined            FrameType = 0x00
	FrameTypeBrokerSetup          FrameType = 0x01
	FrameTypeDestinationSetup     FrameType = 0x02
	FrameTypeGroup                FrameType = 0x03
	FrameTypeBroadcast            FrameType = 0x04
	FrameTypeShard                FrameType = 0x05
	FrameTypeAuthorizationWrapper FrameType = 0x06
)

func (f FrameType) String() string {
	switch f {
	case FrameTypeUndefined:
		return "UNDEFINED"
	case FrameTypeBrokerSetup:
		return "BROKER_SETUP"
	case FrameTypeDestinationSetup:
		return "DESTINATION_SETUP"
	case FrameTypeGroup:
		return "GROUP"
	case FrameTypeBroadcast:
		return "BROADCAST"
	case FrameTypeShard:
		return "SHARD"
	case FrameTypeAuthorizationWrapper:
		return "AUTHORIZATION_WRAPPER"
	default:
		panic("unknown frame type")
	}
}
