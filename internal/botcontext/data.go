package botcontext

import (
	"slices"
)

type Data struct {
	Admins      []int64  `yaml:"admins"`
	SystemUser  string   `yaml:"systemUser"`
	NamesToKill []string `yaml:"namesToKill"`
}

func (d *Data) AddAdmin(id int64) {
	d.Admins = append(d.Admins, id)
}

func (d *Data) IsAdmin(id int64) bool {
	return slices.Contains(d.Admins, id)
}

func (d *Data) RemoveNames(args []string) {
	for _, arg := range args {
		index := slices.Index(d.NamesToKill, arg)
		if index >= 0 {
			d.NamesToKill = slices.Delete(d.NamesToKill, index, index+1)
		}
	}
}

func (d *Data) AddNames(args []string) {
	d.NamesToKill = append(d.NamesToKill, args...)

}

func (d *Data) ClearSystemUser() {
	d.SetSystemUser("")
}

func (d *Data) SetSystemUser(user string) {
	d.SystemUser = user
}
