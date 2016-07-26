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
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"mime/multipart"
	"mothership/models"
	"net/http"
	"net/http/httptest"
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

func TestScoutHeartbeat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ScoutAPI controller Suite")
}

var _ = Describe("ScoutAPI controller", func() {
	BeforeSuite(func() {
		db, err = sql.Open("postgres", "user=cfreeman dbname=mothership_test")
		Ω(err).Should(BeNil())
	})

	AfterEach(func() {
		_, err := db.Exec(`DELETE FROM scout_logs`)
		Ω(err).Should(BeNil())

		_, err = db.Exec(`DELETE FROM scout_healths`)
		Ω(err).Should(BeNil())

		_, err = db.Exec(`DELETE FROM scouts`)
		Ω(err).Should(BeNil())
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

		It("should create a new log if the scout is authorised", func() {
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
