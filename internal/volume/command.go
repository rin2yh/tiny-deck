package volume

// 末尾 \n はライン区切り文字であり定数に含めない。送信側が付与する。
const (
	CmdUp   = "vol:up"
	CmdDown = "vol:down"
	CmdMute = "vol:mute"

	PrefixCurrent = "vol:cur:"
	PrefixMuted   = "vol:muted:"
)
