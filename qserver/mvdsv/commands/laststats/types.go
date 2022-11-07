package laststats

type Entry struct {
	Version  int      `json:"version"`
	Date     string   `json:"date"`
	Map      string   `json:"map"`
	Hostname string   `json:"hostname"`
	Ip       string   `json:"ip"`
	Port     int      `json:"port"`
	Mode     string   `json:"mode"`
	Tl       int      `json:"tl"`
	Fl       int      `json:"fl"`
	Dm       int      `json:"dm"`
	Tp       int      `json:"tp"`
	Duration int      `json:"duration"`
	Demo     string   `json:"demo"`
	Teams    []string `json:"teams"`
	Players  []Player `json:"players"`
}

type Player struct {
	TopColor    int          `json:"top-color"`
	BottomColor int          `json:"bottom-color"`
	Ping        int          `json:"ping"`
	Login       string       `json:"login"`
	Name        string       `json:"name"`
	Team        string       `json:"team"`
	Stats       PlayerStats  `json:"stats"`
	Dmg         PlayerDamage `json:"dmg"`
	XferRL      int          `json:"xferRL"`
	XferLG      int          `json:"xferLG"`
	Spree       struct {
		Max  int `json:"max"`
		Quad int `json:"quad"`
	} `json:"spree"`
	Control float64 `json:"control"`
	Speed   struct {
		Max float64 `json:"max"`
		Avg float64 `json:"avg"`
	} `json:"speed"`
	Weapons struct {
		Axe struct {
			Acc struct {
				Attacks int `json:"attacks"`
				Hits    int `json:"hits"`
			} `json:"acc,omitempty"`
			Deaths  int `json:"deaths"`
			Pickups struct {
				Dropped    int `json:"dropped,omitempty"`
				TotalTaken int `json:"total-taken"`
			} `json:"pickups"`
		} `json:"axe"`
		Sg struct {
			Acc struct {
				Attacks int `json:"attacks"`
				Hits    int `json:"hits"`
			} `json:"acc"`
			Deaths int `json:"deaths"`
			Damage struct {
				Enemy int `json:"enemy"`
				Team  int `json:"team"`
			} `json:"damage"`
			Pickups struct {
				Dropped    int `json:"dropped,omitempty"`
				TotalTaken int `json:"total-taken,omitempty"`
			} `json:"pickups,omitempty"`
			Kills struct {
				Total int `json:"total"`
				Team  int `json:"team"`
				Enemy int `json:"enemy"`
				Self  int `json:"self"`
			} `json:"kills,omitempty"`
		} `json:"sg"`
		Ssg struct {
			Deaths  int `json:"deaths"`
			Pickups struct {
				Taken           int `json:"taken"`
				TotalTaken      int `json:"total-taken"`
				SpawnTaken      int `json:"spawn-taken"`
				SpawnTotalTaken int `json:"spawn-total-taken"`
			} `json:"pickups"`
			Kills struct {
				Total int `json:"total"`
				Team  int `json:"team"`
				Enemy int `json:"enemy"`
				Self  int `json:"self"`
			} `json:"kills,omitempty"`
		} `json:"ssg,omitempty"`
		Ng struct {
			Kills struct {
				Total int `json:"total"`
				Team  int `json:"team"`
				Enemy int `json:"enemy"`
				Self  int `json:"self"`
			} `json:"kills"`
			Deaths  int `json:"deaths"`
			Pickups struct {
				Taken           int `json:"taken"`
				TotalTaken      int `json:"total-taken"`
				SpawnTaken      int `json:"spawn-taken"`
				SpawnTotalTaken int `json:"spawn-total-taken"`
				Dropped         int `json:"dropped,omitempty"`
			} `json:"pickups"`
			Acc struct {
				Attacks int `json:"attacks"`
				Hits    int `json:"hits"`
			} `json:"acc,omitempty"`
			Damage struct {
				Enemy int `json:"enemy"`
				Team  int `json:"team"`
			} `json:"damage,omitempty"`
		} `json:"ng"`
		Rl struct {
			Kills struct {
				Total int `json:"total"`
				Team  int `json:"team"`
				Enemy int `json:"enemy"`
				Self  int `json:"self"`
			} `json:"kills"`
			Deaths  int `json:"deaths"`
			Pickups struct {
				Taken           int `json:"taken"`
				TotalTaken      int `json:"total-taken"`
				Dropped         int `json:"dropped,omitempty"`
				SpawnTaken      int `json:"spawn-taken,omitempty"`
				SpawnTotalTaken int `json:"spawn-total-taken,omitempty"`
			} `json:"pickups"`
			Acc struct {
				Attacks int `json:"attacks"`
				Hits    int `json:"hits"`
				Real    int `json:"real"`
				Virtual int `json:"virtual"`
			} `json:"acc,omitempty"`
			Damage struct {
				Enemy int `json:"enemy"`
				Team  int `json:"team"`
			} `json:"damage,omitempty"`
		} `json:"rl"`
		Lg struct {
			Deaths  int `json:"deaths"`
			Pickups struct {
				Taken           int `json:"taken"`
				TotalTaken      int `json:"total-taken"`
				SpawnTaken      int `json:"spawn-taken"`
				SpawnTotalTaken int `json:"spawn-total-taken"`
				Dropped         int `json:"dropped,omitempty"`
			} `json:"pickups"`
			Acc struct {
				Attacks int `json:"attacks"`
				Hits    int `json:"hits"`
			} `json:"acc,omitempty"`
			Kills struct {
				Total int `json:"total"`
				Team  int `json:"team"`
				Enemy int `json:"enemy"`
				Self  int `json:"self"`
			} `json:"kills,omitempty"`
			Damage struct {
				Enemy int `json:"enemy"`
				Team  int `json:"team"`
			} `json:"damage,omitempty"`
		} `json:"lg"`
		Gl struct {
			Deaths  int `json:"deaths"`
			Pickups struct {
				Taken           int `json:"taken"`
				TotalTaken      int `json:"total-taken"`
				SpawnTaken      int `json:"spawn-taken"`
				SpawnTotalTaken int `json:"spawn-total-taken"`
			} `json:"pickups"`
		} `json:"gl,omitempty"`
	} `json:"weapons"`
	Items PlayerItems `json:"items"`
	Bot   struct {
		Skill      int  `json:"skill"`
		Customised bool `json:"customised"`
	} `json:"bot,omitempty"`
}

type PlayerStats struct {
	Frags      int `json:"frags"`
	Deaths     int `json:"deaths"`
	Tk         int `json:"tk"`
	SpawnFrags int `json:"spawn-frags"`
	Kills      int `json:"kills"`
	Suicides   int `json:"suicides"`
}

type PlayerDamage struct {
	Taken        int `json:"taken"`
	Given        int `json:"given"`
	Team         int `json:"team"`
	Self         int `json:"self"`
	TeamWeapons  int `json:"team-weapons"`
	EnemyWeapons int `json:"enemy-weapons"`
	TakenToDie   int `json:"taken-to-die"`
}

type PlayerItems struct {
	Ga struct {
		Took int `json:"took"`
		Time int `json:"time"`
	} `json:"ga"`
	Ra struct {
		Took int `json:"took"`
		Time int `json:"time"`
	} `json:"ra"`
	Health100 struct {
		Took int `json:"took"`
	} `json:"health_100,omitempty"`
	Ya struct {
		Took int `json:"took"`
		Time int `json:"time"`
	} `json:"ya,omitempty"`
	Health15 struct {
		Took int `json:"took"`
	} `json:"health_15,omitempty"`
	Health25 struct {
		Took int `json:"took"`
	} `json:"health_25,omitempty"`
}
