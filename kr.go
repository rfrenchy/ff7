package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
)

// https://pkg.go.dev/encoding/binary#Read

// hex = 4 bits, 4 bits
// 1 byte is 8 bits
// 2 hex symbols is 1 byte?
// 1 byte is 255 values

// F (hex)  = 15 (dec)
// 10 (hex) = 16 (dec)

// https://ff7-mods.github.io/ff7-flat-wiki/FF7/Text_encoding.html
var fftext = map[int](string){
	// offset 00
	0x00: " ", 0x01: "!", 0x02: `"`, 0x03: "#", 0x04: "$", 0x05: "%",
	0x06: "&", 0x07: "'", 0x08: "(", 0x09: ")", 0x0A: ":", 0x0B: "+",
	0x0C: ",", 0x0D: "-", 0x0E: ".", 0x0F: "/",
	// offset 10
	0x10: "0", 0x11: "1", 0x12: "2", 0x13: "3", 0x14: "4", 0x15: "5",
	0x16: "6", 0x17: "7", 0x18: "8", 0x19: "9", 0x1A: ":", 0x1B: ";",
	0x1C: "<", 0x1D: "=", 0x1E: ">", 0x1F: "?",
	// offset 20
	0x20: "@", 0x21: "A", 0x22: "B", 0x23: "C", 0x24: "D", 0x25: "E",
	0x26: "F", 0x27: "G", 0x28: "H", 0x29: "I", 0x2A: "J", 0x2B: "K",
	0x2C: "L", 0x2D: "M", 0x2E: "N", 0x2F: "O",
	// offset 30
	0x30: "P", 0x31: "Q", 0x32: "R", 0x33: "S", 0x34: "T", 0x35: "U",
	0x36: "V", 0x37: "W", 0x38: "X", 0x39: "Y", 0x3A: "Z", 0x3B: "[",
	0x3C: `\`, 0x3D: "]", 0x3E: "^", 0x3F: "_",
	// offset 40
	0x40: "`", 0x41: "a", 0x42: "b", 0x43: "c", 0x44: "d", 0x45: "e",
	0x46: "f", 0x47: "g", 0x48: "h", 0x49: "i", 0x4A: "j", 0x4B: "k",
	0x4C: "l", 0x4D: "m", 0x4E: "n", 0x4F: "o",
	// offset 50
	0x50: "p", 0x51: "q", 0x52: "r", 0x53: "s", 0x54: "t", 0x55: "u",
	0x56: "v", 0x57: "w", 0x58: "x", 0x59: "y", 0x5A: "z", 0x5B: "{",
	0x5C: "|", 0x5D: "}", 0x5E: "~", 0x5F: "",
	// offset 60
	0x60: "Ä", 0x61: "Á", 0x62: "Ç", 0x63: "É", 0x64: "Ñ", 0x65: "Ö",
	0x66: "Ü", 0x67: "á", 0x68: "à", 0x69: "â", 0x6A: "ä", 0x6B: "ã",
	0x6C: "å", 0x6D: "ç", 0x6E: "é", 0x6F: "è",
	// offset 70
	0x70: "ê", 0x71: "ë", 0x72: "í", 0x73: "ì", 0x74: "î", 0x75: "ï",
	0x76: "ñ", 0x77: "ó", 0x78: "ò", 0x79: "ô", 0x7A: "ö", 0x7B: "õ",
	0x7C: "ú", 0x7D: "ù", 0x7E: "û", 0x7F: "ü",
	// offset 80
	0x80: "⌘", 0x81: "°", 0x82: "¢", 0x83: "£", 0x84: "Ù", 0x85: "Û",
	0x86: "¶", 0x87: "ß", 0x88: "®", 0x89: "©", 0x8A: "™", 0x8B: "´",
	0x8C: "¨", 0x8D: "≠", 0x8E: "Æ", 0x8F: "Ø",
	// offset 90
	0x90: "∞", 0x91: "±", 0x92: "≤", 0x93: "≥", 0x94: "¥", 0x95: "µ",
	0x96: "∂", 0x97: "Σ", 0x98: "Π", 0x99: "π", 0x9A: "⌡", 0x9B: "ª",
	0x9C: "º", 0x9D: "Ω", 0x9E: "æ", 0x9F: "ø",
	// offset A0
	0xA0: "¿", 0xA1: "¡", 0xA2: "¬", 0xA3: "√", 0xA4: "ƒ", 0xA5: "≈",
	0xA6: "∆", 0xA7: "«", 0xA8: "»", 0xA9: "…", 0xAA: "{NOTHING}", 0xAB: "À",
	0xAC: "Ã", 0xAD: "Õ", 0xAE: "Œ", 0xAF: "œ",
	// offset B0
	0xB0: "–", 0xB1: "—", 0xB2: "“", 0xB3: "”", 0xB4: "‘", 0xB5: "’",
	0xB6: "÷", 0xB7: "◊", 0xB8: "ÿ", 0xB9: "Ÿ", 0xBA: "⁄", 0xBB: "¤",
	0xBC: "‹", 0xBD: "›", 0xBE: "ﬁ", 0xBF: "ﬂ",
	// offset C0
	0xC0: "■", 0xC1: "▪", 0xC2: "‚", 0xC3: "„", 0xC4: "‰", 0xC5: "Â",
	0xC6: "Ê", 0xC7: "Ë", 0xC8: "Á", 0xC9: "È", 0xCA: "í", 0xCB: "î",
	0xCC: "ï", 0xCD: "ì", 0xCE: "Ó", 0xCF: "Ô",
	// offset D0
	0xD0: " ", 0xD1: "Ò", 0xD2: "Ù", 0xD3: "Û", 0xD4: "", 0xD5: "",
	0xD6: "", 0xD7: "", 0xD8: "", 0xD9: "", 0xDA: "", 0xDB: "",
	0xDC: "", 0xDD: "", 0xDE: "", 0xDF: "",
	// offset E0
	0xE0: "{Choice}", 0xE1: "{Tab}", 0xE2: ",", 0xE3: ".”", 0xE4: "…”", 0xE5: "",
	0xE6: "", 0xE7: "{EOL}", 0xE8: "{New Scr}", 0xE9: "{New Scr?}", 0xEA: "{Cloud}", 0xEB: "{Barret}",
	0xEC: "{Tifa}", 0xED: "{Aerith}", 0xEE: "{Red XIII}", 0xEF: "{Yuffie}",
	// offset F0
	0xF0: "{Cait Sith}", 0xF1: "{Vincent}", 0xF2: "{Cid}", 0xF3: "{Party #1}", 0xF4: "{Party #2}", 0xF5: "{Party #3}",
	0xF6: "〇", 0xF7: "△", 0xF8: "☐", 0xF9: "✕", 0xFA: "", 0xFB: "",
	0xFC: "", 0xFD: "", 0xFE: "{FUNC}", 0xFF: "{END}",
}

