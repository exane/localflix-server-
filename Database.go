package main

import (
  "github.com/jinzhu/gorm"
  "fmt"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "io/ioutil"
  "encoding/json"
)

var DB gorm.DB

type Config struct {
  Database   struct {
               Root     string
               Password string
               Db       string
             }
  Server     struct {
               Url  string
               Port string
             }
  Fileserver struct {
               Url            string
               Root_directory string
               Port           string
             }
  TMDb       struct {
               API_KEY string
             }
}

var config *Config

func load_config() *Config {
  if config != nil {
    return config
  }
  config = &Config{}
  js, _ := ioutil.ReadFile("./config.json")
  json.Unmarshal(js, config)
  return config
}

func init_db() {
  load_config()
  db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", config.Database.Root, config.Database.Password, config.Database.Db))
  DB = *db
  DB.LogMode(true)

  if err != nil {
    panic("failed to connect database")
  }
}

func create_tables() {
  DB.DropTableIfExists(Episode{}, Serie{}, Season{}, Subtitle{})
  DB.CreateTable(Episode{}, Serie{}, Season{}, Subtitle{})
  //DB.AutoMigrate(Episode{}, Serie{}, Season{})
}

func dump_import() {
  data := loadDump("./fetch/DATA_DUMP.json")

  for _, val := range data {
    DB.Create(&val)
  }
}

func loadDump(file string) []Serie {
  js, _ := ioutil.ReadFile(file)
  data := []Serie{}
  json.Unmarshal(js, &data)
  return data
}
