package main

import (
	// "github.com/gin-gonic/contrib/static"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Device struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Fields       []string `json:"fields"`
	Secret       string   `json:"secret"`
	DisplayOrder int      `json:"display_order"`
}

func (d *Device) Validate() (bool, string) {
	if len(d.Name) == 0 || len(d.Name) > 32 {
		return false, "Invalid Name"
	}
	if len(d.Secret) > 128 {
		return false, "Invalid Secret"
	}
	return true, ""
}

func parseIntDefault(str string, defvalue int) int {
	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return defvalue
	}
	return int(v)
}

func initHttpd(db KVS) *gin.Engine {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"_status": 200, "message": "It works!"})
	})

	r.POST("/device", func(c *gin.Context) {
		newdev := &Device{
			c.PostForm("name"),
			c.PostForm("description"),
			strings.Split(c.PostForm("fields"), ","),
			c.PostForm("secret"),
			parseIntDefault(c.PostForm("display_order"), 0)}
		ok, msg := newdev.Validate()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"_status": 400, "message": msg})
			return
		}
		dev := &Device{}
		found, err := db.get("device", newdev.Name, dev)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"_status": 500, "message": "get error"})
			return
		}
		if found && dev.Secret != newdev.Secret && dev.Secret != c.PostForm("_secret") {
			c.JSON(http.StatusForbidden, gin.H{"_status": 403, "message": "invalid secret"})
			return
		}
		_, err = db.store("device", newdev.Name, newdev)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"_status": 500, "message": "store error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"_status": 201, "message": "created", "device": newdev})
	})

	r.GET("/stats", func(c *gin.Context) {
		offset := parseIntDefault(c.Query("offset"), 0)
		limit := parseIntDefault(c.Query("limit"), 100)

		var devNameFields []*struct {
			Name         string   `json:"name"`
			Fields       []string `json:"fields"`
			Description  string   `json:"description"`
			DisplayOrder int      `json:"display_order"`
		}
		_, _ = db.query("device", &devNameFields, "", "", offset, limit)
		c.JSON(http.StatusOK, gin.H{"_status": 200, "stats": devNameFields})
	})

	r.GET("/stats/:device/values", func(c *gin.Context) {
		dev := &Device{}
		found, _ := db.get("device", c.Param("device"), dev)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"_status": 404, "message": "not found"})
			return
		}
		offset := parseIntDefault(c.Query("offset"), 0)
		limit := parseIntDefault(c.Query("limit"), 100)
		values := []*map[string]float64{}
		_, _ = db.query("values:"+dev.Name, &values, "", "", offset, limit)
		c.JSON(http.StatusOK, gin.H{"_status": 200, "fields": dev.Fields, "description": dev.Description, "values": values})
	})

	r.GET("/stats/:device/values/latest", func(c *gin.Context) {
		dev := &Device{}
		found, _ := db.get("device", c.Param("device"), dev)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"_status": 404, "message": "not found"})
			return
		}
		values := []*map[string]float64{}
		_, _ = db.query("values:"+dev.Name, &values, "", "", 0, 1)
		c.JSON(http.StatusOK, gin.H{"_status": 200, "values": values[0]})
	})

	r.POST("/stats/:device/values", func(c *gin.Context) {
		dev := &Device{}
		found, _ := db.get("device", c.Param("device"), dev)
		if !found {
			c.JSON(http.StatusInternalServerError, gin.H{"_status": 404, "message": "not found"})
			return
		}
		if c.PostForm("_secret") != dev.Secret {
			c.JSON(http.StatusForbidden, gin.H{"_status": 403, "message": "invalid secret"})
			return
		}

		timestamp, _ := strconv.ParseInt(c.DefaultPostForm("_timestamp", "0"), 10, 64)
		if timestamp == 0 {
			timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		}
		stat := map[string]interface{}{"_timestamp": timestamp}
		for _, f := range dev.Fields {
			value, _ := strconv.ParseFloat(c.PostForm(f), 64)
			stat[f] = value
		}
		log.Printf("stat %d %v", timestamp, stat)
		db.store("values", c.Param("device")+":"+strconv.FormatInt(timestamp, 10), stat)
		c.JSON(http.StatusOK, gin.H{"_status": 201, "message": "created"})
	})

	r.DELETE("/stats/:device/values/:timestamp", func(c *gin.Context) {
		dev := &Device{}
		found, _ := db.get("device", c.Param("device"), dev)
		if !found {
			c.JSON(http.StatusInternalServerError, gin.H{"_status": 404, "message": "not found"})
			return
		}
		if c.Query("_secret") != dev.Secret {
			c.JSON(http.StatusInternalServerError, gin.H{"_status": 403, "message": "invalid secret"})
			return
		}
		_, _ = db.del("values", dev.Name+":"+c.Param("timestamp"))
		c.JSON(http.StatusOK, gin.H{"_status": 200, "message": "deleted"})
	})

	// Prometheus metrics
	r.GET("/metrics", func(c *gin.Context) {
		var devs []*struct {
			Name         string   `json:"name"`
			Fields       []string `json:"fields"`
			Description  string   `json:"description"`
			DisplayOrder int      `json:"display_order"`
		}
		_, _ = db.query("device", &devs, "", "", 0, 10000)
		res := ""
		for _, dev := range devs {
			values := []*map[string]float64{}
			_, _ = db.query("values:"+dev.Name, &values, "", "", 0, 1)
			for k, v := range *values[0] {
				if k == "_timestamp" {
					continue
				}
				res = res + "k_on_" + k + "{device=\"" + dev.Name + "\"} " + fmt.Sprint(v)
				if (*values[0])["_timestamp"] != 0.0 {
					res += " " + fmt.Sprint(int((*values[0])["_timestamp"]))
				}
				res = res + "\n"
			}
		}
		c.Data(http.StatusOK, "text/plain", []byte(res))
	})

	r.Static("/static", "./static")
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}

func main() {
	port := flag.Int("p", 8080, "http port")
	dbtype := flag.String("t", "leveldb", "datastore type")
	dbpath := flag.String("d", "data", "datastore uri(s) or path")
	flag.Parse()
	db, _ := NewKVS(*dbtype, *dbpath)
	defer db.Close()
	gin.SetMode(gin.ReleaseMode)
	log.Printf("start server. port: %d", *port)
	initHttpd(db).Run(":" + fmt.Sprint(*port))
}