type CharacterRecords struct {
	Cloud    CharacterRecord
	Barret   CharacterRecord
	Tifa     CharacterRecord
	Aeris    CharacterRecord
	RedXIII  CharacterRecord
	Yuffie   CharacterRecord
	CaitSith CharacterRecord
	Vincent  CharacterRecord
	Cid      CharacterRecord
}

func (cr *CharacterRecords) Get(character string) *CharacterRecord {
	switch character {
	case "Cloud":
		return &cr.Cloud
	case "Barret":
		return &cr.Barret
	case "Tifa":
		return &cr.Tifa
	case "Aeris":
		return &cr.Aeris
	case "RedXIII":
		return &cr.RedXIII
	case "Yuffie":
		return &cr.Yuffie
	case "CaitSith":
		return &cr.CaitSith
	case "Vincent":
		return &cr.Vincent
	case "Cid":
		return &cr.Cid
	}

	return nil
}

// 56 bytes
type CharacterRecord struct {
	StrengthCurve   uint8
	VitalityCurve   uint8
	MagicCurve      uint8
	SpiritCurve     uint8
	DexterityCurve  uint8
	LuckCurve       uint8
	HPCurve         uint8
	MPCurve         uint8
	EXPCurve        uint8
	FF1_            uint8
	StartLVL        uint8
	FF2_            uint8
	Limit1x1        uint8
	Limit1x2        uint8
	Limit1x3_       uint8
	Limit2x1        uint8
	Limit2x2        uint8
	Limit2x3_       uint8
	Limit3x1        uint8
	Limit3x2        uint8
	Limit3x3_       uint8
	Limit4x1        uint8
	Limit4x2        uint8
	Limit4x3_       uint8
	KillsLimit2     uint16
	KillsLimit3     uint16
	UsesLimit1x2    uint16
	UsesLimit1x3_   uint16
	UsesLimit2x2    uint16
	UsesLimit2x3_   uint16
	UsesLimit3x2    uint16
	UsesLimit3x3_   uint16
	HPDivisorLimit1 uint32
	HPDivisorLimit2 uint32
	HPDivisorLimit3 uint32
	HPDivisorLimit4 uint32
}

