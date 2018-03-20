/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Icaro.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"flag"

	"github.com/nethesis/icaro/sun/sun-api/configuration"
	"github.com/nethesis/icaro/sun/sun-tasks/tasks"
)

func main() {
	// read flags
	ConfigFilePtr := flag.String("c", "/opt/icaro/sun-api/conf.json", "Path to configuration file")

	var TaskAction string
	flag.StringVar(&TaskAction, "a", "", "Specify action to execute")
	var TaskMode bool
	flag.BoolVar(&TaskMode, "w", false, "Runs as a background job")

	// parse flags
	flag.Parse()

	// init configuration
	configuration.Init(ConfigFilePtr)

	// init tasks
	tasks.Init(TaskAction, TaskMode)

}
