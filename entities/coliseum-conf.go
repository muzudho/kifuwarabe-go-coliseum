package entities

// ColiseumConf - Tomlファイル。
type ColiseumConf struct {
	Coliseum Coliseum
}

// Coliseum - [Coliseum] 区画。
type Coliseum struct {
	WhiteWorkspace string
	BlackWorkspace string
}