func (cr *CharacterRecord) PrimaryCurve(stat string) int {
	var x uint8
	x = 0

	switch stat {
	case "Strength":
		x = cr.StrengthCurve
	case "Vitality":
		x = cr.VitalityCurve
	case "Magic":
		x = cr.MagicCurve
	case "Spirit":
		x = cr.SpiritCurve
	case "Dexterity":
		x = cr.DexterityCurve
	}

	return int(x)
}

type StatRandomBonus struct {
	Rb0  uint8
	Rb1  uint8
	Rb2  uint8
	Rb3  uint8
	Rb4  uint8
	Rb5  uint8
	Rb6  uint8
	Rb7  uint8
	Rb8  uint8
	Rb9  uint8
	Rb10 uint8
	Rb11 uint8
}

func (s *StatRandomBonus) Get(index int) uint8 {
	switch index {
	case 0:
		return s.Rb0
	case 1:
		return s.Rb1
	case 2:
		return s.Rb2
	case 3:
		return s.Rb3
	case 4:
		return s.Rb4
	case 5:
		return s.Rb5
	case 6:
		return s.Rb6
	case 7:
		return s.Rb7
	case 8:
		return s.Rb8
	case 9:
		return s.Rb9
	case 10:
		return s.Rb10
	case 11:
		return s.Rb11
	}

	// replace with err /panic?
	return 0
}

// Percentage Value
type HPRandomBonus struct {
	Rb1  uint8
	Rb2  uint8
	Rb3  uint8
	Rb4  uint8
	Rb5  uint8
	Rb6  uint8
	Rb7  uint8
	Rb8  uint8
	Rb9  uint8
	Rb10 uint8
	Rb11 uint8
	Rb12 uint8
}

// Percentage Value
type MPRandomBonus struct {
	Rb1  uint8
	Rb2  uint8
	Rb3  uint8
	Rb4  uint8
	Rb5  uint8
	Rb6  uint8
	Rb7  uint8
	Rb8  uint8
	Rb9  uint8
	Rb10 uint8
	Rb11 uint8
	Rb12 uint8
}

// 16 bytes
type StatCurveRecord struct {
	L2_11_Gradient  int8
	L2_11_Base      int8
	L12_21_Gradient int8
	L12_21_Base     int8
	L22_31_Gradient int8
	L22_31_Base     int8
	L32_41_Gradient int8
	L32_41_Base     int8
	L42_51_Gradient int8
	L42_51_Base     int8
	L52_61_Gradient int8
	L52_61_Base     int8
	L62_81_Gradient int8
	L62_81_Base     int8
	L82_99_Gradient int8
	L82_99_Base     int8
}

func (s *StatCurveRecord) Get(level int) (gradient int, base int) {
	if level <= 11 {
		return int(s.L2_11_Gradient), int(s.L2_11_Base)
	} else if level <= 21 {
		return int(s.L12_21_Gradient), int(s.L12_21_Base)
	} else if level <= 31 {
		return int(s.L22_31_Gradient), int(s.L22_31_Base)
	} else if level <= 41 {
		return int(s.L32_41_Gradient), int(s.L32_41_Base)
	} else if level <= 51 {
		return int(s.L42_51_Gradient), int(s.L42_51_Base)
	} else if level <= 61 {
		return int(s.L52_61_Gradient), int(s.L52_61_Base)
	} else if level <= 81 {
		return int(s.L62_81_Gradient), int(s.L62_81_Base)
	} else {
		return int(s.L82_99_Gradient), int(s.L82_99_Base)
	}
}

// 37 different stat curves
type StatCurveRecords struct {
	SCR_1  StatCurveRecord
	SCR_2  StatCurveRecord
	SCR_3  StatCurveRecord
	SCR_4  StatCurveRecord
	SCR_5  StatCurveRecord
	SCR_6  StatCurveRecord
	SCR_7  StatCurveRecord
	SCR_8  StatCurveRecord
	SCR_9  StatCurveRecord
	SCR_10 StatCurveRecord
	SCR_11 StatCurveRecord
	SCR_12 StatCurveRecord
	SCR_13 StatCurveRecord
	SCR_14 StatCurveRecord
	SCR_15 StatCurveRecord
	SCR_16 StatCurveRecord
	SCR_17 StatCurveRecord
	SCR_18 StatCurveRecord
	SCR_19 StatCurveRecord
	SCR_20 StatCurveRecord
	SCR_21 StatCurveRecord
	SCR_22 StatCurveRecord
	SCR_23 StatCurveRecord
	SCR_24 StatCurveRecord
	SCR_25 StatCurveRecord
	SCR_26 StatCurveRecord
	SCR_27 StatCurveRecord
	SCR_28 StatCurveRecord
	SCR_29 StatCurveRecord
	SCR_30 StatCurveRecord
	SCR_31 StatCurveRecord
	SCR_32 StatCurveRecord
	SCR_33 StatCurveRecord
	SCR_34 StatCurveRecord
	SCR_35 StatCurveRecord
	SCR_36 StatCurveRecord
	SCR_37 StatCurveRecord
}

