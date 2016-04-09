package main

import "github.com/robzienert/cqlmigrate"

// I could've used go-bindata (https://github.com/jteeuwen/go-bindata) to
// embed static files into the binary, but I think that is more effort than it's
// actually worth. This method will be simple enough.
var cassandraMigrations = []cqlmigrate.Spec{
	{
		Name: "2016-01-09-initial_release",
		Data: `
    CREATE TABLE features_namespaced (
			key varchar,
      namespace varchar,
      type varchar,
      value varchar,
      gate_value varchar,
      gate_groups list<varchar>,
      gate_actors list<varchar>,
      gate_actor_percent int,
      gate_percent_of_time int,
      date_created timestamp,
      last_updated timestamp,
      PRIMARY KEY(namespace, key)
    );

		CREATE TABLE features (
			key varchar,
      type varchar,
      value varchar,
      gate_value varchar,
      gate_groups list<varchar>,
      gate_actors list<varchar>,
      gate_actor_percent int,
      gate_percent_of_time int,
      date_created timestamp,
      last_updated timestamp,
      PRIMARY KEY(key)
		);

		CREATE TABLE audit (
			action varchar,
			actor varchar,
			fields map<varchar, varchar>,
			date_created timestamp,
			PRIMARY KEY((action, date_created))
		);
    `,
	},
}
