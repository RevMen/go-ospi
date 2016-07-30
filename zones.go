package main

var zones = make([]Zone, 8)

// //var zoneDigits [8]bool
// //var program []Command

type Zone struct {
	Name    string `json:"name"`
	Address int    `json:"address"`
	Num     int    `json:"num"`
}

func initZones() {
	zones[0] = Zone{Name: "unused", Address: 0, Num: 1}
	zones[1] = Zone{Name: "North Side", Address: 1, Num: 2}
	zones[2] = Zone{Name: "blocked", Address: 2, Num: 3}
	zones[3] = Zone{Name: "Back Center", Address: 3, Num: 4}
	zones[4] = Zone{Name: "Back South", Address: 4, Num: 5}
	zones[5] = Zone{Name: "South Side", Address: 5, Num: 6}
	zones[6] = Zone{Name: "Front Center", Address: 6, Num: 7}
	zones[7] = Zone{Name: "unused", Address: 7, Num: 8}
}
