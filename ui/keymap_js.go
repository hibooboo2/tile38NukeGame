// +build js,wasm

package ui

var keyMap = map[int32]KeyCode{
	'"':  K_QUOTEDBL,
	'#':  K_HASH,
	'%':  K_PERCENT,
	'$':  K_DOLLAR,
	'&':  K_AMPERSAND,
	'\'': K_QUOTE,
	'(':  K_LEFTPAREN,
	')':  K_RIGHTPAREN,
	'*':  K_ASTERISK,
	'+':  K_PLUS,
	',':  K_COMMA,
	'-':  K_MINUS,
	'.':  K_PERIOD,
	' ':  K_SPACE,

	'0': K_0,
	'1': K_1,
	'2': K_2,
	'3': K_3,
	'4': K_4,
	'5': K_5,
	'6': K_6,
	'7': K_7,
	'8': K_8,
	'9': K_9,

	':':  K_COLON,
	';':  K_SEMICOLON,
	'<':  K_LESS,
	'=':  K_EQUALS,
	'>':  K_GREATER,
	'?':  K_QUESTION,
	'@':  K_AT,
	'[':  K_LEFTBRACKET,
	'\\': K_BACKSLASH,
	']':  K_RIGHTBRACKET,
	'^':  K_CARET,
	'_':  K_UNDERSCORE,
	'`':  K_BACKQUOTE,

	'a': K_a,
	'b': K_b,
	'c': K_c,
	'd': K_d,
	'e': K_e,
	'f': K_f,
	'g': K_g,
	'h': K_h,
	'i': K_i,
	'j': K_j,
	'k': K_k,
	'l': K_l,
	'm': K_m,
	'n': K_n,
	'o': K_o,
	'p': K_p,
	'q': K_q,
	'r': K_r,
	's': K_s,
	't': K_t,
	'u': K_u,
	'v': K_v,
	'w': K_w,
	'x': K_x,
	'y': K_y,
	'z': K_z,
}
