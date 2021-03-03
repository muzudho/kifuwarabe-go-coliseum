package entities

// ColiseumConf - Tomlファイル。
type ColiseumConf struct {
	White Color
	Black Color
}

// Color - [White], [Black] 区画。
type Color struct {
	Connector string
	Workspace string
}