func main() {
	Kernel3()
	// Kernel4()
}

func (s *StatCurveRecords) Get(index int) *StatCurveRecord {
	switch index {
	case 0:
		return &s.SCR_1
	case 1:
		return &s.SCR_2
	case 2:
		return &s.SCR_3
	case 3:
		return &s.SCR_4
	case 4:
		return &s.SCR_5
	case 5:
		return &s.SCR_6
	case 6:
		return &s.SCR_7
	case 7:
		return &s.SCR_8
	case 8:
		return &s.SCR_9
	case 9:
		return &s.SCR_10
	case 10:
		return &s.SCR_11
	case 11:
		return &s.SCR_12
	case 12:
		return &s.SCR_13
	case 13:
		return &s.SCR_14
	case 14:
		return &s.SCR_15
	case 15:
		return &s.SCR_16
	case 16:
		return &s.SCR_17
	case 17:
		return &s.SCR_18
	case 18:
		return &s.SCR_19
	case 19:
		return &s.SCR_20
	case 20:
		return &s.SCR_21
	case 21:
		return &s.SCR_22
	case 22:
		return &s.SCR_23
	case 23:
		return &s.SCR_24
	case 24:
		return &s.SCR_25
	case 25:
		return &s.SCR_26
	case 26:
		return &s.SCR_27
	case 27:
		return &s.SCR_28
	case 28:
		return &s.SCR_29
	case 29:
		return &s.SCR_30
	case 30:
		return &s.SCR_31
	case 31:
		return &s.SCR_32
	case 32:
		return &s.SCR_33
	case 33:
		return &s.SCR_34
	case 34:
		return &s.SCR_35
	case 35:
		return &s.SCR_36
	case 36:
		return &s.SCR_37
	}
	return nil
}

type K3 struct {
	Records         CharacterRecords
	StatRandomBonus StatRandomBonus
	HPRB            HPRandomBonus
	MPRB            MPRandomBonus
	SCRecords       StatCurveRecords
}

