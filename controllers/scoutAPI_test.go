/*
 * Copyright (C) 2016 Clinton Freeman
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package controllers

import (
	"bytes"
	"database/sql"
	"github.com/MeasureTheFuture/mothership/configuration"
	"github.com/MeasureTheFuture/mothership/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	heartbeatJSON = `{
						"UUID":"59ef7180-f6b2-4129-99bf-970eb4312b4b",
						"Version":"0.1",
						"Health":{
							"IpAddress":"10.1.1.1",
							"CPU":0.4,
							"Memory":0.1,
							"TotalMemory":1233312.0,
							"Storage":0.1
						}
					}`

	db  *sql.DB
	err error
)

func buildPostRequest(fileName string, url string, uuid string, content string) (*http.Request, error) {
	body := bytes.Buffer{}
	w := multipart.NewWriter(&body)
	defer w.Close()

	part, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(part, content)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(echo.POST, url, &body)
	req.Header.Add("Mothership-Authorization", uuid)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return nil, err
	}

	return req, nil
}

func TestScoutAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ScoutAPI controller Suite")
}

var _ = Describe("ScoutAPI controller", func() {

	BeforeSuite(func() {
		config, err := configuration.Parse(os.Getenv("GOPATH") + "/mothership.json")
		Ω(err).Should(BeNil())
		db, err = sql.Open("postgres", "user="+config.DBUserName+" dbname="+config.DBTestName)
		Ω(err).Should(BeNil())
	})

	cleaner := func() {
		_, err := db.Exec(`DELETE FROM scout_logs`)
		Ω(err).Should(BeNil())

		_, err = db.Exec(`DELETE FROM scout_healths`)
		Ω(err).Should(BeNil())

		_, err = db.Exec(`DELETE FROM scouts`)
		Ω(err).Should(BeNil())
	}

	AfterEach(cleaner)
	AfterSuite(cleaner)

	Context("ScoutCalibration", func() {
		It("should drop the request if no authentication is supplied", func() {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/scout_api/calibrated", strings.NewReader(heartbeatJSON))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutCalibrated(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))

			i, err := models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))
		})

		It("should update the calibration frame iff the scout is authorised", func() {
			s := models.Scout{-1, "59ef7180-f6b2-4129-99bf-970eb4312b4b", "192.168.0.1", true, "foo", "calibrating"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			e := echo.New()

			body := bytes.Buffer{}
			w := multipart.NewWriter(&body)

			part, err := w.CreateFormFile("file", "calibrationFrame.jpg")
			Ω(err).Should(BeNil())

			src, err := os.Open("../testdata/calibrationFrame.jpg")
			Ω(err).Should(BeNil())

			_, err = io.Copy(part, src)
			Ω(err).Should(BeNil())
			src.Close()

			req, err := http.NewRequest(echo.POST, "scout_api/calibrated", &body)
			Ω(err).Should(BeNil())
			req.Header.Add("Mothership-Authorization", "59ef7180-f6b2-4129-99bf-970eb4312b4b")
			req.Header.Set("Content-Type", w.FormDataContentType())
			w.Close()
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutCalibrated(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusOK))

			s2, err := models.GetScoutByUUID(db, "59ef7180-f6b2-4129-99bf-970eb4312b4b")
			Ω(err).Should(BeNil())
			Ω(s2.State).Should(Equal(models.ScoutState("calibrated")))
			src, err = os.Open("../testdata/calibrationFrame.jpg")
			Ω(err).Should(BeNil())
			con, err := ioutil.ReadAll(src)
			Ω(err).Should(BeNil())
			frm, err := s2.GetCalibrationFrame(db)
			Ω(err).Should(BeNil())
			Ω(frm).Should(Equal(con))
		})
	})

	Context("ScoutHeartbeat", func() {
		It("should drop the request if no authentication is supplied", func() {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/scout_api/heartbeat", strings.NewReader(heartbeatJSON))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutHeartbeat(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))

			i, err := models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))
		})

		It("should create a new heartbeat iff the scout is authorised", func() {
			e := echo.New()
			req, err := buildPostRequest("scout.log", "/scout_api/heartbeat", "59ef7180-f6b2-4129-99bf-970eb4312b4b", heartbeatJSON)
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutHeartbeat(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))
			i, err := models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutHealths(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))

			err = ScoutHeartbeat(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))
			i, err = models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutHealths(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))

			s, err := models.GetScoutByUUID(db, "59ef7180-f6b2-4129-99bf-970eb4312b4b")
			Ω(err).Should(BeNil())
			s.Authorised = true
			err = s.Update(db)
			Ω(err).Should(BeNil())

			rec = httptest.NewRecorder()
			c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			err = ScoutHeartbeat(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusOK))
			i, err = models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutHealths(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))

			sh, err := models.GetLastScoutHealth(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(sh.CPU).Should(BeNumerically("==", float32(0.4)))
			Ω(sh.Memory).Should(BeNumerically("==", float32(0.1)))
			Ω(sh.TotalMemory).Should(BeNumerically("==", float32(1233312.0)))
			Ω(sh.Storage).Should(BeNumerically("==", float32(0.1)))
		})
	})

	Context("ScoutLog", func() {
		It("should drop the request if no authentication is supplied", func() {
			e := echo.New()
			req, err := http.NewRequest(echo.POST, "/scout_api/log", strings.NewReader(heartbeatJSON))
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutLog(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))

			i, err := models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))
		})

		It("should create a new log iff the scout is authorised", func() {
			e := echo.New()
			req, err := buildPostRequest("scout.log", "/scout_api/log", "59ef7180-f6b2-4129-99bf-970eb4312b4b", "log contents")
			Ω(err).Should(BeNil())
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			err = ScoutLog(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))
			i, err := models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutLogs(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))

			err = ScoutLog(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusNotFound))
			i, err = models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutLogs(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(0)))

			s, err := models.GetScoutByUUID(db, "59ef7180-f6b2-4129-99bf-970eb4312b4b")
			Ω(err).Should(BeNil())
			s.Authorised = true
			err = s.Update(db)
			Ω(err).Should(BeNil())

			rec = httptest.NewRecorder()
			c = e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			err = ScoutLog(db, c)
			Ω(err).Should(BeNil())
			Ω(rec.Code).Should(Equal(http.StatusOK))
			i, err = models.NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))
			i, err = models.NumScoutLogs(db)
			Ω(err).Should(BeNil())
			Ω(i).Should(Equal(int64(1)))

			sl, err := models.GetLastScoutLog(db, s.Id)
			Ω(err).Should(BeNil())
			log := string(sl.Log[:])
			Ω(log).Should(Equal("log contents"))
		})
	})
})
