CREATE TABLE scout_interactions (
	scout_id int REFERENCES scouts(id),
	duration real NOT NULL,
	waypoints path NOT NULL,
	waypoint_widths path NOT NULL,
	waypoint_times real[] NOT NULL,
	processed boolean NOT NULL DEFAULT false,
	entered_at timestamp NOT NULL,
	created_at timestamp NOT NULL
);
CREATE INDEX scout_interactions_idx ON scout_interactions (created_at);