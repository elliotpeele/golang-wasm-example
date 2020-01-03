// Copyright (c) 2019 Elliot Peele <elliot@bentlogic.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/elliotpeele/golang-wasm-example/sampledata/models"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "datagen",
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := cmd.PersistentFlags().GetInt("users")
		if err != nil {
			return err
		}
		projects, err := cmd.PersistentFlags().GetInt("projects")
		if err != nil {
			return err
		}
		if err := generateUsers(users); err != nil {
			return err
		}
		return generateProjects(projects)
	},
}

func randomTime() time.Time {
	return time.Date(
		randomdata.Number(1995, 2019),
		time.Month(randomdata.Number(1, 12)),
		randomdata.Number(1, 30),
		randomdata.Number(0, 23),
		randomdata.Number(0, 59),
		randomdata.Number(0, 59),
		0,
		time.UTC,
	)
}

func generateUsers(count int) error {
	var users []models.User
	for i := 0; i < count; i++ {
		users = append(users, models.User{
			ID:        uuid.NewV4().String(),
			FirstName: randomdata.FirstName(randomdata.RandomGender),
			LastName:  randomdata.LastName(),
			Email:     randomdata.Email(),
			UpdatedAt: randomTime(),
		})
	}
	f, err := os.Create("generated/users.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(users)
}

func generateProjects(count int) error {
	var projects []models.Project
	for i := 0; i < count; i++ {
		projects = append(projects, models.Project{
			ID:        uuid.NewV4().String(),
			Name:      randomdata.SillyName(),
			UpdatedAt: randomTime(),
		})
	}
	f, err := os.Create("generated/projects.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(projects)
}

func init() {
	rootCmd.PersistentFlags().Int("users", 20, "number of users to generate")
	rootCmd.PersistentFlags().Int("projects", 20, "number of projects to generate")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
