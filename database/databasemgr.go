/*
 * databasemgr.go --- Database manager.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package database

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/Asmodai/gohacks/types"
)

type DatabaseMgr struct {
}

func (dbm *DatabaseMgr) Open(driver string, dsn string) (IDatabase, error) {
	return Open(driver, dsn)
}

func (dbm *DatabaseMgr) OpenConfig(conf *Config) (IDatabase, error) {
	db, err := dbm.Open(conf.Driver, conf.ToDSN())
	if err != nil {
		return nil, types.NewErrorAndLog(
			"DATABASE",
			err.Error(),
		)
	}

	if conf.SetPoolLimits == true {
		db.SetMaxIdleConns(conf.MaxIdleConns)
		db.SetMaxOpenConns(conf.MaxOpenConns)
	}

	return db, nil
}

func (dbm *DatabaseMgr) CheckDB(db IDatabase) error {
	if err := db.Ping(); err != nil {
		return types.NewError(
			"DATABASE",
			"Unable to ping database: %s",
			err.Error(),
		)
	}

	return nil
}

/* databasemgr.go ends here. */