// KERNEL.BIN2 i.e. 3rd Kernel bin
// Battle and Growth Data
func Kernel3() {
	k, err := os.Open("KERNEL.bin2")
	if err != nil {
		panic(err)
	}
	defer k.Close()

	var k3 K3
	err = binary.Read(k, binary.LittleEndian, &k3)
	if err != nil {
		panic(err)
	}

	//fmt.Println("*** Cloud")
	// k3.Records.Tifa.Print()
	//fmt.Println("RSB1", k3.StatRB.Rb1)
	//fmt.Println("RSB2", k3.StatRB.Rb2)
	//fmt.Println("RSB5", k3.StatRB.Rb5)
	//fmt.Println("RSB8", k3.StatRB.Rb8)
	//fmt.Println("RSB11", k3.StatRB.Rb11)
	// fmt.Println("Start LVL:", k3.Records.Cloud.StartLVL)
	//fmt.Println("*** Barret")
	//k3.Records.Barret.Print()

	// Cloud
	stats_cloud_normal := struct {
		Strength  int
		Vitality  int
		Magic     int
		Spirit    int
		Dexterity int
	}{20, 16, 19, 17, 6}

	fmt.Println("level strength vitality magic spirit dexterity")
	fmt.Println(1, stats_cloud_normal.Strength, stats_cloud_normal.Vitality, stats_cloud_normal.Magic, stats_cloud_normal.Spirit, stats_cloud_normal.Dexterity)

	// skip first level
	//for i := 1; i <= 1; i++ {
	for lvl := 2; lvl <= 99; lvl++ {
		// get curve or characters stat
		str := k3.Records.Get("Cloud").PrimaryCurve("Strength")
		vit := k3.Records.Get("Cloud").PrimaryCurve("Vitality")
		mag := k3.Records.Get("Cloud").PrimaryCurve("Magic")
		spr := k3.Records.Get("Cloud").PrimaryCurve("Spirit")
		dex := k3.Records.Get("Cloud").PrimaryCurve("Dexterity")

		// get gradient and base for curve for given level bracket
		g_str, b_str := k3.SCRecords.Get(str).Get(lvl)
		g_vit, b_vit := k3.SCRecords.Get(vit).Get(lvl)
		g_mag, b_mag := k3.SCRecords.Get(mag).Get(lvl)
		g_spr, b_spr := k3.SCRecords.Get(spr).Get(lvl)
		g_dex, b_dex := k3.SCRecords.Get(dex).Get(lvl)

		// Baseline Stat = b + (g * lvl /100)
		bs_str := b_str + (g_str * lvl / 100)
		bs_vit := b_vit + (g_vit * lvl / 100)
		bs_mag := b_mag + (g_mag * lvl / 100)
		bs_spr := b_spr + (g_spr * lvl / 100)
		bs_dex := b_dex + (g_dex * lvl / 100)

		// Get different between new baseline stat and current stat
		stat_difference_str := (rand.Intn(8) + bs_str - stats_cloud_normal.Strength) % 11
		stat_difference_vit := (rand.Intn(8) + bs_vit - stats_cloud_normal.Vitality) % 11
		stat_difference_mag := (rand.Intn(8) + bs_mag - stats_cloud_normal.Magic) % 11
		stat_difference_spr := (rand.Intn(8) + bs_spr - stats_cloud_normal.Spirit) % 11
		stat_difference_dex := (rand.Intn(8) + bs_dex - stats_cloud_normal.Dexterity) % 11

		// Get value to increase stat by
		stats_cloud_normal.Strength += int(k3.StatRandomBonus.Get(int(math.Abs(float64(stat_difference_str)))))
		stats_cloud_normal.Vitality += int(k3.StatRandomBonus.Get(int(math.Abs(float64(stat_difference_vit)))))
		stats_cloud_normal.Magic += int(k3.StatRandomBonus.Get(int(math.Abs(float64(stat_difference_mag)))))
		stats_cloud_normal.Spirit += int(k3.StatRandomBonus.Get(int(math.Abs(float64(stat_difference_spr)))))
		stats_cloud_normal.Dexterity += int(k3.StatRandomBonus.Get(int(math.Abs(float64(stat_difference_dex)))))

		// Print
		fmt.Println(lvl, stats_cloud_normal.Strength, stats_cloud_normal.Vitality, stats_cloud_normal.Magic, stats_cloud_normal.Spirit, stats_cloud_normal.Dexterity)
	}
}

// Stat per level
func (k3 *K3) Stat(character string, stat string, lvl int) float64 {
	// get gradient and base

	// how to get current stat

	//	k3.Records
	//		.Get(character).Get(stat)

	//	g, b := k3.SCRecords.
	//		Get(k3.Records.Get(character).PrimaryCurve(stat)).
	//		Get(lvl)

	// change g and b to be gotten based off current stat value?

	// import
	// math/rand
	r := rand.Intn(8)
	// + rand?
	//r := 0

	// Baseline Stat = Base + [Gradient * Level / 100]
	// Stat Difference = Rnd(1..8) + Baseline - Current Stat (capped at 11?)

	return float64(r)
}

func Kernel4() {
	k, err := os.Open("KERNEL.bin4")
	if err != nil {
		panic(err)
	}
	defer k.Close()

	x := make([]byte, 1)
	i := 0
	for {
		_, err := k.Read(x)
		if err != nil {
			break
		}

		fmt.Println(fmt.Sprintf("%x %s", i, fftext[int(x[0])]))
		i++
	}
}

func (cr *CharacterRecord) Print() {
	fmt.Println("Strength:", cr.StrengthCurve)
	fmt.Println("Vitality:", cr.VitalityCurve)
	fmt.Println("Magic:", cr.MagicCurve)
	fmt.Println("Spirit:", cr.SpiritCurve)
	fmt.Println("Dexterity:", cr.DexterityCurve)
	fmt.Println("Luck:", cr.LuckCurve)
	fmt.Println("HP:", cr.HPCurve)
	fmt.Println("MP:", cr.MPCurve)
	fmt.Println("EXP:", cr.EXPCurve)
}

func Testtext() {
	fmt.Println(fftext[0x22] + fftext[0x41] + fftext[0x52] + fftext[0x52] +
		fftext[0x45] + fftext[0x54])

	fmt.Println(fftext[0x34] + fftext[0x49] + fftext[0x46] + fftext[0x41])
	fmt.Println(fftext[0x21] + fftext[0x45] + fftext[0x52] + fftext[0x49] +
		fftext[0x54] + fftext[0x48])
}
