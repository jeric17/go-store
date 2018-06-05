package store

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jinzhu/gorm"
	"spanda-admin/sdata"
)

type StoreConfig struct {
	Limit uint64
}

var Defaults = &StoreConfig{
	Limit: 10,
}

type Store struct {
	Table, Name, ForeignKey string
	Limit                   uint64
}

func (s Store) GetLimit() uint64 {
	if s.Limit == 0 {
		return Defaults.Limit
	}
	return s.Limit
}

func (s Store) List(o url.Values, m interface{}, c *uint64) {
	s.ListObject(o, m, c).Find(m)
}

func (s Store) StrToUint(str string) uint64 {
	i, _ := strconv.Atoi(str)
	return uint64(i)
}

func (s Store) PageOffset(o string) uint64 {
	fmt.Println("Limit", s.GetLimit())
	i, _ := strconv.Atoi(o)
	offset := (s.GetLimit() * uint64(i))

	return offset
}

func (s Store) ListObject(o url.Values, m interface{}, c *uint64) *gorm.DB {
	offset := s.PageOffset(o.Get("offset"))
	return s.Instance().Order(o.Get("sort")).Count(c).Limit(s.GetLimit()).Offset(offset)
}

func (s Store) Where(m interface{}, arg1 interface{}, arg2 ...interface{}) {
	s.Instance().Where(arg1, arg2...).Find(m)
}

func (s Store) Db(m interface{}) *gorm.DB {
	return sdata.Db.Model(m)
}

func (s Store) Create(m interface{}) {
	sdata.Db.Create(m)
}

func (s Store) Update(m interface{}) {
	sdata.Db.Save(m)
}

func (s Store) UpdateWhere(m interface{}, u interface{}, arg1 interface{}, arg2 interface{}) *gorm.DB {
	return sdata.Db.Model(m).Where(arg1, arg2).Update(u).Debug()
}

func (s Store) Instance() *gorm.DB {
	return sdata.Db.Table(s.Table)
}

func (s Store) ToLike(str string) string {
	var buffer bytes.Buffer
	buffer.WriteString("%")
	buffer.WriteString(str)
	buffer.WriteString("%")

	return buffer.String()
}

func (s Store) With(model Store, i *gorm.DB) *gorm.DB {
	return s.WithFk(model, i, "")
}

func (s Store) WithFk(model Store, i *gorm.DB, fk string) *gorm.DB {
	return s.BaseWith(model, i, fk, "left", s.Table, model.Name)
}

func (s Store) WithFkRight(model Store, i *gorm.DB, fk string) *gorm.DB {
	return s.BaseWith(model, i, fk, "right", s.Table, model.Name)
}

func (s Store) BaseWith(model Store, i *gorm.DB, fk string, jt string, tbl string, pl string) *gorm.DB {
	var join bytes.Buffer

	if fk == "" {
		fk = s.ForeignKey
	}

	join.WriteString(jt)
	join.WriteString(" join ")
	join.WriteString(model.Table)
	join.WriteString(" on ")
	join.WriteString(model.Table)
	join.WriteString(".")
	join.WriteString(model.ForeignKey)
	join.WriteString(" = ")
	join.WriteString(tbl)
	join.WriteString(".")
	join.WriteString(fk)

	if pl != "" {
		return i.Joins(join.String()).Preload(pl)
	}
	return i.Joins(join.String())
}

