package util

const (
    ChannelSystem = 1 // 站内信
    ChannelSMS    = 2 // 短信
    ChannelEmail  = 3 // 邮件
    ChannelPush   = 4 // APP推送
)

func GetChannelString(channel int32) string {
    switch channel {
    case ChannelSystem:
        return "system"
    case ChannelSMS:
        return "sms"
    case ChannelEmail:
        return "email"
    case ChannelPush:
        return "push"
    default:
        return "system"
    }
}