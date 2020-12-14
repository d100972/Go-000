package model

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

func openDB(username, password, addr, name string) (*gorm.DB,error){
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		// "Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		err = errors.Wrapf(err, "Database connection failed. Database name: %s", name)
		return nil, err
	}

	// set for db connection
	setupDB(db)

	return db, nil
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	// db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// used for cli
func InitSelfDB() (*gorm.DB, error) {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

func GetSelfDB() (*gorm.DB, error) {
	return InitSelfDB()
}

func InitDockerDB() (*gorm.DB, error) {
	return openDB(viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"))
}

func GetDockerDB() (*gorm.DB, error) {
	return InitDockerDB()
}

func (db *Database) Init() error {
	selfDB, err := GetSelfDB()
	if err != nil {
		return err
	}

	dockerDB, err := GetDockerDB()
	if err != nil {
		return err
	}

	DB = &Database{
		Self:   selfDB,
		Docker: dockerDB,
	}
	return nil
}

func (db *Database) Close() error {
	err := DB.Self.Close(); if err != nil {
		err = errors.Wrap(err, "gorm Close() error ->")
		return err
	}

	err = DB.Docker.Close(); if err != nil {
		err = errors.Wrap(err, "gorm Close() error ->")
		return err
	}
	return nil
}
