package constants

// source: http://mpgedit.org/mpgedit/mpeg_format/mpeghdr.htm
// Define the map with string keys and int values
var BitrateMap = map[byte]map[string]int{
	0b0000: {
		"V1,L1":       0,
		"V1,L2":       0,
		"V1,L3":       0,
		"V2,L1":       0,
		"V2, L2 & L3": 0,
	},
	0b0001: {
		"V1,L1":       32,
		"V1,L2":       32,
		"V1,L3":       32,
		"V2,L1":       32,
		"V2, L2 & L3": 8,
	},
	0b0010: {
		"V1,L1":       64,
		"V1,L2":       48,
		"V1,L3":       40,
		"V2,L1":       48,
		"V2, L2 & L3": 16,
	},
	0b0011: {
		"V1,L1":       96,
		"V1,L2":       56,
		"V1,L3":       48,
		"V2,L1":       56,
		"V2, L2 & L3": 24,
	},
	0b0100: {
		"V1,L1":       128,
		"V1,L2":       64,
		"V1,L3":       56,
		"V2,L1":       64,
		"V2, L2 & L3": 32,
	},
	0b0101: {
		"V1,L1":       160,
		"V1,L2":       80,
		"V1,L3":       64,
		"V2,L1":       80,
		"V2, L2 & L3": 40,
	},
	0b0110: {
		"V1,L1":       192,
		"V1,L2":       96,
		"V1,L3":       80,
		"V2,L1":       96,
		"V2, L2 & L3": 48,
	},
	0b0111: {
		"V1,L1":       224,
		"V1,L2":       112,
		"V1,L3":       96,
		"V2,L1":       112,
		"V2, L2 & L3": 56,
	},
	0b1000: {
		"V1,L1":       256,
		"V1,L2":       128,
		"V1,L3":       112,
		"V2,L1":       128,
		"V2, L2 & L3": 64,
	},
	0b1001: {
		"V1,L1":       288,
		"V1,L2":       160,
		"V1,L3":       128,
		"V2,L1":       144,
		"V2, L2 & L3": 80,
	},
	0b1010: {
		"V1,L1":       320,
		"V1,L2":       192,
		"V1,L3":       160,
		"V2,L1":       160,
		"V2, L2 & L3": 96,
	},
	0b1011: {
		"V1,L1":       352,
		"V1,L2":       224,
		"V1,L3":       192,
		"V2,L1":       176,
		"V2, L2 & L3": 112,
	},
	0b1100: {
		"V1,L1":       384,
		"V1,L2":       256,
		"V1,L3":       224,
		"V2,L1":       192,
		"V2, L2 & L3": 128,
	},
	0b1101: {
		"V1,L1":       416,
		"V1,L2":       320,
		"V1,L3":       256,
		"V2,L1":       224,
		"V2, L2 & L3": 144,
	},
	0b1110: {
		"V1,L1":       448,
		"V1,L2":       384,
		"V1,L3":       320,
		"V2,L1":       256,
		"V2, L2 & L3": 160,
	},
}

// Define the 2D map with string keys and int values
var SamplingRateMap = map[byte]map[string]int{
	0b00: {
		"MPEG1":   44100,
		"MPEG2":   22050,
		"MPEG2.5": 11025,
	},
	0b01: {
		"MPEG1":   48000,
		"MPEG2":   24000,
		"MPEG2.5": 12000,
	},
	0b10: {
		"MPEG1":   32000,
		"MPEG2":   16000,
		"MPEG2.5": 8000,
	},
	0b11: {
		"MPEG1":   0, // Assuming "reserv." means 0
		"MPEG2":   0,
		"MPEG2.5": 0,
	},
}

var VersionMap = map[int]string{
	// 00
	0: "V2",
	// 01
	1: "reserved",
	// 10
	2: "V2",
	// 11
	3: "V1",
}

var LayerMap = map[int]string{
	// 00
	0: "reserved",
	// 01
	1: "L3",
	// 10
	2: "L2",
	// 11
	3: "L1",
}
